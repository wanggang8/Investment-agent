## ADDED Requirements

### Requirement: P4 review fixes for HTTP API behavior

P4 HTTP API fixes SHALL complete the behavior verified by post-archive review without adding new API groups.

#### Scenario: Confirmation requests validate action-specific fields
- **WHEN** confirmation uses `executed_manually`
- **THEN** `operation_type` SHALL be one of `buy`, `sell`, or `reduce`
- **AND** `executed_at` SHALL NOT be later than current server time
- **WHEN** confirmation uses `marked_error`
- **THEN** `actual_outcome`, `root_cause_tag`, and `lesson_learned` SHALL be required
- **AND** `root_cause_tag` SHALL match the documented enum

#### Scenario: Rule proposal rejection is explicit
- **WHEN** user confirms a rule proposal with `confirm=false`
- **THEN** the API SHALL mark the proposal as `rejected`
- **AND** write an audit event

#### Scenario: Market refresh returns stable error codes
- **WHEN** all market symbols fail due to source or input issues
- **THEN** the API SHALL return a stable market-related error code instead of raw internal errors
