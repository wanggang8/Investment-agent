# P48: 数据源质量回归包

## Summary

新增本地数据源质量回归包，用固定 fixture 和当前已保存的 source health 验证 `fresh`、`no_data`、`source_unavailable`、`parse_error`、`stale` 等质量分类、freshness 汇总和脱敏摘要。P48 不新增数据源、调度器、自动修复或交易能力，只提供可重复运行的本地质量检查入口。

## Why

P34 已建立 source health/freshness，P43 已提供数据质量只读面板，但目前缺少一个稳定、可重复的回归入口来证明数据源分类、降级语义和脱敏边界没有被后续改动破坏。真实公开源偶发波动时，用户也需要区分“源变化导致的降级”与“本地解析/展示逻辑退化”。

## What Changes

- 新增数据源质量回归 DTO 与服务：
  - 默认 `fixture` 模式使用确定性本地样本覆盖 `fresh`、`no_data`、`source_unavailable`、`parse_error`、`stale` 和敏感摘要脱敏。
  - `current` 模式只读评估最新市场快照中的 P34 source health，不触发刷新。
- 新增只读 API：`GET /api/v1/data-source-quality/regression?mode=fixture|current&symbol=...`。
  - 返回 case 列表、状态汇总、缺口、脱敏诊断和安全文案。
  - 不写 SQLite、不触发 collector、不创建通知。
- 新增本地 CLI 任务：`go run ./cmd/agent --task data-source-quality-regression --source fixture|current --symbol 000300`。
  - CLI 可写入一条本地 `audit_events` 摘要，摘要必须脱敏且只说明回归结果。
  - 默认 `fixture`，不访问公网；`current` 只读评估现有本地快照。
- 更新文档、OpenSpec delta 与 smoke/单测，保证回归入口可验证且不出现高风险操作入口。

## Scope

- 复用现有 `market_snapshots.market_metrics_json.metadata.p34_source_health` 与 `SourceHealthItem` DTO。
- 可新增只读应用服务、handler、CLI task、单测和文档。
- 可将 source health 提取逻辑从 handler 移到 service 层复用。

## Out of Scope

- 新增券商接口、自动交易、一键交易、代下单。
- 新增外部推送、自动确认、自动规则应用、自动修复承诺。
- 新增登录源、付费源、授权源、Level2、高频源或浏览器抓取。
- 新增数据库 schema、后台调度器、真实公网默认访问或自动刷新。
- 收益承诺、确定性涨跌预测、把回归结果升级为交易建议。

## Validation

- `go test ./...`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `bash scripts/e2e-smoke.sh`
- P48 安全扫描（见 `tasks.md` 7.8）
- `openspec validate p48-data-source-quality-regression-pack --strict`
- `openspec validate --all --strict`
- `git diff --check`
