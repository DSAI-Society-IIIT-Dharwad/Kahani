# 📚 Kahani AI - Complete Documentation Index

## 🎯 Quick Links

| What you want to do       | Document to read                   |
| ------------------------- | ---------------------------------- |
| **Get started quickly**   | [INSTALLATION.md](INSTALLATION.md) |
| **Understand the system** | [README.md](README.md)             |
| **Quick reference guide** | [QUICKSTART.md](QUICKSTART.md)     |
| **Test the API**          | [API_TESTING.md](API_TESTING.md)   |
| **Learn architecture**    | [ARCHITECTURE.md](ARCHITECTURE.md) |
| **See project overview**  | [SUMMARY.md](SUMMARY.md)           |

---

## 📖 Documentation Guide

### For First-Time Users

**Start here** → Follow this order:

1. **[INSTALLATION.md](INSTALLATION.md)** (10 minutes)

   - Prerequisites check
   - Step-by-step installation
   - Troubleshooting common issues
   - Verification checklist

2. **[QUICKSTART.md](QUICKSTART.md)** (5 minutes)

   - Manual setup steps
   - Quick start commands
   - First story creation
   - Web UI guide

3. **[README.md](README.md)** (15 minutes)
   - Complete system overview
   - Features and capabilities
   - Configuration options
   - Example workflows

### For Developers

**Best order for understanding the codebase:**

1. **[ARCHITECTURE.md](ARCHITECTURE.md)** (20 minutes)

   - System architecture diagrams
   - Data flow explanations
   - Component interactions
   - Database schemas

2. **[API_TESTING.md](API_TESTING.md)** (15 minutes)

   - Complete API reference
   - cURL examples
   - Python test scripts
   - Postman collection

3. **[SUMMARY.md](SUMMARY.md)** (10 minutes)
   - Project file structure
   - Technology stack
   - Implementation details
   - Quick reference

### For Users

**Quick guides for using the system:**

1. **Web UI Tutorial** → [QUICKSTART.md](QUICKSTART.md#-using-the-web-ui)
2. **API Examples** → [API_TESTING.md](API_TESTING.md)
3. **Example Script** → Run `python example_usage.py`

---

## 📁 File Reference

### Core Application Files

| File          | Purpose       | Lines | Description                             |
| ------------- | ------------- | ----- | --------------------------------------- |
| `main.py`     | API Server    | ~550  | FastAPI application with all endpoints  |
| `config.py`   | Configuration | ~35   | Environment settings management         |
| `database.py` | Database      | ~30   | SQLAlchemy setup and session management |
| `models.py`   | Data Models   | ~80   | Database table definitions              |
| `schemas.py`  | API Schemas   | ~80   | Request/response validation             |

### Service Layer Files

| File                   | Purpose      | Lines | Description                              |
| ---------------------- | ------------ | ----- | ---------------------------------------- |
| `milvus_service.py`    | Vector DB    | ~120  | Milvus client and operations             |
| `embedding_service.py` | Embeddings   | ~45   | Sentence transformer wrapper             |
| `llm_service.py`       | LLM Services | ~200  | Groq Cloud integration (3 LLM roles)     |
| `rag_service.py`       | RAG Pipeline | ~80   | Retrieval-augmented generation           |
| `background_tasks.py`  | Async Tasks  | ~110  | Periodic lore extraction & summarization |

### Frontend & Assets

| File                | Purpose | Description                        |
| ------------------- | ------- | ---------------------------------- |
| `static/index.html` | Web UI  | Beautiful story creation interface |

### Configuration Files

| File                 | Purpose      | Description                        |
| -------------------- | ------------ | ---------------------------------- |
| `.env`               | Environment  | API keys and settings (configured) |
| `.env.example`       | Template     | Example environment file           |
| `requirements.txt`   | Dependencies | Python packages                    |
| `docker-compose.yml` | Milvus Setup | Docker configuration               |

### Documentation Files

| File              | Purpose     | Pages      | Target Audience    |
| ----------------- | ----------- | ---------- | ------------------ |
| `README.md`       | Main Docs   | ~250 lines | Everyone           |
| `INSTALLATION.md` | Setup Guide | ~300 lines | New users          |
| `QUICKSTART.md`   | Quick Guide | ~400 lines | Quick reference    |
| `API_TESTING.md`  | API Docs    | ~300 lines | Developers/Testers |
| `ARCHITECTURE.md` | Design Docs | ~400 lines | Developers         |
| `SUMMARY.md`      | Overview    | ~350 lines | Project managers   |
| `INDEX.md`        | This File   | ~200 lines | Navigation         |

### Utility Files

| File               | Purpose       | Description                        |
| ------------------ | ------------- | ---------------------------------- |
| `start.sh`         | Launch Script | macOS/Linux quick start            |
| `start.bat`        | Launch Script | Windows quick start                |
| `example_usage.py` | Demo Script   | Complete workflow example          |
| `.gitignore`       | Git Config    | Files to ignore in version control |

---

## 🎓 Learning Paths

### Path 1: "I just want to use it"

1. Install → [INSTALLATION.md](INSTALLATION.md)
2. Quick start → [QUICKSTART.md](QUICKSTART.md)
3. Open UI → http://localhost:8000/ui
4. Start creating stories!

**Time**: ~20 minutes

### Path 2: "I want to understand how it works"

1. Read overview → [README.md](README.md)
2. Study architecture → [ARCHITECTURE.md](ARCHITECTURE.md)
3. Review code → Start with `main.py`
4. Run examples → `python example_usage.py`

**Time**: ~1 hour

### Path 3: "I want to build on it"

1. Setup dev environment → [INSTALLATION.md](INSTALLATION.md)
2. Understand API → [API_TESTING.md](API_TESTING.md)
3. Study architecture → [ARCHITECTURE.md](ARCHITECTURE.md)
4. Review services → Read `*_service.py` files
5. Extend features → Modify and test

**Time**: ~2-3 hours

### Path 4: "I'm integrating it"

1. API reference → [API_TESTING.md](API_TESTING.md)
2. Endpoints overview → [README.md](README.md#-api-endpoints)
3. Test examples → Use Postman/cURL
4. Schema definitions → Review `schemas.py`

**Time**: ~30 minutes

---

## 🔍 Find Information Quickly

### How do I...?

| Question                     | Answer Location                                                |
| ---------------------------- | -------------------------------------------------------------- |
| **Install the system?**      | [INSTALLATION.md](INSTALLATION.md)                             |
| **Start the server?**        | [QUICKSTART.md](QUICKSTART.md#-ready-to-start)                 |
| **Create a story?**          | [QUICKSTART.md](QUICKSTART.md#-using-the-web-ui)               |
| **Use the API?**             | [API_TESTING.md](API_TESTING.md)                               |
| **Fix installation issues?** | [INSTALLATION.md](INSTALLATION.md#troubleshooting)             |
| **Understand RAG flow?**     | [ARCHITECTURE.md](ARCHITECTURE.md#-data-flow---story-creation) |
| **Configure settings?**      | [README.md](README.md#-configuration)                          |
| **Extract lore?**            | [API_TESTING.md](API_TESTING.md#8-extract-lore)                |
| **Canonicalize story?**      | [API_TESTING.md](API_TESTING.md#11-canonicalize-story)         |
| **Change LLM model?**        | Edit `.env` → `LLM_MODEL`                                      |
| **Use different DB?**        | Edit `.env` → `DATABASE_URL`                                   |

### Where is...?

| Component              | File Location         |
| ---------------------- | --------------------- |
| **Main API endpoints** | `main.py`             |
| **RAG implementation** | `rag_service.py`      |
| **Vector search**      | `milvus_service.py`   |
| **LLM calls**          | `llm_service.py`      |
| **Database models**    | `models.py`           |
| **Background tasks**   | `background_tasks.py` |
| **Web interface**      | `static/index.html`   |
| **Configuration**      | `config.py` + `.env`  |

---

## 📊 Quick Stats

### Project Size

- **Python files**: 10
- **Total lines of code**: ~1,800
- **Documentation**: ~2,000 lines
- **API endpoints**: 13
- **Database tables**: 3
- **Background tasks**: 2

### Technologies Used

- **Backend**: FastAPI, Python 3.9+
- **Vector DB**: Milvus 2.3
- **LLM**: Groq Cloud (Llama 3.1)
- **Embeddings**: Sentence Transformers
- **Database**: SQLite (SQLAlchemy)
- **Tasks**: APScheduler

---

## 🎯 Common Workflows

### Story Creation Workflow

```
1. User enters prompt
   ↓
2. System retrieves context (RAG)
   ↓
3. LLM-1 generates suggestion
   ↓
4. User edits & signs
   ↓
5. System stores + creates embedding
   ↓
6. Embedding stored in Milvus
```

**Docs**: [ARCHITECTURE.md](ARCHITECTURE.md#-data-flow---story-creation)

### Lore Extraction Workflow

```
1. Select story lines
   ↓
2. LLM extracts entities
   ↓
3. Store in database
   ↓
4. Create embeddings
   ↓
5. Add to vector DB
```

**Docs**: [API_TESTING.md](API_TESTING.md#8-extract-lore)

### Canonicalization Workflow

```
1. Get all verified lines
   ↓
2. LLM-2 transforms to narrative
   ↓
3. Store canonical version
   ↓
4. Return polished story
```

**Docs**: [API_TESTING.md](API_TESTING.md#11-canonicalize-story)

---

## 🔗 External Resources

### Official Documentation

- [FastAPI Docs](https://fastapi.tiangolo.com/)
- [Milvus Docs](https://milvus.io/docs)
- [Groq Cloud](https://console.groq.com/docs)
- [Sentence Transformers](https://www.sbert.net/)

### Tutorials Referenced

- RAG implementation patterns
- Vector database best practices
- LLM prompt engineering

---

## 📞 Getting Help

### 1. Check Documentation

- Installation issues → [INSTALLATION.md](INSTALLATION.md#troubleshooting)
- API questions → [API_TESTING.md](API_TESTING.md)
- Architecture questions → [ARCHITECTURE.md](ARCHITECTURE.md)

### 2. Review Examples

- Run `python example_usage.py`
- Check [API_TESTING.md](API_TESTING.md) for cURL examples
- Visit http://localhost:8000/docs for interactive API docs

### 3. Debug Steps

1. Check logs in terminal
2. Verify Milvus is running: `docker ps`
3. Test health endpoint: `curl http://localhost:8000/health`
4. Check `.env` configuration

---

## 🎨 Visual Guides

### Architecture Diagram

See: [ARCHITECTURE.md](ARCHITECTURE.md#-high-level-architecture)

### Data Flow Diagrams

See: [ARCHITECTURE.md](ARCHITECTURE.md#-data-flow---story-creation)

### Database Schema

See: [ARCHITECTURE.md](ARCHITECTURE.md#-database-schema)

---

## 📝 Cheat Sheet

### Start System

```bash
docker-compose up -d          # Start Milvus
source venv/bin/activate      # Activate Python env
python main.py                # Start server
```

### Test System

```bash
curl http://localhost:8000/health              # Health check
open http://localhost:8000/ui                  # Open UI
python example_usage.py                        # Run example
```

### Stop System

```bash
Ctrl+C                        # Stop FastAPI server
docker-compose down           # Stop Milvus
deactivate                    # Exit Python env
```

### Common Commands

```bash
pip install -r requirements.txt     # Install deps
docker-compose logs                 # View Milvus logs
rm kahani.db                        # Reset database
docker-compose restart              # Restart Milvus
```

---

## 🎉 You're Ready!

Choose your path:

- 🚀 **Quick start** → [INSTALLATION.md](INSTALLATION.md)
- 📖 **Learn more** → [README.md](README.md)
- 💻 **Use API** → [API_TESTING.md](API_TESTING.md)
- 🏗️ **Understand system** → [ARCHITECTURE.md](ARCHITECTURE.md)

**Happy storytelling! 🎭📖✨**

---

_Last updated: October 29, 2025_
_Project: Kahani AI - RAG-Powered Story Generation_
_Version: 1.0.0_
