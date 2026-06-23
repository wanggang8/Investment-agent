## ADDED Requirements

### Requirement: Governance and ops pages present productized workbench experiences

The frontend SHALL turn rules, audit, notifications, daily reports, daily auto run, local install, local knowledge, and settings into readable governance and operations surfaces while preserving existing backend contracts and local-only safety boundaries.

#### Scenario: Rules page presents governance status before raw rule details

- **WHEN** the user opens `/rules`
- **THEN** the first screen MUST show current rule version, proposal counts, pending user confirmation, pending final confirmation, gatekeeper or validation risk, and next manual governance actions
- **AND** rule proposals MUST show reason, sample count, overfit risk, validation status, guardrail decision, audit summary, related local records, and explicit manual confirmation boundaries
- **AND** raw rule JSON or threshold details MUST NOT dominate the first screen when structured summaries are available
- **AND** the page MUST NOT imply automatic rule application, automatic confirmation, broker connectivity, external push, or trading

#### Scenario: Audit page presents an operational timeline with summary context

- **WHEN** the user opens `/audit`
- **THEN** the first screen MUST show audit event count, recent activity, important event categories, and next local inspection actions before the detailed timeline
- **AND** audit events MUST remain traceable to existing API DTO fields without reading SQLite, local files, VecLite, or raw logs
- **AND** empty, degraded, unknown, and error states MUST be visible and safe

#### Scenario: Notifications page behaves as a local inbox

- **WHEN** the user opens `/notifications`
- **THEN** the first screen MUST show unread count, severity distribution, source categories, local processing status, and next manual actions
- **AND** mark-read controls MUST be described as local application state only
- **AND** the page MUST NOT promise SMS, email, webhook, third-party notification delivery, external push, automatic confirmation, or trading

#### Scenario: Daily reports and daily auto run explain discipline and local runtime health

- **WHEN** the user opens `/daily-discipline/reports` or `/daily-auto-run`
- **THEN** the first screen MUST show current discipline or runtime state, evidence or execution coverage, degraded or missing prerequisites, recent activity, and next manual checks
- **AND** daily auto run MUST distinguish disabled, scheduled, running, success, degraded, failed, and unknown states without styling degraded or unknown as normal success
- **AND** daily auto run diagnostics MUST guide manual recheck and MUST NOT promise automatic repair, automatic source refresh, automatic confirmation, automatic rule application, database overwrite, or trading

#### Scenario: Local install, local knowledge, and settings share safe configuration and diagnostics patterns

- **WHEN** the user opens `/local-install`, `/local-knowledge`, or `/settings`
- **THEN** the page MUST organize configuration, diagnostic status, previews, summaries, and next manual actions into readable sections
- **AND** sensitive values, API keys, private paths, SQL, raw stack traces, complete prompts, and raw vendor payloads MUST NOT be rendered
- **AND** existing local write actions such as knowledge import confirmation or market refresh MUST remain explicit local facts or local refreshes and MUST NOT imply trading, rule application, or external delivery

#### Scenario: Governance and ops experiences remain mobile readable

- **WHEN** P61 pages render at 390px viewport width
- **THEN** primary status, next manual actions, safety boundaries, forms, inbox cards, timelines, diagnostic summaries, and navigation MUST remain readable without page-level horizontal overflow
- **AND** two-dimensional data MAY scroll only inside clearly scoped local containers
- **AND** screenshots or browser evidence MUST be captured for desktop and mobile validation

