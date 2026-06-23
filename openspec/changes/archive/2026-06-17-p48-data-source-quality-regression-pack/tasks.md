## 1. OpenSpec 与范围

- [x] 1.1 确认 P47 已归档，P48 为当前活跃 change。
- [x] 1.2 确认 P48 聚焦数据源质量回归包：source health、freshness、parse_error/no_data/source_unavailable/stale 分类、fixture/current 回归和脱敏摘要。
- [x] 1.3 确认 P48 不新增券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、收益承诺、登录/付费/授权/Level2/高频源。

## 2. 后端 DTO 与服务

- [x] 2.1 新增 `internal/application/dto/data_source_quality.go`，定义 regression response、case item 和 status 字段。
- [x] 2.2 新增 `internal/application/service/source_health.go`，将 P34 source health 从 market snapshot 提取为可复用只读函数。
- [x] 2.3 新增 `internal/application/service/data_source_quality.go`，实现 `fixture` 与 `current` 回归、状态汇总和 missing categories。
- [x] 2.4 实现诊断脱敏：完整 key、私有路径、原始 SQL、prompt、raw HTTP、HTTP status line、private key 和供应商原始响应不得进入返回值或审计摘要。
- [x] 2.5 新增 `internal/application/service/data_source_quality_test.go`，覆盖 fixture 全通过、current 无数据降级、parse_error/no_data/source_unavailable/stale 分类、未知分类失败和脱敏。

## 3. 后端 Handler 与 API

- [x] 3.1 新增 `internal/application/handler/data_source_quality_handler.go`。
- [x] 3.2 在 `internal/application/handler/app.go` 初始化服务并注册 `GET /api/v1/data-source-quality/regression`。
- [x] 3.3 更新 `internal/application/handler/market_handler.go` 使用 service 层 source health 提取逻辑，避免 handler 私有实现漂移。
- [x] 3.4 新增 `internal/application/handler/data_source_quality_handler_test.go`，覆盖 fixture/current、非法 mode、空库降级和响应脱敏。

## 4. CLI 回归入口

- [x] 4.1 在 `cmd/agent/main.go` 增加 `data-source-quality-regression` task，默认 `--source fixture`，支持 `--source current`。
- [x] 4.2 CLI 输出紧凑摘要，并通过既有 `appendTaskAudit` 写入脱敏 audit output ref。
- [x] 4.3 更新 `cmd/agent/main_test.go`，覆盖 help、fixture task、current task、unsupported source 和不写账户/确认/交易表。

## 5. 文档与契约

- [x] 5.1 更新 `docs/api.md`，新增 P48 regression API。
- [x] 5.2 更新 `docs/data-model.md`，说明 P48 只读复用 `market_snapshots.market_metrics_json`，CLI 仅写脱敏任务审计。
- [x] 5.3 更新 `docs/development-plan.md`、`openspec/project.md`、`openspec/PROGRESS.md`、`docs/GOVERNANCE.md` 和 `AGENTS.md` 当前阶段状态。
- [x] 5.4 在 OpenSpec delta 中记录 P48 行为要求。

## 6. 执行前复审

- [x] 6.1 计划完成后执行只读子 agent 复审，确认无 Critical / Important。
- [x] 6.2 复审通过后再执行实现任务。

## 7. 验收

- [x] 7.1 运行 `go test ./...`。
- [x] 7.2 运行 `npm --prefix web test -- --run`。
- [x] 7.3 运行 `npm --prefix web run build`。
- [x] 7.4 运行 `bash scripts/e2e-smoke.sh`。
- [x] 7.5 运行 `openspec validate p48-data-source-quality-regression-pack --strict`。
- [x] 7.6 运行 `openspec validate --all --strict`。
- [x] 7.7 运行 `git diff --check`。
- [x] 7.8 运行安全扫描：`rg -n 'sk-[A-Za-z0-9][A-Za-z0-9_-]{8,}|BEGIN (RSA|OPENSSH|PRIVATE) KEY|/Users/[^[:space:]，；。、]+|(?i:select[[:space:]]+\*[[:space:]]+from)|(?i:raw[[:space:]]+http)|(?i:prompt[[:space:]]*:)|完整[[:space:]]*prompt|HTTP/[0-9.]+[[:space:]]+[0-9]{3}|券商接口|自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|自动修复|收益承诺|登录源|付费源|授权源|Level2|高频源' internal/application/dto/data_source_quality.go internal/application/service/data_source_quality.go internal/application/service/source_health.go internal/application/handler/data_source_quality_handler.go cmd/agent/main.go docs/api.md docs/data-model.md`，人工复核命中项，确认不存在未脱敏敏感内容或高风险操作入口；允许安全边界说明文本命中。

## 8. 归档前复审

- [x] 8.1 执行完成后再次只读子 agent 复审，确认无 Critical / Important。
- [x] 8.2 复审通过后执行 archive，并将 P48 归档。
