# P117 Continuous Product Usability Acceptance

## Summary

Run a seven-day continuous-use acceptance pass that interprets whether the product is usable as a real local investment discipline tool, not merely whether isolated APIs or pages respond. P117 links cold start, daily routine, offline ledger updates, invalid input recovery, data-quality degradation, manual decision confirmation, marked-error review, cross-page readback, restart persistence and safety boundaries.

## Motivation

P115 proved broad real-user scenarios and P116 proved richer multi-fund transaction ledger behavior. The remaining product-readiness question is whether a user can keep using the product over time, understand what changed, recover from bad inputs and trust that pages remain consistent after several days of local facts.

## In Scope

- Seven-day local seeded user story, with day-by-day operations and interpretation.
- Empty/cold-start usability behavior.
- Portfolio onboarding and daily dashboard/workbench/review readback.
- Multi-fund offline transactions, valid/invalid import and correction paths.
- Data-quality gate degradation, explicit user resolution and retirement.
- Manual decision confirmation and marked-error learning loop.
- Restart persistence using the same SQLite database.
- Browser screenshots for key daily surfaces and 390px mobile check.
- Safety negative evidence: no broker/order/push tables, no auto confirmation, no auto rule application, no trading/return-guarantee affordance.

## Out of Scope

- New investment runtime capability.
- Broker integration, one-click trading, order placement, external push or automatic trading.
- Fresh external provider guarantee or fresh real LLM quality claim.
- Docker, installer, release package refresh, Git tag, GitHub Release or physical second-machine validation.
- Archiving P114/P115/P116/P117 without user confirmation.

## Success Criteria

- P117 runner completes with all usability scenarios passed.
- Evidence includes API/SQLite summary, browser summary, final usability interpretation report and screenshots.
- Restart persistence is proven by reading previously written local facts after backend restart.
- Product usability report explicitly distinguishes fresh pass, scoped pass and claim boundaries.
- Regression gates pass, with P93 stale status recorded if still stale.
