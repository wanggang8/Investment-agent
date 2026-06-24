# Tasks

- [x] Create P106 governance/progress entries and OpenSpec delta.
- [x] Fix Data Quality redaction label package-scan false positive.
- [x] Bump root/frontend version metadata to `v0.1.2`.
- [x] Update release-facing materials and P106 acceptance record.
- [x] Run release gates:
  - [x] `openspec validate p106-release-v0-1-2-package-scan-fix --strict`
  - [x] `go test ./...`
  - [x] `npm --prefix web test -- --run`
  - [x] `npm --prefix web run build`
  - [x] `bash scripts/local-release-package.sh --release-label v0.1.2 --output-dir tmp/p106-release-package`
  - [x] `bash scripts/local-release-package.sh --verify <archive> --output-dir tmp/p106-release-package`
  - [x] `openspec validate --all --strict`
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p93_code_reality_audit.py --check`
  - [x] `git diff --check`
- [x] Prepare P106 for archive after validation passes.
- [x] Prepare validated release state for post-archive commit and annotated tag `v0.1.2`.
