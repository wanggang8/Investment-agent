## ADDED Requirements

### Requirement: P8 review frontend safe-state coverage
The frontend SHALL test localized status display, unknown status fallback, safe error states, empty states, and no automatic trading affordances across P0–P8 user flows.

#### Scenario: Localized status coverage
- **WHEN** dashboard, evidence, portfolio, audit, market, rule proposal, or decision statuses are rendered
- **THEN** known values SHALL use Chinese display text and unknown values SHALL use “未知状态”.

#### Scenario: Error and empty page states are visible
- **WHEN** API calls fail or return empty successful data
- **THEN** pages SHALL show safe user-facing states instead of blank or misleading content.

#### Scenario: Confirmation failure is not shown as success
- **WHEN** submitting a user confirmation fails
- **THEN** the page SHALL retain the previous decision state and SHALL NOT show success copy.
