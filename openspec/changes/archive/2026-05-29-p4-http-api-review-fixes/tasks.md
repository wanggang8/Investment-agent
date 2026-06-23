# Tasks: p4-http-api-review-fixes

## 1. 阻塞问题修复

- [x] 1.1 `executed_manually` 校验 `operation_type=buy/sell/reduce`
- [x] 1.2 `executed_manually` 校验 `executed_at` 不晚于当前时间
- [x] 1.3 `marked_error` 校验 `actual_outcome`、`root_cause_tag`、`lesson_learned` 必填
- [x] 1.4 `marked_error` 校验 `root_cause_tag` 枚举
- [x] 1.5 规则提案送审支持 `confirm=false` 拒绝路径
- [x] 1.6 规则提案送审保护空 `GatekeeperAudits`
- [x] 1.7 规则最终确认同事务写入 `audit_events`

## 2. 主要警告修复

- [x] 2.1 增加 `DATA_SOURCE_UNAVAILABLE` 与 `MARKET_SNAPSHOT_WRITE_FAILED` 错误码
- [x] 2.2 市场刷新非法 JSON 返回 `BAD_REQUEST`
- [x] 2.3 市场刷新全部失败返回 `DATA_SOURCE_UNAVAILABLE`
- [x] 2.4 证据索引失败更新 `rag_chunks.index_status=failed`
- [x] 2.5 普通设置保存 `page_preference`
- [x] 2.6 普通设置显式拒绝规则阈值、裁决优先级、SOP 类字段
- [x] 2.7 补齐 `RuleProposalDTO` 前端契约字段
- [x] 2.8 缺省 `request_id` 使用每请求唯一值

## 3. 测试与验收

- [x] 3.1 补充确认字段校验测试
- [x] 3.2 补充规则拒绝路径和最终确认审计测试
- [x] 3.3 补充市场错误码测试
- [x] 3.4 补充证据索引状态测试
- [x] 3.5 补充设置分离测试
- [x] 3.6 验收：`go test ./internal/application/handler/... ./internal/application/dto/... ./internal/pkg/apperr/... ./pkg/httputil/...`
- [x] 3.7 验收：`go test ./...`
