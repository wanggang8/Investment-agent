# Design: P118 Product Usability Edge Scenario Acceptance

## Approach

P118 is an acceptance-only change. It starts the existing local backend with temporary config and SQLite, writes realistic long-history local facts through a mix of public APIs and controlled SQLite seed data, restarts the backend against the same database, then runs a browser pass through the accumulated product surfaces.

The acceptance intentionally excludes release/install/upgrade scenarios. It does not add backend capabilities unless a real blocker is discovered.

## Scenario Groups

1. Long-cycle durability: 30 daily reports, many audit events, notifications, risk records and transaction facts remain readable.
2. Abnormal recovery: invalid imports and invalid transactions are rejected without partial writes; data-quality degradation creates explicit scoped resolution rather than a clean claim.
3. Decision-quality interpretation: rising, falling and volatile contexts produce different local seeded decision records, each with explicit final verdict, evidence status, prohibited actions and disclaimer boundaries.
4. Household ledger: multiple local account tags and several funds/ETFs can coexist, with cash/core/satellite/active-fund state readable across pages.
5. Cross-page and mobile usability: accumulated facts render on primary pages without console errors, page errors or API 5xx failures.
6. Safety: no broker/order/push execution tables, no auto confirmations, no auto rule application and no visible trading/return-guarantee affordance.

## Files

- `scripts/p118-product-usability-edge-scenario-acceptance.sh`: isolated lifecycle runner for backend, restart probe, frontend and Playwright.
- `scripts/p118_product_usability_edge_scenario_acceptance.py`: API/SQLite runner, seed data, restart probe and final JSON merge.
- `web/e2e/p118-product-usability-edge-scenario-acceptance.spec.ts`: browser journey and screenshots.
- `docs/release/acceptance/2026-06-25-p118-product-usability-edge-scenario-acceptance-matrix.md`: scenario matrix.
- `docs/release/acceptance/2026-06-25-p118-product-usability-edge-scenario-acceptance.md`: final acceptance record.

## Evidence Semantics

- `fresh_pass`: verified through the P118 local backend/frontend/SQLite/browser run.
- `scoped_pass`: verified as local seeded/degraded-safe behavior, without claiming external provider cleanliness, fresh LLM quality, broker execution or release readiness.
- `claim_boundary`: explicit negative boundary attached to final summary and Markdown record.

## Safety

Safety evidence is collected from both SQLite and rendered UI:

- SQLite has no broker/order/push execution tables.
- `operation_confirmations` has no auto confirmation rows.
- `audit_events` has no automatic rule/trade application traces.
- Browser-visible primary flows do not expose one-click trading, broker order placement, external push or return-guarantee claims.
