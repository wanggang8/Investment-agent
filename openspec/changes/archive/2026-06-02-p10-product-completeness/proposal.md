## Why

P0-P9 已完成本地验收，但对照总需求、架构、工作流和前端契约后，仍存在产品级完整度缺口。当前需要把这些缺口纳入一个独立阶段，避免把本地 MVP 状态误认为完整产品交付。

## What Changes

- 新增 P10 产品级补全阶段，聚焦真实外部数据源、真实 VecLite、完整 RAG/新闻管道、节点级 Eino 编排、守门人深度审计和前端操作入口。
- 修正文档状态一致性：将 P9 已归档完成的条目与开发计划总清单对齐。
- 明确 P10 只补全现有需求文档中已提出的能力，不新增自动交易、主动荐股或收益承诺能力。
- 保持 DeepSeek 只生成分析材料，最终裁决仍由领域规则完成。

## Capabilities

### New Capabilities
- `product-completeness`: 覆盖 P10 产品级补全阶段，包括真实数据、真实检索、完整工作流编排、前端操作入口和文档状态一致性。

### Modified Capabilities
<!-- No existing requirement is modified in this proposal. P10 adds a cross-cutting `product-completeness` capability that references existing real-data, frontend, and review capabilities without changing their archived requirement text. -->

## In Scope

- 真实行情与情报数据源适配，不写真实密钥。
- VecLite 文件索引、重建与降级路径。
- 新闻/公告/RAG 管道节点拆分与审计。
- Eino 节点级 Graph 编排。
- 守门人审计的规则检查、冲突检查、回测样本和审计理由。
- 前端刷新市场、重建索引、更新设置、账户录入、规则库展示、证据验证面板、单条决策审计时间线。
- `docs/development-plan.md` 中 P9 和总清单状态与归档结果对齐。

## Out of Scope

- 自动下单或券商交易 API。
- 系统主动推荐具体标的。
- 收益承诺或确定性涨跌预测。
- 修改 L1 契约正文；本阶段只在 change delta 中描述变更，archive 时再合并。
- 发明 `docs/development-plan.md` 与总需求以外的新产品能力。

## Impact

- 后端：`internal/application/workflow`、`internal/infrastructure`、`internal/application/service`、`cmd/agent`。
- 前端：`web/src/pages`、`web/src/components`、`web/src/features`、`web/src/services`、`web/src/types`。
- 文档：`docs/development-plan.md`、`docs/configuration.md`、`docs/testing-plan.md`、OpenSpec delta。
- 验证：Go 全量测试、前端测试与构建、OpenSpec 严格校验、本地任务命令。
