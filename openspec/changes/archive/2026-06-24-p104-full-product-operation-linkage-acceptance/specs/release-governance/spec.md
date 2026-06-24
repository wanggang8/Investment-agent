## ADDED Requirements

### Requirement: P104 Product Operation Linkage Acceptance
The project SHALL maintain a repeatable local-source acceptance gate that verifies representative product operations through HTTP APIs, SQLite side effects, downstream readback, audit traceability, and forbidden automation absence.

#### Scenario: Local runner validates linked product behavior
- **GIVEN** the repository source tree is available locally
- **WHEN** the P104 acceptance runner is executed
- **THEN** it SHALL create an isolated temporary SQLite database and config
- **AND** it SHALL start the local backend on localhost
- **AND** it SHALL exercise representative portfolio, decision confirmation, review, audit, notification, risk, and data-quality operations through HTTP APIs
- **AND** it SHALL verify durable SQLite side effects and downstream readback
- **AND** it SHALL fail if forbidden broker/order/push/automatic-confirmation evidence is present.

#### Scenario: Acceptance record stays honest about scope
- **GIVEN** P104 validation has completed
- **WHEN** the release acceptance record is updated
- **THEN** it SHALL distinguish fresh P104 local-source linkage evidence from Docker, installer, package, remote deployment, physical second-machine, broker, automatic trading, automatic confirmation, automatic rule application, and return-guarantee claims.
