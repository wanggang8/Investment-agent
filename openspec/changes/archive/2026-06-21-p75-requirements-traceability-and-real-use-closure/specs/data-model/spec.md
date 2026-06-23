## ADDED Requirements

### Requirement: Configured Symbol Profile Registry For Knowledge Readiness

The knowledge readiness data model SHALL represent accepted local fund/ETF symbol profiles as structured registry entries instead of a single hardcoded accepted symbol. The current accepted local registry SHALL include `510300 -> 000300` and `159915 -> 399006`; unsupported symbols SHALL remain `known=false` with blocked or information-insufficient readiness and SHALL NOT receive fabricated market, valuation, liquidity, formal-evidence, or RAG readiness.

#### Scenario: Accepted non-510300 symbol profile is traceable

- **GIVEN** the user-entered symbol is `159915`
- **WHEN** readiness, collector routing, and analyst LLM context are built
- **THEN** the profile SHALL resolve to tracked index `399006`
- **AND** fund-side collector evidence SHALL stay bound to `159915`
- **AND** index-side collector evidence SHALL stay bound to `399006`
- **AND** readiness dependencies SHALL expose source, `data_date`, freshness, `request_id`, and affected symbol correlation where source-health evidence exists
- **AND** analyst LLM context SHALL include `symbol_profile.159915` and SHALL NOT silently include `symbol_profile.510300` for the `159915` flow.

#### Scenario: Built-in symbol profiles do not replace evidence

- **GIVEN** a configured symbol profile is available
- **WHEN** consultation, alerts, expected-return, data-quality, or release readiness are evaluated
- **THEN** the profile MAY provide symbol/tracked-index routing and background context
- **AND** it SHALL NOT satisfy formal evidence, source verification, current data, valuation, liquidity, public-source availability, expected-return provenance, or trade-like confirmation requirements by itself.
