## ADDED Requirements

### Requirement: Expand readonly public data coverage for daily discipline

The system SHALL expand real public data coverage for daily discipline, expected return, and future risk warning contexts while preserving local-only, readonly, low-frequency, and degradable behavior.

#### Scenario: Extended public data refresh succeeds
- **WHEN** a configured P34 public data source returns valid index, constituent, valuation, financial, capital-flow, margin-financing, or sentiment-proxy data
- **THEN** the system SHALL normalize the payload with source name, source level, source type, symbol or index code, data date, captured time, content hash, and raw metadata
- **AND** it SHALL persist the data through the existing market, evidence, or explicitly added source-health persistence path
- **AND** it SHALL write an audit event with successful source counts and affected symbols.

#### Scenario: Extended public data is unavailable
- **WHEN** a configured P34 public data source has no matching records, unavailable endpoints, incompatible response shape, stale data, or parse failures
- **THEN** the system SHALL classify the result as `no_data`, `source_unavailable`, `parse_error`, `stale`, or an equivalent documented status
- **AND** it SHALL not write fabricated metric values
- **AND** it SHALL preserve other successful source results when partial refresh is possible.

#### Scenario: Extended public data violates safety boundary
- **WHEN** a candidate source requires login, paid access, authorization, CAPTCHA bypass, broker account access, Level2 data, high-frequency polling, or access-control circumvention
- **THEN** the system SHALL exclude that source from P34 runtime implementation
- **AND** it SHALL continue with fixture, stub, or other verified public sources.

### Requirement: Track source health and freshness for expanded data

The system SHALL expose source health and freshness information for P34 data so users and workflows can distinguish usable, stale, missing, and failed inputs.

#### Scenario: Source health is recorded
- **WHEN** a P34 refresh attempts any configured expanded public source
- **THEN** the system SHALL record source name, source type, last success time, last failure time, failure category, data date, affected symbols, and source level where available
- **AND** this health state SHALL be queryable by application services and frontend status views.

#### Scenario: Freshness affects workflow context
- **WHEN** expanded data is stale, missing, or lower-grade only
- **THEN** DailyDisciplineGraph and expected return context SHALL mark the corresponding input as stale, missing, degraded, or insufficient
- **AND** the workflow SHALL not silently treat the input as fresh A-level evidence.

### Requirement: Keep expanded data read-only and non-trading

The system SHALL ensure P34 data refresh, storage, health reporting, and workflow consumption do not create trading side effects.

#### Scenario: P34 refresh completes
- **WHEN** P34 refresh collects or fails to collect expanded public data
- **THEN** it SHALL NOT update local positions, portfolio snapshots, operation confirmations, broker state, orders, external notifications, formal rule versions, local account import batches, or local account corrections
- **AND** it SHALL only write local data facts, source health, and audit records.

#### Scenario: P34 data enters analysis context
- **WHEN** P34 data is used by daily discipline, expected return, or risk context preparation
- **THEN** it SHALL remain analysis material subject to rules and human review
- **AND** it SHALL NOT produce automatic trading actions, one-click order controls, profit guarantees, or deterministic price predictions.
