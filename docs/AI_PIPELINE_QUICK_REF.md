# AI Pipeline Quick Reference

## Model Configuration

```env
OLLAMA_MODEL=llama3.1:8b-instruct-q4_K_M
```

**Size**: 4.7GB  
**Quality**: ⭐⭐⭐⭐ (Production-grade)  
**Speed**: Moderate (quality > speed)

---

## Pull Model

```bash
# Option 1: Using Make
make ollama-pull

# Option 2: Manual
ollama pull llama3.1:8b-instruct-q4_K_M

# Verify
ollama list
```

---

## Processing Pipeline

```
Upload → Extract → Clean → Chunk → Summarize → Merge → Extract → Save
```

**Stages**:
1. PDF Upload (instant)
2. Text Extraction
3. Text Cleaning
4. Semantic Chunking (2000-2500 words)
5. Parallel Summarization (2 workers)
6. Summary Merging
7. Structured Extraction
8. Database Save

---

## Performance Expectations

| Pages | Time | Status |
|-------|------|--------|
| 1-5 | 30-60s | ✅ Fast |
| 10-20 | 2-4 min | ✅ Good |
| 30-50 | 5-10 min | ✅ OK |
| 100+ | 15-30 min | ⚠️ Slow |

**Note**: Slow is acceptable. Quality > speed.

---

## Memory Usage

- **Peak**: ~8GB
- **Safe for**: 16GB RAM systems
- **Workers**: 2 parallel max

---

## Key Settings

```go
// Ollama
num_ctx: 8192          // 8k context
temperature: 0.2       // Deterministic
top_p: 0.9            // Nucleus sampling
repeat_penalty: 1.1   // Reduce repetition

// Chunking
targetWords: 2000-2500
overlapWords: 250

// Workers
maxParallel: 2
```

---

## Troubleshooting

### Model not found
```bash
ollama pull llama3.1:8b-instruct-q4_K_M
```

### Ollama not running
```bash
ollama serve &
```

### Out of memory
- Reduce workers to 1
- Close other apps
- Use smaller model

### Too slow
- Use GPU (5-10x faster)
- Use smaller model (trades quality)

---

## Logs

```bash
# Watch logs
tail -f apps/api/logs/meetext.log

# Check Ollama
curl http://localhost:11434/api/tags

# Check model
ollama list
```

---

## Commands

```bash
# Start everything
make start

# Stop everything
make stop

# Pull model
make ollama-pull

# Check Ollama status
make ollama-status

# Restart services
make stop && make start
```

---

## Supported Files

✅ **PDF only**  
❌ Audio (coming with Whisper)  
❌ Video (coming with Whisper)

---

## Progress Stages

1. Extracting text
2. Cleaning transcript
3. Chunking text
4. Processing chunk X/Y
5. Extracting structured data
6. Generating report
7. Completed

---

## Error Handling

- **Retry**: 2 attempts with backoff
- **Fallback**: Truncated text if LLM fails
- **Recovery**: Panic recovery in async
- **Status**: Mark as `failed` on error

---

## Documentation

- **Architecture**: `docs/AI_PIPELINE_ARCHITECTURE.md`
- **Full Summary**: `docs/PRODUCTION_AI_REFACTOR_SUMMARY.md`
- **Ollama Setup**: `docs/OLLAMA_SETUP.md`
- **Quick Start**: `QUICKSTART.md`
