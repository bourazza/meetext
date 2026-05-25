#!/bin/bash

echo "Testing Ollama WITHOUT format=json parameter"
echo "=============================================="
echo ""

echo "Test 1: Simple prompt"
echo "---------------------"
curl -s http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "You MUST return ONLY valid JSON: {\"summary\": \"<your summary here>\"}. Summarize: The team discussed authentication improvements.",
  "stream": false,
  "options": {"temperature": 0.3, "seed": 12345}
}' | jq -r '.response'

echo ""
echo "---------------------"
echo ""

echo "Test 2: Different content"
echo "-------------------------"
curl -s http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "You MUST return ONLY valid JSON: {\"summary\": \"<your summary here>\"}. Summarize: The team discussed database migration.",
  "stream": false,
  "options": {"temperature": 0.3, "seed": 67890}
}' | jq -r '.response'

echo ""
echo "=============================================="
echo "If both responses show actual JSON (not null):"
echo "  ✅ Ollama is working"
echo ""
echo "If responses are null:"
echo "  ❌ Ollama is broken"
echo ""
