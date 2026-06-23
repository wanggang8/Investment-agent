# local-ops-scheduler Specification

## Purpose
Document safe local scheduler examples and operator runbook requirements for local Investment Agent tasks.
## Requirements
### Requirement: Safe local scheduler examples
The project SHALL provide local scheduler examples for running supported `cmd/agent` tasks while keeping scheduled automation disabled by default.

#### Scenario: Scheduler examples are inert and local
- **WHEN** an operator reviews scheduler examples
- **THEN** examples SHALL require explicit local installation or editing before use
- **THEN** examples SHALL use placeholder paths and SHALL NOT contain secrets, account identifiers, or personal data

#### Scenario: Scheduler examples cannot trade
- **WHEN** scheduler examples invoke `cmd/agent`
- **THEN** they SHALL NOT include automatic trading, broker order, one-click order, automatic portfolio mutation, automatic confirmation, or automatic rule application behavior
- **THEN** they SHALL preserve existing user confirmation and gatekeeper audit requirements

### Requirement: Local operations runbook
The project SHALL document local operations for startup, backup, index recovery, scheduler setup, and fault handling.

#### Scenario: Operator can recover local index state
- **WHEN** VecLite index data is missing, stale, or incompatible
- **THEN** the runbook SHALL describe how to rebuild or verify the index using local task entrypoints without changing portfolio state

#### Scenario: Fault handling remains safe
- **WHEN** data sources, DeepSeek, VecLite, SQLite, or scheduled task execution fails
- **THEN** the runbook SHALL describe safe handling and audit expectations
- **THEN** the runbook SHALL keep automatic trading and automatic rule application prohibited

### Requirement: P39 Local E2E Runbook And Fixture Boundary
The project SHALL document and preserve a fixed-ID local browser E2E fixture and runbook for P39 without introducing real secrets, public network dependency, or persistent local data pollution.

#### Scenario: E2E fixture is local and fixed-ID
- **WHEN** the P39 Playwright or browser smoke fixture starts
- **THEN** it SHALL use temporary local SQLite/config paths and fixed-ID seed data
- **AND** it SHALL NOT include real API keys, broker credentials, personal account identifiers, or public network dependencies
- **AND** generated test artifacts SHALL be ignored or cleaned up according to repository gitignore boundaries

#### Scenario: E2E runbook describes startup and cleanup
- **WHEN** local operation documentation is reviewed
- **THEN** it SHALL describe the P39 E2E command, expected local services, port handling, failure cleanup, and how to avoid using private persistent data
- **AND** it SHALL keep automatic trading, broker integration, external push, and automatic rule application prohibited

### Requirement: P40 Local Deploy Preflight And Recovery Drill
The project SHALL provide a local deploy, operations, and recovery drill that verifies runtime prerequisites, safe startup diagnostics, backup recovery behavior, and generated artifact boundaries without introducing trading or external push capabilities.

#### Scenario: Local preflight reports actionable runtime status
- **WHEN** an operator runs the P40 local preflight or initialization command
- **THEN** it SHALL check Go, Node, npm, Playwright browser availability, SQLite path, VecLite path, config file presence, data directories, and write permissions
- **AND** it SHALL report pass, warning, failed, or skipped states with actionable local remediation hints
- **AND** it SHALL NOT print raw API keys, broker credentials, or personal account data

#### Scenario: Startup diagnostics remain local and safe
- **WHEN** local server or agent startup prerequisites are missing or invalid
- **THEN** diagnostics SHALL expose stable status or error categories for config, migration, SQLite, VecLite, data source, and LLM readiness
- **AND** diagnostics SHALL write only local logs, audit records, or diagnostic files
- **AND** diagnostics SHALL NOT trigger automatic trading, external push, automatic confirmation, or automatic rule application

#### Scenario: Backup recovery drill avoids unsafe overwrite
- **WHEN** a backup recovery smoke or drill is executed
- **THEN** it SHALL default to temporary local restore paths
- **AND** restoring into an existing database SHALL require explicit confirmation or fail safely
- **AND** restored facts SHALL be verified through API or browser-visible behavior rather than direct manual database inspection

#### Scenario: Generated operations artifacts are governed
- **WHEN** preflight, startup diagnostics, recovery smoke, or browser smoke creates files
- **THEN** the outputs SHALL be documented and either cleaned up or ignored by repository gitignore rules
- **AND** failure reproduction SHALL avoid using private persistent databases or real secrets
