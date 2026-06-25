# P117 Continuous Product Usability Acceptance Matrix

> Date: 2026-06-25  
> Change: `p117-continuous-product-usability-acceptance`  
> Status: passed; pending archive only  
> Boundary: local continuous-use acceptance only; no broker integration, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, release package refresh, physical second-machine validation, external provider guarantee, or fresh real LLM claim.

| ID | Day | Scenario | Main Surfaces | Required Operations | Required Evidence | Expected Eligibility | Actual Status |
| --- | --- | --- | --- | --- | --- | --- | --- |
| U01 | Day 0 | Cold start and empty facts | `/positions`, health API | Start temp backend/frontend/SQLite, verify empty current portfolio | Health API, empty portfolio 404, no fabricated facts | `fresh_pass` | `fresh_pass` |
| U02 | Day 1 | First account onboarding | `/positions` | Record cash, total assets and multi-fund holdings | Portfolio API, SQLite positions, browser screenshot | `fresh_pass` | `fresh_pass` |
| U03 | Day 1 | First explanation and next step | `/`, `/workbench` | Open primary daily surfaces after onboarding | Dashboard/workbench API/browser readback | `fresh_pass` | `fresh_pass` |
| U04 | Day 2 | Daily routine readback | `/`, `/workbench`, `/review`, `/audit` | Check daily status, review and audit after local facts | API/browser evidence and no page errors | `fresh_pass` | `fresh_pass` |
| U05 | Day 3 | Offline transaction update | `/positions` | Record buy/sell/reduce across several funds | Transaction API, SQLite transaction count/symbols | `fresh_pass` | `fresh_pass` |
| U06 | Day 3 | Risk and notification action | `/risk-alerts`, `/notifications` | Resolve risk and mark notification read | Risk/notification API, SQLite readback | `fresh_pass` | `fresh_pass` |
| U07 | Day 4 | Invalid import recovery | `/positions` | Validate mixed invalid import and reject confirm | 400 response, unchanged transaction count | `fresh_pass` | `fresh_pass` |
| U08 | Day 4 | Invalid transaction recovery | API, `/positions` | Reject oversell/future/negative fee/missing symbol | 400 responses, no partial writes | `fresh_pass` | `fresh_pass` |
| U09 | Day 4 | Local correction audit | `/positions`, `/audit` | Record correction after user notices wrong fact | Correction API, SQLite correction/audit | `fresh_pass` | `fresh_pass` |
| U10 | Day 5 | Data-quality degradation resolution | `/data-quality` | Create and retire gate resolution | DQ API, SQLite retired resolution, scoped claim boundary | `scoped_pass` | `scoped_pass` |
| U11 | Day 6 | Manual decision execution | `/decisions/:id` | Record executed_manually confirmation | Browser/API confirmation, SQLite transaction/confirmation | `fresh_pass` | `fresh_pass` |
| U12 | Day 6 | Marked-error learning loop | `/decisions/:id`, `/review` | Mark decision wrong and record lesson | API/SQLite marked_error, review/audit readback | `fresh_pass` | `fresh_pass` |
| U13 | Day 7 | Cross-page consistency | `/positions`, `/`, `/workbench`, `/review`, `/audit`, `/decision-loop` | Reopen aggregate pages after seven-day facts | Browser screenshots, API consistency counters | `fresh_pass` | `fresh_pass` |
| U14 | Day 7 | Restart persistence | backend restart | Restart backend with same SQLite, read previous data | Restart probe summary, current portfolio and audit counts | `fresh_pass` | `fresh_pass` |
| U15 | Day 7 | Mobile usability | 390px `/positions`, `/workbench` | Verify primary mobile path after accumulated facts | Screenshots, no console/page/API 5xx | `fresh_pass` | `fresh_pass` |
| U16 | Day 7 | Safety negative evidence | all core routes | Check forbidden tables/actions/claims/secrets | SQLite safety counters and UI text scan | `fresh_pass` | `fresh_pass` |
| U17 | Day 7 | Usability interpretation report | acceptance report | Summarize completion, recoverability, consistency and trust boundaries | Final JSON interpretation and Markdown record | `fresh_pass` | `fresh_pass` |
| U18 | Day 7 | Regression gates | local commands | Run OpenSpec/Go/frontend/P92/P93/diff gates | Command output summary and stale boundary | `fresh_pass` | `fresh_pass_with_p93_stale_boundary` |

## Required Safety Counters

- `forbidden_broker_order_push_tables = 0`.
- `auto_confirmation_rows = 0`.
- `auto_rule_apply_audit_events = 0`.
- `automatic_trading_affordances = 0`.
- `return_guarantee_claims = 0`.
- `secret_or_raw_prompt_leaks_on_primary_ui = 0`.

## Runner Summary

- Runner: `bash scripts/p117-continuous-product-usability-acceptance.sh`.
- Final JSON: `docs/release/ui-audit-assets/2026-06-25-p117-continuous-product-usability-acceptance/p117-usability-summary.json`.
- Runner scenarios: 17.
- Fresh pass: 16.
- Scoped pass: 1.
- Browser console errors: 0.
- Browser page errors: 0.
- Browser API 5xx responses: 0.
- Restart probe: passed.

## Regression Summary

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
