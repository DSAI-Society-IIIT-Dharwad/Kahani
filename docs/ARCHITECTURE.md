# 🏗️ Kahani System Architecture

## System Overview

Kahani is a three-tier collaborative storytelling platform combining AI-powered narrative generation with blockchain-based ownership and verification.

```
┌─────────────────────────────────────────────────────────────────┐
│                     KAHANI ECOSYSTEM                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   Next.js    │───▶│  FastAPI     │───▶│  Blockchain  │      │
│  │   Frontend   │    │  RAG Engine  │    │   Network    │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
│        │                    │                    │               │
│        │                    │                    │               │
│   [User Auth]         [Story Gen]          [NFT Minting]        │
│   [Story UI]          [Lore Extract]       [Wallet Mgmt]        │
│   [PDF Export]        [Vector DB]          [PBFT Consensus]     │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🎯 Three-Layer Architecture

### Layer 1: Frontend (Next.js 14 + React)

**Technology Stack:**
- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS + shadcn/ui
- Radix UI components
- jsPDF (export)

**Responsibilities:**
- User authentication via Supabase
- Story composition interface
- AI suggestion integration
- Real-time collaboration UI
- PDF generation and download
- Wallet display (read-only)
- NFT gallery visualization

**Key Routes:**
```
/                      # Landing page with story carousel
/story?id={n}          # Collaborative story board
/api/kahani/*          # Proxy to FastAPI backend
```

**State Management:**
```typescript
Story State {
  sentences: StorySentence[]
  currentPlayer: ActivePlayer
  suggestedLine: string | null
  canonicalStory: CanonicalStoryResponse | null
}
```

---

### Layer 2: AI Backend (FastAPI + RAG)

**Technology Stack:**
- FastAPI (async Python)
- Milvus Lite (vector DB)
- Groq Cloud (LLM inference)
- Sentence Transformers (embeddings)
- SQLite (persistent storage)
- APScheduler (background tasks)

**Responsibilities:**
- RAG-powered story suggestions
- Context retrieval from vector DB
- Lore extraction (characters, locations, events, items)
- Story canonicalization (polishing)
- Embedding generation (384-dim vectors)
- Scheduled knowledge base updates

**Core Services:**

```python
RAG Pipeline:
  User Prompt
    ↓
  Embedding (sentence-transformers)
    ↓
  Vector Search (Milvus, top-k=10)
    ↓
  Context Assembly
    ↓
  LLM Generation (Groq llama-3.1-8b-instant)
    ↓
  Response + Context Metadata
```

**Data Flow:**
```
POST /api/story/suggest → RAG → LLM-1 → Suggestion
POST /api/story/edit    → Store → Embedding → Milvus
POST /api/lore/extract  → LLM   → Extract Entities
POST /api/story/canonicalize → LLM-2 → Polish
```

**Background Tasks:**
- Every 30 min: Lore extraction from verified lines
- Every 60 min: Story summarization
- Every 5 min: Health checks

---

### Layer 3: Blockchain Network (Go + PBFT Consensus)

> **Status:** Work in Progress

**Technology Stack:**
- Go 1.21+
- BadgerDB (embedded storage)
- IPFS (content addressing)
- Ed25519 (cryptographic signing)
- PBFT consensus algorithm
- PageKite (tunneling)

**Responsibilities:**
- Immutable story contribution ledger
- Automatic wallet generation per Supabase user
- NFT minting with co-authorship tracking
- Decentralized consensus (4+ validator nodes)
- GitHub-style contribution visualization
- On-chain ownership percentages

**Node Types:**

**Validator Nodes (4 minimum):**
- Run PBFT consensus
- Create and validate blocks
- Store full blockchain
- Expose REST API via PageKite
- Poll Supabase for new users

**Observer Nodes (optional):**
- Read-only blockchain access
- Cache popular content
- Geo-distributed read replicas
- WebSocket real-time updates

---

## 🔄 Complete Data Flow

### Story Creation Flow

```
┌─────────────┐
│   User      │
│  (Browser)  │
└──────┬──────┘
       │ 1. Enters prompt
       ▼
┌─────────────────┐
│  Next.js UI     │
│  /story         │
└────────┬────────┘
         │ 2. POST /api/kahani/story/suggest
         ▼
┌──────────────────────┐
│  Next.js API Proxy   │
│  [...kahaniPath]     │
└──────────┬───────────┘
           │ 3. Forward to FastAPI
           ▼
┌────────────────────────────┐
│  FastAPI RAG Service       │
│  ┌──────────────────────┐  │
│  │  Embedding Generator │  │
│  └──────────┬───────────┘  │
│             ▼              │
│  ┌──────────────────────┐  │
│  │  Milvus Vector DB    │  │
│  │  (Semantic Search)   │  │
│  └──────────┬───────────┘  │
│             ▼              │
│  ┌──────────────────────┐  │
│  │  Context Assembly    │  │
│  │  (Top-K=10)          │  │
│  └──────────┬───────────┘  │
│             ▼              │
│  ┌──────────────────────┐  │
│  │  Groq LLM (Llama)    │  │
│  │  Generate Suggestion │  │
│  └──────────┬───────────┘  │
└─────────────┼──────────────┘
              │ 4. Return suggestion
              ▼
       ┌──────────────┐
       │  Next.js UI  │
       │  Display     │
       └──────┬───────┘
              │ 5. User edits & signs
              ▼
       ┌──────────────────┐
       │  POST story/edit │
       └──────┬───────────┘
              │ 6. Store verified line
              ▼
       ┌────────────────────┐
       │  SQLite + Milvus   │
       │  (Future context)  │
       └────────┬───────────┘
                │ 7. Optional: Blockchain
                ▼
       ┌─────────────────────┐
       │  Validator Network  │
       │  (Store immutably)  │
       └─────────────────────┘
```

### NFT Minting Flow (Blockchain Layer)

```
┌──────────────┐
│  User        │
│  (Frontend)  │
└──────┬───────┘
       │ 1. Click "Finish Story"
       ▼
┌────────────────────┐
│  POST              │
│  /story/canonicalize│
└────────┬───────────┘
         │ 2. Aggregate contributions
         ▼
┌─────────────────────────┐
│  FastAPI Backend        │
│  ┌───────────────────┐  │
│  │  LLM-2 Polish     │  │
│  │  Create Narrative │  │
│  └─────────┬─────────┘  │
└────────────┼────────────┘
             │ 3. Return canonical story
             ▼
┌──────────────────────────┐
│  Next.js UI              │
│  ┌────────────────────┐  │
│  │  jsPDF Generator   │  │
│  │  Download PDF      │  │
│  └────────────────────┘  │
└──────────────────────────┘

       [Future: Blockchain NFT]
             │
             ▼
┌─────────────────────────────────┐
│  POST /api/story/:id/mint       │
│  (Blockchain Network)           │
│  ┌───────────────────────────┐  │
│  │  Aggregate Contributions  │  │
│  │  ├─ Main Author (most)    │  │
│  │  └─ Co-authors (%)        │  │
│  └─────────┬─────────────────┘  │
│            ▼                     │
│  ┌───────────────────────────┐  │
│  │  Generate NFT Image       │  │
│  │  (IPFS Upload)            │  │
│  └─────────┬─────────────────┘  │
│            ▼                     │
│  ┌───────────────────────────┐  │
│  │  Create Metadata JSON     │  │
│  │  (Authors, Ownership %)   │  │
│  └─────────┬─────────────────┘  │
│            ▼                     │
│  ┌───────────────────────────┐  │
│  │  mint_nft Transaction     │  │
│  │  PBFT Consensus           │  │
│  └─────────┬─────────────────┘  │
│            ▼                     │
│  ┌───────────────────────────┐  │
│  │  Block Committed          │  │
│  │  NFT Token ID Returned    │  │
│  └───────────────────────────┘  │
└─────────────────────────────────┘
```

---

## 📊 Database & Storage Architecture

### Frontend State
```
Browser LocalStorage/Session:
  - User session (Supabase JWT)
  - Draft story lines
  - UI preferences
```

### AI Backend Storage

**SQLite Schema:**
```sql
story_lines (
  id INTEGER PRIMARY KEY,
  user_id TEXT,
  llm_proposed TEXT,
  line_text TEXT,
  user_edited BOOLEAN,
  verified BOOLEAN,
  embedding_id TEXT,
  created_at TIMESTAMP
)

lore_entries (
  id INTEGER PRIMARY KEY,
  entity_type TEXT, -- character/location/event/item
  entity_name TEXT,
  description TEXT,
  embedding_id TEXT,
  source_line_ids TEXT -- JSON array
)

canonical_stories (
  id INTEGER PRIMARY KEY,
  title TEXT,
  full_text TEXT,
  original_lines_count INTEGER,
  created_at TIMESTAMP
)
```

**Milvus Collections:**
```
story_embeddings:
  - vector: [384 float32]
  - metadata: {line_id, user_id, timestamp}

lore_embeddings:
  - vector: [384 float32]
  - metadata: {entity_name, entity_type}
```

### Blockchain Storage (Future)

**BadgerDB (Key-Value):**
```
blocks/{index} → Block JSON
wallets/{address} → Wallet JSON
nfts/{tokenID} → NFT JSON
state/latest → Current state snapshot
```

**IPFS (Content-Addressed):**
```
Story text files    → QmXXXX...
NFT metadata JSON   → QmYYYY...
NFT images (PNG)    → QmZZZZ...
```

---

## 🔐 Security Architecture

### Frontend Layer
- Supabase Auth (JWT tokens)
- HTTPS-only API calls
- CORS policy enforcement
- XSS protection (React escaping)
- Rate limiting (proxy level)

### AI Backend Layer
```python
Security Measures:
  ├─ JWT verification middleware
  ├─ Input sanitization (prompts)
  ├─ API rate limiting (100 req/min)
  ├─ Groq API key rotation
  └─ HTTPS required (production)
```

### Blockchain Layer (Future)
```go
Cryptographic Stack:
  ├─ Ed25519 signatures (all transactions)
  ├─ SHA-256 block hashing
  ├─ AES-256-GCM wallet encryption
  ├─ PBFT consensus (Byzantine fault tolerance)
  └─ Immutable audit trail
```

---

## 🌐 Network Topology

### Current Deployment (AI System)

```
        Internet
           │
           ▼
    ┌──────────────┐
    │   ngrok      │
    │   Tunnel     │
    └──────┬───────┘
           │
           ▼
    ┌──────────────┐
    │  FastAPI     │
    │  Backend     │
    │  Port: 8000  │
    └──────┬───────┘
           │
           ▼
    ┌──────────────┐
    │  Milvus Lite │
    │  (Embedded)  │
    └──────────────┘
```

### Future Blockchain Network

```
                  Internet
                     │
          ┌──────────┼──────────┐
          │          │          │
          ▼          ▼          ▼
    ┌─────────┐ ┌─────────┐ ┌─────────┐
    │Validator│ │Validator│ │Validator│
    │   #1    │ │   #2    │ │   #3    │
    │(Leader) │ │         │ │         │
    └────┬────┘ └────┬────┘ └────┬────┘
         │           │           │
         └───────────┼───────────┘
              PBFT Consensus
         (Pre-Prepare → Prepare → Commit)
                     │
                     ▼
              ┌────────────┐
              │  Observer  │
              │   Nodes    │
              │ (Regional) │
              └────────────┘
```

---

## 📈 Scalability Considerations

### Frontend Scalability
- Static generation where possible
- CDN distribution (Vercel)
- Code splitting per route
- Image optimization (Next.js)
- Progressive Web App (future)

### AI Backend Scalability
```
Current: Single Instance
  └─ Handles ~100 concurrent users
  
Future: Horizontal Scaling
  ├─ Load Balancer (nginx)
  ├─ Multiple FastAPI instances
  ├─ Shared Milvus cluster
  └─ Redis cache layer
```

### Blockchain Scalability
```
PBFT Limitations:
  ├─ Throughput: ~1000 TPS (4 validators)
  ├─ Latency: ~3-5 seconds (consensus delay)
  └─ Max validators: ~10 (network overhead)
  
Mitigation Strategies:
  ├─ Batch transactions (100 per block)
  ├─ Observer nodes (read scaling)
  └─ IPFS (offload large content)
```

---

## 🔄 State Synchronization

### Frontend ↔ Backend Sync
```javascript
// Optimistic UI updates
const handleSubmit = async (text) => {
  // 1. Update UI immediately
  setSentences(prev => [...prev, newSentence])
  
  // 2. Send to backend
  try {
    await editStoryLine({...})
    // 3. Fetch canonical state
    const lines = await fetchStoryLines()
    setSentences(transformStoryLines(lines))
  } catch (error) {
    // 4. Rollback on error
    setSentences(prev => prev.slice(0, -1))
  }
}
```

### Backend ↔ Blockchain Sync (Future)
```go
// Periodic blockchain sync
func (v *ValidatorNode) SyncWithBackend() {
    ticker := time.NewTicker(60 * time.Second)
    for range ticker.C {
        // Fetch verified lines from FastAPI
        lines := fetchVerifiedLines()
        
        // Create transactions
        txs := convertToTransactions(lines)
        
        // Propose block
        v.consensus.ProposeBlock(txs)
    }
}
```

---

## 🛠️ Development & Deployment

### Local Development Setup
```bash
# 1. Frontend
cd Kahani
pnpm install
pnpm dev  # http://localhost:3000

# 2. AI Backend
cd Kahani_Ai_backend
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python main.py  # http://localhost:8000

# 3. Blockchain (future)
cd blockchain
go mod download
go run cmd/validator/main.go
```

### Production Deployment

**Frontend (Vercel):**
```bash
vercel --prod
# Auto-deploy on git push to main
```

**Backend (Railway/Render):**
```dockerfile
FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY . .
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

**Blockchain (Docker + PageKite):**
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o validator cmd/validator/main.go
CMD ["./validator"]
```

---

## 📊 Monitoring & Observability

### Health Checks
```
GET /health (Frontend Next.js)
  └─ Returns: API connectivity status

GET /health (AI Backend)
  ├─ Database: SQLite connection
  ├─ Vector DB: Milvus ping
  ├─ LLM: Groq API status
  └─ Storage: Disk space

GET /health (Blockchain - future)
  ├─ Peer count
  ├─ Consensus state
  ├─ Block height
  └─ Wallet registry size
```

### Metrics Collection (Future)
```
Prometheus metrics:
  ├─ kahani_stories_created_total
  ├─ kahani_llm_latency_seconds
  ├─ kahani_rag_context_retrieved
  ├─ kahani_nft_minted_total
  └─ kahani_blocks_committed_total
```

---

## 🔮 Future Enhancements

### Phase 1 (Q1 2026)
- [ ] Blockchain validator network (4 nodes)
- [ ] Automatic wallet generation
- [ ] NFT minting with co-authorship
- [ ] IPFS integration

### Phase 2 (Q2 2026)
- [ ] Observer nodes for read scaling
- [ ] WebSocket real-time updates
- [ ] Advanced lore graph visualization
- [ ] Multi-language support

### Phase 3 (Q3 2026)
- [ ] Mobile app (React Native)
- [ ] Collaborative editing (OT/CRDT)
- [ ] Marketplace for NFTs
- [ ] DAO governance for story curation

---

## 📚 References

- [Next.js Documentation](https://nextjs.org/docs)
- [FastAPI Best Practices](https://fastapi.tiangolo.com/)
- [Milvus Vector Database](https://milvus.io/docs)
- [PBFT Consensus Paper](https://pmg.csail.mit.edu/papers/osdi99.pdf)
- [IPFS Documentation](https://docs.ipfs.tech/)
