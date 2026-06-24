# P109 Tasks

## 1. Scope And Governance

- [x] Confirm P109 has an active OpenSpec change before implementation.
- [x] Confirm P109 does not add LLM rerank calls, broker actions, external push, auto-confirm, or auto-rule behavior.
- [x] Confirm SQLite summaries and `rag_chunks` remain authoritative.

## 2. Tests First

- [x] Add failing tests for query rewrite using symbol, original question, and intent/evidence keywords.
- [x] Add failing tests that semantic indexes receive widened candidate topK before local rerank.
- [x] Add failing tests that formal verified same-symbol evidence outranks weak/background evidence.
- [x] Add failing tests for evidence diversity when multiple chunks map to one event type or summary.
- [x] Add failing retrieval quality evaluation tests for expected evidence hit/miss after rerank.

## 3. Implementation

- [x] Add deterministic retrieval query rewrite.
- [x] Add metadata-aware candidate scoring and bounded rerank.
- [x] Add evidence diversity selection.
- [x] Wire strategy into `RetrievalAdapter` while preserving fallback behavior.
- [x] Keep DTO/API compatibility and avoid leaking raw prompts, keys, or paths.

## 4. Documentation And Evidence

- [x] Update architecture/config/release docs for P109 strategy behavior.
- [x] Add P109 acceptance record.
- [x] Archive P109 and update progress/release indexes.

## 5. Validation

- [x] Run focused Go tests for retrieval strategy.
- [x] Run `go test ./...`.
- [x] Run `go vet ./...`.
- [x] Run `npm --prefix web test -- --run`.
- [x] Run `npm --prefix web run build`.
- [x] Run `openspec validate p109-retrieval-quality-strategy-optimization --strict`.
- [x] Run `openspec validate --all --strict`.
- [x] Run `git diff --check`.
- [x] Commit and push after validation passes.
