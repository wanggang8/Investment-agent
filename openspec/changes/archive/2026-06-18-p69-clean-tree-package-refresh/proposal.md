# P69 Clean Tree Package Refresh

## Why

P68 concluded that the product/runtime release is `release_ready_limited_current_data_scope`, but the latest package evidence is stale: P64/P65 package artifacts were generated before P66-P68 and with `source_status=dirty`. A final handoff archive should be regenerated from a clean tree so the manifest identifies a clean source state and includes the P65-P68 governance materials.

Because P69 itself creates documentation and OpenSpec records, the clean package must be generated from a separate clean checkout/worktree at the committed P68 HEAD. The P69 acceptance record will then document that package as a clean P65-P68 distribution snapshot. P69 does not claim that the generated archive includes the later P69 acceptance record unless a follow-up package refresh is performed after P69 commit.

## What Changes

- Create a clean-tree package refresh plan and acceptance record.
- Generate a local release package from a clean detached worktree at the P68 commit.
- Verify the package archive and adjacent manifest.
- Run package repeat acceptance from the generated archive.
- Refresh release packaging, release handoff, release README, and governance/progress materials with the P69 package identity, checksum, source commit, `source_status=clean`, repeat summary, and remaining boundaries.

## In Scope

- OpenSpec change, release docs, package/repeat acceptance records, and governance/progress docs.
- Temporary local worktree/checkouts under project `tmp/`, generated packages under project `tmp/`, and sanitized command evidence.
- Reuse existing P64/P65 scripts without expanding runtime behavior.

## Out of Scope

- No runtime feature work, SQLite schema, HTTP API, Eino workflow, frontend UI, provider calls, LLM calls, source refresh, data repair, migrations, restores, remote publication, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic repair, broker interface, trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, return promise, login-gated source, paid source, authorization-gated source, Level2 data, or high-frequency source.
- No claim that P66 current-data gate passed; current release remains limited by P68 unless a later current-data policy pass is recorded.
- No claim that a separate physical second machine ran the package.

## Impact

- Docs/OpenSpec and local `tmp/` artifacts only.
- Expected committed outputs are sanitized release acceptance and handoff documents, not package archives, temporary worktrees, logs, databases, build output, or generated node modules.
