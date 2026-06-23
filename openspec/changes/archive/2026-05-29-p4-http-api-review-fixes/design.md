# Design: P4 HTTP API Review Fixes

## Approach

本修复保持 P4 handler 架构不变，只在审查指出的边界处补齐校验、事务和 DTO。

## Key Decisions

- 用户确认请求先做字段级校验，再开启 SQLite 事务，避免无效请求写入 `operation_confirmations`。
- 规则提案确认请求使用 `{ "confirm": true|false }`。`false` 表示用户拒绝，写审计并把提案状态设为 `rejected`。
- 规则最终确认在同一事务中归档旧 active、插入新 active、更新提案、写 `audit_events`。
- 市场刷新把请求解析失败返回 `BAD_REQUEST`；全部数据源失败返回 `DATA_SOURCE_UNAVAILABLE`；写入失败由仓储错误转换为 `MARKET_SNAPSHOT_WRITE_FAILED`。
- 证据刷新模拟索引失败时，必须更新已写入 `rag_chunks.index_status=failed`，但保留 SQLite 情报和验证事实。
- 设置接口用宽请求体显式拦截规则阈值、裁决优先级、SOP 类字段，普通偏好只写通知、页面偏好和数据源。

## Testing

- 补充确认字段校验、规则拒绝、最终确认审计、市场错误码、证据索引状态、设置分离、DTO 字段测试。
- 验收命令：`go test ./internal/application/handler/... ./internal/application/dto/... ./internal/pkg/apperr/... ./pkg/httputil/...` 和 `go test ./...`。
