package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/rehttp"
	"github.com/Xe/ln"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/internal/lokahiadminserver"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/caarlos0/env"
	"github.com/codahale/hdrhistogram"
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

	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		log.Fatal(err)
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

	lr := &lokahiadminserver.LocalRun{
		HC:  &http.Client{Transport: tr},
		Cs:  database.ChecksPostgres(db),
		Rs:  database.RunsPostgres(db),
		Ris: database.RunInfosPostgres(db),
	}

	hs := hdrhistogram.New(0, 9999999999999, 1)

	go func() {
		for {
			time.Sleep(time.Minute)

			ln.Log(
				ctx,
				ln.Action("performance data"),
				ln.F{
					"min":  hs.Min(),
					"mean": hs.Mean(),
					"max":  hs.Max(),
					"p95":  hs.ValueAtQuantile(95),
				},
			)
		}
	}()

	sc, err := nc.QueueSubscribe("webhook.egress", "webhookworker", func(m *nats.Msg) {
		st := time.Now()
		defer func() { hs.RecordValue(int64(time.Now().Sub(st))) }()

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		wd := &lokahiadmin.WebhookData{}
		err := proto.Unmarshal(m.Data, wd)
		if err != nil {
			ln.Error(ctx, err, ln.Action("nats webhook.egress handler"))
			return
		}

		c := &lokahi.Check{}
		err = proto.Unmarshal(wd.CheckProto, c)
		if err != nil {
			ln.Error(ctx, err, ln.Action("nats webhook.egress handler"))
			return
		}

		if wd.Health.Healthy {
			c.State = lokahi.Check_UP
		} else {
			c.State = lokahi.Check_DOWN
		}

		if wd.Health.Error != "" {
			c.State = lokahi.Check_ERROR

			nc.Publish("check.errors", m.Data)
		}

		lr.SendWebhook(ctx, c, wd.Health, func() {})
	})
	if err != nil {
		ln.FatalErr(ctx, err, ln.Action("nats webhook.egress subscribe"))
	}
	sc.SetPendingLimits(500, 65535)

	ln.Log(ctx, ln.Action("waiting for work..."))
	for {
		select {}
	}
}
