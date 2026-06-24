# Tasks

## 1. Governance

- [x] Confirm no unrelated active change exists.
- [x] Validate P101:
  - [x] `openspec validate p101-unify-local-config-path --strict`

## 2. Red Check

- [x] Prove current scripts are still split-brain:
  - [x] `! rg -n 'LOCAL_CONFIG="\\$\\{P(63|71|72|75)_LOCAL_CONFIG:-\\$ROOT_DIR/configs/config\\.local\\.yaml\\}"' scripts`

## 3. Implementation

- [x] Update script defaults:
  - [x] `scripts/p63-full-ui-regression.sh`
  - [x] `scripts/p71-real-product-acceptance.sh`
  - [x] `scripts/p72-real-user-fund-scenario-acceptance.sh`
  - [x] `scripts/p75-non-510300-real-ui-journey.sh`
- [x] Align OpenAI-compatible LLM requests:
  - [x] Send `Accept: application/json`.
  - [x] Send a stable `User-Agent`.
  - [x] Retry one transport timeout.
  - [x] Keep endpoint path, body schema, parser, and analyst quality gate unchanged.
- [x] Raise default/example LLM timeout to 60 seconds.
- [x] Update current docs/acceptance notes that refer to current local configuration behavior.
- [x] Keep historical OpenSpec archive records unchanged unless they are not historical facts.

## 4. Validation

- [x] Static path check passes:
  - [x] `! rg -n 'LOCAL_CONFIG="\\$\\{P(63|71|72|75)_LOCAL_CONFIG:-\\$ROOT_DIR/configs/config\\.local\\.yaml\\}"' scripts`
- [x] Config validation:
  - [x] `go run ./cmd/agent --validate-config`
- [x] Focused LLM client tests:
  - [x] `go test ./internal/infrastructure/llm/deepseek`
- [x] Real LLM smoke:
  - [x] `go run ./cmd/agent --task llm-smoke --symbol 510300`
- [x] Real LLM acceptance reruns:
  - [x] `bash scripts/p71-real-product-acceptance.sh`
  - [x] `bash scripts/p72-real-user-fund-scenario-acceptance.sh`
  - [x] `bash scripts/p86-core-goal-knowledge-safety-final-acceptance.sh`
- [x] Core gates:
  - [x] `openspec validate --all --strict`
  - [x] `go test ./...`
  - [x] `npm --prefix web test -- --run`
  - [x] `npm --prefix web run build`
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p93_code_reality_audit.py --check`
  - [x] `git diff --check`

## 5. Acceptance And Archive

- [x] Add `docs/release/acceptance/2026-06-24-p101-unify-local-config-path.md`.
- [x] Update governance/progress materials with P101.
- [x] Archive P101 after validation.
