## ADDED Requirements

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
