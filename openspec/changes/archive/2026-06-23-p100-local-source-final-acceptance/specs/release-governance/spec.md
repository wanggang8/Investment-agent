## ADDED Requirements

### Requirement: P100 local source final acceptance

After P99, the project SHALL provide a final acceptance pass for the local source runtime before making any new local-source release-readiness claim.

#### Scenario: Local source scope is explicit

- **GIVEN** the project has completed P92, P93, and P99
- **WHEN** P100 acceptance is executed
- **THEN** it SHALL use the local Go backend, Vite frontend, SQLite, VecLite, real browser UI, API readback, SQLite/readback, and audit evidence as its acceptance scope
- **AND** it SHALL explicitly exclude Docker Compose, install, upgrade, uninstall, purge, package refresh, Git tag, GitHub Release, and physical second-machine validation.

#### Scenario: Requirement and code reality gates are preserved

- **GIVEN** P92 and P93 are the latest requirement and code-reality audit layers
- **WHEN** P100 evaluates local source readiness
- **THEN** it SHALL run the P92 final requirement audit check
- **AND** it SHALL run the P93 code reality audit check
- **AND** it SHALL block a passing conclusion if any full-release-required row is no longer `real_pass` or any active release-blocking code-reality finding exists.

#### Scenario: Product usability and design are reviewed

- **GIVEN** the local source runtime is started for acceptance
- **WHEN** P100 reviews the product in a real browser
- **THEN** it SHALL verify critical product journeys across workbench, portfolio maintenance, data refresh, consultation, decision detail, review, governance, audit, notifications, and data quality surfaces
- **AND** it SHALL record whether the UI communicates next actions, states, manual confirmation, evidence links, and safety boundaries clearly
- **AND** it SHALL include responsive review for narrow, tablet, and desktop widths.

#### Scenario: Final claim remains bounded

- **GIVEN** P100 produces a final acceptance record
- **WHEN** release readiness is described
- **THEN** it MAY claim local source final acceptance only when the machine gates, local runtime journeys, data-impact evidence, and design rubric pass or have documented non-blocking degradation
- **AND** it SHALL NOT claim Docker installation, package distribution, GitHub Release, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic recovery, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.
