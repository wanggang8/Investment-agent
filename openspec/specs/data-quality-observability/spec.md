# data-quality-observability Specification

## Purpose
Define the read-only frontend surface that helps users inspect data source health, evidence/retrieval quality, LLM quality, and affected workflows without exposing sensitive diagnostics or triggering automated actions.
## Requirements
### Requirement: Provide a data quality observability surface

The frontend SHALL provide a data quality observability surface that aggregates existing source health, evidence freshness, RAG/VecLite retrieval quality, LLM quality, and local diagnostic readiness into one read-only inspection surface.

#### Scenario: User opens data quality observability

- **WHEN** source health, evidence, retrieval, LLM, review, or diagnostic DTOs are available
- **THEN** the surface SHALL show data source health, evidence and retrieval quality, LLM quality status, and affected workflow scope
- **AND** it SHALL provide navigation to authoritative pages such as settings, evidence, review, audit, risk alerts, decisions, or the workbench
- **AND** it SHALL NOT trigger refresh, repair, LLM smoke, rule application, external push, or trading execution.

#### Scenario: Quality states are incomplete or degraded

- **WHEN** source health, market data, evidence, RAG/VecLite, LLM, review, or diagnostics status is `source_unavailable`, `parse_error`, `stale`, `missing`, `unknown`, `degraded`, `quality_failed`, or unavailable
- **THEN** the surface SHALL show the degraded or unknown status explicitly and a safe inspection next step where available
- **AND** it SHALL NOT display the status as successful, imply guaranteed recovery, or generate investment conclusions from incomplete quality data.

### Requirement: Data quality observability remains sanitized

The data quality observability surface SHALL expose stable summaries and safe diagnostics without revealing secrets or sensitive local implementation details.

#### Scenario: Quality diagnostics contain sensitive details

- **WHEN** diagnostics, source errors, LLM failures, retrieval metadata, or local readiness records contain raw details
- **THEN** the surface SHALL display only stable status, category, timestamps, affected scope, and safe summaries
- **AND** it SHALL NOT display full API keys, full prompts, private local paths, SQL errors, vendor raw responses, broker/account secrets, or unnecessary account detail.

### Requirement: Data quality observability uses supported DTOs only

The data quality observability surface SHALL use supported API/service DTOs and frontend mappers rather than direct local storage or filesystem access.

#### Scenario: Data quality panels render

- **WHEN** the surface renders source, evidence, retrieval, LLM, review, audit, settings, or diagnostic information
- **THEN** the data SHALL come from API/service DTOs or a dedicated read-only aggregation DTO
- **AND** the frontend SHALL NOT read SQLite, VecLite, local logs, diagnostic files, private config files, browser localStorage, or browser sessionStorage directly.
