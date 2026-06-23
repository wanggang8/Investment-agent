# Acceptance Run: P53 G0-G9

> Date: 2026-06-17
> Code-under-test commit: `5832477`
> Branch: `main`
> Operator: Codex
> Environment: macOS Darwin arm64, Go `go1.25.0`, Node `v25.8.1`, npm `11.11.0`
> Source matrix: `docs/project-acceptance-gate-matrix.md`
> Result: `release_ready`

This run executed the P52 G0-G9 acceptance gates. P53 changed documentation and OpenSpec governance only; it did not modify runtime code, SQLite schema, HTTP API, Eino workflow, frontend behavior, or scripts.

Temporary artifacts were written under `tmp/acceptance/p53-2026-06-17/` and are not committed. Release materials below reference only redacted summaries and relative artifact paths.

## Summary

| Gate | Status | Release impact | Notes |
| --- | --- | --- | --- |
| G0 Governance | pass | does_not_block | OpenSpec all strict validation and diff check passed. Active change output only showed the expected P53 change during execution. |
| G1 Go all tests | pass | does_not_block | `go test ./...` exit 0. |
| G2 Go focused integration | pass | does_not_block | CLI, server, workflow, handler, and SQLite packages exit 0. |
| G3 Frontend tests/build | pass | does_not_block | Vitest exit 0. First build was killed by the local OS (`Killed: 9`, exit 137); exact command retry exit 0. |
| G4 Browser E2E smoke | pass | does_not_block | First run failed because the Vite process was killed and curl could not connect; exact command retry exit 0. |
| G5 Local fixture/current smoke | degraded | does_not_block | Recovery, retrieval, and fixture regression passed. Current regression returned `status=degraded:cases=1:degraded=1:failed=0`; P52 allows current degraded when classified. |
| G6 Real public source opt-in | pass | does_not_block | Initial temporary config missed the real-mode market prerequisite and failed validation. After enabling real market collectors in the temporary config, public evidence refresh exit 0 and wrote evidence/audit rows. |
| G7 Real LLM opt-in | pass | does_not_block | Real LLM smoke exit 0 with model `gpt-5.4-mini`; audit summary reports parse and quality passed. |
| G8 Local install/release upgrade | pass | does_not_block | Release upgrade check passed. First install diagnostics failed through the same E2E process kill; retry exit 0. |
| G9 Safety/redaction | pass | does_not_block | Scans matched existing test fixtures and prohibitive safety text; manual review found no submitted complete key, raw payload, full prompt, private path, or new prohibited capability. |

## Gate Details

| Gate | Command | Final status | Artifact | Notes |
| --- | --- | --- | --- | --- |
| G0 | `openspec validate --all --strict` | pass | `tmp/acceptance/p53-2026-06-17/logs/g0-openspec-all.log` | 33 passed, 0 failed. |
| G0 | `git diff --check` | pass | `tmp/acceptance/p53-2026-06-17/logs/g0-diff-check.log` | No output. |
| G0 | `find openspec/changes -maxdepth 1 -mindepth 1 -type d ! -name archive -print` | pass | `tmp/acceptance/p53-2026-06-17/logs/g0-active-change.log` | Output was the expected active P53 change. |
| G1 | `go test ./...` | pass | `tmp/acceptance/p53-2026-06-17/logs/g1-go-test-all.log` | Exit 0. |
| G2 | `go test ./cmd/agent ./cmd/server ./internal/application/workflow ./internal/application/handler ./internal/infrastructure/persistence/sqlite` | pass | `tmp/acceptance/p53-2026-06-17/logs/g2-go-focused.log` | Exit 0. |
| G3 | `npm --prefix web test -- --run` | pass | `tmp/acceptance/p53-2026-06-17/logs/g3-web-vitest.log` | Exit 0. |
| G3 | `npm --prefix web run build` | pass | `tmp/acceptance/p53-2026-06-17/logs/g3-web-build-retry.log` | Initial exit 137 from OS kill; exact retry exit 0. |
| G4 | `bash scripts/e2e-smoke.sh` | pass | `tmp/acceptance/p53-2026-06-17/logs/g4-e2e-smoke-retry.log` | Initial exit 7 after Vite was killed; exact retry exit 0. |
| G5 | `bash scripts/recovery-smoke.sh` | pass | `tmp/acceptance/p53-2026-06-17/logs/g5-recovery-smoke.log` | Exit 0. |
| G5 | `go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300` | pass | `tmp/acceptance/p53-2026-06-17/logs/g5-retrieval-quality.log` | Exit 0; audit event written. |
| G5 | `go run ./cmd/agent --task data-source-quality-regression --source fixture --symbol 000300` | pass | `tmp/acceptance/p53-2026-06-17/logs/g5-data-quality-fixture.log` | Exit 0. |
| G5 | `go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300` | degraded | `tmp/acceptance/p53-2026-06-17/logs/g5-data-quality-current.log` | Exit 0 with `status=degraded:cases=1:degraded=1:failed=0`; impact limited to current local DB data quality claim. |
| G6 | `go run ./cmd/agent --config tmp/acceptance/p53-2026-06-17/config.real-public.yaml --task public-evidence-refresh --symbol 000001 --start-date 2026-06-01 --end-date 2026-06-17` | pass | `tmp/acceptance/p53-2026-06-17/logs/g6-real-public-evidence-retry.log` | Retry exit 0 after temporary config fix. SQLite counts: `intelligence_items=1`, `intelligence_summary=1`, `rag_chunks=1`, `source_verifications=1`, `audit_events=4`. |
| G7 | `go run ./cmd/agent --config tmp/acceptance/p53-2026-06-17/config.real-llm.yaml --task llm-smoke --symbol 510300` | pass | `tmp/acceptance/p53-2026-06-17/logs/g7-real-llm-smoke.log` | Exit 0. Audit summary: `llm_smoke:quality=passed:parse=parsed:no_auto_trading`. |
| G8 | `bash scripts/local-install-diagnostics.sh --config configs/config.example.yaml --include-release-upgrade --target-version p53-acceptance --output-dir tmp/acceptance/p53-2026-06-17/install` | pass | `tmp/acceptance/p53-2026-06-17/logs/g8-install-diagnostics-retry.log` | Initial exit 7 from E2E process kill; retry exit 0. |
| G8 | `bash scripts/local-release-upgrade-check.sh --config configs/config.example.yaml --target-version p53-acceptance --output-dir tmp/acceptance/p53-2026-06-17/release-upgrade` | pass | `tmp/acceptance/p53-2026-06-17/logs/g8-release-upgrade.log` | Exit 0. |
| G9 | Safety capability scan | pass | `tmp/acceptance/p53-2026-06-17/logs/g9-safety-capability-scan.log` | Matches are existing prohibitive text, tests, and route names; no new prohibited runtime capability. |
| G9 | Redaction scan | pass | `tmp/acceptance/p53-2026-06-17/logs/g9-redaction-scan.log` | Matches are synthetic test values or policy text; release docs do not include complete keys, private paths, raw HTTP payloads, complete prompts, or raw SQL dumps. |

## Degraded Or Retried Items

| Item | Classification | Handling | Release impact |
| --- | --- | --- | --- |
| G3 initial build exit 137 | local_resource_killed | Exact command retry passed. | does_not_block |
| G4 initial E2E exit 7 | local_resource_killed | Exact command retry passed. | does_not_block |
| G6 initial config validation failure | acceptance_environment_config | Temporary acceptance config was corrected; retry passed. | does_not_block |
| G8 initial diagnostics exit 7 | local_resource_killed | Same E2E process kill pattern; retry passed. | does_not_block |
| G5 current data quality degraded | current_local_data_degraded | Fixture regression passed; current mode had one degraded case and zero failed cases. | does_not_block; limits current local DB quality claim |

## Safety Review

- No P53 file adds broker integration, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair commitment, real DB overwrite, return guarantee, login source, paid source, authorized source, Level2, or high-frequency source.
- Real LLM key remained in local ignored configuration and temporary config only; it is not copied into committed release materials.
- Temporary logs and SQLite files are not committed.
- Real-source and real-LLM results are summarized by counts and audit summaries only.

## Conclusion

Release status: `release_ready`.

The release is ready with one documented non-blocking degradation: current local data-source quality returned degraded with zero failed cases. The project must not claim future public-source availability, future LLM availability, investment returns, trading execution, or automatic decision execution based on this run.
