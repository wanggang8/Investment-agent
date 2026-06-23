## ADDED Requirements

### Requirement: P8 review hardening coverage
The system SHALL include regression tests and documented verification gates for decision confirmation, rule proposal flow, account snapshot consistency, and frontend behavior.

#### Scenario: Confirmation transaction is atomic
- **WHEN** a manual execution confirmation fails during any dependent write
- **THEN** tests SHALL verify that no partial state remains.

#### Scenario: Consult-to-confirm is executable when contract allows it
- **WHEN** a decision detail advertises offline confirmation actions
- **THEN** tests SHALL verify the confirmation endpoint accepts the corresponding request.

#### Scenario: Full validation includes frontend tests
- **WHEN** the project is prepared for review or archive
- **THEN** the documented validation SHALL include `go test ./...`, `cd web && npm run build`, and `cd web && npm test`.
