#!/bin/bash
set -e

echo "🚀 Starting Meetext services..."

# Start all services
docker-compose up -d

echo "⏳ Waiting for Ollama to be ready..."
sleep 10

# Check if model exists, if not pull it
MODEL="llama3.2:1b"
echo "📦 Checking if model $MODEL exists..."

if ! docker exec meetext_ollama ollama list | grep -q "$MODEL"; then
    echo "📥 Pulling model $MODEL (this may take a few minutes)..."
    docker exec meetext_ollama ollama pull "$MODEL"
    echo "✅ Model $MODEL pulled successfully"
else
    echo "✅ Model $MODEL already exists"
fi

echo ""
echo "🎉 Meetext is ready!"
echo ""
echo "📍 Services:"
echo "   - Frontend: http://localhost:3000"
echo "   - API:      http://localhost:8080"
echo "   - Ollama:   http://localhost:11434"
echo ""
echo "📊 View logs:"
echo "   docker-compose logs -f"
echo ""
echo "🛑 Stop services:"
echo "   docker-compose down"
