# Design: P5 前端驾驶舱

## Context

P0-P4 已提供工程骨架、数据底座、领域规则、工作流和 HTTP API。P5 只实现前端驾驶舱和支撑页面，字段、状态、错误码以 `docs/frontend-contract.md`、`docs/api.md`、`docs/ui-design.md`、`docs/ui-flow.md` 为准。

P5 前置架构治理已完成：前端可使用 feature/shared 结构扩展，后端 handler、事务、ID、时间和枚举边界已有统一路径。本阶段重点是把稳定 API 契约转成可审计、可理解、无自动交易入口的 Web 控制台。

## Goals / Non-Goals

**Goals:**

- 建立前端类型和 service 层，让页面只依赖 HTTP API DTO。
- 实现三栏 Agent 决策驾驶舱，优先展示纪律状态、禁止事项、证据和确认边界。
- 实现决策详情、证据、规则、审计、持仓、设置、复盘页面。
- 为 409、500、503 等错误提供一致页面状态，避免展示底层错误细节。
- 对复杂状态映射、交易边界和审计展示写必要中文注释。

**Non-Goals:**

- 不新增后端接口或改变 P4 API 行为。
- 不实现自动交易、券商接入、一键买卖或自动跟随。
- 不把 Agent 分析观点作为最终裁决展示。
- 不引入 `docs/development-plan.md` P5 以外的页面或功能。

## Decisions

### 1. 类型和 service 先行

P5.1 先完成 `web/src/types/*` 与 `web/src/services/*`。页面和组件只消费这些类型与 service 返回值。

- 原因：P5 页面多，先固定类型和错误映射可减少字段漂移。
- 替代方案：页面中临时定义类型。放弃原因是会让字段和错误状态分散，后续审计困难。

### 2. 页面组合与业务组件分离

`web/src/pages/*` 保持路由组合、加载态和错误态；业务展示放在 `web/src/components/*` 或已准备的 `web/src/features/*`。

- 原因：符合 `docs/frontend-contract.md` 和已完成的 P5.0 架构治理。
- 替代方案：把业务逻辑集中在页面文件。放弃原因是页面会过长，也不利于复用和测试。

### 3. 驾驶舱采用三栏信息架构

`CockpitLayout` 组织左侧导航与系统状态、中间裁决工作区、右侧证据与规则面板。首屏优先展示纪律状态、触发规则、最终裁决、账户摘要和证据摘要。

- 原因：符合 `docs/ui-design.md` 三栏结构和 `docs/ui-flow.md` 首页审计要求。
- 替代方案：普通仪表盘卡片网格。放弃原因是风险、证据和裁决链路不够清晰。

### 4. 错误码映射为页面状态

API client 统一解析 `request_id`、`data`、`error`、`meta`。`DATA_REQUIRED` 映射首次使用，`EVIDENCE_NOT_FOUND` 等 409 映射信息不足或冻结观察，500/503 展示通用失败或数据源异常，不展示底层错误文本。

- 原因：契约要求前端只依赖稳定错误码，不解析底层错误细节。
- 替代方案：由每个页面单独处理错误。放弃原因是状态文案容易不一致。

### 5. 用户确认只记录线下动作

确认区只展示 `planned`、`executed_manually`、`watch`、`marked_error`。`executed_manually` 文案必须说明仅记录用户已完成的线下交易。

- 原因：产品边界禁止自动交易。
- 替代方案：按钮文案使用“买入/卖出”。放弃原因是会误导用户以为系统可代为执行。

## Risks / Trade-offs

- 页面字段多 → 先实现共享类型和格式化工具，组件按契约字段消费。
- 错误状态分支多 → 统一错误码映射，页面只展示标准状态对象。
- 审计页信息密度高 → 使用时间线与详情展开，默认折叠内部 ID、hash 和长文本。
- P5.0 已完成但仍在 P5 阶段内 → tasks.md 标为前置确认项，保留验收命令作为基线。
- 中文注释过量会干扰代码阅读 → 只注释非显然的业务边界、状态映射和审计约束。

## Migration Plan

1. 确认 P5.0 前置架构治理已完成，执行 `go test ./...` 与 `cd web && npm run build` 作为基线。
2. 实现 P5.1 类型和 service，前端构建通过后再进入页面实现。
3. 实现 P5.2 今日纪律驾驶舱与确认区。
4. 实现 P5.3 支撑页面和组件。
5. 每个小节完成后执行 `cd web && npm run build`；归档前确认 specs delta 只含 P5 前端增量。

## Open Questions

无。本变更不引入 `docs/development-plan.md` 之外的新需求。