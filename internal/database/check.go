package database

import (
	"github.com/Xe/ln"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/jinzhu/gorm"
)

// Check is an individual HTTP check that gets scheduled every so often.
type Check struct {
	gorm.Model

	// UUID is the unique identifier of this check.
	UUID string `gorm:"unique"`
	// URL is the URL that this check will monitor.
	URL string `gorm:"unique"`
	// WebhookURL is the URL that state changes will be POSTed to.
	WebhookURL string
	// WebhookResponseTimeNanoseconds is the last response time of the
	// webhook URL in nanoseconds.
	WebhookResponseTimeNanoseconds int64
	// Every is how often this check will run in seconds.
	// (minimum: 60, maximum: 600)
	Every int
	// PlaybookURL is a URL for operations staff to respond to downtime of
	// the service that is backed by the URL above.
	PlaybookURL string
	// State is the check's last known state.
	State string
}

func (c Check) F() ln.F {
	return ln.F{
		"check_id":    c.ID,
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
		State:       lokahi.Check_INIT,

		WebhookResponseTimeNanoseconds: c.WebhookResponseTimeNanoseconds,
	}
}
