# P56 UI 验收阻断与产品化设计修复

## Why

P55 已真实启动本地项目并通过浏览器操作验收主要 UI 功能，结论为 `blocked`。阻断点是：真实 LLM-backed consultation 能写入决策，但决策详情页会因 `final_verdict.optional_actions` 为 `null` 而崩溃，导致不能声明前端全功能验收通过。

同时，P55 Product Design audit 发现当前前端更像“功能路由集合/工程验收面板”，还没有形成稳定的本地投资纪律产品体验。主要问题包括：左侧导航过长且缺少任务分组、移动端侧栏挤占内容、`/positions` 与 `/data-quality` 横向溢出、表单控件偏浏览器默认、决策详情偏调试报告而非可解释决策记录。

P56 需要先解除验收阻断，再按 Product Design skill 和产品设计调研结论完成一轮高优先级 UI 产品化修复，并重跑真实 UI 验收。

## What Changes

- 修复真实 LLM 决策详情 nullable DTO 崩溃：
  - `final_verdict.optional_actions`、`final_verdict.prohibited_actions` 等列表字段为 `null`、缺失或非数组时必须安全展示。
  - 增加真实 LLM-like DTO 回归测试。
- 重构前端产品化基础体验：
  - 建立更稳定的 app shell、任务分组导航和移动端导航行为。
  - 调整全局 CSS token、字体层级、卡片/按钮/表单/表格基础样式，使其更符合本地投资纪律工作台而非默认 Vite/demo 页面。
  - 优先修复 `/positions`、`/data-quality` 移动端横向溢出。
  - 优化 `/`、`/workbench`、`/consultation`、`/decisions/:decisionId`、`/positions`、`/data-quality` 的信息层级和操作路径。
- 使用 Product Design audit/research 方法补充设计依据：
  - P56 design 必须记录产品 brief、用户目标、信息架构、视觉原则、交互原则和调研来源。
  - 子 agent 方案复审和执行后复审必须检查 Product Design skill 与调研依据是否被实际引用。
- 重跑 P55 核心 UI 验收矩阵：
  - 真实启动后端与 Vite 前端。
  - 使用浏览器操作验证阻断路径、移动端路径、核心页面路径和禁止自动交易边界。
  - 产出 P56 UI 验收记录和截图资产。

## Scope

- 前端 React/Vite/TypeScript 代码、前端测试、必要的 API DTO adapter 安全处理。
- 与 UI 产品化直接相关的 CSS、组件、布局、导航、表单、表格和页面信息层级。
- 文档与治理更新：P56 change、验收记录、设计审查记录、进度状态。
- 真实 UI 操作验收与截图留证。

## Out of Scope

- 新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺。
- 新增登录源、付费源、授权源、Level2、高频源或需要浏览器登录态的数据源。
- 改写 LLM 决策边界；LLM 仍只能生成分析材料，最终裁决仍由规则/领域逻辑产生。
- 以视觉优化名义降低风险、安全、数据质量或人工复核提示的可见性。
- 将 Product Design 研究材料替代 OpenSpec change 或 `docs/` 契约真源。
- 在 P56 中发布版本标签、打包分发或声明最终 release ready；这些必须以后续独立阶段为准。

## Validation

- 方案阶段：
  - `openspec validate p56-ui-acceptance-blocker-fixes --strict`
  - `openspec validate --all --strict`
  - `git diff --check`
  - 子 agent 方案复审无 Critical / Important，且复审覆盖 Product Design skill 与调研依据。
- 实现阶段：
  - 前端单元测试覆盖 nullable decision DTO、导航/移动端关键状态和禁止自动交易文案。
  - 前端构建通过。
  - 后端测试按受影响范围执行；若仅前端修改，仍至少执行与 DTO/HTTP fixture 相关的现有测试或说明不适用原因。
  - 真实启动本地项目并通过浏览器操作重跑 P55 核心 UI 验收矩阵。
  - 桌面与 390px 移动端截图证明 `/positions`、`/data-quality`、决策详情和核心导航无横向溢出、无崩溃。
  - 扫描提交内容，确认没有 key、完整 prompt、私有 SQLite、raw vendor payload 或敏感路径泄露。
- 归档前：
  - 子 agent 执行后复审无 Critical / Important。
  - archive 后再做提交前复审。
