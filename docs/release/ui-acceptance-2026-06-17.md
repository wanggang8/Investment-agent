# P55 UI Acceptance Run

> Date: 2026-06-17  
> Result: `blocked` for full UI acceptance  
> Evidence folder: `docs/release/ui-audit-assets/2026-06-17-p55/`

## Scope

This run started the local backend and Vite frontend, operated the UI through the in-app Browser, and captured screenshots for the current primary routes and interactions. It used a temporary SQLite database and temporary VecLite directory under `tmp/ui-acceptance/p55-2026-06-17/`; those runtime files are not release artifacts.

The run does not claim investment returns, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, login sources, paid sources, authorized sources, Level2 data, or high-frequency data.

## Environment

| Item | Value |
| --- | --- |
| Backend | `127.0.0.1:18081` |
| Frontend | `127.0.0.1:14174` |
| Data | `cmd/smoke-seed` on temporary SQLite |
| Browser method | in-app Browser with Playwright-style UI operations |
| Screenshot mode | full-page where stable; 1280x720 clip for long pages |
| Real LLM | temporary config only; no key persisted in release docs |

## Overall Result

Full UI acceptance is blocked by one real-LLM decision-detail rendering defect:

| ID | Severity | Area | Finding | Evidence | Release impact |
| --- | --- | --- | --- | --- | --- |
| P55-B1 | blocked | Decision detail | A real LLM UI consultation created `decision_62160bd3494023dd`, but opening its detail page crashed the React UI because `final_verdict.optional_actions` was `null` and `DecisionTrace` calls `.join()` directly. | Vite log points to `web/src/components/decision/DecisionTrace.tsx:44`; API summary showed `optional_actions_type=NoneType`, `prohibited_actions_type=list`, `analyst_report_count=3`. | Blocks claiming every frontend function passes. P56 should harden DTO rendering and/or normalize API arrays before release. |

Supporting LLM evidence:

- Initial real LLM smoke with the provided gateway root returned HTML and failed JSON parsing.
- Direct `/v1/chat/completions` probe returned valid JSON for model `gpt-5.4-mini`.
- After updating only the temporary config to the `/v1` base path and extending timeout to 90 seconds, `cmd/agent --task llm-smoke --symbol 510300` passed and wrote a sanitized `llm_smoke:quality=passed:parse=parsed:no_auto_trading` audit event.
- The UI consultation after the corrected config wrote a decision with successful `value` and `trend_risk` analysis audit events, then the detail render hit P55-B1.

## Route And Feature Matrix

| Route / feature | Result | Operation evidence | Screenshot |
| --- | --- | --- | --- |
| `/` Dashboard | pass | Loaded today's discipline report, risk alert summary, dashboard cards, and no-auto-trading note. | `01-dashboard.png` |
| `/workbench` | pass | Loaded today/portfolio/rules/consultation sections and navigation entry points. | `02-workbench.png` |
| `/decision-loop` | pass | Loaded read-only decision-loop overview and explanation chain. | `03-decision-loop.png` |
| `/data-quality` | pass with design issue | Loaded source health, Evidence/RAG, LLM quality, and followed "查看证据" to `/evidence`. | `04-data-quality.png`, `28-data-quality-evidence-link.png` |
| `/positions` | pass with design issue | Saved local account and holding calibration; verified local fact success message. | `05-positions-before-calibration.png`, `06-positions-after-calibration.png` |
| `/consultation` seeded/local decision path | pass | Submitted consultation, created decision detail, recorded manual plan confirmation. | `08-consultation-confirmed-clip.png` |
| `/consultation` real LLM path | blocked | Real LLM-backed consultation wrote a decision and analysis audits, but detail render crashed on nullable `optional_actions`. | `33-consultation-real-llm-retest-timeout-state.png` |
| `/decisions/decision_smoke_p30` | pass | Loaded fixture decision with expected return and dynamic sell evaluation. | `09-decision-detail-p30.png` |
| `/decisions/decision_smoke_p39_out_of_scope` | pass | Loaded out-of-scope decision and capability rejection. | `10-decision-detail-out-of-scope.png` |
| `/decisions/decision_smoke_p39_llm_degraded` | pass | Loaded degraded LLM fixture decision and degraded status. | `11-decision-detail-llm-degraded.png` |
| `/evidence` | pass | Loaded evidence summary and reference ID. | `12-evidence.png` |
| `/rules` | pass | Loaded rule proposal, gatekeeper result, and final-confirm boundary copy. | `13-rules.png` |
| `/audit` | pass | Loaded audit list and expanded P30 smoke reference in a scoped row. | `14-audit.png`, `14-audit-expanded.png` |
| `/notifications` | pass | Loaded notification center and marked all notifications read. | `15-notifications.png`, `15-notifications-read.png` |
| `/risk-alerts` | pass | Loaded alert center and prohibited-action safety copy. | `16-risk-alerts.png` |
| `/risk-alerts/risk_smoke_p39` | pass | Loaded alert detail and moved SOP status to observing. | `17-risk-alert-detail.png`, `27-risk-alert-observing.png` |
| `/daily-auto-run` | pass | Loaded disabled/failed local daily auto-run state and safety copy. | `18-daily-auto-run.png` |
| `/daily-discipline/reports` | pass | Loaded report history and navigated into detail. | `19-daily-reports.png` |
| `/daily-discipline/reports/daily_report_smoke_p32` | pass | Loaded report detail and no-auto-trading note. | `20-daily-report-detail.png` |
| `/review` | pass | Loaded review summary, tracking, and rule-change safety copy. | `21-review.png` |
| `/local-install` | pass | Loaded installation diagnostics and edited config draft fields. | `22-local-install.png`, `22-local-install-edited.png` |
| `/local-knowledge` | pass | Validated redacted preview and confirmed local background fact import. | `23-local-knowledge-before.png`, `24-local-knowledge-validated.png`, `25-local-knowledge-confirmed.png` |
| `/settings` | pass | Loaded settings and DeepSeek readiness summary without exposing keys. | `26-settings.png` |

## Mobile Checks

Viewport: 390x844.

| Route | Result | Evidence |
| --- | --- | --- |
| `/workbench` | needs_optimization | No horizontal overflow, but the fixed sidebar consumes too much mobile width and leaves content cramped. Screenshot: `29-mobile-workbench.png`. |
| `/data-quality` | needs_optimization | Horizontal overflow observed: body/document width 528px on 390px viewport. Screenshot: `30-mobile-data-quality.png`. |
| `/risk-alerts` | pass | No horizontal overflow; buttons remained reachable. Screenshot: `31-mobile-risk-alerts.png`. |
| `/positions` | needs_optimization | Horizontal overflow observed: body/document width 729px on 390px viewport, mainly from inline form/table layout. Screenshot: `32-mobile-positions.png`. |

## Safety Boundary

Observed UI controls did not expose complete API keys, broker order actions, one-click trading, order delegation, automatic rule application, automatic confirmation, external push, or return promises. The real LLM key was used only in the temporary runtime config and is not included in committed release docs.

## Follow-Up

P56 should fix P55-B1 before another release-ready claim. It should also address the mobile overflow and cramped mobile navigation/design issues recorded in `docs/release/ui-design-audit-2026-06-17.md`.
