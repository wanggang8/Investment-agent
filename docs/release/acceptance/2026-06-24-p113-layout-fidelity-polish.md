# P113 Layout Fidelity Polish Acceptance

Date: 2026-06-24

Change: `p113-layout-fidelity-polish`

Reference image: `/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`

## Result

`pass`

P113 fixes the P112-after layout problems reported by the user: mobile metric rail clipping, desktop report hero compression, tiny action links, raw rule JSON appearing too early, settings error cards pushing the report hero below the fold, and decision-detail report hierarchy in a populated local data state.

P113 is frontend-only. It does not add backend APIs, SQLite schema, Eino workflow behavior, LLM capability, data source capability, investment rules, Docker/package/release behavior, broker integration, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, or return promises.

## Evidence

Evidence directory:

`docs/release/ui-audit-assets/2026-06-24-p113-layout-fidelity-polish/`

Captured screenshots:

- 19 desktop route screenshots at 1492x1068.
- 19 mobile route screenshots at 390x844.
- 2 additional populated decision-detail screenshots, captured after seeding only the temporary P113 SQLite database.

Metrics:

- `layout-metrics-after.json`: 38 route/viewport checks.
- `decision-detail-real-metrics.json`: populated decision detail desktop/mobile checks.

## Page Matrix

| Route | Desktop | Mobile | Notes |
| --- | --- | --- | --- |
| `/` | pass | pass | No horizontal overflow; action link height fixed. |
| `/workbench` | pass | pass | No horizontal overflow; action link height fixed. |
| `/decision-loop` | pass | pass | No clipped first-viewport content. |
| `/data-quality` | pass | pass | Mobile metric rail replaced by compact grid; desktop hero no longer causes overflow. |
| `/positions` | pass | pass | Stable report hero and actions. |
| `/consultation` | pass | pass | Form remains responsive. |
| `/decisions/:id` empty/local missing state | pass | pass | No overflow or tiny controls. |
| `/decisions/:id` populated state | pass | pass | Populated detail hierarchy verified with temporary local data. |
| `/evidence` | pass | pass | No horizontal overflow outside table container. |
| `/rules` | pass | pass | Current-rule JSON and proposal raw details are folded behind explicit detail controls. |
| `/audit` | pass | pass | No tiny mobile action target. |
| `/notifications` | pass | pass | Stable mobile controls. |
| `/risk-alerts` | pass | pass | Stable mobile controls. |
| `/daily-auto-run` | pass | pass | Stable mobile controls. |
| `/daily-discipline/reports` | pass | pass | Stable mobile controls. |
| `/daily-discipline/reports/:id` | pass | pass | Stable empty/detail state. |
| `/review` | pass | pass | Stable report layout. |
| `/local-install` | pass | pass | Mobile metric grid no longer clips. |
| `/local-knowledge` | pass | pass | Mobile metric grid no longer clips. |
| `/settings` | pass | pass | Error notices moved below report hero and compacted; metrics no longer clip. |

## Rendered Checks

Browser path:

- Browser plugin was used for route loading, console health, DOM metrics, and the 38-route screenshot set.
- Browser checks passed for page identity, nonblank content, no framework overlay, console health, screenshot evidence, and rendered DOM metrics.

Fallback:

- Browser screenshot capture timed out only for the extra populated decision-detail screenshot.
- The populated decision-detail desktop/mobile screenshots were captured with the project Playwright CLI as a local fallback.
- Browser DOM metrics were still used for the populated decision-detail no-overflow/touch-target checks.

Metric result:

```json
{
  "checked": 38,
  "issueCount": 0,
  "issues": []
}
```

Populated decision-detail metric result:

```json
[
  {
    "viewportName": "desktop",
    "overflow": false,
    "offscreen": [],
    "smallActions": []
  },
  {
    "viewportName": "mobile",
    "overflow": false,
    "offscreen": [],
    "smallActions": []
  }
]
```

## Commands

- `openspec validate p113-layout-fidelity-polish --strict`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- Browser route screenshot/DOM metric loop against `http://127.0.0.1:14113/`
- Playwright fallback screenshots for populated decision detail

Additional final gates are recorded in the P113 task checklist before archive.
