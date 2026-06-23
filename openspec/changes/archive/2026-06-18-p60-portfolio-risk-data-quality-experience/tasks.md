# Tasks: P60 组合、风险与数据质量体验重构

## 1. 方案与审查

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 `docs/product-experience-polish-roadmap.md`、`docs/frontend-contract.md`、P58/P59 归档设计和当前 `/positions`、`/risk-alerts`、`/data-quality` 实现。
- [x] 1.3 使用 Product Design get-context / P55-P56 audit / P57 roadmap / P58-P59 operational cockpit 明确 P60 brief。
- [x] 1.4 创建 `p60-portfolio-risk-data-quality-experience` OpenSpec change。
- [x] 1.5 运行 OpenSpec 校验、diff check 和敏感扫描。
- [x] 1.6 子 agent 方案复审无 Critical / Important 后执行。

## 2. 前端体验模型

- [x] 2.1 新增或调整 portfolio experience model 测试，覆盖初始化、正常维护、高风险比例、导入待确认、错误态和本地-only 动作。
- [x] 2.2 实现 `PortfolioExperienceModel` 或等价映射，不新增 API 字段。
- [x] 2.3 新增或调整 risk disposition model 测试，覆盖 triggered/active/observing/escalated/resolved/archived 队列映射、最高严重程度和下一步动作。
- [x] 2.4 实现 `RiskDispositionModel` 或等价映射，不新增 API 字段。
- [x] 2.5 新增或调整 data quality experience model 测试，覆盖 source health、证据/RAG、LLM、受影响工作流、降级/未知/缺失状态和脱敏显示。
- [x] 2.6 实现 `DataQualityExperienceModel` 或等价映射，不新增 API 字段。
- [x] 2.7 保证所有模型和组件只接收 DTO/props，不直接访问 SQLite、VecLite、localStorage、sessionStorage、本地文件或临时配置。

## 3. Portfolio 体验

- [x] 3.1 更新 `PortfolioPage.test.tsx`，要求首屏展示组合状态、维护阶段、下一步人工动作和 local-only 安全文案。
- [x] 3.2 调整 `/positions` 标题、首屏状态、维护阶段、下一步人工动作和账户快照信息架构。
- [x] 3.3 把初始化/校准、持仓维护、线下交易、批量导入、错误修正分成清晰模式或区块。
- [x] 3.4 保留现有 portfolio service 调用，确保写入动作仍只记录本地事实或审计事实。
- [x] 3.5 确认空态、错误态、导入未通过和缺少持仓时不会提交无效本地事实。

## 4. Risk Alerts 体验

- [x] 4.1 更新 `RiskAlertPage.test.tsx`，要求风险处置总览、队列分组、严重程度、受影响标的、关联链接和安全文案可见。
- [x] 4.2 调整 `/risk-alerts` 信息架构：先展示处置总览，再展示待看/处理中/需复盘/已记录队列。
- [x] 4.3 风险卡片展示触发摘要、禁止动作、建议人工动作、SOP 状态、更新时间、关联决策/报告/通知/审计。
- [x] 4.4 生命周期按钮文案明确为本地 SOP 记录；resolved/archived 不展示生命周期按钮。
- [x] 4.5 确认风险处置不新增交易、外推、自动确认、自动规则应用或组合写入入口。

## 5. Data Quality 体验

- [x] 5.1 更新 `DataQualityPage.test.tsx`，要求质量总览、四类质量信号、下一步本地检查、受影响工作流和安全脱敏可见。
- [x] 5.2 调整 `/data-quality` 信息架构：先展示整体质量状态和下一步，再展示 source health、证据/RAG、LLM、影响范围。
- [x] 5.3 降级、过期、缺失、解析失败、不可用、失败、未知状态必须显示为待检查或降级，不展示为成功。
- [x] 5.4 保留只读行为；页面不得出现自动修复、自动刷新、自动确认、自动规则应用、外部推送、交易或收益承诺入口。
- [x] 5.5 强化脱敏显示，不渲染 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack。

## 6. 样式、移动端与 E2E

- [x] 6.1 增加或调整 P60 所需 CSS class，复用 P58/P59 operational tokens。
- [x] 6.2 390px 下 `/positions`、`/risk-alerts`、`/data-quality` 无页面级横向溢出。
- [x] 6.3 更新 Playwright smoke，覆盖 P60 路由可达、关键状态、队列/质量信号和 forbidden copy scan。

## 7. 验收

- [x] 7.1 运行 `npm test -- --run src/pages/PortfolioPage.test.tsx src/pages/RiskAlertPage.test.tsx src/pages/DataQualityPage.test.tsx`。
- [x] 7.2 运行相关新增 view model 测试。
- [x] 7.3 运行 `npm test`。
- [x] 7.4 运行 `npm run build`。
- [x] 7.5 运行 `go test ./...`。
- [x] 7.6 启动真实本地后端和 Vite 前端。
- [x] 7.7 使用浏览器操作 `/positions`，执行一次本地校准或维护路径，并确认状态、表格和安全文案。
- [x] 7.8 使用浏览器操作 `/risk-alerts`，查看风险队列、关联链接；如本地 fixture 支持，执行一次本地 SOP 生命周期动作。
- [x] 7.9 使用浏览器操作 `/data-quality`，验证质量总览、受影响工作流、本地导航、脱敏和只读边界。
- [x] 7.10 采集 P60 桌面和 390px 移动截图，并验证 `body.scrollWidth` 和 `documentElement.scrollWidth` 不超过 viewport。
- [x] 7.11 扫描 UI 文案，确认无自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺入口。
- [x] 7.12 运行 `bash scripts/e2e-smoke.sh` 或 `E2E_BASE_URL=<local vite url> npm run test:e2e`。
- [x] 7.13 运行 `openspec validate p60-portfolio-risk-data-quality-experience --strict` 与 `openspec validate --all --strict`。
- [x] 7.14 运行 `git diff --check`。
- [x] 7.15 执行敏感信息扫描。

## 8. 报告、复审与归档

- [x] 8.1 新增 P60 UI 验收报告和截图资产。
- [x] 8.2 更新 `docs/frontend-contract.md`、`docs/product-experience-polish-roadmap.md`、`docs/development-plan.md` 和进度文档。
- [x] 8.3 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 8.4 执行 OpenSpec archive。
- [x] 8.5 archive 后确认无活跃 change，并推进下一阶段 P61。
- [x] 8.6 提交前子 agent 复审无 Critical / Important。
- [x] 8.7 提交 P60。
