# P108 Real sqlite-vec Retrieval Acceptance

> Date: 2026-06-24  
> Change: `p108-real-sqlite-vec-retrieval`  
> Conclusion: `real_sqlite_vec_retrieval_implemented_with_embedding_config_and_fallback`

## Scope

P108 replaces the production retrieval path's JSON-only VecLite behavior with a real sqlite-vec auxiliary index when `embedding.enabled=true`.

SQLite remains the authoritative fact store for `intelligence_summary` and `rag_chunks`. The sqlite-vec file is rebuildable auxiliary data. The legacy FileVectorIndex/VecLite path remains available when embedding is disabled or for local fallback/test compatibility.

## Implemented

| Area | Evidence |
| --- | --- |
| Embedding config | Added `embedding.enabled`, `provider`, `api_key`, `base_url`, `model`, `dimensions`, and `timeout_seconds`; config validation requires OpenAI-compatible `/embeddings` fields when enabled. |
| Embedding client | Added OpenAI-compatible embeddings client that sends `{model,input}` to `/embeddings`, supports bearer auth, validates non-empty embeddings, and does not log vectors or keys. |
| sqlite-vec index | Added separate local sqlite-vec index file using `github.com/asg017/sqlite-vec-go-bindings/cgo` with `mattn/go-sqlite3`; primary business SQLite remains on `modernc.org/sqlite`. |
| Semantic retrieval | Retrieval adapter sends symbol/question text to semantic topK when sqlite-vec is enabled and falls back to existing SQLite summaries when vector search is unavailable. |
| Index rebuild | `evidence-index` now rebuilds through the configured index implementation so sqlite-vec can be rebuilt from authoritative SQLite `rag_chunks`. |
| Smoke command | Added `go run ./cmd/agent --task embedding-smoke --symbol 510300` to verify the embedding endpoint, model, auth, dimensions, and audit record. |
| Docs/config templates | Updated examples, env vars, README/architecture/config docs, Docker env placeholders, and local troubleshooting text. |

## Validation

| Command | Result |
| --- | --- |
| `go test ./internal/infrastructure/config -run 'TestValidate.*Embedding'` | pass |
| `go test ./internal/infrastructure/embedding/openai` | pass |
| `go test ./internal/application/service -run 'TestSQLiteVec|TestRetrievalAdapterUsesSemantic'` | pass |
| `go test ./internal/application/workflow -run TestEvidenceRetrievalPassesUserQuestionToSemanticRetrieval` | pass |
| `go test ./cmd/agent -run 'TestRunEmbeddingSmoke|TestRunHelpShowsLocalTasks'` | pass |
| `go test ./...` | pass |
| `go vet ./...` | pass |
| `npm --prefix web test -- --run` | pass |
| `npm --prefix web run build` | pass |
| `openspec validate p108-real-sqlite-vec-retrieval --strict` | pass |
| `openspec validate --all --strict` | pass |
| `git diff --check` | pass |
| `bash scripts/local-release-package.sh --release-label p108-sqlite-vec-smoke --output-dir tmp/p108-release-smoke-final` | pass; archive SHA `67e342d483400cc15a2d2e0ba0a6aa5d1aab11162ff4f82c704510919c320834` |
| `bash scripts/local-release-package.sh --verify tmp/p108-release-smoke-final/20260624T062716Z/investment-agent-p108-sqlite-vec-smoke.tar.gz` | pass; 1779 entries, 0 errors, 0 warnings |

The Go commands emit macOS CGO warnings from the upstream sqlite-vec binding about deprecated process-global SQLite extension APIs. The tests pass; the warning is not a product behavior failure.

## Boundaries

P108 does not add broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, paid/login/authorized data-source claims, Level2, high-frequency data, or return guarantees.

P108 does not claim Docker validation or physical second-machine validation. Docker-specific work remains out of scope for this change.

Chat/analysis models such as `deepseek.model` or GPT chat models are not accepted as substitutes for embedding tests. Real semantic retrieval requires an embedding model endpoint that supports OpenAI-compatible `/embeddings`.
