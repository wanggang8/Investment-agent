# Delta for P3 Eino Workflows（合并目标：`docs/workflow.md`）

## ADDED Requirements

### Requirement: Workflow nodes must emit auditable execution fragments

系统 SHALL 在应用层工作流节点中统一返回节点状态、错误码和审计事件片段。

#### Scenario: Successful node execution

- **GIVEN** 一个节点成功产生输出
- **WHEN** 节点返回结果
- **THEN** 审计片段必须包含 `action`、`node_name`、`node_action`、`status`、`input_ref_type/input_ref`
- **AND** 若存在输出引用，必须包含 `output_ref_type/output_ref`

#### Scenario: Failed node execution

- **GIVEN** 一个节点执行失败
- **WHEN** 节点返回结果
- **THEN** `status` 必须为 `failed`
- **AND** 必须填写 `error_code`
- **AND** 失败节点必须写入 `audit_events`

### Requirement: Daily and consultation workflows must preserve rule authority

系统 SHALL 通过 DailyDisciplineGraph 和 ConsultationGraph 组织状态快照、证据核查、分析材料、预期收益、规则裁决和记录保存，并保证最终裁决只来自领域规则。

#### Scenario: Consultation out of capability scope

- **GIVEN** 主动咨询标的不在能力圈内
- **WHEN** ConsultationGraph 执行
- **THEN** 工作流必须包含能力圈检查
- **AND** RuleArbitrationNode 必须生成 `rejected` 类型裁决

#### Scenario: LLM unavailable

- **GIVEN** DeepSeek 分析节点不可用
- **WHEN** DailyDisciplineGraph 或 ConsultationGraph 执行
- **THEN** 工作流可降级记录 `ANALYST_UNAVAILABLE`
- **AND** DeepSeek 节点不得写最终裁决

#### Scenario: Expected return precision mapping

- **GIVEN** ExpectedReturnNode 读取到样本数
- **WHEN** `sample_count>=20`
- **THEN** `precision_status` 为 `available` 且可返回概率
- **WHEN** `5<=sample_count<20`
- **THEN** `precision_status` 为 `insufficient` 且不得返回精确概率
- **WHEN** `sample_count<5`
- **THEN** `precision_status` 为 `unavailable` 且 `scenarios` 为空

### Requirement: Evidence, market, evolution and gatekeeper workflows must write only allowed facts

系统 SHALL 将 P3 工作流输出写入对应 SQLite 事实表，并遵守规则演进边界。

#### Scenario: Evidence verification writes source verification

- **GIVEN** 情报刷新和证据核查执行完成
- **WHEN** EvidenceVerificationGraph 保存结果
- **THEN** 必须写入 `intelligence_items`、`intelligence_summary`、`rag_chunks` 与 `source_verifications`

#### Scenario: Gatekeeper approval waits for final confirmation

- **GIVEN** 守门人审计通过
- **WHEN** GatekeeperAuditGraph 保存审计结果
- **THEN** 只生成 `gatekeeper_audits`
- **AND** 提案状态为 `pending_final_confirm`
- **AND** 不写正式 `rule_versions`

## MODIFIED Requirements

（无）

## REMOVED Requirements

（无）
