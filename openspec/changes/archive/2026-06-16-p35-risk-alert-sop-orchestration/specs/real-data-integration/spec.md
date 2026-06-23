## ADDED Requirements

### Requirement: Data health SHALL feed risk alert orchestration

The system SHALL use real data freshness, source health, and degraded public data diagnostics as inputs to risk alert orchestration without treating missing data as complete evidence.

#### Scenario: Source health is stale or failed
- **WHEN** P34 source health indicates stale, missing, no_data, source_unavailable, parse_error, disabled, or stubbed status for a category needed by risk analysis
- **THEN** risk alert orchestration SHALL preserve that condition as degraded or insufficient data context
- **AND** it SHALL NOT silently clear risk alerts that depend on the missing or degraded category.

#### Scenario: Expanded data supports a risk alert
- **WHEN** P34 expanded data is fresh enough to support valuation, liquidity, sentiment, or evidence-insufficiency risk checks
- **THEN** risk alert orchestration SHALL record the source category, freshness, source level, data date, and affected symbols in the risk trigger context where available
- **AND** it SHALL keep lower-grade or stubbed data visibly marked as supporting context rather than formal high-confidence evidence.

### Requirement: Risk orchestration SHALL not expand external data boundaries

Risk alert orchestration SHALL only consume already configured local facts, market snapshots, evidence summaries, and source health records.

#### Scenario: Risk orchestration needs more data
- **WHEN** required risk inputs are missing or stale
- **THEN** the system SHALL produce degraded or insufficient-data risk diagnostics
- **AND** it SHALL NOT introduce login, paid, authorized, Level2, high-frequency, broker, or external push dependencies.
