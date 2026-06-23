## ADDED Requirements

### Requirement: Execute project acceptance gates before release readiness claims

The project SHALL execute the project acceptance gates before making release readiness claims.

#### Scenario: Release candidate materials are prepared

- **GIVEN** the P52 project acceptance gate matrix exists
- **WHEN** P53 release candidate materials are prepared
- **THEN** the materials SHALL include actual G0-G9 gate results or explicit waivers
- **AND** the materials SHALL NOT treat the existence of the P52 matrix as evidence that acceptance passed.

### Requirement: Record blocked, degraded, and skipped gates explicitly

The project SHALL preserve acceptance failures and degraded outcomes in release materials.

#### Scenario: A gate does not pass

- **WHEN** an acceptance gate is blocked, degraded, or skipped
- **THEN** the acceptance record SHALL include the command or workflow entry, artifact, failure or skip reason, classification when applicable, and release impact
- **AND** release candidate materials SHALL state `release_blocked` when any non-waived release-blocking gate is blocked.

### Requirement: Keep acceptance artifacts redacted

The project SHALL keep acceptance and release materials free of sensitive runtime details.

#### Scenario: Acceptance materials are written

- **WHEN** P53 writes acceptance or release candidate materials
- **THEN** the materials SHALL NOT include complete API keys, private paths, raw HTTP responses, complete prompts, raw SQL dumps, or unredacted vendor payloads
- **AND** real-source or real-LLM results SHALL be summarized with redacted artifacts and explicit failure categories.
