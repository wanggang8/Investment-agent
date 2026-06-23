# Tasks: P71 Real Product Acceptance True Pass

## 1. Change Setup And Baseline Evidence

- [x] 1.1 Read `AGENTS.md`, `docs/GOVERNANCE.md`, `openspec/project.md`, and `openspec/PROGRESS.md`.
- [x] 1.2 Create `p71-real-product-acceptance-true-pass` OpenSpec change with proposal, design, release-governance delta, and tasks.
- [x] 1.3 Update governance/progress docs to mark P71 active.
- [x] 1.4 Record baseline command evidence for current strict gate and confirm it is currently blocked.
- [x] 1.5 Run `openspec validate p71-real-product-acceptance-true-pass --strict`, `openspec validate --all --strict`, and `git diff --check`.

## 2. Failing Tests For Strict Acceptance

- [x] 2.1 Add a failing test proving P71 current acceptance requires `policy=passed` / `gate=pass` and rejects scope exclusion as pass evidence.
- [x] 2.2 Add a failing test proving acceptance setup rebuilds/verifies a healthy VecLite file index before consultation.
- [x] 2.3 Add a failing browser or script-level test proving `VECTOR_INDEX_UNAVAILABLE` in the real UI consultation result blocks P71.
- [x] 2.4 Add or tighten tests that reject mock-only pass evidence for P71 real product acceptance.
- [x] 2.5 Run each new focused test and confirm it fails for the expected reason before implementation.

## 3. Current Data True Pass

- [x] 3.1 Diagnose why current `000300` source health is missing or degraded in the real local database.
- [x] 3.2 Use or harden read-only public collector paths to write real current source-health facts for `000300`; do not fabricate freshness or silently fall back to stub.
- [x] 3.3 Preserve P66/P67 semantics: scope exclusion remains available for limited release only, not for P71 pass.
- [x] 3.4 Run `go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate` and require exit 0 with `policy=passed` / `gate=pass`.

## 4. VecLite Acceptance Hardening

- [x] 4.1 Add or update an acceptance entry point that seeds/rebuilds SQLite RAG chunks into the configured VecLite file index before UI consultation.
- [x] 4.2 Verify index health is `healthy`, chunk count is positive, metadata is consistent with SQLite summaries, and freshness is not `unknown` or `stale` for the consultation symbol.
- [x] 4.3 Ensure real UI consultation fails the P71 acceptance run if retrieval quality falls back to `sqlite_summary` due to index unavailability.
- [x] 4.4 Run `go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300` and require `status=hit`, `fallback=veclite`, and `index=healthy`.

## 5. Full Real UI Product Acceptance

- [x] 5.1 Create or update a P71 acceptance script derived from P63 that starts a real local backend and Vite frontend with strict gates.
- [x] 5.2 Cover all P63 primary routes at mobile and desktop widths.
- [x] 5.3 Operate key UI paths: portfolio calibration, consultation submit, generated decision detail open, evidence/index rebuild, data-quality current gate check, settings/market refresh path, local knowledge validate/confirm, and governance/ops routes.
- [x] 5.4 Record sanitized screenshots and browser summary JSON under a P71 evidence directory.
- [x] 5.5 Treat console errors, page errors, unexpected failed API responses, UI overflow, missing LLM material, retrieval degradation, or forbidden capabilities as blockers.

## 6. Release Materials And Package Refresh

- [x] 6.1 Add `docs/release/acceptance/2026-06-18-p71-real-product-acceptance.md` with command evidence, UI evidence, current-data evidence, VecLite evidence, real LLM evidence, safety scan, and result.
- [x] 6.2 Update release candidate, handoff, release README, repeatability, development plan, docs README, governance, AGENTS, OpenSpec project, and progress materials.
- [x] 6.3 If and only if P71 strict gates pass, run post-P70 package refresh and verify/repeat acceptance from the generated archive.
- [x] 6.4 If any strict gate fails, record `release_blocked_*` result and do not claim full real product acceptance.

## 7. Verification, Review, Archive

- [x] 7.1 Run `go test ./...`.
- [x] 7.2 Run `npm --prefix web test`.
- [x] 7.3 Run `npm --prefix web run build`.
- [x] 7.4 Run `bash scripts/e2e-smoke.sh`.
- [x] 7.5 Run P71 strict acceptance script.
- [x] 7.6 Run safety/redaction scans and manually classify expected prohibitive wording.
- [x] 7.7 Run `openspec validate p71-real-product-acceptance-true-pass --strict`, `openspec validate --all --strict`, and `git diff --check`.
- [x] 7.8 Archive P71 only after strict evidence is recorded; do not archive as pass if strict gates remain blocked.
