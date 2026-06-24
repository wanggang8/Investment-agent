## ADDED Requirements

### Requirement: P113 layout fidelity pass blocks archive on visible layout defects

The frontend SHALL preserve the P111/P112 reference-cockpit visual direction while eliminating visible layout defects across desktop and 390px mobile routes. A route SHALL NOT be considered complete when primary content is clipped, horizontally overflowing, overlapped, visibly misaligned, excessively compressed, or presented with touch targets too small for practical mobile use.

#### Scenario: Full-route desktop and mobile layout gate

- **GIVEN** the application is running locally with the P113 frontend
- **WHEN** all covered routes are captured at desktop and approximately 390px mobile widths
- **THEN** each screenshot is checked for horizontal overflow, clipping, overlap, compressed cards, unstable text wrapping, too-small action targets, and first-viewport report hierarchy
- **AND** any P0/P1/P2 layout defect must be fixed and recaptured before archive.

#### Scenario: Mobile metric and action layout

- **GIVEN** a route uses report hero metrics, status cards, compact actions, or next-step links
- **WHEN** the route is loaded at approximately 390px width
- **THEN** primary metrics do not require horizontal scrolling
- **AND** cards and labels remain inside the viewport
- **AND** actionable links or buttons have a visibly tappable target and do not collapse into tiny text.

#### Scenario: Secondary page polish

- **GIVEN** a secondary route such as data quality, settings, local install, local knowledge, decision detail, rules, audit, notifications, daily reports, or daily auto run
- **WHEN** the first viewport is inspected
- **THEN** the page prioritizes a clear report/status/next-action composition
- **AND** raw engineering content, long JSON, file paths, or dense diagnostic payloads do not dominate the first user-facing layer.
