# 📊 Kahani Data Flow Diagrams (DFD)

## DFD Level 0: Context Diagram

```
                        ┌─────────────────┐
                        │                 │
          ┌────────────▶│  Storyteller    │◀──────────┐
          │             │  (User)         │           │
          │             └─────────────────┘           │
          │                                            │
    [Stories]                                    [Auth Token]
    [PDFs]                                       [Contributions]
          │                                            │
          │             ┌─────────────────┐           │
          │             │                 │           │
          └─────────────│  KAHANI SYSTEM  │───────────┘
                        │                 │
                        └────────┬────────┘
                                 │
                    ┌────────────┼────────────┐
                    │            │            │
              [Context]      [Lore]      [Blocks]
                    │            │            │
                    ▼            ▼            ▼
            ┌──────────┐  ┌──────────┐  ┌──────────┐
            │ Vector   │  │   LLM    │  │Blockchain│
            │ Database │  │  (Groq)  │  │ Network  │
            └──────────┘  └──────────┘  └──────────┘
```

**External Entities:**
- **Storyteller**: End user creating collaborative narratives
- **Vector Database**: Milvus for semantic search
- **LLM Service**: Groq Cloud for story generation
- **Blockchain Network**: Decentralized ledger (future)

---

## DFD Level 1: Main System Processes

```
┌────────────────────────────────────────────────────────────────────┐
│                        KAHANI SYSTEM                                │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌─────────┐                                                        │
│  │  User   │                                                        │
│  └────┬────┘                                                        │
│       │                                                             │
│       │ [1] Story Prompt                                            │
│       ▼                                                             │
│  ┌────────────────┐                                                 │
│  │   P1: Story    │──[Context Request]──▶┌──────────┐              │
│  │   Suggestion   │◀─[Relevant Stories]──│ D1:      │              │
│  │   Generator    │                      │ Milvus   │              │
│  └────────┬───────┘                      │ Vectors  │              │
│           │                              └──────────┘              │
│           │ [2] AI Suggestion                                       │
│           ▼                                                         │
│  ┌────────────────┐                                                 │
│  │   P2: Story    │                                                 │
│  │   Editor &     │                                                 │
│  │   Verifier     │                                                 │
│  └────────┬───────┘                                                 │
│           │                                                         │
│           │ [3] Verified Line                                       │
│           ▼                                                         │
│  ┌────────────────┐──[Store]──▶┌──────────┐                        │
│  │   P3: Story    │            │ D2:      │                        │
│  │   Persistence  │◀─[Retrieve]│ SQLite   │                        │
│  │   Manager      │            │ Database │                        │
│  └────────┬───────┘            └──────────┘                        │
│           │                                                         │
│           │ [4] Embedding                                           │
│           ▼                                                         │
│  ┌────────────────┐──[Vector]──▶┌──────────┐                       │
│  │   P4: Vector   │             │ D1:      │                       │
│  │   Embedding    │             │ Milvus   │                       │
│  │   Storage      │             │ Vectors  │                       │
│  └────────┬───────┘             └──────────┘                       │
│           │                                                         │
│           │ [Background Task]                                       │
│           ▼                                                         │
│  ┌────────────────┐──[Extract]──▶┌──────────┐                      │
│  │   P5: Lore     │              │ D3:      │                      │
│  │   Extractor    │◀─[Query]─────│ Lore DB  │                      │
│  └────────┬───────┘              └──────────┘                      │
│           │                                                         │
│           │ [5] Request Canonical                                   │
│           ▼                                                         │
│  ┌────────────────┐──[Polish]───▶┌──────────┐                      │
│  │   P6: Story    │              │ D4:      │                      │
│  │   Canonicalizer│              │ Canonical│                      │
│  └────────┬───────┘              │ Stories  │                      │
│           │                      └──────────┘                      │
│           │ [6] NFT Mint (future)                                   │
│           ▼                                                         │
│  ┌────────────────┐──[Transaction]─▶┌────────────┐                 │
│  │   P7: NFT      │                 │ D5:        │                 │
│  │   Minter       │◀──[Confirm]─────│ Blockchain │                 │
│  │  (Work in      │                 │ Ledger     │                 │
│  │   Progress)    │                 └────────────┘                 │
│  └────────────────┘                                                 │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

**Data Stores:**
- **D1**: Milvus Vectors (semantic search index)
- **D2**: SQLite Database (story lines, metadata)
- **D3**: Lore DB (extracted entities)
- **D4**: Canonical Stories (polished narratives)
- **D5**: Blockchain Ledger (immutable record - future)

---

## DFD Level 2: Story Suggestion Process (P1)

```
┌─────────────────────────────────────────────────────────────┐
│              P1: STORY SUGGESTION GENERATOR                  │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  [Input: User Prompt]                                        │
│         │                                                    │
│         ▼                                                    │
│  ┌────────────────┐                                          │
│  │ P1.1: Prompt   │                                          │
│  │ Sanitization   │                                          │
│  └────────┬───────┘                                          │
│           │ [Cleaned Prompt]                                 │
│           ▼                                                  │
│  ┌────────────────┐──[Query Vector]──▶┌──────────┐          │
│  │ P1.2: Generate │                   │ D1:      │          │
│  │ Query Embedding│                   │ Sentence │          │
│  │ (384-dim)      │                   │Transform │          │
│  └────────┬───────┘                   └──────────┘          │
│           │                                                  │
│           │ [Embedding Vector]                               │
│           ▼                                                  │
│  ┌────────────────┐──[Search]────────▶┌──────────┐          │
│  │ P1.3: Semantic │                   │ D1:      │          │
│  │ Search (Milvus)│◀─[Top-K Results]──│ Milvus   │          │
│  │ top_k=10       │                   │ Vectors  │          │
│  └────────┬───────┘                   └──────────┘          │
│           │                                                  │
│           │ [Context: Past Stories]                          │
│           ▼                                                  │
│  ┌────────────────┐                                          │
│  │ P1.4: Assemble │                                          │
│  │ LLM Prompt     │                                          │
│  │ (Context +     │                                          │
│  │  User Prompt)  │                                          │
│  └────────┬───────┘                                          │
│           │                                                  │
│           │ [Full Prompt]                                    │
│           ▼                                                  │
│  ┌────────────────┐──[API Call]──────▶┌──────────┐          │
│  │ P1.5: LLM      │                   │ Groq     │          │
│  │ Generation     │◀─[Response]───────│ Cloud    │          │
│  │ (llama-3.1)    │                   │ (LLM-1)  │          │
│  └────────┬───────┘                   └──────────┘          │
│           │                                                  │
│           │ [Generated Story Line]                           │
│           ▼                                                  │
│  ┌────────────────┐                                          │
│  │ P1.6: Format   │                                          │
│  │ Response with  │                                          │
│  │ Context Metadata                                          │
│  └────────┬───────┘                                          │
│           │                                                  │
│           ▼                                                  │
│  [Output: {suggestion, context_used, context_count}]        │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

---

## DFD Level 2: Story Verification & Persistence (P2 + P3)

```
┌──────────────────────────────────────────────────────────────┐
│           P2: STORY EDITOR & VERIFIER                         │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  [Input: AI Suggestion + User Edit]                          │
│         │                                                     │
│         ▼                                                     │
│  ┌────────────────┐                                           │
│  │ P2.1: Compare  │                                           │
│  │ Original vs    │                                           │
│  │ Edited Version │                                           │
│  └────────┬───────┘                                           │
│           │ [Delta: user_edited flag]                         │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P2.2: User     │◀─[JWT Token]──[Supabase Auth]            │
│  │ Authentication │                                           │
│  └────────┬───────┘                                           │
│           │ [user_id]                                         │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P2.3: Create   │                                           │
│  │ Verified Story │                                           │
│  │ Line Object    │                                           │
│  └────────┬───────┘                                           │
│           │ [StoryLine{verified=true}]                        │
│           ▼                                                   │
├───────────────────────────────────────────────────────────────┤
│           P3: STORY PERSISTENCE MANAGER                       │
├───────────────────────────────────────────────────────────────┤
│  ┌────────────────┐──[INSERT]────────▶┌──────────┐           │
│  │ P3.1: Store in │                   │ D2:      │           │
│  │ SQLite         │                   │ SQLite   │           │
│  └────────┬───────┘                   │ story_   │           │
│           │                           │ lines    │           │
│           │ [Assigned ID]             └──────────┘           │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P3.2: Generate │                                           │
│  │ Embedding      │                                           │
│  │ (384-dim)      │                                           │
│  └────────┬───────┘                                           │
│           │ [Vector]                                          │
│           ▼                                                   │
│  ┌────────────────┐──[INSERT]────────▶┌──────────┐           │
│  │ P3.3: Store    │                   │ D1:      │           │
│  │ Vector in      │                   │ Milvus   │           │
│  │ Milvus         │                   │ Vectors  │           │
│  └────────┬───────┘                   └──────────┘           │
│           │                                                   │
│           │ [embedding_id]                                    │
│           ▼                                                   │
│  ┌────────────────┐──[UPDATE]────────▶┌──────────┐           │
│  │ P3.4: Link     │                   │ D2:      │           │
│  │ Embedding ID   │                   │ SQLite   │           │
│  │ to Story Line  │                   └──────────┘           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           ▼                                                   │
│  [Output: Stored Story Line with embedding_id]               │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

---

## DFD Level 2: Lore Extraction (P5)

```
┌──────────────────────────────────────────────────────────────┐
│                P5: LORE EXTRACTOR                             │
│                (Background Task: Every 30 min)                │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  [Trigger: Scheduled Task / Manual Request]                  │
│         │                                                     │
│         ▼                                                     │
│  ┌────────────────┐──[SELECT verified]──▶┌──────────┐        │
│  │ P5.1: Fetch    │                      │ D2:      │        │
│  │ Recent Verified│◀─[Story Lines]───────│ SQLite   │        │
│  │ Lines          │                      └──────────┘        │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Array of story texts]                            │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P5.2: Aggregate│                                           │
│  │ Text for       │                                           │
│  │ Analysis       │                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Concatenated story text]                         │
│           ▼                                                   │
│  ┌────────────────┐──[LLM Call]─────────▶┌──────────┐        │
│  │ P5.3: LLM      │                      │ Groq     │        │
│  │ Entity         │◀─[Entities JSON]─────│ Cloud    │        │
│  │ Extraction     │                      │ (LLM-1)  │        │
│  │ Prompt:        │                      └──────────┘        │
│  │ "Extract chars,│                                           │
│  │  locations,    │                                           │
│  │  events, items"│                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Structured Lore JSON]                            │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P5.4: Parse &  │                                           │
│  │ Validate JSON  │                                           │
│  │ (characters,   │                                           │
│  │  locations,    │                                           │
│  │  events, items)│                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Validated Entities]                              │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P5.5: Create   │                                           │
│  │ Embeddings for │                                           │
│  │ Each Entity    │                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Entity + Vector pairs]                           │
│           ▼                                                   │
│  ┌────────────────┐──[INSERT]──────────▶┌──────────┐         │
│  │ P5.6: Store    │                     │ D3:      │         │
│  │ Lore Entries   │                     │ Lore DB  │         │
│  │ in SQLite      │                     │ (SQLite) │         │
│  └────────┬───────┘                     └──────────┘         │
│           │                                                   │
│           │ [lore_entry_ids]                                  │
│           ▼                                                   │
│  ┌────────────────┐──[INSERT Vectors]──▶┌──────────┐         │
│  │ P5.7: Store    │                     │ D1:      │         │
│  │ Lore Vectors   │                     │ Milvus   │         │
│  │ in Milvus      │                     │ lore_    │         │
│  └────────┬───────┘                     │ vectors  │         │
│           │                             └──────────┘         │
│           ▼                                                   │
│  [Output: Enriched Knowledge Base]                           │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

---

## DFD Level 2: Story Canonicalization (P6)

```
┌──────────────────────────────────────────────────────────────┐
│              P6: STORY CANONICALIZER                          │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  [Input: line_ids[] + title]                                 │
│         │                                                     │
│         ▼                                                     │
│  ┌────────────────┐──[SELECT WHERE id IN]─▶┌──────────┐     │
│  │ P6.1: Fetch    │                        │ D2:      │     │
│  │ Story Lines    │◀─[Story Line objects]──│ SQLite   │     │
│  │ by IDs         │                        └──────────┘     │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Array of verified lines]                         │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P6.2: Sort by  │                                           │
│  │ Timestamp      │                                           │
│  │ (Chronological)│                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Ordered story fragments]                         │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P6.3: Assemble │                                           │
│  │ Context for    │                                           │
│  │ LLM-2          │                                           │
│  │ Prompt:        │                                           │
│  │ "Polish into   │                                           │
│  │  coherent      │                                           │
│  │  narrative"    │                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Canonicalization Prompt]                         │
│           ▼                                                   │
│  ┌────────────────┐──[API Call]──────────▶┌──────────┐       │
│  │ P6.4: LLM-2    │                       │ Groq     │       │
│  │ Generation     │◀─[Polished Text]──────│ Cloud    │       │
│  │ (llama-3.1)    │                       │ (LLM-2)  │       │
│  │ temp=0.5       │                       └──────────┘       │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Canonical story text]                            │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P6.5: Calculate│                                           │
│  │ Metadata       │                                           │
│  │ (word count,   │                                           │
│  │  line count)   │                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Metadata object]                                 │
│           ▼                                                   │
│  ┌────────────────┐──[INSERT]──────────────▶┌──────────┐     │
│  │ P6.6: Store    │                         │ D4:      │     │
│  │ Canonical Story│                         │ Canonical│     │
│  │ in Database    │                         │ Stories  │     │
│  └────────┬───────┘                         └──────────┘     │
│           │                                                   │
│           │ [canonical_story_id]                              │
│           ▼                                                   │
│  [Output: {id, title, full_text, original_lines_count}]      │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

---

## DFD Level 2: NFT Minting (P7) - Future Implementation

```
┌──────────────────────────────────────────────────────────────┐
│              P7: NFT MINTER (Work in Progress)                │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  [Input: canonical_story_id]                                 │
│         │                                                     │
│         ▼                                                     │
│  ┌────────────────┐──[SELECT]──────────────▶┌──────────┐     │
│  │ P7.1: Fetch    │                         │ D4:      │     │
│  │ Canonical Story│◀─[Story + line_ids]─────│ Canonical│     │
│  └────────┬───────┘                         │ Stories  │     │
│           │                                 └──────────┘     │
│           │                                                   │
│           │ [Story object]                                    │
│           ▼                                                   │
│  ┌────────────────┐──[SELECT contributions]─▶┌──────────┐    │
│  │ P7.2: Fetch    │                          │ D2:      │    │
│  │ All            │◀─[Contribution array]────│ SQLite   │    │
│  │ Contributions  │                          └──────────┘    │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [user_id → contribution_count mapping]            │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P7.3: Calculate│                                           │
│  │ Co-Authorship  │                                           │
│  │ - Main Author  │                                           │
│  │ - Co-Authors   │                                           │
│  │ - Ownership %  │                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Author metadata]                                 │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P7.4: Generate │                                           │
│  │ NFT Image      │                                           │
│  │ (Story title,  │                                           │
│  │  author graph) │                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Image PNG binary]                                │
│           ▼                                                   │
│  ┌────────────────┐──[Upload]───────────────▶┌──────────┐    │
│  │ P7.5: Upload   │                          │ IPFS     │    │
│  │ Image to IPFS  │◀─[Image CID]─────────────│ Network  │    │
│  └────────┬───────┘                          └──────────┘    │
│           │                                                   │
│           │ [QmXXX... image CID]                              │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P7.6: Create   │                                           │
│  │ NFT Metadata   │                                           │
│  │ JSON           │                                           │
│  │ {title, image, │                                           │
│  │  authors[...]} │                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Metadata JSON]                                   │
│           ▼                                                   │
│  ┌────────────────┐──[Upload]───────────────▶┌──────────┐    │
│  │ P7.7: Upload   │                          │ IPFS     │    │
│  │ Metadata to    │◀─[Metadata CID]──────────│ Network  │    │
│  │ IPFS           │                          └──────────┘    │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [QmYYY... metadata CID]                           │
│           ▼                                                   │
│  ┌────────────────┐                                           │
│  │ P7.8: Create   │                                           │
│  │ mint_nft       │                                           │
│  │ Transaction    │                                           │
│  └────────┬───────┘                                           │
│           │                                                   │
│           │ [Signed Transaction]                              │
│           ▼                                                   │
│  ┌────────────────┐──[Broadcast]────────────▶┌──────────┐    │
│  │ P7.9: Submit to│                          │ D5:      │    │
│  │ Validator      │                          │ Blockchain    │
│  │ Network (PBFT) │◀─[Block Confirmation]────│ Network  │    │
│  └────────┬───────┘                          │ (Validators)  │
│           │                                  └──────────┘    │
│           │ [NFT Token ID]                                    │
│           ▼                                                   │
│  [Output: {token_id, image_cid, metadata_cid, authors[]}]    │
│                                                               │
└───────────────────────────────────────────────────────────────┘
```

---

## Data Dictionary

### Data Flows

| Flow ID | Name | Description | Format |
|---------|------|-------------|--------|
| F1 | Story Prompt | User's creative writing prompt | String (UTF-8, max 500 chars) |
| F2 | AI Suggestion | LLM-generated story continuation | String (JSON: {suggestion, context_used, context_count}) |
| F3 | Verified Line | User-approved story sentence | Object: {id, user_id, line_text, llm_proposed, user_edited, verified, timestamp} |
| F4 | Embedding | 384-dimensional vector representation | float32[384] |
| F5 | Context Request | Semantic search query | Object: {query_text, top_k, filters} |
| F6 | Relevant Stories | Retrieved context from vector DB | Array of {line_id, text, similarity_score} |
| F7 | Lore Entities | Extracted world-building elements | Object: {characters[], locations[], events[], items[]} |
| F8 | Canonical Story | Polished narrative output | Object: {id, title, full_text, original_lines_count, created_at} |
| F9 | NFT Transaction | Blockchain mint request | Object: {story_id, image_cid, metadata_cid, authors[], signatures[]} |
| F10 | Block Confirmation | Consensus approval | Object: {block_index, tx_id, validator_signatures[]} |

### Data Stores

| Store ID | Name | Type | Content | Persistence |
|----------|------|------|---------|-------------|
| D1 | Milvus Vectors | Vector DB | 384-dim embeddings + metadata | RAM + Disk (index) |
| D2 | SQLite Database | Relational DB | story_lines, metadata | Disk (kahani.db) |
| D3 | Lore DB | Relational DB | lore_entries (characters, locations, events, items) | Disk (kahani.db) |
| D4 | Canonical Stories | Relational DB | canonical_stories table | Disk (kahani.db) |
| D5 | Blockchain Ledger | BadgerDB (KV) | Immutable blocks + state | Disk (blockchain/) |
| D6 | IPFS Network | Content-Addressed | Story text, NFT images/metadata | Distributed |

### Processes

| Process ID | Name | Inputs | Outputs | Trigger |
|------------|------|--------|---------|---------|
| P1 | Story Suggestion Generator | User prompt, Vector DB context | AI suggestion + context metadata | User request |
| P2 | Story Editor & Verifier | AI suggestion, User edit, JWT token | Verified story line object | User action |
| P3 | Story Persistence Manager | Verified story line | Stored line + embedding_id | P2 completion |
| P4 | Vector Embedding Storage | Story text | 384-dim vector in Milvus | P3 completion |
| P5 | Lore Extractor | Recent verified lines | Extracted entities (characters, etc.) | Scheduled (30 min) or manual |
| P6 | Story Canonicalizer | line_ids[], title | Polished canonical story | User request |
| P7 | NFT Minter (future) | canonical_story_id | NFT token_id + blockchain entry | User request |

---

## Sequence Diagrams

### Sequence 1: User Creates Story Line

```
User          Frontend       API Proxy      FastAPI       Milvus      Groq LLM
 │                │              │             │            │            │
 │─[1] Enter────▶│              │             │            │            │
 │   prompt      │              │             │            │            │
 │               │              │             │            │            │
 │               │─[2] POST────▶│             │            │            │
 │               │   suggest    │             │            │            │
 │               │              │             │            │            │
 │               │              │─[3] Forward▶│            │            │
 │               │              │             │            │            │
 │               │              │             │─[4] Query─▶│            │
 │               │              │             │   context  │            │
 │               │              │             │            │            │
 │               │              │             │◀[5] Top-K──│            │
 │               │              │             │   results  │            │
 │               │              │             │            │            │
 │               │              │             │─[6] Generate────────────▶│
 │               │              │             │   w/context│            │
 │               │              │             │            │            │
 │               │              │             │◀[7] Story──────────────│
 │               │              │             │   suggestion           │
 │               │              │             │            │            │
 │               │              │◀[8] Return──│            │            │
 │               │              │   JSON      │            │            │
 │               │              │             │            │            │
 │               │◀[9] Response─│             │            │            │
 │               │              │             │            │            │
 │◀[10] Display──│              │             │            │            │
 │    suggestion │              │             │            │            │
 │               │              │             │            │            │
 │─[11] Edit &──▶│              │             │            │            │
 │    Sign       │              │             │            │            │
 │               │              │             │            │            │
 │               │─[12] POST───▶│             │            │            │
 │               │   edit       │             │            │            │
 │               │              │             │            │            │
 │               │              │─[13] Store─▶│            │            │
 │               │              │   in SQLite │            │            │
 │               │              │             │            │            │
 │               │              │             │─[14] Embed▶│            │
 │               │              │             │   & store  │            │
 │               │              │             │            │            │
 │               │              │◀[15] OK─────│            │            │
 │               │              │             │            │            │
 │               │◀[16] Success─│             │            │            │
 │               │              │             │            │            │
 │◀[17] Updated──│              │             │            │            │
 │    story list │              │             │            │            │
```

### Sequence 2: Background Lore Extraction

```
Scheduler      FastAPI      SQLite DB    Groq LLM    Milvus
    │              │             │            │          │
[30 min]           │             │            │          │
    │              │             │            │          │
    │─[1] Trigger─▶│             │            │          │
    │   lore task  │             │            │          │
    │              │             │            │          │
    │              │─[2] SELECT─▶│            │          │
    │              │   verified  │            │          │
    │              │   lines     │            │          │
    │              │             │            │          │
    │              │◀[3] Story───│            │          │
    │              │   texts     │            │          │
    │              │             │            │          │
    │              │─[4] Extract entities────▶│          │
    │              │   prompt    │            │          │
    │              │             │            │          │
    │              │◀[5] Entities│            │          │
    │              │   JSON      │            │          │
    │              │             │            │          │
    │              │─[6] Parse & │            │          │
    │              │   validate  │            │          │
    │              │             │            │          │
    │              │─[7] INSERT─▶│            │          │
    │              │   lore_     │            │          │
    │              │   entries   │            │          │
    │              │             │            │          │
    │              │─[8] Embed──────────────────────────▶│
    │              │   entities  │            │          │
    │              │             │            │          │
    │◀[9] Complete─│             │            │          │
```

### Sequence 3: NFT Minting (Future)

```
User      Frontend    API Proxy   FastAPI   Blockchain   IPFS
 │            │           │          │          │          │
 │─[1] Click─▶│           │          │          │          │
 │  "Mint NFT"│           │          │          │          │
 │            │           │          │          │          │
 │            │─[2] POST─▶│          │          │          │
 │            │  canonicalize        │          │          │
 │            │           │          │          │          │
 │            │           │─[3] ────▶│          │          │
 │            │           │  Polish  │          │          │
 │            │           │          │          │          │
 │            │◀[4] Story────────────│          │          │
 │            │  canonical│          │          │          │
 │            │           │          │          │          │
 │◀[5] Show───│           │          │          │          │
 │  canonical │           │          │          │          │
 │            │           │          │          │          │
 │─[6] POST──▶│           │          │          │          │
 │  mint_nft  │           │          │          │          │
 │            │           │          │          │          │
 │            │───────────────[7] Aggregate────▶│          │
 │            │           │  contributions      │          │
 │            │           │          │          │          │
 │            │           │         [8] Generate│          │
 │            │           │          NFT image  │          │
 │            │           │          │          │          │
 │            │           │          │─[9] Upload image───▶│
 │            │           │          │          │          │
 │            │           │          │◀[10] CID────────────│
 │            │           │          │          │          │
 │            │           │          │─[11] Upload metadata▶│
 │            │           │          │          │          │
 │            │           │          │◀[12] CID────────────│
 │            │           │          │          │          │
 │            │           │          │─[13] mint_nft tx───▶│
 │            │           │          │          │          │
 │            │           │          │          │─[PBFT]   │
 │            │           │          │          │ consensus│
 │            │           │          │          │          │
 │            │           │          │◀[14] Block committed│
 │            │           │          │   token_id          │
 │            │           │          │          │          │
 │            │◀──────────────[15] NFT minted───│          │
 │            │           │  {token_id, CIDs}   │          │
 │            │           │          │          │          │
 │◀[16] Success─          │          │          │          │
 │  + token_id│           │          │          │          │
```

---

## State Transition Diagrams

### Story Line States

```
┌──────────┐
│  DRAFT   │  [User typing in UI]
└────┬─────┘
     │
     │ [User clicks "Summon Suggestion"]
     ▼
┌──────────────┐
│  SUGGESTED   │  [AI returned proposal]
└────┬─────────┘
     │
     │ [User edits or accepts]
     ▼
┌──────────────┐
│   EDITED     │  [Local state in browser]
└────┬─────────┘
     │
     │ [User clicks "Add to Story"]
     ▼
┌──────────────┐
│  SUBMITTED   │  [POST to backend]
└────┬─────────┘
     │
     │ [Backend validates & stores]
     ▼
┌──────────────┐
│   VERIFIED   │  [Persisted in SQLite + Milvus]
└────┬─────────┘
     │
     │ [Background task or manual request]
     ▼
┌──────────────┐
│ IN_CANONICAL │  [Part of polished story]
└────┬─────────┘
     │
     │ [NFT minting (future)]
     ▼
┌──────────────┐
│  MINTED_NFT  │  [Immutable on blockchain]
└──────────────┘
```

### Blockchain Transaction States (Future)

```
┌───────────┐
│  PENDING  │  [Created, not broadcast]
└─────┬─────┘
      │
      │ [Broadcast to validators]
      ▼
┌───────────┐
│PRE_PREPARE│  [Leader proposes block]
└─────┬─────┘
      │
      │ [Validators receive]
      ▼
┌───────────┐
│  PREPARE  │  [Validators send prepare msg]
└─────┬─────┘
      │
      │ [2f+1 prepare messages received]
      ▼
┌───────────┐
│  COMMIT   │  [Validators send commit msg]
└─────┬─────┘
      │
      │ [2f+1 commit messages received]
      ▼
┌───────────┐
│ FINALIZED │  [Block added to chain]
└───────────┘
```

---

## Data Access Patterns

### Read Operations

| Operation | Process | Data Store | Query Pattern | Frequency |
|-----------|---------|------------|---------------|-----------|
| Get story suggestions | P1 | D1 (Milvus) | Vector similarity search | On user prompt |
| Fetch story lines | Frontend | D2 (SQLite) | SELECT * FROM story_lines WHERE verified=true | Page load |
| Retrieve lore | P5, Frontend | D3 (Lore DB) | SELECT * FROM lore_entries WHERE entity_type=? | Background + UI |
| Get canonical story | P6, Frontend | D4 (Canonical) | SELECT * FROM canonical_stories WHERE id=? | User request |
| Fetch NFT metadata | Frontend | D5 (Blockchain) | Get NFT by token_id | NFT gallery view |
| Context retrieval | P1 | D1 (Milvus) | topK(embedding, k=10) | Every suggestion |

### Write Operations

| Operation | Process | Data Store | Write Pattern | Frequency |
|-----------|---------|------------|---------------|-----------|
| Store story line | P3 | D2 (SQLite) | INSERT INTO story_lines | Per contribution |
| Save embedding | P4 | D1 (Milvus) | INSERT vector + metadata | Per contribution |
| Store lore entities | P5 | D3 (Lore DB) | INSERT INTO lore_entries | Every 30 min (batch) |
| Save canonical story | P6 | D4 (Canonical) | INSERT INTO canonical_stories | On finalization |
| Mint NFT transaction | P7 | D5 (Blockchain) | Append block to chain | On NFT mint |
| Upload to IPFS | P7 | D6 (IPFS) | IPFS.add(content) | On NFT creation |

---

## Performance Metrics

### Expected Latencies

| Operation | Target Latency | Bottleneck |
|-----------|----------------|------------|
| Story suggestion | < 3s | Groq LLM API |
| Vector search | < 100ms | Milvus query |
| Store story line | < 200ms | SQLite write + embedding gen |
| Lore extraction | ~10-30s | LLM entity extraction (batch) |
| Canonicalization | ~5-15s | LLM-2 polish |
| NFT minting | ~3-5s | PBFT consensus (future) |

### Throughput Estimates

| Metric | Current (AI Only) | Future (With Blockchain) |
|--------|-------------------|--------------------------|
| Concurrent users | ~100 | ~1000 (with observers) |
| Stories/hour | ~500 | ~500 (write limited by validators) |
| Reads/second | ~50 | ~500 (distributed observers) |
| NFTs minted/hour | N/A | ~200 (PBFT throughput) |

---

This DFD documentation provides complete visibility into data transformations across the Kahani system, from user input through AI processing to future blockchain persistence.
