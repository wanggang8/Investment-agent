## ADDED Requirements

### Requirement: Progress tracking must advance to P7

系统 SHALL 在 `openspec/PROGRESS.md` 中记录 P6 已完成后的下一阶段，使无参数 `/opsx:propose` 能继续选择 P7。

#### Scenario: Next phase points to P7

- **WHEN** P6 已归档且 P7–P9 总计划已补充
- **THEN** `openspec/PROGRESS.md` 必须将 `next_phase` 设置为 `P7`
- **AND** 必须将 `next_change_id` 设置为 `p7-real-data-integration`
- **AND** 必须在阶段状态表中包含 P7、P8、P9

#### Scenario: Completed P6 remains recorded

- **WHEN** 进度文件补充 P7–P9
- **THEN** P6 必须继续标记为 `done`
- **AND** `current_change` 必须保持为空，直到新的 change 被创建
