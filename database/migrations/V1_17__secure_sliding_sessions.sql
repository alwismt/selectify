DELETE FROM user_session;

CREATE TABLE user_device (
    device_id           UUID PRIMARY KEY,
    user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_token_hash   VARCHAR(64) NOT NULL UNIQUE,
    user_agent          TEXT,
    last_ip             TEXT,
    first_seen_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_seen_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_device_user_id ON user_device(user_id);
CREATE INDEX idx_user_device_last_seen_at ON user_device(last_seen_at);

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE user_device TO selectify_rw;
GRANT SELECT ON TABLE user_device TO selectify_ro;
GRANT ALL ON TABLE user_device TO selectify_owner;

ALTER TABLE user_session
    ADD COLUMN session_token_hash VARCHAR(64) NOT NULL,
    ADD COLUMN last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN absolute_expires_at TIMESTAMPTZ NOT NULL,
    ADD COLUMN remember_me BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN device_id UUID REFERENCES user_device(device_id) ON DELETE SET NULL;

CREATE UNIQUE INDEX idx_user_session_token_hash ON user_session(session_token_hash);
CREATE INDEX idx_user_session_device_id ON user_session(device_id);
CREATE INDEX idx_user_session_absolute_expires_at ON user_session(absolute_expires_at);
