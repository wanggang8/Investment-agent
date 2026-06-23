# Design: P1 数据底座

## SQLite migration

- 迁移目录：`internal/infrastructure/persistence/sqlite/migration/`
- 迁移入口：`internal/infrastructure/persistence/sqlite/migrate.go`
- 首版使用标准库 `database/sql` 与 SQLite driver。
- `001_init.sql` 创建 `docs/data-model.md` 中 P1 必需核心表与索引。
- `002_seed_rules.sql` 写入默认 `rule_versions`（`v3.0`）与默认信源等级配置。
- 所有枚举字段使用 `CHECK` 约束。
- 历史表只追加，不设计物理删除接口。

## Repository 分层

- 接口放在 `internal/domain/repository/`。
- SQLite 实现放在 `internal/infrastructure/persistence/sqlite/`。
- Repository 接口只表达业务需要，不暴露 SQL 细节。
- 事务由 SQLite 实现层管理；失败时必须回滚。

## Repository 范围

| Repository | 能力 |
| --- | --- |
| `PortfolioRepository` | 账户快照、当前持仓、持仓快照写读 |
| `DecisionRepository` | 决策记录、证据引用、用户确认写读 |
| `IntelligenceRepository` | 情报摘要、RAG 文本块、多源验证写读 |
| `RuleRepository` | 规则版本、规则提案、守门人审计写读 |
| `AuditRepository` | 审计事件追加写入 |

## Testing

- migration 测试使用临时 SQLite 文件或内存库。
- 每个 Repository 至少覆盖：写入、读取、事务失败回滚。
- 验收命令与 `docs/development-plan.md` P1 保持一致。

## Constraints

- 不实现领域裁决逻辑。
- 不实现 HTTP handler。
- 不接 VecLite，只保留可重建所需的 SQLite 数据结构。
