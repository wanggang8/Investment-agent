## Context

P0-P3 已完成工程骨架、SQLite 事实表、Repository、领域规则和 Eino 工作流。终轮复审确认功能主干可归档，但也暴露出横向基础能力缺少统一契约：错误码散在工作流和仓储层，ID 与时间由业务代码直接拼接，事务边界主要由实现经验决定，审计字段和测试策略没有形成清晰标准。

P4 将开始 HTTP API，P5 将接入前端业务页面。如果不先统一基础治理，错误响应、前端提示、审计追踪和契约测试会在 API/前端阶段重复调整。

约束：

- 不直接编辑 L1 契约文件；本变更只写 delta，archive 时合并。
- 不引入复杂外部依赖。
- 不提前实现 P4/P5 业务接口和真实前端页面。
- 保留 P0-P3 已通过的业务行为，优先做兼容迁移。

## Goals / Non-Goals

**Goals:**

- 建立统一错误包，覆盖错误码、错误类型、分类、包装、HTTP 映射和审计映射。
- 建立统一 ID 与时间基础包，保证生产代码可读、测试可控。
- 明确事务边界和仓储方法命名规则，避免跨表半写入。
- 明确审计事件枚举与字段填写规范，减少节点审计漂移。
- 明确测试策略，约束 P0-P3 回归和 P4/P5 以后契约测试。
- 同步规划总架构与治理文档，确保 `docs/architecture.md`、`docs/development-plan.md`、`docs/GOVERNANCE.md`、`openspec/project.md` 在归档时体现基础治理层。
- 迁移 P0-P3 已有关键路径，不大规模改写业务逻辑。

**Non-Goals:**

- 不实现 P4 HTTP 业务 API。
- 不实现 P5 前端业务页面。
- 不引入分布式追踪、OpenTelemetry 或复杂日志平台。
- 不更换 SQLite、Eino、React 技术选型。
- 不把所有历史测试一次性改成完整矩阵，只补影响基础治理可信度的必要测试。

## Decisions

### 1. 统一错误包放在 `internal/pkg/apperr`

选择 `internal/pkg/apperr`，而不是继续把错误放在各业务包。

理由：错误码会被 workflow、repository、HTTP handler、audit writer 和前端契约共同使用，需要一个不依赖业务层的稳定基础包。

核心类型：

```go
type Code string
type Category string
type AppError struct {
    Code Code
    Category Category
    Message string
    HTTPStatus int
    Retryable bool
    Cause error
}
```

约束：

- 业务层只返回 `error`，但可用 `apperr.As` / `errors.As` 读取结构化错误。
- 原有字符串错误码保留为兼容常量，逐步迁移到 `apperr.Code`。
- 仓储层的状态非法、未找到、冲突、约束失败必须映射为统一分类。
- HTTP 层只根据 `AppError` 生成响应信封，不解析底层 SQL 或 workflow 细节。

备选方案：

- 放在 `pkg/apperr`：对外暴露意味太强，本项目暂不做公共库。
- 每层各自定义错误：短期省事，但 P4/P5 会形成重复映射。

### 2. ID 与时间分别放在 `internal/pkg/idgen` 和 `internal/pkg/clock`

ID 生成与时间获取从业务代码中抽离。

`clock` 提供：

- `Clock` 接口：`Now() time.Time`
- `SystemClock`
- `FixedClock`（测试用）
- `FormatRFC3339UTC(time.Time) string`

`idgen` 提供：

- 实体 ID 生成函数：`DecisionID(requestID)`、`EvidenceRefID(decisionID, index)`、`AuditEventID(requestID, nodeName, index)` 等。
- 规则：可读、稳定、无空值；测试可预测；生产可在后续加入随机后缀。

备选方案：

- 使用 UUID：通用但降低测试可读性，当前本地单机阶段不必强制。
- 继续在业务代码拼接：会导致规则分散，难以审计。

### 3. 事务边界用仓储组合方法表达

跨多张事实表的写入必须由仓储提供组合方法，业务 Graph 不直接串联多个仓储写入形成事实单元。

必须原子化的事实单元：

- 决策记录 + 证据引用 + DecisionRecordNode 审计。
- 用户确认 + 相关事实表 + 确认审计 + 决策确认状态更新。
- 证据核查的 item + summary + rag chunks + source verification。
- 守门人审计 + 提案状态推进。
- 规则应用：旧 active 归档 + 新 active 创建 + 提案 applied 状态。

允许部分成功的场景必须在契约中明确。例如市场刷新多标的可按标的独立写入，但每个标的一次刷新必须原子化。

### 4. 审计契约采用“枚举 + 节点字段规则”

审计事件继续写 `audit_events`，但新增枚举约束文档和代码常量。

规则：

- 每个节点必须填写 `action`、`node_name`、`node_action`、`status`、`input_ref_type/input_ref`。
- 失败和有原因的降级必须填写 `error_code`。
- 产生持久化结果时必须填写 `output_ref_type/output_ref`。
- Graph 级事件可以存在，但不能替代主链路节点审计。
- 审计错误码必须来自统一错误码或兼容映射表。

### 5. 测试策略按层约束

基础治理变更新增最低测试线：

- `apperr`：错误包装、分类、HTTP 映射、审计映射。
- `idgen/clock`：ID 规则、UTC 格式、固定时钟。
- Repository：事务成功、失败回滚、字段级断言。
- Workflow：正常、失败、降级、终态跳过、审计字段。
- API（P4）：响应信封、错误码、HTTP 状态、DTO 字段。
- Frontend（P5）：错误码到 UI 状态映射，不直接依赖底层错误文本。

### 6. 文档同步策略

本变更归档时必须同步以下文档，不只合并代码相关 delta：

- `docs/architecture.md`：补充基础治理包在分层架构中的位置，说明 `apperr`、`idgen`、`clock`、事务边界、审计契约和测试策略如何贯穿 HTTP、workflow、domain、repository、infrastructure。
- `docs/development-plan.md`：补充 `P3-foundation` 阶段，说明它位于 P3 与 P4 之间，是进入 HTTP API 前的治理加固阶段。
- `docs/GOVERNANCE.md`：归档后更新活跃变更表，并说明横向基础治理变更也必须走 OpenSpec delta。
- `openspec/project.md`：更新阶段映射，加入 `p3-foundation-governance` 与 P4 的依赖关系。
- L1 契约：将 foundation、workflow、data-model、api、frontend-contract 的 delta 合并到对应章节。

该同步动作放在任务清单中单独验收，避免只改代码不改总设计文档。

## Risks / Trade-offs

- 迁移范围扩大 → 只迁移 P0-P3 已有关键路径，避免提前做 P4/P5。
- 错误码命名与已有常量冲突 → 保留兼容别名，测试覆盖映射。
- ID 规则后续可能变化 → 统一入口，后续只改 `idgen`。
- 事务组合方法增加仓储接口数量 → 只为真实跨表事实单元增加，不为单表写入增加包装。
- 审计字段测试变多 → 优先覆盖主链路和已出现问题的辅助 Graph。

## Migration Plan

1. 新增基础包：`apperr`、`idgen`、`clock`。
2. 将 workflow 错误码迁移为 `apperr.Code` 兼容常量。
3. 将仓储业务错误迁移为 `apperr.AppError`，保留 `errors.Is` 可判断能力。
4. 将 workflow 中显式拼接 ID 和 `time.Now().UTC()` 的关键路径替换为基础包。
5. 为已存在的跨表写入补齐组合仓储方法和测试。
6. 补充审计字段与错误映射测试。
7. 更新总架构、开发计划、治理说明和 OpenSpec 项目说明的 delta。
8. 全量执行 `go test ./...` 和 `web npm run build`。

## Open Questions

无阻塞问题。默认采用本设计中的保守策略：不新增外部依赖、不提前实现 P4/P5，只固化横向基础能力并迁移 P0-P3 关键路径。
