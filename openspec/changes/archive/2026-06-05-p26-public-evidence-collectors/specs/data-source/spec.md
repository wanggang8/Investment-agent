# 公告与证据源 Collector Delta

## ADDED Requirements

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
