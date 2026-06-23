# Design: P59 决策解释体验重构

## Product Design Brief

P59 的设计目标是让用户在主动咨询和决策详情中快速理解一条完整链路：

1. 我问了什么。
2. 系统生成了什么裁决。
3. 为什么是这个裁决。
4. 证据、LLM、规则和审计分别贡献了什么。
5. 哪些动作被禁止，哪些只能人工复核。
6. 后续闭环是否已确认、记录、复盘或仍有缺口。

本阶段使用 Product Design get-context、P55-P56 UI audit、P57 产品体验路线图和 P58 operational cockpit 风格作为设计输入。页面继续采用安静、密集、可扫描的本地投资纪律产品风格，不做营销页、AI 聊天首页或装饰性 hero。

## Current Problems

- `/consultation` 与 `/decisions/:decisionId` 复用同一页面，但页面标题和信息层级偏“决策详情”，没有把主动咨询输入、结果和查看解释路径明确分开。
- 决策详情现有 `DecisionTrace` 内容完整，但是长技术 trace 堆叠；最终裁决、禁止动作、数据可信度和“为什么”没有形成第一屏主线。
- Evidence、LLM、rules、audit、expected return 与 confirmation 各自存在，但缺少面向用户的贡献分层和默认折叠策略。
- `/evidence` 更像证据清单，缺少“这些证据如何支撑决策解释”的入口说明。
- `/decision-loop` 已能串联闭环事实，但需要更清晰地区分已完成阶段、缺口和只读安全边界。
- P56 修复了真实 LLM nullable blocker，但 P59 需要继续覆盖更广泛 missing/nullable DTO 和真实 LLM UI 路径。

## Approach Options

### Option A: Shared decision explanation model and story components

建立轻量前端 explanation view model，将现有 decision DTO 映射为“输入/裁决/原因/证据/LLM/rules/audit/闭环”的展示结构。`DecisionDetailPage`、`DecisionTrace`、`EvidencePage`、`DecisionLoopPage` 分别复用其中的状态、文案和安全边界。

优点：解释主线一致、null-safe 逻辑集中、测试可覆盖。缺点：需要抽取一层模型并调整多个页面。

### Option B: Only restyle existing cards

保留现有 card stack，只通过标题、CSS 和局部文案强化重点。

优点：改动小。缺点：不能解决主动咨询、详情、证据和闭环之间的故事断裂。

### Option C: Merge consultation/detail/evidence/loop into one wizard

把四个页面整合成一个连续流程。

优点：链路感强。缺点：破坏既有路由矩阵和本地工具型产品的可达性，也会扩大 P59 范围。

P59 采用 Option A：共享解释模型和局部组件，保留既有路由与 API，不新增后端能力。

## Information Architecture

### Consultation `/consultation`

首屏顺序：

1. 主动咨询上下文：标的、场景、问题、输入假设。
2. 生成按钮与安全说明：只生成分析和决策材料，不自动确认、不交易、不应用规则。
3. 成功结果摘要：裁决、关键原因、禁止动作、可选人工动作、数据可信度。
4. 查看解释路径：生成决策详情、证据、决策闭环、本地审计的导航。
5. 失败/降级状态：展示错误分类、可重试建议和不展示交易建议的安全空态。

### Decision Detail `/decisions/:decisionId`

首屏顺序：

1. 决策故事 hero：最终裁决、状态、生成时间、标的、场景、数据可信度。
2. 安全边界：禁止动作、可选人工动作、只读/人工复核说明。
3. 为什么：最多 3-5 条关键原因，来自最终裁决、规则、LLM 和 evidence summary。
4. 可信度：独立信源、最高信源等级、检索质量、LLM 质量/解析状态。
5. 下一步本地导航：证据、审计、决策闭环、确认记录入口。
6. 二级详情：Evidence、LLM、rules、expected return、arbitration、audit timeline 和 confirmation 默认分层，长 trace 折叠。

### Evidence `/evidence`

首屏顺序：

1. 证据可信度概览：已验证、需关注、背景材料、最新时间。
2. 来源等级说明：S/A/B/C、formal/background 的实际含义。
3. 与决策解释相关的导航：返回决策详情、查看决策闭环、查看审计。
4. 证据表格：保留筛选和展开详情，但避免把原始 hash/vendor payload 当作主要信息。

### Decision Loop `/decision-loop`

首屏顺序：

1. 闭环概览：总数、未闭合数量、最近决策、只读说明。
2. 阶段时间线：建议、确认、线下记录、风险/复盘、审计。
3. 缺口：缺少确认、本地流水、风险线索、复盘线索时直接说明下一步人工动作。
4. 本地链接：只导航到风险、复盘、审计或决策详情，不提供生命周期写入动作。

## View Model

P59 允许新增轻量前端 view model，例如 `DecisionExplanationViewModel`。输入仍来自现有 DTO，不新增 API 字段：

- `DecisionDetailResponse`
- `FinalVerdict`
- `AnalystReport[]`
- `RetrievalQualitySummary`
- `EvidenceItem[]`
- `DecisionLoopItem[]`

核心字段：

- `storyTitle`
- `decisionContext`
- `verdictText`
- `statusTone`
- `trustSummary`
- `keyReasons[]`
- `prohibitedActions[]`
- `optionalActions[]`
- `explanationLinks[]`
- `evidenceHighlights[]`
- `analystHighlights[]`
- `ruleHighlights[]`
- `auditHighlights[]`
- `safetyNotes[]`
- `missingDataWarnings[]`

所有 list 字段必须 null-safe。缺失字段展示“暂无/未提供/需人工复核”等安全空态，不把 unknown、degraded、insufficient、missing 或 failed 显示成普通成功。

## Components

P59 允许新增小范围 presentational components，但不做 P62 级组件库抽象：

- `DecisionStoryHero`：裁决、上下文、数据可信度、生成时间。
- `DecisionSafetyPanel`：禁止动作、可选人工动作、人工复核边界。
- `DecisionWhyPanel`：关键原因和贡献来源。
- `DecisionTrustPanel`：证据、检索质量、LLM 质量和规则链可信度。
- `DecisionTraceSections` 或等价分层结构：Evidence、LLM、rules、expected return、audit 等二级详情。
- Evidence/Loop 页面可复用小型 summary/timeline 组件，但组件中不得直接调用 API。

## Visual Direction

- 延续 P58 的 operational cockpit 视觉语言：密集、清晰、低装饰、状态优先。
- 第一屏是决策解释，不是营销 hero，也不是聊天气泡。
- 状态色按语义区分：safe/neutral、warning/degraded、danger/prohibited、unknown/insufficient。
- 长文本和技术 trace 默认折叠；展开内容必须可读，不造成页面级横向滚动。
- 390px 下按照“裁决 -> 安全边界 -> 为什么 -> 导航 -> 详情”纵向堆叠。
- 避免卡片套卡片；页面区块使用单层 cards 或 full-width section。

## Safety Boundary

P59 所有 CTA 都只能是生成分析、查看、复核、记录提示或本地导航。页面不得提供或暗示券商连接、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、覆盖真实库或收益承诺。LLM 贡献只作为分析材料展示，不得被描述为最终裁决来源。

## Validation Plan

- 方案复审：子 agent 必须检查 Product Design brief、P57/P58 alignment、信息架构、真实 LLM 验收、nullable DTO 和安全边界。
- 单元/组件测试：
  - `/consultation` 成功、失败、真实 LLM-like DTO、解释链接和安全文案。
  - `/decisions/:decisionId` 首屏裁决、禁止动作、可选人工动作、可信度、关键原因、长 trace 折叠和 missing/nullable DTO。
  - `/evidence` 可信度概览、筛选/展开、链接可达和空态/错误态。
  - `/decision-loop` 阶段时间线、缺口、本地链接、安全只读边界和无写入动作按钮。
- 构建和回归：
  - `npm test`
  - `npm run build`
  - `go test ./...`
- 浏览器验收：
  - 本地 server + Vite。
  - 真实 LLM consultation 生成决策，并从 UI 打开决策详情。
  - `/consultation`、`/decisions/:decisionId`、`/evidence`、`/decision-loop` 桌面截图。
  - 390px 移动截图与 `body.scrollWidth <= viewport` 检查。
  - UI 文案 forbidden copy 扫描。
