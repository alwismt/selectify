ALTER TABLE merchant
    DROP CONSTRAINT IF EXISTS merchant_currency_fk;

ALTER TABLE merchant
    DROP COLUMN IF EXISTS currency_code;

ALTER TABLE currency
    ADD COLUMN IF NOT EXISTS is_default BOOLEAN NOT NULL DEFAULT FALSE;

CREATE UNIQUE INDEX IF NOT EXISTS uq_currency_default
    ON currency (is_default)
    WHERE is_default = TRUE;

UPDATE currency
SET is_default = FALSE
WHERE is_default = TRUE;

UPDATE currency
SET is_default = TRUE
WHERE code = 'EUR';