## ADDED Requirements

### Requirement: OpenSpec project mapping must include foundation governance phase

系统 SHALL 在 OpenSpec 项目说明中记录基础治理阶段与后续 API 阶段的映射关系。

#### Scenario: Phase mapping is updated

- **WHEN** 基础治理变更创建或归档
- **THEN** `openspec/project.md` 必须说明该 change 的阶段定位
- **AND** 必须说明它与 P4 API 阶段的依赖关系

#### Scenario: Process remains change-driven

- **WHEN** 新阶段开始
- **THEN** 仍必须通过 OpenSpec change 生成 proposal、design、specs 和 tasks
- **AND** 不得依赖临时计划文档作为真源
