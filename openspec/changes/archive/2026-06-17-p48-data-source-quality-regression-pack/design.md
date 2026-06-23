# P48 Design

## Overview

P48 提供“数据源质量回归包”：一个后端服务同时支撑 API 与 CLI，用固定 fixture 和当前本地 source health 验证数据源质量分类、freshness 汇总和脱敏边界。它不是新的数据采集阶段，而是对 P34/P43 已有质量语义的回归保护。

## Regression Modes

### `fixture`

默认模式，不访问公网、不读取用户私有事实。服务内置最小样本：

- `fresh`: A 级公开指数健康样本。
- `no_data`: 可达但窗口无记录。
- `source_unavailable`: 源不可达或请求失败。
- `parse_error`: 响应结构不兼容。
- `stale`: 数据日落后于期望。
- `redaction`: 包含 key、私有路径、SQL、prompt、raw HTTP 和 private key 片段的诊断文本，必须输出脱敏摘要。

### `current`

只读读取最新市场快照或指定 `symbol` 的最新快照，解析其中 `p34_source_health` 并按同一规则评估。该模式不触发 collector、不刷新快照、不创建通知；CLI 仅额外写入本地任务审计摘要。

## API Shape

`GET /api/v1/data-source-quality/regression`

查询参数：

- `mode`: 可选，`fixture` 或 `current`，默认 `fixture`。
- `symbol`: 可选，仅 `current` 模式用于读取指定标的的最新市场快照。

响应 `data`：

- `mode`
- `status`: `passed / degraded / failed`
- `generated_at`
- `summary`
- `cases`: `DataSourceQualityCase[]`
- `missing_categories`
- `safety_note`

`DataSourceQualityCase`：

- `case_id`
- `source_name`
- `source_level`
- `source_type`
- `data_category`
- `expected_freshness`
- `actual_freshness`
- `status`: `passed / degraded / failed`
- `data_date`
- `failure_category`
- `affected_symbols`
- `diagnostic_preview`

## Backend

Add:

- `internal/application/dto/data_source_quality.go`
- `internal/application/service/source_health.go`
- `internal/application/service/data_source_quality.go`
- `internal/application/handler/data_source_quality_handler.go`

Service rules:

1. Treat `fresh`, `stubbed`, `no_data`, `source_unavailable`, `parse_error`, `stale`, `missing`, and `unknown` as recognized freshness categories.
2. Mark fixture cases as `passed` only when actual freshness equals expected freshness and diagnostic previews are sanitized.
3. Mark current mode as `degraded` when no source health exists or when any case is non-fresh but recognized.
4. Mark current mode as `failed` only for unsafe diagnostics, unrecognized categories, or malformed source health that cannot be mapped safely.
5. Never include raw `market_metrics_json`, raw HTTP body, SQL, keys, private paths, prompts or supplier raw response in the API/CLI response.

## CLI

Add task:

```bash
go run ./cmd/agent --task data-source-quality-regression --source fixture --symbol 000300
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300
```

The task prints a compact summary and writes a sanitized `audit_events` record through existing `appendTaskAudit`. The audit output ref should be compact, for example:

```text
data_source_quality:mode=fixture:status=passed:cases=6:degraded=0:failed=0:no_auto_trading
```

## Frontend

No new page is required. P48 may add a small read-only link or smoke assertion later, but the primary product surface is API/CLI because P43 already owns the visual quality dashboard.

## Guardrails

- No new collector, schema, scheduler, external push, rule mutation, confirmation mutation, account mutation, broker integration, or trading action.
- Default fixture mode must be deterministic and offline.
- Current mode must never trigger refresh or network access.
- All diagnostics are previews only and must be sanitized before return or audit.
