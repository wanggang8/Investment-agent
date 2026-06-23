## ADDED Requirements

### Requirement: P43 data quality observability frontend tests

The frontend SHALL test the P43 data quality observability surface across successful, empty, degraded, error, unknown, sanitized, and narrow-viewport paths.

#### Scenario: Data quality component states are tested

- **WHEN** frontend unit tests run
- **THEN** they SHALL cover data quality panels for successful DTOs, empty local facts, source_unavailable, parse_error, stale, missing, unknown, LLM/RAG/VecLite degraded states, API errors, and safe Chinese status text
- **AND** tests SHALL assert that secrets, full prompts, private local paths, SQL errors, automatic trading, one-click order placement, external push, automatic confirmation, and automatic rule application copy is absent.

#### Scenario: Data quality browser smoke is tested

- **WHEN** Playwright smoke runs
- **THEN** it SHALL open the data quality route, verify primary panels and navigation entrypoints, check narrow viewport usability, and scan for forbidden automatic-action or sensitive diagnostic copy
- **AND** it SHALL use fixed local fixture data rather than private persistent data.
