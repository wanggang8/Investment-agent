# P115 Real User Scenario Acceptance Run

> Date: 2026-06-25  
> Change: `p115-real-user-scenario-acceptance`  
> Result: `passed`  
> Evidence root: `docs/release/ui-audit-assets/2026-06-25-p115-real-user-scenario-acceptance/`

## Summary

P115 executed a fresh local real-user scenario acceptance run against the current source tree with an isolated temporary SQLite database, local Go backend, Vite frontend, API/SQLite runner, and Playwright browser journey.

- Scenario count: 34.
- `fresh_pass`: 28.
- `scoped_pass`: 6.
- `blocked`: 0.
- `degraded_expected`: 0.
- Browser status: `passed`.
- Browser console errors: 0.
- Browser page errors: 0.
- Browser API 5xx responses: 0.

The scoped scenarios are `S09`, `S13`, `S14`, `S16`, `S20`, and `S22`; they are intentionally scoped because the run used `local_seeded_linkage` / local deterministic evidence and does not claim future external provider availability or fresh real LLM output.

## Evidence Artifacts

- Final merged summary: `docs/release/ui-audit-assets/2026-06-25-p115-real-user-scenario-acceptance/p115-scenario-summary.json`.
- API/SQLite summary: `docs/release/ui-audit-assets/2026-06-25-p115-real-user-scenario-acceptance/api_sqlite/p115-api-sqlite-summary.json`.
- Browser summary: `docs/release/ui-audit-assets/2026-06-25-p115-real-user-scenario-acceptance/browser/p115-browser-results.json`.
- Browser screenshots: `docs/release/ui-audit-assets/2026-06-25-p115-real-user-scenario-acceptance/browser/*.png`.

Representative browser evidence includes:

- S01 local install and redaction: `browser/s01-local-install.png`.
- S03 portfolio calibration: `browser/s03-portfolio-calibration.png`.
- S05 batch import: `browser/s05-import-confirm.png`.
- S08 rebalance review: `browser/s08-rebalance.png`.
- S09 consultation surface: `browser/s09-consultation.png`.
- S10 decision detail: `browser/s10-decision-detail.png`.
- S11 manual execution confirmation: `browser/s11-manual-confirmation.png`.
- S11B marked-error flow: `browser/s11b-marked-error.png`.
- S15 local knowledge: `browser/s15-local-knowledge.png`.
- S17 data-quality resolution surface: `browser/s17-dq-resolution.png`.
- S19 rules surface: `browser/s19-rules.png`.
- S29 mobile positions path: `browser/s29-mobile-positions.png`.

## Coverage

P115 covered the following real-use scenario groups:

- First launch, local capability boundary, diagnostics, and redacted local install status.
- Empty account guidance, portfolio initialization, holding edit/remove, batch import, offline transaction recording, local correction audit, and quarterly rebalance review.
- Active consultation surface, decision detail, manual offline execution confirmation, marked-error review loop, and decision-loop traceability.
- Evidence list/verification, local knowledge import, knowledge readiness and local retrieval/index readiness surfaces.
- Market/source-health and data-quality gate resolution flows.
- Risk alert lifecycle, rule proposal/effect validation/final confirmation, notification read state, daily reports, daily auto-run readonly status, dashboard/workbench aggregate state, review, audit, settings, and API diagnostics.
- 390px mobile core positions path.
- Failure and safety paths including invalid portfolio input, missing decision id, settings forbidden mutation, and no-forbidden-automation checks.

## Findings Fixed During Run

P115 exposed one frontend runtime quality issue before the final pass:

- `web/src/components/decision/DecisionTrace.tsx` used insufficiently unique React keys for decision trace lists. Seeded or historical arbitration steps can have empty `rule_id` / repeated `priority`, producing duplicate key warnings such as `0-`. The component now includes the list index and stable fallback labels in keys for expected-return scenarios, holding-class coverage, assumption checks, historical contexts, and arbitration steps.

P115 also corrected the browser acceptance fixture to avoid API/browser state contamination:

- API/SQLite confirmation scenarios use `decision_p115_execute` and `decision_p115_error`.
- Browser confirmation scenarios use separate pending decisions `decision_p115_browser_execute` and `decision_p115_browser_error`.

## Safety Boundary

P115 confirmed the following negative evidence in the final merged summary:

- `forbidden_broker_order_push_tables = 0`.
- `auto_confirmation_rows = 0`.
- `auto_rule_apply_audit_events = 0`.
- `automatic_trading_affordances = 0`.
- `return_guarantee_claims = 0`.
- `secret_or_raw_prompt_leaks_on_primary_ui = 0`.

P115 does not add or claim broker integration, automatic trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, automatic repair/recovery, return guarantees, paid/login/authorized sources, Level2 data, high-frequency data, Docker/package refresh, or physical second-machine validation.

## Claim Boundary

P115 validates current-source local real-user scenario linkage and functional reality evidence. It must not be cited as:

- Fresh external provider availability.
- Fresh real LLM output.
- P93 fresh code-reality pass after P114/P115 changes.
- Full release package refresh.
- Physical second-machine validation.

The final merged summary states: `P115 validates local real-user scenario linkage for current source. P93 may remain stale after P114; this summary must not be cited as fresh P93 code-reality pass.`
