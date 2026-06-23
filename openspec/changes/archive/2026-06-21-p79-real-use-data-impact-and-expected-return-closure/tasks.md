# Tasks: P79 Real Use Data-Impact And Expected-Return Closure

## 1. Setup And Governance

- [x] 1.1 Confirm no other active OpenSpec change exists before executing P79.
- [x] 1.2 Create OpenSpec change `p79-real-use-data-impact-and-expected-return-closure`.
- [x] 1.3 Define P79 as a post-P78 evidence layer, not a rewrite of historical matrices.
- [x] 1.4 Mark P79 active in `docs/GOVERNANCE.md`, `openspec/project.md`, `openspec/PROGRESS.md`, `AGENTS.md`, and `docs/development-plan.md`.
- [x] 1.5 Run `openspec validate p79-real-use-data-impact-and-expected-return-closure --strict`, `openspec validate --all --strict`, and `git diff --check`.

## 2. P79 Matrix And Checker

- [x] 2.1 Implement `scripts/p79_real_use_data_impact_and_expected_return_closure.py`.
- [x] 2.2 Read the P78 matrix and preserve row identity, source lines, requirement text, P78 status, remediation group, and batch.
- [x] 2.3 Add P79 columns for `p79_status`, `p79_closure_basis`, `p79_fresh_evidence_command`, `p79_fresh_evidence_artifact`, `p79_remaining_gap`, and `p79_next_action`.
- [x] 2.4 Enforce a conservative P79 upgrade whitelist for rows directly covered by real UI data-impact evidence.
- [x] 2.5 Enforce that broad product-goal, monthly-attribution, and expected-return probability/scenario rows cannot be upgraded without direct field-level UI/readback proof.
- [x] 2.6 Emit the P79 matrix and summary JSON.
- [x] 2.7 Add `--check` mode that fails on missing fresh artifacts, invalid upgrades, private absolute paths, or overbroad package/release claims.
- [x] 2.8 Add field-level checker gates for position symbol/name/quantity/cost/buy-reason/tag, confirmation quantity/price, transaction before/after state, evidence refs, audit events, and non-`510300` dynamic symbol fields.

## 3. Fresh Real UI Evidence

- [x] 3.1 Rerun P72 real-user fund scenario under `docs/release/ui-audit-assets/2026-06-21-p79-real-user-fund/`.
- [x] 3.2 Rerun P75 accepted-local non-`510300` UI journey under `docs/release/ui-audit-assets/2026-06-21-p79-non-510300/`.
- [x] 3.3 Validate SQLite readback for portfolio snapshots, positions, confirmations, transactions, decisions, evidence refs, and audit events.
- [x] 3.4 Validate forbidden broker/order/external-push/auto-confirmation artifacts remain absent.
- [x] 3.5 Record expected-return rows that remain partial with exact remaining gaps.
- [x] 3.6 Fix and test ExpectedReturnNode LLM quality-failure safety fallback, then rerun P72 real UI evidence instead of accepting a degraded consultation.

## 4. Acceptance And Release Materials

- [x] 4.1 Add `docs/release/acceptance/2026-06-21-p79-real-use-data-impact-and-expected-return-closure.md`.
- [x] 4.2 Update `docs/release/README.md` with P79 records and conclusion.
- [x] 4.3 Update `docs/release/acceptance-repeatability.md` with P79 repeat commands.
- [x] 4.4 Update release candidate and handoff wording so P79 progress is not confused with full original-requirement pass or P76 package freshness.
- [x] 4.5 Confirm P79 does not expand P76 package claims.

## 5. Review, Verification, Archive, And Commit

- [x] 5.1 Run `python3 scripts/p79_real_use_data_impact_and_expected_return_closure.py --check`.
- [x] 5.2 Run all fresh evidence commands cited by P79 upgraded rows.
- [x] 5.3 Run `openspec validate p79-real-use-data-impact-and-expected-return-closure --strict`.
- [x] 5.4 Run `openspec validate --all --strict`.
- [x] 5.5 Run `git diff --check`.
- [x] 5.6 Run read-only subagent review; fix Critical/Important findings before archive.
- [x] 5.7 Archive P79 and merge release-governance delta.
- [x] 5.8 Confirm no active change remains after archive.
- [x] 5.9 Run final `openspec validate --all --strict`, `python3 scripts/p79_real_use_data_impact_and_expected_return_closure.py --check`, and `git diff --check`.
- [x] 5.10 Commit P79 materials if all required verification passes.
