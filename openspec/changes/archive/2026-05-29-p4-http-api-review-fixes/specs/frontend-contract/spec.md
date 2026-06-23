## ADDED Requirements

### Requirement: P4 review fixes for frontend DTO compatibility

Rule proposal and settings DTOs SHALL include fields needed by the P5 frontend contract.

#### Scenario: Rule proposal list is renderable
- **WHEN** frontend reads rule proposals
- **THEN** each proposal DTO SHALL expose proposal metadata, before/after rule snapshots, impact scope, risk notes, and optional audit summary fields when available

#### Scenario: Settings preserve page preference
- **WHEN** frontend saves ordinary system settings
- **THEN** page preference SHALL be persisted and returned through the settings API
