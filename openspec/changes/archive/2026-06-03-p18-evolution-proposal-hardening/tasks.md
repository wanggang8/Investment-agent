## 1. Proposal chain tests

- [x] Add tests for review/evolution input creating a durable rule proposal record.
- [x] Assert created proposal status is non-applied and enters the existing review queue.
- [x] Assert proposal source metadata references review period, supporting decisions, audit events, or error cases where available.

## 2. Safety and gate tests

- [x] Add tests that proposal generation does not write a new active rule version.
- [x] Add tests that insufficient samples or missing source facts do not create applied proposals.
- [x] Add tests that gatekeeper audit and final user confirmation remain required before rule activation.

## 3. Implementation

- [x] Implement or refine `EvolutionProposalGraph` or equivalent service path for review-generated proposals.
- [x] Persist proposal metadata using existing repository/status fields where possible.
- [x] Write audit events linking review input to proposal output.
- [x] Preserve existing rule proposal, gatekeeper, and final confirmation behavior.
- [x] Avoid broad rule engine or frontend rewrites.

## 4. Documentation and OpenSpec sync

- [x] Document P18 boundary: review can create proposals, but cannot apply rules.
- [x] Keep P18 delta specs local before archive.

## 5. Validation

- [x] Run `openspec validate p18-evolution-proposal-hardening --strict`.
- [x] Run targeted Go tests for evolution proposal, rule proposal repository/service, and gatekeeper path.
- [x] Run `go test ./...`.
- [x] Run `openspec validate --all --strict`.
