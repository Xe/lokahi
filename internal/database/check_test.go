package database

import (
	"context"
	"os"
	"testing"

	"github.com/Xe/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestChecks(t *testing.T) {
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

	cp := &checksPostgres{db: db}

	chk, err := cp.Create(ctx, Check{
		URL:         "http://ill.mend.your.heart",
		WebhookURL:  "http://with.threads.of.mine",
		PlaybookURL: "http://im.not.a.nurse/but/youll/be/just-fine.md",
		Every:       60,
	})
	if err != nil {
		t.Fatal(err)
	}

	chk2, err := cp.Get(ctx, chk.UUID)
	if err != nil {
		t.Fatal(err)
	}

	if chk.UUID != chk2.UUID {
		t.Fatalf("expected check uuids to match, wanted both to be %s, got: %s, %s", chk.UUID, chk.UUID, chk2.UUID)
	}

	_, err = cp.Delete(ctx, chk.UUID)
	if err != nil {
		t.Fatal(err)
	}

	for range make([]struct{}, 300) {
		_, err := cp.Create(ctx, Check{
			URL:         "http://ill.mend.your.heart?" + uuid.New(),
			WebhookURL:  "http://with.threads.of.mine",
			PlaybookURL: "http://im.not.a.nurse/but/youll/be/just-fine.md",
			Every:       60,
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	cks, err := cp.List(ctx, 300, 0)
	if err != nil {
		t.Fatal(err)
	}

	if len(cks) != 300 {
		t.Fatalf("wanted len(cks) == 300, got: %v", len(cks))
	}

	chk3, err := cp.Create(ctx, Check{
		URL:         "http://ill.mend.your.heart?" + uuid.New(),
		WebhookURL:  "http://with.threads.of.mine",
		PlaybookURL: "http://im.not.a.nurse/but/youll/be/just-fine.md",
		Every:       60,
	})
	if err != nil {
		t.Fatal(err)
	}

	chk3.WebhookURL = "http://boot.fun"

	chk4, err := cp.Put(ctx, *chk3)
	if err != nil {
		t.Fatal(err)
	}

	if chk4.WebhookURL != "http://boot.fun" {
		t.Fatalf("wanted chk4 to have a modified webhook url, wanted http://boot.fun, got: %s", chk4.WebhookURL)
	}

	vals, err := cp.ListByEveryValue(ctx, 60)
	if err != nil {
		t.Fatal(err)
	}

	if len(vals) == 0 {
		t.Fatal("no results fetched from ListByEveryValue")
	}
}
