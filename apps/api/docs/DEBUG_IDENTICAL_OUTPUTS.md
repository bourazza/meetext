# Debugging: Identical AI Outputs Issue

## Problem

AI always generates the same output regardless of PDF content.

## Root Causes

### 1. **Low Temperature (FIXED)**
- **Before**: `temperature: 0.2` - too deterministic
- **After**: `temperature: 0.3` - allows more variation
- **Impact**: Higher temperature = more diverse outputs

### 2. **No Seed Variation (FIXED)**
- **Before**: No seed parameter - same seed every time
- **After**: `seed: time.Now().UnixNano()` - different seed per request
- **Impact**: Different seeds = different outputs for same input

### 3. **PDF Extraction May Be Failing**
- If PDF text extraction returns empty or identical text, AI will produce identical outputs
- **Check logs** for PDF text preview

## Debugging Steps

### Step 1: Check Backend Logs

Start the backend and watch logs:

```bash
cd apps/api
make run
```

Upload a PDF and look for these log entries:

```
meeting: PDF text extracted successfully
  text_length: 5432
  preview: "This is the actual PDF content..."
```

**What to verify**:
- `text_length` should be different for different PDFs
- `preview` should show actual PDF content, not empty or generic text

### Step 2: Check AI Analysis Logs

Look for:

```
meeting: AI analysis completed
  summary_preview: "This meeting discussed..."
  tasks_count: 5
  decisions_count: 3
```

**What to verify**:
- `summary_preview` should be different for different PDFs
- Counts should vary based on content

### Step 3: Check Ollama Logs

Look for:

```
ollama: sending request
  model: llama3.1:8b-instruct-q4_K_M
  prompt_len: 2345
  seed: 847392
```

**What to verify**:
- `seed` should be different for each request
- `prompt_len` should vary based on PDF size

### Step 4: Test with Different PDFs

Upload 3 completely different PDFs:

1. **Technical document** (e.g., API documentation)
2. **Meeting notes** (e.g., team standup)
3. **Business proposal** (e.g., project plan)

Compare outputs:
- Summaries should be completely different
- Tasks should be specific to each document
- Technical notes should reflect actual content

## Expected Behavior

### Different PDFs → Different Outputs

| PDF Type | Expected Summary | Expected Tasks |
|----------|------------------|----------------|
| API Docs | "Technical documentation for REST API..." | "Document authentication flow", "Add code examples" |
| Meeting Notes | "Team discussed sprint planning..." | "Complete user story #123", "Review PR #456" |
| Business Proposal | "Project proposal for client X..." | "Prepare budget estimate", "Schedule kickoff meeting" |

### Same PDF → Similar (but not identical) Outputs

Due to `temperature: 0.3` and random seeds, uploading the same PDF twice should produce:
- **Similar** summaries (same key points)
- **Slightly different** wording
- **Same** tasks/decisions (factual content)

## Quick Test

### Test 1: Upload Same PDF Twice

```bash
# Upload test.pdf
# Upload test.pdf again
# Compare outputs
```

**Expected**: Summaries should be 80-90% similar but not word-for-word identical.

### Test 2: Upload Different PDFs

```bash
# Upload meeting-notes.pdf
# Upload technical-spec.pdf
# Compare outputs
```

**Expected**: Summaries should be completely different.

## Common Issues

### Issue 1: PDF Text Extraction Failing

**Symptom**: All PDFs show `text_length: 0` or very small numbers

**Solution**:
```bash
# Check if pdftotext is installed
which pdftotext

# Install if missing (Ubuntu/Debian)
sudo apt-get install poppler-utils

# Install if missing (macOS)
brew install poppler
```

### Issue 2: Ollama Model Not Loaded

**Symptom**: Ollama requests timeout or fail

**Solution**:
```bash
# Check Ollama status
curl http://localhost:11434/api/tags

# Pull model if missing
ollama pull llama3.1:8b-instruct-q4_K_M

# Verify model loaded
ollama list
```

### Issue 3: Temperature Too Low

**Symptom**: Outputs are identical even with different seeds

**Solution**: Already fixed - temperature increased to 0.3

### Issue 4: Caching Issue

**Symptom**: Backend returns cached results

**Solution**:
```bash
# Restart backend
cd apps/api
make restart

# Clear Redis cache if using Redis
redis-cli FLUSHALL
```

## Verification Checklist

- [ ] Backend logs show different `text_length` for different PDFs
- [ ] Backend logs show different `preview` content for different PDFs
- [ ] Ollama logs show different `seed` values for each request
- [ ] AI summaries are different for different PDFs
- [ ] Task counts vary based on PDF content
- [ ] Same PDF uploaded twice produces similar (not identical) outputs

## Advanced Debugging

### Enable Verbose Logging

Edit `.env`:

```bash
LOG_LEVEL=debug
```

Restart backend:

```bash
cd apps/api
make restart
```

### Check Raw Ollama Response

Add temporary logging in `provider.go`:

```go
// After line: var genResp generateResponse
p.log.Debug().Str("raw_response", string(body)).Msg("ollama: raw response")
```

This will show the exact JSON returned by Ollama.

### Test Ollama Directly

```bash
# Test with same prompt twice
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Summarize: The team discussed authentication improvements.",
  "stream": false,
  "options": {
    "temperature": 0.3,
    "seed": 12345
  }
}'

# Test with different seed
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.1:8b-instruct-q4_K_M",
  "prompt": "Summarize: The team discussed authentication improvements.",
  "stream": false,
  "options": {
    "temperature": 0.3,
    "seed": 67890
  }
}'
```

**Expected**: Responses should be similar but not identical.

## Solution Summary

### Changes Made

1. **Increased temperature**: `0.2` → `0.3`
2. **Added random seed**: `time.Now().UnixNano()`
3. **Added debug logging**: PDF preview, AI summary preview, seed values

### Next Steps

1. Restart backend
2. Upload 2-3 different PDFs
3. Check logs for:
   - Different PDF text previews
   - Different AI summaries
   - Different seed values
4. Verify outputs are different

## Still Having Issues?

If outputs are still identical after these fixes:

1. **Share backend logs** - Copy the full log output when uploading a PDF
2. **Share PDF samples** - Provide 2 different PDFs you're testing with
3. **Check Ollama version** - Run `ollama --version` (should be 0.1.0+)
4. **Test Ollama directly** - Use curl commands above to verify Ollama itself works

The issue is most likely:
- PDF text extraction returning empty/identical text
- Ollama not receiving different prompts
- Caching somewhere in the pipeline
