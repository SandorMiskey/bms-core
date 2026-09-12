@AGENTS.md

<!-- SPDX-License-Identifier: Apache-2.0 -->
# Claude Code Mechanics

> Lifecycle: canonical living document.
> Last updated: 2026-09-13.
> Evidence: proposed delta mechanics under R0/P; the imported owner policy remains authoritative.

## License Notice
Copyright (c) 2026 Sandor Miskey (HA5BMS, sandor@HA5BMS.RADIO)

This document is licensed under Apache-2.0.

## Delta Boundary
A delta may add tool mechanics only; it never relaxes, redefines or supersedes
AGENTS policy. Record an unsupported control as a gap. Do not interpret lack of
a technical prompt as approval.

Read the common policy/current-unit paths named by this repository's AGENTS when
working in a paired workspace. A sibling reference is not an automatic import.

## Worktree Invocation
Use the owner-recorded Git worktree and start Claude normally from its directory.
Do not use `claude --worktree`, `claude -w`, EnterWorktree, or Agent
`isolation: worktree`: automatic tool-owned cleanup conflicts with the common
explicit merge/delete timing. Do not copy ignored secrets/settings from another
worktree as a shortcut.

## State and Permissions
Read project state from the repository, not Claude auto-memory notes. The old
HAMstack notes are inventoried in the R0/P proposal; no deletion or global memory
change is implied. The permission JSON there is a candidate, not installed
settings. Check active settings at the approved setup step; never claim that this
Markdown file mechanically enforces the permission floor.

At a planning-to-implementation handoff, report the selected model/effort and
reconsider its cost. Do not automatically switch tools or enforce an old memory's
model name. Follow the shared handoff format and current owner assignment.
