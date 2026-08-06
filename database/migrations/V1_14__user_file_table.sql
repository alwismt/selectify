CREATE TABLE user_file (
    user_file_id       UUID PRIMARY KEY,
    user_id            BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_type       TEXT NOT NULL,
    created_at         TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP NOT NULL DEFAULT NOW()
);

GRANT INSERT, SELECT, UPDATE, DELETE
    ON TABLE user_file
    TO selectify_rw;

GRANT SELECT
    ON TABLE user_file
    TO selectify_ro;

GRANT ALL
    ON TABLE user_file
    TO selectify_owner;

ALTER TABLE user_file
    ADD CONSTRAINT uq_user_file_user UNIQUE (user_id);