package database

import (
	"context"
	"time"

	"github.com/Xe/uuid"
	"github.com/jmoiron/sqlx"
)

type Runs interface {
	Get(ctx context.Context, rid string) (*Run, error)
	GetForCheck(ctx context.Context, cid string, limit, offset int) ([]Run, error)
	Put(ctx context.Context, rn Run) (*Run, error)
}

type Run struct {
	ID        int       `db:"id"`
	UUID      string    `db:"uuid"`
	CreatedAt time.Time `db:"created_at"`
	Message   string    `db:"message"`
}

type runsPostgres struct {
	db *sqlx.DB
}

func (r *runsPostgres) Get(ctx context.Context, rid string) (*Run, error) {
	var result Run
	err := r.db.Get(&result, "SELECT * FROM runs WHERE uuid = $1", rid)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *runsPostgres) GetForCheck(ctx context.Context, cid string, limit, offset int) ([]Run, error) {
	var result []Run
	err := r.db.Select(&result, "SELECT * FROM runs WHERE $1 = ANY(check_ids) LIMIT $2 OFFSET $3", cid, limit, offset)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *runsPostgres) Put(ctx context.Context, rn Run) (*Run, error) {
	if rn.UUID == "" {
		rn.UUID = uuid.New()
	}

	_, err := r.db.NamedExec(`INSERT INTO runs (uuid, message) VALUES (:uuid, :message)`, rn)
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, rn.UUID)
}

type RunInfo struct {
	ID                             int       `db:"id"`
	UUID                           string    `db:"uuid"`
	CreatedAt                      time.Time `db:"created_at"`
	CheckID                        string    `db:"check_id"`
	ResponseTimeNanoseconds        int64     `db:"response_time_nanoseconds"`
	WebhookResponseTimeNanoseconds int64     `db:"webhook_response_time_nanoseconds"`
}
