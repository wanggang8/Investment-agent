# Tasks: P59 决策解释体验重构

## 1. 方案与审查

- [x] 1.1 阅读 `docs/GOVERNANCE.md`、`openspec/project.md`、`openspec/PROGRESS.md`。
- [x] 1.2 阅读 `docs/product-experience-polish-roadmap.md` 与 `docs/frontend-contract.md` 的决策详情、证据、决策闭环和 P57-P58 契约。
- [x] 1.3 使用 Product Design get-context / P55-P56 audit / P57 roadmap / P58 operational cockpit 明确 P59 brief。
- [x] 1.4 创建 `p59-decision-explainability-experience` OpenSpec change。
- [x] 1.5 运行 OpenSpec 校验、diff check 和敏感扫描。
- [x] 1.6 子 agent 方案复审无 Critical / Important 后执行。

## 2. 决策解释模型与组件

- [x] 2.1 增加 decision explanation view model 测试，覆盖成功、nullable/missing、degraded/unknown、真实 LLM-like DTO。
- [x] 2.2 实现 null-safe `DecisionExplanationViewModel` 或等价映射，不新增 API 字段。
- [x] 2.3 新增或调整 `DecisionStoryHero` / `DecisionSafetyPanel` / `DecisionWhyPanel` / `DecisionTrustPanel` 等小范围组件。
- [x] 2.4 保证所有组件只接收 props，不直接调用 API、SQLite、VecLite、localStorage、sessionStorage 或本地文件。

## 3. Consultation 体验

- [x] 3.1 更新 `DecisionDetailPage.test.tsx` 的 `/consultation` 覆盖，要求输入假设、生成建议、解释路径和安全文案可见。
- [x] 3.2 调整 `/consultation` 标题、说明、表单层级和成功结果摘要。
- [x] 3.3 真实 LLM consultation 成功后提供打开新决策详情的本地导航；失败时展示可恢复错误和不生成交易建议的安全空态。
- [x] 3.4 确认 consultation 不提供自动确认、自动交易、外部推送、自动规则应用或收益承诺入口。

## 4. Decision Detail 体验

- [x] 4.1 更新决策详情测试，要求首屏展示最终裁决、禁止动作、可选人工动作、数据可信度、关键原因和安全边界。
- [x] 4.2 重构 `DecisionDetailPage` / `DecisionTrace`，将长 trace 改为 story-first + layered details。
- [x] 4.3 Evidence、LLM、rules、expected return、audit、confirmation 分层展示；长列表默认折叠或降为二级详情。
- [x] 4.4 Nullable/missing DTO 使用安全空态，不把缺字段解释为允许执行或成功。

## 5. Evidence 与 Decision Loop 体验

- [x] 5.1 更新 `EvidencePage.test.tsx`，覆盖证据可信度概览、来源等级说明、决策解释链接、空态/错误态和安全文案。
- [x] 5.2 调整 `/evidence` 信息架构：可信度概览优先，证据表格保留筛选与展开。
- [x] 5.3 更新 `DecisionLoopPage.test.tsx`，覆盖只读时间线、缺口说明、本地链接和无写入动作按钮。
- [x] 5.4 调整 `/decision-loop` 信息架构：先讲闭环状态，再讲阶段、缺口和链接。

## 6. 样式、移动端与 E2E

- [x] 6.1 增加或调整 P59 所需 CSS class，复用现有 operational tokens。
- [x] 6.2 390px 下 `/consultation`、`/decisions/:decisionId`、`/evidence`、`/decision-loop` 无页面级横向溢出。
- [x] 6.3 更新 Playwright smoke，覆盖 P59 路由可达、解释链接和 forbidden copy scan。

## 7. 验收

- [x] 7.1 运行 `npm test -- --run src/pages/DecisionDetailPage.test.tsx src/pages/EvidencePage.test.tsx src/pages/DecisionLoopPage.test.tsx`。
- [x] 7.2 运行 `npm test`。
- [x] 7.3 运行 `npm run build`。
- [x] 7.4 运行 `go test ./...`。
- [x] 7.5 启动真实本地后端和 Vite 前端。
- [x] 7.6 使用浏览器操作 `/consultation`，执行真实 LLM consultation，并打开生成后的 `/decisions/:decisionId`。
- [x] 7.7 使用浏览器操作 `/evidence`、`/decision-loop`，验证链接、折叠、错误/空态和只读边界。
- [x] 7.8 采集 P59 桌面和 390px 移动截图，并验证 `body.scrollWidth` 和 `documentElement.scrollWidth` 不超过 viewport。
- [x] 7.9 扫描 UI 文案，确认无自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、收益承诺入口。
- [x] 7.10 运行 `E2E_BASE_URL=<local vite url> npm run test:e2e`。
- [x] 7.11 运行 `openspec validate p59-decision-explainability-experience --strict` 与 `openspec validate --all --strict`。
- [x] 7.12 运行 `git diff --check`。
- [x] 7.13 执行敏感信息扫描。

## 8. 报告、复审与归档

- [x] 8.1 新增 P59 UI 验收报告和截图资产。
- [x] 8.2 更新 `docs/frontend-contract.md`、`docs/product-experience-polish-roadmap.md`、`docs/development-plan.md` 和进度文档。
- [x] 8.3 子 agent 执行后复审无 Critical / Important 后归档。
- [x] 8.4 执行 OpenSpec archive。
- [x] 8.5 archive 后确认无活跃 change，并推进下一阶段 P60。
- [x] 8.6 提交前子 agent 复审无 Critical / Important。
- [x] 8.7 提交 P59。
