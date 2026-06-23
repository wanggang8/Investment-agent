# P25 真实公开数据源调研验证任务

## 1. 调研验证候选源

- [x] 验证巨潮资讯公告入口
  - 确认公告列表、基金公告、公告详情是否存在稳定公开请求。
  - 记录分页、筛选参数、字段、日期范围、详情 URL 与 PDF/正文获取方式。
  - 记录首轮限制结论；法律声明、免责声明、robots 和访问频率审查作为 P26 实现前置检查。
  - 结论见 `verification.md`：可进入 P26 首批公告/证据 collector；P26 实现前补充法律声明、robots 和限频审查。

- [x] 验证上交所公开入口
  - 覆盖上市公司公告、基金公告、基金数据、每日交易统计、上证基金网/ETF 入口。
  - 区分公开页面、动态接口、需要授权的信息服务。
  - 记录可采字段、更新频率、失败策略和 source_level。
  - 结论见 `verification.md`：已验证基金成交概况页面字段；稳定 JSON/API 未定位，暂不进入首批核心 collector。

- [x] 验证深交所公开入口
  - 覆盖深市公告、基金入口、公募 REITs、市场概览、行情信息授权入口。
  - 区分公告披露与授权行情，不把“行情信息授权”当作可自由接入 API。
  - 记录可采字段、更新频率、失败策略和 source_level。
  - 结论见 `verification.md`：可进入 P26 首批公告/证据 collector。

- [x] 验证证监会公开入口
  - 覆盖证监会公告、监管规则、行政许可、行政处罚、市场禁入、监管措施。
  - 记录列表页/详情页字段、发布时间、来源机关和归档范围。
  - 标注其适合作为监管背景或正式证据的边界。
  - 结论见 `verification.md`：可进入 P26 首批监管信息 collector。

- [x] 验证基金业协会公开入口
  - 覆盖公募基金市场数据、私募/资管月报、机构/产品公示、自律处分、异常/失联机构。
  - 区分公开公示、登录业务系统和报送系统。
  - 记录可采字段、更新频率和 source_level/evidence_role。
  - 结论见 `verification.md`：行业统计和自律栏目可作为 P26 二线背景源；机构/产品/人员查询需二次验证。

- [x] 验证东方财富基金公开入口
  - 覆盖基金净值、场内基金/ETF、基金档案、基金公告、基金持仓和基金排行。
  - 记录净值更新时间窗口、字段、是否存在稳定 JSON/页面接口。
  - 标注第三方聚合源限制：默认 B 级，不作为唯一正式证据。
  - 结论见 `verification.md`：可进入 P27 首批基金净值/ETF 市场数据 collector，但必须标注 B 级限制。

- [x] 验证新浪财经公开入口
  - 覆盖行情中心、基金频道、ETF 期权、财经新闻、7x24 快讯。
  - 区分公开新闻/背景信息、登录自选、Level2 或授权行情。
  - 标注仅作 B 级背景源，不单独解除信息不足。
  - 结论见 `verification.md`：仅作为 B 级辅助背景或低优先级市场线索源，不进入首批核心 collector。

- [x] 验证中证指数公开入口
  - 使用浏览器 network 或可审计方式确认指数列表、指数详情、样本成分、估值/行情是否有稳定公开请求；法律声明、下载许可和频率限制作为 P27 实现前置检查。
  - 若页面内容为空或无稳定公开请求，记录为 P25 阻塞，不进入 P27 实现范围。
  - 结论见 `verification.md`：可进入 P27 首批指数资料、样本、权重和估值文件 collector。

## 2. 输出接入设计

- [x] 为每个源生成数据源卡片
  - 包含权威性、访问方式、字段、频率、限制、失败策略、source_level、evidence_role、入库路径。
  - 结果见 `verification.md`。

- [x] 设计公告/证据源标准 JSON
  - 映射到 `intelligence_items`、`rag_chunks`、`source_verifications`、`audit_events`。
  - 明确 `source_name`、`source_level`、`source_type`、`title`、`text`、`url`、`published_at`、`captured_at`、`content_hash`。
  - 结果见 `verification.md`。

- [x] 设计行情/净值源标准 JSON
  - 映射到 `market_snapshots` 和现有 market DTO。
  - 明确不得伪造 `pe_percentile` / `pb_percentile`；缺失时进入 `missing` 或 `insufficient_data`。
  - 结果见 `verification.md`。

- [x] 设计补采与刷新策略
  - 第一阶段最近 90 天。
  - 第二阶段最近 1 年。
  - 公告类交易日 30–60 分钟低频轮询；监管/协会每日；基金净值交易日 21:30 后加次日补拉。
  - 结果见 `verification.md`。

- [x] 设计安全与合规边界
  - 只读公开数据。
  - 不登录、不绕过权限、不使用付费/授权行情。
  - 失败进入降级状态并写审计，不阻塞本地应用启动。
  - 结果见 `verification.md`。

## 3. 同步计划文档

- [x] 更新 `docs/development-plan.md`
  - 增加 P25/P26/P27 阶段摘要。
  - 修正 P19/P20 为“基础桥接和 parser 已完成，真实外部源 collector 待 P25+ 验证和实现”。
  - 同步 P25 首轮验证结论和 P26/P27 首批候选范围。

- [x] 更新 `openspec/PROGRESS.md`
  - 将当前活跃变更指向 `p25-real-public-data-sources`。
  - 将 P25 状态设为 `in_progress`，P26/P27 作为后续候选。
  - 同步 P26/P27 首批候选范围。

- [x] 更新 `openspec/project.md`
  - 阶段映射中加入 P25、P26、P27。

- [x] 按需更新 `docs/README.md` 和 `docs/configuration.md`
  - 确保下一轮入口和配置说明不再暗示真实外部源已接通。

## 4. 验收

- [x] 运行 `openspec validate --all --strict`。
- [x] 搜索并确认 P19/P20/P25 口径没有互相冲突。
- [x] 输出 P26/P27 是否可进入实现计划的判断。
