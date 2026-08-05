#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SQL_FILE="$SCRIPT_DIR/../database/migrations/init.sql"

export PGHOST="${PGHOST:-localhost}"
export PGPORT="${PGPORT:-5432}"
export PGUSER="${PGUSER:-postgres}"
export PGDATABASE="${PGDATABASE:-postgres}"
# export PGPASSWORD="${PGPASSWORD:-postgres}"

echo "Running Selectify DB init against $PGHOST:$PGPORT as $PGUSER (db=$PGDATABASE)"
psql -v ON_ERROR_STOP=1 -f "$SQL_FILE"
echo "Done."