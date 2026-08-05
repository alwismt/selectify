ALTER TABLE product_variants
    ADD CONSTRAINT uq_product_variants_product_id_id
        UNIQUE (product_id, id);

CREATE TABLE product_file (
    product_file_id UUID PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES product(product_id) ON DELETE CASCADE,
    variant_id BIGINT,
    position INTEGER NOT NULL DEFAULT 0 CHECK (position >= 0),
    content_type TEXT NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_product_file_variant
      FOREIGN KEY (product_id, variant_id)
          REFERENCES product_variants(product_id, id)
          ON DELETE CASCADE
);

CREATE INDEX idx_product_file_product
    ON product_file (product_id, position);

CREATE INDEX idx_product_file_variant
    ON product_file (variant_id, position)
    WHERE variant_id IS NOT NULL;

CREATE UNIQUE INDEX uq_product_file_primary
    ON product_file (product_id)
    WHERE is_primary = TRUE
        AND variant_id IS NULL;

CREATE UNIQUE INDEX uq_variant_file_primary
    ON product_file (variant_id)
    WHERE is_primary = TRUE
        AND variant_id IS NOT NULL;

CREATE UNIQUE INDEX uq_product_file_position
    ON product_file (product_id, position)
    WHERE variant_id IS NULL;

CREATE UNIQUE INDEX uq_variant_file_position
    ON product_file (variant_id, position)
    WHERE variant_id IS NOT NULL;


GRANT INSERT, SELECT, UPDATE, DELETE
    ON TABLE product_file
    TO selectify_rw;

GRANT SELECT
    ON TABLE product_file
    TO selectify_ro;

GRANT ALL
    ON TABLE product_file
    TO selectify_owner;