# Design: P68 Post-P67 Release Readiness Governance

## Context

P66 and P67 intentionally split two facts:

- P66 strict current-data gate remains the source of truth for clean current-data health.
- P67 resolution records only explain how a release owner handled a blocked or waiver-required gate.

The P68 design is therefore a governance reconciliation, not a feature build. It should read the current release materials as a user or release owner would read them, identify any wording that could overstate readiness, and leave a single decision artifact for the next operator.

## Decision Model

P68 writes a release readiness decision with these possible outcomes:

| Outcome | Meaning |
| --- | --- |
| `release_ready_limited_current_data_scope` | Existing release evidence is sufficient, P67 resolution is active, and materials explicitly exclude current local data health from clean claims. |
| `release_ready_requires_package_refresh` | Product/runtime acceptance remains release-ready, but final distribution artifacts should be regenerated after P65-P67/P68 commits from a clean tree before external handoff. |
| `release_blocked` | A blocking gate has no pass, waiver, or scope exclusion, or release materials contain unresolved safety/redaction issues. |

P68 may choose more than one non-blocking qualifier, for example `release_ready_limited_current_data_scope` with a next-stage recommendation to refresh packaging. It must never describe `scope_exclusion` as a data-quality pass.

## Evidence Inputs

- `docs/release/release-candidate-2026-06-18.md`
- `docs/release/release-handoff-2026-06-18.md`
- `docs/release/README.md`
- `docs/release/acceptance-repeatability.md`
- `docs/release/acceptance/2026-06-18-p66-current-data-policy.md`
- `docs/release/acceptance/2026-06-18-p67-current-data-resolution.md`
- Current command evidence:

```bash
openspec validate p68-post-p67-release-readiness-governance --strict
openspec validate --all --strict
git diff --check
go run ./cmd/agent --task data-source-quality-resolution-check --symbol 000300
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate
```

The P66 command is expected to exit non-zero while reporting `policy=blocked` / `gate=block`; this is evidence, not an implementation failure.

## Document Updates

P68 should create:

- `docs/release/acceptance/2026-06-18-p68-release-readiness-governance.md`

P68 should update if needed:

- `docs/release/release-candidate-2026-06-18.md`
- `docs/release/release-handoff-2026-06-18.md`
- `docs/release/README.md`
- `docs/release/acceptance-repeatability.md`
- `docs/development-plan.md`
- `docs/README.md`
- `docs/GOVERNANCE.md`
- `AGENTS.md`
- `openspec/project.md`
- `openspec/PROGRESS.md`

## Review Rules

The plan review and execution review must check for these failure modes:

- Release materials claim or imply current local data is clean while P66 is blocked.
- P67 `resolved_with_scope_exclusion` is treated as a policy pass.
- P64/P65 package evidence is presented as final after later commits without either regeneration or a clear P69 recommendation.
- Any new wording introduces prohibited capabilities or future provider/return promises.

## Verification

Because P68 is docs/governance-only, the minimum verification is OpenSpec validation, `git diff --check`, P66/P67 command evidence, and targeted safety scans. If P68 changes scripts or runtime code unexpectedly, full Go/frontend/E2E verification becomes mandatory before archive.
