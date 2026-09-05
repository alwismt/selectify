CREATE TABLE category (
    category_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(120) NOT NULL UNIQUE,

    parent_id BIGINT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_category_parent
      FOREIGN KEY (parent_id)
          REFERENCES category(category_id)
          ON DELETE RESTRICT
);

CREATE TABLE product_category (
    product_id BIGINT NOT NULL,
    category_id BIGINT NOT NULL,

    PRIMARY KEY (product_id, category_id),

    CONSTRAINT fk_product_category_product
      FOREIGN KEY (product_id)
          REFERENCES product(product_id)
          ON DELETE CASCADE,

    CONSTRAINT fk_product_category_category
      FOREIGN KEY (category_id)
          REFERENCES category(category_id)
          ON DELETE RESTRICT
);

CREATE INDEX idx_product_category_category
    ON product_category (category_id);

ALTER TABLE category
    OWNER TO selectify_owner;

ALTER TABLE product_category
    OWNER TO selectify_owner;

ALTER SEQUENCE category_category_id_seq
    OWNER TO selectify_owner;

GRANT SELECT, INSERT, UPDATE, DELETE
    ON TABLE category
    TO selectify_rw;

GRANT SELECT, INSERT, UPDATE, DELETE
    ON TABLE product_category
    TO selectify_rw;

GRANT USAGE, SELECT
    ON SEQUENCE category_category_id_seq
    TO selectify_rw;

GRANT SELECT
    ON TABLE category
    TO selectify_ro;

GRANT SELECT
    ON TABLE product_category
    TO selectify_ro;

GRANT ALL
    ON TABLE category
    TO selectify_owner;

GRANT ALL
    ON TABLE product_category
    TO selectify_owner;

GRANT ALL
    ON SEQUENCE category_category_id_seq
    TO selectify_owner;