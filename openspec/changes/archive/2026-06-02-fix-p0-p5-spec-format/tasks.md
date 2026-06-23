## 1. Spec format validation fix

- [x] 1.1 Reproduce `p0-p5-capabilities` strict validation failure.
- [x] 1.2 Add the missing `## Purpose` section without changing P0-P5 requirement text.
- [x] 1.3 Verify `openspec validate p0-p5-capabilities --strict` passes.
- [x] 1.4 Verify `openspec validate --specs --strict` passes.
- [x] 1.5 Request subagent review after the fix.

## 2. Validation record

- [x] `openspec validate p0-p5-capabilities --strict`: passed.
- [x] `openspec validate --specs --strict`: passed, 6 passed / 0 failed.
