## 1. Frontend ops status tests

- [x] Add or strengthen frontend tests for data source/index/review status panels.
- [x] Cover success, degraded, failed, empty, and unknown states.
- [x] Assert ops surfaces use DTO/service data only and expose no local file or database access.

## 2. Review summary and tracking tests

- [x] Add or strengthen review summary tests for period, counts, degradation indicators, and rule suggestions.
- [x] Add or strengthen tests for tracking links to audit events, decisions, rule proposals, and error cases.
- [x] Assert tracking links do not show automatic trading or automatic rule application actions.

## 3. Implementation

- [x] Implement or refine frontend ops status panel components.
- [x] Implement or refine review summary display using existing DTOs where possible.
- [x] Implement or refine tracking entrypoints as safe navigation/filter links only.
- [x] Add minimal DTO or mapper adjustments only if required for visible state.
- [x] Preserve unknown status fallback behavior.

## 4. Documentation and OpenSpec sync

- [x] Keep P16 delta specs local before archive.
- [x] Document P16 boundaries: display and navigation only, no mutation shortcuts.

## 5. Validation

- [x] Run `openspec validate p16-frontend-ops-review-surface --strict`.
- [x] Run `cd web && npm run test && npm run build`.
- [x] Run `go test ./...` if DTO or backend code changes.
- [x] Run `openspec validate --all --strict`.
