# Proposal: P77 Requirements Real Pass Upgrade Gate

## Why

P75 completed the original-requirement traceability audit, but its conclusion was intentionally scoped: 341 atomic requirement rows contained 0 `real_pass` rows, and many rows remained `partial`, `scoped_pass`, or `deterministic_local_evidence`. The user now wants a path toward true product acceptance without relying on mocks, scope exclusions, or broad release wording.

P77 establishes a stricter, repeatable upgrade gate for moving P75 rows toward `real_pass`. It must preserve historical P75 evidence while creating a new P77 evidence layer that makes every upgrade decision auditable.

## What Changes

- Add a P77 real-pass upgrade rule set for atomic requirement rows.
- Generate a P77 upgrade matrix derived from the P75 matrix instead of rewriting P75 historical artifacts.
- Upgrade only rows whose current evidence satisfies the P77 gate.
- Keep all remaining rows explicitly categorized with residual gaps and next remediation direction.
- Refresh release/acceptance materials so they describe the new P77 evidence layer without expanding unsupported claims.

## Scope

In scope:

- Release-governance evidence, acceptance artifacts, and automation for requirement-row status upgrades.
- Fresh reruns of relevant real UI / SQLite / safety evidence where those runs are used as P77 upgrade evidence.
- First-batch upgrade review focused on safety boundaries, deterministic local prohibitions, accepted-local dynamic symbol evidence, and user-action/data-impact evidence that can be revalidated now.

Out of scope:

- Broker interfaces, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, future public-source availability promises, paid/login/authorized/Level2/high-frequency data-source expansion, or investment return promises.
- Rewriting P75 historical matrix rows to look complete.
- Claiming `release_ready_full_requirements_traceable` unless every full-release-required atomic row satisfies the new gate.

## Success Criteria

- P77 OpenSpec validation passes.
- A P77 matrix and acceptance record are generated from the P75 matrix.
- Every upgraded row has a concrete evidence basis and no forbidden overclaim.
- Remaining gaps are still visible row-by-row.
- Final P77 conclusion is honest about whether the product reached full original-requirement pass or remains scoped.
