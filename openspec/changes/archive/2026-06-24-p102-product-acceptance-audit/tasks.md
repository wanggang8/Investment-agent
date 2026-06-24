# Tasks

## 1. Governance

- [x] Validate P102:
  - [x] `openspec validate p102-product-acceptance-audit --strict`

## 2. Machine Gates

- [x] Run real LLM smoke:
  - [x] `go run ./cmd/agent --task llm-smoke --symbol 510300`
- [x] Run focused/full gates as needed:
  - [x] `openspec validate --all --strict`
  - [x] `go test ./...`
  - [x] `npm --prefix web test -- --run`
  - [x] `npm --prefix web run build`
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p93_code_reality_audit.py --check`
  - [x] `git diff --check`

## 3. Product Journey Capture

- [x] Start local Go backend and Vite frontend using `configs/config.yaml`.
- [x] Capture desktop screenshots for key product surfaces.
- [x] Capture at least one mobile/responsive screenshot.
- [x] Exercise at least one real workflow interaction and record state/readback evidence.
- [x] Check console errors and browser-visible failures.

## 4. Product Acceptance Analysis

- [x] Assess strengths, UX risks, accessibility risks, and product boundary risks from current-run evidence.
- [x] Write `docs/release/acceptance/2026-06-24-p102-product-acceptance-audit.md`.
- [x] Save screenshot notes in the audit asset folder.

## 5. Archive

- [x] Update governance/progress status.
- [x] Archive P102 after validation.
