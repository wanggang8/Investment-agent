# P75 Requirements Traceability And Real Use Closure

## Why

P71-P74 proved important local acceptance slices: strict current-data pass, healthy VecLite retrieval, real LLM-backed UI operation, a real `510300` user scenario, product-effectiveness UX validation, and built-in knowledge/data-readiness visibility.

Those stages still do not prove that every original requirement in `docs/requirements.md` has been fulfilled for real product use. The current release claims are scoped to accepted paths and can still hide gaps such as:

- original requirements not mapped one-by-one to code, UI, data, and acceptance evidence;
- a real user adding or holding an arbitrary fund/ETF that is not the `510300` accepted path;
- external data being partially collected but not complete for all required categories;
- built-in master wisdom existing in a registry but not consistently used by runtime LLM/workflow decisions;
- fund-side data and tracked-index-side data not being resolved and joined dynamically from user input;
- expected-return, risk-alert, liquidity, emotion, SOP, and decision-impact calculations not being accuracy-checked end to end for real scenarios;
- UI operation acceptance proving route usability but not proving every task flow, state transition, data mutation, readback, and audit link;
- release materials using scoped pass language that can be mistaken for full product completion.

P75 creates a formal traceability and closure phase before any broader full-product claim is made.

## What Changes

- Create a requirement traceability matrix that maps each requirement block in `docs/requirements.md` to implementation evidence, UI evidence, data evidence, scenario evidence, data-impact evidence, and release-claim status.
- Classify every requirement as `real_pass`, `scoped_pass`, `deterministic_local_evidence`, `partial`, `not_implemented`, or `blocked`.
- Add a real-use acceptance matrix covering the user's concerns: dynamic fund onboarding, external data lookup based on the user-entered fund, built-in knowledge/LLM usage, collected-data completeness, risk/alert/expected-return accuracy, cross-page data impact, UI design reasonableness, and release-claim boundaries.
- Identify and fix or explicitly block gaps where the product behaves like a demo, fixture-only path, hardcoded accepted symbol, or scoped acceptance without full-product evidence.
- Add targeted runtime hardening only where needed to close audited gaps, especially around dynamic symbol profile resolution, fund/index data dependency resolution, LLM knowledge-context unification, and acceptance evidence.
- Produce a final P75 acceptance record that states whether the product is genuinely full-requirement ready, still scoped-ready, or release-blocked.

## In Scope

- OpenSpec change, design, tasks, release-governance delta, traceability evidence, acceptance scripts, docs, tests, and release materials.
- Read-only or safety-preserving runtime improvements needed to close traceability gaps.
- Requirement coverage for:
  - product goals and non-goals;
  - safety boundaries and executable rule criteria;
  - user scenarios and UI operation flows;
  - multi-agent decision workflow and rule-first verdicts;
  - structured/unstructured data requirements and source grading;
  - master wisdom, conflict rules, and LLM context;
  - SOP A-F;
  - expected return, dynamic sell evaluation, risk alerts, and liquidity checks;
  - portfolio/account sync, manual confirmations, data mutations, readback, and audit events;
  - daily/monthly/quarterly review, error cases, evolution proposals, and gatekeeper audit;
  - UI/UX design acceptance and release-claim wording.
- At least one non-`510300` real or locally-supported fund/ETF scenario if public data and symbol profile resolution can be made safe; otherwise the gap must remain `blocked` with evidence and a remediation plan.

## Out of Scope

- No broker interface, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic upgrade, automatic restore, or automatic overwrite of real user databases.
- No investment return promise, future market direction claim, or proof of future investment performance.
- No paid, login-gated, authorization-gated, Level2, high-frequency, or access-control-bypassing source.
- No physical second-machine repeat or new final package refresh unless explicitly added by a later change.
- No fake pass through fixture-only data, scope exclusion, waiver, stale current data, C-level background notes, or hardcoded accepted-symbol shortcuts.

## Impact

P75 may modify traceability docs, release docs, acceptance scripts, tests, backend readiness/resolution services, workflow context construction, frontend readback surfaces, and deterministic acceptance fixtures. It must not change user portfolios, confirmations, rules, market facts, source health, evidence facts, or local databases except through explicit UI/API actions inside acceptance scenarios that verify those mutations.
