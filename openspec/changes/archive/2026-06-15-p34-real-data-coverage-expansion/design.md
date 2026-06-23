## Context

P26 已接入首批公告/监管公开证据 collector，P27 已接入东方财富基金基础市场数据和中证指数基础信息 collector，P29 已补齐公开证据真实 smoke、显式日期窗口、错误分类和 no-data 语义。P33 已完成本地账户/持仓 onboarding，使每日纪律具备用户持仓上下文。

P34 需要把真实公开数据覆盖从“首批可用源”扩展到“能支撑日常纪律和后续风险预警的更完整公开数据上下文”。范围仍保持本地、只读、低频、可降级，不接券商、登录、付费、授权、Level2 或高频源。

## Goals / Non-Goals

**Goals:**

- 扩展中证指数样本、权重、估值文件等公开数据读取能力，并保持 endpoint shape 可校准、可降级。
- 评估并接入一批公开可用的成分股财务、资金流向、融资融券或可替代情绪指标。
- 标准化新增数据的 freshness、missing、stale、no_data、source_unavailable、parse_error 状态。
- 将新增数据纳入每日纪律、expected return 和后续风险预警可读取的上下文。
- 在前端或运维状态中展示数据源健康和最近成功/失败记录。

**Non-Goals:**

- 不接券商账户、交易 API、自动交易、一键交易、撤单、改单或外部推送。
- 不接登录、付费、授权、Level2、高频行情或需要绕过访问控制的数据源。
- 不承诺收益，不预测确定涨跌，不用缺失数据伪造风险或收益结论。
- 不在 P34 完成完整风险 SOP 编排；风险中心属于 P35。
- 不在 P34 完成 RAG 召回质量评估体系；该能力属于 P38。

## Decisions

1. **优先扩展现有 collector 与 market refresh 管道。**
   - 新增数据先进入现有 source fetch → normalize → persistence → audit 模式。
   - 市场型指标优先写入 `market_snapshots.market_metrics_json` 或轻量 source health 结构；只有现有模型无法表达健康状态或历史记录时才新增表。
   - 原因：P26/P27/P29 已验证现有管道可实现幂等、降级和审计。

2. **按数据类别定义 source contract。**
   - 指数样本/权重/估值、成分财务、资金流向、融资融券、情绪替代指标分别定义字段、日期、新鲜度、source level 和失败分类。
   - 对无法稳定公开获取的类别，必须以 `no_data`、`source_unavailable` 或 `parse_error` 表达，不伪造字段。
   - 原因：不同公开源的更新周期和可信等级不同，不能混用成单一“行情成功”。

3. **新增数据只作为分析上下文，不改变交易边界。**
   - DailyDisciplineGraph、ExpectedReturnNode 和未来 P35 风险预警可读取新增上下文。
   - 任何输出仍需规则裁决、人工复核和已有确认机制。
   - 原因：真实数据覆盖增强的是材料质量，不是自动执行能力。

4. **真实 smoke 与 fixture/stub 双轨。**
   - 每个新增真实源需要 fixture 测试和可选真实 smoke。
   - 默认本地验收不得依赖公网；真实 smoke 必须显式启用日期窗口、标的和 source。
   - 原因：保持开发确定性，同时保留真实接口形状校准能力。

5. **健康状态面向运维和前端展示。**
   - 每次刷新记录最近成功时间、失败类别、影响范围、源等级和数据日期。
   - 前端只展示状态和下一步，不外部推送、不触发交易。
   - 原因：P34 的使用价值依赖用户知道哪些公开数据可用、缺失或过期。

## Implementation Notes / Archive Deltas

- API / DTO：P34 暴露 `GET /api/v1/market/source-health`，返回 source category、freshness、source level、source type、data date、last success/failure、failure category 与 affected symbols；每日纪律报告 DTO 增加 `p34_source_coverage`，用于展示 supporting data summary、missing categories 和结构化 source health items。
- Refresh 入口：`cmd/agent --task p34-expanded-refresh --source <source> --symbol <symbol-or-index> --start-date YYYY-MM-DD --end-date YYYY-MM-DD` 作为显式本地刷新入口；当前支持 `sentiment_proxy_fixture`、`csindex_extended` 和 `configured`，未知 source 返回 bad request；fixture smoke 使用 `sentiment_proxy_fixture`，真实源 smoke 需显式选择公开源。
- 数据模型：优先复用 `market_snapshots.market_metrics_json` 保存 `p34_data_categories`、`p34_source_health` 和 source metadata；审计写入 `audit_events`；未新增交易、券商或外部推送表。
- 工作流：Daily Discipline / Expected Return 读取 P34 supporting data summary 和 missing categories；缺失、过期或失败类别保持诊断信息，不转换为收益承诺或交易动作。
- 前端：Settings 能展示 P34 source health；每日纪律报告详情展示 P34 扩展数据覆盖状态，并保留不接券商、不触发交易的边界文案。
- Source 评估：指数样本、权重、估值沿中证指数公开候选 metadata 低频读取；source health 持久化保留 `data_date`、成功/失败时间、失败分类、影响标的和 source level；首批情绪替代指标以本地 fixture/stub 验收，`stubbed` 不写真实成功时间，使用 `failure_category: stubbed` 表示非真实公开源；真实公开源不可用时记录 `no_data`、`source_unavailable` 或 `parse_error`，不伪造数据。
- 验收证据：2026-06-15 已显式运行 `p34-expanded-refresh --source sentiment_proxy_fixture --symbol 000300 --start-date 2026-06-01 --end-date 2026-06-05` 与 `p34-expanded-refresh --source csindex_extended --symbol 000300 --start-date 2026-06-01 --end-date 2026-06-05`；两者均完成并写入 `audit_events`，输出确认不会执行交易。

## Risks / Trade-offs

- [Risk] 公开源 shape 频繁变化。→ 使用 source-specific parser、fixture、真实 smoke 和 `parse_error` 分类，不把解析失败混为无数据。
- [Risk] 扩展过多源导致范围失控。→ 首轮只接 P34 明确类别中的可验证公开源；候选失败可记录为降级，不阻塞可用源。
- [Risk] B 级第三方数据被误当 A 级正式证据。→ source level 必须随 payload 保存，B/C 级不得单独解除信息不足。
- [Risk] 日常刷新被误解为实时行情。→ 保持低频任务、公开数据说明和 stale/missing 标记。
- [Risk] 数据进入 expected return 后被误读为收益承诺。→ expected return 继续输出情景材料和置信/样本说明，不生成确定性涨跌预测。
