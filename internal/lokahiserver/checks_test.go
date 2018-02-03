package lokahiserver

import (
	"context"
	"testing"

	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func TestChecks(t *testing.T) {
	db, err := gorm.Open("postgres", "postgres://postgres:hunter2@127.0.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(&database.Check{}).Error
	if err != nil {
		t.Fatal(err)
	}
	defer db.Exec("DROP TABLE checks")

	c := &Checks{DB: db}
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
