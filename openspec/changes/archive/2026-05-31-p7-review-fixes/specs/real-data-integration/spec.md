## MODIFIED Requirements

### Requirement: RAG and VecLite retrieval
系统 SHALL 提供 VecLite 索引读写、SQLite 重建和不可用降级能力。

#### Scenario: VecLite index is built
- **WHEN** `rag_chunks` 与 `intelligence_summary` 存在可检索文本
- **THEN** 系统 MUST 将其纳入检索构建流程
- **AND** VecLite 索引路径 MUST 来自配置

#### Scenario: VecLite index is rebuilt
- **WHEN** VecLite 索引缺失或需要重建
- **THEN** 系统 MUST 支持从 SQLite 文本块重建索引
- **AND** 重建过程 MUST 可测试

#### Scenario: VecLite is unavailable
- **WHEN** VecLite 不可用或检索失败
- **THEN** 系统 MUST 按既有约定降级到 SQLite 摘要
- **AND** 摘要不足时 MUST 返回信息不足
- **AND** 系统 MUST 记录检索输入、命中证据和降级原因

#### Scenario: Retrieval service records fallback context
- **WHEN** 检索从 VecLite 降级到 SQLite 摘要
- **THEN** 工作流 MUST 保留输入标的、命中摘要或降级原因
- **AND** 审计事件 MUST 可追踪本次检索状态

#### Scenario: C-level source is restricted
- **WHEN** 检索命中 C 级信源
- **THEN** C 级信源 MUST 只能作为 background 材料
- **AND** C 级信源 MUST NOT 作为正式裁决依据

### Requirement: DeepSeek analyst materials
系统 SHALL 接入 DeepSeek 分析服务，并确保 DeepSeek 只生成分析材料，不生成最终裁决。

#### Scenario: DeepSeek analysis succeeds
- **WHEN** DeepSeek 返回可解析输出
- **THEN** 系统 MUST 将输出解析为 `analyst_reports` 或等价结构
- **AND** 规则引擎 MUST 继续负责最终裁决

#### Scenario: DeepSeek input is bounded
- **WHEN** 系统构造 DeepSeek prompt
- **THEN** prompt 输入 MUST 只包含允许使用的证据、持仓上下文和规则边界
- **AND** 非显然 prompt 约束 MUST 有中文注释说明

#### Scenario: DeepSeek supports expected return material
- **WHEN** 预期收益节点执行
- **THEN** 系统 MUST 调用分析服务生成预期收益分析材料
- **AND** 数值情景可继续由本地样本逻辑生成
- **AND** DeepSeek MUST NOT 写最终裁决

#### Scenario: DeepSeek is unavailable
- **WHEN** DeepSeek 缺配置、超时、不可用或输出不可解析
- **THEN** 工作流 MUST 进入降级状态
- **AND** 规则引擎 MUST 继续生成最终裁决
- **AND** 系统 MUST 写入可追踪的降级原因

#### Scenario: No automatic trading from analysis
- **WHEN** DeepSeek 或检索服务生成分析材料
- **THEN** 系统 MUST NOT 生成自动交易动作
- **AND** 系统 MUST NOT 提供一键交易入口

### Requirement: Runtime dependency wiring
系统 SHALL 根据本地配置组装真实数据、检索与分析服务依赖。

#### Scenario: DeepSeek key is configured
- **WHEN** `DEEPSEEK_API_KEY` 或等价配置存在
- **THEN** 生产依赖 MUST 使用 DeepSeek client 作为分析服务

#### Scenario: DeepSeek key is missing
- **WHEN** DeepSeek key 缺失
- **THEN** 生产依赖 MUST 使用可追踪降级服务或本地 stub
- **AND** 不得伪造真实 DeepSeek 响应

#### Scenario: Data source stub setting is applied
- **WHEN** `data_sources.use_stub` 为 true
- **THEN** 生产依赖 MUST 使用本地 stub 数据源
- **AND** 本地验收 MUST 不依赖公网
