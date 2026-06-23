# P27 基金净值与 ETF 市场数据 collector 任务

## 1. 前置确认

- [x] 复核 P25 验证结论
  - 已阅读 `openspec/changes/archive/2026-06-05-p25-real-public-data-sources/verification.md`。
  - 确认 P27 首批范围仅包含中证指数和东方财富基金。
  - 确认上交所基金成交概况、 新浪财经 ETF/市场线索、登录/交易/用户信息接口、Level2 或授权行情不进入首批实现。

- [x] 复核法律声明、robots 和访问频率
  - 已对中证指数、东方财富基金公开页面或公开文件补充只读、低频、公开访问边界检查。
  - 东方财富基金 robots 未见禁止 `/pingzhongdata/` 或基金页面的直接规则；中证指数 robots 返回空内容，许可和文件下载边界仍需在实现中保持低频、只读、可禁用和人工复核约束。
  - 若发现禁止自动化采集、登录限制、验证码、付费、授权或 Level2 限制，将对应源移出首批实现范围。

## 2. Collector 设计

- [x] 定义市场数据 collector 接口
  - 已复用 `MarketDataSource.FetchMarketData(ctx, symbol)` 作为 P27 collector 边界；当前触发路径按 `symbol` 低频只读刷新，扩展端点作为可选 metadata 抓取。
  - 已返回标准基金净值/ETF/指数市场数据 payload，包含 `source_name`、`source_level`、`source_type`、`symbol`、`trade_date`、`nav/close_price`、`accumulated_nav`、`metadata`、`captured_at`、`content_hash` 和 `raw` metadata。

- [x] 设计 freshness 与降级策略
  - 基础行情缺失、净值缺失、字段缺失、下载失败和解析失败统一返回 `DATA_SOURCE_UNAVAILABLE`，由 composite source 降级到后续 collector 或 fallback。
  - 可选扩展 metadata 端点失败时不阻断基础行情，体现为 degraded metadata；B 级东方财富源不伪造估值或百分位字段，也不替代 A 级正式证据。

- [x] 设计去重与幂等策略
  - 已使用 `source_name + symbol + trade_date + source_type` 生成稳定市场快照 ID；重复刷新不会重复写入市场事实。
  - collector payload 同时包含 `content_hash` 以支持后续文件/RAG 去重；内容 metadata 变化保留在 `market_metrics_json` 与审计路径中。

## 3. 首批源实现范围

- [x] 实现中证指数 collector
  - 已实现指数基础信息 JSON collector，输出 `source_name=csindex`、`source_level=A`、`source_type=index_basic` 与 metadata，并可通过 market-refresh 写入 `market_metrics_json`。
  - 已实现可选扩展 metadata 抓取：指数样本、权重和估值文件候选；单个扩展端点失败时不阻断基础行情。

- [x] 实现东方财富基金 collector
  - 已实现 `pingzhongdata/{symbol}.js` 最近净值解析，输出基金名称、代码、单位净值、累计净值、收益率、`source_level=B` 和 `source_type=fund_nav`，并可通过 market-refresh 写入 `market_metrics_json`。
  - 已实现可选扩展 metadata 解析：历史净值、累计净值、资产配置、业绩评价和基金经理基础档案；B 级源仍不伪造估值分位、不替代 A 级正式证据。

- [x] 明确暂缓源
  - 上交所基金成交概况：稳定 JSON 未定位，暂列二次验证。
  - 新浪财经 ETF/市场线索：仅作为 B 级辅助背景源，不进首批核心行情 collector。
  - 登录、交易、用户信息、模拟交易推广、Level2 或授权行情接口一律不实现。

## 4. 本地触发与配置

- [x] 扩展配置项
  - 已增加 P27 市场数据 collector 的启用开关、sources 和 base URL 配置。
  - 默认关闭真实 collector，保留 stub/fixture fallback。
  - 校验 sources 只允许首批已实现源。

- [x] 增加本地任务或复用 market-refresh 验收路径
  - 已明确复用 `cmd/agent --task market-refresh --symbol <symbol>` 触发 P27 collector；无需新增交易或外部通知任务。
  - 失败时沿用 market-refresh 可解释错误与审计路径。

## 5. 测试与验收

- [x] 增加 collector 单元测试
  - 已使用 httptest 覆盖中证指数基础信息 JSON、东方财富基金 `pingzhongdata` JS、source_level/source_type、standard payload 字段、fallback、B 级源不伪造百分位，以及空数据、字段缺失和解析失败场景。

- [x] 增加入库集成测试
  - 已验证 collector source metadata、captured_at 与 content_hash 可写入 `market_snapshots.market_metrics_json`。
  - 已验证同一 source/symbol/trade_date/source_type 重复刷新不会重复写入市场事实，并保留成功审计路径。
  - B 级源不能单独解除需要 A 级权威数据的 `insufficient_data` 约束由现有证据质量规则保持，不在 market collector 中绕过。

- [x] 运行后端测试
  - `go test ./...`

- [x] 运行 OpenSpec 校验
  - `openspec validate --all --strict`

## 6. 文档同步

- [x] 更新 `docs/development-plan.md`
  - 已将 P27 当前进行中范围、首批源和验收目标写清楚。

- [x] 更新 `docs/configuration.md`
  - 已补充 P27 collector 配置项和默认关闭说明；未写真实密钥、登录凭证或授权行情配置。

- [x] 按需更新 `docs/data-model.md`、`docs/api.md`、`docs/frontend-contract.md`
  - 已在 `docs/data-model.md` 说明 P27 source metadata 复用现有 `market_snapshots.market_metrics_json`，以及 source/symbol/trade_date/source_type 幂等边界。
  - 本轮未新增 HTTP API 或前端契约，因此无需更新 `docs/api.md`、`docs/frontend-contract.md`。

- [x] 更新 `openspec/PROGRESS.md`、`docs/GOVERNANCE.md`、`AGENTS.md`
  - 已同步 P27 实现完成、active change 待归档前复审的状态。

## 7. 归档前复审

- [x] 子 agent 规格与代码复审
  - 已修复复审发现的关键问题：P27 collectors 优先于通用 endpoint/stub fallback；P27 source 不回填输入默认 PE/PB 分位；仅启用 market collectors 时配置校验可通过；P27 path-based 请求不追加 `?symbol=` query；wiring sources 做 trim 后匹配。
