# Design: P6 验收加固

## Context

P0–P5 已完成主要工程能力，P6 负责把阶段成果转成可验证的端到端验收材料。项目契约真源仍在 `docs/`，本 change 的 `specs/` 只描述 P6 delta，实际实现以 `docs/development-plan.md` P6 和 `docs/functional-spec.md` A01–A17 为边界。

## Goals / Non-Goals

**Goals:**

- 建立 `docs/testing-plan.md`，逐条覆盖 A01–A17 验收断言。
- 建立或确认配置、migration、seed、本地启动相关文档。
- 使用 `go test ./...` 与 `cd web && npm run build` 验证阶段可交付状态。
- 若验收暴露缺口，只做满足既有契约的实现修正，并为非显然业务逻辑写中文注释。

**Non-Goals:**

- 不新增 P6 计划以外的产品能力。
- 不引入自动交易能力、收益承诺或标的推荐。
- 不新增数据库、工作流、API 或前端契约真源。

## Decisions

- 使用验收清单组织 P6：A01–A17 每条都需要说明前置数据、操作步骤、期望结果和相关验收点，便于人工和自动测试共同使用。
- 配置与 migration 文档分开维护：`docs/configuration.md` 说明运行配置，`docs/migration-plan.md` 说明 schema 初始化、seed 和升级执行方式。
- 保持测试命令与计划一致：后端统一执行 `go test ./...`，前端统一执行 `cd web && npm run build`，不加入计划外工具链要求。
- 代码中文注释只写在非显然处：用于说明业务约束、降级、事务和安全边界，避免重复解释语法。

## Risks / Trade-offs

- A01–A17 跨越后端、前端和文档，可能暴露既有实现缺口 → 仅修正与既有契约不一致的部分，并在 tasks 中记录对应断言。
- 端到端场景可能需要真实外部数据源或 LLM → 测试计划需说明可替代的本地 stub、fixture 或降级验证方式。
- 配置文档容易泄露敏感信息 → 只描述环境变量名和占位值，不写真实密钥。
