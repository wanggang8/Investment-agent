# P74 Built-In Knowledge And Data Readiness

## Why

P71-P73 proved that the local product can pass strict current-data gates, real UI journeys, real LLM consultation, data-impact checks, and product-effectiveness UX validation.

One important product-readiness question remains: whether the product has a complete, auditable preparation layer for the knowledge and data it relies on. The current implementation includes documented master wisdom, rule logic, public collectors, source-health checks, local knowledge import, and LLM prompts, but those pieces are not yet unified into a visible readiness contract that answers:

- which built-in investment principles are available;
- which principles are rules versus LLM context;
- which symbol profile and ETF/index/fund data are required;
- which required data categories are currently fresh, degraded, missing, or background-only;
- which product features are affected by each missing category;
- whether the LLM actually received the relevant principles and data-readiness summary.

P74 creates that missing product layer.

## What Changes

- Add a built-in knowledge and data readiness registry covering master principles, discipline rules, risk SOPs, symbol profile dependencies, ETF/fund/index data dependencies, and LLM context eligibility.
- Add a read-only readiness service/API that summarizes readiness for a requested symbol using built-in registry facts plus existing source health, evidence verification, rule, and market snapshot data.
- Add LLM context hardening so analyst requests include a sanitized, scenario-relevant knowledge/data-readiness summary, while final verdicts remain rule-based.
- Add UI readback on data quality, rules, consultation, or decision detail surfaces so users can see what knowledge/data was used, what is missing, and what safely degraded.
- Add a P74 acceptance runner that verifies realistic knowledge/data scenarios without mocks as pass evidence: complete ETF path, missing valuation data, background-only local knowledge, single-source evidence, multi-source formal evidence, and out-of-scope capability.
- Update release/governance materials with clear claims and remaining boundaries.

## In Scope

- OpenSpec change, docs delta, design/tasks, product/API/frontend contract updates, release acceptance record, and evidence artifacts.
- Runtime additions that are read-only or context-building only:
  - built-in registry code or data files;
  - readiness DTO/service/handler;
  - prompt input summary extension;
  - UI display of readiness and references;
  - tests and acceptance scripts.
- Readiness categories for current local scope: `master_principles`, `discipline_rules`, `risk_sop`, `symbol_profile`, `fund_profile`, `tracked_index`, `market_price`, `valuation_percentiles`, `liquidity`, `sentiment_proxy`, `formal_evidence`, `rag_index`, `llm_context`.
- `510300` / `000300` as the primary real ETF/index scenario, with deterministic fixture scenarios for degraded/missing states.

## Out of Scope

- No broker interface, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic upgrade, automatic restore, or automatic overwrite of real user databases.
- No return promise, future market direction claim, or proof of future investment performance.
- No paid, login-gated, authorization-gated, Level2, or high-frequency sources.
- No promotion of C-level local knowledge, notes, or master commentary to formal evidence.
- No new public-provider dependency that is required for local release pass unless it has a safe degraded path.
- No physical second-machine repeat or package refresh unless a separate change is opened.

## Impact

P74 may modify backend services, DTOs, handlers, workflow LLM request construction, frontend pages/view models, tests, scripts, and release/governance docs. It should not mutate user portfolios, confirmations, rules, source health, market snapshots, or evidence facts while calculating readiness.
