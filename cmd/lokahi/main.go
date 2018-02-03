package main

import (
	"context"

	"github.com/Xe/ln"
	"github.com/caarlos0/env"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	NoPass      bool   `env:"NO_PASS" envDefault:"false"`
	Port        string `env:"PORT" envDefault:"5000"`
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
}
