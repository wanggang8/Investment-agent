## ADDED Requirements

### Requirement: Market and intelligence data sources

系统 SHALL 提供真实行情与情报数据源适配能力，并支持按配置启用真实数据源或本地 stub。

#### Scenario: Market data refresh succeeds

- **WHEN** 行情数据源返回全部目标标的的有效数据
- **THEN** 系统 MUST 写入 `market_snapshots`
- **AND** 系统 MUST 写入成功状态的 `audit_events`

#### Scenario: Market data refresh degrades

- **WHEN** 行情数据源部分失败、全部失败或返回过期数据
- **THEN** 系统 MUST 返回既有错误或降级状态
- **AND** 系统 MUST 记录失败标的、降级原因或错误码
- **AND** 系统 MUST 写入可追踪的审计事件

#### Scenario: Intelligence data is ingested

- **WHEN** 新闻、公告或手工导入情报被采集
- **THEN** 系统 MUST 写入 `intelligence_items`
- **AND** 系统 MUST 保留来源、时间和信源等级信息

#### Scenario: Local stub is available

- **WHEN** 真实数据源未配置或被关闭
- **THEN** 系统 MUST 可使用本地 stub 完成开发与测试
- **AND** 系统 MUST NOT 写入真实密钥或环境私有值到文档或日志

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

#### Scenario: DeepSeek is unavailable

- **WHEN** DeepSeek 缺配置、超时、不可用或输出不可解析
- **THEN** 工作流 MUST 进入降级状态
- **AND** 规则引擎 MUST 继续生成最终裁决
- **AND** 系统 MUST 写入可追踪的降级原因

#### Scenario: No automatic trading from analysis

- **WHEN** DeepSeek 或检索服务生成分析材料
- **THEN** 系统 MUST NOT 生成自动交易动作
- **AND** 系统 MUST NOT 提供一键交易入口
