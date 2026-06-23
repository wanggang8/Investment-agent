## ADDED Requirements

### Requirement: Unified application errors

系统 SHALL 使用统一应用错误体系表达业务错误、基础设施错误和工作流错误，并为 API、审计和测试提供稳定映射。

#### Scenario: Error is created with a stable code

- **WHEN** 业务代码创建应用错误
- **THEN** 错误必须包含稳定 `code`
- **AND** 必须包含 `category`
- **AND** 必须包含可面向用户或日志的默认 `message`

#### Scenario: Error maps to HTTP response

- **WHEN** HTTP 层接收到应用错误
- **THEN** 必须根据错误映射得到 HTTP status
- **AND** 响应信封中的 `error.code` 必须等于应用错误码
- **AND** 不得把底层 SQL、文件路径或外部服务原始错误直接返回给前端

#### Scenario: Error maps to audit event

- **WHEN** 工作流节点失败或降级
- **THEN** 审计事件的 `error_code` 必须来自统一错误码或兼容映射
- **AND** 降级原因必须可通过错误码追踪

#### Scenario: Wrapped error keeps cause

- **WHEN** 基础设施错误被包装为应用错误
- **THEN** 调用方必须能通过标准错误机制识别原始 cause
- **AND** 调用方必须能读取应用错误码和分类

### Requirement: Unified ID and time generation

系统 SHALL 通过统一基础包生成关键实体 ID 和时间字符串，禁止业务路径散落拼接规则和直接格式化时间。

#### Scenario: Entity ID is generated

- **WHEN** 工作流或仓储需要生成决策、证据引用、审计事件或规则应用 ID
- **THEN** 必须通过统一 ID 生成入口生成
- **AND** 生成结果必须非空、可读、可测试

#### Scenario: Time is recorded

- **WHEN** 系统写入 `created_at`、`updated_at`、`executed_at` 或审计时间
- **THEN** 默认必须使用 UTC
- **AND** 输出字符串必须使用 RFC3339
- **AND** 测试必须能注入固定时间

### Requirement: Repository transaction boundaries

系统 SHALL 在仓储边界定义跨表事实写入的事务单元，禁止应用层把同一事实单元拆成多个非原子写入。

#### Scenario: Decision record is persisted

- **WHEN** DecisionRecordNode 保存正式决策
- **THEN** `decision_records`、`evidence_refs` 和该节点 `audit_events` 必须作为一个事实单元成功或失败

#### Scenario: Evidence facts are persisted

- **WHEN** EvidenceVerificationGraph 保存证据事实
- **THEN** `intelligence_items`、`intelligence_summary`、`rag_chunks` 和 `source_verifications` 必须在同一事务中成功或失败

#### Scenario: Gatekeeper audit advances proposal

- **WHEN** GatekeeperAuditGraph 保存审计结论并推进提案状态
- **THEN** `gatekeeper_audits` 与 `rule_proposals.status` 必须在同一事务中成功或失败

#### Scenario: Rule version is applied

- **WHEN** 用户最终确认应用规则提案
- **THEN** 旧 active 规则归档、新 active 规则创建和提案状态更新必须在同一事务中成功或失败
- **AND** 仓储必须保证新规则版本状态为 `active`

### Requirement: Auditable event contract

系统 SHALL 统一审计事件字段和枚举，确保工作流、用户动作、规则演进和系统维护动作可追踪。

#### Scenario: Workflow node emits audit fragment

- **WHEN** 工作流节点完成、降级或失败
- **THEN** 审计片段必须包含 `action`、`node_name`、`node_action`、`status`、`input_ref_type`、`input_ref`
- **AND** 失败或有原因的降级必须包含 `error_code`

#### Scenario: Node creates persisted output

- **WHEN** 节点写入持久化事实
- **THEN** 审计片段必须包含 `output_ref_type` 与 `output_ref`

#### Scenario: Graph level audit is emitted

- **WHEN** 辅助 Graph 只产生 Graph 级审计事件
- **THEN** 事件必须说明 Graph 名称、业务动作、输入引用和输出引用
- **AND** 不得替代 Daily/Consultation 主链路的节点审计

### Requirement: Layered verification strategy

系统 SHALL 按层定义最低验证要求，避免只依赖数量断言或单一 happy path。

#### Scenario: Foundation package is tested

- **WHEN** 新增错误、ID、时间基础包
- **THEN** 必须覆盖错误映射、包装、ID 规则、UTC/RFC3339 和固定时钟测试

#### Scenario: Repository transaction is tested

- **WHEN** 仓储方法写入多个事实表
- **THEN** 必须覆盖成功字段级断言
- **AND** 必须覆盖失败回滚断言

#### Scenario: Workflow branch is tested

- **WHEN** 工作流存在成功、失败、降级或终态跳过分支
- **THEN** 必须断言最终状态、持久化事实、审计事件和错误码

#### Scenario: API and frontend are implemented later

- **WHEN** P4/P5 实现 API 与前端页面
- **THEN** 必须验证 API 响应信封、错误码到 HTTP 状态映射、DTO 字段和前端展示状态映射
