# Tasks: P113 布局高保真精修

## 1. Governance And Baseline

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 创建并校验 `p113-layout-fidelity-polish` OpenSpec change。
- [x] 1.3 更新 `docs/GOVERNANCE.md` 与 `openspec/PROGRESS.md`，将 P113 标记为当前活跃 change。
- [x] 1.4 确认 P113 仍以 P111 第二方案参考图为视觉真源，且不扩展后端或投资能力。

## 2. Layout Audit Baseline

- [x] 2.1 复核 P112 后 fresh audit findings，列出必须修复的布局问题。
- [x] 2.2 定位共享 CSS 与页面结构中的横向溢出、压缩、错位、触控高度不足和 raw 内容暴露根因。
- [x] 2.3 建立 P113 mismatch ledger，覆盖页面、视口、问题等级、修复策略和复验状态。

## 3. Shared Layout Fixes

- [x] 3.1 修复移动端 `daily-hero` / metric grid 横向溢出与裁切。
- [x] 3.2 修复桌面 report hero 多卡片压缩、文字截断和首屏节奏。
- [x] 3.3 提升移动端 action link / compact button 触控高度与可辨识度。
- [x] 3.4 补强文本换行、min-width、grid constraints，避免长字段导致错位。

## 4. Page-Level Polish

- [x] 4.1 修复 `/data-quality` 桌面与移动布局。
- [x] 4.2 修复 `/settings`、`/local-install`、`/local-knowledge` 桌面与移动布局。
- [x] 4.3 修复 `/decisions/:id` 决策详情层级与 report/ledger composition。
- [x] 4.4 修复 `/rules`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run` 的移动触控、内容密度和工程化内容暴露。
- [x] 4.5 巡检其余路由，修复明显错位、重叠、裁切或不精致的布局问题。

## 5. Rendered QA And Review Loop

- [x] 5.1 运行前端测试与构建。
- [x] 5.2 启动本地 backend/frontend，采集全 18 个桌面路由截图。
- [x] 5.3 采集全 18 个 390px 移动路由截图。
- [x] 5.4 对每个页面执行 no-overflow、clip/overlap、touch target、console health 和参考图对比检查。
- [x] 5.5 若发现 P0/P1/P2 布局问题，继续修复并重新截图，直到关闭。

## 6. Validation

- [x] 6.1 运行 `npm --prefix web test -- --run`。
- [x] 6.2 运行 `npm --prefix web run build`。
- [x] 6.3 运行 forbidden affordance scan。
- [x] 6.4 运行 sensitive/redaction scan。
- [x] 6.5 运行 `openspec validate p113-layout-fidelity-polish --strict`。
- [x] 6.6 运行 `openspec validate --all --strict`。
- [x] 6.7 运行 `git diff --check`。

## 7. Documentation, Archive, Commit

- [x] 7.1 新增 P113 acceptance record，包含截图目录、page matrix、mismatch ledger、复审结论和命令结果。
- [x] 7.2 更新相关 UI / frontend contract / roadmap 文档中的 P113 结果。
- [x] 7.3 更新 `docs/GOVERNANCE.md` 与 `openspec/PROGRESS.md`。
- [x] 7.4 执行 OpenSpec archive，合并 delta。
- [x] 7.5 最终验证后提交 P113。
