# Tasks: p4-http-api

> 对齐 `docs/development-plan.md` P4：HTTP API。实现代码必须写中文注释，说明 DTO 字段来源、handler 边界、事务写入和错误处理意图。

## 1. API DTO 与错误处理（P4.1）

参考文档：`docs/api.md`、`docs/frontend-contract.md`。

创建或确认以下文件：

```text
internal/application/dto/common.go
internal/application/dto/dashboard.go
internal/application/dto/decision.go
internal/application/dto/evidence.go
internal/application/dto/rule.go
internal/application/dto/audit.go
internal/application/handler/errors.go
pkg/httputil/response.go
```

任务：

- [x] 1.1 创建 `internal/application/dto/common.go`
- [x] 1.2 创建 `internal/application/dto/dashboard.go`
- [x] 1.3 创建 `internal/application/dto/decision.go`
- [x] 1.4 创建 `internal/application/dto/evidence.go`
- [x] 1.5 创建 `internal/application/dto/rule.go`
- [x] 1.6 创建 `internal/application/dto/audit.go`
- [x] 1.7 创建 `internal/application/handler/errors.go`
- [x] 1.8 确认或扩展 `pkg/httputil/response.go`
- [x] 1.9 所有响应包含 `request_id`
- [x] 1.10 错误响应符合 `docs/api.md` 第 2、3 节
- [x] 1.11 DTO 字段名与 `docs/frontend-contract.md` 一致
- [x] 1.12 `EVIDENCE_NOT_FOUND` 返回 409，并让前端显示信息不足状态
- [x] 1.13 为 P4.1 新增代码补充中文注释
- [x] 1.14 验收：`go test ./internal/application/dto/... ./internal/application/handler/... ./pkg/httputil/...`

## 2. 核心 API（P4.2）

创建或确认以下文件：

```text
internal/application/handler/dashboard_handler.go
internal/application/handler/decision_handler.go
internal/application/handler/portfolio_handler.go
internal/application/handler/evidence_handler.go
internal/application/handler/rule_handler.go
internal/application/handler/audit_handler.go
internal/application/handler/settings_handler.go
internal/application/handler/market_handler.go
internal/application/handler/review_handler.go
```

任务：

- [x] 2.1 创建 `dashboard_handler.go`，实现 `GET /api/v1/dashboard/today`
- [x] 2.2 创建 `decision_handler.go`，实现 `POST /api/v1/decisions/consult`、`GET /api/v1/decisions/{decision_id}`、`GET /api/v1/decisions`、`POST /api/v1/decisions/{decision_id}/confirmations`
- [x] 2.3 创建 `portfolio_handler.go`，实现 `POST /api/v1/portfolio/init`、`GET /api/v1/portfolio/current`、`POST /api/v1/portfolio/adjustments`
- [x] 2.4 创建 `evidence_handler.go`，实现 `POST /api/v1/evidence/refresh`、`GET /api/v1/evidence`、`GET /api/v1/evidence/verification`、`POST /api/v1/evidence/rebuild-index`
- [x] 2.5 创建 `market_handler.go`，实现 `POST /api/v1/market/refresh`、`GET /api/v1/market/snapshots/latest`
- [x] 2.6 创建 `rule_handler.go`，实现 `GET /api/v1/rules/current`、`GET /api/v1/rule-proposals`、`POST /api/v1/rule-proposals/{proposal_id}/confirm`、`POST /api/v1/rule-proposals/{proposal_id}/final-confirm`
- [x] 2.7 创建 `settings_handler.go`，实现 `GET /api/v1/settings/system`、`PUT /api/v1/settings`、`GET /api/v1/settings/capability`、`PUT /api/v1/settings/capability`
- [x] 2.8 创建 `audit_handler.go`，实现 `GET /api/v1/audit-events`
- [x] 2.9 创建 `review_handler.go`，实现 `GET /api/v1/review/summary`
- [x] 2.10 为 P4.2 新增代码补充中文注释

## 3. P4.2 验收断言

- [x] 3.1 `executed_manually` 成功后同时写入 6 类记录：`operation_confirmations`、`position_transactions`、`positions`、`portfolio_snapshots`、`position_snapshots`、`audit_events`
- [x] 3.2 `marked_error` 成功后同一事务写入 `operation_confirmations`、`error_cases`、`audit_events`，响应返回 `error_case_id`
- [x] 3.3 `planned` 与 `watch` 不写 `position_transactions`，不新增账户快照
- [x] 3.4 `planned` 与 `watch` 可互相转换，并可升级为 `executed_manually` 或 `marked_error`；每次成功转换都创建新的 `operation_confirmations`
- [x] 3.5 `executed_manually` 与 `marked_error` 是确认终态，再次确认返回 `BAD_REQUEST`，不得重复写账户快照、交易流水或错误案例
- [x] 3.6 `record_type!=formal_trade_advice` 或 `confirmation_status=not_required` 时，确认接口返回 `BAD_REQUEST`
- [x] 3.7 守门人审计通过只进入 `pending_final_confirm`，不写 `rule_versions`
- [x] 3.8 最终确认应用规则后创建新 active `rule_versions`，旧 active 归档
- [x] 3.9 `POST /api/v1/evidence/refresh` 同步完成情报采集、摘要、`source_verifications` 写入和索引更新；索引失败不得回滚 SQLite 事实数据
- [x] 3.10 市场状态枚举统一使用 `liquidity_state=normal/warning/danger`、`sentiment_state=cold/neutral/hot/extreme`
- [x] 3.11 `PUT /api/v1/settings` 只能保存通知、页面偏好、普通数据源；`PUT /api/v1/settings/capability` 只保存能力圈；规则阈值、裁决优先级和 SOP 必须生成 `rule_proposals`
- [x] 3.12 `POST /api/v1/market/refresh` 覆盖全部成功、部分成功、全部失败、快照写入失败四类响应；部分成功时返回 200，并在 `failed_symbols` 中写明失败标的与原因
- [x] 3.13 `sample_count<3` 的规则提案调用送审接口返回 `BAD_REQUEST`，不得写 `gatekeeper_audits` 或进入 `under_gatekeeper_audit`
- [x] 3.14 `sample_count<3` 的规则提案即使状态异常进入 `pending_final_confirm`，最终确认接口也必须返回 `BAD_REQUEST`，不得写 `rule_versions`

## 4. 验收与归档前

- [x] 4.1 验收：`go test ./internal/application/handler/...`
- [x] 4.2 验收期望：API 契约测试覆盖成功响应、错误响应、状态流转、事务写入
- [x] 4.3 确认本 change 的 specs delta 只包含 P4 HTTP API 相关内容
- [x] 4.4 完成后更新 `docs/development-plan.md` P4 任务状态
- [x] 4.5 完成后更新 `docs/GOVERNANCE.md` 活跃变更状态
- [x] 4.6 完成后更新 `openspec/PROGRESS.md`：P4 标记为 done，下一阶段指向 P5

## Plan alignment

- P4.1 对应任务：1.1–1.14。
- P4.2 核心 API 对应任务：2.1–2.10。
- P4.2 验收断言对应任务：3.1–3.14。
- P4 验收与归档前任务：4.1–4.6。
- 与 `docs/development-plan.md` P4 小节、创建文件、任务列表和验收命令一一对应；仅额外加入中文注释要求，来自本次用户指令。
