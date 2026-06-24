# P104 Full Product Operation Linkage Acceptance

Date: 2026-06-24  
Change: `p104-full-product-operation-linkage-acceptance`  
Verdict: `local_source_product_operation_linkage_acceptance_passed`

## Scope

P104 establishes a repeatable local-source product operation/linkage gate. It validates representative user operations through HTTP APIs, SQLite durable effects, downstream readback, audit traceability, and safety negative evidence.

P104 does not validate Docker, installer scripts, GitHub Release, package refresh, physical second-machine execution, remote deployment, broker interfaces, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, or investment return guarantees.

## Fresh Runner Evidence

Command:

```bash
bash scripts/p104-product-operation-linkage-acceptance.sh
```

Artifact:

- `docs/release/ui-audit-assets/2026-06-24-p104-product-operation-linkage/p104-operation-linkage-summary.json`
- `docs/release/ui-audit-assets/2026-06-24-p104-product-operation-linkage/p104-runner-output.log`

Runner status: `passed`

Covered operation steps:

- `portfolio_invalid_total_rejected`
- `portfolio_adjustment_readback`
- `holding_edit_remove`
- `offline_transaction`
- `batch_import_validate_confirm`
- `correction_and_rebalance`
- `decision_confirmation_loop`
- `notifications_mark_read`
- `risk_lifecycle`
- `data_quality_resolution_create_retire`
- `downstream_dashboard_review_audit`

SQLite readback summary from the fresh run:

| Evidence | Value |
| --- | ---: |
| `portfolio_snapshots` | 6 |
| `positions` | 3 |
| `position_transactions` | 3 |
| `operation_confirmations` | 3 |
| `p104_executed_confirmations` | 1 |
| `import_batches_committed` | 1 |
| `corrections` | 1 |
| `risk_resolved` | 1 |
| `p104_notification_read` | 1 |
| `dq_resolutions_retired` | 1 |
| `audit_events` | 15 |
| `latest_total_assets` | 1790.00 |
| `latest_position_market_value` | 803.75 |
| `latest_cash` | 986.25 |

Safety negative evidence:

| Check | Value |
| --- | ---: |
| `auto_confirmation_rows` | 0 |
| `forbidden_broker_order_push_tables` | 0 |
| `auto_rule_apply_audit_events` | 0 |

## Product Acceptance Meaning

P104 answers the practical product-acceptance question: after a user operation, does the product show the expected result, persist the expected data, expose the result through linked downstream views/APIs, preserve auditability, and keep forbidden automation absent?

The fresh runner confirms the core local product operation chain:

- Invalid portfolio totals are rejected before data mutation.
- Valid portfolio adjustment writes current holdings/snapshots and can be read back.
- Holding edit/remove operations update local account facts.
- Offline transactions and batch imports create durable local transaction facts.
- Corrections and rebalance review are auditable local records/suggestions.
- A formal decision can be manually confirmed as executed, then read through decision detail and decision-loop APIs.
- Dashboard, review summary, audit events, notifications, risk alerts, and data-quality gate resolution read or update the linked state coherently.
- Forbidden automation remains absent.

## Relationship To Existing Evidence

P104 complements, not replaces:

- P92 final original-requirement ledger.
- P93 code reality/design audit.
- P102/P103 real-model product UI acceptance and UX linkage fixes.
- Existing focused Go/frontend tests for branch-level behavior.

P104 is intentionally a repeatable representative linked-flow gate. It is not a proof that every invalid input branch, browser viewport, deployment path, or future external provider condition has been exhausted.
