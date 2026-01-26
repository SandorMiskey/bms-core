# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

.PHONY: build dep fmt help lint migrate-check migrate-down-postgres migrate-down-sqlite migrate-dump-postgres migrate-dump-sqlite migrate-up-postgres migrate-up-sqlite test vet

MIGRATE ?= migrate
MIGRATIONS_SHARED ?= db/migrations/shared
MIGRATIONS_POSTGRES ?= db/migrations/postgres
MIGRATIONS_SQLITE ?= db/migrations/sqlite
SCHEMA_POSTGRES ?= db/schema/postgres
SCHEMA_SQLITE ?= db/schema/sqlite
POSTGRES_DSN ?= postgres://localhost:5432/bms?sslmode=disable
SQLITE_PATH ?= bms.db
SQLITE_DSN ?= sqlite3://$(SQLITE_PATH)

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
	@command -v $(MIGRATE) >/dev/null || { echo "missing migrate CLI"; exit 1; }
	@command -v sqlite3 >/dev/null || { echo "missing sqlite3"; exit 1; }
	@command -v pg_dump >/dev/null || { echo "missing pg_dump"; exit 1; }

migrate-check: dep ## run migration lint checks
	@bash scripts/migrate-check.sh

migrate-down-postgres: dep migrate-check ## rollback one postgres migration
	$(MIGRATE) -database $(POSTGRES_DSN) -path $(MIGRATIONS_POSTGRES) down 1

migrate-up-postgres: dep migrate-check ## apply postgres migrations
	$(MIGRATE) -database $(POSTGRES_DSN) -path $(MIGRATIONS_POSTGRES) up

migrate-down-sqlite: dep migrate-check ## rollback one sqlite migration
	$(MIGRATE) -database $(SQLITE_DSN) -path $(MIGRATIONS_SQLITE) down 1

migrate-up-sqlite: dep migrate-check ## apply sqlite migrations
	$(MIGRATE) -database $(SQLITE_DSN) -path $(MIGRATIONS_SQLITE) up

migrate-dump-postgres: dep ## generate postgres schema dump
	@mkdir -p $(SCHEMA_POSTGRES)
	@version=$$($(MIGRATE) -database $(POSTGRES_DSN) -path $(MIGRATIONS_POSTGRES) version 2>/dev/null | cut -d ' ' -f 1 || true); \
	if [ -z "$$version" ]; then version=none; fi; \
	ts=$$(date +%Y%m%d%H%M%S); \
	file="$(SCHEMA_POSTGRES)/schema_$${ts}_after_$${version}.sql"; \
	pg_dump --schema-only --no-owner --no-privileges "$(POSTGRES_DSN)" > "$${file}"

migrate-dump-sqlite: dep ## generate sqlite schema dump
	@mkdir -p $(SCHEMA_SQLITE)
	@version=$$($(MIGRATE) -database $(SQLITE_DSN) -path $(MIGRATIONS_SQLITE) version 2>/dev/null | cut -d ' ' -f 1 || true); \
	if [ -z "$$version" ]; then version=none; fi; \
	ts=$$(date +%Y%m%d%H%M%S); \
	file="$(SCHEMA_SQLITE)/schema_$${ts}_after_$${version}.sql"; \
	sqlite3 "$(SQLITE_PATH)" ".schema" > "$${file}"

test: ## run all tests
	go test ./...

vet: ## run go vet
	go vet ./...
