package lokahiadminserver

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/Xe/uuid"
	"github.com/codahale/hdrhistogram"
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

	ck.WebhookResponseTimeNanoseconds = int64(diff)
	err = l.DB.Save(&ck).Error
	if err != nil {
		result.Error = err.Error()
		return result, ck
	}

	l.timing.RecordValue(int64(diff))

	if sc := resp.StatusCode % 100; sc == 2 {
		log.Println(sc)
		result.Healthy = true
	}

	ln.Log(ctx, ck, ln.F{"resp_status_code": resp.StatusCode, "resp_time": diff})

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
