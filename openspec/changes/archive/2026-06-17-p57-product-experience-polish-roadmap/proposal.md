# P57 产品体验打磨总规划

## Why

P55 真实 UI 验收证明项目主要功能可通过前端操作完成，但也暴露出产品体验问题：页面更像工程功能集合，缺少清晰的日常投资纪律产品主线。P56 已修复真实 LLM 决策详情崩溃、移动端溢出和基础 UI 产品化问题，当前 full UI acceptance 状态为 `p56_scope_pass`。

用户明确要求继续打磨产品设计、UI 设计和功能设计。因此 P57 不应直接进入发布状态刷新，而应先建立完整的产品体验打磨路线图，把后续改造拆成可审查、可验收、可真实 UI 验证的阶段。原 P57 `p57-release-readiness-refresh` 后移到产品体验打磨完成后再执行。

## What Changes

- 新增产品体验北极星：
  - 本地投资纪律工作台，而不是券商交易终端、AI 聊天 demo 或工程调试后台。
  - 每日核心问题为：今天能不能动、为什么、我需要人工做什么、数据和规则是否可信。
- 固化 P58-P63 后续阶段：
  - P58 今日工作台重构。
  - P59 决策解释体验重构。
  - P60 组合、风险与数据质量体验重构。
  - P61 治理和运维页面产品化，覆盖 Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings。
  - P62 设计系统与可访问性验收。
  - P63 全量真实 UI 回归与发布状态刷新。
- 定义产品设计、UI 设计和功能设计的统一验收门禁：
  - Product Design skill 复核。
  - 子 agent 方案审查、执行后审查、提交前审查。
  - 真实启动后端和前端，通过浏览器操作核心流程。
  - 桌面、平板/窄桌面、390px 移动端截图和无横向溢出检查。
  - 安全边界和敏感信息扫描。
- 更新治理文档，说明 P57 从发布刷新改为产品体验路线图，发布刷新后移。

## Scope

- OpenSpec change、设计路线图、后续阶段拆分、验收门禁和安全边界。
- `docs/development-plan.md`、`openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`AGENTS.md` 的阶段状态更新。
- `docs/frontend-contract.md` / `openspec/specs/frontend-experience-tests/spec.md` 中新增产品体验打磨路线图要求摘要。

## Out of Scope

- P57 不修改运行时代码、SQLite schema、HTTP API、Eino 工作流或 React 页面。
- P57 不直接重构 Dashboard、Workbench、Decision Detail、Positions、Data Quality、Rules、Audit 或 Settings。
- P57 不声明所有产品设计问题已修复，也不刷新最终 release-ready 口径。
- 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
- 不将 Product Design 调研资料替代 OpenSpec change 或 `docs/` 契约真源。

## Validation

- `openspec validate p57-product-experience-polish-roadmap --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 敏感信息扫描：确认没有 key、完整 prompt、私有 SQLite、raw vendor payload、供应商原始响应或敏感路径泄露。
- 子 agent 方案复审必须覆盖：
  - 是否遵守 `docs/GOVERNANCE.md`。
  - 是否使用 Product Design skill 和调研依据。
  - P58-P63 是否边界清晰、可独立验收。
  - 是否保持不交易、不外推、不自动确认、不自动规则应用、不收益承诺边界。
  - 是否没有把发布刷新提前到产品体验打磨前。
