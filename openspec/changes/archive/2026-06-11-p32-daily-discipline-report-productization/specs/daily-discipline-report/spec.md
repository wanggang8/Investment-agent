## ADDED Requirements

### Requirement: Expose today's daily discipline report

The system SHALL expose a local today endpoint and frontend surface for the current local date's daily discipline report, derived from daily workflow or daily auto-run results, without creating trading execution side effects.

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
- **AND** it SHALL NOT fabricate a report summary, evidence, expected return, or trading instruction

### Requirement: List and show historical reports

The system SHALL provide local API and frontend surfaces for listing historical daily discipline reports and showing a selected report's detail, including degraded or missing-prerequisite outcomes.

#### Scenario: History shows prior daily reports and detail

- **GIVEN** one or more daily discipline reports have been recorded for prior local dates or holdings scopes
- **WHEN** the user opens the history surface and selects a report
- **THEN** the system SHALL list reports with date, status, scope summary, generated time, and high-level summary or missing-prerequisite indicator
- **AND** the selected detail SHALL show the report status, full summary, related decision/evidence/audit references, and missing prerequisite or failure diagnostics when present
- **AND** the history and detail surfaces SHALL remain local-only and read-only

### Requirement: Idempotent per local date and holdings scope

The system SHALL keep daily discipline report aggregation idempotent for the same local date and holdings scope so repeated aggregation does not create conflicting duplicate reports.

#### Scenario: Repeated aggregation reuses or updates the same report identity

- **GIVEN** a daily discipline report has already been indexed for a local date and holdings scope
- **WHEN** the same daily workflow result, auto-run result, retry, or manual aggregation is processed again with the same idempotency key
- **THEN** the system SHALL avoid creating a conflicting duplicate report for that date and scope
- **AND** it SHALL reuse the prior report or update it according to documented status transition behavior
- **AND** repeated aggregation SHALL remain visible through timestamps, retry metadata, or audit references without implying a new trading decision was executed
