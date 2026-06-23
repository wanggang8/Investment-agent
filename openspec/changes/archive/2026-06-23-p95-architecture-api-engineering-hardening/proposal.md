# P95 Architecture API Engineering Hardening

## Why

P94 made the public GitHub repository buildable and protected by CI/CD. A follow-up architecture and engineering review found several maintainability risks that should be closed before the project accumulates more public contributions:

- Go package discovery currently includes `web/node_modules` after frontend dependencies are installed.
- P93 code reality scanning can be affected by ignored or untracked local runtime artifacts such as `cmd/agent/tmp/veclite`.
- HTTP API routes are registered by hand and documented by hand, leaving contract drift possible.
- Architecture documentation still contains examples from older directory shapes.
- SQLite runtime settings should be more explicit for local UI plus background task concurrency.
- Docker deployment should support file-based LLM secrets in addition to `.env`.

## What Changes

- Add a repository-owned Go package selection helper and use it from CI/release docs/scripts so `web/node_modules` cannot become part of backend validation.
- Harden P93 scanning so it scans tracked/release-relevant source files only, not ignored local runtime artifacts.
- Add an API route contract check that compares registered backend routes with documented routes and makes the current route surface auditable.
- Update architecture documentation to match the current code layout and known engineering constraints.
- Add SQLite connection pragmas for foreign keys, busy timeout, and WAL-oriented local concurrency where appropriate.
- Add Docker/Compose support for file-based DeepSeek API key secrets while keeping `.env` as the simple local path.

## Out Of Scope

- New investment strategy, new product UI workflows, or new business API capabilities.
- Broker integration, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.
- Rewriting the full backend architecture or migrating away from SQLite, Eino, Vite, or React.
- Publishing a Git tag or GitHub Release.

## Acceptance

P95 is acceptable only if:

- Backend validation no longer discovers Go packages under `web/node_modules`.
- P93 `--check` is stable after frontend dependency installation and local ignored artifacts.
- The API route contract check passes and is wired into CI.
- Architecture docs describe the current real directory layout and route-contract policy.
- SQLite settings are covered by focused tests.
- Docker file-secret support is documented and does not expose secrets in committed files.
- Existing backend, frontend, OpenSpec, P91/P92/P93, and package smoke gates remain passing.
