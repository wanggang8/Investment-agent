# P68 Post-P67 Release Readiness Governance

## Why

P63 refreshed the release candidate to `release_ready`, P64/P65 added packaging and repeat acceptance, P66 made the current data-source quality caveat a strict policy gate, and P67 added an auditable local resolution workflow. The current local state is now explicit: P66 still returns `policy=blocked` / `gate=block`, while P67 records `resolved_with_scope_exclusion` for symbol `000300`.

The release materials now need one final governance pass so the project does not leave ambiguous wording behind. P68 decides whether the release-ready statement remains valid only as a limited release claim, whether release candidate or handoff materials must be refreshed, and whether packaging repeat evidence must be regenerated after P65-P67 commits.

## What Changes

- Create a P68 release readiness decision record that evaluates P63/P65 release evidence against P66/P67 current-data policy and resolution evidence.
- Refresh release candidate, handoff, release README, and repeatability wording if they still imply current local data is clean or if they point at stale package evidence.
- Decide and document whether P69 is needed for post-P67 packaging repeat from a clean tree.
- Update governance/progress documents with the final P68 conclusion and next-stage recommendation.

## In Scope

- Release governance documentation, acceptance decision material, OpenSpec delta, progress/governance updates, and verification commands.
- Running current command evidence needed for the decision: OpenSpec validation, diff checks, P66 strict gate, P67 resolution check, and lightweight test/build smoke if release materials are materially refreshed.
- A clear release status vocabulary: `release_ready_limited_current_data_scope`, `release_ready_requires_package_refresh`, or `release_blocked`.

## Out of Scope

- No runtime feature work, SQLite schema, HTTP API, Eino workflow, frontend page/component changes, provider calls, source refresh, LLM calls, package publication, Git tag creation, automatic upgrade, automatic migration, automatic repair, broker interface, trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, return promise, login-gated source, paid source, authorization-gated source, Level2 data, or high-frequency source.
- No claim that P67 turns P66 `policy=blocked` into `policy=passed`.
- No claim that current local data is clean or healthy unless the P66 strict gate passes.

## Impact

- Docs/OpenSpec only unless the review finds a release-material inconsistency that requires a narrowly scoped script/document update.
- Final project completion assessment should become clearer: either no runtime work remains, or one final P69 package-refresh stage is explicitly queued.
