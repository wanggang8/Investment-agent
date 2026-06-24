# Tasks

## 1. Scope And Governance

- [x] Confirm no unrelated active change exists:
  - [x] `find openspec/changes -maxdepth 1 -mindepth 1 -type d ! -name archive -print`
  - [x] `git status --short`
- [x] Validate the P100 change:
  - [x] `openspec validate p100-local-source-final-acceptance --strict`
- [x] Record commit, version, environment, and local config path for the acceptance report:
  - [x] `git rev-parse HEAD`
  - [x] `cat VERSION`
  - [x] `go version`
  - [x] `node --version`
  - [x] `npm --version`

## 2. Local Source Configuration

- [x] Create or refresh ignored local config from the template:
  - [x] `cp configs/config.example.yaml configs/config.yaml`
- [x] Review `configs/config.yaml` and ensure acceptance paths are local/disposable or intentionally read-only.
- [x] Confirm real-source acceptance claims do not use stub data:
  - [x] `rg -n "use_stub:[[:space:]]*true|runtime.mode|api_key|sqlite|veclite" configs/config.yaml configs/config.example.yaml`
- [x] If real LLM output is part of the claim, configure a local ignored key or environment variable and verify that no full key appears in logs or reports.

## 3. Core Machine Gates

- [x] Run governance and whitespace checks:
  - [x] `openspec validate --all --strict`
  - [x] `git diff --check`
- [x] Run backend gates:
  - [x] `go test ./...`
  - [x] `go vet ./...`
- [x] Run frontend gates:
  - [x] `npm --prefix web test -- --run`
  - [x] `npm --prefix web run build`
- [x] Run requirement and code-reality gates:
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p93_code_reality_audit.py --check`

## 4. Local Source Runtime Acceptance

- [x] Run browser smoke:
  - [x] `bash scripts/e2e-smoke.sh`
- [x] Run local product acceptance runners that do not require Docker or install flows:
  - [x] `bash scripts/p71-real-product-acceptance.sh` (`authentication_or_key`; no real LLM key in local config)
  - [x] `bash scripts/p72-real-user-fund-scenario-acceptance.sh` (`authentication_or_key`; no real LLM key in local config)
  - [x] `bash scripts/p83-governance-traceability-acceptance.sh`
  - [x] `bash scripts/p84-portfolio-confirmation-acceptance.sh`
  - [x] `bash scripts/p85-expected-return-analysis-acceptance.sh`
  - [x] `bash scripts/p86-core-goal-knowledge-safety-final-acceptance.sh` (`authentication_or_key` in nested P81/P75 real-LLM-dependent UI rerun)
  - [x] `bash scripts/p87-portfolio-state-allocation-acceptance.sh`
  - [x] `bash scripts/p88-remaining-full-release-blockers-acceptance.sh`
  - [x] `bash scripts/p89-real-provider-dynamic-probability-acceptance.sh` (`source/provider-current-state` plus stale P89 UI assertion; superseded by P90 capital-flow closure for the remaining P89-chain rows)
  - [x] `bash scripts/p90-capital-flow-provider-acceptance.sh`
- [x] Classify any real provider or LLM failure as `network`, `rate_limit`, `authentication_or_key`, `source_schema_change`, `no_data`, `parse_failure`, `model_unavailable`, `quality_failure`, or `redaction_failure`.

## 5. Real Browser Product And Design Review

- [x] Start local source backend and frontend:
  - [x] `go run ./cmd/server`
  - [x] `npm --prefix web run dev`
- [x] Verify the main product route set in a real browser:
  - [x] `/workbench`
  - [x] `/positions`
  - [x] `/settings`
  - [x] `/consultation`
  - [x] `/review`
  - [x] `/rules`
  - [x] `/audit`
  - [x] `/notifications`
  - [x] `/data-quality`
- [x] Verify a local user journey:
  - [x] Workbench shows current status and next manual action.
  - [x] Positions supports add/edit/import/offline transaction or equivalent seeded acceptance evidence.
  - [x] Settings/data refresh produces API and SQLite/readback evidence when provider access is available.
  - [x] Consultation produces explainable analysis or an honest safe degradation.
  - [x] Decision detail displays evidence, rules, assumptions, expected-return context, and manual confirmation state without nullable-field crashes.
  - [x] Review/rules/audit/notifications show governance traceability.
- [x] Verify product design at 390px, 768px, and 1280px:
  - [x] No incoherent overlap or clipped critical control.
  - [x] Loading/empty/error/degraded/success states are understandable.
  - [x] Manual confirmation and no-trading boundary remain clear.
  - [x] Evidence and audit links are discoverable from critical conclusions.

## 6. Final Acceptance Report

- [x] Create `docs/release/acceptance/2026-06-23-p100-local-source-final-acceptance.md`.
- [x] Record:
  - [x] Commit, version, date, operator, environment, config path, and LLM/provider availability.
  - [x] Each command result with pass/degraded/blocked/skipped status.
  - [x] Browser journey evidence and artifact paths.
  - [x] API/SQLite/readback/audit evidence summary.
  - [x] Product design rubric result.
  - [x] Explicit out-of-scope statement for Docker, install/upgrade/uninstall, GitHub Release, package refresh, and physical second-machine validation.
  - [x] Final conclusion: `local_source_release_acceptance_passed`, `local_source_release_acceptance_passed_with_documented_degradation`, or `local_source_release_acceptance_blocked`.

## 7. Archive

- [x] Update governance/progress materials with the P100 conclusion.
- [x] Run final validation:
  - [x] `openspec validate --all --strict`
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p93_code_reality_audit.py --check`
  - [x] `git diff --check`
- [x] Archive P100 after all blocking tasks pass or the final acceptance report records a blocked conclusion honestly.
