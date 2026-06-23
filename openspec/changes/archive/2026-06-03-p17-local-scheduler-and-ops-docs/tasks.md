## 1. Scheduler safety tests

- [x] Add tests or static checks for scheduler examples covering launchd and cron snippets.
- [x] Assert examples are disabled-by-default or require explicit local installation/editing.
- [x] Assert examples contain no automatic trading, broker order, one-click order, portfolio mutation, automatic confirmation, or automatic rule application commands.

## 2. `cmd/agent` help and task tests

- [x] Add or strengthen tests for `cmd/agent --help`.
- [x] Ensure help output lists supported local tasks and explains scheduler safety boundaries.
- [x] Ensure help output does not imply portfolio mutation, broker integration, or automatic rule application.

## 3. Implementation

- [x] Add macOS launchd example with placeholders only.
- [x] Add cron-compatible local scheduler example with placeholders only.
- [x] Refine `cmd/agent` help text only if tests require it.
- [x] Preserve existing confirmation and gatekeeper audit behavior.
- [x] Add narrow audit coverage only for any local task path found missing.

## 4. Operations documentation

- [x] Document local startup and initialization checklist.
- [x] Document SQLite backup/restore expectations.
- [x] Document VecLite index rebuild and degraded recovery.
- [x] Document data source, DeepSeek, VecLite, SQLite, and scheduled task failure handling.
- [x] Document scheduler setup, verification, disable, and removal procedure.
- [x] Document P17 boundary: no trading, no order placement, no automatic portfolio mutation, no automatic rule application.

## 5. Validation

- [x] Run `openspec validate p17-local-scheduler-and-ops-docs --strict`.
- [x] Run relevant Go tests for `cmd/agent` and any changed packages.
- [x] Run `go test ./...`.
- [x] Run `openspec validate --all --strict`.
