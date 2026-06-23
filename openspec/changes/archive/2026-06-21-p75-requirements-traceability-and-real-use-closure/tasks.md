# Tasks: P75 Requirements Traceability And Real Use Closure

## 1. Change Setup And User-Concern Capture

- [x] 1.1 Confirm no other active OpenSpec change exists before P75 creation.
- [x] 1.2 Create OpenSpec change `p75-requirements-traceability-and-real-use-closure`.
- [x] 1.3 Capture the user's concerns as explicit acceptance dimensions: original requirement comparison, dynamic fund input, external data completeness, built-in master wisdom/LLM usage, real UI operation, function/data linkage, analysis accuracy, UI design, and release-claim honesty.
- [x] 1.4 Run `openspec validate p75-requirements-traceability-and-real-use-closure --strict`, `openspec validate --all --strict`, and targeted whitespace checks for the new change files.
- [x] 1.5 Update active-change governance entries in `docs/GOVERNANCE.md`, `openspec/project.md`, and `openspec/PROGRESS.md` so P75 is not hidden behind "no active change" metadata.
- [x] 1.6 Before any runtime/API/UI/workflow implementation, decide whether P75 requires L1 deltas for `docs/requirements.md`, `docs/api.md`, `docs/data-model.md`, `docs/workflow.md`, and `docs/frontend-contract.md`; add matching OpenSpec deltas before implementation if behavior or contracts change.

## 2. Original Requirement Traceability Matrix

- [x] 2.1 Build `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md` from `docs/requirements.md` sections 1-19.
- [x] 2.2 Split every normative paragraph, bullet, table row, SOP step, acceptance criterion, risk/compliance statement, and glossary/appendix item into a stable atomic `requirement_id` with `source_section`, `source_start_line`, `source_end_line`, and `requirement_text_hash`.
- [x] 2.2a Define the matrix hash method in the artifact header: normalize line endings to LF, trim trailing whitespace per line, join the requirement text with `\n`, and store lowercase SHA-256 as `requirement_text_hash`.
- [x] 2.3 For every atomic requirement, assign one of `real_pass`, `scoped_pass`, `deterministic_local_evidence`, `partial`, `not_implemented`, or `blocked`.
- [x] 2.4 For every atomic requirement, record `criticality`, `criticality_reason`, `full_release_required`, `non_goal_basis`, `optional_basis`, `allowed_release_claim`, `delivered_by_change`, `verification_command`, `acceptance_artifact`, `evidence_freshness`, and `release_impact`.
- [x] 2.5 For every non-`real_pass` row, record the exact gap, product impact, release-claim impact, remediation decision, and whether the final conclusion must downgrade or block.
- [x] 2.6 Specifically mark P72/P73/P74 evidence as scoped where it only proves `510300`, selected workflows, temporary DBs, browser tasks, or readiness scenarios.
- [x] 2.7 Review whether any current release document overstates full-product completion; list required wording fixes.
- [x] 2.8 Replace or extend `coverage-review.md` after the atomic matrix exists so it reports section-by-section coverage and residual risk instead of treating planned coverage as proof.

## 3. Built-In Knowledge And LLM Usage Audit

- [x] 3.1 Compare `docs/requirements.md` master wisdom and conflict rules against the P74 registry and runtime workflow LLM context.
- [x] 3.2 Add failing tests for any mismatch where the workflow uses a hand-written subset instead of the structured registry/readiness context.
- [x] 3.3 Unify runtime LLM knowledge context with the P74 readiness registry or document any blocked limitation.
- [x] 3.4 Verify built-in knowledge remains LLM/background/rule-context only and cannot satisfy formal market evidence.

## 4. Dynamic Fund/ETF Data Resolution Audit And Hardening

- [x] 4.1 Audit current symbol handling to identify hardcoded `510300` or `000300` assumptions in backend, frontend, scripts, tests, and release claims.
- [x] 4.2 Add failing service tests proving a user-entered fund/ETF symbol must resolve fund profile, tracked index, market price, valuation, liquidity, formal evidence, RAG status, and safe degradation dynamically.
- [x] 4.3 Implement or harden a symbol-profile resolver that can safely support configured known symbols and block unknown symbols without fabricated readiness.
- [x] 4.4 Verify external collector/readiness behavior uses the user-entered symbol and its tracked index rather than fixed accepted symbols.
- [x] 4.5 Add at least one non-`510300` real or accepted-local user-entered fund/ETF scenario if safe public/configured data is available; otherwise record the blocker and release impact.
- [x] 4.6 In that non-`510300` scenario, trigger read-only market/evidence refresh and verify collector request parameters or request-construction artifacts, stored facts, `source_health`, `audit_events`, `data_date`, freshness, and readiness response are bound to the user-entered fund/ETF and its tracked index; preseeded facts or readiness rows alone cannot prove dynamic external querying.
- [x] 4.6a Record a correlation key across collector output, stored facts, `source_health`, `audit_events`, and readiness: user symbol, tracked index symbol, data category, source, `data_date`, and `request_id`.
- [x] 4.7 If public source, LLM, or collector execution fails, classify it using P52 categories (`network`, `rate_limit`, `authentication_or_key`, `source_schema_change`, `no_data`, `parse_failure`, `model_unavailable`, `quality_failure`, `redaction_failure`) and downgrade/block only the affected claims.

## 5. External And Built-In Data Completeness Audit

- [x] 5.1 Compare original required data categories against implemented collectors, tables, APIs, UI readbacks, and acceptance evidence.
- [x] 5.2 Classify market price, valuation, fund profile, tracked index, liquidity, formal evidence, RAG/index health, sentiment proxy, funds flow, margin financing, and constituent financials as `real_pass`, `scoped_pass`, `partial`, `not_implemented`, or `blocked`.
- [x] 5.3 Verify unstructured intelligence and anti-fake requirements field-by-field: `intelligence_items`, `intelligence_summary`, `source_verifications`, `rag_chunks`, source grade, URL, hash, summary, VecLite chunk status, source verification, time decay, and C-source background-only behavior.
- [x] 5.4 Verify F-1 through F-5 deterministically, including C-source exclusion from formal verdicts and local structured financial-number precedence.
- [x] 5.5 Verify required data categories cannot be silently substituted by built-in knowledge, local notes, C-level background, stub data, stale source-health entries, or LLM-generated text.
- [x] 5.6 Add readiness/API/UI tests for missing and degraded categories that affect safety margin, valuation, expected return, risk alerts, and trading-style suggestions.
- [x] 5.7 Build a missing-data propagation matrix: no media heat means no "normal emotion" claim; no margin financing means no normal financing claim; no constituent financials means no intact fundamentals claim; no funds flow means no neutral funds-flow claim; no valuation/liquidity/formal evidence means affected safety-margin, expected-return, alert, and SOP claims degrade or block.
- [x] 5.8 Build a field-level fund/index/benchmark join matrix covering fund symbol, fund profile, tracked index, benchmark symbol/profile, fund NAV/price/liquidity, index valuation, index constituents/financials, formal evidence, source, join key, as-of date, freshness, conflict handling, missing-benchmark behavior, and mismatch/stale behavior.

## 6. Analysis Accuracy And Data-Impact Acceptance

- [x] 6.1 Add deterministic checks for valuation-zone, liquidity-risk, source-verification, risk-alert, expected-return sample threshold, dynamic sell trigger, manual confirmation, and portfolio snapshot calculations.
- [x] 6.2 Add deterministic test vectors for every executable criterion in requirements 2.4 and 2.5, every R-1 through R-6 rule, 6.2 priority ordering, 6.5 state transitions, and 6.6 cooldown extension.
- [x] 6.3 Add threshold vectors for liquidity 20-day average/plan amount 20x, single-day 5%, emotion 90%/10% and 3-day abnormality, 2 independent A/S sources, PE/PB valuation zones, expected-return `<5` and `<20` sample gates, and cooldown/state-machine boundaries.
- [x] 6.4 Validate every expected-return output field in 9.4, every 9.5 trigger, 9.6 pressure/dynamic calibration behavior, and 9.7 sample-count downgrade behavior.
- [x] 6.4a For every expected-return probability output, verify sample count, sample interval, screening/filter conditions, cohort construction, and output provenance are displayed and trace back to stored facts.
- [x] 6.5 Validate portfolio allocation and rebalance behavior: core/satellite/cash classification, quarterly ±15% deviation, satellite-over-limit trigger, and "take-profit funds return to core assets" readback.
- [x] 6.6 Extend the real-use acceptance runner so UI actions are followed by SQLite verification and cross-page readback.
- [x] 6.7 Build an action-to-SQLite-table-to-readback matrix for adding, editing, importing, correcting, confirming, marking error, generating proposal, gatekeeper review, daily report, monthly review, and quarterly review.
- [x] 6.8 For every action matrix row, record expected changed tables, prohibited changed tables, required `audit_events`, and required readback pages.
- [x] 6.9 Verify every critical mutation writes a matching `audit_events` row and never creates broker, order, push, external-push, automatic-confirmation, or automatic-trading state.
- [x] 6.10 Verify daily, monthly, and quarterly self-check outputs: daily PE/redline/position facts; monthly P&L attribution, discipline audit, emotion log, and error-case statistics; quarterly benchmark comparison with benchmark source/symbol/as-of date/freshness, rule-effectiveness review, and evolution proposal summary.
- [x] 6.11 Verify evolution proposal types and gatekeeper checks: threshold adjustment, SOP addition, master-weight adjustment, behavior-pattern alert, root-rule violation, sample count `<3`, emotion bias, backtest degradation, rule conflict, and pass/deny/user-review outputs.

## 7. SOP A-F And Real User Scenario Coverage

- [x] 7.1 Map SOP A-F to existing tests, UI journeys, and data-impact evidence.
- [x] 7.2 Add missing real UI acceptance scenarios for holding drop, holding rise, hot-topic chasing, panic sell, macro gray-rhino, and black-swan event flows.
- [x] 7.3 For each SOP, verify rule priority, data prerequisites, LLM role, user confirmation behavior, derived page readback, and safe degradation.
- [x] 7.4 Mark any SOP that cannot be fully verified as `partial` or `blocked` with explicit release impact.
- [x] 7.5 For every SOP with user-visible behavior, require real browser UI operation; API evidence may only supplement rule priority, data prerequisite, and database assertions.

## 8. UI Design And Real-Use UX Review

- [x] 8.1 Review primary user flows from a real user's perspective: onboarding, adding fund, checking data readiness, consulting, reading decision detail, confirming offline action, reviewing alerts, marking errors, reviewing proposals, and inspecting audit trail.
- [x] 8.2 Verify desktop and mobile layouts for misleading copy, hidden next actions, text overflow, state ambiguity, weak error recovery, and trading-like affordances.
- [x] 8.3 Add screenshots and browser result evidence for any UI flow added or changed by P75.
- [x] 8.4 Record UI findings as pass/fix/block in the P75 acceptance record.
- [x] 8.5 Build a critical UI flow matrix covering onboarding, add fund, data readiness, consultation, decision detail, alerts, offline confirmation, error marking, rule proposal, gatekeeper pass/deny/user-review states, monthly review, quarterly review, audit trail, and settings/safety boundaries.
- [x] 8.6 For every critical UI flow matrix row, require `requirement_id`, `ui_flow_id`, browser action, DOM assertion, expected SQLite changes, prohibited SQLite changes, audit event, readback page, mobile result, failure-state result, screenshot path, and `pass`/`fix`/`block` status.
- [x] 8.6a Add one continuous non-`510300` browser journey row: add fund -> data readiness -> consultation or alerts -> SQLite verification -> derived page readback, using the same user-entered symbol and tracked index correlation key from 4.6a.
- [x] 8.7 Add a UX misunderstanding checklist for trading boundary, state language, next action, evidence insufficiency, user offline execution, in-system confirmation, and account-state mutation.
- [x] 8.8 Failure-state UI checks must include unsupported symbol, insufficient data, stale data, degraded source, model unavailable, validation error, gatekeeper deny, and gatekeeper user-review.

## 9. Release Claim And Documentation Closure

- [x] 9.1 Update release materials so they distinguish full-requirement pass, scoped pass, partial implementation, blocked data, and out-of-scope boundaries.
- [x] 9.2 Enumerate every non-`real_pass` atomic requirement in the final release conclusion; grouped gap categories are not sufficient.
- [x] 9.3 Update governance/progress/development-plan docs with P75 status and exact conclusion.
- [x] 9.4 Update `docs/release/README.md`, `docs/release/acceptance-repeatability.md`, release candidate/handoff materials, and acceptance record index with the P75 conclusion.
- [x] 9.5 Do not archive P75 until the traceability matrix and acceptance record agree on the final release conclusion.
- [x] 9.6 If the conclusion is not `release_ready_full_requirements_traceable`, list the next required changes before any broader release claim.

## 10. Verification And Self-Review

- [x] 10.1 Run targeted backend tests for knowledge context, symbol resolution, readiness, calculations, and data-impact behavior.
- [x] 10.2 Run targeted frontend tests for readiness/readback/UI flow changes.
- [x] 10.3 Run full `go test ./...`, `npm --prefix web test -- --run`, `npm --prefix web run build`, P75 acceptance runner, safety scans, `openspec validate --all --strict`, and `git diff --check`.
- [x] 10.3a Record that 10.3 is command execution evidence only; P75 now separately closes deterministic vectors, SOP A-F accepted-local real UI, failure-state UI, mark-error, gatekeeper, and 8.6a non-510300 browser journey evidence, while live-provider/arbitrary-symbol/full-data-domain gaps remain scoped.
- [x] 10.4 Review whether every user-raised concern is represented in P75 tasks and evidence.
- [x] 10.5 Review whether any task result is only mock/fixture/demo evidence and downgrade the related traceability row unless real or accepted-local evidence exists.
- [x] 10.6 Review whether P75 actually completed or only documented gaps; final report must say which one happened.
- [x] 10.7 P75 final release conclusion must cite and satisfy P52 G0-G9, P66 strict current-data policy, P67 resolution state, P71-P74 repeatability rules, and P52 G6/G7 failure classifications, or explicitly downgrade/block each affected claim.
- [x] 10.7a For P71-P74 repeatability, either rerun the concrete repeat commands/scripts from `docs/release/acceptance-repeatability.md` and cite artifacts, or explicitly record inherited historical/scoped evidence as scoped evidence and downgrade affected claims.
- [x] 10.8 Store G9 safety scan output plus human review summary, and confirm forbidden terms only appear in prohibition/boundary contexts; the scan must include P52 terms plus P75 expanded terms: `券商接口`, `自动交易`, `一键交易`, `代下单`, `外部推送`, `自动确认`, `自动应用规则`, `自动规则应用`, `自动修复`, `自动升级`, `自动迁移`, `自动恢复`, `自动覆盖真实库`, `真实库覆盖`, `收益承诺`, `登录源`, `付费源`, `授权源`, `Level2`, and `高频`.
- [x] 10.9 Archive-before-final check: no requirement may be marked `real_pass` if its only evidence is screenshot-only, route-smoke-only, API-only for user-visible behavior, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, temporary-DB-only, or single-symbol-only.
