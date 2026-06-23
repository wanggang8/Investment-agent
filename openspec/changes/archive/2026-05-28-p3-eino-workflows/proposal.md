# Proposal: P3 Eino 工作流

## Why

P0-P2 已完成工程骨架、SQLite 数据底座、领域模型与规则裁决。P3 需要在应用层建立工作流上下文、节点框架和核心 Graph，把数据读取、证据核查、分析材料、预期收益、规则裁决、审计记录组织成可测试的执行链路。

## In Scope

- P3.1：实现 `WorkflowContext` 与节点框架。
- P3.1：节点输出必须包含状态、错误码和审计事件片段。
- P3.1：审计事件必须填写 `action`、`node_name`、`node_action`、`status`、`input_ref_type/input_ref`；产生输出时填写 `output_ref_type/output_ref`。
- P3.1：失败和明确降级必须填写 `error_code`，失败节点必须写入 `audit_events`。
- P3.2：实现 DailyDisciplineGraph 与 ConsultationGraph。
- P3.2：每日纪律读取账户、持仓、市场、证据和规则版本。
- P3.2：主动咨询包含能力圈检查。
- P3.2：DeepSeek 节点只写分析报告，不写最终裁决。
- P3.2：ExpectedReturnNode 按 `sample_count` 映射 `precision_status`，并满足三种 DTO 约束。
- P3.2：RuleArbitrationNode 调用领域规则生成最终裁决。
- P3.2：DecisionRecordNode 写 `decision_records`、`evidence_refs`、`audit_events`，并持久化预期收益情景。
- P3.3：实现 EvidenceVerificationGraph、MarketRefreshGraph、EvolutionProposalGraph、GatekeeperAuditGraph。
- P3.3：情报、RAG、source_verifications、market_snapshots、rule_proposals、gatekeeper_audits、rule_versions 等写入规则按计划实现。
- 已实现代码必须补充中文注释，说明工作流、节点、审计和降级边界。

## Out of Scope

- P4-P6：HTTP API handler、前端驾驶舱页面、E2E 验收加固。
- 不新增 development-plan 以外的工作流需求。
- 不改变 P0-P2 已归档契约，除非 P3 delta 明确需要补充工作流约束。
- 不实现自动交易。

## Plan Alignment

本变更与 `docs/development-plan.md` 的 P3.1、P3.2、P3.3 小节一一对应：每个小节的创建文件、任务列表和验收命令均写入 `tasks.md`。
