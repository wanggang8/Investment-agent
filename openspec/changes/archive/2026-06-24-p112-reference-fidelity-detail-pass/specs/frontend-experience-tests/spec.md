## ADDED Requirements

### Requirement: P112 reference fidelity detail pass blocks archive on residual visual mismatch

The frontend SHALL use the user-approved second reference image as the visual truth for the Investment Agent product UI. All primary and secondary routes SHALL share a rigorous report-cockpit visual language: compact dark navigation, white top status toolbar, tight report hero, clear manual next actions, status metric cards, evidence/checklist modules, progress/ledger surfaces, restrained borders/radii/shadows, and mobile layouts that expose actionable continuation in the first screen.

#### Scenario: P112 page-level reference fidelity gate

- **GIVEN** the application is running locally with the P112 frontend
- **WHEN** all covered desktop routes and key mobile routes are captured as screenshots
- **THEN** each page is compared against the approved reference image for layout rhythm, hierarchy, density, tone, icon treatment, typography, border/radius/shadow, and responsive behavior
- **AND** the page is not considered complete while P0/P1/P2 visual mismatches remain open
- **AND** sub agent review must report no Critical/Important/P0/P1/P2 issues before the change is archived.

#### Scenario: Secondary page report composition

- **GIVEN** a secondary route such as positions, data quality, risk alerts, evidence, governance, notifications, daily reports, local install, local knowledge, or settings
- **WHEN** the page loads at desktop width
- **THEN** the first viewport presents a compact report/status composition rather than a generic old-style card stack
- **AND** the primary status block begins near the top content area rather than being pushed below explanatory banners
- **AND** downstream actions, evidence, checklist, or ledger content is visible or clearly previewed without excessive vertical whitespace.

#### Scenario: Mobile reference efficiency

- **GIVEN** a key route is loaded at approximately 390px width
- **WHEN** the first viewport is inspected
- **THEN** the page shows the status report plus actionable continuation or supporting evidence
- **AND** there is no horizontal overflow, clipped primary content, overlapping text, or oversized hero that hides the rest of the workflow.
