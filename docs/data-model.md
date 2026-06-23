# Investment Agent 数据模型设计

> 文档版本：v1.1
> 最后更新：2026-06-17
> 适用范围：SQLite 表设计、VecLite 索引设计、决策记录、证据引用、用户确认、审计事件、规则进化和状态流转。

## 1. 设计目标

本文档用于把需求、架构、工作流、API 和前端契约中的数据要求固化为可开发的数据模型。

目标：

- SQLite 作为本地事实基准，保存账户、持仓、行情、情报摘要、决策记录、用户确认、规则版本和审计事件。
- VecLite 作为辅助检索索引，可由 SQLite 中的结构化情报摘要和 RAG 文本块重建。
- 每条正式建议必须可追溯到账户快照、持仓快照、市场快照、规则版本、证据引用、分析摘要、裁决链和审计事件。
- 用户确认只记录线下动作，不触发自动交易。
- 规则优化只能生成提案，必须经过用户确认和守门人审计后才可写入正式规则。

## 2. 数据边界

### 2.1 SQLite 事实基准

SQLite 保存所有需要长期追踪、可审计、可恢复的数据。

包括：

- 账户快照、持仓当前态、持仓时点快照。
- 市场行情、估值、流动性和情绪指标快照。
- 用户设置、能力圈配置、数据源配置和规则版本。
- 原始情报、结构化情报摘要、RAG 文本块和多源验证结果。
- 决策记录、证据引用、用户确认、错误案例和审计事件。
- 应用内通知、规则提案、规则效果验证、应用后追踪、守门人审计和规则变更历史。

### 2.2 VecLite 辅助检索索引

VecLite 保存从 `rag_chunks` 生成的向量索引和 BM25 检索索引。

约束：

- VecLite 不作为唯一事实来源。
- VecLite 不保存不可重建的业务事实。
- VecLite 文件不可用或索引版本不兼容时，系统必须从 SQLite 的 `rag_chunks` 和 `intelligence_summary` 重建索引。
- 决策记录中不得只保存 VecLite 向量 ID，必须保存 SQLite 证据引用。

### 2.3 前端访问边界

前端不得访问 SQLite、VecLite 或本地文件。前端只通过 HTTP API 获取 DTO。

### 2.4 P4/P5 事务性写入

P4/P5 写路径必须保持数据模型和 API 契约定义的事务边界。

确认写入规则：

- `POST /api/v1/decisions/{decision_id}/confirmations` 使用 `confirmation_type=executed_manually` 成功时，系统必须原子写入 `operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots` 和 `audit_events`。
- `confirmation_type=marked_error` 成功时，系统必须原子写入 `operation_confirmations`、`error_cases` 和 `audit_events`，响应必须包含 `error_case_id`。
- `confirmation_type=planned` 或 `confirmation_type=watch` 成功时，不得写入 `position_transactions`，不得创建新的账户快照。

工作流写入规则：

- `DecisionRecordNode` 成功时，系统必须原子写入 `decision_records`、`evidence_refs` 和对应 `audit_events`。
- `EvidenceVerificationGraph` 写入事实时，系统必须原子写入 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 和对应 `audit_events`。
- `PublicEvidenceIngestionService` 处理真实公开公告 collector 时也必须写入同一组证据事实表；P29 smoke 使用临时 SQLite 验证 `public-evidence-refresh` 成功后至少生成公开源 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 和成功审计事件。
- `LocalKnowledgeService` 处理 P46 本地知识导入确认时复用同一组证据事实表；默认写入 C 级 background 事实、background_only 核验和 pending RAG 块，不新增 migration。
- `MarketRefreshGraph` 写入事实时，系统必须原子写入 `market_snapshots` 和对应 `audit_events`。
- `EvolutionProposalGraph` 写入事实时，系统必须原子写入 `rule_proposals`、待用户处理的 `notifications` 和对应 `audit_events`。
- `GatekeeperAuditGraph` 写入事实时，系统必须原子写入 `gatekeeper_audits`、`rule_proposals.status` 和对应 `audit_events`。

## 3. 命名规范

领域层必须定义并校验 API、数据模型与工作流中使用的核心枚举；来自契约外的枚举值不得进入正式裁决、状态流转或持久化写入。

| 类型 | 规则 | 示例 |
| --- | --- | --- |
| 表名 | 小写复数，下划线分隔 | `decision_records` |
| 主键 | 业务前缀 + 时间或随机后缀 | `dec_20260522_0001` |
| 外键 | `{entity}_id` | `decision_id` |
| 时间字段 | ISO 8601 字符串或 SQLite datetime | `created_at` |
| 状态字段 | 使用 API 契约枚举 | `confirmation_status` |
| JSON 字段 | 仅保存快照、摘要、非主查询结构 | `payload_json` |

### 3.1 ID 与时间规则

- 工作流和仓储生成的关键实体 ID 必须通过统一 ID 生成入口生成，外部已验证输入可直接作为事实 ID 保存。
- `decision_id`、`evidence_ref_id`、`audit_event_id`、`transaction_id`、规则应用版本等业务 ID 必须非空、可读、可测试。
- 写入 `created_at`、`updated_at`、`executed_at`、审计时间等字段时，默认使用 UTC；跨层传递使用 RFC3339 字符串或 SQLite 可解析的 UTC 时间字符串。

## 4. 核心表设计

### 4.1 账户快照表：`portfolio_snapshots`

用途：保存每次决策读取的账户状态。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `snapshot_id` | TEXT PK | 账户快照 ID |
| `snapshot_time` | DATETIME | 快照生成时间 |
| `cash` | REAL | 现金金额 |
| `total_assets` | REAL | 总资产 |
| `cash_ratio` | REAL | 现金比例 |
| `high_risk_ratio` | REAL | 高风险资产比例 |
| `position_count` | INTEGER | 持仓数量 |
| `source` | TEXT | manual / system |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_portfolio_snapshots_time(snapshot_time)`

### 4.2 当前持仓表：`positions`

用途：保存当前持仓状态和买入逻辑。该表是当前态，不用于复现历史决策时点。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `position_id` | TEXT PK | 持仓 ID |
| `symbol` | TEXT | 标的代码 |
| `name` | TEXT | 标的名称 |
| `quantity` | REAL | 当前数量 |
| `cost_price` | REAL | 成本价 |
| `current_price` | REAL | 当前价格 |
| `market_value` | REAL | 市值 |
| `unrealized_profit_ratio` | REAL | 未实现收益率 |
| `position_state` | TEXT | normal / sell_only / frozen_watch |
| `buy_date` | DATE | 首次买入日期或建仓日期 |
| `buy_reason` | TEXT | 买入理由 |
| `asset_tag` | TEXT | core / satellite / high_risk 等 |
| `updated_at` | DATETIME | 更新时间 |

索引：

- `idx_positions_symbol(symbol)`
- `idx_positions_state(position_state)`

说明：

- 多次买入的精细成本可后续由 `position_lots` 承载；当前文档先保留聚合成本。
- 每次正式建议读取的是 `position_snapshots`，不是直接读取会被后续修改覆盖的当前态。

### 4.3 持仓时点快照表：`position_snapshots`

用途：保存某个账户快照下的持仓明细集合，用于复现 `WorkflowContext.position_snapshots`。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `position_snapshot_id` | TEXT PK | 持仓快照明细 ID |
| `snapshot_id` | TEXT | 账户快照 ID |
| `position_id` | TEXT | 当前持仓 ID，可为空 |
| `symbol` | TEXT | 标的代码 |
| `name` | TEXT | 标的名称 |
| `quantity` | REAL | 快照时数量 |
| `cost_price` | REAL | 快照时成本价 |
| `current_price` | REAL | 快照时价格 |
| `market_value` | REAL | 快照时市值 |
| `unrealized_profit_ratio` | REAL | 快照时未实现收益率 |
| `position_state` | TEXT | normal / sell_only / frozen_watch |
| `buy_date` | DATE | 首次买入日期或建仓日期 |
| `buy_reason` | TEXT | 快照时买入理由 |
| `asset_tag` | TEXT | 资产标签 |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_position_snapshots_snapshot(snapshot_id)`
- `idx_position_snapshots_symbol(symbol)`

约束：

- 一次 `portfolio_snapshots` 可对应多条 `position_snapshots`。
- 每次账户快照必须包含当时完整持仓集合，不得只保存本次变化的持仓。
- `decision_records.portfolio_snapshot_id` 关联的账户快照必须能查到对应持仓快照集合。

### 4.4 操作确认表：`operation_confirmations`

用途：保存用户对建议的处理结果。该表只记录线下动作，不执行交易。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `confirmation_id` | TEXT PK | 确认记录 ID |
| `decision_id` | TEXT | 关联建议 ID |
| `confirmation_type` | TEXT | planned / executed_manually / watch / marked_error |
| `operation_type` | TEXT | buy / sell / reduce，仅 `executed_manually` 使用 |
| `symbol` | TEXT | 标的代码 |
| `quantity` | REAL | 成交数量，仅 executed_manually 必填 |
| `price` | REAL | 成交价格，仅 executed_manually 必填 |
| `executed_at` | DATETIME | 线下执行时间 |
| `error_case_id` | TEXT | 错误案例 ID，可为空 |
| `payload_json` | TEXT | 用户提交原始载荷，可为空 |
| `note` | TEXT | 备注 |
| `created_at` | DATETIME | 创建时间 |

约束：

- `planned`、`watch`、`marked_error` 不更新账户。
- 每次成功确认都创建新的 `operation_confirmations`；`decision_records.confirmation_status` 保存最近一次有效状态。
- `planned` 与 `watch` 可互相转换，也可升级为 `executed_manually` 或 `marked_error`。
- `executed_manually` 与 `marked_error` 是确认终态，再次确认必须返回 `BAD_REQUEST`，不得重复更新账户或错误案例。
- `record_type!=formal_trade_advice` 或 `confirmation_status=not_required` 的记录不得创建确认动作。
- `executed_manually` 的 `operation_type` 只能是 buy / sell / reduce，且必须填写 `quantity`、`price`、`executed_at`。
- `watch` 只能通过 `confirmation_type=watch` 表达，不得作为已执行成交类型。
- 所有确认动作必须写入审计事件。

### 4.5 市场快照表：`market_snapshots`

用途：保存行情、估值、流动性和情绪相关数据，支撑规则裁决和冷静机制审计。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `market_snapshot_id` | TEXT PK | 市场快照 ID |
| `symbol` | TEXT | 标的代码或市场指数 |
| `trade_date` | DATE | 交易日期 |
| `close_price` | REAL | 收盘价 |
| `price_change_ratio` | REAL | 涨跌幅 |
| `volume` | REAL | 成交量 |
| `turnover` | REAL | 成交额 |
| `turnover_rate` | REAL | 换手率 |
| `volatility` | REAL | 波动率 |
| `margin_balance` | REAL | 融资余额 |
| `margin_balance_change` | REAL | 融资余额变化率 |
| `pe` | REAL | PE |
| `pb` | REAL | PB |
| `pe_percentile` | REAL | PE 分位 |
| `pb_percentile` | REAL | PB 分位 |
| `volume_percentile` | REAL | 成交量历史分位 |
| `volatility_percentile` | REAL | 波动率历史分位 |
| `liquidity_state` | TEXT | normal / warning / danger |
| `sentiment_state` | TEXT | cold / neutral / hot / extreme |
| `market_metrics_json` | TEXT | 其他市场指标快照；P27 collector source metadata 存放于此，包括 `source_name`、`source_level`、`source_type`、`trade_date`、`captured_at`、`content_hash` 与 `metadata`；P34/P48 source health 读取 `metadata.p34_source_health` 与 `metadata.p34_data_categories`；P90 capital-flow provider 写入 `metadata.p88_structured_fields.capital_flow.date/net_inflow/net_outflow/raw_net_flow`，其中 `raw_net_flow` 为日净流向原值 |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_market_snapshots_symbol_date(symbol, trade_date)`

约束：

- 情绪极端判定不得只依赖 LLM 文本判断。
- 冷静机制相关裁决必须能从本表字段或 `market_metrics_json` 中复现。
- P27 市场数据 collector 对同一 `source_name + symbol + trade_date + source_type` 的重复刷新必须保持市场事实幂等，不重复写入 `market_snapshots`。

### 4.6 规则配置表：`rule_versions`

用途：保存正式规则版本、阈值和生效状态。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `rule_version` | TEXT PK | 规则版本号 |
| `status` | TEXT | active / archived |
| `rules_json` | TEXT | 根本规则、阈值、SOP 快照 |
| `effective_at` | DATETIME | 生效时间 |
| `created_from_proposal_id` | TEXT | 来源提案 ID，可为空 |
| `created_at` | DATETIME | 创建时间 |

约束：

- 同一时间只能有一个 active 版本。
- 所有决策记录必须保存当时使用的 `rule_version`。

### 4.7 能力圈配置表：`capability_configs`

用途：保存用户允许系统分析的资产类型、行业、地区和策略范围。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `capability_id` | TEXT PK | 能力圈配置 ID |
| `asset_types_json` | TEXT | 资产类型范围 |
| `symbols_json` | TEXT | 明确纳入的标的 |
| `excluded_symbols_json` | TEXT | 明确排除的标的 |
| `strategy_scope_json` | TEXT | 策略范围 |
| `updated_at` | DATETIME | 更新时间 |

### 4.8 用户设置表：`user_settings`

用途：保存设置页中的非规则类配置，以及需要生成规则提案的配置快照。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `settings_id` | TEXT PK | 设置 ID |
| `position_limits_json` | TEXT | 仓位上限配置快照，仅用于展示或提案生成上下文；不得通过 `PUT /api/v1/settings` 直接更新并作为正式裁决规则生效 |
| `cash_min_ratio` | REAL | 最低现金比例配置快照，仅用于展示或提案生成上下文；正式裁决以 `rule_versions.rules_json` 为准 |
| `notification_config_json` | TEXT | 通知配置 |
| `data_sources_json` | TEXT | 数据源配置 |
| `updated_at` | DATETIME | 更新时间 |

约束：

- 修改通知、页面偏好、普通数据源配置可立即保存。
- 修改根本规则、裁决优先级、核心阈值时，必须生成 `rule_proposals`，不得直接修改 active 规则版本。
- `position_limits_json`、`cash_min_ratio` 只作为设置页快照或规则提案上下文，不允许通过 `PUT /api/v1/settings` 直接更新并影响正式裁决；正式裁决必须读取 `rule_versions.rules_json`。

### 4.9 应用内通知表：`notifications`

用途：保存本地应用内通知中心状态，供前端轮询展示和用户标记已读。该表只记录本地提示，不触发邮件、短信、系统 Push、Webhook、WebSocket、交易或规则自动应用。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `notification_id` | TEXT PK | 通知 ID |
| `type` | TEXT | 通知类型，例如 data_source_failure / vector_index_failure / rule_proposal_pending / review_degraded / risk_alert |
| `severity` | TEXT | info / warning / critical |
| `title` | TEXT | 通知标题 |
| `message` | TEXT | 通知正文 |
| `source_type` | TEXT | 关联来源类型，可为空 |
| `source_id` | TEXT | 关联来源 ID，可为空 |
| `read_at` | DATETIME | 已读时间；为空表示未读 |
| `created_at` | DATETIME | 创建或最近刷新时间 |

索引：

- `idx_notifications_created_at(created_at DESC)`
- `idx_notifications_read_at(read_at)`
- `idx_notifications_active_source(type, source_type, source_id) WHERE read_at IS NULL AND source_type IS NOT NULL AND source_id IS NOT NULL`

约束：

- `severity` 只能是 info / warning / critical。
- 同一 `type/source_type/source_id` 的未读通知必须去重；重复写入时刷新严重程度、标题、正文和 `created_at`，避免同一故障刷屏。
- 标记已读只更新 `read_at`，不得修改关联业务事实。
- 通知仅作为本地 UI 提醒，不能作为交易、规则应用或外部推送触发器。

### 4.10 每日自动运行状态表：`daily_auto_run_states`

用途：保存 P31 本地每日自动运行的状态、幂等键和诊断信息，供 server 内 scheduler、API、前端状态展示和审计排查使用。该表只记录本地自动运行状态，不触发交易、外部推送或规则自动应用。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `run_id` | TEXT PK | 自动运行实例 ID |
| `idempotency_key` | TEXT UNIQUE | 幂等键，由本地日期、scope、标的集合 hash 和任务版本组成 |
| `local_date` | TEXT | 按配置时区解释的本地日期 |
| `scope` | TEXT | 运行范围；P31 仅支持 `holdings` |
| `symbol_set_hash` | TEXT | 本次持仓标的集合 hash |
| `status` | TEXT | disabled / scheduled / running / success / degraded / failed |
| `last_run_at` | DATETIME | 最近实际运行时间，可为空 |
| `next_run_at` | DATETIME | 下一次计划运行时间，可为空 |
| `failure_code` | TEXT | 失败或降级错误码，可为空 |
| `failure_reason` | TEXT | 可读诊断说明，可为空 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

索引：

- `idx_daily_auto_run_states_date_scope(local_date, scope)`

约束：

- `scope` 只能是 `holdings`。
- `status` 只能是 disabled / scheduled / running / success / degraded / failed。
- failed 状态必须记录 `failure_code`。
- 同一 `idempotency_key` 只能保留一个状态记录，重复运行应更新原记录或显式记录 retry/rerun 审计，不得生成冲突的每日状态。
- 状态表不得作为交易执行、规则自动应用或外部通知触发器。

### 4.11 每日纪律报告索引表：`daily_discipline_reports`

用途：保存 P32 每日纪律报告的本地日期、持仓 scope、报告状态、关联决策/审计/通知线索和缺前提诊断，供今日报告、历史报告列表和详情回看使用。该表是阅读与追踪索引，不替代 `decision_records` 的正式裁决事实源，不触发交易、外部推送或规则自动应用。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `report_id` | TEXT PK | 稳定报告 ID |
| `local_date` | TEXT | 按配置时区解释的本地日期 |
| `scope` | TEXT | 报告范围；P32 仅支持 `holdings` |
| `symbol_set_hash` | TEXT | 本次持仓标的集合 hash |
| `source_type` | TEXT | auto_run / manual |
| `source_id` | TEXT | daily auto-run 幂等键或手动任务 request id，可为空 |
| `decision_id` | TEXT | 关联 `decision_records`，可为空 |
| `status` | TEXT | not_started / running / success / degraded / failed / insufficient_data |
| `summary` | TEXT | 报告摘要，可为空 |
| `failure_code` | TEXT | 失败、降级或缺前提错误码，可为空 |
| `failure_reason` | TEXT | 用户可读诊断，可为空 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

索引：

- `idx_daily_discipline_reports_date_scope(local_date, scope)`
- `idx_daily_discipline_reports_updated_at(updated_at DESC)`
- `idx_daily_discipline_reports_status(status)`

约束：

- `scope` 只能是 `holdings`。
- `source_type` 只能是 auto_run / manual。
- `status` 只能是 not_started / running / success / degraded / failed / insufficient_data。
- 同一 `local_date + scope + symbol_set_hash` 只能保留一个当前报告索引，重复聚合应复用或更新原报告，不得生成冲突的同日同 scope 报告。
- `failed` 或 `insufficient_data` 状态应记录 `failure_code` 或 `failure_reason`，并不得伪造摘要、证据、预期收益或交易指令。
- 报告索引只保存每日阅读入口和关联线索；正式裁决、证据和审计仍以 `decision_records`、`evidence_refs`、`source_verifications`、`audit_events` 为准。

### 4.12 风险预警事实表：`risk_alerts`

用途：保存 P35 本地风险预警与 SOP 状态。该表只作为本地人工复核和追踪事实，不代表订单、券商指令、账户变更、组合变更或正式规则版本变更。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `alert_id` | TEXT PK | 风险预警 ID |
| `risk_type` | TEXT | valuation_high / buy_thesis_broken / liquidity_danger / sentiment_extreme / position_limit_breach / insufficient_evidence / data_degraded |
| `severity` | TEXT | info / warning / critical |
| `sop_status` | TEXT | triggered / active / observing / escalated / resolved / archived |
| `symbol` | TEXT | 影响标的 |
| `trigger_summary` | TEXT | 触发摘要 |
| `trigger_context_json` | TEXT | 触发上下文 JSON，可包含估值、流动性、情绪、证据或 source health |
| `prohibited_actions_json` | TEXT | 禁止动作列表 JSON |
| `suggested_actions_json` | TEXT | 建议人工动作列表 JSON |
| `related_decision_id` | TEXT | 关联决策 ID，可为空 |
| `related_report_id` | TEXT | 关联每日纪律报告 ID，可为空 |
| `related_notification_id` | TEXT | 关联通知 ID，可为空 |
| `related_audit_event_id` | TEXT | 关联审计事件 ID，可为空 |
| `last_triggered_at` | DATETIME | 最近触发时间 |
| `resolved_at` | DATETIME | 解除或归档时间 |
| `resolution_reason` | TEXT | 解除、归档或状态变更原因 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

索引：

- `idx_risk_alerts_active_identity(risk_type, symbol) WHERE sop_status IN ('triggered','active','observing','escalated')`
- `idx_risk_alerts_status_updated_at(sop_status, updated_at DESC)`
- `idx_risk_alerts_symbol_status(symbol, sop_status)`

约束：

- 同一 active identity（risk_type + symbol）只保留一个 triggered / active / observing / escalated 风险事实，重复触发更新原记录并保留可追溯审计。
- 终态 `resolved` / `archived` 不得转回非终态。
- active / escalated 风险可生成本地 `notifications`，并写入 `audit_events`。
- 风险事实写入不得修改 `positions`、`portfolio_snapshots`、`operation_confirmations`、`position_transactions`、`rule_versions`、broker state、orders 或 external notifications。

### 4.13 原始情报表：`intelligence_items`

用途：保存从外部信源采集并清洗后的情报条目。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `intelligence_id` | TEXT PK | 情报 ID |
| `source_name` | TEXT | 来源名称 |
| `source_level` | TEXT | S / A / B / C |
| `original_url` | TEXT | 原文链接 |
| `published_at` | DATETIME | 发布时间 |
| `captured_at` | DATETIME | 抓取时间 |
| `content_hash` | TEXT | 内容 hash |
| `raw_title` | TEXT | 原始标题 |
| `raw_text_ref` | TEXT | 原文保存引用，可为空 |
| `created_at` | DATETIME | 创建时间 |

约束：

- C 级信源默认不进入正式证据链。
- 无来源信息的情报必须丢弃。
- P29 真实公开公告 collector 写入时，`source_name/source_level/original_url/published_at/captured_at/content_hash/raw_title` 必须足以追溯公开源记录；重复采集同一 source record 不得重复生成事实。
- P46 本地知识导入写入时，`source_name` 使用用户确认的 `source_label`，`source_level=C`，`original_url` 为空，`content_hash` 由规范化内容生成；同一 import batch 内不得重复生成事实。

### 4.12 结构化情报摘要表：`intelligence_summary`

用途：保存可检索、可排序、可用于证据链的结构化摘要。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `summary_id` | TEXT PK | 摘要 ID |
| `intelligence_id` | TEXT | 原始情报 ID |
| `symbol` | TEXT | 关联标的 |
| `entity` | TEXT | 实体 |
| `event_type` | TEXT | 事件类型 |
| `impact_direction` | TEXT | positive / negative / neutral |
| `summary` | TEXT | 摘要文本 |
| `source_level` | TEXT | S / A / B / C |
| `evidence_role` | TEXT | formal / background |
| `time_weight` | REAL | 时效权重（0–1） |
| `relevance_score` | REAL | 检索相关度（0–1） |
| `independent_source_count` | INTEGER | 独立信源数量 |
| `high_grade_independent_source_count` | INTEGER | S/A 级独立信源数量 |
| `source_name` | TEXT | 信源名称 |
| `original_url` | TEXT | 原文链接 |
| `published_at` | DATETIME | 发布时间 |
| `captured_at` | DATETIME | 抓取时间 |
| `content_hash` | TEXT | 内容 hash |
| `verification_group_id` | TEXT | 多源验证分组 |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_intelligence_summary_symbol(symbol)`
- `idx_intelligence_summary_group(verification_group_id)`

约束：

- `source_level=C` 只能作为 background，不得作为正式裁决依据。
- `evidence_role=formal` 的普通证据必须满足 S/A/B 信源要求。
- 重大利好、重大利空、买入逻辑破坏等重大事件对应的 `evidence_role=formal` 必须满足至少 2 个 A 或 S 级独立信源。
- 真实公开公告 collector 的 CNInfo A 级公告可作为 formal 证据入库；深交所或证监会不可用、无数据或解析失败时，只能通过错误审计或 background/degraded 语义表达，不得伪造成已满足的高等级独立信源。
- P46 本地知识导入生成的摘要必须固定 `source_level=C`、`evidence_role=background`，`independent_source_count=1`，`high_grade_independent_source_count=0`；不得提升为正式证据。

### 4.13 RAG 文本块表：`rag_chunks`

用途：保存可重建 VecLite 的文本块和索引元信息。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `chunk_id` | TEXT PK | 文本块 ID |
| `summary_id` | TEXT | 摘要 ID |
| `chunk_text` | TEXT | 文本块内容 |
| `chunk_hash` | TEXT | 文本块 hash |
| `vector_id` | TEXT | VecLite 向量 ID |
| `vector_collection` | TEXT | VecLite 集合或文件内分组 |
| `embedding_model` | TEXT | 嵌入模型 |
| `embedding_version` | TEXT | 嵌入模型版本 |
| `index_version` | TEXT | 索引版本 |
| `index_status` | TEXT | pending / indexed / stale / failed |
| `indexed_at` | DATETIME | 最近索引时间 |
| `metadata_json` | TEXT | source_level、symbol、published_at 等 |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_rag_chunks_summary(summary_id)`
- `idx_rag_chunks_vector(vector_id)`
- `idx_rag_chunks_status(index_status)`

约束：

- 新建情报摘要时，相关文本块先写为 `pending`；VecLite 写入成功后更新为 `indexed`，写入失败后更新为 `failed`，且不得回滚已保存的 SQLite 情报事实。
- 当 `chunk_hash`、`embedding_version` 或 `index_version` 变化时，`index_status` 必须变为 stale 或 pending。
- VecLite 重建只依赖本表与 `intelligence_summary`，不得依赖前端状态。
- P46 本地知识导入只写 `pending` 文本块，确认响应只返回索引计划；实际索引状态仍由后续重建流程更新。

### 4.14 交易流水表：`position_transactions`

用途：记录每次人工执行的买卖、减仓、加仓流水，用于重建持仓变化和成本变化。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `transaction_id` | TEXT PK | 交易流水 ID |
| `confirmation_id` | TEXT | 关联确认记录 ID |
| `symbol` | TEXT | 标的代码 |
| `operation_type` | TEXT | buy / sell / reduce |
| `quantity` | REAL | 数量 |
| `price` | REAL | 成交价格 |
| `fees` | REAL | 费用，可为空 |
| `occurred_at` | DATETIME | 实际发生时间 |
| `before_position_json` | TEXT | 变更前持仓快照 |
| `after_position_json` | TEXT | 变更后持仓快照 |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_position_transactions_confirmation(confirmation_id)`
- `idx_position_transactions_symbol(symbol)`

约束：

- `executed_manually` 且 `operation_type` 为 buy / sell / reduce 时，必须至少写入一条 `position_transactions`。
- `hold`、`watch`、`planned`、`marked_error` 不生成交易流水。
- 交易流水是持仓变化的事实记录，当前持仓表只是聚合结果。

### 4.14.1 P47 决策闭环解释只读查询边界

P47 不新增表、索引或 migration。决策闭环解释只复用现有事实表，并通过只读 repository 查询聚合：

- `decision_records`：读取建议、标的、最终裁决、当前确认状态、工作流状态和生成时间。
- `operation_confirmations`：按 `decision_id` 与 `created_at ASC` 读取用户记录。
- `position_transactions`：按 `confirmation_id` 与 `occurred_at ASC` 读取线下本地流水。
- `error_cases`：读取用户标记错误形成的复盘样本。
- `risk_alerts`：读取与 `decision_id` 或标的相关的风险线索。
- `audit_events`：读取与决策、确认或错误案例相关的审计线索。

P47 聚合层不得返回 `operation_confirmations.payload_json`、`position_transactions.before_position_json`、`position_transactions.after_position_json` 或任何私有路径、完整 key、原始 SQL、完整 prompt、供应商原始响应。`note` 只能作为脱敏后的 `note_preview` 返回。P47 API 不写入上述表，不更新确认状态，不创建通知，不改变风险 SOP，不修改规则版本。

### 4.14.2 P48 数据源质量回归边界

P48 不新增表、索引或 migration。数据源质量回归复用现有 `market_snapshots.market_metrics_json` 中的 P34 source health metadata：

- `metadata.p34_source_health`：按数据类别记录 freshness、last success/failure、failure category、data date 和 affected symbols。
- `metadata.p34_data_categories`：记录本次快照可解释的数据类别顺序。

`GET /api/v1/data-source-quality/regression` 在 `fixture` 模式下只使用内置确定性样本，不读取用户私有事实；在 `current` 模式下只读 `market_snapshots` 最新快照，不触发 collector、不写市场快照、不创建通知、不修改账户、确认、风险 SOP 或规则。`cmd/agent --task data-source-quality-regression` 仅允许通过既有 `audit_events` 写入脱敏任务摘要，`output_ref` 只保存 mode、status、case/degraded/failed 计数和安全边界，不保存 raw payload、完整 key、私有路径、原始 SQL、完整 prompt、raw HTTP、private key 或供应商原始响应。

### 4.14.3 当前数据门禁处置表：`data_quality_gate_resolutions`

用途：保存 P67 本地人工处置记录，用于解释 P66 current data policy gate 的 release claim state。该表只记录人工声明边界，不修改 source health、P66 policy、市场快照、账户、确认、风险 SOP、通知或规则版本。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `resolution_id` | TEXT PK | 处置记录 ID |
| `symbol` | TEXT | 当前 policy 关联标的 |
| `policy_fingerprint` | TEXT | 基于 symbol、policy verdict/gate、blocking/waiver reasons、case category/status/failure category 的 canonical hash |
| `policy_verdict` | TEXT | 创建时复制的 P66 policy verdict |
| `release_gate` | TEXT | 创建时复制的 P66 release gate |
| `policy_summary` | TEXT | 脱敏 policy 摘要 |
| `resolution_type` | TEXT | `waiver` / `scope_exclusion` |
| `status` | TEXT | `active` / `retired` |
| `scope` | TEXT | 处置范围 |
| `reason` | TEXT | 人工处置原因 |
| `release_impact` | TEXT | 发布声明影响 |
| `evidence_ref` | TEXT | 可选本地证据引用 |
| `blocking_reasons_json` | TEXT | 创建时复制的 blocking reasons |
| `waiver_reasons_json` | TEXT | 创建时复制的 waiver reasons |
| `created_by` | TEXT | 创建者，默认 `local_user` |
| `retired_by` | TEXT | 退役者 |
| `created_at` | DATETIME | 创建时间 |
| `retired_at` | DATETIME | 退役时间 |
| `safety_note` | TEXT | 固定安全边界说明 |

索引与约束：

- `idx_data_quality_gate_resolutions_active_policy`：`UNIQUE(symbol, policy_fingerprint) WHERE status='active'`。
- 同一 `symbol + policy_fingerprint` 只允许一个 active resolution；同类型重复请求复用 active record，不同类型 active 请求返回冲突。
- `policy=blocked` 只允许 `scope_exclusion`；`policy=waiver_required` 允许 `waiver` 或 `scope_exclusion`；`policy=passed` 不允许创建 resolution。
- `scope`、`reason`、`release_impact` 必填并脱敏。
- `POST` create/retire 成功必须写入脱敏 `audit_events`；`GET` check/list 不写审计。

### 4.15 本地账户导入批次表：`local_account_import_batches`

用途：保存批量导入校验与确认的批次 metadata。该表不替代账户快照、持仓、交易流水或审计事实。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `import_batch_id` | TEXT PK | 导入批次 ID |
| `request_id` | TEXT | 校验请求 ID |
| `status` | TEXT | validated / committed / rejected |
| `row_count` | INTEGER | 总行数 |
| `valid_count` | INTEGER | 有效行数 |
| `invalid_count` | INTEGER | 无效行数 |
| `validation_summary_json` | TEXT | 校验摘要 |
| `rows_hash` | TEXT | 已校验 rows 的内容 hash，确认时必须匹配 |
| `created_at` | DATETIME | 创建时间 |
| `committed_at` | DATETIME | 确认写入时间，可为空 |

索引：

- `idx_local_account_import_batches_status(status)`
- `idx_local_account_import_batches_created_at(created_at DESC)`

约束：

- `rows_hash` 必须非空。
- 只有 `status=validated` 且 `invalid_count=0` 的批次允许确认写入。
- 确认写入时必须重新计算 rows hash 并与本表记录一致。
- 校验阶段不得写入账户、持仓、交易或快照事实。

### 4.16 本地账户修正表：`local_account_corrections`

用途：记录用户对本地账户、持仓、交易或导入批次输入错误的修正审计。该表只表达修正事实，不静默改写历史。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `correction_id` | TEXT PK | 修正 ID |
| `target_type` | TEXT | portfolio_snapshot / position / position_snapshot / position_transaction / import_batch |
| `target_id` | TEXT | 目标事实 ID |
| `before_json` | TEXT | 修正前引用或摘要 |
| `after_json` | TEXT | 修正后引用或摘要 |
| `correction_reason` | TEXT | 修正原因 |
| `snapshot_id` | TEXT | 可选关联快照 ID |
| `audit_event_id` | TEXT | 关联审计事件 ID |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_local_account_corrections_target(target_type, target_id)`
- `idx_local_account_corrections_created_at(created_at DESC)`

约束：

- 修正记录不得物理删除或覆盖历史快照、交易流水或审计事件。
- 若修正需要改变当前持仓、现金或快照，必须通过持仓编辑或线下交易记录生成新的快照状态。

## 5. 决策与审计表设计

### 5.1 决策记录表：`decision_records`

用途：保存每次正式建议的核心结果。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `decision_id` | TEXT PK | 建议 ID |
| `request_id` | TEXT | 请求 ID |
| `workflow_type` | TEXT | daily_discipline / consultation 等 |
| `symbol` | TEXT | 标的，可为空 |
| `question` | TEXT | 用户问题，可为空 |
| `workflow_status` | TEXT | completed / degraded / failed |
| `record_type` | TEXT | formal_trade_advice / non_trade_record / rejection_record |
| `dashboard_state` | TEXT | first_use / normal / insufficient_data / frozen_watch / high_risk |
| `capability_status` | TEXT | in_scope / out_of_scope / unknown，可为空 |
| `capability_reason` | TEXT | 能力圈检查说明，可为空 |
| `source_verification_status` | TEXT | satisfied / failed / background_only，可为空 |
| `risk_reason_code` | TEXT | 高危状态原因码，可为空 |
| `media_heat_summary_json` | TEXT | 媒体热度摘要 |
| `user_emotion_tags_json` | TEXT | 用户文本情绪标签 |
| `triggered_rules_json` | TEXT | 规则命中结构 |
| `errors_json` | TEXT | 节点错误列表 |
| `final_verdict_status` | TEXT | buy_allowed / hold / reduce / sell_only / frozen_watch / rejected / insufficient_data |
| `final_verdict_text` | TEXT | 最终裁决文案 |
| `prohibited_actions_json` | TEXT | 禁止事项 |
| `optional_actions_json` | TEXT | 可选事项 |
| `confirmation_status` | TEXT | not_required / pending / planned / executed_manually / watch / marked_error |
| `portfolio_snapshot_id` | TEXT | 账户快照 ID |
| `market_snapshot_id` | TEXT | 市场快照 ID |
| `rule_version` | TEXT | 规则版本 |
| `analyst_reports_json` | TEXT | 多 Agent 分析摘要 |
| `expected_return_scenarios_json` | TEXT | 预期收益情景快照，只用于情景概率、动态卖出评估和复核触发展示，不参与最终裁决覆盖 |
| `arbitration_chain_json` | TEXT | 裁决链 |
| `context_snapshot_json` | TEXT | WorkflowContext 的关键可复现快照 |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_decision_records_created_at(created_at)`
- `idx_decision_records_symbol(symbol)`
- `idx_decision_records_status(confirmation_status)`
- `idx_decision_records_record_type(record_type)`
- `idx_decision_records_dashboard_state(dashboard_state)`

约束：

- `DecisionRecordNode` 保存失败时，前端不得展示为正式建议。
- `record_type=formal_trade_advice` 表示正式交易类建议，必须关联 `request_id`、`rule_version`、`portfolio_snapshot_id`、审计事件和可用证据引用；交易类正式建议必须关联 `market_snapshot_id`。
- `record_type=rejection_record` 只用于能力圈外、规则拒绝等拒绝型记录，`final_verdict_status=rejected`，`confirmation_status=not_required`，不得包含交易确认动作。
- `record_type=non_trade_record` 只用于 `insufficient_data`、数据过期、证据不足等非交易型复现记录，`confirmation_status=not_required`，可不关联 `market_snapshot_id`，但必须在 `errors_json` 写明原因。
- 首次使用、缺少账户或缺少持仓时不得创建 `decision_records`，只允许返回 `DATA_REQUIRED` 或 `dashboard_state=first_use` 并写审计事件。
- `dashboard_state` 保存当时页面主状态，避免历史页面因后续实时计算而变化。
- `triggered_rules_json` 应保存规则命中结构，用于驾驶舱和详情页展示。
- `expected_return_scenarios_json` 必须保存当次 `expected_return_scenarios`，用于历史详情复现；该字段不得覆盖最终裁决字段或规则裁决结果。
- P28 后 `expected_return_scenarios_json` 的 JSON shape 包含 `precision_status`、`reason`、`sample_count`、`sample_window`、`screening_condition`、`scenarios[].trigger`、`sell_evaluation.status/triggers/prompts/actions/non_trading_disclaimer` 与 `reassessment_trigger.reason/boundary/current_value`。
- `sell_evaluation` 和 `reassessment_trigger` 只作为人工复核材料保存，不得驱动自动交易、账户更新、确认记录、通知或外部推送；旧 JSON 缺少这些字段时按空/不适用回放。

### 5.2 证据引用表：`evidence_refs`

用途：保存建议引用过哪些证据。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `evidence_ref_id` | TEXT PK | 内部引用 ID |
| `evidence_id` | TEXT | 对外展示证据编号 |
| `decision_id` | TEXT | 建议 ID |
| `summary_id` | TEXT | 情报摘要 ID |
| `source_name` | TEXT | 来源名称 |
| `source_level` | TEXT | S / A / B / C |
| `evidence_role` | TEXT | formal / background |
| `published_at` | DATETIME | 发布时间 |
| `captured_at` | DATETIME | 抓取时间 |
| `original_url` | TEXT | 原文链接 |
| `summary` | TEXT | 展示摘要 |
| `content_hash` | TEXT | 内容 hash |
| `time_weight` | REAL | 时效权重（0–1） |
| `relevance_score` | REAL | 相关度（0–1） |
| `independent_source_count` | INTEGER | 独立信源数量 |
| `high_grade_independent_source_count` | INTEGER | S/A 级独立信源数量 |
| `created_at` | DATETIME | 创建时间 |

约束：

- `evidence_ref_id` 用于内部关联，`evidence_id` 用于 API 和前端展示。
- `evidence_role=formal` 不允许使用 C 级信源。
- C 级信源如被保留，只能作为背景材料或审计说明，不进入正式裁决链。
- 前端默认折叠 hash 和内部 ID，详情中可展开。

索引：

- `idx_evidence_refs_decision(decision_id)`
- `idx_evidence_refs_evidence(evidence_id)`

### 5.3 多源验证表：`source_verifications`

用途：保存同一事件或标的在证据核查后的多源验证结果，供 EvidenceVerificationGraph、决策详情和证据页使用。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `verification_id` | TEXT PK | 验证记录 ID |
| `verification_group_id` | TEXT | 多源验证分组 ID |
| `event_id` | TEXT | 事件 ID |
| `symbol` | TEXT | 标的代码，可为空 |
| `event_type` | TEXT | 事件类型 |
| `evidence_role` | TEXT | formal / background |
| `verification_status` | TEXT | satisfied / failed / background_only |
| `independent_source_count` | INTEGER | 独立信源数量 |
| `high_grade_independent_source_count` | INTEGER | S/A 级独立信源数量 |
| `highest_source_level` | TEXT | 最高信源等级 |
| `latest_published_at` | DATETIME | 最新发布时间 |
| `evidence_ids_json` | TEXT | 关联证据编号列表 |
| `created_at` | DATETIME | 创建时间 |

索引：

- `idx_source_verifications_group(verification_group_id)`
- `idx_source_verifications_symbol(symbol)`
- `idx_source_verifications_symbol_event(symbol, event_id)`

约束：

- 涉及买入逻辑破坏、重大利好或重大利空的信息，必须满足至少 2 个 A 或 S 级独立信源确认。
- 普通正式证据允许 S/A/B 级来源进入 `formal` 证据链；C 级只能作为 `background`。
- `verification_status=failed` 时，RuleArbitrationNode 必须暂停交易类建议。
- `event_id` 是业务事件 ID，查询时优先与 `symbol` 组合使用。
- C 级信源只能以 `evidence_role=background` 参与核查和展示。
- 后续单源 `no_data`、`source_unavailable` 或 `parse_error` 不得覆盖已满足的多源验证；这些状态应通过 `audit_events.status/error_code` 或降级记录表达。
- P46 本地知识导入写入的核验记录必须为 `evidence_role=background`、`verification_status=background_only`，用于检索上下文和审计追踪，不得满足正式多源验证门槛。

### 5.4 审计事件表：`audit_events`

用途：记录节点执行、用户确认、错误标注、规则提案、守门人审计等事件。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `audit_event_id` | TEXT PK | 审计事件 ID |
| `request_id` | TEXT | 请求 ID |
| `decision_id` | TEXT | 关联建议 ID，可为空 |
| `workflow_type` | TEXT | 工作流类型 |
| `node_name` | TEXT | 节点名称 |
| `actor` | TEXT | system / user / gatekeeper |
| `action` | TEXT | 审计动作枚举，见下方说明 |
| `node_action` | TEXT | 节点级动作，可为空 |
| `proposal_id` | TEXT | 规则提案 ID，可为空 |
| `confirmation_id` | TEXT | 确认记录 ID，可为空 |
| `error_case_id` | TEXT | 错误案例 ID，可为空 |
| `status` | TEXT | success / degraded / failed |
| `error_code` | TEXT | 错误码，可为空 |
| `before_state` | TEXT | 操作前状态 |
| `after_state` | TEXT | 操作后状态 |
| `rule_version` | TEXT | 规则版本 |
| `snapshot_id` | TEXT | 账户快照 ID，可为空 |
| `input_ref_type` | TEXT | 输入引用类型 |
| `input_ref` | TEXT | 输入摘要或引用 ID |
| `output_ref_type` | TEXT | 输出引用类型 |
| `output_ref` | TEXT | 输出记录 ID |
| `created_at` | DATETIME | 创建时间 |

说明：API 响应同时返回 `audit_event_id` 与 `event_id`，两者均映射为本表 `audit_event_id`；`event_id` 仅作为前端兼容别名。

`action` 可选值：

- `generate_decision`
- `confirm_operation`
- `mark_error`
- `create_proposal`
- `audit_rule_change`
- `update_rule`
- `refresh_market_data`
- `update_settings`
- `update_capability`
- `rebuild_index`
- `run_local_task`
- `risk_alert`

`node_action` 可选值：

- `load_state_snapshot`
- `check_capability`
- `retrieve_evidence`
- `verify_source`
- `run_value_analyst`
- `run_trend_risk_officer`
- `estimate_expected_return`
- `arbitrate_rule`
- `record_decision`
- `create_rule_proposal`
- `audit_rule_change`
- `apply_rule_version`
- `confirm_operation`
- `mark_error_case`
- `refresh_evidence`
- `refresh_market_data`
- `rebuild_vector_index`
- `degrade_workflow`

索引：

- `idx_audit_events_request(request_id)`
- `idx_audit_events_decision(decision_id)`
- `idx_audit_events_created(created_at)`

约束：

- 每个 Eino 节点完成、降级或失败时必须写入审计事件。
- `input_ref_type` 和 `output_ref_type` 必须指明引用对象，例如 `workflow_context`、`decision_record`、`evidence_ref`、`rule_proposal`。
- `action` 是业务动作，必须填写；`node_name` 和 `node_action` 是工作流节点字段，工作流节点事件必须填写，纯用户动作可为空。
- `status=failed` 时必须填写 `error_code`；`status=degraded` 时如有明确降级原因也应填写 `error_code`。
- 公开证据真实 collector 失败审计必须区分 `no_data`、`source_unavailable`、`parse_error`：`no_data` 代表源可达但窗口无记录，不应被解释为系统不可用；`source_unavailable` 与 `parse_error` 可触发降级或后续修复任务，但不得创建交易、确认或外部通知。
- `public-evidence-refresh` 成功审计应记录输入标的和写入数量，失败审计应记录 source-specific 错误码，便于复盘真实采集是否完成。
- 状态机、设置变更、规则变更、能力圈变更必须填写 `before_state` 与 `after_state`。
- `input_ref`、`output_ref`、`before_state`、`after_state` 不得保存密钥明文或未脱敏本地路径。

## 6. 复盘与规则进化表设计

### 6.1 错误案例表：`error_cases`

用途：保存用户标记错误后的事实结果和经验记录。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `error_case_id` | TEXT PK | 错误案例 ID |
| `decision_id` | TEXT | 关联建议 ID |
| `confirmation_id` | TEXT | 关联确认记录 ID |
| `actual_outcome` | TEXT | 实际结果 |
| `root_cause_tag` | TEXT | evidence_missed / rule_threshold_issue / analyst_error / user_context_missing / market_exception |
| `lesson_learned` | TEXT | 经验记录 |
| `created_at` | DATETIME | 创建时间 |

约束：

- `marked_error` 必须在同一事务中写入 `operation_confirmations`、`error_cases` 和 `audit_events`。
- 标记错误请求中的 `note` 保存在 `operation_confirmations.note`，原始请求保存在 `operation_confirmations.payload_json`。

### 6.2 规则提案表：`rule_proposals`

用途：保存由错误案例生成的规则优化提案。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `proposal_id` | TEXT PK | 提案 ID |
| `proposal_type` | TEXT | threshold / sop / risk_rule / capability |
| `status` | TEXT | draft / pending_user_confirm / under_gatekeeper_audit / pending_final_confirm / rejected / applied |
| `source_error_case_id` | TEXT | 主来源错误案例 ID |
| `title` | TEXT | 提案标题 |
| `proposal_version` | TEXT | 提案版本，例如 v3.1-proposal |
| `before_rule_json` | TEXT | 修改前规则 |
| `after_rule_json` | TEXT | 修改后规则 |
| `reason` | TEXT | 提案原因 |
| `impact_scope_json` | TEXT | 影响范围，例如 consultation / daily_discipline |
| `risk_notes_json` | TEXT | 风险提示 |
| `sample_count` | INTEGER | 相关样本数量 |
| `final_confirmed_at` | DATETIME | 用户最终确认时间，可为空 |
| `final_confirmed_note` | TEXT | 用户最终确认备注，可为空 |
| `applied_rule_version` | TEXT | 实际应用后的规则版本，可为空 |
| `related_error_cases_json` | TEXT | 相关错误案例 |
| `created_at` | DATETIME | 创建时间 |

约束：

- 样本数量少于 3 个时，提案状态只能是 `draft` 或 `pending_user_confirm`，并必须在 `risk_notes_json` 中写入样本不足说明。
- 提案不直接修改 `rule_versions`。
- 主成功路径为 `draft -> pending_user_confirm -> under_gatekeeper_audit -> pending_final_confirm -> applied`；完整状态机以 `docs/functional-spec.md` 与 `docs/api.md` 的规则提案状态机表为准。
- 用户放弃送审、守门人审计拒绝、最终拒绝应用均进入 `rejected`，且不得写入 `rule_versions`。
- `rejected` 与 `applied` 是终态，终态重复确认必须返回 `BAD_REQUEST`，只允许写入失败或拒绝类 `audit_events`，不得修改提案和规则版本。
- 守门人审计结果为 `needs_user_review` 时，提案回到 `pending_user_confirm`；用户可放弃或重新送审；如需修改，由后续 EvolutionGraph 或受控内部任务生成新提案版本。
- `approved`、`rejected`、`needs_user_review` 只作为 `gatekeeper_audits.audit_result`，其中 `approved` 不得作为提案状态。
- `pending_final_confirm` 表示守门人审计通过，但还未写入正式规则。
- P36 规则效果验证不通过、过拟合高、历史回放不利或验证版本与提案版本不一致时，提案不得进入最终应用；缺失验证必须展示为 `not_evaluated`，不得假定安全。

### 6.3 规则效果验证事实表：`rule_effect_validations`

用途：保存规则提案的本地效果验证事实，解释提案来源、样本代表性、过拟合风险、历史回放结果和门禁结论。该表只服务规则治理、守门人审计和前端展示，不代表正式规则、订单、券商指令或账户变更。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `validation_id` | TEXT PK | 验证事实 ID |
| `proposal_id` | TEXT | 规则提案 ID |
| `candidate_rule_version` | TEXT | 候选规则版本 |
| `validation_status` | TEXT | not_evaluated / insufficient / passed / failed / needs_more_samples / needs_user_review |
| `sample_count` | INTEGER | 样本数量 |
| `sample_window` | TEXT | 样本窗口 |
| `representativeness_status` | TEXT | 样本代表性状态，取值同 validation_status |
| `overfit_risk` | TEXT | low / medium / high |
| `replay_result` | TEXT | passed / failed / mixed / unknown |
| `guardrail_decision` | TEXT | passed / rejected / needs_user_review |
| `source_explanation_json` | TEXT | 来源解释，包含错误案例、决策、风险预警等本地事实线索 |
| `metrics_json` | TEXT | 命中、误判、缺证据、降级、风险预警等指标快照 |
| `risk_notes_json` | TEXT | 风险说明 |
| `related_error_cases_json` | TEXT | 关联错误案例 ID |
| `related_decision_ids_json` | TEXT | 关联决策 ID |
| `related_risk_alert_ids_json` | TEXT | 关联风险预警 ID |
| `related_audit_event_ids_json` | TEXT | 关联审计事件 ID |
| `safety_note` | TEXT | 安全文案，说明不会自动应用规则或执行交易 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

索引：

- `idx_rule_effect_validations_proposal(proposal_id, updated_at DESC)`
- `idx_rule_effect_validations_rule_version(candidate_rule_version, updated_at DESC)`
- `idx_rule_effect_validations_status(validation_status, updated_at DESC)`

约束：

- 生成或刷新验证必须写入对应 `audit_events`。
- `validation_status=insufficient/failed/needs_more_samples`、`overfit_risk=high`、`replay_result=failed` 或 `guardrail_decision=rejected` 不得支持最终应用规则。
- 该表不得触发或修改 `rule_versions`、`operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、broker state、orders 或 external notifications。
- source、metrics 和 related JSON 可以为空或降级，但缺失事实必须显式表达，不得伪造成通过验证。

### 6.4 规则应用后追踪事实表：`rule_effect_tracking`

用途：保存已应用规则版本在后续复盘周期中的命中、误判、缺证据、降级、风险预警和趋势事实。该表只读展示，不自动回滚或替换正式规则。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `tracking_id` | TEXT PK | 追踪事实 ID |
| `applied_rule_version` | TEXT | 已应用规则版本 |
| `proposal_id` | TEXT | 来源提案 ID，可为空 |
| `period` | TEXT | 追踪周期 |
| `hit_count` | INTEGER | 规则命中次数 |
| `misjudgment_count` | INTEGER | 误判次数 |
| `missing_evidence_count` | INTEGER | 缺证据次数 |
| `degraded_count` | INTEGER | 降级次数 |
| `risk_alert_count` | INTEGER | 关联风险预警次数 |
| `trend_direction` | TEXT | improved / flat / worsened / unknown |
| `metrics_json` | TEXT | 指标快照 |
| `related_proposal_ids_json` | TEXT | 关联提案 ID |
| `related_audit_event_ids_json` | TEXT | 关联审计事件 ID |
| `related_risk_alert_ids_json` | TEXT | 关联风险预警 ID |
| `safety_note` | TEXT | 安全文案，说明不会自动回滚规则或执行交易 |
| `created_at` | DATETIME | 创建时间 |
| `updated_at` | DATETIME | 更新时间 |

索引：

- `idx_rule_effect_tracking_rule_version(applied_rule_version, period)`
- `idx_rule_effect_tracking_proposal(proposal_id, period)`
- `idx_rule_effect_tracking_trend(trend_direction, updated_at DESC)`

约束：

- 追踪事实可被复盘汇总读取，并可生成 warning 或 draft follow-up suggestion，但不得直接创建、回滚或替换 active `rule_versions`。
- 趋势样本不足时必须使用 `unknown` 或说明信息不足，不得宣称改善或恶化。
- 写入追踪事实时可写入对应审计事件，但不得写入交易、账户、券商、订单或外部推送事实。

### 6.5 守门人审计表：`gatekeeper_audits`

用途：保存规则修改提案的审计结论。

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `gatekeeper_audit_id` | TEXT PK | 审计 ID |
| `proposal_id` | TEXT | 提案 ID |
| `audit_result` | TEXT | approved / rejected / needs_user_review |
| `audit_reason` | TEXT | 审计原因 |
| `required_changes` | TEXT | 需要调整的内容 |
| `violates_fundamental_rule` | INTEGER | 是否违反根本规则，0 / 1 |
| `has_rule_conflict` | INTEGER | 是否存在规则冲突，0 / 1 |
| `backtest_metrics_json` | TEXT | 回测指标变化 |
| `allow_apply` | INTEGER | 是否允许进入规则应用，0 / 1 |
| `audited_rule_version` | TEXT | 审计所基于的规则版本 |
| `created_at` | DATETIME | 创建时间 |

约束：

- 审计通过后仍需用户最终确认。
- 审计否决不得写入正式规则。
- `allow_apply=1` 仍不代表自动应用规则，必须经过用户确认。

## 7. 状态模型

### 7.1 用户确认状态

| 状态 | 来源动作 | 是否更新账户 | 说明 |
| --- | --- | --- | --- |
| `not_required` | 系统裁决 | 否 | 不需要用户处理 |
| `pending` | 系统裁决 | 否 | 等待用户记录处理方式 |
| `planned` | planned | 否 | 仅记录计划 |
| `executed_manually` | executed_manually | 是 | 用户已在线下执行，更新账户和持仓 |
| `watch` | watch | 否 | 标记待观察 |
| `marked_error` | marked_error | 否 | 写入错误案例库 |

### 7.2 页面状态

| 状态 | 触发条件 | 数据来源 |
| --- | --- | --- |
| `first_use` | DATA_REQUIRED | `portfolio_snapshots` / `position_snapshots` 缺失 |
| `normal` | 工作流完成且无高危规则 | `decision_records.dashboard_state` |
| `insufficient_data` | DATA_STALE / EVIDENCE_NOT_FOUND / VECTOR_INDEX_UNAVAILABLE | `decision_records.errors_json` |
| `frozen_watch` | SOURCE_VERIFICATION_FAILED | `source_verifications` / `evidence_refs` |
| `high_risk` | 逻辑破坏、规则缺失、记录失败 | `decision_records` / `audit_events` |

`high_risk` 必须同时保存结构化原因，写入 `decision_records.risk_reason_code`，可选值包括 `BUY_LOGIC_BROKEN / RULE_VERSION_MISSING / DECISION_RECORD_FAILED`。

### 7.3 持仓状态

| 状态 | 触发条件 | 说明 |
| --- | --- | --- |
| `normal` | 买入逻辑完好 | 正常监控 |
| `sell_only` | 买入逻辑破坏 | 禁止新增买入 |
| `frozen_watch` | 多源验证不足或突发事件待核查 | 暂停主动操作 |

## 8. 写入时机

| 事件 | 写入表 | 说明 |
| --- | --- | --- |
| 用户录入持仓 | `positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 初始账户状态 |
| 用户编辑或移除当前持仓 | `positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 生成新的当前态和快照；历史保留 |
| 用户记录线下交易 | `position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 只记录用户已在线下完成的动作，不执行交易 |
| 用户校验批量导入 | `local_account_import_batches` | 只保存校验 metadata 和 rows hash，不写账户事实 |
| 用户确认批量导入 | `local_account_import_batches`、`positions`、`position_transactions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 仅允许与已校验 rows hash 匹配的有效批次 |
| 用户记录错误修正 | `local_account_corrections`、`audit_events` | 只记录修正审计；状态变化另走持仓编辑或线下交易 |
| 用户确认计划 | `operation_confirmations`、`audit_events` | 不更新账户 |
| 用户确认已手动执行 | `operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 更新本地账户状态 |
| 市场数据刷新 | `market_snapshots`、`audit_events` | 保存行情、估值、流动性、情绪指标与刷新审计结果 |
| 情报采集 | `intelligence_items` | 保存来源与 hash |
| 情报清洗 | `intelligence_summary`、`rag_chunks` | 写 SQLite 摘要与文本块 |
| 情报核查 | `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications`、`audit_events` | 保存情报事实、摘要、文本块、多源验证结果与审计事件，供决策工作流引用 |
| VecLite 建索引 | VecLite 文件、`rag_chunks` | 从 `rag_chunks` 写入向量索引，并更新索引状态 |
| 工作流完成 | `decision_records`、`evidence_refs`、`audit_events` | 保存正式建议、证据引用与审计事件 |
| 运行状态需要提示用户 | `notifications` | 写入应用内通知；不执行外部推送、交易或规则应用 |
| 用户标记错误 | `operation_confirmations`、`error_cases`、`audit_events` | 写入错误案例 |
| 规则提案生成 | `rule_proposals`、`notifications`、`audit_events` | 生成提案、写入待用户处理的应用内通知；不修改正式规则 |
| 规则效果验证 | `rule_effect_validations`、`audit_events` | 生成或刷新提案验证事实；不修改正式规则、账户或交易事实 |
| 规则效果追踪 | `rule_effect_tracking`、`audit_events` | 记录应用后规则效果；不自动回滚、替换规则或执行交易 |
| 守门人审计 | `gatekeeper_audits`、`rule_proposals.status`、`audit_events` | 生成审计结论，推进提案状态并写入审计事件 |
| 规则应用 | `rule_versions`、`rule_proposals`、`audit_events` | 仅限 `pending_final_confirm + confirm=true + sample_count>=3`；样本不足不得写入正式规则 |

### 8.1 原子事实单元与仓储错误分类

跨表事实写入必须在 repository 边界定义原子单元。一个原子单元中任一表写入失败时，相关新记录必须回滚，调用方必须收到可分类应用错误。

| 原子事实单元 | 涉及表 | 事务要求 |
| --- | --- | --- |
| 决策记录保存 | `decision_records`、`evidence_refs`、`audit_events` | DecisionRecordNode 的正式输出必须整体成功或失败。 |
| 市场数据刷新 | `market_snapshots`、`audit_events` | 成功刷新时快照与成功审计必须整体成功；快照写入失败时仍必须以独立事务保留失败审计。 |
| 用户确认计划 | `operation_confirmations`、`audit_events` | 不更新账户，但确认记录与审计必须一致。 |
| 用户录入、编辑或移除本地持仓 | `positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 当前态、快照和审计必须整体成功或失败；历史快照保留。 |
| 用户记录线下交易 | `position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 本地交易事实、当前态、快照和审计必须整体成功或失败。 |
| 用户校验批量导入 | `local_account_import_batches` | 只保存校验 metadata 和 rows hash，不写账户事实。 |
| 用户确认批量导入 | `local_account_import_batches`、`positions`、`position_transactions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 批次状态、rows hash、账户事实和审计必须整体成功或失败。 |
| 用户记录错误修正 | `local_account_corrections`、`audit_events` | 修正事实和审计必须一致；不得静默覆盖历史。 |
| 用户确认已手动执行 | `operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events` | 本地账户变化与确认审计必须整体成功或失败。 |
| 用户标记错误 | `operation_confirmations`、`error_cases`、`audit_events` | 错误案例与确认动作必须一致。 |
| 证据事实保存 | `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications`、`audit_events` | 情报、摘要、索引块、多源验证和审计事件必须整体成功或失败。 |
| 规则提案生成 | `rule_proposals`、`notifications`、`audit_events` | 规则提案、待处理通知和审计事件必须整体成功或失败。 |
| 规则效果验证 | `rule_effect_validations`、`audit_events` | 验证事实和审计事件必须整体成功或失败。 |
| 规则效果追踪 | `rule_effect_tracking`、`audit_events` | 追踪事实和审计事件必须整体成功或失败。 |
| 守门人审计 | `gatekeeper_audits`、`rule_proposals.status`、`audit_events` | 审计结论、提案状态推进和审计事件必须整体成功或失败。 |
| 规则应用 | `rule_versions`、`rule_proposals`、`audit_events` | 旧 active 归档、新 active 创建、提案应用状态和审计事件必须整体成功或失败。 |

仓储错误至少分为 `not_found`、`conflict`、`invalid_state`、`constraint`、`internal`。SQLite `sql.ErrNoRows` 必须映射为 `not_found`；唯一约束或 CHECK 约束失败必须映射为 `conflict` 或 `constraint`；非法状态流转必须映射为 `invalid_state`。

## 9. 表关系约束

| 关系 | 基数 | 说明 |
| --- | --- | --- |
| `portfolio_snapshots.snapshot_id -> position_snapshots.snapshot_id` | 1:N | 一个账户快照对应多条持仓快照 |
| `portfolio_snapshots.snapshot_id -> decision_records.portfolio_snapshot_id` | 1:N | 多条建议可引用同一账户快照 |
| `market_snapshots.market_snapshot_id -> decision_records.market_snapshot_id` | 1:N | 多条建议可引用同一市场快照 |
| `decision_records.decision_id -> evidence_refs.decision_id` | 1:N | 一条建议引用多条证据 |
| `decision_records.decision_id -> operation_confirmations.decision_id` | 1:N | 一条建议可被多次记录处理动作 |
| `operation_confirmations.confirmation_id -> position_transactions.confirmation_id` | 1:0..N | 一次已执行确认可产生多条交易流水 |
| `operation_confirmations.confirmation_id -> error_cases.confirmation_id` | 1:0..1 | 标记错误时生成错误案例 |
| `operation_confirmations.error_case_id -> error_cases.error_case_id` | 0..1:1 | 确认记录可直接引用错误案例 |
| `intelligence_items.intelligence_id -> intelligence_summary.intelligence_id` | 1:N | 一条原始情报可生成多条摘要 |
| `intelligence_summary.summary_id -> rag_chunks.summary_id` | 1:N | 一条摘要可拆为多段文本块 |
| `intelligence_summary.summary_id -> evidence_refs.summary_id` | 1:N | 一条摘要可被多条建议引用 |
| `intelligence_summary.verification_group_id -> source_verifications.verification_group_id` | N:1 | 多条摘要归入同一验证分组 |
| `notifications.source_type/source_id -> 业务来源` | N:0..1 | 通知可关联 market_refresh、evidence_refresh、rule_proposal、review_summary 等本地来源 |
| `error_cases.error_case_id -> rule_proposals.source_error_case_id` | 1:N | 一个错误案例可生成多个提案 |
| `rule_proposals.proposal_id -> gatekeeper_audits.proposal_id` | 1:N | 一个提案可经历多次审计 |
| `rule_proposals.proposal_id -> rule_versions.created_from_proposal_id` | 1:0..1 | 提案被应用后生成规则版本 |
| `rule_proposals.applied_rule_version -> rule_versions.rule_version` | 0..1:1 | 提案应用后指向正式规则版本 |

删除规则：

- `decision_records`、`evidence_refs`、`audit_events`、`rule_versions` 不做级联删除。
- 情报原文可按用户清理策略移除，但必须保留摘要、来源、hash 和审计记录。
- 业务历史记录默认只追加，不覆盖。

## 10. JSON 字段结构约束

### 10.1 `analyst_reports_json`

```json
[
  {
    "agent_name": "价值分析师",
    "conclusion": "估值处于观察区，暂不新增买入。",
    "key_reasons": ["PE 分位处于 50%-80%"],
    "risk_warnings": ["继续关注估值上行风险"],
    "confidence": "medium",
    "evidence_ids": ["ev_001"],
    "prompt_version": "p37-analyst-v1",
    "model": "gpt-5.4-mini",
    "input_summary": "value 510300 ...",
    "output_summary": "估值处于观察区...",
    "parse_status": "parsed",
    "quality_status": "passed"
  }
]
```

P37 起，`prompt_version`、`model`、`input_summary`、`output_summary`、`parse_status` 和 `quality_status` 为可选 LLM metadata，用于追踪真实模型调用质量。字段不得包含 API key、完整本地路径、完整账户明细或未经摘要的敏感输入；质量失败的 LLM 输出不得作为最终裁决、交易动作、确认记录、规则更新或外部通知来源。

### 10.2 `arbitration_chain_json`

```json
[
  {
    "priority": 3,
    "rule_id": "R-3",
    "result": "限制新增买入"
  }
]
```

### 10.3 `triggered_rules_json`

```json
[
  {
    "rule_id": "R-3",
    "rule_name": "不超过仓位上限",
    "severity": "warning",
    "description": "高风险资产仓位接近上限。"
  }
]
```

### 10.4 `errors_json`

```json
[
  {
    "code": "VECTOR_INDEX_UNAVAILABLE",
    "node_name": "EvidenceRetrievalNode",
    "message": "VecLite 索引不可用，已使用 SQLite 摘要降级展示。"
  }
]
```

## 11. VecLite 索引设计

### 11.1 索引内容

VecLite 存储来自 `rag_chunks` 的文本块向量。

metadata 必须包含：

- `chunk_id`
- `summary_id`
- `symbol`
- `source_level`
- `evidence_role`
- `published_at`
- `content_hash`
- `chunk_hash`
- `embedding_model`
- `embedding_version`
- `index_version`

### 11.2 重建规则

当 VecLite 索引不可用或版本不兼容时：

1. 从 `rag_chunks` 读取 `index_status=pending` 或 `index_status=stale` 的文本块。
2. 使用 `embedding_model` 与 `embedding_version` 重新生成向量。
3. 写入 `vector_collection` 对应的 VecLite 索引。
4. 更新 `vector_id`、`indexed_at` 和 `index_status`。
5. 若 SQLite 摘要充足，系统可降级使用 SQLite 摘要。
6. 若 SQLite 摘要不足，前端进入信息不足状态。

## 12. 数据保留与审计规则

- `decision_records`、`evidence_refs`、`audit_events` 不允许物理删除，除非用户明确执行数据清理。
- 删除情报原文时，必须保留 `content_hash`、摘要和来源信息。
- 所有规则版本变更必须保留历史版本。
- 所有用户确认动作必须保留审计事件。
- 备份必须同时包含 SQLite 数据文件和 VecLite 向量文件。
- 金额字段当前使用 `REAL`，后续如需要严格对账，可迁移为分单位整数或 decimal 字符串。

## 13. P74 内置知识与准备度模型

P74 不新增 SQLite 表。内置知识注册表由应用代码提供稳定只读条目，并在运行时与现有 SQLite 事实组合：

| 模型 | 来源 | 用途 |
| --- | --- | --- |
| `KnowledgeEntry` | 内置 registry | 大师原则、纪律规则、风险 SOP、标的画像；不得作为正式市场证据 |
| `KnowledgeSymbolProfile` | 内置 registry | `510300` 到 `000300` 的 ETF/指数映射；未知标的返回 blocked |
| `KnowledgeDataDependency` | registry + `market_snapshots` + `source_verifications` + `rule_versions` | 数据依赖矩阵和 ready/degraded/blocked 判定 |
| `KnowledgeFeatureImpact` | dependency 聚合 | 展示降级对 consultation、decision detail、risk alerts、expected return 等功能的影响 |

依赖读取规则：

- `market_snapshots.market_metrics_json.metadata.p34_source_health` 提供 `fund_profile`、`tracked_index`、`market_price`、`valuation_percentiles`、`liquidity`、`sentiment_proxy` 等 freshness/source level。
- `source_verifications` 提供 `formal_evidence` readiness；必须高等级多源满足才可 ready。
- `rule_versions.status='active'` 提供 `active_rule` readiness；缺失时 degraded。
- `rag_index` readiness 只表达检索降级边界，不改变 SQLite 事实。
- `llm_context_summary` 只持有脱敏摘要，不持久化完整 prompt。

## 14. 与其他文档关系

- 系统架构与技术栈见 `docs/architecture.md`。
- Eino 工作流节点和错误处理见 `docs/workflow.md`。
- HTTP API 字段和枚举见 `docs/api.md`。
- 前端展示字段映射见 `docs/frontend-contract.md`。
- 后续数据库表关系图应以本文档为依据。
