# P109 Design

## Strategy

P109 keeps sqlite-vec as the first-stage semantic candidate source and adds a deterministic second-stage strategy inside the application service layer. The strategy is local, testable, and does not call LLMs or external rerank services.

## Query Rewrite

Retrieval queries are rewritten from `RetrievalRequest{Symbol, Query}` into a short structured string:

- symbol and original question;
- inferred intent keywords such as buy discipline, sell/risk discipline, valuation, announcement, capital flow, portfolio state, and source verification;
- stable evidence categories expected by the product.

The rewrite is for retrieval only. It does not alter the user-visible question or final rule verdict.

## Candidate Flow

1. Request a wider candidate set from sqlite-vec, using `max(requestedTopK*4, 8)`.
2. Join candidates back to authoritative SQLite summaries.
3. Drop inconsistent candidates and prefer same-symbol facts.
4. Score candidates with vector rank, keyword overlap, source level, verification status, formal evidence role, event type, and indexed freshness.
5. Apply diversity limits so one event type or one summary cannot dominate the topK.
6. Return the reranked bounded set; if no usable candidates remain, fall back to SQLite summaries with degraded reason.

## Fallback

Fallback remains conservative. If sqlite-vec, embedding, metadata alignment, or reranking cannot provide usable evidence, retrieval falls back to SQLite summaries and reports a degraded retrieval quality summary. Background or unverified evidence is never promoted into formal satisfied evidence by reranking.

## Testing

Tests use deterministic in-memory/fake vector indexes and repository fixtures. No real embedding provider or LLM is required for unit tests. Real embedding provider smoke remains `embedding-smoke` from P108.
