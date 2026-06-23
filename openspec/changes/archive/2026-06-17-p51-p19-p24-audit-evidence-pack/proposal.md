# P51 P19-P24 审计证据包

## Why

P19-P24 已作为可用 MVP 能力交付，并在 `docs/development-plan.md`、`openspec/project.md`、`openspec/PROGRESS.md` 中标记为已交付或 done；但这些阶段没有逐阶段完整 `openspec/changes/archive/` 包。P50 已明确：进入发布候选材料前，应先补一份基于当前仓库事实的 P19-P24 审计证据包。

本阶段目标是把“能力已交付但缺完整 archive 包”的状态整理为可核对证据矩阵，避免后续发布或验收材料引用不可追溯结论。

## What Changes

- 新增 `docs/p19-p24-audit-evidence-pack.md`，按 P19-P24 分阶段记录：
  - 当前交付边界。
  - archive 状态与不可声明事项。
  - 文档证据、代码证据、测试证据、可重跑命令。
  - 依赖的后续阶段补强（如 P25-P29 真实源、P30/P39 E2E、P40/P44/P49 运维验证）。
  - 残余缺口与发布前处理建议。
- 明确 P14-P18 已有标准 archive；若请求范围被写成“P19-P14”，P51 只复核 P14-P18 archive 存在性，重点仍是 P19-P24 缺口。
- 更新治理、进度和开发计划，标记 P51 活跃，下一阶段指向 P52 `p52-project-acceptance-gate-matrix`。
- 增加 OpenSpec 行为摘要，约束 P19-P24 审计证据包不得伪造历史 archive、不得新增运行时能力、不得扩大安全边界。

## Scope

- 仅修改文档、OpenSpec change、OpenSpec specs 摘要和进度状态。
- 不修改后端、前端、数据库 schema、迁移、API 行为、Eino 工作流、脚本行为或测试代码。
- 不接入新的数据源，不跑真实外部源，不新增测试自动化。

## Out of Scope

- 补建伪历史 P19/P20/P21/P22/P23/P24 archive 包。
- 改写 P19-P24 的完成时间或历史任务记录。
- 宣称 P19/P20 已接通所有真实外部源。
- 建立全项目验收门禁矩阵；该工作留给 P52。
- 生成发布候选材料；该工作留给 P53。
- 接入券商 API、自动交易、一键交易、代下单、外部推送、自动确认、自动规则应用、自动修复承诺、自动覆盖真实库、收益承诺、登录源、付费源、授权源、Level2 或高频源。

## Validation

- `openspec validate p51-p19-p24-audit-evidence-pack --strict`
- `openspec validate --all --strict`
- `git diff --check`
- 只读复审：确认 P51 证据包基于仓库事实，不把缺 archive 改写为已 archive，不扩大安全边界。
