# P87 Tasks

## 1. Plan And Inventory

- [x] Confirm P87 covers exactly the 32 P84-after full-release-required non-`real_pass` rows not owned by P85 or P86.
- [x] Reconcile P85, P87, and P86 planned row sets to prove all 157 P84-after rows have exactly one owner.
- [x] Map each P87 row to portfolio state, allocation/rebalance, data-insufficient safety, manual confirmation/audit, public-source/SQLite readiness, UI-readback, or release-safety evidence.

## 2. Acceptance Harness

- [x] Build or extend a P87 runner that performs real UI portfolio/allocation/state scenarios against real local backend/frontend.
- [x] Verify core/satellite/cash allocation, market value, cash, position state, buy date, and derived ratios through UI, API, SQLite, and deterministic checks; quarterly rebalance action flow was evaluated but not upgraded because P87 does not prove the full row.
- [x] Verify sell-only and frozen-watch transitions, data-insufficient states, and multi-source-insufficient states visibly block or qualify trade-like advice.
- [x] Evaluate proposal confirmation/rejection, audit events, structured summaries, and downstream dashboard/review/readback rows; rows whose full breadth was not freshly proven remain deferred/non-`real_pass`.
- [x] Evaluate public-source readiness and SQLite storage claims only with actual collector/readiness and stored-field evidence; rows without full fresh evidence remain deferred/non-`real_pass`.
- [x] Verify release/install safety boundaries and forbidden automatic behaviors remain absent for the P87 scenario; full release/upgrade preflight row remains deferred/non-`real_pass`.

## 3. Runtime Fixes If Needed

- [x] Fix product defects that block required portfolio state, allocation, safety, or readback behavior.
- [x] Add focused Go/frontend tests for any code changes.

## 4. Evidence And Governance

- [x] Generate P87 acceptance record and updated evidence layer.
- [x] Update release/governance docs with P87 row upgrades and remaining full-release-required count.
- [x] Run `openspec validate --all --strict`.
- [x] Run P87 runner and relevant Go/frontend tests.
- [x] Run read-only subagent review before archive and resolve all Critical/Important findings.
- [ ] Archive P87 after validation and review pass.
