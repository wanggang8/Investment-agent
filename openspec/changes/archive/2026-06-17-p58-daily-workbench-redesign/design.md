# Design: P58 今日工作台重构

## Product Design Brief

P58 的设计目标是把 `/` 和 `/workbench` 变成每天打开后可立即判断状态的本地投资纪律 cockpit。用户不应在多个卡片里拼答案，而应在第一屏直接看到：

1. 今天能不能动。
2. 为什么是这个结论。
3. 下一步需要哪些人工动作。
4. 数据、风险和规则是否可信。

本阶段使用 Product Design get-context / P55-P56 audit / P57 roadmap 作为设计输入。设计依据继续采用 dashboard at-a-glance、表单认知负担、WCAG reflow、金融解释性和信任设计资料。契约真源仍是 OpenSpec change 和 `docs/`。

## Current Problems

- Dashboard 以 `CockpitLayout` 三栏展示系统、决策和证据，但“今日能不能动”的主结论不够突出。
- Workbench 已聚合多源数据，但四张卡片并列，用户仍要自行综合今日纪律、组合风险、规则复盘和咨询入口。
- Dashboard 与 Workbench 有重复入口：今日纪律、持仓/风险、规则/复盘、主动咨询都存在，但层级和第一屏重点不同。
- 当前页面缺少明确的“下一步人工动作队列”，不能把查看解释、维护持仓、处理风险、查看数据质量等任务按优先级组织。
- P56 已解决基础移动 reflow，但 P58 需要验证 `/` 和 `/workbench` 的第一屏在 390px 下仍能读出今日状态。

## Approach Options

### Option A: Shared daily cockpit model and shell

建立 shared view model + shared daily cockpit component。Dashboard 和 Workbench 使用同一套今日状态、人工动作和摘要组件；Dashboard 偏今日总览，Workbench 偏任务展开。

优点：减少重复、信息主线一致、测试集中。缺点：需要提取一层 view model，触及两个页面。

### Option B: Only restyle existing cards

保留当前 Dashboard/Workbench 结构，只通过 CSS 和局部文案强调重点。

优点：改动小。缺点：不能解决信息主线割裂，P57 的产品目标达成不足。

### Option C: Merge Dashboard and Workbench into one route

让 `/` 和 `/workbench` 渲染完全相同页面，最大程度减少重复。

优点：一致性最强。缺点：可能破坏既有用户路径和 P39/P55 route matrix 的语义，也会让 Workbench 的任务展开空间不足。

P58 采用 Option A：共享模型和核心组件，但保留两个路由的差异化角色。

## Information Architecture

### Dashboard `/`

Dashboard 是每日入口，首屏顺序为：

1. 今日状态摘要：最终裁决、状态标签、更新时间、数据可信度。
2. 禁止动作和可选人工动作：列表必须清晰，不得暗示自动执行。
3. 下一步人工动作：查看决策详情、查看纪律报告、维护持仓、处理风险、查看数据质量。
4. 数据与风险摘要：证据数量、独立信源、组合风险、活跃风险。
5. 解释入口：规则、证据、审计、决策闭环。

### Workbench `/workbench`

Workbench 是任务视图，首屏顺序为：

1. 与 Dashboard 同源的今日状态摘要。
2. 人工动作队列，按优先级组织：
   - 阻断前提：账户/持仓缺失、数据不足、source degraded。
   - 今日复核：查看纪律报告、查看决策详情。
   - 风险处置：查看风险预警。
   - 解释追踪：查看证据、规则、复盘、审计。
3. 组合与风险、规则与复盘、主动咨询入口作为次级区域。

## View Model

新增或提取一个轻量前端 view model，输入仍来自现有 API DTO：

- `DashboardTodayResponse`
- `DailyDisciplineReport`
- `PortfolioCurrentResponse`
- `RiskAlert[]`
- `RuleProposal[]`
- `ReviewSummary`

核心字段：

- `statusLabel`
- `statusTone`
- `verdictText`
- `trustSummary`
- `riskSummary`
- `updatedAtText`
- `prohibitedActions`
- `optionalActions`
- `nextActions[]`
- `explanationLinks[]`
- `warnings[]`

View model 必须 null-safe。缺失字段展示安全空态，不把 unknown/degraded/insufficient 显示为成功。

## Components

P58 允许新增小范围前端组件，但不做 P62 级组件库抽象：

- `DailyDecisionHero`：今日状态、裁决、可信度、更新时间。
- `ManualActionQueue`：下一步人工动作列表，只渲染本地导航链接。
- `WorkbenchSignalGrid`：数据、风险、规则、复盘摘要。
- `SafetyBoundaryNote`：只读/人工复核边界，文案简短。

组件应放在 `web/src/components/dashboard/` 或 `web/src/features/dashboard/`，视复用范围决定。不得在组件中直接调用 API。

## Visual Direction

- 使用现有 `global.css` 的 operational tokens，不引入新主题。
- 第一屏是工作台摘要，不做 hero marketing。
- 状态色按语义区分：normal/success、warning/degraded、danger/high risk、unknown/insufficient。
- 移动端 390px 下摘要纵向堆叠，动作队列保持可点击和不横向溢出。
- 避免卡片套卡片。摘要 band 和动作队列可以是 page section 或单层 card。

## Safety Boundary

P58 所有 CTA 都只能是查看、维护、记录、复核、进入页面等本地人工动作或只读导航。禁止出现自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、覆盖真实库、收益承诺等入口或暗示。

## Validation Plan

- 方案复审：子 agent 必须检查 Product Design brief、P57 alignment、信息架构、验收门禁和安全边界。
- 单元/组件测试：
  - Dashboard 成功态、缺数据态、降级态、高风险态。
  - Workbench 人工动作队列、数据可信度、风险摘要和安全文案。
  - 禁止自动交易/自动确认/自动规则应用文案。
- 构建和回归：
  - `npm test`
  - `npm run build`
  - `go test ./...`
- 浏览器验收：
  - 本地 server + Vite。
  - `/` 和 `/workbench` 桌面 1280x720 截图。
  - `/` 和 `/workbench` 移动 390x844 截图。
  - 390px 检查 `body.scrollWidth <= viewport`。
  - 页面扫描 forbidden copy。
