# 🔗 Kahani Blockchain Layer

> **Status**: Work in Progress  
> **Implementation**: Go 1.21+ with PBFT Consensus

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Core Components](#core-components)
- [Data Models](#data-models)
- [Consensus Mechanism](#consensus-mechanism)
- [Wallet System](#wallet-system)
- [NFT Implementation](#nft-implementation)
- [Project Structure](#project-structure)
- [Implementation Phases](#implementation-phases)
- [API Reference](#api-reference)
- [Deployment](#deployment)
- [Security](#security)

---

## Overview

The Kahani blockchain layer provides:
- **Immutable story storage** via content-addressed IPFS
- **NFT minting** with automatic co-authorship tracking
- **Decentralized consensus** using Practical Byzantine Fault Tolerance (PBFT)
- **Wallet auto-generation** per Supabase user account
- **Public accessibility** via PageKite tunneling

### Key Features

✅ **Automatic Wallet Creation**: Every Supabase user gets a blockchain wallet  
✅ **Co-Authorship NFTs**: Track contribution percentages on-chain  
✅ **IPFS Integration**: Content-addressed story and metadata storage  
✅ **PBFT Consensus**: Byzantine fault-tolerant block validation  
✅ **Ed25519 Signatures**: Cryptographic transaction signing  
✅ **BadgerDB Storage**: Embedded key-value database for blockchain state  

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   KAHANI BLOCKCHAIN LAYER                    │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌────────────┐   ┌────────────┐   ┌────────────┐          │
│  │ Validator  │   │ Validator  │   │ Validator  │          │
│  │  Node 1    │   │  Node 2    │   │  Node 3    │          │
│  └──────┬─────┘   └──────┬─────┘   └──────┬─────┘          │
│         │                 │                 │                │
│         │◄────PBFT────────┼────────────────▶│                │
│         │   Consensus     │                 │                │
│         │                 │                 │                │
│  ┌──────▼─────────────────▼─────────────────▼──────┐        │
│  │           P2P Network (TCP + PageKite)          │        │
│  └──────┬──────────────────────────────────────────┘        │
│         │                                                    │
│  ┌──────▼─────────────────────────────────────────┐         │
│  │              BadgerDB State Storage             │         │
│  │  - Blocks  - Transactions  - NFTs  - Wallets   │         │
│  └──────┬──────────────────────────────────────────┘         │
│         │                                                    │
│  ┌──────▼──────┐                 ┌────────────┐             │
│  │    IPFS     │                 │  Supabase  │             │
│  │  (Content)  │                 │   (Users)  │             │
│  └─────────────┘                 └────────────┘             │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Responsibility | Technology |
|-------|---------------|------------|
| **Consensus** | Block validation, leader election | PBFT algorithm |
| **Network** | Peer discovery, message broadcast | TCP + PageKite |
| **Storage** | Persistent blockchain state | BadgerDB |
| **Content** | Decentralized story/metadata hosting | IPFS |
| **Identity** | User-wallet mapping | Supabase integration |

---

## Core Components

### 1. Wallet Generator

**File**: `internal/wallet/generator.go`

```go
type WalletGenerator interface {
    GenerateWallet(supabaseUserID string) (*Wallet, error)
    GetWallet(supabaseUserID string) (*Wallet, error)
    StoreWallet(wallet *Wallet) error
}

type Wallet struct {
    Address        string    // Hex-encoded public key
    PublicKey      []byte    // Ed25519 public key (32 bytes)
    PrivateKey     []byte    // Ed25519 private key (64 bytes, encrypted)
    SupabaseUserID string    // Link to Supabase auth
    CreatedAt      time.Time
    EncryptionSalt []byte    // For private key encryption
}
```

**Key Features**:
- Ed25519 key pair generation using `crypto/ed25519`
- AES-256-GCM encryption for private keys at rest
- Deterministic address derivation: `address = hex(SHA256(publicKey)[:20])`
- Automatic creation on first user login

**Encryption Flow**:
```
User Password (Supabase) → PBKDF2(password, salt, 100000) → AES-256 Key
Private Key → AES-256-GCM Encrypt → Encrypted Blob → BadgerDB
```

### 2. Supabase Integration

**File**: `internal/supabase/client.go`

```go
type SupabaseClient interface {
    GetUserByID(userID string) (*User, error)
    ValidateJWT(token string) (*Claims, error)
    OnUserCreated(callback func(*User)) error
}

// Webhook handler for new user registration
func (s *Service) HandleUserCreated(w http.ResponseWriter, r *http.Request) {
    var user User
    json.NewDecoder(r.Body).Decode(&user)
    
    // Generate wallet automatically
    wallet, err := s.walletGen.GenerateWallet(user.ID)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    
    // Store in blockchain state
    s.blockchain.StoreWallet(wallet)
    w.WriteHeader(200)
}
```

**Integration Points**:
- Listen for Supabase webhook: `POST /api/blockchain/user-created`
- Validate JWT tokens from frontend requests
- Map `supabase_user_id` to blockchain `wallet_address`

### 3. PBFT Consensus Engine

**File**: `internal/consensus/pbft.go`

```go
type PBFT struct {
    nodeID        string
    validators    []string        // List of validator addresses
    currentView   int64
    currentSeq    int64
    f             int            // Max Byzantine nodes (n = 3f + 1)
    prepareLog    map[int64]map[string]*PrepareMsg
    commitLog     map[int64]map[string]*CommitMsg
}

// Three-phase PBFT protocol
func (p *PBFT) ProposeBlock(block *Block) error {
    // Phase 1: Pre-Prepare (leader broadcasts)
    prePrepare := &PrePrepareMsg{
        View:      p.currentView,
        Sequence:  p.currentSeq,
        Block:     block,
        Signature: p.sign(block),
    }
    p.broadcast(prePrepare)
    
    // Phase 2: Prepare (validators echo)
    // Wait for 2f + 1 prepare messages
    
    // Phase 3: Commit (validators commit)
    // Wait for 2f + 1 commit messages
    
    return nil
}
```

**Consensus Parameters**:
- **n = 4 validators** (tolerates 1 Byzantine node, f=1)
- **Block time**: ~3-5 seconds
- **View change timeout**: 10 seconds
- **Message signatures**: Ed25519

**PBFT Message Flow**:
```
Leader                Validator 1         Validator 2         Validator 3
  │                        │                   │                   │
  │─[PRE-PREPARE]─────────▶│                   │                   │
  │─[PRE-PREPARE]──────────┼──────────────────▶│                   │
  │─[PRE-PREPARE]──────────┼───────────────────┼──────────────────▶│
  │                        │                   │                   │
  │◀────[PREPARE]──────────│                   │                   │
  │◀────[PREPARE]──────────┼───────────────────│                   │
  │◀────[PREPARE]──────────┼───────────────────┼───────────────────│
  │                        │                   │                   │
  │      [2f+1 PREPARE messages received → Move to COMMIT]         │
  │                        │                   │                   │
  │◀────[COMMIT]───────────│                   │                   │
  │◀────[COMMIT]───────────┼───────────────────│                   │
  │◀────[COMMIT]───────────┼───────────────────┼───────────────────│
  │                        │                   │                   │
  │      [2f+1 COMMIT messages → FINALIZE BLOCK]                   │
  │                        │                   │                   │
  ▼                        ▼                   ▼                   ▼
[Block added to all validator chains]
```

### 4. P2P Network Layer

**File**: `internal/network/p2p.go`

```go
type P2PNode struct {
    nodeID       string
    listenAddr   string
    peers        map[string]*Peer
    messageBus   chan Message
    pagekite     *PageKiteConfig
}

type PageKiteConfig struct {
    FrontendAddr string  // e.g., "kahani-validator-1.pagekite.me"
    BackendPort  int     // Local port to expose
    Secret       string  // PageKite secret
}

func (n *P2PNode) Start() error {
    // Start TCP listener
    listener, _ := net.Listen("tcp", n.listenAddr)
    
    // Start PageKite tunnel
    go n.startPageKite()
    
    // Accept incoming connections
    for {
        conn, _ := listener.Accept()
        go n.handlePeer(conn)
    }
}
```

**Network Topology**:
- **Validators**: Full nodes running PBFT (e.g., 4 nodes)
- **Observers**: Read-only nodes (sync blockchain, serve queries)
- **PageKite Tunnels**: Public endpoints for validators
  - `validator-1.pagekite.me:443` → `localhost:8001`
  - `validator-2.pagekite.me:443` → `localhost:8002`

**Message Types**:
```go
type MessageType int

const (
    MsgTransaction    MessageType = iota
    MsgPrePrepare
    MsgPrepare
    MsgCommit
    MsgViewChange
    MsgNewView
    MsgBlockRequest
    MsgBlockResponse
)
```

### 5. BadgerDB Storage

**File**: `internal/storage/badger.go`

```go
type BadgerStore struct {
    db *badger.DB
}

// Key prefixes for different data types
const (
    PrefixBlock       = "blk:"  // blk:<block_index> → Block
    PrefixTx          = "tx:"   // tx:<tx_hash> → Transaction
    PrefixNFT         = "nft:"  // nft:<token_id> → NFT
    PrefixWallet      = "wal:"  // wal:<address> → Wallet
    PrefixState       = "st:"   // st:<key> → State variable
    PrefixUserWallet  = "uw:"   // uw:<supabase_id> → Wallet address
)

func (s *BadgerStore) StoreBlock(block *Block) error {
    key := fmt.Sprintf("%s%d", PrefixBlock, block.Index)
    data, _ := json.Marshal(block)
    
    return s.db.Update(func(txn *badger.Txn) error {
        return txn.Set([]byte(key), data)
    })
}

func (s *BadgerStore) GetNFT(tokenID string) (*NFT, error) {
    key := fmt.Sprintf("%s%s", PrefixNFT, tokenID)
    var nft NFT
    
    err := s.db.View(func(txn *badger.Txn) error {
        item, err := txn.Get([]byte(key))
        if err != nil {
            return err
        }
        return item.Value(func(val []byte) error {
            return json.Unmarshal(val, &nft)
        })
    })
    
    return &nft, err
}
```

**Storage Schema**:
```
BadgerDB Key-Value Store
├─ blk:0 → Genesis Block
├─ blk:1 → Block #1
├─ tx:abc123 → Transaction{type: "mint_nft", ...}
├─ nft:token_1 → NFT{story_id: 42, image_cid: "Qm...", authors: [...]}
├─ wal:0x1a2b3c → Wallet{pubkey, encrypted_privkey, ...}
├─ uw:user_uuid_123 → "0x1a2b3c"  (Supabase ID → Wallet mapping)
└─ st:latest_block → 157
```

### 6. IPFS Integration

**File**: `internal/ipfs/client.go`

```go
type IPFSClient interface {
    Add(data []byte) (string, error)     // Returns CID
    Get(cid string) ([]byte, error)
    Pin(cid string) error
}

func (i *IPFSClient) UploadStoryMetadata(nft *NFT) (string, error) {
    metadata := map[string]interface{}{
        "name":        nft.Title,
        "description": nft.Description,
        "image":       fmt.Sprintf("ipfs://%s", nft.ImageCID),
        "attributes": []map[string]interface{}{
            {"trait_type": "Word Count", "value": nft.WordCount},
            {"trait_type": "Line Count", "value": nft.LineCount},
            {"trait_type": "Main Author", "value": nft.Authors[0].Address},
        },
        "authors": nft.Authors,  // Co-authorship tracking
    }
    
    jsonData, _ := json.Marshal(metadata)
    cid, _ := i.Add(jsonData)
    
    // Pin to ensure availability
    i.Pin(cid)
    
    return cid, nil
}
```

**IPFS Data Structure**:
```
IPFS Network
├─ QmStory123... → Full canonical story text
├─ QmImage456... → NFT cover image (PNG)
└─ QmMeta789...  → NFT metadata JSON
                   {
                     "name": "The Dragon's Tale",
                     "image": "ipfs://QmImage456...",
                     "authors": [
                       {"address": "0xABC", "percentage": 60},
                       {"address": "0xDEF", "percentage": 40}
                     ]
                   }
```

---

## Data Models

### Block Structure

```go
type Block struct {
    Index         int64         `json:"index"`
    Timestamp     time.Time     `json:"timestamp"`
    Transactions  []Transaction `json:"transactions"`
    PreviousHash  string        `json:"previous_hash"`
    Hash          string        `json:"hash"`
    Validator     string        `json:"validator"`  // Who proposed this block
    Signatures    []Signature   `json:"signatures"` // PBFT validator signatures
}

func (b *Block) CalculateHash() string {
    record := fmt.Sprintf("%d%s%s%s", 
        b.Index, 
        b.Timestamp, 
        b.PreviousHash,
        hashTransactions(b.Transactions),
    )
    h := sha256.New()
    h.Write([]byte(record))
    return hex.EncodeToString(h.Sum(nil))
}
```

### Transaction Types

```go
type Transaction struct {
    TxID        string        `json:"tx_id"`
    Type        TxType        `json:"type"`
    From        string        `json:"from"`         // Sender wallet address
    To          string        `json:"to"`           // Recipient (optional)
    Data        interface{}   `json:"data"`         // Type-specific payload
    Signature   string        `json:"signature"`    // Ed25519 signature
    Timestamp   time.Time     `json:"timestamp"`
}

type TxType string

const (
    TxMintNFT        TxType = "mint_nft"
    TxTransferNFT    TxType = "transfer_nft"
    TxContribute     TxType = "contribute"  // Record story contribution
    TxCreateWallet   TxType = "create_wallet"
)
```

**Transaction Payloads**:

```go
// For TxMintNFT
type MintNFTData struct {
    StoryID      int64    `json:"story_id"`      // Reference to canonical story
    Title        string   `json:"title"`
    ImageCID     string   `json:"image_cid"`     // IPFS content ID
    MetadataCID  string   `json:"metadata_cid"`  // IPFS metadata
    Authors      []Author `json:"authors"`       // Co-authorship details
}

// For TxContribute
type ContributeData struct {
    StoryID      int64  `json:"story_id"`
    LineID       int64  `json:"line_id"`
    LineText     string `json:"line_text"`  // For immutability
    LineHash     string `json:"line_hash"`  // SHA256 of line text
}
```

### NFT Model

```go
type NFT struct {
    TokenID      string    `json:"token_id"`      // Unique identifier
    StoryID      int64     `json:"story_id"`      // Links to AI backend DB
    Title        string    `json:"title"`
    Description  string    `json:"description"`
    ImageCID     string    `json:"image_cid"`     // IPFS image
    MetadataCID  string    `json:"metadata_cid"`  // IPFS metadata
    Authors      []Author  `json:"authors"`
    MintedAt     time.Time `json:"minted_at"`
    MintedBy     string    `json:"minted_by"`     // Wallet address who triggered mint
    TxHash       string    `json:"tx_hash"`       // Minting transaction hash
}

type Author struct {
    Address       string  `json:"address"`         // Wallet address
    SupabaseID    string  `json:"supabase_id"`    // Link to user
    Username      string  `json:"username"`        // Display name
    Contributions int     `json:"contributions"`   // Number of lines contributed
    Percentage    float64 `json:"percentage"`      // Ownership percentage
}
```

**Co-Authorship Calculation**:
```go
func CalculateAuthorship(storyLines []StoryLine) []Author {
    contributionMap := make(map[string]int)
    
    // Count contributions per user
    for _, line := range storyLines {
        contributionMap[line.UserID]++
    }
    
    // Calculate percentages
    total := len(storyLines)
    authors := []Author{}
    for userID, count := range contributionMap {
        wallet := getWalletBySupabaseID(userID)
        authors = append(authors, Author{
            Address:       wallet.Address,
            SupabaseID:    userID,
            Contributions: count,
            Percentage:    float64(count) / float64(total) * 100,
        })
    }
    
    // Sort by contribution (descending)
    sort.Slice(authors, func(i, j int) bool {
        return authors[i].Contributions > authors[j].Contributions
    })
    
    return authors
}
```

---

## Consensus Mechanism

### PBFT Overview

**Practical Byzantine Fault Tolerance** provides consensus in environments with up to `f` Byzantine (malicious) nodes, given `n ≥ 3f + 1` total validators.

**Kahani Configuration**:
- **Validators**: 4 nodes (n=4)
- **Byzantine tolerance**: 1 node (f=1)
- **Quorum size**: 2f + 1 = 3 signatures required

### PBFT Phases

**1. Pre-Prepare Phase**
```go
type PrePrepareMsg struct {
    View      int64  `json:"view"`       // Current view number
    Sequence  int64  `json:"sequence"`   // Block sequence
    Block     *Block `json:"block"`      // Proposed block
    Signature string `json:"signature"`  // Leader's signature
}

// Leader broadcasts to all validators
func (p *PBFT) sendPrePrepare(block *Block) {
    msg := &PrePrepareMsg{
        View:      p.currentView,
        Sequence:  p.currentSeq,
        Block:     block,
        Signature: p.sign(block),
    }
    p.broadcast(msg)
}
```

**2. Prepare Phase**
```go
type PrepareMsg struct {
    View      int64  `json:"view"`
    Sequence  int64  `json:"sequence"`
    BlockHash string `json:"block_hash"`
    NodeID    string `json:"node_id"`
    Signature string `json:"signature"`
}

// Each validator validates and echoes
func (p *PBFT) handlePrePrepare(msg *PrePrepareMsg) {
    if !p.validateBlock(msg.Block) {
        return
    }
    
    prepare := &PrepareMsg{
        View:      msg.View,
        Sequence:  msg.Sequence,
        BlockHash: msg.Block.Hash,
        NodeID:    p.nodeID,
        Signature: p.sign(msg.Block.Hash),
    }
    p.broadcast(prepare)
    p.prepareLog[msg.Sequence][p.nodeID] = prepare
}

// Wait for quorum
func (p *PBFT) isPrepared(seq int64) bool {
    return len(p.prepareLog[seq]) >= 2*p.f + 1
}
```

**3. Commit Phase**
```go
type CommitMsg struct {
    View      int64  `json:"view"`
    Sequence  int64  `json:"sequence"`
    BlockHash string `json:"block_hash"`
    NodeID    string `json:"node_id"`
    Signature string `json:"signature"`
}

func (p *PBFT) handlePrepare(msg *PrepareMsg) {
    p.prepareLog[msg.Sequence][msg.NodeID] = msg
    
    if p.isPrepared(msg.Sequence) {
        commit := &CommitMsg{
            View:      msg.View,
            Sequence:  msg.Sequence,
            BlockHash: msg.BlockHash,
            NodeID:    p.nodeID,
            Signature: p.sign(msg.BlockHash),
        }
        p.broadcast(commit)
        p.commitLog[msg.Sequence][p.nodeID] = commit
    }
}

func (p *PBFT) isCommitted(seq int64) bool {
    return len(p.commitLog[seq]) >= 2*p.f + 1
}

// Finalize block
func (p *PBFT) finalizeBlock(block *Block) {
    block.Signatures = p.collectSignatures(block.Index)
    p.blockchain.AddBlock(block)
    p.currentSeq++
}
```

### View Change (Leader Failure)

```go
func (p *PBFT) startViewChange() {
    p.currentView++
    
    viewChange := &ViewChangeMsg{
        NewView:   p.currentView,
        LastSeq:   p.currentSeq,
        NodeID:    p.nodeID,
        Signature: p.sign(fmt.Sprintf("%d:%d", p.currentView, p.currentSeq)),
    }
    p.broadcast(viewChange)
    
    // Wait for 2f + 1 VIEW-CHANGE messages
    // New leader selected: validators[newView % len(validators)]
}
```

---

## Wallet System

### Wallet Generation Flow

```
┌─────────────┐
│ User signs  │
│ up via      │
│ Supabase    │
└──────┬──────┘
       │
       │ [Webhook: POST /api/blockchain/user-created]
       ▼
┌──────────────┐
│ Blockchain   │
│ Service      │
└──────┬───────┘
       │
       │ [Generate Ed25519 key pair]
       ▼
┌──────────────┐
│ Private Key  │
│ Encryption   │
│ (AES-256)    │
└──────┬───────┘
       │
       │ [Store in BadgerDB]
       ▼
┌──────────────┐
│ Wallet       │
│ Address:     │
│ 0x1a2b3c...  │
└──────────────┘
```

### Wallet Storage Schema

```go
// In BadgerDB
Key:   wal:0x1a2b3c4d5e6f7890...
Value: {
  "address": "0x1a2b3c4d5e6f7890...",
  "public_key": "base64_encoded_pubkey",
  "private_key_encrypted": "AES256_encrypted_blob",
  "supabase_user_id": "user_uuid_from_supabase",
  "created_at": "2024-01-15T10:30:00Z",
  "encryption_salt": "random_salt_bytes"
}

// Reverse lookup
Key:   uw:user_uuid_from_supabase
Value: "0x1a2b3c4d5e6f7890..."
```

### Transaction Signing

```go
func (w *Wallet) SignTransaction(tx *Transaction) (string, error) {
    // Decrypt private key
    privKey, err := w.DecryptPrivateKey(userPassword)
    if err != nil {
        return "", err
    }
    
    // Create canonical transaction representation
    txData := fmt.Sprintf("%s:%s:%s:%v:%s",
        tx.TxID,
        tx.Type,
        tx.From,
        tx.Data,
        tx.Timestamp.Unix(),
    )
    
    // Sign with Ed25519
    signature := ed25519.Sign(privKey, []byte(txData))
    
    return hex.EncodeToString(signature), nil
}

func VerifyTransactionSignature(tx *Transaction, publicKey []byte) bool {
    txData := fmt.Sprintf("%s:%s:%s:%v:%s",
        tx.TxID,
        tx.Type,
        tx.From,
        tx.Data,
        tx.Timestamp.Unix(),
    )
    
    sigBytes, _ := hex.DecodeString(tx.Signature)
    return ed25519.Verify(publicKey, []byte(txData), sigBytes)
}
```

---

## NFT Implementation

### NFT Minting Flow

```
┌──────────┐
│  User    │
│ Frontend │
└────┬─────┘
     │
     │ [1] POST /api/blockchain/mint-nft
     │     {canonical_story_id: 42}
     ▼
┌────────────┐
│ Blockchain │
│  Service   │
└────┬───────┘
     │
     │ [2] Fetch canonical story from AI backend
     ▼
┌────────────┐
│  FastAPI   │
│  Backend   │
└────┬───────┘
     │
     │ [3] Story + line contributions
     ▼
┌────────────┐
│ Calculate  │
│ Co-Author  │
│ Percentages│
└────┬───────┘
     │
     │ [4] Generate NFT cover image
     ▼
┌────────────┐
│   IPFS     │◀── Upload image
│  Network   │
└────┬───────┘
     │
     │ [5] Image CID: QmABC123...
     ▼
┌────────────┐
│  Create    │
│  Metadata  │
│   JSON     │
└────┬───────┘
     │
     │ [6] Upload metadata to IPFS
     ▼
┌────────────┐
│   IPFS     │◀── Upload JSON
│  Network   │
└────┬───────┘
     │
     │ [7] Metadata CID: QmDEF456...
     ▼
┌────────────┐
│   Create   │
│ mint_nft   │
│Transaction │
└────┬───────┘
     │
     │ [8] Sign with user's wallet
     ▼
┌────────────┐
│  PBFT      │
│ Consensus  │
└────┬───────┘
     │
     │ [9] Block finalized
     ▼
┌────────────┐
│  BadgerDB  │
│  nft:42 →  │
│  {NFT}     │
└────┬───────┘
     │
     │ [10] Return token_id
     ▼
┌────────────┐
│  Frontend  │
│  (Success) │
└────────────┘
```

### NFT Metadata Standard

Following OpenSea metadata standards with Kahani extensions:

```json
{
  "name": "The Dragon's Tale",
  "description": "A collaborative fantasy story about a dragon who learns to code",
  "image": "ipfs://QmABC123.../cover.png",
  "external_url": "https://kahani.app/story/42",
  "attributes": [
    {
      "trait_type": "Word Count",
      "value": 3245
    },
    {
      "trait_type": "Line Count",
      "value": 87
    },
    {
      "trait_type": "Main Author",
      "value": "0x1a2b3c4d5e6f"
    },
    {
      "trait_type": "Co-Authors",
      "value": 4
    },
    {
      "trait_type": "Genre",
      "value": "Fantasy"
    },
    {
      "trait_type": "Minted Date",
      "display_type": "date",
      "value": 1642262400
    }
  ],
  "kahani_metadata": {
    "story_id": 42,
    "canonical_story_id": 15,
    "blockchain_tx": "0xabc123...",
    "authors": [
      {
        "address": "0x1a2b3c4d5e6f",
        "supabase_id": "user_uuid_1",
        "username": "DragonWriter",
        "contributions": 52,
        "percentage": 59.77
      },
      {
        "address": "0x7g8h9i0j1k2l",
        "supabase_id": "user_uuid_2",
        "username": "CodeScribe",
        "contributions": 35,
        "percentage": 40.23
      }
    ],
    "full_text_cid": "QmStory789..."
  }
}
```

---

## Project Structure

```
storytelling-blockchain/
├── cmd/
│   ├── validator/            # Validator node binary
│   │   └── main.go
│   ├── observer/             # Observer node binary
│   │   └── main.go
│   └── cli/                  # CLI tools (wallet gen, query)
│       └── main.go
├── internal/
│   ├── blockchain/
│   │   ├── chain.go          # Blockchain data structure
│   │   ├── block.go          # Block model
│   │   └── transaction.go    # Transaction types
│   ├── consensus/
│   │   ├── pbft.go           # PBFT implementation
│   │   ├── messages.go       # PBFT message types
│   │   └── view_change.go    # Leader election
│   ├── network/
│   │   ├── p2p.go            # P2P networking
│   │   ├── pagekite.go       # PageKite integration
│   │   └── discovery.go      # Peer discovery
│   ├── storage/
│   │   ├── badger.go         # BadgerDB interface
│   │   └── schema.go         # Key naming conventions
│   ├── wallet/
│   │   ├── generator.go      # Ed25519 wallet generation
│   │   ├── encryption.go     # Private key encryption
│   │   └── signing.go        # Transaction signing
│   ├── ipfs/
│   │   ├── client.go         # IPFS HTTP client
│   │   └── pinning.go        # Content pinning
│   ├── nft/
│   │   ├── minter.go         # NFT minting logic
│   │   ├── metadata.go       # Metadata generation
│   │   └── authorship.go     # Co-authorship calculation
│   ├── supabase/
│   │   ├── client.go         # Supabase API client
│   │   └── webhooks.go       # User creation handler
│   └── api/
│       ├── server.go         # HTTP API server
│       ├── handlers.go       # Request handlers
│       └── middleware.go     # Auth, CORS, logging
├── pkg/
│   ├── types/                # Shared data types
│   │   ├── block.go
│   │   ├── transaction.go
│   │   ├── nft.go
│   │   └── wallet.go
│   └── utils/                # Helper functions
│       ├── hash.go
│       └── encoding.go
├── config/
│   ├── validator-1.yaml      # Validator node config
│   ├── validator-2.yaml
│   ├── validator-3.yaml
│   ├── validator-4.yaml
│   └── observer.yaml         # Observer node config
├── scripts/
│   ├── start-validator.sh    # Launch validator node
│   ├── start-observer.sh     # Launch observer node
│   └── generate-genesis.go   # Create genesis block
├── docker/
│   ├── Dockerfile.validator
│   └── docker-compose.yml    # Multi-node deployment
├── docs/
│   ├── API.md                # HTTP API reference
│   ├── PBFT.md               # Consensus details
│   └── DEPLOYMENT.md         # Production setup
├── tests/
│   ├── consensus_test.go
│   ├── wallet_test.go
│   └── integration_test.go
├── go.mod
├── go.sum
└── README.md
```

---

## Implementation Phases

### Phase 1: Core Infrastructure (Weeks 1-2)

**Files to Implement**:
1. `pkg/types/*.go` - Define all data models
2. `internal/storage/badger.go` - BadgerDB wrapper
3. `internal/wallet/generator.go` - Ed25519 wallet creation
4. `internal/wallet/encryption.go` - AES-256 encryption
5. `internal/blockchain/chain.go` - Basic blockchain structure
6. `internal/blockchain/block.go` - Block model & hashing

**Validation**:
- [ ] Generate wallet and verify signature
- [ ] Store/retrieve block from BadgerDB
- [ ] Calculate block hash correctly

### Phase 2: Consensus Engine (Weeks 3-4)

**Files to Implement**:
1. `internal/consensus/pbft.go` - PBFT core logic
2. `internal/consensus/messages.go` - Message types
3. `internal/consensus/view_change.go` - Leader election
4. `internal/network/p2p.go` - TCP networking

**Validation**:
- [ ] 4 validator nodes reach consensus on block
- [ ] System tolerates 1 Byzantine node
- [ ] View change works when leader fails

### Phase 3: IPFS & NFT Minting (Weeks 5-6)

**Files to Implement**:
1. `internal/ipfs/client.go` - IPFS integration
2. `internal/nft/minter.go` - Minting logic
3. `internal/nft/metadata.go` - Metadata generation
4. `internal/nft/authorship.go` - Co-authorship calculation
5. `internal/api/handlers.go` - HTTP endpoints

**Validation**:
- [ ] Upload story to IPFS successfully
- [ ] Mint NFT with correct co-authorship
- [ ] Retrieve NFT metadata via API

### Phase 4: Supabase Integration (Week 7)

**Files to Implement**:
1. `internal/supabase/client.go` - Supabase SDK
2. `internal/supabase/webhooks.go` - User creation webhook
3. `internal/api/middleware.go` - JWT validation

**Validation**:
- [ ] New Supabase user auto-generates wallet
- [ ] JWT tokens validated correctly
- [ ] User-wallet mapping works

### Phase 5: PageKite & Deployment (Week 8)

**Files to Implement**:
1. `internal/network/pagekite.go` - PageKite tunneling
2. `docker/Dockerfile.validator` - Containerization
3. `docker/docker-compose.yml` - Multi-node setup
4. `config/*.yaml` - Node configurations

**Validation**:
- [ ] Validators accessible via public URLs
- [ ] Docker Compose starts 4-node network
- [ ] End-to-end NFT minting works

### Phase 6: Testing & Optimization (Weeks 9-10)

**Files to Implement**:
1. `tests/consensus_test.go` - PBFT unit tests
2. `tests/integration_test.go` - Full flow tests
3. `docs/API.md` - API documentation

**Validation**:
- [ ] 90%+ code coverage
- [ ] Load test: 100 TPS for 5 minutes
- [ ] Security audit passes

---

## API Reference

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

#### 1. Wallet Management

**Create Wallet** (Internal, called by Supabase webhook)
```http
POST /wallet/create
Content-Type: application/json

{
  "supabase_user_id": "user_uuid_123",
  "username": "DragonWriter"
}

Response 201:
{
  "address": "0x1a2b3c4d5e6f7890...",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Get Wallet**
```http
GET /wallet/{address}

Response 200:
{
  "address": "0x1a2b3c4d5e6f7890...",
  "supabase_user_id": "user_uuid_123",
  "created_at": "2024-01-15T10:30:00Z"
}
```

#### 2. NFT Operations

**Mint NFT**
```http
POST /nft/mint
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "canonical_story_id": 42,
  "title": "The Dragon's Tale",
  "description": "A collaborative fantasy story"
}

Response 201:
{
  "token_id": "nft_1",
  "tx_hash": "0xabc123...",
  "image_cid": "QmABC123...",
  "metadata_cid": "QmDEF456...",
  "authors": [
    {
      "address": "0x1a2b3c",
      "username": "DragonWriter",
      "percentage": 59.77
    }
  ],
  "minted_at": "2024-01-15T10:35:00Z"
}
```

**Get NFT**
```http
GET /nft/{token_id}

Response 200:
{
  "token_id": "nft_1",
  "story_id": 42,
  "title": "The Dragon's Tale",
  "image_cid": "QmABC123...",
  "metadata_cid": "QmDEF456...",
  "authors": [...],
  "minted_at": "2024-01-15T10:35:00Z"
}
```

**List User NFTs**
```http
GET /nft/user/{address}

Response 200:
{
  "nfts": [
    {
      "token_id": "nft_1",
      "title": "The Dragon's Tale",
      "minted_at": "2024-01-15T10:35:00Z"
    }
  ],
  "total": 1
}
```

#### 3. Blockchain Queries

**Get Block**
```http
GET /block/{index}

Response 200:
{
  "index": 157,
  "timestamp": "2024-01-15T10:40:00Z",
  "transactions": [...],
  "previous_hash": "0xdef456...",
  "hash": "0xghi789...",
  "validator": "validator-1",
  "signatures": [...]
}
```

**Get Transaction**
```http
GET /transaction/{tx_hash}

Response 200:
{
  "tx_id": "tx_abc123",
  "type": "mint_nft",
  "from": "0x1a2b3c",
  "data": {...},
  "signature": "0x...",
  "block_index": 157,
  "timestamp": "2024-01-15T10:40:00Z"
}
```

**Get Blockchain Stats**
```http
GET /stats

Response 200:
{
  "latest_block": 157,
  "total_transactions": 2345,
  "total_nfts": 42,
  "total_wallets": 128,
  "validators": 4,
  "network_status": "healthy"
}
```

---

## Deployment

### Local Development (Single Node)

```bash
# Start IPFS daemon
ipfs daemon &

# Start validator node
go run cmd/validator/main.go \
  --config config/validator-1.yaml \
  --badger-path ./data/validator-1 \
  --listen-addr :8001

# In another terminal, test API
curl http://localhost:8001/api/v1/stats
```

### Production (4 Validator Nodes)

**Using Docker Compose**:
```bash
cd docker/
docker-compose up -d

# Verify all validators are running
docker-compose ps

# Check logs
docker-compose logs -f validator-1
```

**Manual Deployment** (e.g., AWS EC2):

**Validator 1**:
```bash
# On ec2-validator-1.amazonaws.com
./validator \
  --config /etc/kahani/validator-1.yaml \
  --badger-path /var/lib/kahani/blockchain \
  --listen-addr :8001 \
  --pagekite-frontend kahani-v1.pagekite.me:443 \
  --pagekite-secret <SECRET_1>
```

**Validator 2, 3, 4**: Same commands with different config files

**Observer Node**:
```bash
./observer \
  --config /etc/kahani/observer.yaml \
  --badger-path /var/lib/kahani/observer \
  --validator-urls "https://kahani-v1.pagekite.me,https://kahani-v2.pagekite.me,..."
```

### Configuration File Example

**`config/validator-1.yaml`**:
```yaml
node:
  id: validator-1
  type: validator
  listen_addr: ":8001"
  
consensus:
  pbft_enabled: true
  validators:
    - validator-1
    - validator-2
    - validator-3
    - validator-4
  max_byzantine: 1
  block_time_seconds: 5
  view_change_timeout_seconds: 10

network:
  bootstrap_peers:
    - "https://kahani-v2.pagekite.me"
    - "https://kahani-v3.pagekite.me"
    - "https://kahani-v4.pagekite.me"
  pagekite:
    frontend: "kahani-v1.pagekite.me"
    backend_port: 8001
    secret: "${PAGEKITE_SECRET}"

storage:
  badger_path: "/var/lib/kahani/blockchain"
  ipfs_api: "http://localhost:5001"

supabase:
  url: "${SUPABASE_URL}"
  service_role_key: "${SUPABASE_SERVICE_ROLE_KEY}"
  jwt_secret: "${SUPABASE_JWT_SECRET}"

api:
  enabled: true
  port: 8080
  cors_origins:
    - "https://kahani.app"
    - "http://localhost:3000"
```

---

## Security

### Threat Model

| Threat | Mitigation |
|--------|-----------|
| **Private key theft** | AES-256-GCM encryption at rest, password-derived keys |
| **Transaction forgery** | Ed25519 signature verification on all transactions |
| **Byzantine validators** | PBFT tolerates up to f=1 malicious node |
| **MITM attacks** | TLS/HTTPS on all PageKite tunnels |
| **Replay attacks** | Timestamp + nonce in transaction data |
| **DoS attacks** | Rate limiting on API endpoints |

### Best Practices

1. **Key Management**:
   - Never store unencrypted private keys
   - Use hardware security modules (HSM) for validator keys in production
   - Rotate PageKite secrets quarterly

2. **Network Security**:
   - Enable TLS 1.3 for all P2P communication
   - Whitelist validator IP addresses
   - Use VPN for validator-to-validator communication

3. **Smart Contract Audits** (Future):
   - If implementing custom smart contracts, undergo security audit
   - Use OpenZeppelin libraries where applicable

4. **Monitoring**:
   - Alert on consensus view changes (potential leader failure)
   - Monitor transaction validation failures
   - Track BadgerDB disk usage

---

## Roadmap

### Q1 2025
- [x] Blockchain architecture design
- [ ] Core infrastructure implementation
- [ ] PBFT consensus engine
- [ ] Wallet auto-generation

### Q2 2025
- [ ] IPFS integration
- [ ] NFT minting with co-authorship
- [ ] Supabase webhook integration
- [ ] PageKite deployment

### Q3 2025
- [ ] Observer node implementation
- [ ] Public block explorer UI
- [ ] NFT marketplace integration
- [ ] Performance optimization (target: 100 TPS)

### Q4 2025
- [ ] Cross-chain bridge (Ethereum/Polygon)
- [ ] Royalty distribution smart contracts
- [ ] Mobile app blockchain sync

---

## References

- [PBFT Paper](http://pmg.csail.mit.edu/papers/osdi99.pdf) - Original Practical Byzantine Fault Tolerance paper
- [BadgerDB Documentation](https://dgraph.io/docs/badger/)
- [IPFS Documentation](https://docs.ipfs.tech/)
- [Ed25519 Specification](https://ed25519.cr.yp.to/)
- [OpenSea Metadata Standards](https://docs.opensea.io/docs/metadata-standards)
- [PageKite Documentation](https://pagekite.net/wiki/)

---

**For more details**:
- Main architecture: [`docs/ARCHITECTURE.md`](./ARCHITECTURE.md)
- Data flow diagrams: [`docs/DATA_FLOW.md`](./DATA_FLOW.md)
- Frontend documentation: [`README.md`](../README.md)
- AI Backend documentation: [`Kahani_Ai_backend/README.md`](../Kahani_Ai_backend/README.md)
