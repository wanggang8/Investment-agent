# P58 今日工作台重构

## Why

P57 已固化产品体验打磨路线图，明确 P58 的目标是把 Dashboard 和 Workbench 打磨成真正的“每日投资纪律 cockpit”。当前 `/` 和 `/workbench` 虽然已经聚合今日纪律、组合、风险、规则和复盘信息，但仍存在信息主线分散的问题：用户需要在多个卡片之间自行判断今天能不能动、为什么、下一步做什么，以及数据是否可信。

P58 要把这两个入口重构成日常首屏：先给出今日状态、风险、数据可信度和下一步人工动作，再提供解释与导航。它不新增后端能力，不改变交易边界，只重组前端信息架构、展示组件、测试和验收证据。

## What Changes

- 重构 Dashboard / Workbench 的产品信息架构：
  - 首屏展示今日裁决、风险等级、数据可信度、最近更新时间和下一步人工动作。
  - 把“今日纪律报告、决策摘要、组合风险、数据质量、规则/复盘入口”整合为一个可扫描的每日 cockpit。
  - Dashboard 作为每日入口，Workbench 作为同一主线下的任务视图，避免重复和割裂。
- 新增前端展示模型：
  - 从现有 `DashboardTodayResponse`、`DailyDisciplineReport`、portfolio、risk、rule、review DTO 组合出 `DailyWorkbenchViewModel`。
  - 明确 degraded / insufficient / high risk / unknown 状态的安全显示。
- 强化页面级组件：
  - 今日状态 hero / summary band。
  - 下一步人工动作队列。
  - 数据可信度和风险摘要。
  - 解释入口与任务导航。
- 更新测试与真实 UI 验收：
  - Vitest 覆盖成功、空态、降级态、高风险态、需要人工动作、安全文案和禁用自动交易入口。
  - Playwright 覆盖 `/`、`/workbench` 桌面和 390px 移动端无横向溢出。
  - 真实启动本地后端和 Vite 前端，保存 P58 桌面/移动截图和验收报告。

## Scope

- 前端 React/Vite/TypeScript：
  - `web/src/features/dashboard/*`
  - `web/src/pages/WorkbenchPage.tsx`
  - 相关 dashboard/workbench 测试
  - 必要的轻量 UI/helper 组件或 view model 文件
  - 必要 CSS
- 前端 E2E smoke 适配。
- P58 验收报告、截图资产、OpenSpec 和治理文档更新。

## Out of Scope

- 不修改 SQLite schema、HTTP API、Eino 工作流、Go 后端业务逻辑或真实数据 collector。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺。
- 不新增登录源、付费源、授权源、Level2 或高频源。
- 不把 Dashboard/Workbench 改成营销 landing page 或 AI chat 首页。
- 不提前执行 P59 决策详情解释重构、P60 组合/风险/数据质量深度重构、P61 治理/运维产品化或 P62 组件库全面抽象。

## Validation

- 方案阶段：
  - `openspec validate p58-daily-workbench-redesign --strict`
  - `openspec validate --all --strict`
  - `git diff --check`
  - 敏感信息扫描无 key、完整 prompt、私有 SQLite、raw vendor payload 或临时配置泄露。
  - 子 agent 方案复审无 Critical / Important，且复审覆盖 Product Design 依据。
- 实现阶段：
  - `npm test -- --run src/features/dashboard src/pages/WorkbenchPage.test.tsx`
  - `npm test`
  - `npm run build`
  - `go test ./...`，若无后端修改仍执行以证明集成基线。
  - `E2E_BASE_URL=<local vite url> npm run test:e2e`
  - 真实启动本地后端和前端，通过浏览器操作 `/`、`/workbench`、关键导航和移动端视口。
  - 生成 P58 UI 验收报告和桌面/移动截图。
- 归档前：
  - 子 agent 执行后复审无 Critical / Important。
  - archive 后提交前复审无 Critical / Important。
