# Design

## Version Choice

Use `v0.1.0` as the initial release version.

Rationale:

- The project has completed full requirement/code-reality audits and deployment hardening, so `0.0.0` is no longer useful as a release-facing identifier.
- `0.1.0` is conservative: it marks the first local release line without implying a `1.0.0` public stability contract, physical second-machine validation, or remote release publication.

## Version Locations

- `VERSION` contains the canonical repository release label: `v0.1.0`.
- `web/package.json` and `web/package-lock.json` use npm-compatible `0.1.0`.
- Release documentation records the same version and explicitly preserves release-claim boundaries.

No runtime code path reads this version in P99; package-generation and Git tagging can be handled by later explicit release actions.
