# P119 Design

## Evidence Shape

P119 使用三层证据：

1. API/SQLite seed and readback
   - 使用临时 SQLite 和本地 Go backend。
   - 写入或生成足够的真实本地事实，让每个页面都有可见数据、按钮和异常边界。
   - 通过 SQLite 统计确认 UI 操作产生了真实本地事实。

2. Browser control inventory
   - 使用 Playwright 在真实 Vite 页面中遍历生产路由。
   - 逐路由采集 `button`、`a[href]`、`input`、`select`、`textarea`、`summary`。
   - 记录 label、role/tag、disabled、href、route、classification、viewport、布局检查结果。

3. Product and visual QA
   - 桌面 1440px 逐路由检查。
   - 移动 390px 覆盖高密度页面：今日纪律、工作台、持仓、决策详情、数据质量、风险预警、本地知识、设置。
   - DOM 边界扫描包括水平溢出、按钮/输入控件越界、无名称控件、空白页、明显 raw/debug/mock/stub/secret 文案。

## Route Coverage

P119 覆盖当前 `web/src/App.tsx` 中的生产路由：

- `/`
- `/workbench`
- `/decision-loop`
- `/data-quality`
- `/positions`
- `/consultation`
- `/decisions`
- `/decisions/:decisionId`
- `/evidence`
- `/rules`
- `/audit`
- `/notifications`
- `/risk-alerts`
- `/risk-alerts/:alertId`
- `/daily-auto-run`
- `/daily-discipline/reports`
- `/daily-discipline/reports/:reportId`
- `/review`
- `/local-install`
- `/local-knowledge`
- `/settings`
- `/api-diagnostics`

## Control Classification

- `navigation`: internal page links and external-safe read-only links.
- `light_interaction`: details expand/collapse, evidence chain expand, mobile nav toggle.
- `read_action`: refresh/check operations that read or rebuild local material without user portfolio/rule mutation.
- `write_local_fact`: local account, decision confirmation, notification, risk SOP, data-quality resolution, local knowledge import.
- `governance_confirm`: rule proposal user confirmation or final confirmation.
- `disabled_expected`: disabled until prerequisite is satisfied.
- `boundary_notice`: visible text explaining unsupported broker/order/auto paths.

Any visible interactive element that cannot be named and classified fails the run.

## Backend Consistency Rule

If the UI presents an operation as possible, P119 requires one of:

- real API and SQLite readback evidence,
- explicit disabled state with prerequisite,
- explicit product boundary text showing it is not available,
- read-only navigation/light interaction classification.

If this rule fails, P119 should fix the UI wording/state or the backend implementation, depending on which layer is inconsistent.

## Safety Boundary

P119 must not add hidden broker/order/push tables, automatic confirmations, automatic rule application, or any wording that implies guaranteed returns. Negative evidence remains part of the summary.
