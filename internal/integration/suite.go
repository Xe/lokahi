// Package integration contains some helpers for writing integration tests using
// BDD and Cucumber via godog (https://github.com/DATA-DOG/godog).
//
// At a high level, to use this package you need to create a folder under this
// folder named `component` where component is the component you want to test
// with integration testing. Then add the following in `component/main_test.go`:
//
//     package component
//
//     import (
//       "os"
//       "testing"
//       "time"
//
//       "github.com/DATA-DOG/godog"
//     )
//
//     func TestMain(m *testing.M) {
//       status := godog.RunWithOptions("godog", func(s *godog.Suite) {
// 	   FeatureContext(s)
//       }, godog.Options{
// 	   Format:    "progress",
// 	   Paths:     []string{"features"},
// 	   Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
//       })
//
//       if st := m.Run(); st > status {
// 	   status = st
//       }
//       os.Exit(status)
//     }
//
// And then run `$ mkdir component/features`. Create `component_test.go` with the
// following in it:
//
//     package component
//
//     type component struct {
//       *integration.Suite
//     }
//
//     func FeatureContext(s *godog.Suite) {
//       c := &component{
//         Suite: &integration.Suite{},
//       }
//
//       c.Register(s)
//     }
package integration

import (
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
)

// Suite is a group of story fragments. This includes some fun helpers for the
// life of this test suite as well as some setup and teardown code.
type Suite struct {
	DB *sqlx.DB // raw SQL database in case you need it

	// DAOs
	Checks   database.Checks
	Runs     database.Runs
	RunInfos database.RunInfos

	// Meta crap
	Mux *http.ServeMux
	TS  *httptest.Server

	ClientChecks lokahi.Checks

	err error
}

// Setup reads configuration information from the environment and then uses
// this to set up an integration stack of lokahi.
//
// The suggested story fragment for this function is:
//
//     Given a base stack
//
// and it can be used with a function of something like:
//
//     func FeatureContext(s *godog.Suite) {
//       val := &feature{
//         Suite: &integration.Suite{},
//       }
//
//       s.Step(`^a base stack$`, val.Suite.Setup)
//     }
func (s *Suite) Setup() error {
	durl := os.Getenv("DATABASE_URL")
	if durl == "" {
		return errors.New("no DATABASE_URL")
	}

	// Migrate database for you
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

	s.DB = db
	s.Checks = database.ChecksPostgres(db)
	s.Runs = database.RunsPostgres(db)
	s.RunInfos = database.RunInfosPostgres(db)

	// mux
	s.Mux = http.NewServeMux()
	s.Mux.Handle(lokahi.ChecksPathPrefix, lokahi.NewChecksServer(&lokahiserver.Checks{
		DB: s.Checks,
	}, nil))

	s.TS = httptest.NewServer(s.Mux)
	s.ClientChecks = lokahi.NewChecksProtobufClient(s.TS.URL, &http.Client{})

	return nil
}

// Teardown destroys all resources that Setup allocated.
//
// The suggested story fragment for this function is:
//
//     Then tear everything down
//
// and it can be used with a function of something like:
//
//     func FeatureContext(s *godog.Suite) {
//       val := &feature{
//         Suite: &integration.Suite{},
//       }
//
//       s.Step(`^tear everything down$`, val.Suite.Teardown)
//     }
func (s *Suite) Teardown() error {
	err := s.DB.Close()
	if err != nil {
		return err
	}

	s.TS.Close()

	return nil
}

// SetErr sets the suite active error to the given value. This is useful when
// writing tests kinda like `When I try to create the Task`, `Then there was no error` /
// `Then there was an error`.
func (s *Suite) SetErr(err error) {
	s.err = err
}

// GetErr fetches the suite active error.
//
// The suggested story fragement for this function is:
//
//      Then there was no error
//
// and should be followed up with an additional condition or two examining the
// error string for known values. It can be used with a FeatureContext function
// of something like:
//
//     func FeatureContext(s *godog.Suite) {
//       val := &feature{
//         Suite: &integration.Suite{},
//       }
//
//       s.Step(`^there was no error$`, val.Suite.GetErr)
//     }
//
// This doesn't have any branches because the zero value for any go value boxed
// in an interface is nil.
func (s Suite) GetErr() error {
	return s.err
}

// WantAnError asserts that the suite active error is non-nil. If the suite active
// error is nil for some reason, an error is returned.
//
// The suggested story fragment for this function is:
//
//       Then there was an error
//
// and should be followed up with an additional condition or two examining the
// error string for known values. It can be used with a FeatureContext function
// of something like:
//
//     func FeatureContext(s *godog.Suite) {
//       val := &feature{
//         Suite: &integration.Suite{},
//       }
//
//       s.Step(`^there was an error$`, val.Suite.WantAnError)
//     }
func (s Suite) WantAnError() error {
	if s.err == nil {
		return errors.New("expected an error, but there wasn't one")
	}

	return nil
}

// Register is a quick shortcut for registering all of the convenience shortcuts
// defined in the Suite type.
func (s *Suite) Register(gs *godog.Suite) {
	gs.Step(`^a base stack$`, s.Setup)
	gs.Step(`^tear everything down$`, s.Teardown)
	gs.Step(`^there was no error$`, s.GetErr)
	gs.Step(`^there was an error$`, s.WantAnError)
}
