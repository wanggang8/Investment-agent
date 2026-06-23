# Design: P4 HTTP API

## Overview

P4 在应用层增加 DTO 与 handler，把 P1-P3 已完成的仓储、规则和工作流能力封装成本地 HTTP API。API 只服务本地 Web 控制台和本地任务，不提供自动交易入口。

## Architecture

```text
cmd/server
  -> internal/application/handler
      -> internal/application/dto
      -> internal/application/workflow
      -> internal/domain/repository
      -> pkg/httputil
  -> internal/infrastructure/persistence/sqlite
```

## Response and Error Strategy

- 成功响应统一由 `pkg/httputil` 输出：`request_id`、`data`、可选 `meta`。
- 失败响应统一由 `apperr` 映射为错误码和 HTTP 状态。
- handler 不解析底层 SQL 错误文本；仓储错误已在仓储出口转为 `apperr`。
- 未知错误统一返回 `INTERNAL_ERROR`，不暴露 SQL、文件路径或外部服务原始错误。
- `EVIDENCE_NOT_FOUND` 返回 409，供前端映射为信息不足状态。

## DTO Strategy

- DTO 位于 `internal/application/dto`，只表达 API 与前端契约字段。
- DTO 字段名使用 JSON snake_case，对齐 `docs/api.md`。
- handler 内做领域模型到 DTO 的转换，不让前端依赖领域内部结构。
- DTO 新增代码写中文注释，说明字段来源或展示约束。

## Handler Strategy

- handler 位于 `internal/application/handler`。
- 每个 handler 只负责请求解析、调用应用能力、DTO 组装和响应输出。
- 会写入多张表的接口必须通过仓储事务或新增事务服务保证原子性。
- 用户确认接口按 `docs/api.md` 约束区分 `planned`、`watch`、`executed_manually`、`marked_error`。

## Transaction Boundaries

必须保证以下接口事务原子性：

- `POST /api/v1/portfolio/init`：`portfolio_snapshots + positions + position_snapshots + audit_events`。
- `POST /api/v1/portfolio/adjustments`：`positions + portfolio_snapshots + position_snapshots + audit_events`，不得写 `position_transactions`。
- `POST /api/v1/decisions/{decision_id}/confirmations`：
  - `executed_manually` 写 `operation_confirmations + position_transactions + positions + portfolio_snapshots + position_snapshots + audit_events`。
  - `marked_error` 写 `operation_confirmations + error_cases + audit_events`。
  - `planned` 与 `watch` 只写确认记录和审计，不写交易流水或新增账户快照。
- 规则最终确认成功后创建新 active `rule_versions`，旧 active 归档。

## API Groups

- Dashboard：读取今日纪律驾驶舱。
- Decision：咨询、详情、列表、确认。
- Portfolio：初始化、当前状态、手动校准。
- Evidence：刷新、查询、验证、重建索引。
- Market：刷新、最新快照。
- Rule：当前规则、提案列表、送审确认、最终确认。
- Settings：系统设置、能力圈设置。
- Audit：审计事件查询。
- Review：复盘摘要。

## Testing Strategy

- DTO 与错误响应测试：`go test ./internal/application/dto/... ./internal/application/handler/... ./pkg/httputil/...`。
- Handler 契约测试：`go test ./internal/application/handler/...`。
- 测试必须覆盖成功响应、错误响应、状态流转和事务写入。
- 对用户确认、规则确认、证据刷新和市场刷新写字段级断言，不只断言“有记录”。

## Open Questions

无。P4 范围完全来自 `docs/development-plan.md` P4 与既有 L1 契约。
