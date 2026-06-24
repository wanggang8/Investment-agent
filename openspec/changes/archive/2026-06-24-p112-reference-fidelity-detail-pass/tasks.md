# Tasks: P112 参考图高保真细节修复

## 1. Governance And Baseline

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 P111 proposal/tasks/acceptance 和 P111 后 fresh audit findings。
- [x] 1.3 打开参考图并确认 P112 仍以第二方案为唯一视觉真源。
- [x] 1.4 创建并校验 `p112-reference-fidelity-detail-pass` OpenSpec change。
- [x] 1.5 更新 `docs/GOVERNANCE.md` 与 `openspec/PROGRESS.md`，将 P112 标记为当前活跃 change。

## 2. P112 Mismatch Baseline

- [x] 2.1 启动本地 backend/frontend，采集当前 P111 后全 18 个桌面路由截图。
- [x] 2.2 采集关键移动路由截图：`/`、`/workbench`、`/data-quality`、`/evidence`、`/settings`。
- [x] 2.3 建立 P112 mismatch ledger，记录页面、截图、参考图差异、等级、修复策略和复验结果。
- [x] 2.4 将以下已知问题列为必须修复：二级页面旧式卡片结构、hero 过高、主状态区下移、侧栏过碎、移动首屏效率不足、状态/图标/ledger 密度不足。

## 3. Shared Reference Tightening

- [x] 3.1 收敛 reference CSS tokens：hero padding、高度、surface radius、border、shadow、type scale、list density。
- [x] 3.2 重修 `AppLayout` sidebar：减少分组碎片感，提升 active item、section label、底部状态与参考图一致性。
- [x] 3.3 优化 `ReferenceHero`：紧凑高度、icon well、状态栏、右侧 status/action block 与移动端 reflow。
- [x] 3.4 优化 `StatusMetricGrid`：tone colors、icon well、强数值、detail rows 与参考图一致性。
- [x] 3.5 优化 `PriorityActionQueue`、`ProgressTracker`、`EvidenceChecklist`、`LedgerSurface` 的密度、图标和边框节奏。

## 4. Core Cockpit Polish

- [x] 4.1 修复 `/` 今日纪律：状态语义、hero 视觉重心、行动队列密度、metric tone、progress/checklist 细节。
- [x] 4.2 修复 `/workbench`：与首页共享严谨 cockpit rhythm，同时保持工作台语义。
- [x] 4.3 重新截图 `/` 和 `/workbench`，对照参考图修复所有 P0/P1/P2。

## 5. Secondary Page Composition Fixes

- [x] 5.1 修复 `/positions`、`/data-quality`、`/risk-alerts`：压缩首屏，重排为 reference report + next action + ledger/checklist。
- [x] 5.2 修复 `/consultation`、`/decisions/:id`、`/decision-loop`、`/evidence`：统一决策/证据页面的 report hero、progress、ledger 和 checklist 语言。
- [x] 5.3 修复 `/rules`、`/review`、`/audit`：统一治理页面的紧凑 report/ledger 结构。
- [x] 5.4 修复 `/notifications`、`/daily-discipline/reports`、`/daily-auto-run`：让主状态区进入首屏参考图节奏，去除过多前置提示层。
- [x] 5.5 修复 `/local-install`、`/local-knowledge`、`/settings`：压缩过高 hero 和大卡片堆叠，改成紧凑本地运维 report surface。
- [x] 5.6 每个页面修完后截图、填 ledger、复修 P0/P1/P2。

## 6. Responsive And Browser QA

- [x] 6.1 使用 Browser 采集全 18 个桌面页面截图和 console/reflow JSON。
- [x] 6.2 使用 Browser 采集关键移动页面截图和 no-overflow JSON。
- [x] 6.3 使用 `view_image` 对参考图、桌面首页、代表性二级页面、移动页面做最终人工对比。
- [x] 6.4 更新 P112 mismatch ledger，确认无剩余 P0/P1/P2。

## 7. Subagent Review Loop

- [x] 7.1 启动子 agent A：全 18 个桌面页面逐页对比参考图，输出 Critical/Important/P0/P1/P2/P3 findings。
- [x] 7.2 启动子 agent B：关键移动页面和响应式 reflow 审查，输出 Critical/Important/P0/P1/P2/P3 findings。
- [x] 7.3 若任一子 agent 返回 Critical/Important/P0/P1/P2，修复后重新截图并再次发起对应子 agent 审查。
- [x] 7.4 只有当子 agent 复审无 Critical/Important/P0/P1/P2 时，P112 才可进入归档。

## 8. Validation

- [x] 8.1 运行 `npm --prefix web test -- --run`。
- [x] 8.2 运行 `npm --prefix web run build`。
- [x] 8.3 运行 `go test ./...`。
- [x] 8.4 运行 `go vet ./...`。
- [x] 8.5 运行 forbidden affordance scan。
- [x] 8.6 运行 sensitive/redaction scan。
- [x] 8.7 运行 `openspec validate p112-reference-fidelity-detail-pass --strict`。
- [x] 8.8 运行 `openspec validate --all --strict`。
- [x] 8.9 运行 `git diff --check`。

## 9. Documentation, Archive, Commit

- [x] 9.1 新增 P112 acceptance record，包含参考图、截图目录、page matrix、mismatch ledger、子 agent 复审结论和命令结果。
- [x] 9.2 更新 `docs/ui-design.md`、`docs/frontend-contract.md`、`docs/product-experience-polish-roadmap.md`、`docs/development-plan.md` 中的 P112 结果。
- [x] 9.3 更新 `docs/GOVERNANCE.md` 与 `openspec/PROGRESS.md`。
- [x] 9.4 执行 OpenSpec archive，合并 delta。
- [x] 9.5 最终验证后提交 P112。
