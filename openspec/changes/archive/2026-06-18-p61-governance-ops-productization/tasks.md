# Tasks: P61 治理和运维页面产品化

## 1. 方案与审查

- [x] 1.1 阅读 `AGENTS.md`、`docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 `docs/product-experience-polish-roadmap.md`、`docs/development-plan.md`、`docs/frontend-contract.md` 和 P58-P60 归档材料。
- [x] 1.3 使用 Product Design get-context / P57 roadmap / P58-P60 operational cockpit 明确 P61 brief。
- [x] 1.4 阅读当前 `/rules`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run`、`/local-install`、`/local-knowledge`、`/settings` 页面和测试。
- [x] 1.5 创建 `p61-governance-ops-productization` OpenSpec change。
- [x] 1.6 运行 OpenSpec 校验、diff check 和敏感扫描。
- [x] 1.7 子 agent 方案复审无 Critical / Important 后执行。

## 2. 前端治理/运维 view model

- [x] 2.1 新增 `rulesGovernanceModel` 测试，覆盖当前规则、提案计数、待确认、最终确认、门禁风险、下一步人工动作和 forbidden copy。
- [x] 2.2 实现 `rulesGovernanceModel`，只消费现有 rule DTO。
- [x] 2.3 新增 `auditOpsModel` 测试，覆盖事件计数、类型分布、最近事件、空态、错误态和安全摘要。
- [x] 2.4 实现 `auditOpsModel`，只消费现有 audit DTO。
- [x] 2.5 新增 `notificationInboxModel` 测试，覆盖未读、严重程度、来源分类、处理状态、本地 inbox 边界和 forbidden copy。
- [x] 2.6 实现 `notificationInboxModel`，只消费现有 notification DTO。
- [x] 2.7 新增 `dailyOpsModel` 测试，覆盖 daily reports 与 daily auto run 的状态、证据/执行覆盖、缺失前提、下一步人工动作和安全文案。
- [x] 2.8 实现 `dailyOpsModel`，只消费现有 daily report / auto-run DTO。
- [x] 2.9 新增 `localOpsModel` 测试，覆盖 local install、local knowledge、settings 的配置/诊断/脱敏/下一步动作模型。
- [x] 2.10 实现 `localOpsModel`，不访问 SQLite、VecLite、localStorage、sessionStorage、本地文件或临时配置。

## 3. Rules / Audit / Notifications 体验

- [x] 3.1 更新 `RulesPage.test.tsx`，要求首屏展示规则治理状态、提案分组、下一步人工动作和安全边界。
- [x] 3.2 调整 `/rules` 信息架构，减少首屏 raw JSON，突出提案理由、样本、过拟合、守门人、审计和人工确认。
- [x] 3.3 保留现有 `confirmRuleProposal` / `finalConfirmRuleProposal` 调用，按钮文案明确为人工本地规则治理动作。
- [x] 3.4 更新 `AuditPage.test.tsx`，要求首屏展示审计摘要、事件分类、最近活动、空态/错误态和时间线。
- [x] 3.5 调整 `/audit` 信息架构，把 raw event list 包装为摘要 + 可扫描时间线。
- [x] 3.6 更新 `NotificationPage.test.tsx`，要求本地通知收件箱、严重程度/来源分布、未读处理状态和本地-only 边界。
- [x] 3.7 调整 `/notifications` 信息架构，保留轮询和标记已读能力，但不暗示外部推送。

## 4. Daily Reports / Daily Auto Run 体验

- [x] 4.1 更新 `DailyDisciplineReportsPage.test.tsx`，要求纪律复盘总览、证据覆盖、趋势、缺口和报告入口。
- [x] 4.2 调整 `/daily-discipline/reports` 信息架构，使报告历史像复盘列表而非数据 dump。
- [x] 4.3 更新 `DailyAutoRunPage.test.tsx`，要求运行健康总览、计划/最近/下次执行、失败诊断、缺失前提、关联入口和安全边界。
- [x] 4.4 调整 `/daily-auto-run` 信息架构，区分 disabled/scheduled/running/success/degraded/failed/unknown，不把未知或降级显示为成功。
- [x] 4.5 确认 Daily Auto Run 不出现自动修复、自动确认、自动规则应用、覆盖真实库、后台交易或收益承诺入口。

## 5. Local Install / Local Knowledge / Settings 体验

- [x] 5.1 更新 `LocalInstallPage.test.tsx`，要求配置草稿、诊断摘要、失败步骤、下一步人工复验和脱敏边界。
- [x] 5.2 调整 `/local-install` 信息架构，统一配置、命令、诊断摘要和安全说明。
- [x] 5.3 更新 `LocalKnowledgePage.test.tsx`，要求导入草稿、脱敏预览、索引计划、确认理由、本地事实写入边界和阻断态。
- [x] 5.4 调整 `/local-knowledge` 信息架构，保持 validate -> preview -> explicit confirm，不新增自动索引成功承诺。
- [x] 5.5 更新 `SettingsPage.test.tsx`，要求能力圈、系统状态、数据源健康、市场刷新、错误摘要和本地配置边界。
- [x] 5.6 调整 `/settings` 信息架构，统一配置/诊断/数据源健康的状态表达和下一步动作。
- [x] 5.7 强化脱敏显示，不渲染 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack。

## 6. 样式、移动端与 E2E

- [x] 6.1 增加或调整 P61 所需 CSS class，复用 P58-P60 operational tokens。
- [x] 6.2 390px 下 P61 页面无页面级横向溢出；表格和长 JSON 仅在局部容器滚动。
- [x] 6.3 更新 Playwright smoke，覆盖 P61 路由可达、关键状态、下一步动作和 forbidden copy scan。

## 7. 验收

- [x] 7.1 运行 P61 新增 view model 测试。
- [x] 7.2 运行 P61 页面定向测试。
- [x] 7.3 运行 `npm --prefix web test`。
- [x] 7.4 运行 `npm --prefix web run build`。
- [x] 7.5 运行 `go test ./...`。
- [x] 7.6 启动真实本地后端和 Vite 前端。
- [x] 7.7 使用浏览器操作 `/rules`，验证规则治理状态、提案、确认边界和 forbidden copy。
- [x] 7.8 使用浏览器操作 `/audit`，验证审计摘要、时间线、空态/错误态和关联入口。
- [x] 7.9 使用浏览器操作 `/notifications`，验证本地 inbox；如真实数据支持，执行一次标记已读。
- [x] 7.10 使用浏览器操作 `/daily-discipline/reports` 和 `/daily-auto-run`，验证报告复盘和运行诊断。
- [x] 7.11 使用浏览器操作 `/local-install`，上传诊断摘要 fixture 并验证脱敏摘要。
- [x] 7.12 使用浏览器操作 `/local-knowledge`，执行 validate；confirm 仅在测试数据安全时执行。
- [x] 7.13 使用浏览器操作 `/settings`，验证系统状态、数据源健康和市场刷新本地边界。
- [x] 7.14 采集 P61 桌面和 390px 移动截图，并验证 `body.scrollWidth` 和 `documentElement.scrollWidth` 不超过 viewport。
- [x] 7.15 扫描 UI 文案，确认无自动交易、一键交易、代下单、外部推送、短信/邮件/第三方通知承诺、自动确认、自动规则应用、自动修复、覆盖真实库、收益承诺入口。
- [x] 7.16 运行 `bash scripts/e2e-smoke.sh` 或 `E2E_BASE_URL=<local vite url> npm run test:e2e`。
- [x] 7.17 运行 `openspec validate p61-governance-ops-productization --strict` 与 `openspec validate --all --strict`。
- [x] 7.18 运行 `git diff --check`。
- [x] 7.19 执行敏感信息扫描。

## 8. 报告、复审与归档

- [x] 8.1 新增 P61 UI 验收报告和截图资产。
- [x] 8.2 更新 `docs/frontend-contract.md`、`docs/product-experience-polish-roadmap.md`、`docs/development-plan.md` 和进度文档。
- [x] 8.3 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 8.4 执行 OpenSpec archive。
- [x] 8.5 archive 后确认无活跃 change，并推进下一阶段 P62。
- [x] 8.6 提交前子 agent 复审无 Critical / Important。
- [x] 8.7 提交 P61。
