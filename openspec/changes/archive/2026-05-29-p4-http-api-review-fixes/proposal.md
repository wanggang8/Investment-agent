# Proposal: P4 HTTP API Review Fixes

## Summary

修复 P4 HTTP API 子 agent 审查发现的阻塞问题和主要契约偏差：补齐确认字段校验、规则拒绝路径、规则最终确认审计事务、市场错误码、证据索引失败状态、设置保存和 DTO 字段。

## Why

P4 已归档后复审发现部分任务被标记完成但实现不完整。该修复 change 用于在不展开 P5 的前提下，将 P4 API 行为补齐到 `docs/api.md`、`docs/data-model.md`、`docs/frontend-contract.md` 和 `docs/workflow.md` 的当前契约。

## In Scope

- 修复 `confirmation_handler.go` 的 executed_manually 与 marked_error 请求校验。
- 修复 `rule_handler.go` 的 confirm=false 拒绝路径、空审计结果保护、最终确认同事务审计写入。
- 补齐 `apperr` 市场相关错误码。
- 修复 `market_handler.go` 非法 JSON 处理、全部失败错误码和写入失败分类。
- 修复 `evidence_handler.go` 索引失败时 `rag_chunks.index_status` 更新和响应字段。
- 修复 `settings_handler.go` 页面偏好保存、规则类字段显式拒绝。
- 补齐 `RuleProposalDTO` 前端契约字段。
- 补充针对上述问题的契约测试。

## Out of Scope

- 不新增 P5 前端页面。
- 不接入真实外部行情、券商或 VecLite 服务。
- 不重写 P4 架构，不改变已归档 P4 的 API 列表。
