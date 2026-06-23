## ADDED Requirements

### Requirement: Provide release handoff materials after acceptance execution

The project SHALL provide release handoff materials after a release-ready acceptance execution.

#### Scenario: A release candidate is handed off

- **GIVEN** P53 has produced a `release_ready` candidate with actual acceptance results
- **WHEN** P54 release handoff materials are written
- **THEN** the materials SHALL cite the P53 acceptance run and release candidate documents
- **AND** the materials SHALL preserve documented degradations, retries, and non-claims instead of broadening the release-ready statement.

### Requirement: Define acceptance repeatability rules

The project SHALL document how release acceptance can be repeated consistently.

#### Scenario: Acceptance is repeated after P53

- **WHEN** an operator reruns acceptance after P53
- **THEN** the repeatability guidance SHALL define the output directory convention, command order, allowed retry conditions, current-data degraded handling, real-source configuration prerequisites, real-LLM redaction rules, and release status decision rules
- **AND** a second failure after an allowed retry SHALL be treated as blocked unless an explicit waiver is recorded.

### Requirement: Keep release handoff within safety boundaries

Release handoff materials SHALL NOT expand the product safety boundary.

#### Scenario: Handoff materials are reviewed

- **WHEN** release handoff materials describe release readiness
- **THEN** they SHALL NOT claim investment returns, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, future provider availability, login sources, paid sources, authorized sources, Level2 data, or high-frequency data
- **AND** they SHALL NOT include complete API keys, private paths, raw HTTP responses, complete prompts, raw SQL dumps, or unredacted vendor payloads.
