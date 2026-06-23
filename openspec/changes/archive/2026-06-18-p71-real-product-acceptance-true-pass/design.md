# Design: P71 Real Product Acceptance True Pass

## Context

The P70 state is intentionally limited. Fresh command evidence on 2026-06-18 still shows:

```bash
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate
```

exits non-zero with:

```text
data_source_quality:mode=current:status=degraded:policy=blocked:gate=block:cases=1:degraded=1:failed=0:no_auto_trading
```

The current local database has no evaluable P34 source-health facts for `000300`, so the P66 strict gate blocks a clean current-data claim. Separately, P63's real UI consultation succeeded but recorded `VECTOR_INDEX_UNAVAILABLE`, so retrieval-enhanced context was not fully available in that temporary acceptance setup.

## Acceptance Model

P71 introduces a stricter status model:

| Outcome | Meaning |
| --- | --- |
| `release_ready_full_real_product_acceptance` | Current-data strict gate passes, VecLite index is healthy during real UI consultation, real UI/LLM gates pass, safety scans pass, and post-P70 package refresh succeeds from the accepted commit. |
| `release_blocked_current_data` | The P66 strict current-data gate does not return `policy=passed` / `gate=pass`. |
| `release_blocked_retrieval_index` | VecLite/RAG index is missing, corrupted, incompatible, empty when required, stale, or UI consultation returns `VECTOR_INDEX_UNAVAILABLE`. |
| `release_blocked_ui_or_llm` | Real browser operation, decision detail, real LLM parse/quality, console/page/API assertions, or responsive UI checks fail. |
| `release_blocked_safety_or_package` | Safety/redaction/package verification fails. |

P71 does not allow `release_ready_limited_current_data_scope` as its success status. That status remains historically valid for P70 only.

## Work Plan

1. Create the P71 change and record current failing evidence.
2. Add strict acceptance tests before implementation:
   - Current-data true pass test: current source health must produce `policy=passed` / `gate=pass`.
   - VecLite acceptance test: acceptance setup must rebuild/verify a healthy file index and real consultation must not report `VECTOR_INDEX_UNAVAILABLE`.
   - UI acceptance test: the browser summary must fail if the consultation workflow is degraded by retrieval, if LLM materials are missing, if unexpected API failures occur, or if primary operations rely on mock-only evidence.
3. Diagnose and fix the current-data blocker at the source. Preferred fix is to run or harden real read-only public collector paths so `000300` current source health facts exist and are fresh. If external sites are unavailable, P71 records a blocked status rather than fabricating current data.
4. Harden P63/P71 acceptance setup by rebuilding VecLite from SQLite RAG chunks before UI consultation and asserting index health/freshness.
5. Re-run P52 G0-G9 plus the stricter P71 gates.
6. If and only if the strict gates pass, create a post-P70 package from the accepted commit and verify/repeat it.
7. Archive P71 by merging deltas into `docs/` and update release/progress/governance materials.

## Real-UI Evidence Standard

The P71 browser evidence must come from a real local Go server and Vite frontend. Pass evidence may use seeded local prerequisites for account/portfolio/navigation setup, but it must not count frontend mocks, mocked network responses, or fixture-only current data as proof of real product acceptance.

Required UI operations:

- Navigate all P63 primary routes across desktop and mobile widths.
- Submit or operate key page actions where supported: portfolio calibration, consultation submit, generated decision detail open, evidence/index rebuild, data-quality current gate check, settings/market refresh path, local knowledge validate/confirm, and governance/ops routes.
- Record console errors, page errors, all failed `/api/v1/` responses, overflow checks, screenshot artifacts, and sanitized browser summary JSON.
- Treat unexpected API failures, retrieval degradation, missing LLM material, and forbidden trading/automation affordances as blockers.

## Safety Boundaries

All changes must preserve local-only, read-only, non-trading boundaries. P71 may use explicit user-local configs for real public data and real LLM opt-in, but committed materials must not include full keys, private paths, raw SQL dumps, complete prompts, raw vendor payloads, local databases, VecLite files, logs, or Playwright traces.

## Verification

Minimum P71 verification:

```bash
openspec validate p71-real-product-acceptance-true-pass --strict
openspec validate --all --strict
git diff --check
go test ./...
npm --prefix web test
npm --prefix web run build
bash scripts/e2e-smoke.sh
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate
go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300
P71_SERVER_PORT=<port> P71_WEB_PORT=<port> bash scripts/p71-real-product-acceptance.sh
```

Package refresh verification is required only after all strict acceptance gates pass.
