#!/usr/bin/env sh
set -eu

curl -fsS "http://127.0.0.1:8080/api/v1/health" >/dev/null
curl -fsS "http://127.0.0.1:4173/" >/dev/null
