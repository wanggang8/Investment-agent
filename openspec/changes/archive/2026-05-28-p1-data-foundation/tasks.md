# Tasks: p1-data-foundation

> 对齐 `docs/development-plan.md` P1：数据底座。

## 1. SQLite migration（P1.1）

- [x] 1.1 创建 `internal/infrastructure/persistence/sqlite/migration/001_init.sql`
- [x] 1.2 创建 `internal/infrastructure/persistence/sqlite/migration/002_seed_rules.sql`
- [x] 1.3 创建 `internal/infrastructure/persistence/sqlite/migrate.go`
- [x] 1.4 按 `docs/data-model.md` 创建 P1 核心表与字段
- [x] 1.5 按 `docs/data-model.md` 创建索引
- [x] 1.6 对枚举字段添加 CHECK 约束
- [x] 1.7 确保历史表追加式写入，不提供物理删除路径
- [x] 1.8 写入默认 `rule_versions`，版本号为 `v3.0`
- [x] 1.9 写入默认信源等级配置
- [x] 1.10 验收：`go test ./internal/infrastructure/persistence/sqlite/...`

## 2. Repository 层（P1.2）

- [x] 2.1 创建 `internal/domain/repository/portfolio_repo.go`
- [x] 2.2 创建 `internal/domain/repository/decision_repo.go`
- [x] 2.3 创建 `internal/domain/repository/intelligence_repo.go`
- [x] 2.4 创建 `internal/domain/repository/rule_repo.go`
- [x] 2.5 创建 `internal/domain/repository/audit_repo.go`
- [x] 2.6 创建 SQLite 实现：`portfolio_repo_impl.go`、`decision_repo_impl.go`、`intelligence_repo_impl.go`、`rule_repo_impl.go`、`audit_repo_impl.go`
- [x] 2.7 实现账户快照写读
- [x] 2.8 实现持仓当前态与持仓快照写读
- [x] 2.9 实现决策记录与证据引用写读
- [x] 2.10 实现用户确认写入
- [x] 2.11 实现情报摘要、RAG 文本块、多源验证写读
- [x] 2.12 实现规则版本、规则提案、守门人审计写读
- [x] 2.13 实现审计事件追加写入
- [x] 2.14 为每个 Repository 添加写入、读取、事务失败回滚测试
- [x] 2.15 验收：`go test ./internal/domain/repository/... ./internal/infrastructure/persistence/sqlite/...`

## 3. 归档前

- [x] 3.1 确认 `specs/data-model/spec.md` 的 delta 已合并或已被 `docs/data-model.md` 覆盖
- [x] 3.2 勾选 `docs/development-plan.md` P1 相关任务
- [x] 3.3 更新 `openspec/PROGRESS.md`：P1 标记为 `in_progress`

## Plan alignment

- P1.1 对应任务：1.1–1.10，共 10 项。
- P1.2 对应任务：2.1–2.15，共 15 项。
- 归档前治理任务：3.1–3.3，共 3 项。
