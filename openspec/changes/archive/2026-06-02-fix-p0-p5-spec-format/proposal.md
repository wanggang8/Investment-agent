# fix-p0-p5-spec-format

## Why

`openspec validate --specs --strict` failed on `p0-p5-capabilities` because the spec lacked the required `## Purpose` section. This blocks clean OpenSpec validation after P9 archive.

## What Changes

- Add a `## Purpose` section to `openspec/specs/p0-p5-capabilities/spec.md`.
- Keep existing P0-P5 requirements and scenarios unchanged.
- Verify the individual spec and all specs pass strict validation.

## In Scope

- OpenSpec spec formatting only.

## Out of Scope

- Product behavior changes.
- L1 contract changes under `docs/`.
- Backend, frontend, database, or API implementation changes.
