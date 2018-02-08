package lokahiserver

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/Xe/uuid"
	"github.com/gogo/protobuf/proto"
	nats "github.com/nats-io/go-nats"
)

type Runs struct {
	Ah lokahiadmin.Health

	Cs database.Checks
	Rs database.Runs
	Nc *nats.Conn
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

func (r *Runs) Checks(ctx context.Context, cids *lokahi.CheckIDs) (*lokahi.Run, error) {
	rid := uuid.New()
	st := time.Now()

	ctx = ln.WithF(ctx, ln.F{
		"run_id":    rid,
		"cid_count": len(cids.Ids),
	})
	ln.Log(ctx, ln.Action("starting run"))

	result := &lokahi.Run{
		Results:       map[string]*lokahi.Health{},
		Cids:          cids,
		StartTimeUnix: st.Unix(),
	}

	defer func() { result.Finished = true }()

	var lock sync.Mutex
	var cks []database.Check

	for _, shardWork := range split(cids.Ids, 50) {
		go func(inp []string) {
			for _, cid := range inp {
				dck, err := r.Cs.Get(ctx, cid)
				if err != nil {
					ln.Error(ctx, err, ln.F{"cid": cid})
					continue
				}
				ctx := ln.WithF(ctx, dck.F())

				// todo:
				// res, err := r.Ah.Run(ctx, &lokahiadmin.RunRequest{
				// 	Check: lokahiadmin.CheckFromDatabaseCheck(dck),
				// 	RunId: rid,
				// })

				data, err := proto.Marshal(dck.AsProto())
				if err != nil {
					ln.Error(ctx, err, ln.F{"cid": cid})
					continue
				}

				var res lokahi.Health
				reply, err := r.Nc.Request("check.run", data, time.Minute)
				if err != nil {
					ln.Error(ctx, err, ln.F{"cid": cid})
					continue
				}

				err = proto.Unmarshal(reply.Data, &res)

				if res.Error != "" {
					panic(cid + ":" + res.Error)
				}

				dck, err = r.Cs.Get(ctx, cid)
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

	ela := time.Since(st)
	result.ElapsedNanoseconds = int64(ela)

	dbr := database.Run{
		UUID:    rid,
		Message: fmt.Sprintf("%d checks run in %v", len(result.Results), ela),
	}

	ln.Log(ctx, ln.Action("finished running checks"))

	_, err := r.Rs.Put(ctx, dbr)
	if err != nil {
		return nil, err
	}

	return result, nil
}
