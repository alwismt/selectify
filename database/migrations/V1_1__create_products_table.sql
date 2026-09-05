CREATE TABLE product (
    product_id BIGSERIAL PRIMARY KEY,
    sku             TEXT UNIQUE NOT NULL,

    name            TEXT NOT NULL,
    description     TEXT,
    slug            TEXT NOT NULL,

    price_amount    BIGINT NOT NULL,
    currency        TEXT NOT NULL DEFAULT 'EUR',

    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    in_stock        BOOLEAN NOT NULL DEFAULT TRUE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE INDEX idx_products_active
    ON product (is_active)
    WHERE deleted_at IS NULL;

CREATE INDEX idx_products_name_search
    ON product
        USING GIN (to_tsvector('simple', name));

CREATE INDEX idx_products_sku
    ON product (sku);

CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_products_updated_at
    BEFORE UPDATE ON product
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

ALTER DEFAULT PRIVILEGES FOR ROLE selectify_owner
    GRANT USAGE, SELECT ON SEQUENCES TO selectify_rw;

ALTER DEFAULT PRIVILEGES FOR ROLE selectify_owner
    GRANT USAGE, SELECT ON SEQUENCES TO selectify_ro;

GRANT USAGE, SELECT ON SEQUENCE product_product_id_seq TO selectify_rw;
GRANT INSERT, SELECT, UPDATE, DELETE ON product TO selectify_rw;
GRANT SELECT ON product TO selectify_ro;
GRANT ALL ON TABLE product TO selectify_owner;

GRANT EXECUTE ON FUNCTION set_updated_at() TO selectify_rw;
GRANT EXECUTE ON FUNCTION set_updated_at() TO selectify_ro;

