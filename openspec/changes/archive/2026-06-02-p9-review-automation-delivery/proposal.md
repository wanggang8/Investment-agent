## Why

P9 是开发计划中的最后阶段，需要把已完成的每日纪律、数据刷新、证据索引、复盘和交付能力串成本地可长期使用的入口。当前系统已有后端、前端、数据、工作流和测试基础，但还缺少 `cmd/agent` 任务入口、周期复盘自动化、规则有效性评估和本地运维交付说明。

## What Changes

- 新增 `cmd/agent` 本地任务入口，支持手动触发每日纪律、市场刷新、情报索引和复盘任务。
- 增加本地调度配置读取能力，但默认不启用自动交易，也不执行任何下单行为。
- 增强任务执行审计：每次本地任务记录输入摘要、状态、错误码和输出引用。
- 完善月度/季度复盘，汇总确认动作、错误案例、规则提案、规则命中、误判、缺证据、降级和审计事件。
- 将规则有效性评估结果写入规则提案或复盘摘要，但仍必须经过守门人审计和用户最终确认。
- 增强前端复盘页，展示周期摘要、规则建议和追踪入口。
- 补充本地启动、初始化、备份、索引重建、恢复、常见故障和 P7-P9 完整验收命令说明。
- 实现代码中的调度、安全边界和审计写入逻辑需要配中文注释。

## Capabilities

### New Capabilities
- `review-automation-delivery`: 覆盖 P9 的本地任务入口、周期复盘、规则有效性评估、交付说明与安全边界。

### Modified Capabilities
- `real-data-integration`: 增加本地任务入口对市场刷新、情报索引、VecLite 重建和故障恢复的触发与运维要求。
- `frontend-experience-tests`: 增加复盘页周期摘要、规则建议和追踪入口的前端展示要求。
- `e2e-hardening`: 增加 P9 完整验收命令与本地交付验证要求。

## Impact

- 后端：`cmd/agent`、workflow/service/handler/repository 的复盘查询与任务触发逻辑、审计写入。
- 前端：复盘页、规则建议展示、追踪入口、相关类型和服务测试。
- 数据：复用现有 SQLite 表和审计事件，不引入自动交易表或下单接口。
- 文档：`docs/configuration.md`、`docs/migration-plan.md`、`docs/testing-plan.md` 和本地交付说明。
- 验证：`go test ./...`、`go run ./cmd/agent --help`、`go test ./internal/application/workflow/... ./internal/application/handler/...`、`cd web && npm run build`。

## Scope Mapping

- In scope：P9.1、P9.2、P9.3 的全部任务和验收命令。
- Out of scope：自动交易、券商接口、真实下单、绕过守门人审计或用户最终确认的规则应用、development-plan 之外的新功能。
- Plan 对应关系：本 change 的 tasks.md 将逐条覆盖 P9.1、P9.2、P9.3 的任务列表与验收命令，一一对应，不新增计划外需求。
