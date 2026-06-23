## 1. 文档状态一致性

- [x] 1.1 对照 `openspec/changes/archive/2026-06-02-p9-review-automation-delivery/tasks.md`，更新 `docs/development-plan.md` 中 P9.1 全部已完成任务勾选状态。
- [x] 1.2 对照 P9 归档验收记录，更新 `docs/development-plan.md` 中 P9.2 全部已完成任务勾选状态。
- [x] 1.3 对照 P9 归档交付记录，更新 `docs/development-plan.md` 中 P9.3 全部已完成任务勾选状态。
- [x] 1.4 更新 `docs/development-plan.md` 总清单中已完成的 `cmd/agent`、被动标的边界、月度/季度复盘边界勾选状态。
- [x] 1.5 确认文档状态修正不新增需求、不改变禁止自动交易边界。

## 2. 真实数据源与 VecLite 补全

- [x] 2.1 实现配置化行情数据源入口；未接入供应商时返回可识别不可用错误，保留本地 stub 和失败降级。
- [x] 2.2 实现配置化情报源入口与本地情报写入路径，写入 `intelligence_items` 并保留来源、时间、信源等级。
- [x] 2.3 实现可替换 VecLite 的 JSON 文件索引适配层，索引路径来自配置，SQLite 仍为事实基准。
- [x] 2.4 支持从 `rag_chunks` 与 `intelligence_summary` 重建 VecLite 索引。
- [x] 2.5 对索引缺失、损坏、检索失败添加降级路径和审计记录。
- [x] 2.6 为真实数据、索引重建和降级路径添加测试。

## 3. 工作流与守门人补全

- [x] 3.1 将 DailyDisciplineGraph 从单 Lambda 包装改为节点级 Eino Graph 编排。
- [x] 3.2 将 ConsultationGraph 从单 Lambda 包装改为节点级 Eino Graph 编排。
- [x] 3.3 将 EvidenceVerificationGraph 拆分为新闻获取、分类、标准化、嵌入、索引写入、多源验证等可审计步骤。
- [x] 3.4 补强 GatekeeperAuditGraph 的根本规则检查、规则冲突检查、回测样本和审计理由。
- [x] 3.5 确认 DeepSeek 仍只写分析材料，最终裁决仍由领域规则生成。
- [x] 3.6 为节点级工作流、审计事件和降级路径添加测试。

## 4. 前端产品级入口

- [x] 4.1 在设置或数据页面增加市场刷新入口，展示执行状态和无交易说明。
- [x] 4.2 在证据页增加情报刷新和索引重建入口，展示执行状态、错误码和降级原因。
- [x] 4.3 增加账户/持仓录入或校准入口，只记录本地事实，不连接交易接口。
- [x] 4.4 在规则页展示当前规则库、裁决优先级、阈值和规则提案。
- [x] 4.5 在证据页增加多源验证独立面板，展示独立信源数量、最高等级、最新发布时间和证据引用。
- [x] 4.6 在决策详情页增加单条决策审计时间线或内嵌审计区域。
- [x] 4.7 统一金额、比例、时间、分位和状态文案格式化。
- [x] 4.8 为新增入口、错误态、空态和禁止自动交易断言添加前端测试。

## 5. 验收

- [x] 5.1 执行 `go test ./...` 并记录结果。
- [x] 5.2 执行 `cd web && npm run test && npm run build` 并记录结果。
- [x] 5.3 执行 OpenSpec 严格校验并记录结果。
- [x] 5.4 执行 `go run ./cmd/agent --help` 并记录结果。
- [x] 5.5 执行 `go run ./cmd/agent --config configs/config.example.yaml --task market-refresh` 并记录结果。
- [x] 5.6 执行 `go run ./cmd/agent --config configs/config.example.yaml --task evidence-index` 并记录结果。
- [x] 5.7 执行 `go run ./cmd/agent --config configs/config.example.yaml --task review --period monthly` 并记录结果。
- [x] 5.8 确认 P10 没有自动交易入口、主动荐股能力或收益承诺。
