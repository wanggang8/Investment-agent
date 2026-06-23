## ADDED Requirements

### Requirement: Provide a user decision workbench

The frontend SHALL provide a user decision workbench that aggregates the existing local daily discipline, portfolio, risk, rule governance, review, audit, and consultation entrypoints into one safe daily-use surface.

#### Scenario: Workbench shows daily priorities

- **WHEN** dashboard and daily discipline report DTOs are available
- **THEN** the workbench SHALL show the current discipline status, final verdict status, missing prerequisites or degraded reasons, and links to the related report or decision
- **AND** it SHALL NOT present executable trading actions

#### Scenario: Workbench handles empty or degraded local facts

- **WHEN** account, market, evidence, source health, LLM, RAG, VecLite, risk, or review data is missing or degraded
- **THEN** the workbench SHALL show a safe status and an inspectable next step where available
- **AND** it SHALL NOT display missing information as successful or imply guaranteed recovery

### Requirement: Workbench entrypoints remain navigational and non-executing

The user decision workbench SHALL only navigate, filter, or show existing local facts unless a later dedicated change explicitly adds a safe local write flow.

#### Scenario: User opens a related workbench entrypoint

- **WHEN** the user selects a workbench link for daily reports, portfolio, risk alerts, rule proposals, review summaries, audit events, or decision consultation
- **THEN** the frontend SHALL navigate to the relevant page or filtered record
- **AND** it SHALL NOT submit confirmations, apply rules, place orders, call broker APIs, or send external pushes

#### Scenario: Consultation is available from the workbench

- **WHEN** the user chooses the active consultation entrypoint
- **THEN** the workbench SHALL navigate to the existing consultation page or form with safe explanatory copy
- **AND** it SHALL NOT automatically submit a consultation request or treat LLM material as final verdict authority

### Requirement: Workbench uses supported DTOs only

The workbench SHALL use supported API/service DTOs and frontend mappers rather than direct local storage access.

#### Scenario: Workbench renders its panels

- **WHEN** the workbench renders daily, portfolio, risk, rule, review, audit, settings, or diagnostics information
- **THEN** the data SHALL come from API/service DTOs
- **AND** the frontend SHALL NOT read SQLite, VecLite, local logs, diagnostic files, or private config files directly
