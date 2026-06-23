## ADDED Requirements

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
系统 SHALL 为关键页面提供证据、决策链、审计时间线交互，并清楚展示空态、错误态和降级态。

#### Scenario: Evidence and audit interactions are available
- **WHEN** 用户查看证据、决策链或审计时间线
- **THEN** 前端 MUST 支持筛选或展开关键条目
- **AND** 展开内容 MUST 继续来自 API DTO

#### Scenario: Degraded and error states are visible
- **WHEN** 页面遇到信息不足、数据过期、LLM 降级、VecLite 不可用或错误响应
- **THEN** 前端 MUST 展示明确空态或错误态
- **AND** 用户 MUST 能区分信息不足、冻结观察、高危、降级和普通错误

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
