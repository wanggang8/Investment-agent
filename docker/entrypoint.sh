#!/usr/bin/env sh
set -eu

mkdir -p /data/sqlite /data/veclite /logs

if [ -z "${DEEPSEEK_API_KEY:-}" ]; then
  echo "warning: DEEPSEEK_API_KEY is empty; LLM-backed analysis will degrade safely where required" >&2
fi

/usr/local/bin/investment-agent-server > /logs/server.log 2>&1 &
SERVER_PID="$!"

for _ in $(seq 1 120); do
  if curl -fsS "http://127.0.0.1:8080/api/v1/health" >/dev/null 2>&1; then
    break
  fi
  if ! kill -0 "$SERVER_PID" >/dev/null 2>&1; then
    echo "investment-agent server exited before health check passed" >&2
    cat /logs/server.log >&2 || true
    exit 1
  fi
  sleep 1
done

if ! curl -fsS "http://127.0.0.1:8080/api/v1/health" >/dev/null 2>&1; then
  echo "investment-agent server health check failed" >&2
  cat /logs/server.log >&2 || true
  exit 1
fi

exec nginx -g "daemon off;"
