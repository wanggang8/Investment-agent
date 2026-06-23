# P85 Design

## Evidence Strategy

P85 should run at least three acceptance modes:

- Complete enough data: expected-return/scenario fields appear with provenance and deterministic calculations match independent expectations.
- Degraded data: affected analysis is marked degraded, blocked, or qualified with safe user guidance.
- LLM quality failure/unavailable: failed material is discarded or downgraded, deterministic-local fallback remains safe, and final verdict stays rule-governed.

The evidence must include UI screenshots or browser traces, API summaries, workflow metadata, deterministic calculation checks, and safety scans.

## Real-Pass Rule

A row may become `real_pass` only when the current product proves calculation correctness, provenance, degradation behavior, and decision-boundary safety through real local execution. Qualitative LLM prose alone is not sufficient.

