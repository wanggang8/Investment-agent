## ADDED Requirements

### Requirement: P40 Local Recovery Smoke Acceptance
The system SHALL include local smoke coverage for deployment readiness and backup recovery that verifies restored data through supported APIs or browser-visible behavior.

#### Scenario: Recovery smoke verifies readable restored facts
- **WHEN** the P40 recovery smoke restores a backup into temporary local paths
- **THEN** it SHALL run migrations or compatibility checks as needed
- **AND** it SHALL verify at least one restored decision, audit event, position, report, or equivalent historical fact through API or browser-visible behavior
- **AND** it SHALL NOT require manual SQLite inspection

#### Scenario: Recovery smoke preserves safety boundaries
- **WHEN** recovery smoke, startup smoke, or browser smoke executes
- **THEN** it SHALL use temporary local paths and non-secret fixture data by default
- **AND** it SHALL NOT connect to broker APIs, place orders, send external pushes, automatically confirm user actions, or automatically apply rules
