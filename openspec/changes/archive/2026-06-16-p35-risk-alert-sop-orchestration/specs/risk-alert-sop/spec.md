## ADDED Requirements

### Requirement: Risk alerts SHALL persist local SOP state

The system SHALL persist risk alerts as local facts so users can track risk type, severity, SOP state, trigger evidence, prohibited actions, suggested manual actions, and lifecycle changes across daily discipline runs.

#### Scenario: Risk alert is triggered
- **WHEN** a daily discipline or consultation decision detects valuation high risk, broken buy thesis, liquidity danger, extreme sentiment, position limit breach, insufficient evidence, or degraded data health
- **THEN** the system SHALL create or update a local risk alert with risk type, severity, SOP status, affected symbol, trigger reasons, prohibited actions, suggested manual actions, related decision or report ID, and created or updated time
- **AND** it SHALL write an audit event for the trigger without creating orders, confirmations, broker state, or external notifications.

#### Scenario: Repeated risk is detected
- **WHEN** the same active risk type for the same symbol and source decision is detected again
- **THEN** the system SHALL update the existing active risk or append a traceable occurrence rather than creating conflicting duplicate active alerts
- **AND** it SHALL preserve prior audit history.

### Requirement: Risk alert SOP SHALL support lifecycle transitions

The system SHALL support deterministic SOP lifecycle transitions for risk alerts while preserving human review and non-trading boundaries.

#### Scenario: Risk is observed or escalated
- **WHEN** a risk remains active across subsequent daily discipline runs or its severity increases
- **THEN** the system SHALL move the alert to observing or escalated status with updated reason, timestamps, and audit event
- **AND** it SHALL keep prohibited actions visible to the frontend.

#### Scenario: Risk is resolved or archived
- **WHEN** deterministic inputs show the risk condition no longer applies, or the user explicitly archives the local alert
- **THEN** the system SHALL mark the alert resolved or archived with reason and audit event
- **AND** it SHALL NOT modify positions, portfolio snapshots, operation confirmations, rule versions, or account facts.

### Requirement: Risk alerts SHALL notify and remain local-only

The system SHALL surface active or escalated risk alerts through local application notifications and UI links without any external push channel.

#### Scenario: Risk notification is created
- **WHEN** a risk alert enters active or escalated status
- **THEN** the system SHALL create or refresh a local in-app notification with severity, title, message, source type, source ID, and risk link
- **AND** repeated active risk notifications SHALL be deduplicated by type, source type, and source ID.

#### Scenario: Notification is read or dismissed
- **WHEN** the user marks a risk notification read
- **THEN** the system SHALL update only local notification read state
- **AND** it SHALL NOT change the underlying risk alert lifecycle unless a separate explicit risk action is requested.

### Requirement: Risk alert API SHALL expose traceable status

The system SHALL expose local risk alerts through application APIs so frontend surfaces can list active risks, inspect detail, and perform explicit local lifecycle actions.

#### Scenario: Frontend lists risk alerts
- **WHEN** the frontend requests risk alerts
- **THEN** the API SHALL return alert ID, risk type, severity, SOP status, affected symbol, trigger summary, evidence or data-health references, prohibited actions, suggested manual actions, related decision/report/notification/audit links, and timestamps.

#### Scenario: User updates a risk alert lifecycle
- **WHEN** the user requests resolve, continue observing, escalate, or archive for a risk alert
- **THEN** the API SHALL validate the transition, update the local alert, and write an audit event
- **AND** it SHALL NOT create a trade, order, operation confirmation, external notification, or rule version.

### Requirement: Risk alert UI SHALL show safe SOP guidance

The frontend SHALL provide a risk alert center and embedded summaries that explain current risk status and SOP guidance without presenting automatic trading controls.

#### Scenario: Risk center displays active alerts
- **WHEN** active or escalated alerts exist
- **THEN** the risk alert center SHALL display risk type, severity, affected symbol, current SOP status, trigger evidence, source health or missing data diagnostics, prohibited actions, suggested manual actions, and links to related report, decision, notification, and audit records.

#### Scenario: Risk center is empty
- **WHEN** no active risk alerts exist
- **THEN** the frontend SHALL show an empty state with local-only safety wording
- **AND** it SHALL NOT imply that future risk is impossible or that returns are guaranteed.

### Requirement: Risk alert orchestration SHALL preserve safety boundaries

Risk alert orchestration SHALL remain a local, read-only, human-review feature and SHALL NOT perform trading, external push, or autonomous rule changes.

#### Scenario: Risk SOP recommends human action
- **WHEN** a risk alert suggests cooling off, frozen watch, sell-only review, staged profit-taking review, reassessment, or continued observation
- **THEN** the system SHALL describe the recommendation as manual review guidance
- **AND** it SHALL NOT offer one-click order placement, broker connection, automatic execution, profit guarantee, or deterministic price prediction.

### Requirement: Risk alert API contract SHALL be documented for archive

The P35 delta SHALL record the risk alert API, DTO shape, error categories, and transaction boundary to be merged into `docs/api.md` during archive.

#### Scenario: Risk alert endpoints are archived
- **WHEN** P35 is archived
- **THEN** `docs/api.md` SHALL include `GET /api/v1/risk-alerts`, `GET /api/v1/risk-alerts/{alert_id}`, and `POST /api/v1/risk-alerts/{alert_id}/lifecycle`
- **AND** list responses SHALL use the common envelope with `PageResult<RiskAlertDTO>` and detail/action responses SHALL use `RiskAlertDTO`.

#### Scenario: Risk alert DTO is archived
- **WHEN** API docs are updated
- **THEN** `RiskAlertDTO` SHALL include alert ID, risk type, severity, SOP status, symbol, trigger summary, trigger context, prohibited actions, suggested actions, related decision/report/notification/audit IDs and links, lifecycle timestamps, resolution reason, safety note, created time, and updated time.

#### Scenario: Risk alert API errors are archived
- **WHEN** lifecycle API input is invalid, a transition conflicts with current state, or repositories are unavailable
- **THEN** API docs SHALL map those cases to existing error categories such as bad request, invalid state or conflict, not found, and internal error without exposing unsafe operational detail.

#### Scenario: Risk alert transaction boundary is archived
- **WHEN** risk alerts are created or lifecycle state changes
- **THEN** API docs SHALL state that `risk_alerts`, local `notifications`, and `audit_events` may be written in the same local transaction
- **AND** positions, portfolio snapshots, operation confirmations, position transactions, rule versions, broker state, orders, and external notifications SHALL remain untouched.

### Requirement: Risk alert data model contract SHALL be documented for archive

The P35 delta SHALL record the local `risk_alerts` data model, enums, indexes, and non-trading constraints to be merged into `docs/data-model.md` during archive.

#### Scenario: Risk alert table is archived
- **WHEN** data model docs are updated
- **THEN** `docs/data-model.md` SHALL include `risk_alerts` with `alert_id`, `risk_type`, `severity`, `sop_status`, `symbol`, `trigger_summary`, `trigger_context_json`, `prohibited_actions_json`, `suggested_actions_json`, related IDs, `last_triggered_at`, `resolved_at`, `resolution_reason`, `created_at`, and `updated_at`.

#### Scenario: Risk alert enums and indexes are archived
- **WHEN** data model docs are updated
- **THEN** docs SHALL include risk types `valuation_high`, `buy_thesis_broken`, `liquidity_danger`, `sentiment_extreme`, `position_limit_breach`, `insufficient_evidence`, `data_degraded`
- **AND** SOP statuses `triggered`, `active`, `observing`, `escalated`, `resolved`, `archived`
- **AND** indexes for active identity by risk type plus symbol, status plus updated time, and symbol plus status.

#### Scenario: Risk alert model non-trading boundary is archived
- **WHEN** the data model is documented
- **THEN** the docs SHALL state that `risk_alerts` is a local review fact table and SHALL NOT represent an order, broker instruction, portfolio mutation, or rule-version change.
