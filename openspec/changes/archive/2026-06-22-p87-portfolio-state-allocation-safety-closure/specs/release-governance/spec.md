## ADDED Requirements

### Requirement: P87 portfolio state allocation safety closure

After P84, the portfolio/allocation/state/data-impact rows not covered by P85 or P86 SHALL NOT be marked `real_pass` unless fresh real local execution proves the complete row-specific behavior through UI operation, API responses, read-only SQLite evidence, deterministic checks where applicable, and explicit forbidden-capability absence.

#### Scenario: P87 row inventory closes the planning gap

- **GIVEN** P87 starts from the P84 evidence matrix
- **WHEN** the P85, P87, and P86 plans are reviewed together
- **THEN** they SHALL cover exactly 157 P84-after full-release-required non-`real_pass` rows
- **AND** no row SHALL be omitted or owned by two execution batches.

#### Scenario: Portfolio state and allocation evidence is row-specific

- **GIVEN** P87 evaluates account, holding-state, allocation, rebalance, or confirmation requirements
- **WHEN** a row is marked `real_pass`
- **THEN** the acceptance evidence SHALL include real browser UI operation, API/readback, SQLite field checks, and deterministic calculations where the value is deterministic
- **AND** broad rows such as monthly attribution or audit history SHALL only pass if the full stated breadth is proven.

#### Scenario: Data-insufficient and release safety remain hard boundaries

- **GIVEN** P87 evaluates degraded data, insufficient evidence, release checks, or safety boundaries
- **WHEN** evidence is unavailable or degraded
- **THEN** the product SHALL visibly qualify or block the affected advice
- **AND** it SHALL NOT create confirmations, trigger trades, suppress blockers, or imply automatic install/upgrade/migration/repair behavior.
