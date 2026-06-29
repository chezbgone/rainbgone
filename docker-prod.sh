#!/usr/bin/env bash
set -euo pipefail

docker compose -f compose.prod.yaml up --build "$@"
