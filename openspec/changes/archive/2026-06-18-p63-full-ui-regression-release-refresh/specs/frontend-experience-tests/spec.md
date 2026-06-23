## ADDED Requirements

### Requirement: Full real UI regression refreshes the release status

After product experience polishing is complete, the project SHALL execute a full real UI regression and release status refresh before claiming a current release candidate state.

#### Scenario: Full UI regression covers all primary routes

- **WHEN** P63 executes browser acceptance against the local backend and frontend
- **THEN** Dashboard, Workbench, Consultation, Decision Detail, Evidence, Decision Loop, Positions, Data Quality, Risk Alerts, Risk Alert Detail, Rules, Audit, Notifications, Daily Reports, Daily Report Detail, Daily Auto Run, Review, Local Install, Local Knowledge, and Settings MUST be operated through the UI or covered by an equivalent browser assertion
- **AND** each route MUST record whether its primary status, key actions, empty/degraded/error states, and local-only safety boundaries are visible and usable
- **AND** console errors, page errors, and failed API responses MUST be recorded or explicitly ruled non-blocking with reasons

#### Scenario: Real LLM consultation is verified or classified

- **WHEN** P63 runs the consultation journey with a real LLM-backed configuration
- **THEN** the UI MUST submit a real consultation request and attempt to open the resulting decision detail
- **AND** the acceptance record MUST state whether LLM analysis was returned, parsed, quality-gated, and displayed without letting the model write the final rule verdict
- **AND** network, rate limit, authentication, model unavailable, parse, or quality failures MUST be classified and mapped to release impact instead of being treated as success

#### Scenario: Release candidate status is refreshed from current evidence

- **WHEN** P63 produces release materials
- **THEN** the release candidate MUST reference the current acceptance run and current code-under-test commit
- **AND** the status MUST be one of `release_ready`, `release_degraded`, or `blocked`
- **AND** degraded, skipped, retried, or waived gates MUST list their reason, artifact, and release impact
- **AND** the handoff MUST include Not Claimed boundaries for returns, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, future provider availability, login sources, paid sources, authorized sources, Level2 data, and high-frequency data

#### Scenario: Full UI regression preserves safety and redaction boundaries

- **WHEN** P63 scans UI text, release materials, logs, browser evidence, and committed assets
- **THEN** it MUST NOT find new user-facing broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source claims
- **AND** committed artifacts MUST NOT contain complete API keys, private paths, raw SQL dumps, raw stack traces, complete prompts, local database paths, raw vendor payloads, Playwright trace archives, temporary SQLite databases, or unredacted local logs
