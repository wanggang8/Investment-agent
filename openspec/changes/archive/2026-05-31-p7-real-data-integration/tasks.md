# Tasks: P7 真实数据与分析底座

## 1. P7.1 真实行情与情报数据源

- [x] 1.1 增加行情数据源适配层，支持按配置启用真实数据源或本地 stub。
- [x] 1.2 增加情报数据源适配层，支持新闻、公告或手工导入数据进入 `intelligence_items`。
- [x] 1.3 为市场快照刷新写入 `market_snapshots` 与 `audit_events`。
- [x] 1.4 对部分失败、全部失败、数据过期分别返回既有错误或降级状态。
- [x] 1.5 在配置文档中说明数据源开关、凭证环境变量和本地 stub 用法，不写真实密钥。
- [x] 1.6 对数据源降级、审计写入和禁止自动交易边界添加必要中文注释。
- [x] 1.7 执行验收命令：`go test ./internal/infrastructure/... ./internal/application/...`。

## 2. P7.2 RAG/VecLite 检索与索引

- [x] 2.1 实现 VecLite 索引读写适配，索引路径来自配置。
- [x] 2.2 将 `rag_chunks` 与 `intelligence_summary` 纳入检索构建流程。
- [x] 2.3 支持从 SQLite 文本块重建 VecLite 索引。
- [x] 2.4 VecLite 不可用时按既有约定降级到 SQLite 摘要或信息不足。
- [x] 2.5 记录检索输入、命中证据和降级原因到审计事件或可追踪上下文。
- [x] 2.6 确认 C 级信源仍不得作为正式裁决依据。
- [x] 2.7 对检索重建、降级和信源边界添加必要中文注释。
- [x] 2.8 执行验收命令：`go test ./internal/infrastructure/... ./internal/application/workflow/...`。

## 3. P7.3 DeepSeek 分析师材料

- [x] 3.1 增加 DeepSeek 客户端封装，API Key 从环境变量读取。
- [x] 3.2 将价值分析、趋势风险和预期收益节点从占位实现改为可调用分析服务。
- [x] 3.3 明确 prompt 输入只包含允许使用的证据、持仓上下文和规则边界。
- [x] 3.4 解析 DeepSeek 输出为 `analyst_reports` 或等价结构，不写最终裁决。
- [x] 3.5 LLM 不可用、超时或输出不可解析时，工作流进入降级状态，并由规则引擎继续生成最终裁决。
- [x] 3.6 对非显然的 prompt 约束、降级和审计逻辑添加中文注释。
- [x] 3.7 确认 DeepSeek 仅提供分析材料，最终裁决仍由规则引擎负责，LLM 故障不产生自动交易动作。
- [x] 3.8 执行验收命令：`go test ./internal/application/workflow/... ./internal/infrastructure/...`。

## 4. 阶段一致性检查

- [x] 4.1 确认本 change 没有实现 P8 前端图表、交互增强和前端测试。
- [x] 4.2 确认本 change 没有实现 P9 `cmd/agent`、周期复盘和本地交付说明。
- [x] 4.3 确认没有加入 `docs/development-plan.md` P7 以外的新需求。
- [x] 4.4 确认没有真实密钥、账号、token 或个人敏感信息进入代码和文档。
- [x] 4.5 执行 `openspec status --change p7-real-data-integration`，确认 artifacts 完整。
