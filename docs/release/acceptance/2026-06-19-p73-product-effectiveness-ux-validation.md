# P73 Product Effectiveness And UX Validation

> Status: `release_ready_product_effectiveness_ux_acceptance`
> Change: `p73-product-effectiveness-ux-validation`

P73 is intended to validate whether the product supports the actual Investment Agent goal: helping the user execute investment discipline through evidence, rules, traceability, review, and safe manual confirmation. It does not use future investment returns as the pass criterion.

## What Was Added

- OpenSpec change `p73-product-effectiveness-ux-validation`.
- Product-effectiveness metrics:
  - discipline adherence;
  - evidence sufficiency;
  - traceability;
  - review usefulness;
  - UX comprehension.
- Browser UX task spec:
  - `web/e2e/p73-product-effectiveness-ux-validation.spec.ts`
- Effect replay checker:
  - `scripts/p73_effect_replay_check.py`
- Local runner:
  - `scripts/p73-product-effectiveness-ux-validation.sh`
- P73 smoke fixture:
  - C-level `background_only` evidence;
  - `non_trade_record` decision showing that background material cannot produce a trade-like suggestion.

## UX Task Evidence

P73 executed real browser UI tasks against a local Go backend, Vite frontend, temporary SQLite database, and seeded deterministic product facts.

| Task | Result | Evidence |
| --- | --- | --- |
| Daily discipline goal | pass | Home page shows today's report, decision detail entry, and decision-loop navigation |
| Portfolio goal | pass | User records local holding facts and sees the no-broker/no-trading boundary |
| Evidence/data quality goal | pass | Data Quality links to Evidence; evidence page distinguishes formal evidence and P73 background material |
| Decision traceability goal | pass | Decision detail shows Agent material, final verdict metadata, ruling chain, expected-return boundary, and no-auto-trading copy |
| Background-only blocking goal | pass | C-level background material remains insufficient, produces a non-trade record, and exposes no execution confirmation |
| Manual confirmation goal | pass | User records a planned offline action; SQLite replay verifies no portfolio transaction is created |
| Risk/review/rule effect goal | pass | Risk, review, and rules surfaces expose read-only tracking and no-auto-rule-application boundary |
| Mobile goal | pass | 390px core pages have reachable navigation and no horizontal overflow |
| Unsafe input goal | pass | Incomplete consultation is blocked with a clear message and does not create a suggestion |

Browser artifacts:

- `docs/release/ui-audit-assets/2026-06-19-p73/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-19-p73/effect-replay-summary.json`
- screenshots under `docs/release/ui-audit-assets/2026-06-19-p73/`

## Effect Replay Evidence

The replay checker passed with these product-effectiveness facts:

- `decision_smoke_p73_background_only` is `record_type=non_trade_record`, `workflow_status=completed`, `source_verification_status=background_only`, `final_verdict_status=insufficient_data`, and `confirmation_status=not_required`.
- `verification_smoke_p73_background` is `evidence_role=background`, `verification_status=background_only`, `highest_source_level=C`, `independent_source_count=0`, and `high_grade_independent_source_count=0`.
- The planned confirmation created through UI for `decision_smoke_p30` has `planned_confirmation_linked_transactions=0`.
- `risk_smoke_p39` exposes prohibited actions including automatic trading and external push.
- `val_smoke_p39` has `validation_status=passed`, `sample_count=6`, `overfit_risk=low`, and a no-auto-rule-application safety note.
- No broker/order/trade-execution/external-push/webhook tables exist.

## UX Audit

| Area | Result | Evidence |
| --- | --- | --- |
| Information hierarchy | pass | Primary pages expose headings, state labels, next actions, and technical trace sections without blank states |
| Next-action clarity | pass | Daily discipline and portfolio pages route users to local maintenance, data quality, decision detail, risk, and consultation actions |
| State labels | pass | Background-only decision clearly separates final verdict, data-insufficient state, and no-confirmation status |
| Navigation | pass | Desktop side navigation and 390px mobile navigation remain reachable |
| Mobile/reflow | pass | P73 mobile route loop found no horizontal overflow beyond the accepted tolerance |
| Copy safety | pass | UI repeatedly states no broker connection, no automatic trading, no automatic confirmation, and no automatic rule application |

Critical findings: none.

Major findings: none.

Minor findings: none requiring code changes. During execution, P73 tightened the browser spec to match current product labels: the UI uses `查看决策详情`, `Agent 分析材料`, `最终裁决明细`, and `工作流状态：已完成`.

Accepted gaps:

- P73 validates product-goal and UX effectiveness with deterministic local facts; it does not claim future investment returns.
- P73 relies on P71/P72 for live public-source and real LLM availability evidence.
- P73 is not a longitudinal real-world user study.
- P73 does not perform a physical second-machine package repeat.
- No post-P72/P73 package refresh has been generated in this change.

## Verification Completed

```text
openspec validate p73-product-effectiveness-ux-validation --strict
Change 'p73-product-effectiveness-ux-validation' is valid

openspec validate --all --strict
Totals: 34 passed, 0 failed

git diff --check
passed

go test ./cmd/smoke-seed
ok investment-agent/cmd/smoke-seed

go test ./...
ok

npm --prefix web test -- --run
48 files passed, 166 tests passed

npm --prefix web run build
production build passed

bash scripts/p73-product-effectiveness-ux-validation.sh
1 browser test passed
effect replay status=passed

strict secret/redaction scan
0 matches

forbidden-boundary string scan
matches classified as release boundaries, no-claim statements, or test assertions only
```

## Current Result

P73 passes for the accepted local product-effectiveness and UX validation scope. This pass is based on actual browser UI operations and deterministic SQLite replay evidence, not mocks and not safe-degradation-only behavior.

## Boundaries

P73 must not claim:

- improved future investment returns;
- future market prediction;
- future public-source or model-provider availability;
- broker connectivity;
- automatic trading;
- one-click trading;
- order delegation;
- external push;
- automatic confirmation;
- automatic rule application;
- automatic repair, migration, restore, or real DB overwrite.
