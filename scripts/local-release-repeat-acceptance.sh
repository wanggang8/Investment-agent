#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT_DIR="${ROOT_DIR}/tmp/local-release-repeat"
ARCHIVE_PATH=""
SKIP_INSTALL="0"
SKIP_E2E="0"
TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
PACKAGE_ROOT=""
RUN_DIR=""
WORKSPACE_DIR=""
LOG_DIR=""
COMMANDS_TSV=""

usage() {
  cat <<EOF
Usage: bash scripts/local-release-repeat-acceptance.sh [options]

Options:
  --archive PATH     Release package archive to repeat from
  --output-dir PATH  Output directory (default: tmp/local-release-repeat)
  --skip-install     Skip npm ci in the extracted package
  --skip-e2e         Skip E2E smoke; diagnostic reruns only, not P65 main acceptance
  -h, --help         Show this help
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --archive)
      ARCHIVE_PATH="$2"
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    --skip-install)
      SKIP_INSTALL="1"
      shift
      ;;
    --skip-e2e)
      SKIP_E2E="1"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

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

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

record_step() {
  local name="$1"
  local command="$2"
  local status="$3"
  local duration="$4"
  local log_path="$5"
  printf '%s\t%s\t%s\t%s\t%s\n' "$name" "$command" "$status" "$duration" "$log_path" >>"$COMMANDS_TSV"
}

run_step() {
  local name="$1"
  shift
  local log_path="${LOG_DIR}/${name}.log"
  local command_text="$*"
  local started
  local ended
  local duration
  local status
  started="$(date +%s)"
  set +e
  (cd "$PACKAGE_ROOT" && "$@") >"$log_path" 2>&1
  status="$?"
  set -e
  ended="$(date +%s)"
  duration="$((ended - started))"
  record_step "$name" "$command_text" "$status" "$duration" "$log_path"
  return "$status"
}

write_summary() {
  local status="$1"
  local summary_path="${RUN_DIR}/repeat-summary.json"
  local package_sha
  package_sha="$(sha256_file "$ARCHIVE_ABS")"
  python3 - \
    "$summary_path" \
    "$COMMANDS_TSV" \
    "$ROOT_DIR" \
    "$RUN_DIR" \
    "$WORKSPACE_DIR" \
    "$PACKAGE_ROOT" \
    "$ARCHIVE_ABS" \
    "$MANIFEST_PATH" \
    "$package_sha" \
    "$status" \
    "$SKIP_INSTALL" \
    "$SKIP_E2E" \
    "$TIMESTAMP" <<'PY'
import json
import sys
from pathlib import Path

(
    summary_path,
    commands_tsv,
    root_dir,
    run_dir,
    workspace_dir,
    package_root,
    archive_path,
    manifest_path,
    package_sha,
    status,
    skip_install,
    skip_e2e,
    timestamp,
) = sys.argv[1:14]

root_dir = str(Path(root_dir).resolve())
run_dir = str(Path(run_dir).resolve())
workspace_dir = str(Path(workspace_dir).resolve())
package_root = str(Path(package_root).resolve())

def safe(value: str) -> str:
    return (
        str(value)
        .replace(package_root, "<package-root>")
        .replace(workspace_dir, "<repeat-workspace>")
        .replace(run_dir, "<repeat-run>")
        .replace(root_dir, "<repo>")
    )

manifest = json.loads(Path(manifest_path).read_text(encoding="utf-8"))
commands = []
if Path(commands_tsv).exists():
    for raw in Path(commands_tsv).read_text(encoding="utf-8").splitlines():
        if not raw:
            continue
        name, command, exit_code, duration_seconds, log_path = raw.split("\t", 4)
        commands.append(
            {
                "name": name,
                "command": command,
                "exit_code": int(exit_code) if exit_code.isdigit() else exit_code,
                "status": "passed" if exit_code == "0" else ("skipped" if exit_code == "SKIPPED" else "failed"),
                "duration_seconds": int(duration_seconds) if duration_seconds.isdigit() else duration_seconds,
                "log": safe(log_path),
            }
        )

summary = {
    "generated_at": timestamp,
    "status": status,
    "archive": Path(archive_path).name,
    "manifest": Path(manifest_path).name,
    "package_sha256": package_sha,
    "release_label": manifest.get("release_label"),
    "source_commit": manifest.get("commit"),
    "source_status": manifest.get("source_status"),
    "repeat_run_dir": safe(run_dir),
    "repeat_workspace": safe(workspace_dir),
    "package_root": safe(package_root),
    "skip_install": skip_install == "1",
    "skip_e2e": skip_e2e == "1",
    "commands": commands,
    "known_caveats": [
        "This is a cross-machine-equivalent local isolated repeat, not proof that a separate physical machine has run the package.",
        "Public-source and model-provider availability are not guaranteed by this repeat.",
    ],
    "not_claimed": [
        "remote publishing",
        "Git tag creation",
        "installer signing",
        "automatic upgrade",
        "automatic migration",
        "automatic restore",
        "automatic repair",
        "real database overwrite",
        "broker connectivity",
        "automatic trading",
        "one-click trading",
        "order delegation",
        "external push",
        "automatic confirmation",
        "automatic rule application",
        "login-gated sources",
        "paid sources",
        "authorization-gated sources",
        "Level2 data",
        "high-frequency data",
        "future provider availability",
        "investment returns",
    ],
    "safety_note": "Repeat acceptance verifies the package from an isolated local workspace. It does not publish, tag, upgrade, migrate, restore, repair, overwrite databases, trade, push notifications, call public providers, call LLM providers, confirm actions, or apply rules automatically.",
}
Path(summary_path).write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
print(json.dumps(summary, ensure_ascii=False, indent=2))
PY
}

if [[ -z "$ARCHIVE_PATH" ]]; then
  echo "--archive is required" >&2
  usage >&2
  exit 1
fi
if [[ ! -f "$ARCHIVE_PATH" ]]; then
  echo "Archive not found: $ARCHIVE_PATH" >&2
  exit 1
fi

OUTPUT_DIR="$(normalize_output_dir)"
ARCHIVE_ABS="$(cd "$(dirname "$ARCHIVE_PATH")" && pwd)/$(basename "$ARCHIVE_PATH")"
MANIFEST_PATH="$(dirname "$ARCHIVE_ABS")/release-manifest.json"
if [[ ! -f "$MANIFEST_PATH" ]]; then
  echo "Manifest not found next to archive: $MANIFEST_PATH" >&2
  exit 1
fi

RUN_DIR="${OUTPUT_DIR}/${TIMESTAMP}"
WORKSPACE_DIR="${RUN_DIR}/workspace"
LOG_DIR="${RUN_DIR}/logs"
COMMANDS_TSV="${RUN_DIR}/commands.tsv"
mkdir -p "$WORKSPACE_DIR" "$LOG_DIR"
: >"$COMMANDS_TSV"

bash "$ROOT_DIR/scripts/local-release-package.sh" --verify "$ARCHIVE_ABS" --output-dir "$OUTPUT_DIR/package-verify"
tar -xzf "$ARCHIVE_ABS" -C "$WORKSPACE_DIR"

PACKAGE_ROOT_NAME="$(tar -tzf "$ARCHIVE_ABS" | awk -F/ 'NF > 1 {print $1; exit}')"
if [[ -z "$PACKAGE_ROOT_NAME" ]]; then
  echo "Unable to detect package root from archive" >&2
  write_summary "failed"
  exit 1
fi
PACKAGE_ROOT="${WORKSPACE_DIR}/${PACKAGE_ROOT_NAME}"
if [[ ! -d "$PACKAGE_ROOT" ]]; then
  echo "Extracted package root not found: $PACKAGE_ROOT" >&2
  write_summary "failed"
  exit 1
fi

overall_status="passed"

run_step "openspec-validate" "openspec" "validate" "--all" "--strict" || overall_status="failed"
run_step "go-test" "go" "test" "./..." || overall_status="failed"

if [[ "$SKIP_INSTALL" == "1" ]]; then
  record_step "npm-ci" "npm --prefix web ci" "SKIPPED" "0" "${LOG_DIR}/npm-ci.log"
else
  run_step "npm-ci" "npm" "--prefix" "web" "ci" || overall_status="failed"
fi

run_step "npm-test" "npm" "--prefix" "web" "test" || overall_status="failed"
run_step "npm-build" "npm" "--prefix" "web" "run" "build" || overall_status="failed"

if [[ "$SKIP_E2E" == "1" ]]; then
  record_step "e2e-smoke" "bash scripts/e2e-smoke.sh" "SKIPPED" "0" "${LOG_DIR}/e2e-smoke.log"
else
  run_step "e2e-smoke" "env" "E2E_SERVER_PORT=18165" "E2E_WEB_PORT=14265" "bash" "scripts/e2e-smoke.sh" || overall_status="failed"
fi

write_summary "$overall_status"

if [[ "$overall_status" != "passed" ]]; then
  exit 1
fi
