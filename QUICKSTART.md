# Meetext Quick Start

## Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Redis 7+
- Ollama (for AI features)

---

## 1. Install Ollama

```bash
# Install Ollama
curl -fsSL https://ollama.com/install.sh | sh

# Verify installation
ollama --version
```

---

## 2. Start Ollama and Pull Model

```bash
# Start Ollama service in background
ollama serve > /dev/null 2>&1 &

# Pull the lightweight model (1.3GB)
ollama pull llama3.2:1b

# Verify it's running
curl http://localhost:11434/api/tags
```

---

## 3. Setup Database

```bash
# Create PostgreSQL database
make local-setup

# Run migrations
make local-migrate
```

---

## 4. Start Services

### Option A: Using Make (Recommended)

```bash
# Start everything (API + Web + Ollama check)
make start
```

### Option B: Manual

```bash
# Terminal 1: Start API
cd apps/api
go run cmd/api/main.go

# Terminal 2: Start Web
cd apps/web
npm install
npm run dev
```

---

## 5. Access the Application

- **Frontend**: http://localhost:3000
- **API**: http://localhost:8080
- **Ollama**: http://localhost:11434

---

## Test the Upload Pipeline

1. Go to http://localhost:3000
2. Login or register
3. Upload a PDF file
4. Watch the AI processing in real-time
5. View extracted tasks, decisions, and summaries

---

## Troubleshooting

### "Connection refused" to Ollama

```bash
# Check if Ollama is running
pgrep -fl ollama

# If not, start it
ollama serve &

# Pull the model if missing
ollama pull llama3.2:1b
```

### "400 Bad Request" on upload

- Check that the file is a PDF
- Audio/video uploads are not supported yet (coming with Whisper integration)

### "AI processing failed"

- Ensure Ollama is running: `curl http://localhost:11434/api/tags`
- Check API logs: `tail -f apps/api/logs/meetext.log`
- Verify model is pulled: `ollama list`

### Slow AI processing

- Use a smaller model: `llama3.2:1b` (default)
- Close other heavy applications
- Ensure you have at least 8GB RAM available

---

## Stop Services

```bash
# Stop all services
make stop

# Stop Ollama
pkill ollama
```

---

## Docker Setup (Alternative)

If you prefer Docker:

```bash
# Start all services including Ollama
./scripts/start.sh

# Or manually
docker-compose up -d
docker exec meetext_ollama ollama pull llama3.2:1b
```

---

## Next Steps

- Read [OLLAMA_SETUP.md](./docs/OLLAMA_SETUP.md) for advanced Ollama configuration
- Check [README.md](./README.md) for full architecture documentation
- Explore the API at http://localhost:8080/health
