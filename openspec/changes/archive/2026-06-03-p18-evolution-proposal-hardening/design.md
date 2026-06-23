# Design

## Goals
- Make review-derived rule suggestions auditable and durable.
- Route proposals through existing statuses and gatekeeper requirements.
- Preserve a clear source chain from review facts to proposal record to audit event.
- Keep all mutation gates unchanged.

## Approach

### Proposal creation path
Use the existing rule proposal repository and status model where possible. A review-derived proposal should start as a draft or user-confirmation-pending proposal, with source metadata linking it to review period, error cases, decisions, and audit events where available.

### Evolution proposal graph
If `EvolutionProposalGraph` already exists, add tests and minimal implementation to ensure it produces persisted proposal records. If the graph is currently a stub, implement the narrow path required for review suggestion input to proposal output.

### Audit and traceability
Proposal creation must write audit events with input references for review summary or supporting facts. The proposal must remain inspectable by existing API/frontend surfaces.

### Safety boundaries
The path must never write a new active rule version. Rule application still requires gatekeeper audit and final user confirmation. Empty samples or insufficient evidence should produce review output or a rejected/blocked proposal state, not an applied rule.

## Risks and mitigations
- **Risk:** Review effectiveness data could be treated as permission to change rules.  
  **Mitigation:** Tests assert no active rule version is created by review proposal generation.
- **Risk:** Source traceability is weak.  
  **Mitigation:** Persist input refs and audit output refs using existing fields.
- **Risk:** Implementation expands into a full rule-learning system.  
  **Mitigation:** Keep P18 limited to durable proposal creation and queueing.
