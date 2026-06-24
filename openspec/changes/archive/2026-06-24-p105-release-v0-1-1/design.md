# Design

## Versioning

Use `v0.1.1` as a patch release. This is conservative because the delta after `v0.1.0` is primarily:

- Local config and OpenAI-compatible LLM request hardening.
- Product acceptance evidence and UX linkage fixes.
- Repeatable product operation-linkage acceptance runner.

It does not introduce new investment runtime capabilities, broker connectivity, automatic trading, or a `1.0.0` stability claim.

## Metadata

- Root `VERSION` contains `v0.1.1`.
- `web/package.json` and `web/package-lock.json` use npm-compatible `0.1.1`.
- Release docs distinguish current version `v0.1.1` from the historical initial marker `v0.1.0`.

## Tag Boundary

The tag `v0.1.1` should be created only after P105 validation passes and P105 is archived. Pushing the tag may trigger the GitHub release workflow, but P105 does not itself claim that a remote GitHub Release successfully completed unless the push/workflow result is separately verified.
