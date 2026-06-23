# P82 Design

## Evidence Strategy

P82 must prove causality: each accepted UI action should have an expected API result, table mutation or durable readback, audit event, and visible user-facing state after navigation or refresh.

Representative scenarios should include:

- Risk alert SOP open, state transition, note or follow-up, and failure/error classification.
- Manual confirmation or non-trade record where allowed by current product boundaries.
- Notification/readback and daily review surfaces that reflect the action.
- Negative checks that unsafe actions are absent or blocked.

## Real-Pass Rule

A row may become `real_pass` only when a browser-driven operation and read-only SQLite/readback evidence prove the intended state transition or durable record. API-only smoke, route availability, screenshots without data readback, or mocked network evidence are insufficient.

