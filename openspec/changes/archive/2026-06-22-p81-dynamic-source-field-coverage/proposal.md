# P81 Dynamic Source Field Coverage

## Why

P80 closed a real-use review/audit/governance batch, but 273 full-release-required rows still are not `real_pass`. The largest remaining batch is dynamic external/built-in data source field coverage: 59 rows still need fresh proof that the product obtains the fields from the user-selected fund or ETF, exposes readiness and provenance, and safely degrades when formal evidence is unavailable.

## What Changes

- Define the P81 row set as the dynamic source field coverage batch from the P80 matrix.
- Execute fresh local acceptance that starts from real user-selected symbols instead of hard-coded `510300` assumptions.
- Verify source field readback, feature impact, provenance, freshness, degraded/missing status, LLM context use, and safety boundaries through API/UI/read-only SQLite evidence.
- Produce a P81 evidence layer and release acceptance record without rewriting P75-P80 historical matrices.

## In Scope

- 59 rows: REQ-02-003, REQ-02-009, REQ-02-015, REQ-02-016, REQ-02-023, REQ-02-027, REQ-02-028, REQ-02-030, REQ-04-001, REQ-04-002, REQ-04-004, REQ-04-006, REQ-04-009, REQ-04-010, REQ-04-011, REQ-04-012, REQ-04-013, REQ-04-014, REQ-04-015, REQ-04-017, REQ-04-018, REQ-04-021, REQ-04-022, REQ-04-023, REQ-04-024, REQ-04-026, REQ-04-027, REQ-05-001, REQ-05-002, REQ-05-006, REQ-05-007, REQ-05-008, REQ-05-009, REQ-05-011, REQ-05-012, REQ-05-013, REQ-05-014, REQ-05-015, REQ-05-016, REQ-05-017, REQ-05-018, REQ-05-019, REQ-05-020, REQ-06-001, REQ-06-009, REQ-07-013, REQ-07-014, REQ-14-001, REQ-14-002, REQ-14-003, REQ-15-006, REQ-15-008, REQ-16-012, REQ-16-018, REQ-16-020, REQ-17-006, REQ-17-009, REQ-17-013, REQ-17-021.

## Out of Scope

- No broker connection, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real DB overwrite, return promise, paid/login/authorized source, Level2 source, or high-frequency source.
- No claim of full original-requirement pass until all remaining full-release-required rows are resolved by real evidence or explicitly bounded as non-goal/reference-only.
- No package refresh unless a later package change explicitly requests it.

