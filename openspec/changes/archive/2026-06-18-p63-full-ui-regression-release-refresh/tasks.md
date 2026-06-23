# Tasks: P63 全量真实 UI 回归与发布状态刷新

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 `docs/product-experience-polish-roadmap.md`、`docs/project-acceptance-gate-matrix.md`、P53 release candidate / handoff 和 P62 验收记录。
- [x] 1.3 使用 Product Design get-context playback 明确 P63 brief：现有产品、现有设计系统、full interactivity。
- [x] 1.4 创建 `p63-full-ui-regression-release-refresh` OpenSpec change。
- [x] 1.5 更新当前进度文档，标记 P63 active。
- [x] 1.6 运行 `openspec validate p63-full-ui-regression-release-refresh --strict`、`openspec validate --all --strict`、`git diff --check`。
- [x] 1.7 子 agent 方案复审无 Critical / Important 后执行。

## 2. 验收脚本与证据采集准备

- [x] 2.1 盘点现有 `scripts/e2e-smoke.sh`、`web/e2e/local-smoke.spec.ts`、P58-P62 验收脚本和截图资产。
- [x] 2.2 如现有 E2E 不足，新增或扩展 P63 全路由浏览器验收，覆盖页面标题、主要状态、关键动作、console/page error、HTTP 失败和 forbidden copy。
- [x] 2.3 增加或复用 390px、768px、1280px reflow 检查，确保每个主要路由无页面级横向溢出。
- [x] 2.4 准备 `docs/release/ui-audit-assets/2026-06-18-p63/` 资产目录，只保存截图和脱敏 JSON 摘要。
- [x] 2.5 确认验收脚本不会提交临时 DB、Playwright trace、完整 raw response、完整 prompt、完整 key、私有路径或本地日志。

## 3. 自动化门禁执行

- [x] 3.1 执行 P52 G0：`openspec validate --all --strict`、`git diff --check`、确认 `openspec/changes/` 下只有当前 P63 活跃 change。
- [x] 3.2 执行 P52 G1：`go test ./...`。
- [x] 3.3 执行 P52 G2：`go test ./cmd/agent ./cmd/server ./internal/application/workflow ./internal/application/handler ./internal/infrastructure/persistence/sqlite`。
- [x] 3.4 执行 P52 G3：`npm --prefix web test` 与 `npm --prefix web run build`。
- [x] 3.5 执行 P52 G4：`bash scripts/e2e-smoke.sh`。
- [x] 3.6 执行 P52 G5：`bash scripts/recovery-smoke.sh`、`go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300`、fixture/current data-source quality regression。
- [x] 3.7 执行 P52 G6：真实公开源 opt-in `public-evidence-refresh`，使用临时 SQLite、显式日期窗口和真实源配置；若失败按 P52 分类记录。
- [x] 3.8 执行 P52 G7：真实 LLM opt-in `go run ./cmd/agent --task llm-smoke --symbol 510300`，使用测试 key 配置并验证脱敏审计；若失败按 P52 分类记录。
- [x] 3.9 执行 P52 G8：本地安装诊断和 release-upgrade 检查，输出到 `tmp/`。
- [x] 3.10 执行 P52 G9：禁止能力文案扫描和敏感信息扫描，人工复核命中项。

## 4. 真实 UI 与真实 LLM 回归

- [x] 4.1 启动真实本地后端和 Vite 前端，使用临时 SQLite/配置或已确认的安全测试配置。
- [x] 4.2 浏览器操作全主要路由：`/`、`/workbench`、`/consultation`、`/decisions/:decisionId`、`/evidence`、`/decision-loop`、`/positions`、`/data-quality`、`/risk-alerts`、`/risk-alerts/:alertId`、`/rules`、`/audit`、`/notifications`、`/daily-auto-run`、`/daily-discipline/reports`、`/daily-discipline/reports/:reportId`、`/review`、`/local-install`、`/local-knowledge`、`/settings`。
- [x] 4.3 对每个路由记录桌面/平板/移动端 reflow、标题/landmark、关键状态、下一步人工动作、空/错误/降级态和安全文案。
- [x] 4.4 执行真实 LLM-backed consultation UI journey，并打开新生成的 decision detail。
- [x] 4.5 若真实 LLM 或公开源失败，按 P52 分类记录为 network、rate_limit、authentication_or_key、model_unavailable、quality_failure 等，并说明 release impact。
- [x] 4.6 如发现阻断级前后端或 UI 缺陷，只在既有 API/schema/workflow/LLM 能力边界内做最小修复并重新执行对应验收；非阻断问题记录到 P63 验收文档。

## 5. 发布材料刷新

- [x] 5.1 新增 `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`，记录 G0-G9、全路由 UI、真实 LLM、截图、降级、waiver 和结论。
- [x] 5.2 新增或更新 `docs/release/release-candidate-2026-06-18.md`，状态只允许 `release_ready`、`release_degraded` 或 `blocked`。
- [x] 5.3 新增或更新 `docs/release/release-handoff-2026-06-18.md`，说明交付边界、复验入口、已知降级和 Not Claimed。
- [x] 5.4 更新 `docs/product-experience-polish-roadmap.md`、`docs/development-plan.md`、`docs/README.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 5.5 确认 P63 发布材料不承诺收益、未来外部源可用性、券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用或自动修复。

## 6. 复审、归档与提交

- [x] 6.1 运行 `openspec validate p63-full-ui-regression-release-refresh --strict` 与 `openspec validate --all --strict`。
- [x] 6.2 运行 `git diff --check`。
- [x] 6.3 运行最终必要验证：`npm --prefix web test`、`npm --prefix web run build`、`go test ./...`、`bash scripts/e2e-smoke.sh`。
- [x] 6.4 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 6.5 执行 OpenSpec archive。
- [x] 6.6 archive 后确认无活跃 change，并规划 P63 后下一步。
- [x] 6.7 提交前子 agent 复审无 Critical / Important。
- [ ] 6.8 提交 P63。
