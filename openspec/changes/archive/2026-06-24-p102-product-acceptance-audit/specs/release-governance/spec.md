## ADDED Requirements

### Requirement: P102 product acceptance audit

The project SHALL support a product-level acceptance audit after real LLM access is configured.

#### Scenario: Product audit evidence is captured

- **GIVEN** local backend, frontend, SQLite, VecLite, and real LLM config are available
- **WHEN** P102 product acceptance is executed
- **THEN** the audit SHALL capture current-run screenshots for key product workflows
- **AND** it SHALL assess UX, design reasonableness, accessibility risks, data/readback trust, and safety boundaries.

#### Scenario: Release claims remain bounded

- **GIVEN** P102 writes product acceptance findings
- **WHEN** release readiness is described
- **THEN** the project MAY claim product-level local-source acceptance only for the checked local runtime scope
- **AND** it SHALL NOT claim Docker installation, package distribution, GitHub Release, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.
