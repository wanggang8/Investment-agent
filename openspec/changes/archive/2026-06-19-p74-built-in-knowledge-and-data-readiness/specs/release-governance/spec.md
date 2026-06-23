## ADDED Requirements

### Requirement: P74 Built-In Knowledge And Data Readiness

P74 SHALL make built-in investment knowledge and required data readiness explicit, auditable, and visible to runtime users.

#### Scenario: Built-in knowledge is structured and bounded

- **GIVEN** the product includes master principles, discipline rules, risk SOPs, and symbol profile knowledge
- **WHEN** P74 exposes built-in knowledge
- **THEN** each entry SHALL have a stable ID, category, summary, applicability, rule mapping, LLM context eligibility, and safety boundary
- **AND** master principles, local notes, and background knowledge SHALL NOT be classified as formal market evidence.

#### Scenario: Data readiness maps required categories to feature impact

- **GIVEN** a user requests readiness for a symbol
- **WHEN** the readiness service evaluates local facts
- **THEN** it SHALL report required and optional data categories as ready, degraded, missing, background-only, or blocked
- **AND** it SHALL map missing or degraded categories to affected product surfaces and claims
- **AND** it SHALL NOT trigger collectors, rebuild indexes, modify source health, write market snapshots, update rules, create notifications, mutate portfolios, or create confirmations.

#### Scenario: Readiness API is sanitized and safe

- **GIVEN** the frontend requests `GET /api/v1/knowledge-readiness?symbol=510300`
- **WHEN** the backend returns readiness data
- **THEN** the response SHALL include overall status, symbol profile, knowledge references, data dependencies, feature impacts, LLM context summary, and safety notes
- **AND** it SHALL NOT expose full prompts, raw HTTP responses, raw LLM responses, private local paths, API keys, private keys, original SQL, or complete account details.

#### Scenario: LLM analysis receives readiness context without decision authority

- **GIVEN** a workflow invokes an LLM analyst node
- **WHEN** relevant knowledge/data readiness context exists
- **THEN** the analyst request SHALL include a sanitized summary of matched principles and data readiness
- **AND** the prompt SHALL state that background knowledge cannot satisfy formal evidence
- **AND** the LLM SHALL remain limited to analysis material and SHALL NOT generate or override the final rule verdict.

#### Scenario: UI shows readiness and gaps

- **GIVEN** readiness data is available
- **WHEN** the user views data quality, rules, consultation, or decision detail surfaces
- **THEN** the UI SHALL distinguish knowledge available as rules, knowledge available as LLM context, background-only knowledge, ready data, degraded data, missing data, and blocked claims
- **AND** it SHALL show safe next steps without broker actions, automatic trading, one-click trading, external push, automatic confirmation, or return promises.

#### Scenario: P74 acceptance covers complete and degraded data scenarios

- **GIVEN** P74 acceptance runs against local test databases and the real UI/API surfaces
- **WHEN** acceptance completes
- **THEN** it SHALL cover a complete `510300` ETF/index path and degraded paths for missing valuation data, background-only local knowledge, single-source evidence, multi-source formal evidence, and out-of-scope capability
- **AND** it SHALL block pass if readiness is only documented but not available through API/UI/LLM-context evidence.

#### Scenario: P74 claims remain bounded

- **GIVEN** P74 passes
- **WHEN** release materials state the result
- **THEN** they MAY claim built-in knowledge and data readiness traceability for the accepted local scope
- **AND** they SHALL NOT claim future investment returns, future market direction, future public-source or model-provider availability, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, automatic database overwrite, paid/login/authorization-gated sources, Level2 data, or high-frequency data.
