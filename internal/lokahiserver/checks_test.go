package lokahiserver

import (
	"context"
	"os"
	"testing"

	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestChecks(t *testing.T) {
	durl := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", durl)
	if err != nil {
		t.Fatal(err)
	}

	err = database.Migrate(durl)
	if err != nil {
		t.Fatal(err)
	}
	defer database.Destroy(durl)

	c := &Checks{DB: database.ChecksPostgres(db)}
	ctx := context.Background()

	chk, err := c.Create(ctx, &lokahi.CreateOpts{
		Url:        "https://google.memes",
		WebhookUrl: "https://youtube.memes",
		Every:      30,
	})
	if err != nil {
		t.Fatal(err)
	}

	chk2, err := c.Get(ctx, &lokahi.CheckID{Id: chk.Id})
	if err != nil {
		t.Fatal(err)
	}

	if chk.Url != chk2.Url {
		t.Fatalf("wanted %q == %q", chk.Url, chk2.Url)
	}

	chk2.WebhookUrl = "https://something.else"

	chk3, err := c.Put(ctx, chk2)
	if err != nil {
		t.Fatal(err)
	}

	if chk3.WebhookUrl != chk2.WebhookUrl {
		t.Fatalf("wanted %q == %q", chk3.WebhookUrl, chk2.WebhookUrl)
	}

	_, err = c.Delete(ctx, &lokahi.CheckID{Id: chk.Id})
	if err != nil {
		t.Fatal(err)
	}
}
