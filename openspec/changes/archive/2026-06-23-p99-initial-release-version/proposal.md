# P99 Initial Release Version

## Why

The repository is release-ready through the latest requirement and code-reality audits, but it still has no root release version file and the frontend package remains at `0.0.0`. This makes local package identity and release handoff less explicit than the rest of the release materials.

## What Changes

- Introduce `v0.1.0` as the initial repository release version.
- Add a root `VERSION` file as the human-readable release version marker.
- Sync the frontend package metadata to `0.1.0`.
- Update release materials and governance progress to record the initial version.

## Scope Boundaries

- Does not create a Git tag.
- Does not publish a GitHub Release.
- Does not refresh a final distribution package unless done by a later packaging step.
- Does not claim physical second-machine validation.
- Does not add broker connectivity, automatic trading, one-click trading, delegated orders, external push, automatic confirmation, automatic rule application, paid/login/auth/Level2/HFT sources, or return guarantees.
