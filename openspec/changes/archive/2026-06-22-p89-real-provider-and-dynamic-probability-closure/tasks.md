# P89 Tasks

## 1. Plan And Inventory

- [x] Confirm P89 covers exactly the 10 full-release-required rows that remain non-`real_pass` after P88.
- [x] Build `scripts/p89_remaining_real_provider_dynamic_inventory_check.py` to parse the P88 matrix, emit `docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability/p89-inventory.json`, and fail on row-count or row-id drift.
- [x] Request subagent plan review before implementation and resolve every Critical/Important finding before coding.

## 2. Real Structured Provider Verification And Collection

- [x] Build a source-preverification registry/checker for capital-flow, margin-financing, and constituent-financial providers. Each category must record authority, public access shape, stable request/page evidence, fields, update frequency, legal/access limits, rate assumptions, failure behavior, runtime status, and SQLite target path.
- [x] Add provider-level tests that fail unless accepted-local/fixture/stub/manual seed evidence is explicitly excluded from `real_pass` eligibility.
- [x] Implement read-only runtime provider adapters only for sources that pass no-login/no-paid/no-authorization/no-Level2/no-high-frequency checks.
- [ ] Add SQLite/API readback for capital-flow `date`, `net_inflow`, `net_outflow`.
- [x] Add SQLite/API readback for margin-financing `date`, `margin_balance`, `balance_change_rate`.
- [x] Add SQLite/API readback for constituent-financial `revenue`, `net_profit`, `growth`, `disclosure_date`.
- [x] If any category cannot be verified with a real public provider, keep its rows `partial` and record the exact provider blocker in the P89 matrix.

## 3. Extreme Fear Historical Scenario UI

- [x] Add tests and implementation for extreme-fear state that locks active trading advice.
- [x] Add historical similar-scenario context fields and UI display for the extreme-fear state.
- [x] Extend real browser acceptance to trigger the extreme-fear path and verify UI/API/SQLite readback.

## 4. Dynamic Expected-Return Monitoring

- [x] Add tests for baseline and changed valuation/fundamental/market-state inputs that lower affected scenario probabilities.
- [x] Add tests for periodic core-assumption checks.
- [x] Add tests for two consecutive months below expectation producing a scenario-downshift warning.
- [x] Add tests for one month on a pessimistic actual path producing a manual probability-adjustment suggestion.
- [x] Implement minimal DTO/API/UI/readback fields for dynamic monitoring without adding automatic trade or automatic probability application behavior.
- [x] Extend real browser acceptance to run before/after dynamic expected-return scenarios through UI/API/SQLite readback.

## 5. Final P89 Matrix And Claims

- [x] Build `scripts/p89_real_provider_dynamic_probability_closure.py` to generate `docs/release/acceptance/2026-06-22-p89-real-provider-dynamic-probability-closure.md` and `docs/release/acceptance/2026-06-22-p89-real-provider-dynamic-probability-matrix.md`.
- [x] Upgrade only rows with direct P89 evidence; preserve exact blockers for any row still not proven.
- [x] Update release/governance docs with the P89 result and avoid full-pass claims unless the matrix has no full-release-required non-`real_pass` rows.
- [x] Run subagent final review and resolve every Critical/Important finding.
- [x] Run `openspec validate --all --strict`, P89 inventory/closure checks, provider preverification, P89 real browser acceptance, focused Go/frontend tests, frontend build, and `git diff --check`.
- [x] Archive P89 after validation and review pass.
