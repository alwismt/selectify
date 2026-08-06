package repo

import (
	"context"
	"encoding/json"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type EventRepo interface {
	InsertEvent(ctx context.Context, event *model.Event) error
}

type eventRepo struct {
	rwDb *sqlx.DB
	roDb *sqlx.DB
}

func NewEventRepo(dbConn *db.DatabaseConnection) EventRepo {
	return &eventRepo{
		rwDb: dbConn.RwDb,
		roDb: dbConn.RoDb,
	}
}

func (r *eventRepo) InsertEvent(ctx context.Context, event *model.Event) error {
	eventDataJSON, err := json.Marshal(event.Data)
	if err != nil {
		return logger.Error(ctx, err, "failed to marshal event data")
	}

	query := `INSERT INTO events (
			event_id,
			event_data,
			event_received_date
		)
		VALUES ($1, $2, $3)`

	_, err = r.rwDb.ExecContext(ctx, query, event.ID, eventDataJSON, event.ReceivedDate)
	if err != nil {
		return logger.Error(ctx, err, "failed to insert event")
	}

	return nil
}
