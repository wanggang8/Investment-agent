# roadmap-finalization Specification Delta

## ADDED Requirements

### Requirement: Prioritize P19-P24 audit evidence before release materials

The project SHALL prioritize a dedicated P19-P24 audit evidence phase before preparing release candidate materials after P49.

#### Scenario: Post-P49 work is selected

- **GIVEN** P49 has been archived
- **AND** P19-P24 are marked delivered but do not have per-phase complete archive packages
- **WHEN** the team selects the next governance or product stage
- **THEN** the next recommended stage SHALL be P51 `p51-p19-p24-audit-evidence-pack`
- **AND** that stage SHALL produce current-state audit evidence instead of fabricating historical archive packages
- **AND** release candidate materials SHALL remain deferred until the P19-P24 audit evidence and project acceptance gate matrix have been completed.

### Requirement: Define project acceptance gates before release candidate work

The project SHALL define a layered acceptance gate matrix before release candidate materials are prepared.

#### Scenario: Project acceptance is planned

- **WHEN** acceptance planning is performed after the P19-P24 audit evidence phase
- **THEN** the plan SHALL cover unit tests, integration tests, E2E tests, real-source tests, real-LLM tests, local smoke tests, install diagnostics, release-upgrade checks, and safety boundary checks
- **AND** each gate SHALL define command or workflow entry, prerequisites, pass criteria, allowed degradation, artifacts, and whether failure blocks release
- **AND** real-source and real-LLM tests SHALL require explicit opt-in configuration and classified failure reporting.

### Requirement: Keep release candidate materials dependent on audit and acceptance readiness

The project SHALL treat release candidate material preparation as dependent on both historical audit evidence and acceptance gate readiness.

#### Scenario: Release candidate materials are requested

- **GIVEN** P51 and P52 are not both completed
- **WHEN** release candidate materials or release packaging are requested
- **THEN** the project SHALL first complete or explicitly waive the missing audit or acceptance phase in a governance change
- **AND** the project SHALL NOT claim release readiness from P49 alone.
