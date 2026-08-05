DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'selectify_owner') THEN
    CREATE ROLE selectify_owner LOGIN PASSWORD 'passVVord';
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'selectify_rw') THEN
    CREATE ROLE selectify_rw LOGIN PASSWORD 'passVVord';
  END IF;

  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'selectify_ro') THEN
    CREATE ROLE selectify_ro LOGIN PASSWORD 'passVVrd';
  END IF;
END$$;

-- Create database if it doesn't exist (must be outside transaction -> use \gexec trick)
SELECT format('CREATE DATABASE %I OWNER %I', 'selectify', 'selectify_owner')
WHERE NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'selectify')
\gexec

\connect selectify

ALTER SCHEMA public OWNER TO selectify_owner;

GRANT CONNECT ON DATABASE selectify TO selectify_ro, selectify_rw;

REVOKE ALL ON SCHEMA public FROM PUBLIC;

GRANT USAGE ON SCHEMA public TO selectify_ro, selectify_rw;

GRANT CREATE ON SCHEMA public TO selectify_rw;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO selectify_ro;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO selectify_rw;

GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO selectify_ro;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO selectify_rw;

GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO selectify_ro, selectify_rw;

ALTER DEFAULT PRIVILEGES FOR ROLE selectify_owner IN SCHEMA public
  GRANT SELECT ON TABLES TO selectify_ro;

ALTER DEFAULT PRIVILEGES FOR ROLE selectify_owner IN SCHEMA public
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO selectify_rw;

ALTER DEFAULT PRIVILEGES FOR ROLE selectify_owner IN SCHEMA public
  GRANT USAGE, SELECT ON SEQUENCES TO selectify_ro;

ALTER DEFAULT PRIVILEGES FOR ROLE selectify_owner IN SCHEMA public
  GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO selectify_rw;

ALTER DEFAULT PRIVILEGES FOR ROLE selectify_owner IN SCHEMA public
  GRANT EXECUTE ON FUNCTIONS TO selectify_ro, selectify_rw;