package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/caarlos0/env"
)

type impl struct{}

func (i impl) Handle(ctx context.Context, st *lokahi.CheckStatus) (*lokahi.Nil, error) {
	log.Printf("check id: %s, state: %s, latency: %s, status code: %d, playbook url: %s", st.Check.Id, st.Check.State, time.Duration(st.LastResponseTimeNanoseconds), st.RespStatusCode, st.Check.PlaybookUrl)

	return &lokahi.Nil{}, nil
}

type config struct {
	Port string `env:"PORT" envDefault:"9001"`
}

func main() {
	ctx := context.Background()

	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		ln.FatalErr(ctx, err)
	}

	hs := lokahi.NewWebhookServer(impl{}, nil)
	mux := http.NewServeMux()

	mux.Handle(lokahi.WebhookPathPrefix, hs)

	ln.Log(ctx, ln.F{"port": cfg.Port}, ln.Action("listening on HTTP"))
	ln.FatalErr(ctx, http.ListenAndServe(":"+cfg.Port, mux))
}
