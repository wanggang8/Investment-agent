# 项目验收门禁矩阵

> 更新时间：2026-06-17
>
> 作用：定义发布候选材料前必须使用的项目级验收门禁。本文定义门禁，不代表这些门禁已经全部执行或通过。

## 1. 总览

P52 承接 P51 `docs/p19-p24-audit-evidence-pack.md`，把项目已有验证入口整理成发布前可执行矩阵。P53 发布候选材料必须引用实际门禁结果或明确 waiver，不能仅凭本文件声明发布就绪。

发布阻断原则：

1. G0-G5、G8、G9 默认阻断发布。
2. G6 真实公开源和 G7 真实 LLM 为显式 opt-in 门禁；失败必须分类。若失败由网络、限流、额度或外部服务不可用导致，可记录为降级或 waiver，但不得声明对应真实能力通过。
3. 任意安全边界或脱敏问题均阻断发布。
4. 任意命令失败但未记录原因、产物和处理结论，视为阻断。

## 2. 门禁矩阵

| Gate | 分类 | 命令或入口 | 前置条件 | 通过标准 | 允许降级 | 产物位置 | 阻断发布 |
| --- | --- | --- | --- | --- | --- | --- | --- |
| G0 | 治理与文档一致性 | `openspec validate --all --strict`；`git diff --check`；`find openspec/changes -maxdepth 1 -mindepth 1 -type d ! -name archive -print` | 工作树处于待验收提交状态 | OpenSpec 全部通过；无 diff 格式问题；归档后无活跃 change 残留 | 不允许 | 命令输出、最终提交记录 | 是 |
| G1 | Go 全量测试 | `go test ./...` | Go 依赖可用；不需要真实外部源 | 全量测试 exit 0 | 不允许 | 命令输出 | 是 |
| G2 | Go 聚焦集成 | `go test ./cmd/agent ./cmd/server ./internal/application/workflow ./internal/application/handler ./internal/infrastructure/persistence/sqlite` | 本地 SQLite 临时库测试可运行 | CLI、server、workflow、handler、persistence 聚焦测试 exit 0 | 不允许 | 命令输出 | 是 |
| G3 | 前端测试与构建 | `npm --prefix web test -- --run`；`npm --prefix web run build` | `web/node_modules` 可用；必要时先 `npm --prefix web install` | Vitest 与 build exit 0 | 不允许 | 命令输出、build 日志 | 是 |
| G4 | 浏览器 E2E smoke | `bash scripts/e2e-smoke.sh` | 本机具备 Playwright Chromium；端口可用 | Playwright smoke 通过，临时 SQLite/日志/trace 不污染工作树 | 仅当本机缺浏览器依赖时可 `skipped`，但 P53 不能声明浏览器验收通过 | `tmp/` 临时产物或脚本日志 | 是 |
| G5 | 本地 fixture/current smoke | `bash scripts/recovery-smoke.sh`；`go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300`；`go run ./cmd/agent --task data-source-quality-regression --source fixture --symbol 000300`；`go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300` | 使用临时配置或确认只读当前库；不得覆盖真实库 | 恢复 smoke 通过；fixture 回归 passed；current 模式若 degraded 必须解释 | current 可 degraded，但必须分类并说明影响范围 | `tmp/recovery-smoke/`、CLI 输出、审计摘要 | 是 |
| G6 | 真实公开源 opt-in | `go run ./cmd/agent --task public-evidence-refresh --symbol 000001 --start-date YYYY-MM-DD --end-date YYYY-MM-DD` | 显式真实源配置：`data_sources.use_stub=false`、`data_sources.public_evidence.enabled=true`；使用临时 SQLite；窗口明确 | 成功写入 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications`、`audit_events`，或 no_data 被正确分类 | 网络、限流、source_unavailable、no_data、parse_error 可降级；必须分类 | 临时 SQLite、CLI 输出、审计摘要 | 条件阻断 |
| G7 | 真实 LLM opt-in | `go run ./cmd/agent --task llm-smoke --symbol 510300` | 本地配置存在临时/测试 key、base_url、model、timeout；不得输出完整 key | 调用成功，parse/quality 通过，审计摘要脱敏，LLM 不写最终裁决 | 网络、额度、认证、模型不可用可降级；质量失败阻断 LLM 能力声明 | CLI 输出、`audit_events` 脱敏摘要 | 条件阻断 |
| G8 | 本地安装与升级 | `bash scripts/local-install-diagnostics.sh --include-release-upgrade --target-version vNEXT --output-dir <tmp>/install`；`bash scripts/local-release-upgrade-check.sh --target-version vNEXT --output-dir <tmp>/release-upgrade` | 使用示例配置或临时配置；输出目录在 `tmp/` | 生成 install summary、preflight/recovery/e2e/release-upgrade 结果；诊断脱敏 | e2e 可因浏览器环境跳过，但必须记录；preflight failed 阻断 | `<tmp>/install/**`、`<tmp>/release-upgrade/**` | 是 |
| G9 | 安全边界与脱敏 | `rg -n "自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|自动修复|收益承诺|Level2|高频" docs openspec internal cmd web scripts`；`rg -n "sk-[A-Za-z0-9]|PRIVATE KEY|原始 SQL|raw HTTP|完整 prompt" docs openspec internal cmd web scripts` 后人工复核 | 在最终 diff 上执行；允许命中禁止能力的“禁止/不得/不会”说明 | 无新增高风险入口；无完整 key、私有路径、原始 SQL、完整 prompt、供应商 raw 响应泄露 | 只允许安全边界说明命中；需人工复核 | 扫描输出和复核记录 | 是 |

## 3. 真实测试失败分类

真实公开源和真实 LLM 门禁必须显式 opt-in，并按以下分类记录：

| 分类 | 说明 | 默认处理 |
| --- | --- | --- |
| `network` | DNS、TLS、连接、超时等网络问题 | 可降级，需重试建议 |
| `rate_limit` | 公开源或模型供应商限流 | 可降级，需记录等待或低频策略 |
| `authentication_or_key` | key 缺失、错误、额度耗尽、权限不足 | 阻断对应真实能力声明 |
| `source_schema_change` | 真实源响应 shape 改变 | 阻断对应 collector/解析声明 |
| `no_data` | 源可达但窗口无数据 | 可降级，不等于接口失败 |
| `parse_failure` | 响应可达但解析失败 | 阻断对应解析声明 |
| `model_unavailable` | 模型不存在、服务不可用、base_url 错误 | 阻断对应 LLM 能力声明 |
| `quality_failure` | LLM 输出未通过 parse/quality gate | 阻断 LLM 质量声明 |
| `redaction_failure` | 输出或产物泄露 key、私有路径、raw payload | 阻断发布 |

真实测试通过只说明本次测试窗口可用，不代表收益承诺、未来可用性承诺、交易能力或自动决策能力。

## 4. 验收记录格式

P52 只定义格式，不创建实际验收结果。P53 或发布前验收应按以下格式记录：

```markdown
# Acceptance Run: <label>

- Date:
- Commit:
- Operator:
- Environment:
- Config:

| Gate | Status | Command | Artifact | Notes | Release impact |
| --- | --- | --- | --- | --- | --- |
| G0 | pass/degraded/blocked/skipped | `...` | path or log | reason | blocks/does_not_block |
```

建议位置：`docs/release/acceptance/YYYY-MM-DD-<label>.md`，或作为 P53 发布候选材料的一部分。

## 5. P53 使用要求

P53 发布候选材料必须：

1. 引用 P51 审计证据包。
2. 引用本门禁矩阵。
3. 填入实际验收结果或明确 waiver。
4. 对任何 skipped/degraded/blocked 门禁给出发布影响。
5. 不得把 P52 文档本身当作验收已通过的证据。

## 6. 安全边界

P52 不新增运行时能力，也不改变任何安全边界。验收、发布和材料整理均不得引入或暗示以下能力：券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
