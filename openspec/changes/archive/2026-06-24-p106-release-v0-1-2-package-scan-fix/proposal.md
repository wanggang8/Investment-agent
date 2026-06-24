# P106 Release v0.1.2 Package Scan Fix

## Why

The pushed `v0.1.1` tag triggered CI and Release workflows, but both failed at the release package smoke/build step. The package scanner rejected `web/src/pages/DataQualityPage.tsx` because a redaction label used a JSON-like `prompt: "..."` shape that matched the prompt-payload forbidden-content rule.

Because `v0.1.1` was already pushed, the safer release path is a patch release `v0.1.2` rather than moving the published tag.

## What Changes

- Remove the package-scan false positive from the Data Quality page redaction labels.
- Bump current source version metadata from `v0.1.1` to `v0.1.2`.
- Add a P106 acceptance record including the local release package smoke/verify gate.
- Archive P106, commit, tag `v0.1.2`, and push.

## Scope

In scope:

- Release package scan compatibility.
- Patch release metadata and release-facing documentation.

Out of scope:

- New investment runtime capability.
- Moving or deleting the already pushed `v0.1.1` tag.
- Docker/installer/physical second-machine validation.
