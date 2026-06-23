# P75 Real Use Closure Acceptance

> Date: 2026-06-20
> Result: `release_ready_scoped_with_traceability_gaps`

## Conclusion

P75 has produced an atomic traceability matrix and real-use closure evidence scaffolding. It does not claim full original-requirement completion unless every `full_release_required=true` row is `real_pass`. Current generated evidence is conservative and keeps scoped/partial gaps visible. The expanded G9 forbidden-term scan has been reviewed by category and found no forbidden runtime affordance, so the current result may be scoped release-ready, but not full original-requirement pass.

## Status Counts

- `blocked`: 0
- `deterministic_local_evidence`: 17
- `not_implemented`: 0
- `partial`: 291
- `real_pass`: 0
- `scoped_pass`: 33

## Key Release Impact

- P72/P73/P74 remain scoped evidence, especially around `510300`, temporary DBs, and selected UI journeys.
- Dynamic non-`510300` accepted-local request/source-health/readiness evidence now exists for `159915 -> 399006`; live public-provider availability and full arbitrary-symbol coverage remain unclaimed.
- Missing funds flow, margin financing, constituent financials, media heat, benchmark, valuation, liquidity, or formal evidence must propagate to dependent claim downgrades.
- UI pass requires real browser actions, DOM/readback, SQLite truth table, audit events, mobile checks, and failure-state checks.

## Dynamic Non-510300 Accepted-Local Evidence

- Scenario: `159915` 创业板 ETF with tracked index `399006`.
- Verification: `go test ./cmd/agent -run 'TestRunNon510300DynamicAcceptanceBindsCollectorSourceHealthAuditAndReadiness' -count=1`.
- Request proof: the test executes real local CLI tasks for `market-refresh --symbol 159915` and `public-evidence-refresh --symbol 159915 --start-date 2026-06-01 --end-date 2026-06-30`; the local HTTP server asserts market `symbol=159915`, CNInfo `stock=159915`, SZSE `keyword=159915`, and the explicit CNInfo date window.
- Stored-fact proof: SQLite `market_snapshots.market_metrics_json` contains `p34_source_health` for `symbol_profile`, `fund_profile`, `tracked_index`, `market_price`, `valuation_percentiles`, `liquidity`, and `sentiment_proxy`; fund-side categories bind to `159915`, index-side categories bind to `399006`, with `data_date=2026-06-19`, freshness, source, and `request_id`.
- Formal-evidence proof: SQLite `intelligence_summary`, `rag_chunks`, and `source_verifications` contain two formal public-evidence summaries for `159915`; source verification is `satisfied`; chunk metadata carries source and the same public-evidence ingestion `request_id` as the audit event.
- Readiness proof: `KnowledgeReadinessService` returns known profile `159915 -> 399006`; required readiness categories checked by the test include `symbol_profile`, `tracked_index`, `market_price`, `valuation_percentiles`, `liquidity`, `formal_evidence`, and `rag_index`, all without fabricating unknown-symbol support. `999999` remains blocked by the existing unknown-symbol test.
- Correlation proof: source-health-backed readiness dependencies expose `source_name`, `source_type`, `data_date`, `request_id`, and `affected_symbols` through the API/UI contract so a degraded or ready UI state can be traced back to the collector run.
- LLM-context proof: `go test ./internal/application/workflow -run 'TestAnalystRequestsScopeSymbolProfileKnowledgeToWorkflowSymbol|TestAnalystRequestsIncludeKnowledgeReadinessContext' -count=1` verifies analyst requests use the shared registry/readiness context, include `symbol_profile.159915` for `159915`, and do not silently include `symbol_profile.510300` for that non-510300 flow.
- Release boundary: this is accepted-local evidence for dynamic routing, correlation, and safe readiness behavior. It is not a claim that live CNInfo/SZSE/Eastmoney/CSIndex providers will always be reachable, nor a claim that every arbitrary ETF/fund profile is supported.

## P52 Failure Classification Evidence

- Verification: `go test ./internal/application/workflow -run 'TestPublicEvidenceIngestionAuditsPartialSourceFailures|TestPublicEvidenceIngestionAuditsSourceFailures|TestAnalystServiceUnavailableIncludesStableCategory' -count=1`.
- Verification: `go test ./internal/application/workflow -count=1` confirms the mapped classifications do not break collector, evidence, LLM, expected-return, and workflow tests.
- Public evidence failures now keep the source prefix but use P52 categories in audit `error_code`, for example `cninfo:network` and `cninfo:no_data`.
- Analyst failures in audit output refs map internal categories to P52 categories, for example LLM timeout becomes `category=model_unavailable`; missing key maps to `authentication_or_key`; quality gate failure maps to `quality_failure`; parse/empty output maps to `parse_failure`.
- Release boundary: P52 classification only scopes/downgrades affected claims. It does not make a failed provider pass, and it does not retry network/provider failures into false success.

## Degraded Data Propagation Evidence

- Verification: `go test ./internal/application/service -run 'TestKnowledgeReadinessServicePropagatesCriticalDataGapsToFeatureImpacts' -count=1`.
- Verification: `go test ./internal/application/service -run TestKnowledgeReadinessServiceDoesNotSubstituteStubBackgroundOrLLMForRequiredData -count=1`.
- Verification: `go test ./internal/application/handler -run 'TestGetKnowledgeReadinessReturnsDegradedDependencyImpacts' -count=1`.
- Verification: `npm --prefix web test -- --run DataQualityPage.test.tsx`.
- Coverage: valuation percentile parse failure degrades safety-margin and expected-return claims; missing liquidity degrades risk alerts and trade-like sizing suggestions; insufficient formal evidence degrades consultation/decision detail/risk alerts and states that no trading confirmation may be generated.
- No-substitution proof: required categories with `freshness=stubbed` remain degraded, `formal_evidence=background_only` remains degraded, and `llm_context=ready` does not turn missing/stubbed required facts into ready data.
- UI readback: `/data-quality` renders `估值分位 · 降级`, `流动性 · 降级`, `正式证据 · 降级`, plus the safe-degradation text for safety margin, expected-return precision, large/market-style action suggestions, and trade confirmation.
- Release boundary: this closes readiness/API/UI propagation for the tested degraded categories. Deterministic calculation vectors and SOP A-F browser scenarios are closed separately below within accepted-local scope; full live-provider, arbitrary-symbol, and full action-to-SQLite-to-readback breadth remains scoped/partial.

## Anti-Fake Rule Deterministic Evidence

- F-1 source metadata: `go test ./internal/application/workflow -run TestPublicEvidencePayloadEnforcesSourceMetadataAndFormalBoundary -count=1` verifies public evidence without a valid source level or evidence role is rejected, while C-level material is retained only as background.
- F-2 major-event verification: `go test ./internal/application/workflow -run 'TestPublicEvidenceIngestionMajorEventsRequireTwoHighGradeIndependentSources|TestEvidenceVerificationRequiresTwoHighGradeIndependentSources' -count=1` plus `go test ./internal/domain/rule -run TestEvaluatePriorityScenarios -count=1` verifies major-event A+B evidence remains failed and rule arbitration freezes insufficient high-grade major events.
- F-3 structured financial precedence: `go test ./internal/infrastructure/persistence/sqlite -run TestMarketRepositoryPreservesStructuredFinancialFields -count=1` and `go test ./internal/application/workflow -run TestAnalystRequestsPreferStructuredFinancialFacts -count=1` verify local structured market/financial fields survive SQLite roundtrip and are injected into analyst requests with `structured_facts_override_text_claims`.
- F-4 time decay: `go test ./internal/application/workflow -run TestPublicEvidenceIngestionAppliesF4TimeDecayAndBackgroundBoundary -count=1` verifies 0-24h/1-7d/7-30d/>30d weights and >30d background-only treatment before verification writes.
- F-5 objective wording: `go test ./internal/application/workflow -run TestPublicEvidencePayloadNormalizesEmotionalDescriptions -count=1` verifies common emotional wording is converted before hash/RAG/analysis ingestion.
- Release boundary: these rows are deterministic local evidence, not a full-product real-pass. Full pass still requires live-provider coverage, full data category collection/readback, SOP A-F browser scenarios, and the complete UI/action truth table.

## Analysis Deterministic Evidence

- Valuation boundary: `go test ./internal/domain/rule -run TestEvaluateValuationHighRiskBoundaryAtEightyPercent -count=1` verifies PE/PB 80% enters high-risk/no-new-buy treatment.
- Unknown allocation guard: `go test ./internal/domain/rule -run TestEvaluateDoesNotTriggerAllocationWhenRatiosAreUnknown -count=1` verifies missing core/satellite ratios do not fabricate rebalance guidance.
- Rule vectors: `go test ./internal/domain/rule -count=1` covers valuation zones, liquidity prohibitions, source verification, sentiment, take-profit, allocation, expected-return non-override, and proposal state transitions.
- P75 executable criteria vectors: `go test ./internal/domain/rule -run 'TestP75' -count=1` covers 2.4/2.5 sentiment inputs, 20-day liquidity 20x, same-day 5%, R-1 through R-6, rule priority, `normal`/`sell_only`/`frozen_watch` position-state mapping, and 3-trigger/5-day cooldown extension.
- Risk alert boundary: `go test ./internal/application/service -run TestRiskAlertServiceUsesValuationHighRiskBoundaryAtEightyPercent -count=1` verifies UI-facing risk alerts use the same 80% high-valuation boundary as rule arbitration.
- Expected-return vectors: `go test ./internal/application/workflow -run 'TestBuildExpectedReturn|TestExpectedReturnNode' -count=1` covers sample-count precision gates, dynamic sell evaluation, matching symbol position, sample provenance, and missing-price degradation.
- Expected-return detail/API/UI readback: `go test ./internal/application/handler -run 'TestDecisionDetailFromWorkflowExpectedReturnUsesWorkflowSampleCount|TestDecisionDetailFromRecordRestoresMarketContextSnapshot|TestDecisionDetailExpectedReturn' -count=1` and `npm --prefix web test -- --run DecisionTrace.test.tsx` verify symbol, date, current price/NAV, PE/PB percentiles, sample count, sample window, screening condition, scenario range/probability/trigger, sell evaluation, reassessment trigger, and disclaimer are rendered from workflow or stored context snapshot facts.
- Risk-alert vectors: `go test ./internal/application/service -run 'TestRiskAlert|TestSourceHealthRisk' -count=1` covers source-health-backed degraded-data alert inputs and risk alert persistence/readback behavior.
- Portfolio/confirmation vectors: `go test ./internal/application/service -run 'TestPortfolioService|TestConfirmationService' -count=1` covers portfolio snapshot math, edit/remove/import/correction rollback, manual confirmation, stale confirmation rejection, and sell snapshot preservation.
- Portfolio allocation/readback vectors: `go test ./internal/domain/rule -run TestP75PortfolioAllocationAndTakeProfitReadback -count=1` and `npm --prefix web test -- --run DecisionTrace.test.tsx` verify core underweight, satellite over-limit, and take-profit funds returning to core assets are exposed as manual optional-action readback.
- Daily/monthly/quarterly review vectors: `go test ./internal/application/handler -run 'TestTodayDailyDisciplineReport|Test.*DailyDiscipline|TestGetReviewSummary' -count=1`, `go test ./internal/application/workflow -run 'TestDailyAutoRun|TestRunDaily' -count=1`, and `npm --prefix web test -- --run DailyDisciplineReportDetailPage.test.tsx ReviewSummaryPage.test.tsx WorkbenchPage.test.tsx` verify daily discipline, review summaries, rule-effect tracking, degraded review notifications, and UI readback.
- Evolution/gatekeeper vectors: `go test ./internal/application/workflow -run 'TestEvolutionProposalGraph|TestGatekeeper' -count=1`, `go test ./internal/application/handler -run 'TestRuleProposal|TestRuleEffect' -count=1`, and `npm --prefix web test -- --run RulesPage.test.tsx RuleProposalPanel.test.tsx` verify threshold/SOP/capability/risk-rule proposal families with P75 subtypes, sample-count guardrails, gatekeeper pass/deny/user-review states, validation/backtest/conflict handling, and no automatic rule application before final confirmation.
- Release boundary: this closes P75 tasks 6.1, 6.2, 6.3, 6.4, 6.4a, 6.5, 6.9, 6.10, and 6.11 for deterministic local/API/UI checks. It does not convert live-provider, arbitrary-symbol, or full-data-domain gaps into `real_pass`.

## P75 SOP / Failure-State Real UI Evidence

- Verification: `bash scripts/p75-sop-failure-real-ui-acceptance.sh` completed with 1 passed Chromium test.
- Browser proof: `p75-sop-failure-real-ui.spec.ts` covers SOP-A holding drop, SOP-B holding rise, SOP-C hot-topic chasing, SOP-D panic sell, SOP-E macro gray-rhino, SOP-F black-swan event, unsupported symbol, insufficient data, stale/degraded source, model unavailable, validation error, gatekeeper deny, gatekeeper user-review, mark-error, proposal send-to-gatekeeper, and 390px mobile checks.
- SQLite proof: `docs/release/ui-audit-assets/2026-06-20-p75-sop-failure/db-impact-check.log` reports `status=passed`, `lifecycle_audits=6`, `sop_updated_status_count=6`, `mark_error_cases=1`, `mark_error_audits=1`, `gatekeeper_node_audits=6`, `gatekeeper_status=pending_final_confirm`, `after_position_transactions=0`, unchanged `rule_versions`, and `forbidden_broker_order_push_tables=0`.
- UI design finding: P75 fixed two real auditability issues found during browser acceptance: audit rows now expose input/output references and decision/proposal/confirmation/error-case associations; risk-alert cards now expose SOP context, data prerequisites, and LLM role instead of hiding them in backend JSON.
- Release boundary: this closes P75 tasks 6.9, 7.2, 7.3, 7.5, 8.1, 8.2, 8.4, and 8.8 within accepted-local real-browser scope. It still does not claim live external provider completeness, every arbitrary fund/index branch, or every original atomic requirement as `real_pass`.

## Repeated Real UI Scenario Evidence

- Verification: `bash scripts/p72-real-user-fund-scenario-acceptance.sh` reran successfully on 2026-06-20 after P75 hardening.
- Precheck proof: public evidence refresh, P34 expanded refresh, strict current-data gate, and LLM smoke completed successfully with no trading action.
- Browser proof: Playwright completed `p72-real-user-fund-scenario.spec.ts` with 1 passed Chromium test covering portfolio calibration/edit/import/correction/offline transaction, local knowledge import, VecLite rebuild, market/data-quality review, real LLM-backed consultation, decision detail, manual offline confirmation, daily discipline report, risk alerts, notifications, decision-loop/audit/review/rules/workbench readbacks, screenshots, and forbidden-affordance checks.
- SQLite proof: `docs/release/ui-audit-assets/2026-06-18-p72/db-impact-check.log` reports `status=passed`, `workflow_status=completed`, `confirmation_status=executed_manually`, `analyst_report_count=3`, expected cash/asset/position aggregates, committed local import/correction facts, decision-linked/manual offline confirmations, position transactions, user-confirm audit events, daily report, risk alerts, notifications, and no forbidden tables.
- Hardening proof: a real value-analyst quality-gate failure was traced to `quality_failed`; the LLM client now performs one stricter safety reprompt for quality failures while continuing to reject repeated unsafe output. It does not retry network, HTTP, parse, timeout, or missing-key failures into a false pass.

## Safety Scan Summary

- Pattern: `券商接口|自动交易|一键交易|代下单|外部推送|自动确认|自动应用规则|自动规则应用|自动修复|自动升级|自动迁移|自动恢复|自动覆盖真实库|真实库覆盖|收益承诺|登录源|付费源|授权源|Level2|高频`
- Human boundary review status: `reviewed_pass`
- Matches reviewed: 823
- Needs manual follow-up count: 0
- Human review summary: All P75 forbidden-term matches were reviewed by category. Matches are prohibition/boundary copy, negative tests, scan configuration, historical/release governance, captured acceptance evidence, seed prohibited-action fixtures, sanitizer/gatekeeper rules, or non-data use of 高频. No broker, auto-trading, one-click trading, delegated order, external push, auto-confirm, auto-rule, auto-repair, auto-upgrade/migration/restore, real database overwrite, return-promise, login/paid/authorized source, Level2, or high-frequency data-source affordance was found.
- Release impact: no clean full-release claim is allowed because requirement traceability remains scoped/partial, but G9 no longer blocks scoped release-ready wording.

### Safety Classification Counts

| Category | Count |
| --- | --- |
|`acceptance_script_boundary`|7|
|`governance_or_release_boundary_doc`|302|
|`historical_boundary_docs`|411|
|`non_data_frequency_label`|1|
|`runtime_boundary_or_sanitizer`|35|
|`runtime_sanitizer_or_gatekeeper_rule`|9|
|`scan_configuration`|4|
|`test_assertion_or_fixture`|54|

- Sample matches:
  - `scripts/p73_effect_replay_check.py:96:    require("自动交易" in (risk_alert["prohibited_actions_json"] or ""), "risk alert must expose prohibited actions", risk_alert)`
  - `scripts/p73_effect_replay_check.py:107:    require("不自动应用规则" in rule_effect["safety_note"], "rule effect safety note missing no-auto-apply boundary", rule_effect)`
  - `scripts/local-release-repeat-acceptance.sh:225:        "Level2 data",`
  - `docs/superpowers/specs/2026-06-08-p32-daily-discipline-report-productization-design.md:26:- 外部推送、券商接口或自动执行交易。`
  - `docs/superpowers/specs/2026-06-08-p32-daily-discipline-report-productization-design.md:154:P32 不新增交易执行、券商接口、外部推送、登录源、付费源或高频抓取。报告中的行动项只允许人工复核、查看证据、补齐数据、记录线下计划。报告不得承诺收益、不得预测确定涨跌、不得覆盖规则裁决。`
  - `scripts/local-release-package.sh:329:        "Level2 data",`
  - `openspec/PROGRESS.md:114:| 每日纪律报告产品化 | P32 `done` | P32 已归档到 `openspec/changes/archive/2026-06-11-p32-daily-discipline-report-productization/`；已将 Daily workflow 与 P31 auto-run 结果统一为今日纪律报告索引，提供今日报告、历史列表、详情回看和本地 smoke 验证入口；继续保持人工复核、只读展示和不自动交易边界 |`
  - `openspec/PROGRESS.md:116:| 账户与持仓录入/校准体验 | P33 `done` | P33 已归档到 `openspec/changes/archive/2026-06-12-p33-account-position-onboarding/`；已完成本地账户初始化、持仓录入/维护、线下交易流水、一致性校验、批量导入、错误修正、Portfolio 首次使用引导，以及 Dashboard/每日纪律缺前提跳转；继续保持本地事实记录、人工复核和不自动交易边界 |`
  - `openspec/PROGRESS.md:117:| 真实数据覆盖扩展 | P34 `done` | P34 已归档到 `openspec/changes/archive/2026-06-15-p34-real-data-coverage-expansion/`；已扩展中证指数样本/权重/估值、情绪替代指标 fixture、freshness/失败分类、source health、工作流输入和前端状态展示；继续保持只读、低频、可降级和不自动交易边界 |`
  - `openspec/PROGRESS.md:118:| 风险预警与 SOP 编排 | P35 `done` | P35 已归档到 `openspec/changes/archive/2026-06-16-p35-risk-alert-sop-orchestration/`；已建立本地风险预警事实、SOP 状态流转、通知/审计/每日纪律报告联动和风险预警中心；继续保持人工复核、只读追踪、不自动交易和不外部推送边界 |`
  - `openspec/PROGRESS.md:119:| 规则进化效果验证 | P36 `done` | P36 已归档到 `openspec/changes/archive/2026-06-16-p36-rule-evolution-effect-validation/`；已建立规则提案来源解释、样本代表性、过拟合检查、历史回放、守门人门禁接入、应用后追踪和前端展示；继续保持本地只读、守门人审计、用户最终确认、不自动应用规则和不自动交易边界 |`
  - `openspec/PROGRESS.md:121:| RAG / VecLite 检索质量加固 | P38 `done` | P38 已归档到 `openspec/changes/archive/2026-06-16-p38-rag-veclite-retrieval-quality/`；已建立本地 retrieval quality fixture、quality-aware ranking、source verification / RAG metadata 一致性校验、index freshness、审计/API/前端展示和 `retrieval-quality-smoke`；继续保持 source verification、规则裁决、守门人审计、人工确认和不自动交易边界 |`
  - `openspec/PROGRESS.md:122:| 前端完整用户旅程与全路径 E2E | P39 `done` | P39 已归档到 `openspec/changes/archive/2026-06-16-p39-frontend-full-user-journey-e2e/`；已补齐空库 onboarding、配置/账户初始化、市场/证据刷新、每日纪律、主动咨询、确认记录、风险预警、复盘和规则治理的浏览器级验收，并覆盖 console/a11y/窄屏与降级路径；继续禁止券商接口、自动交易、外部推送、自动规则应用和收益承诺 |`
  - `openspec/PROGRESS.md:123:| 本地部署、运维与恢复演练 | P40 `done` | P40 已归档到 `openspec/changes/archive/2026-06-16-p40-local-deploy-ops-recovery-drill/`；已建立本地预检、启动诊断、备份恢复 smoke、数据源健康面板、日志/临时文件/诊断文件治理和恢复演练文档；继续禁止券商接口、自动交易、外部推送、自动规则应用和收益承诺 |`
  - `openspec/PROGRESS.md:126:| 数据质量可观测 | P43 `done` | P43 已归档到 `openspec/changes/archive/2026-06-16-p43-data-quality-observability/`；已新增 `/data-quality` 只读质量面板，聚合数据源健康、证据与检索、LLM 质量和受影响工作流；继续禁止交易、外推、自动修复、自动确认和自动规则能力 |`
  - `openspec/PROGRESS.md:129:| 本地知识导入治理 | P46 `done` | 已归档到 `openspec/changes/archive/2026-06-17-p46-local-knowledge-import-governance/`；已提供 validate、脱敏预览、显式确认、C/background 背景事实写入和 pending 索引计划；继续禁止券商接口、自动交易、外部推送、自动确认、自动应用规则和收益承诺 |`
  - `openspec/PROGRESS.md:130:| 决策闭环解释 | P47 `done` | 已归档到 `openspec/changes/archive/2026-06-17-p47-decision-loop-explainability/`；已实现只读串联建议、确认、线下记录、风险/审计/复盘线索和缺口说明；继续禁止券商接口、自动交易、外部推送、自动确认、自动应用规则和收益承诺 |`
  - `openspec/PROGRESS.md:131:| 数据源质量回归包 | P48 `done` | 已归档到 `openspec/changes/archive/2026-06-17-p48-data-source-quality-regression-pack/`；已提供 fixture/current 数据源质量回归、freshness/失败分类验证、脱敏摘要、只读 API 和 CLI 脱敏审计；继续禁止券商接口、自动交易、外部推送、自动确认、自动应用规则、自动修复承诺和收益承诺 |`
  - `openspec/PROGRESS.md:132:| 本地发布与升级体验 | P49 `done` | 已归档到 `openspec/changes/archive/2026-06-17-p49-local-release-upgrade-experience/`；已提供本地版本检查、升级前备份提醒、迁移前预检、升级后 smoke 汇总、脱敏诊断 JSON 和脚本摘要；继续禁止自动升级、自动迁移、自动修复、覆盖真实库、交易、外推和收益承诺 |`
  - `openspec/PROGRESS.md:147:| 发布打包与版本标记 | P64 `done` | 已归档到 `openspec/changes/archive/2026-06-18-p64-release-packaging-version-tagging/`；已实现本地发布包脚本、sidecar manifest、archive checksum、verify 模式、输出目录限制、tracked/allowlisted source staging、敏感信息/禁止路径扫描和 P64 packaging handoff；不新增运行时交易、券商、外推、自动确认、自动升级、自动迁移、自动恢复或自动修复能力 |`

## Hardcoded Symbol Scan

- Matches for `510300|000300`: 947
- These are not automatically blockers, but any full-product claim must distinguish accepted-path tests from dynamic symbol support.

## Knowledge Context Scan

- Matches for expanded master/context IDs in workflow: 4
- P75 also runs focused Go tests for analyst knowledge readiness context and verifies workflow prompts use the shared P74 registry summary builder.

## Missing Data Propagation Matrix

| Data category | Dependent claim | Required treatment | Note |
| --- | --- | --- | --- |
|media_heat|normal emotion / no extreme sentiment|degrade_or_block|No media heat means no normal-emotion claim.|
|margin_financing|normal financing|degrade_or_block|No margin financing means no normal financing claim.|
|constituent_financials|intact fundamentals / buy logic not broken|degrade_or_block|No constituent financials means fundamentals cannot be declared intact.|
|funds_flow|neutral funds flow|degrade_or_block|No funds flow means no neutral funds-flow claim.|
|valuation|safety margin / expected return|degrade_or_block|No valuation means no reliable valuation or margin-of-safety claim.|
|liquidity|trade-like sizing / market order safety|block|No liquidity means no large or market-style action suggestion.|
|formal_evidence|formal verdict / buy logic breakage|block|No formal evidence means no trade-like confirmation.|

## Field-Level Fund/Index/Benchmark Join Matrix

| Category | Join key | Field | Source | Freshness | Treatment |
| --- | --- | --- | --- | --- | --- |
|fund_profile|fund symbol|fund code|fund source|as-of date required|missing blocks dynamic fund readiness|
|tracked_index|fund profile -> tracked index|tracked index symbol|index source|as-of date required|missing blocks index valuation claims|
|fund_price_liquidity|fund symbol|price/liquidity|fund market source|freshness required|stale degrades alerts/positions|
|index_valuation|tracked index symbol|PE/PB percentiles|index valuation source|freshness required|stale degrades expected return/safety margin|
|benchmark|portfolio benchmark symbol|benchmark return|benchmark source|freshness required|missing downgrades quarterly comparison|

## Deterministic Test Vectors Required

| Vector | Input | Expected output | Forbidden output |
| --- | --- | --- | --- |
|liquidity_20x|20-day average < plan amount * 20|liquidity risk|no market-style action|
|single_day_5pct|plan amount > same-day amount * 5%|liquidity risk|batch/limit/pause only|
|emotion_percentile|sentiment >90% or <10%|cooldown|no active trading suggestion|
|source_verification|<2 independent A/S sources|freeze/observe|no formal trade claim|
|expected_return_samples_lt5|sample_count < 5|qualitative only|no interval|
|expected_return_samples_lt20|sample_count < 20|no precise probability|sample warning required|

## UI/Action Matrix Minimum Columns

`requirement_id`, `ui_flow_id`, `browser_action`, `dom_assertion`, `expected_sqlite_changes`, `prohibited_sqlite_changes`, `audit_event`, `readback_page`, `mobile_result`, `failure_state_result`, `screenshot_path`, `status`

## Action-To-SQLite-To-Readback Matrix

| Action | Expected changed tables | Prohibited changed tables | Required audit action | Readback pages | Status |
| --- | --- | --- | --- | --- | --- |
|add_holding|positions, portfolio_snapshots, position_snapshots, audit_events|decision_records, operation_confirmations, position_transactions, rule_versions|record_local_fact|positions, dashboard, audit|scoped_p72_p75_non510300|
|edit_holding|positions, portfolio_snapshots, position_snapshots, local_account_corrections, audit_events|decision_records, operation_confirmations, rule_versions|correct_local_fact|positions, audit|scoped_p72|
|batch_import_holdings|local_account_import_batches, positions, portfolio_snapshots, position_snapshots, audit_events|decision_records, operation_confirmations, rule_versions|import_local_facts|positions, audit|scoped_p72|
|correct_local_fact|local_account_corrections, positions, portfolio_snapshots, position_snapshots, audit_events|decision_records, operation_confirmations, rule_versions|correct_local_fact|positions, audit|scoped_p72|
|manual_confirmation|operation_confirmations, position_transactions, portfolio_snapshots, position_snapshots, audit_events, decision_records.confirmation_status|broker_orders, external_pushes, rule_versions|confirm_manual_offline_action|decision_detail, decision_loop, audit|scoped_p72|
|mark_error_case|error_cases, operation_confirmations, audit_events|position_transactions, broker_orders|mark_error_case|decision_detail, review/error_cases, audit|scoped_p75_sop_failure_real_ui|
|generate_rule_proposal|rule_proposals, audit_events|rule_versions unless gatekeeper+final confirm|generate_rule_proposal|rules, audit|matrix_required_not_full_ui_executed|
|gatekeeper_review|gatekeeper_audits, rule_proposals.status, audit_events|rule_versions unless final user confirm|audit_rule_change|rules, audit|scoped_p75_sop_failure_real_ui|
|daily_report|daily_discipline_reports, audit_events, decision_records when workflow runs|operation_confirmations unless user confirms offline action|run_daily_discipline|daily_reports, workbench, audit|scoped_p72|
|monthly_review|review artifacts / error statistics read model, audit_events when generated|broker_orders, external_pushes, operation_confirmations|generate_monthly_review|review, audit|matrix_required_not_full_ui_executed|
|quarterly_review|benchmark comparison / rule_effect_tracking / audit_events when generated|broker_orders, external_pushes, auto rule_versions|generate_quarterly_review|review, rule_effect, audit|matrix_required_not_full_ui_executed|

## SOP A-F Real-Use Coverage Matrix

| SOP | Requirement rows | Trigger | Rule priority | Required data prerequisites | LLM role | User confirmation behavior | Readback pages | Current evidence | Gap | Status | Release impact |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
|SOP-A holding_drop|REQ-08-001..REQ-08-005|single-day or short-term drop >5%|buy-logic-break check before valuation/buy-more; fear cooldown overrides active action|position snapshot, buy logic, PE percentile, sentiment proxy, formal evidence, historical analog samples|explain rule result and evidence gaps only; cannot override rule verdict or fabricate formal evidence|only offline/manual confirmation may update account state|risk alerts, audit, mobile readback|`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-A card, data prerequisites, LLM role, lifecycle UI action, audit readback, and no broker/order/push state.|remaining branch-depth variants require broader live-provider and arbitrary-symbol coverage|scoped_p75_real_ui_pass|Closes P75 SOP-A real-browser claim within accepted-local scope.|
|SOP-B holding_rise|REQ-08-006..REQ-08-010|short-term rise >15% or floating profit >20%|PE/PB high valuation and staged take-profit before allocation destination; no repeated 20% stage|position P&L, price/NAV, PE/PB percentile, prior take-profit stage, core/satellite/cash classification|explain staged discipline and remaining-position trailing stop; cannot turn suggestion into order|user records offline sell/rebalance only after acting outside system|risk alerts, audit, mobile readback|`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-B card, data prerequisites, LLM role, lifecycle UI action, audit readback, and no broker/order/push state.|remaining 20%/30%/trailing-stop live variants require broader data coverage|scoped_p75_real_ui_pass|Closes P75 SOP-B real-browser claim within accepted-local scope.|
|SOP-C hot_topic_chasing|REQ-08-011..REQ-08-015|user asks whether to buy a hot-topic asset|circle-of-competence refusal overrides positive signals; high valuation/position limit blocks chasing|symbol profile, capability-circle tag, PE/PB percentile, current allocation, formal evidence, readiness status|translate discipline constraints and ask for missing evidence; cannot recommend capability-outside purchase|no account mutation unless user later records explicit local fact|risk alerts, audit, mobile readback|`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-C card, data prerequisites, LLM role, lifecycle UI action, audit readback, and no broker/order/push state.|capability-outside refusal remains separately covered by readiness/failure-state UI, not every hot-topic provider branch|scoped_p75_real_ui_pass|Closes P75 SOP-C real-browser claim within accepted-local scope.|
|SOP-D panic_sell|REQ-08-016..REQ-08-020|user expresses fear and wants to clear position|cooldown first; objective data and historical analog before any rational reminder|user text risk tag, sentiment percentile, PE percentile, holding valuation, historical analog sample/provenance|calmly summarize facts and uncertainty; cannot confirm trade or bypass cooldown|final decision remains user-owned and confirmation records only offline action|risk alerts, audit, mobile readback|`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-D card, data prerequisites, LLM role, lifecycle UI action, audit readback, and no broker/order/push state.|historical analog live-provider branch remains scoped|scoped_p75_real_ui_pass|Closes P75 SOP-D real-browser claim within accepted-local scope.|
|SOP-E macro_gray_rhino|REQ-08-021..REQ-08-023|known macro risk develops into material threat|buy-logic-break and formal evidence before sell suggestion; volatility disturbance reduces rather than adds exposure|formal evidence, source verification, volatility vs historical average, expected-return scenario probabilities and provenance|reassess scenario assumptions with provenance; cannot invent probability changes without samples|user reviews proposal; any rule/threshold change goes through gatekeeper|risk alerts, rules, audit, mobile readback|`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-E card, data prerequisites, LLM role, lifecycle UI action, gatekeeper readback, audit readback, and no broker/order/push state.|live two-source macro evidence variants remain scoped|scoped_p75_real_ui_pass|Closes P75 SOP-E real-browser claim within accepted-local scope.|
|SOP-F black_swan|REQ-08-024..REQ-08-026|sudden black-swan event|freeze active actions for 24h; require two A-level source impact assessments before reassessment|event timestamp, A-level source verification, freeze expiry, affected holdings, audit trail|state freeze and evidence insufficiency; cannot force reassessment before freeze/source gates pass|no confirmation during freeze except recording external facts after user action; never broker/order/push state|risk alerts, settings, audit, mobile readback|`p75-sop-failure-real-ui-acceptance.sh` verifies SOP-F card, data prerequisites, LLM role, lifecycle UI action, stale/degraded source UI, audit readback, and no broker/order/push state.|24h clock transition and live A-source variants remain scoped|scoped_p75_real_ui_pass|Closes P75 SOP-F real-browser claim within accepted-local scope.|

## Critical UI Flow Matrix

| Requirement | UI flow | Browser action | DOM assertion | Expected SQLite changes | Prohibited SQLite changes | Audit event | Readback page | Mobile result | Failure-state result | Screenshot path | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
|REQ-04-002|onboarding|open first-use positions/dashboard flow|empty state or existing snapshot is explicit|none until user saves|no broker/order/push tables|none or local_fact when saved|positions/dashboard|not rerun in P75|validation errors required|P72 screenshots|scoped_p72|
|REQ-11-001|add_fund|add 159915 holding in browser|holding row and data-quality symbol visible|positions, portfolio_snapshots, position_snapshots|decision_records, operation_confirmations|record_local_fact|positions,data-quality|not rerun mobile in P75|unsupported symbol handled by readiness blocked tests|2026-06-20-p75-non-510300|pass_scoped_non510300|
|REQ-05-021|data_readiness|open /data-quality?symbol=159915|readiness cards show 159915/399006/request_id|none|all mutation tables|none|data-quality|not rerun mobile in P75|degraded categories covered by DataQualityPage tests|2026-06-20-p75-non-510300|pass_scoped_non510300|
|REQ-08-001|consultation|submit real LLM-backed consultation|decision link/detail available|decision_records,evidence_refs,audit_events|broker_orders,external_pushes,rule_versions|generate_decision|decision_detail,decision_loop,audit|not rerun mobile in P75|model unavailable safe degradation covered by workflow tests|2026-06-20-p75-non-510300|pass_scoped_non510300|
|REQ-08-002|decision_detail|open decision detail|evidence, analyst reports, arbitration, expected return displayed|none|mutation tables|none|decision_detail|not rerun mobile in P75|insufficient evidence/detail degradation required|P72/P75 screenshots|scoped|
|REQ-35-001|alerts|open risk alerts and update SOP lifecycle|SOP A-F cards show trigger, data prerequisites, LLM role, safety copy, and updated statuses|risk_alerts.sop_status,audit_events|broker/order/push tables, operation_confirmations, position_transactions|risk_alert|risk_alerts,audit|390px pass in P75|stale/degraded source checked in settings/data-quality|2026-06-20-p75-sop-failure|scoped_p75_real_ui_pass|
|REQ-11-002|offline_confirmation|confirm manual offline action|confirmation status/readback visible|operation_confirmations,position_transactions,portfolio_snapshots,audit_events|broker_orders,external_pushes,auto confirmations|confirm_manual_offline_action|decision_detail,decision_loop,audit|not rerun mobile in P75|stale confirmation rejection tested|P72 screenshots|scoped_p72|
|REQ-12-001|error_marking|mark error case in decision detail|confirmation status, review and audit readback visible|operation_confirmations,error_cases,audit_events|position_transactions,broker_orders|mark_error|decision_detail,review,audit|390px adjacent decision/audit pass in P75|validation error checked on consultation|2026-06-20-p75-sop-failure|scoped_p75_real_ui_pass|
|REQ-13-001|rule_proposal|generate proposal|proposal status visible|rule_proposals,audit_events|rule_versions before final confirm|generate_rule_proposal|rules,audit|not executed|insufficient sample and conflict states required|none|matrix_required_not_full_ui_executed|
|REQ-13-002|gatekeeper_pass_deny_review|send proposal to gatekeeper and review pass/deny/user-review states|approved/rejected/user-review visible; sent proposal stops at pending_final_confirm|gatekeeper_audits,rule_proposals.status,audit_events|auto rule_versions without final confirm|audit_rule_change|rules,audit|390px rules pass in P75|deny and user-review states checked|2026-06-20-p75-sop-failure|scoped_p75_real_ui_pass|
|REQ-10-001|monthly_review|open monthly review|P&L/discipline/emotion/error stats visible|none unless report generated|broker/order/push tables|generate_monthly_review when generated|review,audit|not executed|missing data degradation required|none|matrix_required_not_full_ui_executed|
|REQ-10-002|quarterly_review|open quarterly review|benchmark/rule-effect/evolution summary visible|none unless report generated|auto rule_versions,broker/order/push|generate_quarterly_review when generated|review,rule_effect,audit|not executed|missing benchmark degradation required|none|matrix_required_not_full_ui_executed|
|REQ-18-001|audit_trail|open audit page|recent actions/error codes visible|none|all mutation tables|none|audit|not rerun mobile in P75|classified failures required|P72/P75 screenshots|scoped|
|REQ-18-002|settings_safety|open settings/safety pages|no trading/broker/auto affordance|none|broker/order/push tables|none|settings,audit|not rerun mobile in P75|forbidden affordance scan passed|P72 screenshots|scoped|

## UX Misunderstanding Checklist

| Risk | Required UX treatment | P75 status |
| --- | --- | --- |
|trading_boundary|UI must say suggestions are analysis/discipline records, not orders or broker actions|scoped_scan_pass|
|state_language|ready/degraded/blocked must explain data source and affected features|scoped_data_quality_pass|
|next_action|manual next steps must not imply automatic execution|scoped_p72|
|evidence_insufficient|insufficient evidence must freeze/degrade consultation and decision detail|deterministic_workflow_pass|
|offline_execution|confirmation means user manually executed outside system|scoped_p72|
|in_system_confirmation|confirmation records local facts only and must not create broker/order/push state|deterministic_service_pass|
|account_state_mutation|position/portfolio writes require explicit local action and audit|scoped_p72|

## Continuous Non-510300 UI Flow Evidence

`add fund` -> `data readiness` -> `consultation or alerts` -> `SQLite verification` -> `derived page readback` -> `same user symbol/tracked index correlation key`

- Verification: `bash scripts/p75-non-510300-real-ui-journey.sh` passed on 2026-06-20.
- Browser proof: `web/e2e/p75-non-510300-real-ui-journey.spec.ts` performed real UI actions for `/positions`, `/data-quality?symbol=159915`, `/consultation`, `/decisions/{decision_id}`, and `/decision-loop`.
- Data-quality UI proof: `/data-quality?symbol=159915` displayed `已准备`, `创业板ETF · ETF · 跟踪 399006`, `跟踪指数 · 已准备`, `估值分位 · 已准备`, `request：req_...`, and `标的：399006` after the P75 UI symbol filter fix.
- SQLite/request proof: `scripts/p75_non_510300_sqlite_check.py` verified market request `symbol=159915`, CNInfo request `stock=159915,...` with `seDate=2026-06-01~2026-06-30`, SZSE request `keyword=159915`, position facts, market source-health request correlation, satisfied formal evidence, indexed RAG chunks, completed LLM-backed decision with 3 analyst reports, consultation audit chain, and no forbidden trading/external-push tables.
- Screenshot/artifact path: `docs/release/ui-audit-assets/2026-06-20-p75-non-510300/`.
- Release boundary: the non-510300 journey closes 8.6a; the P75 SOP/failure-state runner closes the SOP, mobile, mark-error, gatekeeper, and failure-state rows listed above within accepted-local real-browser scope. Full action-matrix rows that require live-provider or arbitrary-symbol breadth remain scoped/partial.

## Repeatability Treatment

P75 records inherited P71/P73/P74 evidence as scoped and downgrades affected claims rather than claiming fresh repeatability for those milestones. P72 was rerun after P75 hardening and is cited above as repeated real UI scenario evidence. This satisfies P75 10.7a for a scoped conclusion, but it does not authorize a full original-requirement pass.
