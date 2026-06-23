# Proposal: P80 Review Audit Governance Real-Use Closure

## Why

P79 deliberately did not upgrade broad audit and review rows because the evidence only proved local account data impact. The next highest-value real-use gap is the user governance loop: marking a decision as wrong, seeing it in review, sending a rule proposal through gatekeeper audit, and reading the resulting audit trail.

P80 turns the existing real browser SOP/failure-state journey into a stricter field-level evidence layer. It only upgrades rows directly covered by real UI operations and SQLite readback for error cases, rule proposals, gatekeeper audits, audit event fields, and review readback.

## What Changes

- Add a P80 checker that reads the P79 matrix and emits a P80 matrix, acceptance record, and summary JSON.
- Rerun the P75 SOP/failure-state real browser journey under a fresh P80 artifact directory.
- Extract field-level SQLite evidence from the temporary P75 SOP database into committed redacted JSON.
- Upgrade only rows whose entire claim is covered by the fresh P80 evidence.
- Keep broad monthly attribution, full rule application, data-source, SOP A-F completeness, expected-return, and full original-requirement rows non-`real_pass` unless directly proven.

## Scope

In scope:

- Review and audit readback for error marking, rule proposal confirmation, gatekeeper audit, and audit event fields.
- SQLite readback for `error_cases`, `operation_confirmations`, `rule_proposals`, `gatekeeper_audits`, `audit_events`, `risk_alerts`, and absence of broker/order/external-push tables.
- Release materials and progress docs that honestly record P80 row counts.

Out of scope:

- Full original-requirement pass.
- P76 package refresh or any claim that P76 includes P80 evidence.
- Real external provider expansion, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, provider availability promises, or investment return promises.
- Upgrading broad monthly attribution or rule-application-time rows unless P80 proves those exact fields.

## Acceptance

P80 passes only if:

- Fresh P80 real UI evidence artifacts exist and pass field-level SQLite/readback checks.
- `python3 scripts/p80_review_audit_governance_closure.py --check` passes.
- `openspec validate p80-review-audit-governance-real-use-closure --strict` passes.
- `openspec validate --all --strict` passes.
- `git diff --check` passes.
- A read-only review finds no Critical or Important issue before archive.
