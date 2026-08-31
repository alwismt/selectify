CREATE TABLE country (
     code CHAR(2) PRIMARY KEY,
     name VARCHAR(100) NOT NULL UNIQUE,

     is_active BOOLEAN NOT NULL DEFAULT TRUE,

     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

     CONSTRAINT country_code_uppercase_check
         CHECK (code = UPPER(code))
);

INSERT INTO country (code, name)
VALUES
    ('AT', 'Austria'),
    ('BE', 'Belgium'),
    ('BG', 'Bulgaria'),
    ('HR', 'Croatia'),
    ('CZ', 'Czechia'),
    ('DK', 'Denmark'),
    ('EE', 'Estonia'),
    ('FI', 'Finland'),
    ('FR', 'France'),
    ('DE', 'Germany'),
    ('GR', 'Greece'),
    ('HU', 'Hungary'),
    ('IS', 'Iceland'),
    ('IT', 'Italy'),
    ('LV', 'Latvia'),
    ('LI', 'Liechtenstein'),
    ('LT', 'Lithuania'),
    ('LU', 'Luxembourg'),
    ('MT', 'Malta'),
    ('NL', 'Netherlands'),
    ('NO', 'Norway'),
    ('PL', 'Poland'),
    ('PT', 'Portugal'),
    ('RO', 'Romania'),
    ('SK', 'Slovakia'),
    ('SI', 'Slovenia'),
    ('ES', 'Spain'),
    ('SE', 'Sweden'),
    ('CH', 'Switzerland');

ALTER TABLE merchant
    ADD CONSTRAINT merchant_country_fk
        FOREIGN KEY (country_code)
            REFERENCES country(code)
            ON UPDATE CASCADE
            ON DELETE RESTRICT;

ALTER TABLE country
    OWNER TO selectify_owner;

GRANT SELECT, INSERT, UPDATE, DELETE
    ON TABLE country
    TO selectify_rw;

GRANT SELECT
    ON TABLE country
    TO selectify_ro;