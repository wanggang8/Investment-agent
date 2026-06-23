#!/usr/bin/env python3
"""P91 deployment artifact and safety checker."""

from __future__ import annotations

import re
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
REQUIRED_FILES = [
    ".dockerignore",
    ".env.example",
    ".github/workflows/ci.yml",
    ".github/workflows/release.yml",
    "Dockerfile",
    "docker-compose.yml",
    "docker/entrypoint.sh",
    "docker/healthcheck.sh",
    "configs/config.docker.yaml",
    "docs/deployment.md",
    "scripts/backup.sh",
    "scripts/deploy-lib.sh",
    "scripts/doctor.sh",
    "scripts/install.sh",
    "scripts/status.sh",
    "scripts/uninstall.sh",
    "scripts/upgrade.sh",
]


def require(condition: bool, reason: str) -> None:
    if not condition:
        raise SystemExit(f"status=failed\nreason={reason}")


def text(path: str) -> str:
    p = ROOT / path
    require(p.exists(), f"missing:{path}")
    return p.read_text(encoding="utf-8")


def main() -> None:
    for path in REQUIRED_FILES:
        require((ROOT / path).exists(), f"missing:{path}")

    env_example = text(".env.example")
    for key in ("DEEPSEEK_API_KEY=", "DEEPSEEK_BASE_URL=", "DEEPSEEK_MODEL=", "INVESTMENT_AGENT_DATA_DIR="):
        require(key in env_example, f"env_example_missing:{key}")
    require("sk-" not in env_example.lower(), "env_example_must_not_embed_key")

    compose = text("docker-compose.yml")
    for token in ("investment-agent", "env_file:", ".env", "healthcheck:", "INVESTMENT_AGENT_SQLITE_PATH", "INVESTMENT_AGENT_VECLITE_PATH"):
        require(token in compose, f"compose_missing:{token}")
    require("restart: unless-stopped" in compose, "compose_restart_policy")
    require("127.0.0.1:${INVESTMENT_AGENT_WEB_PORT" in compose, "compose_web_port_must_bind_localhost")
    require("127.0.0.1:${INVESTMENT_AGENT_SERVER_PORT" in compose, "compose_server_port_must_bind_localhost")

    dockerfile = text("Dockerfile")
    for token in ("npm", "go build", "cmd/server", "web/dist", "docker/entrypoint.sh"):
        require(token in dockerfile, f"dockerfile_missing:{token}")
    require("DEEPSEEK_API_KEY" not in dockerfile, "dockerfile_must_not_embed_key")

    install = text("scripts/install.sh")
    for token in ("detect_install_mode", "first_install", "upgrade", "scripts/upgrade.sh", "compose_cmd", "up -d", "wait_for_health"):
        require(token in install, f"install_missing:{token}")
    require("rm -rf" not in install, "install_must_not_delete_data")

    upgrade = text("scripts/upgrade.sh")
    for token in ("scripts/backup.sh", "compose_cmd", "up -d", "wait_for_health", "write_release_state"):
        require(token in upgrade, f"upgrade_missing:{token}")

    uninstall = text("scripts/uninstall.sh")
    for token in ("--purge", "DELETE INVESTMENT AGENT DATA", "compose_cmd", "down"):
        require(token in uninstall, f"uninstall_missing:{token}")
    require(re.search(r'CONFIRMATION_PHRASE=.DELETE INVESTMENT AGENT DATA.', uninstall), "uninstall_confirmation_phrase")
    require("rm -rf" in uninstall, "uninstall_purge_must_be_explicit")
    require("if [[ \"$PURGE\" == \"1\" ]]" in uninstall, "uninstall_rm_guard")

    workflows = text(".github/workflows/ci.yml") + "\n" + text(".github/workflows/release.yml")
    for token in ("openspec validate --all --strict", "go test", "npm --prefix web test", "npm --prefix web run build", "p91_deployment_check.py", "local-release-package.sh"):
        require(token in workflows, f"workflow_missing:{token}")
    require("secrets.DEEPSEEK_API_KEY" not in workflows, "workflow_must_not_reference_runtime_secret")

    config = text("configs/config.docker.yaml")
    for token in ("/data/sqlite/investment-agent.db", "/data/veclite", "api_key: \"\"", "p89_structured_public"):
        require(token in config, f"docker_config_missing:{token}")

    docs = text("docs/deployment.md")
    for token in ("install.sh", "upgrade.sh", "uninstall.sh", "--purge", "DEEPSEEK_API_KEY"):
        require(token in docs, f"deployment_doc_missing:{token}")
    require("127.0.0.1" in docs and "reverse proxy" in docs, "deployment_doc_must_describe_localhost_binding")

    backup = text("scripts/backup.sh")
    require("env.redacted" in backup and "<redacted>" in backup, "backup_must_redact_runtime_secret")

    print("p91_deployment_check:status=passed")


if __name__ == "__main__":
    main()
