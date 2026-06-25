## ADDED Requirements

### Requirement: Product Usability Edge Scenario Acceptance Evidence

The release governance evidence SHALL include a repeatable product-usability edge scenario acceptance pass before declaring the local product resilient for accumulated real-use workflows beyond the primary happy path.

#### Scenario: Long-cycle and abnormal-use acceptance

- **WHEN** an acceptance runner simulates accumulated local use beyond the primary continuous-use path
- **THEN** it SHALL cover long-history data accumulation, abnormal input recovery, data-quality degradation handling, context-sensitive decision interpretation, multi-account household ledger facts, cross-page readback and restart persistence
- **AND** it SHALL write API/SQLite/browser evidence and a product-usability interpretation report
- **AND** it SHALL exclude release/install/upgrade claims unless a separate release deployment change explicitly covers them.

#### Scenario: Edge-scenario safety boundaries

- **WHEN** the edge scenario acceptance completes
- **THEN** it SHALL report negative evidence for broker/order/push tables, automatic confirmation rows, automatic rule-application audit events, trading affordances, return guarantee claims and secret/raw prompt leakage
- **AND** stale historical audit checks SHALL be recorded as stale rather than treated as fresh passes.
