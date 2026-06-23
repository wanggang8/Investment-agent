## ADDED Requirements

### Requirement: Review page displays periodic summaries and tracking
The frontend SHALL display periodic review summaries, rule suggestions, and tracking entrypoints using API/service DTOs rather than direct local storage access.

#### Scenario: Periodic summary is visible
- **WHEN** monthly or quarterly review data is available
- **THEN** the review page shows the period summary, relevant audit status, and supporting counts.

#### Scenario: Rule suggestions are visible but not applied automatically
- **WHEN** a review produces rule suggestions
- **THEN** the frontend displays the suggestions as review output or rule proposal entrypoints and does not present automatic rule application behavior.

#### Scenario: Tracking entrypoint is available
- **WHEN** a review summary references audit events, rule proposals, or error cases
- **THEN** the frontend provides a visible path to inspect the related tracking records.
