# Design: P80 Review Audit Governance Real-Use Closure

## Evidence Model

P80 is an evidence layer over P79. It reads:

- `docs/release/acceptance/2026-06-21-p79-real-use-data-impact-and-expected-return-matrix.md`
- fresh P80 artifacts under `docs/release/ui-audit-assets/2026-06-22-p80-review-audit-governance/`

It emits:

- `docs/release/acceptance/2026-06-22-p80-review-audit-governance-matrix.md`
- `docs/release/acceptance/2026-06-22-p80-review-audit-governance-closure.md`
- `docs/release/ui-audit-assets/2026-06-22-p80-review-audit-governance/review-audit-governance-summary.json`

P80 does not rewrite P75, P77, P78, or P79 history.

## Upgrade Rules

Rows can move to `real_pass` only when the checker proves the row from fresh evidence. For P80 that means:

- browser result status is `passed`;
- error marking created exactly one local confirmation and one `error_cases` row;
- error case readback includes `decision_id`, `confirmation_id`, `actual_outcome`, `root_cause_tag`, `lesson_learned`, and `created_at`;
- rule proposal confirmation creates a gatekeeper audit and keeps final rule application blocked until final user confirmation;
- gatekeeper audit readback covers pass/reject/needs-user-review states and node audit events including fundamental rule, conflict, backtest, decision, and audit record nodes;
- audit events readback includes `action`, `node_action`, `actor`, `status`, `before_state`, `after_state`, `request_id`, and references to decision/proposal/confirmation/error case where applicable;
- review UI/readback exposes the marked error tag and proposal/governance state;
- forbidden broker/order/external-push tables remain absent.

Rows stay non-`real_pass` if they require:

- monthly profit/loss attribution without direct P&L attribution proof;
- final rule application time without final user confirmation applying a rule version;
- every SOP A-F branch beyond the accepted P80 UI lifecycle evidence;
- external provider coverage, expected-return probability/scenario fields, or full product-level claims.

## Safety

P80 continues the local-only safety boundary. It uses a temporary SQLite database and committed redacted summaries/screenshots. It does not commit the temporary DB, raw prompts, raw provider payloads, complete keys, private absolute paths, or logs containing local secrets.
