## ADDED Requirements

### Requirement: Expanded real data SHALL enter workflow input context safely

The system SHALL make P34 expanded public data available to daily discipline, expected return, and future risk warning input context without bypassing rules or human review.

#### Scenario: Daily discipline reads expanded data
- **WHEN** DailyDisciplineGraph prepares context for configured holdings or indexes
- **THEN** it SHALL include available P34 valuation, constituent, financial, capital-flow, margin, or sentiment-proxy summaries with source level and freshness status
- **AND** missing or stale categories SHALL remain explicit in the workflow context.

#### Scenario: Expected return reads expanded data
- **WHEN** ExpectedReturnNode prepares scenario material
- **THEN** it SHALL treat P34 data as supporting sample/context material
- **AND** it SHALL preserve sample limitations, freshness, source level, and missing categories
- **AND** it SHALL NOT convert P34 data into guaranteed return, deterministic price, or automatic trading output.

### Requirement: Expanded real data SHALL preserve stub and degradation behavior

The system SHALL keep deterministic local validation possible when expanded public data sources are disabled or unavailable.

#### Scenario: Stub mode is enabled
- **WHEN** `data_sources.use_stub` or an equivalent test fixture mode is enabled
- **THEN** P34 data services SHALL use deterministic fixture or stub data for tests and local validation
- **AND** they SHALL not require public network access.

#### Scenario: Real source mode is enabled but fails
- **WHEN** a real P34 source is enabled and fails with timeout, unavailable source, parse failure, stale data, or no records
- **THEN** the application SHALL return a stable degraded or insufficient-data state
- **AND** it SHALL write auditable failure metadata
- **AND** it SHALL not silently substitute fabricated real data.

### Requirement: Expanded real data SHALL expose health to application surfaces

The system SHALL expose source health and recent refresh status for P34 data through application service DTOs or existing settings/ops surfaces.

#### Scenario: Frontend requests data source health
- **WHEN** the frontend displays data source, daily discipline, or ops status
- **THEN** it SHALL be able to show each P34 source category as fresh, stale, missing, unavailable, parse-error, disabled, or stubbed
- **AND** `GET /api/v1/market/source-health` or an equivalent settings/ops DTO SHALL expose source category, freshness, source level, source type, data date, last success or failure time, failure category, and affected symbols where available
- **AND** it SHALL show last success or failure time where available.

#### Scenario: Health state is insufficient for formal analysis
- **WHEN** required P34 source categories are stale, missing, or failed
- **THEN** frontend and workflow surfaces SHALL explain the missing category
- **AND** they SHALL not present the analysis as fully evidenced.
