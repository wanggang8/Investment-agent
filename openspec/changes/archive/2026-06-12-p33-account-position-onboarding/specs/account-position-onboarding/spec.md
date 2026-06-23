## ADDED Requirements

### Requirement: Initialize a local account and holdings from an empty database

The system SHALL allow a user to initialize a local account with cash, total assets, holdings, cost basis, buy reasons, and basic risk preference information without connecting to any broker or external trading account.

#### Scenario: Empty database user completes account initialization

- **GIVEN** the local database has no account snapshot and no current holdings
- **WHEN** the user submits a valid account initialization form
- **THEN** the system SHALL write a portfolio snapshot, current positions, position snapshots, and audit events in one transaction
- **AND** the frontend SHALL show that account and holdings prerequisites are satisfied
- **AND** the system SHALL NOT call broker APIs, place orders, or mark any operation as executed by the system

#### Scenario: Invalid initialization is rejected

- **GIVEN** the user submits missing symbols, negative cash, negative quantity, missing cost basis, or inconsistent total assets beyond documented tolerance
- **WHEN** the initialization request is validated
- **THEN** the system SHALL reject the request with structured validation errors
- **AND** it SHALL NOT write partial account, holding, snapshot, transaction, or audit facts that imply a successful initialization

### Requirement: Maintain holdings through append-only local facts

The system SHALL support adding, editing, and removing current holdings through confirmed local facts while preserving historical snapshots and audit traceability.

#### Scenario: User edits a current holding

- **GIVEN** an existing current holding and account snapshot
- **WHEN** the user confirms a valid holding edit with a reason
- **THEN** the system SHALL create a new account snapshot and position snapshot set
- **AND** it SHALL preserve prior historical snapshots
- **AND** it SHALL write an audit event linking the before and after local facts

#### Scenario: User removes a current holding

- **GIVEN** an existing current holding
- **WHEN** the user confirms removal from the current holdings view
- **THEN** the system SHALL represent the removal through a new local fact or snapshot state
- **AND** historical records SHALL remain queryable
- **AND** the UI SHALL explain that removal changes the current local view and does not delete historical facts

### Requirement: Record offline transactions and keep account consistency

The system SHALL allow users to record transactions that they already executed offline, then update local cash, holdings, snapshots, transactions, and audit records consistently.

#### Scenario: Offline buy transaction is recorded

- **GIVEN** the user has an initialized account
- **WHEN** the user records an offline buy transaction with symbol, quantity, price, executed time, fees if available, and note
- **THEN** the system SHALL write a `position_transactions` record
- **AND** it SHALL update current holdings, cash, portfolio snapshot, position snapshots, and audit events in one transaction
- **AND** the UI SHALL label the action as recording an offline action, not system trading

#### Scenario: Offline sell or reduce transaction is recorded

- **GIVEN** the user has enough local quantity for a holding
- **WHEN** the user records an offline sell or reduce transaction
- **THEN** the system SHALL reduce the local holding quantity, update cash and snapshots, and write an audit event
- **AND** it SHALL reject the request if the local quantity would become negative

### Requirement: Validate batch import before writing facts

The system SHALL support batch or table-based import of holdings and historical transactions with row-level validation before any facts are written.

#### Scenario: Batch import contains mixed valid and invalid rows

- **GIVEN** the user provides multiple holding or transaction rows
- **WHEN** the user requests validation
- **THEN** the system SHALL return row-level validation results and a summary of valid and invalid rows
- **AND** it SHALL NOT write account, holding, transaction, or snapshot facts during validation

#### Scenario: User confirms a valid batch import

- **GIVEN** a batch import has no validation errors and the validation response returned an `import_batch_id`
- **WHEN** the user confirms import with that same `import_batch_id`
- **THEN** the system SHALL verify the batch exists, is still `validated`, has zero invalid rows, and matches the validated rows hash
- **AND** it SHALL write the imported facts in a single transaction
- **AND** it SHALL write audit events that identify the import batch and source rows

### Requirement: Correct input mistakes without silently overwriting history

The system SHALL provide an explicit correction flow for input mistakes, preserving before/after references and correction reasons.

#### Scenario: User corrects an entered holding or transaction

- **GIVEN** a user identifies an input mistake in local account, holding, or transaction data
- **WHEN** the user submits a correction with a reason
- **THEN** the system SHALL write a correction fact and an audit event with before and after references
- **AND** it SHALL NOT silently overwrite or physically delete historical facts
- **AND** when the correction should change current holdings, cash, or snapshots, the user SHALL use the confirmed holding edit or offline transaction flow so a new snapshot state is written explicitly

### Requirement: Guide first-time users to the next prerequisite

The frontend SHALL guide first-time users from an empty local database to a state where daily discipline can run, showing missing prerequisites and safe next actions.

#### Scenario: User opens the app before account initialization

- **GIVEN** no account snapshot or current holdings exist
- **WHEN** the user opens the Dashboard or today's daily discipline report
- **THEN** the frontend SHALL show an onboarding state with account initialization as the next action
- **AND** it SHALL explain that no formal daily discipline recommendation is generated until prerequisites are present
- **AND** it SHALL NOT show automatic trading, one-click order, profit guarantee, or deterministic price prediction controls

### Requirement: Expose local account onboarding APIs

The system SHALL expose local-only HTTP APIs for account initialization, holding maintenance, offline transaction recording, batch import validation/confirmation, and correction facts. These APIs SHALL use structured validation errors, require explicit user confirmation for mutations beyond initialization/calibration, and SHALL not connect to broker APIs or mark any action as executed by the system.

The P33 API surface SHALL include:

- `POST /api/v1/portfolio/holdings`
- `POST /api/v1/portfolio/holdings/remove`
- `POST /api/v1/portfolio/offline-transactions`
- `POST /api/v1/portfolio/imports/validate`
- `POST /api/v1/portfolio/imports/confirm`
- `POST /api/v1/portfolio/corrections`

#### Scenario: API records local account facts atomically

- **GIVEN** a valid local account onboarding request
- **WHEN** the user calls the corresponding portfolio API
- **THEN** the response SHALL include local fact references such as snapshot id, position id, transaction id, import batch id, correction id, and audit event ids as applicable
- **AND** the mutation SHALL commit account, holding, snapshot, transaction/import/correction, and audit writes in one transaction when multiple facts are required
- **AND** a failed validation or persistence error SHALL leave no partial successful local account state

#### Scenario: API validation reports safe errors

- **GIVEN** a request contains invalid amount, quantity, date, missing symbol/name/cost/buy reason, inconsistent total assets, or sell quantity above the current local holding
- **WHEN** the API validates the request
- **THEN** it SHALL return a structured validation error suitable for frontend display
- **AND** it SHALL not expose internal stack traces or imply that a trade happened

### Requirement: Preserve append-only account data model semantics

The system SHALL represent P33 local account changes through append-only snapshots, transactions, import batches, corrections, and audit events. Current views MAY be updated for latest-state reads, but historical snapshots and correction records SHALL remain queryable.

#### Scenario: Import batch and correction records are stored separately

- **GIVEN** a user validates or confirms a batch import or records an input correction
- **WHEN** the system stores local facts
- **THEN** import batch metadata SHALL be stored with row counts, validation summary, validated rows hash, status, timestamps, and request id
- **AND** correction facts SHALL store target type, target id, before/after payloads, correction reason, optional snapshot/audit references, and created time
- **AND** neither import nor correction records SHALL replace `decision_records` or formal audit facts as the main trace source

### Requirement: Provide frontend onboarding and correction view models

The frontend SHALL expose a Portfolio onboarding surface that maps API data to view models for initialization, holding edit/remove, offline transaction, batch import validation/confirmation, and correction review. Dashboard and daily discipline report surfaces SHALL link to this onboarding surface when account or holdings prerequisites are missing.

#### Scenario: User reviews local-only actions before submission

- **GIVEN** the user opens the Portfolio onboarding and maintenance surface
- **WHEN** the user enters initialization, holding edit, offline transaction, batch import, or correction data
- **THEN** the view model SHALL display cash, total assets, holding cost, buy reason, asset tag, risk preference, row validation status, before/after correction information, and audit references when available
- **AND** the UI SHALL state that these actions only record local facts and do not place orders, connect to brokers, or produce guaranteed returns
