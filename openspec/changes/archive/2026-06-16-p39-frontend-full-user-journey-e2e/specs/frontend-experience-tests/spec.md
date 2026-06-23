## ADDED Requirements

### Requirement: P39 Browser UX Stability Checks
The frontend SHALL add browser-level stability checks for key local journeys, including console error capture, unhandled rejection capture, narrow viewport smoke, and basic accessibility-oriented assertions.

#### Scenario: Key pages have no unexpected browser errors
- **WHEN** Playwright opens dashboard, portfolio, evidence, decision detail, audit/review, rules, risk alerts, daily discipline report, and settings pages in the P39 fixture
- **THEN** the test SHALL fail on unexpected console errors or unhandled page errors
- **AND** allowed diagnostic logs, if any, SHALL be explicitly scoped and documented in the test fixture

#### Scenario: Narrow viewport keeps primary controls usable
- **WHEN** key pages are rendered under a narrow mobile-like viewport
- **THEN** primary navigation, status labels, form controls, and action buttons SHALL remain visible and non-overlapping
- **AND** critical labels SHALL NOT be hidden in a way that changes the safety meaning of a page

#### Scenario: Basic accessibility expectations are covered
- **WHEN** forms, navigation, and interactive controls are rendered in the P39 browser journey
- **THEN** controls SHALL have accessible names, form inputs SHALL have labels or equivalent accessible descriptions, and navigation landmarks or equivalent page structure SHALL be discoverable
- **AND** these checks SHALL rely on browser-visible semantics rather than direct local file or SQLite reads

#### Scenario: Vitest and Playwright remain separated
- **WHEN** frontend verification runs
- **THEN** Vitest SHALL continue to cover component and mapper behavior
- **AND** Playwright SHALL cover browser journeys with deterministic fixture data
- **AND** the two suites SHALL avoid collecting each other's files or sharing mutable persistent test state
