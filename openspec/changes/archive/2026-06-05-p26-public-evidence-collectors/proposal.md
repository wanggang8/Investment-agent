# P26 公告与证据源 collector

## Why

P25 已验证真实公开数据源的访问形态、字段、合规边界和 P26/P27 实现范围。当前系统已有 P19/P20 的可配置 HTTP bridge、ETF/基金证据 payload parser、fixture/stub fallback 和信源分级能力，但还没有面向真实公开权威源的生产 collector。

下一步应优先实现公告与监管证据源，因为它们直接服务于买入逻辑破坏、多源验证、证据链和 RAG 检索。P25 结论显示巨潮资讯、深交所和证监会具备首批接入条件；AMAC 行业统计/自律栏目可作为二线背景源。上交所公告、AMAC 机构/产品查询和新浪财经不进入 P26 首批范围。

## What Changes

- 新增首批只读公开证据 collector 计划：巨潮资讯、深交所、证监会。
- 将 AMAC 行业统计/自律栏目列为二线背景源候选；P26 首批实现可暂缓 AMAC，不把它作为必需 runtime dependency。
- 统一公告/监管证据标准 JSON，接入 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。
- 支持最近 90 天低频补采、增量刷新、去重、降级、审计和 fixture/stub fallback。
- 明确首批不实现上交所公告、AMAC 机构/产品/人员查询、东方财富基金净值、中证指数市场数据、新浪财经背景源；这些分别留给二次验证或 P27。

## Out of Scope

- 不接券商交易 API。
- 不实现买入、卖出、撤单、改单或任何自动交易。
- 不登录、不绕过权限、不使用付费或授权行情。
- 不高频抓取，不使用 Level2 或授权市场数据。
- 不发送邮件、短信、系统 Push、Webhook、WebSocket 或外部通知。
- 不把第三方聚合源当作 A 级正式证据。
- 不实现 P27 的基金净值、ETF、指数样本、权重或估值 collector。

## Acceptance

- `tasks.md` 明确 P26 首批 collector 的实现、入库、审计、降级和验收任务。
- delta 明确生产 collector 的边界、数据映射、补采策略和失败行为。
- 首批 collector 不依赖登录、付费、授权行情或浏览器自动化绕过。
- P26 实现后，给定 ETF/fund symbol 可抓最近 90 天相关公告/监管证据，生成 `intelligence_items`、`rag_chunks`、`source_verifications` 和 `audit_events`。
- 源不可用时系统可解释降级，不阻塞本地应用启动，不生成无来源正式结论。
