# AGENTS Guidelines (bms-core)

> Lifecycle: canonical living document.
> Last updated: 2026-09-13 (lifecycle/review metadata).
> Evidence: Verified now for repository topology; current owner instructions govern workflow. Previous product choices await reapproval.

<!-- SPDX-License-Identifier: Apache-2.0 -->

## License Notice
Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

This document is licensed under Apache-2.0.

## Scope
- Applies to the `bms-core` repository: the Go server (`bmsd`), the TUI/CLI (`bms`), and everything under `internal/`.
- Workspace-level conventions (planning documents, documentation standards, multi-repo workflow, licensing boundaries) live in `bms-meta/AGENTS.md`, one directory up when this repo is checked out inside the workspace. This file is self-contained so a standalone clone still carries the rules it needs; the two files link to each other instead of duplicating content.
- This file is tool-neutral repository guidance. `CLAUDE.md` imports it and
  adds only tool mechanics, never a policy exception. AGENTS ownership/thinning
  is under review in R0/P; no file consolidation has been assumed approved.

## Repository Status (2026-09-13)
- Verified now: the existing tree still contains config, logging, health and DB
  connection foundations; there is no migration/schema/storage, auth or TUI
  implementation. See Repo Structure below and the dated meta restart review
  for the prior baseline checks and six known lint findings.
- Current owner instruction: product implementation is paused for organization
  and full reapproval of previous decisions. Do not start 1.1.4 from an old handoff.
- Read `../bms-meta/CURRENT-WORK.md` and the revalidation programme in
  `../bms-meta/EXECUTION-PLAN.md` from a paired worktree. The shared policy is
  `../bms-meta/AGENTS.md`; sibling guidance is not automatic inheritance.
- A standalone clone retains its local code guidance. R0/P is deciding how to
  keep a thin entry point and bootstrap the shared policy; do not remove this
  file or create an external symlink before that decision.
- Existing technical rules below remain constraints while new product work is
  paused; their design choices are inputs to the relevant R*/P review units.

## Branches, Worktrees, and Handoff
- Every change, including documentation, uses a dedicated branch and linked
  Git worktree; never edit the long-lived `main` checkout. Use
  `codex/<subject>` or `claude/<subject>`; match suffixes in paired repositories.
- Keep worktrees inside `/Volumes/HAMstack/.worktrees/<lane>/<subject>/` in the
  workspace. A standalone clone may use a sibling task directory.
- One session owns a checkout. Inspect branch, status, worktrees and remote
  ancestry before starting; preserve another session's uncommitted changes.
- Either Claude or Codex may plan, implement or review. Record scope, baseline,
  decisions, checks and next action. /P and /I use distinct branches/worktrees;
  the /P plan must pass the shared Implementer Acceptance Test before closure.
- Stage explicit paths, review the staged diff and commit only owned changes.
  Do not use workspace `commit-all` in parallel lanes. Merge, push and deletion
  each require explicit owner instruction; completion does not authorize them.
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
- Before migration SQL, obtain approval for the DB decision gate and the exact
  migration scope. R6/P revalidates the baseline; historical design is not approval.
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
- Merge, push and deletion require explicit owner instruction. Never infer them
  from local commit, passing checks, or planning approval.

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
