# Tasks

- [x] Create P104 governance/progress entries and OpenSpec delta.
- [x] Add product operation/linkage acceptance matrix.
- [x] Add repeatable local-source P104 runner.
- [x] Execute P104 runner and record fresh evidence.
- [x] Run regression gates:
  - [x] `openspec validate p104-full-product-operation-linkage-acceptance --strict`
  - [x] `go test ./...`
  - [x] `npm --prefix web test -- --run`
  - [x] `npm --prefix web run build`
  - [x] `openspec validate --all --strict`
  - [x] `python3 scripts/p92_final_requirement_audit.py --check`
  - [x] `python3 scripts/p93_code_reality_audit.py --check`
  - [x] `git diff --check`
- [x] Archive P104 after validation passes.
