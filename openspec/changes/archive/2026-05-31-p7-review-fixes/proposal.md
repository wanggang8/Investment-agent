# Proposal: P7 审查问题修复

## Summary

修复 P7 子 agent 审查指出的实现缺口：补齐预期收益分析服务调用、RAG/VecLite 检索降级与审计、运行时配置到工作流依赖的串联。

## Why

`p7-real-data-integration` 已归档，但审查发现部分实现偏占位或未接入生产依赖组装。该 change 只修复 P7 既有范围内的问题，不新增 P8/P9 能力。

## What Changes

- 让预期收益节点也调用 `AnalystService` 生成分析材料，并在不可用时进入降级状态。
- 增加检索服务边界：优先使用 VecLite/RAG 检索，失败时降级到 SQLite 摘要，摘要不足时返回信息不足。
- 将检索输入、命中证据和降级原因写入审计事件或工作流可追踪上下文。
- 将 `DeepSeek`、`data_sources`、`VecLite` 配置接入生产依赖组装，避免运行时总是使用 static/stub。
- 修正实际独立信源数量统计，避免用请求参数数量替代数据源返回数量。
- 补充必要测试和中文注释。

## In Scope

- 仅覆盖 P7 审查问题。
- 修改 P7 相关工作流、服务、基础设施配置与测试。
- 保留已归档 P7 的安全边界：不自动交易、LLM 不写最终裁决、C 级信源不作正式裁决依据。

## Out of Scope

- 不实现 P8 前端图表、交互增强或前端测试。
- 不实现 P9 `cmd/agent`、周期复盘或本地交付说明。
- 不新增真实供应商复杂能力或外部请求频率策略。
- 不写入真实密钥、账号、token 或个人敏感信息。

## Capabilities

### Modified Capabilities

- `real-data-integration`: 补齐 P7 数据、检索、DeepSeek 分析材料的运行时路径与降级审计。
- `e2e-hardening`: 保持 P7 验收命令与安全边界可验证。

## Impact

- 可能影响 `internal/application/workflow`、`internal/application/service`、`internal/infrastructure/config`、`internal/infrastructure/llm/deepseek`、后端入口依赖组装。
- 可能新增或调整测试：工作流、服务、配置、基础设施测试。
