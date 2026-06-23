# Tasks: P67 Current Data Gate Resolution Workflow

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 P66 spec、P66 acceptance、P52/P54 repeatability、现有 data-source-quality service/API/CLI/UI。
- [x] 1.3 使用 Product Design get-context 播放 P67 UI brief：沿用现有 `/data-quality` operational cockpit，全交互本地记录，不做新视觉探索。
- [x] 1.4 创建 `p67-current-data-gate-resolution-workflow` OpenSpec change。
- [x] 1.5 写明 P67 只新增本地人工 resolution 记录和 release claim state，不改变 P66 policy，不新增外部源、provider 调用、自动刷新、自动修复、券商接口或交易能力。
- [x] 1.6 更新当前进度文档，标记 P67 active。
- [x] 1.7 运行 `openspec validate p67-current-data-gate-resolution-workflow --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 1.8 子 agent 方案复审无 Critical / Important 后执行。

## 2. 后端持久化与领域模型

- [x] 2.1 新增 SQLite migration `data_quality_gate_resolutions`。
- [x] 2.2 新增 domain/repository model，覆盖 resolution type/status、policy fingerprint、policy summary、scope/reason/release impact/evidence ref、copied reasons、created_by/retired_by、safety note、timestamps。
- [x] 2.3 实现 repository create-or-reuse active、list newest first、find active matching policy fingerprint、retire。
- [x] 2.4 增加 repository tests 覆盖 create/reuse/list/retire/retired ignored、同一 `symbol + policy_fingerprint` 只允许一个 active resolution。

## 3. 服务、DTO、API 与审计

- [x] 3.1 新增 DTO：resolution record、resolution check、create request、retire response。
- [x] 3.2 新增 service：读取 P66 current policy，计算 canonical `policy_fingerprint`、`release_claim_state`、`clean_data_claim_allowed`、fixed allowed/prohibited claims。
- [x] 3.3 实现 create validation：policy 必须 blocked/waiver_required；blocked 只允许 `scope_exclusion`，waiver_required 允许 `waiver` 或 `scope_exclusion`；scope/reason/release impact 必填；文本脱敏；同类型重复 active record 复用；不同类型 active 冲突拒绝。
- [x] 3.4 实现 retire：只更新 resolution status 并写 sanitized audit。
- [x] 3.5 新增 handlers/routes：GET gate-resolution、GET resolutions、POST resolutions、POST retire。
- [x] 3.6 增加 service/handler tests 覆盖 pass、requires_resolution、waiver_required resolved_with_waiver、blocked resolved_with_scope_exclusion、blocked waiver rejected、policy_fingerprint stable matching、duplicate reuse、conflicting active rejected、retire、redaction、GET 不写 audit、POST 写 sanitized audit。

## 4. CLI 验收入口

- [x] 4.1 新增 `data-source-quality-resolution-check` task。
- [x] 4.2 CLI 输出 compact sanitized `policy=<verdict>`、`gate=<gate>`、`fingerprint=<hash>`、`resolution=<none|waiver|scope_exclusion>`、`claim_state=<state>`。
- [x] 4.3 `requires_resolution` 返回非零；`pass`、`waiver_required + resolved_with_waiver`、`resolved_with_scope_exclusion` 返回 0；blocked policy 只有 matching `scope_exclusion` 可返回 0。
- [x] 4.4 保持 P66 `--strict-quality-gate` 行为不变。
- [x] 4.5 增加 CLI tests 覆盖 unresolved exit 1、resolved exit 0、strict gate 仍 blocked、audit output_ref 脱敏。

## 5. 前端产品化

- [x] 5.1 新增前端 types/service：gate resolution check、resolution list、create、retire。
- [x] 5.2 扩展 `dataQualityExperienceModel`：加入 release claim state、active resolution、manual resolution action。
- [x] 5.3 扩展 `/data-quality`：展示当前数据门禁处置状态、active resolution、allowed/prohibited claims。
- [x] 5.4 新增本地记录表单：resolution type、scope、reason、release impact、evidence ref；blocked policy 只展示 `scope_exclusion`，waiver_required policy 展示 `waiver` 与 `scope_exclusion`。
- [x] 5.5 新增 retire action，只撤销本地 resolution 记录。
- [x] 5.6 增加 Vitest 覆盖 unresolved、waiver_required resolved waiver、blocked resolved scope exclusion、blocked 不展示 waiver、retire、sanitization、forbidden copy scan。

## 6. 发布材料与治理文档

- [x] 6.1 新增 `docs/release/acceptance/2026-06-18-p67-current-data-resolution.md`，记录 resolution check 命令、UI/CLI 结果、发布声明边界和 Not Claimed。
- [x] 6.2 更新 `docs/release/acceptance-repeatability.md`，要求 future release-ready 引用 P66 policy 和 P67 resolution state。
- [x] 6.3 更新 `docs/release/release-handoff-2026-06-18.md` 和 `docs/release/README.md`，说明 P67 resolution state。
- [x] 6.4 更新 `docs/development-plan.md`、`docs/README.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 6.5 确认发布材料不声明当前数据 clean，除非 P66 policy pass；不承诺未来 provider 可用性、收益、自动刷新、自动修复、券商接口、交易、外推或自动规则能力。

## 7. 测试与验证

- [x] 7.1 按 TDD 先写失败测试，再实现后端 repository/service/handler/CLI。
- [x] 7.2 按 TDD 先写失败测试，再实现前端 model/page/service。
- [x] 7.3 运行 focused Go tests：`go test ./internal/infrastructure/persistence/sqlite ./internal/application/service ./internal/application/handler ./cmd/agent`。
- [x] 7.4 运行 `go test ./...`。
- [x] 7.5 运行 `npm --prefix web test`。
- [x] 7.6 运行 `npm --prefix web run build`。
- [x] 7.7 运行 `bash scripts/e2e-smoke.sh`。
- [x] 7.8 运行 P66 strict gate 和 P67 resolution check，记录 `policy_fingerprint`、release claim state 和 blocked scope exclusion 边界。
- [x] 7.9 运行安全扫描，确认无完整 key、私有路径、raw prompt、raw provider payload 或新增禁止能力入口。
- [x] 7.10 运行 `openspec validate p67-current-data-gate-resolution-workflow --strict`、`openspec validate --all --strict`、`git diff --check`。

## 8. 复审、归档与提交

- [x] 8.1 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 8.2 执行 OpenSpec archive。
- [x] 8.3 archive 后确认无活跃 change，并规划 P67 后下一步。
- [x] 8.4 提交前子 agent 复审无 Critical / Important。
- [x] 8.5 提交 P67。
