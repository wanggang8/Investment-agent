# Design: P57 产品体验打磨总规划

## Skill And Research Inputs

P57 使用 Product Design plugin 的 `get-context` 方法确认设计 brief，并沿用 P55/P56 的 audit 和 research 证据。Superpowers `brainstorming` / `writing-plans` 的方法用于组织方案，但根据 `docs/GOVERNANCE.md`，规格真源必须写入 OpenSpec change 和 `docs/`，不得新建 `docs/superpowers/plans/` 作为真源。

调研依据继续采用 P56 已记录的稳定 UX 资料：

- Nielsen Norman Group dashboard guidance: `https://www.nngroup.com/articles/dashboards-preattentive/`
- Nielsen Norman Group form design: `https://www.nngroup.com/articles/web-form-design/`
- Nielsen Norman Group cognitive load / form principles: `https://www.nngroup.com/articles/4-principles-reduce-cognitive-load/`
- WCAG 2.1 Reflow: `https://www.w3.org/WAI/WCAG21/Understanding/reflow.html`
- Material Design data tables: `https://m2.material.io/components/data-tables`
- Explainable robo-advisor research: `https://ceur-ws.org/Vol-3222/paper6.pdf`
- Robo-advisor trust / UX reference: `https://www.zymr.com/blog/how-to-build-robo-advisors-platform`

## Product Brief

Investment Agent 的产品体验应围绕“本地投资纪律工作台”组织。它不是券商交易终端，不替用户下单；不是 AI 聊天工具，不让 LLM 决定最终裁决；也不是工程调试后台，不要求用户理解内部 DTO、trace 或 collector 细节才能完成日常工作。

核心用户每天打开产品时，需要快速回答四个问题：

1. 今天能不能动？
2. 为什么是这个结论？
3. 我需要人工做什么？
4. 数据和规则是否可信？

## Experience Principles

### 1. Decision First

首页和工作台应优先呈现今日结论、禁止动作、允许的人工动作、风险等级和数据可信度。证据、规则、LLM 分析和审计材料应作为解释层逐级展开，而不是压在第一屏。

### 2. Manual Action Clarity

所有动作必须是人工动作、线下记录或只读导航。按钮文案应避免暗示自动交易、自动确认、自动修复或自动应用规则。高风险、未知、降级和信息不足状态不得被视觉上处理成成功。

### 3. Explainability Over Trace Dump

决策详情应讲清楚“裁决是什么、为什么、证据是否可靠、LLM 贡献了什么、规则如何约束、审计链路在哪里”。原始 trace、长 JSON、诊断细节应默认折叠或转换为可扫描结构。

### 4. Operational Density

这是高频使用的本地工作台，不做营销式 landing page。视觉应保持密集、稳定、可扫描；卡片用于承载单个重复项或明确工具，不把页面区块全部做成装饰卡片。

### 5. Reflow And Accessibility

核心路径必须在 390px 移动视口、768px 窄桌面/平板和 1280px 桌面可用。页面级横向滚动只允许出现在明确二维数据容器内。核心表单、导航、折叠区和按钮必须有可访问名称和键盘路径。

## Product Architecture

### Primary Task Areas

| Area | User question | Current routes | Future treatment |
| --- | --- | --- | --- |
| 今日 | 今天能不能动？ | `/`, `/workbench`, daily reports | P58 合并为每日决策中心 |
| 决策 | 为什么这么建议？ | `/consultation`, `/decisions/:id`, `/decision-loop`, `/evidence` | P59 重构成解释链路 |
| 组合与风险 | 我的本地事实和风险是否可靠？ | `/positions`, `/risk-alerts`, `/data-quality` | P60 重构为维护与处置队列 |
| 治理 | 规则、审计和复盘是否可信？ | `/rules`, `/audit`, `/review`, `/notifications` | P61 降低工具感 |
| 系统 | 本地运行是否健康？ | `/settings`, `/local-install`, `/local-knowledge`, `/daily-auto-run` | P61/P62 统一运维体验 |

### Page Priority

一级核心页面：

- Dashboard / Workbench
- Consultation
- Decision Detail
- Positions
- Data Quality
- Risk Alerts

二级解释和复盘页面：

- Evidence
- Decision Loop
- Daily Reports
- Review

三级治理和运维页面：

- Rules
- Audit
- Notifications
- Local Install
- Local Knowledge
- Settings
- Daily Auto Run

## Stage Roadmap

### P58: 今日工作台重构

目标：把 Dashboard 和 Workbench 打磨成真正的每日投资纪律 cockpit。

范围：

- 首屏显示今日裁决、风险、数据可信度、最近决策和下一步人工动作。
- 收敛重复模块，避免 dashboard/workbench 信息割裂。
- 提供清晰入口：查看解释、维护持仓、处理风险、查看数据质量。
- 移动端第一屏可读，不能依赖横向滚动。

验收：

- Vitest 覆盖 dashboard/workbench 空态、降级态、高风险态和安全文案。
- Playwright 覆盖 `/`、`/workbench` 桌面和 390px 移动端。
- 截图证明 5 秒内可识别今日状态和下一步人工动作。

### P59: 决策解释体验重构

目标：把 Consultation、Decision Detail、Evidence 和 Decision Loop 串成可理解的决策故事。

范围：

- Consultation 明确输入假设、生成建议、查看解释的路径。
- Decision Detail 首屏展示最终裁决、禁止动作、可选人工动作、数据可信度和安全边界。
- Evidence / LLM / rules / audit 分层展示，长 trace 默认折叠。
- Decision Loop 只读串联建议、确认、线下记录、风险、复盘和审计。

验收：

- 真实 LLM consultation 后可打开新决策详情。
- Nullable/missing DTO fixture 回归。
- Evidence 链接、audit 链接、decision loop 链接可达。
- 禁止自动交易/自动确认/自动规则应用文案扫描。

### P60: 组合、风险与数据质量体验重构

目标：把 Positions、Risk Alerts 和 Data Quality 打磨成日常维护和处置页面。

范围：

- Positions 区分初始化、编辑、校准、导入、错误修正。
- Risk Alerts 从列表转为处置队列：待看、处理中、已记录、需复盘。
- Data Quality 把 source health、RAG、LLM、affected workflows 做成可扫描质量面板。
- 长 source、diagnostic、decision id 统一换行、截断或局部滚动。

验收：

- 表单组件测试覆盖校验、错误、成功和空态。
- 390px 移动端无页面级横向溢出。
- degraded/current/stale/source_unavailable fixture 覆盖。
- 风险处置不包含自动交易或外部推送入口。

### P61: 治理和运维页面产品化

目标：降低 Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings 的工程工具感。

范围：

- Rules 以规则提案卡片展示原因、样本、过拟合风险、守门人状态和人工确认边界。
- Audit 从 raw event list 优化为时间线和筛选摘要。
- Notifications 形成只读通知收件箱和处理状态概览，明确本地应用内通知边界，不暗示外部推送。
- Daily Reports 形成纪律复盘报告，而不是数据 dump。
- Daily Auto Run 展示本地自动运行配置、最近执行、失败诊断和手动复验入口，保持默认关闭或显式启用边界，不承诺自动修复。
- Local Install / Settings / Local Knowledge 统一诊断、配置、脱敏和安全提示模式。

验收：

- 页面级组件测试覆盖关键状态。
- 浏览器截图审查 dense JSON/pre 是否已被产品化或合理折叠。
- Notifications 必须扫描无外部推送、短信、邮件或第三方通知承诺。
- Daily Auto Run 必须扫描无自动修复、自动确认、自动规则应用、覆盖真实库或后台交易承诺。

### P62: 设计系统与可访问性验收

目标：把重复 UI 规则固化为稳定基础层和可验收标准。

范围：

- 定义并复用 Button、Field、StatusBadge、PageHeader、SummaryCard、DetailSection、ResponsiveTable、EmptyState、ErrorState。
- 统一 status token：success、warning、danger、degraded、unknown、readonly、blocked。
- 键盘路径：主导航、移动菜单、表单、折叠区、关键按钮。
- 可访问语义：landmark、label、aria-expanded、错误提示、状态文本。
- 视觉回归：固定关键页面截图。

验收：

- Vitest 覆盖组件行为。
- Playwright keyboard smoke。
- 390px、768px、1280px 截图。
- WCAG reflow 检查。

### P63: 全量真实 UI 回归与发布状态刷新

目标：产品体验打磨完成后，重新执行全量真实 UI 验收并刷新最终交付口径。

范围：

- 真实启动后端和前端。
- 真实 LLM consultation。
- 全路由浏览器操作。
- 移动端、桌面端、错误态、降级态、安全边界扫描。
- 更新 release-ready / handoff 文档。

验收：

- `npm test`
- `npm run build`
- `go test ./...`
- Playwright 全路由和真实 LLM journey。
- OpenSpec 全量 strict。
- 敏感信息扫描。
- 子 agent 执行后复审和提交前复审。

## Safety Boundary

P57-P63 均不得新增或暗示：

- 券商接口。
- 自动交易。
- 一键交易。
- 代下单。
- 外部推送。
- 自动确认。
- 自动规则应用。
- 自动修复承诺。
- 自动覆盖真实库。
- 收益承诺。
- 登录源、付费源、授权源、Level2 或高频源。

## Review Gates

每个后续阶段必须遵守：

1. 创建独立 OpenSpec change。
2. 方案写完后子 agent 审查。
3. 审查无 Critical / Important 后执行。
4. 执行完成后真实 UI / 自动化验证。
5. 执行后子 agent 审查。
6. 审查无 Critical / Important 后 archive。
7. archive 后提交前子 agent 审查。
8. 审查无 Critical / Important 后提交。
