# P70 Final Release Decision And Risk Closure

## Why

P69 produced clean-tree package evidence from the committed P68 source and repeated package acceptance successfully. The project now has no required next phase in `openspec/PROGRESS.md`, but the release story is split across P63-P69 records: product/UI acceptance, current-data policy, scope exclusion, package refresh, and handoff boundaries.

P70 creates the final release decision and risk-closure record for this milestone. It should make one clear statement about whether the project is complete for the limited local release scope, what is explicitly not claimed, and which future stages are optional rather than required.

## What Changes

- Create a final P70 release decision and risk-closure acceptance record that reconciles P63-P69 evidence.
- Refresh release handoff, release README, repeatability, governance, and progress wording so the current state is not ambiguous.
- Confirm whether any mandatory next phase remains, and record optional future work separately.
- Run verification and targeted safety scans before archive.

## In Scope

- Release governance documentation, acceptance evidence, OpenSpec delta, progress/governance updates, and final decision wording.
- Verification commands needed for a documentation/governance closeout: OpenSpec validation, `git diff --check`, release wording scans, forbidden capability scans, and targeted command evidence if the final decision references current-data or package state.
- A final milestone status vocabulary: `release_ready_limited_current_data_scope`, `release_blocked`, or `needs_follow_up_before_handoff`.

## Out of Scope

- No runtime feature work, SQLite schema changes, HTTP API changes, Eino workflow changes, frontend page/component changes, provider calls, source refresh, LLM calls, package publication, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic repair, broker interface, trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, return promise, login-gated source, paid source, authorization-gated source, Level2 data, or high-frequency source.
- No claim that P67 scope exclusion is a P66 policy pass.
- No claim that current local data is clean or healthy while the P66 strict gate remains blocked.
- No claim that the P69 archive includes P69 or P70 documents unless a later package refresh is performed.

## Impact

- Docs/OpenSpec only unless review finds a narrow release-tooling inconsistency that blocks the final decision.
- The expected outcome is a clear final milestone closeout: no required next phase for the limited local release scope, with optional future stages documented separately.
