# P69 Clean Tree Package Refresh Acceptance

> Date: 2026-06-18
> Change: `p69-clean-tree-package-refresh`
> Status: `passed`
> Release status remains: `release_ready_limited_current_data_scope`

## Scope

P69 regenerates final local package evidence from a clean detached worktree at the committed P68 source commit. It verifies the archive and repeats acceptance from the extracted package workspace.

P69 does not change runtime behavior, publish remotely, create a Git tag, sign an installer, run migrations, upgrade an installation, restore data, repair files, call public providers, call LLM providers, trade, push externally, confirm actions, or apply rules.

## Package Identity

| Field | Value |
| --- | --- |
| Release label | `p69-clean-tree` |
| Source commit | `cc0a64781e199a7745432b63bce26de4402042b5` |
| Source status | `clean` |
| Archive | `tmp/p69-final-release/20260618T084011Z/investment-agent-p69-clean-tree.tar.gz` |
| Sidecar manifest | `tmp/p69-final-release/20260618T084011Z/release-manifest.json` |
| Package SHA-256 | `d764ce5770289b6c174c919923ace181354165f8c8b114cfff444701cf158faa` |
| Archive entry count | 1323 |
| Verify summary | `tmp/p69-final-release/20260618T084023Z-verify/verify-summary.json` |
| Verify status | `passed` |
| Repeat summary | `tmp/p69-final-repeat/20260618T084028Z/repeat-summary.json` |
| Repeat status | `passed` |
| `skip_install` | `false` |
| `skip_e2e` | `false` |

The archive includes committed source through P68. It does not claim to include the P69 or later acceptance records unless a later package refresh is performed.

## Commands Run

```bash
git worktree add --detach tmp/p69-clean-tree-source cc0a64781e199a7745432b63bce26de4402042b5
cd tmp/p69-clean-tree-source
npm --prefix web ci
git status --short
bash scripts/local-release-package.sh --release-label p69-clean-tree --output-dir tmp/p69-release
bash scripts/local-release-package.sh --verify tmp/p69-release/20260618T084011Z/investment-agent-p69-clean-tree.tar.gz --output-dir tmp/p69-release
bash scripts/local-release-repeat-acceptance.sh --archive tmp/p69-release/20260618T084011Z/investment-agent-p69-clean-tree.tar.gz --output-dir tmp/p69-repeat
```

The package and repeat artifacts were copied into the main project `tmp/p69-final-release/` and `tmp/p69-final-repeat/` before the temporary detached worktree was removed.

`npm --prefix web ci` reported one high-severity dependency audit item. P69 did not run `npm audit fix` because that would mutate dependencies and is outside the clean package refresh scope; the package repeat acceptance command matrix passed.

## Repeat Command Matrix

| Command | Result | Duration |
| --- | --- | --- |
| `openspec validate --all --strict` | passed | 2s |
| `go test ./...` | passed | 11s |
| `npm --prefix web ci` | passed | 1s |
| `npm --prefix web test` | passed | 6s |
| `npm --prefix web run build` | passed | 3s |
| `env E2E_SERVER_PORT=18165 E2E_WEB_PORT=14265 bash scripts/e2e-smoke.sh` | passed | 29s |

## Safety And Redaction Checks

Package verification passed with no errors or warnings. The archive excludes `.git/`, `.cursor/`, project `tmp/`, `cmd/agent/tmp/`, UI audit assets, local config, generated frontend output, dependency folders, Playwright reports, test results, SQLite databases, logs, traces, raw payloads, complete prompt payloads, complete API keys, and private local paths.

The repeat flow ran from the extracted package workspace, not from the active repository checkout.

## Caveats

| Item | Position |
| --- | --- |
| P69 documentation inclusion | The package source commit is P68 `cc0a64781e199a7745432b63bce26de4402042b5`; it does not include this P69 acceptance document. |
| Current data policy | P66 remains `policy=blocked` / `gate=block`; release status remains `release_ready_limited_current_data_scope`. |
| Physical second machine | Not performed in P69. This record covers a local isolated repeat from a clean package archive. |
| Real providers | No public-source or LLM provider calls are made by the package or repeat scripts. |

## Not Claimed

P69 does not claim current local data is clean, P66 policy passed, P69 documentation is inside the generated archive, physical second-machine execution, remote publishing, Git tag creation, installer signing, automatic upgrade, automatic migration, automatic restore, automatic repair, real database overwrite, broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, future public-source availability, future model-provider availability, login-gated sources, paid sources, authorization-gated sources, Level2 data, high-frequency data, or investment returns.
