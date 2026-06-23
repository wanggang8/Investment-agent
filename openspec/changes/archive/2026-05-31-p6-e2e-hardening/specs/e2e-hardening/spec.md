## ADDED Requirements

### Requirement: End-to-end acceptance plan
The P6 change SHALL provide an end-to-end acceptance plan that covers all A01-A17 testable acceptance assertions from `docs/functional-spec.md` without adding requirements beyond `docs/development-plan.md`.

#### Scenario: A01-A17 assertions are covered
- **WHEN** `docs/testing-plan.md` is reviewed
- **THEN** it SHALL include acceptance coverage for A01 through A17
- **AND** each assertion SHALL state expected observable outcomes aligned with `docs/development-plan.md` P6.1

#### Scenario: No automatic trading acceptance path
- **WHEN** the A15 acceptance assertion is verified
- **THEN** the plan SHALL confirm that no trade execution API or one-click trading frontend entry exists
- **AND** user actions SHALL remain limited to recording offline actions

#### Scenario: Evidence and degradation paths are represented
- **WHEN** A03, A04, A10, A11, A16, and A17 are verified
- **THEN** the plan SHALL cover evidence insufficiency, VecLite degradation, C-level source handling, LLM degradation, market refresh outcomes, and expected-return display states

### Requirement: Configuration and startup documentation
The P6 change SHALL provide configuration and startup documentation for local operation, migration, and seed data as listed in `docs/development-plan.md` P6.2.

#### Scenario: Runtime configuration is documented
- **WHEN** `docs/configuration.md` is reviewed
- **THEN** it SHALL document SQLite data file path, VecLite index file path, DeepSeek API Key environment variable, data source switches, log level, and local startup commands

#### Scenario: Migration and seed are documented
- **WHEN** `docs/migration-plan.md` is reviewed
- **THEN** it SHALL document migration execution and seed data behavior
- **AND** it SHALL avoid real secrets or environment-specific private values

### Requirement: P6 verification commands
The P6 implementation SHALL preserve the verification commands defined in `docs/development-plan.md` for the acceptance hardening phase.

#### Scenario: Backend and frontend verification pass
- **WHEN** P6 implementation is completed
- **THEN** `go test ./...` SHALL pass
- **AND** `cd web && npm run build` SHALL pass
