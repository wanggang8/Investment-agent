## ADDED Requirements

### Requirement: Visual system redesign preserves discipline-first product semantics

The frontend SHALL upgrade the product visual system after acceptance while preserving the local investment discipline workflow, safety boundaries, and existing API contracts.

#### Scenario: Visual direction is selected before production UI implementation

- **WHEN** P110 begins visual system redesign
- **THEN** the project MUST generate three distinct visual direction options for the core workbench experience before modifying production UI code
- **AND** the selected option MUST become the visual target for implementation
- **AND** unselected direction images MUST NOT be treated as product requirements or runtime evidence

#### Scenario: Core workbench remains decision-first and non-trading

- **WHEN** Dashboard or Workbench is redesigned
- **THEN** the first screen MUST still prioritize current discipline status, prohibited actions, allowed manual actions, data trust, risk summary, and evidence/rule traceability
- **AND** it MUST NOT prioritize return fantasy, market excitement, price-board behavior, or trading execution
- **AND** all actions MUST remain local navigation, local maintenance, local readback, offline record, or manual review actions

#### Scenario: Visual tokens improve scanability without hiding states

- **WHEN** visual tokens, spacing, navigation, cards, panels, tables, forms, or badges are updated
- **THEN** success, warning, danger, degraded, unknown, readonly, blocked, first-use, insufficient-data, frozen-watch, and high-risk states MUST retain readable text labels
- **AND** degraded, unknown, readonly, blocked, missing, stale, failed, and information-insufficient states MUST NOT be styled or worded as ordinary success
- **AND** status meaning MUST NOT rely on color alone

#### Scenario: Redesign keeps frontend contract boundaries

- **WHEN** P110 changes pages or shared UI components
- **THEN** pages MUST continue using existing frontend services or shared API clients
- **AND** P110 MUST NOT require new backend API fields, SQLite schema changes, Eino workflow changes, LLM prompt changes, data source changes, or rule engine changes
- **AND** raw SQLite, VecLite files, local config files, localStorage, sessionStorage, private paths, prompts, API keys, SQL, raw vendor payloads, and stack traces MUST NOT become visible UI dependencies

#### Scenario: Responsive visual QA covers core routes

- **WHEN** P110 implementation is complete
- **THEN** core routes MUST be captured or equivalently checked at 390px, 768px, and 1280px viewport widths
- **AND** primary status, manual action queue, forms, summaries, explanation panels, tables, timelines, diagnostics, empty states, error states, and navigation MUST remain readable without page-level horizontal overflow
- **AND** two-dimensional tables, logs, JSON, and diagnostic text MAY scroll only inside clearly scoped local containers

#### Scenario: Safety copy boundaries survive redesign

- **WHEN** redesigned UI is scanned
- **THEN** it MUST NOT add or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source
- **AND** any copy about suggested actions MUST make clear that final actions are local, manual, offline, or read-only as applicable
