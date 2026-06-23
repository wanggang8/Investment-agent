# P95 Tasks

## 1. Scope And Boundary Review

- [x] Confirm no active OpenSpec change conflicts with P95.
- [x] Confirm P95 does not edit root `README.md`, `docs/README.md` history structure, `docs/product-overview.md`, or `docs/quickstart.md` except for engineering command consistency created by P96.

## 2. Go Package And CI Hardening

- [x] Add a backend package helper that excludes `web/node_modules` and other frontend dependency packages.
- [x] Add tests or command evidence proving backend package selection excludes `web/node_modules`.
- [x] Update CI/release workflows and developer docs to use the helper for Go tests.

## 3. P93 Stable Source Inventory

- [x] Make P93 scan tracked/release-relevant files instead of ignored local runtime artifacts.
- [x] Add evidence that `cmd/agent/tmp/veclite` or equivalent ignored files do not affect P93 classification counts.
- [x] Keep P93 secret and forbidden-path checks active.

## 4. API Route Contract Check

- [x] Add a script that compares registered API routes with normalized documented routes.
- [x] Document the currently missing routes or adjust docs so the check passes.
- [x] Wire the check into CI.

## 5. SQLite And Docker Runtime Hardening

- [x] Add SQLite connection pragmas for foreign keys, busy timeout, and file-backed WAL where compatible.
- [x] Add focused SQLite tests for the new connection behavior.
- [x] Add `DEEPSEEK_API_KEY_FILE` config support and Compose secrets documentation.

## 6. Architecture Documentation And Acceptance

- [x] Update `docs/architecture.md` to reflect current directories and engineering constraints.
- [x] Generate P95 acceptance evidence.
- [x] Run OpenSpec, backend, frontend, P91/P92/P93, API route, package smoke, and security-equivalent checks.
- [ ] Archive P95 only after validation passes.
