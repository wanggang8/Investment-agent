#!/usr/bin/env python3
"""Validate the P121 final release review for v0.1.3."""

from __future__ import annotations

import argparse
import json
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
CHANGE_ID = "p121-final-review-and-v0-1-3-tag-release"
VERSION_TAG = "v0.1.3"
VERSION_NPM = "0.1.3"
RELEASE_NOTES = ROOT / "docs/release/release-v0.1.3.md"
ACCEPTANCE = ROOT / "docs/release/acceptance/2026-06-25-p121-final-review-and-v0.1.3-tag-release.md"
RELEASE_README = ROOT / "docs/release/README.md"
ROOT_IMAGE = ROOT / "ig_05724f56eb7089ab016a3b9109e1848191a87e68883d0c9826.png"

REQUIRED_ARCHIVES = [
    "2026-06-25-p114-visual-productization-alignment-fixes",
    "2026-06-25-p115-real-user-scenario-acceptance",
    "2026-06-25-p116-multi-fund-transaction-ledger-acceptance",
    "2026-06-25-p117-continuous-product-usability-acceptance",
    "2026-06-25-p118-product-usability-edge-scenario-acceptance",
    "2026-06-25-p119-full-ui-control-and-affordance-acceptance",
    "2026-06-25-p120-p114-p119-final-closure-summary",
]

REQUIRED_ACCEPTANCE_DOCS = [
    "docs/release/acceptance/2026-06-24-p114-visual-productization-alignment-fixes.md",
    "docs/release/acceptance/2026-06-25-p115-real-user-scenario-acceptance.md",
    "docs/release/acceptance/2026-06-25-p116-multi-fund-transaction-ledger-acceptance.md",
    "docs/release/acceptance/2026-06-25-p117-continuous-product-usability-acceptance.md",
    "docs/release/acceptance/2026-06-25-p118-product-usability-edge-scenario-acceptance.md",
    "docs/release/acceptance/2026-06-25-p119-full-ui-control-and-affordance-acceptance.md",
    "docs/release/acceptance/2026-06-25-p114-p119-final-closure-summary.md",
]

REQUIRED_BOUNDARY_PHRASES = [
    "P93 remains a historical final code-reality/design audit",
    "does not claim a fresh P93 pass",
    "broker connectivity",
    "automatic trading",
    "one-click trading",
    "order delegation",
    "external push",
    "automatic confirmation",
    "automatic rule application",
    "physical second-machine validation",
    "investment returns",
]

REQUIRED_GATE_LINES = [
    "openspec validate p121-final-review-and-v0-1-3-tag-release --strict",
    "openspec validate --all --strict",
    "go test ./...",
    "go vet ./...",
    "npm --prefix web test -- --run",
    "npm --prefix web run build",
    "python3 scripts/p92_final_requirement_audit.py --check",
    "python3 scripts/p121_final_release_review.py --check",
    "git diff --check",
    "bash scripts/local-release-package.sh --release-label v0.1.3 --output-dir tmp/p121-release-package-final",
]


def run(cmd: list[str]) -> str:
    return subprocess.run(cmd, cwd=ROOT, check=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True).stdout.strip()


def read(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def check_versions(errors: list[str]) -> None:
    version = read(ROOT / "VERSION").strip()
    if version != VERSION_TAG:
        errors.append(f"VERSION is {version!r}, expected {VERSION_TAG!r}")

    package_json = json.loads(read(ROOT / "web/package.json"))
    if package_json.get("version") != VERSION_NPM:
        errors.append(f"web/package.json version is {package_json.get('version')!r}, expected {VERSION_NPM!r}")

    package_lock = json.loads(read(ROOT / "web/package-lock.json"))
    if package_lock.get("version") != VERSION_NPM:
        errors.append(f"web/package-lock.json version is {package_lock.get('version')!r}, expected {VERSION_NPM!r}")
    root_package = package_lock.get("packages", {}).get("", {})
    if root_package.get("version") != VERSION_NPM:
        errors.append(f"web/package-lock root package version is {root_package.get('version')!r}, expected {VERSION_NPM!r}")


def check_governance(errors: list[str]) -> None:
    active_dirs = [
        p.name
        for p in (ROOT / "openspec/changes").iterdir()
        if p.is_dir() and p.name != "archive"
    ]
    if active_dirs not in ([], [CHANGE_ID]):
        errors.append(f"active OpenSpec changes are {active_dirs!r}, expected none or only {CHANGE_ID!r}")

    for archive in REQUIRED_ARCHIVES:
        if not (ROOT / "openspec/changes/archive" / archive).is_dir():
            errors.append(f"missing archive {archive}")

    for doc in REQUIRED_ACCEPTANCE_DOCS:
        if not (ROOT / doc).is_file():
            errors.append(f"missing acceptance doc {doc}")

    governance = read(ROOT / "docs/GOVERNANCE.md")
    progress = read(ROOT / "openspec/PROGRESS.md")
    project = read(ROOT / "openspec/project.md")
    if f"当前活跃变更：`{CHANGE_ID}`" not in governance and "当前活跃变更：无。" not in governance:
        errors.append("docs/GOVERNANCE.md does not identify either active P121 or no active change after P121 archive")
    if f"| **current_change** | `{CHANGE_ID}` |" not in progress and "| **current_change** | `none` |" not in progress:
        errors.append("openspec/PROGRESS.md does not identify either current P121 or no active change after P121 archive")
    if f"`{CHANGE_ID}` | active" not in project and f"`{CHANGE_ID}` | done" not in project:
        errors.append("openspec/project.md does not list P121 as active or done")


def check_release_docs(errors: list[str]) -> None:
    for path in [RELEASE_NOTES, ACCEPTANCE, RELEASE_README]:
        if not path.is_file():
            errors.append(f"missing release document {path.relative_to(ROOT)}")
            return

    release = read(RELEASE_NOTES)
    acceptance = read(ACCEPTANCE)
    readme = read(RELEASE_README)
    combined = "\n".join([release, acceptance, readme])

    if "pending" in acceptance.lower() or "pending" in release.lower():
        errors.append("release notes or P121 acceptance record still contain pending status")
    if "Status: `ready_after_p121_final_review`" not in release:
        errors.append("release notes do not mark ready_after_p121_final_review")
    if "Status: `passed`" not in acceptance:
        errors.append("P121 acceptance record does not mark Status: passed")
    if f"Current source release version: `{VERSION_TAG}`." not in readme:
        errors.append("docs/release/README.md does not advertise v0.1.3 as current source release version")

    for phrase in REQUIRED_BOUNDARY_PHRASES:
        if phrase not in combined:
            errors.append(f"release materials missing boundary phrase: {phrase}")
    for line in REQUIRED_GATE_LINES:
        if line not in acceptance:
            errors.append(f"P121 acceptance record missing gate line: {line}")

    if "Archive | `investment-agent-v0.1.3.tar.gz`" not in acceptance:
        errors.append("P121 acceptance record does not contain v0.1.3 archive identity")
    if "SHA256 | `" not in acceptance:
        errors.append("P121 acceptance record does not contain package SHA256")


def check_tag_and_assets(errors: list[str]) -> None:
    tags = set(run(["git", "tag", "--list", VERSION_TAG]).splitlines())
    if VERSION_TAG in tags:
        errors.append(f"tag {VERSION_TAG} already exists before P121 publication")

    tracked_root_image = run(["git", "ls-files", str(ROOT_IMAGE.relative_to(ROOT))])
    if tracked_root_image:
        errors.append("reference PNG copied into repository root is tracked; release should use documented generated-image path instead")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--check", action="store_true", help="Validate P121 release review state.")
    args = parser.parse_args()

    if not args.check:
        parser.error("use --check")

    errors: list[str] = []
    check_versions(errors)
    check_governance(errors)
    check_release_docs(errors)
    check_tag_and_assets(errors)

    if errors:
        print("P121 final release review: failed")
        for error in errors:
            print(f"- {error}")
        return 1

    print("P121 final release review: passed")
    print(f"- version={VERSION_TAG}")
    print(f"- change={CHANGE_ID}")
    print("- p114_p120_archives=present")
    print("- p93_stale_boundary=explicit")
    print("- release_notes=ready")
    return 0


if __name__ == "__main__":
    sys.exit(main())
