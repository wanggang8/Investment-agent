# review-automation-delivery Specification

## Purpose
This specification defines the local review automation and delivery capabilities added in P9: manual `cmd/agent` task entrypoints, monthly and quarterly review summaries, gated rule-effectiveness evaluation, and safe local operation documentation.
## Requirements
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

### Requirement: Periodic review summaries
The system SHALL generate and expose monthly and quarterly review summaries from existing local facts, including confirmations, error cases, rule proposals, audit events, rule hits, misjudgments, missing evidence, degradation cases, and frontend tracking metadata.

#### Scenario: Monthly review aggregates local facts
- **WHEN** a monthly review task is triggered
- **THEN** the system summarizes confirmation actions, error cases, rule proposals, and audit events for the selected period.

#### Scenario: Quarterly review evaluates rule effectiveness
- **WHEN** a quarterly review task is triggered
- **THEN** the system summarizes rule hits, misjudgments, missing evidence, and degradation cases for the selected period.

#### Scenario: Review summary supports tracking display
- **WHEN** review output references audit events, decisions, rule proposals, or error cases
- **THEN** the system SHALL expose enough DTO metadata for the frontend to show safe tracking entrypoints.

### Requirement: Rule effectiveness evaluation remains gated
The system SHALL write rule effectiveness evaluation output only as a review summary or rule proposal, and any rule change MUST still pass gatekeeper audit and user final confirmation.

#### Scenario: Evaluation creates review output
- **WHEN** rule effectiveness evaluation identifies a potential improvement
- **THEN** the system records the result in a review summary or rule proposal without applying a rule version automatically.

#### Scenario: Review output creates traceable proposal
- **WHEN** a review-generated suggestion is promoted into a rule proposal
- **THEN** the proposal SHALL retain source metadata for related review period, audit events, decisions, or error cases where available.

#### Scenario: Rule proposal remains user confirmed
- **WHEN** a review-generated rule proposal is ready for adoption
- **THEN** the system keeps the existing gatekeeper audit and final user confirmation requirements.

#### Scenario: Frontend tracking does not apply changes
- **WHEN** the frontend displays review-generated suggestions or tracking entrypoints
- **THEN** the entrypoints SHALL NOT apply rule changes, mutate portfolio data, or trigger trading behavior.

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

### Requirement: Review summaries SHALL include rule effect trends
The system SHALL include rule effect validation and applied rule tracking summaries in monthly or quarterly review outputs when relevant local facts exist.

#### Scenario: Review contains applied rule tracking
- **WHEN** a review period includes decisions, confirmations, errors, degradation, or risk alerts related to an applied rule version
- **THEN** the review output SHALL include rule hit count, misjudgment count, missing evidence count, degradation count, risk alert count, trend direction, and links to related proposals and audits.

#### Scenario: Review lacks enough facts
- **WHEN** the review period does not contain enough local facts to evaluate a rule
- **THEN** the review output SHALL show insufficient tracking data
- **AND** it SHALL NOT claim that the applied rule improved or degraded results.

### Requirement: Review-generated follow-up suggestions SHALL remain gated
The system SHALL treat negative rule tracking findings as review information or draft follow-up proposals, not automatic rule rollback or automatic rule replacement.

#### Scenario: Review finds deterioration after rule application
- **WHEN** review tracking finds worse metrics after a rule version was applied
- **THEN** the system SHALL expose the deterioration as a warning or draft follow-up suggestion
- **AND** any rule change SHALL still require a proposal, gatekeeper audit, and user final confirmation.

#### Scenario: Frontend opens review tracking links
- **WHEN** the frontend displays tracking links from review output
- **THEN** those links SHALL navigate to local proposal, audit, risk alert, or decision views
- **AND** they SHALL NOT execute rule changes, trading actions, or external notifications.

