# P58 UI Acceptance Run

> Date: 2026-06-17
> Change: `p58-daily-workbench-redesign`
> Scope: Dashboard `/` and Workbench `/workbench` daily discipline cockpit redesign.

## Summary

P58 UI acceptance passed for the scoped Dashboard and Workbench redesign.

The app was started with a temporary local config and smoke SQLite database under `tmp/p58-ui/`. The run did not use or overwrite the real local database.

## Commands

| Gate | Command | Result |
| --- | --- | --- |
| Target frontend tests | `npm test -- --run src/features/dashboard src/pages/WorkbenchPage.test.tsx` | Pass: 3 files, 11 tests |
| Full frontend tests | `npm test` | Pass: 33 files, 114 tests |
| Frontend build | `npm run build` | Pass |
| Backend tests | `go test ./...` | Pass |
| Real smoke / E2E | `bash scripts/e2e-smoke.sh` | Pass: 2 Playwright tests |

## Browser Acceptance

| Route | Viewport | Screenshot | Result |
| --- | --- | --- | --- |
| `/` | 1280x900 | `docs/release/ui-audit-assets/2026-06-17-p58/dashboard-desktop.png` | Pass |
| `/` | 390x900 | `docs/release/ui-audit-assets/2026-06-17-p58/dashboard-mobile.png` | Pass |
| `/workbench` | 1280x900 | `docs/release/ui-audit-assets/2026-06-17-p58/workbench-desktop.png` | Pass |
| `/workbench` | 390x900 | `docs/release/ui-audit-assets/2026-06-17-p58/workbench-mobile.png` | Pass |

Recorded assets:

- `docs/release/ui-audit-assets/2026-06-17-p58/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-17-p58/navigation-checks.json`
- `docs/release/ui-audit-assets/2026-06-17-p58/safety-scan.json`

## Findings

- Dashboard and Workbench both expose the daily status hero, manual action queue, and signal grid.
- 390px mobile checks reported `body.scrollWidth` and `documentElement.scrollWidth` equal to the viewport width for both routes.
- No text overlap candidates were detected in the checked desktop or mobile viewports.
- Manual action links navigated to `/positions`, `/data-quality`, `/daily-discipline/reports/daily_report_smoke_p32`, and `/consultation` without a blank page.
- CTA scan found no forbidden execution links or buttons. Dashboard still displays existing risk alert boundary text such as `禁止动作：自动交易、外部推送`; this is a prohibited-action explanation, not an execution affordance.

## Notes

- Browser plugin DOM checks were used against the live Vite app. Screenshot capture through the in-app browser timed out, so screenshot files were captured with the project Playwright runtime against the same live URL.
- Console captured expected local conflict responses from existing smoke data flows; they did not block rendering or navigation.
