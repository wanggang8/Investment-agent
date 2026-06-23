# P82 Tasks

## 1. Plan And Inventory

- [x] Confirm the P82 row set contains exactly 53 SOP/action data-impact rows from the P80/P81 matrix for evaluation.
- [x] Map each row to a real UI scenario, expected API result, expected SQLite/readback target, safety negative check, and upgrade-or-defer decision.

## 2. Acceptance Harness

- [x] Build or extend a P82 browser runner for SOP/action scenarios against a real local Go backend, Vite frontend, and temporary SQLite database.
- [x] Verify post-action API responses and read-only SQLite deltas for SOP states, confirmations, notifications, audit events, risk alerts, reviews, and related references.
- [x] Verify page refresh/navigation readback for the resulting user-visible state.
- [x] Verify unsafe or unsupported action affordances remain absent or blocked.

## 3. Runtime Fixes If Needed

- [x] Fix product defects that prevent required real UI-to-SQLite behavior from passing.
- [x] Add focused tests for any code changes.

## 4. Evidence And Governance

- [x] Generate P82 acceptance record and updated evidence layer.
- [x] Update release/governance docs with P82 row upgrades and remaining full-release-required count.
- [x] Run `openspec validate --all --strict`.
- [x] Run P82 runner and relevant Go/frontend tests.
- [x] Run read-only subagent review before archive and resolve all Critical/Important findings.
- [x] Archive P82 after validation and review pass.
