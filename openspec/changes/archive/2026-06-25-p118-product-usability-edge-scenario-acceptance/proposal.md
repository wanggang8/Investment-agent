# P118 Product Usability Edge Scenario Acceptance

## Summary

Run a deeper product-usability acceptance pass for the non-release scenarios that remain after P117: long-period local use, abnormal data and input recovery, decision-quality interpretation, multi-account household ledger behavior and safety negative evidence.

## Motivation

P117 proved a seven-day continuous-use path. The user explicitly asked to exclude release/install/upgrade scenarios for now and continue with the remaining product-use scenarios. P118 therefore focuses on realistic edge cases that determine whether the local discipline tool remains usable when history grows, data quality fluctuates, decisions differ by context and a user maintains several accounts/funds.

## In Scope

- 30-day local durability story with accumulated reports, transactions, audit rows, notifications and risk records.
- Abnormal input and data-quality recovery: invalid imports, duplicate-like transactions, stale/missing source facts and scoped data-quality resolutions.
- Decision-quality interpretation across rising, falling and volatile contexts without claiming prediction accuracy.
- Multi-account/household ledger simulation using explicit local account tags in positions and transaction notes.
- Cross-page consistency for portfolio, dashboard, workbench, review, audit, decision loop, data quality, risk alerts and notifications.
- Browser evidence for core accumulated-state surfaces and 390px mobile checks.
- Safety negative evidence for broker/order/push tables, automatic confirmation, automatic rule application, trading affordances, return guarantees and sensitive leakage.

## Out of Scope

- Docker, install, upgrade, uninstall, release package refresh, Git tag, GitHub Release or physical second-machine validation.
- New investment runtime capability.
- Broker integration, one-click trading, order placement, external push or automatic trading.
- Fresh external provider guarantee or fresh real LLM quality claim.
- Future return, future market direction or model-accuracy guarantee.
- Archiving P114/P115/P116/P117/P118 without user confirmation.

## Success Criteria

- P118 runner completes with all edge usability scenarios passed.
- Evidence includes API/SQLite summary, restart summary, browser summary, final interpretation JSON and screenshots.
- Long-history counts prove data accumulation without page/API failure.
- Decision interpretation scenarios prove context-sensitive, traceable and bounded recommendations.
- Multi-account facts remain readable and auditable without pretending to connect to broker accounts.
- Regression gates pass, with P93 stale status recorded if still stale.
