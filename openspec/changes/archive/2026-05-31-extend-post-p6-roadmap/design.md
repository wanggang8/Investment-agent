# Design: P6 后续路线图补充

## Context

P0–P6 已完成工程骨架、数据底座、领域规则、工作流、HTTP API、前端驾驶舱和端到端验收加固。当前 `openspec/PROGRESS.md` 没有下一阶段，导致无参数 `/opsx:propose` 无法继续推进。

需求和架构文档仍保留若干产品级能力：真实行情与情报数据、RAG/VecLite、DeepSeek 分析材料、图表化前端体验、前端测试、`cmd/agent` 本地任务、月度/季度复盘与规则有效性评估。它们不适合作为一个超大 change 一次实现，应先补入总计划，再按阶段处理。

## Goals / Non-Goals

**Goals:**

- 在总计划中补 P7、P8、P9，并明确每阶段目标、任务和验收方式。
- 让 OpenSpec 进度文件恢复可推进状态，下一阶段指向 `p7-real-data-integration`。
- 在 OpenSpec 项目说明中补阶段与 change id 映射。
- 保持 P0–P6 已完成记录不变。
- 只写计划治理 delta，archive 时再合并到对应文档。

**Non-Goals:**

- 不实现 P7–P9 的功能代码。
- 不新增需求和架构文档之外的产品能力。
- 不修改 L1 契约正文。
- 不改变禁止自动交易、DeepSeek 不写最终裁决、规则提案需用户确认等边界。

## Decisions

### Decision 1: 拆为 P7、P8、P9 三个阶段

选择三阶段拆分：

- P7：真实数据、RAG/VecLite、DeepSeek 分析师。
- P8：前端图表、体验增强、前端测试。
- P9：定时任务、周期复盘、规则有效性评估、本地交付。

原因：这三类工作依赖关系清晰，风险类型不同。P7 先补数据与分析能力，P8 在稳定 API 上增强展示，P9 再做周期化和交付化。

备选方案是合并成一个 P7，但会让任务过大，验收边界不清。另一个方案是细分为更多阶段，但当前项目规模下会增加治理成本。

### Decision 2: 本 change 只修改规划与治理文档

本 change 只让路线图恢复完整，不触碰业务实现。后续每个阶段仍通过独立 change 生成 proposal、design、specs 和 tasks。

原因：计划补充本身是治理工作；若同时实现功能，容易让规划变更和运行时能力混在一起。

### Decision 3: P7 作为下一阶段

`openspec/PROGRESS.md` 将设置：

- `current_phase`: `P6`
- `status`: `done`
- `next_phase`: `P7`
- `next_change_id`: `p7-real-data-integration`

原因：P7 是后续能力的基础，P8/P9 依赖 P7 提供更可信的数据与分析材料。

### Decision 4: specs 使用计划治理 delta

创建或修改以下 delta：

- `development-plan`：总计划必须覆盖 P7–P9。
- `openspec-project`：项目阶段映射必须包含 P7–P9。
- `progress-tracking`：进度文件必须记录 P7 为下一阶段。

这些 delta 只约束治理文档，不替代 requirements、architecture、api、workflow、frontend-contract 的 L1 契约。

## Risks / Trade-offs

- P7–P9 粒度仍可能偏大 → 每阶段 propose 时再按 tasks 拆为可验收小项。
- 外部数据与 DeepSeek 接入存在环境差异 → P7 计划中必须包含可替换 mock、缺配置降级和审计记录。
- 前端增强可能超出契约范围 → P8 只围绕现有 UI 与 frontend-contract 做展示和测试。
- 定时任务涉及本地长期运行 → P9 需要明确本地开发、手动触发和异常可见性，不把它包装成自动交易能力。

## Migration Plan

1. `/opsx:apply` 时更新 `docs/development-plan.md`、`openspec/PROGRESS.md`、`openspec/project.md`。
2. 归档时将 delta 合并入正式文档。
3. 后续无参数 `/opsx:propose` 应读取 P7 并创建 `p7-real-data-integration`。
4. 若发现 P7 范围仍过宽，可在 P7 propose 阶段继续拆分，但不得删除 P8/P9 总计划。

## Open Questions

无。本变更只补总计划；具体供应商、数据源频率、DeepSeek prompt 与前端测试框架细节留到 P7/P8/P9 的独立 change 中决定。
