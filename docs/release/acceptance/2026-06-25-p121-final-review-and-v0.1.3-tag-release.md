# P121 Final Review And v0.1.3 Tag Release

Date: 2026-06-25  
Change: `p121-final-review-and-v0-1-3-tag-release`  
Source version: `v0.1.3`  
Status: `passed`

## Purpose

P121 is the fresh final release review for the current repository state after P114-P120 were archived. It exists because P93 is correctly stale after the later product/UI and scenario-acceptance changes; this release must not pretend that P93 was freshly re-run against the current tree.

## Scope

In scope:

- Validate current governance state and P114-P120 archive presence.
- Synchronize `VERSION`, `web/package.json`, and `web/package-lock.json` to `v0.1.3`.
- Record release notes and package/tag boundaries.
- Run OpenSpec, Go, frontend, P92, P121, whitespace, and local release package gates before tagging.

Out of scope:

- New investment runtime capability.
- Backend API, SQLite schema, Eino workflow, frontend feature expansion, Docker installation validation, upgrade validation, uninstall validation, or physical second-machine validation.

## Verification Results

| Gate | Result |
| --- | --- |
| `openspec validate p121-final-review-and-v0-1-3-tag-release --strict` | passed |
| `openspec validate --all --strict` | passed |
| `go test ./...` | passed |
| `go vet ./...` | passed |
| `npm --prefix web test -- --run` | passed |
| `npm --prefix web run build` | passed |
| `python3 scripts/p92_final_requirement_audit.py --check` | passed |
| `python3 scripts/p121_final_release_review.py --check` | passed |
| `git diff --check` | passed |
| `bash scripts/local-release-package.sh --release-label v0.1.3 --output-dir tmp/p121-release-package-final` | passed |
| `bash scripts/local-release-package.sh --verify tmp/p121-release-package-final/20260625T051439Z/investment-agent-v0.1.3.tar.gz --output-dir tmp/p121-release-package-final-verify` | passed |

## Package Identity

| Field | Value |
| --- | --- |
| Archive | `investment-agent-v0.1.3.tar.gz` |
| SHA256 | `ba08fb46606688239f67ea5534b8b038f389969ae1c8e8e66a7ff83e8a895bbc` |
| Manifest | `tmp/p121-release-package-final/20260625T051439Z/release-manifest.json` |

## Release Boundary

P121 may claim that P114-P120 scoped product/UI/real-use/control acceptance artifacts are archived and that the current tree passed P121 release gates after those changes. It may not claim fresh P93 pass after P114-P120. It may not claim broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, fresh real LLM/provider quality, future provider availability, prediction accuracy, investment returns, Docker installation validation, upgrade validation, uninstall validation, or physical second-machine validation.
