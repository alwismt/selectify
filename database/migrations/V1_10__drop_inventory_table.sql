DROP TABLE IF EXISTS inventory;

ALTER TABLE product_variants
    ADD COLUMN stock_qty    INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN reserved_qty INTEGER NOT NULL DEFAULT 0;

ALTER TABLE product_variants
    ADD CONSTRAINT product_variants_stock_chk
        CHECK (stock_qty >= 0 AND reserved_qty >= 0 AND reserved_qty <= stock_qty);

CREATE INDEX idx_product_variants_available
    ON product_variants (id)
    WHERE is_active = true AND deleted_at IS NULL;