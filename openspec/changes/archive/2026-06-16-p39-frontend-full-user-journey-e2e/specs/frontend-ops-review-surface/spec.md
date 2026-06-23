## ADDED Requirements

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
