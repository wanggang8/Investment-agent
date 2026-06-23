## ADDED Requirements

### Requirement: P75 Requirements Traceability And Real Use Closure

P75 SHALL prevent scoped acceptance evidence from being presented as full original-requirement completion and SHALL require explicit traceability from `docs/requirements.md` to real product evidence before any full-product release claim.

#### Scenario: Original requirements are traced before full-product claims

- **GIVEN** `docs/requirements.md` is the L1 product requirement source
- **WHEN** P75 evaluates release readiness
- **THEN** each original requirement paragraph, bullet, table row, SOP step, acceptance criterion, and safety/compliance statement SHALL be assigned a stable requirement ID, source line range, and requirement text hash
- **AND** each atomic requirement SHALL be mapped to implementation evidence, UI evidence where applicable, data evidence, workflow/rule/LLM evidence where applicable, scenario evidence, data-impact evidence where applicable, and safety-boundary evidence
- **AND** each atomic requirement SHALL be classified as `real_pass`, `scoped_pass`, `deterministic_local_evidence`, `partial`, `not_implemented`, or `blocked`
- **AND** each atomic requirement SHALL include criticality, criticality reason, full-release requirement flag, allowed release claim, evidence freshness, verification command, acceptance artifact, and delivered-by change where available
- **AND** release materials SHALL NOT claim full product completion while any `full_release_required=true` requirement remains `scoped_pass`, `deterministic_local_evidence`, `partial`, `not_implemented`, or `blocked`.

#### Scenario: User-raised real-use concerns are first-class gates

- **GIVEN** the user asks whether the product is truly usable rather than a demo
- **WHEN** P75 defines its acceptance scope
- **THEN** it SHALL include gates for dynamic user-entered fund/ETF symbols, external data lookup based on that symbol, built-in master wisdom usage by workflow/LLM context, external and built-in data completeness, analysis accuracy, UI task-flow design, function-to-data impact, cross-page readback, auditability, and release-claim honesty
- **AND** it SHALL NOT treat route smoke tests, screenshots, fixture-only tests, or a single accepted symbol as sufficient full-product evidence by themselves.

#### Scenario: Dynamic symbol support is not fabricated

- **GIVEN** a user enters a fund or ETF symbol
- **WHEN** the product evaluates readiness, consultation, alerts, or expected-return behavior
- **THEN** it SHALL resolve the symbol profile, tracked index, fund-side data, index-side data, formal evidence, market price, valuation, liquidity, and safe-degradation status from configured or collected facts
- **AND** at least one non-`510300` fund or ETF scenario SHALL trigger read-only market/evidence collection or an accepted-local request-construction equivalent, and SHALL prove collector or bridge request parameters, stored facts, source health, freshness, audit events, and readiness are bound to the user-entered symbol and its tracked index
- **AND** preseeded local facts or readiness rows without request-construction evidence SHALL NOT prove dynamic external querying
- **AND** unknown or unsupported symbols SHALL return blocked or information-insufficient states
- **AND** the product SHALL NOT silently substitute `510300`, `000300`, fixture data, stale data, built-in commentary, or C-level background material to make the flow appear ready.

#### Scenario: Missing data propagates to dependent claims

- **GIVEN** original requirements depend on market price, valuation, liquidity, funds flow, margin financing, constituent financials, media heat, sentiment proxy, formal evidence, RAG/index health, fund profile, and tracked index data
- **WHEN** any category is missing, stale, background-only, or source-unavailable
- **THEN** P75 SHALL record which claims, UI states, expected-return outputs, alerts, SOP steps, and release statements must degrade or block
- **AND** the product SHALL NOT declare normal emotion state, normal financing state, intact fundamentals, neutral funds flow, reliable safety margin, reliable valuation, reliable expected return, or trade-like next action when the data required for that claim is missing.

#### Scenario: Fund-side and index-side facts are joined safely

- **GIVEN** a fund or ETF uses both fund-side data and tracked-index-side data
- **WHEN** P75 verifies readiness, consultation, alerts, or expected-return behavior
- **THEN** it SHALL record join keys, source category, as-of date, freshness, and conflict handling for fund profile, NAV or price, liquidity, tracked index, benchmark symbol, index valuation, constituent or financial data, and formal evidence
- **AND** stale, mismatched, or ambiguous joins SHALL degrade or block affected claims.

#### Scenario: Built-in knowledge can guide analysis but cannot replace evidence

- **GIVEN** the product has built-in master wisdom, discipline rules, SOPs, and symbol profiles
- **WHEN** workflow and LLM analyst requests are constructed
- **THEN** they SHALL use the structured readiness/knowledge context or an explicitly equivalent source
- **AND** they SHALL expose which built-in knowledge was used
- **AND** they SHALL NOT allow built-in knowledge, local notes, prompts, or LLM output to satisfy formal evidence, source-verification, current-data, valuation, liquidity, or expected-return data requirements.

#### Scenario: LLM quality failures are retried safely

- **GIVEN** a real LLM analyst response fails the local quality gate because it contains unsafe trade-like instructions, final-verdict wording, deterministic prediction, or return-promise wording
- **WHEN** the failure category is `quality_failed`
- **THEN** the LLM client MAY perform one stricter safety reprompt that asks only for analysis material, evidence gaps, risks, and manual review questions
- **AND** the retry prompt SHALL explicitly forbid buy/sell instructions, final verdicts, deterministic predictions, and return promises
- **AND** repeated unsafe output SHALL remain `ANALYST_UNAVAILABLE` and degrade the affected analyst node rather than bypassing the quality gate
- **AND** network, HTTP, timeout, parse, empty-response, missing-key, model-unavailable, and provider-unavailable failures SHALL NOT be retried into a false pass.

#### Scenario: Analysis accuracy is checked against deterministic data

- **GIVEN** acceptance creates or uses known local facts
- **WHEN** P75 verifies risk alerts, expected return, valuation zones, liquidity rules, source verification, manual confirmations, portfolio snapshots, and derived page readbacks
- **THEN** it SHALL compare product outputs to deterministic expected values
- **AND** deterministic test vectors SHALL cover every executable threshold in the original requirements, including liquidity 20x and 5% thresholds, emotion 90%/10% and 3-day abnormality thresholds, two independent A/S source verification, PE/PB valuation zones, expected-return `<5` and `<20` sample gates, and cooldown/state-machine boundaries
- **AND** it SHALL verify which SQLite tables changed and which did not change after every user action
- **AND** every critical user action SHALL have matching audit evidence.

#### Scenario: SOP scenarios are mapped to real UI/product behavior

- **GIVEN** the original requirements define SOP A-F
- **WHEN** P75 evaluates scenario coverage
- **THEN** each SOP SHALL have real UI/data-impact evidence or a non-pass status with release impact whenever the SOP has user-visible behavior
- **AND** API evidence MAY only supplement rule priority, prerequisite, and database assertions
- **AND** pass evidence SHALL include rule priority, data prerequisites, LLM role, user confirmation behavior, safe degradation, and readback surfaces.

#### Scenario: UI actions are traced through state and readback

- **GIVEN** a real user performs onboarding, fund addition, data-readiness review, consultation, decision-detail review, alert review, offline-action confirmation, error marking, rule proposal review, gatekeeper review, monthly review, or quarterly review
- **WHEN** P75 claims that flow is ready
- **THEN** the acceptance evidence SHALL include real browser operation, DOM/readback assertions, expected SQLite table changes, prohibited SQLite table changes, audit events, related page readbacks, mobile layout checks, and failure-state checks
- **AND** at least one continuous non-`510300` browser journey SHALL cover add fund, readiness, consultation or alerts, SQLite verification, and derived page readback for the same user-entered symbol and tracked index.

#### Scenario: UI design is part of release readiness

- **GIVEN** a feature is technically implemented
- **WHEN** P75 evaluates real user readiness
- **THEN** the relevant UI flow SHALL be reviewed for discoverability, clear next action, correct state language, mobile/desktop layout, error recovery, and absence of misleading trading affordances
- **AND** UI copy SHALL clearly distinguish system analysis, user offline execution, in-system confirmation, and account-state mutation
- **AND** UI design issues that can cause user misunderstanding SHALL be marked as release-impacting findings unless fixed.

#### Scenario: P75 release conclusion inherits prior release gates

- **GIVEN** P75 may produce a new release conclusion
- **WHEN** P75 reports `release_ready_full_requirements_traceable` or any scoped release-ready conclusion
- **THEN** the evidence SHALL cite and satisfy P52 G0-G9 gates, P66 strict current-data policy, P67 resolution state, P71-P74 repeatability rules, and P52 G6/G7 external-source/LLM failure classifications
- **AND** any skipped, degraded, waived, scope-excluded, source-unavailable, model-unavailable, or redaction-related result SHALL downgrade or block the affected claim.

#### Scenario: Final P75 conclusion is bounded

- **GIVEN** P75 acceptance completes
- **WHEN** release materials report the outcome
- **THEN** the conclusion SHALL be one of `release_ready_full_requirements_traceable`, `release_ready_scoped_with_traceability_gaps`, `release_pending_safety_review_scoped_with_traceability_gaps`, or `release_blocked_requirements_traceability`
- **AND** the conclusion SHALL enumerate every remaining scoped, deterministic-local-only, partial, not implemented, or blocked atomic requirement
- **AND** it SHALL NOT claim future investment returns, future market direction, future public-source or model-provider availability, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, automatic database overwrite, paid/login/authorization-gated sources, Level2 data, high-frequency data, or physical second-machine completion unless separately evidenced.
