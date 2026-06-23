# Proposal: P6 验收加固

## Summary

创建 `docs/development-plan.md` 中 P6 验收加固变更包，用于补齐端到端验收场景、配置与启动文档，并以全量后端测试和前端构建作为阶段验收。

## Why

P0–P5 已完成工程、数据、规则、工作流、HTTP API 与前端驾驶舱，P6 需要把 `docs/functional-spec.md` 中 A01–A17 可测试验收断言形成可执行验收清单。该阶段同时补齐配置、migration、seed 与本地启动文档，便于后续归档和交付检查。

## What Changes

- 创建或确认 `docs/testing-plan.md`，覆盖 A01–A17 可测试验收断言。
- 为首次使用、每日纪律、证据不足、VecLite 降级、能力圈外、用户确认、错误标记、C 级信源、LLM 不可用、规则提案、审计语义、禁止自动交易、市场刷新、预期收益评估建立验收说明。
- 创建或确认 `docs/configuration.md` 与 `docs/migration-plan.md`，覆盖 SQLite、VecLite、DeepSeek API Key、数据源开关、日志级别、migration、seed 和本地启动命令。
- 执行阶段验收命令：`go test ./...` 与 `cd web && npm run build`。
- 本阶段新增或调整的实现代码需写必要中文注释，说明非显然的业务边界、降级原因、事务语义和禁止自动交易约束。

## In Scope

- P6.1：创建或确认 `docs/testing-plan.md`。
- P6.1：逐条覆盖 `docs/functional-spec.md` 的 A01–A17 可测试验收断言。
- P6.1：保留并执行 `go test ./...` 与 `cd web && npm run build` 作为阶段验收。
- P6.2：创建或确认 `docs/configuration.md` 与 `docs/migration-plan.md`。
- P6.2：文档覆盖 SQLite 数据文件路径、VecLite 索引文件路径、DeepSeek API Key 环境变量、数据源开关、日志级别、migration 执行方式、seed 数据说明和本地启动命令。
- 本阶段若需调整实现以满足验收，只限于 P6 验收断言、配置启动文档和既有契约要求范围内。

## Out of Scope

- 不新增 `docs/development-plan.md` P6 以外的产品需求、页面、API、数据表或工作流。
- 不修改自动交易边界；系统仍不提供交易执行接口和一键交易入口。
- 不让 DeepSeek 生成最终裁决，不承诺收益，不主动推荐具体标的。
- 不把 `docs/testing-plan.md`、`docs/configuration.md` 或 `docs/migration-plan.md` 写成新的 L1 契约真源。
- 不在未归档 change 期间直接修改 L1 契约内容；如涉及契约变化，仅在本 change 的 `specs/` 写 delta。

## Capabilities

### New Capabilities

- `e2e-hardening`: 覆盖 P6 端到端验收场景、配置与启动文档、阶段验收命令的本阶段增量契约。

### Modified Capabilities

- 无。当前 `openspec/specs/` 没有可修改的既有 capability；本变更只写 P6 阶段 delta。

## Impact

- 影响文档：`docs/testing-plan.md`、`docs/configuration.md`、`docs/migration-plan.md`。
- 可能影响后端测试、前端构建或少量实现修正，但仅用于满足 A01–A17、配置启动说明和既有契约。
- 验收以 `docs/development-plan.md` P6 命令为准：`go test ./...` 与 `cd web && npm run build`。

## Plan Alignment

本 change 与 `docs/development-plan.md` P6 小节一一对应：

- P6.1 端到端场景 → tasks.md 第 1 节，覆盖 `docs/testing-plan.md` 与 A01–A17 全部验收断言。
- P6.1 验收命令 → tasks.md 第 2 节，保留 `go test ./...` 与 `cd web && npm run build`。
- P6.2 配置与启动文档 → tasks.md 第 3 节，覆盖 `docs/configuration.md`、`docs/migration-plan.md` 和全部必须包含项。

未加入 `docs/development-plan.md` 以外的需求；“实现的代码要写好中文注释”只作为 tasks.md 执行约束，不写入 specs 契约 delta。
