# AGENTS Guidelines (bms-core)

<!-- SPDX-License-Identifier: Apache-2.0 -->

## License Notice
Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

This document is licensed under Apache-2.0.

## Scope
- Applies to the `bms-core` repository: the Go server (`bmsd`), the TUI/CLI (`bms`), and everything under `internal/`.
- Workspace-level conventions (planning documents, documentation standards, multi-repo workflow, licensing boundaries) live in `bms-meta/AGENTS.md`, one directory up when this repo is checked out inside the workspace. This file is self-contained so a standalone clone still carries the rules it needs; the two files link to each other instead of duplicating content.
- This file is the single source of truth for agent guidance in this repo. `CLAUDE.md` only imports it.

## Repository Status (2026-09-12)
- Planning resumption review completed locally. Phase 1 (MVP server + TUI/CLI)
  is in progress; steps 1.2.1 through 1.2.4 are complete. The next product step
  remains 1.1.4, beginning with specification review in
  `bms-meta/docs/plans/storage-foundation.md`. Read `bms-meta/CURRENT-WORK.md`
  for the current lane and handoff, not an old session transcript.
- Audit baseline: build, vet and existing tests pass. `just lint` reports six
  unchecked-close findings in `internal/config/load.go` and
  `internal/db/db_test.go`; the first implementation package includes their
  bounded repair. `just fmt-check` is empty.
- Built and tested: `internal/config` (strict TOML, pointer-based overlays, defaults → file → env → CLI → server precedence, aggregated validation, redaction, warnings), `internal/logging` (slog, canonical fields, request/trace context, startup diagnostics), `internal/health` (`/healthz`, `/readyz`), `internal/db` (sqlite/libsql/postgres open, pooling, SQLite DSN pragmas), `internal/errtext` (centralized error text).
- Entry points: `cmd/bmsd` resolves config, logs diagnostics, serves health, opens the database, marks readiness, and shuts down on signal. `cmd/bms` resolves config and logs diagnostics only; there is no TUI yet.
- Not yet present: `db/` (migrations, schema dumps, queries, `sqlc.yaml`), `proto/`, `pkg/`, `plugins/`. `configs/`, `deploy/`, `docs/`, `gen/` hold `.gitkeep` placeholders.
- `scripts/migrate-check.sh` and every `migrate-*` recipe expect the `db/` tree and fail with `missing directory` until it exists. That is step `1.1.4 Migration and Storage Implementation`, the next step in the Phase 1 execution order (`bms-meta/EXECUTION-PLAN.md`).
- Decisions from the 2026-09-10 plan review that shape the next steps (`bms-meta/docs/reviews/2026-09-10-plan-review.md`): migrations are embedded in `bmsd` and applied through the golang-migrate library (the Homebrew `migrate` CLI has no SQLite driver); SQLite DSN pragmas `foreign_keys(1)`, `journal_mode(WAL)`, `synchronous(NORMAL)` join `busy_timeout`; timestamps are `TIMESTAMP`/`TIMESTAMPTZ` RFC 3339 UTC; `public_id` is ULID via `oklog/ulid/v2`; the TUI targets Bubble Tea v2 (`charm.land/bubbletea/v2`); MVP auth is local trust + optional password + sessions.
- Implementation boundary verified 2026-09-12: only `busy_timeout` is injected
  into SQLite DSNs. Foreign keys/WAL/NORMAL, embedded migrations, schema-aware
  readiness, auth handlers and the new key-auth warnings are still planned.
  The current config accepts remote libSQL URLs, not a local-file libSQL route.
- Update this section when a plan step completes or the tree changes.

## Branches, Worktrees, and Handoff
- Every change, including documentation, uses a dedicated branch and linked
  Git worktree; never edit the long-lived `main` checkout. Use
  `codex/<subject>` or `claude/<subject>`; match suffixes in paired repositories.
- Keep worktrees inside `/Volumes/HAMstack/.worktrees/<lane>/<subject>/` in the
  workspace. A standalone clone may use a sibling task directory.
- One session owns a checkout. Inspect branch, status, worktrees and remote
  ancestry before starting; preserve another session's uncommitted changes.
- Either Claude or Codex may plan, implement or review. Record scope, baseline,
  decisions, checks and next action; use a distinct worktree for implementation.
- Stage explicit paths, review the staged diff and commit only owned changes.
  Do not use workspace `commit-all` in parallel lanes. Push only when asked.
- Shared status and contract edits are serialized; isolated paths alone do not
  make logically dependent work independent. Use disposable validation targets
  and do not stop another session's runtime.
- Workspace detail: `bms-meta/docs/workflow.md`; live handoff:
  `bms-meta/CURRENT-WORK.md`. These sibling references work in paired worktrees.

## Build, Lint, Test
- Task runner: `just` (the Makefile was removed on 2026-09-10). A bare `just` lists the recipes.
- Build all: `just build` (`go build ./...`). Vet: `just vet`. Build, vet, and test together: `just check`.
- Run all tests: `just test`. One package: `just test ./internal/config`.
- One test, verbose: `just test-one ./internal/db TestSQLiteDSN`.
- Format in place: `just fmt` (`gofmt -w .`). List unformatted files: `just fmt-check`.
- Lint: `just lint` (`golangci-lint run ./...`, configured by `.golangci.yml`: the v2 standard set plus gofmt). The binary must be built with the same Go major version as the active toolchain; a mismatch panics on package load. Reinstall with `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`.
- Migrations: `just dep`, `just migrate-check`, `just migrate-up <engine>`, `just migrate-down <engine>`, `just migrate-version <engine>`, `just migrate-dump <engine>`, `just migrate-validate <engine>`, where `<engine>` is `sqlite` or `postgres`. `up` and `down` depend on `dep` and `migrate-check`, so lint always runs first.
- Local settings: copy `.env.example` to `.env` (gitignored; a copy, not a symlink, or edits land in the example); the justfile loads it automatically, and shell environment wins over `.env`, which wins over the recipe defaults. Homebrew keeps `libpq` keg-only, so set `BMS_PG_DUMP=/opt/homebrew/opt/libpq/bin/pg_dump` (the version-independent `opt` symlink, not a `Cellar` path).
- Tooling expected on the machine: `go`, `just`, `migrate` (golang-migrate), `sqlite3`, `pg_dump`, `sqlc`, `golangci-lint`. Still to install when the step arrives: `buf` (proto, Phase 1.4).
- Planned API tooling (Phase 1.5): Swagger UI for the REST gateway, `protoc-gen-doc` for gRPC, Postman collections.
- CI is deferred; recipes are run manually.

## Code Style: General
- Prefer small, cohesive packages with clear ownership.
- Keep functions short; favor explicit over clever.
- Avoid TODOs without a tracking issue or owner.
- Keep public APIs documented with Go doc comments.
- Keep generated code in `gen/` or `pkg/generated/`.
- Indent with tabs (`ts=4`, `sw=4`, `noet`) in every language.
- Go files are `gofmt` canonical: `just fmt` rewrites them and `just fmt-check` must report nothing. `gofmt` aligns trailing comments with spaces, so the tab-alignment rule below applies to non-Go files only.
- In non-Go files (justfile, shell, SQL, TOML, `.env`), align inline comments with tabs.
- Add a Vim modeline to every source file using the appropriate comment syntax (Go: `// vim: set ts=4 sw=4 noet:`; shell, justfile, SQL, TOML: `# vim: set ts=4 sw=4 noet:` or `-- vim: ...`).
- Add a file-level comment to every source file that explains what it contains, what it does, and how it works.
- Leave a blank line between the file-level comment block and the `package` declaration.
- Leave a blank line before the Vim modeline.
- Avoid direct output/logging in internal/library packages; entrypoints (`cmd/*`) own logging and printing.
- Use Vim-style folding markers: `// Block Name {{{` ... `// }}}`. For a logical block, add a header comment followed by `// }}}` to close it. In justfiles the markers are safe inside `#` comments at the top level; never put `{{{` inside a recipe body, where `{{` is interpolation.
- Prefer the Go standard library; external dependencies require clear justification recorded in the commit and the relevant doc. Current exceptions: `BurntSushi/toml` (strict decoding), `modernc.org/sqlite` (cgo-free SQLite), `tursodatabase/go-libsql` (libSQL, cgo, non-Windows), `jackc/pgx/v5` (PostgreSQL).
- For `.gitignore`, use the SPDX header + blank line, block headers with `{{{`/`}}}`, and a blank line after each block header and before each block end.

## Code Style: Go
- Import order: standard library, blank line, third-party, blank line, local module.
- Avoid one-letter names except short-lived loop indexes.
- Prefer `context.Context` as the first argument in public functions.
- Errors: wrap with `%w` and add context; avoid string comparisons. User-facing and log error text lives in `internal/errtext` as constants; do not inline new error strings in packages.
- Return `error` last; avoid panics except in `main`.
- Prefer `time.Time` over strings for timestamps in domain types. Config durations are strings that must parse with `time.ParseDuration` and are validated in `internal/config/validate.go`.
- Use `sqlc` types once `db/` exists; avoid raw SQL strings in services.
- Keep interfaces small; depend on abstractions in `internal/`.
- For config, use explicit structs; avoid `map[string]interface{}`. The runtime `Config` is non-pointer; the file/CLI/server inputs are pointer-based `*Overlay` structs so "unset" and "explicit zero" stay distinguishable. Add a field to both, to `merge.go`, to `redaction.go` if it is a secret, and to `docs/configuration.md` in `bms-meta`.
- Connection-scoped SQLite pragmas go into the DSN (`_pragma=name(value)`), never through `ExecContext`, because `database/sql` pools connections and an executed pragma reaches only one of them. See `internal/db/db.go`.
- Build tags: platform-conditional code gets a `_enabled.go`/`_disabled.go` pair with matching function signatures (see `internal/db/libsql_*.go`).

## Code Style: gRPC/REST
- gRPC: use proto3, explicit field numbers, and clear naming.
- Keep request/response messages small and composable.
- Version APIs by package name or URL prefix.
- REST gateway should mirror gRPC resource names.
- WebSocket events should be versioned and documented.

## Code Style: SQL/Migrations
- One migration per change; keep reversible where possible.
- Deterministic ordering: `YYYYMMDDHHMMSS_<slug>.up.sql` / `.down.sql` (see `docs/data-layer.md` in `bms-meta`).
- Three parallel sets: `db/migrations/shared` is the source; `sqlite/` and `postgres/` must contain the same file names. `scripts/migrate-check.sh` enforces pairs, parity, the `-- IRREVERSIBLE:` header on `DROP`/`TRUNCATE`, `public_id`-scoped seed rollbacks, and a schema dump for the latest migration.
- Prefer indexes for common logbook queries (call, band, time).
- Avoid vendor-specific SQL in shared schemas when possible.
- Require SQLite/PostgreSQL schema parity; use lowest-common-denominator SQL.
- Prefer transactional migrations when supported; always provide down migrations.

## Code Style: TUI/CLI (BubbleTea)
- Keep model updates pure; side effects via commands.
- Separate rendering from state updates.
- Keybindings documented and centralized.
- Themes stored as structured config, not ad-hoc constants.

## Code Style: Lua Plugins
- Keep plugin API stable; document hooks and data contracts.
- Sandbox plugins; avoid file/network access by default.
- Provide versioned APIs and compatibility notes.

## Error Handling and Logging
- Use structured logging (`log/slog`) with the canonical fields in `internal/logging/fields.go`.
- Log at boundaries (RPC, storage, integrations), from entrypoints or boundary packages only.
- Avoid logging secrets or auth tokens; extend `internal/config/redaction.go` when adding a secret field.
- Convert internal errors to user-friendly messages at UI layers.

## Config and Secrets
- Support server and local client configs with clear precedence (documented in `bms-meta/docs/configuration.md`).
- Store secrets in env vars or the OS keychain when available.
- Use XDG base directories on Unix-like systems (`~/.config`, `~/.local/share`, `~/.cache`).
- Provide sample config files with safe defaults.

## Licensing Rules
- Every new file gets an SPDX header (`SPDX-License-Identifier: Apache-2.0`) and the copyright line; templates are in `bms-meta/SPDX_Apache.txt`.
- This repository is Apache-2.0 only. SaaS-critical components (multi-tenant management, cloud sync, realtime feeds, contest coordination, AI/statistics) belong in `bms-cloud`, never here.
- Record license owner: Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO).

## Repo Structure
- Present: `cmd/bmsd`, `cmd/bms`, `internal/{config,db,errtext,health,logging}`, `scripts/`, `configs/`, `deploy/`, `docs/`, `gen/`.
- Planned: `internal/{auth,logbook,contest,integrations,sync}`, `pkg/` (shared SDK/client packages), `proto/` (gRPC contracts), `db/{migrations,schema,queries}` + `sqlc.yaml`, `plugins/` (Lua runtime assets).

## Testing Guidance
- Unit tests next to the implementation (`*_test.go`), in the same package.
- Table-driven tests for parsing and validation; name cases as sentences (`"negative busy timeout fails"`).
- Favor deterministic tests; avoid `time.Now` without injection.
- Database tests use `file:<name>?mode=memory&cache=shared` so several pooled connections share one in-memory database; a bare `:memory:` gives each connection its own.
- Add integration tests for DB migrations and gRPC flows when those layers land.
- Every new test group gets a bullet in the "Testing Notes" section of the relevant `bms-meta/docs/*.md`.

## Commit Guidance
- Keep commits scoped and descriptive: `<area>: <imperative summary>` (`database: complete connection pooling`, `config: add --config override flag`).
- The same message is used for the paired `bms-meta` commit when docs change alongside code.
- Note license classification in the message when adding a new component.
- Never push unless asked.

## Security and Privacy
- Enforce TLS for remote connections when enabled.
- Authentication optional, but secure defaults.
- Rate-limit external integrations.
- Validate all user input from UI and APIs.
- Redact personal data in logs.
- Use least-privilege DB users when supported.

## Performance
- Use pagination for logbook lists.
- Add indexes for call/date/band queries.
- Cache lookups for QRZ/LoTW.

## UX Notes
- Keep keymaps consistent across TUI/web/desktop.
- Provide clear offline/online status indicators.

## When Missing Commands
- If a needed recipe does not exist, add it to `justfile` and document it in this file in the same change.

## Contact
- Primary owner: Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO).
