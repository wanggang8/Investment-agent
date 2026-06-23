## ADDED Requirements

### Requirement: P48 SHALL provide repeatable data source quality regression

The system SHALL provide a local data source quality regression capability that verifies source health freshness, failure classification, and safe summaries without requiring external network access by default.

#### Scenario: Fixture regression runs offline
- **WHEN** P48 data source quality regression runs in `fixture` mode
- **THEN** it SHALL evaluate deterministic local cases for `fresh`, `no_data`, `source_unavailable`, `parse_error`, `stale`, and sensitive diagnostic redaction
- **AND** it SHALL return a stable summary containing case status, expected freshness, actual freshness, data category, affected symbols, and safety note
- **AND** it SHALL NOT access public endpoints, private files, broker systems, external notification channels, or trading APIs.

#### Scenario: Current source health is evaluated read-only
- **WHEN** P48 data source quality regression runs in `current` mode
- **THEN** it SHALL evaluate the latest existing P34 source health from local market snapshots
- **AND** it SHALL NOT trigger collectors, refresh market data, rebuild indexes, call LLMs, create notifications, update rules, update confirmations, or change account and position facts.

### Requirement: P48 SHALL expose sanitized regression summaries

P48 SHALL expose data source quality regression through local API and CLI summaries that do not reveal sensitive diagnostics.

#### Scenario: Regression API is requested
- **WHEN** a user calls `GET /api/v1/data-source-quality/regression`
- **THEN** the API SHALL return `mode`, `status`, `generated_at`, `summary`, `cases`, `missing_categories`, and `safety_note`
- **AND** each case SHALL include only sanitized `diagnostic_preview` text
- **AND** the API SHALL NOT write SQLite records or trigger external side effects.

#### Scenario: Regression CLI task is executed
- **WHEN** a user runs `go run ./cmd/agent --task data-source-quality-regression`
- **THEN** the task SHALL run the same regression service and print a compact local summary
- **AND** it MAY write a local `audit_events` record containing only mode, status, case counts, degraded/failed counts, and safety boundary
- **AND** it SHALL NOT store raw source payloads, complete API keys, private paths, raw SQL, full prompts, raw HTTP exchanges, private keys, or supplier raw responses in output or audit metadata.

### Requirement: P48 SHALL preserve data source safety boundaries

P48 SHALL remain a local quality regression feature and SHALL NOT expand the system's data acquisition, trading, or automation boundaries.

#### Scenario: A requested regression source is outside the allowed boundary
- **WHEN** a requested mode or source would require login, paid access, authorization-gated data, CAPTCHA bypass, broker access, Level2 data, high-frequency polling, browser scraping, or access-control circumvention
- **THEN** P48 SHALL reject or omit that mode/source
- **AND** it SHALL continue to provide fixture regression and current local source health evaluation where available.

#### Scenario: Regression detects degraded data quality
- **WHEN** regression cases detect `no_data`, `source_unavailable`, `parse_error`, `stale`, `missing`, `unknown`, or another degraded but recognized condition
- **THEN** the system SHALL classify the regression as degraded or failed according to documented rules
- **AND** it SHALL NOT automatically repair data, mark sources healthy, execute trades, apply rules, confirm operations, or promise returns.
