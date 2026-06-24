# Index and Retrieval Hardening Specification

## Purpose

Define local JSON index health, rebuild statistics, retrieval fallback observability, and the replacement boundary for a future real VecLite API.
## Requirements
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
The system SHALL preserve retrieval fallback context when local index search is empty, corrupted, incompatible, stale, or unavailable, and SHALL expose enough retrieval quality metadata to explain whether returned evidence came from index search or SQLite fallback.

#### Scenario: SQLite fallback is used
- **WHEN** local index retrieval cannot provide usable evidence and SQLite summaries exist
- **THEN** the system SHALL return evidence from SQLite summaries
- **THEN** the response or workflow context SHALL include a degradation reason and fallback source

#### Scenario: SQLite fallback is insufficient
- **WHEN** local index retrieval fails and SQLite summaries are also insufficient
- **THEN** the system SHALL return information-insufficient or evidence-not-found state
- **THEN** the system SHALL include the index failure reason in audit or response metadata

#### Scenario: Retrieval quality summary is available
- **WHEN** evidence retrieval returns top-k evidence
- **THEN** the system SHALL expose query summary, top-k count, hit/miss status when evaluated, index health, index freshness, fallback source, and degraded reason when known
- **AND** this summary SHALL NOT include secrets or complete local file paths

#### Scenario: Retrieval quality smoke is run locally
- **WHEN** the operator runs the local `retrieval-quality-smoke` task for a symbol
- **THEN** the system SHALL evaluate retrieval through the normal local adapter and write an audit summary containing status, top-k, fallback source, index health, source consistency, and the no-auto-trading boundary
- **AND** the task SHALL NOT write account, confirmation, transaction, broker, or rule-application facts

### Requirement: Retrieval quality evaluation
The system SHALL provide a repeatable local retrieval quality evaluation over representative fixtures.

#### Scenario: Expected evidence is retrieved
- **WHEN** a retrieval quality fixture defines expected evidence ids or expected source constraints
- **THEN** evaluation SHALL report whether top-k retrieval satisfied the expected evidence or constraints
- **AND** misses SHALL include diagnostic fields for index health, fallback source, and source consistency status

#### Scenario: Background-only evidence appears
- **WHEN** retrieval returns C-level, background-only, or unverified evidence for a formal query
- **THEN** evaluation SHALL report it as background or unexpected formal evidence
- **AND** rule arbitration SHALL NOT treat it as satisfied formal evidence

### Requirement: Quality-aware retrieval ranking
The system SHALL rank or filter retrieval results using text relevance together with evidence quality metadata.

#### Scenario: Formal verified evidence outranks weaker background evidence
- **WHEN** multiple retrieval candidates have comparable text relevance
- **THEN** verified formal S/A/B evidence SHOULD rank above C-level, background, stale, or unverified evidence

#### Scenario: Index freshness affects retrieval status
- **WHEN** the index is stale, corrupted, incompatible, or rebuilt from older chunks
- **THEN** retrieval SHALL preserve this status in metadata or audit
- **AND** returned evidence SHALL remain traceable to SQLite facts

#### Scenario: Indexed chunks are stale
- **WHEN** retrieved index chunks are older than the local freshness window
- **THEN** retrieval quality metadata SHALL report `index_freshness=stale`
- **AND** the stale status SHALL NOT promote background or unverified evidence into formal decision evidence

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

### Requirement: P109 SHALL optimize retrieval quality with deterministic local strategy

The system SHALL improve sqlite-vec retrieval quality using deterministic local query rewrite, metadata-aware reranking, and evidence diversity without relying on external rerank services or changing final verdict ownership.

#### Scenario: Retrieval query is rewritten for investment context

- **GIVEN** a consultation retrieval request includes a symbol and user question
- **WHEN** semantic retrieval is executed
- **THEN** the query sent to the semantic index SHALL include the original user question, the symbol, inferred investment intent keywords, and expected evidence categories
- **AND** the rewrite SHALL NOT replace the user-visible question or final verdict text.

#### Scenario: Candidate set is widened before rerank

- **GIVEN** semantic vector search is available
- **WHEN** the requested result count is smaller than the rerank candidate window
- **THEN** sqlite-vec SHALL be queried for a wider candidate set
- **AND** the application SHALL locally rerank and bound the final evidence set.

#### Scenario: Metadata-aware rerank prefers reliable evidence

- **GIVEN** vector candidates include same-symbol verified formal evidence and weaker background evidence
- **WHEN** reranking is applied
- **THEN** verified formal S/A/B evidence SHOULD rank above C-level, background, stale, unverified, or symbol-mismatched evidence when text relevance is comparable
- **AND** background evidence SHALL NOT be promoted into satisfied formal evidence.

#### Scenario: Evidence diversity is preserved

- **GIVEN** multiple top vector candidates map to the same summary, source, or event type
- **WHEN** the final topK evidence set is selected
- **THEN** the system SHOULD avoid duplicate-dominated results when alternative relevant evidence exists
- **AND** returned evidence SHALL remain traceable to authoritative SQLite summaries.

