// Package lokahiadminserver implements twirp package github.xe.lokahi.admin
package lokahiadminserver

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Xe/lokahi/rpc/lokahiadmin"
	nats "github.com/nats-io/go-nats"
)

func HealthPath(method string) string {
	return lokahiadmin.HealthPathPrefix + method
}

func WebhookPath(method string) string {
	return lokahiadmin.WebhookPathPrefix + method
}

type NatsRoundTripper struct {
	NC *nats.Conn
}

func (n *NatsRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 125*time.Millisecond)
	defer cancel()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	msg, err := n.NC.RequestWithContext(ctx, r.RequestURI, data)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(msg.Data)
	resp := &http.Response{
		Body:       ioutil.NopCloser(buf),
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header: http.Header{
			"Content-Length": []string{fmt.Sprint(len(msg.Data))},
			"Content-Type":   []string{"application/protobuf"},
		},
		ContentLength: int64(len(data)),
		Close:         true,
	}

	return resp, nil
}
