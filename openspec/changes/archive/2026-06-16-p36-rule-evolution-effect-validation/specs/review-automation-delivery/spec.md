## ADDED Requirements

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