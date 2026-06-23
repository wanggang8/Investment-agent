## MODIFIED Requirements

### Requirement: RAG and VecLite retrieval
系统 SHALL 提供本地 JSON 文件索引读写、SQLite 重建、健康状态、重建统计和不可用降级能力；真实 VecLite API 替换边界 SHALL 保持可替换但不在 P13 强制接入。

#### Scenario: VecLite index is built
- **WHEN** `rag_chunks` 与 `intelligence_summary` 存在可检索文本
- **THEN** 系统 MUST 将其纳入本地 JSON 文件索引构建流程
- **AND** 索引路径 MUST 来自配置
- **AND** 系统 MUST 记录健康状态和重建统计

#### Scenario: VecLite index is rebuilt
- **WHEN** 本地索引缺失、损坏、不兼容或需要重建
- **THEN** 系统 MUST 支持从 SQLite 文本块重建索引
- **AND** 重建过程 MUST 可测试
- **AND** 重建结果 MUST 暴露 indexed/skipped 数量与降级原因

#### Scenario: VecLite is unavailable
- **WHEN** 本地索引不可用、损坏、不兼容或检索失败
- **THEN** 系统 MUST 按既有约定降级到 SQLite 摘要
- **AND** 摘要不足时 MUST 返回信息不足
- **AND** 系统 MUST 记录检索输入、命中证据和降级原因

#### Scenario: Retrieval service records fallback context
- **WHEN** 检索从本地索引降级到 SQLite 摘要
- **THEN** 工作流 MUST 保留输入标的、命中摘要或降级原因
- **AND** 审计事件 MUST 可追踪本次检索状态

#### Scenario: C-level source is restricted
- **WHEN** 检索命中 C 级信源
- **THEN** C 级信源 MUST 只能作为 background 材料
- **AND** C 级信源 MUST NOT 作为正式裁决依据
