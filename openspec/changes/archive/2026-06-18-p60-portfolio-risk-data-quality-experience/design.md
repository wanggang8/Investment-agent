# Design: P60 组合、风险与数据质量体验重构

## Product Design Brief

P60 的设计目标是把组合维护、风险处置和数据质量从“能用的功能页”打磨成日常操作页：

1. 用户先知道当前组合是否可用于今日纪律和主动咨询。
2. 用户能区分初始化、校准、编辑、导入、纠错这些本地事实动作。
3. 用户打开风险页后能按处置队列理解待看、处理中、需复盘、已记录的风险。
4. 用户打开数据质量页后能看到 source health、证据/RAG、LLM 和受影响工作流之间的关系。
5. 所有动作都保持本地事实记录、人工复核、只读诊断或本地导航，不出现自动交易、外部推送或自动修复暗示。

本阶段使用 Product Design get-context、P55-P56 UI audit、P57 产品体验路线图、P58 daily workbench 和 P59 decision story 的 operational cockpit 风格作为设计输入。页面继续采用安静、密集、可扫描的本地投资纪律产品风格，不做营销页、券商账户页、自动交易工作台或装饰性 hero。

Product Design saved-context preflight 已执行，当前未发现已保存的 Product Design context；本阶段以仓库内现有页面、P57 路线图和 P58/P59 已实现样式作为视觉来源。

## Current Problems

- `/positions` 把首次使用、校准、维护、线下交易、导入、纠错全部堆在同一长页里；用户难以判断当前应该先做哪一步。
- `/positions` 只有账户快照和表格，缺少“组合是否可用”“本地事实是否需要校准”“下一步人工动作”的首屏判断。
- `/risk-alerts` 已能列出风险和 SOP 按钮，但更像记录列表；风险没有形成待处理队列，也缺少处置状态说明和复盘线索。
- `/data-quality` 已聚合 source health、证据、LLM、复盘降级，但首屏仍偏技术分块；质量状态、影响范围和下一步本地导航需要更强关系表达。
- 三页在 390px 下已基本可用，但 P60 需要继续检查长表格、长风险摘要、source health 明细和按钮排列。
- 这三页承担日常维护职责，安全边界必须比普通只读页更明确，尤其要避免“修复”“确认”“处理”被误解为自动执行。

## Approach Options

### Option A: Shared page-level experience models

为组合、风险、数据质量分别新增轻量 view model，将现有 DTO 映射为页面状态、队列、下一步人工动作、质量信号和安全文案。页面只负责加载、提交和布局。

优点：信息架构清晰、测试可以覆盖模型逻辑、不会扩大后端范围。缺点：需要新增少量前端模型文件并调整现有页面。

### Option B: Only restyle existing pages

保留现有页面结构，仅调整标题、文案和 CSS。

优点：改动小。缺点：不能解决动作混杂、风险队列缺失和数据质量关系不清的问题。

### Option C: Merge the three pages into one portfolio operations hub

把组合、风险、数据质量合并为单页。

优点：链路集中。缺点：破坏既有导航和验收路径，也会扩大 P60 范围并影响 P61/P63。

P60 采用 Option A：保留三条既有路由，新增轻量展示模型和局部组件，把每页首屏重构为“状态 -> 下一步 -> 明细”的产品化结构。

## Information Architecture

### Portfolio `/positions`

首屏顺序：

1. 组合维护状态：总资产、现金、持仓数、高风险比例、快照时间。
2. 当前阶段：首次初始化、需要校准、正常维护、导入待确认、错误修正提示。
3. 下一步人工动作：维护账户事实、编辑持仓、补记线下交易、校验导入、记录修正。
4. 本地事实表单：按模式分区展示初始化/校准、持仓维护、线下交易、批量导入、错误修正。
5. 持仓结构和当前持仓：保留图表和表格，但放在状态与动作之后。

### Risk Alerts `/risk-alerts`

首屏顺序：

1. 风险处置总览：活跃/升级/观察/已解除/已归档数量，最严重风险和受影响标的。
2. 处置队列：待看、处理中、需复盘、已记录。
3. 风险卡片：风险类型、严重程度、SOP 状态、触发摘要、禁止动作、建议人工动作、更新时间。
4. 关联材料：决策、纪律报告、通知、审计。
5. 生命周期动作：继续观察、升级复核、解除预警。按钮文案必须强调“本地 SOP 记录”，不暗示交易或外推。

队列映射：

- 待看：`triggered`
- 处理中：`active`、`observing`
- 需复盘：`escalated`
- 已记录：`resolved`、`archived`

### Data Quality `/data-quality`

首屏顺序：

1. 数据质量总览：source health、证据/RAG、LLM、影响工作流四类信号和整体状态。
2. 下一步本地检查：查看设置、查看证据、查看风险预警、查看受影响决策、查看复盘。
3. source health 明细：来源、类别、新鲜度、失败类别、影响标的、最近成功/失败。
4. 证据/RAG 明细：证据数、核验状态、独立信源、VecLite 状态。
5. LLM 与复盘质量：DeepSeek 状态、复盘/数据源/索引状态、降级和缺证据数量。
6. 受影响工作流：决策 ID、标的、降级状态和本地决策链接。

## View Models

P60 允许新增轻量前端 view model。输入仍来自现有 DTO，不新增 API 字段。

### PortfolioExperienceModel

输入：

- `PortfolioCurrentResponse`
- 当前表单状态和导入校验状态
- `PageErrorState`

核心字段：

- `statusLabel`
- `statusTone`
- `summaryMetrics[]`
- `stageLabel`
- `stageDescription`
- `nextActions[]`
- `maintenanceModes[]`
- `safetyNotes[]`
- `warnings[]`

### RiskDispositionModel

输入：

- `RiskAlert[]`
- `PageErrorState`

核心字段：

- `summaryMetrics[]`
- `highestSeverity`
- `affectedSymbols[]`
- `queues[]`
- `nextActions[]`
- `safetyNotes[]`
- `warnings[]`

### DataQualityExperienceModel

输入：

- `SystemStatus`
- `MarketSnapshot`
- `SourceHealthItem[]`
- `EvidenceItem[]`
- `SourceVerification`
- `ReviewSummary`
- 各 API 的 `PageErrorState`

核心字段：

- `overallLabel`
- `overallTone`
- `qualitySignals[]`
- `nextActions[]`
- `sourceHealthRows[]`
- `trustMetrics[]`
- `impactRows[]`
- `safetyNotes[]`
- `warnings[]`

所有 list 字段必须 null-safe。缺失、unknown、degraded、stale、parse_error、missing、failed 或 insufficient 不得展示为普通成功。

## Components

P60 允许新增小范围 presentational components，但不做 P62 级组件库抽象：

- `PortfolioMaintenanceHero` 或等价 hero：组合状态、阶段、下一步人工动作。
- `PortfolioMaintenanceModes` 或等价模式区：初始化、校准、编辑、线下交易、导入、纠错。
- `RiskDispositionSummary` 或等价 summary：风险队列和严重程度。
- `RiskQueueSection` 或等价队列组件：按状态分组展示风险卡片。
- `DataQualitySummary` 或等价 summary：质量信号和下一步本地导航。
- 可复用 P58/P59 的 `.daily-hero`、`.daily-action-queue`、`.daily-signal-grid` 风格，但命名可更通用。

组件不得直接调用 API、SQLite、VecLite、localStorage、sessionStorage 或本地文件。

## Visual Direction

- 延续 P58/P59 的 operational cockpit：状态优先、信息密集、低装饰、可扫描。
- 第一屏必须是维护/处置/质量状态，不是表单海洋，也不是技术日志。
- 状态色按语义区分：success/normal、warning/degraded、danger/escalated、unknown/insufficient。
- 表单与写入动作必须与状态/解释分区，不让用户误以为页面会自动完成交易或修复。
- 390px 下按照“状态 -> 下一步 -> 表单/队列/明细 -> 表格”的纵向顺序堆叠。
- 避免卡片套卡片；页面区块使用单层 cards 或 full-width section。

## Safety Boundary

P60 所有 CTA 只能是本地事实记录、人工复核、查看、校验、继续观察、升级复核、解除预警或本地导航。页面不得提供或暗示券商连接、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、覆盖真实库或收益承诺。

`/data-quality` 中的“下一步”必须表达为检查或查看，不得表达为一键修复。`/risk-alerts` 中的 SOP 动作只更新本地风险状态，不代表交易、通知外推或规则变更。`/positions` 中的写入动作只记录本地账户事实或审计事实。

## Validation Plan

- 方案复审：子 agent 必须检查 Product Design brief、P57-P59 alignment、信息架构、现有 API 边界、真实 UI 验收、移动端和安全边界。
- 单元/组件测试：
  - Portfolio model 和 `/positions`：初始化、正常维护、高风险比例、导入待确认、错误态、所有写入动作安全文案。
  - Risk model 和 `/risk-alerts`：队列分组、严重程度、关联链接、生命周期按钮、resolved/archived 无动作。
  - Data Quality model 和 `/data-quality`：整体状态、source health、证据/RAG、LLM、受影响工作流、安全脱敏和无执行按钮。
- 构建和回归：
  - `npm test`
  - `npm run build`
  - `go test ./...`
- 浏览器验收：
  - 本地 server + Vite。
  - `/positions` 执行一次本地校准或维护路径，并确认表格/状态更新。
  - `/risk-alerts` 查看队列和关联链接；如 fixture 支持，执行一次本地 SOP 生命周期动作。
  - `/data-quality` 查看质量总览、受影响工作流和本地导航。
  - 桌面与 390px 移动截图。
  - `body.scrollWidth <= viewport` 和 forbidden copy scan。
