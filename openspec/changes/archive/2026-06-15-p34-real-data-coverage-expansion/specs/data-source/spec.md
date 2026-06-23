## ADDED Requirements

### Requirement: P34 collectors SHALL classify expanded public data sources by category

P34 collectors SHALL classify every expanded public data source by category, source level, refresh cadence, and safety boundary before runtime use.

#### Scenario: Expanded source is selected
- **WHEN** P34 selects a candidate source for index samples, index weights, valuation files, constituent financials, capital flow, margin financing, or sentiment proxy data
- **THEN** the source SHALL be documented as public, readonly, low-frequency, and not requiring login, paid access, authorization, CAPTCHA bypass, broker access, Level2, or high-frequency polling
- **AND** the source SHALL have a source category and source level used in payload metadata and audit records.

#### Scenario: Expanded source is only partially usable
- **WHEN** a source exposes only some required fields or only B/C-level third-party data
- **THEN** the collector SHALL persist only verified fields with explicit missing markers for unavailable fields
- **AND** it SHALL not promote the source to A-level evidence or use it alone to clear insufficient data.

### Requirement: P34 collectors SHALL normalize expanded market and source-health payloads

P34 collectors SHALL normalize expanded public data before persistence so downstream workflows can consume it without source-specific parsing.

#### Scenario: Expanded data item is collected
- **WHEN** a collector successfully fetches an expanded data item
- **THEN** it SHALL produce a payload with `source_name`, `source_level`, `source_type`, `data_category`, `symbol`, `trade_date` or `data_date`, `captured_at`, `content_hash`, normalized metrics, and raw source metadata where available
- **AND** it SHALL persist or expose the item for persistence into `market_snapshots.market_metrics_json`, source health metadata, audit events, or another explicitly documented P34 storage path.

#### Scenario: Expanded collector is repeated
- **WHEN** the same source record, file, symbol, date, or content is collected again
- **THEN** the collector SHALL deduplicate by source identity, symbol, date, file identity, or content hash
- **AND** repeated refreshes SHALL NOT create conflicting market facts or duplicate health records.

### Requirement: P34 refresh SHALL remain bounded and low-frequency

P34 data refresh SHALL be bounded, local-triggered, and low-frequency for all expanded public sources.

#### Scenario: P34 refresh runs
- **WHEN** a user or local task triggers expanded data refresh
- **THEN** the refresh SHALL use explicit source configuration, bounded symbol or index scope, and documented date windows
- **AND** `cmd/agent --task p34-expanded-refresh` SHALL accept source, symbol or index, start date, and end date inputs for explicit local refresh
- **AND** it SHALL support fixture or stub fallback when stub mode is explicitly enabled
- **AND** it SHALL NOT perform high-frequency market scraping or realtime polling.

#### Scenario: P34 source fails during refresh
- **WHEN** one expanded source fails with `no_data`, `source_unavailable`, `parse_error`, stale data, or write failure
- **THEN** other enabled sources MAY continue
- **AND** the failed source SHALL write source-specific audit and health metadata without hiding successful source results.
