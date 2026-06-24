# P108 Tasks

## 1. Scope And Dependency Review

- [x] Confirm no active OpenSpec change conflicts with P108.
- [x] Confirm sqlite-vec official Go CGO binding does not work with `modernc.org/sqlite` and therefore P108 uses a separate vector index database.
- [x] Confirm P108 does not modify final verdict, broker, trading, push, auto-confirm, or auto-rule boundaries.

## 2. Tests First

- [x] Add failing config tests for required embedding fields when `embedding.enabled=true`.
- [x] Add failing embedding client/provider tests for OpenAI-compatible `/embeddings` request/response handling.
- [x] Add failing sqlite-vec vector index tests for upsert, topK query ordering, health, and fallback behavior.
- [x] Add failing retrieval adapter tests proving sqlite-vec hits use semantic topK and unavailable vector search falls back to SQLite summaries.

## 3. Implementation

- [x] Add embedding config structs, validation, examples, and redacted diagnostics.
- [x] Add embedding provider interface and OpenAI-compatible embeddings client.
- [x] Add sqlite-vec vector index implementation using a separate local vector index SQLite file.
- [x] Wire embedding provider and sqlite-vec index into workflow dependencies.
- [x] Preserve FileVectorIndex compatibility only as fallback/test support where needed.

## 4. Documentation And Evidence

- [x] Update architecture/config/README docs to distinguish chat model config from embedding model config.
- [x] Add P108 acceptance record.
- [x] Update governance/progress/release indexes after archive.

## 5. Validation

- [x] Run focused red/green Go tests.
- [x] Run `go test ./...`.
- [x] Run `openspec validate p108-real-sqlite-vec-retrieval --strict`.
- [x] Run `openspec validate --all --strict`.
- [x] Run README/config local checks where changed.
- [x] Run `git diff --check`.
- [x] Archive P108 after validation passes.
