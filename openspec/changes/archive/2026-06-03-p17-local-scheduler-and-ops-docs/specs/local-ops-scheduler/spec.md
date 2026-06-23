## ADDED Requirements

### Requirement: Safe local scheduler examples
The project SHALL provide local scheduler examples for running supported `cmd/agent` tasks while keeping scheduled automation disabled by default.

#### Scenario: Scheduler examples are inert and local
- **WHEN** an operator reviews scheduler examples
- **THEN** examples SHALL require explicit local installation or editing before use
- **THEN** examples SHALL use placeholder paths and SHALL NOT contain secrets, account identifiers, or personal data

#### Scenario: Scheduler examples cannot trade
- **WHEN** scheduler examples invoke `cmd/agent`
- **THEN** they SHALL NOT include automatic trading, broker order, one-click order, automatic portfolio mutation, or automatic rule application behavior
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
