# Real Data Integration Specification

## Purpose
Document real data, retrieval, analyst, and runtime wiring requirements for Investment Agent integrations.
## Requirements
### Requirement: Market and intelligence data sources
系统 SHALL 提供最小可用真实只读行情与情报数据源适配能力，并支持按配置启用真实数据源或本地 stub。P12 范围 SHALL 明确排除完整财务源、完整情绪源、实时性 SLA、券商交易 API、自动交易、主动荐股和收益承诺。

#### Scenario: Market data refresh succeeds
- **WHEN** configured readonly market data source returns valid data for target symbols
- **THEN** 系统 MUST 写入 `market_snapshots`
- **AND** 系统 MUST 写入成功状态的 `audit_events`

#### Scenario: Market data refresh degrades
- **WHEN** 行情数据源部分失败、全部失败、超时、解析失败或返回过期数据
- **THEN** 系统 MUST 返回既有错误或降级状态
- **AND** 系统 MUST 记录失败标的、降级原因或错误码
- **AND** 系统 MUST 写入可追踪的审计事件

#### Scenario: Intelligence data is ingested
- **WHEN** configured readonly intelligence source returns valid news, announcement, or manually imported intelligence
- **THEN** 系统 MUST 写入 `intelligence_items`
- **AND** 系统 MUST 保留来源、时间和信源等级信息

#### Scenario: Local stub is available
- **WHEN** 真实数据源未配置或被关闭
- **THEN** 系统 MUST 可使用本地 stub 完成开发与测试
- **AND** 系统 MUST NOT 写入真实密钥或环境私有值到文档或日志

### Requirement: RAG and VecLite retrieval
系统 SHALL 提供本地 JSON 文件索引读写、SQLite 重建、健康状态、重建统计和不可用降级能力；真实 VecLite API 替换边界 SHALL 保持可替换但不在 P13 强制接入。

#### Scenario: VecLite index is built
- **WHEN** `rag_chunks` 与 `intelligence_summary` 存在可检索文本
- **THEN** 系统 MUST 将其纳入本地 JSON 文件索引构建流程
- **AND** 索引路径 MUST 来自配置
- **AND** 系统 MUST 记录健康状态和重建统计

#### Scenario: VecLite index is rebuilt
- **WHEN** 本地索引缺失、损坏、不兼容或需要重建
- **THEN** 系统 MUST 支持从 SQLite 文本块重建索引
- **AND** 重建过程 MUST 可测试
- **AND** 重建结果 MUST 暴露 indexed/skipped 数量与降级原因

#### Scenario: VecLite is unavailable
- **WHEN** 本地索引不可用、损坏、不兼容或检索失败
- **THEN** 系统 MUST 按既有约定降级到 SQLite 摘要
- **AND** 摘要不足时 MUST 返回信息不足
- **AND** 系统 MUST 记录检索输入、命中证据和降级原因

#### Scenario: Retrieval service records fallback context
- **WHEN** 检索从本地索引降级到 SQLite 摘要
- **THEN** 工作流 MUST 保留输入标的、命中摘要或降级原因
- **AND** 审计事件 MUST 可追踪本次检索状态

#### Scenario: C-level source is restricted
- **WHEN** 检索命中 C 级信源
- **THEN** C 级信源 MUST 只能作为 background 材料
- **AND** C 级信源 MUST NOT 作为正式裁决依据

### Requirement: DeepSeek analyst materials
系统 SHALL 接入可配置 DeepSeek/OpenAI-compatible 分析服务，并确保 LLM 只生成分析材料，不生成最终裁决。

#### Scenario: DeepSeek analysis succeeds
- **WHEN** DeepSeek 或兼容端点返回可解析输出
- **THEN** 系统 MUST 将输出解析为 `analyst_reports` 或等价结构
- **AND** 系统 MUST 记录 prompt version、model、input summary、output summary、parse status 和 quality status
- **AND** 规则引擎 MUST 继续负责最终裁决

#### Scenario: DeepSeek input is bounded
- **WHEN** 系统构造 DeepSeek prompt
- **THEN** prompt 输入 MUST 只包含允许使用的证据、持仓上下文和规则边界
- **AND** 非显然 prompt 约束 MUST 有中文注释说明
- **AND** 审计与响应 MUST NOT 包含 API key、完整敏感 prompt、券商/账户密钥或不必要本地文件路径

#### Scenario: DeepSeek supports expected return material
- **WHEN** 预期收益节点执行
- **THEN** 系统 MUST 调用分析服务生成预期收益分析材料
- **AND** 数值情景可继续由本地样本逻辑生成
- **AND** DeepSeek MUST NOT 写最终裁决

#### Scenario: DeepSeek is unavailable
- **WHEN** DeepSeek 缺配置、超时、HTTP 不可用、空响应、输出不可解析或质量检查失败
- **THEN** 工作流 MUST 进入降级状态
- **AND** 规则引擎 MUST 继续生成最终裁决
- **AND** 系统 MUST 写入可追踪的降级原因和稳定错误分类

#### Scenario: No automatic trading from analysis
- **WHEN** DeepSeek 或检索服务生成分析材料
- **THEN** 系统 MUST NOT 生成自动交易动作
- **AND** 系统 MUST NOT 提供一键交易入口

### Requirement: Runtime dependency wiring
系统 SHALL 根据本地配置组装真实数据、检索与分析服务依赖。

#### Scenario: DeepSeek key is configured
- **WHEN** `deepseek.api_key`、`deepseek.base_url` 和 `deepseek.model` 或等价配置存在
- **THEN** 生产依赖 MUST 使用配置的 DeepSeek/OpenAI-compatible client 作为分析服务

#### Scenario: DeepSeek key is missing
- **WHEN** DeepSeek key 缺失
- **THEN** 生产依赖 MUST 使用可追踪降级服务或本地 stub
- **AND** 不得伪造真实 DeepSeek 响应

#### Scenario: Data source stub setting is applied
- **WHEN** `data_sources.use_stub` 为 true
- **THEN** 生产依赖 MUST 使用本地 stub 数据源
- **AND** 本地验收 MUST 不依赖公网

### Requirement: P8 review data metadata preservation
The system SHALL preserve source metadata, evidence quality metadata, verification status, market freshness, and rule proposal audit facts across persistence, DTOs, and frontend display.

#### Scenario: Evidence metadata is preserved from intelligence item to API
- **WHEN** intelligence items and summaries are listed as evidence
- **THEN** source name, original URL, published time, captured time, content hash, time weight, relevance score, source level, evidence role, verification status, independent source count, and high-grade independent source count SHALL be returned without placeholder substitution.

#### Scenario: Evidence quality is preserved in decision references
- **WHEN** evidence is persisted as decision evidence refs
- **THEN** formal/background role, source level, time weight, relevance score, content hash, and high-grade independent source count SHALL remain available for decision-detail reconstruction.

#### Scenario: Market snapshot reports freshness
- **WHEN** market data is refreshed or queried
- **THEN** the market DTO SHALL include enough date and status information for the frontend to distinguish fresh, stale, and missing data.

#### Scenario: Rule proposal audit facts gate final application
- **WHEN** final rule application is requested
- **THEN** the latest approved gatekeeper audit allowing application SHALL be present before an active rule is written.

### Requirement: Local tasks can trigger real data maintenance
The system SHALL allow the local `cmd/agent` entrypoint to trigger market refresh, intelligence indexing, and VecLite-related maintenance tasks while preserving the existing real data degradation and audit behavior.

#### Scenario: Market refresh can be triggered locally
- **WHEN** the user triggers market refresh through `cmd/agent`
- **THEN** the system uses configured or stub data sources, records audit events, and reports readable errors on source or write failure.

#### Scenario: Intelligence indexing can be triggered locally
- **WHEN** the user triggers intelligence indexing through `cmd/agent`
- **THEN** the system updates local intelligence/RAG data using existing repositories and records task execution in audit events.

#### Scenario: Index recovery is documented
- **WHEN** VecLite index data is unavailable or damaged
- **THEN** the system documentation describes how to rebuild or recover from local persisted data.

### Requirement: Expanded real data SHALL enter workflow input context safely

The system SHALL make P34 expanded public data available to daily discipline, expected return, and future risk warning input context without bypassing rules or human review.

#### Scenario: Daily discipline reads expanded data
- **WHEN** DailyDisciplineGraph prepares context for configured holdings or indexes
- **THEN** it SHALL include available P34 valuation, constituent, financial, capital-flow, margin, or sentiment-proxy summaries with source level and freshness status
- **AND** missing or stale categories SHALL remain explicit in the workflow context.

#### Scenario: Expected return reads expanded data
- **WHEN** ExpectedReturnNode prepares scenario material
- **THEN** it SHALL treat P34 data as supporting sample/context material
- **AND** it SHALL preserve sample limitations, freshness, source level, and missing categories
- **AND** it SHALL NOT convert P34 data into guaranteed return, deterministic price, or automatic trading output.

### Requirement: Expanded real data SHALL preserve stub and degradation behavior

The system SHALL keep deterministic local validation possible when expanded public data sources are disabled or unavailable.

#### Scenario: Stub mode is enabled
- **WHEN** `data_sources.use_stub` or an equivalent test fixture mode is enabled
- **THEN** P34 data services SHALL use deterministic fixture or stub data for tests and local validation
- **AND** they SHALL not require public network access.

#### Scenario: Real source mode is enabled but fails
- **WHEN** a real P34 source is enabled and fails with timeout, unavailable source, parse failure, stale data, or no records
- **THEN** the application SHALL return a stable degraded or insufficient-data state
- **AND** it SHALL write auditable failure metadata
- **AND** it SHALL not silently substitute fabricated real data.

### Requirement: Expanded real data SHALL expose health to application surfaces

The system SHALL expose source health and recent refresh status for P34 data through application service DTOs or existing settings/ops surfaces.

#### Scenario: Frontend requests data source health
- **WHEN** the frontend displays data source, daily discipline, or ops status
- **THEN** it SHALL be able to show each P34 source category as fresh, stale, missing, unavailable, parse-error, disabled, or stubbed
- **AND** `GET /api/v1/market/source-health` or an equivalent settings/ops DTO SHALL expose source category, freshness, source level, source type, data date, last success or failure time, failure category, and affected symbols where available
- **AND** it SHALL show last success or failure time where available.

#### Scenario: Health state is insufficient for formal analysis
- **WHEN** required P34 source categories are stale, missing, or failed
- **THEN** frontend and workflow surfaces SHALL explain the missing category
- **AND** they SHALL not present the analysis as fully evidenced.

### Requirement: Data health SHALL feed risk alert orchestration

The system SHALL use real data freshness, source health, and degraded public data diagnostics as inputs to risk alert orchestration without treating missing data as complete evidence.

#### Scenario: Source health is stale or failed
- **WHEN** P34 source health indicates stale, missing, no_data, source_unavailable, parse_error, disabled, or stubbed status for a category needed by risk analysis
- **THEN** risk alert orchestration SHALL preserve that condition as degraded or insufficient data context
- **AND** it SHALL NOT silently clear risk alerts that depend on the missing or degraded category.

#### Scenario: Expanded data supports a risk alert
- **WHEN** P34 expanded data is fresh enough to support valuation, liquidity, sentiment, or evidence-insufficiency risk checks
- **THEN** risk alert orchestration SHALL record the source category, freshness, source level, data date, and affected symbols in the risk trigger context where available
- **AND** it SHALL keep lower-grade or stubbed data visibly marked as supporting context rather than formal high-confidence evidence.

### Requirement: Risk orchestration SHALL not expand external data boundaries

Risk alert orchestration SHALL only consume already configured local facts, market snapshots, evidence summaries, and source health records.

#### Scenario: Risk orchestration needs more data
- **WHEN** required risk inputs are missing or stale
- **THEN** the system SHALL produce degraded or insufficient-data risk diagnostics
- **AND** it SHALL NOT introduce login, paid, authorized, Level2, high-frequency, broker, or external push dependencies.

