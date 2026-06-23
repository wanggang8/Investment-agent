## 1. Evidence quality tests

- [x] Add or strengthen tests for time weight and relevance score preservation through normalization, retrieval, DTOs, and decision refs.
- [x] Add or strengthen tests that C-level sources are background-only across refresh, retrieval, and decision evidence paths.
- [x] Add or strengthen tests that fewer than two S/A independent sources cannot produce satisfied major-event evidence.

## 2. Structured facts vs analyst materials

- [x] Add tests showing rule arbitration uses structured facts and source verification status.
- [x] Add tests showing analyst materials cannot override the final verdict.
- [x] Keep expected return materials explanatory and non-decisive.

## 3. Implementation

- [x] Preserve quality metadata in evidence normalization and retrieval paths.
- [x] Preserve quality metadata in decision evidence refs and response DTOs.
- [x] Ensure source-role restrictions are enforced before final rule arbitration.
- [x] Keep entity extraction, expanded event classification, and complex relevance scoring out of scope.

## 4. Documentation and OpenSpec sync

- [x] Update change-local specs only before archive.
- [x] Document P15 boundaries: required quality fields now, optional enrichment later.

## 5. Validation

- [x] Run `openspec validate p15-evidence-quality-enrichment --strict`.
- [x] Run `go test ./...`.
- [x] Run `cd web && npm run test && npm run build`.
- [x] Run `openspec validate --all --strict`.
- [x] Run relevant `cmd/agent` local task validation where applicable.
