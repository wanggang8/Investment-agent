## ADDED Requirements

### Requirement: P77 SHALL govern post-P75 real-pass upgrades

P77 SHALL create a conservative, auditable upgrade gate for moving P75 atomic requirement rows toward `real_pass` without rewriting historical P75 acceptance evidence or expanding unsupported release claims.

#### Scenario: P77 upgrade evidence is a new layer

- **GIVEN** P75 produced `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md`
- **WHEN** P77 evaluates row-level status upgrades
- **THEN** P77 SHALL generate a new matrix that preserves P75 row IDs, source line ranges, requirement text hashes, original statuses, full-release-required flags, and release impacts
- **AND** P77 SHALL NOT mutate the historical P75 matrix to make prior evidence appear stronger than it was
- **AND** each P77 row SHALL record `p77_status`, upgrade basis, gate dimensions, fresh evidence command, fresh evidence artifact, residual gap, and next remediation.

#### Scenario: Real-pass upgrade requires all applicable evidence dimensions

- **GIVEN** a P75 row is being considered for `real_pass`
- **WHEN** P77 evaluates that row
- **THEN** the row SHALL have implementation evidence
- **AND** user-visible behavior SHALL have real UI evidence
- **AND** mutating behavior SHALL have SQLite changed-table, prohibited-table, audit-event, and readback evidence
- **AND** data-source, collector, workflow, rule, LLM, RAG, and scenario-dependent behavior SHALL have direct evidence for each applicable dependency
- **AND** safety evidence SHALL confirm no broker interface, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, future provider-availability promise, or investment return promise is introduced
- **AND** the row SHALL remain non-`real_pass` if its only evidence is screenshot-only, route-smoke-only, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, temporary-DB-only, or incompatible single-symbol-only.

#### Scenario: P77 release conclusion remains bounded

- **WHEN** P77 reports its final conclusion
- **THEN** it SHALL claim `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`
- **AND** otherwise it SHALL use a scoped conclusion that lists the remaining non-`real_pass` rows or grouped categories with row-level matrix reference
- **AND** it SHALL preserve P76 package boundaries unless a separate package refresh change is executed.
