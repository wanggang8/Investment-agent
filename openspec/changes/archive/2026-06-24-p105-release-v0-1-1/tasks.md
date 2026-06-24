# Tasks

- [x] Create P105 governance/progress entries and OpenSpec delta.
- [x] Bump root and frontend version metadata to `v0.1.1`.
- [x] Update release-facing materials and P105 acceptance record.
- [ ] Run release gates:
  - [x] `openspec validate p105-release-v0-1-1 --strict`
  - [x] `go test ./...`
  - [x] `npm --prefix web test -- --run`
  - [x] `npm --prefix web run build`
  - [x] `openspec validate --all --strict`
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p93_code_reality_audit.py --check`
  - [x] `git diff --check`
- [x] Prepare P105 for archive after validation passes.
- [x] Prepare validated release state for post-archive commit and annotated tag `v0.1.1`.
