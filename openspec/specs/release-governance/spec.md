# release-governance Specification

## Purpose
定义发布状态治理规则：在 P66 当前数据 strict gate 与 P67 本地处置记录之后，发布材料必须明确 release-ready 的适用范围、当前数据 clean claim 边界、候选包证据新鲜度和后续 clean-tree package refresh 要求。
## Requirements
### Requirement: Initial release version marker

The project SHALL preserve `v0.1.0` as the historical initial local release version marker recorded by P99, while allowing the current repository release marker to advance in later release changes.

#### Scenario: Initial version history is preserved

- **GIVEN** release history or P99 acceptance materials are inspected
- **WHEN** the initial local release marker is described
- **THEN** the materials SHALL identify `v0.1.0` as the initial local release version recorded by P99
- **AND** the materials SHALL NOT treat that historical marker as proof that the current root `VERSION` file must remain `v0.1.0`.

### Requirement: Post-P67 release readiness decision

After P67, release materials SHALL include a P68 release readiness decision before any new final release handoff claim is made.

#### Scenario: P66 remains blocked and P67 scope exclusion is active

- **GIVEN** the P66 strict current-data gate reports `policy=blocked` and `gate=block`
- **AND** the P67 resolution check reports `claim_state=resolved_with_scope_exclusion`
- **WHEN** release readiness is described
- **THEN** the release status SHALL explicitly exclude current local data health from clean-data claims
- **AND** the materials SHALL NOT describe current local data as clean or healthy
- **AND** the materials SHALL NOT describe the P66 policy as passed

#### Scenario: Final distribution package evidence is stale after later commits

- **GIVEN** release package repeat evidence was generated before later P66, P67, or P68 commits
- **WHEN** final distribution readiness is described
- **THEN** the materials SHALL either require a package refresh stage or explicitly limit the package evidence to the earlier candidate archive
- **AND** the materials SHALL NOT imply that the earlier package artifact includes later commits.

#### Scenario: No active P67 resolution exists

- **GIVEN** the P66 strict current-data gate reports `gate=block` or `gate=waiver_required`
- **AND** the P67 resolution check reports `claim_state=requires_resolution`
- **WHEN** release readiness is described
- **THEN** the result SHALL be release-blocking for any claim depending on current local data health
- **AND** the materials SHALL direct the operator to record a valid waiver or scope exclusion before making a limited release-ready claim.

### Requirement: Final release decision and risk closure

After P69 clean-tree package refresh, release materials SHALL include a final release decision and risk-closure record before the project is described as having no mandatory next phase.

#### Scenario: Limited local release evidence is sufficient

- **GIVEN** P63 full UI regression and P65 repeat acceptance evidence remain passing
- **AND** P67 reports an active scope exclusion for the blocked current-data gate
- **AND** P69 clean-tree package verify and repeat acceptance have passed
- **WHEN** the release handoff describes final milestone status
- **THEN** the status SHALL be `release_ready_limited_current_data_scope`
- **AND** the handoff SHALL state that no mandatory next phase remains for that limited scope
- **AND** optional future stages SHALL be separated from release blockers.

#### Scenario: Current data remains blocked

- **GIVEN** the P66 strict current-data gate reports `policy=blocked` and `gate=block`
- **WHEN** the final release decision is written
- **THEN** it SHALL NOT claim current local data is clean or healthy
- **AND** it SHALL NOT describe P67 `resolved_with_scope_exclusion` as a P66 policy pass
- **AND** it SHALL preserve the current-data limitation in the release status.

#### Scenario: Package evidence does not include later documentation

- **GIVEN** the P69 package source commit predates P69 and P70 documentation
- **WHEN** package evidence is described in final handoff material
- **THEN** the material SHALL state the exact covered source commit or phase boundary
- **AND** it SHALL NOT imply that the P69 package archive includes P69 or P70 documents.

### Requirement: Full real product acceptance true pass

After P70, if the project claims full real product acceptance instead of the limited current-data release scope, it SHALL execute a strict P71 acceptance run that treats current-data gate failure, VecLite retrieval degradation, real UI failure, real LLM failure, safety failure, and package verification failure as blockers.

#### Scenario: Current data must pass without scope exclusion

- **GIVEN** P71 is evaluating full real product acceptance
- **WHEN** the current data-source quality strict gate is executed for `000300`
- **THEN** the command MUST return `policy=passed` and `gate=pass`
- **AND** P67 `resolved_with_scope_exclusion`, fixture-only regression, or waiver documentation MUST NOT be accepted as a P71 current-data pass.

#### Scenario: VecLite degradation blocks full acceptance

- **GIVEN** P71 is running real UI consultation or retrieval acceptance
- **WHEN** VecLite/RAG index health is missing, corrupted, incompatible, stale, empty when required, or the workflow reports `VECTOR_INDEX_UNAVAILABLE`
- **THEN** P71 SHALL be marked blocked for retrieval index readiness
- **AND** the acceptance record SHALL NOT describe retrieval-enhanced context as passed.

#### Scenario: Full UI acceptance uses real local operation

- **GIVEN** P71 runs the frontend acceptance suite
- **WHEN** primary product routes and key actions are checked
- **THEN** the browser MUST operate against a real local Go backend and Vite frontend
- **AND** frontend mocks, mocked network responses, or fixture-only current data MUST NOT be used as pass evidence for real product acceptance
- **AND** unexpected API failures, console errors, page errors, mobile/desktop overflow, missing generated decision detail, missing LLM material, or forbidden trading/automation affordances MUST block the P71 pass.

#### Scenario: Post-P70 package refresh follows strict acceptance

- **GIVEN** strict P71 current-data, VecLite, real UI, real LLM, safety, and redaction gates pass
- **WHEN** final distribution package evidence is generated
- **THEN** the package SHALL be built from the accepted post-P70 commit
- **AND** package verify and repeat acceptance SHALL pass
- **AND** the package manifest SHALL state that P69, P70, and P71 acceptance materials are included only if they are present in the packaged source commit.

### Requirement: Real user fund scenario data-impact acceptance

After P71, if the project claims the product is ready for practical real-user use, it SHALL execute a P72 real user fund/ETF scenario acceptance that verifies UI operations, API responses, data side effects, auditability, derived readbacks, deterministic accuracy, and safety boundaries.

#### Scenario: Reviewed scenario matrix precedes execution

- **GIVEN** P72 is evaluating practical real-user readiness
- **WHEN** acceptance execution starts
- **THEN** a scenario matrix SHALL exist before execution
- **AND** the matrix SHALL cover real fund setup, holding maintenance, offline transaction recording, local knowledge/RAG, current data, daily discipline, risk alerts, real consultation, manual confirmation, review/readback, failure handling, and safety boundaries
- **AND** execution SHALL NOT proceed if the matrix is known to omit a primary user workflow.

#### Scenario: Data impact is verified after real UI operation

- **GIVEN** a P72 browser scenario performs a user operation
- **WHEN** the operation reports success in the UI
- **THEN** P72 SHALL verify the expected local data impact through API responses and read-only SQLite evidence
- **AND** the acceptance record SHALL include sanitized evidence for the relevant tables and audit events
- **AND** page refresh or navigation SHALL show the same resulting state where the product exposes it.

#### Scenario: Deterministic accuracy is separated from investment prediction

- **GIVEN** P72 evaluates a real fund/ETF scenario
- **WHEN** the product calculates market value, unrealized profit ratio, cash/asset ratios, risk trigger state, report links, or counts
- **THEN** those deterministic values SHALL be checked against independently computed expectations
- **AND** LLM or expected-return material SHALL only be checked for traceability, parse/quality status, rule-consistent final verdict, and risk disclosure
- **AND** P72 SHALL NOT claim future return or market-direction accuracy.

#### Scenario: Safety boundaries remain active in real use

- **GIVEN** P72 runs realistic user scenarios
- **WHEN** pages, API results, browser summaries, and SQLite-derived artifacts are scanned
- **THEN** they SHALL NOT expose complete keys, full prompts, raw provider payloads, private local paths, SQL dumps, broker/order capabilities, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, automatic overwrite, or return promises.

#### Scenario: P72 gaps are recorded honestly

- **GIVEN** a scenario cannot be executed because of provider availability, product behavior, or missing coverage
- **WHEN** P72 writes release or acceptance materials
- **THEN** the result SHALL be blocked or gap-qualified
- **AND** the materials SHALL NOT describe the missing scenario as passed.

### Requirement: P73 Product Effectiveness And UX Validation

P73 SHALL validate whether the product supports its stated investment-discipline assistant goal, not merely whether features execute.

#### Scenario: Product goal metrics are recorded

- **GIVEN** P73 evaluates product effectiveness
- **WHEN** acceptance materials are written
- **THEN** they SHALL record discipline adherence, evidence sufficiency, traceability, review usefulness, and UX comprehension checks
- **AND** they SHALL NOT use future investment returns as the required pass metric.

#### Scenario: Real UX task validation is required

- **GIVEN** a user operates the local product
- **WHEN** P73 runs browser acceptance
- **THEN** it SHALL cover first-use or missing-prerequisite guidance, daily discipline, portfolio maintenance, data quality/evidence review, consultation, decision detail, manual confirmation, risk/notification/audit/rules/review readback, and invalid or unsafe input
- **AND** page errors, unexpected API failures, console errors, forbidden affordances, or critical UX confusion SHALL block pass.

#### Scenario: Effect replay validates discipline behavior

- **GIVEN** local facts, evidence records, decisions, confirmations, risk alerts, rule-effect facts, and audit events exist
- **WHEN** P73 runs effect replay checks
- **THEN** C-level background-only material SHALL NOT satisfy formal evidence
- **AND** insufficient evidence SHALL result in safe blocking, gap qualification, or non-trade records
- **AND** manual confirmation SHALL be the only accepted path that mutates local portfolio facts
- **AND** rule proposal/effect validation SHALL expose sample, overfit, gate, or tracking state when available.

#### Scenario: UX audit findings are dispositioned

- **GIVEN** P73 captures representative UI screenshots and task results
- **WHEN** the UX audit is written
- **THEN** findings SHALL be classified as critical, major, minor, or accepted gap
- **AND** critical findings SHALL block product-effectiveness pass until fixed and rerun.

#### Scenario: P73 claims remain bounded

- **GIVEN** P73 passes
- **WHEN** release materials state the result
- **THEN** they MAY claim product-effectiveness and UX validation for the accepted local scope
- **AND** they SHALL NOT claim investment return improvement, future market prediction, future public-source or model-provider availability, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, or automatic real DB overwrite.

### Requirement: P74 Built-In Knowledge And Data Readiness

P74 SHALL make built-in investment knowledge and required data readiness explicit, auditable, and visible to runtime users.

#### Scenario: Built-in knowledge is structured and bounded

- **GIVEN** the product includes master principles, discipline rules, risk SOPs, and symbol profile knowledge
- **WHEN** P74 exposes built-in knowledge
- **THEN** each entry SHALL have a stable ID, category, summary, applicability, rule mapping, LLM context eligibility, and safety boundary
- **AND** master principles, local notes, and background knowledge SHALL NOT be classified as formal market evidence.

#### Scenario: Data readiness maps required categories to feature impact

- **GIVEN** a user requests readiness for a symbol
- **WHEN** the readiness service evaluates local facts
- **THEN** it SHALL report required and optional data categories as ready, degraded, missing, background-only, or blocked
- **AND** it SHALL map missing or degraded categories to affected product surfaces and claims
- **AND** it SHALL NOT trigger collectors, rebuild indexes, modify source health, write market snapshots, update rules, create notifications, mutate portfolios, or create confirmations.

#### Scenario: Readiness API is sanitized and safe

- **GIVEN** the frontend requests `GET /api/v1/knowledge-readiness?symbol=510300`
- **WHEN** the backend returns readiness data
- **THEN** the response SHALL include overall status, symbol profile, knowledge references, data dependencies, feature impacts, LLM context summary, and safety notes
- **AND** it SHALL NOT expose full prompts, raw HTTP responses, raw LLM responses, private local paths, API keys, private keys, original SQL, or complete account details.

#### Scenario: LLM analysis receives readiness context without decision authority

- **GIVEN** a workflow invokes an LLM analyst node
- **WHEN** relevant knowledge/data readiness context exists
- **THEN** the analyst request SHALL include a sanitized summary of matched principles and data readiness
- **AND** the prompt SHALL state that background knowledge cannot satisfy formal evidence
- **AND** the LLM SHALL remain limited to analysis material and SHALL NOT generate or override the final rule verdict.

#### Scenario: UI shows readiness and gaps

- **GIVEN** readiness data is available
- **WHEN** the user views data quality, rules, consultation, or decision detail surfaces
- **THEN** the UI SHALL distinguish knowledge available as rules, knowledge available as LLM context, background-only knowledge, ready data, degraded data, missing data, and blocked claims
- **AND** it SHALL show safe next steps without broker actions, automatic trading, one-click trading, external push, automatic confirmation, or return promises.

#### Scenario: P74 acceptance covers complete and degraded data scenarios

- **GIVEN** P74 acceptance runs against local test databases and the real UI/API surfaces
- **WHEN** acceptance completes
- **THEN** it SHALL cover a complete `510300` ETF/index path and degraded paths for missing valuation data, background-only local knowledge, single-source evidence, multi-source formal evidence, and out-of-scope capability
- **AND** it SHALL block pass if readiness is only documented but not available through API/UI/LLM-context evidence.

#### Scenario: P74 claims remain bounded

- **GIVEN** P74 passes
- **WHEN** release materials state the result
- **THEN** they MAY claim built-in knowledge and data readiness traceability for the accepted local scope
- **AND** they SHALL NOT claim future investment returns, future market direction, future public-source or model-provider availability, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, automatic database overwrite, paid/login/authorization-gated sources, Level2 data, or high-frequency data.

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

### Requirement: P76 SHALL refresh package evidence after P75

P76 SHALL regenerate final local package evidence after P75 so package freshness claims are not based on the earlier P71 archive.

#### Scenario: Package source is the clean post-P75 package commit

- **GIVEN** P75 has been committed and archived
- **WHEN** P76 generates package evidence
- **THEN** the package SHALL be generated from a clean source commit that includes committed P72-P75 evidence and any P76 acceptance-harness correction required to make repeat acceptance deterministic
- **AND** the package manifest SHALL record `source_status=clean`
- **AND** the package manifest SHALL record the package source commit
- **AND** release materials SHALL state whether P72-P75 acceptance Markdown and OpenSpec archives are included in the packaged source.

#### Scenario: Package verify and repeat acceptance pass

- **WHEN** the P76 package archive is generated
- **THEN** package verification SHALL confirm archive checksum consistency, required entries, forbidden path exclusions, and manifest safety boundaries
- **AND** repeat acceptance SHALL run from an extracted package workspace rather than from the active repository checkout
- **AND** repeat acceptance SHALL cover OpenSpec validation, Go tests, frontend dependency installation, frontend tests, frontend build, and local E2E smoke.

#### Scenario: Package handoff remains bounded

- **WHEN** P76 updates release materials
- **THEN** the materials SHALL include package identity, source commit, source status, checksum, archive entry count, verify result, repeat result, known caveats, and Not Claimed boundaries
- **AND** the materials SHALL NOT claim that the archive includes P76 package-after-the-fact evidence or `docs/release/ui-audit-assets/`
- **AND** the materials SHALL preserve P75 `release_ready_scoped_with_traceability_gaps`
- **AND** the materials SHALL NOT claim physical second-machine execution, remote publishing, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic restore, automatic repair, real database overwrite, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, future provider availability, or investment returns.

### Requirement: P77 SHALL govern post-P75 real-pass upgrades

P77 SHALL create a conservative, auditable upgrade gate for moving P75 atomic requirement rows toward `real_pass` without rewriting historical P75 acceptance evidence or expanding unsupported release claims.

#### Scenario: P77 upgrade evidence is a new layer

- **GIVEN** P75 produced `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md`
- **WHEN** P77 evaluates row-level status upgrades
- **THEN** P77 SHALL generate a new matrix that preserves P75 row IDs, source line ranges, requirement text hashes, original statuses, full-release-required flags, and release impacts
- **AND** P77 SHALL NOT mutate the historical P75 matrix to make prior evidence appear stronger than it was
- **AND** each P77 row SHALL record `p77_status`, upgrade basis, gate dimensions, fresh evidence command, fresh evidence artifact, residual gap, and next remediation.

#### Scenario: Real-pass upgrade requires all applicable evidence dimensions

- **GIVEN** a P75 row is being considered for `real_pass`
- **WHEN** P77 evaluates that row
- **THEN** the row SHALL have implementation evidence
- **AND** user-visible behavior SHALL have real UI evidence
- **AND** mutating behavior SHALL have SQLite changed-table, prohibited-table, audit-event, and readback evidence
- **AND** data-source, collector, workflow, rule, LLM, RAG, and scenario-dependent behavior SHALL have direct evidence for each applicable dependency
- **AND** safety evidence SHALL confirm no broker interface, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, future provider-availability promise, or investment return promise is introduced
- **AND** the row SHALL remain non-`real_pass` if its only evidence is screenshot-only, route-smoke-only, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, temporary-DB-only, or incompatible single-symbol-only.

#### Scenario: P77 release conclusion remains bounded

- **WHEN** P77 reports its final conclusion
- **THEN** it SHALL claim `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`
- **AND** otherwise it SHALL use a scoped conclusion that lists the remaining non-`real_pass` rows or grouped categories with row-level matrix reference
- **AND** it SHALL preserve P76 package boundaries unless a separate package refresh change is executed.

### Requirement: P78 SHALL close real-pass gaps in conservative batches

P78 SHALL continue post-P77 full-requirement acceptance by classifying remaining non-`real_pass` rows into remediation batches and upgrading only rows that meet the P77 evidence gate with fresh, row-specific evidence.

#### Scenario: P78 classifies remaining gaps before upgrading rows

- **GIVEN** P77 produced `docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md`
- **WHEN** P78 evaluates the remaining full-release-required non-`real_pass` rows
- **THEN** P78 SHALL generate a new matrix that preserves P77 row IDs, source ranges, original status, P77 status, full-release-required flag, and release impact
- **AND** each non-`real_pass` row SHALL receive a remediation group, batch assignment, remaining gap, and next action
- **AND** P78 SHALL NOT mutate P75 or P77 historical matrices.

#### Scenario: P78 batch upgrades require direct evidence

- **GIVEN** a P78 batch proposes a row for `real_pass`
- **WHEN** the P78 checker evaluates the row
- **THEN** implementation behavior SHALL be backed by fresh deterministic tests or direct runtime evidence
- **AND** user-visible behavior SHALL have real browser UI evidence when applicable
- **AND** data-bearing behavior SHALL have SQLite readback evidence for the exact fields claimed
- **AND** expected-return or analysis rows SHALL show sample count, sample window, screening condition, source/provenance fields, precision/degradation status, and non-trading disclaimer when applicable
- **AND** the row SHALL remain non-`real_pass` if evidence is only inherited, screenshot-only, route-smoke-only, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, temporary-DB-only, or incompatible single-symbol-only.

#### Scenario: P78 release conclusion remains bounded

- **WHEN** P78 reports its conclusion
- **THEN** it SHALL claim `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`
- **AND** otherwise it SHALL use a scoped conclusion that records upgraded row count, remaining non-`real_pass` row count, remediation groups, and package freshness boundaries
- **AND** it SHALL NOT claim P76 package inclusion unless a separate package refresh change is executed.

### Requirement: P79 Real-Use Data-Impact Closure

After P78, any P79 claim that portfolio, confirmation, local-account, or expected-return rows have moved to `real_pass` SHALL be backed by fresh real UI execution and SQLite/readback evidence.

#### Scenario: P79 upgrades require action-to-data proof

- **GIVEN** a P79 row is upgraded to `real_pass`
- **WHEN** the P79 checker evaluates the row
- **THEN** the row SHALL have fresh P79 evidence from a real browser journey or direct runtime readback
- **AND** data-bearing rows SHALL include SQLite readback for expected changed tables
- **AND** local-account rows SHALL include field-level readback for the relevant position, confirmation, transaction, evidence, and audit fields rather than table counts alone
- **AND** data-bearing rows SHALL include negative evidence that prohibited broker/order/external-push/automatic-confirmation tables or claims were not created
- **AND** the row SHALL remain non-`real_pass` if the evidence is inherited-only, screenshot-only, route-smoke-only, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, or incompatible single-action-only.

#### Scenario: P79 expected-return rows remain bounded

- **GIVEN** a P79 expected-return row requires probabilities, scenario ranges, sell-evaluation triggers, valuation fields, sample counts, sample windows, screening conditions, source/provenance fields, or non-trading disclaimers
- **WHEN** fresh P79 evidence lacks any required field
- **THEN** that row SHALL remain non-`real_pass`
- **AND** P79 SHALL record the missing field-level evidence as the remaining gap.

#### Scenario: P79 release conclusion remains scoped

- **WHEN** P79 reports its conclusion
- **THEN** it SHALL claim `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`
- **AND** otherwise it SHALL use a scoped conclusion that records upgraded row count, remaining non-`real_pass` row count, and package freshness boundaries
- **AND** it SHALL NOT claim P76 package inclusion unless a separate package refresh change is executed.

#### Scenario: Expected-return quality failure uses safe local material

- **GIVEN** the expected-return LLM material is parseable but fails the analyst safety quality gate
- **WHEN** deterministic local expected-return scenarios have been generated
- **THEN** the failed LLM material SHALL be discarded
- **AND** ExpectedReturnNode SHALL emit safe deterministic local expected-return material with metadata showing `model=deterministic-local`, `parse_status=parsed`, `quality_status=passed`, and `fallback_reason=llm_quality_failure`
- **AND** ordinary analyst timeout, authentication, or model-unavailable errors SHALL continue to degrade the analyst node
- **AND** this fallback SHALL NOT be used to upgrade expected-return probability or scenario rows without separate field-level UI/readback evidence.

### Requirement: P80 Review Audit Governance Closure

After P79, any P80 claim that review, audit, error-case, rule-proposal, or gatekeeper-governance rows have moved to `real_pass` SHALL be backed by fresh real UI execution and field-level SQLite/readback evidence.

#### Scenario: P80 upgrades require review and audit field proof

- **GIVEN** a P80 row is upgraded to `real_pass`
- **WHEN** the P80 checker evaluates the row
- **THEN** the row SHALL have fresh browser evidence from the P80 review/audit/governance journey
- **AND** data-bearing rows SHALL include SQLite readback for the expected tables and fields
- **AND** audit rows SHALL include `action`, `node_action`, `actor`, `status`, `before_state`, `after_state`, and `request_id` when those fields are part of the row claim
- **AND** governance rows SHALL include rule proposal, gatekeeper audit, and audit-event references when those fields are part of the row claim
- **AND** the row SHALL remain non-`real_pass` if the evidence is count-only, screenshot-only, route-smoke-only, fixture-only without UI operation, mock-only, waiver-only, or only partially covers the row text.

#### Scenario: P80 broad monthly and final-application rows remain bounded

- **GIVEN** a row requires monthly attribution, full quarterly review, final rule application time, or every SOP/data-impact branch
- **WHEN** P80 does not prove the exact required fields through fresh UI and readback
- **THEN** that row SHALL remain non-`real_pass`
- **AND** P80 SHALL record the missing field-level evidence as the remaining gap.

#### Scenario: P80 release conclusion remains scoped

- **WHEN** P80 reports its conclusion
- **THEN** it SHALL claim `release_ready_full_requirements_traceable` only if every `full_release_required=true` row is `real_pass`
- **AND** otherwise it SHALL use a scoped conclusion that records upgraded row count, remaining non-`real_pass` row count, and package freshness boundaries
- **AND** it SHALL NOT claim P76 package inclusion unless a separate package refresh change is executed.

### Requirement: P81 dynamic source field coverage closure

After P80, dynamic source field coverage rows SHALL NOT be marked `real_pass` unless fresh acceptance proves the current product obtains, evaluates, displays, or safely blocks the relevant data fields for a user-selected symbol through real local product paths.

#### Scenario: P81 row inventory is complete before execution

- **GIVEN** P81 starts from the P80 evidence matrix
- **WHEN** execution begins
- **THEN** the P81 plan SHALL enumerate exactly 59 dynamic source field coverage rows
- **AND** the plan SHALL preserve the previous status and target evidence type for each row.

#### Scenario: User-selected symbol drives source evidence

- **GIVEN** P81 evaluates data source coverage
- **WHEN** a browser or API scenario requests readiness or analysis for a symbol
- **THEN** the evidence SHALL show that the requested symbol drives the source/readiness result
- **AND** hard-coded `510300`-only evidence SHALL NOT be sufficient for P81 `real_pass`.

#### Scenario: Formal evidence unavailable

- **GIVEN** an external or built-in data category is missing, degraded, stale, or background-only
- **WHEN** P81 evaluates impacted features
- **THEN** the product SHALL safely degrade, qualify, or block affected claims
- **AND** the acceptance result SHALL NOT mark formal evidence requirements as passed by background knowledge alone.

#### Scenario: P81 claims remain bounded

- **GIVEN** P81 passes some or all rows
- **WHEN** release materials are updated
- **THEN** they SHALL state the exact upgraded rows and remaining non-`real_pass` count
- **AND** they SHALL NOT claim full original-requirement pass, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real DB overwrite, return promises, paid/login/authorized source, Level2 source, or high-frequency source.

### Requirement: P82 SOP action UI-to-SQLite closure

After P80, SOP/action data-impact rows SHALL NOT be marked `real_pass` unless a real browser workflow proves the user operation, API result, durable local data impact, auditability, and UI readback.

#### Scenario: P82 row inventory is complete before execution

- **GIVEN** P82 starts from the P81 evidence matrix while preserving P80 classification provenance
- **WHEN** execution begins
- **THEN** the P82 plan SHALL enumerate exactly 53 SOP/action rows for evaluation
- **AND** each row SHALL map to a real UI scenario, readback target, safety negative check, and upgrade-or-defer decision.

#### Scenario: P82 upgrades only directly proven rows

- **GIVEN** a P82 planned row is broader than the fresh SOP/action evidence
- **WHEN** P82 generates the evidence layer
- **THEN** that row SHALL remain non-`real_pass`
- **AND** the acceptance record SHALL name the exact remaining gap and next-batch owner.

#### Scenario: UI operation creates expected local evidence

- **GIVEN** a P82 user action reports success in the UI
- **WHEN** P82 validates the result
- **THEN** it SHALL check API response state, read-only SQLite evidence, audit events, and visible readback after navigation or refresh.

#### Scenario: Unsupported automation remains blocked

- **GIVEN** a scenario concerns SOP, confirmations, notifications, or rule governance
- **WHEN** the UI is inspected
- **THEN** it SHALL NOT expose automatic trading, one-click trading, order delegation, external push, automatic confirmation, or automatic rule application as available product actions.

### Requirement: P83 governance traceability backfill

After P80, governance and release traceability rows SHALL NOT be upgraded unless each row has exact artifact links, fresh validation where needed, and an honest status classification.

#### Scenario: P83 row inventory is complete before execution

- **GIVEN** P83 starts from the latest P82 evidence matrix
- **WHEN** execution begins
- **THEN** the P83 plan SHALL enumerate exactly 43 governance traceability rows
- **AND** each row SHALL have a target evidence or classification path.

#### Scenario: Evidence links are concrete

- **GIVEN** P83 marks a row as upgraded
- **WHEN** the evidence layer is reviewed
- **THEN** it SHALL include exact files, commands, tests, package manifests, acceptance records, UI/API evidence, or safety scans
- **AND** narrative-only assertions SHALL NOT be sufficient.

#### Scenario: Historical gaps remain honest

- **GIVEN** a historical archive or physical repeat acceptance was never performed
- **WHEN** P83 writes governance materials
- **THEN** it SHALL preserve that limitation
- **AND** it SHALL NOT fabricate historical archives, physical second-machine evidence, package refreshes, remote release, or Git tag evidence.

### Requirement: P84 portfolio confirmation data-impact closure

After P83, portfolio and manual-confirmation rows SHALL NOT be marked `real_pass` unless real local UI workflows prove user action, local data mutation, downstream readback, deterministic value accuracy where applicable, and safety boundaries.

#### Scenario: P84 row inventory is complete before execution

- **GIVEN** P84 starts from the P83 evidence matrix
- **WHEN** execution begins
- **THEN** the P84 plan SHALL enumerate exactly 35 portfolio/confirmation rows
- **AND** each row SHALL map to before/after data impact and downstream readback evidence.

#### Scenario: Portfolio mutation is local and manual

- **GIVEN** a P84 scenario changes portfolio-related state
- **WHEN** the change is accepted
- **THEN** the evidence SHALL show local user-driven action, API response, SQLite before/after delta, audit event, and UI readback
- **AND** it SHALL NOT depend on broker synchronization, automatic trading, order placement, or automatic confirmation.

#### Scenario: Derived values are checked independently

- **GIVEN** P84 evidence includes market value, cost, quantity, ratio, cash, or profit/loss values
- **WHEN** the row is evaluated
- **THEN** P84 SHALL compare product output with independently computed expectations
- **AND** future return or market-direction accuracy SHALL NOT be claimed.

### Requirement: P85 expected return analysis accuracy closure

After P84, expected-return and analysis-accuracy rows SHALL NOT be marked `real_pass` unless fresh real local execution proves deterministic calculation correctness, provenance, degradation safety, LLM boundary safety, and user-visible readback.

#### Scenario: P85 row inventory is complete before execution

- **GIVEN** P85 starts from the P84 evidence matrix
- **WHEN** execution begins
- **THEN** the P85 plan SHALL enumerate exactly 31 expected-return and analysis-accuracy rows
- **AND** each row SHALL map to a concrete acceptance mode and evidence target.

#### Scenario: Expected-return calculations are deterministic checks

- **GIVEN** P85 evaluates expected-return or scenario fields
- **WHEN** a row is marked `real_pass`
- **THEN** the product output SHALL be compared with independently computed deterministic expectations where the value is deterministic
- **AND** future return, future market direction, or investment performance accuracy SHALL NOT be claimed.

#### Scenario: LLM remains analysis-only

- **GIVEN** a P85 scenario includes real LLM output, unavailable LLM, or LLM quality failure
- **WHEN** the decision workflow completes or degrades
- **THEN** LLM material SHALL remain analysis-only
- **AND** it SHALL NOT override final rule verdict, create confirmations, trigger trades, or suppress required data-quality blockers.

### Requirement: P87 portfolio state allocation safety closure

After P84, the portfolio/allocation/state/data-impact rows not covered by P85 or P86 SHALL NOT be marked `real_pass` unless fresh real local execution proves the complete row-specific behavior through UI operation, API responses, read-only SQLite evidence, deterministic checks where applicable, and explicit forbidden-capability absence.

#### Scenario: P87 row inventory closes the planning gap

- **GIVEN** P87 starts from the P84 evidence matrix
- **WHEN** the P85, P87, and P86 plans are reviewed together
- **THEN** they SHALL cover exactly 157 P84-after full-release-required non-`real_pass` rows
- **AND** no row SHALL be omitted or owned by two execution batches.

#### Scenario: Portfolio state and allocation evidence is row-specific

- **GIVEN** P87 evaluates account, holding-state, allocation, rebalance, or confirmation requirements
- **WHEN** a row is marked `real_pass`
- **THEN** the acceptance evidence SHALL include real browser UI operation, API/readback, SQLite field checks, and deterministic calculations where the value is deterministic
- **AND** broad rows such as monthly attribution or audit history SHALL only pass if the full stated breadth is proven.

#### Scenario: Data-insufficient and release safety remain hard boundaries

- **GIVEN** P87 evaluates degraded data, insufficient evidence, release checks, or safety boundaries
- **WHEN** evidence is unavailable or degraded
- **THEN** the product SHALL visibly qualify or block the affected advice
- **AND** it SHALL NOT create confirmations, trigger trades, suppress blockers, or imply automatic install/upgrade/migration/repair behavior.

### Requirement: P86 core goal knowledge safety final closure

After P81-P85 and P87, P86 SHALL reconcile the remaining core-goal, source/data transition, knowledge/LLM/RAG, expected-return, review/audit, implementation, release-safety, and unclassified rows into a final row-level matrix and SHALL only claim full original-requirement pass if every full-release-required row is resolved by valid evidence or explicitly reclassified with documented rationale.

#### Scenario: P86 row inventory completes the post-P87 remainder

- **GIVEN** P86 starts from the P87 evidence matrix
- **WHEN** the P86 plan is reviewed
- **THEN** P86 SHALL cover exactly the 137 remaining full-release-required non-`real_pass` rows from P87
- **AND** no P87 remaining row SHALL be omitted from the final P86 inventory.

#### Scenario: Integrated real user acceptance

- **GIVEN** P86 evaluates the product goal
- **WHEN** end-to-end acceptance runs
- **THEN** it SHALL use real local UI operation, API responses, workflow metadata, read-only SQLite evidence, and deterministic checks where applicable
- **AND** it SHALL cover setup, portfolio/account state, data readiness, knowledge/RAG, consultation, expected return, risk/SOP, manual confirmation, review, audit, release governance, and safety.

#### Scenario: Row upgrade requires direct evidence

- **GIVEN** P86 generates the final matrix
- **WHEN** a row is upgraded to `real_pass`
- **THEN** the matrix SHALL cite direct row-level evidence from P86 or cumulative P81-P87 artifacts
- **AND** it SHALL NOT rely only on seeded SQLite rows, route smoke, screenshots, fixture/mock/stub data, or broad narrative.

#### Scenario: Full-pass claim is evidence gated

- **GIVEN** P86 writes final release materials
- **WHEN** any full-release-required row remains partial, blocked, scoped-only, unsupported, or unverified
- **THEN** P86 SHALL NOT claim full original-requirement pass
- **AND** it SHALL list the exact remaining rows and blockers.

#### Scenario: Forbidden capabilities remain out of product scope

- **GIVEN** P86 passes integrated acceptance
- **WHEN** final claims are written
- **THEN** they SHALL NOT introduce or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real DB overwrite, paid/login/authorized source, Level2 source, high-frequency source, future provider availability, or return promises.

### Requirement: P88 remaining full release blockers closure

After P86, P88 SHALL resolve or explicitly preserve the 27 remaining full-release-required rows by adding row-specific implementation and real UI/API/SQLite/workflow evidence for source-verified transitions, structured public-data fields, expected-return historical/probability behavior, quarterly rebalance, and SOP addendum proposals.

#### Scenario: P88 row inventory starts from the P86 remainder

- **GIVEN** P88 starts from the P86 matrix
- **WHEN** the inventory gate runs
- **THEN** it SHALL find exactly 27 full-release-required non-`real_pass` rows
- **AND** the row IDs SHALL be `REQ-02-022`, `REQ-02-025`, `REQ-04-016`, `REQ-04-025`, `REQ-05-003`, `REQ-05-004`, `REQ-05-005`, `REQ-06-023`, `REQ-06-024`, `REQ-08-004`, `REQ-08-023`, `REQ-09-001`, `REQ-09-003`, `REQ-09-004`, `REQ-09-006`, `REQ-09-007`, `REQ-09-008`, `REQ-09-009`, `REQ-09-010`, `REQ-09-013`, `REQ-09-023`, `REQ-09-024`, `REQ-09-025`, `REQ-09-027`, `REQ-10-004`, `REQ-13-010`, and `REQ-17-015`.

#### Scenario: Source-verified state transitions are proven by evidence counts

- **GIVEN** formal source-verification evidence exists for a held symbol
- **WHEN** at least two independent A/S formal sources confirm buy-logic break
- **THEN** P88 SHALL prove the workflow enters `sell_only`, prohibits buy/add actions, and records source-count provenance, API/readback, SQLite facts, and audit evidence
- **AND** when fewer than two independent A/S formal sources exist for buy-logic questioned, major positive, or major negative information, P88 SHALL prove the workflow enters `frozen_watch` with source-count provenance and pause guidance.

#### Scenario: Structured data fields require preverified public-source evidence

- **GIVEN** P88 expands structured data evidence for capital flow, margin financing, and constituent financials
- **WHEN** collectors are used for `real_pass` structured-data evidence
- **THEN** P88 SHALL record source preverification before claiming production readiness
- **AND** it SHALL prove field-level readback for capital-flow date/net-inflow/net-outflow, margin-financing date/balance/change-rate, and constituent-financial revenue/profit/growth/disclosure-date.
- **AND** accepted-local, fixture, stub, or manually seeded evidence SHALL NOT upgrade structured-data collector rows to `real_pass`.

#### Scenario: Expected-return report uses historical/probability evidence and safe degradation

- **GIVEN** expected-return analysis runs through real UI/API/SQLite acceptance
- **WHEN** sufficient historical similar samples exist
- **THEN** P88 SHALL prove probabilities are derived from sample proportions, the base scenario is the highest-frequency path, pessimistic scenario is displayed, and the report shows target name/code, future-12-month ranges, sample metadata, triggers, and disclaimer
- **AND** it SHALL prove a representative holding-class coverage matrix covering broad ETF/index fund, sector/growth ETF or fund, and equity/security-like constituent-financial path before upgrading `REQ-09-001`
- **AND** when samples are fewer than five, P88 SHALL prove no return range is generated and a supplement-data list is displayed.

#### Scenario: Expected-return dynamic monitoring is proven

- **GIVEN** valuation, fundamentals, market state, assumptions, or actual path data change
- **WHEN** the expected-return monitoring path runs
- **THEN** P88 SHALL prove affected scenario probabilities are lowered when applicable
- **AND** it SHALL prove periodic assumption checks, two-month below-expectation downshift warning, and one-month pessimistic-path manual probability-adjustment suggestion.

#### Scenario: Quarterly rebalance remains manual and auditable

- **GIVEN** a portfolio drifts beyond quarterly +/-15% target bands
- **WHEN** the rebalance flow runs
- **THEN** P88 SHALL prove manual buy/sell recommendation amounts through UI/API/SQLite/audit readback
- **AND** it SHALL NOT create broker orders, trades, automatic confirmations, or external push events.

#### Scenario: SOP addendum proposal is generated without automatic rule application

- **GIVEN** repeated review/error-case evidence identifies a high-frequency uncovered scenario
- **WHEN** P88 runs the SOP addendum path
- **THEN** it SHALL create a pending `sop` rule proposal, notification, and audit event
- **AND** it SHALL NOT modify active rules unless the existing gatekeeper and final user-confirmation flow is explicitly completed.

#### Scenario: P88 final claims remain evidence gated

- **GIVEN** P88 writes final release materials
- **WHEN** any full-release-required row remains partial, blocked, scoped-only, unsupported, or unverified
- **THEN** P88 SHALL NOT claim full original-requirement pass
- **AND** it SHALL list exact remaining rows and blockers.
- **AND** any row reclassification SHALL require explicit L1/OpenSpec rationale and SHALL NOT be reported as equivalent to 27/27 `real_pass`.

### Requirement: P90 capital-flow provider closure

After P89, P90 SHALL resolve or explicitly preserve the two remaining capital-flow related full-release-required rows using real public provider evidence and product UI/API/SQLite readback.

#### Scenario: P90 row inventory starts from the P89 remainder

- **GIVEN** P90 starts from the P89 matrix
- **WHEN** the inventory gate runs
- **THEN** it SHALL find exactly two full-release-required non-`real_pass` rows
- **AND** the row IDs SHALL be `REQ-04-016` and `REQ-05-003`.

#### Scenario: Capital-flow rows require a verified public runtime provider

- **GIVEN** P90 evaluates capital-flow fields
- **WHEN** it claims a row as `real_pass`
- **THEN** it SHALL prove a no-login/no-paid/no-authorization/no-Level2/no-high-frequency runtime provider was verified and used
- **AND** it SHALL prove `date`, `net_inflow`, and `net_outflow` were persisted in SQLite and read back through product APIs or UI
- **AND** parser-only, fixture, stub, accepted-local, or manually seeded evidence SHALL NOT upgrade those rows.

#### Scenario: Directional net-flow semantics are explicit

- **GIVEN** the public H5 capital-flow history endpoint exposes a daily net-flow value
- **WHEN** P90 stores the value
- **THEN** positive net flow SHALL map to `net_inflow`
- **AND** negative net flow SHALL map to `net_outflow`
- **AND** the raw daily net-flow value SHALL be preserved as `raw_net_flow`.

#### Scenario: P90 final claims remain evidence gated

- **GIVEN** P90 writes final release materials
- **WHEN** any full-release-required row remains partial, blocked, scoped-only, unsupported, or unverified
- **THEN** P90 SHALL NOT claim full original-requirement pass
- **AND** it SHALL list exact remaining rows and blockers.

### Requirement: P91 GitHub release and Docker deployment

After P90, the project SHALL provide a GitHub-ready release and Docker Compose deployment path that can initialize and run the product without embedding secrets or overwriting user data.

#### Scenario: Install script detects first install versus upgrade

- **GIVEN** a user runs `bash scripts/install.sh`
- **WHEN** no local deployment state or data directory exists
- **THEN** the script SHALL initialize local deployment directories, create `.env` from `.env.example` when needed, start Docker Compose, and run health checks
- **AND** when existing deployment state or data is present it SHALL route through the upgrade path instead of deleting or reinitializing user data.

#### Scenario: Runtime secrets are supplied outside the package

- **GIVEN** the release package and Docker image are built
- **WHEN** the user configures LLM credentials
- **THEN** `DEEPSEEK_API_KEY`, base URL, model, and timeout SHALL come from `.env` or environment variables
- **AND** no complete API key SHALL be committed, baked into the image, or written to release manifests.

#### Scenario: Uninstall preserves data by default

- **GIVEN** a deployed instance has local SQLite, VecLite, backup, log, and `.env` data
- **WHEN** the user runs `bash scripts/uninstall.sh`
- **THEN** containers and networks MAY be removed
- **AND** local data SHALL be preserved by default
- **AND** deleting local data SHALL require `--purge` and an exact confirmation phrase.

#### Scenario: GitHub release automation remains evidence gated

- **GIVEN** GitHub Actions creates a release artifact
- **WHEN** CI or release packaging runs
- **THEN** it SHALL run OpenSpec validation, Go tests, frontend tests/build, deployment checks, and package verification
- **AND** it SHALL upload release artifacts without claiming physical second-machine validation, broker connectivity, trading, automatic confirmation, external push, Level2 data, paid/login sources, or return guarantees.

### Requirement: P92 final original requirement audit ledger

After P91, the project SHALL provide a final original-requirement audit ledger that independently summarizes whether every original requirement row is covered by final acceptance evidence.

#### Scenario: Full-release rows must all be final real pass

- **GIVEN** P75 generated the original requirement traceability matrix
- **AND** P88 produced the latest full 341-row evidence matrix
- **AND** P89 and P90 produced final blocker overlays
- **WHEN** P92 generates the final audit ledger
- **THEN** every full-release-required row SHALL have final status `real_pass`
- **AND** reference-only rows SHALL remain separated from product pass claims
- **AND** the ledger SHALL fail validation if any full-release-required row is missing, stale, or non-`real_pass`.

#### Scenario: Ledger includes operational review dimensions

- **GIVEN** an original requirement row is included in the final ledger
- **WHEN** P92 writes the row
- **THEN** it SHALL include the requirement id, source section, requirement text, final status, feature area, UI/product surface, expected behavior or data impact, readback or audit evidence, acceptance command or artifact, and boundary notes.

#### Scenario: Final audit claims remain bounded

- **GIVEN** P92 summarizes final release readiness
- **WHEN** it describes accepted scope
- **THEN** it MAY claim original product requirements are accepted for the local/GitHub-Docker release scope
- **AND** it SHALL NOT claim physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.

### Requirement: P93 final code reality and design audit

After P92, the project SHALL provide a code-facing release audit that checks whether original requirements are backed by real implementation rather than demo, placeholder, hardcoded, or dead-code behavior.

#### Scenario: Production implementation evidence is mapped

- **GIVEN** the original requirements are accepted by P92
- **WHEN** P93 audits code reality
- **THEN** it SHALL map original requirement sections to concrete production Go, React, configuration, script, and release files
- **AND** it SHALL cross-check the P92 row-level ledger so every original requirement row resolves to a current code/evidence bundle through its source section
- **AND** it SHALL identify whether each requirement area is backed by runtime code, UI, tests, and evidence
- **AND** it SHALL keep P92 as the 341-row row-level artifact rather than replacing it with a coarser P93 claim.

#### Scenario: Demo and hardcoding risks are classified

- **GIVEN** suspicious terms such as `demo`, `mock`, `stub`, `placeholder`, `fake`, `dummy`, `TODO`, `FIXME`, or hardcoded values appear in the repository
- **WHEN** P93 evaluates them
- **THEN** it SHALL classify each material occurrence as test-only, config-only, documentation-only, accepted local fallback, or release-blocking
- **AND** release-blocking occurrences SHALL be fixed or reported as blockers.

#### Scenario: Secret literals are blocked

- **GIVEN** current non-test source or configuration files may contain local credentials
- **WHEN** P93 evaluates hardcoding risk
- **THEN** it SHALL scan for unredacted `sk-...` API key literals using a bounded token pattern
- **AND** any such literal in scanned non-test source/config files SHALL be release-blocking until removed or replaced with an empty/user-supplied runtime value.

#### Scenario: Final claims remain bounded

- **GIVEN** P93 passes
- **WHEN** release readiness is described
- **THEN** it MAY claim the implementation has passed final code reality and design audit for the local/GitHub-Docker release scope
- **AND** it SHALL NOT claim physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.

### Requirement: P95 public engineering validation hardening

P95 SHALL make public repository validation stable across clean checkouts and local developer checkouts with frontend dependencies installed.

#### Scenario: Backend package discovery excludes frontend dependencies

- **GIVEN** frontend dependencies have been installed under `web/node_modules`
- **WHEN** backend validation selects Go packages for tests
- **THEN** packages below `web/node_modules` SHALL NOT be included
- **AND** the selection helper SHALL fail if a package from frontend dependency trees is selected.

#### Scenario: P93 source scan ignores local runtime artifacts

- **GIVEN** ignored local runtime artifacts exist under project paths such as `cmd/agent/tmp/`
- **AND** nonignored new source files may exist before they are committed
- **WHEN** P93 code reality audit runs in check mode
- **THEN** ignored local runtime artifacts SHALL NOT change the report
- **AND** tracked plus nonignored untracked release-relevant source files SHALL be eligible for scanning
- **AND** tracked release-relevant files SHALL still be scanned for secrets, demo/stub risks, and release-boundary violations.

#### Scenario: API route contract is checked

- **GIVEN** backend handlers register `/api/v1` routes
- **WHEN** the API route contract check runs
- **THEN** every registered route SHALL be documented in `docs/api.md` or `docs/frontend-contract.md`
- **AND** documented route examples with query strings SHALL normalize to their path identity
- **AND** the check SHALL fail when docs reference a route that is no longer registered.

#### Scenario: Local SQLite runtime is concurrency-aware

- **GIVEN** the local server opens a SQLite database
- **WHEN** the database is file-backed
- **THEN** the connection SHALL enable foreign key enforcement and a bounded busy timeout
- **AND** it SHALL attempt WAL mode for local UI/background-task read-write concurrency
- **AND** in-memory tests SHALL remain supported.

#### Scenario: Docker deployment supports file-based LLM secrets

- **GIVEN** an operator supplies a `DEEPSEEK_API_KEY_FILE`
- **WHEN** the application loads runtime configuration
- **THEN** the key SHALL be read from that file if `DEEPSEEK_API_KEY` is not set
- **AND** committed configuration and documentation SHALL NOT contain real keys.

### Requirement: P96 public documentation front door

P96 SHALL make the public repository understandable to a new reader without requiring them to start from archived phase logs.

#### Scenario: Root README introduces the product honestly

- **GIVEN** a user opens the GitHub repository
- **WHEN** they read the root `README.md`
- **THEN** they SHALL see the product purpose, supported local workflows, architecture/data-flow visuals, installation entrypoint, documentation map, CI/release status, and safety boundaries
- **AND** the README SHALL NOT claim broker connectivity, automatic trading, one-click trading, delegated order placement, external push, automatic confirmation, automatic rule application, return guarantees, paid/login/authorization-only sources, Level2 data, or high-frequency data.

#### Scenario: Documentation map is concise

- **GIVEN** a maintainer opens `docs/README.md`
- **WHEN** they use it as a navigation page
- **THEN** it SHALL point to product, architecture, API, data model, workflow, frontend, deployment, governance, release evidence, and history documents
- **AND** it SHALL NOT require reading the full P0-P96 phase log to find normal documentation.

#### Scenario: Historical release evidence remains available

- **GIVEN** P96 moves or summarizes phase history
- **WHEN** a reader needs release caveats or acceptance history
- **THEN** the history SHALL remain discoverable under `docs/release/`
- **AND** P96 SHALL NOT erase historical limitations, scoped claims, or Not Claimed boundaries.

#### Scenario: Requirements truth source remains stable

- **GIVEN** `docs/requirements.md` is the L1 product requirement truth source
- **WHEN** P96 adds public-facing docs
- **THEN** public docs SHALL link to requirements for full details
- **AND** P96 SHALL NOT rewrite L1 requirement semantics as marketing copy.

### Requirement: P97 default local config file

The local runtime SHALL treat `configs/config.yaml` as the preferred local configuration file and keep `configs/config.example.yaml` as a committed template.

#### Scenario: Local server uses config.yaml by default

- **GIVEN** `INVESTMENT_AGENT_CONFIG` is unset
- **AND** `configs/config.yaml` exists
- **WHEN** the server or shared config loader starts with no explicit config path
- **THEN** it SHALL load `configs/config.yaml`
- **AND** it SHALL NOT load `configs/config.example.yaml` instead.

#### Scenario: Fresh checkout remains runnable

- **GIVEN** `INVESTMENT_AGENT_CONFIG` is unset
- **AND** `configs/config.yaml` does not exist
- **WHEN** the shared config loader starts with no explicit config path
- **THEN** it SHALL fall back to `configs/config.example.yaml`.

#### Scenario: Local config is not committed

- **GIVEN** a user creates `configs/config.yaml`
- **WHEN** Git status is checked
- **THEN** the file SHALL be ignored by default
- **AND** real local keys SHALL remain outside committed source.

### Requirement: P98 SHALL harden release runtime mode and frontend redaction reuse

P98 SHALL add release-mode guardrails and shared frontend redaction without changing investment runtime capabilities.

#### Scenario: Release runtime rejects stub data

- **GIVEN** runtime mode is configured as `release`
- **AND** `data_sources.use_stub` is `true`
- **WHEN** configuration validation runs
- **THEN** validation SHALL fail with an actionable message
- **AND** release/Docker defaults SHALL keep `data_sources.use_stub=false`.

#### Scenario: Development fallback remains available

- **GIVEN** runtime mode is omitted or configured as `development`
- **WHEN** local example or test configuration enables stub data
- **THEN** validation SHALL continue to allow local stub data
- **AND** this SHALL NOT create a release claim for real provider operation.

#### Scenario: Frontend diagnostic redaction is shared

- **GIVEN** frontend pages or components display diagnostic or failure text
- **WHEN** the text contains key-shaped tokens, SQL fragments, prompt fragments, raw diagnostic payloads, stack traces, or local paths
- **THEN** the text SHALL be redacted through a shared utility
- **AND** current page/component tests SHALL continue to prove sensitive details are not displayed.

### Requirement: P101 unified local config path

The project SHALL use `configs/config.yaml` as the default ignored local config path for current local runtime and current local-source acceptance scripts.

#### Scenario: Historical acceptance scripts use the runtime default

- **GIVEN** a user configures LLM and local runtime settings in `configs/config.yaml`
- **WHEN** current local-source acceptance scripts are run without explicit override variables
- **THEN** they SHALL read `configs/config.yaml` by default
- **AND** they SHALL NOT require a separate `configs/config.local.yaml` file.

#### Scenario: Explicit overrides remain available

- **GIVEN** an operator needs a one-off private config file
- **WHEN** a script-specific variable such as `P71_LOCAL_CONFIG`, `P72_LOCAL_CONFIG`, `P75_LOCAL_CONFIG`, or `P63_LOCAL_CONFIG` is provided
- **THEN** that script SHALL use the explicit path
- **AND** this override SHALL NOT change the default documented local config path.

### Requirement: P101 OpenAI-compatible local LLM request compatibility

The local analyst LLM client SHALL remain compatible with OpenAI Chat Completions gateways that expect JSON accept headers, stable user-agent identification, and longer bounded response times.

#### Scenario: Compatible headers are sent

- **GIVEN** a configured OpenAI-compatible LLM gateway
- **WHEN** the analyst client sends a chat completion request
- **THEN** it SHALL send `Accept: application/json`
- **AND** it SHALL send a stable `User-Agent`
- **AND** it SHALL continue using the configured `<base_url>/chat/completions` path.

#### Scenario: Transport timeout is retried once

- **GIVEN** the first LLM request times out before receiving response headers
- **WHEN** the retry succeeds
- **THEN** the analyst client SHALL return parsed analysis material
- **AND** it SHALL mark metadata with a bounded timeout retry
- **AND** it SHALL NOT loosen the local parser or quality gate.

#### Scenario: Default timeout allows slower compatible gateways

- **GIVEN** a local config omits `deepseek.timeout_seconds`
- **WHEN** defaults are applied
- **THEN** the configured timeout SHALL be 60 seconds.

#### Scenario: Release claims stay bounded

- **GIVEN** P101 changes script defaults and LLM request compatibility
- **WHEN** release readiness is described
- **THEN** the project MAY claim local config path consistency for source-runtime validation
- **AND** it MAY claim OpenAI-compatible LLM request compatibility for Chat Completions style gateways
- **AND** it SHALL NOT claim Docker installation, package distribution, GitHub Release, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.

### Requirement: P102 product acceptance audit

The project SHALL support a product-level acceptance audit after real LLM access is configured.

#### Scenario: Product audit evidence is captured

- **GIVEN** local backend, frontend, SQLite, VecLite, and real LLM config are available
- **WHEN** P102 product acceptance is executed
- **THEN** the audit SHALL capture current-run screenshots for key product workflows
- **AND** it SHALL assess UX, design reasonableness, accessibility risks, data/readback trust, and safety boundaries.

#### Scenario: Release claims remain bounded

- **GIVEN** P102 writes product acceptance findings
- **WHEN** release readiness is described
- **THEN** the project MAY claim product-level local-source acceptance only for the checked local runtime scope
- **AND** it SHALL NOT claim Docker installation, package distribution, GitHub Release, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.

### Requirement: P103 product acceptance UX linkage fixes

The product SHALL address P102 non-blocking UX findings without expanding investment runtime capabilities.

#### Scenario: Portfolio empty state is onboarding-safe

- **GIVEN** no local portfolio snapshot exists
- **WHEN** the user opens the portfolio page
- **THEN** the page SHALL present first-use onboarding and local account calibration guidance instead of a generic system failure.

#### Scenario: Decision analysis remains auditable without overwhelming the page

- **GIVEN** a decision contains real LLM analyst reports
- **WHEN** the user opens the decision detail page
- **THEN** the page SHALL show the final verdict and safety boundary first
- **AND** the full analysis material SHALL remain available through explicit expansion.

#### Scenario: Decision loop deep link focuses the target

- **GIVEN** a decision-loop URL includes `decision_id`
- **WHEN** the linked decision is present in the loop response
- **THEN** the page SHALL focus that decision's loop record and keep trace links read-only.

#### Scenario: Release claims remain bounded

- **GIVEN** P103 fixes P102 UX findings
- **WHEN** release readiness is described
- **THEN** the project MAY claim the checked local product UX issues were fixed
- **AND** it SHALL NOT claim Docker installation, package distribution, GitHub Release, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns.

### Requirement: P104 Product Operation Linkage Acceptance
The project SHALL maintain a repeatable local-source acceptance gate that verifies representative product operations through HTTP APIs, SQLite side effects, downstream readback, audit traceability, and forbidden automation absence.

#### Scenario: Local runner validates linked product behavior
- **GIVEN** the repository source tree is available locally
- **WHEN** the P104 acceptance runner is executed
- **THEN** it SHALL create an isolated temporary SQLite database and config
- **AND** it SHALL start the local backend on localhost
- **AND** it SHALL exercise representative portfolio, decision confirmation, review, audit, notification, risk, and data-quality operations through HTTP APIs
- **AND** it SHALL verify durable SQLite side effects and downstream readback
- **AND** it SHALL fail if forbidden broker/order/push/automatic-confirmation evidence is present.

#### Scenario: Acceptance record stays honest about scope
- **GIVEN** P104 validation has completed
- **WHEN** the release acceptance record is updated
- **THEN** it SHALL distinguish fresh P104 local-source linkage evidence from Docker, installer, package, remote deployment, physical second-machine, broker, automatic trading, automatic confirmation, automatic rule application, and return-guarantee claims.

### Requirement: P105 current release version v0.1.1

The repository SHALL declare `v0.1.1` as the current local source release version after P100-P104 validation has passed and P105 release gates have completed.

#### Scenario: Current version metadata is synchronized

- **GIVEN** P105 release validation has passed
- **WHEN** a user or release operator inspects version metadata
- **THEN** the root `VERSION` file SHALL contain `v0.1.1`
- **AND** `web/package.json` SHALL declare version `0.1.1`
- **AND** the root package entry in `web/package-lock.json` SHALL declare version `0.1.1`.

#### Scenario: P105 release claims stay bounded

- **GIVEN** `v0.1.1` is described in release materials
- **WHEN** release readiness is communicated
- **THEN** the project MAY claim local source product acceptance through P104 and current source version metadata synchronization
- **AND** it SHALL NOT claim Docker installation validation, package refresh, GitHub Release workflow success, physical second-machine validation, broker connectivity, trading, one-click trading, order placement, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability, or investment returns unless separately validated.

### Requirement: P106 release package scan compatibility for v0.1.2

The repository SHALL keep release package prompt-payload scanning strict while avoiding source-level false positives in redacted UI labels.

#### Scenario: Package scanner passes redacted UI labels

- **GIVEN** the release package script scans tracked source files
- **WHEN** frontend redaction labels are inspected
- **THEN** caller-specific replacement labels SHALL NOT use long JSON-like `prompt: "..."` payload shapes that match prompt-payload forbidden-content rules
- **AND** the release package smoke and verify steps SHALL pass before `v0.1.2` is tagged.

#### Scenario: v0.1.2 patch release stays bounded

- **GIVEN** `v0.1.2` is described in release materials
- **WHEN** release readiness is communicated
- **THEN** the project MAY claim the release-package scan compatibility fix and current source version metadata synchronization
- **AND** it SHALL NOT claim Docker installation validation, physical second-machine validation, broker connectivity, trading, automatic confirmation, automatic rule application, future provider availability, or investment returns unless separately validated.
