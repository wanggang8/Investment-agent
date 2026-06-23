## MODIFIED Requirements

### Requirement: P0-P5 capability summary has OpenSpec purpose
The `p0-p5-capabilities` summary SHALL include a `## Purpose` section and a `## Requirements` section so strict OpenSpec validation can parse the spec.

#### Scenario: Strict specs validation includes P0-P5 summary
- **WHEN** `openspec validate --specs --strict` is executed
- **THEN** `p0-p5-capabilities` MUST pass validation
- **AND** existing P0-P5 capability requirements MUST remain unchanged
