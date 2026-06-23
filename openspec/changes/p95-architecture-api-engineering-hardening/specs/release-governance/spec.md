## ADDED Requirements

### Requirement: P95 SHALL harden public engineering validation

P95 SHALL make public repository validation stable across clean checkouts and local developer checkouts with frontend dependencies installed.

#### Scenario: Backend package discovery excludes frontend dependencies

- **GIVEN** frontend dependencies have been installed under `web/node_modules`
- **WHEN** backend validation selects Go packages for tests
- **THEN** packages below `web/node_modules` SHALL NOT be included
- **AND** the selection helper SHALL fail if a package from frontend dependency trees is selected.

#### Scenario: P93 source scan ignores local runtime artifacts

- **GIVEN** ignored local runtime artifacts exist under project paths such as `cmd/agent/tmp/`
- **AND** nonignored new source files may exist before they are committed
- **WHEN** P93 code reality audit runs in check mode
- **THEN** ignored local runtime artifacts SHALL NOT change the report
- **AND** tracked plus nonignored untracked release-relevant source files SHALL be eligible for scanning
- **AND** tracked release-relevant files SHALL still be scanned for secrets, demo/stub risks, and release-boundary violations.

#### Scenario: API route contract is checked

- **GIVEN** backend handlers register `/api/v1` routes
- **WHEN** the API route contract check runs
- **THEN** every registered route SHALL be documented in `docs/api.md` or `docs/frontend-contract.md`
- **AND** documented route examples with query strings SHALL normalize to their path identity
- **AND** the check SHALL fail when docs reference a route that is no longer registered.

#### Scenario: Local SQLite runtime is concurrency-aware

- **GIVEN** the local server opens a SQLite database
- **WHEN** the database is file-backed
- **THEN** the connection SHALL enable foreign key enforcement and a bounded busy timeout
- **AND** it SHALL attempt WAL mode for local UI/background-task read-write concurrency
- **AND** in-memory tests SHALL remain supported.

#### Scenario: Docker deployment supports file-based LLM secrets

- **GIVEN** an operator supplies a `DEEPSEEK_API_KEY_FILE`
- **WHEN** the application loads runtime configuration
- **THEN** the key SHALL be read from that file if `DEEPSEEK_API_KEY` is not set
- **AND** committed configuration and documentation SHALL NOT contain real keys.
