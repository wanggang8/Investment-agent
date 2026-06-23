# Frontend Experience Tests Specification

## Purpose
Document frontend interaction, display, and test requirements for Investment Agent user-facing flows.
## Requirements
### Requirement: Frontend charts use API DTOs
系统 SHALL 在今日纪律、持仓、复盘页面提供图表展示，并确保图表数据只来自 API DTO。

#### Scenario: Dashboard charts render from DTO data
- **WHEN** 今日纪律、持仓或复盘页面收到 API DTO
- **THEN** 页面 MUST 展示仓位、风险、证据覆盖或复盘摘要图表
- **AND** 图表组件 MUST NOT 直接读取 SQLite、VecLite 或本地文件

#### Scenario: Chart data mapping is explainable
- **WHEN** 前端把 API DTO 转换为图表展示模型
- **THEN** 非显然转换逻辑 MUST 有中文注释说明

### Requirement: Frontend interaction states are explicit
系统 SHALL 为关键页面提供证据、决策链、审计时间线、运维状态和复盘摘要交互，并清楚展示空态、错误态、降级态和成功态。

#### Scenario: Evidence and audit interactions are available
- **WHEN** 用户查看证据、决策链或审计时间线
- **THEN** 前端 MUST 支持筛选或展开关键条目
- **AND** 展开内容 MUST 继续来自 API DTO

#### Scenario: Degraded and error states are visible
- **WHEN** 页面遇到信息不足、数据过期、LLM 降级、VecLite 不可用、复盘数据为空或错误响应
- **THEN** 前端 MUST 展示明确空态、错误态、降级态或成功态
- **AND** 用户 MUST 能区分信息不足、冻结观察、高危、降级、成功和普通错误

#### Scenario: Ops and review states are visible
- **WHEN** 运维状态、索引健康或复盘摘要 DTO 可用
- **THEN** 前端 MUST 展示对应状态、计数和安全说明
- **AND** 未知状态 MUST 使用安全显示，不得展示为成功

### Requirement: Frontend tests cover P8 behavior
系统 SHALL 建立前端测试，覆盖 API client、关键状态、用户确认、规则提案最终确认和禁止自动交易入口。

#### Scenario: API client and state rendering are tested
- **WHEN** 前端测试执行
- **THEN** 测试 MUST 覆盖 API client 的 `request_id`、`data`、`error` 处理
- **AND** 测试 MUST 覆盖信息不足、冻结观察、高危、降级和错误响应状态

#### Scenario: Confirmation and rule proposal flows are tested
- **WHEN** 前端测试执行
- **THEN** 测试 MUST 覆盖用户确认流程只记录线下动作
- **AND** 测试 MUST 覆盖规则提案 `pending_final_confirm` 可见且不会自动应用规则

#### Scenario: No automatic trading entry exists
- **WHEN** 前端测试检查核心页面和确认流程
- **THEN** 页面 MUST NOT 出现自动交易、一键交易或代下单入口

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

### Requirement: Review page displays periodic summaries and tracking
The frontend SHALL display periodic review summaries, rule suggestions, ops status, and tracking entrypoints using API/service DTOs rather than direct local storage access.

#### Scenario: Periodic summary is visible
- **WHEN** monthly or quarterly review data is available
- **THEN** the review page shows the period summary, relevant audit status, supporting counts, and degradation indicators.

#### Scenario: Rule suggestions are visible but not applied automatically
- **WHEN** a review produces rule suggestions
- **THEN** the frontend displays the suggestions as review output or rule proposal entrypoints and does not present automatic rule application behavior.

#### Scenario: Tracking entrypoint is available
- **WHEN** a review summary references audit events, rule proposals, error cases, or decisions
- **THEN** the frontend provides a visible path to inspect the related tracking records.

#### Scenario: Ops status is visible from review surface
- **WHEN** review or ops summary data contains data source, index, or degradation status
- **THEN** the frontend displays the status without reading local files, SQLite, or VecLite directly.

### Requirement: P39 Browser UX Stability Checks
The frontend SHALL add browser-level stability checks for key local journeys, including console error capture, unhandled rejection capture, narrow viewport smoke, and basic accessibility-oriented assertions.

#### Scenario: Key pages have no unexpected browser errors
- **WHEN** Playwright opens dashboard, portfolio, evidence, decision detail, audit/review, rules, risk alerts, daily discipline report, and settings pages in the P39 fixture
- **THEN** the test SHALL fail on unexpected console errors or unhandled page errors
- **AND** allowed diagnostic logs, if any, SHALL be explicitly scoped and documented in the test fixture

#### Scenario: Narrow viewport keeps primary controls usable
- **WHEN** key pages are rendered under a narrow mobile-like viewport
- **THEN** primary navigation, status labels, form controls, and action buttons SHALL remain visible and non-overlapping
- **AND** critical labels SHALL NOT be hidden in a way that changes the safety meaning of a page

#### Scenario: Basic accessibility expectations are covered
- **WHEN** forms, navigation, and interactive controls are rendered in the P39 browser journey
- **THEN** controls SHALL have accessible names, form inputs SHALL have labels or equivalent accessible descriptions, and navigation landmarks or equivalent page structure SHALL be discoverable
- **AND** these checks SHALL rely on browser-visible semantics rather than direct local file or SQLite reads

#### Scenario: Vitest and Playwright remain separated
- **WHEN** frontend verification runs
- **THEN** Vitest SHALL continue to cover component and mapper behavior
- **AND** Playwright SHALL cover browser journeys with fixed-ID local fixture data
- **AND** the two suites SHALL avoid collecting each other's files or sharing mutable persistent test state

### Requirement: P42 workbench frontend tests

The frontend SHALL test the P42 user decision workbench across successful, empty, degraded, error, and narrow-viewport paths.

#### Scenario: Workbench component states are tested

- **WHEN** frontend unit tests run
- **THEN** they SHALL cover workbench panels for successful DTOs, empty local facts, degraded source/LLM/RAG status, API errors, and safe Chinese status text
- **AND** tests SHALL assert that automatic trading, one-click order placement, external push, automatic confirmation, and automatic rule application copy is absent

#### Scenario: Workbench browser smoke is tested

- **WHEN** Playwright smoke runs
- **THEN** it SHALL open the workbench route, verify primary panels and navigation entrypoints, check narrow viewport usability, and scan for forbidden automatic-action copy
- **AND** it SHALL use fixed local fixture data rather than private persistent data

### Requirement: P43 data quality observability frontend tests

The frontend SHALL test the P43 data quality observability surface across successful, empty, degraded, error, unknown, sanitized, and narrow-viewport paths.

#### Scenario: Data quality component states are tested

- **WHEN** frontend unit tests run
- **THEN** they SHALL cover data quality panels for successful DTOs, empty local facts, source_unavailable, parse_error, stale, missing, unknown, LLM/RAG/VecLite degraded states, API errors, and safe Chinese status text
- **AND** tests SHALL assert that secrets, full prompts, private local paths, SQL errors, automatic trading, one-click order placement, external push, automatic confirmation, and automatic rule application copy is absent.

#### Scenario: Data quality browser smoke is tested

- **WHEN** Playwright smoke runs
- **THEN** it SHALL open the data quality route, verify primary panels and navigation entrypoints, check narrow viewport usability, and scan for forbidden automatic-action or sensitive diagnostic copy
- **AND** it SHALL use fixed local fixture data rather than private persistent data.

### Requirement: Real LLM decision details tolerate nullable frontend DTO fields

The frontend SHALL render decision detail pages safely when real LLM-backed decisions contain nullable, missing, or empty list fields in final verdict and trace DTOs.

#### Scenario: Optional actions are null

- **WHEN** a decision detail DTO contains `final_verdict.optional_actions` as `null`
- **THEN** the decision detail page MUST render without a page-level crash
- **AND** the optional-action section MUST show an empty or unavailable safe state instead of calling array-only methods on the value

#### Scenario: Prohibited actions are null or missing

- **WHEN** a decision detail DTO contains `final_verdict.prohibited_actions` as `null` or omits the field
- **THEN** the decision detail page MUST render without a page-level crash
- **AND** the frontend MUST NOT display a false success or false permission to trade

#### Scenario: Real LLM-like fixture is covered by tests

- **WHEN** frontend tests run
- **THEN** they MUST include a real LLM-like decision fixture with nullable verdict list fields
- **AND** they MUST assert that the decision trace, safety boundary, and final verdict remain visible

### Requirement: Productized task-based frontend shell

The frontend SHALL present Investment Agent as a task-based local investment discipline product rather than a flat route list.

#### Scenario: Navigation is grouped by user task

- **WHEN** the app shell renders on desktop
- **THEN** navigation MUST group routes by user task such as today, decision, portfolio, evidence, governance, and system operations
- **AND** all existing primary routes from the P55 route matrix MUST remain reachable

#### Scenario: Mobile navigation does not consume the reading viewport

- **WHEN** the app shell renders on a 390px wide viewport
- **THEN** navigation MUST avoid a permanently visible wide sidebar that compresses main content
- **AND** primary navigation controls MUST remain keyboard and pointer accessible

### Requirement: Productized operational UI system

The frontend SHALL use consistent operational UI primitives for buttons, forms, status labels, cards, and tabular or key-value data displays.

#### Scenario: Forms use consistent field structure

- **WHEN** users operate consultation, positions, local install, local knowledge, or settings forms
- **THEN** labels, hints, inputs, error/success states, and primary/secondary actions MUST use consistent styling and spacing
- **AND** browser-default controls MUST NOT be the dominant visual treatment for critical workflows

#### Scenario: Status and safety states remain visible

- **WHEN** pages show high risk, frozen watch, information insufficient, degraded, unknown, or success states
- **THEN** the UI MUST use consistent semantic styling and Chinese labels
- **AND** unknown or degraded states MUST NOT be styled as ordinary success

### Requirement: Mobile reflow for core acceptance pages

The frontend SHALL reflow core acceptance pages so the page itself does not horizontally overflow on mobile-sized viewports, except for explicitly scoped two-dimensional data containers.

#### Scenario: Positions page reflows on mobile

- **WHEN** `/positions` renders at 390px viewport width
- **THEN** account and holding forms MUST remain visible without page-level horizontal overflow
- **AND** holdings data MUST be readable through stacked cards, key-value rows, or a clearly scoped local table scroller

#### Scenario: Data quality page reflows on mobile

- **WHEN** `/data-quality` renders at 390px viewport width
- **THEN** source health, evidence/RAG, LLM quality, and affected workflow sections MUST remain visible without page-level horizontal overflow
- **AND** long source identifiers, status tokens, and diagnostic text MUST wrap, truncate safely, or scroll only within a scoped local container

### Requirement: Product design evidence is linked to UI acceptance fixes

The P56 implementation SHALL document how Product Design skill guidance and product design research were applied to acceptance-blocking UI fixes.

#### Scenario: Design rationale is traceable

- **WHEN** P56 is reviewed
- **THEN** the change materials MUST include the product brief, design principles, research inputs, and page-level UI plan
- **AND** the implementation report MUST map material UI changes back to those inputs

#### Scenario: Subagent review covers design and safety

- **WHEN** P56 plan review, execution review, or pre-commit review runs
- **THEN** the subagent review MUST check Product Design skill usage, research-backed rationale, mobile usability, real UI acceptance evidence, and prohibited automatic-action boundaries

### Requirement: Product experience polish roadmap governs post-P56 UI work

The frontend product experience SHALL be polished through a staged roadmap before a new final release-ready refresh is claimed.

#### Scenario: Product north star is explicit

- **WHEN** planning post-P56 frontend work
- **THEN** the product MUST be treated as a local investment discipline workbench
- **AND** it MUST NOT be treated as a broker trading terminal, AI chat demo, marketing landing page, or engineering debug console
- **AND** the core daily questions MUST be: can I act today, why, what manual action is needed, and whether data and rules are trustworthy

#### Scenario: Product polish is staged

- **WHEN** post-P56 UI/product improvements are planned
- **THEN** the work MUST be split into independent OpenSpec changes for daily workbench, decision explainability, portfolio/risk/data quality, governance/ops productization, design system/accessibility, and final real UI regression
- **AND** governance/ops productization MUST explicitly include rules, audit, notifications, daily reports, daily auto run, local install, local knowledge, and settings surfaces
- **AND** each stage MUST define scope, out-of-scope safety boundaries, Product Design evidence, browser validation, and subagent review gates

#### Scenario: Release refresh is sequenced after product polish

- **WHEN** P57 product experience roadmap is accepted
- **THEN** release-readiness refresh MUST be deferred until the product polish stages have either completed or been explicitly waived
- **AND** documentation MUST NOT claim that all product design, UI design, or frontend issues are fully fixed before the corresponding stages pass validation

#### Scenario: Safety boundaries remain visible in polished UI

- **WHEN** any polished UI adds or changes a control, page, CTA, state, workflow, or report
- **THEN** it MUST NOT add or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source
- **AND** high risk, unknown, degraded, stale, missing, information-insufficient, and blocked states MUST NOT be styled or worded as ordinary success

#### Scenario: Real UI validation remains required

- **WHEN** a product polish stage changes frontend behavior, layout, information architecture, component primitives, or user-facing copy
- **THEN** validation MUST include frontend unit or component tests, frontend build, relevant backend tests when touched, browser-operated local UI verification, desktop and mobile screenshots, mobile reflow checks, safety copy scans, sensitive information scans, and subagent reviews
- **AND** real LLM validation MUST be included for stages that alter consultation, decision detail, evidence explanation, LLM quality display, or decision-loop surfaces

### Requirement: Daily workbench presents an at-a-glance investment discipline cockpit

The frontend SHALL make the dashboard and workbench answer the daily investment discipline questions before secondary details.

#### Scenario: Dashboard first screen shows daily decision state

- **WHEN** the dashboard route `/` renders with dashboard and daily discipline data
- **THEN** the first screen MUST show the current verdict or safe unavailable state, status tone, data trust summary, last update context, prohibited actions, optional manual actions, and next manual actions
- **AND** the first screen MUST provide local navigation to decision detail or daily report when such links are available

#### Scenario: Workbench first screen shows task queue

- **WHEN** the workbench route `/workbench` renders
- **THEN** it MUST use the same daily state model as the dashboard
- **AND** it MUST present a prioritized manual action queue before secondary portfolio, risk, rule, review, or consultation sections
- **AND** each action MUST be a local navigation or manual review action, not an execution action

#### Scenario: Daily workbench handles degraded and insufficient states safely

- **WHEN** dashboard, daily report, portfolio, risk, rule, or review data is unavailable, degraded, stale, high risk, unknown, or insufficient
- **THEN** dashboard and workbench MUST show safe Chinese status text and a clear next manual step
- **AND** they MUST NOT style or describe degraded, unknown, high risk, stale, missing, or insufficient states as ordinary success

#### Scenario: Daily workbench remains mobile readable

- **WHEN** `/` or `/workbench` renders at 390px viewport width
- **THEN** the daily status, manual action queue, and primary task links MUST remain visible without page-level horizontal overflow
- **AND** screenshots MUST be captured for desktop and mobile validation

#### Scenario: Daily workbench preserves safety boundaries

- **WHEN** dashboard or workbench UI adds or changes a CTA, status, summary, or task link
- **THEN** it MUST NOT add or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source

### Requirement: Decision explanation surfaces present a readable decision story

The frontend SHALL connect consultation, decision detail, evidence, and decision-loop pages into a readable decision story while preserving existing backend contracts and safety boundaries.

#### Scenario: Consultation route makes the generated decision path explicit

- **WHEN** the user opens `/consultation`
- **THEN** the page MUST clearly show the consultation inputs, assumptions, generation state, and safety boundary before the generated result
- **AND** after a successful consultation it MUST show the generated verdict or safe unavailable state, key reasons, prohibited actions, optional manual actions, data trust context, and local navigation to the generated decision detail when a decision id exists
- **AND** the page MUST NOT present the generated result as an automatic confirmation, trade, rule application, external push, or return promise

#### Scenario: Decision detail first screen explains the verdict before technical trace

- **WHEN** the user opens `/decisions/:decisionId`
- **THEN** the first screen MUST show the final verdict or safe unavailable state, generated context, prohibited actions, optional manual actions, key reasons, data trust summary, and safety boundary before long technical details
- **AND** Evidence, LLM, rules, expected return, arbitration, audit, and confirmation details MUST be grouped into readable layers
- **AND** long traces MUST not obscure the first-screen verdict, safety boundary, or explanation path

#### Scenario: Decision explanation handles nullable and degraded DTOs safely

- **WHEN** decision, verdict, evidence, analyst, expected-return, retrieval-quality, audit, or confirmation fields are null, missing, empty, degraded, unknown, failed, or insufficient
- **THEN** the frontend MUST render safe Chinese empty/degraded text without page-level crashes
- **AND** nullable or missing fields MUST NOT be described as permission to trade, auto-confirm, auto-apply rules, or treat the decision as successful without review

#### Scenario: Evidence page explains source trust and links back to decision explanation

- **WHEN** the user opens `/evidence`
- **THEN** the page MUST prioritize a source trust summary, source-level explanation, verification status, and local navigation back to decision explanation surfaces before raw evidence detail
- **AND** the evidence table MUST preserve filtering and expansion without exposing raw vendor payload, complete prompt, private path, key, or local database content

#### Scenario: Decision loop page reads as a read-only decision lifecycle

- **WHEN** the user opens `/decision-loop`
- **THEN** the page MUST show a read-only lifecycle from recommendation to confirmation, manual record, risk/review, and audit links
- **AND** it MUST clearly identify missing links or open gaps as manual follow-up items
- **AND** it MUST NOT provide controls that create confirmations, trades, risk lifecycle changes, rule applications, notifications, or settings changes

#### Scenario: Decision explanation remains mobile readable and safe

- **WHEN** `/consultation`, `/decisions/:decisionId`, `/evidence`, or `/decision-loop` renders at 390px viewport width
- **THEN** primary verdict, safety boundary, trust context, key reasons, and local navigation MUST remain readable without page-level horizontal overflow
- **AND** screenshots or browser evidence MUST be captured for desktop and mobile validation

### Requirement: Portfolio, risk, and data-quality pages present operational maintenance experiences

The frontend SHALL turn positions, risk alerts, and data quality into readable maintenance and disposition surfaces while preserving existing backend contracts and local-only safety boundaries.

#### Scenario: Positions page explains the portfolio maintenance state before local fact forms

- **WHEN** the user opens `/positions`
- **THEN** the first screen MUST show portfolio status, snapshot context, total assets, cash ratio, position count, high-risk ratio, current maintenance stage, and next manual actions before detailed forms or tables
- **AND** it MUST distinguish initialization, calibration, holding edit, offline transaction recording, batch import validation/confirmation, and correction audit paths
- **AND** every write action MUST be described as a local fact or audit record, not broker sync, order placement, automatic portfolio management, or return optimization

#### Scenario: Positions page preserves explicit local-only write boundaries

- **WHEN** the user submits calibration, holding edit, holding removal, offline transaction, batch import, or correction audit actions
- **THEN** the frontend MUST call only the existing portfolio service methods and display safe success or error messages
- **AND** disabled or unavailable actions MUST explain the missing local prerequisite instead of implying automatic recovery
- **AND** the page MUST NOT expose broker login, automatic trade, one-click order, delegated order, external push, automatic confirmation, automatic rule application, automatic repair, database overwrite, or return promise controls

#### Scenario: Risk alerts render as a disposition queue

- **WHEN** the user opens `/risk-alerts`
- **THEN** the page MUST show risk disposition summary, severity, affected symbols, and queues for pending review, in progress, needs review, and recorded risks
- **AND** each alert MUST show risk type, severity, SOP status, trigger summary, prohibited actions, suggested manual actions, related local links, updated time, and safety note
- **AND** the queue mapping MUST treat `triggered` as pending review, `active` and `observing` as in progress, `escalated` as needs review, and `resolved` or `archived` as recorded

#### Scenario: Risk SOP actions remain explicit local lifecycle records

- **WHEN** an unresolved risk alert is shown
- **THEN** lifecycle controls MAY allow continue observing, escalate for review, or resolve locally through the existing risk alert lifecycle service
- **AND** resolved or archived risks MUST NOT show lifecycle controls
- **AND** SOP controls MUST NOT imply automatic trading, external notification, automatic confirmation, rule application, or portfolio mutation

#### Scenario: Data quality page explains quality signals and affected workflows

- **WHEN** the user opens `/data-quality`
- **THEN** the first screen MUST show an overall quality state, source health signal, evidence/RAG signal, LLM signal, affected workflow signal, and next local inspection actions
- **AND** source health, evidence verification, VecLite, DeepSeek, review degradation, missing evidence, and affected decision/workflow details MUST remain visible in readable layers
- **AND** degraded, stale, missing, parse_error, unavailable, failed, unknown, or insufficient states MUST not be displayed as normal success

#### Scenario: Data quality diagnostics remain read-only and sanitized

- **WHEN** data quality APIs return source health failures, evidence summaries, review explanations, system paths, or unexpected diagnostic values
- **THEN** the frontend MUST render safe Chinese summaries and local navigation without exposing API keys, private paths, SQL, complete prompts, raw vendor payloads, local database paths, or raw stack traces
- **AND** the page MUST NOT offer automatic repair, automatic source refresh, automatic confirmation, rule application, external push, database overwrite, or trading controls

#### Scenario: Portfolio, risk, and data-quality experiences remain mobile readable

- **WHEN** `/positions`, `/risk-alerts`, or `/data-quality` renders at 390px viewport width
- **THEN** primary status, next manual actions, local safety boundary, form controls, queue cards, quality signals, and navigation MUST remain readable without page-level horizontal overflow
- **AND** screenshots or browser evidence MUST be captured for desktop and mobile validation

### Requirement: Governance and ops pages present productized workbench experiences

The frontend SHALL turn rules, audit, notifications, daily reports, daily auto run, local install, local knowledge, and settings into readable governance and operations surfaces while preserving existing backend contracts and local-only safety boundaries.

#### Scenario: Rules page presents governance status before raw rule details

- **WHEN** the user opens `/rules`
- **THEN** the first screen MUST show current rule version, proposal counts, pending user confirmation, pending final confirmation, gatekeeper or validation risk, and next manual governance actions
- **AND** rule proposals MUST show reason, sample count, overfit risk, validation status, guardrail decision, audit summary, related local records, and explicit manual confirmation boundaries
- **AND** raw rule JSON or threshold details MUST NOT dominate the first screen when structured summaries are available
- **AND** the page MUST NOT imply automatic rule application, automatic confirmation, broker connectivity, external push, or trading

#### Scenario: Audit page presents an operational timeline with summary context

- **WHEN** the user opens `/audit`
- **THEN** the first screen MUST show audit event count, recent activity, important event categories, and next local inspection actions before the detailed timeline
- **AND** audit events MUST remain traceable to existing API DTO fields without reading SQLite, local files, VecLite, or raw logs
- **AND** empty, degraded, unknown, and error states MUST be visible and safe

#### Scenario: Notifications page behaves as a local inbox

- **WHEN** the user opens `/notifications`
- **THEN** the first screen MUST show unread count, severity distribution, source categories, local processing status, and next manual actions
- **AND** mark-read controls MUST be described as local application state only
- **AND** the page MUST NOT promise SMS, email, webhook, third-party notification delivery, external push, automatic confirmation, or trading

#### Scenario: Daily reports and daily auto run explain discipline and local runtime health

- **WHEN** the user opens `/daily-discipline/reports` or `/daily-auto-run`
- **THEN** the first screen MUST show current discipline or runtime state, evidence or execution coverage, degraded or missing prerequisites, recent activity, and next manual checks
- **AND** daily auto run MUST distinguish disabled, scheduled, running, success, degraded, failed, and unknown states without styling degraded or unknown as normal success
- **AND** daily auto run diagnostics MUST guide manual recheck and MUST NOT promise automatic repair, automatic source refresh, automatic confirmation, automatic rule application, database overwrite, or trading

#### Scenario: Local install, local knowledge, and settings share safe configuration and diagnostics patterns

- **WHEN** the user opens `/local-install`, `/local-knowledge`, or `/settings`
- **THEN** the page MUST organize configuration, diagnostic status, previews, summaries, and next manual actions into readable sections
- **AND** sensitive values, API keys, private paths, SQL, raw stack traces, complete prompts, and raw vendor payloads MUST NOT be rendered
- **AND** existing local write actions such as knowledge import confirmation or market refresh MUST remain explicit local facts or local refreshes and MUST NOT imply trading, rule application, or external delivery

#### Scenario: Governance and ops experiences remain mobile readable

- **WHEN** P61 pages render at 390px viewport width
- **THEN** primary status, next manual actions, safety boundaries, forms, inbox cards, timelines, diagnostic summaries, and navigation MUST remain readable without page-level horizontal overflow
- **AND** two-dimensional data MAY scroll only inside clearly scoped local containers
- **AND** screenshots or browser evidence MUST be captured for desktop and mobile validation

### Requirement: Design system primitives and accessibility gates are standardized

The frontend SHALL provide reusable UI primitives and browser-level accessibility/reflow gates for the productized local investment discipline workbench before the final full UI regression refresh.

#### Scenario: Shared UI primitives expose accessible semantics

- **WHEN** P62 introduces Button, Field, StatusBadge, PageHeader, SummaryCard, DetailSection, ResponsiveTable, EmptyState, and ErrorState primitives
- **THEN** each primitive MUST expose stable text or accessible names for its user-visible purpose
- **AND** form controls MUST be associated with labels, hints, and errors when those are present
- **AND** collapsible details MUST use button semantics and maintain `aria-expanded`
- **AND** table-like content MUST provide a caption or accessible name and preserve column meaning under mobile reflow

#### Scenario: Status tokens remain consistent and impossible to confuse with success

- **WHEN** pages render success, warning, danger, degraded, unknown, readonly, or blocked states
- **THEN** those states MUST use consistent frontend tone names
- **AND** every state MUST include readable text, not only color
- **AND** degraded, unknown, readonly, blocked, missing, stale, failed, and information-insufficient states MUST NOT be styled or worded as ordinary success

#### Scenario: Keyboard paths cover primary local workflows

- **WHEN** the local frontend is operated with a keyboard
- **THEN** the primary navigation, mobile menu, representative forms, collapsible detail sections, and critical local-only buttons MUST be reachable and operable
- **AND** focus indicators MUST remain visible on desktop, tablet-width, and mobile-width layouts
- **AND** disabled or working controls MUST expose their state through text, attributes, or equivalent accessible semantics

#### Scenario: Reflow and visual evidence cover desktop, tablet-width, and mobile layouts

- **WHEN** representative P58-P61 pages render at 390px, 768px, and 1280px viewport widths
- **THEN** primary status, next manual actions, forms, summaries, tables, timelines, diagnostics, empty states, error states, and navigation MUST remain readable without page-level horizontal overflow
- **AND** two-dimensional tables, JSON, logs, or diagnostic text MAY scroll only inside clearly scoped local containers
- **AND** screenshots or equivalent browser evidence MUST be captured for the three viewport classes

#### Scenario: Design system hardening preserves local-only safety boundaries

- **WHEN** P62 changes components, page layout, UI text, keyboard behavior, or validation evidence
- **THEN** it MUST NOT add or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source
- **AND** UI output MUST NOT render complete API keys, private paths, SQL, raw stack traces, complete prompts, local database paths, or raw vendor payloads

### Requirement: Full real UI regression refreshes the release status

After product experience polishing is complete, the project SHALL execute a full real UI regression and release status refresh before claiming a current release candidate state.

#### Scenario: Full UI regression covers all primary routes

- **WHEN** P63 executes browser acceptance against the local backend and frontend
- **THEN** Dashboard, Workbench, Consultation, Decision Detail, Evidence, Decision Loop, Positions, Data Quality, Risk Alerts, Risk Alert Detail, Rules, Audit, Notifications, Daily Reports, Daily Report Detail, Daily Auto Run, Review, Local Install, Local Knowledge, and Settings MUST be operated through the UI or covered by an equivalent browser assertion
- **AND** each route MUST record whether its primary status, key actions, empty/degraded/error states, and local-only safety boundaries are visible and usable
- **AND** console errors, page errors, and failed API responses MUST be recorded or explicitly ruled non-blocking with reasons

#### Scenario: Real LLM consultation is verified or classified

- **WHEN** P63 runs the consultation journey with a real LLM-backed configuration
- **THEN** the UI MUST submit a real consultation request and attempt to open the resulting decision detail
- **AND** the acceptance record MUST state whether LLM analysis was returned, parsed, quality-gated, and displayed without letting the model write the final rule verdict
- **AND** network, rate limit, authentication, model unavailable, parse, or quality failures MUST be classified and mapped to release impact instead of being treated as success

#### Scenario: Release candidate status is refreshed from current evidence

- **WHEN** P63 produces release materials
- **THEN** the release candidate MUST reference the current acceptance run and current code-under-test commit
- **AND** the status MUST be one of `release_ready`, `release_degraded`, or `blocked`
- **AND** degraded, skipped, retried, or waived gates MUST list their reason, artifact, and release impact
- **AND** the handoff MUST include Not Claimed boundaries for returns, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, future provider availability, login sources, paid sources, authorized sources, Level2 data, and high-frequency data

#### Scenario: Full UI regression preserves safety and redaction boundaries

- **WHEN** P63 scans UI text, release materials, logs, browser evidence, and committed assets
- **THEN** it MUST NOT find new user-facing broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source claims
- **AND** committed artifacts MUST NOT contain complete API keys, private paths, raw SQL dumps, raw stack traces, complete prompts, local database paths, raw vendor payloads, Playwright trace archives, temporary SQLite databases, or unredacted local logs
