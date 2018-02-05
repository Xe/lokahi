package database

import (
	"context"
	"time"

	"github.com/Xe/ln"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/jmoiron/sqlx"
)

// Checks is the set of calls that can be made to the database regarding checks
type Checks interface {
	Create(ctx context.Context, c Check) (*Check, error)
	Delete(ctx context.Context, cid string) (*Check, error)
	Get(ctx context.Context, cid string) (*Check, error)
	ListByEveryValue(ctx context.Context, time int) ([]Check, error)
	List(ctx context.Context, count, offset int) ([]Check, error)
	Put(ctx context.Context, c Check) (*Check, error)
}

type checksPostgres struct {
	db *sqlx.DB
}

func (c *checksPostgres) Create(ctx context.Context, ch Check) (*Check, error) {
	_, err := c.db.NamedExec(`INSERT INTO checks(url, webhook_url, playbook_url, every) VALUES (:url, :webhook_url, :playbook_url, :every)`, ch)
	if err != nil {
		return nil, err
	}

	var result Check
	err = c.db.Get(&result, "SELECT * FROM checks WHERE url = $1", ch.URL)
	if err != nil { // unlikely
		return nil, err
	}

	return &result, nil
}

func (c *checksPostgres) Get(ctx context.Context, cid string) (*Check, error) {
	var result Check
	err := c.db.Get(&result, "SELECT * FROM checks WHERE uuid = $1", cid)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *checksPostgres) Delete(ctx context.Context, cid string) (*Check, error) {
	result, err := c.Get(ctx, cid)
	if err != nil {
		return nil, err
	}

	_, err = c.db.Exec(`DELETE FROM checks WHERE uuid = $1`, cid)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *checksPostgres) List(ctx context.Context, count, offset int) ([]Check, error) {
	var result []Check

	err := c.db.Select(&result, `SELECT * FROM checks LIMIT $1 OFFSET $2`, count, offset)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *checksPostgres) Put(ctx context.Context, ch Check) (*Check, error) {
	ch.EditedAt = time.Now()

	_, err := c.db.NamedExec(`UPDATE checks SET edited_at = :edited_at, url = :url, webhook_url = :webhook_url, webhook_response_time_nanoseconds = :webhook_response_time_nanoseconds, every = :every, playbook_url = :playbook_url, state = :state WHERE uuid = :uuid`, ch)
	if err != nil {
		return nil, err
	}

	return c.Get(ctx, ch.UUID)
}

func (c *checksPostgres) ListByEveryValue(ctx context.Context, every int) ([]Check, error) {
	var result []Check

	err := c.db.Select(&result, `SELECT * FROM CHECKS WHERE every = $1`, every)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Check is an individual HTTP check that gets scheduled every so often.
type Check struct {
	ID   int    `db:"id"`
	UUID string `db:"uuid"`

	CreatedAt time.Time `db:"created_at"`
	EditedAt  time.Time `db:"edited_at"`

	// URL is the URL that this check will monitor.
	URL string `db:"url"`
	// WebhookURL is the URL that state changes will be POSTed to.
	WebhookURL string `db:"webhook_url"`
	// WebhookResponseTimeNanoseconds is the last response time of the
	// webhook URL in nanoseconds.
	WebhookResponseTimeNanoseconds int64 `db:"webhook_response_time_nanoseconds"`
	// Every is how often this check will run in seconds.
	// (minimum: 60, maximum: 600)
	Every int `db:"every"`
	// PlaybookURL is a URL for operations staff to respond to downtime of
	// the service that is backed by the URL above.
	PlaybookURL string `db:"playbook_url"`
	// State is the check's last known state.
	State string `db:"state"`
}

// F ields for logging.
func (c Check) F() ln.F {
	return ln.F{
		"check_uuid":  c.UUID,
		"check_url":   c.URL,
		"check_every": c.Every,
		"check_state": c.State,
	}
}

// AsProto converts this to the protobuf representation.
func (c Check) AsProto() *lokahi.Check {
	return &lokahi.Check{
		Id:          c.UUID,
		Url:         c.URL,
		WebhookUrl:  c.WebhookURL,
		Every:       int32(c.Every),
		PlaybookUrl: c.PlaybookURL,
		State:       lokahi.Check_State(lokahi.Check_State_value[c.State]),

		WebhookResponseTimeNanoseconds: c.WebhookResponseTimeNanoseconds,
	}
}
