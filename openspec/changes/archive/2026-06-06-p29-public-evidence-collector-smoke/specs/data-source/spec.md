## ADDED Requirements

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
