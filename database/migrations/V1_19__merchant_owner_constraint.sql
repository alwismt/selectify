CREATE UNIQUE INDEX user_role_one_owner_per_merchant
    ON user_role (merchant_id)
    WHERE merchant_role = 'owner';

ALTER TABLE merchant OWNER TO selectify_owner;
ALTER TABLE user_role OWNER TO selectify_owner;
ALTER SEQUENCE merchant_merchant_id_seq OWNER TO selectify_owner;

GRANT SELECT, INSERT, UPDATE, DELETE
    ON TABLE merchant
    TO selectify_rw;

GRANT SELECT, INSERT, UPDATE, DELETE
    ON TABLE user_role
    TO selectify_rw;

GRANT USAGE, SELECT
    ON SEQUENCE merchant_merchant_id_seq
    TO selectify_rw;

GRANT SELECT
    ON TABLE merchant
    TO selectify_ro;

GRANT SELECT
    ON TABLE user_role
    TO selectify_ro;

GRANT SELECT
    ON SEQUENCE merchant_merchant_id_seq
    TO selectify_ro;