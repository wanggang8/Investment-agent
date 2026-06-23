## 1. 后端阻塞修复

- [x] 1.1 补 consult-to-confirm 集成测试，覆盖正式建议可确认、非可操作裁决不可误确认
- [x] 1.2 修复 decision record 的 record_type、confirmation_status、available actions 与详情接口一致性
- [x] 1.3 补能力圈排除场景测试，并让主动咨询读取 capability_configs
- [x] 1.4 补规则提案生成后确认测试，修复 draft/pending_user_confirm 状态链路
- [x] 1.5 补守门人审计链路测试，修复用户确认后审计推进与最终确认前置校验

## 2. 后端重要数据一致性修复

- [x] 2.1 补手动执行买入、减仓、清仓现金与资产快照测试
- [x] 2.2 修复 manual execution 后 cash、cash_ratio、total_assets、positions、position_snapshots 一致性
- [x] 2.3 补确认事务中途失败回滚测试
- [x] 2.4 补决策详情回放测试，返回 analyst_reports、expected_return_scenarios、arbitration_chain、account_snapshot
- [x] 2.5 补市场快照 DTO 测试，返回 trade_date、data_status、market_metrics
- [x] 2.6 补证据列表真实信源字段测试，修复 source_name、original_url、published_at、captured_at、content_hash 映射
- [x] 2.7 补规则输入合并 SourceVerificationStatus 测试，并修复误判风险

## 3. 前端展示与交互修复

- [x] 3.1 补 Dashboard 状态映射和正常态测试，修复原始枚举与误导文案
- [x] 3.2 补 EvidenceSummary、PortfolioTable、market/data 状态中文映射与未知值测试
- [x] 3.3 补 DecisionTrace 免责声明、confidence/scenario 展示测试并修复
- [x] 3.4 补 RulesPage APIClientError 测试并接入统一错误展示
- [x] 3.5 补 DecisionDetail API 错误与确认失败测试，失败时不显示成功确认
- [x] 3.6 补 Evidence/Audit/Portfolio/ReviewSummary 页面成功空态测试
- [x] 3.7 补 EvidenceTable 展开字段测试，展示 hash、time_weight、relevance_score

## 4. 治理与质量门禁修复

- [x] 4.1 同步 docs/development-plan.md 中 P6/P7/P8 已完成项和总清单状态
- [x] 4.2 更新 docs/testing-plan.md 或等效验收文档，纳入 go test、前端 build、前端 test
- [x] 4.3 补齐 OpenSpec specs 轻量摘要，避免 P0–P5 capability 查询缺失且不复制 L1 全文
- [x] 4.4 检查测试命名、fetch mock 清理与质量细节

## 5. 验证与全仓复审

- [x] 5.1 执行 go test ./...
- [x] 5.2 执行 cd web && npm run build && npm test
- [x] 5.3 启动后端、前端、治理、数据工作流、测试质量 5 个同范围全仓子 agent 复审
- [x] 5.4 根据复审结果修复阻塞和重要问题，直到全仓复审通过或仅剩明确可接受轻微项

## 6. 复审发现修复

- [x] 6.1 修复守门人审计使用固定样本数推进的问题，并补真实 sample_count 不足测试
- [x] 6.2 将流动性状态纳入规则裁决，禁止危险流动性下新增买入类确认建议
- [x] 6.3 修复市场刷新到查询链路的核心行情指标持久化与 DTO 映射
- [x] 6.4 修复手动买入资金不足校验与费用字段一致性
- [x] 6.5 修复证据 content_hash/chunk_hash 由占位值生成的问题，并补稳定哈希测试
- [x] 6.6 修复重大事件 A/S 独立信源数量与规则输入映射
- [x] 6.7 补能力圈 excluded_symbols 集成测试
- [x] 6.8 修复前端 RULE_VERSION_MISSING 状态映射、Dashboard 确认入口、Portfolio 空态与买入理由展示
- [x] 6.9 清理 fetch mock、更新验收记录、同步开发计划日期并记录本轮复审结果
- [x] 6.10 执行 go test ./... 与 cd web && npm run build && npm test
- [x] 6.11 再次启动同范围复审，处理阻塞和重要问题

## 7. 最新复审重要问题修复

- [x] 7.1 修复 Dashboard 确认入口、主动咨询入口与前端枚举表单校验
- [x] 7.2 修复审计时间线契约字段展示与前端状态文案偏差
- [x] 7.3 修复手动执行费用字段、交易流水与账户快照一致性
- [x] 7.4 修复市场核心结构列写入与查询验证
- [x] 7.5 修复低样本规则提案风险说明与重大事件证据链路验证
- [x] 7.6 修复决策详情账户快照回放、确认状态流转与 consult scenario 校验
- [x] 7.7 修复架构 import 边界与 wiring 归属说明
- [x] 7.8 补 A/S 持久化、RULE_VERSION_MISSING、hash 差异、流动性动作等回归测试
- [x] 7.9 更新复审记录并执行全量后端与前端验证

## 8. 第三轮复审重要问题修复

- [x] 8.1 统一 consult scenario 契约值并补前后端测试
- [x] 8.2 修正买入费用计入持仓成本价并补费用端到端断言
- [x] 8.3 修复市场刷新写入失败独立失败审计并补测试
- [x] 8.4 补规则提案页契约字段展示与测试
- [x] 8.5 补审计字段、root_cause_tag、RULE_VERSION_MISSING、market 结构列、hash 差异等测试证据
- [x] 8.6 补 OpenSpec delta 与复审记录，执行全量验证

## 9. 第五轮复审重要问题修复

- [x] 9.1 补 tasks 与 OpenSpec delta 覆盖第四、五轮前端和数据契约修复
- [x] 9.2 补审计列表 API DTO 字段、市场失败审计持久化、source count 仓储与合法 scenario 测试
- [x] 9.3 修复 EvidenceDTO / 前端 Evidence 类型与 verification 服务返回类型
- [x] 9.4 补前端 event_id 稳定展开、费用校验、规则提案内容与非默认 scenario 测试
- [x] 9.5 收敛 consult scenario 枚举与边界测试覆盖
- [x] 9.6 执行全量验证并更新 review-notes / testing-plan
