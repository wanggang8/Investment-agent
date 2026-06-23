## ADDED Requirements

### Requirement: P66 SHALL classify current data source quality with a release policy

P66 SHALL extend data source quality regression with a current-data policy verdict that separates clean release evidence, explicit waiver evidence, and release-blocking current data quality states.

#### Scenario: Current regression returns policy verdict

- **WHEN** data source quality regression runs in `current` mode
- **THEN** the response SHALL include a `policy` object with `verdict`, `release_gate`, degraded/failed/blocking/waiver counts, sanitized reasons, manual next actions, and a safety note
- **AND** the policy SHALL be derived only from existing local source-health facts
- **AND** it SHALL NOT refresh data, call public providers, call LLM providers, write rules, confirm actions, repair files, overwrite databases, or execute trades.

#### Scenario: Current policy passes

- **WHEN** current regression has source-health facts and every case is `passed`
- **THEN** the policy verdict SHALL be `passed`
- **AND** the release gate SHALL be `pass`.

#### Scenario: Current policy requires waiver

- **WHEN** current regression contains only recognized degraded optional categories and no failed cases, no missing source-health facts, no explicit `freshness=missing`, no unrecognized failure category, and no core category degradation
- **THEN** the policy verdict SHALL be `waiver_required`
- **AND** the release gate SHALL be `waiver_required`
- **AND** release materials SHALL explicitly document the waiver before making a release-ready claim.
- **AND** an optional degraded category MAY appear in the legacy `missing_categories` list without automatically becoming `blocked`.

#### Scenario: Current policy blocks release

- **WHEN** current regression has no source-health facts, explicit `freshness=missing`, failed/unrecognized freshness, unrecognized failure category, or degraded core categories
- **THEN** the policy verdict SHALL be `blocked`
- **AND** the release gate SHALL be `block`
- **AND** release materials SHALL NOT claim current data-source quality is clean.

#### Scenario: Fixture policy remains deterministic

- **WHEN** data source quality regression runs in `fixture` mode
- **THEN** the policy verdict SHALL be `passed`
- **AND** fixture mode SHALL remain an offline deterministic classification and redaction regression, not evidence that current local data is healthy.

### Requirement: P66 SHALL expose current data policy through existing local surfaces

P66 SHALL expose the current data-source quality policy through the existing local API, CLI, and Data Quality UI without expanding automation boundaries.

#### Scenario: API response is requested

- **WHEN** a user calls `GET /api/v1/data-source-quality/regression?mode=current`
- **THEN** the API SHALL return the policy object in the response
- **AND** it SHALL remain read-only and SHALL NOT write audit records.

#### Scenario: CLI task is executed

- **WHEN** a user runs the data-source-quality CLI task
- **THEN** the compact audit output SHALL include sanitized `policy` and `gate` fields
- **AND** a strict current-data gate SHALL be available for release acceptance and SHALL return non-zero when the policy release gate is `block`.

#### Scenario: Data Quality page is viewed

- **WHEN** the Data Quality page loads current data-source quality policy
- **THEN** it SHALL display `passed`, `waiver_required`, or `blocked` with appropriate tone, sanitized reasons, and manual next actions
- **AND** it SHALL NOT provide automatic refresh, automatic repair, automatic confirmation, automatic rule application, external push, broker, order, or trading actions.

### Requirement: P66 SHALL preserve source and release safety boundaries

P66 SHALL remain a local policy and acceptance hardening phase.

#### Scenario: Policy detects degradation

- **WHEN** the policy detects degraded or blocked current data
- **THEN** it SHALL explain release impact and manual next actions
- **AND** it SHALL NOT claim future public-source availability, current data health, provider availability, investment returns, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic upgrade, automatic migration, real database overwrite, login-gated sources, paid sources, authorization-gated sources, Level2 data, or high-frequency data.
