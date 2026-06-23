# Tasks: p3-eino-workflows

> 对齐 `docs/development-plan.md` P3：Eino 工作流。实现代码必须写中文注释，说明工作流、节点、审计和降级边界。

## 1. WorkflowContext 与节点框架（P3.1）

参考文档：`docs/workflow.md`。

创建文件：

```text
internal/application/workflow/context.go
internal/application/workflow/node.go
internal/application/workflow/audit_writer.go
internal/application/workflow/errors.go
```

任务：

- [x] 1.1 创建 `internal/application/workflow/context.go`
- [x] 1.2 创建 `internal/application/workflow/node.go`
- [x] 1.3 创建 `internal/application/workflow/audit_writer.go`
- [x] 1.4 创建 `internal/application/workflow/errors.go`
- [x] 1.5 定义 WorkflowContext
- [x] 1.6 定义节点输入输出约定
- [x] 1.7 每个节点必须返回状态、错误码、审计事件片段
- [x] 1.8 工作流节点审计必须填写 `action`、`node_name`、`node_action`、`status`、`input_ref_type/input_ref`；产生输出时填写 `output_ref_type/output_ref`
- [x] 1.9 `status=failed` 时必须填写 `error_code`；降级有明确原因时填写 `error_code`
- [x] 1.10 失败节点必须写入 `audit_events`
- [x] 1.11 为 P3.1 新增代码补充中文注释
- [x] 1.12 验收：`go test ./internal/application/workflow/...`

## 2. DailyDisciplineGraph 与 ConsultationGraph（P3.2）

参考文档：`docs/workflow.md` 第 2、3、5、6 节。

创建文件：

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

任务：

- [x] 2.1 创建 `internal/application/workflow/daily_discipline_graph.go`
- [x] 2.2 创建 `internal/application/workflow/consultation_graph.go`
- [x] 2.3 创建 `internal/application/workflow/nodes/state_snapshot_node.go`
- [x] 2.4 创建 `internal/application/workflow/nodes/capability_check_node.go`
- [x] 2.5 创建 `internal/application/workflow/nodes/evidence_retrieval_node.go`
- [x] 2.6 创建 `internal/application/workflow/nodes/value_analyst_node.go`
- [x] 2.7 创建 `internal/application/workflow/nodes/trend_risk_officer_node.go`
- [x] 2.8 创建 `internal/application/workflow/nodes/expected_return_node.go`
- [x] 2.9 创建 `internal/application/workflow/nodes/rule_arbitration_node.go`
- [x] 2.10 创建 `internal/application/workflow/nodes/decision_record_node.go`
- [x] 2.11 每日纪律读取账户、持仓、市场、证据和规则版本
- [x] 2.12 主动咨询必须包含能力圈检查
- [x] 2.13 DeepSeek 节点只写分析报告
- [x] 2.14 ExpectedReturnNode 按 `sample_count` 映射 `precision_status`：`>=20 available`、`5~19 insufficient`、`<5 unavailable`
- [x] 2.15 `expected_return_scenarios` DTO 符合三种状态约束：available 可返回概率，insufficient 不返回精确概率且写样本不足说明，unavailable 返回空 `scenarios` 且写定性原因
- [x] 2.16 RuleArbitrationNode 生成最终裁决
- [x] 2.17 DecisionRecordNode 写 `decision_records`、`evidence_refs`、`audit_events`，并把预期收益情景保存到 `decision_records.expected_return_scenarios_json`
- [x] 2.18 为 P3.2 新增代码补充中文注释
- [x] 2.19 验收：`go test ./internal/application/workflow/... -run 'Daily|Consultation'`
- [x] 2.20 验收期望：正常、信息不足、能力圈外、LLM 不可用、预期收益样本不足五类场景均通过

## 3. Evidence、Evolution、Gatekeeper 工作流（P3.3）

参考文档：`docs/workflow.md`、`docs/data-model.md` 第 5、6 节。

创建文件：

```text
internal/application/workflow/evidence_verification_graph.go
internal/application/workflow/market_refresh_graph.go
internal/application/workflow/evolution_proposal_graph.go
internal/application/workflow/gatekeeper_audit_graph.go
```

任务：

- [x] 3.1 创建 `internal/application/workflow/evidence_verification_graph.go`
- [x] 3.2 创建 `internal/application/workflow/market_refresh_graph.go`
- [x] 3.3 创建 `internal/application/workflow/evolution_proposal_graph.go`
- [x] 3.4 创建 `internal/application/workflow/gatekeeper_audit_graph.go`
- [x] 3.5 情报写入 `intelligence_items`、`intelligence_summary`、`rag_chunks`
- [x] 3.6 VecLite 索引从 `rag_chunks` 构建，可由 SQLite 重建
- [x] 3.7 多源验证写 `source_verifications`
- [x] 3.8 市场刷新实现为独立 MarketRefreshGraph：读取数据源、标准化市场状态、写入 `market_snapshots` 与 `audit_events`
- [x] 3.9 错误案例生成规则提案，但不改正式规则
- [x] 3.10 守门人审计只生成 `gatekeeper_audits`
- [x] 3.11 审计通过后状态为 `pending_final_confirm`，不写正式规则
- [x] 3.12 用户最终确认后才写 `rule_versions`
- [x] 3.13 为 P3.3 新增代码补充中文注释
- [x] 3.14 验收：`go test ./internal/application/workflow/... -run 'Evidence|Market|Evolution|Gatekeeper'`

## 4. 归档前

- [x] 4.1 确认 `specs/workflow/spec.md` delta 已合并或已被 `docs/workflow.md` 覆盖
- [x] 4.2 勾选 `docs/development-plan.md` P3 相关任务
- [x] 4.3 更新 `docs/GOVERNANCE.md` 活跃变更状态
- [x] 4.4 更新 `openspec/PROGRESS.md`：P3 标记为 `in_progress`

## Plan alignment

- P3.1 对应任务：1.1–1.12，共 12 项。
- P3.2 对应任务：2.1–2.20，共 20 项。
- P3.3 对应任务：3.1–3.14，共 14 项。
- 归档前治理任务：4.1–4.4，共 4 项。
- 与 `docs/development-plan.md` P3 小节、创建文件、任务列表和验收命令一一对应；仅额外加入中文注释要求，来自本次用户指令。
