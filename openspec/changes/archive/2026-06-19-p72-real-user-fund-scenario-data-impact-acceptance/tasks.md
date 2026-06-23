# Tasks: P72 Real User Fund Scenario Data Impact Acceptance

## 1. Change Setup And Coverage Review

- [x] 1.1 Read `docs/GOVERNANCE.md`, `openspec/project.md`, and `openspec/PROGRESS.md`.
- [x] 1.2 Create branch `codex/p72-real-user-fund-scenario-data-impact-acceptance`.
- [x] 1.3 Create OpenSpec change `p72-real-user-fund-scenario-data-impact-acceptance`.
- [x] 1.4 Update governance/progress docs to mark P72 active.
- [x] 1.5 Run `openspec validate p72-real-user-fund-scenario-data-impact-acceptance --strict`, `openspec validate --all --strict`, and `git diff --check`.
- [x] 1.6 Perform pre-execution coverage review: verify the matrix covers real fund setup, holding maintenance, offline transaction, local knowledge/RAG, current data, daily discipline, risk alerts, real consultation, manual confirmation, review/readback, failure handling, safety, and deterministic accuracy.

## 2. Failing/Blocking Acceptance Tests

- [x] 2.1 Add a P72 browser acceptance test that fails unless portfolio calibration for real fund `510300` produces correct UI/API state.
- [x] 2.2 Add a P72 browser acceptance test that fails unless local knowledge import writes RAG data and VecLite rebuild returns healthy.
- [x] 2.3 Add a P72 browser acceptance test that fails unless daily/risk/consultation/manual confirmation readbacks are present.
- [x] 2.4 Add a P72 SQLite data-impact checker that fails unless expected tables and deterministic calculations match.
- [x] 2.5 Add blocker checks for forbidden capabilities, page errors, unexpected API failures, LLM degradation, retrieval degradation, raw secrets, raw prompts, raw payloads, and private paths.

## 3. P72 Real Scenario Runner

- [x] 3.1 Add `scripts/p72-real-user-fund-scenario-acceptance.sh` based on P71 with temp SQLite, temp VecLite, real LLM config, `use_stub=false`, and real current-data collector configuration.
- [x] 3.2 Ensure the script runs P34 current-data refresh and strict current-data gate before UI acceptance.
- [x] 3.3 Ensure the script starts the real local Go server and Vite frontend and runs only the P72 spec.
- [x] 3.4 Ensure the script runs read-only SQLite impact verification after the browser test and writes sanitized JSON artifacts.

## 4. Real UI Scenario Coverage

- [x] 4.1 Operate portfolio calibration for `510300` with deterministic values and verify page refresh consistency.
- [x] 4.2 Operate holding edit, import validate/confirm, correction record, and offline transaction through UI.
- [x] 4.3 Operate local knowledge validate/confirm and VecLite rebuild through UI/API.
- [x] 4.4 Operate current-data/data-quality/market refresh path and verify strict gate state.
- [x] 4.5 Operate daily discipline/risk alert/review/notification/audit/rules/workbench readback pages.
- [x] 4.6 Submit real LLM consultation, open generated detail, and verify parsed/passed analyst reports, rule final verdict, evidence chain, and healthy/fresh VecLite retrieval.
- [x] 4.7 Record manual confirmation from generated decision and verify decision loop/readback state.
- [x] 4.8 Trigger representative invalid-input or blocked-operation states and verify safe error behavior.

## 5. Data Impact And Accuracy Verification

- [x] 5.1 Verify `portfolio_snapshots`, `positions`, `position_transactions`, import/correction records, and `audit_events` reflect UI operations.
- [x] 5.2 Verify deterministic calculations: market value, unrealized profit ratio, cash ratio, total assets, and position count.
- [x] 5.3 Verify `intelligence_items`, `intelligence_summary`, `rag_chunks`, `source_verifications`, and VecLite health after local knowledge.
- [x] 5.4 Verify `daily_discipline_reports`, `decision_records`, `risk_alerts`, and `notifications` link to the scenario.
- [x] 5.5 Verify `operation_confirmations` and `position_transactions` link to the generated decision and are visible in decision-loop readback.
- [x] 5.6 Verify no automatic trading/order/broker records or external-push behavior exist.

## 6. Release Materials And Post-Execution Review

- [x] 6.1 Add `docs/release/acceptance/2026-06-18-p72-real-user-fund-scenario.md` with matrix, command evidence, browser evidence, SQLite data-impact summary, accuracy checks, safety scan, gaps, and result.
- [x] 6.2 Store sanitized browser screenshots/results and DB impact summary under `docs/release/ui-audit-assets/2026-06-18-p72/`.
- [x] 6.3 Perform post-execution coverage review: compare executed evidence against the matrix and list missing or intentionally out-of-scope scenarios.
- [x] 6.4 Update release README, release candidate/handoff, repeatability, development plan, docs README, governance, AGENTS, OpenSpec project, and progress materials.

## 7. Verification, Review, Archive

- [x] 7.1 Run `go test ./...`.
- [x] 7.2 Run `npm --prefix web test`.
- [x] 7.3 Run `npm --prefix web run build`.
- [x] 7.4 Run `bash scripts/e2e-smoke.sh`.
- [x] 7.5 Run `bash scripts/p72-real-user-fund-scenario-acceptance.sh`.
- [x] 7.6 Run safety/redaction scans and manually classify expected test strings.
- [x] 7.7 Run `openspec validate p72-real-user-fund-scenario-data-impact-acceptance --strict`, `openspec validate --all --strict`, and `git diff --check`.
- [x] 7.8 Archive P72 only after real scenario evidence and post-execution coverage review are recorded.
