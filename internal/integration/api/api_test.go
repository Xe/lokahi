package api

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/Xe/lokahi/internal/integration"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/uuid"
	_ "github.com/lib/pq"
)

type api struct {
	*integration.Suite

	checkListOpts   *lokahi.ListOpts
	checkCreateOpts *lokahi.CreateOpts
	rc              *lokahi.Check
}

func (a *api) iTryToCreateTheCheck() error {
	ck, err := a.ClientChecks.Create(context.Background(), a.checkCreateOpts)
	a.rc = ck
	a.SetErr(err)

	return nil
}

func (a *api) iTryToDeleteTheCheck() error {
	log.Printf("deleting check %s", a.rc.Id)
	_, err := a.ClientChecks.Delete(context.Background(), &lokahi.CheckID{
		Id: a.rc.Id,
	})

	a.SetErr(err)

	return nil
}

func (a *api) iTryToFetchTheCheck() error {
	ck, err := a.ClientChecks.Get(context.Background(), &lokahi.CheckID{
		Id: a.rc.Id,
	})
	a.rc = ck
	a.SetErr(err)

	return nil
}

func (a *api) iTryToListChecks() error {
	_, err := a.ClientChecks.List(context.Background(), a.checkListOpts)
	a.SetErr(err)

	return nil
}

func (a *api) iTryToPutTheCheck() error {
	ck, err := a.Suite.ClientChecks.Put(context.Background(), a.rc)
	a.rc = ck
	a.SetErr(err)

	return nil
}

func (a *api) anExampleCheck() error {
	o := &lokahi.CreateOpts{
		Url:         "https://google.com?" + uuid.New(),
		WebhookUrl:  "http://sample_hook:9001/twirp/github.xe.lokahi.Webhook/Handle",
		Every:       60,
		PlaybookUrl: "https://figureit.out",
	}

	ck, err := a.Suite.ClientChecks.Create(context.Background(), o)
	if err != nil {
		return err
	}

	a.rc = ck

	return nil
}

func (a *api) theCheckCannotBeFetched() error {
	err := a.GetErr()

	if e := err.Error(); !strings.Contains(e, sql.ErrNoRows.Error()) {
		return err
	}

	return nil
}

func (a *api) aRandomCheckID() error {
	a.rc.Id = uuid.New()

	return nil
}

func (a *api) iWantToListChecks() error {
	a.checkListOpts = &lokahi.ListOpts{}

	return nil
}

func (a *api) checkListCountIs(count int) error {
	a.checkListOpts.Count = int32(count)

	return nil
}

func (a *api) checkListOffsetIs(offset int) error {
	a.checkListOpts.Offset = int32(offset)

	return nil
}

func (a *api) iWantToCreateACheck() error {
	a.checkCreateOpts = &lokahi.CreateOpts{}

	return nil
}

func (a *api) aCheckMonitoringUrlOf(curl string) error {
	a.checkCreateOpts.Url = curl

	return nil
}

func (a *api) aCheckWebhookUrlOf(wurl string) error {
	a.checkCreateOpts.WebhookUrl = wurl

	return nil
}

func (a *api) aCheckEveryOf(every int) error {
	a.checkCreateOpts.Every = int32(every)

	return nil
}

func (a *api) aCheckPlaybookUrlOf(purl string) error {
	a.checkCreateOpts.PlaybookUrl = purl

	return nil
}

func (a *api) iCreateTheCheck() error {
	ck, err := a.ClientChecks.Create(context.Background(), a.checkCreateOpts)
	if err != nil {
		return err
	}

	a.rc = ck

	return nil
}

func (a *api) aRandomUrlInTheLastCheck() error {
	a.rc.Url = "https://google.com?" + uuid.New()

	return nil
}

func (a *api) theResultingCheckShouldHaveAnID() error {
	if a.rc.Id == "" {
		return errors.New("the check doesn't have an ID")
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	a := &api{
		Suite: &integration.Suite{},
	}

	a.Register(s)

	s.Step(`^I want to create a check$`, a.iWantToCreateACheck)
	s.Step(`^a check monitoring url of "([^"]*)"$`, a.aCheckMonitoringUrlOf)
	s.Step(`^a check webhook url of "([^"]*)"$`, a.aCheckWebhookUrlOf)
	s.Step(`^a check every of (\d+)$`, a.aCheckEveryOf)
	s.Step(`^a check playbook url of "([^"]*)"$`, a.aCheckPlaybookUrlOf)
	s.Step(`^I try to create the check$`, a.iTryToCreateTheCheck)
	s.Step(`^I can fetch the check$`, a.iTryToFetchTheCheck)
	s.Step(`^an example check$`, a.anExampleCheck)
	s.Step(`^I try to delete the check$`, a.iTryToDeleteTheCheck)
	s.Step(`^a random check ID$`, a.aRandomCheckID)
	s.Step(`^I try to fetch the check$`, a.iTryToFetchTheCheck)
	s.Step(`^the resulting check should have an ID$`, a.theResultingCheckShouldHaveAnID)
	s.Step(`^the check cannot be fetched$`, a.theCheckCannotBeFetched)
	s.Step(`^I want to list checks$`, a.iWantToListChecks)
	s.Step(`^I try to list checks$`, a.iTryToListChecks)
	s.Step(`^check list count is (\d+)$`, a.checkListCountIs)
	s.Step(`^check list offset is (\d+)$`, a.checkListOffsetIs)
	s.Step(`^a random url in the last check$`, a.aRandomUrlInTheLastCheck)
	s.Step(`^I try to put the check$`, a.iTryToPutTheCheck)
}
