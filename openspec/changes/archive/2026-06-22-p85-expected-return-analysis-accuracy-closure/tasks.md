# P85 Tasks

## 1. Plan And Inventory

- [x] Confirm the P85 row set contains exactly 31 expected-return and analysis-accuracy rows from the P84 matrix.
- [x] Map each row to complete-data, degraded-data, LLM-quality-failure, UI-readback, deterministic-check, or safety-boundary evidence.

## 2. Acceptance Harness

- [x] Build or extend a P85 runner that performs real UI consultation/decision scenarios against real local backend/frontend.
- [x] Verify expected-return/scenario API fields, UI readback, provenance, sample/window metadata, and disclaimers.
- [x] Independently recompute deterministic expected-return and portfolio-derived values used in the evidence.
- [x] Verify degraded/missing data and LLM quality failure safely block, discard, or qualify analysis without creating trade confirmation.
- [x] Verify LLM material does not override final rule verdict or trigger automatic actions.

## 3. Runtime Fixes If Needed

- [x] Fix product defects that block required expected-return or analysis-boundary behavior.
- [x] Add focused Go/frontend tests for any code changes.

## 4. Evidence And Governance

- [x] Generate P85 acceptance record and updated evidence layer.
- [x] Update release/governance docs with P85 row upgrades and remaining full-release-required count.
- [x] Run `openspec validate --all --strict`.
- [x] Run P85 runner and relevant Go/frontend tests.
- [x] Run read-only subagent review before archive and resolve all Critical/Important findings.
- [ ] Archive P85 after validation and review pass.
