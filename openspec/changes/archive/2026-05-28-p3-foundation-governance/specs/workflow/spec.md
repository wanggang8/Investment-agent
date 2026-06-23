## ADDED Requirements

### Requirement: Workflow errors must use foundation error codes

系统 SHALL 将工作流节点错误码对齐统一应用错误体系，并保留审计与 API 映射一致性。

#### Scenario: Workflow node fails

- **WHEN** 工作流节点返回 failed
- **THEN** `NodeResult.ErrorCode` 必须能映射到统一应用错误码
- **AND** `AuditFragment.ErrorCode` 必须与节点错误码一致

#### Scenario: Workflow node degrades

- **WHEN** 工作流节点以 degraded 完成且存在明确原因
- **THEN** 节点必须返回统一应用错误码
- **AND** 后续规则裁决必须能读取降级原因

### Requirement: Workflow IDs and timestamps must be generated centrally

系统 SHALL 通过统一 ID 与时间基础包生成工作流相关 ID 和时间。

#### Scenario: Decision record ID is needed

- **WHEN** DecisionRecordNode 需要生成 `decision_id`
- **THEN** 必须通过统一 ID 生成入口生成

#### Scenario: Audit event time is recorded

- **WHEN** AuditWriter 创建审计事件
- **THEN** 必须通过统一 Clock 获取时间
- **AND** 必须写入 UTC RFC3339 字符串

### Requirement: Workflow persistence must respect transaction boundaries

系统 SHALL 在工作流写事实表时调用仓储事务方法，不得在 Graph 中拆开同一事实单元。

#### Scenario: DecisionRecordNode persists output

- **WHEN** DecisionRecordNode 保存决策记录
- **THEN** 必须通过仓储事务方法写入决策、证据引用和审计事件

#### Scenario: Supporting graph persists evidence facts

- **WHEN** EvidenceVerificationGraph 保存证据事实
- **THEN** 必须调用证据事实事务方法

#### Scenario: Supporting graph advances rule proposal

- **WHEN** GatekeeperAuditGraph 推进提案状态
- **THEN** 必须调用守门人审计与状态更新事务方法
