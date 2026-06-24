# Tasks: P111 高保真参考图视觉重构

## 1. Governance And Reference Lock

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 `docs/frontend-contract.md`、`docs/ui-design.md`、`docs/product-experience-polish-roadmap.md`。
- [x] 1.3 打开并记录参考图 `/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png`。
- [x] 1.4 创建并校验 `p111-high-fidelity-reference-redesign` OpenSpec change。
- [x] 1.5 更新 `docs/GOVERNANCE.md` 与 `openspec/PROGRESS.md`，将 P111 标记为当前活跃 change。
- [x] 1.6 确认 P111 不使用 `docs/superpowers/plans/` 作为规格真源。

## 2. Reference Extraction And Page Gates

- [x] 2.1 建立 reference extraction ledger，逐项记录 sidebar、topbar、hero、action queue、metric grid、snapshot strip、progress tracker、evidence checklist、ledger surface 的目标规则。
- [x] 2.2 建立 page fidelity matrix，列出 `/`、`/workbench`、`/positions`、`/data-quality`、`/risk-alerts`、`/evidence`、`/consultation`、`/decisions/:id`、`/decision-loop`、`/rules`、`/review`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run`、`/local-install`、`/local-knowledge`、`/settings`。
- [x] 2.3 定义 mismatch 等级 P0/P1/P2/P3/pass，并写入 acceptance 模板。
- [x] 2.4 定义每页完成门禁：截图、ledger、修复 P0/P1/P2、重新截图、无横向溢出、无 console app error。

## 3. Tests First

- [x] 3.1 新增 reference component tests，覆盖 report hero、priority action queue、metric grid、snapshot strip、progress tracker、evidence checklist、ledger surface。
- [x] 3.2 更新 `AppLayout` tests，覆盖 reference topbar、sidebar active 状态、本地模式、安全边界和原路由名称。
- [x] 3.3 更新 Dashboard/Workbench tests，断言 reference-style first viewport 模块存在，且不再是 P110 generic vertical card stack。
- [x] 3.4 更新 DataQuality/RiskAlert/Evidence/DecisionLoop tests，断言页面使用 reference component language。
- [x] 3.5 增加或更新 forbidden affordance / sensitive scan 测试或脚本。
- [x] 3.6 运行 targeted tests，确认新增测试在实现前失败或覆盖缺口明确。

## 4. Shared Reference Components

- [x] 4.1 新增 `web/src/components/reference/referenceTypes.ts`。
- [x] 4.2 新增 reference topbar shell，支持页面标题、日期、本地模式、数据截至、刷新/查看类动作。
- [x] 4.3 新增 `ReferenceHero`，支持左 icon well、主标题、状态行、说明、当前纪律状态、禁止动作。
- [x] 4.4 新增 `PriorityActionQueue`，支持编号、优先级 chip、detail、meta、右侧 outline action。
- [x] 4.5 新增 `StatusMetricGrid`，支持 icon well、强数值、status chip、底部 detail rows。
- [x] 4.6 新增 `SnapshotStrip`，支持资金/持仓/仓位/更新时间等横向快照。
- [x] 4.7 新增 `ProgressTracker`，支持桌面横向步骤和移动纵向 timeline。
- [x] 4.8 新增 `EvidenceChecklist`，支持证据、LLM、规则、审计、source health checklist。
- [x] 4.9 新增 `LedgerSurface`，统一 table/list/panel 密度。
- [x] 4.10 更新 `web/src/styles/global.css` P111 reference styles，避免 one-off 页面样式。

## 5. Core Cockpit High-Fidelity Implementation

- [x] 5.1 重构 `AppLayout` 为 reference shell：sidebar + topbar + command content gutters。
- [x] 5.2 重构 `DashboardFeature` (`/`) 为参考图结构：report hero、priority queue、status overview、snapshot、recent progress、evidence checklist。
- [x] 5.3 重构 `WorkbenchPage` 为同一结构，但侧重聚合工作台和跨页面任务。
- [x] 5.4 对 `/` 采集桌面截图，对照参考图填写 mismatch ledger；修复所有 P0/P1/P2。
- [x] 5.5 对 `/workbench` 采集桌面截图，对照参考图填写 mismatch ledger；修复所有 P0/P1/P2。

## 6. Maintenance And Evidence Pages

- [x] 6.1 重构 `/positions`：reference status hero、local fact action queue、portfolio metrics、ledger/form surfaces。
- [x] 6.2 重构 `/data-quality`：reference quality hero、source/RAG/LLM metric cards、next local checks、evidence checklist、ledger panels。
- [x] 6.3 重构 `/risk-alerts`：reference risk hero、SOP priority queue、risk metrics、alert ledger。
- [x] 6.4 重构 `/evidence`：reference evidence hero、source verification checklist、evidence ledger。
- [x] 6.5 重构 `/consultation` 和 `/decisions/:id`：reference decision hero、progress tracker、analysis panels、manual confirmation boundary。
- [x] 6.6 重构 `/decision-loop`：reference progress tracker、loop ledger、audit/readback checklist。
- [x] 6.7 每个页面完成后采集桌面截图，填写 mismatch ledger，修复 P0/P1/P2。

## 7. Governance And Ops Pages

- [x] 7.1 重构 `/rules`、`/review`、`/audit` 为 reference governance/ledger language。
- [x] 7.2 重构 `/notifications`、`/daily-discipline/reports`、`/daily-auto-run` 为 reference inbox/report/ops language。
- [x] 7.3 重构 `/local-install`、`/local-knowledge`、`/settings` 为 reference local ops language。
- [x] 7.4 每个页面完成后采集桌面截图，填写 mismatch ledger，修复 P0/P1/P2。

## 8. Responsive And Browser QA

- [x] 8.1 启动真实本地 Go backend、Vite frontend 和临时 SQLite。
- [x] 8.2 捕获所有 P111 覆盖页面桌面截图和 console/reflow JSON。
- [x] 8.3 捕获核心路由 390px 和桌面截图及 no-overflow JSON。
- [x] 8.4 使用 `view_image` 对参考图与最新渲染图做最终人工对照。
- [x] 8.5 生成 `docs/release/ui-audit-assets/2026-06-24-p111-high-fidelity-reference-redesign/visual-mismatch-ledger.md`。
- [x] 8.6 若发现任何 P0/P1/P2 mismatch，回到对应任务修复并重新截图。

## 9. Validation

- [x] 9.1 运行 `npm --prefix web test`。
- [x] 9.2 运行 `npm --prefix web run build`。
- [x] 9.3 运行 `go test ./...`。
- [x] 9.4 运行 `go vet ./...`。
- [x] 9.5 运行 forbidden affordance scan。
- [x] 9.6 运行 sensitive/redaction scan。
- [x] 9.7 运行 `openspec validate p111-high-fidelity-reference-redesign --strict`。
- [x] 9.8 运行 `openspec validate --all --strict`。
- [x] 9.9 运行 `git diff --check`。

## 10. Documentation, Archive, Commit

- [x] 10.1 新增 P111 acceptance record，包含参考图、页面矩阵、mismatch ledger、截图目录、命令结果、安全边界。
- [x] 10.2 更新 `docs/ui-design.md`、`docs/frontend-contract.md`、`docs/product-experience-polish-roadmap.md`、`docs/development-plan.md`。
- [x] 10.3 复审 P111 diff，确认无 Critical / Important；若可使用 subagent 且用户明确要求，则执行子 agent 复审，否则记录本地等价复审。
- [x] 10.4 执行 OpenSpec archive，合并 delta。
- [x] 10.5 最终验证后提交 P111。
