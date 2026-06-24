# P101 Unified Local Config Path And LLM Compatibility Acceptance

Date: 2026-06-24

Change: `p101-unify-local-config-path`

Conclusion: `local_config_unified_and_openai_compatible_llm_smoke_passed`

## Scope

P101 aligned current local-source acceptance scripts with the runtime default ignored config file `configs/config.yaml`, preserving script-specific override variables. During real-provider validation, P101 also fixed the OpenAI-compatible LLM request compatibility gap found by comparing this project with the sibling `ai-agent` project:

- The local LLM client continues to call `POST <base_url>/chat/completions`.
- The request body remains OpenAI Chat Completions compatible: `model` plus `messages`.
- The response parser still reads `choices[0].message.content`.
- P101 adds `Accept: application/json` and a stable `User-Agent`.
- P101 retries one transport timeout and keeps parser/quality-gate behavior unchanged.
- Default/example `deepseek.timeout_seconds` is now 60 seconds.

No API key is committed or printed. Docker, installation, upgrade, uninstall, GitHub Release, package refresh, and physical second-machine validation remain out of scope.

## Findings Closed

- Initial request with missing `/v1` produced `/chat/completions` and failed.
- After using `/v1`, a minimal non-SDK request still received gateway 403 / Cloudflare 1010.
- Adding SDK-style JSON accept plus user-agent headers returned HTTP 200 from the configured OpenAI-compatible gateway.
- The original 15-second timeout was too tight for the configured model gateway.
- A 60-second bounded timeout plus one timeout retry made `llm-smoke` pass through the normal CLI path.

## Evidence

Commands executed and passed:

- `openspec validate p101-unify-local-config-path --strict`
- Static config path check for P63/P71/P72/P75 script defaults
- `go test ./internal/infrastructure/config ./internal/infrastructure/llm/deepseek`
- `go run ./cmd/agent --validate-config`
- `go run ./cmd/agent --task llm-smoke --symbol 510300`
- `bash scripts/p71-real-product-acceptance.sh`
- `bash scripts/p72-real-user-fund-scenario-acceptance.sh`
- `bash scripts/p86-core-goal-knowledge-safety-final-acceptance.sh`
- `openspec validate --all --strict`
- `go test ./...`
- `npm --prefix web test -- --run`
- `npm --prefix web run build`
- `python3 scripts/p92_final_requirement_audit.py --check`
- `python3 scripts/p93_code_reality_audit.py --check`
- `git diff --check`

Latest local audit evidence includes successful `llm-smoke:symbol=510300:model=gpt-5.4-mini` rows with `llm_smoke:quality=passed:parse=parsed:no_auto_trading`.

## Boundaries

P101 does not add broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, paid/login/auth-only sources, Level2 data, high-frequency data, future provider availability guarantees, or investment return guarantees.
