## ADDED Requirements

### Requirement: Real LLM decision details tolerate nullable frontend DTO fields

The frontend SHALL render decision detail pages safely when real LLM-backed decisions contain nullable, missing, or empty list fields in final verdict and trace DTOs.

#### Scenario: Optional actions are null

- **WHEN** a decision detail DTO contains `final_verdict.optional_actions` as `null`
- **THEN** the decision detail page MUST render without a page-level crash
- **AND** the optional-action section MUST show an empty or unavailable safe state instead of calling array-only methods on the value

#### Scenario: Prohibited actions are null or missing

- **WHEN** a decision detail DTO contains `final_verdict.prohibited_actions` as `null` or omits the field
- **THEN** the decision detail page MUST render without a page-level crash
- **AND** the frontend MUST NOT display a false success or false permission to trade

#### Scenario: Real LLM-like fixture is covered by tests

- **WHEN** frontend tests run
- **THEN** they MUST include a real LLM-like decision fixture with nullable verdict list fields
- **AND** they MUST assert that the decision trace, safety boundary, and final verdict remain visible

### Requirement: Productized task-based frontend shell

The frontend SHALL present Investment Agent as a task-based local investment discipline product rather than a flat route list.

#### Scenario: Navigation is grouped by user task

- **WHEN** the app shell renders on desktop
- **THEN** navigation MUST group routes by user task such as today, decision, portfolio, evidence, governance, and system operations
- **AND** all existing primary routes from the P55 route matrix MUST remain reachable

#### Scenario: Mobile navigation does not consume the reading viewport

- **WHEN** the app shell renders on a 390px wide viewport
- **THEN** navigation MUST avoid a permanently visible wide sidebar that compresses main content
- **AND** primary navigation controls MUST remain keyboard and pointer accessible

### Requirement: Productized operational UI system

The frontend SHALL use consistent operational UI primitives for buttons, forms, status labels, cards, and tabular or key-value data displays.

#### Scenario: Forms use consistent field structure

- **WHEN** users operate consultation, positions, local install, local knowledge, or settings forms
- **THEN** labels, hints, inputs, error/success states, and primary/secondary actions MUST use consistent styling and spacing
- **AND** browser-default controls MUST NOT be the dominant visual treatment for critical workflows

#### Scenario: Status and safety states remain visible

- **WHEN** pages show high risk, frozen watch, information insufficient, degraded, unknown, or success states
- **THEN** the UI MUST use consistent semantic styling and Chinese labels
- **AND** unknown or degraded states MUST NOT be styled as ordinary success

### Requirement: Mobile reflow for core acceptance pages

The frontend SHALL reflow core acceptance pages so the page itself does not horizontally overflow on mobile-sized viewports, except for explicitly scoped two-dimensional data containers.

#### Scenario: Positions page reflows on mobile

- **WHEN** `/positions` renders at 390px viewport width
- **THEN** account and holding forms MUST remain visible without page-level horizontal overflow
- **AND** holdings data MUST be readable through stacked cards, key-value rows, or a clearly scoped local table scroller

#### Scenario: Data quality page reflows on mobile

- **WHEN** `/data-quality` renders at 390px viewport width
- **THEN** source health, evidence/RAG, LLM quality, and affected workflow sections MUST remain visible without page-level horizontal overflow
- **AND** long source identifiers, status tokens, and diagnostic text MUST wrap, truncate safely, or scroll only within a scoped local container

### Requirement: Product design evidence is linked to UI acceptance fixes

The P56 implementation SHALL document how Product Design skill guidance and product design research were applied to acceptance-blocking UI fixes.

#### Scenario: Design rationale is traceable

- **WHEN** P56 is reviewed
- **THEN** the change materials MUST include the product brief, design principles, research inputs, and page-level UI plan
- **AND** the implementation report MUST map material UI changes back to those inputs

#### Scenario: Subagent review covers design and safety

- **WHEN** P56 plan review, execution review, or pre-commit review runs
- **THEN** the subagent review MUST check Product Design skill usage, research-backed rationale, mobile usability, real UI acceptance evidence, and prohibited automatic-action boundaries
