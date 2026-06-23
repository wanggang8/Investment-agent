# P82 SOP Action UI SQLite Closure

## Why

P80 proved several review/audit/governance flows, but 53 SOP/action rows still require direct real UI-to-SQLite evaluation. These rows cover the logic that turns risk alerts, SOP state changes, manual actions, notifications, confirmations, and reviews into durable local facts and user-visible readback. P82 may only upgrade rows whose original requirement is directly supported by fresh evidence; broader rows must remain deferred to later batches.

## What Changes

- Define the P82 row set as the SOP/action data-impact batch carried forward from the P81 evidence matrix, preserving its P80 classification provenance.
- Execute representative real browser workflows for SOP state transitions, failure states, manual follow-up, notifications, confirmations, and review readback.
- Verify each user action through API responses and read-only SQLite evidence.
- Produce a P82 evidence layer and acceptance record that separates upgraded rows from evaluated-but-deferred rows.

## In Scope

- 53 evaluated rows: REQ-02-004, REQ-02-018, REQ-02-019, REQ-04-005, REQ-07-011, REQ-08-001, REQ-08-002, REQ-08-003, REQ-08-005, REQ-08-006, REQ-08-007, REQ-08-008, REQ-08-009, REQ-08-010, REQ-08-011, REQ-08-012, REQ-08-013, REQ-08-014, REQ-08-015, REQ-08-016, REQ-08-017, REQ-08-019, REQ-08-021, REQ-08-022, REQ-08-024, REQ-08-025, REQ-08-026, REQ-10-001, REQ-10-002, REQ-10-003, REQ-10-004, REQ-10-005, REQ-12-001, REQ-12-002, REQ-12-003, REQ-13-001, REQ-13-002, REQ-13-003, REQ-13-004, REQ-13-005, REQ-13-007, REQ-13-008, REQ-13-009, REQ-13-011, REQ-13-012, REQ-13-015, REQ-13-016, REQ-13-017, REQ-13-019, REQ-16-016, REQ-16-029, REQ-17-004, REQ-17-010.
- Rows whose original requirement is broader than P82 SOP/action evidence must be deferred with exact next-batch ownership rather than upgraded.

## Out of Scope

- No automatic execution of SOP actions, automatic confirmation, automatic trading, external notification push, or automatic rule application.
- No full original-requirement pass claim from this batch alone.
