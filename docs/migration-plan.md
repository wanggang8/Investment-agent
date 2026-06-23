# Migration 与 Seed 说明

> 适用范围：SQLite schema 初始化、默认规则版本与信源等级 seed。  
> 执行边界：migration 只初始化和演进本地事实库，不连接外部交易系统。

## 1. 文件位置

```text
internal/infrastructure/persistence/sqlite/migration/001_init.sql
internal/infrastructure/persistence/sqlite/migration/002_seed_rules.sql
internal/infrastructure/persistence/sqlite/migrate.go
```

## 2. 执行方式

HTTP 服务启动时会打开 `sqlite.path` 指定的数据库文件，并执行内嵌 migration：

```bash
go run ./cmd/server
```

执行链路：

1. `configs/config.yaml`、`INVESTMENT_AGENT_CONFIG` 或 fresh-checkout fallback `configs/config.example.yaml` 提供配置。
2. `INVESTMENT_AGENT_SQLITE_PATH` 可覆盖 `sqlite.path`。
3. `cmd/server/main.go` 调用 `sqlite.Open(cfg.SQLite.Path)`。
4. `sqlite.Migrate(ctx, db)` 按文件名顺序执行 `migration/*.sql`。

## 3. SQLite 数据文件路径

默认路径：

```text
./data/investment-agent.db
```

可通过环境变量覆盖：

```bash
export INVESTMENT_AGENT_SQLITE_PATH=./data/investment-agent.db
```

## 4. Migration 行为

- `001_init.sql` 创建核心表、索引与枚举 CHECK 约束。
- `002_seed_rules.sql` 写入默认规则版本和信源等级配置。
- migration 文件通过 Go embed 打入二进制，按文件名排序执行。
- SQL 使用 `IF NOT EXISTS` 与 `INSERT OR IGNORE`，重复启动保持幂等。

## 5. Seed 数据说明

`002_seed_rules.sql` 默认写入：

- active 规则版本：`v3.0`。
- 信源等级配置：`S`、`A`、`B`、`C`。
- C 级信源 `formal_allowed=0`，只能作为背景材料。

这些 seed 数据仅用于初始化本地规则与证据边界，不代表规则不可演进；规则提案仍必须经过守门人审计与用户最终确认。

## 6. 验证命令

```bash
go test ./internal/infrastructure/persistence/sqlite/...
go test ./...
```

## 7. 回滚与恢复

当前 migration 采用追加式事实表和幂等初始化，不提供自动降级脚本。开发期需要恢复时，可停止服务后备份或移除本地 SQLite 数据文件，再重新启动服务执行 migration。

```bash
cp ./data/investment-agent.db ./data/investment-agent.backup.db
# 如需重建开发库，再删除本地 db 文件并重新启动服务
```

请勿在生产或真实数据环境中删除数据库文件；正式迁移策略应在独立 change 中定义。

## 8. P9 本地任务恢复说明

`cmd/agent` 会在任务执行前打开 SQLite 并执行幂等 migration。若任务因 SQLite 写入失败中断：

1. 停止 `cmd/server` 和正在执行的本地任务。
2. 备份当前 SQLite 文件。
3. 检查 `sqlite.path`、目录权限和磁盘空间。
4. 修复后重新执行 `go run ./cmd/agent --task <task>`。

VecLite 索引损坏时，不需要删除 SQLite 事实库，可执行：

```bash
go run ./cmd/agent --task evidence-index
```

该操作只重建本地检索辅助数据，不触发交易。
