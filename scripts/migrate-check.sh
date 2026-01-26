#!/bin/sh
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)
#
# migrate-check: migration lint and parity validation
#
# This script is POSIX sh compatible for maximum portability.
# It validates the migration layout and rollback rules before applying migrations.
#
# Checks performed:
# - Required directories exist (shared/sqlite/postgres migrations, schema dumps)
# - Every *.up.sql has a matching *.down.sql and vice versa
# - shared/sqlite/postgres migration files match exactly (timestamp + slug)
# - Destructive migrations require an IRREVERSIBLE header
# - Seed down migrations delete by public_id (no blanket deletes)
# - Schema dumps exist for the latest migration timestamp

set -eu

# Resolve repository root for relative paths.
ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)

MIGRATIONS_SHARED=${MIGRATIONS_SHARED:-"$ROOT_DIR/db/migrations/shared"}
MIGRATIONS_SQLITE=${MIGRATIONS_SQLITE:-"$ROOT_DIR/db/migrations/sqlite"}
MIGRATIONS_POSTGRES=${MIGRATIONS_POSTGRES:-"$ROOT_DIR/db/migrations/postgres"}
SCHEMA_SQLITE=${SCHEMA_SQLITE:-"$ROOT_DIR/db/schema/sqlite"}
SCHEMA_POSTGRES=${SCHEMA_POSTGRES:-"$ROOT_DIR/db/schema/postgres"}

fail() {
  echo "migrate-check: $*" >&2
  exit 1
}

require_dir() {
  dir=$1
  [ -d "$dir" ] || fail "missing directory: $dir"
}

# Ensure required directories exist.
require_dir "$MIGRATIONS_SHARED"
require_dir "$MIGRATIONS_SQLITE"
require_dir "$MIGRATIONS_POSTGRES"
require_dir "$SCHEMA_SQLITE"
require_dir "$SCHEMA_POSTGRES"

check_pairs() {
  dir=$1
  for up in "$dir"/*.up.sql; do
    [ -e "$up" ] || continue
    base=${up%.up.sql}
    [ -f "$base.down.sql" ] || fail "missing down migration for $up"
  done
  for down in "$dir"/*.down.sql; do
    [ -e "$down" ] || continue
    base=${down%.down.sql}
    [ -f "$base.up.sql" ] || fail "missing up migration for $down"
  done
}

# Verify every migration has a matching up/down pair.
check_pairs "$MIGRATIONS_SHARED"
check_pairs "$MIGRATIONS_SQLITE"
check_pairs "$MIGRATIONS_POSTGRES"

# Ensure shared migrations exist in both sqlite and postgres.
for shared in "$MIGRATIONS_SHARED"/*.up.sql "$MIGRATIONS_SHARED"/*.down.sql; do
  [ -e "$shared" ] || continue
  file=$(basename "$shared")
  [ -f "$MIGRATIONS_SQLITE/$file" ] || fail "missing sqlite migration: $file"
  [ -f "$MIGRATIONS_POSTGRES/$file" ] || fail "missing postgres migration: $file"
done

# Enforce IRREVERSIBLE header for destructive statements and parity with shared.
for dbdir in "$MIGRATIONS_SQLITE" "$MIGRATIONS_POSTGRES"; do
  for up in "$dbdir"/*.up.sql; do
    [ -e "$up" ] || continue
    file=$(basename "$up")
    [ -f "$MIGRATIONS_SHARED/$file" ] || fail "extra migration not in shared: $file"
    if grep -Eiq '\b(drop|truncate)\b' "$up"; then
      if ! grep -Eiq '^--[[:space:]]*IRREVERSIBLE:' "$up"; then
        fail "missing IRREVERSIBLE header in $file"
      fi
    fi
  done
done

# Seed rollbacks must delete by public_id, not blanket deletes.
for down in "$MIGRATIONS_SHARED"/*_seed*.down.sql "$MIGRATIONS_SQLITE"/*_seed*.down.sql "$MIGRATIONS_POSTGRES"/*_seed*.down.sql; do
  [ -e "$down" ] || continue
  if ! grep -q "public_id" "$down"; then
    fail "seed rollback missing public_id filter: $(basename "$down")"
  fi
done

# Identify the latest migration timestamp for schema dump validation.
latest=""
for up in "$MIGRATIONS_SHARED"/*.up.sql; do
  [ -e "$up" ] || continue
  base=$(basename "$up")
  ts=${base%%_*}
  if [ -z "$latest" ] || [ "$ts" \> "$latest" ]; then
    latest="$ts"
  fi
done

# Ensure schema dumps exist for the latest migration.
if [ -n "$latest" ]; then
  if ! ls "$SCHEMA_SQLITE"/schema_*_after_${latest}*.sql >/dev/null 2>&1; then
    fail "missing sqlite schema dump for migration $latest"
  fi
  if ! ls "$SCHEMA_POSTGRES"/schema_*_after_${latest}*.sql >/dev/null 2>&1; then
    fail "missing postgres schema dump for migration $latest"
  fi
fi

# Success.
echo "migrate-check: ok"
