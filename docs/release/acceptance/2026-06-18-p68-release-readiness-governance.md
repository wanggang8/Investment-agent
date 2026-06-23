# P68 Release Readiness Governance

> Date: 2026-06-18
> Change: `p68-post-p67-release-readiness-governance`
> Status: `release_ready_limited_current_data_scope`
> Package status: `release_ready_requires_package_refresh`

## Scope

P68 reconciles the P63 release-ready evidence, P64/P65 package evidence, P66 current-data policy gate, and P67 current-data resolution workflow. It is a governance and release-material refresh only. It does not change runtime behavior, refresh data, call providers, call LLMs, create package artifacts, publish remotely, create a Git tag, trade, push externally, confirm actions, apply rules, repair files, migrate data, or restore data.

## Evidence

| Check | Result | Evidence | Release impact |
| --- | --- | --- | --- |
| P67 resolution check | pass | `claim_state=resolved_with_scope_exclusion:policy=blocked:gate=block:fingerprint=e63f5ffed4a4307a4f791a4d68a9fc4a6f37bd34188dcacea1fd3d1f9ac4da48:resolution=scope_exclusion:clean_data_claim=false:no_auto_trading` | Supports a limited release claim that excludes current local data health from the clean claim. |
| P66 strict current-data gate | expected block | `policy=blocked:gate=block:cases=1:degraded=1:failed=0:no_auto_trading`, exit 1 | Blocks any clean current-data release claim. |
| P63 release evidence | pass with caveats | `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md` | Still supports product/runtime release readiness with documented caveats. |
| P65 package repeat evidence | pass for candidate archive | Source commit `ef2f55acfcd2ee5e96676a59014e3766282b876a`, `source_status=dirty` | Validates the P65 candidate archive workflow, but final distribution should be regenerated after P65-P68 from a clean tree. |

## Decision

The project remains release-ready for local product/runtime handoff under a limited current-data scope:

- P63/P65 acceptance evidence remains valid as historical execution evidence.
- P66 still blocks clean current local data claims.
- P67 records an active scope exclusion for symbol `000300`.
- P68 refreshes the release statement to `release_ready_limited_current_data_scope`.
- Final distribution packaging should be refreshed in P69 from a clean tree because P64/P65 package evidence predates later P66-P68 commits and was generated with `source_status=dirty`.

## Allowed Claims

- The local product/runtime release is ready under the documented safety boundaries.
- Current local data health is excluded from the clean release claim through the active P67 `scope_exclusion`.
- P64/P65 package workflows are validated for candidate artifacts and can be repeated.
- A P69 clean-tree package refresh is the recommended final distribution step.

## Not Claimed

- The current local data snapshot is not claimed as clean.
- The P66 policy is not passed.
- The P67 resolution is not a provider repair, data refresh, source-health fix, or policy pass.
- The P65 candidate archive is not claimed to include P66-P68 commits.
- No physical second-machine package repeat has been performed.
- No future public-source availability, model-provider availability, investment return, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic upgrade, automatic migration, automatic restore, real database overwrite, login-gated source, paid source, authorized source, Level2 data, or high-frequency data is claimed.

## Verification Commands

```bash
go run ./cmd/agent --task data-source-quality-resolution-check --symbol 000300
go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate
openspec validate p68-post-p67-release-readiness-governance --strict
openspec validate --all --strict
git diff --check
```

The strict P66 command is expected to exit non-zero while reporting `policy=blocked` / `gate=block`. That expected block is the basis for the limited current-data scope.
