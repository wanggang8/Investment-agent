# P121 Final Review And v0.1.3 Tag Release

## Why

P114-P120 have completed the post-redesign product/UI acceptance sweep and have been archived. The repository now needs one bounded release-governance pass that rechecks the current tree, records the P93 stale boundary honestly, synchronizes version metadata, writes release-facing notes, and publishes a new tag only if the fresh gates pass.

## What Changes

- Add a P121 final release review record for the current source tree after P114-P120.
- Bump current source version metadata from `v0.1.2` to `v0.1.3`.
- Update release-facing materials with clear scope, verification, and safety boundaries.
- Sanitize local absolute paths while copying text files into the release package so historical evidence can stay intact in source without leaking workstation paths in the distributed archive.
- Archive P121, commit the reviewed release state, create annotated tag `v0.1.3`, and push the release commit/tag.

## Scope

In scope:

- OpenSpec/release governance for `v0.1.3`.
- Fresh validation gates against the current source tree.
- Release notes and tag message quality.
- Local release package smoke/verify if the source tree passes earlier gates.
- Release package text redaction for local filesystem paths.

Out of scope:

- New investment runtime capability.
- Backend API, SQLite schema, Eino workflow, or frontend feature expansion beyond release metadata.
- Docker/installer/upgrade/physical second-machine validation unless separately run.
- Broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, future provider availability, prediction accuracy, or investment returns.
