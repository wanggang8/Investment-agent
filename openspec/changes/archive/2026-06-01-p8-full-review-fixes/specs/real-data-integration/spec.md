## ADDED Requirements

### Requirement: P8 review data metadata preservation
The system SHALL preserve source metadata, verification status, market freshness, and rule proposal audit facts across persistence, DTOs, and frontend display.

#### Scenario: Evidence metadata is preserved from intelligence item to API
- **WHEN** intelligence items and summaries are listed as evidence
- **THEN** source name, original URL, published time, captured time, content hash, time weight, relevance score, source level, evidence role, and verification status SHALL be returned without placeholder substitution.

#### Scenario: Market snapshot reports freshness
- **WHEN** market data is refreshed or queried
- **THEN** the market DTO SHALL include enough date and status information for the frontend to distinguish fresh, stale, and missing data.

#### Scenario: Rule proposal audit facts gate final application
- **WHEN** final rule application is requested
- **THEN** the latest approved gatekeeper audit allowing application SHALL be present before an active rule is written.
