# Tasks: P80 Review Audit Governance Real-Use Closure

## 1. Setup And Governance

- [x] 1.1 Confirm no active OpenSpec change exists before creating P80.
- [x] 1.2 Create OpenSpec change `p80-review-audit-governance-real-use-closure`.
- [x] 1.3 Define P80 as a post-P79 evidence layer, not a rewrite of historical matrices.
- [x] 1.4 Mark P80 active in `docs/GOVERNANCE.md`, `openspec/project.md`, `openspec/PROGRESS.md`, `AGENTS.md`, and `docs/development-plan.md`.

## 2. P80 Checker And Matrix

- [x] 2.1 Implement `scripts/p80_review_audit_governance_closure.py`.
- [x] 2.2 Read the P79 matrix and preserve row identity, source lines, requirement text, P79 status, remediation group, and P79 evidence fields.
- [x] 2.3 Add P80 columns for status, closure basis, command, artifact, remaining gap, and next action.
- [x] 2.4 Enforce a conservative P80 upgrade whitelist for rows fully covered by review/audit/governance field evidence.
- [x] 2.5 Keep broad monthly attribution and final rule-application-time rows non-`real_pass` unless exact proof exists.
- [x] 2.6 Emit P80 matrix, acceptance Markdown, and summary JSON.
- [x] 2.7 Add `--check` mode that fails on missing fresh artifacts, invalid upgrades, private absolute paths, raw payload/key leakage, or overbroad release claims.

## 3. Fresh Real UI Evidence

- [x] 3.1 Rerun P75 SOP/failure-state real browser journey under `docs/release/ui-audit-assets/2026-06-22-p80-review-audit-governance/`.
- [x] 3.2 Extract field-level SQLite readback for `error_cases`, `operation_confirmations`, `rule_proposals`, `gatekeeper_audits`, `audit_events`, `risk_alerts`, and forbidden table absence.
- [x] 3.3 Validate review/audit/rules UI readback via browser results and screenshots.
- [x] 3.4 Record remaining gaps for monthly attribution, final rule application time, and full original-requirement pass.

## 4. Acceptance And Release Materials

- [x] 4.1 Add `docs/release/acceptance/2026-06-22-p80-review-audit-governance-closure.md`.
- [x] 4.2 Add `docs/release/acceptance/2026-06-22-p80-review-audit-governance-matrix.md`.
- [x] 4.3 Update `docs/release/README.md`, `docs/release/acceptance-repeatability.md`, release candidate, and handoff wording with exact P80 counts and boundaries.
- [x] 4.4 Confirm P80 does not expand P76 package claims.

## 5. Review, Verification, Archive, And Commit

- [x] 5.1 Run `python3 scripts/p80_review_audit_governance_closure.py --check`.
- [x] 5.2 Run all fresh evidence commands cited by P80 upgraded rows.
- [x] 5.3 Run `openspec validate p80-review-audit-governance-real-use-closure --strict`.
- [x] 5.4 Run `openspec validate --all --strict`.
- [x] 5.5 Run `git diff --check`.
- [x] 5.6 Run read-only review; fix Critical/Important findings before archive.
- [x] 5.7 Archive P80 and merge release-governance delta.
- [x] 5.8 Confirm no active change remains after archive.
- [x] 5.9 Run final `openspec validate --all --strict`, `python3 scripts/p80_review_audit_governance_closure.py --check`, and `git diff --check`.
- [x] 5.10 Commit P80 materials if all required verification passes.
