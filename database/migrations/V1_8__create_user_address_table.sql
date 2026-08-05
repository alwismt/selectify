CREATE TABLE user_addresses (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id),
    label         TEXT NULL, -- Home, Office, etc.
    phone         TEXT NULL,
    line1         TEXT NOT NULL,
    line2         TEXT NULL,
    city          TEXT NOT NULL,
    region        TEXT NULL,
    postal_code   TEXT NOT NULL,
    country_code  CHAR(2) NOT NULL,
    is_default    BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_addresses_user_id ON user_addresses(user_id);

CREATE TABLE order_addresses (
     id            BIGSERIAL PRIMARY KEY,
     order_id      BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
     type          TEXT NOT NULL CHECK (type IN ('shipping','billing')),
     phone         TEXT NULL,
     line1         TEXT NOT NULL,
     line2         TEXT NULL,
     city          TEXT NOT NULL,
     region        TEXT NULL,
     postal_code   TEXT NOT NULL,
     country_code  CHAR(2) NOT NULL,

     created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX ux_order_addresses_order_type
    ON order_addresses(order_id, type);
