package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"alwis.dev/selectify/internal/db"
	"alwis.dev/selectify/internal/logger"
	"alwis.dev/selectify/internal/model"
)

type EventRepo interface {
	InsertEvent(ctx context.Context, event *model.Event) error
	GetByID(ctx context.Context, eventID string) (*model.Event, error)
	MarkProcessed(ctx context.Context, eventID string, processedAt time.Time) error
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

func (r *eventRepo) GetByID(ctx context.Context, eventID string) (*model.Event, error) {
	var row struct {
		ID            string       `db:"event_id"`
		DataJSON      []byte       `db:"event_data"`
		ReceivedDate  time.Time    `db:"event_received_date"`
		ProcessedDate sql.NullTime `db:"event_processed_date"`
	}

	err := r.roDb.QueryRowxContext(ctx, `
		SELECT event_id, event_data, event_received_date, event_processed_date
		FROM events
		WHERE event_id = $1`, eventID).StructScan(&row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, logger.Errorf(ctx, err, "failed to get event %q", eventID)
	}

	event := &model.Event{
		ID:           row.ID,
		ReceivedDate: row.ReceivedDate,
	}
	if row.ProcessedDate.Valid {
		event.ProcessedDate = &row.ProcessedDate.Time
	}

	var data model.EventData
	if err = json.Unmarshal(row.DataJSON, &data); err != nil {
		return nil, logger.Errorf(ctx, err, "failed to unmarshal event data for %q", eventID)
	}
	event.Data = &data

	return event, nil
}

func (r *eventRepo) MarkProcessed(ctx context.Context, eventID string, processedAt time.Time) error {
	result, err := r.rwDb.ExecContext(ctx, `
		UPDATE events
		SET event_processed_date = $2
		WHERE event_id = $1`, eventID, processedAt)
	if err != nil {
		return logger.Errorf(ctx, err, "failed to mark event %q as processed", eventID)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return logger.Errorf(ctx, err, "failed to read rows affected for event %q", eventID)
	}
	if rows == 0 {
		return logger.Errorf(ctx, sql.ErrNoRows, "event %q not found when marking processed", eventID)
	}

	return nil
}
