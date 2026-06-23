# Tasks: P53 验收执行与发布候选材料

## 1. 范围与计划

- [x] 1.1 确认当前无活跃 change，P53 为下一阶段。
- [x] 1.2 创建 `p53-acceptance-execution-and-release-candidate-materials` OpenSpec change。
- [x] 1.3 确认 P53 执行 P52 G0-G9，并基于结果生成验收记录和发布候选材料；若门禁失败，必须记录阻断，不修复运行时代码。
- [x] 1.4 子 agent 复审 P53 计划，无 Critical / Important 后继续。

## 2. 验收环境准备

- [x] 2.1 创建 `tmp/acceptance/p53-2026-06-17/logs/` 和 `tmp/acceptance/p53-2026-06-17/data/`。
- [x] 2.2 记录当前 commit、branch、Go/Node/npm 版本和操作系统信息。
- [x] 2.3 为 G6 写入临时真实公开源配置 `tmp/acceptance/p53-2026-06-17/config.real-public.yaml`。
- [x] 2.4 为 G7 写入临时真实 LLM 配置 `tmp/acceptance/p53-2026-06-17/config.real-llm.yaml`，从 `configs/config.local.yaml` 派生，但验收记录不得包含完整 key。

## 3. 执行 P52 G0-G9 门禁

- [x] 3.1 执行 G0：`openspec validate --all --strict`、`git diff --check`、活跃 change 检查，并记录当前 P53 change 为预期活跃项。
- [x] 3.2 执行 G1：`go test ./...`。
- [x] 3.3 执行 G2：`go test ./cmd/agent ./cmd/server ./internal/application/workflow ./internal/application/handler ./internal/infrastructure/persistence/sqlite`。
- [x] 3.4 执行 G3：`npm --prefix web test -- --run`、`npm --prefix web run build`。
- [x] 3.5 执行 G4：`bash scripts/e2e-smoke.sh`。
- [x] 3.6 执行 G5：`bash scripts/recovery-smoke.sh`、`go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300`、`go run ./cmd/agent --task data-source-quality-regression --source fixture --symbol 000300`、`go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300`。
- [x] 3.7 执行 G6：`go run ./cmd/agent --config tmp/acceptance/p53-2026-06-17/config.real-public.yaml --task public-evidence-refresh --symbol 000001 --start-date 2026-06-01 --end-date 2026-06-17`；若失败，同时记录实际错误码和 P52 失败分类。
- [x] 3.8 执行 G7：`go run ./cmd/agent --config tmp/acceptance/p53-2026-06-17/config.real-llm.yaml --task llm-smoke --symbol 510300`；若失败，明确是否成立 waiver，否则写入 `release_blocked`。
- [x] 3.9 执行 G8：`bash scripts/local-install-diagnostics.sh --config configs/config.example.yaml --include-release-upgrade --target-version p53-acceptance --output-dir tmp/acceptance/p53-2026-06-17/install` 和 `bash scripts/local-release-upgrade-check.sh --config configs/config.example.yaml --target-version p53-acceptance --output-dir tmp/acceptance/p53-2026-06-17/release-upgrade`。
- [x] 3.10 执行 G9 安全边界与脱敏扫描，并人工复核命中是否为禁止性说明或真实泄露。

## 4. 生成验收与发布材料

- [x] 4.1 新增 `docs/release/acceptance/2026-06-17-p53-acceptance-run.md`，按 G0-G9 记录命令、状态、产物、说明和发布影响。
- [x] 4.2 新增 `docs/release/release-candidate-2026-06-17.md`，引用 P51、P52、P53 结果，并写明 `release_ready` 或 `release_blocked`。
- [x] 4.3 若存在 degraded/skipped/blocked，逐项写明分类、原因、waiver 或后续修复建议。
- [x] 4.4 确认验收材料和将提交的 release 文档不包含完整 key、私有路径、raw HTTP 响应、完整 prompt 或原始 SQL；临时配置和日志不得被复制进提交材料。

## 5. 文档与进度同步

- [x] 5.1 更新 `openspec/PROGRESS.md`，标记 P53 活跃并根据验收结论规划下一阶段。
- [x] 5.2 更新 `openspec/project.md`、`docs/GOVERNANCE.md`、`docs/development-plan.md`、`AGENTS.md`。
- [x] 5.3 更新 `docs/README.md`，加入 release/acceptance 文档入口。

## 6. 验证与复审

- [x] 6.1 执行 `openspec validate p53-acceptance-execution-and-release-candidate-materials --strict`。
- [x] 6.2 执行 `openspec validate --all --strict`。
- [x] 6.3 执行 `git diff --check`。
- [x] 6.4 子 agent 复审验收结果和发布材料，无 Critical / Important 后归档。

## 7. 归档

- [x] 7.1 执行 OpenSpec archive，将 delta 合并入对应 `docs/` 真源，并按需同步 `openspec/specs/` 摘要。
- [x] 7.2 归档后确认无活跃 change。
- [x] 7.3 提交前子 agent 复审无 Critical / Important。
- [x] 7.4 提交 P53。
