# Proposal: P7 真实数据与分析底座

## Summary

P7 将 P0-P6 已完成的本地骨架接入可替换的真实行情、情报、RAG/VecLite 和 DeepSeek 分析材料能力，同时保留外部依赖缺失时的降级路径与审计记录。

## Why

P6 已完成端到端验收与配置文档，但当前系统仍主要依赖本地数据和占位分析。按照 `docs/development-plan.md` P7，需要补齐真实数据与分析底座，让后续 P8 前端体验和 P9 周期复盘有可靠输入。

## What Changes

- 增加行情数据源适配层，支持按配置启用真实数据源或本地 stub。
- 增加情报数据源适配层，支持新闻、公告或手工导入数据进入 `intelligence_items`。
- 为市场快照刷新写入 `market_snapshots` 与 `audit_events`，覆盖部分失败、全部失败和数据过期场景。
- 实现 VecLite 索引读写适配，索引路径来自配置。
- 将 `rag_chunks` 与 `intelligence_summary` 纳入检索构建流程，并支持从 SQLite 重建索引。
- VecLite 不可用时降级到 SQLite 摘要或信息不足。
- 增加 DeepSeek 客户端封装，API Key 从环境变量读取。
- 将价值分析、趋势风险和预期收益节点从占位实现改为可调用分析服务。
- 解析 DeepSeek 输出为 `analyst_reports` 或等价结构，不写最终裁决。
- LLM 不可用、超时或输出不可解析时，工作流进入降级状态，并由规则引擎继续生成最终裁决。
- 对非显然的 prompt 约束、降级、审计和禁止自动交易边界写中文注释。

## In Scope

- 覆盖 `docs/development-plan.md` P7.1 的全部任务和验收命令：`go test ./internal/infrastructure/... ./internal/application/...`。
- 覆盖 `docs/development-plan.md` P7.2 的全部任务和验收命令：`go test ./internal/infrastructure/... ./internal/application/workflow/...`。
- 覆盖 `docs/development-plan.md` P7.3 的全部任务和验收命令：`go test ./internal/application/workflow/... ./internal/infrastructure/...`。
- 保留 P7 中的降级要求：外部依赖缺失时展示信息不足或定性说明，不伪造证据。
- 保留信源边界：C 级信源不得作为正式裁决依据。
- 保留 LLM 边界：DeepSeek 只生成分析材料，最终裁决由规则引擎负责。
- 保留产品边界：不新增自动交易或一键交易入口。

## Out of Scope

- 不实现 P8 的前端图表、交互增强和前端测试。
- 不实现 P9 的 `cmd/agent` 本地任务、月度/季度复盘和本地交付说明。
- 不新增 `docs/development-plan.md` P7 之外的供应商能力、自动化交易能力或收益承诺能力。
- 不修改 L1 契约正文；如发现契约不一致，先在本 change 的 specs delta 中说明并经归档合并。
- 不把 DeepSeek 输出作为最终裁决，不让 LLM 越过规则引擎。

## Capabilities

### New Capabilities

- `real-data-integration`: 真实行情、情报数据源、RAG/VecLite 检索与 DeepSeek 分析师材料接入。

### Modified Capabilities

- `e2e-hardening`: 承接 P6 已定义的降级、安全边界与验收约束，补充 P7 真实数据与分析路径下的可验证要求。

## Impact

- 可能影响后端入口与配置：`configs/config.example.yaml`、`internal/config` 或等价配置加载模块。
- 可能影响基础设施层：行情数据源、情报数据源、VecLite/RAG、DeepSeek 客户端。
- 可能影响应用层工作流：市场刷新、证据检索、价值分析、趋势风险、预期收益节点。
- 可能影响持久化与审计写入：`market_snapshots`、`intelligence_items`、`rag_chunks`、`intelligence_summary`、`audit_events`。
- 可能影响测试：基础设施测试、应用层测试、工作流降级测试。

## Plan Alignment

本 proposal 与 `docs/development-plan.md` P7 一一对应：

- P7.1 真实行情与情报数据源：对应数据源适配、market/intelligence 写入、失败降级、配置说明与验收命令。
- P7.2 RAG/VecLite 检索与索引：对应 VecLite 读写、SQLite 重建、不可用降级、检索审计与验收命令。
- P7.3 DeepSeek 分析师材料：对应客户端封装、分析节点替换、prompt 输入边界、输出解析、LLM 降级、中文注释与验收命令。

未加入 P7 计划以外的需求。实现代码中的非显然业务约束、降级路径、审计写入和禁止自动交易边界需要写中文注释。
