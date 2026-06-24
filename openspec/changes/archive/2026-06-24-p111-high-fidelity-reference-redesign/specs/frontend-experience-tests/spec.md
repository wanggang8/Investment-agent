## ADDED Requirements

### Requirement: P111 high-fidelity reference redesign blocks archive on visual mismatch

The frontend SHALL implement the selected Calm Command Center reference image as a high-fidelity visual target across the product, and SHALL block archive when covered pages have unresolved material visual mismatches.

#### Scenario: Reference image is the visual source of truth

- **WHEN** P111 begins implementation
- **THEN** `/Users/vick/.codex/generated_images/019ef8a7-f5c0-7442-95b9-e72bcdc89de1/ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png` MUST be treated as the selected visual source of truth
- **AND** P111 MUST extract concrete layout, density, navigation, topbar, hero, action queue, metric grid, snapshot, progress tracker, checklist, ledger, typography, border, radius, color, icon, and spacing rules from that image
- **AND** P111 MUST NOT treat the image as a loose moodboard once implementation starts

#### Scenario: Core cockpit pages match reference layout language

- **WHEN** `/` or `/workbench` is rendered after P111
- **THEN** the first desktop viewport MUST include a reference-style sidebar, top status toolbar, report hero, current discipline state, prohibited action block, priority manual action queue, status metric overview, portfolio or funds snapshot, recent decision or consultation progress preview, and evidence/rules checklist when source data or safe fallback data exists
- **AND** the page MUST NOT regress to P110-style generic vertical card stacking as the primary first viewport structure
- **AND** visible actions MUST remain local navigation, local maintenance, offline record, or manual review actions

#### Scenario: All covered pages have page-level fidelity ledgers

- **WHEN** a P111 covered page is marked complete
- **THEN** the acceptance evidence MUST include a screenshot for that route and viewport
- **AND** the acceptance evidence MUST include a mismatch ledger entry for that page
- **AND** the mismatch ledger MUST compare reference evidence, render evidence, mismatch level, fix made, and any intentional deviation
- **AND** unresolved P0, P1, or P2 mismatches MUST block marking that page complete

#### Scenario: Reference component language applies beyond the homepage

- **WHEN** P111 covers maintenance, evidence, decision, governance, or ops routes
- **THEN** those routes MUST reuse the reference component language for topbar, status hero, action queue, metric grid, snapshot strip, progress tracker, evidence checklist, and ledger surfaces as applicable
- **AND** pages MUST NOT only receive color/token changes while keeping an unrelated engineering-admin layout
- **AND** every intentional page-specific deviation MUST be recorded in the mismatch ledger

#### Scenario: Visual fidelity QA runs page by page

- **WHEN** P111 implementation proceeds across pages
- **THEN** desktop screenshot comparison MUST happen after each covered page or page group, before moving to the next group
- **AND** a final all-route pass MUST still capture desktop and responsive evidence
- **AND** 390px, 768px, and 1280px checks MUST confirm no page-level horizontal overflow, except scoped table/log/diagnostic containers

#### Scenario: High-fidelity redesign preserves product and safety contracts

- **WHEN** P111 changes shared components or pages
- **THEN** pages MUST continue using existing services, API DTOs, route semantics, manual confirmation flows, and redaction utilities
- **AND** P111 MUST NOT require new backend API fields, SQLite schema changes, Eino workflow changes, LLM prompt changes, data source changes, rule engine changes, release package changes, or physical second-machine verification
- **AND** P111 MUST NOT add or imply broker connectivity, automatic trading, one-click trading, order delegation, external push, automatic confirmation, automatic rule application, automatic repair, real database overwrite, return promise, login source, paid source, authorized source, Level2 data, or high-frequency source
