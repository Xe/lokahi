package database

import (
	"context"
	"os"
	"testing"

	"github.com/Xe/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestRuns(t *testing.T) {
	ctx := context.Background()
	durl := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", durl)
	if err != nil {
		t.Fatal(err)
	}

	err = Migrate(durl)
	if err != nil {
		t.Fatal(err)
	}
	defer Destroy(durl)

	rns := runsPostgres{db: db}

	rn, err := rns.Put(ctx, Run{
		Message: "i love you",
	})
	if err != nil {
		t.Fatal(err)
	}

	rn2, err := rns.Get(ctx, rn.UUID)
	if err != nil {
		t.Fatal(err)
	}

	if rn2.Message != "i love you" {
		t.Fatal("message not intact")
	}
}

func TestRunInfos(t *testing.T) {
	ctx := context.Background()
	durl := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", durl)
	if err != nil {
		t.Fatal(err)
	}

	err = Migrate(durl)
	if err != nil {
		t.Fatal(err)
	}
	defer Destroy(durl)

	ris := runInfoPostgres{db: db}
	rid := uuid.New()

	err = ris.Put(ctx, RunInfo{
		UUID:                           uuid.New(),
		RunID:                          rid,
		CheckID:                        uuid.New(),
		ResponseTimeNanoseconds:        42069,
		WebhookResponseTimeNanoseconds: 9001,
	})
	if err != nil {
		t.Fatal(err)
	}

	riset, err := ris.GetRun(ctx, rid)
	if err != nil {
		t.Fatal(err)
	}

	if len(riset) == 0 {
		t.Fatal("expected riset to have len != 0")
	}
}
