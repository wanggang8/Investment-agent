## ADDED Requirements

### Requirement: Governance docs must reflect active changes and archive rules

系统 SHALL 在治理文档和 OpenSpec 项目说明中同步活跃变更、归档规则和阶段映射。

#### Scenario: Active change is updated

- **WHEN** 基础治理变更成为活跃变更
- **THEN** `docs/GOVERNANCE.md` 必须反映当前活跃 change
- **AND** `openspec/project.md` 必须反映阶段与 change 映射

#### Scenario: Archive completes

- **WHEN** 变更完成归档
- **THEN** 活跃变更必须从治理文档中移除
- **AND** 归档位置必须可追溯
