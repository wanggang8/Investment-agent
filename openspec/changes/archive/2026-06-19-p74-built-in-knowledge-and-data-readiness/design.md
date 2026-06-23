# P74 Design

## Product Question

P74 answers whether the product is prepared to support real use with explicit built-in knowledge and required data, instead of relying on scattered docs, implicit rules, or broad prompts.

The goal is not to add more prediction power. The goal is to make the product's knowledge and data dependencies visible, auditable, and safe:

- built-in principles are structured;
- data dependencies are named;
- missing or degraded categories have known feature impact;
- LLM analysis receives only the relevant summarized context;
- UI users can see what is ready and what is not.

## Readiness Model

The readiness service evaluates a symbol against three layers.

### Layer 1: Built-In Knowledge

Knowledge entries have stable IDs and types:

| Type | Examples | Formal role |
| --- | --- | --- |
| `master_principle` | Graham margin of safety, Buffett circle of competence, Livermore trend discipline | LLM context and rule mapping, never formal market evidence |
| `discipline_rule` | no single-source decision, no extreme-emotion trading, no buying after broken thesis | rule engine and UI explanation |
| `risk_sop` | valuation high-risk action, evidence-insufficient action, liquidity caution | UI next-action guidance |
| `symbol_profile` | ETF/index/fund identity and mapping | data dependency selection |

Each entry includes:

- `knowledge_id`
- `title`
- `category`
- `summary`
- `applies_to`
- `rule_mapping`
- `llm_context_allowed`
- `formal_evidence_allowed=false` for master/local/background knowledge
- `safety_boundary`

### Layer 2: Data Dependency Matrix

For the accepted local ETF/index scope, P74 defines required and optional categories:

| Category | Required for | Safe degradation |
| --- | --- | --- |
| `symbol_profile` | capability and data routing | block formal trade advice if unknown |
| `tracked_index` | ETF/index explanation and valuation mapping | degrade expected-return precision |
| `market_price` | position valuation and daily discipline | data stale/insufficient |
| `valuation_percentiles` | Graham-style margin-of-safety discipline | no low/high valuation claim |
| `liquidity` | liquidity risk checks | prohibit large/market-style action wording |
| `sentiment_proxy` | cool-down mechanism | do not claim sentiment pass |
| `formal_evidence` | major event and thesis-break checks | freeze/watch or non-trade record |
| `rag_index` | detailed evidence retrieval | fallback to SQLite summary or degraded retrieval |
| `llm_context` | analyst material quality | degraded analyst material, never final verdict override |

### Layer 3: Product Feature Impact

The readiness response maps each category to affected surfaces:

- Workbench / daily discipline
- Consultation
- Decision detail
- Risk alerts
- Evidence / RAG
- Rules
- Data quality
- Review / decision loop

This lets the UI say "valuation missing affects margin-of-safety and expected-return precision" instead of only showing a generic data failure.

## Backend Architecture

### Built-In Registry

Add a small in-code registry under application service or domain support. It should be deterministic, versioned, and read-only. A future change may move it to SQLite or YAML, but P74 keeps it close to code for testability and release repeatability.

### Readiness Service

`KnowledgeReadinessService` combines:

- built-in registry entries;
- latest market snapshot/source health;
- source verification summary;
- retrieval/index quality if available;
- active rule version;
- capability/symbol profile inference.

The service returns:

- overall `status`: `ready`, `degraded`, or `blocked`;
- `symbol_profile`;
- `knowledge_references`;
- `data_dependencies`;
- `feature_impacts`;
- `llm_context_summary`;
- `safety_notes`.

The service is read-only. It does not trigger collectors, rebuild indexes, write market snapshots, change rules, create notifications, or mutate portfolios.

### API

Add:

`GET /api/v1/knowledge-readiness?symbol=510300`

The response is sanitized and suitable for UI display. It must not include full prompts, API keys, private paths, raw HTTP, raw LLM responses, or complete account details.

## LLM Context

Analyst requests gain a sanitized `KnowledgeContextSummary` field. The DeepSeek prompt builder includes only:

- relevant principle IDs and short summaries;
- data categories that are ready/degraded/missing;
- rule boundary that final verdict remains rule-based;
- safety note that background knowledge cannot satisfy formal evidence.

LLM output metadata should show that a knowledge/data readiness context was attached, but release artifacts must avoid storing full prompts.

## Frontend

The UI should expose readiness without creating a new marketing page. Recommended surfaces:

- `/data-quality`: add a "知识与数据准备度" panel.
- `/rules`: show master/discipline principle mapping.
- decision detail or consultation result: show the readiness context used by the analyst request when available.

The UI must distinguish:

- "available as rule";
- "available as LLM context";
- "background only";
- "missing data blocks or degrades this feature".

## Acceptance

P74 acceptance must include:

- unit tests for registry and readiness status calculation;
- handler tests for sanitized API output;
- LLM prompt tests proving knowledge/data readiness context is included and cannot override final verdict;
- frontend model/page tests for ready/degraded/blocked states;
- browser or scripted acceptance for `510300` complete path and degraded data paths;
- safety scans for forbidden capabilities and sensitive strings.

Pass is blocked if P74 can only show master/data readiness in docs but not through API/UI/LLM context evidence.

## Safety Boundaries

P74 must preserve existing safety boundaries. Built-in knowledge is decision support, not formal external evidence. Missing formal data must degrade or block affected claims. The rule engine remains the source of final verdicts.
