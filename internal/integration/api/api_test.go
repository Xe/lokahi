package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/internal/lokahiserver"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type apiCtx struct {
	ts *httptest.Server
	db *sqlx.DB

	checks lokahi.Checks

	// check to create on iCreateTheCheck
	checkCreateOpts *lokahi.CreateOpts
	// resulting check
	rc *lokahi.Check
}

func (a *apiCtx) aBaseStack() error {
	durl := os.Getenv("DATABASE_URL")
	if durl == "" {
		return errors.New("no DATABASE_URL")
	}

	err := database.Destroy(durl)
	if err != nil && !strings.Contains(err.Error(), "no change") {
		return err
	}

	err = database.Migrate(durl)
	if err != nil {
		return err
	}

	db, err := sqlx.Open("postgres", durl)
	if err != nil {
		return err
	}

	a.db = db

	cks := &lokahiserver.Checks{
		DB: database.ChecksPostgres(db),
	}

	mux := http.NewServeMux()
	mux.Handle(lokahi.ChecksPathPrefix, lokahi.NewChecksServer(cks, nil))

	a.ts = httptest.NewServer(mux)
	a.checks = lokahi.NewChecksProtobufClient(a.ts.URL, &http.Client{})

	return nil
}

func (a *apiCtx) iWantToCreateACheck() error {
	a.checkCreateOpts = &lokahi.CreateOpts{}

	return nil
}

func (a *apiCtx) aCheckMonitoringUrlOf(curl string) error {
	a.checkCreateOpts.Url = curl

	return nil
}

func (a *apiCtx) aCheckWebhookUrlOf(wurl string) error {
	a.checkCreateOpts.WebhookUrl = wurl

	return nil
}

func (a *apiCtx) aCheckEveryOf(every int) error {
	a.checkCreateOpts.Every = int32(every)

	return nil
}

func (a *apiCtx) aCheckPlaybookUrlOf(purl string) error {
	a.checkCreateOpts.PlaybookUrl = purl

	return nil
}

func (a *apiCtx) iCreateTheCheck() error {
	ck, err := a.checks.Create(context.Background(), a.checkCreateOpts)
	if err != nil {
		return err
	}

	a.rc = ck

	return nil
}

func (a *apiCtx) theResultingCheckShouldHaveAnID() error {
	if a.rc.Id == "" {
		return errors.New("the check doesn't have an ID")
	}

	return nil
}

func (a *apiCtx) tearEverythingDown() error {
	err := a.db.Close()
	if err != nil {
		return err
	}

	a.ts.Close()

	err = database.Destroy(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &apiCtx{}

	s.Step(`^a base stack$`, a.aBaseStack)
	s.Step(`^I want to create a check$`, a.iWantToCreateACheck)
	s.Step(`^a check monitoring url of "([^"]*)"$`, a.aCheckMonitoringUrlOf)
	s.Step(`^a check webhook url of "([^"]*)"$`, a.aCheckWebhookUrlOf)
	s.Step(`^a check every of (\d+)$`, a.aCheckEveryOf)
	s.Step(`^a check playbook url of "([^"]*)"$`, a.aCheckPlaybookUrlOf)
	s.Step(`^I create the check$`, a.iCreateTheCheck)
	s.Step(`^the resulting check should have an ID$`, a.theResultingCheckShouldHaveAnID)
	s.Step(`^tear everything down$`, a.tearEverythingDown)
}
