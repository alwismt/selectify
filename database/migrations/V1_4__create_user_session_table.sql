CREATE TABLE user_session (
    session_id  UUID PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    issued_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    user_agent  TEXT,
    ip_address  TEXT
);

CREATE INDEX idx_user_session_user_id ON user_session(user_id);
CREATE INDEX idx_user_session_expires_at ON user_session(expires_at);
CREATE INDEX idx_user_session_active ON user_session(expires_at) WHERE revoked_at IS NULL;
CREATE INDEX idx_user_email ON users(email);

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE user_session TO selectify_rw;
GRANT SELECT ON TABLE user_session TO selectify_ro;
GRANT ALL ON TABLE user_session TO selectify_owner;