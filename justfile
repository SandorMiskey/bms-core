# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)
#
# bms-core task runner.
# This justfile replaces the previous Makefile. It exposes the Go build, test,
# and lint entry points plus the golang-migrate workflow for the sqlite and
# postgres migration sets.
#
# Engine-specific recipes take the engine as an argument (`just migrate-up
# sqlite`) so both engines share one definition and cannot drift apart.
#
# Local overrides belong in `.env`, which is loaded automatically; see
# `.env.example` for the supported variables. Precedence is shell environment,
# then `.env`, then the defaults below.

set dotenv-load := true
set export := true

# Tooling {{{

BMS_MIGRATE := env('BMS_MIGRATE', 'migrate')
BMS_SQLITE3 := env('BMS_SQLITE3', 'sqlite3')

# Homebrew keeps libpq keg-only, so `pg_dump` is usually not on PATH. Point
# BMS_PG_DUMP at the `opt` symlink rather than a `Cellar` path, which pins a
# version and breaks on the next upgrade.
BMS_PG_DUMP := env('BMS_PG_DUMP', 'pg_dump')

# }}}
# Migration paths {{{

BMS_MIGRATIONS_SHARED := env('BMS_MIGRATIONS_SHARED', 'db/migrations/shared')
BMS_MIGRATIONS_SQLITE := env('BMS_MIGRATIONS_SQLITE', 'db/migrations/sqlite')
BMS_MIGRATIONS_POSTGRES := env('BMS_MIGRATIONS_POSTGRES', 'db/migrations/postgres')
BMS_SCHEMA_SQLITE := env('BMS_SCHEMA_SQLITE', 'db/schema/sqlite')
BMS_SCHEMA_POSTGRES := env('BMS_SCHEMA_POSTGRES', 'db/schema/postgres')

# }}}
# Database targets {{{

BMS_SQLITE_PATH := env('BMS_SQLITE_PATH', 'bms.db')
BMS_SQLITE_DSN := env('BMS_SQLITE_DSN', 'sqlite3://' + BMS_SQLITE_PATH)
BMS_DSN_POSTGRES := env('BMS_DSN_POSTGRES', 'postgres://localhost:5432/bms?sslmode=disable')

# }}}
# Go recipes {{{

# show all available recipes
default:
	@just --list --unsorted

# build all packages
build:
	go build ./...

# run all tests, or one package when a path is given
test package='./...':
	go test {{ package }}

# run a single test by name, verbosely
test-one package name:
	go test {{ package }} -run {{ name }} -v

# run go vet
vet:
	go vet ./...

# format Go sources in place
fmt:
	gofmt -w .

# list Go sources that are not gofmt clean
fmt-check:
	@gofmt -l .

# run the linter
lint:
	golangci-lint run ./...

# build, vet, and test in one pass
check: build vet test

# }}}
# Migration recipes {{{

# check that the migration prerequisites are installed
dep:
	#!/bin/sh
	set -eu
	missing=0
	report() {
		if command -v "$1" >/dev/null 2>&1; then
			printf '%-9s %s\n' "$2:" "$(command -v "$1")"
		else
			printf '%-9s missing -- %s\n' "$2:" "$3" >&2
			missing=1
		fi
	}
	report "$BMS_MIGRATE" migrate 'install golang-migrate'
	report "$BMS_SQLITE3" sqlite3 'install sqlite3'
	report "$BMS_PG_DUMP" pg_dump 'Homebrew keeps libpq keg-only; set BMS_PG_DUMP=/opt/homebrew/opt/libpq/bin/pg_dump in .env'
	[ "$missing" -eq 0 ] || exit 1

# run the migration lint checks
migrate-check: dep
	@sh scripts/migrate-check.sh

# apply every pending migration for an engine (sqlite or postgres)
migrate-up engine: dep migrate-check
	@just _migrate {{ engine }} up

# roll back one migration for an engine (sqlite or postgres)
migrate-down engine: dep migrate-check
	@just _migrate {{ engine }} down 1

# report the applied migration version for an engine (sqlite or postgres)
migrate-version engine: dep
	@just _migrate {{ engine }} version

# write a schema dump for an engine (sqlite or postgres)
migrate-dump engine: dep
	#!/bin/sh
	set -eu
	case "{{ engine }}" in
		sqlite)   target="$BMS_SCHEMA_SQLITE" ;;
		postgres) target="$BMS_SCHEMA_POSTGRES" ;;
		*) echo "unknown engine: {{ engine }} (expected sqlite or postgres)" >&2; exit 2 ;;
	esac
	# A dump taken before any migration is still useful, so an unknown version
	# is recorded as `none` rather than treated as an error.
	version=$(just _migrate {{ engine }} version 2>/dev/null | cut -d ' ' -f 1 || true)
	[ -n "$version" ] || version=none
	mkdir -p "$target"
	file="$target/schema_$(date +%Y%m%d%H%M%S)_after_${version}.sql"
	case "{{ engine }}" in
		sqlite)   "$BMS_SQLITE3" "$BMS_SQLITE_PATH" '.schema' > "$file" ;;
		postgres) "$BMS_PG_DUMP" --schema-only --no-owner --no-privileges "$BMS_DSN_POSTGRES" > "$file" ;;
	esac
	echo "wrote $file"

# apply migrations and regenerate the schema dump for an engine
migrate-validate engine: (migrate-up engine) (migrate-dump engine)
	@echo "{{ engine }} migration validation complete"

# resolve engine-specific migrate arguments in one place
[private]
_migrate engine *args:
	#!/bin/sh
	set -eu
	case "{{ engine }}" in
		sqlite)
			database="$BMS_SQLITE_DSN"
			path="$BMS_MIGRATIONS_SQLITE"
			;;
		postgres)
			database="$BMS_DSN_POSTGRES"
			path="$BMS_MIGRATIONS_POSTGRES"
			;;
		*)
			echo "unknown engine: {{ engine }} (expected sqlite or postgres)" >&2
			exit 2
			;;
	esac
	exec "$BMS_MIGRATE" -database "$database" -path "$path" {{ args }}

# }}}

# vim: set ts=4 sw=4 noet:
