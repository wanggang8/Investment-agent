# daily-auto-run Specification

## Purpose
TBD - created by archiving change p31-daily-auto-run-loop. Update Purpose after archive.
## Requirements
### Requirement: Local daily auto-run shall be explicitly enabled and safely schedulable

The system SHALL provide a local daily auto-run capability that is disabled by default, can be explicitly enabled through local configuration, and executes only within the project's non-trading safety boundaries.

#### Scenario: Auto-run remains disabled by default

- **GIVEN** a default local configuration
- **WHEN** the backend server or local agent starts
- **THEN** no daily automatic refresh or decision generation SHALL run without an explicit enablement setting
- **AND** the system SHALL NOT create decision records, notifications, or audit events solely because the server started

#### Scenario: Enabled auto-run schedules the next local run

- **GIVEN** a local configuration with daily auto-run explicitly enabled and a valid local run time
- **WHEN** the server or scheduler-capable entrypoint starts
- **THEN** the system SHALL expose or record the next planned run time
- **AND** the configured schedule SHALL remain local-only
- **AND** it SHALL NOT require broker APIs, external notification channels, login-only sources, paid sources, or Level2/user-identity market data

### Requirement: Daily auto-run shall orchestrate refresh, decision, notification, and audit steps

The system SHALL orchestrate a daily run by using the configured portfolio or symbol scope, refreshing available local market/evidence data, preparing evidence for decision workflows, running the daily discipline decision path, and recording both user-visible status and audit events.

#### Scenario: Successful daily auto-run creates a traceable daily result

- **GIVEN** auto-run is enabled, the account/portfolio prerequisites are present, and required local data sources are available
- **WHEN** the configured daily run time is reached or the daily auto-run is manually triggered for validation
- **THEN** the system SHALL refresh market data and public evidence according to local configuration
- **AND** it SHALL run the daily discipline workflow for the configured scope
- **AND** it SHALL write traceable `decision_records`, related evidence references, application notifications, and `audit_events`
- **AND** the frontend SHALL be able to show the latest run summary, status, and result link

#### Scenario: Partial data refresh does not fabricate a recommendation

- **GIVEN** auto-run is enabled and some configured data sources fail or return no data
- **WHEN** the daily run executes
- **THEN** the run SHALL classify the failure or no-data state
- **AND** it SHALL continue only for paths whose prerequisites remain satisfied
- **AND** it SHALL record degraded or insufficient-data status where prerequisites are missing
- **AND** it SHALL NOT fabricate market data, evidence, expected returns, or trading instructions

### Requirement: Daily auto-run shall be idempotent per run scope and date

The system SHALL prevent repeated automatic runs from creating conflicting duplicate daily decision records for the same local date, portfolio scope, and configured symbol set.

#### Scenario: Re-running the same daily scope is idempotent

- **GIVEN** a daily auto-run has already completed or produced a recorded degraded result for a local date and scope
- **WHEN** the scheduler retries or the same run is triggered again with the same idempotency key
- **THEN** the system SHALL avoid creating conflicting duplicate decision records
- **AND** it SHALL either reuse the prior result, record a retry attempt, or create an explicitly versioned rerun according to documented behavior
- **AND** all behavior SHALL be visible in audit events

### Requirement: Daily auto-run status shall be visible and diagnosable

The system SHALL make daily auto-run status visible through local API/UI surfaces and application notifications so users can understand whether the run succeeded, partially succeeded, failed, or is disabled.

#### Scenario: Frontend displays last run, next run, and failure reason

- **GIVEN** the user opens the local frontend after daily auto-run is configured
- **WHEN** the frontend requests scheduler or daily run status
- **THEN** it SHALL display whether auto-run is disabled, scheduled, running, succeeded, partially succeeded, failed, or degraded
- **AND** it SHALL show the last run time, next planned run time when available, and a user-readable failure reason when applicable
- **AND** it SHALL provide a link to the relevant decision, notification, or audit details when available

#### Scenario: Missing prerequisites are reported instead of hidden

- **GIVEN** auto-run is enabled but the account snapshot, positions, market snapshots, evidence, rules, or configuration prerequisites are missing
- **WHEN** the daily run attempts to execute
- **THEN** the system SHALL record the missing prerequisites as a structured status
- **AND** it SHALL notify the user in-app
- **AND** it SHALL NOT generate a formal trade-advice decision that implies sufficient data exists

### Requirement: Daily auto-run shall preserve non-trading and non-prediction boundaries

The system SHALL keep daily auto-run inside the same safety boundaries as manual workflows: no automatic trading, no external push channels, no unauthorized data access, no guaranteed returns, and no deterministic price predictions.

#### Scenario: Auto-run cannot execute trading operations

- **GIVEN** a daily auto-run is scheduled, running, retrying, or failing
- **WHEN** any orchestration step evaluates recommendations, confirmations, or notifications
- **THEN** it SHALL NOT call brokerage APIs or create buy/sell/cancel/modify order requests
- **AND** it SHALL NOT mark operations as executed manually without user input
- **AND** it SHALL only create local records, analysis materials, notifications, and audit events

