#!/bin/sh
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/apps/api"
mkdir -p logs "$ROOT/.pids"
stdout_log="/tmp/meetext-api.stdout.log"
: > logs/meetext.log
: > "$stdout_log"
if command -v setsid >/dev/null 2>&1; then
  setsid "$ROOT/bin/api" > "$stdout_log" 2>&1 &
else
  "$ROOT/bin/api" > "$stdout_log" 2>&1 &
fi
pid=$!
echo "$pid" > "$ROOT/.pids/api.pid"

sleep 1
if ! kill -0 "$pid" 2>/dev/null; then
  echo "ERROR: API process exited during startup. Check $ROOT/apps/api/logs/meetext.log"
  tail -50 "$ROOT/apps/api/logs/meetext.log" 2>/dev/null || true
  if [ -s "$stdout_log" ]; then
    echo "---- process stdout/stderr ----"
    tail -50 "$stdout_log" 2>/dev/null || true
  fi
  exit 1
fi
