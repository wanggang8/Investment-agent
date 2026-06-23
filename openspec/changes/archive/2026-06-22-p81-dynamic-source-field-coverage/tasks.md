# P81 Tasks

## 1. Plan And Inventory

- [x] Confirm P80 remains the baseline matrix and that the P81 row set contains exactly 59 full-release-required rows.
- [x] Record the row list, previous statuses, target evidence type, and expected pass boundary in a P81 matrix artifact.

## 2. Acceptance Harness

- [x] Build or extend a P81 runner that uses a user-selected non-`510300` fund/ETF symbol and rejects hard-coded `510300`-only evidence.
- [x] Verify source/readiness API responses include provenance, freshness, degraded/missing state, feature impact, and sanitized safety notes.
- [x] Verify UI readback for data-quality/readiness/consultation surfaces through real local browser operation.
- [x] Verify read-only SQLite evidence for relevant source facts, health, evidence references, RAG/index records, and audit events.

## 3. Runtime Fixes If Needed

- [x] Fix any discovered product gaps required for the 59 rows to pass, while preserving forbidden-capability boundaries.
- [x] Add focused Go/frontend tests for any code changes.

## 4. Evidence And Governance

- [x] Generate P81 acceptance record and updated P81 evidence layer without rewriting P75-P80 historical matrices.
- [x] Update release/governance docs with the new P81 status and remaining full-release-required count.
- [x] Run `openspec validate --all --strict`.
- [x] Run the P81 acceptance runner and all relevant Go/frontend tests.
- [x] Run a read-only subagent review before archive and resolve all Critical/Important findings.
- [x] Archive P81 after validation and review pass.
