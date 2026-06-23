# P88 Tasks

## 1. Plan And Inventory

- [x] Confirm P88 covers exactly the 27 full-release-required rows that remain non-`real_pass` after P86.
- [x] Build `scripts/p88_remaining_blockers_inventory_check.py` to parse the P86 matrix, emit `docs/release/ui-audit-assets/2026-06-22-p88-remaining-full-release-blockers/p88-inventory.json`, and fail on row-count or row-id drift.
- [x] Request subagent plan review before implementation and resolve every Critical/Important finding before coding.

## 2. Source-Verified State Transitions

- [x] Add focused tests for buy-logic-broken >=2 A/S formal sources -> `sell_only`, prohibited buy/add action text, decision metadata, portfolio state readback, source-count provenance, and audit event.
- [x] Add focused tests for multi-source-insufficient <2 A/S formal sources -> `frozen_watch`, pause action text, decision metadata, portfolio state readback, source-count provenance, and audit event.
- [x] Implement the minimal workflow/service/API behavior needed for those tests without broker/order/auto-confirm behavior.
- [x] Extend real browser acceptance to trigger both transition paths and verify UI/API/SQLite/readback.

## 3. Structured Data And Source Preverification

- [x] Add a source-preverification registry/artifact for P88 capital-flow, margin-financing, and constituent-financial fields, including authority, public access shape, stable request/page evidence, fields, update frequency, legal/access limits, rate-limit assumption, failure behavior, and SQLite target path.
- [x] Add tests for collector/readback normalization of capital-flow `date`, `net_inflow`, `net_outflow`.
- [x] Add tests for collector/readback normalization of margin-financing `date`, `margin_balance`, `balance_change_rate`.
- [x] Add tests for collector/readback normalization of constituent-financial `revenue`, `net_profit`, `growth`, `disclosure_date`.
- [x] Implement the minimal collector/readback path and explicit source-unavailable fallback; P88 records that no verified non-mock runtime provider exists, so `REQ-05-003`, `REQ-05-004`, and `REQ-05-005` remain `partial`.
- [x] Extend P88 source-preverification and closure acceptance to display field-level readiness and persisted parser/readback paths; do not claim real-provider field completion.

## 4. Expected Return Historical/Probability Closure

- [x] Add tests for historical-sample probability calculation: optimistic probability from similar-sample proportion, base scenario as highest-frequency path, pessimistic scenario always displayed, and sample metadata preserved.
- [x] Add a representative holding-class coverage test/matrix for broad ETF/index fund, sector/growth ETF or fund, and equity/security-like constituent-financial path; keep `REQ-09-001` partial if any class is not proven.
- [x] Add tests for dynamic update by valuation/fundamental/market-state changes and scenario rerun that lowers affected probabilities.
- [x] Add tests for expected-return report fields: fund/security name and code, explicit future-12-month label, scenario ranges, probability basis, sample count/window/screening condition, sell-evaluation trigger, disclaimer, and supplement-data list.
- [x] Add tests for periodic assumption checks, two-month below-expectation downshift warning, and one-month pessimistic-path probability-adjustment suggestion.
- [x] Implement the minimal deterministic expected-return engine and DTO/API/UI readback needed for those tests.
- [x] Extend real browser acceptance to run complete and sample-below-5 scenarios through UI/API/SQLite/readback; dynamic-downshift and periodic assumption rows remain `partial` because P88 did not prove them through real UI/API/SQLite.

## 5. Quarterly Rebalance

- [x] Add tests for quarterly +/-15% drift calculation across core/satellite/cash targets and manual buy/sell recommendation amounts.
- [x] Implement API/service/UI readback for rebalance recommendations as offline manual actions only.
- [x] Extend real browser acceptance to create an out-of-band portfolio, trigger quarterly rebalance, and verify UI/API/SQLite/audit readback plus forbidden broker/order table absence.

## 6. SOP Addendum Proposal

- [x] Add tests for high-frequency uncovered-scenario detection generating a `sop` rule proposal, notification, and audit event.
- [x] Implement the minimal pipeline using existing `rule_proposals`, `notifications`, and `audit_events`; do not write active rule versions without existing gatekeeper/final-confirm rules.
- [x] Extend real browser acceptance to verify proposal visibility, pending status, notification linkage, audit readback, and no automatic rule application.

## 7. Final P88 Matrix And Claims

- [x] Build `scripts/p88_remaining_full_release_blockers_closure.py` to generate `docs/release/acceptance/2026-06-22-p88-remaining-full-release-blockers-closure.md` and `docs/release/acceptance/2026-06-22-p88-remaining-full-release-blockers-matrix.md`.
- [x] Upgrade only rows with direct P88 evidence; preserve exact blockers for any row still not proven. Treat any reclassification as a separate L1/OpenSpec rationale, not equivalent to `real_pass`.
- [x] Update release/governance docs with the P88 result and avoid full-pass claims unless the matrix has no full-release-required non-`real_pass` rows.
- [x] Run subagent final review and resolve every Critical/Important finding.
- [x] Run `openspec validate --all --strict`, P88 inventory/closure checks, focused Go/frontend tests, P88 real browser acceptance, frontend build, and `git diff --check`.
- [ ] Archive P88 after validation and review pass.

## 8. P88 Closure Note

P88 closure matrix currently upgrades 17 of the 27 P86 remaining full-release-required blockers to `real_pass` and keeps 10 rows `partial`: `REQ-04-016`, `REQ-05-003`, `REQ-05-004`, `REQ-05-005`, `REQ-08-004`, `REQ-08-023`, `REQ-09-004`, `REQ-09-023`, `REQ-09-024`, and `REQ-09-025`. This is intentional and prevents accepted-local/parser/unit evidence from being reported as real provider or real UI dynamic-probability completion.
