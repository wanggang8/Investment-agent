## MODIFIED Requirements

### Requirement: DeepSeek analyst materials
系统 SHALL 接入可配置 DeepSeek/OpenAI-compatible 分析服务，并确保 LLM 只生成分析材料，不生成最终裁决。

#### Scenario: DeepSeek analysis succeeds
- **WHEN** DeepSeek 或兼容端点返回可解析输出
- **THEN** 系统 MUST 将输出解析为 `analyst_reports` 或等价结构
- **AND** 系统 MUST 记录 prompt version、model、input summary、output summary、parse status 和 quality status
- **AND** 规则引擎 MUST 继续负责最终裁决

#### Scenario: DeepSeek input is bounded
- **WHEN** 系统构造 DeepSeek prompt
- **THEN** prompt 输入 MUST 只包含允许使用的证据、持仓上下文和规则边界
- **AND** 非显然 prompt 约束 MUST 有中文注释说明
- **AND** 审计与响应 MUST NOT 包含 API key、完整敏感 prompt、券商/账户密钥或不必要本地文件路径

#### Scenario: DeepSeek supports expected return material
- **WHEN** 预期收益节点执行
- **THEN** 系统 MUST 调用分析服务生成预期收益分析材料
- **AND** 数值情景可继续由本地样本逻辑生成
- **AND** DeepSeek MUST NOT 写最终裁决

#### Scenario: DeepSeek is unavailable
- **WHEN** DeepSeek 缺配置、超时、HTTP 不可用、空响应、输出不可解析或质量检查失败
- **THEN** 工作流 MUST 进入降级状态
- **AND** 规则引擎 MUST 继续生成最终裁决
- **AND** 系统 MUST 写入可追踪的降级原因和稳定错误分类

#### Scenario: No automatic trading from analysis
- **WHEN** DeepSeek 或检索服务生成分析材料
- **THEN** 系统 MUST NOT 生成自动交易动作
- **AND** 系统 MUST NOT 提供一键交易入口

### Requirement: Runtime dependency wiring
系统 SHALL 根据本地配置组装真实数据、检索与分析服务依赖。

#### Scenario: DeepSeek key is configured
- **WHEN** `deepseek.api_key`、`deepseek.base_url` 和 `deepseek.model` 或等价配置存在
- **THEN** 生产依赖 MUST 使用配置的 DeepSeek/OpenAI-compatible client 作为分析服务

#### Scenario: DeepSeek key is missing
- **WHEN** DeepSeek key 缺失
- **THEN** 生产依赖 MUST 使用可追踪降级服务或本地 stub
- **AND** 不得伪造真实 DeepSeek 响应

#### Scenario: Data source stub setting is applied
- **WHEN** `data_sources.use_stub` 为 true
- **THEN** 生产依赖 MUST 使用本地 stub 数据源
- **AND** 本地验收 MUST 不依赖公网
