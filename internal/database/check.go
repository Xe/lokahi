package database

import "github.com/jinzhu/gorm"

// Check is an individual HTTP check that gets scheduled every so often.
type Check struct {
	gorm.Model

	// UUID is the unique identifier of this check.
	UUID string `gorm:"unique"`
	// URL is the URL that this check will monitor.
	URL string `gorm:"unique"`
	// WebhookURL is the URL that state changes will be POSTed to.
	WebhookURL string
	// WebhookLastResponseTimeNanoseconds is the last response time of the
	// webhook URL in nanoseconds.
	WebhoolLastResponseTimeNanoseconds int64
	// Every is how often this check will run in seconds.
	// (minimum: 60, maximum: 600)
	Every int
	// PlaybookURL is a URL for operations staff to respond to downtime of
	// the service that is backed by the URL above.
	PlayboolURL string
}
