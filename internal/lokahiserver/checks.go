package lokahiserver

import (
	"context"
	"errors"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/jinzhu/gorm"
)

var errNotImpl = errors.New("not implemented")

// Checks implements service Checks.
type Checks struct {
	DB *gorm.DB
}

// Create creates a new health check with the given options.
func (c *Checks) Create(ctx context.Context, opts *lokahi.CreateOpts) (*lokahi.Check, error) {
	return nil, errNotImpl
}

// Delete removes a check by ID and returns the data that was deleted.
func (c *Checks) Delete(ctx context.Context, cid *lokahi.CheckID) (*lokahi.Check, error) {
	return nil, errNotImpl
}

// List returns a page of checks based on a few options.
func (c *Checks) List(ctx context.Context, opts *lokahi.ListOpts) (*lokahi.ChecksPage, error) {
	return nil, errNotImpl
}

// Put updates a Check.
func (c *Checks) Put(ctx context.Context, chk *lokahi.Check) (*lokahi.Check, error) {
	return nil, errNotImpl
}

// Status returns the detailed histogram status of a check.
func (c *Checks) Status(ctx context.Context, cid *lokahi.CheckID) (*lokahi.CheckStatus, error) {
	return nil, errNotImpl
}
