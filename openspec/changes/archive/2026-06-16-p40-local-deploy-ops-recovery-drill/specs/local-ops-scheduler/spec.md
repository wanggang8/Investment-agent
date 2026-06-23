## ADDED Requirements

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
