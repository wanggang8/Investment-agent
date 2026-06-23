## Why

P16 improves the operator-facing frontend. P12-P15 added real-data sources, index health, local tasks, review summaries, and evidence-quality semantics, but the frontend still needs a clearer surface for operational state, review summaries, and cross-page tracking.

## What Changes

- Add or strengthen frontend surfaces for data source status, index health, review summary, and tracking entrypoints.
- Show empty, failed, degraded, and successful states distinctly.
- Preserve no-automatic-trading and no-automatic-rule-application boundaries in frontend flows.
- Keep backend API shape changes minimal and only add DTO fields if the frontend cannot represent existing facts.

## Capabilities

### New Capabilities
- `frontend-ops-review-surface`: Defines UI requirements for ops status, review summaries, and cross-page tracking.

### Modified Capabilities
- `frontend-experience-tests`: Extends frontend state and test coverage for P16 surfaces.
- `review-automation-delivery`: Ensures review summary display remains DTO-based and safe.

## Impact

- React pages/components for evidence, review, audit, settings, or dashboard status panels.
- Frontend tests for success/degraded/error/empty states.
- Optional DTO mapping adjustments for status and tracking links.
- No direct frontend access to SQLite, VecLite, local files, or automatic trade/order behavior.
