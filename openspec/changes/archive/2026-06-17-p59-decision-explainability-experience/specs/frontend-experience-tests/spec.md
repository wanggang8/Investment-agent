## ADDED Requirements

### Requirement: Decision explanation surfaces present a readable decision story

The frontend SHALL connect consultation, decision detail, evidence, and decision-loop pages into a readable decision story while preserving existing backend contracts and safety boundaries.

#### Scenario: Consultation route makes the generated decision path explicit

- **WHEN** the user opens `/consultation`
- **THEN** the page MUST clearly show the consultation inputs, assumptions, generation state, and safety boundary before the generated result
- **AND** after a successful consultation it MUST show the generated verdict or safe unavailable state, key reasons, prohibited actions, optional manual actions, data trust context, and local navigation to the generated decision detail when a decision id exists
- **AND** the page MUST NOT present the generated result as an automatic confirmation, trade, rule application, external push, or return promise

#### Scenario: Decision detail first screen explains the verdict before technical trace

- **WHEN** the user opens `/decisions/:decisionId`
- **THEN** the first screen MUST show the final verdict or safe unavailable state, generated context, prohibited actions, optional manual actions, key reasons, data trust summary, and safety boundary before long technical details
- **AND** Evidence, LLM, rules, expected return, arbitration, audit, and confirmation details MUST be grouped into readable layers
- **AND** long traces MUST not obscure the first-screen verdict, safety boundary, or explanation path

#### Scenario: Decision explanation handles nullable and degraded DTOs safely

- **WHEN** decision, verdict, evidence, analyst, expected-return, retrieval-quality, audit, or confirmation fields are null, missing, empty, degraded, unknown, failed, or insufficient
- **THEN** the frontend MUST render safe Chinese empty/degraded text without page-level crashes
- **AND** nullable or missing fields MUST NOT be described as permission to trade, auto-confirm, auto-apply rules, or treat the decision as successful without review

#### Scenario: Evidence page explains source trust and links back to decision explanation

- **WHEN** the user opens `/evidence`
- **THEN** the page MUST prioritize a source trust summary, source-level explanation, verification status, and local navigation back to decision explanation surfaces before raw evidence detail
- **AND** the evidence table MUST preserve filtering and expansion without exposing raw vendor payload, complete prompt, private path, key, or local database content

#### Scenario: Decision loop page reads as a read-only decision lifecycle

- **WHEN** the user opens `/decision-loop`
- **THEN** the page MUST show a read-only lifecycle from recommendation to confirmation, manual record, risk/review, and audit links
- **AND** it MUST clearly identify missing links or open gaps as manual follow-up items
- **AND** it MUST NOT provide controls that create confirmations, trades, risk lifecycle changes, rule applications, notifications, or settings changes

#### Scenario: Decision explanation remains mobile readable and safe

- **WHEN** `/consultation`, `/decisions/:decisionId`, `/evidence`, or `/decision-loop` renders at 390px viewport width
- **THEN** primary verdict, safety boundary, trust context, key reasons, and local navigation MUST remain readable without page-level horizontal overflow
- **AND** screenshots or browser evidence MUST be captured for desktop and mobile validation
