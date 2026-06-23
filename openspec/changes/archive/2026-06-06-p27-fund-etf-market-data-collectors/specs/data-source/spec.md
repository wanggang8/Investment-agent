## ADDED Requirements

### Requirement: P27 SHALL implement only verified read-only market data collectors

P27 SHALL implement production collectors only for public market data sources that passed P25 verification and remain inside the local-only safety boundary.

#### Scenario: First-batch market data collectors are selected

- **WHEN** P27 selects first-batch market data collectors
- **THEN** the default scope SHALL include CSIndex index data and Eastmoney fund data
- **AND** CSIndex MAY provide A-level index basics, constituents, weights, valuation files, and related file metadata when licensing and access boundaries remain public and low-frequency
- **AND** Eastmoney fund data SHALL be treated as B-level third-party fund NAV, accumulated NAV, fund profile, holdings, and asset-allocation data
- **AND** SSE fund trading summaries, Sina Finance ETF or market signals, broker data, Level2 data, paid data, login-only data, user account data, and trading endpoints SHALL NOT be required P27 runtime dependencies.

#### Scenario: A source violates the safety boundary

- **WHEN** a source requires login, paid access, CAPTCHA bypass, broker access, Level2 or authorized market data, user identity, trading access, or access-control circumvention
- **THEN** P27 SHALL exclude that source from implementation
- **AND** the system SHALL continue with fixture/stub fallback or other verified sources.

#### Scenario: A market collector output could trigger external side effects

- **WHEN** P27 collectors refresh fund NAV, ETF data, index data, valuation files, or market metadata
- **THEN** they SHALL NOT buy, sell, cancel, amend, or otherwise operate any brokerage account
- **AND** they SHALL NOT send email, SMS, system Push, Webhook, WebSocket, or any other external notification by default
- **AND** any resulting rule or decision state SHALL remain subject to existing user confirmation and gatekeeper audit requirements.

### Requirement: P27 collectors SHALL produce standardized market data payloads

P27 collectors SHALL normalize fund NAV, ETF, index, and valuation data before persistence.

#### Scenario: A fund NAV or index item is collected

- **WHEN** a collector successfully fetches a fund NAV, accumulated NAV, fund profile, index basic item, constituent, weight, valuation file, or related market data item
- **THEN** it SHALL produce a payload with `source_name`, `source_level`, `source_type`, `symbol`, `trade_date`, `nav`, `accumulated_nav`, `close_price`, `metadata`, `captured_at`, `content_hash`, and raw source metadata where applicable
- **AND** it SHALL persist or expose the item for persistence into `market_snapshots.metadata_json`, necessary `intelligence_items` or `rag_chunks`, and `audit_events`.

#### Scenario: A collector runs repeatedly over the same market data

- **WHEN** the same source record, trade date, file, or identical content is collected again
- **THEN** the system SHALL deduplicate by source identity, `symbol`, `trade_date`, `source_type`, or `content_hash`
- **AND** repeated runs SHALL NOT create duplicate market facts or duplicate RAG chunks.

### Requirement: P27 freshness and failures SHALL be explicit and auditable

P27 collectors SHALL record collection attempts, freshness state, and failures without blocking local application startup.

#### Scenario: A market source request fails

- **WHEN** a source is unavailable, pagination fails, file download fails, required fields are missing, date parsing fails, or parsing fails
- **THEN** the collector SHALL write an `audit_events` entry with a source-specific error code and counts where available
- **AND** the workflow SHALL degrade that source without blocking other sources or local startup.

#### Scenario: Market data is stale, missing, or lower-grade only

- **WHEN** collected market data is stale, missing, or only available from B-level third-party sources
- **THEN** the system SHALL keep the relevant market or decision state as `stale`, `missing`, `degraded`, `insufficient_data`, or `frozen_watch` where applicable
- **AND** it SHALL NOT use Eastmoney or any other B-level third-party source as A-level formal evidence or the only basis for clearing insufficient data.

#### Scenario: Percentile or valuation fields are unavailable

- **WHEN** the collector lacks verified source fields for `pe_percentile`, `pb_percentile`, `volume_percentile`, `volatility_percentile`, or equivalent valuation percentile data
- **THEN** the system SHALL leave those fields absent, null, or marked missing
- **AND** it SHALL NOT infer, fabricate, or backfill those percentile fields from unrelated B-level data.

### Requirement: P27 backfill and refresh SHALL be bounded and low-frequency

P27 SHALL support bounded backfill and low-frequency refresh for verified market data sources.

#### Scenario: A first-batch market collector performs initial backfill

- **WHEN** P27 runs initial backfill for a verified market data source
- **THEN** it SHALL default to the latest 90 days or the nearest publicly available historical NAV/file range
- **AND** it SHALL support idempotency, audit logging, safe interruption, and fixture/stub fallback.

#### Scenario: Incremental market refresh runs after backfill

- **WHEN** incremental refresh runs on a trading day
- **THEN** fund NAV sources SHOULD refresh after 21:30 local exchange time and perform one next-day follow-up refresh when data may settle late
- **AND** index basics, constituent, weight, and valuation file sources SHOULD refresh daily or by published file update time
- **AND** refresh SHALL NOT perform high-frequency market data scraping or real-time trading data polling.
