## 1. Provider and index contract tests

- [x] Add tests for local JSON index health: missing, corrupted, incompatible, healthy.
- [x] Add tests for rebuild statistics: indexed count, skipped count, last rebuild time, degradation reason.
- [x] Add tests for retrieval fallback context when index search is empty or fails.

## 2. Local index hardening

- [x] Add versioned JSON index metadata/envelope while preserving rebuildability from SQLite.
- [x] Implement health inspection for configured local index path.
- [x] Ensure incompatible legacy or future index versions are reported and rebuildable.
- [x] Ensure corrupted index files do not prevent SQLite fallback.

## 3. Rebuild and retrieval observability

- [x] Return rebuild statistics from the application service.
- [x] Preserve degradation reason in retrieval results and audit context.
- [x] Keep C-level evidence as background-only during index and SQLite fallback paths.

## 4. API and documentation

- [x] Expose index health and rebuild statistics through existing service/API DTOs.
- [x] Update configuration or recovery documentation for local index rebuild and degraded states.
- [x] Document that real VecLite API integration remains future replacement work.

## 5. Validation

- [x] Run `openspec validate p13-index-and-retrieval-hardening --strict`.
- [x] Run `go test ./...`.
- [x] Run `cd web && npm run test && npm run build`.
- [x] Run `openspec validate --all --strict`.
- [x] Run relevant `cmd/agent` local task validation for index rebuild/health where available.
