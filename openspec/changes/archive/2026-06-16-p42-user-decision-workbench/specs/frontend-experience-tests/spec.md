## ADDED Requirements

### Requirement: P42 workbench frontend tests

The frontend SHALL test the P42 user decision workbench across successful, empty, degraded, error, and narrow-viewport paths.

#### Scenario: Workbench component states are tested

- **WHEN** frontend unit tests run
- **THEN** they SHALL cover workbench panels for successful DTOs, empty local facts, degraded source/LLM/RAG status, API errors, and safe Chinese status text
- **AND** tests SHALL assert that automatic trading, one-click order placement, external push, automatic confirmation, and automatic rule application copy is absent

#### Scenario: Workbench browser smoke is tested

- **WHEN** Playwright smoke runs
- **THEN** it SHALL open the workbench route, verify primary panels and navigation entrypoints, check narrow viewport usability, and scan for forbidden automatic-action copy
- **AND** it SHALL use fixed local fixture data rather than private persistent data
