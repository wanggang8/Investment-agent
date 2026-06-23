# Tasks: P62 设计系统与可访问性验收

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 `docs/product-experience-polish-roadmap.md`、`docs/development-plan.md`、`docs/frontend-contract.md`、`docs/ui-design.md` 和 P58-P61 相关契约。
- [x] 1.3 使用 Product Design get-context / P57 roadmap / P58-P61 operational cockpit 明确 P62 brief。
- [x] 1.4 创建 `p62-design-system-accessibility-hardening` OpenSpec change。
- [x] 1.5 运行 OpenSpec 校验、diff check 和敏感扫描。
- [x] 1.6 子 agent 方案复审无 Critical / Important 后执行。

## 2. 设计系统 primitives

- [x] 2.1 盘点现有按钮、字段、状态、summary card、详情区、表格、空态和错误态重复模式。
- [x] 2.2 为 `Button` 编写组件测试，覆盖 variant、disabled、type、可访问名称和 working/disabled copy。
- [x] 2.3 实现 `Button`。
- [x] 2.4 为 `Field` 编写组件测试，覆盖 label、hint、error、required 和 input/select/textarea 插槽关联。
- [x] 2.5 实现 `Field`。
- [x] 2.6 为 `StatusBadge` 编写组件测试，覆盖 success/warning/danger/degraded/unknown/readonly/blocked 和文本语义。
- [x] 2.7 实现 `StatusBadge`。
- [x] 2.8 为 `PageHeader`、`SummaryCard`、`DetailSection` 编写组件测试，覆盖标题层级、下一步人工动作、折叠区 `aria-expanded` 和键盘可操作性。
- [x] 2.9 实现 `PageHeader`、`SummaryCard`、`DetailSection`。
- [x] 2.10 为 `ResponsiveTable`、`EmptyState`、`ErrorState` 编写组件测试，覆盖 caption/aria-label、局部滚动、移动 reflow、安全空态和脱敏错误态。
- [x] 2.11 实现 `ResponsiveTable`、`EmptyState`、`ErrorState` 和 `web/src/components/ui/index.ts`。

## 3. 样式与页面接入

- [x] 3.1 在 `web/src/styles/global.css` 中补充 primitives 样式、focus-visible 规则和 390px/768px reflow 约束。
- [x] 3.2 在 Dashboard/Workbench 代表性区域接入 `PageHeader`、`StatusBadge` 或 `SummaryCard`，保持 P58 信息架构不回退。
- [x] 3.3 在 Positions、Risk Alerts、Data Quality 中接入状态、空态、错误态、按钮、字段和响应式表格 primitives。
- [x] 3.4 在 Rules、Audit、Notifications 中接入状态、详情区、summary card、空/错误态和响应式表格 primitives。
- [x] 3.5 在 Local Install、Local Knowledge、Settings 中接入字段、按钮、错误态和诊断详情 primitives。
- [x] 3.6 确认接入后仍只消费现有 service/API，不读取 SQLite、VecLite、localStorage、sessionStorage、本地文件或临时配置。

## 4. 可访问性、键盘与 E2E

- [x] 4.1 更新或新增 Playwright keyboard smoke，覆盖主导航、移动菜单、表单控件、折叠区和关键按钮。
- [x] 4.2 增加 390px、768px、1280px viewport reflow 检查；页面级 `scrollWidth` 不得超过 viewport。
- [x] 4.3 更新 `bash scripts/e2e-smoke.sh` 覆盖 P62 keyboard/reflow 证据，或在 E2E suite 中接入等价检查。
- [x] 4.4 扫描 UI 可见文案和按钮，不出现自动交易、一键交易、代下单、外部推送、短信/邮件/第三方推送承诺、自动确认、自动规则应用、自动修复、覆盖真实库、收益承诺入口。
- [x] 4.5 扫描敏感信息，不渲染完整 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack。

## 5. 验收

- [x] 5.1 运行新增 primitives 组件测试。
- [x] 5.2 运行接入页面定向测试。
- [x] 5.3 运行 `npm --prefix web test`。
- [x] 5.4 运行 `npm --prefix web run build`。
- [x] 5.5 运行 `go test ./...`。
- [x] 5.6 启动真实本地后端和 Vite 前端。
- [x] 5.7 使用浏览器操作代表性页面，验证状态、下一步人工动作、表单、折叠区、表格、空态和错误态。
- [x] 5.8 使用键盘完成主导航、移动菜单、表单控件、折叠区和关键按钮路径。
- [x] 5.9 采集 390px、768px、1280px 截图或等价浏览器证据，并记录 `body.scrollWidth`、`documentElement.scrollWidth` 与 viewport。
- [x] 5.10 运行 `bash scripts/e2e-smoke.sh` 或等价 Playwright 命令。
- [x] 5.11 运行 `openspec validate p62-design-system-accessibility-hardening --strict` 与 `openspec validate --all --strict`。
- [x] 5.12 运行 `git diff --check`。
- [x] 5.13 执行敏感信息和 forbidden copy 扫描。

## 6. 报告、复审与归档

- [x] 6.1 新增 P62 UI 验收报告和截图/浏览器证据资产。
- [x] 6.2 更新 `docs/frontend-contract.md`、`docs/ui-design.md`、`docs/product-experience-polish-roadmap.md`、`docs/development-plan.md` 和进度文档。
- [x] 6.3 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 6.4 执行 OpenSpec archive。
- [x] 6.5 archive 后确认无活跃 change，并推进下一阶段 P63。
- [x] 6.6 提交前子 agent 复审无 Critical / Important。
- [x] 6.7 提交 P62。
