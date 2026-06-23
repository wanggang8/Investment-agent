# data-source Specification

## Purpose
TBD - created by archiving change p26-public-evidence-collectors. Update Purpose after archive.
## Requirements
### Requirement: P26 SHALL implement only verified public evidence collectors

P26 SHALL implement production collectors only for public evidence sources that passed P25 verification and remain inside the local-only safety boundary.

#### Scenario: First-batch public evidence collectors are selected

- **WHEN** P26 selects first-batch collectors
- **THEN** the default scope SHALL include CNInfo announcements, SZSE announcements, and CSRC regulatory information
- **AND** AMAC public industry statistics or self-discipline pages MAY be included only as background evidence when stable access and safety boundaries are confirmed
- **AND** SSE announcements, AMAC institution/product/person queries, Eastmoney fund market data, CSIndex market/index files, and Sina Finance background data SHALL NOT be required P26 runtime dependencies.

#### Scenario: A source violates the safety boundary

- **WHEN** a source requires login, paid access, CAPTCHA bypass, broker access, Level2 or authorized market data, or access-control circumvention
- **THEN** P26 SHALL exclude that source from implementation
- **AND** the system SHALL continue with fixture/stub fallback or other verified sources.

#### Scenario: A collector output could trigger external side effects

- **WHEN** P26 collectors refresh public evidence or produce evidence verification results
- **THEN** they SHALL NOT buy, sell, cancel, amend, or otherwise operate any brokerage account
- **AND** they SHALL NOT send email, SMS, system Push, Webhook, WebSocket, or any other external notification by default
- **AND** any resulting rule or decision state SHALL remain subject to existing user confirmation and gatekeeper audit requirements.

### Requirement: P26 collectors SHALL produce standardized evidence payloads

P26 collectors SHALL normalize public announcements and regulatory materials before persistence.

#### Scenario: A public evidence item is collected

- **WHEN** a collector successfully fetches an announcement or regulatory item
- **THEN** it SHALL produce a payload with `source_name`, `source_level`, `source_type`, `evidence_role`, `symbol`, `title`, `text`, `url`, `attachment_url`, `published_at`, `captured_at`, `content_hash`, and raw source metadata
- **AND** it SHALL persist or expose the item for persistence into `intelligence_items`, `rag_chunks`, `source_verifications`, and `audit_events`.

#### Scenario: A collector runs repeatedly over the same source records

- **WHEN** the same source record or identical content is collected again
- **THEN** the system SHALL deduplicate by source record identity or `content_hash`
- **AND** repeated runs SHALL NOT create duplicate facts or duplicate RAG chunks.

### Requirement: P26 collectors SHALL be auditable and degradable

P26 collectors SHALL record collection attempts and failures without blocking local application startup.

#### Scenario: A source request fails

- **WHEN** a source is unavailable, pagination fails, an attachment fails, required fields are missing, or parsing fails
- **THEN** the collector SHALL write an `audit_events` entry with a source-specific error code and counts where available
- **AND** the workflow SHALL degrade that source without blocking other sources or local startup.

#### Scenario: Evidence remains insufficient

- **WHEN** collected evidence does not meet the configured independent A/S source requirement
- **THEN** the system SHALL keep the relevant decision or evidence verification state as `insufficient_data`, `degraded`, or `frozen_watch`
- **AND** it SHALL NOT generate a high-confidence formal conclusion from a single source or B-level source alone.

### Requirement: P26 backfill and refresh SHALL be low-frequency and local-only

P26 SHALL support bounded backfill and low-frequency refresh for verified public evidence sources.

#### Scenario: A first-batch collector performs initial backfill

- **WHEN** P26 runs initial backfill for a verified evidence source
- **THEN** it SHALL default to the latest 90 days
- **AND** it SHALL support pagination, idempotency, audit logging, and safe interruption.

#### Scenario: Incremental refresh runs after backfill

- **WHEN** incremental refresh runs on a trading day
- **THEN** announcement sources SHOULD refresh at low frequency such as 30-60 minutes
- **AND** regulatory or industry background sources SHOULD refresh daily or less frequently
- **AND** refresh SHALL NOT perform high-frequency market data scraping.

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
- **AND** it SHALL persist or expose the item for persistence into `market_snapshots.market_metrics_json` and `audit_events`; market collector data SHALL NOT be represented as public evidence RAG chunks unless a later change explicitly adds that evidence mapping.

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

### Requirement: P27 refresh SHALL remain bounded and low-frequency

P27 SHALL support bounded, low-frequency refresh for verified market data sources through explicit local refresh tasks.

#### Scenario: A first-batch market collector runs through market refresh

- **WHEN** P27 runs `market-refresh` for a verified market data source and symbol
- **THEN** it SHALL collect the latest publicly available fund NAV, index item, or source metadata that the current collector can safely read
- **AND** it SHALL support idempotency, audit logging, safe interruption, and fixture/stub fallback only when stub mode is explicitly enabled.

#### Scenario: Incremental market refresh is scheduled by a future local scheduler

- **WHEN** incremental refresh is scheduled on a trading day
- **THEN** fund NAV sources SHOULD be scheduled no earlier than the normal public NAV publication window and MAY perform a later follow-up refresh when data may settle late
- **AND** index basics, constituent, weight, and valuation file sources SHOULD refresh daily or by published file update time after their real public endpoints are calibrated
- **AND** refresh SHALL NOT perform high-frequency market data scraping or real-time trading data polling.

### Requirement: P29 public evidence ingestion SHALL distinguish no-data, unavailable, and parse failures

P29 public evidence collectors SHALL provide reproducible real-source smoke validation and explicit source-level diagnostics for public evidence ingestion.

#### Scenario: Real public evidence smoke writes evidence tables

- **WHEN** `public-evidence-refresh` runs with an explicitly enabled real public evidence source, symbol, and date window that contains public evidence
- **THEN** the system SHALL fetch read-only public evidence without login, paid access, browser scraping, or high-frequency requests
- **AND** it SHALL write collected payloads into `intelligence_items`, `intelligence_summary`, and `rag_chunks`
- **AND** it SHALL write `source_verifications` when verification inputs are available
- **AND** it SHALL write `audit_events` with source-specific success, degraded, or failed status
- **AND** it SHALL NOT update positions, portfolio snapshots, confirmations, transactions, orders, broker state, or external notifications.

#### Scenario: CNInfo requires an orgId for a symbol

- **WHEN** CNInfo collection runs for a symbol whose public endpoint requires `stock=<symbol>,<orgId>`
- **THEN** the collector SHALL support `data_sources.public_evidence.cninfo_org_ids` mapping
- **AND** direct `symbol,orgId` input MAY be used for explicit smoke validation
- **AND** persisted evidence summaries and RAG chunks SHALL store the clean security code as `symbol`.

#### Scenario: All enabled public evidence sources have no records in the requested window

- **WHEN** every enabled public evidence source is reachable but returns no matching records for the requested symbol and window
- **THEN** `public-evidence-refresh` SHALL complete as a successful empty refresh
- **AND** it SHALL write a success audit entry with `count=0`
- **AND** it SHALL keep source-specific `no_data` degraded audit diagnostics.

#### Scenario: Public evidence endpoint is unavailable or incompatible

- **WHEN** a public evidence endpoint returns an HTTP failure, DNS/client failure, or incompatible response shape
- **THEN** the collector SHALL classify the condition as `source_unavailable` or `parse_error`
- **AND** other enabled collectors MAY continue
- **AND** if no enabled collector returns usable payloads and at least one source failed with `source_unavailable` or `parse_error`, `public-evidence-refresh` SHALL fail with a clear data source error and failed audit metadata.

### Requirement: P29 public evidence collectors SHALL prove real ingestion with smoke validation

The system SHALL provide a reproducible real-source smoke validation for public evidence collectors before claiming the collectors are stable for real ingestion.

#### Scenario: Real public evidence smoke writes evidence tables

- **WHEN** `public-evidence-refresh` is run with an explicitly enabled real public evidence source, a documented smoke symbol, and a documented smoke date window that contains public evidence
- **THEN** the system SHALL fetch read-only public evidence without login, paid access, browser scraping, or high-frequency requests
- **AND** it SHALL write at least one collected payload into `intelligence_items`, `intelligence_summary`, and `rag_chunks`
- **AND** it SHALL write `source_verifications` when verification inputs are available
- **AND** it SHALL write `audit_events` with source-specific success, degraded, or failed status
- **AND** it SHALL NOT update positions, portfolio snapshots, confirmations, transactions, orders, broker state, or external notifications.

#### Scenario: All enabled public evidence sources have no records in the requested window

- **WHEN** all enabled real public evidence sources are reachable but return no matching records for the requested symbol and window
- **THEN** `public-evidence-refresh` SHALL complete as a successful empty refresh
- **AND** it SHALL write a success audit entry with `count=0`
- **AND** it SHALL keep source-specific `no_data` degraded audit diagnostics
- **AND** it SHALL NOT misclassify the reachable no-record result as parser failure or trading-related state.

#### Scenario: Public evidence source endpoint is unavailable or incompatible

- **WHEN** a real public evidence endpoint returns an incompatible response, HTTP failure, or parse failure
- **THEN** the collector SHALL record a source-specific failure code
- **AND** other enabled collectors MAY continue so that one source failure does not hide successful data from other sources
- **AND** if no enabled collector returns usable payloads and at least one source failed with `source_unavailable` or `parse_error`, `public-evidence-refresh` SHALL fail with a clear data source error and write failed audit metadata.

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

