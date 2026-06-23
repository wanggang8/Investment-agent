# Release Candidate: 2026-06-17

> Status: `release_ready`
> Acceptance run: `docs/release/acceptance/2026-06-17-p53-acceptance-run.md`
> Code-under-test commit: `5832477`
> Change: P53 `p53-acceptance-execution-and-release-candidate-materials`

## Basis

- P51 audit evidence: `docs/p19-p24-audit-evidence-pack.md`
- P52 acceptance matrix: `docs/project-acceptance-gate-matrix.md`
- P53 acceptance execution: `docs/release/acceptance/2026-06-17-p53-acceptance-run.md`

P52 defined the gates; P53 executed them. This release candidate status is based on the P53 execution record, not on the existence of the P52 matrix.

## Acceptance Summary

| Area | Result | Release impact |
| --- | --- | --- |
| Governance/OpenSpec | pass | does_not_block |
| Go tests | pass | does_not_block |
| Frontend tests/build | pass | does_not_block |
| Browser E2E smoke | pass | does_not_block |
| Local fixture/current smoke | degraded | does_not_block |
| Real public source opt-in | pass | does_not_block |
| Real LLM opt-in | pass | does_not_block |
| Local install/release upgrade | pass | does_not_block |
| Safety/redaction | pass | does_not_block |

## Release Notes

- Full Go test suite passed.
- Focused integration packages passed.
- Frontend Vitest and production build passed after retry.
- Browser E2E smoke passed after retry.
- Recovery smoke, retrieval quality smoke, fixture data-source regression, and current data-source regression executed.
- Real public evidence refresh executed against a temporary SQLite database and wrote evidence, RAG, verification, and audit rows.
- Real LLM smoke executed with model `gpt-5.4-mini`; parse and quality gates passed.
- Local install diagnostics and release upgrade checks passed after retry.
- Safety and redaction review found no committed complete key, raw payload, full prompt, private path, raw SQL dump, or new prohibited capability.

## Known Degradations

| Item | Status | Impact |
| --- | --- | --- |
| Current data-source quality regression | degraded | The current local DB had one degraded case and zero failed cases. Fixture regression passed, so this does not block release but limits claims about the current local data snapshot. |
| First frontend build/E2E/install attempts | retried | Initial failures were local process kills; exact command retries passed. This does not block release, but future release runs should preserve retry evidence. |
| Initial G6 temporary config | corrected | The first real public source run used an incomplete temporary config. The corrected temporary config passed. |

## Not Claimed

This release candidate does not claim:

- Future availability of public websites or model providers.
- Investment returns or deterministic market outcomes.
- Broker connectivity, automatic trading, one-click trading, order delegation, or external push.
- Automatic confirmation, automatic rule application, automatic repair, automatic migration, or automatic overwrite of real user databases.
- Login, paid, authorized, Level2, or high-frequency data source coverage.

## Decision

`release_ready`.

Proceed only with the documented safety boundaries above. If future release policy requires zero retries or no degraded current-data result, open a follow-up hardening phase before distribution.
