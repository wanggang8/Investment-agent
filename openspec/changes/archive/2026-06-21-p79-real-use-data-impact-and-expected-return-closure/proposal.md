# Proposal: P79 Real Use Data-Impact And Expected-Return Closure

## Summary

P79 continues the P78 real-pass batch closure by executing the next real-use batch against the P78 matrix. It focuses on two user-critical surfaces:

- portfolio / confirmation / local-account data impact, where a real browser action must be tied to SQLite readback, audit evidence, and prohibited-table checks;
- expected-return UI coverage, where P79 records the remaining gaps honestly and only upgrades rows if fresh UI/readback evidence covers the exact probability, sample, scenario, and disclaimer fields;
- expected-return LLM quality hardening, where an unsafe/low-quality expected-return LLM material is discarded and replaced with deterministic local scenario material instead of turning a fully evidenced real UI flow into an unstable degraded consultation.

## Motivation

P78 left 310 full-release-required rows non-`real_pass`. The largest near-term risk is that a real user performs a portfolio, confirmation, or decision-detail action and the UI appears successful while SQLite, audit, or derived readbacks do not match. P79 turns the existing scoped P72/P75 evidence into a stricter batch gate and only promotes rows that are directly covered by fresh real UI execution.

## Scope

In scope:

- Generate a P79 matrix from the P78 matrix without rewriting P75/P77/P78 history.
- Rerun fresh real UI portfolio/data-impact acceptance under P79 artifacts.
- Rerun fresh non-`510300` real UI journey under P79 artifacts to keep dynamic-symbol coverage attached.
- Validate SQLite readback for portfolio snapshots, positions, operation confirmations, position transactions, decision records, evidence refs, and audit events.
- Harden ExpectedReturnNode so LLM `quality_failure` for expected-return material uses safe deterministic local scenario material with explicit fallback metadata, while ordinary analyst unavailability still degrades.
- Validate prohibited broker/order/external-push/auto-confirmation artifacts remain absent.
- Update release materials and progress docs with the exact upgraded count and remaining gap count.

Out of scope:

- Full original-requirement pass.
- P76 package refresh or any claim that P76 includes P79 evidence.
- Broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, automatic migration, automatic restore, real database overwrite, provider availability promises, or investment return promises.
- Upgrading expected-return probability/scenario rows unless P79 produces direct UI/readback evidence for the exact fields.
- Claiming monthly attribution rows as `real_pass` from daily/local account snapshot evidence alone.

## Acceptance

P79 passes only if:

- `scripts/p79_real_use_data_impact_and_expected_return_closure.py --check` passes.
- Fresh P79 real UI evidence artifacts exist and pass their SQLite/readback checks.
- ExpectedReturnNode quality-failure fallback is covered by unit/integration tests and P72 real UI rerun.
- `openspec validate p79-real-use-data-impact-and-expected-return-closure --strict` passes.
- `openspec validate --all --strict` passes.
- `git diff --check` passes.
- A read-only subagent review finds no Critical or Important issue before archive.
