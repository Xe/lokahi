package lokahiserver

import (
	"context"
	"errors"

	"github.com/Xe/lokahi/internal/database"
	"github.com/Xe/lokahi/rpc/lokahi"
)

var errNotImpl = errors.New("not implemented")

// Checks implements service Checks.
type Checks struct {
	DB database.Checks
}

// Create creates a new health check with the given options.
func (c *Checks) Create(ctx context.Context, opts *lokahi.CreateOpts) (*lokahi.Check, error) {
	dck := database.Check{
		URL:         opts.Url,
		WebhookURL:  opts.WebhookUrl,
		Every:       int(opts.Every),
		PlaybookURL: opts.PlaybookUrl,
		State:       lokahi.Check_INIT.String(),
	}

	result, err := c.DB.Create(ctx, dck)
	if err != nil {
		return nil, err
	}

	return result.AsProto(), nil
}

// Delete removes a check by ID and returns the data that was deleted.
func (c *Checks) Delete(ctx context.Context, cid *lokahi.CheckID) (*lokahi.Check, error) {
	_, err := c.DB.Get(ctx, cid.Id)
	if err != nil {
		return nil, err
	}

	result, err := c.DB.Delete(ctx, cid.Id)
	if err != nil {
		return nil, err
	}

	return result.AsProto(), nil
}

// Get returns information on a check by ID.
func (c *Checks) Get(ctx context.Context, cid *lokahi.CheckID) (*lokahi.Check, error) {
	dck, err := c.DB.Get(ctx, cid.Id)
	if err != nil {
		return nil, err
	}

	return dck.AsProto(), nil
}

// List returns a page of checks based on a few options.
func (c *Checks) List(ctx context.Context, opts *lokahi.ListOpts) (*lokahi.ChecksPage, error) {
	dchecks, err := c.DB.List(ctx, int(opts.Count), int(opts.Offset))
	if err != nil {
		return nil, err
	}

	result := &lokahi.ChecksPage{}

	for _, dc := range dchecks {
		result.Results = append(result.Results, &lokahi.ChecksPage_Result{Check: dc.AsProto()})
	}

	return result, nil
}

// Put updates a Check.
func (c *Checks) Put(ctx context.Context, chk *lokahi.Check) (*lokahi.Check, error) {
	dck, err := c.DB.Get(ctx, chk.Id)
	if err != nil {
		return nil, err
	}

	if a := chk.Url; a != "" {
		dck.URL = a
	}

	if a := dck.WebhookURL; a != "" {
		dck.WebhookURL = chk.WebhookUrl
	}

	if a := dck.PlaybookURL; a != "" {
		dck.PlaybookURL = chk.PlaybookUrl
	}

	result, err := c.DB.Put(ctx, *dck)
	if err != nil {
		return nil, err
	}

	return result.AsProto(), nil
}

// Status returns the detailed histogram status of a check.
func (c *Checks) Status(ctx context.Context, cid *lokahi.CheckID) (*lokahi.CheckStatus, error) {
	return nil, errNotImpl
}
