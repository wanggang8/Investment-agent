## ADDED Requirements

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
