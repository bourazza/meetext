#!/bin/sh
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/apps/api"
"$ROOT/bin/api" > logs/meetext.log 2>&1 &
echo $! > "$ROOT/.pids/api.pid"
