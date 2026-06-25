# P116 Multi-Fund Transaction Ledger Acceptance Matrix

> Date: 2026-06-25  
> Change: `p116-multi-fund-transaction-ledger-acceptance`  
> Status: passed by isolated P116 runner  
> Boundary: local ledger acceptance only; no broker integration, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, release package refresh, physical second-machine validation, external provider guarantee, or fresh real LLM claim.

| ID | Scenario | Main Surfaces | Required Operations | Required Evidence | Expected Eligibility | Actual Status |
| --- | --- | --- | --- | --- | --- | --- |
| L01 | Fresh runtime and empty ledger | `/positions`, health API | Start temp backend/frontend/SQLite, verify empty current portfolio behavior | Health API, empty portfolio 404, no fake rows | `fresh_pass` | `fresh_pass` |
| L02 | Multi-fund initial portfolio | `/positions` | Initialize `510300`, `159915`, `588000`, `512000`, `110022` with core/satellite/fund tags | Portfolio API, SQLite positions, browser table | `fresh_pass` | `fresh_pass` |
| L03 | Multi-date offline ledger | `/positions` | Record buy/sell/reduce operations with fees and dates | Transaction APIs, SQLite `position_transactions`, cash/position readback | `fresh_pass` | `fresh_pass` |
| L04 | Mixed batch import | `/positions` | Validate mixed holding/transaction rows with valid and invalid rows, then confirm valid subset | Validate summary, confirm API, SQLite delta, invalid-row evidence | `fresh_pass` | `fresh_pass` |
| L05 | Invalid transaction rejection | `/positions`, API | Reject cash insufficient, oversell, future time, negative fees, missing symbol, invalid state | 400 responses, unchanged transaction count, no partial writes | `fresh_pass` | `fresh_pass` |
| L06 | Edit/remove/correction | `/positions`, `/audit` | Edit one holding, remove one holding, record correction | API, SQLite position/correction/audit readback | `fresh_pass` | `fresh_pass` |
| L07 | Decision manual execution | `/decisions/:id`, `/decision-loop` | Record executed_manually confirmation for seeded decision | Browser/API confirmation, SQLite confirmation + transaction, loop readback | `fresh_pass` | `fresh_pass` |
| L08 | Marked-error decision loop | `/decisions/:id`, `/review`, `/audit` | Mark another decision as `marked_error` | Confirmation API, SQLite marked_error, review/audit readback | `fresh_pass` | `fresh_pass` |
| L09 | Quarterly rebalance | `/positions`, `/review` | Run rebalance across core/satellite/cash buckets | Rebalance API, audit event, no automatic adjustment | `fresh_pass` | `fresh_pass` |
| L10 | Risk and notification linkage | `/risk-alerts`, `/notifications` | Resolve seeded risk alert, mark notification read/all-read | Risk API, notification API, SQLite readback | `fresh_pass` | `fresh_pass` |
| L11 | Data-quality gate lifecycle | `/data-quality` | Create and retire gate resolution for a symbol | DQ APIs, resolution row retired, no current-data clean overclaim | `scoped_pass` | `scoped_pass` |
| L12 | Aggregate readback | `/`, `/workbench`, `/review`, `/audit` | Reopen aggregate pages after all writes | Dashboard/review/audit APIs and browser DOM | `scoped_pass` | `scoped_pass` |
| L13 | Browser multi-fund positions | `/positions` | Browser initialize/import/offline/rebalance, inspect table | Screenshots, DOM assertions, console health | `fresh_pass` | `fresh_pass` |
| L14 | Browser decision confirmation | `/decisions/:id` | Browser manual execution confirmation | Screenshot, confirmation success, no auto action | `fresh_pass` | `fresh_pass` |
| L15 | Mobile portfolio rendering | 390px `/positions` | Verify multi-fund portfolio on mobile | Mobile screenshot, no console/page errors | `fresh_pass` | `fresh_pass` |
| L16 | Safety negative evidence | all core routes | Check forbidden tables/actions/claims/secrets | SQLite safety counters, UI text scan | `fresh_pass` | `fresh_pass` |

## Required Safety Counters

- `forbidden_broker_order_push_tables = 0`.
- `auto_confirmation_rows = 0`.
- `auto_rule_apply_audit_events = 0`.
- `automatic_trading_affordances = 0`.
- `return_guarantee_claims = 0`.
- `secret_or_raw_prompt_leaks_on_primary_ui = 0`.

## Final Runner Summary

- Runner: `bash scripts/p116-multi-fund-transaction-ledger-acceptance.sh`.
- Final JSON: `docs/release/ui-audit-assets/2026-06-25-p116-multi-fund-transaction-ledger-acceptance/p116-scenario-summary.json`.
- Scenario count: 16.
- Fresh pass: 14.
- Scoped pass: 2.
- Browser console errors: 0.
- Browser page errors: 0.
- Browser API 5xx responses: 0.
