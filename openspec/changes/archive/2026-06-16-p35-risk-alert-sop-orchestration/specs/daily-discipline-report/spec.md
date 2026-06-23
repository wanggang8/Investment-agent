## ADDED Requirements

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
