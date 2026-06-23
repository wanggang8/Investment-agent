## ADDED Requirements

### Requirement: Local index health status
The system SHALL expose health for the local JSON index adapter, including healthy, missing, corrupted, incompatible, rebuilding, and degraded states.

#### Scenario: Index file missing
- **WHEN** the configured local index file does not exist
- **THEN** the system SHALL report `missing`
- **THEN** the system SHALL allow rebuild from SQLite `rag_chunks`

#### Scenario: Index file corrupted
- **WHEN** the configured local index file cannot be parsed as valid JSON
- **THEN** the system SHALL report `corrupted`
- **THEN** retrieval SHALL fall back to SQLite summaries when available

#### Scenario: Index file incompatible
- **WHEN** the configured local index file has an unsupported metadata version
- **THEN** the system SHALL report `incompatible`
- **THEN** rebuild SHALL replace it from SQLite-derived chunks

### Requirement: Rebuild statistics
The system SHALL report rebuild statistics for the local index, including indexed chunk count, skipped chunk count, last rebuild time, and last degradation reason when applicable.

#### Scenario: Rebuild succeeds
- **WHEN** the local index is rebuilt from SQLite chunks
- **THEN** the system SHALL report indexed chunk count and last rebuild time
- **THEN** the system SHALL report a healthy index status

#### Scenario: Rebuild cannot write index
- **WHEN** the local index cannot be written
- **THEN** the system SHALL report degraded status and the write failure reason
- **THEN** SQLite facts SHALL remain unchanged

### Requirement: Observable retrieval fallback
The system SHALL preserve retrieval fallback context when local index search is empty, corrupted, incompatible, or unavailable.

#### Scenario: SQLite fallback is used
- **WHEN** local index retrieval cannot provide usable evidence and SQLite summaries exist
- **THEN** the system SHALL return evidence from SQLite summaries
- **THEN** the response or workflow context SHALL include a degradation reason

#### Scenario: SQLite fallback is insufficient
- **WHEN** local index retrieval fails and SQLite summaries are also insufficient
- **THEN** the system SHALL return information-insufficient or evidence-not-found state
- **THEN** the system SHALL include the index failure reason in audit or response metadata
