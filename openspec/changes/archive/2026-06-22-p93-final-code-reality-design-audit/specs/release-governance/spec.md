## ADDED Requirements

### Requirement: P93 final code reality and design audit

After P92, the project SHALL provide a code-facing release audit that checks whether original requirements are backed by real implementation rather than demo, placeholder, hardcoded, or dead-code behavior.

#### Scenario: Production implementation evidence is mapped

- **GIVEN** the original requirements are accepted by P92
- **WHEN** P93 audits code reality
- **THEN** it SHALL map original requirement sections to concrete production Go, React, configuration, script, and release files
- **AND** it SHALL cross-check the P92 row-level ledger so every original requirement row resolves to a current code/evidence bundle through its source section
- **AND** it SHALL identify whether each requirement area is backed by runtime code, UI, tests, and evidence
- **AND** it SHALL keep P92 as the 341-row row-level artifact rather than replacing it with a coarser P93 claim.

#### Scenario: Demo and hardcoding risks are classified

- **GIVEN** suspicious terms such as `demo`, `mock`, `stub`, `placeholder`, `fake`, `dummy`, `TODO`, `FIXME`, or hardcoded values appear in the repository
- **WHEN** P93 evaluates them
- **THEN** it SHALL classify each material occurrence as test-only, config-only, documentation-only, accepted local fallback, or release-blocking
- **AND** release-blocking occurrences SHALL be fixed or reported as blockers.

#### Scenario: Secret literals are blocked

- **GIVEN** current non-test source or configuration files may contain local credentials
- **WHEN** P93 evaluates hardcoding risk
- **THEN** it SHALL scan for unredacted `sk-...` API key literals using a bounded token pattern
- **AND** any such literal in scanned non-test source/config files SHALL be release-blocking until removed or replaced with an empty/user-supplied runtime value.

#### Scenario: Final claims remain bounded

- **GIVEN** P93 passes
- **WHEN** release readiness is described
- **THEN** it MAY claim the implementation has passed final code reality and design audit for the local/GitHub-Docker release scope
- **AND** it SHALL NOT claim physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.
