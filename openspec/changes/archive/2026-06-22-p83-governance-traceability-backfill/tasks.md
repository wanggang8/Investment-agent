# P83 Tasks

## 1. Plan And Inventory

- [x] Confirm the P83 row set contains exactly 43 governance/review traceability rows from the latest matrix, including P82-deferred `REQ-12-002`, `REQ-12-003`, and `REQ-13-011`.
- [x] Classify each row as real product behavior, release/governance evidence, reference-only, scoped, partial, or non-goal candidate.

## 2. Evidence Backfill

- [x] Link each row to exact docs, acceptance files, scripts, commands, tests, package manifests, or screenshots where applicable.
- [x] Execute fresh checks where current evidence is stale, ambiguous, or missing.
- [x] Add a machine-readable P83 evidence layer with row status, artifact links, and rationale.

## 3. Runtime Fixes If Needed

- [x] Fix only product defects that block rows requiring current product behavior.
- [x] Add focused tests for any code changes.

## 4. Governance

- [x] Update release/governance docs with P83 row upgrades and remaining full-release-required count.
- [x] Run `openspec validate --all --strict`.
- [x] Run all P83 evidence checkers and relevant tests.
- [x] Run read-only subagent review before archive and resolve all Critical/Important findings.
- [x] Archive P83 after validation and review pass.
