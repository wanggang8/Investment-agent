# roadmap-finalization Specification

## Purpose

This specification records roadmap governance rules from the P33-P40 feature queue through the post-P44 planning refresh, and defines how dependencies, historical traceability, future candidates, and safety boundaries should be handled before implementation.
## Requirements
### Requirement: Maintain the P33-P40 roadmap as the remaining planned feature queue

The project SHALL treat P33 through P40 as the current remaining planned feature queue after P32 archive, unless a later OpenSpec governance change explicitly revises the roadmap.

#### Scenario: Next feature work is selected after P32 archive

- **GIVEN** P32 has been archived and no earlier active change remains
- **WHEN** the team selects the next feature implementation stage
- **THEN** the selected stage SHALL come from P33 through P40
- **AND** the selected stage SHALL have its own OpenSpec change before implementation
- **AND** the system SHALL NOT skip required governance by directly editing L1 contracts

#### Scenario: Work outside P33-P40 is requested

- **GIVEN** a requested enhancement is not covered by P33 through P40
- **WHEN** the team decides to pursue it
- **THEN** the project SHALL create a separate governance change or future roadmap change
- **AND** it SHALL NOT be silently folded into an unrelated P33-P40 stage

### Requirement: Preserve roadmap execution dependencies

The project SHALL preserve the documented execution dependencies between P33 through P40 so downstream stages have the required product, data, and verification foundations.

#### Scenario: User journey validation depends on account initialization

- **GIVEN** P39 requires a real user journey from empty database to the first daily discipline report
- **WHEN** P39 is planned
- **THEN** P33 account and holdings initialization capability SHALL be completed or explicitly scoped as a prerequisite

#### Scenario: Risk and evidence quality work depends on data coverage

- **GIVEN** P35 risk warnings and P38 retrieval quality depend on reliable data and evidence inputs
- **WHEN** P35 or P38 is planned
- **THEN** P34 data coverage gaps and freshness/failure semantics SHALL be considered as prerequisites or documented assumptions

### Requirement: Keep non-feature governance work separate from P33-P40 implementation

The project SHALL keep P19-P24 historical archive traceability and P40-after roadmap expansion separate from P33-P40 feature implementation.

#### Scenario: Historical archive traceability is requested

- **GIVEN** P19 through P24 are marked delivered but do not have full archive packages
- **WHEN** audit-grade traceability is required
- **THEN** the project SHALL create a dedicated governance change
- **AND** it SHALL NOT rewrite or fabricate historical archive packages

#### Scenario: P40-after product work is requested

- **GIVEN** P40 is the last stage in the current roadmap
- **WHEN** new product work beyond P40 is requested
- **THEN** the project SHALL create a new roadmap or proposal change
- **AND** it SHALL define scope, dependencies, and acceptance before implementation

### Requirement: Govern all post-P40 product work through a new roadmap change

The project SHALL treat P40 as the end of the completed P33-P40 planned feature queue and SHALL require a new roadmap or governance change before any post-P40 product implementation begins.

#### Scenario: Post-P40 product work is requested

- **GIVEN** P40 has been archived and no active implementation change exists
- **WHEN** the team chooses a new product direction
- **THEN** the project SHALL first create a roadmap or governance OpenSpec change
- **AND** the roadmap SHALL define candidate stages, dependencies, out-of-scope boundaries, and acceptance criteria before implementation

#### Scenario: A concrete post-P40 feature is selected

- **GIVEN** a post-P40 roadmap or governance change has identified candidate work
- **WHEN** the team decides to implement one candidate feature
- **THEN** that feature SHALL receive its own OpenSpec change with `proposal.md`, `tasks.md`, and relevant delta files
- **AND** implementation SHALL follow the selected change instead of editing contracts directly

### Requirement: Separate historical audit traceability from new product features

The project SHALL keep P19-P24 historical archive traceability separate from post-P40 product feature implementation.

#### Scenario: P19-P24 archive traceability is requested after P40

- **GIVEN** P19 through P24 are delivered but do not have full archive packages
- **WHEN** audit-grade traceability is requested
- **THEN** the project SHALL create a dedicated governance change for historical traceability
- **AND** it SHALL NOT rewrite or fabricate historical archive packages
- **AND** it SHALL NOT block unrelated post-P40 product roadmap planning unless explicitly prioritized

### Requirement: Preserve safety boundaries in post-P40 roadmap planning

The project SHALL preserve the existing investment safety boundaries unless a dedicated future governance change explicitly reopens them with review.

#### Scenario: Post-P40 roadmap candidates are evaluated

- **WHEN** roadmap candidates are added or prioritized
- **THEN** the roadmap SHALL keep broker APIs, automatic trading, external push, automatic rule application, paid or login-only data sources, Level2 data, high-frequency sources, return promises, and deterministic price predictions out of scope by default
- **AND** LLM usage SHALL remain limited to analysis material rather than final rule verdicts

### Requirement: Govern all post-P44 work through a refreshed roadmap change

The project SHALL treat P44 as the end of the P41-selected post-P40 queue and SHALL require a refreshed roadmap or governance OpenSpec change before any post-P44 implementation begins.

#### Scenario: Post-P44 work is requested

- **GIVEN** P42, P43, P44, and P19-P24 historical traceability have been archived
- **AND** no active implementation change exists
- **WHEN** the team chooses a new direction
- **THEN** the project SHALL first create a roadmap or governance OpenSpec change
- **AND** the roadmap SHALL define candidate stages, dependencies, out-of-scope boundaries, and acceptance criteria before implementation

#### Scenario: A concrete post-P44 feature is selected

- **GIVEN** a post-P44 roadmap or governance change has identified candidate work
- **WHEN** the team decides to implement one candidate feature
- **THEN** that feature SHALL receive its own OpenSpec change with `proposal.md`, `design.md`, `tasks.md`, and relevant delta files
- **AND** implementation SHALL follow the selected change instead of editing contracts directly

### Requirement: Prioritize post-P44 candidates by local safety and verifiability

The project SHALL prefer post-P44 candidates that improve local-only safety, explainability, data quality, and operability without expanding trading or external delivery capabilities.

#### Scenario: Post-P44 candidates are evaluated

- **WHEN** roadmap candidates are added or prioritized
- **THEN** the roadmap SHALL classify them by category, suggested change id, dependencies, acceptance approach, and suitability
- **AND** the roadmap SHALL recommend a next candidate while keeping implementation out of the roadmap change

### Requirement: Preserve safety boundaries in post-P44 roadmap planning

The project SHALL preserve the existing investment safety boundaries unless a dedicated future governance change explicitly reopens them with review.

#### Scenario: Post-P44 roadmap candidates are evaluated

- **WHEN** roadmap candidates are added or prioritized
- **THEN** the roadmap SHALL keep broker APIs, automatic trading, one-click trading, external push, automatic confirmation, automatic rule application, automatic repair promises, paid, login-only, or authorization-gated data sources, Level2 data, high-frequency sources, return promises, and deterministic price predictions out of scope by default
- **AND** LLM usage SHALL remain limited to analysis material rather than final rule verdicts

### Requirement: Prioritize P19-P24 audit evidence before release materials

The project SHALL prioritize a dedicated P19-P24 audit evidence phase before preparing release candidate materials after P49.

#### Scenario: Post-P49 work is selected

- **GIVEN** P49 has been archived
- **AND** P19-P24 are marked delivered but do not have per-phase complete archive packages
- **WHEN** the team selects the next governance or product stage
- **THEN** the next recommended stage SHALL be P51 `p51-p19-p24-audit-evidence-pack`
- **AND** that stage SHALL produce current-state audit evidence instead of fabricating historical archive packages
- **AND** release candidate materials SHALL remain deferred until the P19-P24 audit evidence and project acceptance gate matrix have been completed.

### Requirement: Define project acceptance gates before release candidate work

The project SHALL define a layered acceptance gate matrix before release candidate materials are prepared.

#### Scenario: Project acceptance is planned

- **WHEN** acceptance planning is performed after the P19-P24 audit evidence phase
- **THEN** the plan SHALL cover unit tests, integration tests, E2E tests, real-source tests, real-LLM tests, local smoke tests, install diagnostics, release-upgrade checks, and safety boundary checks
- **AND** each gate SHALL define command or workflow entry, prerequisites, pass criteria, allowed degradation, artifacts, and whether failure blocks release
- **AND** real-source and real-LLM tests SHALL require explicit opt-in configuration and classified failure reporting.

### Requirement: Keep release candidate materials dependent on audit and acceptance readiness

The project SHALL treat release candidate material preparation as dependent on both historical audit evidence and acceptance gate readiness.

#### Scenario: Release candidate materials are requested

- **GIVEN** P51 and P52 are not both completed
- **WHEN** release candidate materials or release packaging are requested
- **THEN** the project SHALL first complete or explicitly waive the missing audit or acceptance phase in a governance change
- **AND** the project SHALL NOT claim release readiness from P49 alone.

### Requirement: Maintain a current-state P19-P24 audit evidence pack

The project SHALL provide a current-state audit evidence pack for P19-P24 before project acceptance gates or release candidate materials are prepared.

#### Scenario: P19-P24 audit evidence is reviewed

- **GIVEN** P19-P24 are marked delivered but do not have per-phase complete archive packages
- **WHEN** audit evidence is requested for these phases
- **THEN** the project SHALL provide a document that lists each phase's delivery boundary, archive status, documentation evidence, code evidence, test evidence, rerunnable commands, residual gaps, and claims that must not be made
- **AND** the document SHALL state that it is not a replacement for missing historical archive packages.

### Requirement: Preserve historical integrity in P19-P24 audit evidence

The project SHALL NOT fabricate P19-P24 historical archive packages, completion timestamps, or unverified implementation claims when creating audit evidence.

#### Scenario: Missing archive packages are documented

- **WHEN** a P19-P24 phase lacks a complete archive package
- **THEN** the audit evidence SHALL use missing or not present archive language
- **AND** it SHALL cite current repository evidence instead of backfilled historical tasks or fabricated archive paths.

### Requirement: Keep P19-P24 audit evidence separate from acceptance gates

The P19-P24 audit evidence pack SHALL remain separate from the project acceptance gate matrix.

#### Scenario: Acceptance planning follows the audit evidence pack

- **GIVEN** the P19-P24 audit evidence pack has been completed
- **WHEN** the next stage is selected
- **THEN** the next recommended stage SHALL be P52 `p52-project-acceptance-gate-matrix`
- **AND** P52 SHALL convert the evidence into release-blocking and non-blocking validation gates instead of treating P51 alone as release readiness.

### Requirement: Maintain a project acceptance gate matrix before release materials

The project SHALL provide a project acceptance gate matrix before release candidate materials are prepared.

#### Scenario: Acceptance readiness is reviewed

- **GIVEN** P51 has provided the P19-P24 audit evidence pack
- **WHEN** project acceptance readiness is reviewed
- **THEN** the project SHALL provide a matrix covering governance checks, Go tests, integration tests, frontend tests and build, E2E smoke, local smoke, real-source opt-in tests, real-LLM opt-in tests, install diagnostics, release-upgrade checks, and safety boundary checks
- **AND** each gate SHALL define command or workflow entry, prerequisites, pass criteria, allowed degradation, artifacts, and whether failure blocks release.

### Requirement: Classify real-source and real-LLM acceptance failures

The project SHALL classify real-source and real-LLM acceptance failures instead of treating every failure as the same release outcome.

#### Scenario: Opt-in real acceptance test fails

- **WHEN** a real-source or real-LLM acceptance gate fails
- **THEN** the result SHALL classify the failure as network, rate limit, authentication or key, source schema change, no data, parse failure, model unavailable, quality failure, or another explicit category
- **AND** the result SHALL state whether the failure blocks release, blocks only a real-source or LLM claim, or is accepted as a documented degradation.

### Requirement: Keep release readiness separate from matrix definition

The project SHALL NOT claim release readiness merely because the acceptance gate matrix exists.

#### Scenario: Release candidate materials are prepared after P52

- **GIVEN** the acceptance gate matrix has been defined
- **WHEN** P53 release candidate materials are prepared
- **THEN** P53 SHALL reference actual gate results or explicitly documented waivers
- **AND** P52 alone SHALL NOT be treated as proof that all gates have passed.

### Requirement: Execute project acceptance gates before release readiness claims

The project SHALL execute the project acceptance gates before making release readiness claims.

#### Scenario: Release candidate materials are prepared

- **GIVEN** the P52 project acceptance gate matrix exists
- **WHEN** P53 release candidate materials are prepared
- **THEN** the materials SHALL include actual G0-G9 gate results or explicit waivers
- **AND** the materials SHALL NOT treat the existence of the P52 matrix as evidence that acceptance passed.

### Requirement: Record blocked, degraded, and skipped gates explicitly

The project SHALL preserve acceptance failures and degraded outcomes in release materials.

#### Scenario: A gate does not pass

- **WHEN** an acceptance gate is blocked, degraded, or skipped
- **THEN** the acceptance record SHALL include the command or workflow entry, artifact, failure or skip reason, classification when applicable, and release impact
- **AND** release candidate materials SHALL state `release_blocked` when any non-waived release-blocking gate is blocked.

### Requirement: Keep acceptance artifacts redacted

The project SHALL keep acceptance and release materials free of sensitive runtime details.

#### Scenario: Acceptance materials are written

- **WHEN** P53 writes acceptance or release candidate materials
- **THEN** the materials SHALL NOT include complete API keys, private paths, raw HTTP responses, complete prompts, raw SQL dumps, or unredacted vendor payloads
- **AND** real-source or real-LLM results SHALL be summarized with redacted artifacts and explicit failure categories.

### Requirement: Provide release handoff materials after acceptance execution

The project SHALL provide release handoff materials after a release-ready acceptance execution.

#### Scenario: A release candidate is handed off

- **GIVEN** P53 has produced a `release_ready` candidate with actual acceptance results
- **WHEN** P54 release handoff materials are written
- **THEN** the materials SHALL cite the P53 acceptance run and release candidate documents
- **AND** the materials SHALL preserve documented degradations, retries, and non-claims instead of broadening the release-ready statement.

### Requirement: Define acceptance repeatability rules

The project SHALL document how release acceptance can be repeated consistently.

#### Scenario: Acceptance is repeated after P53

- **WHEN** an operator reruns acceptance after P53
- **THEN** the repeatability guidance SHALL define the output directory convention, command order, allowed retry conditions, current-data degraded handling, real-source configuration prerequisites, real-LLM redaction rules, and release status decision rules
- **AND** a second failure after an allowed retry SHALL be treated as blocked unless an explicit waiver is recorded.

### Requirement: Keep release handoff within safety boundaries

Release handoff materials SHALL NOT expand the product safety boundary.

#### Scenario: Handoff materials are reviewed

- **WHEN** release handoff materials describe release readiness
- **THEN** they SHALL NOT claim investment returns, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, future provider availability, login sources, paid sources, authorized sources, Level2 data, or high-frequency data
- **AND** they SHALL NOT include complete API keys, private paths, raw HTTP responses, complete prompts, raw SQL dumps, or unredacted vendor payloads.

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
