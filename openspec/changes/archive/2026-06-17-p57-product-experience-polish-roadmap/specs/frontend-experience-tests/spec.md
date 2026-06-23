## ADDED Requirements

### Requirement: Product experience polish roadmap governs post-P56 UI work

The frontend product experience SHALL be polished through a staged roadmap before a new final release-ready refresh is claimed.

#### Scenario: Product north star is explicit

- **WHEN** planning post-P56 frontend work
- **THEN** the product MUST be treated as a local investment discipline workbench
- **AND** it MUST NOT be treated as a broker trading terminal, AI chat demo, marketing landing page, or engineering debug console
- **AND** the core daily questions MUST be: can I act today, why, what manual action is needed, and whether data and rules are trustworthy

#### Scenario: Product polish is staged

- **WHEN** post-P56 UI/product improvements are planned
- **THEN** the work MUST be split into independent OpenSpec changes for daily workbench, decision explainability, portfolio/risk/data quality, governance/ops productization, design system/accessibility, and final real UI regression
- **AND** governance/ops productization MUST explicitly include rules, audit, notifications, daily reports, daily auto run, local install, local knowledge, and settings surfaces
- **AND** each stage MUST define scope, out-of-scope safety boundaries, Product Design evidence, browser validation, and subagent review gates

#### Scenario: Release refresh is sequenced after product polish

- **WHEN** P57 product experience roadmap is accepted
- **THEN** release-readiness refresh MUST be deferred until the product polish stages have either completed or been explicitly waived
- **AND** documentation MUST NOT claim that all product design, UI design, or frontend issues are fully fixed before the corresponding stages pass validation

#### Scenario: Safety boundaries remain visible in polished UI

- **WHEN** any polished UI adds or changes a control, page, CTA, state, workflow, or report
- **THEN** it MUST NOT add or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source
- **AND** high risk, unknown, degraded, stale, missing, information-insufficient, and blocked states MUST NOT be styled or worded as ordinary success

#### Scenario: Real UI validation remains required

- **WHEN** a product polish stage changes frontend behavior, layout, information architecture, component primitives, or user-facing copy
- **THEN** validation MUST include frontend unit or component tests, frontend build, relevant backend tests when touched, browser-operated local UI verification, desktop and mobile screenshots, mobile reflow checks, safety copy scans, sensitive information scans, and subagent reviews
- **AND** real LLM validation MUST be included for stages that alter consultation, decision detail, evidence explanation, LLM quality display, or decision-loop surfaces
