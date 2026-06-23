## ADDED Requirements

### Requirement: P4 review fixes for transactional writes

P4 review fixes SHALL ensure rule application and confirmation corrections preserve documented transaction boundaries.

#### Scenario: Final rule confirmation writes audit atomically
- **WHEN** final confirmation applies a rule proposal
- **THEN** archiving old active rules, creating the new active rule, updating proposal state, and writing `audit_events` SHALL happen in one transaction

#### Scenario: Invalid confirmation requests write no facts
- **WHEN** a confirmation request fails action-specific validation
- **THEN** the system SHALL NOT write `operation_confirmations`, account facts, error cases, or audit events
