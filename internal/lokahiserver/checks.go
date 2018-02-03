package lokahiserver

import (
	"context"
	"errors"

	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/uuid"
	"github.com/jinzhu/gorm"
)

var errNotImpl = errors.New("not implemented")

// Checks implements service Checks.
type Checks struct {
	DB *gorm.DB
}

// Create creates a new health check with the given options.
func (c *Checks) Create(ctx context.Context, opts *lokahi.CreateOpts) (*lokahi.Check, error) {
	dck := database.Check{
		UUID:        uuid.New(),
		URL:         opts.WebhoolUrl,
		Every:       opts.Every,
		PlaybookURL: opts.PlaybookURL,
		State:       lokahi.Check_State_INIT,
	}

	err := c.DB.Create(&dck)
	if err != nil {
		return nil, err
	}

	return dck.AsProto(), nil
}

// Delete removes a check by ID and returns the data that was deleted.
func (c *Checks) Delete(ctx context.Context, cid *lokahi.CheckID) (*lokahi.Check, error) {
	dck, err := c.getCheck(ctx, cid.Id)
	if err != nil {
		return nil, err
	}

	err = c.DB.Where("uuid = ?", dck.UUID).Delete(database.Check{})
	if err != nil {
		return nil, err
	}

	return dck.AsProto(), nil
}

// Get returns information on a check by ID.
func (c *Checks) Get(ctx context.Context, cid *lokahi.CheckID) (*lokahi.Check, error) {
	dck, err := c.getCheck(ctx, cid.Id)
	if err != nil {
		return nil, err
	}

	return dck.AsProto(), nil
}

// getCheck gets a check from the database
func (c *Checks) getCheck(ctx context.Context, id string) (*database.Check, error) {
	var ck database.Check
	err := c.DB.Where("uuid = ?", id).First(&ck)
	if err != nil {
		return nil, err
	}

	return &ck, nil
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
