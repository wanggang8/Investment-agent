## ADDED Requirements

### Requirement: Frontend ops status surface
The frontend SHALL display operator-facing status for local data, index, and review readiness using API/service DTOs only.

#### Scenario: Data and index status are visible
- **WHEN** data source, market freshness, index health, or review readiness data is available
- **THEN** the frontend SHALL show a visible status panel with success, degraded, failed, or empty state labels
- **THEN** the panel SHALL NOT read SQLite, VecLite, or local files directly

#### Scenario: Degraded state is distinguishable
- **WHEN** VecLite, data source, DeepSeek, or review data is degraded or unavailable
- **THEN** the frontend SHALL distinguish degraded from successful and failed states
- **THEN** the user SHALL see a safe explanation that does not imply automatic recovery or trading action

### Requirement: Frontend review summary surface
The frontend SHALL display periodic review summaries, rule suggestions, and related counts using DTO data.

#### Scenario: Review summary is visible
- **WHEN** monthly or quarterly review summary data is available
- **THEN** the frontend SHALL show period, confirmation counts, error counts, audit counts, degradation counts, and rule suggestion counts where present

#### Scenario: Review empty state is safe
- **WHEN** review summary data is empty
- **THEN** the frontend SHALL show an empty state that indicates no review facts are available yet
- **THEN** it SHALL NOT present rule application or trading actions

### Requirement: Cross-page tracking entrypoints
The frontend SHALL provide tracking entrypoints from review or ops surfaces to related records.

#### Scenario: Tracking links are visible
- **WHEN** a review summary references audit events, rule proposals, decisions, or error cases
- **THEN** the frontend SHALL expose visible tracking links or entrypoints
- **THEN** those entrypoints SHALL only navigate or filter related records

#### Scenario: Tracking entrypoints are safe
- **WHEN** a tracking entrypoint references a rule proposal or confirmation record
- **THEN** the frontend SHALL NOT expose automatic rule application, automatic trading, one-click order placement, or broker order behavior
