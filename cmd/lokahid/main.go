package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/internal/lokahiadminserver"
	"github.com/Xe/lokahi/internal/lokahiserver"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/caarlos0/env"
	"github.com/heroku/x/scrub"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	nats "github.com/nats-io/go-nats"
)

type config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	NatsURL     string `env:"NATS_URL,required"`
	NoPass      bool   `env:"NO_PASS" envDefault:"false"`
	Port        string `env:"PORT" envDefault:"5000"`
}

func (c config) F() ln.F {
	result := ln.F{
		"env_NO_PASS": c.NoPass,
		"env_PORT":    c.Port,
	}

	u, err := url.Parse(c.DatabaseURL)
	if err != nil {
		result["env_DATABASE_URL_err"] = err
	} else {
		result["env_DATABASE_URL"] = scrub.URL(u)
	}

	u, err = url.Parse(c.NatsURL)
	if err != nil {
		result["env_NATS_URL_err"] = err
	} else {
		result["env_NATS_URL"] = scrub.URL(u)
	}

	return result
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

	lahc := &http.Client{
		Transport: &lokahiadminserver.NatsRoundTripper{
			NC: nc,
		},
	}

	ctx = ln.WithF(ctx, cfg.F())

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

	cks := &lokahiserver.Checks{
		DB: database.ChecksPostgres(db),
	}

	rs := &lokahiserver.Runs{
		Ah: lokahiadmin.NewHealthProtobufClient("http://nats", lahc),

		Cs: database.ChecksPostgres(db),
		Rs: database.RunsPostgres(db),
	}

	mux := http.NewServeMux()
	mux.Handle(lokahi.RunsPathPrefix, lokahi.NewRunsServer(rs, makeLnHooks()))
	mux.Handle(lokahi.ChecksPathPrefix, lokahi.NewChecksServer(cks, makeLnHooks()))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := db.Exec("SELECT 1+1")
		if err != nil {
			ln.Error(r.Context(), err)
			http.Error(w, "database error", http.StatusInternalServerError)
		}
	})

	ln.Log(ctx, ln.F{"port": os.Getenv("PORT")}, ln.Action("Listening on http"))
	ln.FatalErr(ctx, http.ListenAndServe(":"+cfg.Port, metaInfo(mux)), ln.Action("http server stopped for some reason"))
}

func metaInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, _ := net.SplitHostPort(r.RemoteAddr)
		f := ln.F{
			"remote_ip":       host,
			"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
		}
		ctx := ln.WithF(r.Context(), f)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
