ALTER TABLE product
    ADD COLUMN IF NOT EXISTS merchant_id BIGINT NOT NULL DEFAULT 2;

ALTER TABLE product
    ADD CONSTRAINT product_merchant_fk
        FOREIGN KEY (merchant_id)
            REFERENCES merchant(merchant_id)
            ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_product_merchant_id
    ON product (merchant_id);

ALTER TABLE product
    ALTER COLUMN merchant_id DROP DEFAULT;