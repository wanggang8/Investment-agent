# P100 Local Source Final Acceptance

> Date: 2026-06-23 18:02:46 CST  
> Conclusion: `local_source_release_acceptance_passed_with_documented_degradation`  
> Scope: local source runtime only.

## Scope

P100 validates the local source product using the Go backend, Vite frontend, local SQLite/VecLite, real browser UI, API/readback evidence, and existing P92/P93 audit gates.

Out of scope:

- Docker Compose.
- Install, upgrade, uninstall, purge, or deployment scripts.
- GitHub Release, Git tag, package refresh, or remote publication.
- Physical second-machine validation.
- Any new investment runtime capability.

P100 does not claim broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic recovery, paid/login/auth-only sources, Level2/HFT sources, future provider availability, or return guarantees.

## Environment

| Item | Value |
| --- | --- |
| Commit | `060d21a03373978ba87c954d46f70e006da986ce` |
| Version | `v0.1.0` |
| Go | `go version go1.25.0 darwin/arm64` |
| Node | `v25.8.1` |
| npm | `11.11.0` |
| Config | `configs/config.yaml` ignored local config |
| Runtime mode | `release` |
| Stub data | `data_sources.use_stub=false` |
| SQLite | `./tmp/p100-local-source-final-acceptance/investment-agent.db` |
| VecLite | `./tmp/p100-local-source-final-acceptance/veclite` |
| LLM key | `DEEPSEEK_API_KEY_ABSENT` at P100 execution time; no fresh real LLM output claim |

## Machine Gates

| Gate | Status | Evidence |
| --- | --- | --- |
| P100 OpenSpec change | pass | `openspec validate p100-local-source-final-acceptance --strict` |
| OpenSpec all | pass | `openspec validate --all --strict` |
| Whitespace | pass | `git diff --check` |
| Go tests | pass | `go test ./...` |
| Go vet | pass | `go vet ./...` |
| Frontend tests | pass | `npm --prefix web test -- --run` |
| Frontend build | pass | `npm --prefix web run build` |
| P92 final requirement audit | pass | `python3 scripts/p92_final_requirement_audit.py --check` |
| P93 code reality/design audit | pass | `python3 scripts/p93_code_reality_audit.py --check` |

P92 remains the final row-level requirement ledger: all full-release-required rows remain `real_pass`. P93 remains the code-reality/design audit gate and reports no active release-blocking findings.

## Local Runtime Acceptance

| Command | Status | Notes |
| --- | --- | --- |
| `bash scripts/e2e-smoke.sh` | pass | 4 Playwright tests passed. |
| `bash scripts/p71-real-product-acceptance.sh` | degraded | `authentication_or_key`: local DeepSeek key is absent. No fresh real LLM claim. |
| `bash scripts/p72-real-user-fund-scenario-acceptance.sh` | degraded | `authentication_or_key`: local DeepSeek key is absent. No fresh real LLM claim. |
| `bash scripts/p83-governance-traceability-acceptance.sh` | pass | Real UI/API/readback governance traceability passed. |
| `bash scripts/p84-portfolio-confirmation-acceptance.sh` | pass | Portfolio, confirmation, offline transaction, audit, and downstream readback passed. |
| `bash scripts/p85-expected-return-analysis-acceptance.sh` | pass with bounded claim | Passed with `static_fallback_no_real_llm_claim` because no real LLM key is configured. |
| `bash scripts/p86-core-goal-knowledge-safety-final-acceptance.sh` | degraded | Nested P81/P75 dynamic-source UI rerun requires a local real LLM key; classified as `authentication_or_key`. |
| `bash scripts/p87-portfolio-state-allocation-acceptance.sh` | pass | Portfolio state/allocation/safe degradation UI path passed. |
| `bash scripts/p88-remaining-full-release-blockers-acceptance.sh` | pass | Remaining blocker UI paths passed. |
| `bash scripts/p89-real-provider-dynamic-probability-acceptance.sh` | degraded | Current run failed on stale P89 UI/provider assertion: expected old `P89 结构化字段`/capital-flow-empty state, while current UI uses `结构化字段` and P90 capital-flow provider now reads back. Current provider state also showed financing/financial partial external-source failure for `159915`. |
| `bash scripts/p90-capital-flow-provider-acceptance.sh` | pass | Eastmoney H5 public capital-flow provider, Settings UI refresh, API readback, SQLite readback, and directional mapping passed. |

P89 degradation does not override P92/P93 final requirement status because P90 is the later closure for the remaining capital-flow rows and P90 passed in this run. It does mean P100 does not claim the historical P89 script is still a clean fresh pass without maintenance.

## Browser And Design QA

The local backend and frontend were started from source:

```bash
INVESTMENT_AGENT_CONFIG=configs/config.yaml go run ./cmd/server
VITE_API_PROXY_TARGET=http://127.0.0.1:8080 npm --prefix web run dev -- --host 127.0.0.1 --port 5173 --strictPort
```

Browser artifact:

```text
docs/release/ui-audit-assets/2026-06-23-p100-local-source-final-acceptance/p100-browser-design-summary.json
```

Routes checked:

- `/workbench`
- `/positions`
- `/settings`
- `/consultation`
- `/review`
- `/rules`
- `/audit`
- `/notifications`
- `/data-quality`

Viewports checked:

- `390x844`
- `768x900`
- `1280x900`

Design and usability rubric:

| Check | Status |
| --- | --- |
| Meaningful content and page identity | pass |
| No framework error overlay | pass |
| No horizontal overflow | pass |
| No forbidden action controls | pass |
| No critical clipped controls | pass |
| No obvious interactive overlap | pass |
| Settings market refresh interaction | pass |
| Console error/warn health | pass; 0 relevant logs |

Settings market refresh produced visible success text: `市场刷新完成；只更新本地行情事实和审计记录，不会执行交易。`

Screenshots:

- `docs/release/ui-audit-assets/2026-06-23-p100-local-source-final-acceptance/p100-workbench-desktop-1280.png`
- `docs/release/ui-audit-assets/2026-06-23-p100-local-source-final-acceptance/p100-workbench-mobile-390.png`
- `docs/release/ui-audit-assets/2026-06-23-p100-local-source-final-acceptance/p100-settings-refresh-1280.png`

## Data Impact Evidence

| Area | Evidence |
| --- | --- |
| Governance traceability | `docs/release/ui-audit-assets/2026-06-22-p83-governance-traceability/governance-traceability-summary.json` reports browser and DB readback `passed`. |
| Portfolio and confirmation | `docs/release/ui-audit-assets/2026-06-22-p84-portfolio-confirmation/portfolio-confirmation-summary.json` reports browser and DB readback `passed`. |
| Expected return | `docs/release/ui-audit-assets/2026-06-22-p85-expected-return-analysis/expected-return-summary.json` reports browser and DB readback `passed`, with no fresh real LLM claim. |
| Portfolio state/allocation | `docs/release/ui-audit-assets/2026-06-22-p87-portfolio-state-allocation-safety/portfolio-state-allocation-summary.json` reports browser and DB readback `passed`. |
| Capital flow provider | `docs/release/ui-audit-assets/2026-06-22-p90-capital-flow-provider/p90-acceptance-summary.json` reports browser, API, SQLite readback, source preverification, and directional net-flow mapping `passed`. |

## Final Assessment

P100 passes local source final acceptance with documented degradation:

- Product local source runtime is usable.
- Product design and responsive route QA pass for the checked routes and viewports.
- P92 final requirement and P93 code-reality/design gates remain valid.
- Core local UI/API/SQLite/readback paths pass.
- Fresh real LLM-backed P71/P72 claims were not made during P100 because the local real LLM key was absent at execution time.
- P89 legacy regression requires maintenance because P90 changed the provider/UI reality; P90 capital-flow closure is the current passing evidence for the final capital-flow rows.

Release statement for this scope:

```text
local_source_release_acceptance_passed_with_documented_degradation
```
