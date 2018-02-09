package lokahiadminserver

import (
	"context"
	"net/http"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/Xe/uuid"
	"google.golang.org/api/support/bundler"
)

type Health struct {
	ribdl *bundler.Bundler

	HC  *http.Client
	Ris database.RunInfos
	Cs  database.Checks
}

func (h *Health) Run(ctx context.Context, ri *lokahiadmin.RunRequest) (*lokahiadmin.ServiceHealth, error) {
	if h.ribdl == nil {
		h.ribdl = bundler.NewBundler(database.RunInfo{}, func(i interface{}) {
			data, ok := i.([]database.RunInfo)
			if !ok {
				return
			}

			for _, d := range data {
				err := h.Ris.Put(context.Background(), d)
				if err != nil {
					ln.Error(context.Background(), err, ln.Action("putting deferred runinfo"), ln.F{"run_id": d.RunID, "check_id": d.CheckID})
				}
			}
		})
		h.ribdl.BundleCountThreshold = 300
		h.ribdl.DelayThreshold = 5 * time.Second
		h.ribdl.BundleByteThreshold = 1024 * 1024 * 1024
	}

	ck := ri.Check.DatabaseCheck()
	dck := &ck

	result, err := h.doCheck(ctx, dck, ri.RunId)
	if err != nil {
		return nil, err
	}

	_, err = h.Cs.Put(ctx, *dck)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *Health) doCheck(ctx context.Context, ck *database.Check, rid string) (*lokahiadmin.ServiceHealth, error) {
	result := &lokahiadmin.ServiceHealth{}

	st := time.Now()
	req, err := http.NewRequest("GET", ck.URL, nil)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	resp, err := h.HC.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}
	ed := time.Now()
	diff := ed.Sub(st)

	result.StatusCode = int32(resp.StatusCode)
	result.Url = ck.URL
	result.ResponseTimeNanoseconds = int64(diff)

	if sc := resp.StatusCode / 100; sc == 2 {
		result.Healthy = true
		ck.State = lokahi.Check_UP.String()
	} else {
		ck.State = lokahi.Check_DOWN.String()
	}

	if result.Error != "" {
		ck.State = lokahi.Check_ERROR.String()
	}

	err = h.ribdl.Add(database.RunInfo{
		UUID:                    uuid.New(),
		RunID:                   rid,
		CheckID:                 ck.UUID,
		ResponseStatus:          resp.StatusCode,
		ResponseTimeNanoseconds: int64(diff),
	}, 50)
	if err != nil {
		result.Error = err.Error()
		return result, err
	}

	return result, nil
}
