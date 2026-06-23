## ADDED Requirements

### Requirement: Local agent task entrypoint
The system SHALL provide a local `cmd/agent` entrypoint for manually triggering daily discipline, market refresh, intelligence indexing, and review tasks without enabling automatic trading.

#### Scenario: Help command is available
- **WHEN** the user runs `go run ./cmd/agent --help`
- **THEN** the system displays the supported local tasks and does not mutate portfolio state.

#### Scenario: Manual task execution is audited
- **WHEN** a local task is executed through `cmd/agent`
- **THEN** the system records an `audit_events` entry containing the input summary, execution status, and error code when applicable.

#### Scenario: Local scheduling is safe by default
- **WHEN** local scheduling configuration exists
- **THEN** the system MUST keep automatic trading disabled and MUST NOT expose any order placement behavior.

#### Scenario: Failure keeps existing data consistent
- **WHEN** a local task fails
- **THEN** the system returns a readable error and preserves previously committed data consistency.

### Requirement: Periodic review summaries
The system SHALL generate monthly and quarterly review summaries from existing local facts, including confirmations, error cases, rule proposals, audit events, rule hits, misjudgments, missing evidence, and degradation cases.

#### Scenario: Monthly review aggregates local facts
- **WHEN** a monthly review task is triggered
- **THEN** the system summarizes confirmation actions, error cases, rule proposals, and audit events for the selected period.

#### Scenario: Quarterly review evaluates rule effectiveness
- **WHEN** a quarterly review task is triggered
- **THEN** the system summarizes rule hits, misjudgments, missing evidence, and degradation cases for the selected period.

### Requirement: Rule effectiveness evaluation remains gated
The system SHALL write rule effectiveness evaluation output only as a review summary or rule proposal, and any rule change MUST still pass gatekeeper audit and user final confirmation.

#### Scenario: Evaluation creates review output
- **WHEN** rule effectiveness evaluation identifies a potential improvement
- **THEN** the system records the result in a review summary or rule proposal without applying a rule version automatically.

#### Scenario: Rule proposal remains user confirmed
- **WHEN** a review-generated rule proposal is ready for adoption
- **THEN** the system keeps the existing gatekeeper audit and final user confirmation requirements.

### Requirement: Delivery documentation is complete and safe
The system SHALL document local startup, initialization, data backup, index rebuild, recovery, common fault handling, and P7-P9 acceptance commands without including real secrets or personal sensitive information.

#### Scenario: Delivery docs cover local operation
- **WHEN** the user follows the delivery documentation
- **THEN** the documentation describes local startup, initialization, backup, VecLite index rebuild, and recovery procedures.

#### Scenario: Fault handling is documented
- **WHEN** a common fault occurs for data sources, VecLite, DeepSeek, or SQLite writes
- **THEN** the documentation provides a handling path and keeps the no-automatic-trading boundary clear.

#### Scenario: Documentation contains no secrets
- **WHEN** delivery documentation is reviewed
- **THEN** it MUST NOT contain real API keys, account identifiers, tokens, or personal sensitive information.
