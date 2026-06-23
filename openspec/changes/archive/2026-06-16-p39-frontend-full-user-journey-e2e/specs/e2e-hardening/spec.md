## ADDED Requirements

### Requirement: P39 Browser Full User Journey Acceptance
The system SHALL provide browser-level E2E coverage for a complete local user journey from first use through daily report review, active consultation, offline confirmation, audit inspection, periodic review, rule governance, and risk alert inspection.

#### Scenario: Empty local state reaches first daily report safely
- **WHEN** the E2E fixture starts from an empty or missing-prerequisite local state
- **THEN** the browser journey SHALL expose safe onboarding or prerequisite guidance for configuration, account, and position setup
- **AND** the journey SHALL reach a first daily discipline report using deterministic local fixture data
- **AND** it SHALL NOT require public network access, real secrets, broker credentials, or manual database inspection

#### Scenario: Consultation to confirmation remains a local record flow
- **WHEN** the browser journey performs an active consultation and opens the resulting decision detail
- **THEN** the page SHALL expose decision trace, evidence, retrieval quality, and audit references where available
- **AND** any confirmation action SHALL record only an offline user fact
- **AND** the journey SHALL NOT expose automatic trading, broker order placement, one-click order placement, or portfolio mutation without user-recorded offline confirmation

#### Scenario: Review and rule governance are inspectable but not automatic
- **WHEN** the browser journey opens periodic review and rule governance surfaces
- **THEN** it SHALL show review summaries, rule proposal status, gatekeeper or final confirmation boundaries, and tracking entrypoints where available
- **AND** pending proposals SHALL remain visible as review/governance facts
- **AND** the journey SHALL NOT automatically apply rules or bypass gatekeeper audit and final user confirmation

#### Scenario: Existing P34 P35 P38 statuses are included
- **WHEN** the browser journey traverses dashboard, evidence, decision, review, rules, and risk alert pages
- **THEN** it SHALL include source health, risk alert/SOP, and retrieval quality states from the existing API/service DTOs
- **AND** degraded, empty, missing, or unknown states SHALL be visible as safe non-success states
