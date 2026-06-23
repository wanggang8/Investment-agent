## ADDED Requirements

### Requirement: Daily discipline reports SHALL surface expanded data coverage state

Daily discipline report surfaces SHALL show whether P34 expanded public data was fresh, stale, missing, unavailable, or degraded when the report was generated.

#### Scenario: Expanded data is available for a report
- **WHEN** a daily discipline report is generated with available P34 expanded data
- **THEN** the report context SHALL include source category, source level, data date, freshness state, and affected symbols or indexes
- **AND** the frontend SHALL be able to display that the report used expanded public data as analysis context.

#### Scenario: Expanded data is missing or stale for a report
- **WHEN** a daily discipline report lacks required P34 data categories or receives stale data
- **THEN** the report SHALL include missing or stale categories in its diagnostics
- **AND** it SHALL not mark those categories as satisfied by unrelated or lower-grade data.

#### Scenario: Expanded source fails during report preparation
- **WHEN** a P34 source fails with no data, source unavailable, parse error, timeout, or write failure during report preparation or refresh
- **THEN** the report SHALL preserve a degraded or insufficient-data explanation
- **AND** it SHALL not imply that a broker trade, order, external notification, or guaranteed return was produced.
