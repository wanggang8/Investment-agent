# P25 真实公开数据源验证结果

> 验证时间：2026-06-05  
> 范围：只读公开网页与公开 JSON/文件请求；不登录、不绕过权限、不使用付费或授权行情；P25 只输出验证和计划，不实现生产 collector。

## 总体结论

P19/P20 已完成的是本地可配置 HTTP JSON bridge、ETF/基金证据 payload parser、fixture/stub fallback 与信源分级能力；它们不等同于已经接通真实外部公开源。P25 验证后，后续实现应拆成两个阶段：

- P26 公告与证据 collector：优先进入实现计划，首批源建议为巨潮资讯、深交所、证监会；AMAC 行业统计可作为监管/行业背景源；上交所公告入口需继续定位稳定接口后再纳入首批。
- P27 基金净值、ETF 与指数市场数据 collector：可进入实现计划，但应区分权威源与第三方聚合源。首批建议为中证指数基础资料/样本/权重/估值文件与东方财富基金净值/基金档案；新浪财经仅作为 B 级辅助背景或低优先级 ETF 线索源，不作为正式证据或唯一解除信息不足的数据源。

所有源默认低频、只读、可降级；抓取失败写 `audit_events`，不得阻塞本地应用启动，不得触发交易或外部推送。

## 数据源卡片

### 巨潮资讯

- 权威性：上市公司/基金公告披露平台，适合 A 级正式披露证据。
- 访问方式：公开页面与 `POST https://www.cninfo.com.cn/new/hisAnnouncement/query`。
- 已验证请求参数：`pageNum`、`pageSize`、`column`、`tabName`、`stock`、`searchkey`、`category`、`seDate`、`sortName`、`sortType`、`isHLtitle`。
- 已验证字段：`totalAnnouncement`、`totalRecordNum`、`announcements[].secCode`、`secName`、`orgId`、`announcementId`、`announcementTitle`、`announcementTime`、`adjunctUrl`、`adjunctSize`、`adjunctType`、`announcementType`、`hasMore`、`totalpages`。
- 数据频率：公告类按交易日低频轮询；首期建议 30–60 分钟，不做高频。
- 限制与免责声明：需在 P26 实现前补充法律声明、robots 和限频审查；不得绕过反爬或下载限制。
- 失败策略：请求失败、字段缺失或 PDF 获取失败时记录 `source_unavailable` 或 `parse_error` 审计事件，并保留已采集事实。
- 映射：`source_level=A`；公告正文/PDF 为 `formal`，公告标题或列表摘要为 `background`；入库 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。
- P26 判断：可进入首批公告/证据 collector 实现计划。

### 上交所

- 权威性：交易所公开信息，适合 A 级源。
- 访问方式：已验证基金成交概况公开页面 `https://www.sse.com.cn/market/funddata/overview/day/`；页面展示每日基金、ETF、公募 REITs、LOF、交易型货币基金、基金回购等汇总数据。
- 已验证字段：数据日期、挂牌数、成交量、成交金额等页面级字段。
- 未决点：本轮未定位到稳定 JSON 请求；公告、基金公告和更多历史数据入口需要 P26/P27 前继续验证。
- 数据频率：若后续接入，交易日收盘后或每日低频；不得把授权行情服务作为自由 API。
- 限制与免责声明：必须区分公开网页数据与行情信息授权入口。
- 失败策略：未找到稳定接口时不进入自动 collector；若采用页面解析，需标注解析风险并可降级。
- 映射：稳定公告接口确认后可为 `source_level=A`；当前基金成交概况仅建议作为 `background` 或二次验证项。
- P26/P27 判断：暂不作为首批核心 collector；保留为二次验证候选。

### 深交所

- 权威性：交易所公告披露源，适合 A 级正式证据。
- 访问方式：公开页面与 `GET https://www.szse.cn/api/disc/announcement/detailinfo`、`GET https://www.szse.cn/api/disc/announcement/searchQuery`。
- 已验证请求参数：`random`、`pageSize`、`pageNum`、`plateCode`、`annType`。
- 已验证字段：`companyCount`、`announceCount`、`disclosureTip`、`recordCount`、`data[].secCode`、`secName`、`announList[].id`、`title`、`attachPath`、`attachFormat`、`attachSize`、`annId`、`bigCategoryId`、`bigCategoryName`、`publishTime`、`importantRatio`；分类字段包括 `categoryInfo[].value/text`、`plate[].value/text`、`industry[].value/text`。
- 数据频率：公告类交易日 30–60 分钟低频轮询；首期最近 90 天补采。
- 限制与免责声明：P26 实现前补充法律声明、robots 和限频审查；不得把行情授权入口纳入自由采集。
- 失败策略：分页失败、附件失败、分类接口失败分别写审计；列表存在但附件缺失时保留列表证据并标记不完整。
- 映射：`source_level=A`；交易所公告正文/附件为 `formal`，分类和列表摘要为 `background`；入库 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。
- P26 判断：可进入首批公告/证据 collector 实现计划。

### 证监会

- 权威性：监管机构公开信息，适合 A 级监管证据或背景材料。
- 访问方式：公开页面与 `GET https://www.csrc.gov.cn/getChannelList?channelCode=...`、`GET https://www.csrc.gov.cn/searchList/...?_isAgg=true&_isJson=true...`。
- 已验证字段：`data.page`、`data.rows`、`data.channelId`、`data.total`、`data.results[].title`、`content`、`contentHtml`、`memo`、`url`、`publishedTime`、`publishedTimeStr`、`channelName`、`channelCodeName`、`manuscriptId`、`resList[]`、`domainMetaList[]`；metadata 包含发文日期、索引号、文号、发文单位、发布机构、来源、部门等。
- 数据频率：监管规则、公告、处罚和市场禁入类每日低频即可。
- 限制与免责声明：只采公开政府信息页面；不进入登录、办事或报送系统。
- 失败策略：列表失败或详情失败写审计；附件失败不阻塞正文入库。
- 映射：监管规则/处罚/市场禁入可为 `formal`；政策新闻和部门动态为 `background`；入库 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。
- P26 判断：可进入首批公告/证据 collector 实现计划。

### 基金业协会 AMAC

- 权威性：基金业协会行业统计、公示和自律管理源，适合 A 级行业/自律背景源；具体机构/产品查询仍需二次验证。
- 访问方式：已验证数据详情页面 `https://www.amac.org.cn/sjtj/datastatistics/comprehensive/` 存在公开 JSON：`/portal/front/management/assetManage/findDatas`、`/portal/front/management/assetManage/getAllTimes`、`/portal/front/management/assetManage/findDataByTime`。
- 已验证字段：`managerBussTotalStockVOList[].excelTime`、`productNum`、`productScale`；区域分布 `areaName`、`uploadTime`、`quarter`；分类型规模字段如 `raisedFund`、`companySubsidiaryPlan`、`fundManageCompanyPlan`、`fundManageCompanyPension`、`fundSubsidiaryManagePlan`、`futureSubsidiaryCompanyPlan`、`privateEquity`、`corporateAssetSecBuss`。
- 未决点：机构/产品/人员查询页面本轮未暴露稳定 JSON；报送、考试和业务系统不在采集范围。
- 数据频率：行业统计按季度或每日检查更新；自律处分/异常/失联机构可每日低频。
- 限制与免责声明：只采公开统计和公开自律管理栏目；不采登录业务系统、报送系统、考试成绩或个人信息查询。
- 失败策略：行业统计缺字段时保留页面文本与可用 JSON；查询系统不可达时不阻塞其他源。
- 映射：行业统计为 `background`；自律处分、异常经营、失联机构公开公告可按内容标为 `formal` 或 `background`；入库 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。
- P26 判断：行业统计和自律栏目可作为 P26 二线背景源；机构/产品查询需二次验证后再实现。

### 东方财富基金

- 权威性：第三方基金聚合源，默认 B 级；适合净值、基金档案、持仓和辅助市场数据，不作为唯一正式证据。
- 访问方式：基金页面 `https://fund.eastmoney.com/510300.html` 与 `GET https://fund.eastmoney.com/pingzhongdata/510300.js?v=...`。
- 已验证字段：`fS_name`、`fS_code`、`Data_netWorthTrend`、`Data_ACWorthTrend`、`Data_grandTotal`、`Data_rateInSimilarType`、`Data_fluctuationScale`、`Data_holderStructure`、`Data_assetAllocation`、`Data_performanceEvaluation`、`Data_currentFundManager`、`Data_buySedemption`、`Data_fundSharesPositions`、`syl_1n`、`syl_6y`、`syl_3y`、`syl_1y`。
- 样例字段：`Data_netWorthTrend[].x/y/equityReturn/unitMoney`；资产配置包含股票占净比、债券占净比、现金占净比、净资产。
- 数据频率：净值类交易日 21:30 后拉取，并在次日补拉；不采实时估算净值作为正式净值。
- 限制与免责声明：页面同时请求交易/用户相关接口，P27 必须明确排除登录、交易、用户信息和模拟交易推广接口。
- 失败策略：净值缺失时进入 `missing` 或 `insufficient_data`，不得伪造估值百分位；第三方源异常不得覆盖 A 级证据。
- 映射：`source_level=B`；净值/基金档案进入 `market_snapshots.metadata_json` 或行情标准 JSON；公告类只作辅助背景。
- P27 判断：可进入首批基金净值/ETF 市场数据 collector 实现计划，但必须标注 B 级限制。

### 新浪财经

- 权威性：第三方财经门户，默认 B 级辅助背景源。
- 访问方式：主页、行情中心、基金频道公开页面；已验证公开请求包括 `https://statistic.cj.sina.com.cn/api/macd/stocks_by_rule` 和基金推广 JSONP `https://finance.sina.com.cn/tgdata/fund_tg_api.json?...`。
- 已验证字段：`stocks_by_rule` 返回 `dt`、`symbol`、`name`、`type`、`diff`、`macd`、`vol`、`volchg`、`volatility`、`dea`、`bullish`、`day_return`、`year_return`、`created_at`、`updated_at`、`creator`、`is_del`，并包含 ETF 分组。
- 限制：页面包含登录、自选股、Level2、港股 Level2、APP、模拟交易和广告推广入口；这些均不进入采集范围。
- 数据频率：仅在需要市场背景线索时低频拉取；不作为正式净值、公告或权威行情。
- 失败策略：源不可用时直接降级，不影响 A 级源和本地应用。
- 映射：`source_level=B`；只能作为 `background` 或候选市场线索；不得单独解除 `insufficient_data`。
- P26/P27 判断：暂不进入首批核心 collector；可作为后续低优先级背景源候选。

### 中证指数

- 权威性：指数编制与指数资料源，适合 A 级指数基础资料、样本、权重、估值文件候选。
- 访问方式：公开页面与多个 JSON/文件请求，包括 `POST https://www.csindex.com.cn/csindex-home/index-list/query-index-item`、`/csindex-home/indexInfo/index-basic-info/{indexCode}`、`/csindex-home/indexInfo/index-nicons`、`/csindex-home/indexInfo/index-details-data`、`/csindex-home/perf/index-perf-oneday`、`/csindex-home/index/weight/top10new/{indexCode}`、`/csindex-home/index/weight/market-weight/{indexCode}`、`/csindex-home/index/weight/industry-weight-two-new/{indexCode}`。
- 已验证字段：指数列表包含 `indexCode`、`indexName`、`indexNameEn`、`consNumber`、`latestClose`、`monthlyReturn`、`ifTracked`、`indexSeries`、`assetsClassify`、`region`、`currency`、`indexClassify`、`publishDate`、`indexCompliance`；基础信息包含指数全称/简称、RIC、Bloomberg ID、基日、基点、发布日期、调整频率、币种、指数类型、描述、IOSCO 合规标识。
- 已验证文件：编制方案 PDF、factsheet PDF、样本权重 XLS、样本列表 XLS、指数估值 XLS。
- 数据频率：指数资料和样本权重按日或按文件更新检查；盘中行情类只低频或暂不采，避免被误用为实时行情。
- 限制与免责声明：P27 实现前必须补充法律声明、下载许可和频率限制；估值/行情数据需标注来源和用途。
- 失败策略：文件下载失败不阻塞基础信息入库；估值缺失时进入 `insufficient_data`。
- 映射：指数基础资料、样本和权重可为 `source_level=A`；行情/估值按许可边界决定；入库 `market_snapshots.metadata_json`、`intelligence_items`、`rag_chunks` 和审计记录。
- P27 判断：可进入首批指数/ETF 辅助市场数据 collector 实现计划。

## 标准 JSON 设计

### 公告/证据源标准 JSON

```json
{
  "source_name": "cninfo",
  "source_level": "A",
  "source_type": "public_disclosure",
  "evidence_role": "formal",
  "symbol": "510300",
  "title": "公告标题",
  "text": "正文或摘要",
  "url": "https://...",
  "attachment_url": "https://...pdf",
  "published_at": "2026-06-05T00:00:00+08:00",
  "captured_at": "2026-06-05T13:00:00+08:00",
  "content_hash": "sha256:...",
  "raw": {}
}
```

映射要求：

- `intelligence_items`：保存 `source_name`、`source_level`、`title`、`text`、`url`、`published_at`。
- `rag_chunks`：对正文、摘要和附件文本分块；保存 `content_hash` 和来源 URL。
- `source_verifications`：记录源级别、校验状态、失败原因和最近验证时间。
- `audit_events`：记录刷新任务、输入摘要、成功/失败状态、错误码和计数。

### 行情/净值源标准 JSON

```json
{
  "source_name": "eastmoney_fund",
  "source_level": "B",
  "source_type": "fund_nav",
  "symbol": "510300",
  "trade_date": "2026-06-05",
  "nav": 1.2345,
  "accumulated_nav": 1.5678,
  "close_price": null,
  "turnover_rate": null,
  "metadata": {
    "fund_name": "沪深300ETF华泰柏瑞",
    "asset_allocation": {}
  },
  "captured_at": "2026-06-05T21:35:00+08:00",
  "raw": {}
}
```

映射要求：

- `market_snapshots`：只写已有真实字段；第三方净值和文件资料放入 `metadata_json` 或后续明确字段。
- 不得伪造 `pe_percentile`、`pb_percentile`、`volume_percentile`、`volatility_percentile`。
- 缺少权威行情或估值时返回 `missing`、`degraded` 或 `insufficient_data`，不得用 B 级源覆盖 A 级缺失。

## 补采与刷新策略

- 第一阶段补采：最近 90 天公告、监管、自律和净值数据。
- 第二阶段补采：P26/P27 稳定后扩展到最近 1 年。
- 公告类：交易日 30–60 分钟低频轮询；非交易日降频或每日一次。
- 监管/协会：每日低频检查即可。
- 基金净值：交易日 21:30 后拉取，次日补拉一次；不采实时估算净值作为正式净值。
- 指数资料/样本/权重：每日或按文件更新时间检查。
- 所有任务应支持断点、去重、幂等、审计和降级。

## 安全与合规边界

- 只读公开数据，不登录、不绕过权限、不使用付费或授权行情。
- 不接券商交易 API，不实现买入、卖出、撤单、改单或账户自动变更。
- 不高频抓取；不得把页面中的 Level2、用户、自选股、模拟交易、交易/推广接口纳入采集范围。
- 不发送邮件、短信、系统 Push、Webhook、WebSocket 或外部通知。
- 不把第三方聚合源当作唯一正式证据。
- 失败只影响对应数据源，写审计并降级，不阻塞本地应用启动。

## P26/P27 实现判断

### P26 公告与证据 collector

建议进入实现计划。

首批范围：

1. 巨潮资讯公告 collector。
2. 深交所公告 collector。
3. 证监会监管信息 collector。
4. AMAC 行业统计/自律栏目作为可暂缓的二线背景源，不作为 P26 首批必需 collector 或 runtime dependency。

暂缓范围：

- 上交所公告和基金公告：继续定位稳定 JSON 或明确页面解析风险后再实现。
- AMAC 机构/产品/人员查询：本轮未验证稳定查询接口，暂不实现。
- 新浪财经：仅 B 级背景，不进 P26 首批正式证据源。

P26 验收目标：90 天低频补采、标准证据 JSON、入库 `intelligence_items`/`rag_chunks`、source verification、审计、fixture/stub fallback、失败降级。

### P27 基金净值与 ETF 市场数据 collector

建议进入实现计划，但拆分权威指数资料和第三方基金净值。

首批范围：

1. 中证指数基础资料、样本、权重、估值文件 collector。
2. 东方财富基金净值、累计净值、基金档案和资产配置 collector。

暂缓范围：

- 上交所基金成交概况：稳定 JSON 未定位，暂列二次验证。
- 新浪财经 ETF/市场线索：只作为 B 级辅助背景源，不进首批核心行情 collector。

P27 验收目标：90 天或可用历史净值补采、交易日 21:30 后净值刷新、指数资料日级刷新、缺失字段降级、不伪造估值百分位、不把 B 级源当 A 级正式证据。
