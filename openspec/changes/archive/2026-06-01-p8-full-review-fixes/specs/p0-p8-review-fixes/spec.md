## ADDED Requirements

### Requirement: P0-P8 review fixes preserve decision confirmation consistency
The system SHALL persist decision records with record type, confirmation status, available actions, and detail DTO fields that are consistent with the decision verdict and confirmation service.

#### Scenario: Formal advice can be confirmed manually
- **WHEN** a consult decision creates a formal trade advice that requires user action
- **THEN** the decision detail SHALL expose a pending confirmation state with matching available offline actions, and manual confirmation SHALL be accepted.

#### Scenario: Non-actionable verdicts are not exposed as confirmable trades
- **WHEN** a decision verdict is insufficient data, rejected, or frozen watch
- **THEN** the persisted record type and confirmation status SHALL prevent misleading manual trade confirmation.

#### Scenario: Consult scenario values are contract-bound
- **WHEN** a consult request provides a scenario
- **THEN** only `hold_review`, `buy_review`, `sell_review`, and `rebalance_review` SHALL be accepted, and invalid values SHALL fail before writing a decision record.

### Requirement: P0-P8 review fixes preserve account facts after manual execution
The system SHALL update positions, position transactions, portfolio snapshots, position snapshots, decision status, and audit events atomically when a user records an offline manual execution, including transaction fees in cash and buy-side cost basis.

#### Scenario: Manual buy updates cash and assets
- **WHEN** a user records an offline buy from an account with an existing cash balance and transaction fees
- **THEN** the new snapshot SHALL reflect updated cash, position market value, total assets, cash ratio, and fee-inclusive position cost price.

#### Scenario: Manual sell updates or clears position and cash
- **WHEN** a user records an offline sell or reduce operation with transaction fees
- **THEN** the system SHALL reject insufficient quantity and otherwise update cash net of fees, remaining position, total assets, snapshots, and `position_transactions.fees` consistently.

#### Scenario: Manual execution write failure rolls back
- **WHEN** any write in the confirmation transaction fails
- **THEN** no partial confirmation, transaction, position, snapshot, or audit state SHALL remain.

### Requirement: P0-P8 review fixes enforce capability and rule proposal chains
The system SHALL use real capability configuration for consult decisions and SHALL preserve the user-confirmation, gatekeeper-audit, final-confirmation rule proposal chain.

#### Scenario: Out-of-scope consult is rejected before analysis
- **WHEN** the capability configuration excludes a consulted symbol
- **THEN** the decision SHALL be rejected or marked out-of-scope without treating analyst material as a final decision.

#### Scenario: Rule proposal reaches gatekeeper audit after user confirmation
- **WHEN** a pending user confirmation rule proposal is confirmed by the user
- **THEN** the system SHALL create or run the gatekeeper audit path, move the proposal according to the audit result, expose proposal reason, before/after rules, impact scope, risk notes, and audit summary to the frontend, render known audit results in Chinese, and prefer `before_rule.content` / `after_rule.content` over raw JSON when present.

#### Scenario: Final rule application requires approved audit
- **WHEN** a proposal has no approved gatekeeper audit allowing application
- **THEN** final confirmation SHALL fail and no active rule version SHALL be written.

### Requirement: P0-P8 review fixes return contract-complete DTOs
The system SHALL expose decision detail, market snapshot, evidence list, and audit DTOs with contract-required fields and safe defaults.

#### Scenario: Decision detail replays persisted reasoning
- **WHEN** a saved decision contains analyst reports, expected return scenarios, arbitration chain, evidence refs, and portfolio snapshot id
- **THEN** the detail endpoint SHALL return those fields instead of empty placeholders.

#### Scenario: Market snapshot includes freshness and metrics
- **WHEN** the latest market snapshot is queried
- **THEN** the API SHALL return trade date, data status, dedicated close price, dedicated turnover rate, and market metrics required by the frontend contract.

#### Scenario: Evidence list preserves source metadata
- **WHEN** evidence summaries are listed
- **THEN** source name, original URL, published time, captured time, content hash, time weight, relevance score, and high-grade independent source count SHALL come from persisted intelligence data.

#### Scenario: Evidence hashes are content-derived
- **WHEN** evidence content and metadata are unchanged across refreshes
- **THEN** `content_hash` and `chunk_hash` SHALL remain stable; when content changes, both hashes SHALL change.

#### Scenario: Market write failure has one failure audit
- **WHEN** market snapshot persistence fails
- **THEN** the API SHALL return `MARKET_SNAPSHOT_WRITE_FAILED`, write exactly one failed market refresh audit event, and SHALL NOT write a success event.

#### Scenario: Audit event DTO preserves event identity and detail fields
- **WHEN** audit events are listed or rendered
- **THEN** clients SHALL accept either `audit_event_id` or `event_id` as the event identity, and SHALL preserve workflow type, node name, node action, status transition, rule version, snapshot id, input ref, and output ref.

### Requirement: P0-P8 review fixes provide safe frontend states
The frontend SHALL render known backend statuses as Chinese text, render unknown statuses as “未知状态”, and avoid misleading normal-state warnings.

#### Scenario: Dashboard normal state has no missing-data warning
- **WHEN** dashboard state is normal
- **THEN** the dashboard SHALL NOT show missing-item or pause-reason copy.

#### Scenario: Frontend statuses are localized safely
- **WHEN** dashboard, evidence, portfolio, market, audit, rule proposal, or decision enums are displayed
- **THEN** known values SHALL appear in Chinese and unknown values SHALL appear as a safe unknown placeholder appropriate to the field.

#### Scenario: Expected return disclaimer remains visible
- **WHEN** an expected-return scenario contains both reason and disclaimer
- **THEN** the disclaimer SHALL remain visible.

#### Scenario: Manual execution form submits validated fees
- **WHEN** the user records an offline manual execution in the frontend
- **THEN** the form SHALL allow an optional non-negative finite fee and include it in the confirmation payload when present.

### Requirement: P0-P8 review fixes synchronize governance and quality gates
The repository SHALL reflect P0–P8 completion consistently across OpenSpec progress, development plan, testing plan, and lightweight OpenSpec summary specs.

#### Scenario: P6 and P7 plan status matches archives
- **WHEN** P6 and P7 archive tasks are complete
- **THEN** the development plan SHALL mark the corresponding completed items consistently.

#### Scenario: Full quality gate includes backend and frontend checks
- **WHEN** the project is verified before completion
- **THEN** the documented gate SHALL include Go tests, frontend build, and frontend tests.

#### Scenario: Full re-review repeats all review scopes
- **WHEN** fixes complete
- **THEN** subagents SHALL review backend, frontend, governance, data workflow, and test quality across the full current repository, not only changed files.
