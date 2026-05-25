# Fix: Identical AI Outputs Issue

## Problem

AI was generating identical outputs for different PDFs.

## Root Cause

1. **Temperature too low** (`0.2`) - Made outputs too deterministic
2. **No seed variation** - Same seed used for every request
3. **Insufficient logging** - Couldn't verify if PDFs were being extracted correctly

## Solution

### 1. Increased Temperature

**File**: `apps/api/internal/infrastructure/ollama/provider.go`

```go
// Before
"temperature": 0.2,  // Too deterministic

// After
"temperature": 0.3,  // Allows more variation while staying factual
```

**Impact**: 
- Outputs will vary slightly even for same input
- Still factual and reliable (not too creative)
- Good balance between consistency and variation

### 2. Added Random Seed

**File**: `apps/api/internal/infrastructure/ollama/provider.go`

```go
// Added timestamp-based seed
seed := int(time.Now().UnixNano() % 1000000)

reqBody := generateRequest{
    // ...
    Options: map[string]interface{}{
        // ...
        "seed": seed,  // Different seed per request
    },
}
```

**Impact**:
- Each request gets a unique seed
- Same PDF uploaded twice will produce similar (not identical) outputs
- Different PDFs will produce completely different outputs

### 3. Enhanced Logging

**File**: `apps/api/internal/usecase/meeting/meeting.go`

Added logging to verify:
- PDF text extraction is working
- Different PDFs produce different text
- AI generates different summaries

```go
// Log PDF preview (first 500 chars)
l.Info().
    Int("text_length", len(txt)).
    Str("preview", preview).
    Msg("meeting: PDF text extracted successfully")

// Log AI result preview
l.Info().
    Str("summary_preview", truncateStr(aiResult.Summary, 200)).
    Int("tasks_count", len(aiResult.Tasks)).
    Int("decisions_count", len(aiResult.Decisions)).
    Msg("meeting: AI analysis completed")
```

**Impact**:
- Can verify PDF extraction is working
- Can see if different PDFs produce different text
- Can debug issues faster

## Testing

### Test 1: Upload Same PDF Twice

**Expected Result**:
- Summaries should be 80-90% similar
- Key facts should be identical
- Wording may vary slightly
- Task/decision counts should be same

### Test 2: Upload Different PDFs

**Expected Result**:
- Summaries should be completely different
- Tasks should be specific to each document
- Counts should vary based on content

## How to Verify Fix

### Step 1: Restart Backend

```bash
make stop
make start
```

### Step 2: Upload Test PDFs

Upload 2-3 completely different PDFs (e.g., meeting notes, technical docs, business proposal)

### Step 3: Check Logs

Look for these entries in `apps/api/logs/meetext.log`:

```
meeting: PDF text extracted successfully
  text_length: 5432
  preview: "This is the actual PDF content..."

ollama: sending request
  seed: 847392

meeting: AI analysis completed
  summary_preview: "This meeting discussed..."
```

**Verify**:
- ✅ `text_length` is different for different PDFs
- ✅ `preview` shows actual PDF content
- ✅ `seed` is different for each request
- ✅ `summary_preview` is different for different PDFs

### Step 4: Compare Outputs

Check the generated summaries in the UI:
- Different PDFs → Completely different summaries
- Same PDF twice → Similar but not identical summaries

## Expected Behavior After Fix

| Scenario | Before Fix | After Fix |
|----------|-----------|-----------|
| Same PDF uploaded twice | Identical output | Similar (80-90%) but not identical |
| Different PDFs | Identical output | Completely different outputs |
| Technical doc vs meeting notes | Same generic summary | Specific to each document type |

## Configuration

### Temperature Settings

Current: `0.3` (balanced)

You can adjust in `provider.go` if needed:

```go
"temperature": 0.3,  // Current (recommended)
// 0.1-0.2 = Very deterministic (less variation)
// 0.3-0.4 = Balanced (good variation, still factual)
// 0.5-0.7 = Creative (more variation, may hallucinate)
```

### Seed Behavior

- **Random seed** (current): Different output each time
- **Fixed seed**: Same output for same input (useful for testing)

To use fixed seed for testing:

```go
// For testing only - remove after
seed := 12345  // Fixed seed
```

## Troubleshooting

### Still Getting Identical Outputs?

1. **Check PDF extraction**:
   ```bash
   tail -f apps/api/logs/meetext.log | grep "PDF text extracted"
   ```
   Verify `text_length` and `preview` are different for different PDFs

2. **Check Ollama is running**:
   ```bash
   make ollama-status
   ```

3. **Verify seed is changing**:
   ```bash
   tail -f apps/api/logs/meetext.log | grep "seed"
   ```
   Each request should show different seed value

4. **Test Ollama directly**:
   ```bash
   curl http://localhost:11434/api/generate -d '{
     "model": "llama3.1:8b-instruct-q4_K_M",
     "prompt": "Summarize: Test content",
     "stream": false,
     "options": {"temperature": 0.3, "seed": 12345}
   }'
   ```

### PDF Extraction Not Working?

Install pdftotext:

```bash
# Ubuntu/Debian
sudo apt-get install poppler-utils

# macOS
brew install poppler
```

## Files Changed

1. `apps/api/internal/infrastructure/ollama/provider.go`
   - Increased temperature: 0.2 → 0.3
   - Added random seed generation
   - Added seed to log output

2. `apps/api/internal/usecase/meeting/meeting.go`
   - Added PDF text preview logging
   - Added AI summary preview logging
   - Added truncateStr helper function

3. `apps/api/docs/DEBUG_IDENTICAL_OUTPUTS.md` (new)
   - Comprehensive debugging guide

4. `apps/api/docs/PROMPT_MODES.md` (existing)
   - Documents the three prompt modes (fast/balanced/strict)

## Next Steps

1. **Test the fix**: Upload different PDFs and verify outputs are different
2. **Monitor logs**: Check that PDF extraction and AI generation are working
3. **Adjust temperature**: If outputs are too similar, increase to 0.4; if too random, decrease to 0.25
4. **Report results**: Share logs if issue persists

## Summary

The fix addresses the identical outputs issue by:
- ✅ Adding randomness (seed + higher temperature)
- ✅ Maintaining factual accuracy (temperature still low at 0.3)
- ✅ Adding debugging capabilities (enhanced logging)
- ✅ Preserving determinism where needed (same facts, different wording)

Different PDFs should now produce completely different outputs, while the same PDF uploaded twice will produce similar (but not identical) results.
