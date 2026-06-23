## Context

`GatekeeperAuditGraph` is responsible for auditing rule proposals before they can reach final user confirmation. The workflow document defines six nodes, but the current implementation records a single graph-level audit event.

## Goals / Non-Goals

**Goals:**
- Expose `NodeNames()` and `RegisteredNodeNames()` for the six documented nodes.
- Run gatekeeper audit through node-level Eino graph orchestration.
- Write one audit event per node.
- Keep insufficient samples, conflicts, and fundamental-rule violations from producing an applicable audit approval.
- Keep final rule application behind explicit user confirmation.

**Non-Goals:**
- No frontend page expansion.
- No EvolutionProposalGraph implementation.
- No P15 evidence enrichment.
- No automatic trading or direct broker integration.

## Decisions

1. Use the existing CloudWeGo Eino graph pattern.
   - Rationale: Daily, consultation, and evidence graphs already expose registered node names.
   - Alternative considered: only add `NodeNames()` without Eino compilation. Rejected because P14 specifically requires node-level graph alignment.

2. Keep rule proposal writes transactional.
   - Rationale: `gatekeeper_audits`, proposal status update, and audit events must remain consistent.
   - Alternative considered: write each node independently. Rejected because partial node persistence would complicate rollback semantics.

3. Treat rule application as out-of-scope.
   - Rationale: Gatekeeper approval only moves proposals toward `pending_final_confirm`; final application remains a separate user-confirmed action.

## Risks / Trade-offs

- More node events increase audit volume, but improve traceability.
- Existing tests expecting one audit event may need to assert node-level events instead.
- Backtest remains a sample-count gate in this stage; richer metrics can be handled later.
