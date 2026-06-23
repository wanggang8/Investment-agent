# P84 Tasks

## 1. Plan And Inventory

- [x] Confirm the P84 row set contains exactly 35 portfolio/allocation/confirmation data-impact rows from the latest matrix, including P82-deferred `REQ-10-001` through `REQ-10-004`.
- [x] Map each row to a before/after UI scenario, database table/readback target, and downstream surface.

## 2. Acceptance Harness

- [x] Build or extend a P84 real browser runner with temporary SQLite and real local backend/frontend.
- [x] Execute account, holding, offline transaction, manual confirmation, review, and decision-loop scenarios where supported.
- [x] Verify API responses, read-only SQLite before/after deltas, audit events, and UI readback after navigation or refresh.
- [x] Independently recompute deterministic portfolio values used in evidence.
- [x] Verify forbidden broker/order/automatic confirmation affordances remain absent.

## 3. Runtime Fixes If Needed

- [x] Fix product defects that block required portfolio/confirmation readback or data impact.
- [x] Add focused tests for any code changes.

## 4. Evidence And Governance

- [x] Generate P84 acceptance record and updated evidence layer.
- [x] Update release/governance docs with P84 row upgrades and remaining full-release-required count.
- [x] Run `openspec validate --all --strict`.
- [x] Run P84 runner and relevant Go/frontend tests.
- [x] Run read-only subagent review before archive and resolve all Critical/Important findings.
- [ ] Archive P84 after validation and review pass.
