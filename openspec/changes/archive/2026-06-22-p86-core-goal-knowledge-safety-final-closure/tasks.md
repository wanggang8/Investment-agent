# P86 Tasks

## 1. Plan And Inventory

- [x] Confirm P86 covers exactly the 137 full-release-required rows that remain non-`real_pass` after P87.
- [x] Run `python3 scripts/p85_p87_p86_plan_inventory_check.py` to generate a machine-checkable P86 inventory artifact listing all 137 rows by section, and fail if the P87 matrix count, row IDs, or ownership set drift.
- [x] Request subagent plan review before implementation; resolve every Critical/Important finding before starting integrated execution.

## 2. Integrated Acceptance

- [x] Build or extend a P86 end-to-end real browser runner covering setup, portfolio/account state, data readiness, knowledge/RAG, consultation, expected return, risk/SOP, manual confirmation, review, audit, release governance, and safety.
- [x] Verify each scenario through UI, API, read-only SQLite, workflow metadata, and deterministic checks where applicable.
- [x] Verify desktop and mobile UI clarity for the integrated journey, including navigation, status copy, data impact readback, and absence of forbidden action affordances.
- [x] Verify final product-goal metrics: discipline adherence support, evidence sufficiency, traceability, review usefulness, safe degradation, data-source honesty, and user comprehension.
- [x] Record any product defects discovered during integrated acceptance and fix only defects required to make the real scenario truthful and usable.

## 3. Final Matrix And Claims

- [x] Generate final row-level matrix after P81-P85, P87, and P86 evidence.
- [x] Upgrade rows to `real_pass` only when P86 or cumulative P81-P87 evidence directly proves the row; do not use seeded-only data, route smoke, screenshots, fixture/mock/stub data, or broad narrative as upgrade evidence.
- [x] If any full-release-required row remains non-`real_pass`, record exact blockers, missing implementation, unavailable external data condition, or non-goal/reference-only rationale and avoid full-pass claims.
- [x] If every full-release-required row is `real_pass` or validly non-goal/reference-only, prepare the full original-requirement pass statement with exact evidence references.

## 4. Runtime Fixes If Needed

- [x] Fix product defects that block final integrated acceptance.
- [x] Add focused Go/frontend tests for any code changes.
- [x] Re-run the affected real UI/API/SQLite scenario after every product fix.

## 5. Evidence And Governance

- [x] Generate P86 final acceptance record and evidence layer.
- [x] Update release/governance docs with final P86 conclusion.
- [x] Run `openspec validate --all --strict`.
- [x] Run P86 runner and relevant Go/frontend tests.
- [x] Run read-only subagent review before archive and resolve all Critical/Important findings.
- [x] Archive P86 after validation and review pass.
