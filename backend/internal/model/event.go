package model

import (
	"encoding/json"
	"time"

	"alwis.dev/selectify/internal/types"
)

type Event struct {
	ID            string     `db:"event_id" json:"-"`
	Data          *EventData `db:"event_data" json:"-"`
	ReceivedDate  time.Time  `db:"event_received_date" json:"-"`
	ProcessedDate *time.Time `db:"event_processed_date" json:"-"`
}

type EventData struct {
	Type    types.EventType `json:"type" valid:"required~EVENT_TYPE_REQUIRED"`
	Payload json.RawMessage `json:"data" valid:"required~EVENT_DATA_REQUIRED"`
	Date    *time.Time      `json:"date,omitempty" valid:"-"`
}
