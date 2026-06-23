# P38: RAG / VecLite 检索质量加固

## Why

P13 已建立本地 JSON index health、rebuild 和 SQLite fallback，P15/P26/P27/P34 已补齐证据质量、公开源和数据健康基础。但当前检索仍偏“能检索、能降级”：缺少可重复的质量测试集、召回/重排质量指标、引用一致性检查和索引新鲜度展示。P38 需要把检索链路提升为“可评估、可解释、可复现、可降级”的本地能力。

## What

- 建立本地检索质量测试集，覆盖公告、监管文件、基金/ETF 信息、市场背景材料和 C 级背景材料。
- 增加检索质量评估服务或等价本地任务，输出 expected evidence、actual evidence、hit/miss、source level、evidence role、freshness 和 degradation reason。
- 增强 retrieval ranking：在现有 VecLite/JSON index 与 SQLite fallback 基础上，加入质量元数据、时间权重、信源等级、formal/background 边界和 source verification 状态参与排序或过滤。
- 增加证据引用一致性检查，确保返回 evidence 与 `source_verifications`、信源等级、时效权重、RAG chunk metadata 对齐。
- 增加索引新鲜度、损坏、重建、版本兼容和检索降级的审计/API/前端展示。

## Out of Scope

- 不接券商 API、不自动交易、不外部推送。
- 不绕过 source verification、规则裁决、守门人审计或用户最终确认。
- 不把 C 级或未验证背景材料升级为 formal 裁决证据。
- 不引入付费、登录、授权、Level2 或高频数据源。
- 不要求真实外部 VecLite 服务；当前仍以可替换本地 index adapter 和 SQLite 事实为准。

## Impact

- 后端：retrieval service、index adapter、workflow evidence retrieval、audit、query/service DTO 和本地任务。
- 前端：ops/status 或 decision evidence 展示中的 retrieval quality、index freshness、fallback reason。
- 文档：`docs/workflow.md`、`docs/api.md`、`docs/data-model.md`、`docs/frontend-contract.md`、`docs/configuration.md` 的检索质量与降级契约。
- 验收：Go tests、frontend tests/build、OpenSpec validation，以及本地检索质量 smoke。
