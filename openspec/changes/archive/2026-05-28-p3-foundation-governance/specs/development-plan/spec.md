## ADDED Requirements

### Requirement: Development plan must track foundation governance phase

系统 SHALL 在开发计划中显式记录基础治理阶段，避免 P4/P5 在缺少统一基础能力时提前推进。

#### Scenario: Foundation phase exists

- **WHEN** 项目进入 P4 前的治理加固
- **THEN** `docs/development-plan.md` 必须包含 `P3-foundation` 或等价阶段
- **AND** 必须说明该阶段依赖前置 P3 归档结果

### Requirement: Governance process must remain OpenSpec-driven

系统 SHALL 继续以 OpenSpec change 作为文档与任务的唯一变更入口。

#### Scenario: Foundation governance change is active

- **WHEN** `p3-foundation-governance` 处于 active 状态
- **THEN** 变更范围、任务、delta 与归档说明必须全部在 OpenSpec 变更包内维护
- **AND** 不得只修改 `docs/` 而绕过 delta
