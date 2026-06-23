# P56 UI Acceptance Run

> Date: 2026-06-17
> Result: `pass` for P56 scope
> Evidence folder: `docs/release/ui-audit-assets/2026-06-17-p56/`

## Scope

P56 re-ran the UI acceptance paths blocked in P55 after fixing the real LLM decision-detail nullable DTO crash and productizing the highest-risk frontend surfaces.

This run used a temporary SQLite database and temporary runtime config outside the repository. The temporary config included the test LLM gateway supplied for this task, but no key or raw prompt is included in committed documents or screenshots.

## Result

| Gate | Result | Evidence |
| --- | --- | --- |
| Real LLM smoke | pass | `cmd/agent --task llm-smoke --symbol 510300` completed and wrote audit events; no trade execution. |
| Real UI consultation | pass | Browser submitted consultation for 510300 and rendered `decision_041e6e401a86c9d6`. Screenshot: `08-consultation-real-llm-result.png`. |
| Real LLM decision detail | pass | `/decisions/decision_041e6e401a86c9d6` rendered final verdict, prohibited actions and optional actions without page crash. Screenshot: `09-decision-detail-real-llm.png`. |
| Nullable DTO regression | pass | `DecisionTrace.test.tsx` covers `optional_actions: null`, `prohibited_actions: null`, and nullable analyst arrays. |
| Mobile overflow | pass | `/workbench`, `/data-quality`, `/positions`, `/risk-alerts` each reported `body=390`, `root=390`, `viewport=390` at 390px width. |
| Mobile navigation | pass | Mobile navigation is collapsed behind the `导航` button and opens grouped primary navigation. |
| Safety scan | pass | Browser scan found no one-click trading, auto order, order delegation, broker order, return-promise or exposed-key affordance. |
| Browser errors | pass | Final browser run recorded no unexpected `pageerror` or console error. |

## Route Matrix

| Route / feature | Result | Screenshot |
| --- | --- | --- |
| `/` Dashboard | pass | `01-dashboard.png` |
| `/workbench` | pass | `02-workbench.png` |
| `/decision-loop` | pass | `03-decision-loop.png` |
| `/data-quality` | pass | `04-data-quality.png` |
| `/positions` | pass | `05-positions-before.png`, `06-positions-after-calibration.png` |
| `/consultation` | pass | `07-consultation-form.png`, `08-consultation-real-llm-result.png` |
| `/decisions/decision_041e6e401a86c9d6` | pass | `09-decision-detail-real-llm.png` |
| `/decisions/decision_smoke_p30` | pass | `10-decision-detail-fixture.png` |
| `/evidence` | pass | `11-evidence.png` |
| `/rules` | pass | `12-rules.png` |
| `/audit` | pass | `13-audit.png` |
| `/notifications` | pass | `14-notifications.png` |
| `/risk-alerts` | pass | `15-risk-alerts.png` |
| `/daily-auto-run` | pass | `16-daily-auto-run.png` |
| `/daily-discipline/reports` | pass | `17-daily-reports.png` |
| `/daily-discipline/reports/daily_report_smoke_p32` | pass | `18-daily-report-detail.png` |
| `/review` | pass | `19-review.png` |
| `/local-install` | pass | `20-local-install.png` |
| `/local-knowledge` | pass | `21-local-knowledge.png` |
| `/settings` | pass | `22-settings.png` |

## Mobile Evidence

| Route | Result | Screenshot |
| --- | --- | --- |
| `/workbench` | pass, no page-level horizontal overflow | `23-mobile-workbench.png`, `23-mobile-workbench-nav-open.png` |
| `/data-quality` | pass, no page-level horizontal overflow | `24-mobile-data-quality.png`, `24-mobile-data-quality-nav-open.png` |
| `/positions` | pass, no page-level horizontal overflow | `25-mobile-positions.png`, `25-mobile-positions-nav-open.png` |
| `/risk-alerts` | pass, no page-level horizontal overflow | `26-mobile-risk-alerts.png`, `26-mobile-risk-alerts-nav-open.png` |

## Verification Commands

| Command | Result |
| --- | --- |
| `npm test -- --run src/components/decision/DecisionTrace.test.tsx src/pages/PortfolioPage.test.tsx src/pages/DataQualityPage.test.tsx src/pages/DecisionDetailPage.test.tsx` | pass, 21 tests |
| `npm test` | pass, 32 files / 111 tests |
| `npm run build` | pass |
| `go test ./...` | pass |
| `E2E_BASE_URL=http://127.0.0.1:14175 npm run test:e2e` | pass, 2 tests |

## Notes

- The first Playwright E2E rerun exposed a real-LLM timing mismatch: the P39 journey used a 30s total timeout while real consultation can exceed that. P56 updated the test to allow 180s for the real LLM journey and 150s for final verdict wait.
- P56 does not change the safety boundary: no broker connection, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair promise, real database overwrite, or return promise was added.
