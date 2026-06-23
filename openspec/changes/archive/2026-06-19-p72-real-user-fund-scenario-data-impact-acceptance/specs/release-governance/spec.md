## ADDED Requirements

### Requirement: Real user fund scenario data-impact acceptance

After P71, if the project claims the product is ready for practical real-user use, it SHALL execute a P72 real user fund/ETF scenario acceptance that verifies UI operations, API responses, data side effects, auditability, derived readbacks, deterministic accuracy, and safety boundaries.

#### Scenario: Reviewed scenario matrix precedes execution

- **GIVEN** P72 is evaluating practical real-user readiness
- **WHEN** acceptance execution starts
- **THEN** a scenario matrix SHALL exist before execution
- **AND** the matrix SHALL cover real fund setup, holding maintenance, offline transaction recording, local knowledge/RAG, current data, daily discipline, risk alerts, real consultation, manual confirmation, review/readback, failure handling, and safety boundaries
- **AND** execution SHALL NOT proceed if the matrix is known to omit a primary user workflow.

#### Scenario: Data impact is verified after real UI operation

- **GIVEN** a P72 browser scenario performs a user operation
- **WHEN** the operation reports success in the UI
- **THEN** P72 SHALL verify the expected local data impact through API responses and read-only SQLite evidence
- **AND** the acceptance record SHALL include sanitized evidence for the relevant tables and audit events
- **AND** page refresh or navigation SHALL show the same resulting state where the product exposes it.

#### Scenario: Deterministic accuracy is separated from investment prediction

- **GIVEN** P72 evaluates a real fund/ETF scenario
- **WHEN** the product calculates market value, unrealized profit ratio, cash/asset ratios, risk trigger state, report links, or counts
- **THEN** those deterministic values SHALL be checked against independently computed expectations
- **AND** LLM or expected-return material SHALL only be checked for traceability, parse/quality status, rule-consistent final verdict, and risk disclosure
- **AND** P72 SHALL NOT claim future return or market-direction accuracy.

#### Scenario: Safety boundaries remain active in real use

- **GIVEN** P72 runs realistic user scenarios
- **WHEN** pages, API results, browser summaries, and SQLite-derived artifacts are scanned
- **THEN** they SHALL NOT expose complete keys, full prompts, raw provider payloads, private local paths, SQL dumps, broker/order capabilities, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, automatic overwrite, or return promises.

#### Scenario: P72 gaps are recorded honestly

- **GIVEN** a scenario cannot be executed because of provider availability, product behavior, or missing coverage
- **WHEN** P72 writes release or acceptance materials
- **THEN** the result SHALL be blocked or gap-qualified
- **AND** the materials SHALL NOT describe the missing scenario as passed.
