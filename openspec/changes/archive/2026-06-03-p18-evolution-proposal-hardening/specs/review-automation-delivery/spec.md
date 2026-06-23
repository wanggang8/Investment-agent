## MODIFIED Requirements

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
