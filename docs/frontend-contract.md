# Investment Agent 前端数据契约

> 文档版本：v1.0
> 最后更新：2026-06-18
> 适用范围：React + Vite + TypeScript 前端页面、业务组件、HTTP API DTO 映射。

## 1. 契约定位

本文档定义 UI 页面、业务组件和后端 API 字段之间的映射关系。前端展示必须使用人能理解的文案，不直接把内部技术字段作为页面主信息。

原则：

- API 字段是数据来源，UI 字段是展示结果。
- 页面组件只依赖 `web/src/services/` 或 `web/src/shared/api/` 返回的 DTO。
- `web/src/pages/` 只负责路由组合、布局、加载态和错误态，不承载大段业务逻辑。
- dashboard、decision、evidence、portfolio、rules、audit、settings、market、review 等业务 UI 应放在对应 `web/src/features/<feature>/`。
- 跨业务域复用的 API client、类型、工具和通用组件放在 `web/src/shared/`。
- 通用组件不直接调用 API。
- hash、内部 ID、向量 ID 默认折叠展示。
- 所有状态都必须有文字说明，不只依赖颜色。

## 2. 通用类型

### 2.1 DisplayMeta

| 字段 | 类型 | 来源 | UI 用法 |
| --- | --- | --- | --- |
| requestId | string | `request_id` | 审计详情中展示 |
| generatedAtText | string | `meta.generated_at` | “本次结果生成于 2026-05-22 09:30” |
| ruleVersionText | string | `meta.rule_version` | “当前规则版本 v3.0” |
| dataUpdatedAtText | string | `data.data_updated_at` | 顶部状态条展示 |

### 2.2 Severity

| API 值 | UI 文案 | 颜色倾向 | 图标语义 |
| --- | --- | --- | --- |
| normal | 正常 | 冷灰、低饱和蓝 | 信息 |
| warning | 观察 | 琥珀色 | 注意 |
| danger | 高危 | 深红棕 | 风险 |
| frozen_watch | 冻结观察 | 深蓝灰 | 暂停 |
| insufficient | 信息不足 | 中性灰 | 数据不足 |

前端错误状态只依赖稳定错误码和展示状态，不得解析后端底层错误文本。`DATA_REQUIRED` 映射为 `first_use`；`DATA_STALE`、`EVIDENCE_NOT_FOUND`、`VECTOR_INDEX_UNAVAILABLE`、`ANALYST_UNAVAILABLE`、`DECISION_RECORD_FAILED` 映射为 `insufficient_data`；`SOURCE_VERIFICATION_FAILED` 映射为 `frozen_watch`；`INTERNAL_ERROR` 展示通用失败状态且不得显示 SQL、文件路径或外部服务原始错误。

P4 API DTO 必须使用与本文件页面契约兼容的 JSON 字段名。dashboard、decision、evidence、portfolio、rule、audit、settings、market、review API 返回数据时，字段名必须匹配对应页面引用的 API 字段。

### 2.3 P5 DTO 与 API service

P5 前端必须为 dashboard、portfolio、decision、evidence、rule、audit、settings、market、review、notification 定义类型化 DTO 与 service。service 必须解析统一响应信封中的 `request_id`、`data`、`meta` 与 `error`，调用方只能获得类型化成功数据或稳定前端错误状态。

前端页面和组件需要应用数据时，必须使用 `web/src/services/` 或 `web/src/shared/api/`，不得直接访问 SQLite、VecLite 或本地文件。

## 3. 页面状态契约

| 页面状态 | 触发条件 | 主文案 | 允许动作 |
| --- | --- | --- | --- |
| first_use | API 返回 `DATA_REQUIRED` | 需要先录入账户和持仓，系统才能生成纪律报告。 | 录入持仓、配置能力圈 |
| normal | 工作流完成且无高危规则 | 今日未触发纪律红线。 | 查看证据、查看详情 |
| insufficient_data | API 返回 `DATA_STALE`、`EVIDENCE_NOT_FOUND`、`VECTOR_INDEX_UNAVAILABLE`、`ANALYST_UNAVAILABLE` 或 `DECISION_RECORD_FAILED` | 当前信息不足，系统暂停生成交易类建议。 | 查看缺失数据、刷新数据 |
| frozen_watch | API 返回 `SOURCE_VERIFICATION_FAILED` 或黑天鹅冻结 | 信息仍在核查，进入冻结观察。 | 查看证据、标记待观察 |
| high_risk | 触发逻辑破坏、规则缺失、记录失败等高危状态 | 已触发高危纪律状态，禁止新增买入。 | 查看规则、记录线下结果 |

## 4. 驾驶舱页契约

页面：`web/src/pages/DashboardPage.tsx`
API：`GET /api/v1/dashboard/today`
核心组件：`DisciplineStatusCard`、`RuleTriggerBadge`、`DecisionVerdictPanel`、`SourceVerificationPanel`、`ConfirmActionPanel`

### 4.1 DashboardViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| pageState | `data.dashboard_state` 或错误码 | 映射为页面状态契约 |
| disciplineStatusText | `data.discipline_status` | 展示为首屏主状态 |
| dataUpdatedAtText | `data.data_updated_at` | 格式化为“数据更新于...” |
| totalAssetsText | `data.portfolio_summary.total_assets` | 格式化为人民币金额 |
| cashRatioText | `data.portfolio_summary.cash_ratio` | 格式化为百分比 |
| highRiskRatioText | `data.portfolio_summary.high_risk_ratio` | 格式化为百分比 |
| pePercentileText | `data.market_summary.pe_percentile` | 展示为“PE 分位 63%” |
| triggeredRules | `data.triggered_rules` | 渲染为 `RuleTriggerBadge[]` |
| verdictText | `data.decision_summary.verdict` | 展示为裁决主文案 |
| evidenceSummaryText | `data.evidence_summary` | 可选。存在时展示为“引用 3 个独立信源，最高 A 级”；缺失时隐藏模块或展示“证据不足”状态 |

### 4.2 DisciplineStatusCard

| UI 字段 | API 字段 | 必填 |
| --- | --- | --- |
| statusName | `discipline_status` | 是 |
| statusDescription | 由 `dashboard_state` 派生 | 是 |
| riskLevel | `triggered_rules[].severity` 最大值 | 是 |
| updatedAtText | `data_updated_at` | 是 |
| primaryActionText | 由 `decision_summary.action_required` 派生 | 是 |

### 4.3 P58 DailyWorkbenchViewModel

Dashboard `/` 和 Workbench `/workbench` 必须复用同一个 null-safe daily workbench view model。该 view model 只能组合现有 Dashboard、每日纪律报告、组合、风险、规则和复盘 DTO，不新增 HTTP API 字段。

| UI 字段 | 来源 | 展示规则 |
| --- | --- | --- |
| statusLabel | dashboard state / report status / error state | 映射为成功、降级、数据不足、高风险或未知；不得把 degraded、unknown、missing、stale、高风险展示为普通成功 |
| statusTone | 派生 | `success` / `warning` / `danger` / `degraded` / `unknown`，所有状态必须有文字说明 |
| verdictText | `daily_report.final_verdict`、`daily_report.summary`、`decision_summary.verdict` | 作为首屏主要结论；缺失时展示安全不可用文案 |
| trustSummary | `daily_report.evidence` | 展示证据数、独立信源数；不足时展示证据不足 |
| riskSummary | `portfolio_summary`、`risk_alerts` | 展示组合风险或待复核风险数量 |
| updatedAtText | `daily_report.updated_at` 或 dashboard 更新时间 | 格式化为本地可读时间，缺失时展示暂无 |
| prohibitedActions | 决策、报告或风险预警 | 只解释禁止边界，不得成为执行 CTA |
| optionalActions | 决策、报告或风险预警 | 只展示人工动作或线下记录，不得暗示自动执行 |
| nextActions | 派生本地导航 | 仅允许跳转到 `/positions`、`/data-quality`、`/risk-alerts`、`/decisions/:id`、`/daily-discipline/reports/:id`、`/consultation` 等本地页面 |
| signals | 派生摘要 | 展示数据可信度、组合风险、风险处置、规则与复盘 |
| warnings | API error display state | 展示安全中文错误，不显示 SQL、堆栈、路径或外部服务原始错误 |

Dashboard 首屏顺序为：今日纪律状态、下一步人工动作、今日信号摘要、状态/报告提示、详细驾驶舱。Workbench 首屏顺序为：页面标题、今日纪律状态、下一步人工动作、今日信号摘要，再展示组合/风险、规则/复盘、主动咨询等次级区域。

P58 CTA 安全边界：

- 允许：查看、进入、维护本地账户与持仓、查看数据质量、处理风险预警、查看决策详情、查看每日纪律报告、发起主动咨询。
- 禁止：券商接口、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、覆盖真实库或收益承诺入口。

## 5. 每日纪律报告页契约

页面：`web/src/pages/DailyDisciplineReportsPage.tsx`、`web/src/pages/DailyDisciplineReportDetailPage.tsx`
API：`GET /api/v1/daily-discipline/reports/today`、`GET /api/v1/daily-discipline/reports`、`GET /api/v1/daily-discipline/reports/{report_id}`

### 5.1 DailyDisciplineReportViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| reportStatusText | `data.status` | 映射为未开始 / 运行中 / 已完成 / 降级 / 失败 / 数据不足 |
| localDateText | `data.local_date` | 展示报告对应本地日期 |
| summaryText | `data.summary` | 成功或降级时展示；缺失时展示缺前提提示 |
| missingCategories | `data.missing_categories` | 数据不足时展示缺少账户、持仓、行情、证据、规则或配置等类别 |
| finalVerdictText | `data.final_verdict` | 仅在有关联决策且字段存在时展示 |
| verdictStatusText | `data.verdict_status` | 转换为观察、暂停买入、信息不足等用户可读状态 |
| evidenceCoverageText | `data.evidence` | 展示证据数、独立信源数和高等级独立信源数 |
| trendSummary | `data.trend` | 展示最近成功、降级、失败、数据不足次数 |
| decisionLink | `data.decision_link` | 跳转到决策详情；为空时隐藏 |
| auditLink | `data.audit_link` | 跳转到审计页筛选；为空时隐藏 |
| notificationLink | `data.notification_link` | 跳转到通知页筛选；为空时隐藏 |
| autoRunLink | `data.auto_run_link` | 跳转到每日自动运行状态页；为空时隐藏 |
| riskAlerts | `data.risk_alerts` | 展示风险类型、严重程度、SOP 状态、触发摘要、禁止动作、建议人工动作和风险详情链接 |
| safetyNote | `data.safety_note` | 固定展示，不得隐藏 |

### 5.2 报告状态展示

| API status | UI 状态 | 展示要求 |
| --- | --- | --- |
| not_started | 未开始 | 展示初始化或运行提示，不生成正式建议 |
| running | 运行中 | 展示等待状态和自动运行入口 |
| success | 已完成 | 展示摘要、裁决、证据覆盖和关联材料 |
| degraded | 降级完成 | 展示可用摘要，同时突出降级原因和缺失项 |
| failed | 失败 | 展示 `failure_code` / `failure_reason`，不展示交易类建议 |
| insufficient_data | 数据不足 | 展示 `missing_categories`，不伪造证据、预期收益或交易指令 |

报告页面、历史列表和详情页均保持本地只读。页面不得出现自动下单、一键交易、确定性涨跌预测或收益承诺文案。

### 5.3 风险摘要展示

每日纪律报告详情和 Dashboard 今日纪律报告卡片可展示 `risk_alerts` 摘要。摘要展示 risk type、severity、SOP status、trigger summary、prohibited actions、suggested actions 和 `link`。摘要只作为本地追踪入口；不得展示生命周期动作、自动交易动作或规则自动应用动作。完整 SOP 生命周期动作只允许在风险预警中心页面中显式执行。

## 6. 决策详情页契约

页面：`web/src/pages/DecisionDetailPage.tsx`
API：`GET /api/v1/decisions/{decision_id}`

### 6.1 DecisionDetailViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| displayTitle | `decision_id` + `generated_at` | “2026-05-22 第 2 条建议” |
| generatedAtText | `generated_at` | “生成于...” |
| accountSnapshot | `account_snapshot` | 展示现金、仓位、持仓状态 |
| triggeredRules | `triggered_rules` | 按严重程度排序 |
| evidenceCards | `evidence_chain` | 渲染为 `EvidenceCard[]` |
| agentOpinions | `analyst_reports` | 渲染为 `AgentOpinionPanel[]` |
| expectedReturnScenarios | `expected_return_scenarios` | 渲染为情景收益卡；仅展示概率与样本信息 |
| arbitrationSteps | `arbitration_chain` | 展示裁决优先级命中过程 |
| finalVerdict | `final_verdict` | 展示最终裁决和禁止项 |
| confirmation | `user_confirmation` | 渲染确认区 |

### 6.2 AgentOpinionPanel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| agentName | `analyst_reports[].agent_name` | 价值分析师 / 趋势与风控官 |
| conclusionText | `analyst_reports[].conclusion` | 一段明确结论 |
| keyReasons | `analyst_reports[].key_reasons` | 最多展示 3 条，详情展开 |
| riskWarnings | `analyst_reports[].risk_warnings` | 优先展示 |
| confidenceText | `analyst_reports[].confidence` | 转成低 / 中 / 高 |
| evidenceRefs | `analyst_reports[].evidence_ids` | 链接到证据卡片 |
| promptVersionText | `analyst_reports[].prompt_version` | 可选；详情中展示 |
| modelText | `analyst_reports[].model` | 可选；只展示模型名，不展示 key |
| qualityStatusText | `analyst_reports[].quality_status` / `parse_status` | 可选；展示解析/质量状态摘要 |

### 6.3 ExpectedReturnScenarioPanel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| sampleCountText | `expected_return_scenarios.sample_count` | 少于 20 个时标记样本不足；少于 5 个时展示不可生成收益区间原因 |
| sampleWindowText | `expected_return_scenarios.sample_window` | 展示样本区间；缺失时展示“暂无”或后端 reason |
| screeningConditionText | `expected_return_scenarios.screening_condition` | 展示筛选条件；缺失时展示“暂无”或后端 reason |
| precisionStatusText | `expected_return_scenarios.precision_status` | available：可展示概率；insufficient：展示样本不足且不展示精确概率；unavailable：不展示收益区间 |
| scenarioCards | `expected_return_scenarios.scenarios` | available：展示上行、基准、下行情景、收益区间、概率和 `trigger`；insufficient：可展示区间和 `trigger` 但不得展示精确概率；unavailable：空状态；`trigger` 只解释区间边界，动态卖出状态由后端按 upside 下沿、base 上沿、downside 下沿等情景边界计算 |
| sellEvaluationPanel | `expected_return_scenarios.sell_evaluation` | 展示状态、触发因素、人工提示、建议动作和 non-trading disclaimer；只作为人工复核材料，不提供自动交易入口 |
| reassessmentTriggerPanel | `expected_return_scenarios.reassessment_trigger` | 展示复核原因、边界和当前值；缺字段时按空/不适用处理 |
| disclaimerText | `expected_return_scenarios.disclaimer` | 固定展示，不得隐藏；必须说明不构成收益承诺且最终裁决以规则链为准 |

## 7. 证据页契约

页面：`web/src/pages/EvidencePage.tsx`
API：`GET /api/v1/evidence`

### 7.1 EvidenceCard

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| sourceLevel | `source_level` | S/A/B 标签；C 级仅作为背景材料弱化展示 |
| sourceName | `source_name` | 来源名称 |
| evidenceRole | `evidence_role` | formal / background；background 弱化展示 |
| publishedAtText | `published_at` | “发布于...” |
| capturedAtText | `captured_at` | “抓取于...” |
| timeWeightText | `time_weight` | “时效权重 0.8” |
| summary | `summary` | 首屏展示摘要 |
| originalUrl | `original_url` | 详情中展示 |
| contentHash | `content_hash` | 默认折叠 |
| relevanceScoreText | `relevance_score` | 详情中展示 |

### 7.2 SourceVerificationPanel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| verificationText | `verification_status` | 已满足 / 未满足 / 仅作背景 |
| independentSourceText | `independent_source_count` | “引用 3 个独立信源” |
| highestLevelText | `highest_source_level` | “最高信源等级 A” |
| latestPublishedAtText | `latest_published_at` | “最新发布于...” |
| evidenceLinks | `evidence_ids` | 链接到证据详情 |

## 8. 持仓页契约

页面：`web/src/pages/PortfolioPage.tsx`
API：`GET /api/v1/portfolio/current`

### 8.1 PortfolioTable

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| symbolText | `symbol` | 标的代码 |
| nameText | `name` | 标的名称 |
| quantityText | `quantity` | 数量格式化 |
| costPriceText | `cost_price` | 金额格式化 |
| currentPriceText | `current_price` | 金额格式化 |
| marketValueText | `market_value` | 金额格式化 |
| profitRatioText | `unrealized_profit_ratio` | 百分比格式化 |
| positionStateTag | `position_state` | 映射为状态标签 |
| buyDateText | `buy_date` | 展示买入日期；缺失时展示暂无买入日期 |
| buyReasonSummary | `buy_reason` | 长文本折叠 |

### 8.2 PositionStateTag

| API 值 | UI 文案 | 说明 |
| --- | --- | --- |
| normal | 正常持有 | 买入逻辑完好，正常监控 |
| sell_only | 只卖不买 | 买入逻辑确认破坏 |
| frozen_watch | 冻结观察 | 暂停主动操作，等待信息明朗 |

### 8.3 Portfolio onboarding 与本地事实维护

Portfolio 页必须提供从空库到本地账户事实可用的入口。所有动作只记录本地事实，不连接券商、不下单、不承诺收益。

| UI 区域 | API | 必填与展示 |
| --- | --- | --- |
| 首次使用引导 / 本地校准 | `POST /api/v1/portfolio/adjustments` | 现金、总资产、标的代码、名称、数量、成本价、现价、买入日期、纪律状态、买入理由；风险偏好仅作为页面基础信息展示，不作为已持久化字段表达 |
| 持仓新增/编辑 | `POST /api/v1/portfolio/holdings` | position、reason、confirmation；position 可包含 `buy_date` 与 `position_state`；成功后展示 snapshot/audit 引用 |
| 移除当前持仓 | `POST /api/v1/portfolio/holdings/remove` | position_id、reason、confirmation；页面必须说明历史事实保留 |
| 线下交易记录 | `POST /api/v1/portfolio/offline-transactions` | operation_type=buy/sell/reduce、symbol、quantity、price、executed_at；页面文案必须说明“仅记录已在线下完成的动作” |
| 批量导入校验 | `POST /api/v1/portfolio/imports/validate` | 展示逐行 validation result 和 import_batch_id；校验成功且无错误后才允许确认 |
| 批量导入确认 | `POST /api/v1/portfolio/imports/confirm` | import_batch_id、confirm_reason、rows；确认按钮默认禁用，直到 validation 成功且 invalid_count=0 |
| 错误修正审计 | `POST /api/v1/portfolio/corrections` | before_json、after_json、correction_reason；页面必须说明修正只记录审计，若需改变当前状态需使用持仓编辑或线下交易 |

Dashboard 与今日纪律报告在缺账户或缺持仓时，必须链接到 `/positions`，并说明缺前提时不会生成正式每日纪律建议。

## 9. 确认操作区契约

组件：`ConfirmActionPanel`
API：`POST /api/v1/decisions/{decision_id}/confirmations`

### 9.1 操作类型

| 操作 | API 值 | 是否更新账户 | UI 提示 |
| --- | --- | --- | --- |
| 记录计划 | planned | 否 | 只记录你的计划，不改变本地账户。 |
| 已手动执行 | executed_manually | 是 | 仅记录你已在线下完成的交易，系统不会替你交易。 |
| 标记待观察 | watch | 否 | 暂不确认结果，后续复盘继续跟踪。 |
| 标记错误 | marked_error | 否 | 写入错误案例库，用于后续复盘和提案。 |

### 9.2 已手动执行必填字段

| 字段 | UI 标签 | 校验 |
| --- | --- | --- |
| operation_type | 操作类型 | buy / sell / reduce；仅“已手动执行”可选 |
| symbol | 标的代码 | 必填 |
| quantity | 成交数量 | 必须大于 0 |
| price | 成交价格 | 必须大于 0 |
| executed_at | 成交时间 | 不得晚于当前时间 |
| note | 备注 | 可为空 |

### 9.3 标记错误必填字段

| 字段 | UI 标签 | 校验 |
| --- | --- | --- |
| actual_outcome | 实际结果 | 必填，描述事实结果 |
| root_cause_tag | 错误原因标签 | evidence_missed / rule_threshold_issue / analyst_error / user_context_missing / market_exception |
| lesson_learned | 经验记录 | 必填，用于后续复盘和规则提案 |
| note | 备注 | 可为空 |

## 10. 审计页契约

页面：`web/src/pages/AuditPage.tsx`
API：`GET /api/v1/audit-events`

### 10.1 AuditTimeline

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| auditEventIdText | `audit_event_id` / `event_id` | 两者值相同；优先使用 `audit_event_id`，保留 `event_id` 兼容旧展示 |
| eventTimeText | `created_at` | 时间线节点时间 |
| actorText | `actor` | 系统 / 用户 / 守门人 |
| workflowTypeText | `workflow_type` | 工作流类型中文名 |
| nodeNameText | `node_name` | 节点名称 |
| actionText | `action` | 建议生成 / 用户确认 / 错误标注 / 提案生成 / 守门人审计 / 规则更新 |
| nodeActionText | `node_action` | 节点动作，详情中展示 |
| statusText | `status` | 成功 / 降级 / 失败 |
| errorCodeText | `error_code` | 失败时展示 |
| beforeStateText | `before_state` | 操作前状态，详情中展示 |
| afterStateText | `after_state` | 操作后状态，详情中展示 |
| ruleVersionText | `rule_version` | 执行时规则版本 |
| snapshotIdText | `snapshot_id` | 关联账户快照 ID |
| inputRefText | `input_ref` | 输入摘要或引用 ID，详情中展示 |
| relatedRecordText | `output_ref` | 关联记录 ID |

## 11. 规则提案页契约

页面：`web/src/pages/RulesPage.tsx`
API：`GET /api/v1/rule-proposals`、`GET /api/v1/rule-proposals/{proposal_id}/effect-validation`、`POST /api/v1/rule-proposals/{proposal_id}/effect-validation`、`POST /api/v1/rule-proposals/{proposal_id}/confirm`、`POST /api/v1/rule-proposals/{proposal_id}/final-confirm`

### 11.1 RuleProposalViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| proposalTitle | `title` | 提案标题 |
| proposalStatusText | `status` | 映射为待用户确认 / 审计中 / 待最终确认 / 已应用 / 已拒绝 |
| sourceErrorCaseText | `source_error_case_id` | 可跳转到错误案例 |
| sampleCountText | `sample_count` | 少于 3 个时标记样本不足 |
| sourceExplanationText | `effect_validation.source_explanation` | 展示来源错误案例、决策、风险预警和审计线索；缺失时显示待评估 |
| representativenessText | `effect_validation.representativeness_status` | 映射为样本代表性充足 / 不足 / 需更多样本 |
| overfitRiskText | `effect_validation.overfit_risk` | low / medium / high 映射为低 / 中 / 高，high 必须突出风险 |
| replayResultText | `effect_validation.replay_result` | 映射为回放通过 / 不利 / 混合 / 未知 |
| guardrailDecisionText | `effect_validation.guardrail_decision` | 映射为可进入审计 / 拒绝 / 需要用户复核 |
| validationLink | `effect_validation.validation_link` | 本地验证详情入口；只导航，不执行规则变更 |
| validationSafetyNote | `effect_validation.safety_note` | 固定展示，不得隐藏 |
| beforeRuleText | `before_rule.content` | 变更前规则 |
| afterRuleText | `after_rule.content` | 变更后规则 |
| auditResultText | `audit_result` | 审计通过 / 审计否决 / 需要用户复核 |
| auditSummaryText | `audit_summary` | 守门人审计摘要 |
| finalConfirmAction | `status=pending_final_confirm` | 展示“确认应用到正式规则”与“拒绝应用” |

### 11.2 规则提案状态展示

| API 状态 | UI 文案 | 允许动作 |
| --- | --- | --- |
| draft | 草稿 | 查看 |
| pending_user_confirm | 待用户确认 | 确认送审 / 拒绝 |
| under_gatekeeper_audit | 守门人审计中 | 查看审计进度 |
| pending_final_confirm | 待最终确认 | 确认应用 / 拒绝应用 |
| rejected | 已拒绝 | 查看原因 |
| applied | 已应用 | 查看生成的规则版本 |

规则处于 `pending_final_confirm` 时，页面必须展示最终确认与拒绝动作，并明确说明规则在最终确认成功前不会成为正式规则。

规则效果验证展示规则：

- `effect_validation` 缺失时显示 `not_evaluated` 或待评估，不得展示为通过。
- `validation_status=insufficient/failed/needs_more_samples`、`overfit_risk=high`、`replay_result=failed` 或 `guardrail_decision=rejected` 时，不得展示“确认应用”类主动作。
- 即使 `validation_status=passed`，页面仍必须说明还需要守门人审计与用户最终确认。
- 刷新验证动作只调用本地效果验证 API，不得自动应用规则、交易、回滚规则、连接券商或触发外部通知。

## 12. 设置页契约

页面：`web/src/pages/SettingsPage.tsx`
API：`GET /api/v1/settings/capability`、`PUT /api/v1/settings/capability`、`GET /api/v1/settings/system`、`PUT /api/v1/settings`、`POST /api/v1/market/refresh`、`GET /api/v1/market/snapshots/latest`、`POST /api/v1/evidence/rebuild-index`

### 12.1 CapabilitySettingsViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| assetTypes | `asset_types` | 资产类型范围 |
| includedSymbols | `symbols` | 明确纳入标的 |
| excludedSymbols | `excluded_symbols` | 明确排除标的 |
| strategyScope | `strategy_scope` | 策略范围 |
| updatedAtText | `updated_at` | 最近更新时间 |

### 12.2 SystemStatusViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| sqliteStatusText | `sqlite_status` | 本地事实库状态 |
| vecliteStatusText | `veclite_status` | 向量索引状态 |
| deepseekStatusText | `deepseek_status` | 模型配置状态，不展示完整 Key |
| deepseekModelText | `deepseek_model` | 模型名，可为空，不展示 key |
| dataSourceStatus | `data_sources` | 数据源开关与状态 |
| rebuildIndexAction | `veclite_status` | 索引异常时展示重建入口 |
| updateSettingsAction | `PUT /api/v1/settings` | 仅用于通知、页面偏好和非规则类配置；规则类设置必须生成提案 |

设置页展示完整配置时，必须包含能力圈配置、系统状态、市场快照状态、通知设置说明和索引状态。页面不得展示完整密钥；若后端只提供更新入口而不返回通知配置详情，前端必须明确说明当前后端未返回该详情。

### 12.3 MarketStatusViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| marketSnapshotId | `market_snapshot_id` | 默认折叠，详情中展示 |
| symbolText | `symbol` | 标的代码 |
| tradeDateText | `trade_date` | 交易日期 |
| pePercentileText | `pe_percentile` | 展示为“PE 分位 63%” |
| pbPercentileText | `pb_percentile` | 展示为“PB 分位 42%” |
| liquidityStateText | `liquidity_state` | normal / warning / danger 映射为中文状态 |
| sentimentStateText | `sentiment_state` | cold / neutral / hot / extreme 映射为中文状态 |
| dataStatusText | `data_status` | fresh / stale / missing；非 fresh 时展示信息不足 |
| marketMetricsDetail | `market_metrics` | 原始市场指标，默认折叠，详情中展示 |
| capitalFlowText | `market_metrics.metadata.p88_structured_fields.capital_flow` | 若存在，展示 `date`、`net_inflow`、`net_outflow` 和 `raw_net_flow`；`raw_net_flow` 为日净流向原值 |
| refreshMarketAction | `POST /api/v1/market/refresh` | 手动刷新市场数据 |

## 13. 风险预警中心契约

页面：`web/src/pages/RiskAlertPage.tsx`
API：`GET /api/v1/risk-alerts`、`GET /api/v1/risk-alerts/{alert_id}`、`POST /api/v1/risk-alerts/{alert_id}/lifecycle`

### 13.1 RiskAlertViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| riskTypeText | `risk_type` | 映射为估值高位、买入逻辑破坏、流动性危险、情绪极端、仓位超限、证据不足、数据降级 |
| severityText | `severity` | info / warning / critical 映射为提示、预警、严重 |
| sopStatusText | `sop_status` | triggered / active / observing / escalated / resolved / archived 映射为中文 SOP 状态 |
| symbolText | `symbol` | 展示影响标的 |
| triggerSummary | `trigger_summary` | 展示触发依据 |
| prohibitedActions | `prohibited_actions` | 展示禁止动作；为空时展示暂无 |
| suggestedActions | `suggested_actions` | 展示建议人工动作；为空时展示暂无 |
| relatedLinks | `decision_link` / `report_link` / `notification_link` / `audit_link` | 展示本地追踪入口 |
| lifecycleActions | lifecycle API | active / observing / escalated 等非终态展示继续观察、升级复核、解除预警等本地 SOP 动作；终态不展示动作 |
| safetyNote | `safety_note` | 固定展示，不得隐藏 |

### 13.2 风险预警交互规则

风险预警中心默认展示 triggered / active / observing / escalated 风险，也可通过详情路径 `/risk-alerts/:alertId` 查看单条风险。生命周期动作只更新本地 `risk_alerts` 状态并写审计，不执行交易、不连接券商、不外部推送、不自动应用规则。

## 14. 通知页契约

页面：`web/src/pages/NotificationPage.tsx`
API：`GET /api/v1/notifications`、`POST /api/v1/notifications/{notification_id}/read`、`POST /api/v1/notifications/read-all`

### 13.1 NotificationCenterViewModel

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| unreadCountText | `unread_count` | 展示为“未读通知：N” |
| notificationId | `items[].notification_id` | 内部 ID，默认不直接展示 |
| notificationTypeText | `items[].type` | 映射为数据源、索引、证据核验、规则提案、复盘等通知类型 |
| severityText | `items[].severity` | info / warning / critical 映射为中文严重程度 |
| titleText | `items[].title` | 通知标题 |
| messageText | `items[].message` | 通知正文 |
| sourceTypeText | `items[].source_type` | 关联来源类型，详情中展示 |
| sourceIdText | `items[].source_id` | 关联来源 ID，详情中展示 |
| readStateText | `items[].read_at` | 空值展示“未读”；非空展示“已读 {read_at}” |
| createdAtText | `items[].created_at` | 通知创建或最近刷新时间 |
| markOneReadAction | `POST /api/v1/notifications/{notification_id}/read` | 仅未读通知展示“标记已读” |
| markAllReadAction | `POST /api/v1/notifications/read-all` | 未读数为 0 时禁用“全部标记已读” |

### 13.2 通知页交互规则

通知页通过本地 API 轮询刷新应用内通知，必须展示未读数量、通知列表、单条标记已读和全部标记已读动作。当前前端默认约 30 秒轮询一次；轮询只读取本地 `notifications` 表；标记已读只更新本地通知状态。

通知页不得发送邮件、短信、系统 Push、Webhook 或 WebSocket 消息；不得执行交易；不得自动应用规则提案；不得把通知动作扩展为外部告警渠道。

## 15. 复盘页契约

页面：`web/src/pages/ReviewSummaryPage.tsx`
API：`GET /api/v1/review/summary`

| UI 字段 | API 字段 | 展示规则 |
| --- | --- | --- |
| periodText | `period` | monthly → 月度复盘；quarterly → 季度复盘 |
| decisionCountText | `decision_count` | 建议数量 |
| confirmationCountText | `confirmation_count` | 用户确认数量 |
| executedCountText | `executed_manually_count` | 已手动执行数量 |
| plannedCountText | `planned_count` | 记录计划数量，不计为实际交易 |
| errorCaseCountText | `error_case_count` | 错误案例数量 |
| ruleProposalCountText | `rule_proposal_count` | 规则提案数量 |
| auditEventCountText | `audit_event_count` | 审计事件数量 |
| ruleHitCountText | `rule_hit_count` | 规则命中次数 |
| misjudgmentCountText | `misjudgment_count` | 误判案例数量 |
| missingEvidenceCountText | `missing_evidence_count` | 缺失证据数量 |
| degradedCountText | `degraded_count` | 降级记录数量 |
| opsStatusPanel | `ops_status` | 运维状态面板（数据源 / 索引 / 复盘状态），仅展示，不执行操作 |
| ruleSuggestionsList | `rule_suggestions` | 规则建议列表；只读展示，不可自动应用规则；规则变更须守门人审计与用户最终确认 |
| ruleEffectTrackingList | `rule_effect_tracking` | 展示已应用规则版本的命中、误判、缺证据、降级、风险预警和趋势；只读展示 |
| trackingLinksList | `tracking_links` | 追踪入口，指向审计事件 / 规则提案 / 决策 / 错误案例 / 规则效果追踪；仅导航，不执行任何变更 |

`ops_status` 展示规则：
- `data_source_status`：success / degraded / failed / empty / unknown，使用 `opsStatusText` 映射为中文
- `index_status`：success / degraded / failed / missing / unknown，使用 `opsStatusText` 映射为中文
- `review_status`：success / degraded / failed / empty / unknown，使用 `opsStatusText` 映射为中文
- 未知状态必须显示安全兜底文案，不得展示为成功状态
- 面板底部须显示「仅展示状态与追踪入口，不执行交易，也不自动应用规则」

`rule_effect_tracking` 展示规则：

- `trend_direction` 使用 improved / flat / worsened / unknown 映射为改善 / 持平 / 恶化 / 信息不足。
- `worsened` 时展示风险提示或后续提案入口文案，但不得展示自动回滚、自动替换规则或交易动作。
- 样本不足或缺事实时展示信息不足，不得宣称规则效果改善或恶化。
- 每条追踪必须展示 `safety_note` 或等价安全文案。

P5 页面必须展示建议数量、确认动作、错误案例、规则提案和审计事件汇总。P16 起须同时展示运维状态、规则建议与追踪入口。

## 16. P42 用户决策工作台契约

页面：`web/src/pages/WorkbenchPage.tsx`
路由：`/workbench`
API：`GET /api/v1/dashboard/today`、`GET /api/v1/daily-discipline/reports/today`、`GET /api/v1/portfolio/current`、`GET /api/v1/risk-alerts?status=active,escalated`、`GET /api/v1/rule-proposals`、`GET /api/v1/review/summary`

P42 工作台是只读聚合与导航面板。页面只能通过 `web/src/services/` 消费现有 DTO，不得直接读取 SQLite、VecLite、本地日志、配置文件或临时诊断文件；不得新增券商接口、交易执行、外部推送、自动确认或自动规则应用入口。

### 16.1 区域与数据来源

| UI 区域 | 数据来源 | 展示规则 | 导航 |
| --- | --- | --- | --- |
| 今日先看 | dashboard today、today daily discipline report | 展示今日纪律报告摘要、状态、最终裁决、证据覆盖和触发规则；缺失或错误时展示安全状态，不伪造报告 | `/daily-discipline/reports/{report_id}`、`/` |
| 组合与风险 | portfolio current、active/escalated risk alerts | 展示总资产、现金占比、高风险占比、持仓数量、活跃风险数量和前两条风险摘要；空库时提示先完成本地账户校准 | `/positions`、`/risk-alerts` |
| 规则与复盘 | rule proposals、review summary | 展示待确认规则数量、复盘决策数量、错误样本和审计事件；降级时展示稳定错误码映射后的安全文案 | `/rules`、`/review`、`/audit` |
| 主动咨询入口 | 固定安全文案 | 只说明用户可手动提交咨询问题；不得在工作台自动提交咨询请求或写入确认 | `/consultation` |

### 16.2 状态与降级

- 任一聚合接口失败时，其他区域仍应正常展示；失败区域使用 `StatusNotice` 和 `toPageErrorState` 映射后的安全文案。
- 空库、缺账户、缺持仓、数据过期、source health 降级、LLM/RAG/VecLite 不可用时，工作台必须展示明确的安全状态与可导航的检查入口。
- 页面不得把未知、失败或降级状态渲染为成功；不得展示内部 SQL、文件路径、密钥、HTTP 原始错误或供应商原始报错。
- 页面可展示已有 `prohibited_actions` 或风险摘要，但不能把它们变成交易、外推、确认或规则生效操作。

### 16.3 安全与测试

- 工作台只允许链接到每日纪律报告、今日纪律、持仓、风险预警、规则提案、复盘摘要、审计和主动咨询页面。
- 工作台不得出现自动下单、一键交易、代下单、券商下单、券商接口、自动规则应用、外部推送或收益承诺入口。
- Vitest 必须覆盖成功聚合、空库、降级/错误和安全文案；Playwright smoke 必须覆盖 `/workbench` 可达、四个核心区域、窄屏可用和禁止入口扫描。

## 17. P43 数据质量可观测页面契约

页面：`web/src/pages/DataQualityPage.tsx`
路由：`/data-quality`
API：`GET /api/v1/settings/system`、`GET /api/v1/market/snapshots/latest`、`GET /api/v1/market/source-health`、`GET /api/v1/evidence`、`GET /api/v1/evidence/verification`、`GET /api/v1/review/summary`、`GET /api/v1/data-source-quality/regression`、`GET /api/v1/data-source-quality/gate-resolution`、`GET /api/v1/data-source-quality/resolutions`、`POST /api/v1/data-source-quality/resolutions`、`POST /api/v1/data-source-quality/resolutions/{resolution_id}/retire`

P43 数据质量可观测页以只读质量事实为核心；P67 允许新增本地人工 gate resolution 记录与退役动作。页面只能通过 `web/src/services/` 消费 DTO，不得直接读取 SQLite、VecLite、本地日志、localStorage、sessionStorage、配置文件或临时诊断文件；不得新增券商接口、交易执行、外部推送、自动确认、自动修复或自动规则应用入口。

### 17.1 区域与数据来源

| UI 区域 | 数据来源 | 展示规则 | 导航 |
| --- | --- | --- | --- |
| 数据源健康 | market latest、market source health | 展示市场数据状态、影响标的、数据源类别、新鲜度、数据日、最近成功、最近失败、失败类别和影响范围；未知状态使用安全兜底文案 | `/settings`、`/workbench` |
| 当前数据门禁处置 | data-source-quality regression、gate-resolution、resolutions | 展示 P66 policy、release claim state、clean data claim 是否允许、active resolution、allowed/prohibited claims；允许用户本地记录 `scope_exclusion` 或 `waiver`，以及退役 active resolution | 留在 `/data-quality`，可查看 `/audit` |
| 证据与检索 | evidence list、evidence verification、system settings | 展示证据数量、验证状态、独立信源数量、高等级独立信源数量、最高信源等级和 VecLite 状态；不展示证据长摘要、向量 ID 或本地索引路径 | `/evidence`、`/audit` |
| LLM 质量 | system settings、review summary | 展示 DeepSeek 配置状态、复盘状态、数据源状态、索引状态、降级数量和缺证据数量；只展示已脱敏诊断摘要 | `/review` |
| 影响范围与下一步 | review summary | 展示受影响 workflow 的决策 ID、标的和状态；下一步只导航到权威页面检查，不触发刷新、修复、确认或规则生效 | `/risk-alerts`、`/decisions/{decision_id}`、`/review` |

### 17.2 状态、脱敏与安全边界

- source_unavailable、parse_error、stale、missing、unknown、degraded、quality_failed、VecLite/RAG/LLM 不可用时，页面必须展示明确安全状态，不得渲染为成功。
- 页面不得展示完整 key、完整 prompt、私有本地路径、SQL、供应商原始响应、HTTP 原始错误、账户敏感明细、SQLite 路径或 VecLite 路径。
- 失败分类、健康状态、系统状态和运维状态必须使用稳定 mapper 或安全兜底文案；未知原始值不得直接出现在页面主文案中。
- 页面只允许导航到设置、证据、复盘、审计、风险预警、决策详情和工作台；P67 表单只允许写入本地 resolution 记录或退役 resolution。不得出现自动刷新修复、重建索引、外部推送、自动确认、自动应用规则、自动交易、一键交易、代下单、券商接口或收益承诺入口。
- `policy=blocked` 时 UI 只能允许 `scope_exclusion`，不得允许 `waiver` 创建；`policy=waiver_required` 可允许 `waiver` 或 `scope_exclusion`；`policy=passed` 不显示可用创建动作。
- Vitest 必须覆盖成功聚合、空库或缺数据、降级/错误、unknown 状态、脱敏和禁止入口；Playwright smoke 必须覆盖 `/data-quality` 可达、四个核心区域、窄屏可用和禁止入口扫描。

## 18. P39 浏览器完整用户旅程契约

P39 浏览器级 E2E 必须从本地临时 fixture 观察页面和 API 表现，不直接读取 SQLite、VecLite 或本地文件内容。旅程覆盖 Dashboard、持仓、主动咨询、决策详情、证据、每日纪律报告、审计、复盘、规则治理、风险预警和设置页。

### 18.1 旅程覆盖

- 空库或缺前提状态必须展示安全 onboarding / 初始化入口，并能通过本地账户与持仓校准进入可验收状态。
- 主动咨询必须能生成可渲染决策详情；用户确认动作只记录线下事实或计划，不改变持仓、不连接券商、不产生交易执行语义。
- 每日纪律、复盘、规则、风险预警和设置页必须展示 source health、retrieval quality、risk alert / SOP、规则提案、守门人或最终确认边界等跨页上下文。
- 缺市场数据、证据不足、VecLite/RAG 降级、LLM 降级、能力圈外和规则提案待确认必须展示为安全降级或拒绝状态，不得伪装成成功。

### 18.2 浏览器稳定性与安全扫描

- Playwright 必须捕获 unexpected console error 和 unhandled page error；允许诊断日志时必须在 fixture 内明确限定。
- 窄屏 smoke 必须确认主导航、关键标题、状态标签和主要动作按钮可见且不改变安全语义。
- 表单和按钮必须具备可访问名称，输入控件必须具备 label 或等价描述，导航 landmark 必须可被浏览器语义查询发现。
- 安全扫描必须在页面异步数据就绪后检查关键页面按钮和链接，不得出现自动下单、一键交易、代下单、券商下单、券商接口、自动规则应用、外部推送或收益承诺入口。
- Vitest 只覆盖组件、mapper 和 API client 行为；Playwright 只覆盖浏览器旅程，两者不得互相收集测试文件或共享持久可变测试状态。

## 19. P40 本地运行就绪与健康面板契约

设置页或等价运维页必须展示 P40 本地运行就绪信息，数据来源只能是 API/service DTO 或本地 CLI 诊断输出摘要，不得直接读取 SQLite、VecLite 或本地文件。

### 19.1 运行就绪展示

| UI 字段 | 来源 | 展示规则 |
| --- | --- | --- |
| sqliteReadinessText | `system.sqlite_status` | 映射为可用、重建中、降级、失败、缺失、不可用或未知状态 |
| vecliteReadinessText | `system.veclite_status` | 映射为已配置、可用、重建中、降级、失败、缺失、不可用或未知状态 |
| llmReadinessText | `system.deepseek_status` | 只展示配置、可用、重建中、降级、失败、缺失、不可用或未知状态，不展示 key |
| preflightCommandText | 固定文案 | 展示 `go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json` |
| readinessSafetyText | 固定文案 | 说明只展示本地诊断和修复提示，不执行交易、不外部推送、不自动应用规则 |

### 19.2 数据源健康展示

- 数据源健康必须展示最近成功时间、最近失败时间、失败分类、新鲜度和影响标的或范围。
- fresh / stale / failed / missing / unknown 等状态必须有文字区分，不得只用颜色表示。
- 健康面板只允许查看、过滤、刷新本地行情事实或导航到相关审计/诊断，不得展示交易执行、外部推送、自动恢复承诺或自动规则应用入口。
- 恢复 smoke 和预检失败时，前端只展示安全修复提示；不得建议使用真实私有数据库复现。

## 20. P44 本地安装与诊断页面契约

页面：`web/src/pages/LocalInstallPage.tsx`
路由：`/local-install`

P44 本地安装页是只读运维页面，展示安装向导、脚本命令和本地诊断摘要。

### 20.1 安装命令与状态展示

| 区域 | 展示规则 |
| --- | --- |
| 配置向导 | 允许用户在页面输入 server、端口、SQLite 路径、VecLite 路径和 DeepSeek 基础配置，生成可复制的启动配置草稿；不得直接写入本地文件。 |
| 关键命令 | 展示 `--preflight`、`recovery-smoke`、`e2e-smoke`、`local-install-diagnostics` 及其说明文本，默认可读且保持最小权限。 |
| 摘要导入 | 允许上传本地 JSON 摘要并只读渲染步骤状态、失败项与日志路径；不展示数据库私有路径、完整 key、原始 SQL、供应商原始 HTTP 响应。 |
| 安全边界 | 明示页面职责仅为引导与观察，不提供交易、自动确认、自动规则应用、外部推送、收益承诺。 |

### 20.2 约束与测试

- 页面仅消费本地文件输入和静态文案，不调用新增后端任务 API（除非显示文档链接）。
- 页面不得出现 `自动交易`、`一键交易`、`券商`、`代下单`、`应用规则` 等高风险入口。
- Vitest 必须覆盖配置向导、失败项展示与禁止项扫描；Playwright smoke 必须覆盖 `/local-install` 可达与文本扫描。

## 21. P46 本地知识导入页面契约

页面：`web/src/pages/LocalKnowledgePage.tsx`
路由：`/local-knowledge`
API：`POST /api/v1/local-knowledge/imports/validate`、`POST /api/v1/local-knowledge/imports/confirm`

P46 本地知识导入页用于把用户自有研究记录写入本地背景事实。页面必须先展示脱敏预览、批次 ID、逐行风险和索引计划；只有校验无阻断时才允许用户写入。

### 21.1 区域与数据来源

| 区域 | API 字段 | 展示规则 |
| --- | --- | --- |
| 导入草稿 | `source_label`、`default_symbol`、`rows[]` | 支持 JSON 数组输入；字段包括 title、text、symbol、as_of_date、tags；页面不得直接读取本地文件路径。 |
| 脱敏预览 | `data.rows[].title_preview`、`text_preview`、`risks` | 只展示后端返回的 preview；不得展示完整 key、私有路径、原始 SQL、私钥或完整 prompt。 |
| 批次摘要 | `data.import_batch_id`、`summary` | 展示总数、可写入、需关注和阻断数量；`blocking_count > 0` 时确认按钮禁用。 |
| 索引计划 | `data.index_plan` | 展示预计 RAG 片段数量和 pending 状态；不得在页面承诺已经完成索引。 |
| 确认结果 | confirm response | 展示写入计数、核验计数、审计事件数量和索引状态；不得把 C/background 材料展示为正式证据。 |

### 21.2 状态、脱敏与安全边界

- 页面必须通过 `web/src/services/localKnowledge.ts` 调用 API，不直接访问 SQLite、VecLite、localStorage、sessionStorage 或本地文件。
- confirm 请求必须回传 `import_batch_id`、`source_label`、`default_symbol`、`confirm_reason` 和原始 rows，由后端重算批次绑定。
- 页面可展示本地背景事实写入结果，但不得提供交易执行、规则生效、外部通知或收益保证入口。
- Vitest 必须覆盖 validate 调用、blocking 展示、确认启用条件、confirm 批次回传和安全文案；Playwright smoke 必须覆盖 `/local-knowledge` 可达与禁止入口扫描。

## 22. P47 决策闭环解释页面契约

页面：`web/src/pages/DecisionLoopPage.tsx`
路由：`/decision-loop`
API：`GET /api/v1/decision-loops`

P47 决策闭环解释页用于只读串联建议、用户确认、本地线下记录、风险线索、复盘线索和审计线索。页面不得提供确认、交易、风险生命周期、规则生效、通知发送或设置修改控件。

### 22.1 区域与数据来源

| 区域 | API 字段 | 展示规则 |
| --- | --- | --- |
| 闭环概览 | `total`、`items[].loop_status`、`items[0]` | 展示闭环条数、未闭合条数和最近决策；空列表展示安全空态。 |
| 阶段链路 | `items[].stages` | 展示 recommendation、confirmation、manual_record、risk_review、review 的 label、status 和 summary；只读显示引用 ID。 |
| 缺口 | `items[].missing_links` | 展示缺少确认、本地流水、风险线索或复盘线索；无缺口时展示空态。 |
| 人工记录 | `items[].manual_actions` | 展示确认 ID、确认类型、操作类型、数量、价格、费用、流水 ID 和脱敏备注预览。 |
| 追踪链接 | `risk_links`、`review_links`、`audit_links` | 只渲染本地导航链接，指向风险预警、复盘摘要和审计页面。 |
| 安全边界 | `safety_note` | 固定展示只读说明，不得隐藏。 |

### 22.2 状态、脱敏与安全边界

- 页面必须通过 `web/src/services/decisionLoop.ts` 调用 API，不直接访问 SQLite、VecLite、localStorage、sessionStorage 或本地文件。
- 页面只能展示 `note_preview`，不得展示 raw payload、完整 key、私有路径、原始 SQL、完整 prompt、供应商原始响应或外部订单类信息。
- 页面可从工作台、复盘摘要和主导航进入；入口只用于导航，不改变任何本地事实。
- Vitest 必须覆盖成功展示、缺口展示、空态/错误态和无写入动作按钮；Playwright smoke 必须覆盖 `/decision-loop` 可达与安全文本扫描。

## 23. P56 UI 验收阻断与产品化设计修复契约

P56 修复 P55 真实 UI 验收 blocker，并把核心前端体验升级为任务型本地投资纪律工作台。P56 只改变前端渲染、布局、测试和验收材料，不新增交易、外推、自动确认或自动规则应用能力。

### 23.1 真实 LLM 决策详情 nullable DTO

- 决策详情页必须能安全渲染真实 LLM-backed decision DTO 中 `final_verdict.optional_actions`、`final_verdict.prohibited_actions` 为 `null`、缺失或空数组的情况。
- 缺失或 nullable 的 list 字段必须展示安全空态，例如“暂无”，不得直接调用 array-only 方法导致页面级崩溃。
- nullable list 不得被解释为允许交易、允许自动执行或安全成功。
- 前端测试必须包含 real LLM-like DTO fixture，覆盖 nullable 与缺失 verdict list 字段，并确认最终裁决、安全边界和决策编号仍可见。

### 23.2 任务型 app shell

- App shell 必须按用户任务组织主导航，而不是平铺路由清单。当前任务组为：今日、决策、组合、证据、治理、系统。
- 所有 P55/P56 主要路由仍必须可达：Dashboard、Workbench、Decision Loop、Data Quality、Positions、Consultation、Decision Detail、Evidence、Rules、Audit、Notifications、Risk Alerts、Daily Auto Run、Daily Reports、Review、Local Install、Local Knowledge、Settings。
- 390px 移动视口不得保留永久宽侧栏；必须使用可达的导航按钮或等价折叠导航，展开后仍能访问主导航 landmark。

### 23.3 操作型 UI 基础层

- 表单、按钮、状态标签、卡片、表格和 key-value 数据展示必须使用一致的 spacing、边框、字体层级和 focus 样式。
- 关键表单不得以浏览器默认控件作为主视觉处理；label、hint、error/success 状态和 primary/secondary action 必须清晰。
- 状态展示必须区分正常、观察、信息不足、降级、失败、高风险和未知状态；未知/降级不得被展示为普通成功。

### 23.4 移动端 reflow

- `/positions` 在 390px 视口下不得出现页面级横向滚动；账户/持仓表单必须纵向堆叠，持仓表必须转为可读的 labeled rows 或局部二维容器。
- `/data-quality` 在 390px 视口下不得出现页面级横向滚动；source name、source category、freshness、decision id、diagnostic text 必须在页面宽度内换行或在局部容器内处理。
- 核心移动 smoke 必须验证 `/workbench`、`/data-quality`、`/positions`、`/risk-alerts` 的 `body.scrollWidth` 和 `documentElement.scrollWidth` 不超过 viewport 宽度。

### 23.5 设计与验收证据

- P56 验收必须记录 Product Design get-context/audit/research 输入与实际 UI 改造的对应关系。
- 真实 UI 验收必须启动本地后端和 Vite 前端，通过浏览器操作真实页面，并保存桌面/移动截图。
- 使用真实 LLM 时，key 只能存在于临时运行配置或用户明确允许的本地配置中；提交材料不得包含完整 key、raw prompt、私有 SQLite、供应商原始 payload 或临时配置文件。

## 24. P57 产品体验打磨路线图契约

P57 将 P56 后的产品设计、UI 设计和功能设计打磨拆成 P58-P63 后续阶段。P57 不修改运行时代码，不声明所有产品体验问题已修复，也不刷新最终 release-ready 口径。

### 24.1 产品北极星

- 前端必须被规划为本地投资纪律工作台，不得被规划为券商交易终端、AI chat demo、营销落地页或工程调试台。
- 后续产品体验改造必须优先回答：今天能不能动、为什么、需要什么人工动作、数据和规则是否可信。
- 高风险、未知、降级、过期、缺失、信息不足和 blocked 状态不得被文案或视觉处理成普通成功。

### 24.2 后续阶段拆分

| 阶段 | 方向 | 必须覆盖 |
| --- | --- | --- |
| P58 | 今日工作台重构 | Dashboard、Workbench、今日状态、下一步人工动作、移动端第一屏 |
| P59 | 决策解释体验重构 | Consultation、Decision Detail、Evidence、Decision Loop、真实 LLM 决策解释 |
| P60 | 组合、风险与数据质量体验重构 | Positions、Risk Alerts、Data Quality、数据质量降级和移动端 reflow |
| P61 | 治理和运维页面产品化 | Rules、Audit、Notifications、Daily Reports、Daily Auto Run、Local Install、Local Knowledge、Settings |
| P62 | 设计系统与可访问性验收 | Button、Field、StatusBadge、PageHeader、SummaryCard、DetailSection、ResponsiveTable、EmptyState、ErrorState、键盘路径、WCAG reflow |
| P63 | 全量真实 UI 回归与发布状态刷新 | 真实后端、真实前端、真实 LLM、全路由、移动端、发布口径 |

### 24.3 阶段验收门禁

- 每个产品打磨阶段必须创建独立 OpenSpec change，并在方案中写明 scope、out of scope、安全边界、Product Design 依据和浏览器验收方式。
- 涉及前端行为、布局、信息架构、组件原语或用户可见文案的阶段，必须包含前端单元或组件测试、前端构建、必要的后端测试、浏览器操作本地 UI、桌面和移动端截图、移动端 reflow 检查、安全文案扫描、敏感信息扫描和子 agent 审查。
- 涉及 Consultation、Decision Detail、Evidence、LLM quality 或 Decision Loop 的阶段，必须包含真实 LLM 验证。
- 发布状态刷新必须后置到产品体验打磨阶段完成或明确豁免后执行；不得在 P58-P62 未验收前声明所有产品设计、UI 设计或前端问题已完全修复。

### 24.4 禁止事项

任何产品打磨阶段不得新增或暗示券商连接、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。

## 25. P59 决策解释体验契约

P59 将 `/consultation`、`/decisions/:decisionId`、`/evidence` 和 `/decision-loop` 串成可读的决策解释链。该阶段只改变前端展示、view model、测试和验收材料，不新增后端 API、SQLite schema、Eino 节点、交易能力或规则裁决能力。

### 25.1 主动咨询 `/consultation`

- 页面标题必须表达主动咨询语义，输入区域必须清晰呈现标的、场景、问题和输入假设。
- 生成前必须展示安全边界：系统只生成本地分析材料和规则裁决记录，不自动交易、不自动确认、不自动应用规则。
- 生成成功后必须展示解释路径，并提供本地导航到生成的决策详情、证据、决策闭环和审计页面。
- 生成失败或降级时必须展示可恢复错误或安全空态，不得输出交易建议、收益承诺或自动执行提示。

### 25.2 决策详情 `/decisions/:decisionId`

- 首屏必须先展示决策故事：最终裁决或安全空态、标的、问题、生成时间、关键原因、可信度和安全边界。
- 禁止动作、可选人工动作和人工复核边界必须位于长技术 trace 之前。
- Evidence、LLM、rules、expected return、arbitration、audit 和 confirmation 只能作为分层详情展示；长列表必须保持可扫描或可折叠。
- `final_verdict`、`analyst_reports`、`retrieval_quality`、`evidence_chain`、`expected_return_scenarios`、`audit` 或 `confirmation` 中的 null/missing/degraded/unknown/failed/insufficient 字段必须渲染安全中文空态，不得被解释为允许交易或成功完成。

### 25.3 Evidence `/evidence`

- 页面必须优先展示证据可信度概览、来源等级说明和决策解释入口，再展示证据明细表。
- 证据表必须保留筛选和展开能力；移动端必须使用可读 reflow 或局部容器，避免页面级横向滚动。
- 页面不得展示 raw vendor payload、完整 prompt、完整 key、私有路径、原始 SQL、私有 SQLite 内容或外部订单类信息。

### 25.4 Decision Loop `/decision-loop`

- 页面必须表达只读决策生命周期：建议生成、用户确认、线下记录、风险/复盘和审计。
- 最近决策存在时必须提供本地决策详情链接；缺口必须表达为人工 follow-up，而不是自动写入动作。
- 页面不得提供确认创建、交易、风险生命周期变更、规则应用、通知发送或设置修改按钮。

### 25.5 P59 验收门禁

- Vitest 必须覆盖 decision explanation view model、真实 LLM-like nullable DTO、咨询解释路径、决策详情首屏解释、Evidence 可信度入口和 Decision Loop 只读生命周期。
- Playwright smoke 或浏览器验收必须覆盖 P59 路由可达、解释链接、安全文案和 forbidden copy scan。
- 真实本地后端和 Vite 前端必须启动，通过前端 UI 提交一次 consultation，并打开生成的决策详情。
- 桌面和 390px 移动端必须采集截图或等价浏览器证据；`body.scrollWidth` 和 `documentElement.scrollWidth` 不得超过 viewport。
- 真实 LLM 如遇外部依赖不可用，验收记录必须区分远端失败和本地 UI 降级行为，不得把降级伪装为 LLM 成功。

验收记录：`docs/release/acceptance/2026-06-18-p59-ui-acceptance.md`。

## 26. P60 组合、风险与数据质量体验契约

P60 将 `/positions`、`/risk-alerts` 和 `/data-quality` 从基础表单/列表打磨为日常维护、风险处置和质量检查工作台。P87 后 `/positions` 可通过既有 portfolio write service 录入买入日期与纪律状态；仍不得新增交易能力、外部推送能力、券商连接或自动确认能力。

### 26.1 组合维护 `/positions`

- 首屏必须先展示组合维护状态、当前阶段、关键账户指标和下一步人工动作，再展示表单和明细表。
- 前端必须使用 `GET /api/v1/portfolio/current` 与 portfolio write service；可读写本地事实字段 `buy_date` 与 `position_state`，不得新增交易、券商或外部推送字段。
- `PortfolioExperienceModel` 只能接收 portfolio DTO、导入状态和页面错误态；不得直接访问 SQLite、VecLite、本地文件、localStorage、sessionStorage 或临时配置。
- 页面必须区分初始化/校准、持仓维护、线下交易补记、批量导入和错误修正五类模式。
- 空态、错误态、导入待确认、高风险比例和正常维护状态都必须有明确中文状态和人工下一步。
- 校准、编辑、移除、导入、线下交易和修正动作只能写入本地账户事实或审计事实；不得出现券商连接、下单、收益承诺或自动处理入口。

### 26.2 风险处置 `/risk-alerts`

- 首屏必须展示风险处置总览、最高严重程度、影响标的和下一步人工动作。
- 风险列表必须按 SOP 队列分组：`待看`、`处理中`、`需复盘`、`已记录`。
- `RiskDispositionModel` 只能接收 risk alert DTO；队列映射必须覆盖 `triggered`、`active`、`observing`、`escalated`、`resolved`、`archived`。
- 风险卡片必须展示风险类型、严重程度、SOP 状态、更新时间、触发摘要、禁止动作、建议人工动作和本地关联链接。
- 生命周期按钮只允许记录本地 SOP 状态；`resolved` 和 `archived` 风险不得展示生命周期按钮。
- 页面不得新增交易、外部通知发送、自动确认、自动规则应用或组合写入入口。

### 26.3 数据质量 `/data-quality`

- 首屏必须展示数据质量总览、下一步本地检查和四类质量信号：数据源健康、证据与 RAG、LLM 分析、影响范围。
- `DataQualityExperienceModel` 只能组合 system、market、source health、evidence、verification、review、data-source-quality regression/gate DTO 和页面错误态。
- `degraded`、`stale`、`missing`、`parse_error`、`source_unavailable`、`failed`、`unknown`、`disabled`、`no_data`、`unavailable` 等状态不得展示为普通成功。
- 页面质量事实必须保持只读；P67 只允许本地 gate resolution 记录和退役，不发起后台数据刷新、规则确认、规则生效或资金动作。
- 页面不得渲染完整 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack；有诊断时只能展示脱敏摘要。
- 数据质量页面主体和可操作元素必须通过 forbidden copy scan，尤其不得出现自动修复、自动确认、自动规则应用、一键交易、代下单、收益承诺或密钥形态。

### 26.4 P60 验收门禁

- Vitest 必须覆盖三个 experience model 和三个页面：状态映射、空态/错误态、降级态、导入待确认、风险队列、脱敏和禁止入口。
- Playwright smoke 必须覆盖 P60 路由可达、关键状态、风险队列、质量信号和 forbidden copy scan。
- 必须真实启动本地后端和 Vite 前端，使用浏览器操作 `/positions`、`/risk-alerts`、`/data-quality`。
- `/positions` 必须执行一次本地校准或维护路径；`/risk-alerts` 如存在 eligible SOP 按钮应执行一次本地 SOP 生命周期动作；`/data-quality` 必须验证质量总览、本地导航、脱敏和只读边界。
- 桌面和 390px 移动端必须采集截图或等价浏览器证据；`body.scrollWidth` 和 `documentElement.scrollWidth` 不得超过 viewport。

验收记录：`docs/release/acceptance/2026-06-18-p60-ui-acceptance.md`。

## 27. P61 治理和运维页面产品化

P61 覆盖 `/rules`、`/audit`、`/notifications`、`/daily-discipline/reports`、`/daily-auto-run`、`/local-install`、`/local-knowledge`、`/settings`。这些页面必须从工程工具页升级为治理/运维工作台：首屏先展示状态、指标、下一步人工动作，再展示详情。

### 27.1 通用约束

- 页面只能消费现有 API/service DTO，不得新增后端能力、SQLite schema、Eino workflow、券商接口、交易执行或外部通知发送能力。
- 页面必须复用 P58-P60 operational cockpit tokens：`daily-hero`、`daily-signal-grid`、`daily-hero-side`、`cockpit-card`、`table-wrap`。
- 390px 下不得出现页面级横向溢出；长表格、JSON 或诊断摘要只能在局部容器滚动。
- 页面 body 和可操作元素不得渲染完整 key、私有路径、SQL、完整 prompt、raw vendor payload、本地数据库路径或 raw stack。
- 页面不得出现自动下单、一键交易、代下单、券商接口、外部推送、短信、邮件、Webhook、第三方推送、自动确认、自动修复、自动规则应用、收益承诺或密钥形态。

### 27.2 页面要求

- Rules 必须展示规则治理状态、当前规则、提案数量、待确认/最终确认/门禁风险、规则提案原因、样本、过拟合、守门人、审计和人工确认边界；来自后端的高风险负向 safety note 必须安全改写后展示。
- Audit 必须展示审计检查状态、事件计数、分类摘要、最近活动、空态/错误态和时间线。
- Notifications 必须展示本地通知收件箱、未读/严重/预警/总数指标、来源分类、已读处理状态和本地-only 边界。
- Daily Reports 必须展示每日纪律复盘状态、最新报告、证据覆盖、自动运行、执行范围、缺口和报告入口；空态不得伪造成成功。
- Daily Auto Run 必须展示运行健康、计划/最近/下次执行、失败诊断、缺失前提、关联入口和默认关闭或显式启用边界。
- Local Install 必须展示配置草稿、关键命令、诊断摘要上传、失败步骤和本地复验入口。
- Local Knowledge 必须保持 validate -> preview -> explicit confirm，展示脱敏预览、索引计划、确认理由和本地事实写入边界。
- Settings 必须统一展示本地配置与诊断状态、能力圈、系统状态、数据源健康、市场刷新和错误摘要。

### 27.3 验收门禁

- Vitest 必须覆盖 governance view model 和八个页面的关键状态、空态、错误态、脱敏和 forbidden copy。
- Playwright smoke 必须覆盖 P61 路由可达、关键状态、下一步动作、390px reflow 和 forbidden copy scan。
- 必须真实启动本地后端和 Vite 前端，使用浏览器操作八个 P61 页面；Local Install 需上传诊断摘要 fixture，Local Knowledge 需执行 validate，Settings 如有入口需执行市场刷新。
- 验收记录：`docs/release/acceptance/2026-06-18-p61-ui-acceptance.md`。

## 28. P62 设计系统与可访问性验收

P62 将 P58-P61 反复出现的 operational cockpit UI 规则固化为轻量设计系统 primitives，并补齐键盘路径、可访问语义、390px/768px/1280px reflow 和视觉证据门禁。该阶段不新增后端 API、SQLite schema、Eino workflow、LLM 能力、数据源能力、交易能力或发布状态刷新。

### 28.1 基础组件

- Button 必须提供稳定可访问名称，支持 primary、secondary、ghost、danger、link-like 语义；disabled 或 working 状态必须通过文本、属性或等价语义暴露。
- Field 必须把 label、hint、error 和控件关联；错误态必须可被文本识别。
- StatusBadge 必须覆盖 success、warning、danger、degraded、unknown、readonly、blocked，且不得只依赖颜色。
- PageHeader、SummaryCard、DetailSection 必须保持标题层级、状态、下一步人工动作和折叠语义稳定。
- ResponsiveTable 必须提供 caption 或可访问名称；移动端必须通过 `data-label` 或等价方式保留列含义。
- EmptyState / ErrorState 必须使用安全中文空态或错误态，不得渲染 raw stack、完整 key、私有路径、SQL、完整 prompt 或 raw vendor payload。

### 28.2 页面接入

- Workbench、Positions、Data Quality、Risk Alerts、Rules、Audit、Notifications、Local Install、Local Knowledge、Settings 的新增或改动区域应优先复用 P62 primitives。
- 页面业务逻辑仍归属原 feature/page；primitives 只能处理展示、可访问语义和局部布局，不得读取 SQLite、VecLite、localStorage、sessionStorage、本地文件或临时配置。
- 旧 cockpit class 可以保留，但新增状态、字段、按钮、摘要、表格、空态和错误态不得绕过 P62 可访问性要求。

### 28.3 键盘与 reflow 门禁

- Playwright smoke 必须覆盖主导航、移动菜单、代表性表单、折叠详情区和关键本地按钮的键盘路径。
- 390px、768px、1280px viewport 必须覆盖代表性页面，且 `body.scrollWidth` 与 `documentElement.scrollWidth` 不得超过 viewport。
- 页面级横向滚动禁止；二维表格、JSON、日志或诊断文本只能在明确局部容器滚动。
- 验收必须采集三档 viewport 的截图或等价浏览器证据。

### 28.4 安全边界

- P62 页面和组件不得新增或暗示券商连接、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复、覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。
- 高风险、未知、降级、过期、缺失、信息不足、只读和 blocked 状态不得被展示为普通成功。

验收记录：`docs/release/acceptance/2026-06-18-p62-ui-acceptance.md`。

## 29. P74 知识与数据准备度前端契约

### 29.1 `/data-quality`

API：`GET /api/v1/knowledge-readiness?symbol=...`

数据质量页必须展示“知识与数据准备度”区域：

| UI 区域 | API 字段 | 展示规则 |
| --- | --- | --- |
| 总状态 | `status` | `ready=已准备`、`degraded=降级`、`blocked=阻断`；未知状态展示待检查 |
| 标的画像 | `symbol_profile` | 已知时展示名称、资产类型和跟踪指数；未知时展示不能伪造成 ready |
| 知识引用 | `knowledge_references` | 展示完整内置知识列表，包含大师原则、纪律规则、风险 SOP、标的画像和“是否可作正式证据”边界 |
| 数据依赖矩阵 | `data_dependencies` | 展示完整矩阵，不截断；包含生效规则、正式证据、RAG 索引、LLM 上下文 |
| 影响说明 | `feature_impacts` | 降级或阻断时展示安全影响，不得转成执行 CTA |
| LLM 上下文 | `llm_context_summary` | 只展示是否已附加摘要，不展示原文 |

页面不得直接读取 SQLite、VecLite、本地日志、localStorage、sessionStorage、配置文件或临时证据文件。页面不得出现自动刷新修复、自动确认、自动规则应用、自动交易、一键交易、代下单、券商接口、外部推送或收益承诺入口。

### 29.2 `/decisions/:decisionId`

当 `analyst_reports[].input_summary` 表示已使用知识与数据准备度摘要时，决策详情只展示：

- `LLM 已参考知识与数据准备度摘要`
- 脱敏 prompt version
- “仅展示脱敏摘要”类边界

页面不得展示完整 prompt、`principles=`、`data_readiness=`、持仓上下文原文、密钥、私有路径或 raw provider payload。

### 29.3 验收

- Vitest 必须覆盖 ready、degraded、blocked 三态和决策详情脱敏回显。
- Playwright 必须真实启动本地后端和 Vite 前端，覆盖 `/data-quality`、决策详情、blocked API、390px reflow 和 forbidden affordance scan。
- 验收记录：`docs/release/acceptance/2026-06-19-p74-built-in-knowledge-and-data-readiness.md`。

## 30. 格式化规则

| 数据类型 | 格式 |
| --- | --- |
| 金额 | `¥120,000.00` |
| 比例 | `8.00%` |
| 分位 | `PE 分位 63%` |
| 日期时间 | `2026-05-22 09:30` |
| 信源数量 | `引用 3 个独立信源` |
| 规则版本 | `当前规则版本 v3.0` |
| 内部 ID | 默认折叠，详情中展示 |

## 31. 与其他文档关系

- API 字段来源见 `docs/api.md`。
- Eino 工作流见 `docs/workflow.md`。
- UI 页面布局见 `docs/ui-design.md`。
- 产品体验打磨路线图见 `docs/product-experience-polish-roadmap.md`。
- 架构分层见 `docs/architecture.md`。
