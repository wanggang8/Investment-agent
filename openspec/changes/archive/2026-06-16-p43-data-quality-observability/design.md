# P43 Data Quality Observability Design

## Context

P34、P37、P38、P40 已分别建立真实数据覆盖、真实 LLM 质量评估、RAG/VecLite 检索质量和本地运维恢复能力。P42 已提供日常工作台，但工作台只展示简要状态，不承担诊断解释。P43 需要把已有质量信号整理成一个可检查、可审计、可导航的只读视图。

## Decision

P43 采用前端聚合优先的实现方式：

1. 新增 `/data-quality` 页面，并在主导航中加入“数据质量”。
2. 页面区域：
   - **数据源健康**：展示 source health、新鲜度、失败分类、影响标的或范围。
   - **证据与检索**：展示证据数量、独立信源、检索质量、index freshness、fallback source 和 degraded reason。
   - **LLM 质量**：展示模型配置/可用状态、parse status、quality status 和可安全检查的 smoke/质量门禁摘要。
   - **影响范围与下一步**：列出受影响的纪律报告、决策、风险预警、规则提案或审计入口。
3. 数据来源优先使用现有 services/API DTO：settings/system、settings/capability、market/source health、review summary、evidence summary/table、dashboard/report DTO 中已有质量字段。
4. 如现有 DTO 无法表达 P43 验收所需字段，可新增只读聚合 endpoint 或 DTO；不得新增 migration 或写入工作流。
5. 为避免 P42 曾遇到的 dev StrictMode 双挂载并发问题，聚合页面需要顺序加载或 single-flight 保护。

## Safety

- 页面只展示状态、影响和导航；不触发刷新、重建索引、真实 LLM smoke、外部推送、规则应用、账户变更或交易。
- 所有错误展示必须使用稳定错误码和安全文案，不显示密钥、完整 prompt、私有本地路径、SQL、供应商原始报错或用户账户敏感明细。
- `unknown`、`missing`、`stale`、`parse_error`、`source_unavailable`、`quality_failed` 等状态不得渲染为成功。

## Validation

- Vitest 覆盖数据完整、空库、source_unavailable、parse_error、stale、missing、unknown、LLM/RAG/VecLite 降级和安全脱敏。
- Playwright smoke 覆盖 `/data-quality` 可达、主要区域可见、窄屏可用和禁止入口扫描。
- OpenSpec 和现有本地 E2E 保持通过。
