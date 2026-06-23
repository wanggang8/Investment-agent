## ADDED Requirements

### Requirement: P4 review fixes for workflow-backed APIs

Workflow-backed P4 APIs SHALL expose stable state when degraded behavior occurs.

#### Scenario: Evidence refresh marks index failure
- **WHEN** SQLite evidence facts are written but vector indexing fails
- **THEN** affected RAG chunks SHALL be marked with `index_status=failed`
- **AND** SQLite intelligence and source verification facts SHALL remain committed
