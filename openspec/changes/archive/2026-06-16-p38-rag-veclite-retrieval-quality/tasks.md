## 1. OpenSpec 与范围

- [x] 1.1 确认 P38 只覆盖 RAG / VecLite 检索质量、引用一致性、索引新鲜度/重建和降级展示。
- [x] 1.2 确认 P38 不接券商 API、不自动交易、不外部推送、不绕过 source verification、不把背景材料升级为 formal 裁决证据。
- [x] 1.3 对齐 P13 index health、P15 evidence quality、P26/P27 public evidence/market collectors、P34 source health 的既有契约。

## 2. 检索质量测试集

- [x] 2.1 以 TDD 增加本地 retrieval quality fixture，覆盖公告、监管文件、基金/ETF 信息、市场背景材料和 C 级背景材料。
- [x] 2.2 fixture 包含 query、symbol、expected evidence ids 或 expected source/evidence-role constraints。
- [x] 2.3 增加质量评估结果结构：top-k、hit/miss、missing expected、unexpected background-only、fallback reason、index health/freshness。
- [x] 2.4 确认测试集不依赖公网、不含真实密钥、不伪造 formal source verification。

## 3. Quality-Aware Retrieval

- [x] 3.1 以 TDD 增强 retrieval ranking，纳入 relevance、source level、evidence role、time weight、freshness 和 source verification 状态。
- [x] 3.2 C 级、background 或 verification 不满足的 evidence 不得成为 formal 裁决依据。
- [x] 3.3 VecLite/JSON index 不可用时保持 SQLite fallback 或信息不足语义。
- [x] 3.4 新增 edge tests：空 index、损坏 index、版本不兼容、过期 index、SQLite 摘要不足。

## 4. 引用一致性与审计

- [x] 4.1 校验 retrieved evidence 与 `source_verifications`、`intelligence_summary`、`rag_chunks.metadata_json`、source level、evidence role 和 freshness 一致。
- [x] 4.2 不一致时写 degraded reason 或跳过结果，不静默返回误导性 evidence。
- [x] 4.3 audit_events 或等价追踪事实记录 query summary、index health、fallback source、quality status 和 degraded reason，且不含密钥或完整本地路径。
- [x] 4.4 工作流 EvidenceRetrievalNode 保留检索质量摘要，不改变最终规则裁决边界。

## 5. API / 前端展示

- [x] 5.1 更新 API/DTO，暴露 retrieval quality summary、index freshness、fallback reason 和 source consistency status。
- [x] 5.2 前端 ops/decision 视图展示检索降级原因、召回摘要和重建提示。
- [x] 5.3 前端不得读取 SQLite、VecLite 或本地文件，不展示自动交易、自动修复、自动规则应用入口。
- [x] 5.4 增加 frontend tests，覆盖 healthy/degraded/empty retrieval quality 状态。

## 6. 文档与验收

- [x] 6.1 在 P38 delta 中记录待归档合并到 `docs/workflow.md` 的检索质量、引用一致性和降级契约。
- [x] 6.2 在 P38 delta 中记录待归档合并到 `docs/api.md`、`docs/data-model.md` 的 DTO、审计和事实字段边界。
- [x] 6.3 在 P38 delta 中记录待归档合并到 `docs/frontend-contract.md`、`docs/configuration.md` 的 UI 和配置说明。
- [x] 6.4 更新 `docs/development-plan.md`、`openspec/PROGRESS.md`、`AGENTS.md`、`docs/GOVERNANCE.md` 的 P38 active 状态。
- [x] 6.5 运行 `go test ./...`。
- [x] 6.6 运行 `npm --prefix web test -- --run`。
- [x] 6.7 运行 `npm --prefix web run build`。
- [x] 6.8 运行本地 retrieval quality smoke。
- [x] 6.9 运行 archive 前只读子 agent 复审，且无 Critical / Important 问题。
- [x] 6.10 运行 `openspec validate p38-rag-veclite-retrieval-quality --strict`。
- [x] 6.11 运行 `openspec validate --all --strict`。
