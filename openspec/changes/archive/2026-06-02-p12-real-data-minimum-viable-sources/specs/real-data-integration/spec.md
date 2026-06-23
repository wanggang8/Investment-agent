## MODIFIED Requirements

### Requirement: Market and intelligence data sources
系统 SHALL 提供最小可用真实只读行情与情报数据源适配能力，并支持按配置启用真实数据源或本地 stub。P12 范围 SHALL 明确排除完整财务源、完整情绪源、实时性 SLA、券商交易 API、自动交易、主动荐股和收益承诺。

#### Scenario: Market data refresh succeeds
- **WHEN** configured readonly market data source returns valid data for target symbols
- **THEN** 系统 MUST 写入 `market_snapshots`
- **AND** 系统 MUST 写入成功状态的 `audit_events`

#### Scenario: Market data refresh degrades
- **WHEN** 行情数据源部分失败、全部失败、超时、解析失败或返回过期数据
- **THEN** 系统 MUST 返回既有错误或降级状态
- **AND** 系统 MUST 记录失败标的、降级原因或错误码
- **AND** 系统 MUST 写入可追踪的审计事件

#### Scenario: Intelligence data is ingested
- **WHEN** configured readonly intelligence source returns valid news, announcement, or manually imported intelligence
- **THEN** 系统 MUST 写入 `intelligence_items`
- **AND** 系统 MUST 保留来源、时间和信源等级信息

#### Scenario: Local stub is available
- **WHEN** 真实数据源未配置或被关闭
- **THEN** 系统 MUST 可使用本地 stub 完成开发与测试
- **AND** 系统 MUST NOT 写入真实密钥或环境私有值到文档或日志
