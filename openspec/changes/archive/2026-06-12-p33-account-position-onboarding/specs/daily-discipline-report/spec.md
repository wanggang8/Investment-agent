## MODIFIED Requirements

### Requirement: Expose today's daily discipline report

The system SHALL expose a local today endpoint and frontend surface for the current local date's daily discipline report, derived from daily workflow or daily auto-run results, without creating trading execution side effects. When account or holdings prerequisites are missing, the report surface SHALL link to the local account and holdings onboarding flow.

#### Scenario: Successful report is available for today's local date

- **GIVEN** the daily discipline workflow or daily auto-run has produced a reportable result for today's local date and holdings scope
- **WHEN** the user opens the today daily discipline report surface or the frontend requests the today API
- **THEN** the system SHALL return the current local date report status, summary, holdings scope, generated time, and related decision/evidence/audit references
- **AND** the frontend SHALL display the report with a clear manual review and non-trading boundary
- **AND** the system SHALL NOT call broker APIs, create order requests, or mark any operation as executed

#### Scenario: Missing prerequisites are reported for today's local date

- **GIVEN** today's daily discipline report cannot be produced because account, holdings, market data, evidence, rules, configuration, or prior workflow prerequisites are missing
- **WHEN** the user opens the today daily discipline report surface or the frontend requests the today API
- **THEN** the system SHALL return a structured missing-prerequisites status for today's local date
- **AND** it SHALL list the missing prerequisite categories in user-readable form
- **AND** it SHALL provide a local onboarding link when account or holdings are missing
- **AND** it SHALL NOT fabricate a report summary, evidence, expected return, or trading instruction
