# Design: P79 Real Use Data-Impact And Expected-Return Closure

## Evidence Model

P79 is an evidence layer over P78. It does not edit historical P75/P77/P78 matrices. Instead, it reads the P78 matrix and emits:

- `docs/release/acceptance/2026-06-21-p79-real-use-data-impact-and-expected-return-matrix.md`
- `docs/release/acceptance/2026-06-21-p79-real-use-data-impact-and-expected-return-closure.md`
- `docs/release/ui-audit-assets/2026-06-21-p79/real-use-data-impact-summary.json`

Rows can become `real_pass` only when the checker can point to fresh P79 evidence.

## Upgrade Rules

Portfolio / confirmation rows may be upgraded when all are true:

- P79 reruns the P72 real-user fund scenario in a fresh artifact directory.
- P79 reruns the accepted-local non-`510300` UI journey in a fresh artifact directory.
- SQLite readback confirms portfolio snapshots, positions, operation confirmations, position transactions, decision records, evidence refs, and audit events.
- Field-level readback confirms position symbol, name, quantity, cost price, buy reason, asset tag, confirmation quantity/price, transaction quantity/price, and before/after transaction state.
- The latest local account values match independently computed expectations.
- Forbidden broker/order/external-push tables and automatic confirmation claims are absent.
- The row text is directly covered by the action/readback evidence.

Expected-return rows remain `partial` unless fresh P79 evidence proves available or degraded UI output for the exact required fields. P78 deliberately did not upgrade broad probability/scenario rows from low-sample evidence; P79 keeps that guard.

Broad monthly attribution rows remain non-`real_pass` unless the evidence proves monthly attribution and discipline-audit readback directly. P79 daily/local account snapshot readback is not enough for `REQ-14-005`.

## Expected-Return Quality Fallback

The real P72 UI rerun exposed a quality-gate instability: the expected-return LLM material can be parseable but fail the safety quality evaluator. ExpectedReturnNode now treats only analyst `quality_failure` differently from ordinary service unavailability:

- discard the failed LLM output;
- preserve deterministic local expected-return scenarios and sell-evaluation boundaries;
- write an `expected_return` analyst material generated from deterministic local scenarios;
- record metadata with `model=deterministic-local`, `parse_status=parsed`, `quality_status=passed`, and `fallback_reason=llm_quality_failure`;
- keep timeout/model-unavailable/authentication errors degraded as before.

This is a safety fallback for analysis material only. It does not upgrade expected-return probability/scenario requirement rows, does not create a trade instruction, and does not change final verdict ownership from the rule engine.

## Safety

The checker scans P79 acceptance materials and text artifacts for private absolute paths and overbroad release claims. It fails if materials claim:

- full original-requirement pass while any full-release-required row remains non-`real_pass`;
- P79 evidence inside the existing P76 distribution archive;
- broker/order/automatic trading/external-push capabilities;
- future return or provider-availability promises.
