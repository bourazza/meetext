#!/bin/bash
set -e

echo "🚀 Starting Meetext (Local Development Mode)"
echo ""

# Check if Ollama is installed
if ! command -v ollama &> /dev/null; then
    echo "❌ Ollama is not installed"
    echo ""
    echo "Install Ollama:"
    echo "  curl -fsSL https://ollama.com/install.sh | sh"
    echo ""
    echo "Or visit: https://ollama.com/download"
    exit 1
fi

# Start Ollama in background if not running
if ! pgrep -x "ollama" > /dev/null; then
    echo "🔄 Starting Ollama service..."
    ollama serve > /dev/null 2>&1 &
    sleep 3
else
    echo "✅ Ollama is already running"
fi

# Pull model if not exists
MODEL="llama3.2:1b"
echo "📦 Checking model $MODEL..."
if ! ollama list | grep -q "$MODEL"; then
    echo "📥 Pulling model $MODEL (this may take a few minutes)..."
    ollama pull "$MODEL"
    echo "✅ Model pulled successfully"
else
    echo "✅ Model $MODEL ready"
fi

echo ""
echo "🎉 Ollama is ready at http://localhost:11434"
echo ""
echo "Now start your services:"
echo "  - API:  cd apps/api && go run cmd/api/main.go"
echo "  - Web:  cd apps/web && npm run dev"
