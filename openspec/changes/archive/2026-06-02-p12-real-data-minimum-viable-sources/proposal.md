## Why

P10 已提供可配置数据源入口，但真实只读 provider 仍以不可用状态为主。P12 需要接入最小可用真实只读行情与情报来源，同时保留 stub、降级和审计，避免把“真实数据”范围扩大到完整财务、完整情绪或交易能力。

## What Changes

- Add one minimum viable readonly market data provider.
- Add one minimum viable readonly intelligence provider.
- Preserve local stub mode for development and offline validation.
- Return stable degraded errors and audit events for timeout, unavailable source, stale data, and parse failure.
- Explicitly keep full financial data sources, full sentiment sources, realtime SLA, broker trading APIs, automatic trading, active recommendations, and return guarantees out of scope.

## Capabilities

### New Capabilities
- `real-data-minimum-viable-sources`: Covers minimum viable readonly providers, fallback behavior, auditability, and out-of-scope boundaries for P12.

### Modified Capabilities
- `real-data-integration`: Narrows existing real data integration requirements into a minimum viable provider requirement with explicit P12 exclusions.

## Impact

- Backend data source adapters and wiring.
- Configuration examples and source status reporting.
- Tests for provider success, timeout/unavailable/parse failure, stub fallback, and audit events.
- No frontend behavior is required beyond any existing source status display.
- No trading integration or account mutation is introduced.
