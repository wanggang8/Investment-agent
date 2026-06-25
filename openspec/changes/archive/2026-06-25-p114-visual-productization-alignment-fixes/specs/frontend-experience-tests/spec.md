## ADDED Requirements

### Requirement: Productized Visual Alignment Review

The frontend SHALL pass a productized visual alignment review for forms, actions, same-level cards, and first-layer technical content before a reference-driven UI polish change is considered complete.

#### Scenario: Form and action alignment

- **GIVEN** a page contains form fields, selects, textareas, helper text, error text, and action buttons
- **WHEN** the page is inspected at desktop and approximately 390px mobile widths
- **THEN** fields and actions align to a coherent visual baseline
- **AND** action buttons do not float at inconsistent heights or appear detached from the form section
- **AND** mobile actions use stable full-width or grouped layouts with readable text and sufficient spacing.

#### Scenario: Same-level card consistency

- **GIVEN** a page displays cards in the same semantic group or grid row
- **WHEN** cards have different amounts of content
- **THEN** the cards maintain consistent title treatment, padding, visual weight, and action placement
- **AND** content length does not create a ragged, unintentional hierarchy among equal-priority cards.

#### Scenario: Productized first-layer content

- **GIVEN** a page receives raw diagnostics, commands, local paths, JSON, internal enum values, or detailed provider/readiness material
- **WHEN** the first user-facing layer is rendered
- **THEN** the first layer presents productized status, explanation, and next human action
- **AND** technical material is either mapped to product language, folded into secondary details, or explicitly classified as requiring backend summary support.

#### Scenario: Backend ownership classification

- **GIVEN** a visual issue is caused by unproductized data content
- **WHEN** the issue is added to the visual finding ledger
- **THEN** it is classified as `frontend-mapping`, `backend-summary-needed`, or `intentional-technical-secondary`
- **AND** backend changes are made only when existing DTOs cannot support a safe productized summary.
