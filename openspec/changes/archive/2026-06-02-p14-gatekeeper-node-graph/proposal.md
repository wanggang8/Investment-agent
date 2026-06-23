## Why

`docs/workflow.md` 已定义 `GatekeeperAuditGraph` 的节点顺序，但当前实现仍以单个 `GatekeeperAuditGraph` 审计事件表达整体行为。P14 需要把守门人审计对齐为节点级 Eino Graph，让样本不足、规则冲突、根本规则违反和审计决策都能被逐节点追踪。

## What Changes

- Implement `GatekeeperAuditGraph` as a node-level Eino graph.
- Register nodes matching `docs/workflow.md`: `ProposalLoadNode`, `FundamentalRuleCheckNode`, `ConflictCheckNode`, `BacktestNode`, `AuditDecisionNode`, `AuditRecordNode`.
- Emit audit events for every gatekeeper node.
- Ensure insufficient samples, fundamental rule violations, and conflicts cannot produce an applicable approval.
- Preserve the boundary that rule application still requires user final confirmation.

## Capabilities

### New Capabilities
- `gatekeeper-node-graph`: Covers node-level gatekeeper audit orchestration, audit events, and rejection checks.

### Modified Capabilities
- `real-data-integration`: No changes.
- `product-completeness`: Continues the node-level workflow orchestration baseline for the gatekeeper graph.

## Impact

- `internal/application/workflow/gatekeeper_audit_graph.go` and tests.
- Rule proposal status transition behavior and audit event assertions.
- No frontend scope in this change.
- No automatic rule application, trading API, automatic trading, active recommendation, or return guarantee.
