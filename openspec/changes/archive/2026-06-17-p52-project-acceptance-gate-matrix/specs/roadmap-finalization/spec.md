# roadmap-finalization Specification Delta

## ADDED Requirements

### Requirement: Maintain a project acceptance gate matrix before release materials

The project SHALL provide a project acceptance gate matrix before release candidate materials are prepared.

#### Scenario: Acceptance readiness is reviewed

- **GIVEN** P51 has provided the P19-P24 audit evidence pack
- **WHEN** project acceptance readiness is reviewed
- **THEN** the project SHALL provide a matrix covering governance checks, Go tests, integration tests, frontend tests and build, E2E smoke, local smoke, real-source opt-in tests, real-LLM opt-in tests, install diagnostics, release-upgrade checks, and safety boundary checks
- **AND** each gate SHALL define command or workflow entry, prerequisites, pass criteria, allowed degradation, artifacts, and whether failure blocks release.

### Requirement: Classify real-source and real-LLM acceptance failures

The project SHALL classify real-source and real-LLM acceptance failures instead of treating every failure as the same release outcome.

#### Scenario: Opt-in real acceptance test fails

- **WHEN** a real-source or real-LLM acceptance gate fails
- **THEN** the result SHALL classify the failure as network, rate limit, authentication or key, source schema change, no data, parse failure, model unavailable, quality failure, or another explicit category
- **AND** the result SHALL state whether the failure blocks release, blocks only a real-source or LLM claim, or is accepted as a documented degradation.

### Requirement: Keep release readiness separate from matrix definition

The project SHALL NOT claim release readiness merely because the acceptance gate matrix exists.

#### Scenario: Release candidate materials are prepared after P52

- **GIVEN** the acceptance gate matrix has been defined
- **WHEN** P53 release candidate materials are prepared
- **THEN** P53 SHALL reference actual gate results or explicitly documented waivers
- **AND** P52 alone SHALL NOT be treated as proof that all gates have passed.
