# Proposal: P4 HTTP API

## Summary

实现 `docs/development-plan.md` 中 P4 HTTP API 阶段：补齐 API DTO、统一响应/错误处理、核心 handler、事务性用户确认、证据/市场/规则/设置/审计/复盘接口契约测试。

## Why

P0-P3 已完成工程骨架、数据底座、领域规则与工作流。P4 需要把这些能力通过本地 HTTP API 暴露给后续 P5 前端驾驶舱，同时保持 `docs/api.md`、`docs/frontend-contract.md`、`docs/data-model.md`、`docs/workflow.md` 的行为一致。

## In Scope

- 创建 `internal/application/dto/*`，字段名对齐 `docs/frontend-contract.md` 与 `docs/api.md`。
- 创建 `internal/application/handler/*`，实现 P4.2 列出的所有核心 API。
- 所有业务响应使用包含 `request_id` 的统一 JSON 信封。
- 错误响应由 `internal/pkg/apperr` 与 `pkg/httputil` 统一生成，不暴露 SQL、文件路径或外部服务原始错误。
- 实现用户确认接口的状态门禁与事务写入规则。
- 实现证据刷新、市场刷新、规则提案确认、最终确认、设置更新、审计查询和复盘摘要接口。
- 新增 API 契约测试，覆盖成功响应、错误响应、状态流转和事务写入。
- P4 新增代码必须写中文注释，说明 handler、DTO、事务边界和错误处理意图。

## Out of Scope

- 不实现 P5 前端页面和组件。
- 不新增 development-plan 以外的 API。
- 不接入真实券商交易 API，不提供自动下单能力。
- 不修改 P0-P3 已归档需求，除非为了满足 P4 API 契约发现明确缺口。
- 不直接编辑 L1 契约文档；若发现契约需修订，只在本 change 的 `specs/` 写 delta。

## Plan Alignment

本 change 与 `docs/development-plan.md` P4 一一对应：

- P4.1 API DTO 与错误处理 → 本变更任务 1.1–1.8。
- P4.2 核心 API → 本变更任务 2.1–2.9。
- P4.2 验收断言 → 本变更任务 3.1–3.14。
- P4 验收命令 → 本变更任务 4.1–4.4。

未加入 `development-plan.md` 之外的新需求；仅将“实现代码要写好中文注释”作为本次用户补充约束纳入任务。
