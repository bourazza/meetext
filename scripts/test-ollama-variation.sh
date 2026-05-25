#!/bin/bash

# Test if Ollama is returning different outputs for different inputs

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Testing Ollama Response Variation"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Test 1: Same input, different seeds
echo "Test 1: Same input with different seeds"
echo "----------------------------------------"

curl -s http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Summarize this meeting: The team discussed authentication improvements for the mobile app.",
  "stream": false,
  "format": "json",
  "options": {
    "temperature": 0.3,
    "seed": 12345
  }
}' | jq -r '.response' | jq -r '.summary' > /tmp/test1_seed1.txt

echo "Response 1 (seed 12345):"
cat /tmp/test1_seed1.txt
echo ""

curl -s http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Summarize this meeting: The team discussed authentication improvements for the mobile app.",
  "stream": false,
  "format": "json",
  "options": {
    "temperature": 0.3,
    "seed": 67890
  }
}' | jq -r '.response' | jq -r '.summary' > /tmp/test1_seed2.txt

echo "Response 2 (seed 67890):"
cat /tmp/test1_seed2.txt
echo ""

if diff -q /tmp/test1_seed1.txt /tmp/test1_seed2.txt > /dev/null; then
  echo "❌ FAIL: Responses are identical (should be different)"
else
  echo "✅ PASS: Responses are different"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Test 2: Different inputs
echo "Test 2: Different inputs"
echo "----------------------------------------"

curl -s http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Summarize this meeting: The team discussed authentication improvements for the mobile app.",
  "stream": false,
  "format": "json",
  "options": {
    "temperature": 0.3,
    "seed": 11111
  }
}' | jq -r '.response' | jq -r '.summary' > /tmp/test2_input1.txt

echo "Response 1 (authentication topic):"
cat /tmp/test2_input1.txt
echo ""

curl -s http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Summarize this meeting: The team discussed database migration from MongoDB to PostgreSQL.",
  "stream": false,
  "format": "json",
  "options": {
    "temperature": 0.3,
    "seed": 11111
  }
}' | jq -r '.response' | jq -r '.summary' > /tmp/test2_input2.txt

echo "Response 2 (database topic):"
cat /tmp/test2_input2.txt
echo ""

if diff -q /tmp/test2_input1.txt /tmp/test2_input2.txt > /dev/null; then
  echo "❌ FAIL: Responses are identical (should be completely different)"
else
  echo "✅ PASS: Responses are different"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Test 3: Check if model is actually loaded
echo "Test 3: Model status"
echo "----------------------------------------"

curl -s http://localhost:11434/api/tags | jq -r '.models[] | select(.name | contains("llama3.1:8b")) | .name'

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Diagnosis Complete"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "If both tests FAIL:"
echo "  → Ollama is caching or broken"
echo "  → Try: ollama stop && ollama serve"
echo ""
echo "If Test 1 PASS but Test 2 FAIL:"
echo "  → Ollama is ignoring the prompt content"
echo "  → Check Ollama version: ollama --version"
echo ""
echo "If both tests PASS:"
echo "  → Ollama is working correctly"
echo "  → Problem is in the backend code"
echo ""
