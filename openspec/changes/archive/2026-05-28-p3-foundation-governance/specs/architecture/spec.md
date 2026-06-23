## ADDED Requirements

### Requirement: Architecture document must describe foundation governance

系统 SHALL 在总架构文档中说明统一错误、ID、时间、事务、审计和测试策略在分层架构中的位置。

#### Scenario: Architecture is updated for foundation governance

- **WHEN** 新的基础治理能力被引入
- **THEN** `docs/architecture.md` 必须说明 `apperr`、`idgen`、`clock`、事务边界、审计契约和测试策略的分层位置
- **AND** 必须说明这些能力如何贯穿 HTTP、workflow、domain、repository 和 infrastructure

### Requirement: Development plan must include foundation phase

系统 SHALL 在开发计划中增加基础治理阶段，并说明它位于 P3 与 P4 之间。

#### Scenario: Plan is updated

- **WHEN** 基础治理变更包完成并归档
- **THEN** `docs/development-plan.md` 必须记录 `P3-foundation` 或等价阶段
- **AND** 必须说明该阶段的目标、产物和验收方式

### Requirement: Governance docs must reflect active change history

系统 SHALL 在治理文档和 OpenSpec 项目说明中记录活跃变更和归档规则。

#### Scenario: Active change is updated

- **WHEN** `p3-foundation-governance` 成为活跃变更
- **THEN** `docs/GOVERNANCE.md` 和 `openspec/project.md` 必须更新活跃变更与阶段映射

#### Scenario: Change is archived

- **WHEN** 该变更归档
- **THEN** 活跃变更表必须移除该 change
- **AND** 归档记录必须可追溯
