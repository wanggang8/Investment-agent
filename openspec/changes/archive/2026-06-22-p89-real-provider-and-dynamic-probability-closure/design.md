# P89 Design

## Row Ownership

P89 starts from `docs/release/acceptance/2026-06-22-p88-remaining-full-release-blockers-matrix.md` and owns only rows where:

- `full_release_required == True`
- `p88_status != real_pass`

The expected set is:

```text
REQ-04-016
REQ-05-003
REQ-05-004
REQ-05-005
REQ-08-004
REQ-08-023
REQ-09-004
REQ-09-023
REQ-09-024
REQ-09-025
```

Any drift fails the inventory gate.

## Structured Provider Strategy

P89 treats provider verification and provider runtime proof as separate gates:

1. **Preverification** records candidate authority, public access shape, stable request/page evidence, fields, update frequency, legal/access limits, rate assumptions, failure behavior, and SQLite target path.
2. **Runtime collection** fetches a low-frequency public source without login, payment, authorization, Level2, or high-frequency terms.
3. **SQLite readback** proves stored values exist in product storage and are surfaced through API/UI readiness where applicable.
4. **Closure classification** upgrades a row only when the relevant field set has all three gates. Parser tests alone are not enough.

If any provider cannot be verified in this environment, P89 must record `provider_status=blocked` or `provider_status=not_verified` and keep the affected row `partial`.

## Dynamic Expected-Return Strategy

P89 adds product-visible dynamic monitoring evidence:

- A baseline expected-return run stores probabilities and assumptions.
- A changed valuation/fundamental/market-state input reruns the report and lowers the affected probabilities.
- A two-month assumption miss produces a downshift warning.
- A one-month pessimistic actual path produces a manual probability-adjustment suggestion.
- Extreme fear locks active trading advice and displays historical similar-scenario context.

The acceptance must use real UI actions against the local backend and SQLite readback, not only Go tests.

## Safety

P89 remains read-only toward external markets and local-only toward user data. It must not add or claim:

- broker/order/execution tables or UI affordances
- automatic confirmation
- external push
- automatic rule application
- automatic repair/migration/recovery
- future return accuracy
- future provider availability

## Evidence Outputs

P89 should produce:

- `docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability/p89-inventory.json`
- `docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability/p89-source-preverification.json`
- `docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability/p89-acceptance-summary.json`
- `docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability/db-readback-check.log`
- `docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability/frontend-build.log`
- `docs/release/ui-audit-assets/2026-06-22-p89-real-provider-dynamic-probability/final-validation.log`
- `docs/release/acceptance/2026-06-22-p89-real-provider-dynamic-probability-closure.md`
- `docs/release/acceptance/2026-06-22-p89-real-provider-dynamic-probability-matrix.md`
