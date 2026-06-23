# Design: P65 Cross-Machine Release Repeat Acceptance

## Design Brief

P65 turns the P64 release package workflow into a repeatable handoff flow for a fresh P65 candidate archive. The core question is: can the project be verified from the packaged source in an isolated location, with no dependency on the active repository working tree?

The design uses a local isolated repeat as a cross-machine-equivalent gate. A physical second-machine run remains useful, but it should not be a required blocker for this local development phase.

## Repeat Flow

Add `scripts/local-release-repeat-acceptance.sh` as the repeat entrypoint. The script should accept:

- `--archive PATH`
- `--output-dir PATH` defaulting under project `tmp/`
- `--skip-install` for controlled reruns when dependencies already exist
- `--skip-e2e` only for diagnostic reruns, never for the main P65 acceptance record

The script should:

1. Resolve the archive and require an adjacent `release-manifest.json`.
2. Run `scripts/local-release-package.sh --verify <archive>` before extraction.
3. Extract the archive into `tmp/local-release-repeat/<timestamp>/workspace/`.
4. Detect the package root directory from the archive.
5. Run commands from the extracted package root:
   - `openspec validate --all --strict`
   - `go test ./...`
   - `npm --prefix web ci`
   - `npm --prefix web test`
   - `npm --prefix web run build`
   - `bash scripts/e2e-smoke.sh` with ports supplied by the repeat script
6. Write `repeat-summary.json` with package basename, package sha, release label, source commit, source status, extracted root placeholder, command statuses, timestamps, `skip_install`, `skip_e2e`, and safety note.

## Isolation Boundaries

The repeat script should write only under project `tmp/`. It should not write to real user databases, global config paths, home directories, or the active repository source tree. The extracted package's existing smoke scripts already use temporary SQLite and VecLite paths under the extracted workspace `tmp/`.

The repeat script should sanitize absolute paths in JSON summaries using placeholders such as `<repo>` and `<repeat-workspace>`.

## Documentation

Add `docs/release/acceptance/2026-06-18-p65-cross-machine-repeat.md` after execution. The document should record:

- package archive and manifest basenames;
- package SHA-256;
- repeat output summary path;
- command matrix;
- known caveats;
- physical second-machine follow-up command set;
- Not Claimed boundaries.

Update the release README and handoff so P65 becomes the recommended repeat handoff evidence after P64.

## Safety Boundaries

P65 must remain a verification phase:

- no remote publishing;
- no Git tag creation;
- no installer signing;
- no automatic upgrade, migration, restore, repair, rollback, or overwrite;
- no broker interface, automatic trading, one-click trading, order delegation, delegated order placement, external push, automatic confirmation, or automatic rule application;
- no login-gated, paid, authorization-gated, Level2, or high-frequency data source;
- no provider availability or return guarantees.

## Review Strategy

P65 follows the standard phase cadence:

1. Create change and plan.
2. Sub agent reviews plan.
3. Execute only if no Critical or Important findings remain.
4. Run repeat acceptance and regression gates.
5. Sub agent reviews execution.
6. Archive.
7. Sub agent reviews submit diff.
8. Commit.
