## ADDED Requirements

### Requirement: Local tasks can trigger real data maintenance
The system SHALL allow the local `cmd/agent` entrypoint to trigger market refresh, intelligence indexing, and VecLite-related maintenance tasks while preserving the existing real data degradation and audit behavior.

#### Scenario: Market refresh can be triggered locally
- **WHEN** the user triggers market refresh through `cmd/agent`
- **THEN** the system uses configured or stub data sources, records audit events, and reports readable errors on source or write failure.

#### Scenario: Intelligence indexing can be triggered locally
- **WHEN** the user triggers intelligence indexing through `cmd/agent`
- **THEN** the system updates local intelligence/RAG data using existing repositories and records task execution in audit events.

#### Scenario: Index recovery is documented
- **WHEN** VecLite index data is unavailable or damaged
- **THEN** the system documentation describes how to rebuild or recover from local persisted data.
