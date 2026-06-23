#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_PATH="${ROOT_DIR}/configs/config.example.yaml"
OUTPUT_DIR="${ROOT_DIR}/tmp/local-install-diagnostics"
RUN_RECOVERY="1"
RUN_E2E="1"
RUN_RELEASE_UPGRADE="0"
TARGET_VERSION=""
TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --config)
      CONFIG_PATH="$2"
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    --skip-recovery)
      RUN_RECOVERY="0"
      shift
      ;;
    --skip-e2e)
      RUN_E2E="0"
      shift
      ;;
    --include-release-upgrade)
      RUN_RELEASE_UPGRADE="1"
      shift
      ;;
    --target-version)
      TARGET_VERSION="$2"
      shift 2
      ;;
    -h|--help)
      cat <<EOF
Usage: bash scripts/local-install-diagnostics.sh [options]

Options:
  --config PATH         使用的配置文件路径（默认: configs/config.example.yaml）
  --output-dir PATH     输出目录（默认: tmp/local-install-diagnostics）
  --skip-recovery       跳过 scripts/recovery-smoke.sh
  --skip-e2e            跳过 scripts/e2e-smoke.sh
  --include-release-upgrade
                        显式纳入 P49 release/upgrade 检查
  --target-version VALUE
                        P49 release/upgrade 检查目标版本或 release label
  -h, --help           显示本帮助
EOF
      exit 0
      ;;
    *)
      echo "未知参数: $1" >&2
      exit 1
      ;;
  esac
done

RUN_DIR="${OUTPUT_DIR}/${TIMESTAMP}"
PRE_FLIGHT_FILE="${RUN_DIR}/preflight.json"
RELEASE_UPGRADE_FILE="${RUN_DIR}/release-upgrade.json"
SUMMARY_FILE="${RUN_DIR}/install-summary.json"
STEP_RECORD_FILE="${RUN_DIR}/steps.jsonl"

if [[ ! -f "$CONFIG_PATH" ]]; then
  echo "配置文件不存在: $CONFIG_PATH" >&2
  exit 1
fi

mkdir -p "$RUN_DIR"

sanitize_json_file() {
  local path="$1"
  if [[ ! -f "$path" ]]; then
    return 0
  fi
  python3 - "$path" "$ROOT_DIR" "${HOME:-}" <<'PY'
import json
import sys

path, root_dir, home_dir = sys.argv[1:4]

def scrub(value):
    if isinstance(value, str):
        value = value.replace(root_dir, "<repo>")
        if home_dir:
            value = value.replace(home_dir, "<home>")
        return value
    if isinstance(value, list):
        return [scrub(item) for item in value]
    if isinstance(value, dict):
        return {key: scrub(item) for key, item in value.items()}
    return value

with open(path, encoding="utf-8") as handle:
    data = json.load(handle)
with open(path, "w", encoding="utf-8") as handle:
    json.dump(scrub(data), handle, ensure_ascii=False, indent=2)
    handle.write("\n")
PY
}

run_step() {
  local name="$1"
  local artifact="$2"
  local skip_if="$3"
  local display_command="$4"
  local log_file="$5"
  local status="skipped"
  local exit_code="0"
  shift 5

  if [[ "$skip_if" == "skip" ]]; then
    status="skipped"
    command="-"
    exit_code=""
    local log_line
    log_line="$(printf '%s\036%s\036%s\036%s\036%s\n' \
      "$name" \
      "$status" \
      "" \
      "$artifact" \
      "-")"
    printf '%s\n' "$log_line" >>"$STEP_RECORD_FILE"
    printf '[%s] skip: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$name"
    return 0
  fi

  if [[ -z "$log_file" ]]; then
    log_file="$RUN_DIR/${name}.log"
  fi
  set +e
  "$@" >"$log_file" 2>&1
  exit_code=$?
  set -e

  if [[ $exit_code -eq 0 || "$exit_code" == "0" ]]; then
    status="pass"
  else
    status="failed"
  fi

  local log_line
  log_line="$(printf '%s\036%s\036%s\036%s\036%s\n' \
    "$name" \
    "$status" \
    "$exit_code" \
    "$artifact" \
    "$display_command")"
  printf '%s\n' "$log_line" >>"$STEP_RECORD_FILE"
  printf '[%s] %s -> %s (code=%s)\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$name" "$status" "$exit_code"
  return "$exit_code"
}

preflight_status=0
if run_step "preflight" \
  "$PRE_FLIGHT_FILE" \
  "" \
  "INVESTMENT_AGENT_CONFIG=\"${CONFIG_PATH//$ROOT_DIR/<repo>}\" go run ./cmd/agent --preflight --diagnostics \"${PRE_FLIGHT_FILE//$ROOT_DIR/<repo>}\"" \
  "$RUN_DIR/preflight.log" \
  env "INVESTMENT_AGENT_CONFIG=$CONFIG_PATH" go run ./cmd/agent --preflight --diagnostics "$PRE_FLIGHT_FILE"; then
  preflight_status=0
else
  preflight_status=$?
fi
sanitize_json_file "$PRE_FLIGHT_FILE"

if [[ "$RUN_RECOVERY" == "1" ]]; then
  run_step "recovery_smoke" "$RUN_DIR/recovery_smoke.log" "" "bash scripts/recovery-smoke.sh" "$RUN_DIR/recovery_smoke.log" bash scripts/recovery-smoke.sh
else
  run_step "recovery_smoke" "" "skip" "-" ""
fi

if [[ "$RUN_RELEASE_UPGRADE" == "1" ]]; then
  release_args=(env "INVESTMENT_AGENT_CONFIG=$CONFIG_PATH" go run ./cmd/agent --release-upgrade-check --diagnostics "$RELEASE_UPGRADE_FILE")
  release_display="INVESTMENT_AGENT_CONFIG=\"${CONFIG_PATH//$ROOT_DIR/<repo>}\" go run ./cmd/agent --release-upgrade-check --diagnostics \"${RELEASE_UPGRADE_FILE//$ROOT_DIR/<repo>}\""
  if [[ -n "$TARGET_VERSION" ]]; then
    release_args+=(--target-version "$TARGET_VERSION")
    release_display="${release_display} --target-version \"<provided>\""
  fi
  run_step "release_upgrade_check" "$RELEASE_UPGRADE_FILE" "" "$release_display" "$RUN_DIR/release_upgrade_check.log" "${release_args[@]}"
  sanitize_json_file "$RELEASE_UPGRADE_FILE"
else
  run_step "release_upgrade_check" "" "skip" "-" ""
fi

if [[ "$RUN_E2E" == "1" ]]; then
  run_step "e2e_smoke" "$RUN_DIR/e2e_smoke.log" "" "bash scripts/e2e-smoke.sh" "$RUN_DIR/e2e_smoke.log" bash scripts/e2e-smoke.sh
else
  run_step "e2e_smoke" "" "skip" "-" ""
fi

python3 - "$STEP_RECORD_FILE" "$SUMMARY_FILE" "$TIMESTAMP" "$CONFIG_PATH" "$RUN_DIR" "$PRE_FLIGHT_FILE" "$RELEASE_UPGRADE_FILE" "$ROOT_DIR" <<'PY'
import json
import sys

step_record_path, summary_path, timestamp, config_path, run_dir, preflight_path, release_upgrade_path, root_dir = sys.argv[1:9]
steps = []

def safe(value):
    if not value:
        return None
    return str(value).replace(root_dir, "<repo>")

with open(step_record_path, encoding="utf-8") as handle:
    for raw in handle:
        line = raw.strip("\n")
        if not line:
            continue
        name, status, exit_code, artifact, command = line.split("\036", 4)
        exit_code = None if exit_code == "" else int(exit_code)
        steps.append(
            {
                "name": name,
                "status": status,
                "exit_code": int(exit_code) if exit_code is not None else None,
                "command": command,
                "artifact": safe(artifact) if artifact else None,
            }
        )

summary = {
    "generated_at": timestamp,
    "generated_dir": safe(run_dir),
    "config_path": safe(config_path),
    "preflight_diagnostics": safe(preflight_path),
    "release_upgrade_diagnostics": safe(release_upgrade_path),
    "steps": steps,
}

with open(summary_path, "w", encoding="utf-8") as handle:
    json.dump(summary, handle, ensure_ascii=False, indent=2)
    handle.write("\n")
PY

cat "$SUMMARY_FILE"

if [[ "$preflight_status" -ne 0 ]]; then
  echo "预检失败。请查看：${SUMMARY_FILE//$ROOT_DIR/<repo>}" >&2
  exit 1
fi

echo "诊断摘要已生成：${SUMMARY_FILE//$ROOT_DIR/<repo>}"
