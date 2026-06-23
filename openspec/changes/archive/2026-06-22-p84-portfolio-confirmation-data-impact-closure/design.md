# P84 Design

## Evidence Strategy

P84 should use a temporary local database and real UI flows to prove local portfolio causality:

- Account or holding setup affects portfolio views and daily/workbench prerequisites.
- Offline transaction or manual confirmation affects position/readback where the current product supports it.
- Review/decision-loop/audit surfaces link to the local fact.
- Derived values are independently computed and compared when deterministic.
- Unsafe broker/order affordances remain absent.

## Real-Pass Rule

A row may become `real_pass` only when the evidence proves the before/after state, the user operation, the local data impact, and at least one downstream readback surface. Broad monthly/quarterly attribution rows may remain partial unless fresh current product evidence fully covers them.

