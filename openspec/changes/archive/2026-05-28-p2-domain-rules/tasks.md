# Tasks: p2-domain-rules

> 对齐 `docs/development-plan.md` P2：领域规则。

## 1. 核心模型与枚举（P2.1）

- [x] 1.1 创建 `internal/domain/model/enums.go`
- [x] 1.2 创建 `internal/domain/model/portfolio.go`
- [x] 1.3 创建 `internal/domain/model/market.go`
- [x] 1.4 创建 `internal/domain/model/evidence.go`
- [x] 1.5 创建 `internal/domain/model/decision.go`
- [x] 1.6 创建 `internal/domain/model/rule.go`
- [x] 1.7 创建 `internal/domain/model/audit.go`
- [x] 1.8 定义 `dashboard_state`、`workflow_status`、`position_state`、`verification_status`
- [x] 1.9 定义 `confirmation_status`、`confirmation_type`、`final_verdict.status`
- [x] 1.10 定义 `rule_proposal.status`、`audit_result`、`audit.action`、`audit.status`
- [x] 1.11 定义 `WorkflowContext` 对应领域结构
- [x] 1.12 编写枚举合法性测试
- [x] 1.13 验收：`go test ./internal/domain/model/...`

## 2. 规则裁决引擎（P2.2）

- [x] 2.1 创建 `internal/domain/rule/rules_engine.go`
- [x] 2.2 创建 `internal/domain/rule/source_policy.go`
- [x] 2.3 创建 `internal/domain/rule/capability_policy.go`
- [x] 2.4 创建 `internal/domain/rule/risk_policy.go`
- [x] 2.5 创建 `internal/domain/rule/expectation_engine.go`
- [x] 2.6 创建 `internal/domain/rule/gatekeeper_logic.go`
- [x] 2.7 实现能力圈外返回 `rejected`，并拒绝交易类分析
- [x] 2.8 实现证据不足返回 `insufficient_data`
- [x] 2.9 实现普通正式证据允许 S/A/B 级来源，C 级只能作为 `background`
- [x] 2.10 实现重大利好、重大利空、买入逻辑破坏至少 2 个 A 或 S 独立信源，不满足返回 `frozen_watch`
- [x] 2.11 实现买入逻辑破坏返回 `sell_only`
- [x] 2.12 实现情绪极端时暂停主动交易建议
- [x] 2.13 实现 PE/PB 分位区间规则
- [x] 2.14 实现移动止盈规则
- [x] 2.15 实现 R-5 现金冗余规则
- [x] 2.16 实现核心-卫星仓位规则
- [x] 2.17 实现预期收益评估，且不覆盖最终裁决
- [x] 2.18 实现规则提案完整状态机
- [x] 2.19 按 development-plan P2 场景表补充规则测试
- [x] 2.20 验收：`go test ./internal/domain/rule/...`

## 3. 归档前

- [x] 3.1 确认 `specs/domain-rules/spec.md` delta 已合并或已被 `docs/workflow.md` / `docs/data-model.md` 覆盖
- [x] 3.2 勾选 `docs/development-plan.md` P2 相关任务
- [x] 3.3 更新 `openspec/PROGRESS.md`：P2 标记为 `in_progress`

## Plan alignment

- P2.1 对应任务：1.1–1.13，共 13 项。
- P2.2 对应任务：2.1–2.20，共 20 项。
- 归档前治理任务：3.1–3.3，共 3 项。
