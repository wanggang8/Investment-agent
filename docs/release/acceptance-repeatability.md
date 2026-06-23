# Acceptance Repeatability

> Updated: 2026-06-22
> Applies after: P53 `release_ready`; P66/P67/P68/P70 current-data policy and release-governance rules apply to future release-ready claims; P71 adds strict full real product acceptance rules; P72 adds real user fund scenario and SQLite data-impact acceptance rules; P73 adds product-effectiveness/UX validation rules; P75 adds original-requirement traceability and real-use closure rules; P76 adds post-P75 clean package refresh rules; P77 adds post-P75 real-pass upgrade gate rules; P78 adds requirements real-pass batch closure rules; P79 adds real-use data-impact closure rules; P80 adds review/audit/governance field-readback closure rules; P81 adds dynamic source field coverage closure rules; P82 adds SOP/action UI-to-SQLite closure rules; P83 adds governance traceability UI/API/SQLite closure rules; P84 adds portfolio/confirmation real UI data-impact closure rules; P85 adds expected-return analysis-accuracy UI/API/SQLite closure rules.

This document defines how to repeat release acceptance consistently after P53. It preserves the P52 G0-G9 gate model and the P53 rule that release readiness must be based on actual results or explicit waivers.

## Output Directory

Use a fresh label for every repeat run:

```bash
RUN_DIR="tmp/acceptance/<label>"
mkdir -p "$RUN_DIR/logs" "$RUN_DIR/data"
```

Do not commit `tmp/acceptance/**`, temporary SQLite databases, raw logs, private configs, or full provider responses.

## Command Order

Run gates in this order:

1. G0 governance and diff checks.
2. G1 Go full test suite.
3. G2 focused Go integration packages.
4. G3 frontend tests and build.
5. G4 browser E2E smoke.
6. G5 local fixture/current smoke.
7. G6 real public source opt-in.
8. G7 real LLM opt-in.
9. G8 local install and release upgrade.
10. G9 safety and redaction review.

## Retry Rule

One exact-command retry is allowed only for transient local execution conditions:

- OS process kill.
- Port startup race.
- Browser or Vite startup race.
- Short-lived local resource pressure.

The repeat record must keep both the first failure and the retry result. A second failure of the same command is `blocked` unless a release owner records an explicit waiver with scope and release impact.

Retries are not allowed to hide:

- Test assertion failures.
- Build type errors.
- OpenSpec validation errors.
- Safety or redaction failures.
- Parse failures from reachable real-source responses.
- LLM quality failures.

## G5 Current Data Rule

`fixture` mode is the deterministic regression baseline.

| Result | Release treatment |
| --- | --- |
| fixture pass + current pass | pass |
| fixture pass + current degraded + failed=0 | degraded, non-blocking, limit current local DB health claims |
| fixture pass + current failed>0 | blocked unless explicitly waived |
| fixture failed | blocked |

Current-data degraded status must include the command output summary and release impact. It must not be described as full current-data health.

After P66, future release-ready claims must also run the current data policy gate:

```bash
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate
```

| Policy gate | Release treatment |
| --- | --- |
| `gate=pass` | Supports a clean current-data claim for this run. |
| `gate=waiver_required` | Requires explicit waiver reason, impact, and scope; must not be called clean. |
| `gate=block` | Blocks release-ready claims unless the release scope explicitly excludes current local data health. |

The non-strict command remains useful for diagnostics, but it is not sufficient for a clean release-ready claim when the strict gate blocks.

After P67, if the strict policy gate is blocked or waiver-required, repeat acceptance must also run the local resolution check:

```bash
go run ./cmd/agent --task data-source-quality-resolution-check --symbol 000300
```

| Resolution state | Release treatment |
| --- | --- |
| `claim_state=pass` | Supports a clean current-data gate claim for this run. |
| `claim_state=resolved_with_waiver` | May state a waiver was recorded; must not describe current data as clean. |
| `claim_state=resolved_with_scope_exclusion` | May state current local data health is excluded from the clean claim; must not describe current data as healthy. |
| `claim_state=requires_resolution` | Blocks future release-ready claims that depend on current local data health. |

After P68, release records must also state the release readiness scope:

| P66/P67 state | Required release status wording |
| --- | --- |
| P66 `gate=pass` and P67 `claim_state=pass` | `release_ready` may include a clean current-data gate claim for that run. |
| P66 blocked or waiver-required with P67 `resolved_with_scope_exclusion` | Use `release_ready_limited_current_data_scope`; do not include current local data health in clean claims. |
| P66 waiver-required with P67 `resolved_with_waiver` | Use a waiver-qualified release status and document scope, reason, and release impact. |
| P67 `requires_resolution` | Use `release_blocked` for any release claim depending on current local data health. |

If package evidence predates later committed release-governance or runtime changes, the release record must either limit the package evidence to that candidate archive or open a package-refresh stage from a clean tree before final distribution.

After P70, the milestone decision was `release_ready_limited_current_data_scope` with no mandatory next phase for that limited scope. Future repeat records must preserve the P66/P67 current-data limitation unless a fresh strict gate passes or a new resolution record changes the claim state.

After P71, a full real product acceptance claim requires:

```bash
bash scripts/p71-real-product-acceptance.sh
```

The P71 strict run must pass all of these gates:

| Gate | Required pass evidence |
| --- | --- |
| Current data | `policy=passed` and `gate=pass` for `000300`; not P67 scope exclusion |
| Gate resolution | `claim_state=pass` and clean-data gate claim allowed |
| VecLite | rebuild/index health positive and healthy; consultation retrieval `fallback_source=veclite`, `index_health=healthy`, and freshness not `unknown` or `stale` |
| Real UI | primary routes and key operations pass against a real local Go backend and Vite frontend |
| Real LLM | LLM smoke passes and consultation displays parsed, quality-passed LLM analyst material |
| Safety | no forbidden trading, broker, push, auto-confirm, auto-rule, auto-repair, migration, restore, real DB overwrite, return-promise, login/paid/authorized source, Level2, or high-frequency affordance |

If any P71 strict gate fails, use `release_blocked_current_data`, `release_blocked_retrieval_index`, `release_blocked_ui_or_llm`, or `release_blocked_safety_or_package` instead of full real product acceptance.

After P72, a real user scenario acceptance claim requires:

```bash
bash scripts/p72-real-user-fund-scenario-acceptance.sh
```

The P72 run must pass all of these gates:

| Gate | Required pass evidence |
| --- | --- |
| Real scenario setup | `510300` portfolio calibration, holding edit, import confirmation, correction, and offline transaction are operated through UI |
| Formal evidence | Public evidence refresh writes formal evidence for the scenario; safe-degraded formal evidence collection is not sufficient |
| Local knowledge/RAG | Local knowledge validate/confirm writes intelligence rows, rebuilds VecLite, and SQLite `rag_chunks.index_status` includes an indexed P72 chunk |
| Real consultation | LLM consultation completes with parsed/passed analyst reports, `retrieval.status=hit`, `fallback_source=veclite`, and healthy index |
| Manual confirmation | Generated decision accepts a manual offline confirmation and records `executed_manually` without automatic trading |
| Data impact | Read-only SQLite checker verifies deterministic cash, total assets, positions, transactions, confirmations, daily report, risk alert, notifications, source verification, and forbidden-table absence |
| UI runtime | Page errors, unexpected API failures, failed resources, and console errors are blockers |
| Safety | No broker/order/trade execution/external push tables or visible forbidden affordances |

If public evidence collection is unavailable, P72 must not be marked pass by relying on consultation safe degradation. Record `release_blocked_formal_evidence_collection` or implement a verified read-only source/collection fix before rerunning.

After P73, a product-effectiveness/UX acceptance claim requires:

```bash
bash scripts/p73-product-effectiveness-ux-validation.sh
```

The P73 run must pass all of these gates:

| Gate | Required pass evidence |
| --- | --- |
| Real browser UX tasks | Browser results and screenshots exist under `docs/release/ui-audit-assets/2026-06-19-p73/` and cover daily discipline, portfolio, evidence/data quality, decision traceability, manual confirmation, risk/review/rules, mobile, and unsafe-input tasks |
| Effect replay | `scripts/p73_effect_replay_check.py` returns a sanitized pass summary for background-only blocking, manual-confirmation-only portfolio mutation, risk/readback links, rule-effect gate state, and forbidden-table absence |
| Safe degradation | Background-only or insufficient evidence produces a non-trade/insufficient-data result and is not counted as a product-effectiveness pass |
| UX audit | Information hierarchy, next-action clarity, state labels, navigation, mobile/reflow, and copy safety are reviewed and recorded |
| Safety | No broker/order/trade execution/external push affordances or tables; no automatic confirmation, automatic rule application, or future return guarantee |

If local server startup is blocked, record `release_blocked_local_port_binding`. If real browser evidence is missing, record `release_blocked_missing_real_ux_evidence`. SQLite-only replay, safe degradation, or fixture preparation is not sufficient for a P73 pass.

## P69 Clean Package Repeat Rule

To repeat the current P69 package evidence, use the P69 package archive and write new repeat output under project `tmp/`:

```bash
bash scripts/local-release-package.sh --verify tmp/p69-final-release/20260618T084011Z/investment-agent-p69-clean-tree.tar.gz --output-dir tmp/p69-final-release-rerun
bash scripts/local-release-repeat-acceptance.sh --archive tmp/p69-final-release/20260618T084011Z/investment-agent-p69-clean-tree.tar.gz --output-dir tmp/p69-final-repeat-rerun
```

The current P69 package source commit is `cc0a64781e199a7745432b63bce26de4402042b5` with `source_status=clean`. The package includes committed source through P68 and does not include P69/P70 acceptance documentation unless a later post-P70 package refresh is performed.

If a future handoff archive must include P69 or P70 documentation, run a new package-refresh change and record the new source commit, checksum, verify summary, and repeat summary. That refresh is optional for the current P70 limited local release decision.

## G6 Real Public Source Rule

Use a temporary SQLite database and a temporary config. When `data_sources.use_stub=false`, satisfy the config prerequisite before running public evidence refresh:

- Provide a valid `data_sources.market_endpoint`, or
- Enable `data_sources.market_collectors.enabled=true` with valid collector sources and base URLs.

For public evidence refresh, record:

- Command.
- Date window.
- Symbol.
- Exit code.
- Redacted SQLite counts for `intelligence_items`, `intelligence_summary`, `rag_chunks`, `source_verifications`, and `audit_events`.
- Failure classification when applicable.

Do not commit the temporary database or raw provider payloads.

## P74 Built-In Knowledge/Data Readiness Rule

To repeat P74 readiness acceptance, run:

```bash
bash scripts/p74-built-in-knowledge-data-readiness.sh
```

The runner must use a temporary SQLite database and local HTTP server/frontend. It must cover:

- `510300` complete readiness.
- Missing valuation data.
- Background-only evidence/knowledge.
- Single-source evidence.
- Multi-source formal evidence.
- Out-of-scope symbol profile.
- `/data-quality` UI readiness panel.
- Decision detail LLM readiness readback.
- 390px mobile reflow and forbidden affordance scan.

Artifacts belong under `docs/release/ui-audit-assets/2026-06-19-p74/` or a new dated rerun directory. Do not commit temporary SQLite databases, full prompts, raw provider payloads, local configs, private paths, or complete keys.

## P75 Original-Requirement Traceability Rule

To repeat P75 traceability acceptance, run:

```bash
python3 scripts/p75_requirements_traceability_check.py --check
```

The runner must regenerate:

- `docs/release/acceptance/2026-06-20-p75-requirements-traceability-matrix.md`
- `docs/release/acceptance/2026-06-20-p75-real-use-closure.md`
- `docs/release/ui-audit-assets/2026-06-20-p75/traceability-summary.json`

The traceability matrix must be built from `docs/requirements.md` sections 1-19 and must keep stable atomic rows with source line ranges, lowercase SHA-256 `requirement_text_hash`, evidence status, release impact, gap, remediation decision, and acceptance artifact fields.

Use `release_ready_full_requirements_traceable` only when every `full_release_required=true` row is `real_pass` and the run also satisfies or freshly reruns the applicable P52 G0-G9, P66, P67, P71, P72, P73, and P74 gates. Inherited historical evidence, single-symbol evidence, temporary-DB evidence, screenshot-only evidence, route-smoke evidence, fixture/mock/stub evidence, waiver-only evidence, or scope-exclusion evidence must be treated as scoped or partial rather than full pass.

Use `release_pending_safety_review_scoped_with_traceability_gaps` when the atomic matrix contains scoped, deterministic-local-only, or partial rows and the expanded G9 scan still needs human boundary review. Use `release_ready_scoped_with_traceability_gaps` only when no blocking safety failure is found and the human boundary review confirms all forbidden-term matches are prohibition or boundary contexts. Use `release_blocked_requirements_traceability` when a product-critical requirement is blocked, not implemented, cannot safely degrade, or the safety review finds a forbidden affordance.

P75 must also record:

- hardcoded `510300`/`000300` scan impact,
- missing-data propagation for media heat, margin financing, constituent financials, funds flow, valuation, liquidity, and formal evidence,
- fund/index/benchmark join-key requirements,
- deterministic vectors for threshold and accuracy checks,
- UI/action matrix minimum columns,
- continuous non-`510300` real browser journey requirement,
- expanded G9 forbidden-term scan with human boundary-review summary.

## P76 Post-P75 Package Refresh Rule

To repeat P76 package evidence, use the generated archive and write new output under project `tmp/`:

```bash
bash scripts/local-release-package.sh --verify tmp/p76-final-release/20260621T030713Z/investment-agent-p76-post-p75-final.tar.gz --output-dir tmp/p76-final-release-rerun
bash scripts/local-release-repeat-acceptance.sh --archive tmp/p76-final-release/20260621T030713Z/investment-agent-p76-post-p75-final.tar.gz --output-dir tmp/p76-final-repeat-rerun
```

The current P76 package source commit is `8a317f25917b8ff18ec9b5049e6a6188206a22d3` with `source_status=clean`. The archive has SHA-256 `7540429d0b6c3cdd09dad2ebb10e2356580faf0b05e6acd92bc3bd9763a3dcb7`, package verify passed, and isolated repeat acceptance passed from the extracted package workspace.

The P76 package includes committed P72-P75 acceptance Markdown and OpenSpec archives. It does not claim to include P76 package-after-the-fact evidence or `docs/release/ui-audit-assets/` screenshots/assets.

## P85 Expected-Return Analysis-Accuracy Rule

To repeat P85 expected-return analysis-accuracy acceptance, run:

```bash
P85_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-22-p85-expected-return-analysis" bash scripts/p85-expected-return-analysis-acceptance.sh
python3 scripts/p85_expected_return_analysis_accuracy_closure.py --check
```

The P85 run must use a temporary SQLite database, a real local Go backend, a Vite frontend, and browser UI operations through `/consultation` and decision detail readback. It must verify all of these gates:

| Gate | Required pass evidence |
| --- | --- |
| Complete-data expected return | UI consultation writes and reads sufficient-sample scenarios, sample/window metadata, target-return input, previous-base-midpoint input, and deterministic expected-return values. |
| Downside boundary | A downside scenario triggers deterministic downside evidence and remains non-trading/non-confirming. |
| Unavailable sample | Insufficient sample data is marked unavailable and does not fabricate scenarios or trade guidance. |
| SQLite readback | Read-only checker verifies persisted decisions, `expected_return_scenarios_json`, precision/sample fields, triggers, and forbidden confirmation/order/push absence. |
| Focused tests | Expected-return workflow and decision-detail handler tests pass for both API-shaped and workflow-persisted JSON. |
| Safety | LLM material cannot override the final rule verdict, create automatic confirmation, create broker/order/push state, or imply guaranteed returns. |

If `DEEPSEEK_API_KEY` is unavailable, the P85 record must explicitly use `llm_mode=static_fallback_no_real_llm_claim` and must not claim fresh real LLM output. If the deterministic workflow, UI/API/SQLite readback, or safety-negative checks fail, P85 must remain non-pass until a product fix and rerun provide new evidence.

P85 may upgrade only directly proven expected-return/readback rows. It must not upgrade rows that require future return accuracy, future market-direction accuracy, a longitudinal real-world outcome study, a real historical backtest model, automatic probability downshift, complete allocation/rebalance policy coverage, broker connectivity, automatic trading, external push, automatic confirmation, automatic rule application, login/paid/authorized source coverage, Level2, or high-frequency source coverage.

P76 preserves P75 `release_ready_scoped_with_traceability_gaps`. It must not be used as evidence for `release_ready_full_requirements_traceable`; that requires every applicable full-release row in the P75 traceability matrix to become `real_pass` and the applicable real UI/data/LLM/source/safety gates to pass freshly or with valid current evidence.

## P77 Requirements Real-Pass Upgrade Rule

To repeat P77 upgrade evidence, run the fresh evidence commands first and then regenerate/check the P77 matrix:

```bash
mkdir -p docs/release/ui-audit-assets/2026-06-21-p77
{ go test -v ./internal/infrastructure/llm/deepseek -run 'TestClientRetriesQualityFailureWithStricterBoundary|TestClientRejectsProhibitedLLMOutput|TestEvaluateQualityAllowsNormalAnalysisAndRejectsUnsafeClaims' -count=1 && go test -v ./internal/application/workflow -run 'TestExpectedReturnMaterialIsExplanatoryOnly' -count=1 && go test -v ./internal/domain/rule -run 'TestExpectedReturnDoesNotOverrideVerdict|TestP75RulePriorityAndRootRules' -count=1; } 2>&1 | tee docs/release/ui-audit-assets/2026-06-21-p77/safety-and-boundary-go-tests.log
{ go test -v ./internal/application/workflow -run 'TestPublicEvidencePayloadEnforcesSourceMetadataAndFormalBoundary|TestPublicEvidenceIngestionMajorEventsRequireTwoHighGradeIndependentSources|TestEvidenceVerificationRequiresTwoHighGradeIndependentSources|TestAnalystRequestsPreferStructuredFinancialFacts|TestPublicEvidenceIngestionAppliesF4TimeDecayAndBackgroundBoundary|TestPublicEvidencePayloadNormalizesEmotionalDescriptions' -count=1 && go test -v ./internal/infrastructure/persistence/sqlite -run TestMarketRepositoryPreservesStructuredFinancialFields -count=1 && go test -v ./internal/domain/rule -run TestEvaluatePriorityScenarios -count=1; } 2>&1 | tee docs/release/ui-audit-assets/2026-06-21-p77/f1-f5-go-tests.log
P75_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-21-p77-sop-failure" bash scripts/p75-sop-failure-real-ui-acceptance.sh
P75_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-21-p77-non-510300" bash scripts/p75-non-510300-real-ui-journey.sh
python3 scripts/p77_requirements_real_pass_upgrade.py --check
```

The runner must generate:

- `docs/release/acceptance/2026-06-21-p77-requirements-real-pass-upgrade-matrix.md`
- `docs/release/acceptance/2026-06-21-p77-real-pass-upgrade-acceptance.md`
- `docs/release/ui-audit-assets/2026-06-21-p77/real-pass-upgrade-summary.json`
- `docs/release/ui-audit-assets/2026-06-21-p77/safety-scan.txt`
- `docs/release/ui-audit-assets/2026-06-21-p77/safety-scan-review.json`
- `docs/release/ui-audit-assets/2026-06-21-p77/safety-and-boundary-go-tests.json`
- `docs/release/ui-audit-assets/2026-06-21-p77/f1-f5-go-tests.json`

The P77 runner validates the safety scan review status, concrete verbose Go test names, package pass lines, and generated sidecar metadata before allowing candidate rows to become `real_pass`.

P77 may upgrade a row to `real_pass` only when its applicable implementation, UI, data-impact, workflow/rule/LLM, scenario, and safety dimensions are directly evidenced. It must not use screenshot-only, route-smoke-only, fixture-only, mock/stub-only, waiver-only, scope-exclusion-only, temporary-DB-only, or incompatible single-symbol-only evidence as a real-pass basis.

If any full-release-required row remains non-`real_pass`, use `release_ready_scoped_with_p77_real_pass_progress`, not `release_ready_full_requirements_traceable`. P77 does not refresh the P76 package; a separate package-refresh change is required before claiming a distribution archive includes P77 evidence.

## P78 Requirements Real-Pass Batch Closure Rule

To repeat P78 batch A evidence, run:

```bash
mkdir -p docs/release/ui-audit-assets/2026-06-21-p78
{ go test -v ./internal/application/workflow -run 'TestBuildExpectedReturnIncludesSampleContextForAllPrecisionStates|TestBuildExpectedReturnProducesAdvisorySellEvaluation|TestBuildExpectedReturnDoesNotTriggerTargetWithoutConfiguredTarget|TestBuildExpectedReturnUsesScenarioBoundsForSellTriggers|TestBuildExpectedReturnCoversAllSellEvaluationTriggers|TestExpectedReturnNodeUsesWorkflowPricesForSellEvaluation|TestExpectedReturnNodeUsesMatchingSymbolPosition|TestExpectedReturnNodeUsesWorkflowDynamicSellInputs|TestExpectedReturnNodeIncludesP34SupportingDataContext|TestExpectedReturnSampleCountFromWorkflowDataUsesMarketHistory|TestExpectedReturnSampleCountFromWorkflowDataDoesNotInventSamples|TestBuildExpectedReturnExplainsMissingPriceContext' -count=1 && go test -v ./internal/domain/rule -run TestExpectedReturnDoesNotOverrideVerdict -count=1; } 2>&1 | tee docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.log
P75_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-21-p78-non-510300" bash scripts/p75-non-510300-real-ui-journey.sh
python3 scripts/p78_requirements_real_pass_batch_closure.py --check
```

The runner must generate:

- `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-matrix.md`
- `docs/release/acceptance/2026-06-21-p78-requirements-real-pass-batch-closure.md`
- `docs/release/ui-audit-assets/2026-06-21-p78/real-pass-batch-summary.json`
- `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-go-tests.json`
- `docs/release/ui-audit-assets/2026-06-21-p78/expected-return-ui-readback.json`

P78 may upgrade a row to `real_pass` only when its applicable implementation, UI, data/readback, workflow/rule/LLM, scenario, and safety dimensions are directly evidenced. It must not upgrade broad expected-return, SOP, portfolio, data-source, or product-goal rows using batch A evidence unless the exact row is covered.

If any full-release-required row remains non-`real_pass`, use `release_ready_scoped_with_p78_real_pass_batch_progress`, not `release_ready_full_requirements_traceable`. P78 does not refresh the P76 package; a separate package-refresh change is required before claiming a distribution archive includes P78 evidence.

## P79 Real-Use Data-Impact And Expected-Return Closure Rule

To repeat P79 evidence, run:

```bash
P72_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-21-p79-real-user-fund" bash scripts/p72-real-user-fund-scenario-acceptance.sh
P75_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-21-p79-non-510300" bash scripts/p75-non-510300-real-ui-journey.sh
python3 scripts/p79_real_use_data_impact_and_expected_return_closure.py --check
```

The runner must generate:

- `docs/release/acceptance/2026-06-21-p79-real-use-data-impact-and-expected-return-matrix.md`
- `docs/release/acceptance/2026-06-21-p79-real-use-data-impact-and-expected-return-closure.md`
- `docs/release/ui-audit-assets/2026-06-21-p79/real-use-data-impact-summary.json`

P79 may upgrade portfolio/local-account/confirmation rows only when fresh real UI operations, SQLite readback, audit evidence, and prohibited-table checks directly cover the row. Expected-return rows remain non-`real_pass` unless fresh P79 evidence proves the exact probability, scenario, sell-trigger, valuation, sample, provenance, and disclaimer fields.

If any full-release-required row remains non-`real_pass`, use `release_ready_scoped_with_p79_real_use_data_impact_progress`, not `release_ready_full_requirements_traceable`. P79 does not refresh the P76 package; a separate package-refresh change is required before claiming a distribution archive includes P79 evidence.

## P80 Review Audit Governance Closure Rule

To repeat P80 evidence, run:

```bash
P75_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-22-p80-review-audit-governance" bash scripts/p75-sop-failure-real-ui-acceptance.sh
python3 scripts/p80_review_audit_governance_closure.py --check
```

The runner must generate:

- `docs/release/acceptance/2026-06-22-p80-review-audit-governance-matrix.md`
- `docs/release/acceptance/2026-06-22-p80-review-audit-governance-closure.md`
- `docs/release/ui-audit-assets/2026-06-22-p80-review-audit-governance/review-audit-governance-summary.json`
- `docs/release/ui-audit-assets/2026-06-22-p80-review-audit-governance/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-22-p80-review-audit-governance/db-impact-check.log`

P80 may upgrade review/audit/governance rows only when fresh real UI operations, browser readback, SQLite field readback, gatekeeper audit fields, SOP lifecycle audits, and prohibited-table checks directly cover the row. Broad monthly attribution, full proposal-generation, final rule application time, and complete original-requirement rows remain non-`real_pass` unless a dedicated scenario proves the exact behavior.

If any full-release-required row remains non-`real_pass`, use `release_ready_scoped_with_p80_review_audit_governance_progress`, not `release_ready_full_requirements_traceable`. P80 does not refresh the P76 package; a separate package-refresh change is required before claiming a distribution archive includes P80 evidence.

## P81 Dynamic Source Field Coverage Rule

To repeat P81 evidence, run:

```bash
go test -v ./cmd/agent -run TestRunNon510300DynamicAcceptanceBindsCollectorSourceHealthAuditAndReadiness -count=1
P75_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-22-p81-dynamic-source-field-coverage" bash scripts/p75-non-510300-real-ui-journey.sh
python3 scripts/p81_dynamic_source_field_coverage.py --check
```

The runner must generate:

- `docs/release/acceptance/2026-06-22-p81-dynamic-source-field-coverage-matrix.md`
- `docs/release/acceptance/2026-06-22-p81-dynamic-source-field-coverage.md`
- `docs/release/ui-audit-assets/2026-06-22-p81-dynamic-source-field-coverage/dynamic-source-field-coverage-summary.json`
- `docs/release/ui-audit-assets/2026-06-22-p81-dynamic-source-field-coverage/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-22-p81-dynamic-source-field-coverage/non-510300-db-impact-summary.json`
- `docs/release/ui-audit-assets/2026-06-22-p81-dynamic-source-field-coverage/dynamic-source-go-test.log`

P81 may upgrade dynamic source/readiness rows only when fresh non-`510300` real UI/API/SQLite/readback evidence proves user-symbol request binding, tracked-index readiness, source-health provenance, formal evidence, RAG indexing, LLM context readback, auditability, and forbidden-capability absence. It does not claim paid/login/authorized sources, Level2, high-frequency data, future provider availability, or full original-requirement pass.

If any full-release-required row remains non-`real_pass`, use `release_ready_scoped_with_p81_dynamic_source_progress`, not `release_ready_full_requirements_traceable`. P81 does not refresh the P76 package; a separate package-refresh change is required before claiming a distribution archive includes P81 evidence.

## P82 SOP Action UI-To-SQLite Rule

To repeat P82 evidence, run:

```bash
P75_FINAL_RULE_APPLY=1 P75_ARTIFACT_DIR=$(pwd)/docs/release/ui-audit-assets/2026-06-22-p82-sop-action-ui-sqlite bash scripts/p75-sop-failure-real-ui-acceptance.sh
python3 scripts/p82_sop_action_ui_sqlite_closure.py --check
```

P82 may upgrade SOP/action rows only when fresh real browser UI, API/readback, read-only SQLite evidence, audit events, explicit user final confirmation where applicable, and forbidden-capability absence directly prove the original row. Rows broader than the SOP/action evidence must remain deferred with an exact next owner.

If any full-release-required row remains non-`real_pass`, use `release_ready_scoped_with_p82_sop_action_progress`, not `release_ready_full_requirements_traceable`. P82 does not refresh the P76 package; a separate package-refresh change is required before claiming a distribution archive includes P82 evidence.

## P83 Governance Traceability Backfill Rule

To repeat P83 evidence, run:

```bash
P83_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability" bash scripts/p83-governance-traceability-acceptance.sh
python3 scripts/p83_governance_traceability_backfill.py --check
```

The runner must generate:

- `docs/release/acceptance/2026-06-22-p83-governance-traceability-matrix.md`
- `docs/release/acceptance/2026-06-22-p83-governance-traceability-backfill.md`
- `docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability/governance-traceability-summary.json`
- `docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability/db-readback-check.log`
- `docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability/go-handler-tests.log`
- `docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability/go-workflow-tests.log`
- `docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability/go-agent-tests.log`

P83 may upgrade review/governance/release traceability rows only when fresh real browser UI, API readback, read-only SQLite evidence, focused Go test evidence, redaction/safety evidence, and forbidden-capability absence directly prove the original row. It must not upgrade portfolio allocation, expected-return analysis, crawler/VecLite/multi-agent implementation breadth, dashboard/product-goal, or knowledge/RAG rows using governance-only evidence.

If any full-release-required row remains non-`real_pass`, use `release_ready_scoped_with_p83_governance_traceability_progress`, not `release_ready_full_requirements_traceable`. P83 does not refresh the P76 package; a separate package-refresh change is required before claiming a distribution archive includes P83 evidence.

## P84 Portfolio Confirmation Data Impact Rule

To repeat P84 evidence, run:

```bash
P84_ARTIFACT_DIR="$PWD/docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation" bash scripts/p84-portfolio-confirmation-acceptance.sh
python3 scripts/p84_portfolio_confirmation_data_impact_closure.py --check
```

The runner must generate:

- `docs/release/acceptance/2026-06-22-p84-portfolio-confirmation-data-impact-matrix.md`
- `docs/release/acceptance/2026-06-22-p84-portfolio-confirmation-data-impact-closure.md`
- `docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation/portfolio-confirmation-summary.json`
- `docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation/db-readback-check.log`
- `docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation/go-handler-tests.log`

P84 may upgrade portfolio/confirmation rows only when fresh real browser UI, API readback, read-only SQLite field evidence, downstream UI readback, focused Go test evidence, and forbidden broker/order/push/auto-confirm absence directly prove the original row. It must not upgrade broad allocation policy, quarterly rebalance, sell-only/frozen-watch transition, monthly attribution, full audit/proposal application-time, source readiness, dashboard/product-goal, or release-governance rows using portfolio-confirmation-only evidence.

If any full-release-required row remains non-`real_pass`, use `release_ready_scoped_with_p84_portfolio_confirmation_progress`, not `release_ready_full_requirements_traceable`. P84 does not refresh the P76 package; a separate package-refresh change is required before claiming a distribution archive includes P84 evidence.

## G7 Real LLM Rule

Use a local private config or a temporary config derived from it. Release materials may include:

- Model name.
- Exit code.
- Parse status.
- Quality status.
- Redacted audit summary.

Release materials must not include complete API keys, full prompts, raw model responses, private paths, or vendor payloads.

If G7 fails:

- `authentication_or_key` blocks LLM capability claims.
- `model_unavailable` blocks LLM capability claims.
- `quality_failure` blocks LLM quality claims.
- Network or quota issues may be waived only with explicit release impact.

## Release Status Decision

Use these rules:

- Any non-waived G0-G5, G8, or G9 `blocked` result -> `release_blocked`.
- Any safety or redaction failure -> `release_blocked`.
- A second failure after an allowed retry -> `release_blocked` unless explicitly waived.
- G6 or G7 failure must block the corresponding real-source or LLM claim; whether it blocks the whole release must be documented with release impact.
- `release_ready` or a qualified release-ready status such as `release_ready_limited_current_data_scope` is allowed only when all blocking gates pass or have explicit waivers/scope exclusions and all degraded/skipped results are documented.

## Required Record Format

Every repeated run should create a new acceptance record:

```markdown
# Acceptance Run: <label>

- Date:
- Code-under-test commit:
- Branch:
- Operator:
- Environment:
- Result: release_ready/release_ready_limited_current_data_scope/release_ready_scoped_with_traceability_gaps/release_ready_scoped_with_p81_dynamic_source_progress/release_pending_safety_review_scoped_with_traceability_gaps/release_ready_full_requirements_traceable/release_blocked/release_blocked_requirements_traceability

| Gate | Status | Command | Artifact | Notes | Release impact |
| --- | --- | --- | --- | --- | --- |
```

The record must cite `docs/project-acceptance-gate-matrix.md` and this repeatability document.
