# P73 Design

## Product Goal Being Validated

The product goal is not "predict the market." It is a personal AI investment discipline assistant that helps the user:

- maintain local account and holding facts;
- inspect current discipline, risk, evidence, and rule state;
- ask about self-selected in-scope symbols;
- receive traceable LLM analysis plus rule-based final verdicts;
- manually record offline actions;
- review errors and evolve rules through controlled gates.

P73 validates whether the product experience supports that goal.

## Validation Layers

### Layer 1: Product Goal Metrics

P73 records metrics that map directly to the goal:

| Metric | Meaning | Pass signal |
| --- | --- | --- |
| Discipline adherence | Unsafe or unsupported actions are blocked or made non-trade records | blocked/gap states appear where required |
| Evidence sufficiency | Formal decisions distinguish formal evidence from background | satisfied/background_only/failed shown correctly |
| Traceability | User can follow current facts -> evidence -> LLM reports -> rules -> final verdict -> manual confirmation | UI readbacks expose the chain |
| Review usefulness | Errors, confirmations, risk alerts, and audit facts can be read back | review/audit/decision-loop pages link scenario facts |
| UX comprehension | Page hierarchy and copy help the user choose the next safe action | no critical confusion findings |

### Layer 2: UX Task Matrix

The browser journey should operate tasks a real target user would naturally attempt:

1. First-use or missing-prerequisite guidance.
2. Current daily discipline review.
3. Portfolio setup/maintenance.
4. Data quality and evidence sufficiency review.
5. Active consultation and decision detail review.
6. Manual offline confirmation and decision-loop readback.
7. Risk alert, notification, audit, rules, and review readback.
8. Rule-effect or proposal understanding where available.
9. Invalid or unsafe action handling.

### Layer 3: Effect Replay

Effect replay is not a return backtest. It verifies discipline behavior:

- C-level background-only material must not satisfy formal evidence.
- Insufficient evidence must produce safe blocking or non-trade records.
- Rule proposal validation must expose sample/overfit/gate status.
- Risk alerts must surface in daily/review/notification/audit surfaces.
- Manual confirmations must be the only path that mutates local portfolio facts.

### Layer 4: UX Audit

The audit checks:

- whether the first viewport communicates state and next action;
- whether decision detail distinguishes LLM analysis, formal evidence, and rule final verdict;
- whether data quality and evidence pages make insufficiency understandable;
- whether navigation supports repeated operational use;
- whether mobile/reflow paths remain usable;
- whether copy avoids broker/trading/return-promise ambiguity.

## Evidence Artifacts

P73 writes:

- `docs/release/acceptance/2026-06-19-p73-product-effectiveness-ux-validation.md`
- `docs/release/ui-audit-assets/2026-06-19-p73/`
- sanitized browser result JSON;
- UX audit notes;
- effect replay summary JSON;
- screenshots for representative tasks.

## Safety

P73 must preserve all P71/P72 safety boundaries. When evidence, LLM, RAG, source health, or UI comprehension is inadequate, the result is blocked or gap-qualified. It must not be called pass by relying on safe degradation alone.
