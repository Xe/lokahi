package lokahiadminserver

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/gogo/protobuf/proto"
)

type Webhook struct {
	HC *http.Client
}

func (w *Webhook) Execute(ctx context.Context, wr *lokahiadmin.WebhookRequest) (*lokahiadmin.Nil, error) {
	w.sendWebhook(ctx, wr.Check.DatabaseCheck().AsProto(), wr.Health)

	return &lokahiadmin.Nil{}, nil
}

func (w *Webhook) sendWebhook(ctx context.Context, ck *lokahi.Check, health *lokahiadmin.ServiceHealth) {
	cid := ck.Id

	logErr := func(err error, cid, u string) {
		ln.Error(ctx, err, ln.F{"check_id": cid, "url": u})
	}

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

	resp, err := w.HC.Do(req)
	if err != nil {
		logErr(err, cid, ck.WebhookUrl)
		return
	}

	if s := resp.StatusCode / 100; s != 2 {
		logErr(fmt.Errorf("lokahiadminserver: %s gave HTTP status %d(%d)", ck.WebhookUrl, resp.StatusCode, s), cid, ck.WebhookUrl)
	}
}
