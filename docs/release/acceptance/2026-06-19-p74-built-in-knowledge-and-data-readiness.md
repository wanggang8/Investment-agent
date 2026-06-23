# P74 Built-In Knowledge And Data Readiness Acceptance

> Status: `release_ready_built_in_knowledge_data_readiness`
> Change: `p74-built-in-knowledge-and-data-readiness`

P74 validates whether built-in investment knowledge, required collected data, active rule boundaries, formal evidence, and LLM context readiness are explicit and traceable in the product. It does not validate future returns or future market direction.

## What Was Added

- Deterministic built-in knowledge registry for 7 master principles, discipline rules, risk SOPs, and `510300` symbol profile.
- Read-only `GET /api/v1/knowledge-readiness?symbol=...` API.
- Data readiness dependency matrix covering symbol profile, fund profile, tracked index, market price, valuation percentiles, liquidity, sentiment proxy, active rule, formal evidence, RAG index, and LLM context.
- Sanitized LLM knowledge/data readiness context in analyst requests and prompt construction.
- `/data-quality` UI panel for knowledge references, dependency matrix, safe degradation, and LLM context presence.
- Decision detail readback showing that LLM used the readiness summary without exposing full prompts.
- P74 acceptance runner:
  - `scripts/p74-built-in-knowledge-data-readiness.sh`
  - `scripts/p74_readiness_api_check.py`
  - `web/e2e/p74-built-in-knowledge-data-readiness.spec.ts`

## Scenario Evidence

P74 executed against a local Go backend, Vite frontend, temporary SQLite database, real HTTP API, and real Playwright browser operations.

| Scenario | Result | Evidence |
| --- | --- | --- |
| Complete readiness for `510300` | pass | all required dependencies ready; formal evidence and active rule ready |
| Missing valuation data | pass | `valuation_percentiles=degraded`; no safety-margin claim is fabricated |
| Background-only local knowledge/evidence | pass | formal evidence degrades; built-in knowledge remains non-formal evidence |
| Single-source evidence | pass | formal evidence degrades and does not permit trade confirmation |
| Multi-source formal evidence | pass | formal evidence returns ready |
| Out-of-scope symbol profile | pass | unknown symbol returns `blocked` without fabricating profile |
| `/data-quality` UI | pass | panel shows ready status, master/discipline references, full dependency matrix, active rule, formal evidence, and no-formal-evidence boundary |
| Decision detail readback | pass | shows LLM used knowledge/data readiness summary and only exposes sanitized metadata |
| 390px mobile reflow | pass | data-quality readiness panel has reachable navigation and no horizontal overflow |

The refreshed API evidence reports `knowledge_reference_count=12`, covering Graham, Buffett, Livermore, Dalio, Marks, Lynch, Templeton, two discipline rules, two risk SOPs, and the `510300` symbol profile.

Artifacts:

- `docs/release/ui-audit-assets/2026-06-19-p74/api-readiness-results.json`
- `docs/release/ui-audit-assets/2026-06-19-p74/browser-results.json`
- `docs/release/ui-audit-assets/2026-06-19-p74/api-readiness-check.log`
- `docs/release/ui-audit-assets/2026-06-19-p74/data-quality-readiness.png`
- `docs/release/ui-audit-assets/2026-06-19-p74/decision-readiness-readback.png`
- `docs/release/ui-audit-assets/2026-06-19-p74/mobile-data-quality-readiness.png`

## Verification Completed

```text
go test ./internal/application/service -run 'Test(BuiltInKnowledgeRegistry|KnowledgeReadiness)'
ok

go test ./internal/application/handler -run 'TestGetKnowledgeReadiness'
ok

go test ./internal/infrastructure/llm/deepseek -run TestBuildPromptIncludesKnowledgeReadinessContext
ok

go test ./internal/application/workflow -run TestAnalystRequestsIncludeKnowledgeReadinessContext
ok

npm --prefix web test -- --run src/pages/DataQualityPage.test.tsx
8 tests passed

npm --prefix web test -- --run src/components/decision/DecisionTrace.test.tsx
12 tests passed

bash scripts/p74-built-in-knowledge-data-readiness.sh
API scenario matrix passed
1 Playwright browser test passed
```

## Current Result

P74 passes for the accepted built-in knowledge and data readiness scope. The pass is based on real local API calls, real browser UI operations, temporary SQLite scenario mutations, screenshots, and sanitized JSON evidence.

## Boundaries

P74 does not add or claim:

- broker connectivity;
- automatic trading;
- one-click trading;
- delegated orders;
- external push;
- automatic confirmation;
- automatic rule application;
- automatic repair, migration, restore, or real DB overwrite;
- login, paid, authorized, Level2, or high-frequency data sources;
- future public-source or model-provider availability;
- future market prediction;
- future returns or investment performance.
