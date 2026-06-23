# P32 每日纪律报告产品化

## Why

P31 已完成本地每日自动运行闭环：系统可以在显式启用后串联市场/证据刷新、每日纪律 workflow、通知、审计和前端状态展示。但 P31 的结果仍偏向运行状态与最新结果入口，尚未形成稳定的报告产品层：用户难以直接看到“今日纪律报告”、历史报告列表、报告详情、缺前提诊断和可回溯索引。

每日纪律报告是 P31 之后的产品化层。P32 将把 daily workflow 与 P31 auto-run 结果沉淀为轻量 `daily_discipline_reports` 索引，并提供 today/list/detail API、今日纪律页升级、历史报告与详情页，以及本地 E2E smoke 验收，使用户能以报告为中心完成每日复核。

## What Changes

- 新增轻量 `daily_discipline_reports` 索引，记录本地日期、持仓 scope、报告状态、关联决策/证据/审计、摘要、缺失前提和幂等键。
- 新增聚合 API：today、list、detail，用于查询今日纪律报告、历史报告列表和报告详情。
- 升级今日纪律页，使其优先展示今日报告状态、摘要、缺失前提、证据/审计链接和人工复核边界。
- 新增历史报告与详情展示，支持用户按日期回看每日纪律结果、降级/失败原因和关联材料。
- 增加 seed、handler/service/repository tests、前端 tests 和本地 E2E smoke，验证成功报告、缺前提、历史列表和详情路径。

## In Scope

- 后端模型、migration、repository、wiring 和幂等写入/复用行为。
- 报告聚合 DTO 与 today/list/detail handlers。
- 今日纪律页产品化升级、历史报告列表/详情页、路由与导航入口。
- smoke seed 与 Playwright E2E，覆盖今日报告、历史报告和详情可达。
- 文档、进度和验收命令同步。

## Out of Scope

- 不新增交易执行、券商接口、买入/卖出/撤单/改单请求或任何自动交易能力。
- 不新增外部推送渠道，包括邮件、短信、系统 Push、Webhook、WebSocket 或第三方通知。
- 不新增登录源、付费源、授权源、Level2 或用户身份行情源。
- 不新增高频抓取；报告聚合必须保持低频、本地、可审计和幂等。
- 不承诺收益，不预测确定涨跌，不让模型覆盖规则裁决。
- 不扩展真实数据源广度；缺失数据只记录为报告状态或缺前提诊断。

## Impact

- 将 P31 的自动运行结果转化为可持续使用的报告产品体验。
- 增强每日纪律结果的可发现性、历史可回溯性和缺前提透明度。
- 需要修改 persistence、workflow/service/handler、frontend pages/services/types、seed、E2E 和文档/进度。
