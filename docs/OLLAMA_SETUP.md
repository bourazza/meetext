# Ollama Setup Guide

Meetext uses [Ollama](https://ollama.com) for local AI inference. This guide covers installation and configuration.

---

## Quick Start (Local Development)

```bash
# 1. Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# 2. Start Ollama and pull the model
make ollama-start

# 3. Verify it's running
make ollama-status
```

---

## Manual Installation

### Linux / macOS

```bash
curl -fsSL https://ollama.com/install.sh | sh
```

### Windows

Download from: https://ollama.com/download

---

## Starting Ollama

### Option 1: Automatic (Recommended)

```bash
make ollama-start
```

This will:
- Start Ollama service in the background
- Pull the `llama3.2:1b` model automatically
- Verify everything is working

### Option 2: Manual

```bash
# Start Ollama service
ollama serve &

# Pull the model
ollama pull llama3.2:1b

# Verify
curl http://localhost:11434/api/tags
```

---

## Docker Setup

If running the full stack with Docker Compose:

```bash
# Start all services including Ollama
./scripts/start.sh

# Or manually
docker-compose up -d
docker exec meetext_ollama ollama pull llama3.2:1b
```

---

## Model Configuration

The default model is `llama3.2:1b` (lightweight, fast, good for development).

To change the model, edit `apps/api/.env`:

```env
OLLAMA_MODEL=llama3.2:1b
```

### Recommended Models

| Model | Size | Speed | Quality | Use Case |
|-------|------|-------|---------|----------|
| `llama3.2:1b` | 1.3GB | ⚡⚡⚡ | ⭐⭐ | Development, testing |
| `llama3.2:3b` | 2.0GB | ⚡⚡ | ⭐⭐⭐ | Balanced |
| `qwen2.5:3b` | 2.0GB | ⚡⚡ | ⭐⭐⭐ | Strong structured extraction |
| `llama3.1:8b` | 4.7GB | ⚡ | ⭐⭐⭐⭐ | Production quality |

---

## Troubleshooting

### "Connection refused" error

```bash
# Check if Ollama is running
pgrep -fl ollama

# If not running, start it
ollama serve &
```

### Model not found

```bash
# List installed models
ollama list

# Pull the model
ollama pull llama3.2:1b
```

### Slow inference

- Use a smaller model (`llama3.2:1b` instead of `llama3.1:8b`)
- Ensure you have at least 8GB RAM available
- Close other heavy applications

### GPU Support (Optional)

If you have an NVIDIA GPU:

```bash
# Install NVIDIA Container Toolkit
# https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html

# Uncomment GPU section in docker-compose.yml
# Then restart:
docker-compose down
docker-compose up -d
```

---

## API Endpoints

Ollama exposes a REST API on `http://localhost:11434`:

```bash
# List models
curl http://localhost:11434/api/tags

# Generate text
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2:1b",
  "prompt": "Why is the sky blue?",
  "stream": false
}'
```

---

## Performance Tips

1. **Use the smallest model that works** — `llama3.2:1b` is 10x faster than `llama3.1:8b`
2. **Keep context small** — Meetext automatically chunks large PDFs to prevent context overflow
3. **Monitor RAM usage** — Each model loads into RAM. Close unused models with `ollama stop <model>`
4. **Preload models** — Run `ollama pull <model>` before starting the API to avoid first-request delays

---

## Uninstall

```bash
# Stop Ollama
pkill ollama

# Remove Ollama (Linux/macOS)
sudo rm -rf /usr/local/bin/ollama
sudo rm -rf ~/.ollama

# Remove models only
rm -rf ~/.ollama/models
```

---

## Resources

- [Ollama Documentation](https://github.com/ollama/ollama/blob/main/docs/README.md)
- [Model Library](https://ollama.com/library)
- [API Reference](https://github.com/ollama/ollama/blob/main/docs/api.md)
