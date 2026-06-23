# Design

## Context

The `p0-p5-capabilities` OpenSpec summary already contains valid `## Requirements` and scenario blocks, but strict validation also requires `## Purpose`. The issue is structural formatting, not missing product behavior.

## Approach

Add a concise `## Purpose` section above `## Requirements`. Preserve all existing requirements exactly so the change remains limited to OpenSpec validation compatibility.

## Risks

- Changing requirement text could unintentionally alter historical capability meaning. This change avoids that by only adding Purpose text.
