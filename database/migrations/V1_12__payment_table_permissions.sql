GRANT INSERT, SELECT, UPDATE, DELETE
    ON TABLE payments
    TO selectify_rw;

GRANT SELECT
    ON TABLE payments
    TO selectify_ro;

GRANT ALL
    ON TABLE payments
    TO selectify_owner;