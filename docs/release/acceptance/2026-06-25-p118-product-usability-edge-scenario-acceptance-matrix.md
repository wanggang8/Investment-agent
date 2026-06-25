# P118 Product Usability Edge Scenario Acceptance Matrix

> Date: 2026-06-25  
> Change: `p118-product-usability-edge-scenario-acceptance`  
> Status: passed; pending archive only  
> Boundary: local product-use edge scenario acceptance only; excludes Docker/install/upgrade/uninstall/release package/GitHub Release/physical second-machine validation. No broker integration, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, external provider guarantee, fresh real LLM claim or return guarantee.

| ID | Group | Scenario | Main Surfaces | Required Operations | Required Evidence | Expected Eligibility | Actual Status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| E01 | Long-cycle | 30-day accumulated local use | daily reports, audit, notifications | Seed/read 30 reports and many local facts | SQLite/API counts and browser list readback | `fresh_pass` | `fresh_pass` |
| E02 | Long-cycle | Long transaction history | `/positions` | Record/import many multi-fund transactions | Transaction API and SQLite counts | `fresh_pass` | `fresh_pass` |
| E03 | Long-cycle | Long audit/history surfaces | `/audit`, `/review` | Read accumulated audit, review and reports | API/browser evidence, no 5xx | `fresh_pass` | `fresh_pass` |
| E04 | Recovery | Invalid mixed import recovery | `/positions` | Validate invalid import and reject confirm | HTTP 400, unchanged transaction count | `fresh_pass` | `fresh_pass` |
| E05 | Recovery | Illegal transaction recovery | API, `/positions` | Reject oversell, future date, negative fee and missing symbol | HTTP 400, unchanged transaction count | `fresh_pass` | `fresh_pass` |
| E06 | Recovery | Duplicate/conflict transaction traceability | `/positions`, `/audit` | Record correction after suspicious duplicate-like fact | Correction API and audit readback | `fresh_pass` | `fresh_pass` |
| E07 | Data quality | Stale/missing data degradation | `/data-quality` | Seed stale source health, create scoped resolution | DQ API, SQLite resolution | `scoped_pass` | `scoped_pass` |
| E08 | Data quality | Degradation retirement | `/data-quality`, `/audit` | Retire resolution after local review | DQ retire API, audit evidence | `scoped_pass` | `scoped_pass` |
| E09 | Decision quality | Rising-context decision interpretation | `/decisions/:id` | Read seeded rising-context decision | Final verdict/evidence/prohibited actions | `fresh_pass` | `fresh_pass` |
| E10 | Decision quality | Falling-context decision interpretation | `/decisions/:id` | Read seeded falling-context decision | Different verdict and warning chain | `fresh_pass` | `fresh_pass` |
| E11 | Decision quality | Volatile/insufficient evidence interpretation | `/decisions/:id`, `/decision-loop` | Read seeded volatile decision and loop | Frozen/insufficient-data boundary | `fresh_pass` | `fresh_pass` |
| E12 | Household ledger | Multiple local account tags | `/positions` | Maintain account tags across several positions | SQLite asset/account tag evidence | `fresh_pass` | `fresh_pass` |
| E13 | Household ledger | Cash/core/satellite/active-fund consistency | `/positions`, `/`, `/workbench` | Read portfolio allocation after many facts | API/browser consistency evidence | `fresh_pass` | `fresh_pass` |
| E14 | Household ledger | Manual confirmation and marked error | `/decisions/:id`, `/review` | Execute manually and mark another decision wrong | SQLite confirmation/error cases | `fresh_pass` | `fresh_pass` |
| E15 | Cross-page | Accumulated-state cross-page readback | `/positions`, `/`, `/workbench`, `/review`, `/audit`, `/risk-alerts`, `/notifications`, `/daily-discipline/reports`, `/decision-loop` | Open all accumulated primary pages | Browser screenshots, console/page/API 5xx = 0 | `fresh_pass` | `fresh_pass` |
| E16 | Persistence | Restart persistence after accumulated history | backend restart | Restart backend with same SQLite and reread facts | Restart probe summary | `fresh_pass` | `fresh_pass` |
| E17 | Mobile | 390px accumulated-state usability | 390px `/positions`, `/workbench`, `/decision-loop` | Verify core mobile pages after history accumulation | Screenshots, no page/API failures | `fresh_pass` | `fresh_pass` |
| E18 | Safety | Safety negative evidence | all core routes and SQLite | Check forbidden tables/actions/claims/secrets | SQLite safety counters and UI text scan | `fresh_pass` | `fresh_pass` |
| E19 | Regression | Regression gates | local commands | Run P118/OpenSpec/Go/frontend/P92/P93/diff gates | Command output summary and stale boundary | `fresh_pass` | `fresh_pass_with_p93_stale_boundary` |

## Required Safety Counters

- `forbidden_broker_order_push_tables = 0`.
- `auto_confirmation_rows = 0`.
- `auto_rule_apply_audit_events = 0`.
- `automatic_trading_affordances = 0`.
- `return_guarantee_claims = 0`.
- `secret_or_raw_prompt_leaks_on_primary_ui = 0`.

## Runner Summary

- Runner: `bash scripts/p118-product-usability-edge-scenario-acceptance.sh`.
- Final JSON: `docs/release/ui-audit-assets/2026-06-25-p118-product-usability-edge-scenario-acceptance/p118-edge-usability-summary.json`.
- Runner scenarios: 18.
- Fresh pass: 16.
- Scoped pass: 2.
- Browser console errors: 0.
- Browser page errors: 0.
- Browser API 5xx responses: 0.
- Restart probe: passed.

## Regression Summary

- `bash scripts/p118-product-usability-edge-scenario-acceptance.sh`: passed.
- `openspec validate p118-product-usability-edge-scenario-acceptance --strict`: passed.
- `openspec validate --all --strict`: passed, 39 items.
- `go test ./...`: passed, with sqlite-vec macOS deprecation warnings.
- `go vet ./...`: passed, with sqlite-vec macOS deprecation warnings.
- `npm --prefix web test -- --run`: passed, 53 files / 191 tests.
- `npm --prefix web run build`: passed.
- `python3 scripts/p92_final_requirement_audit.py --check`: passed.
- `python3 scripts/p93_code_reality_audit.py --check`: stale, `docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md`.
- `git diff --check`: passed.
