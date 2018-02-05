package main

import (
	"log"

	"github.com/Xe/lokahi/internal/database"
	"github.com/caarlos0/env"
	_ "github.com/mattes/migrate/database/postgres"
)

var cfg = struct {
	DBURL string `env:"DATABASE_URL,required"`
}{}

func main() {
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrate(cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
}
