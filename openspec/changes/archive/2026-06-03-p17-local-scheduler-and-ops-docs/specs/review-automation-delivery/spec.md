## MODIFIED Requirements

### Requirement: Local agent task entrypoint
The system SHALL provide a local `cmd/agent` entrypoint for manually or locally scheduled daily discipline, market refresh, intelligence indexing, and review tasks without enabling automatic trading.

#### Scenario: Help command is available
- **WHEN** the user runs `go run ./cmd/agent --help`
- **THEN** the system displays the supported local tasks, local scheduling boundary, and safety notes
- **THEN** it does not mutate portfolio state

#### Scenario: Manual or scheduled task execution is audited
- **WHEN** a local task is executed through `cmd/agent`, including from a local scheduler example
- **THEN** the system records an `audit_events` entry containing the input summary, execution status, and error code when applicable.

#### Scenario: Local scheduling is safe by default
- **WHEN** local scheduling configuration exists
- **THEN** the system MUST keep automatic trading disabled and MUST NOT expose any order placement behavior.
- **THEN** scheduled examples MUST NOT bypass user confirmation or gatekeeper audit requirements.

#### Scenario: Failure keeps existing data consistent
- **WHEN** a local task fails
- **THEN** the system returns a readable error and preserves previously committed data consistency.

### Requirement: Delivery documentation is complete and safe
The system SHALL document local startup, initialization, data backup, index rebuild, recovery, scheduler setup/removal, common fault handling, and P7-P17 acceptance commands without including real secrets or personal sensitive information.

#### Scenario: Delivery docs cover local operation
- **WHEN** the user follows the delivery documentation
- **THEN** the documentation describes local startup, initialization, backup, VecLite index rebuild, scheduler setup/removal, and recovery procedures.

#### Scenario: Fault handling is documented
- **WHEN** a common fault occurs for data sources, VecLite, DeepSeek, SQLite writes, or local scheduled task execution
- **THEN** the documentation provides a handling path and keeps the no-automatic-trading boundary clear.

#### Scenario: Documentation contains no secrets
- **WHEN** delivery documentation is reviewed
- **THEN** it MUST NOT contain real API keys, account identifiers, tokens, or personal sensitive information.
