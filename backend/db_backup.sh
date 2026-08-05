#!/usr/bin/env bash
set -euo pipefail

psql -U selectify_owner -h localhost -c "DELETE FROM user_session" selectifytestdb
pg_dump -U selectify_owner -h localhost selectifytestdb -F plain > ../database/backups/selectifytestdb.dump

