# P53 设计：验收执行与发布候选材料

## 设计目标

P53 是一次发布前验收执行，不是功能开发。它把 P52 的门禁矩阵转化为可审计结果，并根据结果生成发布候选材料。若门禁失败，P53 的正确输出是“记录失败和发布影响”，不是修改代码绕过失败。

## 产物结构

新增目录：

- `docs/release/acceptance/`
- `docs/release/`

新增文档：

- `docs/release/acceptance/2026-06-17-p53-acceptance-run.md`
  - 记录日期、commit、环境、配置、G0-G9 命令、状态、产物、说明、发布影响。
  - 状态只允许 `pass`、`degraded`、`blocked`、`skipped`。
  - 对 G6/G7 真实测试失败使用 P52 的失败分类。
- `docs/release/release-candidate-2026-06-17.md`
  - 引用 P51 审计证据包和 P52 门禁矩阵。
  - 引用 P53 实际验收记录。
  - 汇总 release status：`release_ready` 或 `release_blocked`。
  - 若存在阻断项，明确列出后续修复阶段建议，不得宣称可发布。

临时产物目录：

- `tmp/acceptance/p53-2026-06-17/`
  - `logs/*.log`
  - `install/**`
  - `release-upgrade/**`
  - `config.real-public.yaml`
  - `config.real-llm.yaml`
  - `data/*.db`
  - 不提交。

## 门禁执行策略

### G0 治理与文档一致性

命令：

```bash
openspec validate --all --strict
git diff --check
find openspec/changes -maxdepth 1 -mindepth 1 -type d ! -name archive -print
```

计划阶段允许 `find` 输出当前 P53 change；执行后归档前必须只剩当前 change，归档后必须无活跃 change。

### G1 Go 全量测试

命令：

```bash
go test ./...
```

失败默认阻断发布，P53 只记录失败，不修复。

### G2 Go 聚焦集成

命令：

```bash
go test ./cmd/agent ./cmd/server ./internal/application/workflow ./internal/application/handler ./internal/infrastructure/persistence/sqlite
```

失败默认阻断发布。

### G3 前端测试与构建

命令：

```bash
npm --prefix web test -- --run
npm --prefix web run build
```

失败默认阻断发布。

### G4 浏览器 E2E smoke

命令：

```bash
bash scripts/e2e-smoke.sh
```

若本机缺 Playwright Chromium 或端口资源导致无法运行，记录为 `skipped` 或 `blocked`，且不得声明浏览器验收通过。

### G5 本地 fixture/current smoke

命令：

```bash
bash scripts/recovery-smoke.sh
go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300
go run ./cmd/agent --task data-source-quality-regression --source fixture --symbol 000300
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300
```

`current` 可为 `degraded`，但必须记录分类和影响范围。

### G6 真实公开源 opt-in

使用临时配置 `tmp/acceptance/p53-2026-06-17/config.real-public.yaml`：

- SQLite 指向 `tmp/acceptance/p53-2026-06-17/data/real-public.db`
- VecLite 指向 `tmp/acceptance/p53-2026-06-17/data/veclite-real-public`
- `data_sources.use_stub=false`
- `data_sources.public_evidence.enabled=true`
- `data_sources.public_evidence.sources` 使用 `cninfo`、`szse`、`csrc`

命令：

```bash
go run ./cmd/agent --config tmp/acceptance/p53-2026-06-17/config.real-public.yaml --task public-evidence-refresh --symbol 000001 --start-date 2026-06-01 --end-date 2026-06-17
```

失败必须同时记录实际错误码和 P52 分类。当前 collector 常见实际错误码包括 `source_unavailable`、`parse_error`、`no_data`；验收记录需映射到 P52 的 `network`、`rate_limit`、`authentication_or_key`、`source_schema_change`、`no_data`、`parse_failure` 或其他明确分类，避免把实际结果改写为未发生的错误类型。

### G7 真实 LLM opt-in

使用临时配置 `tmp/acceptance/p53-2026-06-17/config.real-llm.yaml`，从 `configs/config.local.yaml` 的测试 LLM 配置派生，SQLite/VecLite 改到 `tmp/acceptance/p53-2026-06-17/data/`。不得在验收记录中写完整 key。

命令：

```bash
go run ./cmd/agent --config tmp/acceptance/p53-2026-06-17/config.real-llm.yaml --task llm-smoke --symbol 510300
```

如果模型不可用或 key 失败，记录为 `model_unavailable` 或 `authentication_or_key`，阻断 LLM 能力声明。若不阻断全项目发布，验收记录和发布材料必须写明 waiver 是否成立、waiver 理由和不可声明的 LLM 能力范围；否则按 `release_blocked` 处理。

### G8 本地安装与升级

命令：

```bash
bash scripts/local-install-diagnostics.sh --config configs/config.example.yaml --include-release-upgrade --target-version p53-acceptance --output-dir tmp/acceptance/p53-2026-06-17/install
bash scripts/local-release-upgrade-check.sh --config configs/config.example.yaml --target-version p53-acceptance --output-dir tmp/acceptance/p53-2026-06-17/release-upgrade
```

失败默认阻断发布，除非仅 E2E 浏览器依赖缺失且已在 G4/G8 明确记录为 `skipped`。

### G9 安全边界与脱敏

命令：

```bash
rg -n "自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|自动修复|收益承诺|Level2|高频" docs openspec internal cmd web scripts
rg -n "sk-[A-Za-z0-9]|PRIVATE KEY|原始 SQL|raw HTTP|完整 prompt" docs openspec internal cmd web scripts
```

允许命中禁止能力的“禁止/不得/不会”说明。任何完整 key、私有路径、raw payload 或完整 prompt 泄漏均阻断发布。

执行后还必须人工检查 `docs/release/**` 和 `tmp/acceptance/p53-2026-06-17/**` 摘要文件，确认最终提交材料没有引用完整 key、raw body、完整 prompt 或未脱敏私有路径。`configs/config.local.yaml` 和 `tmp/` 不提交，但其内容不得被复制进 release 文档。

## 记录策略

每个门禁记录：

- `Status`
- `Command`
- `Artifact`
- `Notes`
- `Release impact`

最终 release status 规则：

- 任意 G0-G5/G8/G9 为 `blocked` -> `release_blocked`
- G6/G7 失败但可归类为外部依赖问题 -> 不声明对应真实能力通过；是否 `release_blocked` 由材料明确说明
- 任意 `redaction_failure` -> `release_blocked`
- 全部阻断门禁通过，且 G6/G7 通过或有可接受 waiver -> `release_ready`

## 安全边界

P53 只执行验收和整理材料，不改变运行时代码，不引入自动修复、自动迁移、自动交易、外部推送或收益承诺。
