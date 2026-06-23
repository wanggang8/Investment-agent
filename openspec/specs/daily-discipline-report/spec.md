# daily-discipline-report Specification

## Purpose

This specification records the local daily discipline report product layer introduced by P32. The report layer turns DailyDisciplineGraph and daily auto-run outcomes into a local, read-only today/history/detail surface without changing trading, data-source, or rule-arbitration boundaries.
## Requirements
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

### Requirement: Daily discipline reports SHALL surface expanded data coverage state

Daily discipline report surfaces SHALL show whether P34 expanded public data was fresh, stale, missing, unavailable, or degraded when the report was generated.

#### Scenario: Expanded data is available for a report
- **WHEN** a daily discipline report is generated with available P34 expanded data
- **THEN** the report context SHALL include source category, source level, data date, freshness state, and affected symbols or indexes
- **AND** the frontend SHALL be able to display that the report used expanded public data as analysis context.

#### Scenario: Expanded data is missing or stale for a report
- **WHEN** a daily discipline report lacks required P34 data categories or receives stale data
- **THEN** the report SHALL include missing or stale categories in its diagnostics
- **AND** it SHALL not mark those categories as satisfied by unrelated or lower-grade data.

#### Scenario: Expanded source fails during report preparation
- **WHEN** a P34 source fails with no data, source unavailable, parse error, timeout, or write failure during report preparation or refresh
- **THEN** the report SHALL preserve a degraded or insufficient-data explanation
- **AND** it SHALL not imply that a broker trade, order, external notification, or guaranteed return was produced.

### Requirement: Daily discipline reports SHALL surface risk alert SOP state

Daily discipline report surfaces SHALL include related risk alert summaries so users can understand which risks were active when the report was generated and which SOP state applies.

#### Scenario: Report has related risk alerts
- **WHEN** a daily discipline report is generated and related risk alerts are active, observing, or escalated
- **THEN** the report detail API and frontend SHALL show risk type, severity, SOP status, affected symbol, trigger summary, prohibited actions, suggested manual actions, and links to related risk alert detail
- **AND** the report SHALL keep the wording local-only and manual-review oriented.

#### Scenario: Report has no active risk alerts
- **WHEN** a daily discipline report has no related active risk alerts
- **THEN** the report SHALL show either an empty risk summary or omit the risk section
- **AND** it SHALL NOT imply that future losses are impossible or that the system guarantees safety.

### Requirement: Daily discipline risk summaries SHALL remain read-only

Daily discipline report risk summaries SHALL only expose local reading and traceability links, not trading or rule-application controls.

#### Scenario: User opens risk summary from report
- **WHEN** the frontend displays a risk summary inside a daily discipline report
- **THEN** it SHALL link to risk alert detail, decision detail, notification, or audit records where available
- **AND** it SHALL NOT display automatic trade, broker connection, one-click sell, one-click buy, or automatic rule update controls.

### Requirement: Daily discipline workflow risk integration SHALL be documented for archive

The P35 delta SHALL record the workflow and frontend contract changes for risk alert orchestration to be merged into `docs/workflow.md` and `docs/frontend-contract.md` during archive.

#### Scenario: Workflow docs are archived
- **WHEN** workflow docs are updated
- **THEN** they SHALL state that daily discipline report generation calls risk alert orchestration after decision/report persistence and passes decision ID, report ID, request ID, market snapshot, and source health context
- **AND** the workflow SHALL be described as local-only and non-trading.

#### Scenario: Frontend contract docs are archived
- **WHEN** frontend contract docs are updated
- **THEN** they SHALL describe the risk alert center, dashboard risk summary, daily discipline report risk summary, notification risk link, SOP status labels, severity labels, prohibited actions, suggested manual actions, and local-only safety wording.

#### Scenario: Read-only summary contract is archived
- **WHEN** frontend risk summaries are documented
- **THEN** the contract SHALL state that embedded risk summaries are read-only traceability surfaces
- **AND** only the risk alert center MAY expose explicit local SOP lifecycle actions.

