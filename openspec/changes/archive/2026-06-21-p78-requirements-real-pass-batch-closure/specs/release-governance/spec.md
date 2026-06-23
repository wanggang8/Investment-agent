## ADDED Requirements

### Requirement: P78 SHALL close real-pass gaps in conservative batches

P78 SHALL continue post-P77 full-requirement acceptance by classifying remaining non-`real_pass` rows into remediation batches and upgrading only rows that meet the P77 evidence gate with fresh, row-specific evidence.

#### Scenario: P78 classifies remaining gaps before upgrading rows

- **GIVEN** P77 produced `docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md`
- **WHEN** P78 evaluates the remaining full-release-required non-`real_pass` rows
- **THEN** P78 SHALL generate a new matrix that preserves P77 row IDs, source ranges, original status, P77 status, full-release-required flag, and release impact
- **AND** each non-`real_pass` row SHALL receive a remediation group, batch assignment, remaining gap, and next action
- **AND** P78 SHALL NOT mutate P75 or P77 historical matrices.

#### Scenario: P78 batch upgrades require direct evidence

- **GIVEN** a P78 batch proposes a row for `real_pass`
- **WHEN** the P78 checker evaluates the row
- **THEN** implementation behavior SHALL be backed by fresh deterministic tests or direct runtime evidence
- **AND** user-visible behavior SHALL have real browser UI evidence when applicable
- **AND** data-bearing behavior SHALL have SQLite readback evidence for the exact fields claimed
- **AND** expected-return or analysis rows SHALL show sample count, sample window, screening condition, source/provenance fields, precision/degradation status, and non-trading disclaimer when applicable
- **AND** the row SHALL remain non-`real_pass` if evidence is only inherited, screenshot-only, route-smoke-only, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, temporary-DB-only, or incompatible single-symbol-only.

#### Scenario: P78 release conclusion remains bounded

- **WHEN** P78 reports its conclusion
- **THEN** it SHALL claim `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`
- **AND** otherwise it SHALL use a scoped conclusion that records upgraded row count, remaining non-`real_pass` row count, remediation groups, and package freshness boundaries
- **AND** it SHALL NOT claim P76 package inclusion unless a separate package refresh change is executed.
