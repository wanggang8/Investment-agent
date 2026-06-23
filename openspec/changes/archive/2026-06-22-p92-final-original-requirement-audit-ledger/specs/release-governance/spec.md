## ADDED Requirements

### Requirement: P92 final original requirement audit ledger

After P91, the project SHALL provide a final original-requirement audit ledger that independently summarizes whether every original requirement row is covered by final acceptance evidence.

#### Scenario: Full-release rows must all be final real pass

- **GIVEN** P75 generated the original requirement traceability matrix
- **AND** P88 produced the latest full 341-row evidence matrix
- **AND** P89 and P90 produced final blocker overlays
- **WHEN** P92 generates the final audit ledger
- **THEN** every full-release-required row SHALL have final status `real_pass`
- **AND** reference-only rows SHALL remain separated from product pass claims
- **AND** the ledger SHALL fail validation if any full-release-required row is missing, stale, or non-`real_pass`.

#### Scenario: Ledger includes operational review dimensions

- **GIVEN** an original requirement row is included in the final ledger
- **WHEN** P92 writes the row
- **THEN** it SHALL include the requirement id, source section, requirement text, final status, feature area, UI/product surface, expected behavior or data impact, readback or audit evidence, acceptance command or artifact, and boundary notes.

#### Scenario: Final audit claims remain bounded

- **GIVEN** P92 summarizes final release readiness
- **WHEN** it describes accepted scope
- **THEN** it MAY claim original product requirements are accepted for the local/GitHub-Docker release scope
- **AND** it SHALL NOT claim physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.
