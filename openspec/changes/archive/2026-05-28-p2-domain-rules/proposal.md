# Proposal: P2 领域规则

## Intent

实现 Investment Agent 的核心领域模型、枚举和规则裁决引擎，让系统能够基于能力圈、证据、多源验证、估值区间、止盈、现金冗余、核心-卫星仓位与规则提案状态机给出可测试的规则裁决。

## Scope

### In scope

- P2.1 核心模型与枚举：定义 API、数据模型、工作流中使用的核心枚举和领域结构。
- P2.2 规则裁决引擎：实现能力圈、证据、信源、重大事件验证、买入逻辑破坏、情绪、估值、移动止盈、现金冗余、核心-卫星仓位、预期收益和规则提案状态机。
- 为 P2 规则场景编写单元测试。

### Out of scope

- P3 Eino Graph、节点框架与实际工作流编排。
- P4 HTTP API handler 与 DTO 输出。
- P5 前端页面真实数据接入。
- LLM 调用、VecLite 检索、外部新闻采集。
- 自动交易或券商接口。

## Source documents

- `docs/development-plan.md`：P2：领域规则（P2.1、P2.2）
- `docs/api.md` 第 4 节：核心枚举
- `docs/data-model.md` 第 3、7 节：命名规范与状态模型
- `docs/workflow.md` 第 4、5 节：WorkflowContext 与节点边界
- `docs/requirements.md` 第 2、6 节：核心原则与纪律边界

## Expected outcome

- 领域模型可被后续 P3/P4 复用。
- 规则引擎不依赖数据库、HTTP 或 LLM。
- `go test ./internal/domain/model/... ./internal/domain/rule/...` 通过。

## Plan alignment

本 change 对应 `docs/development-plan.md` 的 P2 全部内容：P2.1、P2.2。
