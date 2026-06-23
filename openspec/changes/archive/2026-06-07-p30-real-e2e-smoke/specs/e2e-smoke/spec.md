## ADDED Requirements

### Requirement: Local real-environment E2E smoke shall validate critical read-only UI/API paths

The system SHALL provide a local E2E smoke verification path that starts or targets the local backend and frontend, uses temporary local state, and verifies that critical read-only user paths are reachable without frontend runtime crashes.

#### Scenario: Smoke verifies local health and UI availability

- **GIVEN** a local test configuration with temporary SQLite storage and no real secrets
- **WHEN** the E2E smoke runs against the local backend and frontend
- **THEN** it SHALL verify the backend health endpoint is reachable
- **AND** it SHALL verify at least one primary frontend route renders successfully
- **AND** it SHALL NOT require broker APIs, paid data sources, login-only sources, or external notification channels

#### Scenario: Smoke covers decision trace and expected-return display safety

- **GIVEN** local decision data that may include empty expected-return scenarios or degraded analysis materials
- **WHEN** the E2E smoke opens the decision-related UI path
- **THEN** the page SHALL render without JavaScript runtime errors caused by null/empty array fields
- **AND** the smoke SHALL verify expected-return, scenario, sell-evaluation, or degradation text is visible when fixture data provides it
- **AND** the smoke SHALL NOT treat expected-return analysis as a trading instruction or guaranteed outcome

#### Scenario: Smoke covers evidence or audit read-only visibility

- **GIVEN** local evidence, source verification, audit, or degraded-source state produced from controlled test data
- **WHEN** the E2E smoke opens the corresponding UI or API path
- **THEN** the smoke SHALL verify that the read-only state is visible or queryable
- **AND** source degradation SHALL be displayed as diagnostic information rather than silently converted into fabricated data

### Requirement: E2E smoke shall avoid persistent local artifact pollution

The system SHALL keep Playwright/browser smoke artifacts, traces, screenshots, logs, temporary databases, and MCP browser logs out of the committed working tree by default.

#### Scenario: Smoke cleanup leaves the working tree clean of generated artifacts

- **GIVEN** the E2E smoke has completed successfully or failed after collecting diagnostics
- **WHEN** a developer checks `git status --short`
- **THEN** generated smoke artifacts such as `.playwright-mcp/`, Playwright output, temporary SQLite files, screenshots, traces, and logs SHALL NOT appear as untracked files unless explicitly requested for debugging

#### Scenario: Smoke remains inside project safety boundaries

- **GIVEN** the E2E smoke runs in a local development environment
- **WHEN** it exercises backend tasks, frontend pages, or browser automation
- **THEN** it SHALL NOT initiate brokerage operations, automatic trading, external push notifications, credentialed scraping, paid-data access, high-frequency crawling, or Level2/user-identity market data access
