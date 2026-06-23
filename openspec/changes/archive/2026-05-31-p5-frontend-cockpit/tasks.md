# Tasks: p5-frontend-cockpit

> 对齐 `docs/development-plan.md` P5：前端驾驶舱。实现代码必须写必要中文注释，说明复杂状态映射、确认边界、审计语义和禁止自动交易等业务意图；不得写只复述代码的注释。

## 1. 架构治理准备（P5.0 基线确认）

参考文档：`docs/architecture.md`、`docs/data-model.md`、`docs/workflow.md`、`docs/frontend-contract.md`。

说明：`docs/development-plan.md` 中 P5.0 已标记完成，本节作为进入 P5.1 前的基线确认，不作为新增开发任务。

确认项：

- [x] 1.1 确认应用层和工作流包只依赖仓储接口与事务协调接口，不依赖 SQLite 具体实现
- [x] 1.2 确认 HTTP handler 保持请求解析、用例调用和响应写入职责，不直接访问数据库或管理 SQL 事务
- [x] 1.3 确认跨表业务事实和对应 `audit_events` 使用统一事务协调路径
- [x] 1.4 确认关键业务 ID、时间和契约枚举复用共享实现，并支持确定性测试
- [x] 1.5 确认前端已按 feature 和 shared 目录准备 P5 页面扩展结构
- [x] 1.6 验收：

```bash
go test ./...
```

- [x] 1.7 验收：

```bash
cd web
npm run build
```

## 2. 类型与 API client（P5.1）

参考文档：`docs/frontend-contract.md`、`docs/api.md`。

创建或确认以下文件：

```text
web/src/types/api.ts
web/src/types/dashboard.ts
web/src/types/portfolio.ts
web/src/types/decision.ts
web/src/types/evidence.ts
web/src/types/rule.ts
web/src/types/audit.ts
web/src/types/settings.ts
web/src/types/market.ts
web/src/types/review.ts
web/src/services/client.ts
web/src/services/dashboard.ts
web/src/services/portfolio.ts
web/src/services/decision.ts
web/src/services/evidence.ts
web/src/services/rule.ts
web/src/services/audit.ts
web/src/services/settings.ts
web/src/services/market.ts
web/src/services/review.ts
```

任务：

- [x] 2.1 创建或确认 `web/src/types/api.ts`，定义通用响应类型
- [x] 2.2 创建或确认 `web/src/types/dashboard.ts`，定义驾驶舱 DTO
- [x] 2.3 创建或确认 `web/src/types/portfolio.ts`，定义持仓 DTO
- [x] 2.4 创建或确认 `web/src/types/decision.ts`，定义决策 DTO
- [x] 2.5 创建或确认 `web/src/types/evidence.ts`，定义证据 DTO
- [x] 2.6 创建或确认 `web/src/types/rule.ts`，定义规则 DTO
- [x] 2.7 创建或确认 `web/src/types/audit.ts`，定义审计 DTO
- [x] 2.8 创建或确认 `web/src/types/settings.ts`，定义设置 DTO
- [x] 2.9 创建或确认 `web/src/types/market.ts`，定义市场 DTO
- [x] 2.10 创建或确认 `web/src/types/review.ts`，定义复盘 DTO
- [x] 2.11 创建或确认 `web/src/services/client.ts`，统一处理 `request_id`、`data`、`meta`、`error`
- [x] 2.12 创建或确认 dashboard、portfolio、decision、evidence、rule、audit、settings、market、review service 文件
- [x] 2.13 统一处理 409、500、503 错误，并映射为前端可展示状态
- [x] 2.14 确认前端不访问 SQLite、VecLite、本地文件
- [x] 2.15 为复杂错误映射和数据源边界补充必要中文注释
- [x] 2.16 验收：

```bash
cd web
npm run build
```

## 3. Agent 决策驾驶舱（P5.2）

参考文档：`docs/ui-design.md`、`docs/ui-flow.md`。

创建或确认以下文件：

```text
web/src/pages/DashboardPage.tsx
web/src/components/layout/CockpitLayout.tsx
web/src/components/dashboard/DisciplineStatus.tsx
web/src/components/dashboard/PortfolioSummary.tsx
web/src/components/dashboard/TriggeredRules.tsx
web/src/components/dashboard/EvidenceSummary.tsx
web/src/components/dashboard/FinalVerdictCard.tsx
web/src/components/dashboard/UserConfirmationPanel.tsx
```

任务：

- [x] 3.1 创建或确认 `web/src/pages/DashboardPage.tsx`
- [x] 3.2 创建或确认 `web/src/components/layout/CockpitLayout.tsx`
- [x] 3.3 创建或确认 `DisciplineStatus.tsx`
- [x] 3.4 创建或确认 `PortfolioSummary.tsx`
- [x] 3.5 创建或确认 `TriggeredRules.tsx`
- [x] 3.6 创建或确认 `EvidenceSummary.tsx`
- [x] 3.7 创建或确认 `FinalVerdictCard.tsx`
- [x] 3.8 创建或确认 `UserConfirmationPanel.tsx`
- [x] 3.9 实现三栏驾驶舱布局
- [x] 3.10 首屏展示纪律状态、风险红线、今日建议、账户摘要、证据摘要
- [x] 3.11 信息不足状态展示缺失项和暂停原因
- [x] 3.12 冻结观察状态展示等待条件
- [x] 3.13 用户确认区只允许记录计划、已手动执行、待观察、标记错误
- [x] 3.14 页面不得出现自动交易入口
- [x] 3.15 为确认动作边界、信息不足和冻结观察状态补充必要中文注释
- [x] 3.16 验收：

```bash
cd web
npm run build
```

## 4. 决策详情、证据、规则与审计页面（P5.3）

创建或确认以下文件：

```text
web/src/pages/DecisionDetailPage.tsx
web/src/pages/EvidencePage.tsx
web/src/pages/RulesPage.tsx
web/src/pages/AuditPage.tsx
web/src/pages/PortfolioPage.tsx
web/src/pages/SettingsPage.tsx
web/src/pages/ReviewSummaryPage.tsx
web/src/components/decision/DecisionTrace.tsx
web/src/components/evidence/EvidenceTable.tsx
web/src/components/rules/RuleProposalPanel.tsx
web/src/components/audit/AuditEventTimeline.tsx
web/src/components/portfolio/PortfolioTable.tsx
web/src/components/settings/CapabilitySettingsPanel.tsx
web/src/components/review/ReviewSummaryPanel.tsx
```

任务：

- [x] 4.1 创建或确认 `DecisionDetailPage.tsx`
- [x] 4.2 创建或确认 `EvidencePage.tsx`
- [x] 4.3 创建或确认 `RulesPage.tsx`
- [x] 4.4 创建或确认 `AuditPage.tsx`
- [x] 4.5 创建或确认 `PortfolioPage.tsx`
- [x] 4.6 创建或确认 `SettingsPage.tsx`
- [x] 4.7 创建或确认 `ReviewSummaryPage.tsx`
- [x] 4.8 创建或确认 `DecisionTrace.tsx`
- [x] 4.9 创建或确认 `EvidenceTable.tsx`
- [x] 4.10 创建或确认 `RuleProposalPanel.tsx`
- [x] 4.11 创建或确认 `AuditEventTimeline.tsx`
- [x] 4.12 创建或确认 `PortfolioTable.tsx`
- [x] 4.13 创建或确认 `CapabilitySettingsPanel.tsx`
- [x] 4.14 创建或确认 `ReviewSummaryPanel.tsx`
- [x] 4.15 决策详情按 `docs/ui-flow.md` 第 6 节展示
- [x] 4.16 持仓页使用 `GET /api/v1/portfolio/current`，不直接访问 SQLite 或 VecLite
- [x] 4.17 证据页展示 `source_level`、`evidence_role`、`verification_status`
- [x] 4.18 规则提案页展示 `pending_final_confirm` 状态和最终确认动作
- [x] 4.19 设置页展示能力圈配置、系统状态、市场快照状态、通知配置和索引状态，不展示完整密钥
- [x] 4.20 复盘页展示建议数量、确认动作、错误案例、规则提案和审计事件汇总
- [x] 4.21 审计页区分 `action`、`node_name` 与 `node_action`，并按 `status`、`error_code`、输入引用、输出引用展示审计详情
- [x] 4.22 为审计字段语义、规则最终确认和密钥展示边界补充必要中文注释
- [x] 4.23 验收：

```bash
cd web
npm run build
```

## 5. Plan alignment、实现完成与归档检查

- [x] 5.1 确认 `proposal.md` in scope / out of scope 与 `docs/development-plan.md` P5 一致
- [x] 5.2 确认 `specs/frontend-cockpit/spec.md` 只包含 P5 阶段 delta
- [x] 5.3 确认 P5.0 对应任务 1.1–1.7
- [x] 5.4 确认 P5.1 对应任务 2.1–2.16
- [x] 5.5 确认 P5.2 对应任务 3.1–3.16
- [x] 5.6 确认 P5.3 对应任务 4.1–4.23
- [x] 5.7 实现完成后更新 `docs/development-plan.md` P5 任务状态
- [x] 5.8 archive 前确认 `docs/GOVERNANCE.md` 活跃变更表仍指向 `p5-frontend-cockpit`
- [x] 5.9 archive 后从 `docs/GOVERNANCE.md` 活跃变更表移除本 change
- [x] 5.10 archive 后更新 `openspec/PROGRESS.md`：P5 标记为 done，下一阶段指向 P6

## Plan alignment

- P5.0 架构治理准备对应任务：1.1–1.7。
- P5.1 类型与 API client 对应任务：2.1–2.16。
- P5.2 Agent 决策驾驶舱对应任务：3.1–3.16。
- P5.3 决策详情、证据、规则与审计页面对应任务：4.1–4.23。
- P5 归档相关检查对应任务：5.1–5.10。
- 与 `docs/development-plan.md` P5 小节、创建文件、任务列表和验收命令一一对应；仅额外加入中文注释要求，来自本次用户指令。
