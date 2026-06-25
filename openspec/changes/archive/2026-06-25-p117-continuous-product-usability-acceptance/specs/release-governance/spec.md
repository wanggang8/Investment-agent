## ADDED Requirements

### Requirement: Continuous Product Usability Acceptance Evidence

The release governance evidence SHALL include a repeatable continuous-use acceptance pass before declaring the local product broadly usable for real user workflows.

#### Scenario: Seven-day local usability journey

- **WHEN** an acceptance runner simulates a seven-day local user journey
- **THEN** it SHALL cover cold start, account onboarding, daily routine, offline transaction updates, invalid input recovery, data-quality degradation handling, manual decision confirmation, marked-error review, cross-page readback and restart persistence
- **AND** it SHALL write API/SQLite/browser evidence and a usability interpretation report
- **AND** it SHALL distinguish local seeded evidence from external provider, LLM quality, broker execution, release packaging or physical second-machine claims.

#### Scenario: Continuous-use safety boundaries

- **WHEN** the continuous-use acceptance completes
- **THEN** it SHALL report negative evidence for broker/order/push tables, automatic confirmation rows, automatic rule-application audit events, trading affordances, return guarantee claims and secret/raw prompt leakage
- **AND** any stale historical audit check SHALL be recorded as stale rather than treated as a fresh pass.
