# Design: P117 Continuous Product Usability Acceptance

## Approach

P117 is an acceptance-only change. It uses the existing local backend, Vite frontend, temporary SQLite database and Playwright browser runner. It does not modify product behavior unless the acceptance uncovers a real blocker.

The runner simulates a seven-day user story:

1. Day 0: cold start and empty-state checks.
2. Day 1: local portfolio onboarding and first readback.
3. Day 2: daily routine across dashboard/workbench/review/audit.
4. Day 3: offline ledger update and risk/notification handling.
5. Day 4: invalid import and invalid transaction recovery, then correction audit.
6. Day 5: data-quality degradation, explicit resolution and retirement.
7. Day 6: decision manual confirmation and marked-error loop.
8. Day 7: final review, audit, decision loop and restart persistence readback.

## Files

- `scripts/p117-continuous-product-usability-acceptance.sh`: starts isolated backend/frontend, runs API/SQLite checks, restarts backend, runs browser checks and merges evidence.
- `scripts/p117_continuous_product_usability_acceptance.py`: API/SQLite runner, restart persistence probe and final usability interpretation merger.
- `web/e2e/p117-continuous-product-usability-acceptance.spec.ts`: browser journey and screenshots.
- `docs/release/acceptance/2026-06-25-p117-continuous-product-usability-acceptance-matrix.md`: scenario matrix.
- `docs/release/acceptance/2026-06-25-p117-continuous-product-usability-acceptance.md`: final acceptance record.

## Evidence Semantics

- `fresh_pass`: verified through live local backend/frontend/SQLite/browser in this run.
- `scoped_pass`: verified as local seeded/degraded-safe behavior, without claiming external provider, real LLM quality, broker execution or release packaging.
- `claim_boundary`: explicit negative boundary attached to summary and final record.

## Safety

The safety evidence is both database and UI based:

- SQLite `sqlite_master` contains no broker/order/push execution tables.
- `operation_confirmations` contains no auto confirmation rows.
- `audit_events` contains no auto rule/trade application events.
- Browser-visible primary flows do not expose one-click trading, broker order placement, external push or return-guarantee claims.
