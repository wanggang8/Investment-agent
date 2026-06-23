# P27 基金净值与 ETF 市场数据 collector

## Why

P25 已验证 P27 可以进入实现计划，但必须拆分权威指数资料和第三方基金净值：中证指数适合作为 A 级指数基础资料、样本、权重和估值文件候选；东方财富基金适合作为 B 级基金净值、累计净值、基金档案和资产配置辅助市场数据源。

当前系统已有 P19 的可配置 HTTP market bridge 和 P26 的公告/监管证据 collector，但还没有面向基金净值、ETF 与指数日频资料的生产 collector。P27 应补齐只读、低频、本地入库的市场数据 collector，让 `market-refresh` 能在公开源可用时写入真实市场事实，并在缺失或过期时保持 `missing`、`stale`、`degraded` 或 `insufficient_data`，不得伪造估值百分位或把 B 级第三方源当 A 级正式证据。

## What Changes

- 新增首批只读市场数据 collector 计划：中证指数、东方财富基金。
- 将中证指数映射为 A 级指数资料、样本、权重和估值文件候选源。
- 将东方财富基金映射为 B 级基金净值、累计净值、基金档案和资产配置辅助源。
- 标准化基金净值、ETF 与指数市场数据 payload，写入 `market_snapshots.metadata_json`、必要的 `intelligence_items` / `rag_chunks` 和 `audit_events`。
- 支持最近 90 天或公开可得历史净值补采、交易日 21:30 后净值刷新、次日补拉、指数资料日级刷新、去重、幂等、降级和 fixture/stub fallback。
- 明确首批不实现上交所基金成交概况、登录/交易/用户信息接口、Level2 或授权行情、实时估算净值正式化、券商接口和外部通知。

## Out of Scope

- 不接券商交易 API。
- 不实现买入、卖出、撤单、改单或任何自动交易。
- 不登录、不绕过权限、不使用付费、授权、Level2 或需用户身份的行情。
- 不高频抓取，不采集实时估算净值作为正式净值。
- 不发送邮件、短信、系统 Push、Webhook、WebSocket 或外部通知。
- 不把东方财富基金等第三方 B 级聚合源当作 A 级正式证据或唯一解除信息不足的依据。
- 不实现 P26 已完成的公告/监管 collector 重构。

## Acceptance

- `tasks.md` 明确 P27 首批 collector、入库、审计、降级和验收任务。
- delta 明确生产 market collector 的安全边界、数据映射、补采策略和失败行为。
- 首批 collector 不依赖登录、付费、授权行情、Level2、券商接口或浏览器自动化绕过。
- P27 实现后，给定 ETF/fund symbol 可在 fixture/httptest 中获取最近交易日净值或指数资料，写入本地事实并生成审计。
- 缺失字段、源不可用或 B 级源不足时系统可解释降级，不伪造 `pe_percentile`、`pb_percentile`、`volume_percentile` 或 `volatility_percentile`，不生成无来源正式结论。
