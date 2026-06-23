## ADDED Requirements

### Requirement: Daily workbench presents an at-a-glance investment discipline cockpit

The frontend SHALL make the dashboard and workbench answer the daily investment discipline questions before secondary details.

#### Scenario: Dashboard first screen shows daily decision state

- **WHEN** the dashboard route `/` renders with dashboard and daily discipline data
- **THEN** the first screen MUST show the current verdict or safe unavailable state, status tone, data trust summary, last update context, prohibited actions, optional manual actions, and next manual actions
- **AND** the first screen MUST provide local navigation to decision detail or daily report when such links are available

#### Scenario: Workbench first screen shows task queue

- **WHEN** the workbench route `/workbench` renders
- **THEN** it MUST use the same daily state model as the dashboard
- **AND** it MUST present a prioritized manual action queue before secondary portfolio, risk, rule, review, or consultation sections
- **AND** each action MUST be a local navigation or manual review action, not an execution action

#### Scenario: Daily workbench handles degraded and insufficient states safely

- **WHEN** dashboard, daily report, portfolio, risk, rule, or review data is unavailable, degraded, stale, high risk, unknown, or insufficient
- **THEN** dashboard and workbench MUST show safe Chinese status text and a clear next manual step
- **AND** they MUST NOT style or describe degraded, unknown, high risk, stale, missing, or insufficient states as ordinary success

#### Scenario: Daily workbench remains mobile readable

- **WHEN** `/` or `/workbench` renders at 390px viewport width
- **THEN** the daily status, manual action queue, and primary task links MUST remain visible without page-level horizontal overflow
- **AND** screenshots MUST be captured for desktop and mobile validation

#### Scenario: Daily workbench preserves safety boundaries

- **WHEN** dashboard or workbench UI adds or changes a CTA, status, summary, or task link
- **THEN** it MUST NOT add or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source
