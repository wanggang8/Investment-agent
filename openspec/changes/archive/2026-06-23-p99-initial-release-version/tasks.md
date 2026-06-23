# Tasks

- [x] Add root `VERSION` with `v0.1.0`.
- [x] Update frontend package metadata from `0.0.0` to `0.1.0`.
- [x] Update release materials, governance notes, and progress to record P99.
- [x] Add P99 acceptance record.
- [x] Run validation:
  - [x] `openspec validate p99-initial-release-version --strict`
  - [x] `npm --prefix web test`
  - [x] `npm --prefix web run build`
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p93_code_reality_audit.py --check`
  - [x] `python3 scripts/p91_deployment_check.py --check`
  - [x] `git diff --check`
- [x] Archive P99 and run `openspec validate --all --strict`.
