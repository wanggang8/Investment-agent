# P47: 决策闭环解释

## Summary

新增只读决策闭环解释能力，把一条建议从“生成建议 -> 用户确认/观察/标记 -> 线下记录 -> 风险与复盘线索”的链路聚合为可阅读的解释视图。P47 不新增交易、确认、规则生效或外部推送能力，只从现有本地事实表读取并展示。

## Why

P32 已有每日纪律报告，P35 已有风险预警和 SOP，P36 已有规则效果追踪，P42 已有用户决策工作台，P46 已允许本地知识作为背景检索上下文。用户下一步需要的不只是“看见结果”，而是能快速理解某个建议后来被如何处理、有没有线下记录、哪些审计和复盘事实能解释结果。

## What Changes

- 新增只读 API：`GET /api/v1/decision-loops`。
  - 聚合最近决策的最终裁决、确认状态、确认记录、线下交易记录、错误案例、风险预警链接、审计链和复盘结论。
  - 支持 `symbol` 与 `limit` 查询参数。
  - 不写入 SQLite，不更新状态，不创建通知。
- 新增可选详情 API：`GET /api/v1/decision-loops/{decision_id}`。
  - 返回单条决策闭环的完整 stages、links、missing_links 和 safety_note。
- 新增 `/decision-loop` 前端页面。
  - 展示闭环阶段、缺口、人工处理记录、审计线索和复盘线索。
  - 从工作台和复盘页导航进入。
- 更新文档与 smoke，确保页面可达且不出现高风险入口。

## Scope

- 只读聚合：`decision_records`、`operation_confirmations`、`position_transactions`、`error_cases`、`risk_alerts`、`audit_events`、规则效果追踪和复盘摘要。
- 后端可新增只读 repository 查询方法，但不得新增 migration 或写路径。
- 前端只通过 service 调 API，不直接访问 SQLite、localStorage、sessionStorage 或本地文件。

## Out of Scope

- 券商接口、自动交易、一键交易、代下单。
- 外部推送、自动确认、自动规则应用、自动修复承诺。
- 收益承诺、确定性涨跌预测。
- 新增确认动作、交易写入、规则生效、风险 SOP 状态修改。
- 新增数据库 schema、后台任务、调度器、云同步、登录源、付费源、授权源、Level2 或高频源。

## Validation

- `go test ./...`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `bash scripts/e2e-smoke.sh`
- P47 安全扫描（见 `tasks.md` 7.8）
- `openspec validate p47-decision-loop-explainability --strict`
- `openspec validate --all --strict`
- `git diff --check`
