# Tasks: P78 Requirements Real-Pass Batch Closure

## 1. Setup And Governance

- [x] 1.1 Confirm no other active OpenSpec change exists before executing P78.
- [x] 1.2 Create OpenSpec change `p78-requirements-real-pass-batch-closure`.
- [x] 1.3 Define P78 as a post-P77 batch closure evidence layer, not a rewrite of P75/P77 history.
- [x] 1.4 Mark P78 active in `docs/GOVERNANCE.md`, `openspec/project.md`, `openspec/PROGRESS.md`, `AGENTS.md`, and `docs/development-plan.md`.
- [x] 1.5 Run `openspec validate p78-requirements-real-pass-batch-closure --strict`, `openspec validate --all --strict`, and `git diff --check`.

## 2. Batch Classifier And Matrix Generator

- [x] 2.1 Implement `scripts/p78_requirements_real_pass_batch_closure.py`.
- [x] 2.2 Read the P77 matrix and preserve row identity, source section, source lines, requirement text, original status, P77 status, full-release-required flag, and release impact.
- [x] 2.3 Classify remaining full-release-required non-`real_pass` rows by remediation group and batch.
- [x] 2.4 Add P78 columns for `p78_status`, `remediation_group`, `batch`, `closure_basis`, `fresh_evidence_command`, `fresh_evidence_artifact`, `remaining_gap`, and `next_action`.
- [x] 2.5 Enforce that P78 cannot claim full original-requirement pass while any full-release-required row remains non-`real_pass`.
- [x] 2.6 Emit `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-matrix.md`.
- [x] 2.7 Emit `docs/release/ui-audit-assets/2026-06-21-p78/real-pass-batch-summary.json`.
- [x] 2.8 Add `--check` mode that fails on missing evidence artifacts, invalid upgrades, missing expected-return readback fields, or overbroad package/release claims.

## 3. P78 Batch A Evidence

- [x] 3.1 Rerun fresh expected-return Go tests and store verbose output under `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.log`.
- [x] 3.2 Generate `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.json` metadata and validate concrete test names.
- [x] 3.3 Rerun the accepted-local non-`510300` real UI journey under `docs/release/ui-audit-assets/2026-06-21-p78-non-510300/`.
- [x] 3.4 Read the real UI SQLite decision record and emit `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-ui-readback.json`.
- [x] 3.5 Upgrade only the batch A expected-return degradation/disclaimer rows whose deterministic and real UI/readback evidence is complete.

## 4. Acceptance And Release Materials

- [x] 4.1 Add `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-closure.md`.
- [x] 4.2 Update `docs/release/README.md` with P78 records and conclusion.
- [x] 4.3 Update `docs/release/acceptance-repeatability.md` with P78 repeat commands.
- [x] 4.4 Update release candidate and handoff wording so P78 progress is not confused with full original-requirement pass or P76 package freshness.
- [x] 4.5 Confirm P78 does not expand P76 package claims.

## 5. Review, Verification, Archive, And Commit

- [x] 5.1 Run `python3 scripts/p78_requirements_real_pass_batch_closure.py --check`.
- [x] 5.2 Run all fresh evidence commands cited by P78 upgraded rows.
- [x] 5.3 Run `openspec validate p78-requirements-real-pass-batch-closure --strict`.
- [x] 5.4 Run `openspec validate --all --strict`.
- [x] 5.5 Run `git diff --check`.
- [x] 5.6 Run read-only subagent review; fix Critical/Important findings before archive.
- [x] 5.7 Archive P78 and merge release-governance delta.
- [x] 5.8 Confirm no active change remains after archive.
- [x] 5.9 Run final `openspec validate --all --strict`, `python3 scripts/p78_requirements_real_pass_batch_closure.py --check`, and `git diff --check`.
- [x] 5.10 Commit P78 materials if all required verification passes.
