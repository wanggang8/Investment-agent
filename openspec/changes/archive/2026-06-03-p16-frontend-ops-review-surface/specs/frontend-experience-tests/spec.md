## MODIFIED Requirements

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
- **AND** 未知状态 MUST 使用安全兜底显示，不得展示为成功

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
