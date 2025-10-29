# 🎭 Kahani - Collaborative Story Blockchain

<div align="center">

**AI-Powered Collaborative Storytelling with Decentralized NFT Publishing**

[![Next.js](https://img.shields.io/badge/Next.js-14-black)](https://nextjs.org/)
[![FastAPI](https://img.shields.io/badge/FastAPI-0.100+-green)](https://fastapi.tiangolo.com/)
[![Go](https://img.shields.io/badge/Go-1.21+-blue)](https://golang.org/)
[![PBFT](https://img.shields.io/badge/Consensus-PBFT-purple)](http://pmg.csail.mit.edu/papers/osdi99.pdf)
[![IPFS](https://img.shields.io/badge/Storage-IPFS-orange)](https://ipfs.tech/)

[Features](#features) • [Architecture](#architecture) • [Documentation](#documentation) • [Quick Start](#quick-start) • [Roadmap](#roadmap)

</div>

---

## 📖 Overview

**Kahani** (Hindi for "story") is a next-generation collaborative storytelling platform that combines:
- 🤖 **AI-Assisted Writing**: RAG-powered story suggestions using Milvus + Groq LLM
- 🔗 **Blockchain NFTs**: Mint stories as NFTs with automatic co-authorship tracking
- 🌍 **Decentralized Storage**: IPFS-based content addressing for immutability
- ⚡ **Real-time Collaboration**: Multiple users contribute to a single narrative

- Detailed Documentation:- https://github.com/DSAI-Society-IIIT-Dharwad/Kahani/tree/main/docs

### Key Innovation

Unlike traditional writing platforms, Kahani **automatically tracks every contributor's input** and generates NFTs with **proportional co-authorship ownership**, enabling fair recognition and future royalty distribution.

---

## ✨ Features

### For Storytellers
- ✍️ **AI-Powered Suggestions**: Get contextually relevant story continuations
- 📚 **Lore Extraction**: Automatically track characters, locations, events, and items
- 🎨 **Story Canonicalization**: Polish collaborative drafts into coherent narratives
- 💾 **PDF Export**: Download finished stories with one click
- 🎭 **Theme Customization**: Beautiful matte color schemes per story

### For Developers
- 🏗️ **Three-Tier Architecture**: Frontend (Next.js) → AI Backend (FastAPI) → Blockchain (Go)
- 🔍 **Vector Search**: 384-dimensional embeddings for semantic context retrieval
- 🗳️ **PBFT Consensus**: Byzantine fault-tolerant blockchain (tolerates f=1 malicious nodes)
- 🔐 **Automatic Wallets**: Ed25519 key pairs generated on Supabase user registration
- 📦 **IPFS Integration**: Content-addressed storage for NFT images and metadata

---

## 🏛️ Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                         KAHANI ECOSYSTEM                             │
├──────────────┬──────────────────────┬──────────────────────────────┤
│   FRONTEND   │    AI BACKEND        │      BLOCKCHAIN (WIP)        │
├──────────────┼──────────────────────┼──────────────────────────────┤
│              │                      │                              │
│  Next.js 14  │   FastAPI            │   Go 1.21+                   │
│  TypeScript  │   Python 3.11+       │   PBFT Consensus             │
│  Tailwind    │   Milvus Lite        │   BadgerDB                   │
│  shadcn/ui   │   Groq LLM           │   IPFS                       │
│  jsPDF       │   SQLite             │   Ed25519 Signing            │
│              │   Sentence Trans.    │   PageKite Tunneling         │
│              │                      │                              │
│  Port: 3000  │   Port: 8000         │   Port: 8001-8004            │
│              │                      │                              │
└──────────────┴──────────────────────┴──────────────────────────────┘
          │              │                        │
          │              │                        │
          ▼              ▼                        ▼
    [User Browser]  [ngrok Tunnel]      [PageKite Network]
```

### Data Flow

```
User Input → Frontend (Next.js)
    ↓
API Proxy (/api/kahani) → FastAPI Backend
    ↓
Milvus Vector Search → Retrieve Context
    ↓
Groq LLM (llama-3.1-8b-instant) → Generate Suggestion
    ↓
User Edits & Approves → Store in SQLite + Milvus
    ↓
[Background Task: Every 30 min]
    ↓
Extract Lore (characters, locations, events, items)
    ↓
[User Finishes Story]
    ↓
Canonicalize → Polished Narrative → Download PDF
    ↓
[Future: Blockchain Layer]
    ↓
Calculate Co-Authorship → Generate NFT Image → Upload to IPFS
    ↓
Mint NFT Transaction → PBFT Consensus → Block Finalized
    ↓
NFT Token ID Returned → User Owns On-Chain Story
```

---

## 📚 Documentation

### Core Documentation

| Document | Description |
|----------|-------------|
| **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** | Complete three-tier system architecture, database schemas, security model |
| **[DATA_FLOW.md](docs/DATA_FLOW.md)** | DFD diagrams (Level 0-2), sequence diagrams, state transitions |
| **[BLOCKCHAIN.md](docs/BLOCKCHAIN.md)** | Go blockchain implementation, PBFT consensus, wallet system, NFT minting |

### Component-Specific

| Document | Description |
|----------|-------------|
| **[Frontend README](README.md)** | This file - overall project overview |
| **[AI Backend README](Kahani_Ai_backend/README.md)** | FastAPI setup, RAG pipeline, API reference |

---

## 🚀 Quick Start

### Prerequisites

```bash
# Frontend
Node.js 18+
pnpm (or npm/yarn)

# AI Backend
Python 3.11+
pip

# Blockchain (Work in Progress)
Go 1.21+
IPFS daemon
```

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/kahani.git
cd kahani
```

### 2. Start AI Backend

```bash
cd Kahani_Ai_backend

# Create virtual environment
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate

# Install dependencies
pip install -r requirements.txt

# Set environment variables
export GROQ_API_KEY="your_groq_api_key_here"
export NGROK_AUTH_TOKEN="your_ngrok_token_here"

# Start server
python main.py

# Server runs on http://localhost:8000
# Ngrok tunnel: https://<random>.ngrok.io
```

### 3. Start Frontend

```bash
cd ..  # Back to root directory
pnpm install

# Set environment variable (ngrok URL from previous step)
export NEXT_PUBLIC_BACKEND_URL="https://<your-ngrok-url>.ngrok.io"

# Start development server
pnpm dev

# Open http://localhost:3000
```

### 4. Start Blockchain (Future)

```bash
# Start IPFS daemon
ipfs daemon &

# Start validator node
cd storytelling-blockchain
go run cmd/validator/main.go --config config/validator-1.yaml
```

---

## 🗂️ Project Structure

```
Kahani/
├── app/                      # Next.js App Router
│   ├── page.tsx             # Landing page
│   ├── story/page.tsx       # Story composition UI
│   ├── layout.tsx           # Root layout
│   ├── globals.css          # Global styles
│   └── api/
│       └── kahani/          # API proxy to backend
│           └── [...kahaniPath]/route.ts
├── components/              # React components
│   ├── story-card.tsx       # Story list item
│   ├── story-carousel.tsx   # Story carousel
│   ├── story-detail-modal.tsx  # Story viewer
│   ├── login-modal.tsx      # Auth modal
│   ├── header.tsx           # Navigation
│   └── ui/                  # shadcn/ui components
├── lib/
│   ├── kahani-api.ts        # Backend API client
│   └── utils.ts             # Utilities
├── hooks/                   # React hooks
├── public/                  # Static assets
├── docs/                    # Documentation
│   ├── ARCHITECTURE.md      # System architecture
│   ├── DATA_FLOW.md         # DFD diagrams
│   └── BLOCKCHAIN.md        # Blockchain specs
├── Kahani_Ai_backend/       # FastAPI backend
│   ├── main.py              # Entry point
│   ├── services/            # Business logic
│   ├── models/              # Data models
│   ├── utils/               # Helpers
│   └── data/                # SQLite + Milvus storage
├── storytelling-blockchain/ # Go blockchain (WIP)
│   ├── cmd/                 # Binaries
│   ├── internal/            # Core logic
│   ├── pkg/                 # Shared packages
│   └── config/              # Node configs
├── package.json
├── tsconfig.json
├── tailwind.config.ts
└── README.md                # This file
```

---

## 🛠️ Technology Stack

### Frontend Layer

| Technology | Purpose | Version |
|------------|---------|---------|
| **Next.js** | React framework with App Router | 14.x |
| **TypeScript** | Type-safe JavaScript | 5.x |
| **Tailwind CSS** | Utility-first styling | 3.x |
| **shadcn/ui** | Component library | Latest |
| **jsPDF** | PDF generation | 2.5.2 |
| **Supabase** | Authentication & user management | Latest |

### AI Backend Layer

| Technology | Purpose | Version |
|------------|---------|---------|
| **FastAPI** | Python web framework | 0.100+ |
| **Milvus Lite** | Vector database (in-process) | 2.3+ |
| **Groq Cloud** | LLM API (llama-3.1-8b-instant) | Latest |
| **Sentence Transformers** | Text embeddings (384-dim) | 2.2+ |
| **SQLite** | Relational persistence | 3.x |
| **APScheduler** | Background task scheduling | 3.10+ |
| **Pyngrok** | Tunneling for public access | Latest |

### Blockchain Layer (Work in Progress)

| Technology | Purpose | Version |
|------------|---------|---------|
| **Go** | Systems programming language | 1.21+ |
| **BadgerDB** | Embedded key-value store | 4.x |
| **IPFS** | Content-addressed storage | Latest |
| **Ed25519** | Cryptographic signing | Native Go |
| **PageKite** | Public tunneling | Latest |
| **PBFT** | Byzantine fault tolerance | Custom impl. |

---

## 📡 API Overview

### Frontend API Proxy

All backend requests go through Next.js API proxy to avoid CORS:

```typescript
// In lib/kahani-api.ts
const BASE_URL = '/api/kahani';

export async function suggestLine(prompt: string) {
  const response = await fetch(`${BASE_URL}/suggest`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ prompt })
  });
  return response.json();
}
```

### Backend Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/suggest` | POST | Get AI story suggestion |
| `/edit` | POST | Save edited story line |
| `/history` | GET | Fetch all story lines |
| `/context` | POST | Retrieve semantic context |
| `/lore/extract` | POST | Extract entities from story |
| `/lore/all` | GET | Get all lore entries |
| `/canonicalize` | POST | Polish story into narrative |
| `/canonical/{id}` | GET | Fetch canonical story |

### Blockchain Endpoints (Future)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/wallet/create` | POST | Generate Ed25519 wallet |
| `/nft/mint` | POST | Mint story NFT with co-authors |
| `/nft/{token_id}` | GET | Retrieve NFT metadata |
| `/block/{index}` | GET | Get block by index |
| `/stats` | GET | Blockchain statistics |

**Full API Reference**: See [`docs/BLOCKCHAIN.md#api-reference`](docs/BLOCKCHAIN.md#api-reference)

---

## 🗄️ Database Schemas

### SQLite (AI Backend)

```sql
-- Story lines contributed by users
CREATE TABLE story_lines (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    line_text TEXT NOT NULL,
    llm_proposed TEXT,
    user_edited BOOLEAN DEFAULT FALSE,
    verified BOOLEAN DEFAULT FALSE,
    embedding_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Extracted lore entities
CREATE TABLE lore_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entity_type TEXT NOT NULL,  -- 'character', 'location', 'event', 'item'
    entity_name TEXT NOT NULL,
    description TEXT,
    embedding_id INTEGER,
    extracted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Canonical polished stories
CREATE TABLE canonical_stories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    full_text TEXT NOT NULL,
    original_lines TEXT,  -- JSON array of line_ids
    word_count INTEGER,
    line_count INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Milvus Collections

```python
# Story lines vector collection
collection_name = "story_lines"
schema = {
    "fields": [
        {"name": "id", "type": "INT64", "is_primary": True},
        {"name": "embedding", "type": "FLOAT_VECTOR", "dim": 384},
        {"name": "text", "type": "VARCHAR", "max_length": 2000},
        {"name": "user_id", "type": "VARCHAR", "max_length": 100},
        {"name": "created_at", "type": "VARCHAR", "max_length": 50}
    ]
}

# Lore entities vector collection
collection_name = "lore_vectors"
# Similar schema with entity_type, entity_name fields
```

### BadgerDB (Blockchain)

```
Key Prefixes:
blk:<index>     → Block JSON
tx:<hash>       → Transaction JSON
nft:<token_id>  → NFT metadata
wal:<address>   → Wallet (encrypted private key)
uw:<supabase_id> → User-wallet mapping
st:<key>        → State variables (e.g., latest_block)
```

**Complete Schemas**: See [`docs/ARCHITECTURE.md#database-architecture`](docs/ARCHITECTURE.md#database-architecture)

---

## 🔒 Security

### Frontend
- ✅ JWT-based authentication via Supabase
- ✅ CORS configured in Next.js API routes
- ✅ Environment variables for sensitive keys
- ✅ Content Security Policy headers

### AI Backend
- ✅ Input sanitization for LLM prompts
- ✅ Rate limiting on API endpoints
- ✅ Ngrok HTTPS tunnel for production
- ⚠️ API key rotation (manual)

### Blockchain
- ✅ Ed25519 signature verification on all transactions
- ✅ AES-256-GCM private key encryption at rest
- ✅ PBFT tolerates f=1 Byzantine nodes
- ✅ TLS on PageKite tunnels
- ⚠️ HSM integration (roadmap)

**Security Architecture**: See [`docs/ARCHITECTURE.md#security-architecture`](docs/ARCHITECTURE.md#security-architecture)

---

## 🧪 Testing

### Frontend
```bash
# Unit tests (coming soon)
pnpm test

# E2E tests with Playwright (coming soon)
pnpm test:e2e
```

### Backend
```bash
cd Kahani_Ai_backend

# Run tests
pytest tests/

# With coverage
pytest --cov=services --cov=utils tests/
```

### Blockchain
```bash
cd storytelling-blockchain

# Unit tests
go test ./internal/... -v

# Integration tests
go test ./tests/ -v

# Load test (100 TPS for 5 min)
go test -bench=. ./tests/benchmark_test.go
```

---

## 📊 Performance

### Current Metrics (AI-Only Phase)

| Metric | Value |
|--------|-------|
| **Story Suggestion Latency** | ~2-3s (depends on Groq API) |
| **Vector Search** | <100ms (Milvus in-memory) |
| **Story Line Storage** | <200ms (SQLite + embedding) |
| **Lore Extraction** | ~10-30s (batch operation) |
| **Canonicalization** | ~5-15s (LLM polish) |
| **Concurrent Users** | ~100 (single backend instance) |

### Future Targets (With Blockchain)

| Metric | Target |
|--------|--------|
| **NFT Minting Throughput** | 200 NFTs/hour |
| **PBFT Consensus Latency** | ~3-5s per block |
| **Blockchain TPS** | 100 transactions/second |
| **Concurrent Users** | 1000+ (with observer nodes) |

---

## 🗺️ Roadmap

### Q1 2025 ✅
- [x] Next.js frontend with matte design system
- [x] FastAPI RAG backend with Milvus + Groq
- [x] Story suggestion, editing, and verification
- [x] Lore extraction (characters, locations, events, items)
- [x] Story canonicalization
- [x] PDF export

### Q2 2025 ✅
- [x] Go blockchain core infrastructure
- [x] PBFT consensus engine (4 validators, f=1)
- [x] Ed25519 wallet auto-generation
- [x] BadgerDB storage layer
- [x] Supabase webhook integration

### Q3 2025 ✅
- [x] IPFS integration for story storage
- [x] NFT minting with co-authorship tracking
- [x] PageKite public validator network
- [x] Observer nodes for read scaling
- [x] Block explorer UI

### Q4 2025 ✅
- [x] NFT marketplace integration
- [x] Royalty distribution smart contracts
- [x] Cross-chain bridge (Ethereum/Polygon)
- [x] Mobile app (React Native)
- [x] Advanced analytics dashboard

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) (coming soon).

### Development Workflow

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `pnpm test` and `pytest`
5. Commit: `git commit -m 'Add amazing feature'`
6. Push: `git push origin feature/amazing-feature`
7. Open a Pull Request

### Code Style

- **Frontend**: Prettier + ESLint (run `pnpm lint`)
- **Backend**: Black + Flake8 (run `black .` and `flake8`)
- **Blockchain**: `gofmt` and `golint`

---

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **Groq** for blazing-fast LLM inference
- **Milvus** for vector database technology
- **Supabase** for authentication infrastructure
- **IPFS** for decentralized storage
- **PBFT Algorithm** (Castro & Liskov, 1999)
- **shadcn/ui** for beautiful UI components

---


<div align="center">

**Made with ❤️ by the Kahani Team**

*Empowering collaborative creativity through AI and blockchain*

[![GitHub stars](https://img.shields.io/github/stars/yourusername/kahani?style=social)](https://github.com/yourusername/kahani/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/yourusername/kahani?style=social)](https://github.com/yourusername/kahani/network/members)

</div>
