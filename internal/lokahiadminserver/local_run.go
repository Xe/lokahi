package lokahiadminserver

import (
	"bytes"
	"context"
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
	nats "github.com/nats-io/go-nats"
	"google.golang.org/api/support/bundler"
)

type Runner interface {
	Run(context.Context, *lokahiadmin.CheckIDs) (*lokahiadmin.Run, error)
}

type LocalRun struct {
	HC *http.Client

	Cs  database.Checks
	Rs  database.Runs
	Ris database.RunInfos

	Nc *nats.Conn

	ribdl *bundler.Bundler
	ckbdl *bundler.Bundler

	lastState map[string]int32
	timing    *hdrhistogram.Histogram
}

func (l *LocalRun) Minutely() error {
	ctx := context.Background()
	ctx = ln.WithF(ctx, ln.F{"at": "localRun Minutely cron"})
	ln.Log(ctx)

	if l.ckbdl == nil {
		l.ckbdl = bundler.NewBundler(database.Check{}, func(i interface{}) {
			data, ok := i.([]database.Check)
			if !ok {
				return
			}

			for _, d := range data {
				_, err := l.Cs.Put(context.Background(), d)
				if err != nil {
					ln.Error(context.Background(), err, ln.Action("putting deferred check"), ln.F{"check_id": d.UUID})
				}
			}

			//ln.Log(context.Background(), ln.Action("put deferred check"), ln.F{"count": len(data)})
		})
		l.ckbdl.BundleCountThreshold = 300
		l.ckbdl.DelayThreshold = time.Second
		l.ckbdl.BundleByteThreshold = 1024 * 1024 * 1024
	}

	if l.lastState == nil {
		l.lastState = map[string]int32{}
	}

	checks, err := l.Cs.ListByEveryValue(ctx, 60)
	if err != nil {
		return err
	}

	cids := &lokahiadmin.CheckIDs{}

	for _, ck := range checks {
		cids.Ids = append(cids.Ids, ck.UUID)
	}

	ctx, cancel := context.WithTimeout(ctx, 40*time.Second)
	defer cancel()

	result, err := l.Run(ctx, cids)
	if err != nil {
		return err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		for cid, health := range result.Results {
			cst, ok := l.lastState[cid]
			if ok {
				if cst == health.StatusCode {
					continue
				}
			}

			var c database.Check

			for _, cc := range checks {
				if cc.UUID == cid {
					c = cc
				}
			}

			l.lastState[cid] = health.StatusCode

			cdata, _ := proto.Marshal(c.AsProto())
			data, _ := proto.Marshal(&lokahiadmin.WebhookData{
				RunId:      result.Id,
				CheckProto: cdata,
				Health:     health,
			})
			err := l.Nc.Publish("webhook.egress", data)
			if err != nil {
				ln.Error(ctx, err, c)
			}
		}

		ln.Log(ctx, ln.Action("done sending webhooks"))
	}()

	return nil
}

// SendWebhook sends a webhook to a given target by check id.
func (l *LocalRun) SendWebhook(ctx context.Context, ck *lokahi.Check, health *lokahiadmin.Health, done func()) {
	cid := ck.Id

	logErr := func(err error, cid, u string) {
		ln.Error(ctx, err, ln.F{"check_id": cid, "url": u})
	}

	defer done()

	cs := &lokahi.CheckStatus{
		Check: ck,
		LastResponseTimeNanoseconds: health.ResponseTimeNanoseconds,
	}

	data, err := proto.Marshal(cs)
	if err != nil {
		logErr(err, cid, ck.WebhookUrl)
		return
	}

	buf := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", ck.WebhookUrl, buf)
	if err != nil {
		logErr(err, cid, ck.WebhookUrl)
		return
	}

	req = req.WithContext(ctx)

	req.Header.Add("Content-Type", "application/protobuf")
	req.Header.Add("Accept", "application/protobuf")
	req.Header.Add("User-Agent", "lokahi/dev (+https://github.com/Xe/lokahi)")

	resp, err := l.HC.Do(req)
	if err != nil {
		logErr(err, cid, ck.WebhookUrl)
		return
	}

	if s := resp.StatusCode / 100; s != 2 {
		logErr(fmt.Errorf("lokahiadminserver: %s gave HTTP status %d(%d)", ck.WebhookUrl, resp.StatusCode, s), cid, ck.WebhookUrl)
	}
}

// DoCheck executes a HTTP healthcheck given a run id and check.
func (l *LocalRun) DoCheck(ctx context.Context, rid string, ck *database.Check) (*lokahiadmin.Health, database.Check) {
	if l.ribdl == nil {
		l.ribdl = bundler.NewBundler(database.RunInfo{}, func(i interface{}) {
			data, ok := i.([]database.RunInfo)
			if !ok {
				return
			}

			for _, d := range data {
				err := l.Ris.Put(context.Background(), d)
				if err != nil {
					ln.Error(context.Background(), err, ln.Action("putting deferred runinfo"), ln.F{"run_id": d.RunID, "check_id": d.CheckID})
				}
			}

			//ln.Log(context.Background(), ln.Action("put deferred runinfo"), ln.F{"count": len(data)})
		})
		l.ribdl.BundleCountThreshold = 300
		l.ribdl.DelayThreshold = 5 * time.Second
		l.ribdl.BundleByteThreshold = 1024 * 1024 * 1024
	}

	if l.timing == nil {
		l.timing = hdrhistogram.New(0, 30*10000000000000, 1)
	}

	result := &lokahiadmin.Health{}

	st := time.Now()
	req, err := http.NewRequest("GET", ck.URL, nil)
	if err != nil {
		result.Error = err.Error()
		return result, *ck
	}

	resp, err := l.HC.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result, *ck
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

	if result.Error != "" {
		ck.State = lokahi.Check_ERROR.String()
	}

	ck, err = l.Cs.Put(ctx, *ck)
	if err != nil {
		result.Error = err.Error()
		return result, *ck
	}

	err = l.ribdl.Add(database.RunInfo{
		UUID:                    uuid.New(),
		RunID:                   rid,
		CheckID:                 ck.UUID,
		ResponseStatus:          resp.StatusCode,
		ResponseTimeNanoseconds: int64(diff),
	}, 50)
	if err != nil {
		result.Error = err.Error()
	}

	return result, *ck
}

func split(buf []string, lim int) [][]string {
	var chunk []string
	chunks := make([][]string, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

func (l *LocalRun) Run(ctx context.Context, cids *lokahiadmin.CheckIDs) (*lokahiadmin.Run, error) {
	if l.timing == nil {
		l.timing = hdrhistogram.New(0, 30*10000000000000, 1)
	}

	rid := uuid.New()

	result := &lokahiadmin.Run{
		Results: map[string]*lokahiadmin.Health{},
		Cids:    cids,
	}
	defer func() { result.Finished = true }()
	st := time.Now()
	defer func() { result.StartTimeUnix = st.Unix() }()

	var lock sync.Mutex
	var cks []database.Check

	for _, shardWork := range split(cids.Ids, 50) {
		go func(inp []string) {
			for _, cid := range inp {
				dck, err := l.Cs.Get(ctx, cid)
				if err != nil {
					ln.Error(ctx, err, ln.F{"cid": cid})
					continue
				}

				data, err := proto.Marshal(dck.AsProto())
				if err != nil {
					ln.Error(ctx, err, ln.F{"cid": cid})
					continue
				}

				var res lokahiadmin.Health
				reply, err := l.Nc.Request("check.run", data, time.Minute)
				if err != nil {
					ln.Error(ctx, err, ln.F{"cid": cid})
					continue
				}

				err = proto.Unmarshal(reply.Data, &res)

				if res.Error != "" {
					panic(cid + ":" + res.Error)
				}

				dck, err = l.Cs.Get(ctx, cid)
				if err != nil {
					ln.Error(ctx, err, ln.F{"cid": cid})
					continue
				}

				lock.Lock()
				result.Results[cid] = &res
				cks = append(cks, *dck)
				lock.Unlock()
			}
		}(shardWork)
	}

	for len(cks) != len(cids.Ids) {
		time.Sleep(125 * time.Millisecond)
	}

	ela := time.Now().Sub(st)
	result.ElapsedNanoseconds = int64(ela)

	dbr := database.Run{
		UUID:    rid,
		Message: fmt.Sprintf("%d checks run in %v", len(result.Results), ela),
	}

	_, err := l.Rs.Put(ctx, dbr)
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
