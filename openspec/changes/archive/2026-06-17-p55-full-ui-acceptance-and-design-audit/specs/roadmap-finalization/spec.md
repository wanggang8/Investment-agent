## ADDED Requirements

### Requirement: Validate release readiness with real UI operation

The project SHALL support a real frontend UI acceptance pass after command-level release readiness.

#### Scenario: Full UI acceptance is requested

- **GIVEN** command-level acceptance and release handoff materials exist
- **WHEN** full UI acceptance is requested
- **THEN** the project SHALL start the local backend and frontend, operate the UI through a browser, and record page-level results for the major routes
- **AND** command-level smoke alone SHALL NOT be treated as proof that every frontend function has been manually exercised.

### Requirement: Preserve UI acceptance evidence

The project SHALL preserve UI acceptance evidence in release materials.

#### Scenario: A UI route is accepted

- **WHEN** a route or feature is checked through the browser
- **THEN** the acceptance record SHALL include the URL, screenshot or blocker, visible state, key interaction result, safety boundary check, and release impact
- **AND** screenshots SHALL be stored in a local release audit assets folder without complete keys, private paths, raw payloads, or unredacted sensitive data.

### Requirement: Perform Product Design review from captured evidence

The project SHALL review frontend design quality from captured UI evidence.

#### Scenario: Design optimization is assessed

- **WHEN** the frontend UI is reviewed for design quality
- **THEN** the review SHALL use captured screenshots and browser-observed behavior to assess UX, visual hierarchy, consistency, accessibility risks, and optimization opportunities
- **AND** the review SHALL distinguish blocked issues, optimization needs, minor polish, and evidence limits without claiming full WCAG compliance from screenshots alone.
