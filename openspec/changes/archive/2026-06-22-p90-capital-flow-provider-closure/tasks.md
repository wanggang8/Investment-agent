# P90 Tasks

## 1. Plan And Inventory

- [x] Confirm P90 owns exactly `REQ-04-016` and `REQ-05-003` from the P89 remaining rows.
- [x] Build `scripts/p90_capital_flow_inventory_check.py` and emit `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/p90-inventory.json`.
- [x] Request subagent plan review before implementation and resolve every Critical/Important finding.

## 2. Provider Verification And Runtime Collection

- [x] Build P90 source preverification for Eastmoney H5 capital-flow endpoint and exclude fixture/stub/accepted-local/manual seed evidence.
- [x] Write failing Go tests for Eastmoney H5 capital-flow parsing and directional `net_inflow` / `net_outflow` mapping.
- [x] Implement runtime H5 capital-flow fetch in the structured public collector.
- [x] Verify runtime collector does not depend on broker, login, paid, authorization, Level2, or high-frequency source.

## 3. Real UI/API/SQLite Acceptance

- [x] Add real browser acceptance that clicks Settings UI market refresh and reads capital-flow fields from market snapshot API.
- [x] Add SQLite readback checker for the runtime snapshot `capital_flow.date`, `net_inflow`, `net_outflow`, and `raw_net_flow`.
- [x] Run P90 acceptance runner without manually seeding capital-flow fields.

## 4. Final P90 Matrix And Claims

- [x] Generate P90 closure and matrix with `REQ-04-016` and `REQ-05-003` upgraded only if direct evidence exists.
- [x] Update governance docs to show P90 result and avoid package/physical-machine claims.
- [x] Run subagent final review and resolve every Critical/Important finding.
- [x] Run `openspec validate --all --strict`, P90 checks, P90 real browser acceptance, Go/frontend tests, frontend build, and `git diff --check`.
- [x] Archive P90 after validation and review pass.
