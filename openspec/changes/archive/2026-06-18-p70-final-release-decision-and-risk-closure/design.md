# Design: P70 Final Release Decision And Risk Closure

## Context

P63 established product/runtime release readiness after full UI regression. P66 made current data-source quality a strict gate, P67 added a local resolution workflow, P68 narrowed the release claim to a limited current-data scope, and P69 regenerated clean-tree package evidence through the P68 source commit.

The remaining risk is not a missing feature. It is ambiguity: a future operator could confuse the P67 scope exclusion with clean current-data health, or confuse the P69 clean package with a package that includes P69/P70 documentation.

## Decision Model

P70 writes one final milestone decision with these possible outcomes:

| Outcome | Meaning |
| --- | --- |
| `release_ready_limited_current_data_scope` | P63-P69 evidence supports local handoff, current local data clean health is excluded, P69 package evidence is accepted through P68 source, and no mandatory next phase remains. |
| `needs_follow_up_before_handoff` | The product/runtime evidence is not blocked, but a required handoff artifact or wording correction is missing. |
| `release_blocked` | A blocking gate lacks a valid pass, waiver, or scope exclusion, or release materials contain unresolved safety/redaction issues. |

P70 must keep optional work separate from required release work. Optional examples include a physical second-machine repeat, a true P66 current-data pass, a package refresh that includes P69/P70 docs, or VecLite acceptance hardening.

## Evidence Inputs

- `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`
- `docs/release/acceptance/2026-06-18-p65-cross-machine-repeat.md`
- `docs/release/acceptance/2026-06-18-p66-current-data-policy.md`
- `docs/release/acceptance/2026-06-18-p67-current-data-resolution.md`
- `docs/release/acceptance/2026-06-18-p68-release-readiness-governance.md`
- `docs/release/acceptance/2026-06-18-p69-clean-tree-package-refresh.md`
- `docs/release/release-candidate-2026-06-18.md`
- `docs/release/release-handoff-2026-06-18.md`
- `docs/release/release-packaging-2026-06-18.md`
- `docs/release/README.md`
- `docs/release/acceptance-repeatability.md`

## Document Updates

P70 should create:

- `docs/release/acceptance/2026-06-18-p70-final-release-decision.md`

P70 should update:

- `docs/release/release-handoff-2026-06-18.md`
- `docs/release/README.md`
- `docs/release/acceptance-repeatability.md`
- `docs/development-plan.md`
- `docs/README.md`
- `docs/GOVERNANCE.md`
- `AGENTS.md`
- `openspec/project.md`
- `openspec/PROGRESS.md`

P70 should update `docs/release/release-candidate-2026-06-18.md` only if it still contains stale P68-era next-stage wording that contradicts P69/P70.

## Review Rules

The plan review and execution review must check for these failure modes:

- Any wording claims current local data is clean, healthy, or policy-passed while P66 is blocked.
- Any wording treats P67 `resolved_with_scope_exclusion` as provider repair, source refresh, or P66 pass.
- Any wording implies the P69 archive includes P69/P70 docs.
- Optional future stages are described as mandatory blockers for the limited local release scope.
- Any new text introduces prohibited capabilities or future provider/return promises.

## Verification

Minimum verification:

```bash
openspec validate p70-final-release-decision-and-risk-closure --strict
openspec validate --all --strict
git diff --check
go run ./cmd/agent --task data-source-quality-resolution-check --symbol 000300
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate
```

The strict current-data command is expected to exit non-zero while reporting `policy=blocked` / `gate=block`. That is evidence for the limited release boundary, not a P70 implementation failure.

If P70 modifies runtime code, scripts, frontend components, or package behavior unexpectedly, add:

```bash
go test ./...
npm --prefix web test
npm --prefix web run build
bash scripts/e2e-smoke.sh
```
