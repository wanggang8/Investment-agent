# Proposal: 补充 P6 后续总计划

## Summary

补充 P6 之后的剩余阶段计划，把真实数据与 RAG/DeepSeek、前端图表与测试、复盘与本地任务交付纳入 `docs/development-plan.md`、`openspec/PROGRESS.md` 和 `openspec/project.md`。

## Why

P0–P6 已完成 MVP 骨架和验收加固，但需求与架构文档仍包含真实数据源、RAG/VecLite、DeepSeek 分析师、前端图表测试、定时任务和周期复盘等未完全覆盖的产品级能力。需要先把 P7–P9 写入总计划，后续 `/opsx:propose` 才能按阶段自动推进。

## What Changes

- 在 `docs/development-plan.md` 阶段总览、任务分解、依赖关系、总清单和提交节奏中补充 P7–P9。
- P7 聚焦真实数据、RAG/VecLite 与 DeepSeek 分析师接入。
- P8 聚焦前端图表、交互体验和前端测试。
- P9 聚焦 `cmd/agent` 定时任务、月度/季度复盘、规则有效性评估和本地交付。
- 在 `openspec/PROGRESS.md` 中把下一阶段设置为 P7：`p7-real-data-integration`。
- 在 `openspec/project.md` 阶段映射中补充 P7、P8、P9。
- 本变更只补计划和治理文档，不实现功能代码。

## In Scope

- 更新 `docs/development-plan.md`：加入 P7–P9 阶段总览、任务小节、验收命令、依赖关系和建议提交节奏。
- 更新 `openspec/PROGRESS.md`：记录 P7 为下一阶段，补齐 P7–P9 阶段状态。
- 更新 `openspec/project.md`：补充 P7–P9 与建议 change id 映射。
- 对应 specs 只写计划治理 delta，说明总计划必须覆盖 P6 后续剩余能力。
- 后续实现代码仍要求写必要中文注释，说明非显然业务约束、降级、审计和禁止自动交易边界。

## Out of Scope

- 不在本 change 中实现真实数据源、RAG、DeepSeek、前端图表、测试或定时任务。
- 不修改 `docs/requirements.md`、`docs/architecture.md`、`docs/api.md`、`docs/workflow.md`、`docs/frontend-contract.md` 的既有契约内容。
- 不新增 development-plan 之外的产品需求；P7–P9 仅承接 requirements 与 architecture 中已存在但未充分实现的能力。
- 不改变禁止自动交易、DeepSeek 不写最终裁决、规则提案需审计和用户确认等安全边界。

## Capabilities

### New Capabilities

- 无。本变更不引入运行时能力，只补充总计划和治理映射。

### Modified Capabilities

- `development-plan`: 补充 P7–P9 后续阶段计划。
- `openspec-project`: 补充 P7–P9 change 映射。
- `progress-tracking`: 将下一阶段推进到 P7。

## Impact

- 影响文档：`docs/development-plan.md`、`openspec/PROGRESS.md`、`openspec/project.md`。
- 后续 `/opsx:propose` 无参数时可继续读取 `next_change_id=p7-real-data-integration` 创建 P7。
- 不影响现有 Go、前端代码和已归档 P0–P6 change。

## Plan Alignment

本 change 本身是计划补充变更，与现有 `docs/development-plan.md` P0–P6 不冲突：

- P7 真实数据、RAG、DeepSeek → 来自 requirements 的外部数据源、非结构化情报、RAG、DeepSeek 分析材料，以及 architecture 的 `infrastructure/rag`、`llm`、`news`、`analyst`。
- P8 前端图表与测试 → 来自 architecture 的 ECharts、前端测试要求，以及 UI 文档中的驾驶舱展示目标。
- P9 复盘、定时任务、交付化 → 来自 requirements 的每日/月度/季度自检与反馈，以及 architecture 的 `cmd/agent` 任务定位。

未加入需求和架构文档之外的新产品能力；本次只让总计划能够覆盖剩余工作。
