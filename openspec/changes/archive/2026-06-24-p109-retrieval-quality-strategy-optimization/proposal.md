# P109 Retrieval Quality Strategy Optimization

## Why

P108 added real sqlite-vec semantic retrieval, but retrieval quality still depends mostly on the raw user question and vector distance. Real product usefulness needs deterministic program strategies around embeddings: query rewrite, metadata-aware filtering, hybrid reranking, evidence diversity, and regression evaluation.

## What Changes

- Add deterministic query rewrite for investment consultation retrieval so semantic search uses symbol, user intent, discipline keywords, and evidence categories.
- Add metadata-aware candidate scoring using symbol match, source level, verification status, evidence role, event type, keyword overlap, and freshness.
- Fetch a wider sqlite-vec candidate set, rerank locally, and return a bounded diversified topK evidence set.
- Preserve safe fallback to SQLite summaries when vector retrieval is unavailable or metadata consistency fails.
- Extend retrieval quality metadata and regression tests so strategy changes are repeatable and auditable.

## Out Of Scope

- New embedding providers or embedding model selection logic.
- External rerank APIs, external search services, or new paid/login/authorized data sources.
- Changing final verdict ownership, rules, broker connectivity, automatic trading, external push, automatic confirmation, or automatic rule application.
- Docker, installer, physical second-machine, GitHub Release, or tag validation.

## Validation

- TDD focused Go tests for query rewrite, metadata-aware rerank, candidate widening, diversity, and fallback behavior.
- Focused retrieval quality evaluation tests.
- Full Go test suite.
- Frontend test/build because decision DTO/retrieval quality metadata remains user-facing.
- OpenSpec strict validation and whitespace check.
