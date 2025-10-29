# 📐 Kahani Visual Diagrams

This document contains Mermaid.js diagrams that can be rendered in GitHub, VS Code, or any Mermaid-compatible viewer.

---

## System Context Diagram

```mermaid
graph TB
    User[👤 Storyteller] -->|Write Stories| Frontend
    Frontend[Next.js Frontend] -->|API Requests| Backend
    Backend[FastAPI Backend] -->|Vector Search| Milvus
    Backend -->|Generate Text| Groq[Groq LLM]
    Backend -->|Store Data| SQLite[(SQLite DB)]
    Backend -->|Future: Mint NFT| Blockchain
    Blockchain[Go Blockchain] -->|Store Content| IPFS
    Blockchain -->|Consensus| PBFT[PBFT Network]
    Blockchain -->|Persist State| BadgerDB[(BadgerDB)]
    
    style Frontend fill:#3b82f6,color:#fff
    style Backend fill:#10b981,color:#fff
    style Blockchain fill:#8b5cf6,color:#fff
    style Milvus fill:#f59e0b,color:#fff
    style Groq fill:#ef4444,color:#fff
```

---

## Component Architecture

```mermaid
graph LR
    subgraph Frontend["Frontend Layer (Port 3000)"]
        UI[React Components]
        API[API Proxy]
        State[State Management]
    end
    
    subgraph Backend["AI Backend (Port 8000)"]
        FastAPI[FastAPI Server]
        RAG[RAG Pipeline]
        Scheduler[Background Tasks]
        Vector[Vector Embeddings]
    end
    
    subgraph Blockchain["Blockchain Layer (Ports 8001-8004)"]
        V1[Validator 1]
        V2[Validator 2]
        V3[Validator 3]
        V4[Validator 4]
        PBFT_Consensus[PBFT Consensus]
    end
    
    subgraph Storage["Data Layer"]
        Milvus[(Milvus Vectors)]
        SQLite[(SQLite)]
        BadgerDB[(BadgerDB)]
        IPFS[(IPFS)]
    end
    
    UI --> API
    API --> FastAPI
    FastAPI --> RAG
    RAG --> Vector
    Vector --> Milvus
    FastAPI --> SQLite
    Scheduler --> SQLite
    
    V1 --> PBFT_Consensus
    V2 --> PBFT_Consensus
    V3 --> PBFT_Consensus
    V4 --> PBFT_Consensus
    PBFT_Consensus --> BadgerDB
    PBFT_Consensus --> IPFS
    
    style Frontend fill:#3b82f6,color:#fff
    style Backend fill:#10b981,color:#fff
    style Blockchain fill:#8b5cf6,color:#fff
    style Storage fill:#f59e0b,color:#fff
```

---

## Story Creation Sequence Diagram

```mermaid
sequenceDiagram
    actor User
    participant Frontend
    participant APIProxy as API Proxy
    participant FastAPI
    participant Milvus
    participant Groq
    participant SQLite

    User->>Frontend: Enter prompt
    Frontend->>APIProxy: POST /api/kahani/suggest
    APIProxy->>FastAPI: Forward request
    FastAPI->>FastAPI: Generate embedding (384-dim)
    FastAPI->>Milvus: Vector search (top_k=10)
    Milvus-->>FastAPI: Relevant stories
    FastAPI->>Groq: Generate with context
    Groq-->>FastAPI: AI suggestion
    FastAPI-->>APIProxy: Response JSON
    APIProxy-->>Frontend: Story suggestion
    Frontend-->>User: Display suggestion
    
    User->>Frontend: Edit & approve
    Frontend->>APIProxy: POST /api/kahani/edit
    APIProxy->>FastAPI: Save story line
    FastAPI->>SQLite: INSERT INTO story_lines
    SQLite-->>FastAPI: ID assigned
    FastAPI->>Milvus: Store embedding
    Milvus-->>FastAPI: Success
    FastAPI-->>APIProxy: Stored
    APIProxy-->>Frontend: Success
    Frontend-->>User: Updated story list
```

---

## Lore Extraction Flow

```mermaid
sequenceDiagram
    participant Scheduler
    participant FastAPI
    participant SQLite
    participant Groq
    participant Milvus

    Note over Scheduler: Every 30 minutes
    Scheduler->>FastAPI: Trigger lore extraction
    FastAPI->>SQLite: SELECT verified lines
    SQLite-->>FastAPI: Story texts
    FastAPI->>FastAPI: Aggregate text
    FastAPI->>Groq: Extract entities prompt
    Groq-->>FastAPI: Entities JSON
    FastAPI->>FastAPI: Parse & validate
    FastAPI->>SQLite: INSERT lore_entries
    SQLite-->>FastAPI: IDs assigned
    FastAPI->>Milvus: Store entity embeddings
    Milvus-->>FastAPI: Success
    FastAPI-->>Scheduler: Complete
```

---

## NFT Minting Sequence (Future)

```mermaid
sequenceDiagram
    actor User
    participant Frontend
    participant Blockchain
    participant AIBackend as AI Backend
    participant IPFS
    participant PBFT
    participant BadgerDB

    User->>Frontend: Click "Mint NFT"
    Frontend->>Blockchain: POST /nft/mint {story_id}
    Blockchain->>AIBackend: GET canonical story
    AIBackend-->>Blockchain: Story + contributions
    Blockchain->>Blockchain: Calculate co-authorship
    Blockchain->>Blockchain: Generate NFT image
    Blockchain->>IPFS: Upload image
    IPFS-->>Blockchain: Image CID
    Blockchain->>Blockchain: Create metadata JSON
    Blockchain->>IPFS: Upload metadata
    IPFS-->>Blockchain: Metadata CID
    Blockchain->>Blockchain: Create mint_nft tx
    Blockchain->>PBFT: Submit transaction
    
    Note over PBFT: PBFT Consensus
    PBFT->>PBFT: Pre-Prepare
    PBFT->>PBFT: Prepare (2f+1 votes)
    PBFT->>PBFT: Commit (2f+1 votes)
    PBFT->>BadgerDB: Finalize block
    
    BadgerDB-->>Blockchain: Block committed
    Blockchain-->>Frontend: NFT minted {token_id}
    Frontend-->>User: Success + token_id
```

---

## PBFT Consensus Flow

```mermaid
sequenceDiagram
    participant Leader as Leader (V1)
    participant V2 as Validator 2
    participant V3 as Validator 3
    participant V4 as Validator 4

    Note over Leader: Has new transaction
    Leader->>Leader: Create block
    
    Leader->>V2: PRE-PREPARE {block}
    Leader->>V3: PRE-PREPARE {block}
    Leader->>V4: PRE-PREPARE {block}
    
    V2->>Leader: PREPARE
    V2->>V3: PREPARE
    V2->>V4: PREPARE
    
    V3->>Leader: PREPARE
    V3->>V2: PREPARE
    V3->>V4: PREPARE
    
    V4->>Leader: PREPARE
    V4->>V2: PREPARE
    V4->>V3: PREPARE
    
    Note over Leader,V4: 2f+1 = 3 PREPARE msgs received
    
    Leader->>V2: COMMIT
    Leader->>V3: COMMIT
    Leader->>V4: COMMIT
    
    V2->>Leader: COMMIT
    V2->>V3: COMMIT
    V2->>V4: COMMIT
    
    V3->>Leader: COMMIT
    V3->>V2: COMMIT
    V3->>V4: COMMIT
    
    V4->>Leader: COMMIT
    V4->>V2: COMMIT
    V4->>V3: COMMIT
    
    Note over Leader,V4: 2f+1 COMMIT msgs received
    Note over Leader,V4: Block finalized on all nodes
```

---

## Data Model ER Diagram

```mermaid
erDiagram
    USER ||--o{ STORY_LINE : contributes
    USER ||--o{ WALLET : owns
    USER {
        string supabase_id PK
        string username
        string email
    }
    
    STORY_LINE ||--o{ CANONICAL_STORY : "includes in"
    STORY_LINE {
        int id PK
        string user_id FK
        text line_text
        text llm_proposed
        boolean user_edited
        boolean verified
        int embedding_id FK
        timestamp created_at
    }
    
    CANONICAL_STORY ||--o{ NFT : "minted as"
    CANONICAL_STORY {
        int id PK
        string title
        text full_text
        json original_lines
        int word_count
        int line_count
        timestamp created_at
    }
    
    LORE_ENTRY {
        int id PK
        string entity_type
        string entity_name
        text description
        int embedding_id FK
        timestamp extracted_at
    }
    
    WALLET ||--o{ TRANSACTION : signs
    WALLET {
        string address PK
        bytes public_key
        bytes private_key_encrypted
        string supabase_user_id FK
        timestamp created_at
    }
    
    TRANSACTION }o--|| BLOCK : "included in"
    TRANSACTION {
        string tx_id PK
        string type
        string from FK
        string to
        json data
        string signature
        timestamp created_at
    }
    
    BLOCK {
        int index PK
        timestamp timestamp
        json transactions
        string previous_hash
        string hash
        string validator
        json signatures
    }
    
    NFT {
        string token_id PK
        int story_id FK
        string title
        string image_cid
        string metadata_cid
        json authors
        timestamp minted_at
        string tx_hash FK
    }
```

---

## Deployment Architecture

```mermaid
graph TB
    subgraph Internet
        Users[Users]
    end
    
    subgraph Vercel["Vercel (Frontend)"]
        NextJS[Next.js App]
        CDN[Edge CDN]
    end
    
    subgraph Railway["Railway (Backend)"]
        FastAPI[FastAPI Server]
        Ngrok[Ngrok Tunnel]
    end
    
    subgraph AWS["AWS EC2 (Blockchain)"]
        V1[Validator 1<br/>:8001]
        V2[Validator 2<br/>:8002]
        V3[Validator 3<br/>:8003]
        V4[Validator 4<br/>:8004]
        Observer[Observer Node<br/>:8080]
    end
    
    subgraph External["External Services"]
        Supabase[Supabase Auth]
        Groq[Groq LLM API]
        IPFS_Network[IPFS Network]
        PageKite[PageKite Tunnels]
    end
    
    Users -->|HTTPS| CDN
    CDN --> NextJS
    NextJS -->|API Proxy| Ngrok
    Ngrok --> FastAPI
    FastAPI --> Groq
    
    NextJS -.->|Future| V1
    V1 <-->|PBFT| V2
    V2 <-->|PBFT| V3
    V3 <-->|PBFT| V4
    V1 --> IPFS_Network
    V2 --> IPFS_Network
    V3 --> IPFS_Network
    V4 --> IPFS_Network
    
    V1 -.->|Public Access| PageKite
    V2 -.->|Public Access| PageKite
    V3 -.->|Public Access| PageKite
    V4 -.->|Public Access| PageKite
    
    Users -->|Read-Only| Observer
    Observer -->|Sync| V1
    
    NextJS --> Supabase
    
    style Vercel fill:#3b82f6,color:#fff
    style Railway fill:#10b981,color:#fff
    style AWS fill:#8b5cf6,color:#fff
    style External fill:#f59e0b,color:#fff
```

---

## State Transition: Story Line Lifecycle

```mermaid
stateDiagram-v2
    [*] --> Draft: User types
    Draft --> Suggested: Request AI
    Suggested --> Edited: User modifies
    Edited --> Submitted: User approves
    Submitted --> Verified: Backend validates
    Verified --> InCanonical: Story finalized
    InCanonical --> MintedNFT: Blockchain mint
    MintedNFT --> [*]
    
    Suggested --> Draft: User rejects
    Edited --> Draft: User discards
```

---

## State Transition: PBFT Block States

```mermaid
stateDiagram-v2
    [*] --> Pending: Transaction created
    Pending --> PrePrepare: Leader proposes
    PrePrepare --> Prepare: Validators receive
    Prepare --> Commit: 2f+1 PREPARE
    Commit --> Finalized: 2f+1 COMMIT
    Finalized --> [*]
    
    PrePrepare --> ViewChange: Timeout
    Prepare --> ViewChange: Leader failure
    ViewChange --> PrePrepare: New leader elected
```

---

## Network Topology

```mermaid
graph TB
    subgraph Public Internet
        Browser[User Browser]
        IPFS_Nodes[IPFS Network]
    end
    
    subgraph Frontend Network
        Vercel[Vercel Edge]
        NextJS[Next.js App]
    end
    
    subgraph Backend Network
        Railway[Railway PaaS]
        FastAPI[FastAPI + Milvus]
        Ngrok_Tunnel[Ngrok Tunnel]
    end
    
    subgraph Blockchain Network
        PageKite_1[pagekite.me/v1]
        PageKite_2[pagekite.me/v2]
        PageKite_3[pagekite.me/v3]
        PageKite_4[pagekite.me/v4]
        
        Validator_1[Validator 1<br/>AWS EC2]
        Validator_2[Validator 2<br/>AWS EC2]
        Validator_3[Validator 3<br/>AWS EC2]
        Validator_4[Validator 4<br/>AWS EC2]
    end
    
    Browser --> Vercel
    Vercel --> NextJS
    NextJS --> Ngrok_Tunnel
    Ngrok_Tunnel --> FastAPI
    
    Browser -.Future.-> PageKite_1
    PageKite_1 --> Validator_1
    PageKite_2 --> Validator_2
    PageKite_3 --> Validator_3
    PageKite_4 --> Validator_4
    
    Validator_1 <-->|PBFT TCP| Validator_2
    Validator_2 <-->|PBFT TCP| Validator_3
    Validator_3 <-->|PBFT TCP| Validator_4
    Validator_4 <-->|PBFT TCP| Validator_1
    
    Validator_1 --> IPFS_Nodes
    Validator_2 --> IPFS_Nodes
    Validator_3 --> IPFS_Nodes
    Validator_4 --> IPFS_Nodes
    
    style Frontend Network fill:#3b82f6,color:#fff
    style Backend Network fill:#10b981,color:#fff
    style Blockchain Network fill:#8b5cf6,color:#fff
```

---

## Class Diagram: Core Blockchain Types

```mermaid
classDiagram
    class Block {
        +int64 Index
        +time.Time Timestamp
        +[]Transaction Transactions
        +string PreviousHash
        +string Hash
        +string Validator
        +[]Signature Signatures
        +CalculateHash() string
    }
    
    class Transaction {
        +string TxID
        +TxType Type
        +string From
        +string To
        +interface Data
        +string Signature
        +time.Time Timestamp
        +Verify(publicKey) bool
    }
    
    class Wallet {
        +string Address
        +[]byte PublicKey
        +[]byte PrivateKey
        +string SupabaseUserID
        +time.Time CreatedAt
        +[]byte EncryptionSalt
        +Sign(data) string
        +Decrypt(password) []byte
    }
    
    class NFT {
        +string TokenID
        +int64 StoryID
        +string Title
        +string ImageCID
        +string MetadataCID
        +[]Author Authors
        +time.Time MintedAt
        +string TxHash
    }
    
    class Author {
        +string Address
        +string SupabaseID
        +string Username
        +int Contributions
        +float64 Percentage
    }
    
    class PBFT {
        +string NodeID
        +[]string Validators
        +int64 CurrentView
        +int64 CurrentSeq
        +int f
        +ProposeBlock(block) error
        +HandlePrePrepare(msg)
        +HandlePrepare(msg)
        +HandleCommit(msg)
    }
    
    Block "1" *-- "many" Transaction
    NFT "1" *-- "many" Author
    Wallet "1" --> "many" Transaction : signs
    Transaction --> NFT : creates
    PBFT --> Block : validates
```

---

## Technology Stack Layers

```mermaid
graph TB
    subgraph Presentation["Presentation Layer"]
        React[React Components]
        Tailwind[Tailwind CSS]
        shadcn[shadcn/ui]
    end
    
    subgraph Application["Application Layer"]
        NextJS[Next.js 14]
        TypeScript[TypeScript]
        APIProxy[API Proxy]
    end
    
    subgraph Business["Business Logic Layer"]
        FastAPI_Server[FastAPI Server]
        RAG[RAG Pipeline]
        Scheduler[Background Tasks]
    end
    
    subgraph AI["AI/ML Layer"]
        Embeddings[Sentence Transformers]
        LLM[Groq LLM API]
    end
    
    subgraph Consensus["Consensus Layer"]
        PBFT_Algorithm[PBFT Algorithm]
        P2P[P2P Network]
        Wallets[Wallet Manager]
    end
    
    subgraph Data["Data Layer"]
        Milvus_DB[(Milvus Vectors)]
        SQLite_DB[(SQLite)]
        BadgerDB_Store[(BadgerDB)]
        IPFS_Network[(IPFS)]
    end
    
    React --> NextJS
    Tailwind --> React
    shadcn --> React
    NextJS --> APIProxy
    TypeScript --> NextJS
    
    APIProxy --> FastAPI_Server
    FastAPI_Server --> RAG
    FastAPI_Server --> Scheduler
    
    RAG --> Embeddings
    RAG --> LLM
    
    PBFT_Algorithm --> P2P
    Wallets --> PBFT_Algorithm
    
    Embeddings --> Milvus_DB
    FastAPI_Server --> SQLite_DB
    PBFT_Algorithm --> BadgerDB_Store
    PBFT_Algorithm --> IPFS_Network
    
    style Presentation fill:#3b82f6,color:#fff
    style Application fill:#06b6d4,color:#fff
    style Business fill:#10b981,color:#fff
    style AI fill:#f59e0b,color:#fff
    style Consensus fill:#8b5cf6,color:#fff
    style Data fill:#ef4444,color:#fff
```

---

## Diagram Sources

All diagrams in this file are defined using [Mermaid.js](https://mermaid.js.org/) syntax for version control and easy editing.

For more architectural details, see:
- [`ARCHITECTURE.md`](./ARCHITECTURE.md) - Complete system design
- [`DATA_FLOW.md`](./DATA_FLOW.md) - Detailed DFD documentation
- [`BLOCKCHAIN.md`](./BLOCKCHAIN.md) - Blockchain implementation specs
