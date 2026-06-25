#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT_DIR="${ROOT_DIR}/tmp/local-release-package"
RELEASE_LABEL=""
VERIFY_ARCHIVE=""
SKIP_BUILD="0"
TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"

usage() {
  cat <<EOF
Usage: bash scripts/local-release-package.sh [options]

Options:
  --release-label VALUE  Release label to write into the manifest
  --output-dir PATH      Output directory (default: tmp/local-release-package)
  --verify ARCHIVE       Verify an existing local release archive
  --skip-build           Skip frontend build sanity check before packaging
  -h, --help             Show this help
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --release-label)
      RELEASE_LABEL="$2"
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    --verify)
      VERIFY_ARCHIVE="$2"
      shift 2
      ;;
    --skip-build)
      SKIP_BUILD="1"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "未知参数: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

safe_label() {
  python3 - "$1" <<'PY'
import re
import sys

label = sys.argv[1].strip()
label = re.sub(r"[^A-Za-z0-9._-]+", "-", label).strip("-._")
print(label or "local-release")
PY
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

normalize_output_dir() {
  python3 - "$ROOT_DIR" "$OUTPUT_DIR" <<'PY'
import sys
from pathlib import Path

root = Path(sys.argv[1]).resolve()
raw = Path(sys.argv[2])
candidate = raw if raw.is_absolute() else root / raw
resolved = candidate.resolve(strict=False)
tmp_root = root / "tmp"
try:
    resolved.relative_to(tmp_root)
except ValueError:
    print(f"--output-dir must be inside {tmp_root}", file=sys.stderr)
    sys.exit(1)
print(resolved)
PY
}

copy_release_files() {
  local package_root="$1"
  local list_file="$package_root/.release-source-files.txt"
  local untracked_file="$package_root/.release-untracked-files.txt"
  git -C "$ROOT_DIR" ls-files --cached >"$list_file"
  git -C "$ROOT_DIR" ls-files --others --exclude-standard >"$untracked_file"
  python3 - "$ROOT_DIR" "$package_root" "$list_file" "$untracked_file" <<'PY'
import os
import re
import shutil
import sys
from pathlib import Path

root = Path(sys.argv[1])
package_root = Path(sys.argv[2])
list_file = Path(sys.argv[3])
untracked_file = Path(sys.argv[4])

allowed_untracked_exact = {
    ".dockerignore",
    ".env.example",
    "Dockerfile",
    "docker-compose.yml",
    "configs/config.docker.yaml",
    "docs/deployment.md",
    "docs/release/acceptance/2026-06-22-p91-github-release-docker-deployment.md",
    "docs/release/acceptance/2026-06-22-p92-final-original-requirement-audit-ledger.md",
    "docs/release/acceptance/2026-06-22-p92-final-original-requirement-audit-summary.md",
    "docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md",
    "scripts/backup.sh",
    "scripts/deploy-lib.sh",
    "scripts/doctor.sh",
    "scripts/install.sh",
    "scripts/p91_deployment_check.py",
    "scripts/p92_final_requirement_audit.py",
    "scripts/p93_code_reality_audit.py",
    "scripts/status.sh",
    "scripts/uninstall.sh",
    "scripts/upgrade.sh",
}
allowed_untracked_prefixes = (
    ".github/workflows/",
    "docker/",
    "openspec/changes/p91-github-release-docker-deployment/",
    "openspec/changes/p92-final-original-requirement-audit-ledger/",
    "openspec/changes/p93-final-code-reality-design-audit/",
    "openspec/changes/archive/2026-06-22-p91-github-release-docker-deployment/",
    "openspec/changes/archive/2026-06-22-p92-final-original-requirement-audit-ledger/",
    "openspec/changes/archive/2026-06-22-p93-final-code-reality-design-audit/",
)

excluded_exact = {
    "configs/config.local.yaml",
}
excluded_prefixes = (
    ".cursor/",
    ".git/",
    "tmp/",
    "cmd/agent/tmp/",
    "docs/release/ui-audit-assets/",
    "web/node_modules/",
    "web/dist/",
    "playwright-report/",
    "test-results/",
)
excluded_suffixes = (
    ".db",
    ".sqlite",
    ".sqlite3",
    ".log",
    ".trace",
)
excluded_names = {
    ".DS_Store",
}
forbidden_name_markers = (
    "secret",
    "credential",
    "private-key",
    "raw-payload",
    "vendor-payload",
    "sql-dump",
    "prompt-payload",
)
forbidden_content = [
    re.compile(r"(?<![A-Za-z0-9])sk-[A-Za-z0-9_-]{20,}", re.IGNORECASE),
    re.compile(r"/Users/(?!private\b)[A-Za-z0-9._-]+/", re.IGNORECASE),
    re.compile(r"BEGIN (RSA|OPENSSH|PRIVATE) KEY", re.IGNORECASE),
    re.compile(r"Authorization:\s*Bearer", re.IGNORECASE),
    re.compile(r"(?:^|[,{]\s*)[\"']?prompt[\"']?\s*:\s*[\"'][^\"'\n]{16,}", re.IGNORECASE),
    re.compile(r"[\"']raw_(provider|vendor)_payload[\"']\s*:\s*[\[{\"]", re.IGNORECASE),
]

included = []
excluded = []
errors = []
tracked = list_file.read_text(encoding="utf-8").splitlines()
untracked = untracked_file.read_text(encoding="utf-8").splitlines()
allowed_untracked = []

for raw in untracked:
    rel = raw.strip("\n")
    if not rel:
        continue
    if rel in allowed_untracked_exact or any(rel.startswith(prefix) for prefix in allowed_untracked_prefixes):
        allowed_untracked.append(rel)
    else:
        excluded.append(f"{rel}\t(untracked_not_release_allowlisted)")

def sanitize_text(text: str) -> str:
    text = re.sub(r"/Users/[A-Za-z0-9._-]+/Desktop/project/Investment-agent", "<repo>", text)
    text = re.sub(r"/Users/[A-Za-z0-9._-]+/\.codex/generated_images/", "<codex-generated-images>/", text)
    text = re.sub(r"/Users/[A-Za-z0-9._-]+/", "<user-home>/", text)
    return text

def scan_and_sanitize_content(rel: str, source: Path):
    try:
        data = source.read_bytes()
    except OSError as exc:
        errors.append(f"failed to read {rel}: {exc}")
        return None
    if b"\x00" in data:
        return None
    try:
        text = data.decode("utf-8")
    except UnicodeDecodeError:
        return None
    text = sanitize_text(text)
    for pattern in forbidden_content:
        if pattern.search(text):
            errors.append(f"forbidden content pattern in {rel}: {pattern.pattern}")
            return text
    return text

for raw in tracked + allowed_untracked:
    rel = raw.strip("\n")
    if not rel:
        continue
    name = os.path.basename(rel)
    lower = rel.lower()
    should_exclude = (
        rel in excluded_exact
        or any(rel.startswith(prefix) for prefix in excluded_prefixes)
        or any(lower.endswith(suffix) for suffix in excluded_suffixes)
        or name in excluded_names
        or any(marker in lower for marker in forbidden_name_markers)
    )
    if should_exclude:
        excluded.append(rel)
        continue
    source = root / rel
    if not source.is_file():
        excluded.append(rel)
        continue
    sanitized_text = scan_and_sanitize_content(rel, source)
    target = package_root / rel
    target.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(source, target)
    if sanitized_text is not None:
        target.write_text(sanitized_text, encoding="utf-8")
    included.append(rel)

if errors:
    (package_root / "release-file-list.txt").write_text("\n".join(included) + ("\n" if included else ""), encoding="utf-8")
    (package_root / "release-excluded-list.txt").write_text("\n".join(excluded) + ("\n" if excluded else ""), encoding="utf-8")
    for error in errors:
        print(error, file=sys.stderr)
    sys.exit(1)

(package_root / "release-file-list.txt").write_text("\n".join(included) + "\n", encoding="utf-8")
(package_root / "release-excluded-list.txt").write_text("\n".join(excluded) + ("\n" if excluded else ""), encoding="utf-8")
list_file.unlink(missing_ok=True)
untracked_file.unlink(missing_ok=True)
PY
}

write_manifest() {
  local manifest_path="$1"
  local release_label="$2"
  local commit="$3"
  local source_status="$4"
  local archive_path="$5"
  local archive_sha="$6"
  local generated_at="$7"
  python3 - "$manifest_path" "$release_label" "$commit" "$source_status" "$archive_path" "$archive_sha" "$generated_at" <<'PY'
import json
import os
import sys

manifest_path, release_label, commit, source_status, archive_path, archive_sha, generated_at = sys.argv[1:8]
archive_name = os.path.basename(archive_path) if archive_path else None
data = {
    "release_label": release_label,
    "commit": commit,
    "generated_at": generated_at,
    "package_archive": archive_name,
    "package_sha256": archive_sha or None,
    "source_status": source_status,
    "included_roots": [
        "AGENTS.md",
        ".gitignore",
        ".dockerignore",
        ".env.example",
        ".github/workflows/",
        "Dockerfile",
        "cmd/",
        "configs/config.docker.yaml",
        "configs/config.example.yaml",
        "docker-compose.yml",
        "docker/",
        "docs/",
        "examples/",
        "internal/",
        "openspec/",
        "pkg/",
        "scripts/",
        "web/",
        "go.mod",
        "go.sum",
    ],
    "excluded_patterns": [
        ".git/",
        ".cursor/",
        "tmp/",
        "cmd/agent/tmp/",
        "docs/release/ui-audit-assets/",
        "configs/config.local.yaml",
        "web/node_modules/",
        "web/dist/",
        "playwright-report/",
        "test-results/",
        "*.db",
        "*.sqlite",
        "*.sqlite3",
        "*.log",
        "*.trace",
        "raw provider payloads",
        "complete prompt payload files",
        "complete API keys",
        "private local paths",
    ],
    "package_metadata_entries": [
        "release-file-list.txt",
        "release-excluded-list.txt",
    ],
    "verification_commands": [
        "bash scripts/local-release-package.sh --verify <archive>",
        "openspec validate --all --strict",
        "git diff --check",
        "go test $(bash scripts/go-packages.sh)",
        "npm --prefix web test",
        "npm --prefix web run build",
        "bash scripts/e2e-smoke.sh",
    ],
    "acceptance_references": [
        "docs/release/acceptance/2026-06-22-p91-github-release-docker-deployment.md",
        "docs/release/acceptance/2026-06-22-p92-final-original-requirement-audit-summary.md",
        "docs/release/acceptance/2026-06-22-p92-final-original-requirement-audit-ledger.md",
        "docs/release/acceptance/2026-06-22-p93-final-code-reality-design-audit.md",
    ],
    "known_degradations": [],
    "not_claimed": [
        "future public-source availability",
        "future model-provider availability",
        "investment returns",
        "broker connectivity",
        "automatic trading",
        "one-click trading",
        "order delegation",
        "external push",
        "automatic confirmation",
        "automatic rule application",
        "automatic repair",
        "automatic upgrade",
        "automatic migration",
        "automatic restore",
        "real database overwrite",
        "login sources",
        "paid sources",
        "authorized sources",
        "Level2 data",
        "high-frequency data",
    ],
    "safety_note": "This package is a local source handoff artifact. It does not execute upgrades, migrations, repairs, restores, trades, external pushes, confirmations, rule applications, public-source calls, or LLM calls.",
}
with open(manifest_path, "w", encoding="utf-8") as handle:
    json.dump(data, handle, ensure_ascii=False, indent=2)
    handle.write("\n")
PY
}

verify_archive() {
  local archive_path="$1"
  if [[ ! -f "$archive_path" ]]; then
    echo "Archive not found: $archive_path" >&2
    exit 1
  fi
  local archive_abs
  archive_abs="$(cd "$(dirname "$archive_path")" && pwd)/$(basename "$archive_path")"
  local archive_dir
  archive_dir="$(dirname "$archive_abs")"
  local manifest_path="${archive_dir}/release-manifest.json"
  if [[ ! -f "$manifest_path" ]]; then
    echo "Manifest not found next to archive: $manifest_path" >&2
    exit 1
  fi
  local verify_dir="${OUTPUT_DIR}/${TIMESTAMP}-verify"
  local listing_file="${verify_dir}/archive-listing.txt"
  local summary_file="${verify_dir}/verify-summary.json"
  mkdir -p "$verify_dir"
  tar -tzf "$archive_abs" >"$listing_file"
  python3 - "$manifest_path" "$archive_abs" "$listing_file" "$summary_file" "$(sha256_file "$archive_abs")" <<'PY'
import json
import re
import sys
from pathlib import Path

manifest_path, archive_path, listing_path, summary_path, actual_sha = sys.argv[1:6]
manifest_text = Path(manifest_path).read_text(encoding="utf-8")
manifest = json.loads(manifest_text)
listing = Path(listing_path).read_text(encoding="utf-8").splitlines()
archive_name = Path(archive_path).name

errors = []
warnings = []

if manifest.get("package_archive") != archive_name:
    errors.append("manifest package_archive does not match archive basename")
if manifest.get("package_sha256") != actual_sha:
    errors.append("manifest package_sha256 does not match archive checksum")

required_suffixes = [
    "/AGENTS.md",
    "/.dockerignore",
    "/.env.example",
    "/Dockerfile",
    "/cmd/",
    "/configs/config.docker.yaml",
    "/configs/config.example.yaml",
    "/docker-compose.yml",
    "/docker/",
    "/docs/",
    "/internal/",
    "/openspec/",
    "/scripts/",
    "/web/",
    "/go.mod",
    "/go.sum",
    "/release-manifest.json",
]

def has_entry(suffix: str) -> bool:
    if suffix.endswith("/"):
        needle = suffix
        return any(item.endswith(needle) or f"{needle}" in item for item in listing)
    return any(item.endswith(suffix) for item in listing)

for suffix in required_suffixes:
    if not has_entry(suffix):
        errors.append(f"missing required package entry: {suffix}")

forbidden_listing = [
    r"(^|/)\.git(/|$)",
    r"(^|/)\.cursor(/|$)",
    r"(^|/)tmp(/|$)",
    r"(^|/)cmd/agent/tmp(/|$)",
    r"(^|/)docs/release/ui-audit-assets(/|$)",
    r"(^|/)configs/config\.local\.yaml$",
    r"(^|/)web/node_modules(/|$)",
    r"(^|/)web/dist(/|$)",
    r"(^|/)playwright-report(/|$)",
    r"(^|/)test-results(/|$)",
    r"\.(db|sqlite|sqlite3|log|trace)$",
]
for item in listing:
    for pattern in forbidden_listing:
        if re.search(pattern, item, re.IGNORECASE):
            errors.append(f"forbidden archive entry: {item}")

forbidden_text = [
    r"(?<![A-Za-z0-9])sk-[A-Za-z0-9_-]{20,}",
    r"/Users/(?!private\b)[A-Za-z0-9._-]+/",
    r"BEGIN (RSA|OPENSSH|PRIVATE) KEY",
    r"Authorization:\s*Bearer",
    r"(?:^|[,{]\s*)[\"']?prompt[\"']?\s*:\s*[\"'][^\"'\n]{16,}",
]
for pattern in forbidden_text:
    if re.search(pattern, manifest_text, re.IGNORECASE):
        errors.append(f"forbidden manifest text pattern: {pattern}")

summary = {
    "archive": archive_name,
    "manifest": str(Path(manifest_path).name),
    "package_sha256": actual_sha,
    "archive_entry_count": len(listing),
    "status": "passed" if not errors else "failed",
    "errors": errors,
    "warnings": warnings,
    "safety_note": "Verification only inspects the local archive and manifest; it does not execute runtime tasks, migrations, repairs, restores, trades, external pushes, public-source calls, or LLM calls.",
}
Path(summary_path).write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
print(json.dumps(summary, ensure_ascii=False, indent=2))
if errors:
    sys.exit(1)
PY
  echo "Package verification summary: ${summary_file//$ROOT_DIR/<repo>}"
}

OUTPUT_DIR="$(normalize_output_dir)"

if [[ -n "$VERIFY_ARCHIVE" ]]; then
  verify_archive "$VERIFY_ARCHIVE"
  exit 0
fi

if [[ -z "$RELEASE_LABEL" ]]; then
  echo "--release-label is required when creating a package" >&2
  usage >&2
  exit 1
fi

if [[ "$SKIP_BUILD" != "1" ]]; then
  npm --prefix "$ROOT_DIR/web" run build
fi

SAFE_LABEL="$(safe_label "$RELEASE_LABEL")"
RUN_DIR="${OUTPUT_DIR}/${TIMESTAMP}"
PACKAGE_ROOT_NAME="investment-agent-${SAFE_LABEL}"
STAGE_DIR="${RUN_DIR}/stage"
PACKAGE_ROOT="${STAGE_DIR}/${PACKAGE_ROOT_NAME}"
ARCHIVE_PATH="${RUN_DIR}/${PACKAGE_ROOT_NAME}.tar.gz"
MANIFEST_PATH="${RUN_DIR}/release-manifest.json"
INTERNAL_MANIFEST_PATH="${PACKAGE_ROOT}/release-manifest.json"
SUMMARY_PATH="${RUN_DIR}/package-summary.json"
COMMIT="$(git -C "$ROOT_DIR" rev-parse HEAD)"
if [[ -n "$(git -C "$ROOT_DIR" status --short)" ]]; then
  SOURCE_STATUS="dirty"
else
  SOURCE_STATUS="clean"
fi

mkdir -p "$PACKAGE_ROOT"
copy_release_files "$PACKAGE_ROOT"
write_manifest "$INTERNAL_MANIFEST_PATH" "$RELEASE_LABEL" "$COMMIT" "$SOURCE_STATUS" "" "" "$TIMESTAMP"
tar -czf "$ARCHIVE_PATH" -C "$STAGE_DIR" "$PACKAGE_ROOT_NAME"
ARCHIVE_SHA="$(sha256_file "$ARCHIVE_PATH")"
write_manifest "$MANIFEST_PATH" "$RELEASE_LABEL" "$COMMIT" "$SOURCE_STATUS" "$ARCHIVE_PATH" "$ARCHIVE_SHA" "$TIMESTAMP"

python3 - "$SUMMARY_PATH" "$TIMESTAMP" "$RUN_DIR" "$ARCHIVE_PATH" "$MANIFEST_PATH" "$ARCHIVE_SHA" "$ROOT_DIR" <<'PY'
import json
import sys
from pathlib import Path

summary_path, timestamp, run_dir, archive_path, manifest_path, archive_sha, root_dir = sys.argv[1:8]

def safe(value: str) -> str:
    return value.replace(root_dir, "<repo>")

summary = {
    "generated_at": timestamp,
    "generated_dir": safe(run_dir),
    "archive": safe(archive_path),
    "manifest": safe(manifest_path),
    "package_sha256": archive_sha,
    "verify_command": f"bash scripts/local-release-package.sh --verify {safe(archive_path)}",
    "safety_note": "Package creation stages source files only; it does not execute upgrades, migrations, repairs, restores, trades, external pushes, confirmations, rule applications, public-source calls, or LLM calls.",
}
Path(summary_path).write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
print(json.dumps(summary, ensure_ascii=False, indent=2))
PY

echo "Release package created: ${ARCHIVE_PATH//$ROOT_DIR/<repo>}"
echo "Release manifest created: ${MANIFEST_PATH//$ROOT_DIR/<repo>}"
