# Tasks: P77 Requirements Real Pass Upgrade Gate

## 1. Setup And Governance

- [x] 1.1 Confirm no active OpenSpec change exists before creating P77.
- [x] 1.2 Create OpenSpec change `p77-requirements-real-pass-upgrade-gate`.
- [x] 1.3 Define P77 as a post-P75 upgrade evidence layer, not a rewrite of P75 history.
- [x] 1.4 Mark P77 active in `docs/GOVERNANCE.md`, `openspec/project.md`, `openspec/PROGRESS.md`, `AGENTS.md`, and `docs/development-plan.md`.
- [x] 1.5 Run `openspec validate p77-requirements-real-pass-upgrade-gate --strict`, `openspec validate --all --strict`, and `git diff --check`.

## 2. Upgrade Gate And Matrix Generator

- [x] 2.1 Implement a P77 matrix generator that reads `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md`.
- [x] 2.2 Preserve each P75 `requirement_id`, source section, source lines, requirement text hash, original status, full-release-required flag, and release impact.
- [x] 2.3 Add P77 columns for `p77_status`, `upgrade_basis`, `gate_dimensions`, `fresh_evidence_command`, `fresh_evidence_artifact`, `residual_gap`, and `next_remediation`.
- [x] 2.4 Enforce that a row cannot be `real_pass` when evidence is screenshot-only, route-smoke-only, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, temporary-DB-only, or incompatible single-symbol-only.
- [x] 2.5 Emit `docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md`.
- [x] 2.6 Emit `docs/release/ui-audit-assets/2026-06-21-p77/real-pass-upgrade-summary.json`.
- [x] 2.7 Add a `--check` mode that fails on missing required columns, invalid status transitions, overbroad full-pass conclusions, or missing evidence artifacts for upgraded rows.

## 3. Evidence Reruns And First-Batch Upgrade

- [x] 3.1 Rerun real UI SOP/failure-state acceptance if it is used as P77 evidence.
- [x] 3.2 Rerun non-`510300` accepted-local real UI journey if it is used as P77 evidence.
- [x] 3.3 Rerun safety/forbidden capability scans if safety rows are upgraded.
- [x] 3.4 Rerun targeted deterministic backend checks for source verification, data-impact, and no-forbidden-runtime behavior when those rows are upgraded.
- [x] 3.5 Record all command outputs, summary artifacts, and limitations in the P77 acceptance record.

## 4. Acceptance And Release Materials

- [x] 4.1 Add `docs/release/acceptance/2026-06-21-p77-real-pass-upgrade-acceptance.md`.
- [x] 4.2 Update `docs/release/README.md` to include the P77 acceptance record and conclusion.
- [x] 4.3 Update `docs/release/acceptance-repeatability.md` with the P77 rerun command.
- [x] 4.4 Update `docs/release/release-candidate-2026-06-18.md` and `docs/release/release-handoff-2026-06-18.md` only if wording needs to avoid stale full-pass or package claims.
- [x] 4.5 Confirm P77 does not expand P76 package claims.

## 5. Review, Verification, And Archive

- [x] 5.1 Run `python3 scripts/p77_requirements_real_pass_upgrade.py --check`.
- [x] 5.2 Run all fresh evidence commands cited by upgraded rows.
- [x] 5.3 Run `openspec validate p77-requirements-real-pass-upgrade-gate --strict`.
- [x] 5.4 Run `openspec validate --all --strict`.
- [x] 5.5 Run `git diff --check`.
- [x] 5.6 Run subagent review; fix Critical/Important findings before archive.
- [x] 5.7 Archive P77 and merge release-governance delta.
- [x] 5.8 Confirm no active change remains after archive.
- [x] 5.9 Run final `openspec validate --all --strict` and `git diff --check`.
- [x] 5.10 Commit P77 materials if all required verification passes.
