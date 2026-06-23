## ADDED Requirements

### Requirement: P4 transactional API writes

P4 write APIs SHALL preserve the transaction boundaries defined by `docs/data-model.md` and `docs/api.md`.

#### Scenario: Executed manual confirmation writes all required facts

- **WHEN** `POST /api/v1/decisions/{decision_id}/confirmations` succeeds with `confirmation_type=executed_manually`
- **THEN** the system SHALL atomically write `operation_confirmations`, `position_transactions`, `positions`, `portfolio_snapshots`, `position_snapshots`, and `audit_events`

#### Scenario: Marked error confirmation writes error facts

- **WHEN** `POST /api/v1/decisions/{decision_id}/confirmations` succeeds with `confirmation_type=marked_error`
- **THEN** the system SHALL atomically write `operation_confirmations`, `error_cases`, and `audit_events`
- **AND** the response SHALL include `error_case_id`

#### Scenario: Planned and watch confirmations do not mutate account facts

- **WHEN** confirmation succeeds with `confirmation_type=planned` or `confirmation_type=watch`
- **THEN** the system SHALL NOT write `position_transactions`
- **AND** the system SHALL NOT create a new account snapshot
