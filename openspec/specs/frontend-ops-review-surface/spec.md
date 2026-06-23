# frontend-ops-review-surface Specification

## Purpose
Document operator-facing frontend surfaces for local data status, index health, review readiness, periodic review summaries, and safe tracking entrypoints.
## Requirements
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

### Requirement: P39 Cross Feature Operational Journey Surfaces
The frontend SHALL expose cross-page operational and review context required for the P39 full user journey using API/service DTOs only.

#### Scenario: Operational status is reachable during the journey
- **WHEN** the user moves from daily discipline or dashboard surfaces into evidence, decision, review, and risk alert pages
- **THEN** source health, index or retrieval quality, risk alert/SOP status, and review readiness SHALL be reachable without direct SQLite, VecLite, or local file reads
- **AND** degraded or missing status SHALL remain distinguishable from healthy status

#### Scenario: Tracking entrypoints preserve safe boundaries
- **WHEN** a review, decision, risk alert, audit event, confirmation record, or rule proposal is referenced in the journey
- **THEN** the frontend SHALL expose a visible path to inspect related facts or filtered records
- **AND** these entrypoints SHALL only navigate, filter, or record local user facts
- **AND** they SHALL NOT trigger automatic trading, automatic confirmation, external push, or automatic rule application

#### Scenario: Degraded states explain safe next steps
- **WHEN** account data, market data, evidence, VecLite/RAG retrieval, LLM output, capability scope, or rule proposal status is degraded or incomplete
- **THEN** the frontend SHALL show a safe explanation and an inspectable next step where one exists
- **AND** it SHALL NOT imply guaranteed recovery, investment return, price prediction certainty, or executable brokerage action

### Requirement: P40 Data Source Health And Runtime Readiness Surface
The frontend SHALL expose local runtime readiness and data source health states required by the P40 operations drill using API/service DTOs only.

#### Scenario: Data source health shows freshness and failures
- **WHEN** data source health facts are available
- **THEN** the frontend SHALL show last success time, last failure time, failure category, freshness, and affected symbols or scopes
- **AND** fresh, stale, failed, missing, and unknown states SHALL remain visually and textually distinguishable

#### Scenario: Runtime readiness is safe to inspect
- **WHEN** local SQLite, VecLite, data source, or LLM readiness is degraded or unavailable
- **THEN** the frontend SHALL show safe next steps for inspection or local repair
- **AND** it SHALL NOT imply automatic recovery, guaranteed investment return, external notification delivery, or executable brokerage action

#### Scenario: Operations drill entrypoints are non-executing
- **WHEN** a readiness or health panel links to logs, diagnostics, recovery smoke, or related facts
- **THEN** those entrypoints SHALL only navigate, filter, or show local diagnostic facts
- **AND** they SHALL NOT mutate portfolio facts, apply rules, send external pushes, or place orders

### Requirement: P42 workbench aggregates ops and review status safely

The frontend SHALL aggregate risk, rule, review, source health, and runtime readiness states on the user decision workbench without changing their underlying workflows.

#### Scenario: Workbench shows ops and review follow-up

- **WHEN** risk alerts, rule proposals, review summaries, source health, or runtime readiness facts are available
- **THEN** the workbench SHALL show a concise status, safe next step, and navigation path to the authoritative page
- **AND** it SHALL NOT imply automatic repair, automatic rule application, external notification delivery, or trading execution

### Requirement: P43 data quality surface aggregates operational quality safely

The frontend SHALL expose a P43 data quality surface that summarizes source health, evidence freshness, retrieval quality, LLM quality, and local diagnostic readiness without changing their underlying workflows.

#### Scenario: Data quality status is visible

- **WHEN** data quality facts are available from existing frontend services or read-only aggregation DTOs
- **THEN** the frontend SHALL show source health, evidence freshness, retrieval/index freshness, fallback source, LLM parse/quality status, and affected workflow scope where available
- **AND** it SHALL distinguish success, degraded, failed, missing, stale, parse_error, source_unavailable, and unknown states with text, not color alone
- **AND** it SHALL NOT expose automatic repair, automatic rule application, external notification delivery, or trading execution.

#### Scenario: Data quality navigation is safe

- **WHEN** the user selects a data quality link
- **THEN** the frontend SHALL navigate to settings, evidence, review, audit, risk alert, decision, or workbench pages
- **AND** it SHALL NOT trigger refresh, rebuild, smoke tests, confirmation submission, rule application, external push, or account mutation.

