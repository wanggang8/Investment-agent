# Proposal: P5 前端驾驶舱

## Summary

创建 `docs/development-plan.md` 中 P5 前端驾驶舱变更包，用于实现前端类型、API client、今日纪律驾驶舱，以及决策详情、证据、规则、审计、持仓、设置、复盘页面。

## Why

P4 HTTP API 已归档，P5 需要把已有后端契约转成可使用的本地 Web 控制台。前端必须以纪律、证据、规则和审计为主线，避免出现自动交易入口或误导性收益承诺。

## What Changes

- 定义前端通用响应类型，以及 dashboard、portfolio、decision、evidence、rule、audit、settings、market、review DTO。
- 统一 API client 响应信封解析、错误码映射和 409/500/503 展示状态。
- 实现三栏 Agent 决策驾驶舱，展示纪律状态、风险红线、今日建议、账户摘要、证据摘要和用户确认区。
- 实现决策详情、证据、规则、审计、持仓、设置、复盘页面及对应业务组件。
- 保持前端只通过 HTTP API 获取数据，不访问 SQLite、VecLite 或本地文件。
- P5 新增实现代码编写必要中文注释，用于说明非显然的业务约束、状态映射和确认边界；这是执行约束，不作为契约 delta。

## In Scope

- P5.0：确认已完成的架构治理前置项，包括 feature/shared 目录准备、handler 分层、事务协调、共享 ID/时间与枚举复用。
- P5.1：创建或确认 `web/src/types/*` 与 `web/src/services/*` 文件，定义 DTO 与 API client。
- P5.1：统一处理 409、500、503 错误，并映射到前端页面状态。
- P5.2：实现 `DashboardPage`、`CockpitLayout` 与驾驶舱组件。
- P5.2：展示信息不足、冻结观察、用户确认动作边界，并确保没有自动交易入口。
- P5.3：实现决策详情、证据、规则、审计、持仓、设置、复盘页面及组件。
- P5.3：展示规则提案 `pending_final_confirm` 和最终确认动作。
- P5.3：审计页区分 `action`、`node_name`、`node_action`，并展示 `status`、`error_code`、输入引用、输出引用。
- 按 P5 每个小节执行 `cd web && npm run build`；P5.0 保留 `go test ./...` 与前端构建验收。

## Out of Scope

- 不新增 `docs/development-plan.md` 以外的前端页面、后端 API 或业务流程。
- 不修改 P4 已归档 HTTP API 行为。
- 不访问 SQLite、VecLite 或本地文件作为前端数据源。
- 不接入券商交易、自动下单、自动跟随或一键交易能力。
- 不把 DeepSeek 或 Agent 分析观点展示为最终裁决来源。
- 不展示完整密钥、SQL、文件路径或外部服务原始错误文本。

## Capabilities

### New Capabilities

- `frontend-cockpit`: 覆盖 P5 前端驾驶舱、类型与 API client、决策详情、证据、规则、审计、持仓、设置、复盘页面的本阶段增量契约。

### Modified Capabilities

- 无。当前 `openspec/specs/` 没有可修改的既有 capability；本变更仅写 P5 阶段 delta。

## Impact

- 影响前端目录：`web/src/types/`、`web/src/services/`、`web/src/pages/`、`web/src/components/`，以及已准备的 `web/src/features/` 与 `web/src/shared/`。
- 依赖后端已归档的 P4 HTTP API 响应信封、DTO 字段和稳定错误码。
- 验收以 `docs/development-plan.md` P5 各小节命令为准：P5.0 执行 `go test ./...` 与前端构建；P5.1–P5.3 执行前端构建。

## Plan Alignment

本 change 与 `docs/development-plan.md` P5 小节一一对应：

- P5.0 架构治理准备 → tasks.md 第 1 节，作为已完成前置确认与验收基线。
- P5.1 类型与 API client → tasks.md 第 2 节，覆盖全部类型、service 文件、错误状态映射和数据源边界。
- P5.2 Agent 决策驾驶舱 → tasks.md 第 3 节，覆盖全部页面、布局、组件、状态和确认动作要求。
- P5.3 决策详情、证据、规则与审计页面 → tasks.md 第 4 节，覆盖全部页面、组件和字段展示要求。
- P5 验收命令 → tasks.md 第 5 节，逐条保留 development-plan 中的验收命令。

未加入 `docs/development-plan.md` 之外的需求；本次指令中的“实现代码要写好中文注释”只作为 tasks.md 执行约束，不写入 specs 契约 delta。