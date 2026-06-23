## ADDED Requirements

### Requirement: P42 workbench aggregates ops and review status safely

The frontend SHALL aggregate risk, rule, review, source health, and runtime readiness states on the user decision workbench without changing their underlying workflows.

#### Scenario: Workbench shows ops and review follow-up

- **WHEN** risk alerts, rule proposals, review summaries, source health, or runtime readiness facts are available
- **THEN** the workbench SHALL show a concise status, safe next step, and navigation path to the authoritative page
- **AND** it SHALL NOT imply automatic repair, automatic rule application, external notification delivery, or trading execution
