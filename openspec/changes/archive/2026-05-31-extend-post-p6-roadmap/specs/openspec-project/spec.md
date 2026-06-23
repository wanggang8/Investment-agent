## ADDED Requirements

### Requirement: OpenSpec project mapping must include post-P6 phases

系统 SHALL 在 `openspec/project.md` 中记录 P7、P8、P9 与建议 change id 的映射关系，确保后续阶段仍按 OpenSpec change 推进。

#### Scenario: Post-P6 mapping exists

- **WHEN** P6 之后的总计划被补充
- **THEN** `openspec/project.md` 必须包含 P7、P8、P9 的阶段名称和建议 change id
- **AND** P7 必须映射到 `p7-real-data-integration`
- **AND** P8 必须映射到 `p8-frontend-experience-tests`
- **AND** P9 必须映射到 `p9-review-automation-delivery`

#### Scenario: Process remains OpenSpec-driven

- **WHEN** P7、P8 或 P9 开始实施
- **THEN** 每个阶段仍必须通过独立 OpenSpec change 生成 proposal、design、specs 和 tasks
- **AND** 不得把未归档 change 的 design 或临时计划当作验收真源
