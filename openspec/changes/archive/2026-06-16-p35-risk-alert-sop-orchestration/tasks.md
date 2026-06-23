## 1. OpenSpec 与范围

- [x] 1.1 确认 P35 只覆盖风险预警、SOP 状态流转、通知/审计/每日纪律报告联动和前端展示。
- [x] 1.2 确认 P35 不接券商 API、不自动交易、不外部推送、不登录/付费/授权/Level2/高频源、不承诺收益、不预测确定涨跌。
- [x] 1.3 对齐 P33 本地账户事实、P34 source health/freshness、`docs/requirements.md` 风险/SOP 条款、`docs/api.md`、`docs/data-model.md`、`docs/workflow.md`、`docs/frontend-contract.md` 的现有契约。

## 2. 数据模型与领域状态

- [x] 2.1 定义 `risk_alerts` 本地事实模型：alert ID、risk type、severity、SOP status、symbol、trigger summary、trigger context JSON、prohibited actions、suggested actions、related IDs、timestamps。
- [x] 2.2 增加 SQLite migration、domain model、repository interface 和 sqlite repository tests。
- [x] 2.3 定义风险类型：valuation_high、buy_thesis_broken、liquidity_danger、sentiment_extreme、position_limit_breach、insufficient_evidence、data_degraded。
- [x] 2.4 定义 SOP 状态：triggered、active、observing、escalated、resolved、archived，并校验合法流转。
- [x] 2.5 确认风险事实写入不会更新 positions、portfolio_snapshots、operation_confirmations、position_transactions、rule_versions、broker state、orders 或 external notifications。

## 3. 风险编排服务

- [x] 3.1 实现 RiskAlertService，从 decision verdict、triggered rules、portfolio/position、market/source health、evidence status 生成风险预警草稿。
- [x] 3.2 实现幂等 upsert：同一 active risk type + symbol + source scope 不生成冲突重复 active alert。
- [x] 3.3 实现生命周期动作：continue observing、escalate、resolve、archive，写入 audit_events。
- [x] 3.4 将 active/escalated 风险写入本地 notifications，并复用未读通知去重规则。
- [x] 3.5 增加服务 tests，覆盖触发、重复触发、升级、解除、归档、通知去重和非交易边界。

## 4. 工作流与每日纪律接入

- [x] 4.1 在 DailyDisciplineGraph 或其应用层聚合后调用风险编排服务，关联 decision_id、report_id、request_id。
- [x] 4.2 将 P34 source health/freshness 纳入 data_degraded / insufficient_evidence 风险触发上下文。
- [x] 4.3 将风险摘要接入每日纪律报告 DTO 与 query service，保留 risk links、SOP status、prohibited actions 和 suggested actions。
- [x] 4.4 增加 workflow/service tests，覆盖估值高位、买入逻辑破坏、流动性危险、情绪极端、仓位超限、证据不足和数据降级。

## 5. HTTP API 与前端

- [x] 5.1 新增风险预警 API：列表、详情、生命周期动作；统一响应信封、错误码和只读安全文案。
- [x] 5.2 更新 app routing/handler wiring，并增加 handler tests。
- [x] 5.3 新增前端 risk alert types、services、status mappers 和路由。
- [x] 5.4 新增风险预警中心页面，展示风险类型、严重程度、SOP 状态、触发依据、禁止动作、建议人工动作、关联报告/决策/通知/审计。
- [x] 5.5 更新 Dashboard / Daily Discipline report detail / Notification 相关入口，展示风险摘要或链接。
- [x] 5.6 增加前端 tests，覆盖 active/escalated/resolved/archived 状态、空状态、禁止自动交易文案和生命周期动作。

## 6. 文档与验收

- [x] 6.1 在 P35 delta 中记录待归档合并到 `docs/api.md` 的 risk alert API/DTO/错误分类和事务边界。
- [x] 6.2 在 P35 delta 中记录待归档合并到 `docs/data-model.md` 的 risk_alerts 模型、状态枚举、索引和非交易约束。
- [x] 6.3 在 P35 delta 中记录待归档合并到 `docs/workflow.md` 与 `docs/frontend-contract.md` 的风险编排、每日纪律接入和前端展示。
- [x] 6.4 更新 `docs/development-plan.md`、`openspec/PROGRESS.md`、`AGENTS.md`、`docs/GOVERNANCE.md` 的 P35 active 状态。
- [x] 6.5 运行 `go test ./...`。
- [x] 6.6 运行 `npm --prefix web test -- --run`。
- [x] 6.7 运行 `npm --prefix web run build`。
- [x] 6.8 运行 P35 风险场景 smoke，覆盖 active/escalated/resolved/archived 和通知/审计写入。
- [x] 6.9 运行 `openspec validate p35-risk-alert-sop-orchestration --strict`。
- [x] 6.10 运行 `openspec validate --all --strict`。
- [x] 6.11 运行 `git status --short`，确认只包含预期修改且无临时产物。
