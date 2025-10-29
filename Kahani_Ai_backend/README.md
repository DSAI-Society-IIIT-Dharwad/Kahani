# Kahani AI - RAG-Powered Story Generation System

A sophisticated interactive storytelling system using Retrieval-Augmented Generation (RAG) with Milvus vector database and Groq Cloud LLM.

## 🎭 System Architecture

```
User asks → RAG pulls context → LLM-1 proposes → User picks/edits & signs →
Backend verifies & stores + creates embedding → Periodic extractors & summarizers update vector DB →
When finalizing, LLM-2 canonicalizes
```

## ✨ Features

- **RAG-Based Context Retrieval**: Uses Milvus vector DB to retrieve relevant story context
- **Dual LLM System**:
  - LLM-1: Generates story suggestions based on context
  - LLM-2: Canonicalizes final story into polished narrative
- **Automatic Lore Extraction**: Extracts characters, locations, events, and items
- **Periodic Background Tasks**: Auto-updates vector DB with summaries and lore
- **User Verification**: Sign and verify each story line
- **Embedding Storage**: All content embedded and stored for semantic search
- **RESTful API**: Complete FastAPI backend with auto-generated docs
- **Simple Web UI**: Minimal HTML interface for testing

## 🛠️ Tech Stack

- **FastAPI**: Web framework
- **Milvus**: Vector database
- **Groq Cloud**: LLM provider (Llama models)
- **Sentence Transformers**: Embedding generation
- **SQLAlchemy**: ORM for relational data
- **APScheduler**: Background tasks

## 📋 Prerequisites

### Manual Setup Required:

1. **Install Milvus** (Docker recommended):

```bash
# Using Docker Compose (recommended)
wget https://github.com/milvus-io/milvus/releases/download/v2.3.0/milvus-standalone-docker-compose.yml -O docker-compose.yml
docker-compose up -d

# Or use Milvus Lite (lighter alternative)
pip install milvus

# Verify Milvus is running at localhost:19530
```

2. **Python 3.9+** installed

3. **Groq API Key** (already configured in .env)

## 🚀 Installation & Setup

### Step 1: Install Dependencies

```bash
# Create virtual environment
python -m venv venv
source venv/bin/activate  # On macOS/Linux
# or
venv\Scripts\activate  # On Windows

# Install requirements
pip install -r requirements.txt
```

### Step 2: Configure Environment

The `.env` file is already configured with your Groq API key. Verify settings:

```bash
cat .env
```

### Step 3: Start Milvus (if not running)

```bash
# Check if Milvus is running
curl http://localhost:19530

# If not, start with Docker:
docker-compose up -d
```

### Step 4: Run the Application

```bash
# Start the FastAPI server
python main.py

# Or use uvicorn directly
uvicorn main:app --reload --host 0.0.0.0 --port 8000
```

## 🌐 Access Points

Once running:

- **Web UI**: http://localhost:8000/ui
- **API Docs**: http://localhost:8000/docs
- **ReDoc**: http://localhost:8000/redoc
- **Health Check**: http://localhost:8000/health

## 📚 API Endpoints

### Story Creation Flow

1. **Get AI Suggestion** (LLM-1 + RAG)

   ```
   POST /api/story/suggest
   Body: { "user_prompt": "The hero enters...", "user_id": "user123" }
   ```

2. **Edit & Sign Story Line**

   ```
   POST /api/story/edit
   Body: { "llm_proposed": "...", "final_text": "...", "user_id": "user123" }
   ```

3. **Verify Story Line**
   ```
   POST /api/story/verify/{line_id}
   Body: { "line_id": 1, "signature": "user_signed" }
   ```

### Lore & Context

4. **Extract Lore**

   ```
   POST /api/lore/extract
   Body: { "line_ids": [1, 2, 3] }
   ```

5. **Get All Lore**

   ```
   GET /api/lore/all
   ```

6. **Retrieve Context**
   ```
   POST /api/context/retrieve
   Body: { "query": "castle", "top_k": 5 }
   ```

### Canonicalization

7. **Finalize Story** (LLM-2)

   ```
   POST /api/story/canonicalize
   Body: { "title": "My Story", "line_ids": [1,2,3] }
   ```

8. **Get Canonical Story**
   ```
   GET /api/story/canonical/{story_id}
   ```

### Story Lines

9. **Get All Lines**
   ```
   GET /api/story/lines?verified_only=false
   ```

## 🎨 Using the Web UI

1. Open http://localhost:8000/ui
2. Enter a story prompt (e.g., "A knight discovers a hidden door")
3. Click "Get AI Suggestion"
4. Edit the suggestion if needed
5. Click "Sign & Add to Story"
6. Repeat to build your story
7. Click "Extract Lore" to analyze characters, locations, etc.
8. Click "Finalize Story" to create canonical version

## 🔄 Background Tasks

The system automatically runs periodic tasks:

- **Every 30 minutes**: Extract lore from recent story lines
- **Every hour**: Generate and store summaries of story sections

These tasks update the vector database to improve future context retrieval.

## 📊 Database Schema

### StoryLine

- Stores user-signed story lines
- Tracks LLM proposals vs user edits
- Links to vector embeddings

### LoreEntry

- Characters, locations, events, items
- Extracted from story lines
- Embedded for semantic search

### CanonicalStory

- Final polished versions
- Created by LLM-2
- Preserves original line count

## 🧪 Testing the System

### Example Workflow:

```bash
# 1. Get a suggestion
curl -X POST http://localhost:8000/api/story/suggest \
  -H "Content-Type: application/json" \
  -d '{"user_prompt": "A mysterious traveler arrives at the castle", "user_id": "test_user"}'

# 2. Sign and store
curl -X POST http://localhost:8000/api/story/edit \
  -H "Content-Type: application/json" \
  -d '{"llm_proposed": "The traveler...", "final_text": "The traveler approached the gates...", "user_id": "test_user"}'

# 3. Get all lines
curl http://localhost:8000/api/story/lines

# 4. Extract lore
curl -X POST http://localhost:8000/api/lore/extract \
  -H "Content-Type: application/json" \
  -d '{"line_ids": [1, 2]}'

# 5. Canonicalize
curl -X POST http://localhost:8000/api/story/canonicalize \
  -H "Content-Type: application/json" \
  -d '{"title": "The Castle Mystery"}'
```

## 🔧 Configuration

Edit `.env` to customize:

- `GROQ_API_KEY`: Your Groq Cloud API key
- `MILVUS_HOST`: Milvus server host (default: localhost)
- `MILVUS_PORT`: Milvus port (default: 19530)
- `LLM_MODEL`: Groq model to use (default: llama-3.1-70b-versatile)
- `EMBEDDING_MODEL`: Sentence transformer model
- `MAX_CONTEXT_LINES`: Number of context items to retrieve

## 🐛 Troubleshooting

### Milvus Connection Failed

```bash
# Check if Milvus is running
docker ps | grep milvus

# Restart Milvus
docker-compose restart

# Check logs
docker-compose logs milvus-standalone
```

### Import Errors

```bash
# Reinstall dependencies
pip install -r requirements.txt --force-reinstall
```

### Database Issues

```bash
# Delete and recreate database
rm kahani.db
# Restart the app to auto-create tables
```

## 📝 Project Structure

```
Arnav_Kahani_Ai/
├── main.py                 # FastAPI application
├── config.py               # Configuration management
├── database.py             # Database setup
├── models.py               # SQLAlchemy models
├── schemas.py              # Pydantic schemas
├── milvus_service.py       # Milvus vector DB client
├── embedding_service.py    # Embedding generation
├── llm_service.py          # Groq LLM services
├── rag_service.py          # RAG pipeline
├── background_tasks.py     # Periodic tasks
├── static/
│   └── index.html          # Web UI
├── requirements.txt        # Dependencies
├── .env                    # Environment variables
└── README.md              # This file
```

## 🚀 Next Steps

1. **Start Milvus**: `docker-compose up -d`
2. **Install deps**: `pip install -r requirements.txt`
3. **Run app**: `python main.py`
4. **Open UI**: http://localhost:8000/ui
5. **Start writing**: Create your first story!

## 📄 License

MIT License - feel free to use and modify!

## 🤝 Contributing

Built with ❤️ using FastAPI, Milvus, and Groq Cloud.

---

**Happy Storytelling! 🎭📖✨**
