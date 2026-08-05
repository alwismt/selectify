#!/usr/bin/env bash
set -euo pipefail

# Directory of this script
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Where backups are stored (alongside migrations directory)
BACKUP_DIR="$SCRIPT_DIR/../database/backups"

export PGHOST="${PGHOST:-localhost}"
export PGPORT="${PGPORT:-5432}"
export PGUSER="${PGUSER:-selectify_owner}"
export PGPASSWORD="${PGPASSWORD:-passVVord}"

# Name of your local test DB to restore
TEST_DB_NAME="${TEST_DB_NAME:-selectifytestdb}"

BACKUP_FILE="$BACKUP_DIR/${TEST_DB_NAME}.dump"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file '$BACKUP_FILE' not found."
    echo "Available backups:"
    ls -lh "$BACKUP_DIR"/*.dump 2>/dev/null || echo "  (no backups found)"
    exit 1
fi

echo "Restoring database '$TEST_DB_NAME' from '$BACKUP_FILE'"
echo "  host=$PGHOST port=$PGPORT user=$PGUSER"


# Force close existing connections before dropping
psql \
  --dbname="postgres" \
  --command="SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$TEST_DB_NAME' AND pid <> pg_backend_pid();" \
  > /dev/null 2>&1 || true

# Connect to 'postgres' database to drop/create the target database
psql \
  --dbname="postgres" \
  --command="DROP DATABASE IF EXISTS \"$TEST_DB_NAME\";" \
  --command="CREATE DATABASE \"$TEST_DB_NAME\";"

# Now restore into the target database
psql \
  --dbname="$TEST_DB_NAME" \
  --file="$BACKUP_FILE"

echo "Restore completed."