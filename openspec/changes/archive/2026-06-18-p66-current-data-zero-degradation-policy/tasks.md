# Tasks: P66 Current Data Zero-Degradation Policy

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 P52 gate matrix、P63/P65 acceptance records、P48 data-source-quality spec、现有 data-source-quality service/API/CLI/UI。
- [x] 1.3 创建 `p66-current-data-zero-degradation-policy` OpenSpec change。
- [x] 1.4 写明 P66 只新增当前数据质量 policy verdict 和 release gate，不新增外部源、真实 provider 调用、SQLite schema、Eino workflow、券商接口、自动交易或自动修复。
- [x] 1.5 更新当前进度文档，标记 P66 active。
- [x] 1.6 运行 `openspec validate p66-current-data-zero-degradation-policy --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 1.7 子 agent 方案复审无 Critical / Important 后执行。

## 2. 后端 Policy DTO 与服务

- [x] 2.1 扩展 `internal/application/dto/data_source_quality.go`，新增 `DataSourceQualityPolicy` 字段和结构。
- [x] 2.2 在 `internal/application/service/data_source_quality.go` 中实现 current policy 分类：`passed`、`waiver_required`、`blocked`。
- [x] 2.3 将无 source-health facts、明确 `freshness=missing`、failed/unrecognized freshness、unrecognized failure category、core category degradation 分类为 `blocked`。
- [x] 2.4 将仅可识别的 optional category degradation 分类为 `waiver_required`；即使该 optional category 出现在现有 `missing_categories` 降级列表里，也不得误判为 missing source-health metadata。
- [x] 2.5 保持 fixture mode policy 为 `passed`，用于证明分类器和脱敏逻辑稳定。
- [x] 2.6 扩展 `DataSourceQualityAuditOutputRef`，包含 `policy=<verdict>` 和 `gate=<release_gate>`，继续只输出脱敏摘要。
- [x] 2.7 增加/更新 service tests 覆盖 passed、waiver_required、blocked、无 source-health facts、`freshness=missing`、unknown freshness、unknown failure category、optional degraded + `missing_categories` 非空仍 waiver、redaction。

## 3. CLI/API 与严格门禁

- [x] 3.1 更新 handler/API tests，确认 `policy` 出现在 fixture/current response 且不写 audit。
- [x] 3.2 为 CLI 增加当前数据严格门禁入口，支持在 policy `blocked` 时返回非零 exit code。
- [x] 3.3 严格门禁仍只能读取本地当前库，不刷新数据、不调用 provider、不修复、不改规则、不触发交易。
- [x] 3.4 更新 CLI tests：fixture 仍成功；current blocked/waiver 输出 policy 和 gate；严格门禁 blocked 时失败并写脱敏审计。

## 4. 前端数据质量展示

- [x] 4.1 更新前端 data-source-quality 类型和 service 调用，读取 regression policy。
- [x] 4.2 更新 `/data-quality` loader，将 current regression policy 纳入页面模型。
- [x] 4.3 更新 `dataQualityExperienceModel`，把 `blocked` 显示为 danger，`waiver_required` 显示为 warning，`passed` 显示为 success。
- [x] 4.4 页面只展示 policy verdict、reason 和 manual next action，不提供自动刷新、修复、确认、规则应用或交易按钮。
- [x] 4.5 增加/更新 Vitest 覆盖 policy 三态和安全文案。

## 5. 发布材料与验收文档

- [x] 5.1 新增 `docs/release/acceptance/2026-06-18-p66-current-data-policy.md`，记录 policy 命令、结果、是否 blocked/waiver、发布影响和 Not Claimed。
- [x] 5.2 更新 `docs/release/acceptance-repeatability.md`，要求后续 release-ready 声明引用 P66 policy evidence。
- [x] 5.3 更新 `docs/release/release-handoff-2026-06-18.md` 和 `docs/release/README.md`，说明 current-data policy gate。
- [x] 5.4 更新 `docs/development-plan.md`、`docs/README.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 5.5 确认发布材料不承诺未来 provider 可用性、当前数据永远健康、收益、自动刷新、自动修复、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、登录源、付费源、授权源、Level2 或高频源。

## 6. 测试与验证

- [x] 6.1 运行 focused Go tests：`go test ./internal/application/service ./internal/application/handler ./cmd/agent`。
- [x] 6.2 运行 `go test ./...`。
- [x] 6.3 运行 `npm --prefix web test`。
- [x] 6.4 运行 `npm --prefix web run build`。
- [x] 6.5 运行 `bash scripts/e2e-smoke.sh`。
- [x] 6.6 运行 P66 current-data policy gate，记录 passed/waiver/blocked 结论和发布影响。
- [x] 6.7 运行安全扫描，确认无完整 key、私有路径、raw prompt、raw provider payload 或新增禁止能力入口。
- [x] 6.8 运行 `openspec validate p66-current-data-zero-degradation-policy --strict`、`openspec validate --all --strict`、`git diff --check`。

## 7. 复审、归档与提交

- [x] 7.1 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 7.2 执行 OpenSpec archive。
- [x] 7.3 archive 后确认无活跃 change，并规划 P66 后下一步。
- [x] 7.4 提交前子 agent 复审无 Critical / Important。
- [x] 7.5 提交 P66。
