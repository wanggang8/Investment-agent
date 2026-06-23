# P78 Requirements Real-Pass Batch Closure

## Why

P77 established the conservative `real_pass` upgrade gate and upgraded the first 17 atomic rows, but it left 313 full-release-required rows non-`real_pass`. Continuing toward true product acceptance now requires a batch closure workflow that turns those gaps into auditable repair batches instead of broad release language.

The immediate product risk is that positive functionality rows can remain described only by scoped evidence, inherited screenshots, temporary journey artifacts, or generic "needs evidence" notes. P78 starts the next closure batch by targeting rows where direct evidence can be made concrete without expanding unsupported claims: expected-return degradation/disclaimer behavior from `REQ-09`, backed by deterministic Go tests and real UI SQLite readback from a non-`510300` journey.

## What Changes

- Create a P78 evidence layer derived from the P77 matrix without rewriting P75 or P77 history.
- Classify every remaining full-release-required non-`real_pass` row into remediation groups and execution batches.
- Run fresh expected-return deterministic tests and a fresh accepted-local non-`510300` real UI journey.
- Upgrade only the first batch of expected-return degradation/disclaimer rows whose applicable dimensions are directly evidenced by:
  - implementation tests,
  - real browser UI operation,
  - SQLite decision readback,
  - no precise probability under insufficient samples,
  - non-trading disclaimer and safety boundaries.
- Generate P78 acceptance materials and a machine-readable summary.

## Out Of Scope

- No broker interface, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, or real database overwrite.
- No new promise of future public-provider, LLM-provider, or market-data availability.
- No investment return promise, price prediction, or future market-direction claim.
- No rewrite of P75/P77 history and no claim that scoped, partial, fixture-only, mock/stub-only, screenshot-only, route-smoke-only, waiver-only, scope-exclusion-only, temporary-DB-only, or incompatible single-symbol-only evidence is full original-requirement pass.
- No P76 package refresh; a separate package change is required if distribution archives must include P78 materials.

## Acceptance

P78 is accepted only if:

- `python3 scripts/p78_requirements_real_pass_batch_closure.py --check` passes.
- Fresh P78 evidence commands cited by upgraded rows pass and their logs/artifacts exist.
- `openspec validate p78-requirements-real-pass-batch-closure --strict` passes.
- `openspec validate --all --strict` passes.
- `git diff --check` passes.
- A read-only subagent review finds no Critical or Important issues, or all such findings are fixed before archive.
