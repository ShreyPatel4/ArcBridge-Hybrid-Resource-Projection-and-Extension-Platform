#!/usr/bin/env bash
set -euo pipefail

HOST=${1:-localhost}
PORT=${2:-6379}

echo "PING" | nc -w 2 "$HOST" "$PORT" || {
  echo "redis ping failed" >&2
  exit 1
}

echo "redis ping sent"
