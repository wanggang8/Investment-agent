# Proposal: P8 前端体验与测试

## Why

P7 已完成真实数据、RAG/VecLite 与 DeepSeek 分析材料底座，P8 需要在既有 API 与前端契约基础上提升驾驶舱可读性、交互可验证性和前端测试覆盖。该阶段只增强前端展示与测试，不改变产品边界。

## What Changes

- 在今日纪律、持仓、复盘页面增加图表组件，用于展示仓位、风险、证据覆盖和复盘摘要。
- 图表数据只来自 API DTO，不直接读取 SQLite、VecLite 或本地文件。
- 增强证据、决策链、审计时间线的筛选和展开交互。
- 为信息不足、数据过期、LLM 降级、VecLite 不可用等状态补充明确空态和错误态。
- 建立前端测试脚本，覆盖 API client、关键状态渲染、用户确认流程、规则提案最终确认流程和禁止自动交易入口断言。
- 保持用户确认区为线下动作记录，不增加一键交易或自动下单入口。
- 代码中非显然实现逻辑需要写中文注释。

## In Scope

- 对齐 `docs/development-plan.md` P8.1：驾驶舱图表与关键交互。
- 对齐 `docs/development-plan.md` P8.2：前端测试与契约校验。
- 修改范围限于前端页面、组件、API client、前端测试、必要的前端类型适配。
- 如需规格 delta，仅覆盖前端体验、状态展示、测试验收与无自动交易入口约束。

## Out of Scope

- 不修改 P7 数据源、RAG/VecLite、DeepSeek 后端集成能力。
- 不实现 P9 `cmd/agent`、月度/季度复盘自动化或本地交付说明。
- 不新增后端 API 契约之外的数据读取路径。
- 不直接读取 SQLite、VecLite、本地文件或真实密钥。
- 不新增自动交易、一键交易、代下单或收益承诺能力。

## Plan Alignment

- P8.1 驾驶舱图表与关键交互：一一对应。
- P8.1 验收命令 `cd web && npm run build`：一一对应。
- P8.2 前端测试与契约校验：一一对应。
- P8.2 验收命令 `cd web && npm run build && npm test`：一一对应。
- 未加入 `docs/development-plan.md` P8 以外的新需求。

## Capabilities

### New Capabilities
- `frontend-experience-tests`: 覆盖 P8 前端图表、关键交互、状态展示、前端测试和禁止自动交易入口断言。

### Modified Capabilities
- `e2e-hardening`: 增加 P8 前端构建、测试和产品边界验收要求。

## Impact

- 主要影响 `web/src/pages/`、`web/src/components/`、`web/src/services/`、`web/src/types/`、`web/src/styles/` 与前端测试配置。
- 可能调整前端 mock、fixture 或测试脚本，但不改变后端 API 真源契约。
- OpenSpec delta 目标：新增 `frontend-experience-tests`，修改 `e2e-hardening`。
