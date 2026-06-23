## ADDED Requirements

### Requirement: Data writes must declare atomic fact units

系统 SHALL 在数据模型和仓储接口中明确跨表事实写入的原子边界。

#### Scenario: Atomic write succeeds

- **WHEN** 一个原子事实单元写入成功
- **THEN** 该事实单元涉及的所有表必须全部可查询

#### Scenario: Atomic write fails

- **WHEN** 一个原子事实单元中任一表写入失败
- **THEN** 该事实单元涉及的所有新记录必须回滚
- **AND** 调用方必须收到结构化应用错误或可映射错误

### Requirement: Data records must use stable ID and time conventions

系统 SHALL 对事实表 ID 与时间字段采用统一规则。

#### Scenario: Record is inserted

- **WHEN** 插入事实表记录
- **THEN** 主键 ID 必须来自统一 ID 规则或外部已验证输入
- **AND** 时间字段必须使用 UTC RFC3339 或 SQLite 可解析的 UTC 时间字符串

### Requirement: Repository errors must be classified

系统 SHALL 将仓储错误分为 not_found、conflict、invalid_state、constraint、internal 等分类。

#### Scenario: Record is not found

- **WHEN** 仓储读取不存在的记录
- **THEN** 返回错误必须能映射为 not_found

#### Scenario: State transition is invalid

- **WHEN** 仓储拒绝非法状态流转
- **THEN** 返回错误必须能映射为 invalid_state

#### Scenario: Unique constraint fails

- **WHEN** SQLite 唯一约束或主键冲突
- **THEN** 返回错误必须能映射为 conflict 或 constraint
