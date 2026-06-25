# P119 Full UI Control And Affordance Acceptance Matrix

Date: 2026-06-25

Change: `p119-full-ui-control-and-affordance-acceptance`

Status: `passed_pending_archive`

## Scope

P119 answers whether all visible operation buttons, forms, links and UI states need acceptance after P114-P118. It explicitly covers the whole current frontend route map rather than only the homepage.

## Route Matrix

| ID | Route | Primary surface | Required evidence |
| --- | --- | --- | --- |
| R01 | `/` | 今日纪律 | desktop inventory, mobile layout, nav/link classification |
| R02 | `/workbench` | 用户决策工作台 | desktop inventory, mobile layout, cross-page links, topbar refresh toggle |
| R03 | `/decision-loop` | 决策闭环解释 | desktop inventory, loop links, no raw JSON |
| R04 | `/data-quality` | 数据质量可观测 | desktop/mobile inventory, symbol filter, details toggles, resolution create/retire/check |
| R05 | `/positions` | 组合与持仓维护 | desktop/mobile inventory, form alignment, select toggles, portfolio writes/readback |
| R06 | `/consultation` | 主动咨询 | scenario select, form validation, generated decision link, no automatic action |
| R07 | `/decisions` | 决策入口 | navigation classification |
| R08 | `/decisions/:decisionId` | 决策详情 | confirmation buttons, evidence/analysis expand-collapse, confirmation readback |
| R09 | `/evidence` | 情报与证据 | refresh/rebuild buttons, evidence role filter, row expand-collapse |
| R10 | `/rules` | 规则与纪律 | details toggles, SOP proposal, confirm/final-confirm controls |
| R11 | `/audit` | 审计检查 | status filter, detail toggles, no raw secret/debug output |
| R12 | `/notifications` | 本地通知 | mark-one/mark-all readback |
| R13 | `/risk-alerts` | 风险预警中心 | SOP lifecycle buttons/readback |
| R14 | `/risk-alerts/:alertId` | 风险预警详情 | detail scoped lifecycle buttons/readback |
| R15 | `/daily-auto-run` | 每日自动运行 | read-only links and status copy |
| R16 | `/daily-discipline/reports` | 每日纪律报告历史 | report links and list layout |
| R17 | `/daily-discipline/reports/:reportId` | 每日纪律报告详情 | detail links and mobile layout |
| R18 | `/review` | 复盘摘要 | review links, no raw debug output |
| R19 | `/local-install` | 本地安装诊断 | config/details toggles, clear summary, startup draft visibility boundary |
| R20 | `/local-knowledge` | 本地知识导入 | structured-record details toggle, validate/confirm flow and SQLite readback |
| R21 | `/settings` | 设置 | market refresh button/readback |
| R22 | `/api-diagnostics` | 接口诊断 | links and product boundary |

## Control Categories

| Category | Must prove |
| --- | --- |
| Navigation | Link has visible label and internal href opens a production route or known detail route. |
| Light interaction | Details/toggle controls work without page error or layout overflow. |
| Read action | Refresh/check operations call implemented API and surface result safely. |
| Local fact write | Operation creates or updates SQLite/API-backed local facts. |
| Governance confirmation | Rule proposal actions transition through explicit user/gatekeeper/final confirmation states. |
| Expected disabled | Disabled controls have prerequisite state and do not imply broken functionality. |
| Boundary notice | Unsupported broker/order/auto paths are explicit as unavailable. |

## Initial Assertions

- No unclassified visible interactive control.
- No unnamed visible interactive control.
- No blank route.
- No obvious viewport overflow in desktop or mobile scans.
- No browser console errors, page errors, or API 5xx responses.
- Upstream/light-toggle interactions have before/after assertions rather than inventory-only evidence.
- No visible one-click trading, broker order placement, order delegation, external push execution, automatic confirmation, automatic rule application, or return guarantee affordance.
- Key UI writes have SQLite readback; otherwise UI or backend must be fixed before P119 can pass.

## Completion Record

Final execution evidence:

- `docs/release/ui-audit-assets/2026-06-25-p119-full-ui-control-and-affordance-acceptance/p119-ui-control-summary.json`
- `docs/release/ui-audit-assets/2026-06-25-p119-full-ui-control-and-affordance-acceptance/browser/p119-browser-results.json`
- `docs/release/acceptance/2026-06-25-p119-full-ui-control-and-affordance-acceptance.md`

Runner result:

- Routes: 22/22.
- Desktop control inventory: 603 controls.
- Mobile route checks: 8.
- Unnamed controls: 0.
- Unclassified controls: 0.
- Layout issues: 0.
- Product-copy issues: 0.
- Upstream/light-toggle interactions exercised: 24.
- Toggle issues: 0.
- Browser console errors: 0.
- Page errors: 0.
- API 5xx responses: 0.
- SQLite readback passed for portfolio, offline transaction, decision confirmations, marked error, risk SOP, notification read state, data-quality resolution, rule proposal, local knowledge import, evidence chunks and safety counters.
