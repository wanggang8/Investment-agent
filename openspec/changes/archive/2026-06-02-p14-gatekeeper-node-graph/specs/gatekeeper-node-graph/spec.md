## ADDED Requirements

### Requirement: Gatekeeper node-level graph
The system SHALL implement `GatekeeperAuditGraph` as a node-level Eino graph matching the workflow document.

#### Scenario: Gatekeeper graph exposes documented nodes
- **WHEN** the gatekeeper graph is constructed
- **THEN** `NodeNames()` SHALL return `ProposalLoadNode`, `FundamentalRuleCheckNode`, `ConflictCheckNode`, `BacktestNode`, `AuditDecisionNode`, and `AuditRecordNode`
- **THEN** `RegisteredNodeNames()` SHALL return the same node sequence

#### Scenario: Each gatekeeper node writes audit event
- **WHEN** a gatekeeper audit runs
- **THEN** the system SHALL write an audit event for each registered node
- **THEN** each event SHALL include workflow type, node name, node action, input reference, output reference, and status

### Requirement: Gatekeeper rejection gates
The system SHALL prevent unsafe or under-supported rule proposals from becoming applicable through gatekeeper audit.

#### Scenario: Insufficient samples cannot pass
- **WHEN** a proposal has fewer than three samples
- **THEN** the gatekeeper audit SHALL fail or reject before allowing application
- **THEN** no active rule version SHALL be written

#### Scenario: Fundamental rule violation cannot pass
- **WHEN** a proposal contains automatic trading, broker API, active recommendation, or return guarantee behavior
- **THEN** `FundamentalRuleCheckNode` SHALL mark the audit as not applicable
- **THEN** no active rule version SHALL be written

#### Scenario: Conflicting rule change cannot pass
- **WHEN** a proposal has a rule conflict or no effective rule change
- **THEN** `ConflictCheckNode` SHALL mark the audit as not applicable
- **THEN** no active rule version SHALL be written

### Requirement: User final confirmation boundary
The system SHALL keep gatekeeper approval separate from final rule application.

#### Scenario: Gatekeeper approval waits for final confirmation
- **WHEN** gatekeeper audit approves a valid proposal
- **THEN** the proposal SHALL move to `pending_final_confirm`
- **THEN** the system SHALL NOT create an active rule version until the user performs final confirmation
