# CRITICAL: Hardcoded AI Outputs Issue

## Problem

AI is returning **identical hardcoded example data** for all PDFs:
- "Sarah" building authentication endpoint
- "PostgreSQL vs MongoDB" decision
- "Acme Corp API keys" blocker
- Same tasks, decisions, risks every time

This means **Ollama is not processing your actual PDF content**.

## Diagnosis Steps

### Step 1: Test Ollama Directly

```bash
./scripts/test-ollama-variation.sh
```

This will test if Ollama itself is working correctly.

**Expected Results**:
- ✅ Test 1 PASS: Different seeds → different outputs
- ✅ Test 2 PASS: Different inputs → different outputs

**If tests FAIL**: Ollama is broken or caching responses.

### Step 2: Check Backend Logs

```bash
# Restart backend
make stop
make start

# In another terminal, watch logs
tail -f apps/api/logs/meetext.log
```

Upload a PDF and look for:

```
meeting: PDF text extracted successfully
  text_length: 5432
  preview: "Your actual PDF content here..."

ai: sending to ollama for structured extraction
  summaries_preview: "Your actual summary here..."

ai: received from ollama
  response_preview: "Your actual AI response..."
```

**Critical Check**:
- Is `preview` showing YOUR PDF content or generic text?
- Is `summaries_preview` showing YOUR content or generic text?
- Is `response_preview` showing YOUR data or the hardcoded example?

## Root Causes

### Cause 1: Ollama is Caching Responses

**Symptom**: Same output regardless of input

**Fix**:
```bash
# Stop Ollama
pkill ollama

# Clear cache
rm -rf ~/.ollama/cache

# Restart Ollama
ollama serve &

# Reload model
ollama pull llama3.1:8b-instruct-q4_K_M
```

### Cause 2: Ollama is Not Receiving Actual Content

**Symptom**: Logs show correct PDF text but wrong AI output

**Fix**: The prompt is not including the actual content. Check prompts.

### Cause 3: PDF Extraction is Failing

**Symptom**: `text_length: 0` or very small

**Fix**:
```bash
# Install pdftotext
sudo apt-get install poppler-utils  # Ubuntu/Debian
brew install poppler                # macOS

# Test extraction
pdftotext your-test.pdf - | head -50
```

### Cause 4: Ollama Model is Corrupted

**Symptom**: Ollama returns same response for all inputs

**Fix**:
```bash
# Remove model
ollama rm llama3.1:8b-instruct-q4_K_M

# Re-download
ollama pull llama3.1:8b-instruct-q4_K_M

# Verify
ollama list
```

## Emergency Fixes

### Fix 1: Force Ollama Restart

```bash
# Kill all Ollama processes
pkill -9 ollama

# Wait 5 seconds
sleep 5

# Start fresh
ollama serve > /tmp/ollama.log 2>&1 &

# Wait for startup
sleep 3

# Pull model
ollama pull llama3.1:8b-instruct-q4_K_M

# Test
curl http://localhost:11434/api/tags
```

### Fix 2: Clear All Caches

```bash
# Stop everything
make stop
pkill ollama

# Clear Ollama cache
rm -rf ~/.ollama/cache
rm -rf ~/.ollama/models/manifests/*
rm -rf ~/.ollama/models/blobs/*

# Clear backend logs
rm -f apps/api/logs/*.log

# Restart
ollama serve &
sleep 5
ollama pull llama3.1:8b-instruct-q4_K_M
make start
```

### Fix 3: Test with Minimal Example

Create `test-ollama.sh`:

```bash
#!/bin/bash

echo "Testing Ollama with your actual content..."

# Test 1: Simple prompt
curl -s http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Extract tasks from: John will fix the bug by Friday.",
  "stream": false,
  "format": "json",
  "options": {"temperature": 0.3, "seed": 99999}
}' | jq -r '.response'

echo ""
echo "---"
echo ""

# Test 2: Different content
curl -s http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Extract tasks from: Mary will review the PR by Monday.",
  "stream": false,
  "format": "json",
  "options": {"temperature": 0.3, "seed": 99999}
}' | jq -r '.response'
```

Run:
```bash
chmod +x test-ollama.sh
./test-ollama.sh
```

**Expected**: First mentions "John" and "bug", second mentions "Mary" and "PR".

**If both mention "Sarah" and "authentication"**: Ollama is completely broken.

## Debugging Checklist

- [ ] Ollama is running: `curl http://localhost:11434/api/tags`
- [ ] Model is loaded: `ollama list | grep llama3.1`
- [ ] PDF extraction works: Check logs for `preview` with actual content
- [ ] Ollama receives different inputs: Check logs for `summaries_preview`
- [ ] Ollama returns different outputs: Check logs for `response_preview`
- [ ] Test script passes: `./scripts/test-ollama-variation.sh`

## Advanced Debugging

### Enable Full Ollama Logging

```bash
# Stop Ollama
pkill ollama

# Start with debug logging
OLLAMA_DEBUG=1 ollama serve > /tmp/ollama-debug.log 2>&1 &

# Upload a PDF

# Check Ollama logs
tail -f /tmp/ollama-debug.log
```

Look for:
- Is Ollama receiving the full prompt?
- Is the prompt including your PDF content?
- Is Ollama generating new responses or returning cached ones?

### Check Ollama API Directly

```bash
# Send your actual PDF text
curl -X POST http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Summarize: [PASTE YOUR ACTUAL PDF TEXT HERE]",
  "stream": false,
  "format": "json",
  "options": {"temperature": 0.3, "seed": 12345}
}' | jq -r '.response'
```

If this returns the hardcoded "Sarah/authentication" example, **Ollama is completely broken**.

### Check for Proxy/Cache

```bash
# Check if there's a proxy
env | grep -i proxy

# Check if there's a cache service
ps aux | grep -i cache

# Check network
netstat -an | grep 11434
```

## Nuclear Option: Complete Reset

If nothing works:

```bash
# 1. Stop everything
make stop
pkill -9 ollama

# 2. Remove Ollama completely
rm -rf ~/.ollama

# 3. Reinstall Ollama
curl https://ollama.ai/install.sh | sh

# 4. Start fresh
ollama serve &
sleep 5

# 5. Pull model
ollama pull llama3.1:8b-instruct-q4_K_M

# 6. Test
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Say hello",
  "stream": false
}' | jq -r '.response'

# 7. Restart backend
make start
```

## What to Share for Help

If issue persists, share:

1. **Ollama test results**:
   ```bash
   ./scripts/test-ollama-variation.sh > ollama-test.txt 2>&1
   ```

2. **Backend logs** (last 100 lines after uploading a PDF):
   ```bash
   tail -100 apps/api/logs/meetext.log > backend-logs.txt
   ```

3. **Ollama version**:
   ```bash
   ollama --version
   ```

4. **System info**:
   ```bash
   uname -a
   free -h
   ```

5. **PDF sample** (first page text):
   ```bash
   pdftotext your-test.pdf - | head -50
   ```

## Expected vs Actual

### Expected Behavior

Upload PDF about "Database Migration" → Get tasks about database migration

Upload PDF about "UI Redesign" → Get tasks about UI redesign

### Actual Behavior (BROKEN)

Upload ANY PDF → Always get:
- Sarah building authentication
- PostgreSQL decision
- Acme Corp blocker

This is **NOT temperature/seed issue**. This is **Ollama not processing inputs**.

## Immediate Action

1. Run: `./scripts/test-ollama-variation.sh`
2. If tests FAIL → Ollama is broken, follow "Nuclear Option"
3. If tests PASS → Check backend logs for what's being sent to Ollama
4. Share results

The problem is **NOT in the backend code**. The problem is **Ollama is returning cached/hardcoded responses**.
