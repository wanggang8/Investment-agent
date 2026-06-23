## ADDED Requirements

### Requirement: P49 SHALL provide a local release and upgrade readiness report

The system SHALL provide a local release and upgrade readiness report that helps an operator review version intent, backup readiness, migration precheck status, and post-upgrade smoke commands before changing a local installation.

#### Scenario: Upgrade readiness check is requested
- **WHEN** a user runs the P49 local release/upgrade check
- **THEN** the system SHALL report current version, target version, backup reminder, migration file precheck status, post-upgrade smoke commands, and safety note
- **AND** it SHALL make clear that the check did not perform an upgrade, did not run migrations, did not create a backup, and did not restore or overwrite a database.

#### Scenario: Target version is missing
- **WHEN** a user runs the release/upgrade check without a target version
- **THEN** the report SHALL return a warning status that asks the user to provide a target version or release label
- **AND** it SHALL still provide backup and smoke checklist guidance without modifying local data.

### Requirement: P49 SHALL keep release and upgrade diagnostics local, read-only, and sanitized

P49 SHALL expose release and upgrade diagnostics through local CLI and script summaries without leaking sensitive local details or creating side effects.

#### Scenario: Diagnostics JSON is written
- **WHEN** the release/upgrade check writes a diagnostics JSON file
- **THEN** the file SHALL contain only sanitized status, commands, counts, basenames or placeholders, and safety notes
- **AND** it SHALL NOT contain complete API keys, private absolute paths, raw SQL, full prompts, raw HTTP exchanges, private keys, supplier raw responses, account secrets, or broker credentials.

#### Scenario: Local install diagnostics includes release upgrade check
- **WHEN** the local install diagnostics script is explicitly run with release/upgrade inclusion enabled
- **THEN** it SHALL include the release/upgrade check as an additional local step
- **AND** the default local install diagnostics behavior SHALL remain unchanged unless the user explicitly enables that step.

### Requirement: P49 SHALL preserve upgrade safety boundaries

P49 SHALL remain a local upgrade planning and verification aid and SHALL NOT expand automation, trading, data acquisition, or repair boundaries.

#### Scenario: Upgrade readiness detects missing backup or migration uncertainty
- **WHEN** the readiness report detects a missing SQLite file, missing target version, missing migration files, or another warning condition
- **THEN** it SHALL classify the report as warning or blocked according to the check severity
- **AND** it SHALL NOT automatically create backups, run migrations, repair files, restore data, mark the installation healthy, execute trades, apply rules, confirm operations, call LLMs, access public endpoints, or promise returns.

#### Scenario: Post-upgrade smoke plan is shown
- **WHEN** the report lists post-upgrade smoke commands
- **THEN** those commands SHALL remain manual commands for the user to run explicitly
- **AND** they SHALL NOT include broker APIs, external push channels, one-click trading, delegated order placement, login-gated data sources, paid data sources, authorization-gated sources, Level2 feeds, or high-frequency polling.
