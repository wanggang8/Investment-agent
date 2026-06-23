# P18 Evolution Proposal Hardening

## Summary
Complete and harden the review-to-rule-proposal path so review outputs can create traceable rule proposal drafts or records while preserving gatekeeper audit and user final confirmation.

## Why
P9 introduced review summaries and rule effectiveness evaluation. Later phases added stronger evidence, audit, and frontend tracking surfaces. The remaining gap is the proposal chain: review-derived suggestions must become explicit proposal artifacts that enter the existing review queue, remain traceable to source review facts, and never apply rules automatically.

## Scope
- Harden `EvolutionProposalGraph` or equivalent service path for review-generated rule proposals.
- Persist enough proposal metadata for audit and frontend tracking.
- Ensure proposal status enters the existing gated queue, not an applied rule version.
- Add tests for proposal creation, audit events, rejected/insufficient sample cases, and no automatic application.
- Keep frontend behavior display-only unless existing rule proposal entrypoints already support safe navigation.

## Non-goals
- No automatic rule version application.
- No trading or portfolio mutation.
- No new investment recommendation behavior.
- No broad rewrite of the rule engine or gatekeeper graph.
