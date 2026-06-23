# decision-loop-explainability Specification

## Purpose
TBD - created by archiving change p47-decision-loop-explainability. Update Purpose after archive.
## Requirements
### Requirement: Decision loop explanations SHALL be read-only and traceable

The system SHALL provide a read-only decision loop explanation view that links decisions to user confirmations, local manual records, risk/review/audit traces, and missing loop gaps.

#### Scenario: User lists recent decision loops

- **WHEN** the user opens the decision loop surface
- **THEN** the system SHALL return recent decision loop items with decision id, symbol, generated time, final verdict, confirmation status, loop status, stages, related manual actions, risk links, review links, audit links, missing links, and a safety note
- **AND** the API SHALL support `symbol` filtering and a bounded `limit`
- **AND** the API SHALL NOT write `decision_records`, `operation_confirmations`, `position_transactions`, `positions`, `portfolio_snapshots`, `position_snapshots`, `risk_alerts`, `notifications`, `rule_versions`, `audit_events`, broker state, orders, or external notifications

#### Scenario: User opens a single decision loop

- **WHEN** the user requests a decision loop by `decision_id`
- **THEN** the system SHALL return the same read-only loop explanation for that decision
- **AND** missing decision ids SHALL return `NOT_FOUND`
- **AND** the response SHALL include missing loop gaps when expected confirmations, local manual records, risk traces, or review traces are absent

### Requirement: Decision loop explanations SHALL preserve safety boundaries

Decision loop explanations SHALL make existing local facts easier to inspect without adding execution, external delivery, or rule application capabilities.

#### Scenario: Decision loop includes manual action facts

- **WHEN** a loop includes user confirmation or local manual record facts
- **THEN** the response SHALL show only safe summaries, ids, timestamps, quantities, prices, fees, statuses, and links
- **AND** it SHALL NOT expose raw payload JSON, complete secrets, private local paths, raw SQL, full prompts, broker state, order ids, or provider raw responses
- **AND** it SHALL NOT provide broker API, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, automatic repair promise, return promise, login-only source, paid source, authorization-gated source, Level2 source, or high-frequency source capabilities

#### Scenario: Decision loop controls are rendered in the frontend

- **WHEN** the decision loop page renders
- **THEN** the page SHALL show explanation, navigation, and read-only trace links only
- **AND** it SHALL NOT render mutating controls for confirmations, transactions, risk lifecycle, rule application, notifications, or settings

