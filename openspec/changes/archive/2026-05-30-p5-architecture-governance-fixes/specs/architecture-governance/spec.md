## ADDED Requirements

### Requirement: Application dependency direction
The system SHALL keep application and workflow packages independent from concrete SQLite implementations.

#### Scenario: Workflow dependency wiring stays outside workflow package
- **WHEN** workflow dependencies are constructed for production use
- **THEN** the construction SHALL happen in the composition root or bootstrap package
- **AND** `internal/application/workflow` SHALL depend on repository interfaces rather than importing `internal/infrastructure/persistence/sqlite`

#### Scenario: Domain remains independent
- **WHEN** domain packages are built or tested
- **THEN** they SHALL NOT import application, infrastructure, HTTP, or database driver packages

### Requirement: Thin HTTP handlers
HTTP handlers SHALL only parse requests, call application use cases or workflow entrypoints, and write response envelopes.

#### Scenario: Handler performs a cross-table write
- **WHEN** an API action writes multiple tables such as confirmations, portfolio updates, rule application, or audit events
- **THEN** the handler SHALL delegate the write to an application service
- **AND** the handler SHALL NOT directly manage SQL transactions

### Requirement: Transaction coordination
The system SHALL provide one reusable transaction coordination path for multi-repository writes.

#### Scenario: Multi-table write succeeds
- **WHEN** a use case writes related facts across multiple repositories
- **THEN** all facts and related audit events SHALL commit as one unit when the use case succeeds

#### Scenario: Multi-table write fails
- **WHEN** any write inside the transactional use case fails
- **THEN** the transaction SHALL roll back all facts written in that transaction
- **AND** the use case SHALL return a mapped application error

### Requirement: ID and Clock generation
The system SHALL use injectable ID and Clock helpers for generated business IDs and timestamps.

#### Scenario: Production ID and time generation
- **WHEN** handlers, workflows, or application services create request-linked entities, audit events, or records
- **THEN** they SHALL use the shared ID and Clock helpers
- **AND** persisted times SHALL use UTC RFC3339 strings unless the target table requires another documented format

#### Scenario: Deterministic tests
- **WHEN** tests verify generated IDs or timestamps
- **THEN** they SHALL be able to inject deterministic ID and Clock implementations

### Requirement: Contract enum validation
The system SHALL define contract enum values in domain-level types and reuse them for validation.

#### Scenario: Handler validates a contract enum
- **WHEN** a handler validates fields such as confirmation type, operation type, rule proposal status, or error-case root cause tag
- **THEN** validation SHALL use domain model constants and `Valid()` methods
- **AND** local string-only enum switches SHALL NOT be the source of truth

### Requirement: P5 frontend feature organization
The P5 frontend SHALL organize business UI by feature before cockpit pages expand.

#### Scenario: Feature page implementation
- **WHEN** a P5 page is implemented for dashboard, decision, evidence, rules, audit, settings, market, portfolio, or review
- **THEN** page-specific components, DTO mappers, and API access SHALL live under the corresponding `web/src/features/<feature>` area
- **AND** `web/src/pages` SHALL compose feature modules rather than holding most business logic

#### Scenario: Shared frontend code
- **WHEN** code is reused across multiple frontend features
- **THEN** it SHALL live in shared service, type, utility, or component locations rather than inside a single feature
