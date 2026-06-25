# P117 Continuous Product Usability Acceptance

> Date: 2026-06-25  
> Change: `p117-continuous-product-usability-acceptance`  
> Result: `passed_local_continuous_usability_acceptance`

## Summary

P117 executed a seven-day local continuous-use story with isolated backend, Vite frontend, temporary SQLite, backend restart and Playwright browser evidence. The goal was to answer whether the product is usable over time as a local investment discipline assistant, not merely whether isolated pages or APIs respond.

Final evidence:

- `docs/release/ui-audit-assets/2026-06-25-p117-continuous-product-usability-acceptance/p117-usability-summary.json`
- `docs/release/ui-audit-assets/2026-06-25-p117-continuous-product-usability-acceptance/api_sqlite/p117-api-sqlite-summary.json`
- `docs/release/ui-audit-assets/2026-06-25-p117-continuous-product-usability-acceptance/restart/p117-restart-summary.json`
- `docs/release/ui-audit-assets/2026-06-25-p117-continuous-product-usability-acceptance/browser/p117-browser-results.json`

## Result

- Runner scenarios: 17/17 passed.
- Fresh pass: 16.
- Scoped pass: 1 (`U10` data-quality seeded degradation handling).
- Blocked: 0.
- Browser status: passed.
- Browser console errors: 0.
- Browser page errors: 0.
- Browser API 5xx responses: 0.
- Restart probe: passed.

## Usability Interpretation

- Cold start: usable. Empty portfolio is explicit and does not fabricate holdings.
- Onboarding: usable. Local account and holdings can be created without broker connection.
- Daily routine: usable. Dashboard, workbench, review and audit remain readable after local facts.
- Recovery: usable. Invalid imports and invalid transactions are rejected without partial ledger writes.
- Traceability: usable. Corrections, manual confirmations, marked errors and audit rows provide a review trail.
- Persistence: usable. Restart probe reads portfolio, audit and decision-loop facts from the same SQLite database.
- Mobile: usable within checked scope. 390px portfolio/workbench paths render without console/page/API 5xx failures.
- Safety: usable only as a local discipline assistant. No broker/order/push tables, auto confirmations or auto rule-apply audit events were found.

## Seven-Day Coverage

- Day 0: cold start and empty local facts.
- Day 1: account onboarding, holdings, first dashboard/workbench explanation.
- Day 2: daily routine readback across dashboard, review, reports and audit.
- Day 3: multi-fund offline buy/sell/reduce ledger updates, risk and notification handling.
- Day 4: invalid import recovery, invalid transaction recovery and correction audit.
- Day 5: seeded data-quality degradation, explicit resolution and retirement.
- Day 6: manual decision execution confirmation and marked-error learning loop.
- Day 7: cross-page consistency, backend restart persistence, mobile usability and safety negative evidence.

## Safety Boundary

P117 is local continuous-use usability acceptance only. It does not claim broker integration, real exchange execution, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, fresh external provider quality, fresh real LLM output, release package refresh or physical second-machine validation.

Safety counters:

- Forbidden broker/order/push tables: 0.
- Auto confirmation rows: 0.
- Auto rule-apply audit events: 0.
- Automatic trading affordances: 0.
- Return guarantee claims: 0.
- Secret/raw prompt leaks: 0.

## Notes

P117 intentionally treats cold start as API/SQLite evidence because the browser phase runs after the seven-day story has accumulated local facts. The browser phase instead validates user-visible continuity and cross-page readback after those facts exist.

During regression, one full `npm --prefix web test -- --run` attempt produced a transient `DataQualityPage` assertion failure. The targeted test, the whole `DataQualityPage.test.tsx` file and the full Vitest suite all passed on immediate rerun; no code change was made for that non-reproducible symptom.

P93 remains a potentially stale historical audit after later P114/P115/P116/P117 work; P117 must not be cited as fresh P93 code-reality pass.

## Verification

- `bash scripts/p117-continuous-product-usability-acceptance.sh`: passed.
- `openspec validate p117-continuous-product-usability-acceptance --strict`: passed.
- `openspec validate --all --strict`: passed, 38 items.
- `go test ./...`: passed, with sqlite-vec macOS deprecation warnings.
- `go vet ./...`: passed, with sqlite-vec macOS deprecation warnings.
- `npm --prefix web test -- --run`: passed on rerun, 53 files / 191 tests.
- `npm --prefix web run build`: passed.
- `python3 scripts/p92_final_requirement_audit.py --check`: passed.
- `python3 scripts/p93_code_reality_audit.py --check`: stale, `docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`.
- `git diff --check`: passed.
