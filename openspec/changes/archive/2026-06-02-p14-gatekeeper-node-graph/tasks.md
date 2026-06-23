## 1. Gatekeeper graph contract tests

- [x] Add tests for `GatekeeperAuditGraph.NodeNames()` and `RegisteredNodeNames()` matching `docs/workflow.md`.
- [x] Add tests that a successful gatekeeper audit writes one audit event per documented node.
- [x] Add tests that node audit events include workflow type, node name, node action, input/output refs, and status.

## 2. Gatekeeper rejection tests

- [x] Add tests that proposals with fewer than three samples cannot pass gatekeeper audit.
- [x] Add tests that fundamental rule violations cannot allow application.
- [x] Add tests that conflicting or ineffective rule changes cannot allow application.
- [x] Add tests that gatekeeper approval moves valid proposals to `pending_final_confirm` without creating active rule versions.

## 3. Node-level implementation

- [x] Implement `GatekeeperAuditGraph` as a node-level Eino graph.
- [x] Register `ProposalLoadNode`, `FundamentalRuleCheckNode`, `ConflictCheckNode`, `BacktestNode`, `AuditDecisionNode`, and `AuditRecordNode`.
- [x] Preserve existing transactional writes for `gatekeeper_audits`, proposal status updates, and audit events.
- [x] Preserve final user confirmation boundary.

## 4. Documentation and OpenSpec sync

- [x] Keep `docs/workflow.md` node list as the source for P14 node names.
- [x] Update only change-local specs before archive; merge delta during archive.

## 5. Validation

- [x] Run `openspec validate p14-gatekeeper-node-graph --strict`.
- [x] Run `go test ./...`.
- [x] Run `cd web && npm run test && npm run build`.
- [x] Run `openspec validate --all --strict`.
- [x] Run relevant `cmd/agent` local task validation where applicable.
