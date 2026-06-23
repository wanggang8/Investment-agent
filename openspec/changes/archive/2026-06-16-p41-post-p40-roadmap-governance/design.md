# P41 Post-P40 Roadmap Governance Design

## Context

P40 已归档后，`openspec/PROGRESS.md` 进入无活跃 change 状态。既有 P33-P40 队列已经全部完成，后续任何产品愿景、功能增强或审计追溯都需要先通过新的 OpenSpec change 进入治理流程。

## Decision

本 change 只做路线图治理，不做运行时实现。它会把 P40 后工作分为四类，并要求每一类在进入实现前单独创建 OpenSpec change：

| 类别 | 建议 change | 依赖 | 验收重点 | 适合度 |
| --- | --- | --- | --- | --- |
| 产品能力增强 | `p42-user-decision-workbench` | P32、P33、P35、P36、P39 | 每日工作台、主动咨询、组合复盘和规则治理入口；浏览器 E2E 覆盖空库、完整数据、降级和窄屏路径 | 高 |
| 数据质量增强 | `p43-data-quality-observability` | P34、P37、P38、P40 | source health、证据新鲜度、RAG 命中率、LLM 质量门禁和脱敏审计 | 中高 |
| 运维体验增强 | `p44-local-install-diagnostics-packaging` | P31、P39、P40 | 本地安装、配置向导、诊断导出、备份恢复演练和 smoke 汇总，不污染真实私有数据库 | 中 |
| 历史审计追溯 | `p19-p24-historical-archive-traceability` | P19-P24 已交付但 archive 不完整 | 整理历史证据、验收记录和文档一致性；不伪造历史 change | 条件触发 |

默认推荐顺序是先做产品能力增强，再按数据质量风险选择数据质量增强，随后推进运维体验增强；历史审计追溯按审计需求独立排期。

## Boundaries

- 不接券商 API，不自动交易，不外部推送，不自动应用规则。
- 不新增登录源、付费源、授权源、Level2 或高频源。
- 不承诺收益，不预测确定涨跌。
- LLM 仍只生成分析材料，最终裁决仍由规则和守门人链路处理。
- P19-P24 历史 archive 追溯只能作为独立治理 change，不伪造历史归档。

## Verification

本 change 是文档和 OpenSpec 治理变更，验收以 `openspec validate p41-post-p40-roadmap-governance --strict`、`openspec validate --all --strict` 和文档一致性检查为主；不需要运行 Go、前端或 E2E 测试，除非后续修改运行时代码。
