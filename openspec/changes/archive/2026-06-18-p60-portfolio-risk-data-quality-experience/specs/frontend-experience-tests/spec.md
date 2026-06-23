## ADDED Requirements

### Requirement: Portfolio, risk, and data-quality pages present operational maintenance experiences

The frontend SHALL turn positions, risk alerts, and data quality into readable maintenance and disposition surfaces while preserving existing backend contracts and local-only safety boundaries.

#### Scenario: Positions page explains the portfolio maintenance state before local fact forms

- **WHEN** the user opens `/positions`
- **THEN** the first screen MUST show portfolio status, snapshot context, total assets, cash ratio, position count, high-risk ratio, current maintenance stage, and next manual actions before detailed forms or tables
- **AND** it MUST distinguish initialization, calibration, holding edit, offline transaction recording, batch import validation/confirmation, and correction audit paths
- **AND** every write action MUST be described as a local fact or audit record, not broker sync, order placement, automatic portfolio management, or return optimization

#### Scenario: Positions page preserves explicit local-only write boundaries

- **WHEN** the user submits calibration, holding edit, holding removal, offline transaction, batch import, or correction audit actions
- **THEN** the frontend MUST call only the existing portfolio service methods and display safe success or error messages
- **AND** disabled or unavailable actions MUST explain the missing local prerequisite instead of implying automatic recovery
- **AND** the page MUST NOT expose broker login, automatic trade, one-click order, delegated order, external push, automatic confirmation, automatic rule application, automatic repair, database overwrite, or return promise controls

#### Scenario: Risk alerts render as a disposition queue

- **WHEN** the user opens `/risk-alerts`
- **THEN** the page MUST show risk disposition summary, severity, affected symbols, and queues for pending review, in progress, needs review, and recorded risks
- **AND** each alert MUST show risk type, severity, SOP status, trigger summary, prohibited actions, suggested manual actions, related local links, updated time, and safety note
- **AND** the queue mapping MUST treat `triggered` as pending review, `active` and `observing` as in progress, `escalated` as needs review, and `resolved` or `archived` as recorded

#### Scenario: Risk SOP actions remain explicit local lifecycle records

- **WHEN** an unresolved risk alert is shown
- **THEN** lifecycle controls MAY allow continue observing, escalate for review, or resolve locally through the existing risk alert lifecycle service
- **AND** resolved or archived risks MUST NOT show lifecycle controls
- **AND** SOP controls MUST NOT imply automatic trading, external notification, automatic confirmation, rule application, or portfolio mutation

#### Scenario: Data quality page explains quality signals and affected workflows

- **WHEN** the user opens `/data-quality`
- **THEN** the first screen MUST show an overall quality state, source health signal, evidence/RAG signal, LLM signal, affected workflow signal, and next local inspection actions
- **AND** source health, evidence verification, VecLite, DeepSeek, review degradation, missing evidence, and affected decision/workflow details MUST remain visible in readable layers
- **AND** degraded, stale, missing, parse_error, unavailable, failed, unknown, or insufficient states MUST not be displayed as normal success

#### Scenario: Data quality diagnostics remain read-only and sanitized

- **WHEN** data quality APIs return source health failures, evidence summaries, review explanations, system paths, or unexpected diagnostic values
- **THEN** the frontend MUST render safe Chinese summaries and local navigation without exposing API keys, private paths, SQL, complete prompts, raw vendor payloads, local database paths, or raw stack traces
- **AND** the page MUST NOT offer automatic repair, automatic source refresh, automatic confirmation, rule application, external push, database overwrite, or trading controls

#### Scenario: Portfolio, risk, and data-quality experiences remain mobile readable

- **WHEN** `/positions`, `/risk-alerts`, or `/data-quality` renders at 390px viewport width
- **THEN** primary status, next manual actions, local safety boundary, form controls, queue cards, quality signals, and navigation MUST remain readable without page-level horizontal overflow
- **AND** screenshots or browser evidence MUST be captured for desktop and mobile validation
