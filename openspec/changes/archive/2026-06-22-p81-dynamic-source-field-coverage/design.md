# P81 Design

## Evidence Strategy

P81 treats dynamic source coverage as a user-symbol-driven product behavior, not a fixture inventory. The acceptance runner must start from at least one fresh non-`510300` symbol selected through the product path and must prove how the symbol drives collector/readiness/API/UI/LLM context outputs.

The evidence layer should include:

- P80 matrix row references and previous status for all 59 rows.
- Fresh command output summaries for source quality/readiness checks.
- API and UI evidence for user-selected symbol field coverage.
- Read-only SQLite/readback evidence for persisted source facts, health, audit events, and indexed references where applicable.
- Degraded/missing evidence that blocks or qualifies claims when formal data is unavailable.
- A forbidden-capability scan proving the acceptance did not introduce trading, broker, external push, or automation affordances.

## Real-Pass Rule

A row may be upgraded to `real_pass` only when the P81 evidence proves the current product obtains or evaluates the relevant field for the selected symbol from formal accepted sources or safely blocks the affected claim. Background knowledge, fixture-only data, screenshots without readback, or hard-coded symbol assumptions are insufficient.

## Execution Boundary

P81 may add scripts, tests, evidence records, or minor hardening needed to make the acceptance repeatable. It must not broaden product capability beyond existing source/readiness behavior unless a missing requirement requires a scoped implementation fix recorded in tasks.

