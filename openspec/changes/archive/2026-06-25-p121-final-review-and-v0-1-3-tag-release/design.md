# Design

P121 is a release-governance change. It does not introduce new product behavior; it makes the post-P114-P120 repository state shippable as `v0.1.3` only after a fresh review.

The release decision uses three layers:

1. Governance and traceability checks: OpenSpec strict validation, active-change state, P114-P120 archive presence, release materials, and version consistency.
2. Runtime build/test checks: Go test/vet, frontend unit tests, frontend production build, and whitespace checks.
3. Release packaging/tag checks: local package smoke/verify, commit, annotated tag, and push.

P93 is intentionally not re-labeled as fresh after P114-P120 because its checker reports stale evidence. P121 records this as a release boundary and supplies a fresh P121-specific review for the current tree.

The local release package workflow keeps strict forbidden-content scanning. For text files copied into the archive, it first replaces local workstation paths with stable placeholders such as `<repo>`, `<user-home>`, and `<codex-generated-images>/`. This avoids changing archived source evidence while ensuring the distributed package does not contain local absolute paths.
