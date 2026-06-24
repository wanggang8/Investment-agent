# P108 Real sqlite-vec Retrieval

## Why

The project currently exposes a VecLite/vector-index boundary, but the production implementation stores RAG chunks in a local JSON file and filters primarily by symbol. This is not real embedding-based vector retrieval. P108 upgrades the retrieval layer to use a real sqlite-vec backed vector index while preserving SQLite as the authoritative fact store and keeping all investment safety boundaries unchanged.

## What Changes

- Add a separate embedding model configuration for OpenAI-compatible `/embeddings` requests.
- Add a sqlite-vec backed vector index implementation in a separate local vector index file so the existing `modernc.org/sqlite` business database does not need to change drivers.
- Generate embeddings for RAG chunks during evidence indexing/rebuild and query embeddings during retrieval.
- Use sqlite-vec topK vector search before falling back to SQLite summaries.
- Preserve the current `VectorIndex` boundary, index health reporting, SQLite rebuildability, retrieval-quality metadata, and no-auto-trading safety posture.

## Out Of Scope

- Switching the primary business SQLite driver.
- Replacing final rule-based verdict logic.
- Broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorized data-source claims, Level2, or high-frequency data.
- Docker or physical second-machine validation unless separately executed.

## Validation

- TDD unit tests for embedding config validation, deterministic embedding provider behavior, sqlite-vec topK retrieval, and SQLite fallback.
- Focused Go tests for retrieval/vector services and config.
- Full Go test suite.
- OpenSpec strict validation.
- README/docs local consistency check and whitespace check.
