## ADDED Requirements

### Requirement: Govern all post-P40 product work through a new roadmap change

The project SHALL treat P40 as the end of the completed P33-P40 planned feature queue and SHALL require a new roadmap or governance change before any post-P40 product implementation begins.

#### Scenario: Post-P40 product work is requested

- **GIVEN** P40 has been archived and no active implementation change exists
- **WHEN** the team chooses a new product direction
- **THEN** the project SHALL first create a roadmap or governance OpenSpec change
- **AND** the roadmap SHALL define candidate stages, dependencies, out-of-scope boundaries, and acceptance criteria before implementation

#### Scenario: A concrete post-P40 feature is selected

- **GIVEN** a post-P40 roadmap or governance change has identified candidate work
- **WHEN** the team decides to implement one candidate feature
- **THEN** that feature SHALL receive its own OpenSpec change with `proposal.md`, `tasks.md`, and relevant delta files
- **AND** implementation SHALL follow the selected change instead of editing contracts directly

### Requirement: Separate historical audit traceability from new product features

The project SHALL keep P19-P24 historical archive traceability separate from post-P40 product feature implementation.

#### Scenario: P19-P24 archive traceability is requested after P40

- **GIVEN** P19 through P24 are delivered but do not have full archive packages
- **WHEN** audit-grade traceability is requested
- **THEN** the project SHALL create a dedicated governance change for historical traceability
- **AND** it SHALL NOT rewrite or fabricate historical archive packages
- **AND** it SHALL NOT block unrelated post-P40 product roadmap planning unless explicitly prioritized

### Requirement: Preserve safety boundaries in post-P40 roadmap planning

The project SHALL preserve the existing investment safety boundaries unless a dedicated future governance change explicitly reopens them with review.

#### Scenario: Post-P40 roadmap candidates are evaluated

- **WHEN** roadmap candidates are added or prioritized
- **THEN** the roadmap SHALL keep broker APIs, automatic trading, external push, automatic rule application, paid or login-only data sources, Level2 data, high-frequency sources, return promises, and deterministic price predictions out of scope by default
- **AND** LLM usage SHALL remain limited to analysis material rather than final rule verdicts
