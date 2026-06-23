# 本地调度与运维说明

> P17 提供外部本机计划任务示例；P31 新增 server 内置每日自动运行闭环。两者默认都不自动运行，需要用户显式安装、编辑或启用。不会执行交易，不会自动应用规则。

## 安全边界

- `cmd/agent --task` 支持 `daily`、`market-refresh`、`evidence-index`、`public-evidence-refresh`、`p34-expanded-refresh`、`llm-smoke`、`review` 等本地 workflow task；`public-evidence-refresh` 和 `p34-expanded-refresh` 可用 `--symbol`、`--start-date YYYY-MM-DD`、`--end-date YYYY-MM-DD` 指定标的和日期窗口。
- `llm-smoke` 只用于显式验证本地私有配置中的真实 LLM endpoint、model、timeout 和 key 是否可生成分析材料；它不读取账户、不写持仓、不生成裁决、不创建交易确认或交易流水。
- `cmd/agent` 的运维 flag 支持 `--validate-config`、`--preflight`、`--diagnostics`、`--backup`、`--restore`、`--restore-confirm` 与 `--recovery-smoke`，只用于本地配置诊断、依赖预检和 SQLite 文件备份恢复。
- 调度示例只调用本地任务入口，不连接券商，不创建订单，不修改交易状态。
- 账户变化仍只能来自用户记录的线下动作。
- 规则变化仍需守门人审计和用户最终确认。
- 任务结果应写入 `audit_events`，用于追踪输入摘要、状态和错误码。

## 前置条件

1. 准备本地配置文件，并用环境变量指向它：

```bash
export INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml
```

2. 确认帮助输出：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
go run ./cmd/agent --help
```

3. 校验配置并手动验证单个任务：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --validate-config
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --preflight --diagnostics /ABSOLUTE/PATH/TO/Investment-agent/tmp/preflight.json
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --task market-refresh
```

## P40 本地预检

本地预检用于部署前确认 Go、Node、npm、Playwright browser、SQLite path、VecLite path、配置文件、数据源和 DeepSeek 配置状态：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --preflight --diagnostics /ABSOLUTE/PATH/TO/Investment-agent/tmp/preflight.json
```

预检输出状态枚举为 `pass`、`warning`、`failed`、`skipped`。诊断文件只写本地 JSON，不包含 key 原文；失败项应先按修复提示处理，再启动 server 或本地任务。预检只读取本地状态，不会执行交易，不外部推送，也不会自动应用规则。

## P44 本地安装与诊断打包

P44 在本地运维场景里补充一体化诊断脚本，便于新环境快速跑通配置、恢复和（可选）E2E 验收，并产出可追溯摘要：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
bash scripts/local-install-diagnostics.sh --output-dir /ABSOLUTE/PATH/TO/Investment-agent/tmp/local-install-diagnostics
```

可选参数：

```bash
--skip-recovery
--skip-e2e
--include-release-upgrade
--target-version VALUE
```

脚本默认输出到给定目录下：

- `install-summary.json`：步骤序列、状态（`pass`/`failed`/`skipped`）、退出码、命令和产物。
- `preflight.json`、`recovery_smoke.log`、`e2e_smoke.log`：用于复核与排障。
- `release-upgrade.json`：显式启用 `--include-release-upgrade` 时生成的 P49 本地发布/升级检查报告。

边界提醒：

- 本流程不连接券商，不下单，不创建交易执行语义，不外部推送。
- 摘要用于本地核对，不承诺收益，不替代正式回归；
- 如出现 `failed`，建议按日志修复后重跑；`skipped` 步骤用于显式短路场景。

## P49 本地发布与升级检查

升级前先运行只读发布/升级检查，确认目标版本、备份提醒、迁移文件可读性和升级后 smoke 计划：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --release-upgrade-check --target-version vNEXT --diagnostics /ABSOLUTE/PATH/TO/Investment-agent/tmp/release-upgrade.json
bash scripts/local-release-upgrade-check.sh --target-version vNEXT --skip-preflight
```

升级检查只读本地配置和迁移文件状态，不会执行升级、不会运行迁移、不会创建备份、不会恢复或覆盖数据库。若报告提示 `backup_reminder:warning`，应在升级前手动执行：

```bash
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --backup /ABSOLUTE/PATH/TO/Investment-agent/data/backups
```

升级后建议手动运行：

```bash
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --preflight --diagnostics /ABSOLUTE/PATH/TO/Investment-agent/tmp/preflight-after-upgrade.json
bash scripts/recovery-smoke.sh
bash scripts/e2e-smoke.sh
bash scripts/local-install-diagnostics.sh --skip-e2e --include-release-upgrade --target-version vNEXT
```

P49 不自动下载、升级、迁移、修复、覆盖真实库、执行交易、外部推送、代为确认或自动应用规则；报告不得包含 API key 原文、完整私有路径、raw SQL、完整 prompt、raw HTTP 或供应商原始响应。

## P31 内置每日自动运行

`cmd/server` 会读取 `daily_auto_run` 配置并创建本地每日自动运行入口，但默认配置为关闭。只有显式设置 `daily_auto_run.enabled: true` 后，server 生命周期内才会按 `daily_auto_run.run_time` 和 `daily_auto_run.timezone` 计算下一次触发并低频执行每日自动运行；启动 server 本身不会立即运行。

内置自动运行只使用本地 SQLite、公开/已配置只读数据源和应用内通知：

- scope 当前以本地持仓为准；缺持仓时记录 `missing_prerequisites`，不生成正式交易建议。
- 结果写入 `daily_auto_run_states`、`audit_events` 和应用内 `notifications`。
- 同一 local date、scope、symbol set 和 task version 使用幂等 key；执行副作用前先写入 `running` 状态，后续重复触发会复用已有状态，避免重复通知刷屏和重复生成每日结果。
- 失败、重试、超时都会写入审计诊断；诊断包含 `safety=no_auto_trading`。
- 不连接券商，不创建订单，不写入交易执行状态，不发送邮件、短信、Webhook 或系统 Push。

外部 launchd/cron 示例仍适用于手动本地任务入口；内置每日自动运行适用于打开本地 server 时的应用内状态展示和低频自动刷新。二者不要同时配置为高频重复执行同一任务。

## macOS launchd 示例

示例文件：`examples/scheduler/launchd/com.example.investment-agent.plist`。

安装前必须把 `/ABSOLUTE/PATH/TO/Investment-agent` 改为本机仓库路径，并创建日志目录。默认仅作为模板保存，不会自动安装。示例中的任务名是安全本地任务示例，路径与配置为占位符。

```bash
mkdir -p /ABSOLUTE/PATH/TO/Investment-agent/data/logs
cp examples/scheduler/launchd/com.example.investment-agent.plist ~/Library/LaunchAgents/com.example.investment-agent.plist
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.example.investment-agent.plist
launchctl kickstart -k gui/$(id -u)/com.example.investment-agent
```

停用与移除：

```bash
launchctl bootout gui/$(id -u)/com.example.investment-agent
rm ~/Library/LaunchAgents/com.example.investment-agent.plist
```

## cron 示例

示例文件：`examples/scheduler/cron/investment-agent.cron`。

安装前必须把 `/ABSOLUTE/PATH/TO/Investment-agent` 改为本机仓库路径，并确认日志目录存在。示例中的任务名是安全本地任务示例，路径与配置为占位符。

```bash
crontab examples/scheduler/cron/investment-agent.cron
crontab -l
```

停用本项目任务时，先备份当前 cron，再编辑删除本项目行，避免影响其他计划任务：

```bash
crontab -l > /ABSOLUTE/PATH/TO/Investment-agent/data/backups/crontab.before-investment-agent.txt
crontab -e
```

仅当确认当前用户所有 cron 项都可以移除时，才使用全量移除命令：

```bash
crontab -r
```

## SQLite 备份与恢复

备份前暂停 HTTP 服务和本地任务，并使用本地 CLI 创建备份：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --backup /ABSOLUTE/PATH/TO/Investment-agent/data/backups
```

恢复时先停止服务；若 `sqlite.path` 已存在，先把旧库移动到安全位置，再显式确认恢复到该路径：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --restore /ABSOLUTE/PATH/TO/Investment-agent/data/backups/agent-YYYYMMDDTHHMMSSZ.db --restore-confirm
```

不带 `--restore-confirm` 时恢复会拒绝执行；目标 `sqlite.path` 已存在时也会拒绝覆盖，需要先手动移走旧库。

恢复 smoke 会使用临时路径验证备份可恢复，并通过本地 HTTP API 读取恢复后的代表性事实：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
bash scripts/recovery-smoke.sh
```

如端口冲突，可显式指定：

```bash
RECOVERY_SMOKE_SERVER_PORT=18181 bash scripts/recovery-smoke.sh
```

`scripts/recovery-smoke.sh` 会创建 `tmp/recovery-smoke/`、写入临时配置、生成 fixture、备份、恢复到另一临时 DB、启动 server 并调用 `/api/v1/health` 与决策详情 API。脚本正常退出会清理临时目录；失败时查看 `tmp/recovery-smoke/server.log` 后手动删除临时目录。不要改用真实私有数据库复现 smoke。

## VecLite 索引恢复

索引缺失、过期或版本不兼容时，可用本地任务重建辅助索引：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
INVESTMENT_AGENT_CONFIG=/ABSOLUTE/PATH/TO/Investment-agent/configs/config.example.yaml go run ./cmd/agent --task evidence-index
```

索引是可重建辅助数据；事实数据以 SQLite 中的摘要、文本块和审计记录为准。

## P39 本地 E2E fixture

前端完整用户旅程验收通过 `scripts/e2e-smoke.sh` 运行。该脚本只创建临时本地资源：

- `tmp/e2e-smoke/config.yaml`：临时配置，`deepseek.api_key` 为空，数据源使用 stub/fixture。
- `tmp/e2e-smoke/investment-agent-smoke.db`：临时 SQLite，由 `cmd/smoke-seed` 写入代表性事实。
- `tmp/e2e-smoke/veclite/`：临时索引目录。
- `tmp/e2e-smoke/server.log` 与 `web.log`：失败排查用日志。

运行命令：

```bash
cd /ABSOLUTE/PATH/TO/Investment-agent
bash scripts/e2e-smoke.sh
```

如端口冲突，可显式指定：

```bash
E2E_SERVER_PORT=18081 E2E_WEB_PORT=14174 bash scripts/e2e-smoke.sh
```

脚本会在退出时停止本地 server/Vite；成功时清理 `tmp/e2e-smoke/`，失败时保留日志和临时目录供排查。Playwright 失败截图、trace 等输出落在 `tmp/playwright-output/`，该目录已被 `.gitignore` 忽略。

P39 fixture 覆盖空库/缺前提提示、账户校准、每日纪律报告、决策详情、线下确认记录、复盘摘要、规则提案、规则效果追踪、风险预警、source health、retrieval quality、console error 捕获和窄屏 smoke。它不会读取私有持久库，不包含真实密钥，不访问券商接口，不会执行交易，不外部推送，也不会自动应用规则。

## 常见故障

| 故障 | 处理 |
| --- | --- |
| 数据源不可用 | 检查配置中的数据源开关；必要时启用 stub；查看 `audit_events.error_code`。 |
| DeepSeek 未配置 | 设置本地环境变量；未配置时分析材料降级，最终裁决仍由规则生成。 |
| VecLite 索引异常 | 先备份 SQLite，再执行 `evidence-index`。 |
| SQLite 写入失败 | 检查数据目录、权限和磁盘空间；修复后重新执行任务。 |
| 调度任务失败 | 查看本地日志和 `audit_events`；不要把失败任务改成确认绕过或交易脚本。 |
| P39 E2E 端口冲突 | 设置 `E2E_SERVER_PORT` / `E2E_WEB_PORT` 后重跑；确认旧 server 或 Vite 进程已停止。 |
| P39 E2E 浏览器失败 | 查看 `tmp/playwright-output/` 中的截图和 trace；不要改用真实私有数据库复现。 |
| P40 预检失败 | 查看 `--diagnostics` JSON，先修复 failed 项；warning 项可按本地场景判断是否继续。 |
| P40 恢复 smoke 失败 | 查看 `tmp/recovery-smoke/server.log`；确认端口、临时目录权限和备份文件路径。 |

## 验收命令

P17 相关验收命令见 `docs/configuration.md` 的“本地验收命令”，以及 `openspec/changes/p17-local-scheduler-and-ops-docs/tasks.md`。
