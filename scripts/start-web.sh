#!/bin/sh
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/apps/web"
mkdir -p "$ROOT/.pids"
if command -v setsid >/dev/null 2>&1; then
  setsid npm run dev -- --port 3000 > /tmp/meetext-web.log 2>&1 &
else
  npm run dev -- --port 3000 > /tmp/meetext-web.log 2>&1 &
fi
echo $! > "$ROOT/.pids/web.pid"
