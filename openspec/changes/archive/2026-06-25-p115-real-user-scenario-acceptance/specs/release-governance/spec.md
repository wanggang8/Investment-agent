## ADDED Requirements

### Requirement: P115 Real User Scenario Acceptance
The project SHALL maintain a repeatable real-user-scenario acceptance gate that verifies broad product journeys through visible UI entry points, HTTP APIs, SQLite side effects, downstream readback, audit traceability, degradation behavior, and forbidden automation absence.

#### Scenario: Scenario matrix covers visible product operations
- **GIVEN** P115 validation is prepared
- **WHEN** the P115 scenario matrix is written
- **THEN** it SHALL cover first launch, empty account onboarding, portfolio initialization, holding maintenance, batch import, offline transactions, local fact correction, rebalance review, consultation, decision detail, manual confirmation, decision error marking, decision loop, evidence refresh, RAG/index readiness, local knowledge import, market refresh, data-quality resolution, risk lifecycle, rule governance, notifications, daily reports, daily auto-run status, dashboard/workbench summaries, review, audit, settings updates, forbidden settings-based rule/SOP mutation, API diagnostics, local install diagnostics, mobile operation, browser-level interaction parity, and failure/degradation cases
- **AND** each scenario SHALL specify expected UI/browser or API evidence, SQLite readback, downstream readback, and safety negative evidence
- **AND** actual pass/fail status SHALL remain pending until execution evidence exists.

#### Scenario: Runner records scenario-level evidence
- **GIVEN** the repository source tree is available locally
- **WHEN** the P115 runner is executed
- **THEN** it SHALL use an isolated temporary SQLite database and local backend/frontend where required
- **AND** it SHALL produce a scenario-level summary with statuses such as `fresh_pass`, `scoped_pass`, `degraded_expected`, or `blocked`
- **AND** it SHALL record config mode, runtime mode, stub/provider/LLM mode, HTTP method/path/status/request ids, SQLite table/field before-after readback, downstream page or endpoint, screenshots or DOM checks when used, console errors, redaction results, and safety counters
- **AND** it SHALL preserve separate evidence for API/SQLite, browser, and degradation layers.

#### Scenario: Safety boundary remains explicit
- **GIVEN** P115 validation has completed
- **WHEN** the acceptance record is written
- **THEN** it SHALL fail or explicitly block release claims if broker/order/push capability, automatic trading, automatic confirmation, automatic rule application, or return-guarantee evidence is present
- **AND** it SHALL not treat stub, fixture, deterministic-local, single-symbol, or degraded-provider evidence as proof of future external provider availability
- **AND** P104-derived local seeded evidence SHALL be labeled as local functional linkage evidence rather than external provider or real LLM evidence.
