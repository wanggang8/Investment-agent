# Acceptance Run: P63 Full UI Regression And Release Refresh

> Date: 2026-06-18
> Code-under-test commit: `7d8707f`
> Branch: `main`
> Operator: Codex
> Environment: macOS Darwin arm64, Go `go1.25.0`, Node `v25.8.1`, npm `11.11.0`
> Source matrix: `docs/project-acceptance-gate-matrix.md`
> Change: P63 `p63-full-ui-regression-release-refresh`
> Result: `release_ready`

This run re-executed the P52 G0-G9 gates after the P58-P62 product experience polish sequence. P63 added full-route browser regression coverage, captured UI screenshots, refreshed release materials, and did not add broker integration, trading execution, backend business APIs, SQLite schema, Eino workflow, or new LLM capability.

Temporary artifacts were written under `tmp/acceptance/p63-2026-06-18/` and `tmp/p63-full-ui-regression/`; those paths are not committed. Committed UI evidence is limited to screenshots and redacted browser summary JSON under `docs/release/ui-audit-assets/2026-06-18-p63/`.

## Summary

| Gate | Status | Release impact | Notes |
| --- | --- | --- | --- |
| G0 Governance | pass | does_not_block | OpenSpec all strict validation and diff check passed. Active change output only showed the expected P63 change during execution. |
| G1 Go all tests | pass | does_not_block | `go test ./...` exit 0. |
| G2 Go focused integration | pass | does_not_block | CLI, server, workflow, handler, and SQLite packages exit 0. |
| G3 Frontend tests/build | pass | does_not_block | Vitest exit 0 with 157 tests; production build exit 0. |
| G4 Browser E2E smoke | pass | does_not_block | `bash scripts/e2e-smoke.sh` exit 0, including the P63 full UI regression spec with failed API response classification. |
| G5 Local fixture/current smoke | degraded | does_not_block | Recovery, retrieval, and fixture regression passed. Current regression returned `status=degraded:cases=1:degraded=1:failed=0`; P52 allows current degraded when classified. |
| G6 Real public source opt-in | pass | does_not_block | Public evidence refresh exit 0 with a temporary SQLite database and explicit 2026-06-01 to 2026-06-18 window. |
| G7 Real LLM opt-in | pass | does_not_block | LLM smoke exit 0 with model `gpt-5.4-mini`; audit summary reports parse and quality passed. |
| G8 Local install/release upgrade | pass | does_not_block | Install diagnostics and release-upgrade check both passed. |
| G9 Safety/redaction | pass | does_not_block | Manual review found no committed complete key, raw payload, full prompt, raw DB, private log, or new prohibited runtime capability. |

## Gate Details

| Gate | Command | Final status | Evidence | Notes |
| --- | --- | --- | --- | --- |
| G0 | `openspec validate --all --strict` | pass | terminal output | 33 passed, 0 failed. |
| G0 | `git diff --check` | pass | terminal output | No whitespace errors. |
| G0 | `find openspec/changes -maxdepth 1 -mindepth 1 -type d ! -name archive -print` | pass | terminal output | Output was the expected active P63 change during execution. |
| G1 | `go test ./...` | pass | terminal output | Exit 0. |
| G2 | `go test ./cmd/agent ./cmd/server ./internal/application/workflow ./internal/application/handler ./internal/infrastructure/persistence/sqlite` | pass | terminal output | Exit 0. |
| G3 | `npm --prefix web test` | pass | terminal output | 47 files and 157 tests passed. |
| G3 | `npm --prefix web run build` | pass | terminal output | Production build exit 0. |
| G4 | `bash scripts/e2e-smoke.sh` | pass | terminal output | Four Playwright specs passed, including P30, P39, P62, and P63. |
| G5 | `bash scripts/recovery-smoke.sh` | pass | terminal output | Exit 0. |
| G5 | `go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300` | pass | terminal output | Exit 0 after sequential retry. |
| G5 | `go run ./cmd/agent --task data-source-quality-regression --source fixture --symbol 000300` | pass | terminal output | `status=passed:cases=6:degraded=0:failed=0:no_auto_trading`. |
| G5 | `go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300` | degraded | terminal output | Exit 0 with `status=degraded:cases=1:degraded=1:failed=0:no_auto_trading`. |
| G6 | `INVESTMENT_AGENT_CONFIG=tmp/acceptance/p63-2026-06-18/config.real-public.yaml go run ./cmd/agent --task public-evidence-refresh --symbol 000001 --start-date 2026-06-01 --end-date 2026-06-18` | pass | temporary SQLite counts | Counts: `intelligence_items=1`, `intelligence_summary=1`, `rag_chunks=1`, `source_verifications=1`, `audit_events=2`; audit included `source=public_evidence count=1` and `no_auto_trading`. |
| G7 | `INVESTMENT_AGENT_CONFIG=tmp/p63-full-ui-regression/config.p63.yaml go run ./cmd/agent --task llm-smoke --symbol 510300` | pass | temporary SQLite audit summary | Audit summary: `llm-smoke:symbol=510300:model=gpt-5.4-mini` and `llm_smoke:quality=passed:parse=parsed:no_auto_trading`. |
| G8 | `bash scripts/local-install-diagnostics.sh --config configs/config.example.yaml --include-release-upgrade --target-version p63-acceptance --output-dir tmp/acceptance/p63-2026-06-18/install` | pass | `tmp/acceptance/p63-2026-06-18/install/20260618T043022Z/install-summary.json` | Preflight, recovery smoke, release upgrade check, and E2E smoke all passed. |
| G8 | `bash scripts/local-release-upgrade-check.sh --config configs/config.example.yaml --target-version p63-acceptance --output-dir tmp/acceptance/p63-2026-06-18/release-upgrade` | pass | `tmp/acceptance/p63-2026-06-18/release-upgrade/20260618T043022Z/release-upgrade-summary.json` | Preflight and release upgrade check passed. |
| G9 | Forbidden capability scan | pass | terminal output and manual review | Matches were existing prohibitive text, tests, seed text, or UI safety copy; P63 E2E also checks forbidden button/link names. |
| G9 | Sensitive information scan | pass | terminal output and manual review | Matches in `configs/config.local.yaml` are ignored local test configuration and not tracked by git. Committed P63 UI artifacts contain no complete key or raw prompt. |

## Full UI Regression

| Area | Result |
| --- | --- |
| Script | `P63_SERVER_PORT=18084 P63_WEB_PORT=14179 bash scripts/p63-full-ui-regression.sh` |
| Browser spec | `web/e2e/p63-full-ui-regression.spec.ts` |
| Result | 1 Playwright test passed in 1.2m after the failed API response classification fix |
| Routes | 20 primary routes |
| Viewports | 390px, 768px, 1280px |
| Screenshots | 60 PNG files under `docs/release/ui-audit-assets/2026-06-18-p63/` |
| Summary JSON | `docs/release/ui-audit-assets/2026-06-18-p63/browser-results.json` |
| Page-level overflow | 0 failed overflow checks across 60 route/viewport combinations |
| Console errors | 0 |
| Page errors | 0 |
| Failed API responses | 20 classified expected client-state responses |
| Unexpected failed API responses | 0 |

Routes covered:

`/`, `/workbench`, `/consultation`, `/decisions/:decisionId`, `/evidence`, `/decision-loop`, `/positions`, `/data-quality`, `/risk-alerts`, `/risk-alerts/:alertId`, `/rules`, `/audit`, `/notifications`, `/daily-auto-run`, `/daily-discipline/reports`, `/daily-discipline/reports/:reportId`, `/review`, `/local-install`, `/local-knowledge`, `/settings`.

## Real LLM-Backed UI Journey

| Field | Result |
| --- | --- |
| UI path | Browser-filled `/positions`, then submitted `/consultation` |
| Consultation question | `P63 真实 UI 回归：510300 当前是否继续持有？` |
| HTTP status | 200 |
| Decision id | `decision_e6f6d404bb554d61` |
| Decision detail | `/decisions/decision_e6f6d404bb554d61` opened successfully |
| Analyst reports | 3 |
| Parse statuses | `parsed`, `parsed`, `parsed` |
| Quality statuses | `passed`, `passed`, `passed` |
| LLM material displayed in UI | yes |
| Workflow status | `degraded` |
| Source verification status | `satisfied` |
| Degraded reason | `VECTOR_INDEX_UNAVAILABLE` in the temporary P63 SQLite/VecLite acceptance setup |

The degraded workflow status does not mean the LLM call failed. The browser journey produced a decision, displayed LLM material, and opened the decision detail. The degradation limits claims about retrieval-enhanced context availability in this temporary run.

## Failed API Response Classification

| Status | Method | Classification | Count | Endpoint |
| --- | --- | --- | --- | --- |
| 404 | GET | `expected_client_state` | 11 | `/api/v1/portfolio/current` |
| 409 | GET | `expected_client_state` | 9 | `/api/v1/dashboard/today` |

The P63 browser regression records all `/api/v1/` responses with status `>=400`. The classified 404/409 responses come from expected empty or precondition states while the test moves through seeded and temporary local data. There were zero unexpected failed API responses.

## Degraded Or Retried Items

| Item | Classification | Handling | Release impact |
| --- | --- | --- | --- |
| G5 initial parallel retrieval/current execution | acceptance_runner_concurrency | Parallel current-DB commands produced `database is locked`. The exact gate commands were rerun sequentially and reached the expected pass/degraded outcomes. | does_not_block |
| G5 current data quality degraded | current_local_data_degraded | Fixture regression passed; current mode had one degraded case and zero failed cases. | does_not_block; limits current local DB quality claim |
| P63 UI consultation workflow degraded | retrieval_index_unavailable | LLM reports parsed and passed quality; source verification was satisfied; workflow degradation came from temporary VecLite index availability. | does_not_block; limits retrieval-enhanced context claim for this UI run |
| P63 classified 404/409 API responses | expected_client_state | Browser JSON now records all `/api/v1/` `>=400` responses and classifies expected client-state 404/409 responses. | does_not_block; unexpected failed API responses were zero |

## Safety Review

- P63 does not add broker integration, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair commitment, real DB overwrite, return guarantee, login source, paid source, authorized source, Level2, or high-frequency source.
- Real LLM key remained in ignored local configuration and temporary config only. `configs/config.local.yaml` is not tracked by git.
- P63 committed artifacts include screenshots and redacted browser summary JSON only.
- Temporary logs, Playwright traces, SQLite files, complete raw responses, complete prompts, and local private paths are not committed.
- The UI E2E checks that prohibited trading or automation labels are not exposed as actionable buttons or links.

## Conclusion

Release status: `release_ready`.

The release is ready with two documented non-blocking degradations: current local data-source quality returned degraded with zero failed cases, and the P63 temporary UI consultation recorded `VECTOR_INDEX_UNAVAILABLE` while the real LLM reports parsed and passed quality. P63 also records 20 classified expected client-state API responses and zero unexpected failed API responses. The project must not claim future public-source availability, future LLM availability, investment returns, trading execution, broker connectivity, or automatic decision execution based on this run.
