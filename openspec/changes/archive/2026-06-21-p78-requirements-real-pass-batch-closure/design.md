# P78 Design

## Approach

P78 follows the same conservative evidence-layer pattern as P77. The historical P75 and P77 matrices remain immutable evidence records. P78 reads the P77 matrix, creates a new row-level closure matrix, and upgrades only rows that satisfy all applicable dimensions for the current batch.

## Batch Strategy

The remaining P77 gaps are grouped by remediation shape:

- `core_product_goal`: broad product-purpose rows that need decomposition into scenario-level checks.
- `data_source_dynamic`: symbol, source-health, collector, and field propagation rows.
- `expected_return`: probability, sample-count, scenario, sell-evaluation, and disclaimer rows.
- `sop_action_data_impact`: SOP lifecycle, action-to-table, audit, and readback rows.
- `portfolio_confirmation_data`: local account, holding, snapshot, confirmation, and transaction rows.
- `knowledge_llm_rag`: built-in knowledge, local knowledge, LLM context, and retrieval rows.
- `governance_traceability`: release, roadmap, packaging, and traceability rows.
- `safety_boundary`: negative capability rows already handled by P77 unless new evidence is needed.

P78 batch A targets the low-sample expected-return degradation/disclaimer slice because it can be proven with a narrow but complete evidence combination:

- deterministic Go tests for all precision states, sell-evaluation triggers, sample count, sample window, screening condition, missing-price context, and workflow dynamic inputs;
- a real accepted-local non-`510300` UI journey with `use_stub=false` and a real LLM provider;
- SQLite readback of the generated decision's `expected_return_scenarios_json`;
- verification that low-sample output does not expose precise probabilities and includes sample count, sample window, screening condition, source-health context, and non-trading disclaimer.

## Upgrade Rules

A row can become `real_pass` in P78 only when:

- it is present in the P77 matrix and is still full-release-required;
- its P78 evidence artifacts exist and pass schema checks;
- the applicable implementation, UI, data/readback, workflow/rule/LLM, scenario, and safety dimensions are directly evidenced;
- the generated P78 conclusion remains scoped unless every full-release-required row is `real_pass`.

P78 must not upgrade broad rows that are only partially represented by batch A evidence. For example, a deterministic expected-return trigger test does not prove every user-facing sell-evaluation workflow row unless real UI and readback evidence cover that exact behavior.

## Generated Artifacts

- `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-closure.md`
- `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-matrix.md`
- `docs/release/ui-audit-assets/2026-06-21-p78/real-pass-batch-summary.json`
- `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.log`
- `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.json`
- `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-ui-readback.json`
- `docs/release/ui-audit-assets/2026-06-21-p78-non-510300/`
