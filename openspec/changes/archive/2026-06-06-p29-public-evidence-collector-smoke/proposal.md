# Proposal: P29 公开证据 collector 真实采集验收修复

## 背景

P26 已实现巨潮资讯、深交所、证监会公开证据 collector 及入库链路，但本地真实 smoke 复核发现：

- `market-refresh --symbol 510300` 在真实 market collectors 配置下可写入 `market_snapshots`，且来源为 `eastmoney_fund`，不是 stub。
- `public-evidence-refresh --symbol 510300` 在真实公开证据配置下未能写入证据表：巨潮/深交所返回无公告，证监会当前接口返回 HTTP 404。

因此 P26 不能继续只以单元测试/fixture 判断为“真实采集稳定可用”。本变更聚焦公开证据 collector 的真实可运行性、可诊断性和 smoke 验收。

## 目标

- 重新核验巨潮资讯、深交所、证监会当前公开接口参数和响应 shape。
- 修复至少一个 A 级公开证据源，使其在固定 smoke 标的和时间窗口内可真实采集并写入 SQLite。
- 明确区分“源接口不可用”“当前窗口无公告/无监管信息”“解析失败”。
- 增加可复现 smoke 验收方式，证明真实采集结果进入 `intelligence_items`、`intelligence_summary`、`rag_chunks`、`source_verifications` 和 `audit_events`。

## 非目标

- 不接券商交易 API。
- 不自动交易，不创建订单、确认、交易或外部通知。
- 不登录，不绕过权限，不使用付费/授权行情、Level2、浏览器爬虫或高频抓取。
- 不把 B 级第三方来源当作 A 级正式证据。
- 不要求所有公开源都在任何标的/任何窗口内有数据；无数据必须可解释。
