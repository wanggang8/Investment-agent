## Context

P33 已完成本地账户与持仓事实录入，P34 已完成真实公开数据扩展、结构化 source health、freshness 与失败分类。现有 RuleArbitrationNode 已能输出 `hold`、`reduce`、`sell_only`、`frozen_watch`、`insufficient_data` 等裁决，但缺少一个可查询、可追踪、可解除的风险预警事实层。

P35 的目标是把已有规则裁决和 P34 数据健康信号组织成风险预警中心与 SOP 状态，不改变最终裁决边界，不接交易执行能力。

## Goals / Non-Goals

**Goals:**
- 建立本地风险预警事实模型，记录风险类型、严重程度、SOP 状态、触发依据、禁止动作、建议人工动作、关联 decision/report/notification/audit。
- 支持风险状态从 triggered / active / observing / escalated / resolved / archived 流转，并提供解除、持续观察和升级审计。
- 将风险预警接入每日纪律报告和应用内通知，前端可查看风险中心、触发证据、当前状态与安全边界。
- 复用 P34 source health/freshness，把数据缺失或降级显式展示为风险输入。

**Non-Goals:**
- 不接券商交易 API，不自动下单、撤单、改单或确认成交。
- 不发送邮件、短信、系统 Push、Webhook、WebSocket 等外部推送。
- 不新增登录、付费、授权、Level2 或高频数据源。
- 不让 LLM 覆盖最终规则裁决，不承诺收益，不输出确定性涨跌预测。
- 不在本阶段实现复杂机器学习风险评分或跨账户多用户权限。

## Decisions

1. **风险预警作为独立本地事实表，而不是只放在 `decision_records` JSON 中。**
   - 决策：新增 `risk_alerts` 本地表与 repository/service。
   - 理由：风险需要跨日报告持续追踪、解除和归档，单个 decision JSON 难以表达生命周期。
   - 备选：仅扩展 `decision_records.triggered_rules_json`。放弃原因是无法稳定查询 active 风险和未读通知关联。

2. **SOP 状态由规则与服务确定，LLM 只能提供分析材料。**
   - 决策：RiskAlertService 基于 `final_verdict_status`、triggered rules、market/source health、portfolio ratios 和 evidence status 生成风险类型与 SOP 状态。
   - 理由：风险状态属于纪律执行，不应由模型自由裁量。
   - 备选：让 TrendRiskOfficerNode 输出风险状态。放弃原因是会弱化规则优先边界。

3. **通知与审计复用既有表。**
   - 决策：风险触发、升级、解除、归档写 `audit_events`；active/escalated 风险写应用内 `notifications` 并按 `type/source_type/source_id` 去重。
   - 理由：P21/P31 已建立本地通知与审计边界，无需引入新推送系统。
   - 备选：新增独立风险通知表。放弃原因是重复状态管理。

4. **每日纪律报告保留索引职责，只增加风险摘要和关联链接。**
   - 决策：报告 DTO 聚合 risk alerts，报告表不替代 `risk_alerts` 事实源。
   - 理由：保持 P32 报告只读阅读入口语义。

5. **P35 先覆盖确定性规则风险，不扩展真实数据源。**
   - 决策：风险输入来自现有 portfolio、position、market snapshot、decision、evidence/source health。
   - 理由：P34 已提供数据健康基础，P35 聚焦编排和展示。

## Risks / Trade-offs

- [Risk] 风险类型和裁决状态重复表达。→ Mitigation：`decision_records` 保存单次裁决，`risk_alerts` 保存跨次生命周期，字段命名区分 `verdict_status` 与 `sop_status`。
- [Risk] 风险解除条件过度自动化。→ Mitigation：解除只改变本地预警状态，不改变持仓、确认、规则版本或账户事实；必要时要求人工确认。
- [Risk] 通知刷屏。→ Mitigation：复用 notifications 未读去重约束，同一 active 风险重复触发只刷新通知。
- [Risk] 数据缺失被误解为风险已消失。→ Mitigation：source health 缺失或 stale 时生成/保留 `insufficient_evidence` 或 `data_degraded` 类型风险，不自动解除。

## Migration Plan

1. 新增 `risk_alerts` migration，保持向后兼容。
2. 新增 repository/service 和 API；旧页面不依赖风险表时仍可运行。
3. DailyDisciplineGraph 或手动任务成功后调用风险编排服务，写入本地事实、通知和审计。
4. 前端风险页面在无数据时展示空状态，不影响既有报告与驾驶舱。
5. 回滚时可停止调用风险服务，既有决策、账户、通知和审计表不受影响。

## Open Questions

- 首版风险解除是否允许用户手动解除，或只由下一次每日纪律自动评估解除；实现时优先支持显式手动解除并写审计。
- 风险等级是否分三档 warning/critical/info，还是与 notification severity 完全一致；首版建议复用 notification severity。
