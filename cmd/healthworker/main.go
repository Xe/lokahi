package main

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/Xe/ln"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/internal/lokahiadminserver"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/Xe/uuid"
	"github.com/caarlos0/env"
	"github.com/gogo/protobuf/proto"
	"github.com/jmoiron/sqlx"
	nats "github.com/nats-io/go-nats"
)

type config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	NatsURL     string `env:"NATS_URL,required"`
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx = ln.WithF(ctx, ln.F{"in": "main"})

	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		ln.FatalErr(ctx, err)
	}

	err = database.Migrate(cfg.DatabaseURL)
	if err != nil && err.Error() != "no change" {
		ln.FatalErr(ctx, err)
	}

	// wait for postgres
	time.Sleep(2 * time.Second)
	db, err := sqlx.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		ln.FatalErr(ctx, err)
	}

	db.SetMaxOpenConns(30)

	tr := rehttp.NewTransport(
		nil, // will use http.DefaultTransport
		rehttp.RetryAll(rehttp.RetryMaxRetries(3), rehttp.RetryTemporaryErr()), // max 3 retries for Temporary errors
		rehttp.ConstDelay(time.Second),                                         // wait 1s between retries
	)

	h := &lokahiadminserver.Health{
		HC:  &http.Client{Transport: tr},
		Cs:  database.ChecksPostgres(db),
		Ris: database.RunInfosPostgres(db),
	}

	for range make([]struct{}, runtime.NumCPU()*4) {
		nc, err := nats.Connect(cfg.NatsURL)
		if err != nil {
			log.Fatal(err)
		}

		sc, err := nc.QueueSubscribe(lokahiadminserver.HealthPath("Run"), "healthworker", func(m *nats.Msg) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			rr := &lokahiadmin.RunRequest{}
			err := proto.Unmarshal(m.Data, rr)
			if err != nil {
				ln.Error(ctx, err, ln.Action("nats check.run handler unmarshal check"))
				return
			}

			rid := uuid.New()

			hlt, err := h.Run(ctx, rr)
			if err != nil {
				ln.Error(ctx, err, ln.Action("nats check.run handler"))
				return
			}

			data, err := proto.Marshal(hlt)
			if err != nil {
				ln.Error(ctx, err, ln.Action("nats check.run handler"))
				return
			}

			err = nc.Publish(m.Reply, data)
			if err != nil {
				ln.Error(ctx, err)
			}
		})
		if err != nil {
			log.Fatal(err)
		}
		sc.SetPendingLimits(5000, 65535)
	}

	ln.Log(ctx, ln.Action("waiting for work..."), ln.F{"threads": runtime.NumCPU() * 4})
	for {
		select {}
	}
}
