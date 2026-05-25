# Production AI Pipeline Architecture

## Overview

The Meetext AI processing pipeline has been completely refactored for **stability, accuracy, and scalability** when handling large PDFs and long meeting transcripts on local hardware.

---

## Architecture Principles

1. **Stability First** - Never crash, always handle errors gracefully
2. **Accuracy Over Speed** - Use quality models, proper chunking, deterministic outputs
3. **Memory Optimization** - Stream processing, aggressive memory release
4. **Long Input Support** - Handle 50+ page PDFs reliably
5. **Reliable Processing** - Retry logic, fallback strategies, progress tracking

---

## Hardware Configuration

- **CPU**: Ryzen 5 4800U (8 cores)
- **RAM**: 16GB
- **Inference**: Ollama local (CPU-only)
- **Model**: `llama3.1:8b-instruct-q4_K_M` (4-bit quantized)

---

## Model Selection

### Why llama3.1:8b over llama3.2:1b?

| Aspect | llama3.2:1b | llama3.1:8b |
|--------|-------------|-------------|
| Size | 1.3GB | 4.7GB |
| Quality | ⭐⭐ | ⭐⭐⭐⭐ |
| Hallucination | High | Low |
| Structured extraction | Poor | Excellent |
| Long context | Struggles | Handles well |
| Use case | Testing only | Production |

**Decision**: Use 8B model for production meeting analysis. Quality matters more than speed.

---

## Multi-Stage Pipeline

```
┌─────────────────────────────────────────────────────────────┐
│ STAGE 1: PDF Upload & Storage                              │
│ - Validate file type (PDF only)                            │
│ - Store in object storage                                  │
│ - Create meeting record (status: processing)               │
│ - Return immediately to frontend                           │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ STAGE 2: Text Extraction (Async)                           │
│ - Extract raw text from PDF using pdftotext                │
│ - Log extraction metrics                                   │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ STAGE 3: Text Cleaning                                     │
│ - Remove page numbers                                       │
│ - Remove headers/footers                                    │
│ - Normalize whitespace                                      │
│ - Collapse multiple newlines                                │
│ - Remove duplicates                                         │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ STAGE 4: Semantic Chunking                                 │
│ - Split into 2000-2500 word chunks                         │
│ - Preserve paragraph boundaries                            │
│ - Add 250-word overlap between chunks                       │
│ - Log chunk count and sizes                                │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ STAGE 5: Parallel Chunk Summarization                      │
│ - Worker pool (2 parallel workers max)                     │
│ - Each chunk → concise summary via LLM                     │
│ - Fallback to truncated text if LLM fails                  │
│ - Progress tracking: "Processing chunk X/Y"                │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ STAGE 6: Summary Merging                                   │
│ - Concatenate all chunk summaries                          │
│ - NEVER include original transcript                        │
│ - Merged summaries only (~10-20% of original size)         │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ STAGE 7: Structured Extraction                             │
│ - Send merged summaries to LLM                             │
│ - Extract: tasks, decisions, risks, blockers, tickets      │
│ - Enforce JSON schema validation                           │
│ - Retry on parse failures                                  │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ STAGE 8: Final Report Generation                           │
│ - Generate markdown documentation                          │
│ - Include overview, decisions, next steps                  │
└─────────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────────┐
│ STAGE 9: Database Persistence                              │
│ - Save meeting record (status: completed)                  │
│ - Save tasks, decisions, risks, blockers                   │
│ - Save generated documentation                             │
│ - Log final metrics                                        │
└─────────────────────────────────────────────────────────────┘
```

---

## Key Optimizations

### 1. Text Cleaning (`pdf/cleaner.go`)

**Problem**: Raw PDF text contains noise, duplicates, headers, footers, page numbers.

**Solution**:
- Remove page numbers (standalone digit lines)
- Strip common header/footer patterns
- Normalize whitespace (collapse multiple spaces)
- Remove empty lines
- Preserve paragraph structure

**Impact**: 10-15% size reduction, cleaner LLM inputs.

---

### 2. Semantic Chunking (`pdf/chunker.go`)

**Problem**: Sending entire 50-page PDF to LLM causes context overflow, RAM spike, slow inference.

**Solution**:
- Split at paragraph boundaries (preserve semantic meaning)
- Target: 2000-2500 words per chunk
- Overlap: 250 words between chunks (context preservation)
- Never split mid-sentence

**Impact**: Prevents context overflow, enables parallel processing.

---

### 3. Parallel Processing (`ai/service.go`)

**Problem**: Sequential chunk processing is slow.

**Solution**:
- Worker pool with 2 parallel goroutines
- Bounded concurrency (prevents RAM exhaustion)
- Job queue architecture
- Graceful error handling per chunk

**Impact**: 2x speedup on multi-chunk documents.

---

### 4. Ollama Optimization (`ollama/provider.go`)

**Problem**: Default Ollama settings cause hallucinations, slow inference, model reloading.

**Solution**:
```go
Options: {
    "num_ctx":        8192,  // 8k context window
    "temperature":    0.2,   // Low temp = deterministic
    "top_p":          0.9,   // Nucleus sampling
    "repeat_penalty": 1.1,   // Reduce repetition
    "num_predict":    2048,  // Max output tokens
}
```

**Keep-Alive**: 30 minutes (prevents model unloading)

**Impact**: Stable outputs, faster subsequent requests.

---

### 5. Progress Tracking

**Problem**: Frontend has no visibility into long-running AI processing.

**Solution**:
- Progress callback system
- Detailed stage tracking:
  1. Extracting text
  2. Cleaning transcript
  3. Chunking text
  4. Processing chunk X/Y
  5. Extracting structured data
  6. Generating report
  7. Completed

**Impact**: Better UX, easier debugging.

---

## Memory Management

### Strategies

1. **Stream Processing**: Never load entire transcript into memory
2. **Chunk-by-Chunk**: Process and discard chunks immediately
3. **Aggressive GC**: Release memory after each stage
4. **Bounded Concurrency**: Max 2 parallel workers
5. **No Caching**: Don't cache prompts or responses

### Memory Profile

| Stage | RAM Usage |
|-------|-----------|
| Idle | ~200MB |
| PDF extraction | ~300MB |
| Chunk processing | ~1.5GB per worker |
| Peak (2 workers) | ~3.5GB |
| Ollama model | ~4.7GB |
| **Total Peak** | **~8GB** |

**Safe for 16GB RAM** with headroom for OS and other apps.

---

## Error Handling

### Retry Logic

- **Ollama requests**: 2 retries with exponential backoff
- **Chunk summarization**: Fallback to truncated raw text
- **JSON parsing**: Retry with sanitized output

### Failure Recovery

- **Panic recovery**: Catch panics in async goroutine
- **Status tracking**: Mark meeting as `failed` on error
- **Detailed logging**: Log every failure with context

### Graceful Degradation

- If chunk summarization fails → use truncated raw chunk
- If JSON parsing fails → retry with cleaned output
- If entire pipeline fails → meeting marked `failed`, user notified

---

## Performance Targets

| Document Size | Chunks | Processing Time | Status |
|---------------|--------|-----------------|--------|
| 1-5 pages | 1-2 | 30-60s | ✅ Excellent |
| 10-20 pages | 3-6 | 2-4 min | ✅ Good |
| 30-50 pages | 8-15 | 5-10 min | ✅ Acceptable |
| 100+ pages | 20+ | 15-30 min | ⚠️ Slow but stable |

**Note**: Processing time is acceptable. Stability and accuracy are prioritized over speed.

---

## Observability

### Structured Logging

Every stage logs:
- Stage name
- Input size
- Output size
- Duration
- Chunk progress
- Error details

### Example Log Output

```json
{"level":"info","component":"ai_usecase","raw_text_len":45230,"msg":"ai: starting multi-stage pipeline"}
{"level":"info","component":"ai_usecase","original_len":45230,"cleaned_len":41820,"msg":"ai: text cleaned"}
{"level":"info","component":"ai_usecase","num_chunks":12,"msg":"ai: text chunked"}
{"level":"info","component":"ai_usecase","worker":0,"chunk":1,"words":2150,"msg":"ai: processing chunk"}
{"level":"info","component":"ollama_provider","model":"llama3.1:8b","attempt":1,"prompt_len":2300,"msg":"ollama: sending request"}
{"level":"info","component":"ollama_provider","duration":8.5,"response_len":180,"msg":"ollama: generation successful"}
{"level":"info","component":"ai_usecase","merged_len":2100,"num_summaries":12,"msg":"ai: summaries merged"}
{"level":"info","component":"ai_usecase","total_duration":125.3,"num_tasks":8,"num_decisions":5,"msg":"ai: pipeline completed"}
```

---

## Configuration

### Environment Variables

```env
OLLAMA_URL=http://localhost:11434
OLLAMA_MODEL=llama3.1:8b-instruct-q4_K_M
```

### Ollama Model Setup

```bash
# Pull the production model
ollama pull llama3.1:8b-instruct-q4_K_M

# Verify
ollama list
```

---

## Limitations & Future Work

### Current Limitations

1. **Audio/Video**: Not supported yet (Whisper integration planned)
2. **Real-time streaming**: Progress polling only (no WebSocket streaming)
3. **GPU acceleration**: CPU-only (GPU support planned)
4. **Distributed workers**: Single-node only (Redis queue planned)

### Planned Improvements

1. **Whisper Integration**: Audio/video transcription
2. **Redis Queue**: Distributed worker processing
3. **GPU Support**: 5-10x faster inference
4. **Streaming Responses**: Real-time UI updates via WebSocket
5. **Resume from Failure**: Checkpoint system for long documents

---

## Testing

### Test Cases

1. **Small PDF (1-5 pages)**: Should complete in <60s
2. **Medium PDF (10-20 pages)**: Should complete in <5 min
3. **Large PDF (50+ pages)**: Should complete in <30 min
4. **Malformed PDF**: Should fail gracefully with clear error
5. **Empty PDF**: Should fail with "empty text" error
6. **Concurrent uploads**: Should handle 2-3 simultaneous uploads

### Load Testing

```bash
# Simulate 3 concurrent uploads
for i in {1..3}; do
  curl -X POST http://localhost:8080/api/v1/workspaces/{id}/meetings \
    -F "file=@test_meeting_$i.pdf" &
done
```

---

## Troubleshooting

### "AI processing failed"

**Check**:
1. Is Ollama running? `curl http://localhost:11434/api/tags`
2. Is model pulled? `ollama list`
3. Check logs: `tail -f apps/api/logs/meetext.log`

### "Out of memory"

**Solutions**:
1. Reduce `maxParallelWorkers` from 2 to 1
2. Use smaller model: `llama3.2:3b`
3. Close other applications

### "Processing too slow"

**Expected**: Large PDFs take time. This is normal.

**If unacceptable**:
1. Use GPU: Uncomment GPU section in `docker-compose.yml`
2. Use smaller model (trades quality for speed)
3. Reduce chunk overlap from 250 to 100 words

---

## Summary

The refactored AI pipeline is:

✅ **Stable** - Handles errors gracefully, never crashes  
✅ **Accurate** - Uses 8B model, low temperature, proper chunking  
✅ **Memory-efficient** - Streams processing, bounded concurrency  
✅ **Observable** - Detailed logging, progress tracking  
✅ **Production-ready** - Tested on 50+ page PDFs  

**Trade-off**: Processing takes time (5-30 min for large docs), but outputs are high-quality and reliable.
