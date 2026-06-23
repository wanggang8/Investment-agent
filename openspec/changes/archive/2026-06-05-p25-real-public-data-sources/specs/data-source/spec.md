# 数据源真实接入调研 Delta

## ADDED Requirements

### Requirement: P25 SHALL verify public data source access before implementation

P25 SHALL treat real public data source integration as a research and verification phase before production collectors are implemented.

#### Scenario: Candidate source is promoted to implementation scope

- **WHEN** a candidate source is considered for P26 or P27 implementation
- **THEN** the project SHALL record its authority level, public access shape, stable request or endpoint evidence, available fields, update cadence, legal or disclaimer constraints, rate limit assumptions, failure behavior, and target persistence path
- **AND** the source SHALL NOT be marked implementation-ready until this record exists.

#### Scenario: Page exists but stable access is not verified

- **WHEN** a public page exists but stable request parameters, pagination, response fields, or access constraints are unknown
- **THEN** the source SHALL remain in a research or blocked state
- **AND** follow-up implementation SHALL NOT rely on that source as a required runtime dependency.

### Requirement: P19/P20 SHALL be described as infrastructure, not completed external source integration

P19/P20 documentation SHALL distinguish existing HTTP bridge and payload parsing capabilities from unimplemented real external source collectors.

#### Scenario: Planning documents describe P19/P20

- **WHEN** planning or configuration documents summarize P19/P20
- **THEN** they SHALL state that configurable HTTP bridge, parser, fixture/stub fallback, source level mapping, and degradation behavior are complete
- **AND** they SHALL state that real public authority source collectors, historical backfill, and live smoke verification remain P25+ work.

### Requirement: Public data integration SHALL preserve local-only safety boundaries

Real public data source plans SHALL preserve the system's safety and governance boundaries.

#### Scenario: A data source requires login, paid access, broker access, or authorized market data

- **WHEN** a candidate source requires login, paid access, broker trading access, Level2 or authorized market data, or bypassing access controls
- **THEN** it SHALL be excluded from the default implementation scope
- **AND** the system SHALL NOT use it unless a later approved change explicitly redefines the boundary.

#### Scenario: A third-party aggregation source provides market or fund data

- **WHEN** a third-party aggregation source such as a finance portal provides market, fund, or news data
- **THEN** it SHALL default to background or B-level evidence unless the project records a stronger authority rationale
- **AND** it SHALL NOT be the sole source used to satisfy high-confidence formal evidence requirements.
