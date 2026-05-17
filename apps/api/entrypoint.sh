#!/bin/sh
set -e

echo "→ Running database migrations..."
./migrate -direction=up

echo "→ Starting API server..."
exec ./api
