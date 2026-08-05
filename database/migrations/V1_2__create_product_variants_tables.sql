CREATE TABLE product_variants (
    id            BIGSERIAL PRIMARY KEY,
    product_id    BIGINT NOT NULL REFERENCES product(product_id) ON DELETE CASCADE,

    sku           TEXT UNIQUE NOT NULL,

    price_amount  NUMERIC(12,2),
    currency      TEXT NOT NULL DEFAULT 'EUR',

    is_active     BOOLEAN NOT NULL DEFAULT TRUE,

    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX idx_product_variants_product_id
    ON product_variants(product_id);

CREATE TABLE product_variant_attributes (
    id         BIGSERIAL PRIMARY KEY,
    variant_id BIGINT NOT NULL REFERENCES product_variants(id) ON DELETE CASCADE,

    name       TEXT NOT NULL,
    value      TEXT NOT NULL
);

CREATE INDEX idx_variant_attrs_variant_id
    ON product_variant_attributes(variant_id);

CREATE UNIQUE INDEX uq_variant_attrs
    ON product_variant_attributes(variant_id, name);


CREATE TABLE inventory (
    variant_id   BIGINT PRIMARY KEY REFERENCES product_variants(id) ON DELETE CASCADE,
    stock_qty    INTEGER NOT NULL DEFAULT 0 CHECK (stock_qty >= 0),
    reserved_qty INTEGER NOT NULL DEFAULT 0 CHECK (reserved_qty >= 0),

    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

GRANT USAGE, SELECT ON SEQUENCE product_variants_id_seq TO selectify_rw;
GRANT USAGE, SELECT ON SEQUENCE product_variant_attributes_id_seq TO selectify_rw;

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE product_variants TO selectify_rw;
GRANT SELECT ON TABLE product_variants TO selectify_ro;

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE product_variant_attributes TO selectify_rw;
GRANT SELECT ON TABLE product_variant_attributes TO selectify_ro;

GRANT INSERT, SELECT, UPDATE, DELETE ON TABLE inventory TO selectify_rw;
GRANT SELECT ON TABLE inventory TO selectify_ro;

GRANT ALL ON TABLE product_variants TO selectify_owner;
GRANT ALL ON TABLE product_variant_attributes TO selectify_owner;
GRANT ALL ON TABLE inventory TO selectify_owner;

