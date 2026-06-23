#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ARTIFACT_DIR="${P86_ARTIFACT_DIR:-$ROOT_DIR/docs/release/ui-audit-assets/2026-06-22-p86-core-goal-knowledge-safety-final}"
SUMMARY_PATH="$ARTIFACT_DIR/p86-integrated-summary.json"
RUN_LOG="$ARTIFACT_DIR/p86-run.log"

mkdir -p "$ARTIFACT_DIR"
: >"$RUN_LOG"

run_step() {
  local name="$1"
  shift
  local log="$ARTIFACT_DIR/$name.log"
  {
    echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] START $name"
    "$@"
    echo "[$(date -u +%Y-%m-%dT%H:%M:%SZ)] PASS $name"
  } >"$log" 2>&1
  echo "$name=passed log=$log" | tee -a "$RUN_LOG"
}

cd "$ROOT_DIR"

if [[ "${P86_REUSE_EXISTING:-0}" != "1" ]]; then
  run_step inventory python3 scripts/p85_p87_p86_plan_inventory_check.py
  run_step p74-built-in-knowledge env P74_ARTIFACT_DIR="$ARTIFACT_DIR/p74-built-in-knowledge-data-readiness" bash scripts/p74-built-in-knowledge-data-readiness.sh
  run_step p81-dynamic-source-ui env P75_ARTIFACT_DIR="$ARTIFACT_DIR/p81-dynamic-source-field-coverage" bash scripts/p75-non-510300-real-ui-journey.sh
  run_step p81-dynamic-source-go go test -v ./cmd/agent -run TestRunNon510300DynamicAcceptanceBindsCollectorSourceHealthAuditAndReadiness -count=1
  run_step p82-sop-action env P75_FINAL_RULE_APPLY=1 P75_ARTIFACT_DIR="$ARTIFACT_DIR/p82-sop-action-ui-sqlite" bash scripts/p75-sop-failure-real-ui-acceptance.sh
  run_step p83-governance env P83_ARTIFACT_DIR="$ARTIFACT_DIR/p83-governance-traceability" bash scripts/p83-governance-traceability-acceptance.sh
  run_step p84-portfolio-confirmation env P84_ARTIFACT_DIR="$ARTIFACT_DIR/p84-portfolio-confirmation" bash scripts/p84-portfolio-confirmation-acceptance.sh
  run_step p85-expected-return env P85_ARTIFACT_DIR="$ARTIFACT_DIR/p85-expected-return-analysis" bash scripts/p85-expected-return-analysis-acceptance.sh
  run_step p87-portfolio-state env P87_ARTIFACT_DIR="$ARTIFACT_DIR/p87-portfolio-state-allocation-safety" bash scripts/p87-portfolio-state-allocation-acceptance.sh
fi

python3 - "$ARTIFACT_DIR" "$SUMMARY_PATH" <<'PY'
import json
import re
import sys
from pathlib import Path

artifact_dir = Path(sys.argv[1])
summary_path = Path(sys.argv[2])

def read_json(path: Path) -> dict:
    if not path.exists():
        return {"status": "missing", "path": str(path)}
    data = json.loads(path.read_text(encoding="utf-8"))
    return data if isinstance(data, dict) else {"status": "invalid", "path": str(path)}

def read_kv(path: Path) -> dict[str, str]:
    out: dict[str, str] = {}
    if not path.exists():
        return out
    text = path.read_text(encoding="utf-8", errors="replace")
    stripped = text.lstrip()
    if stripped.startswith("{"):
        data = json.loads(text)
        return data if isinstance(data, dict) else {}
    for line in text.splitlines():
        if "=" in line:
            key, value = line.split("=", 1)
            out[key.strip()] = value.strip()
    return out

def go_log_passed(path: Path) -> bool:
    if not path.exists():
        return False
    text = path.read_text(encoding="utf-8", errors="replace")
    return "FAIL" not in text and bool(re.search(r"(?m)^ok\s+", text))

checks: dict[str, object] = {}
checks["inventory"] = read_json(artifact_dir / "p86-inventory.json")
checks["p74_browser"] = read_json(artifact_dir / "p74-built-in-knowledge-data-readiness" / "browser-results.json")
checks["p74_api_log"] = {"status": "passed" if (artifact_dir / "p74-built-in-knowledge-data-readiness" / "api-readiness-check.log").exists() else "missing"}
checks["p81_browser"] = read_json(artifact_dir / "p81-dynamic-source-field-coverage" / "browser-results.json")
checks["p81_db"] = read_kv(artifact_dir / "p81-dynamic-source-field-coverage" / "db-impact-check.log")
checks["p81_go"] = {"status": "passed" if go_log_passed(artifact_dir / "p81-dynamic-source-go.log") else "failed"}
checks["p82_browser"] = read_json(artifact_dir / "p82-sop-action-ui-sqlite" / "browser-results.json")
checks["p82_db"] = read_kv(artifact_dir / "p82-sop-action-ui-sqlite" / "db-impact-check.log")
checks["p83_summary"] = read_json(artifact_dir / "p83-governance-traceability" / "governance-traceability-summary.json")
checks["p84_summary"] = read_json(artifact_dir / "p84-portfolio-confirmation" / "portfolio-confirmation-summary.json")
checks["p85_summary"] = read_json(artifact_dir / "p85-expected-return-analysis" / "expected-return-summary.json")
checks["p87_summary"] = read_json(artifact_dir / "p87-portfolio-state-allocation-safety" / "portfolio-state-allocation-summary.json")

passed = (
    checks["inventory"].get("status") == "passed"
    and checks["inventory"].get("remaining_full_release_required_non_real_pass_rows") == 137
    and checks["p74_browser"].get("status") == "passed"
    and checks["p74_api_log"].get("status") == "passed"
    and checks["p81_browser"].get("status") == "passed"
    and checks["p81_db"].get("status") == "passed"
    and checks["p81_go"].get("status") == "passed"
    and checks["p82_browser"].get("status") == "passed"
    and checks["p82_db"].get("status") == "passed"
    and checks["p83_summary"].get("status") == "passed"
    and checks["p84_summary"].get("status") == "passed"
    and checks["p85_summary"].get("status") == "passed"
    and checks["p87_summary"].get("status") == "passed"
)

payload = {
    "status": "passed" if passed else "failed",
    "artifact_dir": str(artifact_dir),
    "checks": checks,
    "claim_boundary": (
        "P86 integrated runner proves cumulative local UI/API/SQLite/workflow acceptance across P74/P81/P82/P83/P84/P85/P87. "
        "It does not by itself prove historical backtest accuracy, complete external data-source breadth, broker connectivity, automatic trading, automatic confirmation, or return promises."
    ),
}
summary_path.write_text(json.dumps(payload, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
print(f"p86_integrated_acceptance:status={payload['status']}:artifact={summary_path}")
if not passed:
    raise SystemExit(1)
PY
