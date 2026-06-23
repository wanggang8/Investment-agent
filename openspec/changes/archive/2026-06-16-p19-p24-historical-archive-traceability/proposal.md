# P19-P24 历史审计追溯治理

## Why

P19–P24（公开 HTTP 数据桥接、A 股 ETF/基金证据解析、应用内通知、规则提案增强、复盘深化、本地运行硬化）在当前文档中为可用能力基线，但未全部保留完整 `openspec/changes/archive/` 包。

在需要可追溯审计时，不能把缺失的 archive 过程伪装成已完成版本，也不能无序回写历史文档。该阶段新增一份专用追溯说明，明确：

- 已交付范围。
- 与当前实施一致的验收边界。
- 当前缺口与后续独立补齐条件。

## What Changes

- 新增 `docs/p19-p24-historical-archive-traceability.md`，集中记录 P19–P24 交付状态、验收边界、历史缺口与审计提示。
- 在该说明中为每个阶段明确当前 archive 状态（已交付 / 未归档 / 可验证项），并禁止新增任何历史能力。
- 补齐 `openspec/changes/p19-p24-historical-archive-traceability/specs/traceability/spec.md`，用于记录本治理变更的行为约束（只读追溯、不可伪造历史、禁止补写假设完成时间）。

## Scope

- 仅补充历史审计追溯说明，不实现运行时功能。
- 不改动 L1 契约（`requirements`、`api`、`data-model`、`workflow`、`frontend-contract`）。
- 不改造既有 P19–P24 功能实现、规则、工作流、API、数据库、前端。

## Out of Scope

- 接入券商 API。
- 自动交易、一键交易、代下单。
- 外部推送（邮件/短信/Webhook/system push）。
- 自动规则确认或自动规则应用。
- 新增收益承诺文本。
- 重新补写/重排 P19–P24 的历史 archive 内容。

## Validation

- `openspec validate p19-p24-historical-archive-traceability --strict`
- `openspec validate --all --strict`
- `git diff --check`
