# Design

## Root Causes

1. Empty portfolio: `GET /api/v1/portfolio/current` can return `NOT_FOUND`, and the generic API client maps that to `generic_failure`, so `/positions` shows "读取失败" before onboarding guidance.
2. Decision detail density: `DecisionTrace` initializes analysis material as expanded, so real LLM output dominates the page.
3. Decision-loop deep link: `/decision-loop?decision_id=<id>` is ignored; the page always renders the full recent list with the latest item as the summary target.

## Approach

- Treat `NOT_FOUND` on portfolio read as `first_use` with onboarding-safe copy. Keep other `NOT_FOUND` contexts generic because not all missing records are onboarding.
- Keep the full analyst report content available, but default `showAnalysis` to collapsed and add a compact count/expand affordance.
- Read `decision_id` via `useSearchParams`, filter/focus the matching loop item, and show a clear focused-state message. If not found, keep the list empty with a safe message.
- Add Vitest coverage for all three behaviors and record a P103 linked-workflow acceptance report.

## Validation

- `openspec validate p103-product-acceptance-ux-linkage-fixes --strict`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `go test ./...`
- `openspec validate --all --strict`
- `python3 scripts/p92_final_requirement_audit.py --check`
- `python3 scripts/p93_code_reality_audit.py --check`
- `git diff --check`
