# data-source-quality-regression Specification

## Purpose
TBD - created by archiving change p48-data-source-quality-regression-pack. Update Purpose after archive.
## Requirements
### Requirement: P48 SHALL provide repeatable data source quality regression

The system SHALL provide a local data source quality regression capability that verifies source health freshness, failure classification, and safe summaries without requiring external network access by default.

#### Scenario: Fixture regression runs offline
- **WHEN** P48 data source quality regression runs in `fixture` mode
- **THEN** it SHALL evaluate deterministic local cases for `fresh`, `no_data`, `source_unavailable`, `parse_error`, `stale`, and sensitive diagnostic redaction
- **AND** it SHALL return a stable summary containing case status, expected freshness, actual freshness, data category, affected symbols, and safety note
- **AND** it SHALL NOT access public endpoints, private files, broker systems, external notification channels, or trading APIs.

#### Scenario: Current source health is evaluated read-only
- **WHEN** P48 data source quality regression runs in `current` mode
- **THEN** it SHALL evaluate the latest existing P34 source health from local market snapshots
- **AND** it SHALL NOT trigger collectors, refresh market data, rebuild indexes, call LLMs, create notifications, update rules, update confirmations, or change account and position facts.

### Requirement: P48 SHALL expose sanitized regression summaries

P48 SHALL expose data source quality regression through local API and CLI summaries that do not reveal sensitive diagnostics.

#### Scenario: Regression API is requested
- **WHEN** a user calls `GET /api/v1/data-source-quality/regression`
- **THEN** the API SHALL return `mode`, `status`, `generated_at`, `summary`, `cases`, `missing_categories`, and `safety_note`
- **AND** each case SHALL include only sanitized `diagnostic_preview` text
- **AND** the API SHALL NOT write SQLite records or trigger external side effects.

#### Scenario: Regression CLI task is executed
- **WHEN** a user runs `go run ./cmd/agent --task data-source-quality-regression`
- **THEN** the task SHALL run the same regression service and print a compact local summary
- **AND** it MAY write a local `audit_events` record containing only mode, status, case counts, degraded/failed counts, and safety boundary
- **AND** it SHALL NOT store raw source payloads, complete API keys, private paths, raw SQL, full prompts, raw HTTP exchanges, private keys, or supplier raw responses in output or audit metadata.

### Requirement: P48 SHALL preserve data source safety boundaries

P48 SHALL remain a local quality regression feature and SHALL NOT expand the system's data acquisition, trading, or automation boundaries.

#### Scenario: A requested regression source is outside the allowed boundary
- **WHEN** a requested mode or source would require login, paid access, authorization-gated data, CAPTCHA bypass, broker access, Level2 data, high-frequency polling, browser scraping, or access-control circumvention
- **THEN** P48 SHALL reject or omit that mode/source
- **AND** it SHALL continue to provide fixture regression and current local source health evaluation where available.

#### Scenario: Regression detects degraded data quality
- **WHEN** regression cases detect `no_data`, `source_unavailable`, `parse_error`, `stale`, `missing`, `unknown`, or another degraded but recognized condition
- **THEN** the system SHALL classify the regression as degraded or failed according to documented rules
- **AND** it SHALL NOT automatically repair data, mark sources healthy, execute trades, apply rules, confirm operations, or promise returns.

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

### Requirement: P67 SHALL record manual current-data gate resolutions locally

P67 SHALL add a local, auditable workflow for explicit current-data gate waiver or scope exclusion records while preserving the P66 policy verdict as the underlying fact.

#### Scenario: Resolution record is created

- **WHEN** the current P66 data-source policy is `blocked` or `waiver_required`
- **AND** a user submits a valid resolution type, scope, reason, release impact, and optional evidence reference
- **THEN** the system SHALL create or reuse an active local resolution record tied to the current policy fingerprint
- **AND** the record SHALL store only sanitized text and copied policy reasons
- **AND** it SHALL NOT refresh data, repair source health, call public providers, call LLM providers, rebuild indexes, apply rules, confirm actions, overwrite data, or execute trades.

#### Scenario: Resolution record is retired

- **WHEN** a user retires an active current-data gate resolution record
- **THEN** the system SHALL mark only that local resolution record as `retired`
- **AND** subsequent release checks SHALL ignore the retired record
- **AND** it SHALL NOT mutate source health, market snapshots, evidence, rules, confirmations, accounts, positions, or decisions.

#### Scenario: Resolution text is unsafe

- **WHEN** submitted scope, reason, release impact, or evidence text contains complete keys, private paths, raw SQL, full prompts, raw HTTP exchanges, provider raw payloads, stack traces, or local database paths
- **THEN** the system SHALL persist and return only sanitized previews
- **AND** local audit output SHALL contain only compact sanitized references.

### Requirement: P67 SHALL derive release claim state without weakening P66

P67 SHALL expose a release claim state that combines the current P66 policy with the latest matching active manual resolution.

#### Scenario: Policy fingerprint is computed

- **WHEN** current data-source quality policy is evaluated
- **THEN** the system SHALL compute a stable `policy_fingerprint` from symbol, policy verdict, release gate, blocking reasons, waiver reasons, degraded count, failed count, blocking count, waiver count, and normalized current regression case categories
- **AND** matching SHALL use `policy_fingerprint`, not display summary text.

#### Scenario: Current policy passes

- **WHEN** the current P66 policy verdict is `passed`
- **THEN** the release claim state SHALL be `pass`
- **AND** `clean_data_claim_allowed` SHALL be true.

#### Scenario: Current policy needs resolution and no active record exists

- **WHEN** the current P66 policy verdict is `blocked` or `waiver_required`
- **AND** no active matching resolution record exists
- **THEN** the release claim state SHALL be `requires_resolution`
- **AND** future release-ready claims SHALL remain blocked or limited until a resolution is recorded or the current policy passes.

#### Scenario: Current waiver-required policy has an active waiver

- **WHEN** the current P66 policy verdict is `waiver_required`
- **AND** a matching active `waiver` record exists
- **THEN** the release claim state SHALL be `resolved_with_waiver`
- **AND** `clean_data_claim_allowed` SHALL remain false unless the P66 policy itself is `passed`.

#### Scenario: Current blocked policy attempts waiver

- **WHEN** the current P66 policy verdict is `blocked`
- **AND** a user submits `resolution_type=waiver`
- **THEN** the system SHALL reject the request
- **AND** it SHALL explain that blocked current data requires source-health recovery or explicit scope exclusion.

#### Scenario: Current policy has an active scope exclusion

- **WHEN** the current P66 policy verdict is `blocked` or `waiver_required`
- **AND** a matching active `scope_exclusion` record exists
- **THEN** the release claim state SHALL be `resolved_with_scope_exclusion`
- **AND** release materials SHALL explicitly state that current local data health is excluded from clean-data claims.

#### Scenario: Duplicate active resolution would be created

- **WHEN** an active resolution already exists for the same symbol and `policy_fingerprint`
- **AND** the user submits the same resolution type
- **THEN** the system SHALL reuse the active record
- **AND** it SHALL NOT create a duplicate active record.

#### Scenario: Conflicting active resolution would be created

- **WHEN** an active resolution already exists for the same symbol and `policy_fingerprint`
- **AND** the user submits a different resolution type
- **THEN** the system SHALL reject the request until the existing resolution is retired
- **AND** release claim state SHALL remain based on the single active record.

### Requirement: P67 SHALL expose resolution workflow through local API, CLI, and Data Quality UI

P67 SHALL expose current-data gate resolution state through local surfaces without adding automation or external side effects.

#### Scenario: API resolution check is requested

- **WHEN** a user calls `GET /api/v1/data-source-quality/gate-resolution?symbol=000300`
- **THEN** the API SHALL return the current P66 policy, active resolution if any, release claim state, allowed/prohibited claim labels, and a safety note
- **AND** the API SHALL be read-only and SHALL NOT write audit records.

#### Scenario: API creates or retires a resolution

- **WHEN** a user calls the local create or retire resolution API
- **THEN** the API SHALL write only the local resolution record and a sanitized audit event
- **AND** it SHALL NOT trigger collectors, source refresh, repair, external push, rule application, confirmation, broker, order, or trading actions.

#### Scenario: CLI resolution check is executed

- **WHEN** a user runs `go run ./cmd/agent --task data-source-quality-resolution-check --source current --symbol 000300`
- **THEN** the command SHALL print compact sanitized policy, gate, resolution, and claim state fields
- **AND** it SHALL return non-zero only when the claim state is `requires_resolution`
- **AND** a blocked policy SHALL return zero only when a matching active `scope_exclusion` exists
- **AND** it SHALL NOT change P66 strict policy gate behavior.

#### Scenario: Data Quality page is viewed

- **WHEN** the Data Quality page loads current-data gate resolution state
- **THEN** it SHALL show whether the current gate is unresolved, resolved with waiver, resolved with scope exclusion, or passed
- **AND** it SHALL provide only local record and retire-record controls
- **AND** it SHALL NOT provide automatic refresh, automatic repair, automatic confirmation, automatic rule application, external push, broker, order, or trading actions.

