## ADDED Requirements

### Requirement: Development plan must include post-P6 roadmap

系统 SHALL 在 `docs/development-plan.md` 中记录 P6 之后的剩余阶段，确保真实数据、RAG/DeepSeek、前端体验和周期复盘等已知剩余能力有可执行计划。

#### Scenario: Post-P6 phases are planned

- **WHEN** P6 已完成且项目仍存在 requirements 与 architecture 中定义但未充分实现的能力
- **THEN** `docs/development-plan.md` 必须包含 P7、P8、P9 三个后续阶段
- **AND** 每个阶段必须包含目标、主要产物、任务分解和验收方式

#### Scenario: P7 covers data and analyst foundation

- **WHEN** 总计划补充 P7
- **THEN** P7 必须覆盖真实行情与情报数据接入、RAG/VecLite 检索、DeepSeek 分析师材料、配置降级和审计事件
- **AND** 必须继续声明 DeepSeek 不生成最终裁决

#### Scenario: P8 covers frontend experience and tests

- **WHEN** 总计划补充 P8
- **THEN** P8 必须覆盖驾驶舱图表、关键页面交互体验、前端错误/空态、前端测试和构建验收
- **AND** 不得加入自动交易入口

#### Scenario: P9 covers review and local delivery

- **WHEN** 总计划补充 P9
- **THEN** P9 必须覆盖 `cmd/agent` 本地任务、每日任务入口、月度/季度复盘、规则有效性评估、配置示例和交付说明
- **AND** 必须保持所有关键动作可审计

### Requirement: Development plan must preserve completed phase records

系统 SHALL 在补充 P7–P9 时保留 P0–P6 已完成记录，避免重写已归档阶段的验收事实。

#### Scenario: Completed phases remain unchanged

- **WHEN** P7–P9 被加入开发计划
- **THEN** P0–P6 的阶段状态、已完成任务和已归档事实必须保持一致
- **AND** 新增内容必须作为后续阶段追加到计划中
