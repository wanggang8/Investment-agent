# Tasks: P74 Built-In Knowledge And Data Readiness

## 1. Change Setup And Scope Review

- [x] 1.1 Confirm no active OpenSpec change is present before opening P74.
- [x] 1.2 Create OpenSpec change `p74-built-in-knowledge-and-data-readiness`.
- [x] 1.3 Define proposal, design, tasks, and release-governance delta.
- [x] 1.4 Update governance/progress docs to mark P74 active.
- [x] 1.5 Run `openspec validate p74-built-in-knowledge-and-data-readiness --strict`, `openspec validate --all --strict`, and `git diff --check`.

## 2. Built-In Knowledge Registry

- [x] 2.1 Add failing tests for stable built-in knowledge IDs, categories, rule mappings, LLM eligibility, and formal-evidence safety boundaries.
- [x] 2.2 Implement deterministic read-only registry entries for master principles, discipline rules, risk SOPs, and the primary ETF/index symbol profile scope.
- [x] 2.3 Verify C/background/master knowledge cannot be marked as formal evidence.

## 3. Data Readiness Service And API

- [x] 3.1 Add failing service tests for `ready`, `degraded`, and `blocked` readiness outputs.
- [x] 3.2 Implement `KnowledgeReadinessService` that combines registry facts with latest source health, evidence verification, active rule, and market snapshot facts without writing local state.
- [x] 3.3 Add DTOs and `GET /api/v1/knowledge-readiness?symbol=...` handler tests for sanitized output.
- [x] 3.4 Register the handler and verify unsupported/missing symbols return safe degraded or blocked states instead of fabricated readiness.

## 4. LLM Context Hardening

- [x] 4.1 Add failing tests proving analyst prompts include sanitized knowledge/data readiness context when provided.
- [x] 4.2 Extend analyst request construction with `KnowledgeContextSummary`.
- [x] 4.3 Attach scenario-relevant readiness context in consultation/daily workflow analyst calls.
- [x] 4.4 Verify LLM context cannot override final rule verdict and full prompts are not persisted in release/audit output.

## 5. Frontend Readiness Experience

- [x] 5.1 Add frontend types/service for knowledge readiness.
- [x] 5.2 Add failing view-model/page tests for ready, degraded, and blocked readiness states.
- [x] 5.3 Add "知识与数据准备度" display to `/data-quality` and targeted readback on rules or decision detail surfaces.
- [x] 5.4 Verify 390px reflow and copy safety for the new panel.

## 6. P74 Acceptance Runner And Evidence

- [x] 6.1 Add a P74 acceptance script that starts local backend/frontend against a temporary SQLite DB.
- [x] 6.2 Add browser or API/UI checks for `510300` complete readiness, missing valuation data, background-only local knowledge, single-source evidence, multi-source formal evidence, and out-of-scope capability.
- [x] 6.3 Write sanitized JSON and screenshot evidence under `docs/release/ui-audit-assets/2026-06-19-p74/`.
- [x] 6.4 Add `docs/release/acceptance/2026-06-19-p74-built-in-knowledge-and-data-readiness.md`.

## 7. Documentation, Release Materials, And Archive

- [x] 7.1 Update `docs/requirements.md`, `docs/data-model.md`, `docs/api.md`, `docs/workflow.md`, `docs/frontend-contract.md`, and `docs/development-plan.md` with P74 contract changes.
- [x] 7.2 Update release README/handoff/repeatability, governance, AGENTS, OpenSpec project, and progress materials.
- [x] 7.3 Run targeted Go tests for registry/service/handler/LLM prompt.
- [x] 7.4 Run targeted frontend tests for readiness UI.
- [x] 7.5 Run `go test ./...`, `npm --prefix web test -- --run`, `npm --prefix web run build`, P74 acceptance runner, safety scans, `openspec validate --all --strict`, and `git diff --check`.
- [x] 7.6 Review whether the task list fully covers built-in knowledge, collected data, LLM reference, UI traceability, real scenarios, and safety boundaries; add missing tasks before final report.
- [x] 7.7 Archive P74 only after all verification and scope review pass.
