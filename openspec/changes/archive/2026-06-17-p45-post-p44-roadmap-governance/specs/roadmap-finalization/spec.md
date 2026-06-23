## ADDED Requirements

### Requirement: Govern all post-P44 work through a refreshed roadmap change

The project SHALL treat P44 as the end of the P41-selected post-P40 queue and SHALL require a refreshed roadmap or governance OpenSpec change before any post-P44 implementation begins.

#### Scenario: Post-P44 work is requested

- **GIVEN** P42, P43, P44, and P19-P24 historical traceability have been archived
- **AND** no active implementation change exists
- **WHEN** the team chooses a new direction
- **THEN** the project SHALL first create a roadmap or governance OpenSpec change
- **AND** the roadmap SHALL define candidate stages, dependencies, out-of-scope boundaries, and acceptance criteria before implementation

#### Scenario: A concrete post-P44 feature is selected

- **GIVEN** a post-P44 roadmap or governance change has identified candidate work
- **WHEN** the team decides to implement one candidate feature
- **THEN** that feature SHALL receive its own OpenSpec change with `proposal.md`, `design.md`, `tasks.md`, and relevant delta files
- **AND** implementation SHALL follow the selected change instead of editing contracts directly

### Requirement: Prioritize post-P44 candidates by local safety and verifiability

The project SHALL prefer post-P44 candidates that improve local-only safety, explainability, data quality, and operability without expanding trading or external delivery capabilities.

#### Scenario: Post-P44 candidates are evaluated

- **WHEN** roadmap candidates are added or prioritized
- **THEN** the roadmap SHALL classify them by category, suggested change id, dependencies, acceptance approach, and suitability
- **AND** the roadmap SHALL recommend a next candidate while keeping implementation out of the roadmap change

### Requirement: Preserve safety boundaries in post-P44 roadmap planning

The project SHALL preserve the existing investment safety boundaries unless a dedicated future governance change explicitly reopens them with review.

#### Scenario: Post-P44 roadmap candidates are evaluated

- **WHEN** roadmap candidates are added or prioritized
- **THEN** the roadmap SHALL keep broker APIs, automatic trading, one-click trading, external push, automatic confirmation, automatic rule application, automatic repair promises, paid, login-only, or authorization-gated data sources, Level2 data, high-frequency sources, return promises, and deterministic price predictions out of scope by default
- **AND** LLM usage SHALL remain limited to analysis material rather than final rule verdicts
