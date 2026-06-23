## Why

P0-P3 已完成工程骨架、数据底座、领域规则与 Eino 工作流，但错误、ID、时间、事务、审计和测试策略仍分散在各层实现中。若直接进入 P4 HTTP API 与 P5 前端，错误响应、前端文案、审计追踪和事务一致性会继续分化。

本变更在 P4 前补齐横向工程治理契约，统一基础能力边界，降低后续 API、前端和验收阶段的返工风险。

## What Changes

- 新增统一错误体系：定义错误码、错误类型、错误分类、HTTP 映射、审计映射和包装规范。
- 新增统一 ID 与时间规范：定义实体 ID 生成规则、时间来源、UTC/RFC3339 输出、测试时钟注入边界。
- 新增事务边界规范：明确哪些仓储写入必须原子化，哪些允许部分成功，失败时如何回滚和审计。
- 新增审计事件契约：统一 `action`、`node_name`、`node_action`、`input/output ref`、`error_code` 的枚举和填写规则。
- 新增测试策略：明确 P0-P3 已有层级与 P4/P5 以后必须满足的单测、集成测、契约测最低要求。
- 补充总架构与开发治理文档更新范围：`docs/architecture.md`、`docs/development-plan.md`、`docs/GOVERNANCE.md`、`openspec/project.md` 需要在归档时同步体现基础治理层。
- 代码新增基础包并迁移 P0-P3 已有路径，不提前实现 P4/P5 业务 API 或真实前端页面。

## Capabilities

### New Capabilities
- `foundation-governance`: 覆盖统一错误、ID、时间、事务边界、审计事件和测试策略等横向工程治理契约。

### Modified Capabilities
- `workflow`: 明确工作流节点错误、审计和事务边界与统一基础能力的关系。
- `data-model`: 明确事务一致性、ID、时间和审计字段的写入约束。
- `api`: 明确统一错误类型到 HTTP 状态与响应信封的映射，供 P4 实现使用。
- `frontend-contract`: 明确前端只消费稳定错误码和展示状态，不直接解析底层错误。
- `architecture`: 明确基础治理包在分层架构中的位置，以及错误、ID、时间、事务、审计、测试策略如何贯穿各层。
- `development-plan`: 增加 P3-foundation 阶段说明，并确保 P4/P5 依赖该基础治理完成。
- `governance`: 明确横向基础治理变更属于 OpenSpec 管理范围，归档时必须同步 L1/L2 文档。

## Impact

- 文档：新增 `foundation-governance` delta，并修改 workflow、data-model、api、frontend-contract、architecture、development-plan、governance 的相关要求。
- 后端：新增 `internal/pkg/apperr`、`internal/pkg/idgen`、`internal/pkg/clock` 等基础包；迁移 workflow 与 repository 的关键错误和 ID/时间生成路径。
- 测试：新增基础包单测，补充错误映射、ID/时间确定性、事务回滚和审计字段一致性测试。
- API/前端：不在本变更实现业务接口和页面，只为 P4/P5 提供契约边界。
