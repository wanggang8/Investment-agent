# P38 Design

## Context

现有检索链路已有本地 JSON index、SQLite fallback、index health、source verification 和 evidence quality metadata。P38 的核心不是换数据库或接外部向量服务，而是让当前本地检索可测、可解释、可审计，并确保检索结果不会违反 formal/background 和规则优先边界。

## Goals

1. 建立可重复的本地 retrieval quality evaluation。
2. 让检索排序考虑 source level、evidence role、time weight、relevance score、verification status 和 index freshness。
3. 确保证据引用与 SQLite 事实、RAG chunk metadata、source verification 一致。
4. 让 API/审计/前端能说明检索是否 degraded、为什么 degraded、fallback 使用了什么来源。

## Non-Goals

- 不实现真实外部 VecLite API。
- 不让检索结果直接改变最终裁决。
- 不新增交易、外部推送或规则自动应用入口。
- 不扩大数据源范围。

## Approach

### 1. Retrieval Quality Fixture

新增小型 fixture/test set，包含 query、symbol、expected evidence ids 或 expected source constraints。测试集优先使用本地 SQLite seed 和 public evidence fixture，保证离线可重复。

### 2. Quality-Aware Ranking

在现有检索结果上增加轻量 score composition：文本相关度继续是基础分，source level、formal evidence role、time weight、source verification satisfied、freshness 增加或减少排序权重。C 级/background 不得变成 formal evidence，只能作为背景说明。

### 3. Consistency Check

检索返回前校验 evidence id、summary id、chunk metadata、source level、evidence role、verification group 和 freshness 状态。如果发现不一致，写 degraded reason 或跳过该结果，不静默返回误导性 evidence。

### 4. Observability

本地任务或 API 输出 retrieval quality summary：query、expected/actual hit、top-k、fallback source、index health/freshness、degradation reason。前端只展示摘要和安全解释，不读取本地文件。

## Risks

- 质量分数过复杂会导致不可解释。P38 只做可测轻量权重，避免黑盒化。
- fixture 太少会过拟合。P38 以 representative cases 开始，并记录不足。
- 降级展示可能被误解为数据恢复或交易建议。所有文案继续强调只读与不自动交易。

## Verification

- `go test ./...`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `openspec validate p38-rag-veclite-retrieval-quality --strict`
- `openspec validate --all --strict`
- 本地 retrieval quality smoke，确认 top-k、fallback reason 和 audit 不含交易/自动规则动作。
