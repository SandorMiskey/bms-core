# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

.PHONY: build fmt help lint migrate-postgres-down migrate-postgres-up migrate-sqlite-down migrate-sqlite-up test vet

MIGRATE ?= migrate
MIGRATIONS_POSTGRES ?= db/migrations/postgres
MIGRATIONS_SQLITE ?= db/migrations/sqlite
POSTGRES_DSN ?= postgres://localhost:5432/bms?sslmode=disable
SQLITE_DSN ?= sqlite3://bms.db

build: ## build all packages
	go build ./...

fmt: ## format Go sources
	gofmt -w .

help: ## show all available targets
	@echo "available targets:"
	@awk -F ':.*## ' '/^[a-zA-Z0-9_-]+:.*## /{ printf "    %-24s%s\n", $$1, $$2 }' $(MAKEFILE_LIST) | sort

lint: ## run linter
	golangci-lint run ./...

migrate-postgres-down: ## rollback one postgres migration
	$(MIGRATE) -database $(POSTGRES_DSN) -path $(MIGRATIONS_POSTGRES) down 1

migrate-postgres-up: ## apply postgres migrations
	$(MIGRATE) -database $(POSTGRES_DSN) -path $(MIGRATIONS_POSTGRES) up

migrate-sqlite-down: ## rollback one sqlite migration
	$(MIGRATE) -database $(SQLITE_DSN) -path $(MIGRATIONS_SQLITE) down 1

migrate-sqlite-up: ## apply sqlite migrations
	$(MIGRATE) -database $(SQLITE_DSN) -path $(MIGRATIONS_SQLITE) up

test: ## run all tests
	go test ./...

vet: ## run go vet
	go vet ./...
