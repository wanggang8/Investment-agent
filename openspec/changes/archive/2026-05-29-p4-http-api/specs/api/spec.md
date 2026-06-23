## ADDED Requirements

### Requirement: P4 HTTP API response envelope

P4 HTTP API handlers SHALL return JSON responses using the common API envelope defined by `docs/api.md`.

#### Scenario: Successful business response

- **WHEN** a P4 business API succeeds
- **THEN** the response SHALL include `request_id`
- **AND** the response SHALL include `data`
- **AND** the response MAY include `meta.generated_at` and `meta.rule_version` when applicable

#### Scenario: Failed business response

- **WHEN** a P4 business API fails
- **THEN** the response SHALL include `request_id`
- **AND** the response SHALL include `error.code` and `error.message`
- **AND** unknown internal failures SHALL be returned as `INTERNAL_ERROR` with HTTP 500
- **AND** the response SHALL NOT expose SQL, file paths, or external service raw error text

### Requirement: P4 core API surface

P4 SHALL implement only the HTTP API endpoints listed in `docs/development-plan.md` P4.2 and `docs/api.md`.

#### Scenario: Core endpoints are implemented

- **WHEN** P4 is complete
- **THEN** handlers SHALL exist for dashboard, decision, portfolio, evidence, market, rule, settings, audit, and review API groups
- **AND** each handler SHALL use DTOs aligned with `docs/frontend-contract.md`

### Requirement: P4 confirmation state transitions

The decision confirmation API SHALL enforce the state transition and transaction rules defined by `docs/api.md`.

#### Scenario: Terminal confirmation cannot be repeated

- **WHEN** a decision is already `executed_manually` or `marked_error`
- **THEN** another confirmation request SHALL return `BAD_REQUEST`
- **AND** it SHALL NOT write duplicate account snapshots, transactions, or error cases
