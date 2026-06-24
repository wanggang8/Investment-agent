# P104 Product Operation Linkage Acceptance Matrix

Date: 2026-06-24  
Scope: local source, isolated SQLite, local HTTP backend  
Change: `p104-full-product-operation-linkage-acceptance`

## Verdict Model

P104 treats a product operation as accepted only when the expected user action, API result, durable data effect, downstream readback, and safety boundary all agree.

Coverage labels:

- `runner`: covered by `scripts/p104-product-operation-linkage-acceptance.sh`.
- `focused_tests`: covered by existing Go or frontend unit/component tests in this source tree.
- `ledger`: covered by P92/P93 original-requirement and code-reality ledgers.

## Matrix

| Product surface | Operation class | Expected API/UI behavior | Expected SQLite effect | Downstream linkage | Safety boundary | P104 coverage |
| --- | --- | --- | --- | --- | --- | --- |
| Portfolio onboarding/current | Empty or initialized portfolio readback | `/portfolio/current` returns a coherent first-use or current-state payload | `portfolio_snapshots` and `positions` agree on cash, market value, ratios and count | Dashboard/workbench/review can read the current account state | No broker sync or account login claim | `runner`, `focused_tests`, `ledger` |
| Portfolio adjustment | Manual account calibration accepts consistent totals and rejects inconsistent totals | Valid `/portfolio/adjustments` returns latest snapshot; invalid total returns 400 | `portfolio_snapshots`, `position_snapshots`, `positions`, `audit_events` update atomically | `/portfolio/current`, `/review/summary`, `/audit-events` read the same operation | Local fact record only; no trade/order side effect | `runner`, `focused_tests`, `ledger` |
| Holding maintenance | Edit/remove a local holding | `/portfolio/holdings` updates current holding; remove clears selected holding | `positions`, `position_snapshots`, `portfolio_snapshots`, `audit_events` reflect edit/remove | Portfolio and review endpoints see changed holdings | Requires explicit local confirmation text; no automatic correction | `runner`, `focused_tests`, `ledger` |
| Offline transaction | Record a manually executed external action | `/portfolio/offline-transactions` records buy/sell note and safety text | `position_transactions`, `positions`, `portfolio_snapshots`, `audit_events` update | Portfolio/review/audit expose the manual fact | Does not connect to broker or place an order | `runner`, `focused_tests`, `ledger` |
| Batch import | Validate before write, confirm after review | Validate returns valid/invalid row counts; confirm writes only confirmed rows | `local_account_import_batches`, `positions` or `position_transactions`, `audit_events` update | Portfolio/current and audit show imported facts | Validation alone must not mutate holdings | `runner`, `focused_tests`, `ledger` |
| Correction audit | Record a manual correction | `/portfolio/corrections` stores before/after context | `local_account_corrections`, `audit_events` update | Audit/review can trace the correction | Correction is local audit, not automatic repair | `runner`, `focused_tests`, `ledger` |
| Rebalance review | Evaluate allocation drift | `/portfolio/rebalance-review` returns drift and action prompts | Review facts are computed from latest snapshot/positions | Workbench/portfolio can guide next manual action | Suggestion only; no automatic rebalance | `runner`, `focused_tests`, `ledger` |
| Decision detail | Read a formal decision | `/decisions/{id}` returns verdict, evidence, reports, expected-return context and user confirmation state | `decision_records` remains source of truth | `/decision-loop`, `/review`, `/audit-events` can refer to the same decision | LLM material is context only; final verdict remains rule/user controlled | `runner`, `focused_tests`, `ledger` |
| Manual confirmation | Planned or executed-manually user decision | `/decisions/{id}/confirmations` records explicit user action | `operation_confirmations`, optional `position_transactions`, `decision_records`, `audit_events` update | Decision detail, decision loop, review and audit agree on status | No automatic confirmation rows allowed | `runner`, `focused_tests`, `ledger` |
| Decision loop | Explain closed-loop state for one decision | `/decision-loops/{id}` and query filtered list return the focused item | Reads decision, confirmation, transaction, risk/review/audit facts | Decision detail -> loop -> review/audit navigation is coherent | Read-only explanation, no side effects | `runner`, `focused_tests`, `ledger` |
| Review summary | Aggregate recent operations | `/review/summary` reflects confirmations, errors, rules, risk and audit counts | Reads existing durable facts without hidden writes | Dashboard/workbench governance cards align with summary | No performance or return guarantee | `runner`, `focused_tests`, `ledger` |
| Audit | Inspect operation trace | `/audit-events` lists recent user/system events | `audit_events` includes request/action/input references | Deep links to related decisions/notifications remain usable | Sensitive material must stay summarized/redacted | `runner`, `focused_tests`, `ledger` |
| Notifications | Mark one or all read | `/notifications/{id}/read` and `/read-all` update read state | `notifications.read_at` is set | Notification center/read counts update | In-app only; no external push | `runner`, `focused_tests`, `ledger` |
| Risk alerts | Resolve a local SOP alert | `/risk-alerts/{id}/lifecycle` changes SOP status | `risk_alerts.sop_status`, `resolved_at`, `resolution_reason` update | Risk list/detail/review/audit can read the state | Resolution does not create trade confirmations | `runner`, `focused_tests`, `ledger` |
| Data-quality gate | Record and retire local release-claim resolution | `/data-source-quality/resolutions` creates manual scope/waiver record; retire clears it | `data_quality_gate_resolutions`, `audit_events` update | Data-quality page and release claim check reflect state | Does not refresh data or convert degraded data into clean data | `runner`, `focused_tests`, `ledger` |
| Settings/data refresh | Runtime provider configuration and manual refresh | Settings page/API expose configured provider and refresh status | `settings`/market facts read back as applicable | Data-quality and market snapshot endpoints reflect data state | No paid/login/Level2/high-frequency source claim | `focused_tests`, `ledger` |
| Rules/governance | Rule proposals and final confirmation | Rules UI/API list proposals, gatekeeper audits and final explicit confirmation | `rule_proposals`, `gatekeeper_audits`, `audit_events` update | Review/audit/workbench show governance state | No automatic rule application | `focused_tests`, `ledger` |
| Local knowledge | Validate and confirm local knowledge import | Validate previews sanitized facts; confirm writes background facts/index plan | Local knowledge and evidence tables update | Retrieval/decision context can cite sanitized background | No raw secret/path/prompt exposure | `focused_tests`, `ledger` |
| Daily discipline/auto run | Manual report and disabled-by-default scheduler state | Daily report and auto-run pages show status and diagnostics | Daily run/report tables and notifications update on explicit run | Dashboard/workbench/notifications read the report | Auto-run disabled by default; no trading side effects | `focused_tests`, `ledger` |

## Pass Criteria

P104 passes when:

1. The runner completes with `status=passed`.
2. For each runner-covered operation, API readback and SQLite readback agree.
3. Downstream endpoints return coherent linked state after the write operations.
4. Safety negative checks pass: no broker/order/external-push/trade-execution tables, no automatic confirmation rows, and no unexpected auto-rule/trading evidence.
5. Existing regression gates pass.

## Boundary

P104 does not replace P92/P93. It adds fresh product-linkage evidence on top of the final requirement/code-reality ledgers. It also does not claim Docker, installer, package refresh, GitHub Release, physical second-machine validation, broker integration, automatic trading, automatic confirmation, automatic rule application, or investment return guarantees.
