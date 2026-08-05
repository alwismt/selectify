CREATE TABLE users (
    id                BIGSERIAL PRIMARY KEY,
    email             TEXT NOT NULL UNIQUE,
    password_hash     TEXT NOT NULL,
    first_name        TEXT NOT NULL,
    last_name         TEXT NOT NULL,
    phone             TEXT NOT NULL UNIQUE,
    status            TEXT NOT NULL,
    email_verified_at TIMESTAMPTZ,
    last_login_at     TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);


CREATE TABLE user_role (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role        TEXT NOT NULL,
    scope_type  TEXT,
    scope_id    BIGINT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uq_user_role_scope ON user_role(user_id, role, scope_type, scope_id);
CREATE INDEX idx_user_role_user_id ON user_role(user_id);
CREATE INDEX idx_user_role_role ON user_role(role);
CREATE INDEX idx_user_role_scope ON user_role(scope_type, scope_id);


GRANT USAGE, SELECT ON SEQUENCE users_id_seq TO selectify_rw;
GRANT USAGE, SELECT ON SEQUENCE user_role_id_seq TO selectify_rw;

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE users TO selectify_rw;
GRANT SELECT ON TABLE users TO selectify_ro;

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE user_role TO selectify_rw;
GRANT SELECT ON TABLE user_role TO selectify_ro;

GRANT ALL ON TABLE users TO selectify_owner;
GRANT ALL ON TABLE user_role TO selectify_owner;