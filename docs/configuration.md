# 配置与启动说明

> 适用范围：本地 Go 后端、React 前端、SQLite、VecLite、DeepSeek 与数据源配置。  
> 阶段说明：P0–P40 已完成；公开 HTTP 数据桥接、A 股 ETF/基金证据 payload 解析、应用内通知、配置校验、本地备份恢复、公开数据 collector、风险预警、真实 LLM 质量评估、RAG / VecLite 检索质量加固、前端完整用户旅程 E2E、本地预检与恢复演练均为当前可用能力。见 `docs/README.md`、`docs/development-plan.md` §1.1 和 `openspec/PROGRESS.md`。
> 安全边界：文档只写环境变量名和示例占位值，不写真实密钥。

## 配置文件

默认读取本地 `configs/config.yaml`。首次本机启动前从示例复制：

```bash
cp configs/config.example.yaml configs/config.yaml
```

`configs/config.yaml` 会被 Git 忽略，可放本机端口、SQLite/VecLite 路径和可选 LLM key。`configs/config.example.yaml` 只作为提交到仓库的模板；如果本地 `configs/config.yaml` 不存在，程序会回退读取 example，方便 fresh checkout smoke。

可通过环境变量覆盖路径：

```bash
export INVESTMENT_AGENT_CONFIG=/path/to/config.yaml
```

### 配置项

| 字段 | 说明 | 默认示例 |
| --- | --- | --- |
| `server.host` / `server.port` | HTTP 监听地址 | `127.0.0.1:8080` |
| `sqlite.path` | SQLite 数据文件路径 | `./data/investment-agent.db` |
| `veclite.path` | VecLite 索引文件或目录路径 | `./data/veclite` |
| `deepseek.api_key` / `deepseek.base_url` / `deepseek.model` / `deepseek.timeout_seconds` | DeepSeek 或 OpenAI-compatible API Key、API 地址、模型名与超时秒数；常规运行可缺 key，真实 smoke 必须配置 key | 空字符串 / `https://api.deepseek.com` / `deepseek-chat` / `15` |
| `data_sources.enabled` | 启用的数据源名称列表 | `stub` |
| `data_sources.use_stub` | 是否启用本地 stub 数据源 | `true` |
| `data_sources.market_endpoint` | 最小只读行情 HTTP/JSON endpoint，系统会附加 `symbol` 查询参数 | 空字符串 |
| `data_sources.intelligence_endpoint` | 最小只读情报 HTTP/JSON endpoint，系统会附加 `symbol` 查询参数 | 空字符串 |
| `data_sources.public_evidence.enabled` / `sources` | P26 公告/监管证据 collector 手动任务开关与源列表；`false` 时 `public-evidence-refresh` 会拒绝执行 | `false` / `cninfo,szse,csrc` |
| `data_sources.public_evidence.*_base_url` / `cninfo_org_ids` | 巨潮资讯、深交所、证监会公开源基础 URL；`cninfo_org_ids` 可配置 `symbol -> orgId`，避免只依赖内置少量标的映射 | 公开站点 URL / 示例映射 |
| `data_sources.market_collectors.enabled` / `sources` | P27 中证指数与东方财富基金只读市场数据 collector 开关与源列表；默认关闭；启用后无需配置自备 `market_endpoint`，collector 优先于通用 endpoint；`use_stub=false` 时真实源失败不会自动回退到 stub | `false` / `csindex,eastmoney_fund` |
| `data_sources.market_collectors.*_base_url` | 中证指数、东方财富基金公开源基础 URL；测试可指向本地 httptest，生产默认公开站点 | 公开站点 URL |
| `daily_auto_run.enabled` | P31 本地每日自动运行开关；默认关闭，只有显式设为 `true` 才允许 server 内 scheduler 运行 | `false` |
| `daily_auto_run.run_time` / `timezone` | 每日本地计划时间与 IANA 时区；按低频本地调度解释，不用于高频抓取 | `08:30` / `Asia/Shanghai` |
| `daily_auto_run.scope` | 每日运行范围；P31 仅支持本地账户/组合当前持仓，缺持仓时记录缺前提状态 | `holdings` |
| `daily_auto_run.retry` / `timeout_seconds` / `max_symbols` | 自动运行有限重试、单次超时和最大标的数，用于控制本地资源和公开源访问频率 | `1` / `900` / `20` |
| `log.level` | 日志级别：`debug` / `info` / `warn` / `error` | `info` |

### 环境变量覆盖

| 变量 | 覆盖字段 | 说明 |
| --- | --- | --- |
| `INVESTMENT_AGENT_CONFIG` | 配置文件路径 | 优先于默认 `configs/config.yaml` 的显式本地 YAML 配置文件 |
| `INVESTMENT_AGENT_SERVER_PORT` | `server.port` | 本地 HTTP 服务端口 |
| `INVESTMENT_AGENT_SQLITE_PATH` | `sqlite.path` | SQLite 数据文件路径 |
| `INVESTMENT_AGENT_VECLITE_PATH` | `veclite.path` | VecLite 索引文件或目录路径 |
| `DEEPSEEK_API_KEY` | `deepseek.api_key` | DeepSeek API Key，只在本地环境变量中设置 |
| `INVESTMENT_AGENT_DATA_SOURCES` | `data_sources.enabled` | 逗号分隔的数据源名称，例如 `stub,official` |
| `INVESTMENT_AGENT_MARKET_DATA_ENDPOINT` | `data_sources.market_endpoint` | 最小只读行情 HTTP/JSON endpoint |
| `INVESTMENT_AGENT_INTELLIGENCE_ENDPOINT` | `data_sources.intelligence_endpoint` | 最小只读情报 HTTP/JSON endpoint |
| `INVESTMENT_AGENT_USE_STUB_DATA` | `data_sources.use_stub` | 本地 stub 开关，示例值 `true` / `false` |
| `INVESTMENT_AGENT_LOG_LEVEL` | `log.level` | 日志级别 |

## 数据源开关

### 阶段边界

| 阶段 | 数据源范围 | 当前状态 |
| --- | --- | --- |
| P12（已归档） | 最小只读 provider 适配层：HTTP/JSON endpoint 或 stub；不含完整财务、完整情绪源、券商 API | **已实现** |
| P19–P20（已交付） | 可配置公开 HTTP JSON endpoint、ETF/基金证据 payload 解析、中文信源等级映射、fixture/stub fallback；真实外部公开源 collector、历史补采和 smoke 验证属于 P25+ 独立范围 | **基础能力已实现** |
| P25（已归档） | 验证巨潮、交易所、证监会、基金业协会、东方财富基金、新浪财经、中证指数等公开源的访问方式、字段、频率、限制和后续实现范围 | **已完成验证** |
| P26（已完成） | 首批只读公告/监管证据 collector：巨潮资讯、深交所、证监会；AMAC 暂缓为二线背景候选；不接交易、登录、付费或授权行情源 | **已实现代码能力** |
| P27（已完成） | 首批只读市场数据 collector：东方财富基金净值、历史净值、资产配置和基金档案基础字段已通过真实公开源 smoke；中证指数基础信息 collector 已校准到当前公开 `index-basic-info` shape，样本/权重/估值文件扩展接口仍按候选 metadata 低频读取；默认关闭，可通过 market-refresh 写入 `market_snapshots.market_metrics_json`，真实源失败不伪造 stub 数据 | **已完成并归档；扩展文件能力按候选 metadata 读取** |
| P29（已完成） | 公开证据 collector 真实 smoke 修复：巨潮资讯当前接口参数、no_data/source_unavailable/parse_error 诊断、临时 SQLite 真实入库验收；深交所和证监会无公告或无法稳定按标的返回数据时按 no_data/source_unavailable 降级记录，不阻塞可用 A 级公告源 | **CNInfo 真实 smoke 已写入证据表；SZSE/CSRC 分类与降级边界已修复** |
| P39（已完成） | 前端完整用户旅程与全路径 E2E：空库 onboarding、账户/持仓初始化、每日纪律、决策详情、线下确认记录、复盘、规则治理、风险预警、source health、retrieval quality、console/a11y/窄屏 smoke | **使用本地临时 fixture 验收；不接券商、不自动交易、不外部推送、不自动应用规则、不收益承诺** |

`docs/requirements.md` 描述全量数据源清单；本文件与开发计划只描述当前可配置能力和下一轮候选增强。

### 默认与接入方式

`configs/config.example.yaml` 使用 `data_sources.use_stub=true` 的 stub 数据源作为模板默认值；真实本机运行建议复制到 `configs/config.yaml` 后切换到只读公开源或结构化公开 collector。接入方式：

1. 准备只读 HTTP JSON 数据源 endpoint。
2. 在配置中设置 `data_sources.use_stub=false`，并填写 endpoint，例如：

```yaml
data_sources:
  enabled:
    - "public-http"
  use_stub: false
  market_endpoint: "http://localhost:9090/market"
  intelligence_endpoint: "http://localhost:9090/intelligence"
```

凭证只放在本地环境变量或私有配置文件中，不写入仓库。以下能力**始终不在范围内**：券商交易 API、自动交易、主动荐股、收益承诺。

公开证据 collector 默认关闭。真实 smoke 或本地采集需要使用私有配置显式开启，例如只验证巨潮资讯时：

```yaml
data_sources:
  enabled: []
  use_stub: false
  public_evidence:
    enabled: true
    sources:
      - "cninfo"
    cninfo_base_url: "https://www.cninfo.com.cn"
    cninfo_org_ids:
      "510300": "9900000091"
      "000001": "gssz0000001"
    szse_base_url: "https://www.szse.cn"
    csrc_base_url: "https://www.csrc.gov.cn"
```

当前 `public-evidence-refresh` 默认窗口为最近 90 天，CLI 通过 `--symbol` 指定标的，也可通过 `--start-date YYYY-MM-DD --end-date YYYY-MM-DD` 指定显式采集窗口；只传一端时另一端仍按默认推导。P29 smoke 使用 `000001` 与显式日期窗口验证 CNInfo 公告真实入库；`510300` 等 ETF 标的是否有公告取决于真实源窗口，不得把空窗口伪造成接口失败。错误语义：`no_data` 表示源可达但窗口无匹配记录；全源均为 `no_data` 时记为成功空刷新并保留 source-specific degraded 审计；`source_unavailable` 表示 HTTP/接口不可用；`parse_error` 表示响应可达但 shape 不兼容。

| 数据源 | 当前状态 | 验收方式 |
| --- | --- | --- |
| SQLite | 启用 | 由 `sqlite.path` 指定数据文件；启动前配置校验会检查路径可用性 |
| VecLite | 可降级 | 由 `veclite.path` 指定索引位置；不可用时按 P6 A04 验证降级状态 |
| DeepSeek | 可为空；P37 支持 `model` 与 `timeout_seconds` | `DEEPSEEK_API_KEY` 或私有配置 key 为空时按 P6 A11 验证分析节点降级；配置后可执行真实调用和 `llm-smoke` |
| 市场数据 | stub 默认；可接自备公开 HTTP JSON 源；P27 提供东方财富基金真实净值 collector 与中证指数候选 collector，默认关闭，启用后通过 composite source 参与 `market-refresh`；`use_stub=false` 时真实源失败进入错误/降级，不自动写入 stub 行情 | `POST /api/v1/market/refresh` 或 `cmd/agent --task market-refresh` |
| 情报数据 | stub 默认；可接自备 A 股 ETF、基金公开 HTTP 情报 JSON；P26 提供巨潮、深交所、证监会只读 collector，P29 已验证 CNInfo 真实 smoke 可通过 `cmd/agent --task public-evidence-refresh --symbol 000001 --start-date YYYY-MM-DD --end-date YYYY-MM-DD` 写入本地证据表；CNInfo 可用 `data_sources.public_evidence.cninfo_org_ids` 配置标的 orgId 映射；真实采集必须显式设置 `data_sources.use_stub=false`、`data_sources.public_evidence.enabled=true` 并选择源 | 验证 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 与 `audit_events` |

## 启动 HTTP 服务

```bash
go run ./cmd/server
```

健康检查：

```bash
curl http://127.0.0.1:8080/api/v1/health
# {"status":"ok"}
```

## 启动前端

```bash
cd web
npm install
npm run dev
```

生产构建验收：

```bash
cd web
npm run build
```

## P39 本地浏览器 E2E

P39 使用 `scripts/e2e-smoke.sh` 运行本地浏览器级验收。脚本会创建 `tmp/e2e-smoke/` 下的临时 SQLite、VecLite 目录和配置文件，调用 `cmd/smoke-seed` 写入固定 ID 的临时 fixture，再启动本地 Go server 与 Vite，最后运行 Playwright。

```bash
bash scripts/e2e-smoke.sh
```

可选端口：

```bash
E2E_SERVER_PORT=18080 E2E_WEB_PORT=14173 bash scripts/e2e-smoke.sh
```

边界：

- fixture 使用空 key、stub/本地事实和 `example.invalid` 链接，不包含真实 API key、券商凭证或个人账户标识。
- E2E 只通过 HTTP API 与浏览器页面观察状态，不直接读取 SQLite、VecLite 或本地文件内容。
- 临时输出位于 `tmp/` 与 Playwright 输出目录，已由 `.gitignore` 忽略；脚本成功时会清理 `tmp/e2e-smoke/`，失败时保留日志和临时目录供排查。
- 测试覆盖线下确认记录、风险预警、规则提案、复盘追踪、source health 和 retrieval quality，但不会执行交易、不会外部推送，也不会自动应用规则。

## P40 本地预检与恢复 smoke

P40 使用 `cmd/agent` 的本地预检入口检查运行依赖、配置、SQLite/VecLite 路径、数据源和 LLM 配置状态。预检只读取本地状态，不执行任务；诊断 JSON 不包含 API key 原文。

```bash
go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json
```

备份恢复 smoke 使用临时路径创建 fixture、备份、恢复并启动本地 server，通过 HTTP API 验证恢复后的事实可读：

```bash
bash scripts/recovery-smoke.sh
```

可选端口：

```bash
RECOVERY_SMOKE_SERVER_PORT=18181 bash scripts/recovery-smoke.sh
```

边界：

- `--recovery-smoke` 默认拒绝覆盖已有 `sqlite.path`，恢复目标应使用临时路径或先手动移走旧库。
- `scripts/recovery-smoke.sh` 的输出位于 `tmp/recovery-smoke/`，脚本成功时会清理；失败时可查看临时 server log 后删除。
- 预检、恢复 smoke 和健康面板只展示本地诊断、修复提示和只读事实，不会执行交易、不会外部推送，也不会自动应用规则。

## P44 本地安装与诊断打包

P44 提供脚本入口，把本地预检、恢复 smoke 与 e2e smoke 做成统一可复现流程，并生成汇总文件，便于新环境快速验证。

```bash
bash scripts/local-install-diagnostics.sh --output-dir /tmp/investment-agent-diagnostics
```

可选参数：

```bash
--config PATH          配置文件路径（默认：configs/config.yaml，缺失时回退 configs/config.example.yaml）
--output-dir PATH      输出目录（默认：tmp/local-install-diagnostics）
--skip-recovery        跳过 scripts/recovery-smoke.sh
--skip-e2e             跳过 scripts/e2e-smoke.sh
--include-release-upgrade
                       显式纳入 P49 release/upgrade 检查
--target-version VALUE P49 release/upgrade 检查目标版本或 release label
```

脚本生成：

- `install-summary.json`：步骤名、状态（pass/failed/skipped）、退出码、命令与产物路径。
- `preflight.json`：`cmd/agent --preflight --diagnostics` 结果与诊断详情。
- `release-upgrade.json`：显式开启 `--include-release-upgrade` 时生成的 P49 升级检查报告。
- `recovery_smoke.log`、`e2e_smoke.log`（如执行）。

安全边界：

- 仅输出本地诊断事实，不执行交易、不连接券商、不外部推送。
- 摘要不承诺收益，不包含完整私钥；用于核对与交接的审计资料。

## P49 本地发布与升级检查

P49 提供只读发布/升级检查报告，用于升级前确认目标版本、备份提醒、迁移文件可读性和升级后 smoke 计划。该检查不会执行升级、不会运行迁移、不会创建备份、不会恢复或覆盖数据库。

```bash
go run ./cmd/agent --release-upgrade-check --target-version vNEXT --diagnostics ./tmp/release-upgrade.json
bash scripts/local-release-upgrade-check.sh --target-version vNEXT --skip-preflight
```

需要把 P49 检查纳入 P44 诊断打包时，必须显式开启：

```bash
bash scripts/local-install-diagnostics.sh --skip-recovery --skip-e2e --include-release-upgrade --target-version vNEXT
```

报告内容：

- `current_version`、`target_version`、总体 `status`。
- `version_check`、`config_validation`、`backup_reminder`、`migration_precheck` 和 `smoke_plan`。
- 升级前建议手动运行 `--preflight` 与 `--backup`。
- 升级后建议手动运行 `--preflight`、`recovery-smoke`、`e2e-smoke` 和安装诊断打包。

安全边界：

- 报告不包含 API key 原文、完整私有路径、raw SQL、完整 prompt、raw HTTP 或供应商原始响应。
- P49 只提供本地升级计划与检查，不自动下载、升级、迁移、修复、覆盖真实库、执行交易、外部推送、自动确认或自动应用规则。

## 启动 Agent 本地任务入口

```bash
go run ./cmd/agent --help
go run ./cmd/agent --validate-config
go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json
go run ./cmd/agent --task daily
go run ./cmd/agent --task market-refresh
go run ./cmd/agent --task evidence-index
go run ./cmd/agent --task public-evidence-refresh --symbol 510300 --start-date YYYY-MM-DD --end-date YYYY-MM-DD
go run ./cmd/agent --task llm-smoke --symbol 510300
go run ./cmd/agent --task review --period monthly
go run ./cmd/agent --task review --period quarterly
go run ./cmd/agent --backup ./data/backups
go run ./cmd/agent --restore ./data/backups/agent-YYYYMMDDTHHMMSSZ.db --restore-confirm
go run ./cmd/agent --recovery-smoke ./data/backups/agent-YYYYMMDDTHHMMSSZ.db
go run ./cmd/agent --release-upgrade-check --target-version vNEXT --diagnostics ./tmp/release-upgrade.json
```

`llm-smoke` 只读取本地配置中的 `deepseek.api_key`、`deepseek.base_url`、`deepseek.model` 和 `deepseek.timeout_seconds`，执行一次真实分析材料调用，并写入脱敏审计摘要。缺少 `deepseek.api_key` 时该任务会明确拒绝；常规本地任务仍允许在无 key 时按分析材料降级运行。

安全边界：`cmd/agent` 只触发本地分析、刷新、索引、复盘、配置诊断和本地文件备份恢复；不提供买入、卖出、撤单、改单或任何自动交易能力。任务执行会写入 `audit_events`，用于追踪输入摘要、执行状态和错误码。复盘可以生成规则提案，但不能让规则自动生效；规则变更仍需守门人审计和用户最终确认。

### 本地调度说明

当前调度只提供显式配置说明，不默认启用后台任务：

```bash
go run ./cmd/agent --schedule
```

P17 提供本地 launchd 与 cron 示例，详见 `docs/ops-local-scheduler.md`、`examples/scheduler/launchd/com.example.investment-agent.plist` 和 `examples/scheduler/cron/investment-agent.cron`。示例需要用户显式安装或编辑，不会执行交易，不会自动应用规则，也不会绕过用户确认或守门人审计。

如需在本机使用系统计划任务，请让计划任务调用上面的 `--task` 命令，并保留日志输出。不要把任何券商交易脚本挂到本项目任务后面。

## 数据备份、索引重建与恢复

### SQLite 备份

```bash
go run ./cmd/agent --validate-config
go run ./cmd/agent --backup ./data/backups
```

恢复时先停止 `cmd/server` 和本地任务；若 `sqlite.path` 已存在，先把旧库移动到安全位置，再显式确认恢复到该路径：

```bash
go run ./cmd/agent --restore ./data/backups/agent-YYYYMMDDTHHMMSSZ.db --restore-confirm
```

不带 `--restore-confirm` 时恢复会拒绝执行；目标 `sqlite.path` 已存在时也会拒绝覆盖，需要先手动移走旧库。

### VecLite / RAG 索引重建

```bash
go run ./cmd/agent --task evidence-index
# 或通过 HTTP API
curl -X POST http://127.0.0.1:8080/api/v1/evidence/rebuild-index
```

索引是可重建辅助数据；事实数据以 SQLite 中的 `intelligence_summary`、`rag_chunks` 和相关审计记录为准。

P13 后，本地索引继续使用 JSON 文件适配器。文件包含版本化 envelope，用于区分 `healthy`、`missing`、`corrupted`、`incompatible`、`degraded` 等状态。索引缺失、损坏或版本不兼容时，可通过上面的命令由 SQLite `rag_chunks` 重新生成。重建响应会返回 `indexed_count`、`skipped_count`、`last_rebuild_at` 和 `index_health`，便于前端展示索引状态与降级原因。

真实 VecLite API 仍是后续可替换实现，本阶段不强制接入外部 VecLite 服务。

## 常见故障处理

| 故障 | 处理方式 |
| --- | --- |
| 数据源不可用 | 确认 `data_sources.enabled` 与 `INVESTMENT_AGENT_USE_STUB_DATA`；必要时启用 stub，并检查 `audit_events.error_code`。 |
| VecLite 索引损坏 | 备份 SQLite 后执行 `go run ./cmd/agent --task evidence-index` 或调用重建 API。 |
| DeepSeek 缺配置 | 常规工作流可设置 `DEEPSEEK_API_KEY` 或私有配置 `deepseek.api_key`；未配置时系统只降级分析材料，最终裁决仍由规则生成。`llm-smoke` 需要私有配置中的 key。 |
| SQLite 写入失败 | 检查 `sqlite.path`、目录权限、磁盘空间；修复后重新执行任务，已提交事实不应被覆盖。 |
| 本地调度失败 | 查看 launchd/cron 日志和 `audit_events.error_code`；按 `docs/ops-local-scheduler.md` 停用、修复后再启用。 |

## 本地验收命令

```bash
go test ./...
go run ./cmd/agent --help
go test ./internal/application/workflow/... ./internal/application/handler/...
cd web && npm run build && npm test
```
