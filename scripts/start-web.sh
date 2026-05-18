#!/bin/sh
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/apps/web"
npm run dev -- --port 3000 > /tmp/meetext-web.log 2>&1 &
echo $! > "$ROOT/.pids/web.pid"
