# Tasks

- [x] Create P121 governance/progress entries and OpenSpec delta.
- [x] Add a P121 final release review checker and acceptance record.
- [x] Bump root/frontend version metadata to `v0.1.3`.
- [x] Update release-facing materials and release notes.
- [x] Sanitize local absolute paths during release package text copy.
- [x] Run release gates:
  - [x] `openspec validate p121-final-review-and-v0-1-3-tag-release --strict`
  - [x] `openspec validate --all --strict`
  - [x] `go test ./...`
  - [x] `go vet ./...`
  - [x] `npm --prefix web test -- --run`
  - [x] `npm --prefix web run build`
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p121_final_release_review.py --check`
  - [x] `git diff --check`
  - [x] local release package smoke/verify for `v0.1.3`
- [x] Prepare P121 for archive after validation passes.
- [x] Prepare the reviewed release state for post-archive commit, annotated tag `v0.1.3`, and push.
