CREATE TABLE checkout_session (
    id         BIGSERIAL PRIMARY KEY,
    order_id   BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status     TEXT NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE checkout_session_items (
    checkout_session_id BIGINT NOT NULL REFERENCES checkout_session(id) ON DELETE CASCADE,
    cart_item_id        BIGINT NOT NULL REFERENCES cart_items(id) ON DELETE SET NULL,
    PRIMARY KEY (checkout_session_id, cart_item_id)
);