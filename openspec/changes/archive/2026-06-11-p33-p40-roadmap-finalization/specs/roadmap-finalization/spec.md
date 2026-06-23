## ADDED Requirements

### Requirement: Maintain the P33-P40 roadmap as the remaining planned feature queue

The project SHALL treat P33 through P40 as the current remaining planned feature queue after P32 archive, unless a later OpenSpec governance change explicitly revises the roadmap.

#### Scenario: Next feature work is selected after P32 archive

- **GIVEN** P32 has been archived and no earlier active change remains
- **WHEN** the team selects the next feature implementation stage
- **THEN** the selected stage SHALL come from P33 through P40
- **AND** the selected stage SHALL have its own OpenSpec change before implementation
- **AND** the system SHALL NOT skip required governance by directly editing L1 contracts

#### Scenario: Work outside P33-P40 is requested

- **GIVEN** a requested enhancement is not covered by P33 through P40
- **WHEN** the team decides to pursue it
- **THEN** the project SHALL create a separate governance change or future roadmap change
- **AND** it SHALL NOT be silently folded into an unrelated P33-P40 stage

### Requirement: Preserve roadmap execution dependencies

The project SHALL preserve the documented execution dependencies between P33 through P40 so downstream stages have the required product, data, and verification foundations.

#### Scenario: User journey validation depends on account initialization

- **GIVEN** P39 requires a real user journey from empty database to the first daily discipline report
- **WHEN** P39 is planned
- **THEN** P33 account and holdings initialization capability SHALL be completed or explicitly scoped as a prerequisite

#### Scenario: Risk and evidence quality work depends on data coverage

- **GIVEN** P35 risk warnings and P38 retrieval quality depend on reliable data and evidence inputs
- **WHEN** P35 or P38 is planned
- **THEN** P34 data coverage gaps and freshness/failure semantics SHALL be considered as prerequisites or documented assumptions

### Requirement: Keep non-feature governance work separate from P33-P40 implementation

The project SHALL keep P19-P24 historical archive追溯 and P40-after roadmap expansion separate from P33-P40 feature implementation.

#### Scenario: Historical archive traceability is requested

- **GIVEN** P19 through P24 are marked delivered but do not have full archive packages
- **WHEN** audit-grade traceability is required
- **THEN** the project SHALL create a dedicated governance change
- **AND** it SHALL NOT rewrite or fabricate historical archive packages

#### Scenario: P40-after product work is requested

- **GIVEN** P40 is the last stage in the current roadmap
- **WHEN** new product work beyond P40 is requested
- **THEN** the project SHALL create a new roadmap or proposal change
- **AND** it SHALL define scope, dependencies, and acceptance before implementation
