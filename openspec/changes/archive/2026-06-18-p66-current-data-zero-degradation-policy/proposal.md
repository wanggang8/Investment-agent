# Proposal: P66 Current Data Zero-Degradation Policy

## Why

P63 and P65 kept the project release-ready, but both records still carry the same caveat: G5 current data-source quality may be `degraded` while fixture regression passes. P52 allowed that as long as it was classified, which was appropriate for release execution, but it leaves the product without a clear policy for when current local data quality should block a future release, require an explicit waiver, or pass cleanly.

P66 closes that ambiguity. It turns current data-source quality from a descriptive regression result into a policy verdict that can be reused by CLI, API, UI, and release materials.

This phase must remain local, read-only, and policy-only. It must not add new data sources, refresh current data automatically, call public providers or LLM providers, repair data, or weaken the existing safety boundaries.

## What Changes

- Extend the existing data-source-quality regression response with a `policy` object.
- Classify current data-source quality into:
  - `passed`: no degraded/failed current source-health cases.
  - `waiver_required`: only recognized, non-core degraded categories are present and must be explicitly documented before release claims.
  - `blocked`: failed/unrecognized cases, missing source-health facts, or core source categories degraded.
- Add machine-readable release gate fields, counts, reasons, and manual next actions.
- Expose the policy through the existing API and CLI audit output.
- Add a strict local gate command or script for release acceptance that exits non-zero when policy is `blocked`.
- Surface the policy on the existing Data Quality page without adding new product flows.
- Update release handoff and acceptance repeatability docs so future release-ready claims cannot silently ignore current data degradation.

## In Scope

- Existing P48 data-source-quality service, DTO, API, CLI task, and tests.
- Existing Data Quality page model and display.
- A local strict gate entrypoint for current data quality policy.
- P66 acceptance evidence under `docs/release/acceptance/`.
- Updates to governance/progress documents and release handoff materials.
- OpenSpec delta for `data-source-quality-regression`.

## Out Of Scope

- Adding or enabling new public data sources.
- Calling real public providers, LLM providers, login-gated sources, paid sources, authorization-gated sources, Level2 feeds, or high-frequency feeds.
- Automatically refreshing market data, repairing source health, rebuilding indexes, applying rules, confirming actions, writing final decisions, creating broker orders, or changing account/position facts.
- Changing SQLite schema, adding Eino workflows, or adding new backend business APIs beyond the existing data-source-quality response surface.
- Claiming current local data is healthy unless the policy verdict actually passes.

## Selected Design

Three approaches were considered:

1. **Doc-only release rule.** Low risk, but leaves CLI/API/UI unable to enforce the policy.
2. **Hard fail on any current degradation.** Strongest zero-degradation stance, but too brittle for known optional sources and makes every classified external gap a binary failure.
3. **Policy verdict with strict release gate.** Recommended. The product keeps local read-only diagnostics, but future release materials must treat `blocked` as blocking and `waiver_required` as an explicit documented exception.

P66 uses approach 3.

## Validation

- Plan review by sub agent before execution.
- `openspec validate p66-current-data-zero-degradation-policy --strict`
- `openspec validate --all --strict`
- `git diff --check`
- Focused Go tests for data-source-quality service, handler, and CLI.
- Frontend unit tests for Data Quality policy display.
- Strict current-data policy gate against fixture/current test data.
- `go test ./...`
- `npm --prefix web test`
- `npm --prefix web run build`
- `bash scripts/e2e-smoke.sh`
- Execution review by sub agent before archive.
- Submit review by sub agent before commit.
