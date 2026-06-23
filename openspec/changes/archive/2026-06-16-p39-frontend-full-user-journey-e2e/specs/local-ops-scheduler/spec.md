## ADDED Requirements

### Requirement: P39 Local E2E Runbook And Fixture Boundary
The project SHALL document and preserve a deterministic local browser E2E fixture and runbook for P39 without introducing real secrets, public network dependency, or persistent local data pollution.

#### Scenario: E2E fixture is local and deterministic
- **WHEN** the P39 Playwright or browser smoke fixture starts
- **THEN** it SHALL use temporary local SQLite/config paths and deterministic seed data
- **AND** it SHALL NOT include real API keys, broker credentials, personal account identifiers, or public network dependencies
- **AND** generated test artifacts SHALL be ignored or cleaned up according to repository gitignore boundaries

#### Scenario: E2E runbook describes startup and cleanup
- **WHEN** local operation documentation is reviewed
- **THEN** it SHALL describe the P39 E2E command, expected local services, port handling, failure cleanup, and how to avoid using private persistent data
- **AND** it SHALL keep automatic trading, broker integration, external push, and automatic rule application prohibited
