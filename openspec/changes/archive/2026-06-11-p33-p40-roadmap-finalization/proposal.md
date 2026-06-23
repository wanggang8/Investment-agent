## Why

P32 已归档后，项目进入 P33–P40 的最后一组产品化阶段。当前 `docs/development-plan.md` 已列出阶段标题和高层目标，但还缺少统一的执行顺序、依赖关系、验收边界和“是否仍有遗漏”的治理判断。

本变更用于把 P33–P40 固化为后续唯一计划内功能队列，并明确归档、历史追溯和 P40 后新路线图的边界。

## What Changes

- 将 P33–P40 标记为当前剩余计划内功能阶段，并给出推荐执行顺序。
- 明确 P33–P40 每阶段的范围、依赖、验收方向和安全边界。
- 明确不属于 P33–P40 的治理事项：P19–P24 历史 archive 追溯、P40 后新产品路线图。
- 更新进度与治理文档，使无参数 `/opsx:propose` 可进入下一阶段规划。
- 不修改交易、安全、数据源等级或 LLM 裁决边界。

## Capabilities

### New Capabilities
- `roadmap-finalization`: 覆盖 P33–P40 后续开发路线图、阶段依赖、验收边界和遗漏判断。

### Modified Capabilities

无。

## Impact

- 影响文档：`openspec/PROGRESS.md`、`docs/development-plan.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/project.md`。
- 影响 OpenSpec：新增 `p33-p40-roadmap-finalization` 活跃 change，并作为启动 P33 前的治理入口。
- 不影响后端、前端、数据库 migration 或运行时行为。
