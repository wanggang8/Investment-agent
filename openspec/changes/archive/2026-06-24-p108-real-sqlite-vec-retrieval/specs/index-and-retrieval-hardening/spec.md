## ADDED Requirements

### Requirement: P108 SHALL use real sqlite-vec retrieval when embeddings are configured

P108 SHALL replace the JSON-only vector-index path with a real sqlite-vec backed retrieval path when embedding configuration is enabled.

#### Scenario: Evidence chunks are indexed with embeddings

- **GIVEN** embedding configuration is enabled and RAG chunks exist in SQLite
- **WHEN** evidence indexing or vector-index rebuild runs
- **THEN** the system SHALL generate embeddings for chunk text and upsert them into a sqlite-vec vector table
- **AND** SQLite `rag_chunks` SHALL remain the authoritative rebuild source.

#### Scenario: Retrieval uses semantic topK

- **GIVEN** sqlite-vec has indexed chunk embeddings
- **WHEN** a consultation retrieval request is evaluated
- **THEN** the system SHALL generate a query embedding and return topK chunks ordered by sqlite-vec vector distance
- **AND** returned chunks SHALL be checked against authoritative SQLite summaries before becoming evidence.

#### Scenario: Embedding configuration is separate from chat analysis

- **GIVEN** chat analysis configuration is present
- **WHEN** embedding retrieval is enabled
- **THEN** the system SHALL require explicit embedding provider, base URL, model, dimensions, and timeout configuration
- **AND** chat completion output SHALL NOT be accepted as a substitute for embedding vectors.

#### Scenario: Retrieval degrades safely

- **GIVEN** embedding generation or sqlite-vec retrieval is unavailable
- **WHEN** retrieval is requested
- **THEN** the system SHALL fall back to SQLite summary retrieval with a degraded retrieval-quality reason
- **AND** it SHALL NOT change final verdict, execute trades, auto-confirm, push externally, or auto-apply rules.
