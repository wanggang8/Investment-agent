## ADDED Requirements

### Requirement: Evolution proposals SHALL consume effect validation
The system SHALL attach rule effect validation output to rule proposals before guardrail decisions or final application decisions are made.

#### Scenario: Proposal has validation context
- **WHEN** a rule proposal is listed or opened after validation has run
- **THEN** the proposal SHALL expose validation status, sample count, overfit risk, replay result, guardrail decision, and validation link
- **AND** missing validation SHALL be shown as not_evaluated rather than assumed safe.

#### Scenario: Guardrail blocks weak validation
- **WHEN** effect validation is insufficient, overfit risk is high, or historical replay is unfavorable
- **THEN** the proposal SHALL remain draft, pending_user_confirm, rejected, or needs_user_review according to the existing state machine
- **AND** it SHALL NOT become an active rule version.

#### Scenario: Validation passes but final confirmation is still required
- **WHEN** effect validation passes and gatekeeper audit approves the proposal
- **THEN** the proposal MAY move to pending_final_confirm
- **AND** it SHALL still require explicit user final confirmation before any rule version is created.

### Requirement: Rule proposal source explanation SHALL be complete
The system SHALL explain why a rule proposal exists and which local facts support or weaken it.

#### Scenario: Proposal source is inspected
- **WHEN** a user inspects a rule proposal
- **THEN** the system SHALL expose source error cases, review period, related decisions, confirmations, risk alerts, audit events, and impacted workflows where available
- **AND** it SHALL identify missing source facts rather than hiding them.

#### Scenario: Source facts are narrow
- **WHEN** a proposal is derived from too few or too narrow source cases
- **THEN** the proposal SHALL include a sample representativeness warning
- **AND** the warning SHALL be considered by effect validation and gatekeeper audit.