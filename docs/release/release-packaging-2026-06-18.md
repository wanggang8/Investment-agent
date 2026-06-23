# Release Packaging: 2026-06-18

> Change: P64 `p64-release-packaging-version-tagging`
> Basis: P63 `release_ready`
> Package mode: local source handoff
> Status: `p76_post_p75_package_passed`

## Summary

P64 adds a local release package workflow for Investment Agent. It packages release-safe project source, writes a sidecar `release-manifest.json`, records an archive SHA-256 checksum, and verifies the archive without executing runtime behavior.

This packaging workflow does not publish to remote storage, create a Git tag, execute upgrades, run migrations, restore data, repair files, call public providers, call LLM providers, trade, push notifications, confirm operations, or apply rules.

## Commands

Create a package:

```bash
bash scripts/local-release-package.sh --release-label p64-rc --output-dir tmp/p64-release
```

Verify a package:

```bash
bash scripts/local-release-package.sh --verify tmp/p64-release/20260618T052407Z/investment-agent-p64-rc.tar.gz --output-dir tmp/p64-release
```

Final distribution packages should be regenerated after acceptance milestones from a clean working tree so the manifest records `source_status: clean`. This document includes the historical P64 packaging record, the P69 clean-tree package refresh, and the P76 post-P75 package refresh. P76 is the current package freshness evidence for committed source through P75 plus the P76 acceptance-harness locator correction.

## P76 Post-P75 Package Refresh

| Field | Value |
| --- | --- |
| Release label | `p76-post-p75-final` |
| Source commit | `8a317f25917b8ff18ec9b5049e6a6188206a22d3` |
| Source status | `clean` |
| Archive | `tmp/p76-final-release/20260621T030713Z/investment-agent-p76-post-p75-final.tar.gz` |
| Manifest | `tmp/p76-final-release/20260621T030713Z/release-manifest.json` |
| SHA-256 | `7540429d0b6c3cdd09dad2ebb10e2356580faf0b05e6acd92bc3bd9763a3dcb7` |
| Archive entries | 1417 |
| Archive size | 3.0M |
| Verify summary | `tmp/p76-final-release/20260621T030723Z-verify/verify-summary.json` |
| Verify status | `passed` |
| Repeat summary | `tmp/p76-final-repeat/20260621T030727Z/repeat-summary.json` |
| Repeat status | `passed` |

P76 supersedes P69/P71 package freshness for final local source handoff through P75 evidence. Direct package file-list checks confirmed committed P72-P75 acceptance Markdown and OpenSpec archives are included in the archive. P76 itself is package-after-the-fact evidence and is not claimed to be included in this archive. The release status remains `release_ready_scoped_with_traceability_gaps`.

P76 also corrected stale Playwright smoke locators exposed by package repeat acceptance so evidence assertions target the visible P30 evidence row and P30 audit item. This is an acceptance-harness correction only and does not change runtime product behavior.

## P69 Clean Tree Package Refresh

| Field | Value |
| --- | --- |
| Release label | `p69-clean-tree` |
| Source commit | `cc0a64781e199a7745432b63bce26de4402042b5` |
| Source status | `clean` |
| Archive | `tmp/p69-final-release/20260618T084011Z/investment-agent-p69-clean-tree.tar.gz` |
| Manifest | `tmp/p69-final-release/20260618T084011Z/release-manifest.json` |
| SHA-256 | `d764ce5770289b6c174c919923ace181354165f8c8b114cfff444701cf158faa` |
| Archive entries | 1323 |
| Archive size | 2.8M |
| Verify summary | `tmp/p69-final-release/20260618T084023Z-verify/verify-summary.json` |
| Verify status | `passed` |
| Repeat summary | `tmp/p69-final-repeat/20260618T084028Z/repeat-summary.json` |
| Repeat status | `passed` |

P69 superseded P64/P65 dirty candidate package freshness for final local handoff through the P68 source commit. It is now historical package evidence; use P76 for the latest post-P75 local handoff package.

## P64 Acceptance Package

| Field | Value |
| --- | --- |
| Manifest | `tmp/p64-release/20260618T052407Z/release-manifest.json` |
| Archive | `tmp/p64-release/20260618T052407Z/investment-agent-p64-rc.tar.gz` |
| SHA-256 | `98ada08745c97596c7391518dcb4580e77f6497da0f4af96ef9c346f9fd3751a` |
| Archive entries | 1282 |
| Archive size | 2.7M |
| Verify summary | `tmp/p64-release/20260618T052420Z-verify/verify-summary.json` |
| Verify status | `passed` |
| Source status | `dirty` during active P64 implementation; regenerate after commit for distribution |

The P64 acceptance package is a local `tmp/` artifact and is not committed.
Because this document records the local acceptance artifact itself, final distribution packages should be regenerated from the committed clean tree and verified with their adjacent sidecar manifest. For final distribution after P65-P68, use the P69 clean-tree package evidence rather than treating the P64/P65 dirty candidate artifacts as current.

## Manifest Contract

`release-manifest.json` records:

- release label;
- source commit;
- generated timestamp;
- archive basename;
- archive SHA-256;
- source status;
- included roots and package metadata entries;
- excluded patterns;
- verification commands;
- acceptance references;
- known non-blocking degradations;
- Not Claimed boundaries;
- safety note.

The manifest uses archive basenames and relative project references. It must not contain complete API keys, private paths, raw SQL dumps, raw stack traces, complete prompt payload files, raw vendor payloads, local DB paths, or unredacted logs.

## Included Roots

- `AGENTS.md`
- `.gitignore`
- `cmd/`
- `configs/config.example.yaml`
- `docs/`
- `examples/`
- `internal/`
- `openspec/`
- `pkg/`
- `scripts/`
- `web/`
- `go.mod`
- `go.sum`

## Excluded Patterns

- `.git/`
- `.cursor/`
- `tmp/`
- `cmd/agent/tmp/`
- `docs/release/ui-audit-assets/`
- `configs/config.local.yaml`
- `web/node_modules/`
- `web/dist/`
- `playwright-report/`
- `test-results/`
- `*.db`
- `*.sqlite`
- `*.sqlite3`
- `*.log`
- `*.trace`
- raw provider payloads
- complete prompt payload files
- complete API keys
- private local paths

## Verification Result

The P64 package verification checks:

- sidecar manifest exists and parses;
- archive checksum matches manifest `package_sha256`;
- required package entrypoints exist;
- forbidden path patterns are absent;
- manifest text does not include complete key patterns, private path patterns, bearer tokens, private key headers, or prompt payload markers.

The acceptance run passed package verification. Earlier package attempts failed because tracked `cmd/agent/tmp/veclite` and raw UI audit assets would have entered the archive, the first output-dir contract allowed locations outside project `tmp/`, and broad scanning treated policy or fixture text as leaks. The script now rejects output directories outside project `tmp/`, excludes local editor/UI audit/raw runtime artifacts, packages tracked files plus only explicit active-P64 allowlisted new files before commit, and scans for leak-shaped complete keys, private paths, prompt payloads, private keys, bearer tokens, and raw payload markers.

## Known Release Context

P64 did not replace P63 acceptance. It packaged the project based on P63 release evidence:

- `docs/release/acceptance/2026-06-18-p63-full-ui-regression.md`
- `docs/release/release-candidate-2026-06-18.md`
- `docs/release/release-handoff-2026-06-18.md`

Known non-blocking degradations remain:

- G5 current data-source quality was degraded with zero failed cases.
- P63 real UI consultation recorded `VECTOR_INDEX_UNAVAILABLE` while LLM reports parsed and passed quality.
- P63 full UI regression recorded classified expected 404/409 API responses and zero unexpected failed API responses.

Later package refreshes supersede this historical context for package freshness. P76 package evidence includes committed P72-P75 acceptance and OpenSpec archive materials while preserving P75's scoped release conclusion.

## Not Claimed

This packaging workflow does not claim:

- future public-source availability;
- future model-provider availability;
- investment returns;
- broker connectivity;
- automatic trading;
- one-click trading;
- order delegation;
- external push;
- automatic confirmation;
- automatic rule application;
- automatic repair;
- automatic upgrade;
- automatic migration;
- automatic restore;
- real database overwrite;
- login sources;
- paid sources;
- authorized sources;
- Level2 data;
- high-frequency data.
