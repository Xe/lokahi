package lokahiadminserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/Xe/uuid"
	"github.com/codahale/hdrhistogram"
	"github.com/gogo/protobuf/proto"
	"github.com/jinzhu/gorm"
)

type Runner interface {
	Run(context.Context, *lokahiadmin.CheckIDs) (*lokahiadmin.Run, error)
}

type LocalRun struct {
	HC *http.Client
	DB *gorm.DB

	timing *hdrhistogram.Histogram
}

func (l *LocalRun) Minutely() error {
	ln.Log(context.Background(), ln.Action("minutelyCron"))

	var checks []database.Check

	err := l.DB.Where("every = 60").Find(&checks).Error
	if err != nil {
		return err
	}

	cids := &lokahiadmin.CheckIDs{}

	for _, ck := range checks {
		cids.Ids = append(cids.Ids, ck.UUID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()

	result, err := l.Run(ctx, cids)
	if err != nil {
		return err
	}

	logErr := func(err error, cid, u string) {
		ln.Error(ctx, err, ln.F{"check_id": cid, "url": u})
	}

	for cid, health := range result.Results {
		var ck database.Check
		err = l.DB.Where("uuid = ?", cid).First(&ck).Error
		if err != nil {
			logErr(err, cid, ck.WebhookURL)
			continue
		}

		cs := &lokahi.CheckStatus{
			Check: ck.AsProto(),
			LastResponseTimeNanoseconds: health.ResponseTimeNanoseconds,
		}

		data, err := proto.Marshal(cs)
		if err != nil {
			logErr(err, cid, ck.WebhookURL)
			continue
		}

		buf := bytes.NewBuffer(data)

		req, err := http.NewRequest("POST", ck.WebhookURL, buf)
		if err != nil {
			logErr(err, cid, ck.WebhookURL)
			continue
		}
		req.Header.Add("Content-Type", "application/protobuf")
		req.Header.Add("Accept", "application/protobuf")

		st := time.Now()
		var ed time.Time
		var diff time.Duration
		resp, err := l.HC.Do(req)
		if err != nil {
			logErr(err, cid, ck.WebhookURL)
			goto save
		}
		ed = time.Now()
		diff = ed.Sub(st)

		if s := resp.StatusCode / 100; s != 2 {
			logErr(fmt.Errorf("lokahiadminserver: %s gave HTTP status %d(%d)", ck.WebhookURL, resp.StatusCode, s), cid, ck.WebhookURL)
		}

		ck.WebhookResponseTimeNanoseconds = int64(diff)

	save:
		err = l.DB.Save(&ck).Error
		if err != nil {
			logErr(err, cid, ck.WebhookURL)
		}

	}

	return nil
}

func (l *LocalRun) doCheck(ctx context.Context, cid string) (*lokahiadmin.Run_Health, database.Check) {
	result := &lokahiadmin.Run_Health{}

	var ck database.Check
	err := l.DB.Where("uuid = ?", cid).First(&ck).Error
	if err != nil {
		result.Error = err.Error()
		return result, ck
	}

	ln.Log(ctx, ck)

	st := time.Now()
	req, err := http.NewRequest("GET", ck.URL, nil)
	if err != nil {
		result.Error = err.Error()
		return result, ck
	}

	resp, err := l.HC.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result, ck
	}
	ed := time.Now()
	diff := ed.Sub(st)

	result.StatusCode = int32(resp.StatusCode)
	result.Url = ck.URL
	result.ResponseTimeNanoseconds = int64(diff)

	l.timing.RecordValue(int64(diff))

	if sc := resp.StatusCode / 100; sc == 2 {
		result.Healthy = true
		ck.State = lokahi.Check_UP.String()
	} else {
		ck.State = lokahi.Check_DOWN.String()
	}

	ln.Log(ctx, ck, ln.Action("doCheckHTTPDone"), ln.F{"resp_status_code": resp.StatusCode, "resp_time": diff})

	err = l.DB.Save(&ck).Error
	if err != nil {
		panic(err)
		result.Error = err.Error()
		return result, ck
	}

	return result, ck

}

func (l *LocalRun) Run(ctx context.Context, cids *lokahiadmin.CheckIDs) (*lokahiadmin.Run, error) {
	if l.timing == nil {
		l.timing = hdrhistogram.New(0, 300000000000, 1)
	}

	result := &lokahiadmin.Run{
		Results: map[string]*lokahiadmin.Run_Health{},
		Cids:    cids,
		Id:      uuid.New(),
	}
	defer func() { result.Finished = true }()
	st := time.Now()
	defer func() { result.StartTimeUnix = st.Unix() }()

	lock := sync.Mutex{}

	var cks []database.Check
	for _, cid := range cids.Ids {
		go func(cid string) {
			res, ck := l.doCheck(ctx, cid)

			if res.Error != "" {
				panic(cid + ":" + res.Error)
			}

			lock.Lock()
			result.Results[cid] = res
			cks = append(cks, ck)
			lock.Unlock()
		}(cid)
	}

	for len(cks) != len(cids.Ids) {
		time.Sleep(125 * time.Millisecond)
	}

	ela := time.Now().Sub(st)
	result.ElapsedNanoseconds = int64(ela)

	data, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	dbr := database.Run{
		Finished: true,
		Message:  string(data),
		Checks:   cks,
	}

	err = l.DB.Save(&dbr).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (l *LocalRun) Checks(ctx context.Context, cids *lokahiadmin.CheckIDs) (*lokahiadmin.Run, error) {
	return l.Run(ctx, cids)
}

func (l *LocalRun) Stats(ctx context.Context, _ *lokahiadmin.Nil) (*lokahiadmin.HistogramData, error) {
	result := &lokahiadmin.HistogramData{
		MaxNanoseconds:  l.timing.Max(),
		MinNanoseconds:  l.timing.Min(),
		MeanNanoseconds: int64(l.timing.Mean()),
		Stddev:          int64(l.timing.StdDev()),
		P50Nanoseconds:  int64(l.timing.ValueAtQuantile(50)),
		P75Nanoseconds:  int64(l.timing.ValueAtQuantile(75)),
		P95Nanoseconds:  int64(l.timing.ValueAtQuantile(95)),
		P99Nanoseconds:  int64(l.timing.ValueAtQuantile(99)),
	}

	return result, nil
}
