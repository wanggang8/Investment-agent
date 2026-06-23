# Proposal: P1 数据底座

## Intent

实现 Investment Agent 的 SQLite 事实基准与 Repository 访问层，为后续领域规则、Eino 工作流、HTTP API 和前端驾驶舱提供可审计、可恢复的数据基础。

## Scope

### In scope

- P1.1 SQLite migration：按 `docs/data-model.md` 创建核心表、索引、枚举 CHECK 约束。
- P1.1 种子数据：写入默认 `rule_versions`（`v3.0`）与默认信源等级配置。
- P1.2 Repository 层：定义并实现账户、持仓、决策、证据、确认、情报、规则、审计等写读能力。
- Repository 测试：覆盖写入、读取、事务失败回滚。

### Out of scope

- P2 领域规则、状态机、裁决引擎。
- P3 Eino Graph 与节点编排。
- P4 HTTP API handler。
- P5 前端页面真实数据接入。
- P6 端到端验收与加固。

## Source documents

- `docs/development-plan.md`：P1：数据底座（P1.1、P1.2）
- `docs/data-model.md`：SQLite 表、索引、约束、状态字段
- `docs/workflow.md`：Repository 需要支撑的工作流读写边界
- `openspec/project.md`、`docs/GOVERNANCE.md`

## Expected outcome

- migration 可创建空库，重复启动不会破坏已有数据。
- Repository 层有接口、SQLite 实现与测试。
- 后续 P2/P3 可以依赖 Repository 接口，不直接访问 SQLite。

## Plan alignment

本 change 对应 `docs/development-plan.md` 的 P1 全部内容：P1.1、P1.2。
