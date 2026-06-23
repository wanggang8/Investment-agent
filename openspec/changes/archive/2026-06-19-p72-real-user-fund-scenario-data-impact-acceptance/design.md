# Design: P72 Real User Fund Scenario Data Impact Acceptance

## Context

P71 established that the product can pass a strict real acceptance gate. It intentionally focused on release blockers: current data, VecLite, real UI, real LLM, package verify, and repeat acceptance. P72 focuses on real use: does a normal fund/ETF workflow change the intended local facts, preserve auditability, and keep the UI consistent after refresh or navigation?

## Acceptance Model

P72 success status is `real_user_scenario_data_impact_passed`.

Blocking statuses:

| Outcome | Meaning |
| --- | --- |
| `blocked_scenario_coverage` | Scenario matrix is incomplete or unreviewed before execution. |
| `blocked_ui_operation` | Real UI operation fails, page errors occur, unexpected API failures occur, or data does not re-render after refresh/navigation. |
| `blocked_data_impact` | SQLite side effects are missing, inconsistent, duplicated unexpectedly, or not linked by audit/decision/report/risk records. |
| `blocked_accuracy` | Deterministic calculations such as market value, P/L ratio, cash, total assets, risk trigger, or count/link consistency fail. |
| `blocked_real_provider` | Required real current-data or real LLM provider is unavailable. |
| `blocked_safety` | Forbidden trading, broker, auto-action, raw secret, raw prompt, raw payload, or return-promise wording appears. |

## Scenario Matrix

The P72 matrix is the execution contract:

| ID | Scenario | UI operation | Data impact to verify |
| --- | --- | --- | --- |
| S1 | Real fund setup | Calibrate `510300` portfolio with cash, total assets, quantity, cost, current price, buy reason, asset tag | `portfolio_snapshots`, `positions`, `audit_events`; market value = quantity × current price; P/L ratio = `(current-cost)/cost`; page refresh remains consistent |
| S2 | Holding maintenance | Edit holding, validate batch import, confirm import, record correction | current position updated as expected; import/correction records exist; audit links written; no broker/order data appears |
| S3 | Offline transaction | Record a user-completed offline transaction | `position_transactions`, optional current `positions`, audit event, no automatic trade/order affordance |
| S4 | Local knowledge and retrieval | Validate/confirm local note for `510300`, rebuild VecLite | `intelligence_items`, `intelligence_summary`, `rag_chunks`, `source_verifications`, `audit_events`; VecLite health healthy; retrieval quality hit |
| S5 | Current data and market refresh | Run current-data refresh/gate and UI market refresh path | `source_health` facts and quality gate pass; market refresh audit/notification behavior is safe |
| S6 | Daily discipline and risk alerts | Generate daily discipline after high valuation seed/refresh | `daily_discipline_reports`, `decision_records`, `risk_alerts`, `notifications`, `audit_events`; risk trigger matches deterministic valuation threshold |
| S7 | Real consultation and decision detail | Submit real LLM consultation for `510300`, open generated detail | decision workflow completed; analyst reports parsed/passed; final verdict from rules; retrieval from healthy VecLite; no return guarantee |
| S8 | Manual confirmation and loop | Record manual follow-up from generated decision | `operation_confirmations`, `position_transactions`, `decision_records.confirmation_status`, decision-loop UI, audit link |
| S9 | Review/governance/readback | Visit review, audit, notifications, rules, data quality, workbench | derived counts/links reflect scenario; rule proposals are not auto-applied; notifications can be read without external push |
| S10 | Failure/safety boundaries | Trigger representative validation errors and scan visible text | invalid inputs fail safely; no secret/prompt/raw payload/private path; no auto-trade/broker/return promise wording |

## Execution Architecture

P72 will add:

- `scripts/p72-real-user-fund-scenario-acceptance.sh`
  - Builds a temp config from local real LLM config.
  - Uses a temp SQLite DB and temp VecLite file.
  - Seeds the normal smoke dataset.
  - Runs P34 current-data refresh and strict current-data gate for `000300`.
  - Starts the local Go server and Vite frontend.
  - Runs `web/e2e/p72-real-user-fund-scenario.spec.ts`.
  - Runs a Python read-only SQLite impact checker and writes `db-impact-summary.json`.
- `web/e2e/p72-real-user-fund-scenario.spec.ts`
  - Drives the actual UI and records API/UI results to `browser-results.json`.
  - Does not mock browser network calls.
  - Blocks on page errors, unexpected API failures, overflow, missing expected UI state, LLM degradation, retrieval degradation, or forbidden affordances.

P72 does not add a product API for internal data-impact checks. SQLite inspection stays inside the local acceptance script and writes sanitized aggregate evidence only.

## Accuracy Boundaries

P72 can strictly verify deterministic product logic:

- market value,
- unrealized profit ratio,
- cash / total assets / position count,
- table row counts and foreign-key-like links,
- risk threshold trigger presence,
- report/decision/notification/audit consistency,
- retrieval health/freshness,
- LLM parse/quality status.

P72 cannot claim predictive investment accuracy. It can only verify that analysis uses real inputs, cites traceable evidence, preserves rule-based final verdict, and avoids return promises.

## Review Gates

Before execution:

1. Review the matrix against all primary user workflows.
2. Confirm every scenario has UI operation, API/DB impact, and safety expectations.
3. Confirm the plan does not rely on mock-only pass evidence.

After execution:

1. Compare browser results and DB impact summary against the matrix.
2. Record any missing scenario as a gap, not as a pass.
3. Add follow-up scenarios only if they are necessary for practical user confidence and can be executed without changing the product boundary.

## Evidence

P72 writes sanitized artifacts under:

```text
docs/release/ui-audit-assets/2026-06-18-p72/
docs/release/acceptance/2026-06-18-p72-real-user-fund-scenario.md
```

Temporary SQLite databases, VecLite files, logs, raw provider payloads, full prompts, and private configs stay under `tmp/` and are not committed.
