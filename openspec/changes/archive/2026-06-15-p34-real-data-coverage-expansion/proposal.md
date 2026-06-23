## Why

P33 已让空库用户建立本地账户与持仓事实，但每日纪律、风险判断和 expected return 仍主要依赖 P26/P27 已接入的首批公开源。P34 需要扩展真实公开数据覆盖面，让指数样本/权重/估值、成分股财务、资金或可替代情绪指标进入同一套只读、低频、可降级的数据与审计路径。

## What Changes

- 扩展真实公开数据 collector 范围：继续校准中证指数样本、权重、估值文件等公开数据，并评估成分股财务、资金流向、融资融券或可替代情绪指标。
- 为新增数据统一定义 source metadata、source level、freshness、missing、stale、no_data、source_unavailable、parse_error 等状态。
- 将新增数据写入现有市场/证据/审计模型，或在必要时新增轻量表结构；不得伪造缺失数据。
- 将可用的新数据接入每日纪律、expected return 输入上下文，并为 P35 风险预警预留清晰数据边界。
- 增加数据源健康状态、最近成功/失败记录和可验证 smoke，保持默认可用 stub/fixture。
- 保持安全边界：不接券商交易 API、不读取登录/付费/授权/Level2/高频源、不自动交易、不外部推送、不承诺收益、不预测确定涨跌。

## Capabilities

### New Capabilities
- `real-data-coverage-expansion`: 覆盖 P34 新增真实公开数据的采集、标准化、刷新健康、失败分类、工作流输入和安全边界。

### Modified Capabilities
- `data-source`: 扩展 P26/P27 公开证据与市场数据 collector 的来源范围、payload 字段、失败分类、低频刷新和审计要求。
- `real-data-integration`: 要求新增真实数据可进入每日纪律、expected return 和审计上下文，并保持降级与 stub 行为。
- `daily-discipline-report`: 扩展每日纪律报告读取新增数据后的缺失/过期/失败诊断展示要求。

## Impact

- 后端：collector、normalizer、market refresh / evidence refresh 服务、DTO、repository、migration、审计写入和 `cmd/agent` 本地任务。
- 数据模型：可能扩展 `market_snapshots.market_metrics_json`、source health metadata，必要时新增轻量持久化结构。
- 工作流：DailyDisciplineGraph、ExpectedReturnNode、未来 P35 风险预警输入上下文。
- 前端：数据源健康、每日纪律报告、设置/运维入口的状态展示。
- 测试：collector fixture、真实源 smoke、失败分类、stub 降级、OpenSpec strict 校验和 E2E smoke。
