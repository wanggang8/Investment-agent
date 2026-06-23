## MODIFIED Requirements

### Requirement: Frontend ops status surface
The frontend SHALL display operator-facing status for local data, index, retrieval quality, and review readiness using API/service DTOs only.

#### Scenario: Data, index, and retrieval quality status are visible
- **WHEN** data source, market freshness, index health, retrieval quality, or review readiness data is available
- **THEN** the frontend SHALL show a visible status panel with success, degraded, failed, or empty state labels
- **THEN** the panel SHALL NOT read SQLite, VecLite, or local files directly

#### Scenario: Degraded state is distinguishable
- **WHEN** VecLite, retrieval quality, data source, DeepSeek, or review data is degraded or unavailable
- **THEN** the frontend SHALL distinguish degraded from successful and failed states
- **THEN** the user SHALL see a safe explanation that does not imply automatic recovery or trading action

#### Scenario: Retrieval quality summary is visible
- **WHEN** API data includes retrieval quality summary, fallback reason, index freshness, or source consistency status
- **THEN** the frontend SHALL display the summary in an operator-readable way
- **AND** it SHALL NOT expose automatic rule application, automatic trading, broker order behavior, or hidden local file paths
