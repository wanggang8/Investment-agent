# P0-P5 Capability Summary

## Purpose

本摘要只记录已归档 P0-P5 阶段的能力边界，便于 OpenSpec 查询和校验；权威契约仍以 `docs/` 为准。本文件不替代 L1 契约。

## Requirements

### Requirement: P0 engineering skeleton capability
系统 SHALL 提供可本地启动的 Go 后端与 React 前端骨架。

#### Scenario: Local skeleton is available
- **WHEN** 开发者启动本地服务
- **THEN** 后端 MUST 提供健康检查
- **AND** 前端 MUST 提供基础页面路由
- **AND** 配置 MUST 从本地配置或环境变量读取

### Requirement: P1 data foundation capability
系统 SHALL 使用 SQLite 保存账户、持仓、行情、情报、决策、确认、规则和审计数据。

#### Scenario: Data can be migrated and queried
- **WHEN** migration 执行完成
- **THEN** Repository MUST 支持核心表写读
- **AND** 事务失败 MUST 回滚关键写入

### Requirement: P2 domain rules capability
系统 SHALL 使用领域规则生成最终裁决。

#### Scenario: Rules own final verdict
- **WHEN** 系统生成决策
- **THEN** 最终裁决 MUST 来自规则引擎
- **AND** 外部分析材料 MUST NOT 覆盖最终裁决

### Requirement: P3 workflow capability
系统 SHALL 使用工作流组织每日纪律、主动咨询、证据核验、规则提案和守门人审计。

#### Scenario: Workflows are auditable
- **WHEN** 工作流执行关键节点
- **THEN** 系统 MUST 写入可追踪审计事件
- **AND** 降级原因 MUST 可由审计或上下文追踪

### Requirement: P4 HTTP API capability
系统 SHALL 暴露统一信封格式的 HTTP API。

#### Scenario: API envelope is stable
- **WHEN** 前端请求 API
- **THEN** 响应 MUST 使用 `request_id`、`data`、`error` 信封
- **AND** 错误码 MUST 可映射到安全展示状态

### Requirement: P5 frontend cockpit capability
系统 SHALL 提供本地决策驾驶舱。

#### Scenario: Frontend shows decision context
- **WHEN** 用户查看驾驶舱、决策、证据、规则、审计或设置页面
- **THEN** 页面 MUST 展示来自 API DTO 的状态和字段
- **AND** 页面 MUST NOT 提供自动交易入口
