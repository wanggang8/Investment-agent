# P45: P44 后路线图治理

## Summary

在 P42、P43、P44 与 P19-P24 历史追溯均已完成归档后，重新建立 P44 后的下一轮候选队列、优先级、依赖关系、验收边界和安全边界。该阶段只做路线图与治理，不新增运行时代码、不修改数据库 schema、不新增 HTTP API 或前端页面。

## Why

`openspec/PROGRESS.md` 当前显示 P44 已完成，但 `next_change_id` 为空；部分治理文档仍停留在“下一阶段建议 P44”的历史描述。继续推进新功能前，需要先把 P44 后路线图固化，避免无活跃 change 时直接实现功能，也避免旧文档继续误导后续阶段选择。

## What Changes

- 新增 P45 路线图治理 change，明确 P44 后必须先通过独立 OpenSpec change 再进入任何功能实现。
- 更新 P44 后候选方向，建议拆分为：
  - P46：本地知识库与数据导入治理。
  - P47：组合复盘与决策闭环可解释性增强。
  - P48：数据源覆盖与质量回归包。
  - P49：运维发布与本地升级体验。
- 更新 `openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`openspec/project.md` 和 `docs/development-plan.md`，清理 P44 作为下一阶段的过期描述。
- 明确所有候选仍保持不接券商、不自动交易、不外部推送、不自动确认、不自动应用规则、不承诺收益、不引入登录/付费/授权/Level2/高频源的边界。

## Scope

- 仅限 OpenSpec 治理、路线图、候选阶段和文档一致性。
- 本 change 可以修改 L2/L4 治理与计划文档；不得直接修改运行时功能。

## Out of Scope

- 新增券商接口、自动交易、一键交易、代下单。
- 外部推送（邮件、短信、Webhook、System Push）。
- 自动确认、自动规则应用、自动修复承诺。
- 收益承诺、确定性涨跌预测。
- 新增登录源、付费源、授权源、Level2 或高频数据源。
- 改动 Go/React 运行时代码、SQLite schema、HTTP API 或 Eino 工作流。

## Validation

- `openspec validate p45-post-p44-roadmap-governance --strict`
- `openspec validate --all --strict`
- `git diff --check`
