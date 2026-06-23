## Why

P40 已归档，P33-P40 当前计划内功能队列已经完成。项目需要在继续开发前重新固化 P40 后的路线图、优先级、依赖关系和安全边界，避免在无活跃 change 的状态下直接实现新功能或把新产品愿景混入既有阶段。

本变更用于建立 P41+ 的治理入口：先确定下一轮候选方向、哪些事项只能作为独立治理追溯、哪些边界仍然禁止，然后再为每个后续阶段单独创建 OpenSpec change。

## What Changes

- 新增 P41 后路线图治理规则，明确 P40 后不得直接实现新功能。
- 在 `docs/development-plan.md` 中增加 P41+ 候选队列：产品能力增强、数据质量增强、运维体验增强、历史审计追溯等方向。
- 明确 P19-P24 历史 archive 追溯仍是独立治理事项，不伪造历史归档。
- 更新 `openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`AGENTS.md` 和 `openspec/project.md`，将当前活跃 change 指向本治理包。
- 不修改运行时代码、数据库 schema、HTTP API、前端页面或 L1 契约。

## Capabilities

### New Capabilities

- `roadmap-finalization`: 扩展为覆盖 P40 后路线图治理、候选队列、依赖关系和安全边界。

### Modified Capabilities

无运行时能力修改。

## Impact

- 影响文档：`docs/development-plan.md`、`docs/GOVERNANCE.md`、`AGENTS.md`、`openspec/PROGRESS.md`、`openspec/project.md`。
- 影响 OpenSpec：新增活跃 change `p41-post-p40-roadmap-governance`，并在 `roadmap-finalization` 摘要中记录 P40 后治理规则。
- 不影响后端、前端、数据库 migration、本地调度、LLM 调用或 E2E 行为。
