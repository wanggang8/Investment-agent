## ADDED Requirements

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
