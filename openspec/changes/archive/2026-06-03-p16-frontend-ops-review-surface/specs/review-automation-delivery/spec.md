## MODIFIED Requirements

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

#### Scenario: Rule proposal remains user confirmed
- **WHEN** a review-generated rule proposal is ready for adoption
- **THEN** the system keeps the existing gatekeeper audit and final user confirmation requirements.

#### Scenario: Frontend tracking does not apply changes
- **WHEN** the frontend displays review-generated suggestions or tracking entrypoints
- **THEN** the entrypoints SHALL NOT apply rule changes, mutate portfolio data, or trigger trading behavior.
