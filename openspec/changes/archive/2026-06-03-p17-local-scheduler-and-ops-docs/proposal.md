# P17 Local Scheduler and Ops Docs

## Summary
Add a safe local scheduling and operations documentation layer for the existing `cmd/agent` tasks. The change provides default-disabled launchd/cron examples, clearer task help, audit expectations, and operational runbooks without enabling automatic trading or bypassing manual confirmation.

## Why
P9 introduced manual local agent tasks and P16 exposed review/ops status in the frontend. Operators still need a documented way to run those tasks periodically on a local machine while preserving the core safety boundary: scheduled tasks may refresh data, rebuild indexes, or generate review facts, but must never place orders or apply portfolio/rule mutations without existing confirmation gates.

## Scope
- Provide local scheduler examples for macOS launchd and cron/system-style usage with disabled-by-default semantics.
- Strengthen `cmd/agent` help and supported task descriptions where needed.
- Document audit behavior for scheduled/manual local tasks.
- Document startup, backup, index rebuild, failure handling, and safety boundaries.
- Add tests that assert scheduler examples and help text do not expose trading/order behavior.

## Non-goals
- No cloud scheduler or daemon service.
- No automatic trading, broker integration, or order placement.
- No bypass of user confirmation for portfolio changes or rule application.
- No broad rewrite of `cmd/agent` task execution.
