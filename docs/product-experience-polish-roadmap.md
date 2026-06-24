# 产品体验打磨路线图

> 日期：2026-06-17
> 来源：P57 `p57-product-experience-polish-roadmap`
> 状态：P57 规划已归档；后续改造必须逐阶段创建独立 OpenSpec change。

## 1. 产品定位

Investment Agent 是本地投资纪律工作台，不是券商交易终端、AI 聊天 demo、营销落地页或工程调试后台。

用户每天打开产品时，核心问题是：

1. 今天能不能动？
2. 为什么是这个结论？
3. 我需要人工做什么？
4. 数据和规则是否可信？

P56 已修复真实 LLM 决策详情崩溃、移动端溢出和基础产品化问题。P57 不声明所有产品设计、UI 设计或功能设计问题已修复，而是固化后续 P58-P63 的打磨路线图和验收门禁。

## 2. 设计原则

### 2.1 决策优先

首页和工作台优先展示今日结论、禁止动作、允许的人工动作、风险等级和数据可信度。证据、规则、LLM 分析和审计材料作为解释层逐级展开。

### 2.2 人工动作清晰

所有动作必须是人工动作、线下记录或只读导航。按钮、状态和提示不得暗示自动交易、自动确认、自动修复、自动规则应用或收益承诺。

### 2.3 解释优先于 trace dump

决策详情要讲清楚裁决、原因、证据可靠性、LLM 贡献、规则约束和审计链路。长 JSON、原始 trace 和诊断细节应默认折叠或转换成可扫描结构。

### 2.4 高密度但可扫描

产品应保持安静、密集、可信、可扫描；不做营销式 hero。卡片只用于承载单个重复项、摘要或明确工具，不把所有页面区块都做成装饰卡片。

### 2.5 Reflow 和可访问性

核心路径必须在 390px 移动视口、768px 窄桌面/平板和 1280px 桌面可用。页面级横向滚动只允许出现在明确二维数据容器内。主导航、表单、折叠区和关键按钮必须有可访问名称和键盘路径。

## 3. 页面分层

| 层级 | 页面 | 后续阶段 |
| --- | --- | --- |
| 一级核心页面 | Dashboard、Workbench、Consultation、Decision Detail、Positions、Data Quality、Risk Alerts | P58-P60 |
| 二级解释和复盘页面 | Evidence、Decision Loop、Daily Reports、Review | P59、P61 |
| 三级治理和运维页面 | Rules、Audit、Notifications、Local Install、Local Knowledge、Settings、Daily Auto Run | P61-P62 |

## 4. 后续阶段

### P58 今日工作台重构

目标：把 Dashboard 和 Workbench 打磨成真正的每日投资纪律 cockpit。

状态：P58 已按 `p58-daily-workbench-redesign` 实现并归档，归档路径为 `openspec/changes/archive/2026-06-17-p58-daily-workbench-redesign/`。

范围：

- 首屏显示今日裁决、风险、数据可信度、最近决策和下一步人工动作。
- 收敛重复模块，避免 dashboard/workbench 信息割裂。
- 提供清晰入口：查看解释、维护持仓、处理风险、查看数据质量。
- 移动端第一屏可读，不依赖横向滚动。

验收：

- Vitest 覆盖 dashboard/workbench 空态、降级态、高风险态和安全文案。
- Playwright 覆盖 `/`、`/workbench` 桌面和 390px 移动端。
- 截图证明 5 秒内可识别今日状态和下一步人工动作。

验收记录：`docs/release/acceptance/2026-06-17-p58-ui-acceptance.md`。

### P59 决策解释体验重构

目标：把 Consultation、Decision Detail、Evidence 和 Decision Loop 串成可理解的决策故事。

状态：P59 已按 `p59-decision-explainability-experience` 实现并归档，归档路径为 `openspec/changes/archive/2026-06-17-p59-decision-explainability-experience/`。

范围：

- Consultation 明确输入假设、生成建议、查看解释路径。
- Decision Detail 首屏展示最终裁决、禁止动作、可选人工动作、数据可信度和安全边界。
- Evidence、LLM、rules、audit 分层展示，长 trace 默认折叠。
- Decision Loop 只读串联建议、确认、线下记录、风险、复盘和审计。

验收：

- 真实 LLM consultation 后可打开新决策详情。
- Nullable/missing DTO fixture 回归。
- Evidence 链接、audit 链接、decision loop 链接可达。
- 禁止自动交易、自动确认、自动规则应用文案扫描。

验收记录：`docs/release/acceptance/2026-06-18-p59-ui-acceptance.md`。

备注：2026-06-18 真实 LLM consultation 已通过前端 UI 触发 Analyst 节点；远端 `/v1/chat/completions` 返回 HTTP 503，UI 按安全降级展示 `LLM 材料 0 份`，并保留规则裁决、证据、审计和闭环解释。该结果作为外部依赖失败记录，不伪装为 LLM 成功返回。

### P60 组合、风险与数据质量体验重构

目标：把 Positions、Risk Alerts 和 Data Quality 打磨成日常维护和处置页面。

状态：P60 已按 `p60-portfolio-risk-data-quality-experience` 实现并归档到 `openspec/changes/archive/2026-06-18-p60-portfolio-risk-data-quality-experience/`。

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

验收记录：`docs/release/acceptance/2026-06-18-p60-ui-acceptance.md`。

备注：2026-06-18 真实本地 UI 验收已覆盖 `/positions` 本地校准、`/risk-alerts` 风险队列空态、`/data-quality` 本地设置导航、桌面/390px 截图、横向溢出检查和 forbidden copy scan；P60 结论不代表 P63 最终 release-ready 刷新。

### P61 治理和运维页面产品化

目标：降低 Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings 的工程工具感。

状态：P61 已按 `p61-governance-ops-productization` 实现并归档到 `openspec/changes/archive/2026-06-18-p61-governance-ops-productization/`。

范围：

- Rules 以规则提案卡片展示原因、样本、过拟合风险、守门人状态和人工确认边界。
- Audit 从 raw event list 优化为时间线和筛选摘要。
- Notifications 形成只读通知收件箱和处理状态概览，明确本地应用内通知边界，不暗示外部推送。
- Daily Reports 形成纪律复盘报告，而不是数据 dump。
- Daily Auto Run 展示本地自动运行配置、最近执行、失败诊断和手动复验入口，保持默认关闭或显式启用边界，不承诺自动修复。
- Local Install、Settings、Local Knowledge 统一诊断、配置、脱敏和安全提示模式。

验收：

- 页面级组件测试覆盖关键状态，`npm test` 通过。
- `bash scripts/e2e-smoke.sh` 已覆盖 P61 路由、关键状态、下一步动作、390px reflow 和 forbidden copy scan。
- 真实本地后端 + Vite 前端浏览器验收覆盖八个 P61 页面；Local Install 上传诊断摘要 fixture，Local Knowledge 执行 validate/安全样本 confirm，Settings 执行真实市场刷新。
- 截图与结果见 `docs/release/acceptance/2026-06-18-p61-ui-acceptance.md` 与 `docs/release/ui-audit-assets/2026-06-18-p61/browser-results.json`。

### P62 设计系统与可访问性验收

目标：把重复 UI 规则固化为稳定基础层和可验收标准。

状态：P62 已按 `p62-design-system-accessibility-hardening` 实现并归档到 `openspec/changes/archive/2026-06-18-p62-design-system-accessibility-hardening/`。

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

当前验收记录：`docs/release/acceptance/2026-06-18-p62-ui-acceptance.md`。P62 结果只覆盖设计系统与可访问性加固，不代表 P63 全量真实 UI 回归或最终 release-ready 刷新。

### P63 全量真实 UI 回归与发布状态刷新

目标：产品体验打磨完成后，重新执行全量真实 UI 验收并刷新最终交付口径。

状态：P63 已按 `p63-full-ui-regression-release-refresh` 实现并归档到 `openspec/changes/archive/2026-06-18-p63-full-ui-regression-release-refresh/`；已完成 G0-G9、全路由真实 UI 回归、真实 LLM-backed consultation UI journey、截图和发布材料刷新。

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

验收记录：`docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`；发布候选：`docs/release/release-candidate-2026-06-18.md`；交付说明：`docs/release/release-handoff-2026-06-18.md`。

### P111 高保真参考图视觉重构

目标：在 P110 后把用户确认的第二套参考图落到整个产品，而不是只做首页或 moodboard。P111 参考图真源为 `/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`。

范围：

- 全局 shell：深色左侧导航、白色顶部状态栏、本地模式与数据截至控件。
- 核心页面：Dashboard、Workbench、Positions、Data Quality、Risk Alerts、Consultation、Decision Detail、Decision Loop、Evidence。
- 治理和运维页面：Rules、Review、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings。
- 共享 reference components：report hero、priority action queue、status metric grid、snapshot strip、progress tracker、evidence checklist、ledger surface。

验收：

- 每个覆盖路由都必须采集桌面与移动截图。
- 每个页面必须填写视觉 mismatch ledger；P0/P1/P2 未清零时不得标记完成。
- 验收记录和截图目录见 `docs/release/acceptance/2026-06-24-p111-high-fidelity-reference-redesign.md` 与 `docs/release/ui-audit-assets/2026-06-24-p111-high-fidelity-reference-redesign/`。
- P111 不新增投资业务能力，不改变交易安全边界，不声称券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复或收益承诺。

### P64 发布打包与版本标记候选

目标：在 P63 `release_ready` 基础上，规划本地发布包、版本标记、跨机复验入口和交付校验。

状态：候选阶段，尚未创建 OpenSpec change。建议 change id：`p64-release-packaging-version-tagging`。

范围：

- 固化本地发布包内容清单、版本号/commit 标记和校验入口。
- 复用 P63 acceptance 与 P52 G0-G9 作为发布前门禁，不扩大 release_ready 声明。
- 明确跨机复验步骤、临时配置边界、脱敏诊断和回滚说明。
- 不新增后端业务能力、SQLite schema、Eino workflow、LLM 能力或任何交易/外推能力。

### P110 视觉系统重设计

目标：在 P102-P104 后验收稳定的前端基础上，将整体审美从后台台账质感升级为冷静的投资纪律研究终端。

状态：P110 已创建为 `p110-visual-system-redesign` 活跃 change。

范围：

- 视觉方向采用 Calm Command Center 为主，Ledger Pro 用于 Evidence、Decision Loop、Audit 等证据与闭环页面。
- 统一 AppLayout、Workbench、Data Quality、Risk Alerts、Evidence、Decision Loop 等核心页面的视觉 token、surface、状态层级和响应式质量。
- 保持原 API DTO、路由语义、人工确认流程和安全边界，不新增后端能力或交易/外推/自动规则能力。

验收：

- 前端 tests/build、Go tests、OpenSpec strict、forbidden copy scan、敏感信息 scan 和 `git diff --check`。
- 真实本地后端 + Vite 前端浏览器验收，覆盖 390px、768px、1280px 和桌面核心路由截图。
- 验收材料记录到 `docs/release/acceptance/2026-06-24-p110-visual-system-redesign.md` 与 `docs/release/ui-audit-assets/2026-06-24-p110-visual-system-redesign/`。

## 5. 统一安全边界

P58-P63 均不得新增或暗示：

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

## 6. 固定审查门禁

每个后续阶段必须遵守：

1. 创建独立 OpenSpec change。
2. 方案写完后子 agent 审查。
3. 审查无 Critical / Important 后执行。
4. 执行完成后执行真实 UI / 自动化验证。
5. 执行后子 agent 审查。
6. 审查无 Critical / Important 后 archive。
7. archive 后提交前子 agent 审查。
8. 审查无 Critical / Important 后提交。
