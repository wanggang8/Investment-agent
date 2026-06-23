## ADDED Requirements

### Requirement: Design system primitives and accessibility gates are standardized

The frontend SHALL provide reusable UI primitives and browser-level accessibility/reflow gates for the productized local investment discipline workbench before the final full UI regression refresh.

#### Scenario: Shared UI primitives expose accessible semantics

- **WHEN** P62 introduces Button, Field, StatusBadge, PageHeader, SummaryCard, DetailSection, ResponsiveTable, EmptyState, and ErrorState primitives
- **THEN** each primitive MUST expose stable text or accessible names for its user-visible purpose
- **AND** form controls MUST be associated with labels, hints, and errors when those are present
- **AND** collapsible details MUST use button semantics and maintain `aria-expanded`
- **AND** table-like content MUST provide a caption or accessible name and preserve column meaning under mobile reflow

#### Scenario: Status tokens remain consistent and impossible to confuse with success

- **WHEN** pages render success, warning, danger, degraded, unknown, readonly, or blocked states
- **THEN** those states MUST use consistent frontend tone names
- **AND** every state MUST include readable text, not only color
- **AND** degraded, unknown, readonly, blocked, missing, stale, failed, and information-insufficient states MUST NOT be styled or worded as ordinary success

#### Scenario: Keyboard paths cover primary local workflows

- **WHEN** the local frontend is operated with a keyboard
- **THEN** the primary navigation, mobile menu, representative forms, collapsible detail sections, and critical local-only buttons MUST be reachable and operable
- **AND** focus indicators MUST remain visible on desktop, tablet-width, and mobile-width layouts
- **AND** disabled or working controls MUST expose their state through text, attributes, or equivalent accessible semantics

#### Scenario: Reflow and visual evidence cover desktop, tablet-width, and mobile layouts

- **WHEN** representative P58-P61 pages render at 390px, 768px, and 1280px viewport widths
- **THEN** primary status, next manual actions, forms, summaries, tables, timelines, diagnostics, empty states, error states, and navigation MUST remain readable without page-level horizontal overflow
- **AND** two-dimensional tables, JSON, logs, or diagnostic text MAY scroll only inside clearly scoped local containers
- **AND** screenshots or equivalent browser evidence MUST be captured for the three viewport classes

#### Scenario: Design system hardening preserves local-only safety boundaries

- **WHEN** P62 changes components, page layout, UI text, keyboard behavior, or validation evidence
- **THEN** it MUST NOT add or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source
- **AND** UI output MUST NOT render complete API keys, private paths, SQL, raw stack traces, complete prompts, local database paths, or raw vendor payloads
