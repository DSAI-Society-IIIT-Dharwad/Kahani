# System Architecture

## 🏗️ High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         USER INTERFACE                           │
│  ┌──────────────────┐              ┌──────────────────┐         │
│  │   Web UI (HTML)  │              │   REST API       │         │
│  │  - Story Writer  │              │   - cURL         │         │
│  │  - Lore Viewer   │              │   - Postman      │         │
│  └────────┬─────────┘              └────────┬─────────┘         │
└───────────┼────────────────────────────────┼───────────────────┘
            │                                │
            └────────────────┬───────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      FastAPI APPLICATION                         │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                    main.py (API Layer)                    │  │
│  │  • /api/story/suggest      • /api/lore/extract           │  │
│  │  • /api/story/edit         • /api/lore/all               │  │
│  │  • /api/story/verify       • /api/context/retrieve       │  │
│  │  • /api/story/canonicalize • /health                     │  │
│  └──────────────────────────────────────────────────────────┘  │
│              │              │              │                    │
│              ▼              ▼              ▼                    │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ RAG Service │  │  LLM Service │  │ Embedding    │          │
│  │             │  │              │  │ Service      │          │
│  │ • retrieve  │  │ • LLM-1      │  │              │          │
│  │ • generate  │  │ • LLM-2      │  │ • generate   │          │
│  │             │  │ • extract    │  │ • batch      │          │
│  └─────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
              │              │              │
              ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     EXTERNAL SERVICES                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Milvus DB   │  │  Groq Cloud  │  │  SQLite DB   │          │
│  │              │  │              │  │              │          │
│  │ • Vectors    │  │ • Llama 3.1  │  │ • StoryLine  │          │
│  │ • Search     │  │ • Context.   │  │ • LoreEntry  │          │
│  │ • IVF Index  │  │ • Creative   │  │ • Canonical  │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   BACKGROUND TASKS                               │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              APScheduler (background_tasks.py)            │  │
│  │                                                           │  │
│  │  ⏱️  Every 30 min: Extract Lore                           │  │
│  │  ⏱️  Every 60 min: Generate Summaries                     │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📊 Data Flow - Story Creation

```
┌──────────┐
│  USER    │
│  "A hero │
│  enters" │
└────┬─────┘
     │
     │ 1. POST /api/story/suggest
     ▼
┌─────────────────┐
│  RAG Service    │
│                 │
│  ┌───────────┐  │
│  │ Embedding │  │ ← 2. Convert prompt to vector
│  └─────┬─────┘  │
│        │        │
│        ▼        │
│  ┌───────────┐  │
│  │  Milvus   │  │ ← 3. Search similar vectors
│  │  Search   │  │    (story_line, lore, summary)
│  └─────┬─────┘  │
└────────┼────────┘
         │
         │ 4. Retrieved context (top 5-10)
         ▼
┌─────────────────┐
│  LLM Service    │
│  (Groq Cloud)   │
│                 │
│  ┌───────────┐  │
│  │  LLM-1    │  │ ← 5. Generate suggestion
│  │  Llama    │  │    with context
│  └─────┬─────┘  │
└────────┼────────┘
         │
         │ 6. Return suggestion + context
         ▼
┌──────────────────┐
│  USER            │
│  Reviews & Edits │
│  Signs ✅        │
└────┬─────────────┘
     │
     │ 7. POST /api/story/edit
     ▼
┌─────────────────────────────────────┐
│  Backend (Background Task)          │
│                                     │
│  ┌──────────────┐                   │
│  │  SQLite DB   │ ← 8. Store line  │
│  │  StoryLine   │                   │
│  └──────────────┘                   │
│                                     │
│  ┌──────────────┐                   │
│  │  Embedding   │ ← 9. Generate    │
│  │  Service     │     embedding    │
│  └─────┬────────┘                   │
│        │                            │
│        ▼                            │
│  ┌──────────────┐                   │
│  │  Milvus DB   │ ← 10. Store      │
│  │  Insert      │      vector      │
│  └──────────────┘                   │
└─────────────────────────────────────┘
```

---

## 🔄 Background Processing Flow

```
┌─────────────────────────────────────────────┐
│        APScheduler Triggers                 │
│                                             │
│  ⏱️  Every 30 minutes                        │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌──────────────────────────────────────┐
│  Lore Extraction Task                │
│                                      │
│  1. Get recent story lines (20)     │
│     ↓                                │
│  2. Extract with LLM:                │
│     • Characters                     │
│     • Locations                      │
│     • Events                         │
│     • Items                          │
│     ↓                                │
│  3. Create embeddings                │
│     ↓                                │
│  4. Store in:                        │
│     • SQLite (LoreEntry)             │
│     • Milvus (vectors)               │
└──────────────────────────────────────┘

┌──────────────────────────────────────┐
│  Summary Generation Task             │
│                                      │
│  1. Get all verified lines           │
│     ↓                                │
│  2. Chunk into groups (10 lines)    │
│     ↓                                │
│  3. Generate summary per chunk       │
│     ↓                                │
│  4. Create embedding                 │
│     ↓                                │
│  5. Store in Milvus                  │
│     (content_type: "summary")        │
└──────────────────────────────────────┘
```

---

## 🎭 LLM Roles

### LLM-1: Story Suggester

```
Input:
  • User prompt
  • Retrieved context from RAG

Process:
  • Model: llama-3.1-70b-versatile
  • Temperature: 0.7 (creative)
  • Max tokens: 200

Output:
  • Short story line(s) (1-3 sentences)
  • Contextually relevant
  • Natural continuation
```

### LLM-2: Canonicalizer

```
Input:
  • All verified story lines

Process:
  • Model: llama-3.1-70b-versatile
  • Temperature: 0.5 (balanced)
  • Max tokens: 2000

Output:
  • Polished narrative
  • Proper paragraphs
  • Fixed inconsistencies
  • Professional quality
```

### LLM Extractor

```
Input:
  • Story lines to analyze

Process:
  • Model: llama-3.1-70b-versatile
  • Temperature: 0.3 (precise)
  • Output format: JSON

Output:
  • Structured lore data
  • Characters, locations, events, items
  • With descriptions
```

---

## 🗄️ Database Schema

### SQLite Tables

```sql
-- Story Lines
CREATE TABLE story_lines (
    id INTEGER PRIMARY KEY,
    user_id VARCHAR,
    line_text TEXT NOT NULL,
    line_number INTEGER,
    context_used TEXT,
    llm_proposed TEXT,
    user_edited BOOLEAN,
    verified BOOLEAN,
    signature VARCHAR,
    embedding_id VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Lore Entries
CREATE TABLE lore_entries (
    id INTEGER PRIMARY KEY,
    entity_type VARCHAR,  -- character, location, event, item
    entity_name VARCHAR,
    description TEXT,
    source_lines TEXT,
    confidence FLOAT,
    embedding_id VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Canonical Stories
CREATE TABLE canonical_stories (
    id INTEGER PRIMARY KEY,
    title VARCHAR,
    full_text TEXT NOT NULL,
    canonicalized_by VARCHAR,
    original_lines_count INTEGER,
    version INTEGER,
    created_at TIMESTAMP,
    finalized_at TIMESTAMP
);
```

### Milvus Collection Schema

```python
Collection: "story_embeddings"

Fields:
  - id: INT64 (primary, auto_id)
  - embedding_id: VARCHAR(100)
  - text: VARCHAR(5000)
  - content_type: VARCHAR(50)  # story_line, lore_character, summary
  - embedding: FLOAT_VECTOR(384)  # all-MiniLM-L6-v2

Index:
  - Type: IVF_FLAT
  - Metric: L2 distance
  - nlist: 128
  - nprobe: 10
```

---

## 🔐 Security & Verification

```
User signs story line
        ↓
Backend creates signature:
  SHA256(line_text)[:16]
        ↓
Stores in database:
  - verified = True
  - signature = hash
        ↓
Future verification:
  - Compare stored hash
  - Ensure integrity
```

---

## 📈 Scalability Considerations

### Current Setup (Single Instance)

- ✅ Perfect for personal use
- ✅ 1-10 concurrent users
- ✅ Thousands of story lines

### Future Scaling Options

- 🔄 PostgreSQL instead of SQLite
- 🔄 Milvus cluster (distributed)
- 🔄 Redis for caching
- 🔄 Celery for background tasks
- 🔄 Load balancer for API
- 🔄 Vector DB sharding

---

## 🎯 Performance Optimization

### Vector Search

- **Index Type**: IVF_FLAT (good for <1M vectors)
- **For larger**: Switch to IVF_PQ or HNSW
- **nprobe**: Balance speed vs accuracy

### Embedding Generation

- **Batch processing**: Process multiple texts at once
- **Caching**: Cache embeddings for repeated queries
- **Model**: all-MiniLM-L6-v2 (fast, 384 dim)

### LLM Calls

- **Groq**: Ultra-fast inference (~500 tokens/sec)
- **Context limits**: Keep under 4096 tokens
- **Batching**: Queue multiple requests

---

**Architecture Complete! 🏗️✨**
