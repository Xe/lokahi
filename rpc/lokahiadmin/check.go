package lokahiadmin

import (
	"time"

	"github.com/Xe/lokahi/internal/database"
)

func CheckFromDatabaseCheck(dck database.Check) *Check {
	return &Check{
		Id:            int32(dck.ID),
		Uuid:          dck.UUID,
		CreatedAtUnix: dck.CreatedAt.Unix(),
		EditedAtUnix:  dck.EditedAt.Unix(),
		Url:           dck.URL,
		WebhookUrl:    dck.WebhookURL,
		Every:         int32(dck.Every),
		PlaybookUrl:   dck.PlaybookURL,
		State:         dck.State,
	}
}

func (c Check) DatabaseCheck() database.Check {
	return database.Check{
		ID:          int(c.Id),
		UUID:        c.Uuid,
		CreatedAt:   time.Unix(c.CreatedAtUnix, 0),
		EditedAt:    time.Unix(c.EditedAtUnix, 0),
		URL:         c.Url,
		WebhookURL:  c.WebhookUrl,
		Every:       int(c.Every),
		PlaybookURL: c.PlaybookUrl,
		State:       c.State,
	}
}
