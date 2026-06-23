#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CONFIG_PATH="${ROOT_DIR}/configs/config.example.yaml"
OUTPUT_DIR="${ROOT_DIR}/tmp/local-release-upgrade"
TARGET_VERSION=""
RUN_PREFLIGHT="1"
TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --config)
      CONFIG_PATH="$2"
      shift 2
      ;;
    --target-version)
      TARGET_VERSION="$2"
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR="$2"
      shift 2
      ;;
    --skip-preflight)
      RUN_PREFLIGHT="0"
      shift
      ;;
    -h|--help)
      cat <<EOF
Usage: bash scripts/local-release-upgrade-check.sh [options]

Options:
  --config PATH           使用的配置文件路径（默认: configs/config.example.yaml）
  --target-version VALUE  目标版本或 release label
  --output-dir PATH       输出目录（默认: tmp/local-release-upgrade）
  --skip-preflight        跳过 cmd/agent --preflight
  -h, --help              显示本帮助
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
PREFLIGHT_FILE="${RUN_DIR}/preflight.json"
REPORT_FILE="${RUN_DIR}/release-upgrade.json"
SUMMARY_FILE="${RUN_DIR}/release-upgrade-summary.json"
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

record_step() {
  local name="$1"
  local status="$2"
  local exit_code="$3"
  local artifact="$4"
  local display_command="$5"
  printf '%s\036%s\036%s\036%s\036%s\n' "$name" "$status" "$exit_code" "$artifact" "$display_command" >>"$STEP_RECORD_FILE"
}

run_step() {
  local name="$1"
  local artifact="$2"
  local log_file="$3"
  local display_command="$4"
  local exit_code="0"
  local status="pass"
  shift 4

  set +e
  "$@" >"$log_file" 2>&1
  exit_code=$?
  set -e

  if [[ "$exit_code" -ne 0 ]]; then
    status="failed"
  fi
  record_step "$name" "$status" "$exit_code" "$artifact" "$display_command"
  printf '[%s] %s -> %s (code=%s)\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" "$name" "$status" "$exit_code"
  return "$exit_code"
}

if [[ "$RUN_PREFLIGHT" == "1" ]]; then
  run_step "preflight" \
    "$PREFLIGHT_FILE" \
    "$RUN_DIR/preflight.log" \
    "INVESTMENT_AGENT_CONFIG=\"${CONFIG_PATH//$ROOT_DIR/<repo>}\" go run ./cmd/agent --preflight --diagnostics \"${PREFLIGHT_FILE//$ROOT_DIR/<repo>}\"" \
    env "INVESTMENT_AGENT_CONFIG=$CONFIG_PATH" go run ./cmd/agent --preflight --diagnostics "$PREFLIGHT_FILE"
  sanitize_json_file "$PREFLIGHT_FILE"
else
  record_step "preflight" "skipped" "" "" "-"
  printf '[%s] skip: preflight\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
fi

release_args=(env "INVESTMENT_AGENT_CONFIG=$CONFIG_PATH" go run ./cmd/agent --release-upgrade-check --diagnostics "$REPORT_FILE")
release_display="INVESTMENT_AGENT_CONFIG=\"${CONFIG_PATH//$ROOT_DIR/<repo>}\" go run ./cmd/agent --release-upgrade-check --diagnostics \"${REPORT_FILE//$ROOT_DIR/<repo>}\""
if [[ -n "$TARGET_VERSION" ]]; then
  release_args+=(--target-version "$TARGET_VERSION")
  release_display="${release_display} --target-version \"<provided>\""
fi
run_step "release_upgrade_check" "$REPORT_FILE" "$RUN_DIR/release-upgrade.log" "$release_display" "${release_args[@]}"
sanitize_json_file "$REPORT_FILE"

python3 - "$STEP_RECORD_FILE" "$SUMMARY_FILE" "$TIMESTAMP" "$RUN_DIR" "$REPORT_FILE" "$ROOT_DIR" <<'PY'
import json
import sys

step_record_path, summary_path, timestamp, run_dir, report_path, root_dir = sys.argv[1:7]
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
        steps.append(
            {
                "name": name,
                "status": status,
                "exit_code": None if exit_code == "" else int(exit_code),
                "command": command,
                "artifact": safe(artifact) if artifact else None,
            }
        )

summary = {
    "generated_at": timestamp,
    "generated_dir": safe(run_dir),
    "release_upgrade_report": safe(report_path),
    "steps": steps,
    "safety_note": "本脚本只运行本地预检和发布升级检查；不会执行升级、迁移、备份、恢复、交易、外部推送、自动确认、自动应用规则或自动修复。",
}

with open(summary_path, "w", encoding="utf-8") as handle:
    json.dump(summary, handle, ensure_ascii=False, indent=2)
    handle.write("\n")
PY

cat "$SUMMARY_FILE"
echo "发布升级检查摘要已生成：${SUMMARY_FILE//$ROOT_DIR/<repo>}"
