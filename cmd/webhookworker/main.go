package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/internal/lokahiadminserver"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
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

	lr := &lokahiadminserver.LocalRun{
		HC:  &http.Client{},
		Cs:  database.ChecksPostgres(db),
		Rs:  database.RunsPostgres(db),
		Ris: database.RunInfosPostgres(db),
	}

	sc, err := nc.Subscribe("webhook.egress", func(m *nats.Msg) {
		ln.Log(ctx, ln.Action("handler started"), ln.F{"channel": "webhook.egress"})

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
