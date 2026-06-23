# Design: P7 审查问题修复

## Context

P7 归档后审查发现三类问题：预期收益节点未使用分析服务、检索路径缺少 VecLite 查询到 SQLite 摘要的降级、配置未串到生产依赖组装。本 change 在 P7 既有范围内修复，不扩大产品能力。

## Goals / Non-Goals

**Goals:**

- `expectedReturnStep` 通过 `AnalystService` 获取预期收益分析材料。
- 检索服务提供 VecLite 优先、SQLite 摘要降级、信息不足返回和审计上下文。
- 运行时依赖基于配置创建，DeepSeek API Key 存在时使用 DeepSeek client；缺失时使用降级服务或本地 stub。
- 审计记录包含检索输入、命中引用和降级原因。
- 实际独立来源按数据源返回项去重计算。

**Non-Goals:**

- 不接入具体真实行情供应商。
- 不实现自动交易或一键交易。
- 不实现 P8/P9 范围。

## Decisions

### Decision 1: 检索服务作为工作流依赖

新增 `RetrievalService` 接口。工作流证据节点优先使用该接口；接口返回命中证据、降级原因和状态。默认实现可从 SQLite 摘要构造 EvidenceSet，VecLite 不可用时不阻断规则流程。

### Decision 2: 预期收益分析与数值情景分离

`expectedReturnStep` 保留 `BuildExpectedReturn` 的数值情景，同时调用 `AnalystService` 写入 `analyst_reports[expected_return]`。LLM 失败时节点降级，但最终裁决仍交由规则引擎。

### Decision 3: 配置只决定适配器选择，不写供应商密钥

生产组装读取配置字段。`DEEPSEEK_API_KEY` 存在时创建 DeepSeek client；为空时使用降级服务。数据源配置默认 stub；关闭 stub 且无真实实现时返回可追踪降级。

## Risks / Trade-offs

- 当前没有真实 VecLite SDK：以接口与本地实现承载检索路径，保留替换点。
- LLM 降级可能让分析材料缺失：保留规则引擎裁决和审计原因。
- 配置接入可能影响启动：默认 stub 保持本地测试可用。

## Verification

- `go test ./internal/infrastructure/... ./internal/application/...`
- `go test ./internal/infrastructure/... ./internal/application/workflow/...`
- `go test ./internal/application/workflow/... ./internal/infrastructure/...`
