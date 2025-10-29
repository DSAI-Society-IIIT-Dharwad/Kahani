# 🚀 QUICK START GUIDE - Kahani AI

## What You Need to Do Manually

### 1. Install Milvus Vector Database

**Option A: Docker (Recommended)**

```bash
cd /Users/abhangsudhirpawar/Desktop/Arnav_Kahani_Ai
docker-compose up -d
```

**Option B: Milvus Lite (Simpler, but limited)**

```bash
pip install milvus
```

**Verify Milvus is Running:**

```bash
curl http://localhost:19530
```

### 2. Install Python Dependencies

```bash
cd /Users/abhangsudhirpawar/Desktop/Arnav_Kahani_Ai

# Create virtual environment
python3 -m venv venv
source venv/bin/activate  # On macOS/Linux

# Install dependencies
pip install -r requirements.txt
```

### 3. Start the Application

```bash
python main.py
```

### 4. Access the System

- **Web UI**: http://localhost:8000/ui
- **API Documentation**: http://localhost:8000/docs
- **Health Check**: http://localhost:8000/health

---

## 📖 System Flow

### User Story Creation Flow:

1. **User asks** → Enter prompt: "A knight discovers a hidden door"

2. **RAG pulls context** → System searches vector DB for relevant story lines

3. **LLM-1 proposes** → Groq generates suggestion based on context

4. **User picks/edits & signs** → User approves or modifies the suggestion

5. **Backend verifies & stores** → Story line saved to database

6. **Creates embedding** → Sentence transformer creates vector embedding

7. **Stores in Milvus** → Embedding stored for future retrieval

8. **Periodic extractors** → Background tasks extract lore (characters, locations, events)

9. **Lore summarizers** → System generates summaries and stores in vector DB

10. **LLM-2 canonicalizes** → When ready, LLM-2 creates polished final version

---

## 🎯 Key API Endpoints

### Create Story

- `POST /api/story/suggest` - Get AI suggestion with RAG context
- `POST /api/story/edit` - Edit and sign story line
- `GET /api/story/lines` - Get all story lines

### Lore & Context

- `POST /api/lore/extract` - Extract characters, locations, events
- `GET /api/lore/all` - View all extracted lore
- `POST /api/context/retrieve` - Search vector DB for context

### Finalization

- `POST /api/story/canonicalize` - Create final polished version
- `GET /api/story/canonical/{id}` - Get canonical story

---

## 🛠️ Configuration

All settings are in `.env` file:

- ✅ **GROQ_API_KEY** - Already configured with your key
- ✅ **MILVUS_HOST** - localhost (default)
- ✅ **MILVUS_PORT** - 19530 (default)
- ✅ **LLM_MODEL** - llama-3.1-70b-versatile
- ✅ **EMBEDDING_MODEL** - all-MiniLM-L6-v2 (384 dimensions)

---

## 📂 Project Structure

```
Arnav_Kahani_Ai/
├── main.py                 # FastAPI app with all endpoints
├── config.py               # Environment configuration
├── database.py             # SQLAlchemy setup
├── models.py               # Database models (StoryLine, LoreEntry, etc.)
├── schemas.py              # Pydantic request/response schemas
├── milvus_service.py       # Vector DB client
├── embedding_service.py    # Sentence transformer embeddings
├── llm_service.py          # Groq LLM (LLM-1 & LLM-2)
├── rag_service.py          # RAG pipeline (retrieve + generate)
├── background_tasks.py     # Periodic lore extraction & summarization
├── static/
│   └── index.html          # Beautiful web UI
├── requirements.txt        # Python dependencies
├── docker-compose.yml      # Milvus setup
├── .env                    # Environment variables (configured)
├── README.md              # Full documentation
└── QUICKSTART.md          # This file
```

---

## 🎨 Using the Web UI

1. Open http://localhost:8000/ui

2. **Enter a prompt**:

   ```
   "A mysterious traveler arrives at the ancient castle"
   ```

3. **Click "Get AI Suggestion"**

   - System retrieves relevant context from previous story
   - LLM generates contextual suggestion
   - Shows context used

4. **Edit if needed**, then **"Sign & Add to Story"**

   - Creates embedding
   - Stores in vector DB
   - Available for future context

5. **Continue writing** - Each new prompt gets better context!

6. **Extract Lore** - Click to analyze characters, locations, events

7. **Finalize Story** - Click to create canonical polished version

---

## 🧪 Testing with cURL

```bash
# 1. Get suggestion
curl -X POST http://localhost:8000/api/story/suggest \
  -H "Content-Type: application/json" \
  -d '{
    "user_prompt": "The hero discovers a magical sword",
    "user_id": "test_user"
  }'

# 2. Sign and add to story
curl -X POST http://localhost:8000/api/story/edit \
  -H "Content-Type: application/json" \
  -d '{
    "llm_proposed": "The sword gleamed...",
    "final_text": "The ancient sword gleamed with ethereal light.",
    "user_id": "test_user"
  }'

# 3. Get all lines
curl http://localhost:8000/api/story/lines

# 4. Extract lore
curl -X POST http://localhost:8000/api/lore/extract \
  -H "Content-Type: application/json" \
  -d '{"line_ids": [1, 2, 3]}'

# 5. Create canonical version
curl -X POST http://localhost:8000/api/story/canonicalize \
  -H "Content-Type: application/json" \
  -d '{"title": "The Sword of Destiny"}'
```

---

## ⚙️ Background Tasks (Automatic)

The system automatically runs:

- **Every 30 minutes**: Extract lore from recent story lines

  - Identifies characters, locations, events, items
  - Creates embeddings for each
  - Stores in vector DB for context retrieval

- **Every hour**: Generate summaries
  - Creates summaries of story chunks (10 lines each)
  - Embeds and stores for high-level context

---

## 🐛 Troubleshooting

### Milvus not connecting?

```bash
# Check Milvus status
docker ps | grep milvus

# Restart Milvus
docker-compose restart

# View logs
docker-compose logs milvus-standalone
```

### Module import errors?

```bash
# Activate venv
source venv/bin/activate

# Reinstall
pip install -r requirements.txt --force-reinstall
```

### Database errors?

```bash
# Delete and recreate
rm kahani.db
# Restart app (auto-creates tables)
python main.py
```

---

## 🎯 Example Story Creation Session

```
User: "A wizard lives in a tower"
AI Suggests: "In a tower of gleaming crystal, high above the misty valleys,
              lived the ancient wizard Eldrin."
User: ✅ Signs → Stored with embedding

User: "A young adventurer arrives"
AI: [Retrieves context: "wizard Eldrin, crystal tower"]
AI Suggests: "One stormy evening, a young adventurer named Aria
              approached the tower's gates, seeking Eldrin's wisdom."
User: Edits to: "One stormy evening, a brave warrior named Kai
                  approached the crystal tower, desperate for help."
User: ✅ Signs → Stored with embedding

User: "Extract Lore"
System:
  - Characters: Eldrin (ancient wizard), Kai (brave warrior)
  - Locations: Crystal tower, misty valleys
  - Events: Kai seeks help from wizard

User: "Finalize Story"
LLM-2: Creates polished canonical version with proper narrative flow
```

---

## 🚀 Ready to Start!

```bash
# 1. Start Milvus
docker-compose up -d

# 2. Activate environment
source venv/bin/activate

# 3. Run the app
python main.py

# 4. Open browser
open http://localhost:8000/ui
```

---

## 💡 Tips

- **Context improves over time** - The more you write, the better suggestions get
- **Edit freely** - AI suggestions are starting points, make them your own
- **Extract lore regularly** - Helps system understand your story world
- **Use specific prompts** - "The hero fights the dragon" better than "something happens"
- **Check /docs** - Interactive API documentation at http://localhost:8000/docs

---

**Happy Storytelling! 🎭📖✨**
