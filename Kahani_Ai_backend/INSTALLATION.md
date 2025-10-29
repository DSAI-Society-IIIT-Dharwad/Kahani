# 🚀 INSTALLATION GUIDE

## Prerequisites Check

Before starting, verify you have:

### 1. Python 3.9 or higher

```bash
python3 --version
```

Expected: `Python 3.9.x` or higher

### 2. Docker (for Milvus)

```bash
docker --version
docker-compose --version
```

If not installed:

- **macOS**: Download [Docker Desktop for Mac](https://www.docker.com/products/docker-desktop/)
- **Windows**: Download [Docker Desktop for Windows](https://www.docker.com/products/docker-desktop/)
- **Linux**: Install via package manager
  ```bash
  sudo apt-get install docker.io docker-compose  # Ubuntu/Debian
  sudo yum install docker docker-compose          # CentOS/RHEL
  ```

### 3. Git (optional, for cloning)

```bash
git --version
```

---

## Installation Steps

### Step 1: Navigate to Project

```bash
cd /Users/abhangsudhirpawar/Desktop/Arnav_Kahani_Ai
```

### Step 2: Start Milvus Vector Database

```bash
# Start Milvus with Docker Compose
docker-compose up -d

# Verify it's running
docker ps | grep milvus
```

You should see 3 containers:

- `milvus-standalone`
- `milvus-etcd`
- `milvus-minio`

**Wait 15-20 seconds** for Milvus to fully start.

Verify connection:

```bash
curl http://localhost:19530
```

### Step 3: Create Python Virtual Environment

```bash
# Create virtual environment
python3 -m venv venv

# Activate it
source venv/bin/activate  # macOS/Linux
# OR
venv\Scripts\activate     # Windows
```

Your prompt should now show `(venv)`.

### Step 4: Install Python Dependencies

```bash
# Upgrade pip first
pip install --upgrade pip

# Install all requirements
pip install -r requirements.txt
```

This will install:

- FastAPI & Uvicorn
- pymilvus (Milvus client)
- groq (Groq Cloud SDK)
- sentence-transformers (embeddings)
- SQLAlchemy (database)
- APScheduler (background tasks)
- And more...

**Note**: This may take 3-5 minutes as it downloads models.

### Step 5: Verify Configuration

```bash
# Check .env file exists
cat .env
```

Should show:

```
GROQ_API_KEY
MILVUS_HOST=localhost
MILVUS_PORT=19530
...
```

✅ Everything is already configured!

### Step 6: Initialize Database

The database will be created automatically on first run, but you can verify:

```bash
# This will be created when you first run the app
ls -la kahani.db  # Won't exist yet - that's OK!
```

### Step 7: Start the Application

```bash
python main.py
```

You should see:

```
INFO:     Started server process
INFO:     Waiting for application startup.
INFO:     Starting Kahani AI...
INFO:     Connected to Milvus at localhost:19530
INFO:     Application started successfully
INFO:     Uvicorn running on http://0.0.0.0:8000
```

### Step 8: Test the Installation

Open a new terminal and run:

```bash
# Test health endpoint
curl http://localhost:8000/health
```

Expected response:

```json
{
  "status": "healthy",
  "milvus_connected": true,
  "database_ok": true
}
```

✅ **Installation Complete!**

---

## Access the System

### Web Interface

Open your browser:

- **Story Writer UI**: http://localhost:8000/ui
- **API Documentation**: http://localhost:8000/docs
- **Alternative Docs**: http://localhost:8000/redoc

### Test API

```bash
# Create a story suggestion
curl -X POST http://localhost:8000/api/story/suggest \
  -H "Content-Type: application/json" \
  -d '{"user_prompt": "A wizard discovers a magical book", "user_id": "test"}'
```

---

## Quick Test Workflow

### Using the Web UI

1. **Open**: http://localhost:8000/ui

2. **Enter a prompt**:

   ```
   A mysterious traveler arrives at an ancient castle
   ```

3. **Click**: "Get AI Suggestion"

4. **Review** the AI-generated suggestion

5. **Edit** if needed, then **"Sign & Add to Story"**

6. **Add more lines** to build your story

7. **Click "Extract Lore"** to see characters and locations

8. **Click "Finalize Story"** to create a polished version

### Using Python Script

```bash
# Run the example workflow
python example_usage.py
```

This will:

- Create 5 story lines
- Show RAG context retrieval in action
- Extract lore (characters, locations, events)
- Create a canonical story

---

## Troubleshooting

### Issue: "Cannot connect to Milvus"

**Solution 1**: Check if Milvus is running

```bash
docker ps | grep milvus
```

**Solution 2**: Restart Milvus

```bash
docker-compose restart
sleep 20  # Wait for startup
```

**Solution 3**: Check logs

```bash
docker-compose logs milvus-standalone
```

**Solution 4**: Stop and start fresh

```bash
docker-compose down
docker-compose up -d
```

### Issue: "Module not found" errors

**Solution**: Make sure virtual environment is activated

```bash
# Check if (venv) appears in prompt
# If not, activate it:
source venv/bin/activate  # macOS/Linux
venv\Scripts\activate     # Windows

# Reinstall dependencies
pip install -r requirements.txt --force-reinstall
```

### Issue: "Port 8000 already in use"

**Solution**: Kill the process or use a different port

```bash
# Kill process on port 8000
lsof -ti:8000 | xargs kill -9

# Or run on different port
uvicorn main:app --host 0.0.0.0 --port 8001
```

### Issue: Groq API errors

**Solution 1**: Verify API key in `.env`

```bash
grep GROQ_API_KEY .env
```

**Solution 2**: Test API key

```bash
curl -X POST https://api.groq.com/openai/v1/chat/completions \
  -H "Authorization: Bearer " \
  -H "Content-Type: application/json" \
  -d '{"model":"llama-3.1-8b-instant","messages":[{"role":"user","content":"test"}]}'
```

### Issue: Database locked

**Solution**: Close all instances of the app

```bash
# Kill all Python processes
pkill -f "python main.py"

# Delete database and restart
rm kahani.db
python main.py  # Will recreate database
```

### Issue: Slow embedding generation

**Expected**: First run downloads the embedding model (~80MB)

- This only happens once
- Subsequent runs are fast
- Model is cached in: `~/.cache/torch/sentence_transformers/`

---

## Verification Checklist

After installation, verify:

- [ ] Docker is running: `docker ps`
- [ ] Milvus containers are up (3 containers)
- [ ] Virtual environment is activated: `(venv)` in prompt
- [ ] Dependencies installed: `pip list | grep fastapi`
- [ ] `.env` file exists with API key
- [ ] App starts without errors
- [ ] Health check returns "healthy"
- [ ] Web UI loads at http://localhost:8000/ui
- [ ] Can create a story suggestion

---

## Optional: Development Setup

### Install Dev Dependencies

```bash
pip install pytest httpx black flake8 mypy
```

### Run Tests

```bash
# If tests are added later
pytest tests/
```

### Code Formatting

```bash
black *.py
```

### Type Checking

```bash
mypy main.py
```

---

## Uninstallation

If you need to remove everything:

```bash
# Stop Milvus
docker-compose down

# Remove Milvus data (optional)
rm -rf volumes/

# Remove Python virtual environment
rm -rf venv/

# Remove database
rm kahani.db

# Remove cached models (optional)
rm -rf ~/.cache/torch/sentence_transformers/
```

---

## Next Steps

1. ✅ Installation complete
2. 🎨 Try the Web UI: http://localhost:8000/ui
3. 📚 Read the docs: http://localhost:8000/docs
4. 🧪 Run example: `python example_usage.py`
5. 📖 Read guides:
   - `QUICKSTART.md` - Quick start guide
   - `API_TESTING.md` - API examples
   - `ARCHITECTURE.md` - System design

---

## System Requirements

**Minimum**:

- CPU: 2 cores
- RAM: 4 GB
- Disk: 2 GB free space
- Python: 3.9+
- Docker: 20.10+

**Recommended**:

- CPU: 4+ cores
- RAM: 8 GB
- Disk: 5 GB free space
- Python: 3.11+
- Docker: 24.0+

---

**Installation Complete! Happy Storytelling! 🎭📖✨**
