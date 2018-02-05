package database

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestChecks(t *testing.T) {
	durl := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", durl)
	if err != nil {
		t.Fatal(err)
	}

	err = Migrate(durl)
	if err != nil {
		t.Fatal(err)
	}
	defer destroy(durl)

	cp := &checksPostgres{db: db}

	chk, err := cp.Create(context.Background(), Check{
		URL:         "http://ill.mend.your.heart",
		WebhookURL:  "http://with.threads.of.mine",
		PlaybookURL: "http://im.not.a.nurse/but/youll/be/just-fine.md",
		Every:       60,
	})
	if err != nil {
		t.Fatal(err)
	}
}
