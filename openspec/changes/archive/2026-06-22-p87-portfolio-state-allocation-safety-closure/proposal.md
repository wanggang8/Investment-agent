# P87 Portfolio State Allocation Safety Closure

## Why

After P84, P85 covers 31 expected-return and analysis-accuracy rows, and P86 covers 94 final core-goal/knowledge/safety rows. A direct inventory of the P84 matrix shows 32 full-release-required non-`real_pass` rows still unowned. P87 exists to close that gap before P86 performs final consolidation.

These rows are mostly P84-deferred portfolio/allocation/state/data-impact requirements: account state, core/satellite/cash allocation, quarterly rebalance triggers, sell-only/frozen-watch transitions, public-source and SQLite readiness, data-insufficient safety, user confirmation, audit/readback, and release safety boundaries.

## What Changes

- Define P87 as the ownership plan for the 32 P84-after rows not covered by P85 or P86.
- Execute row-specific real local UI/API/SQLite acceptance for portfolio state, allocation, rebalance, data-insufficient behavior, manual confirmation, audit, and release-safety boundaries.
- Upgrade only rows whose complete requirement text is proven by fresh evidence; defer any broad row whose evidence remains partial.
- Preserve all existing safety boundaries and forbidden-capability exclusions.

## In Scope

- 32 rows: REQ-01-001, REQ-01-006, REQ-02-006, REQ-02-022, REQ-02-024, REQ-02-025, REQ-02-031, REQ-03-004, REQ-03-005, REQ-03-006, REQ-04-003, REQ-04-008, REQ-04-016, REQ-04-025, REQ-05-010, REQ-06-023, REQ-06-024, REQ-07-006, REQ-07-015, REQ-08-018, REQ-08-020, REQ-10-001, REQ-10-002, REQ-10-003, REQ-10-004, REQ-11-005, REQ-14-005, REQ-14-007, REQ-16-028, REQ-16-033, REQ-17-015, REQ-17-024.

## Out of Scope

- No broker connection, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real DB overwrite, paid/login/authorized source, Level2 source, high-frequency source, future provider availability, or return promise.
- No full original-requirement pass claim; P86 owns the final consolidated claim after P85 and P87 are complete.
