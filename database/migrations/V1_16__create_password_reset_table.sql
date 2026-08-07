CREATE TABLE password_reset (
    password_reset_id UUID PRIMARY KEY,
    user_id BIGSERIAL NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_password_reset_user
        FOREIGN KEY (user_id)
            REFERENCES "users" (id)
            ON DELETE CASCADE
);

CREATE INDEX idx_password_reset_user_id
    ON password_reset (user_id);

CREATE INDEX idx_password_reset_expires_at
    ON password_reset (expires_at);


GRANT INSERT, SELECT, UPDATE, DELETE
    ON TABLE password_reset
    TO selectify_rw;

GRANT SELECT
    ON TABLE password_reset
    TO selectify_ro;

GRANT ALL
    ON TABLE password_reset
    TO selectify_owner;