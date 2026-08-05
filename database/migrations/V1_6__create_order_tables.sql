CREATE TABLE orders (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id),
    status        TEXT NOT NULL DEFAULT 'pending', -- pending|paid|cancelled|shipped|refunded
    currency      TEXT NOT NULL,

    subtotal      NUMERIC(12,2) NOT NULL,
    shipping      NUMERIC(12,2) NOT NULL DEFAULT 0,
    discount      NUMERIC(12,2) NOT NULL DEFAULT 0,
    total         NUMERIC(12,2) NOT NULL,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_items (
     id            BIGSERIAL PRIMARY KEY,
     order_id      BIGINT NOT NULL,
     variant_id    BIGINT NOT NULL,

     sku           TEXT NOT NULL,
     unit_price    NUMERIC(12,2) NOT NULL,
     currency      CHAR(3) NOT NULL,
     quantity      INT NOT NULL CHECK (quantity > 0),

     attributes    JSONB NOT NULL, -- snapshot
     created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

     CONSTRAINT order_items_order_fk
         FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,

     CONSTRAINT order_items_variant_fk
         FOREIGN KEY (variant_id) REFERENCES product_variants(id)
);
