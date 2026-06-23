# Design: P3 Eino 工作流

## 1. 应用层工作流框架

目标文件：

```text
internal/application/workflow/context.go
internal/application/workflow/node.go
internal/application/workflow/audit_writer.go
internal/application/workflow/errors.go
```

设计要点：

- `WorkflowContext` 对齐 `docs/workflow.md` 第 4 节，作为应用层节点间传递的统一上下文。
- 节点接口统一返回节点状态、错误码和审计事件片段。
- `audit_writer.go` 负责把节点审计片段转换为 `audit_events` 写入。
- `errors.go` 统一定义 `DATA_REQUIRED`、`DATA_STALE`、`RULE_VERSION_MISSING`、`EVIDENCE_NOT_FOUND`、`SOURCE_VERIFICATION_FAILED`、`VECTOR_INDEX_UNAVAILABLE`、`ANALYST_UNAVAILABLE`、`DECISION_RECORD_FAILED` 等错误码。
- 所有新代码写中文注释，说明节点职责、输入输出和降级行为。

## 2. DailyDisciplineGraph 与 ConsultationGraph

目标文件：

```text
internal/application/workflow/daily_discipline_graph.go
internal/application/workflow/consultation_graph.go
internal/application/workflow/nodes/state_snapshot_node.go
internal/application/workflow/nodes/capability_check_node.go
internal/application/workflow/nodes/evidence_retrieval_node.go
internal/application/workflow/nodes/value_analyst_node.go
internal/application/workflow/nodes/trend_risk_officer_node.go
internal/application/workflow/nodes/expected_return_node.go
internal/application/workflow/nodes/rule_arbitration_node.go
internal/application/workflow/nodes/decision_record_node.go
```

设计要点：

- DailyDisciplineGraph 按 `StateSnapshotNode -> EvidenceRetrievalNode -> ValueAnalystNode -> TrendRiskOfficerNode -> ExpectedReturnNode -> RuleArbitrationNode -> DecisionRecordNode` 执行。
- ConsultationGraph 在 StateSnapshotNode 后增加 CapabilityCheckNode。
- DeepSeek 相关节点只产出分析报告；最终裁决由 `domain/rule` 生成。
- ExpectedReturnNode 按样本数映射：`>=20 available`、`5~19 insufficient`、`<5 unavailable`。
- DecisionRecordNode 以事务写入 `decision_records`、`evidence_refs` 和 `audit_events`，并保存 `expected_return_scenarios_json`。

## 3. Evidence、Market、Evolution、Gatekeeper 工作流

目标文件：

```text
internal/application/workflow/evidence_verification_graph.go
internal/application/workflow/market_refresh_graph.go
internal/application/workflow/evolution_proposal_graph.go
internal/application/workflow/gatekeeper_audit_graph.go
```

设计要点：

- EvidenceVerificationGraph 写入 `intelligence_items`、`intelligence_summary`、`rag_chunks` 和 `source_verifications`。
- VecLite 索引只作为可重建索引，事实数据以 SQLite 为准。
- MarketRefreshGraph 独立处理市场数据刷新，写入 `market_snapshots` 与 `audit_events`。
- EvolutionProposalGraph 从错误案例生成规则提案，但不得修改正式规则。
- GatekeeperAuditGraph 只写 `gatekeeper_audits`；审计通过后进入 `pending_final_confirm`，最终确认后才写 `rule_versions`。

## 4. 测试策略

- P3.1：`go test ./internal/application/workflow/...`
- P3.2：`go test ./internal/application/workflow/... -run 'Daily|Consultation'`
- P3.3：`go test ./internal/application/workflow/... -run 'Evidence|Market|Evolution|Gatekeeper'`

测试必须覆盖正常、信息不足、能力圈外、LLM 不可用、预期收益样本不足，以及 Evidence/Market/Evolution/Gatekeeper 的计划场景。
