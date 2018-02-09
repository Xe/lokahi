package main

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	UserPass    string `env:"USERPASS"`
	DatabaseURL string `env:"DATABASE_URL,required"`
	NatsURL     string `env:"NATS_URL,required"`
	Port        string `env:"PORT" envDefault:"5000"`
}

func (c config) F() ln.F {
	result := ln.F{
		"env_PORT": c.Port,
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

	var h http.Handler
	h = metaInfo(mux)
	if cfg.UserPass != "" {
		pair := strings.SplitN(cfg.UserPass, ":", 2)
		if len(pair) != 2 {
			log.Fatalf("expected %s to have one colon", cfg.UserPass)
			return
		}

		h = auth(pair[0], pair[1])(h)
	}

	ln.Log(ctx, ln.F{"port": os.Getenv("PORT")}, ln.Action("Listening on http"))
	ln.FatalErr(ctx, http.ListenAndServe(":"+cfg.Port, h), ln.Action("http server stopped for some reason"))
}

func auth(user, pass string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				http.Error(w, "Not authorized", 401)
				return
			}

			b, err := base64.StdEncoding.DecodeString(s[1])
			if err != nil {
				http.Error(w, err.Error(), 401)
				return
			}

			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				http.Error(w, "Not authorized", 401)
				return
			}

			userCmp := subtle.ConstantTimeCompare([]byte(pair[0]), []byte(user))
			passCmp := subtle.ConstantTimeCompare([]byte(pair[1]), []byte(pass))

			if userCmp == 1 {
				if passCmp != 1 {
					http.Error(w, "Not authorized", 401)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
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
