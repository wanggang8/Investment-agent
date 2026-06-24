# P108 Design

## Architecture

SQLite remains the authoritative store for evidence facts and `rag_chunks`. P108 adds a separate sqlite-vec index file, configured independently from the primary SQLite database. This avoids replacing the existing `modernc.org/sqlite` driver while still using the official sqlite-vec CGO binding with `mattn/go-sqlite3`, as recommended by upstream.

## Embedding Configuration

P108 adds an `embedding` configuration block:

```yaml
embedding:
  enabled: true
  provider: openai_compatible
  api_key: ""
  base_url: ""
  model: ""
  dimensions: 1536
  timeout_seconds: 60
```

Chat/analysis model config and embedding config are intentionally separate. Unit tests use a deterministic local embedding provider; real smoke tests use `/embeddings` only when a configured provider supports it. Chat completions are not accepted as a substitute for embedding tests.

## Vector Index

The existing `VectorIndex` interface is extended only as much as needed to support query-vector retrieval while keeping fallback compatibility. The sqlite-vec implementation stores chunk embeddings in a `vec0` virtual table and stores chunk metadata in a normal SQLite table inside the vector index file.

## Data Flow

Evidence refresh/rebuild:

1. Load authoritative RAG chunks from primary SQLite.
2. Generate embeddings for chunk text.
3. Upsert chunk metadata and vectors into the sqlite-vec index.
4. Mark primary SQLite chunks as indexed after successful vector write.

Consultation retrieval:

1. Generate query embedding from the symbol/question query text.
2. Query sqlite-vec topK nearest chunks.
3. Verify chunk metadata against authoritative SQLite summaries.
4. Return evidence summaries and retrieval quality metadata.
5. Fall back to SQLite summary retrieval if embeddings or sqlite-vec are unavailable.

## Safety

The vector index only changes evidence retrieval order and relevance. It does not change final verdict ownership, does not trigger trades, does not auto-confirm, and does not auto-apply rules.
