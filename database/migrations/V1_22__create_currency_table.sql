CREATE TABLE currency (
    code CHAR(3) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    minor_unit SMALLINT NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT currency_code_uppercase_check
      CHECK (code = UPPER(code)),

    CONSTRAINT currency_minor_unit_check
      CHECK (minor_unit >= 0)
);

INSERT INTO currency (code, name, minor_unit)
VALUES
    ('EUR', 'Euro', 2),
    ('USD', 'US Dollar', 2),
    ('GBP', 'Pound Sterling', 2),
    ('JPY', 'Japanese Yen', 0),
    ('CHF', 'Swiss Franc', 2),
    ('CAD', 'Canadian Dollar', 2),
    ('AUD', 'Australian Dollar', 2),
    ('NZD', 'New Zealand Dollar', 2);

ALTER TABLE product
    DROP COLUMN currency;

ALTER TABLE merchant
    ADD COLUMN currency_code CHAR(3) NOT NULL DEFAULT 'EUR';

ALTER TABLE merchant
    ADD CONSTRAINT merchant_currency_fk
        FOREIGN KEY (currency_code)
            REFERENCES currency(code)
            ON UPDATE CASCADE
            ON DELETE RESTRICT;

ALTER TABLE merchant
    ALTER COLUMN currency_code DROP DEFAULT;

ALTER TABLE currency
    OWNER TO selectify_owner;

GRANT SELECT, INSERT, UPDATE, DELETE
    ON TABLE currency
    TO selectify_rw;

GRANT SELECT
    ON TABLE currency
    TO selectify_ro;