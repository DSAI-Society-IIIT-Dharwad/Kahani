# 🚀 COMPLETE SETUP GUIDE - Kahani AI

## Current Status

✅ You've installed Milvus Lite successfully!  
⚠️ Need to create virtual environment and install dependencies

---

## Step-by-Step Installation

### Step 1: Fix Python/Pip Setup

You're using macOS with Python 3.9. Let's set up properly:

```bash
# Update pip first
/Library/Developer/CommandLineTools/usr/bin/python3 -m pip install --upgrade pip

# Add Python bin to PATH (for this session)
export PATH="/Users/abhangsudhirpawar/Library/Python/3.9/bin:$PATH"

# Or permanently add to your shell profile
echo 'export PATH="/Users/abhangsudhirpawar/Library/Python/3.9/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Step 2: Create Virtual Environment (IMPORTANT!)

```bash
# Navigate to project directory
cd /Users/abhangsudhirpawar/Desktop/Arnav_Kahani_Ai

# Create virtual environment
python3 -m venv venv

# Activate it
source venv/bin/activate

# You should see (venv) in your prompt now
```

### Step 3: Install All Dependencies

```bash
# Make sure you're in the virtual environment (venv should show in prompt)
# Upgrade pip in the virtual environment
pip install --upgrade pip

# Install all requirements
pip install -r requirements.txt
```

This will install:

- FastAPI (web framework)
- pymilvus (Milvus client)
- groq (Groq Cloud SDK)
- sentence-transformers (embeddings)
- SQLAlchemy (database)
- And 15+ other packages

**Note**: This may take 5-10 minutes as it downloads AI models.

### Step 4: Choose Milvus Option

You have 2 options for Milvus:

#### Option A: Use Milvus Lite (Simpler, Already Installed!)

✅ **Recommended for getting started quickly**

```bash
# You already installed this - nothing more needed!
# Milvus Lite runs in-process, no Docker required
```

#### Option B: Use Full Milvus with Docker (More Features)

```bash
# Install Docker Desktop first from: https://www.docker.com/products/docker-desktop/

# Then start Milvus
docker-compose up -d

# Wait 20 seconds for startup
sleep 20

# Verify
docker ps | grep milvus
```

### Step 5: Start the Application

```bash
# Make sure you're in project directory with venv activated
cd /Users/abhangsudhirpawar/Desktop/Arnav_Kahani_Ai
source venv/bin/activate  # If not already activated

# Start the app
python main.py
```

You should see:

```
INFO: Starting Kahani AI...
INFO: Connected to Milvus at localhost:19530
INFO: Application started successfully
INFO: Uvicorn running on http://0.0.0.0:8000
```

### Step 6: Test the System

Open a **new terminal** and test:

```bash
# Test health check
curl http://localhost:8000/health
```

Expected response:

```json
{ "status": "healthy", "milvus_connected": true, "database_ok": true }
```

### Step 7: Access the Web Interface

Open your browser:

- **Story Writer**: http://localhost:8000/ui
- **API Docs**: http://localhost:8000/docs

---

## Quick Commands (Copy & Paste)

Here's everything in order:

```bash
# 1. Navigate to project
cd /Users/abhangsudhirpawar/Desktop/Arnav_Kahani_Ai

# 2. Create and activate virtual environment
python3 -m venv venv
source venv/bin/activate

# 3. Install dependencies
pip install --upgrade pip
pip install -r requirements.txt

# 4. Start the application
python main.py
```

---

## If You Encounter Issues

### Issue: "ModuleNotFoundError"

**Solution**: Make sure virtual environment is activated

```bash
source venv/bin/activate  # You should see (venv) in prompt
pip install -r requirements.txt
```

### Issue: "Cannot connect to Milvus"

**Solution**: You're using Milvus Lite, so update the config:

Edit `.env` file and change:

```
MILVUS_HOST=localhost
MILVUS_PORT=19530
```

To:

```
MILVUS_HOST=localhost
MILVUS_PORT=19530
# Milvus Lite will auto-configure
```

### Issue: Slow installation

**Expected**: First time downloads AI models (~200MB total)

- sentence-transformers model: ~80MB
- Other dependencies: ~120MB
- This only happens once

### Issue: "Port 8000 in use"

```bash
# Find and kill process using port 8000
lsof -ti:8000 | xargs kill -9

# Or use different port
uvicorn main:app --host 0.0.0.0 --port 8001
```

---

## Verification Steps

After installation, verify everything works:

1. ✅ **Virtual environment active**: See `(venv)` in terminal prompt
2. ✅ **Dependencies installed**: `pip list | grep fastapi`
3. ✅ **App starts**: `python main.py` runs without errors
4. ✅ **Health check**: `curl http://localhost:8000/health` returns "healthy"
5. ✅ **Web UI loads**: http://localhost:8000/ui opens in browser
6. ✅ **Can create story**: Enter prompt, get AI suggestion

---

## Test the Complete Workflow

Once running, try this:

1. **Open**: http://localhost:8000/ui
2. **Enter prompt**: "A wizard discovers a magical crystal"
3. **Click**: "Get AI Suggestion"
4. **Review suggestion**, then click "Sign & Add to Story"
5. **Add another line**: "The crystal reveals an ancient prophecy"
6. **Notice**: Second suggestion uses context from first line!
7. **Extract lore**: Click "Extract Lore" to see characters/locations
8. **Finalize**: Click "Finalize Story" for polished version

---

## Expected Timeline

- **Virtual environment setup**: 1 minute
- **Dependency installation**: 5-10 minutes (downloads models)
- **App startup**: 30 seconds
- **First story creation**: 2-3 minutes

---

## What's Already Configured

✅ **Groq API Key**: Already in `.env` file  
✅ **All Python files**: Complete FastAPI application  
✅ **Database models**: SQLite will be created automatically  
✅ **Web interface**: Beautiful HTML UI ready  
✅ **API endpoints**: Full REST API with documentation

---

## Need Help?

If you run into issues:

1. **Check the logs**: Look at terminal output for error messages
2. **Verify venv**: Make sure `(venv)` shows in your prompt
3. **Check dependencies**: `pip list` should show fastapi, groq, etc.
4. **Test health**: `curl http://localhost:8000/health`

---

**Ready to start? Copy and paste the Quick Commands above! 🚀**
