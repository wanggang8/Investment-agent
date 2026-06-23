## ADDED Requirements

### Requirement: Minimum viable readonly providers
The system SHALL support at least one configurable readonly market data provider and at least one configurable readonly intelligence provider while preserving local stub mode.

#### Scenario: Readonly market provider succeeds
- **WHEN** the configured readonly market provider returns valid data for a requested symbol
- **THEN** the system SHALL validate the payload
- **THEN** the system SHALL write a `market_snapshots` fact through the existing persistence path
- **THEN** the system SHALL write a success audit event

#### Scenario: Readonly intelligence provider succeeds
- **WHEN** the configured readonly intelligence provider returns valid intelligence for a requested symbol or topic
- **THEN** the system SHALL preserve source name, source URL when available, published time, captured time, and source level
- **THEN** the system SHALL write local intelligence facts through the existing persistence path
- **THEN** the system SHALL write a success audit event

#### Scenario: Stub mode remains available
- **WHEN** stub mode is enabled or no real provider is configured
- **THEN** the system SHALL use local deterministic stub data for development and validation
- **THEN** the system SHALL NOT require public network access

### Requirement: Stable provider degradation
The system SHALL return stable degraded states and audit events for provider timeout, source unavailable, stale data, and parse failure.

#### Scenario: Provider timeout
- **WHEN** a configured provider exceeds the configured timeout
- **THEN** the system SHALL return a stable source-unavailable or degraded result
- **THEN** the system SHALL write an audit event with the timeout reason

#### Scenario: Provider parse failure
- **WHEN** a provider response cannot be parsed into the expected local fact shape
- **THEN** the system SHALL return a stable parse-failure or degraded result
- **THEN** the system SHALL NOT write invalid facts
- **THEN** the system SHALL write an audit event with the provider name and failure category

### Requirement: P12 out-of-scope boundaries
The system SHALL NOT expand P12 into complete financial data integration, complete sentiment data integration, realtime SLA, broker trading APIs, automatic trading, active recommendations, or return guarantees.

#### Scenario: P12 provider is enabled
- **WHEN** a real readonly provider is enabled
- **THEN** the system SHALL only refresh, ingest, index, summarize, or audit local facts
- **THEN** the system SHALL NOT place trades, expose broker actions, or mutate account state without user-recorded offline confirmation
