# Proposal: P0 工程骨架

## Intent

建立可本地启动的 Go HTTP 服务与 React 前端空壳，为 P1 数据底座与后续 API 提供工程基础。不实现业务领域逻辑。

## Scope

### In scope

- P0.1：Go module、`cmd/server`、`/api/v1/health`、配置读取、logger、`docs/configuration.md`
- P0.2：React + Vite + TS、基础路由骨架、API client 封装、全局状态样式占位

### Out of scope

- SQLite / migration / Repository（P1）
- 领域规则、Eino Graph（P2–P3）
- `docs/api.md` 中除 health 外的业务 API（P4）
- 真实页面数据与 Agent 工作流

## References

- `docs/development-plan.md` § P0
- `docs/architecture.md` § 目录结构与技术栈
- `docs/ui-design.md`、`docs/frontend-contract.md`（P0.2 路由命名）

## Approval

- [ ] 提案已审阅（范围与 out of scope 确认）
