## ADDED Requirements

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
