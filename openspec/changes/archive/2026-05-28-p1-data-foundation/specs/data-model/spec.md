# Delta for Data Model（合并目标：`docs/data-model.md`）

## ADDED Requirements

### Requirement: SQLite migration creates the P1 fact baseline

系统 SHALL 提供 SQLite migration，用于创建 `docs/data-model.md` 中 P1 数据底座需要的核心表、索引与约束。

#### Scenario: Create empty database

- **GIVEN** 一个不存在或为空的 SQLite 数据库
- **WHEN** 执行 P1 migration
- **THEN** 数据库包含以下核心表：`portfolio_snapshots`、`positions`、`position_snapshots`、`operation_confirmations`、`position_transactions`、`market_snapshots`、`rule_versions`、`decision_records`、`intelligence_items`、`intelligence_summary`、`rag_chunks`、`evidence_refs`、`source_verifications`、`capability_configs`、`user_settings`、`audit_events`、`error_cases`、`rule_proposals`、`gatekeeper_audits`
- **AND** 索引与枚举 CHECK 约束按 `docs/data-model.md` 定义创建

#### Scenario: Repeat migration

- **GIVEN** 一个已完成 P1 migration 的 SQLite 数据库
- **WHEN** 再次执行 migration
- **THEN** 已有表结构与种子数据不被破坏

### Requirement: Default rules and source levels are seeded

系统 SHALL 在 P1 初始化时写入默认规则版本与默认信源等级配置。

#### Scenario: Seed default rule version

- **WHEN** 执行 P1 migration
- **THEN** `rule_versions` 中存在 active 规则版本 `v3.0`
- **AND** 同一时间只有一个 active 版本

#### Scenario: Seed default source levels

- **WHEN** 执行 P1 migration
- **THEN** 默认信源等级配置可被后续证据验证流程读取

### Requirement: Repository layer persists and reads core facts

系统 SHALL 提供 Repository 接口与 SQLite 实现，用于读写账户、持仓、决策、证据、确认、情报、规则和审计数据。

#### Scenario: Write and read repository data

- **WHEN** Repository 写入一条核心事实数据
- **THEN** 可通过对应读取方法取回同等业务含义的数据

#### Scenario: Roll back failed transaction

- **GIVEN** 一个 Repository 事务中出现错误
- **WHEN** 事务返回失败
- **THEN** 事务内已写入的数据不应部分持久化

## MODIFIED Requirements

（无）

## REMOVED Requirements

（无）
