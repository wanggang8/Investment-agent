## ADDED Requirements

### Requirement: P82 SOP action UI-to-SQLite closure

After P80, SOP/action data-impact rows SHALL NOT be marked `real_pass` unless a real browser workflow proves the user operation, API result, durable local data impact, auditability, and UI readback.

#### Scenario: P82 row inventory is complete before execution

- **GIVEN** P82 starts from the P81 evidence matrix while preserving P80 classification provenance
- **WHEN** execution begins
- **THEN** the P82 plan SHALL enumerate exactly 53 SOP/action rows for evaluation
- **AND** each row SHALL map to a real UI scenario, readback target, safety negative check, and upgrade-or-defer decision.

#### Scenario: P82 upgrades only directly proven rows

- **GIVEN** a P82 planned row is broader than the fresh SOP/action evidence
- **WHEN** P82 generates the evidence layer
- **THEN** that row SHALL remain non-`real_pass`
- **AND** the acceptance record SHALL name the exact remaining gap and next-batch owner.

#### Scenario: UI operation creates expected local evidence

- **GIVEN** a P82 user action reports success in the UI
- **WHEN** P82 validates the result
- **THEN** it SHALL check API response state, read-only SQLite evidence, audit events, and visible readback after navigation or refresh.

#### Scenario: Unsupported automation remains blocked

- **GIVEN** a scenario concerns SOP, confirmations, notifications, or rule governance
- **WHEN** the UI is inspected
- **THEN** it SHALL NOT expose automatic trading, one-click trading, order delegation, external push, automatic confirmation, or automatic rule application as available product actions.
