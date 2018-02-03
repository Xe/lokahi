package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/jinzhu/gorm"
)

type Runner interface {
	Run(context.Context, *lokahiadmin.CheckIDs) (*lokahiadmin.Run, error)
}

type LocalRunner struct {
	db  *gorm.DB
	cli *http.Client
}

func (l *LocalRunner) doCheck(cid string) *lokahiadmin.Run_Health {
	result := &lokahiadmin.Run_Health{}

	var ck database.Check
	err := l.db.Where("uuid = ?", cid).First(&ck).Error
	if err != nil {
		result.Error = err.Error()
		return result
	}

	st := time.Now()
	req, err := http.NewRequest("GET", ck.URL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	resp, err := l.cli.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	ed := time.Now()
	diff := ed.Sub(st)

	result.StatusCode = int32(resp.StatusCode)
	result.Url = ck.URL
	result.ResponseTimeNanoseconds = int64(diff)

	if sc := resp.StatusCode % 100; sc == 2 {
		result.Healthy = true
	}

	return result

}

func (l *LocalRunner) Checks(ctx context.Context, cids *lokahiadmin.CheckIDs) (*lokahiadmin.Run, error) {
	result := &lokahiadmin.Run{}
	defer func() { result.Finished = true }()

	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	for _, cid := range cids.Ids {
		go func(cid string) {
			wg.Add(1)
			res := l.doCheck(cid)

			lock.Lock()
			result.Results[cid] = res
			lock.Unlock()

			wg.Done()
		}(cid)
	}

	wg.Wait()

	return result, nil
}
