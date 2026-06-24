# P109 Retrieval Quality Strategy Optimization Acceptance

> Date: 2026-06-24  
> Change: `p109-retrieval-quality-strategy-optimization`  
> Conclusion: `retrieval_quality_strategy_optimized_with_deterministic_local_rerank`

## Scope

P109 improves retrieval quality after P108 sqlite-vec integration by adding deterministic local strategy around embeddings. It does not add new embedding providers, LLM reranking, external search services, trading actions, or final verdict changes.

SQLite `intelligence_summary` and `rag_chunks` remain authoritative. sqlite-vec remains a rebuildable auxiliary index.

## Implemented

| Area | Evidence |
| --- | --- |
| Query rewrite | Semantic index queries now include symbol, original question, inferred investment intent, and expected evidence categories such as announcement, risk, valuation, source verification, buy/sell discipline, portfolio state, and capital-flow terms. |
| Candidate widening | Semantic retrieval requests a wider candidate window before local rerank, using at least 8 candidates or 4x requested topK. |
| Metadata-aware rerank | Candidates are scored with vector rank, keyword overlap, symbol match, source level, evidence role, verification status, source count, relevance score, time weight, and indexed freshness. |
| Evidence diversity | Final topK deduplicates summary IDs and avoids event-type dominated results when alternative relevant evidence exists. |
| Audit traceability | `OutputRef` follows the reranked first evidence chunk instead of the raw first vector candidate. |
| Regression evaluation | Retrieval quality evaluation verifies reranked expected evidence can satisfy formal-only fixtures and still flags C/background evidence as diagnostic-only. |

## Validation

| Command | Result |
| --- | --- |
| `go test ./internal/application/service -run 'TestRetrievalAdapter(Rewrites|Widens|Hybrid|Diversifies)'` | pass |
| `go test ./internal/application/service` | pass |
| `go test ./internal/application/workflow -run 'TestEvidenceRetrieval|TestRetrieval'` | pass |
| `go test ./...` | pass |
| `go vet ./...` | pass |
| `npm --prefix web test -- --run` | pass; 49 files / 182 tests |
| `npm --prefix web run build` | pass |
| `openspec validate p109-retrieval-quality-strategy-optimization --strict` | pass |
| `openspec validate --all --strict` | pass; 35 items |
| `git diff --check` | pass |

The Go commands emit macOS CGO warnings from the upstream sqlite-vec binding about deprecated process-global SQLite extension APIs. The tests pass; the warning is not a product behavior failure.

## Boundaries

P109 does not claim improved investment returns, future market prediction accuracy, future provider availability, Docker validation, installer validation, GitHub Release success, or physical second-machine validation.

P109 does not add broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, paid/login/authorized data-source claims, Level2, high-frequency data, or return guarantees.
