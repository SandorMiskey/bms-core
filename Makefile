# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

.PHONY: build dep fmt help lint migrate-check migrate-down-postgres migrate-down-sqlite migrate-dump-postgres migrate-dump-sqlite migrate-up-postgres migrate-up-sqlite migrate-validate-postgres migrate-validate-sqlite test vet

BMS_MIGRATE ?= migrate
BMS_MIGRATIONS_SHARED ?= db/migrations/shared
BMS_MIGRATIONS_SQLITE ?= db/migrations/sqlite
BMS_MIGRATIONS_POSTGRES ?= db/migrations/postgres
BMS_SCHEMA_SQLITE ?= db/schema/sqlite
BMS_SCHEMA_POSTGRES ?= db/schema/postgres
BMS_DSN_POSTGRES ?= postgres://localhost:5432/bms?sslmode=disable
BMS_SQLITE_PATH ?= bms.db
BMS_SQLITE_DSN ?= sqlite3://$(BMS_SQLITE_PATH)

export BMS_MIGRATIONS_SHARED BMS_MIGRATIONS_SQLITE BMS_MIGRATIONS_POSTGRES BMS_SCHEMA_SQLITE BMS_SCHEMA_POSTGRES

build: ## build all packages
	go build ./...

fmt: ## format Go sources
	gofmt -w .

help: ## show all available targets
	@echo "available targets:"
	@awk -F ':.*## ' '/^[a-zA-Z0-9_-]+:.*## /{ printf "    %-24s%s\n", $$1, $$2 }' $(MAKEFILE_LIST) | sort

lint: ## run linter
	golangci-lint run ./...

dep: ## check migration prerequisites
	@echo "checking migration prerequisites..."
	@command -v $(BMS_MIGRATE) >/dev/null || { echo "missing migrate CLI"; exit 1; }
	@version=$$($(BMS_MIGRATE) -version 2>&1 || true); \
	if [ -z "$$version" ]; then version="unknown"; fi; \
	echo "migrate: $$version"
	@command -v sqlite3 >/dev/null || { echo "missing sqlite3"; exit 1; }
	@echo "sqlite3: $$(sqlite3 --version)"
	@command -v pg_dump >/dev/null || { echo "missing pg_dump"; exit 1; }
	@echo "pg_dump: $$(pg_dump --version)"

migrate-check: dep ## run migration lint checks
	@sh scripts/migrate-check.sh

migrate-down-postgres: dep migrate-check ## rollback one postgres migration
	$(BMS_MIGRATE) -database $(BMS_DSN_POSTGRES) -path $(BMS_MIGRATIONS_POSTGRES) down 1

migrate-up-postgres: dep migrate-check ## apply postgres migrations
	$(BMS_MIGRATE) -database $(BMS_DSN_POSTGRES) -path $(BMS_MIGRATIONS_POSTGRES) up

migrate-down-sqlite: dep migrate-check ## rollback one sqlite migration
	$(BMS_MIGRATE) -database $(BMS_SQLITE_DSN) -path $(BMS_MIGRATIONS_SQLITE) down 1

migrate-up-sqlite: dep migrate-check ## apply sqlite migrations
	$(BMS_MIGRATE) -database $(BMS_SQLITE_DSN) -path $(BMS_MIGRATIONS_SQLITE) up

migrate-dump-postgres: dep ## generate postgres schema dump
	@mkdir -p $(BMS_SCHEMA_POSTGRES)
	@version=$$($(BMS_MIGRATE) -database $(BMS_DSN_POSTGRES) -path $(BMS_MIGRATIONS_POSTGRES) version 2>/dev/null | cut -d ' ' -f 1 || true); \
	if [ -z "$$version" ]; then version=none; fi; \
	ts=$$(date +%Y%m%d%H%M%S); \
	file="$(BMS_SCHEMA_POSTGRES)/schema_$${ts}_after_$${version}.sql"; \
	pg_dump --schema-only --no-owner --no-privileges "$(BMS_DSN_POSTGRES)" > "$${file}"


migrate-dump-sqlite: dep ## generate sqlite schema dump
	@mkdir -p $(BMS_SCHEMA_SQLITE)
	@version=$$($(BMS_MIGRATE) -database $(BMS_SQLITE_DSN) -path $(BMS_MIGRATIONS_SQLITE) version 2>/dev/null | cut -d ' ' -f 1 || true); \
	if [ -z "$$version" ]; then version=none; fi; \
	ts=$$(date +%Y%m%d%H%M%S); \
	file="$(BMS_SCHEMA_SQLITE)/schema_$${ts}_after_$${version}.sql"; \
	sqlite3 "$(BMS_SQLITE_PATH)" ".schema" > "$${file}"

migrate-validate-postgres: dep migrate-check migrate-up-postgres migrate-dump-postgres ## validate postgres migrations
	@echo "postgres migration validation complete"

migrate-validate-sqlite: dep migrate-check migrate-up-sqlite migrate-dump-sqlite ## validate sqlite migrations
	@echo "sqlite migration validation complete"

test: ## run all tests
	go test ./...

vet: ## run go vet
	go vet ./...
