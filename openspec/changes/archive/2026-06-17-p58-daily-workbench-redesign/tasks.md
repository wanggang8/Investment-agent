# Tasks: P58 今日工作台重构

## 1. 方案与审查

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 `docs/product-experience-polish-roadmap.md` 与 `docs/frontend-contract.md` P57 契约。
- [x] 1.3 使用 Product Design get-context / P55-P56 audit / P57 roadmap 明确 P58 brief。
- [x] 1.4 创建 `p58-daily-workbench-redesign` OpenSpec change。
- [x] 1.5 运行 OpenSpec 校验、diff check 和敏感扫描。
- [x] 1.6 子 agent 方案复审无 Critical / Important 后执行。

## 2. View Model 与组件设计

- [x] 2.1 增加 Dashboard/Workbench view model 测试，覆盖成功、缺数据、降级、高风险、unknown 状态。
- [x] 2.2 实现 null-safe daily workbench view model，不新增 API 字段。
- [x] 2.3 增加 `DailyDecisionHero` 或等价组件，展示今日裁决、状态、可信度和更新时间。
- [x] 2.4 增加 `ManualActionQueue` 或等价组件，展示下一步人工动作，只允许本地导航链接。
- [x] 2.5 增加 `WorkbenchSignalGrid` 或等价组件，展示组合、风险、证据、规则/复盘摘要。

## 3. Dashboard 重构

- [x] 3.1 更新 Dashboard 测试，要求第一屏可见今日裁决、禁止动作、可选人工动作、下一步人工动作和数据可信度。
- [x] 3.2 重构 `DashboardFeature`，把当前三栏 layout 调整为每日 cockpit 主线。
- [x] 3.3 保留今日纪律报告、风险预警、决策详情、审计和账户初始化入口。
- [x] 3.4 确认错误态和 report 加载失败不导致白屏，不泄露内部错误。

## 4. Workbench 重构

- [x] 4.1 更新 Workbench 测试，要求首屏出现同源今日状态和人工动作队列。
- [x] 4.2 重构 `WorkbenchPage`，避免四张卡片并列造成信息割裂。
- [x] 4.3 将组合/风险、规则/复盘、主动咨询入口降为次级区域。
- [x] 4.4 保留所有 P42/P55/P56 路由入口和禁止自动交易文案扫描。

## 5. 样式与移动端

- [x] 5.1 增加或调整 P58 所需 CSS class，复用现有 operational tokens。
- [x] 5.2 390px 下 Dashboard/Workbench 第一屏纵向堆叠，无页面级横向溢出。
- [x] 5.3 桌面 1280px 下摘要、动作队列和信号 grid 可扫描，无文本重叠。

## 6. 验收

- [x] 6.1 运行 `npm test -- --run src/features/dashboard src/pages/WorkbenchPage.test.tsx`。
- [x] 6.2 运行 `npm test`。
- [x] 6.3 运行 `npm run build`。
- [x] 6.4 运行 `go test ./...`。
- [x] 6.5 启动真实本地后端和 Vite 前端。
- [x] 6.6 使用浏览器操作 `/`、`/workbench`，采集桌面和 390px 移动截图。
- [x] 6.7 验证 `/`、`/workbench` 移动端 `body.scrollWidth` 和 `documentElement.scrollWidth` 不超过 viewport。
- [x] 6.8 扫描 UI 文案，确认无自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺入口。
- [x] 6.9 运行 `E2E_BASE_URL=<local vite url> npm run test:e2e`。
- [x] 6.10 运行 `openspec validate p58-daily-workbench-redesign --strict` 与 `openspec validate --all --strict`。
- [x] 6.11 运行 `git diff --check`。
- [x] 6.12 执行敏感信息扫描。

## 7. 报告、复审与归档

- [x] 7.1 新增 P58 UI 验收报告和截图资产。
- [x] 7.2 更新 `docs/frontend-contract.md`、`docs/product-experience-polish-roadmap.md`、`docs/development-plan.md` 和进度文档。
- [x] 7.3 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 7.4 执行 OpenSpec archive。
- [x] 7.5 archive 后确认无活跃 change，并推进下一阶段 P59。
- [ ] 7.6 提交前子 agent 复审无 Critical / Important。
- [ ] 7.7 提交 P58。
