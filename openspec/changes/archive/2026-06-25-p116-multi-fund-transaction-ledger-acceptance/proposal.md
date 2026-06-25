# P116 Multi-Fund Transaction Ledger Acceptance

## Summary

Run a focused acceptance pass for complex real-user portfolio ledger scenarios: multiple funds/ETFs, multiple offline transactions, mixed batch import, invalid transaction rejection, decision-to-manual-execution linkage, risk/data-quality/review/audit aggregation, and safety negative evidence.

## Motivation

P115 proved broad product linkage across 34 scenarios, but the portfolio ledger portion still looked closer to "multi-symbol smoke coverage" than a realistic multi-fund transaction history. P116 closes that gap by treating the local portfolio ledger as the primary acceptance object and verifying that realistic multi-fund operations remain coherent across API, SQLite, browser UI, dashboard/workbench, review, audit, notifications, risk, and data-quality surfaces.

## Scope

In scope:

- Fresh isolated local backend/frontend/SQLite runner.
- Multi-fund initial portfolio with at least five symbols.
- Multiple offline buy/sell/reduce/clear-like operations across different dates and fees.
- Mixed batch import with holdings, transactions, valid rows, invalid rows, and confirm-after-validate behavior.
- Rejection checks for insufficient cash, oversell, future execution time, negative fees, invalid state, missing symbol, and invalid settings mutation.
- Manual decision confirmation writing an offline execution record.
- Marked-error decision review path.
- Risk alert lifecycle, notification readback, data-quality gate resolution, rebalance review, review summary, dashboard/workbench, and audit readback.
- Browser proof for core multi-fund UI paths and mobile portfolio rendering.
- Safety checks for no broker/order/push tables, no automatic confirmation, no automatic rule application, no trading affordance, and no return guarantee language.

Out of scope:

- Broker integration, one-click trading, order placement, automatic trading, external push, automatic confirmation, automatic rule application, automatic repair/recovery, return guarantees.
- Paid/login/authorized data sources, Level2 data, high-frequency data.
- Docker, installer, release package refresh, GitHub Release, or physical second-machine validation.
- Fresh external provider or real LLM claims.

## Acceptance

P116 passes only if the runner produces a final merged summary with:

- All P116 scenarios classified as `fresh_pass` or explicitly bounded `scoped_pass`.
- Zero `blocked` scenarios.
- SQLite readback for portfolio snapshots, positions, position transactions, operation confirmations, corrections, risk lifecycle, notifications, data-quality resolutions, and audit events.
- Browser proof for `/positions`, `/decisions/:id`, `/decision-loop`, `/`, `/workbench`, `/review`, `/audit`, `/risk-alerts`, `/notifications`, `/data-quality`, and mobile `/positions`.
- Safety counters all zero for forbidden automation and secret/raw prompt leakage.
