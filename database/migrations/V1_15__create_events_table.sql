CREATE TABLE events (
    event_id              UUID PRIMARY KEY,
    event_data            JSONB NOT NULL,
    event_received_date   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    event_processed_date  TIMESTAMPTZ
);

CREATE INDEX idx_events_received_date
    ON events (event_received_date);

CREATE INDEX idx_events_type
    ON events ((event_data->>'type'));

GRANT INSERT, SELECT, UPDATE, DELETE
    ON TABLE events
    TO selectify_rw;

GRANT SELECT
    ON TABLE events
    TO selectify_ro;

GRANT ALL
    ON TABLE events
    TO selectify_owner;
