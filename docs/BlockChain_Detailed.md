# Storytelling Blockchain

Collaborative storytelling playground where contributions are collected on-chain and minted into NFTs. A lightweight PBFT-style consensus keeps validator nodes in sync, Supabase powers authentication and wallet provisioning, and IPFS backs NFT media.

## Highlights
- Go 1.24 service with modular packages (`internal/api`, `internal/blockchain`, `internal/consensus`).
- PBFT consensus with in-memory gossip transport for local clusters.
- Supabase auth middleware and poller that bridges Supabase user wallets on-chain.
- IPFS client with automatic in-memory fallback for local development.
- REST + websocket APIs for contributions, minting, health, and observability.

## Repository Layout
- `cmd/devnode` – entrypoint for validator HTTP node and consensus bootstrap.
- `internal/api` – REST handlers, middleware, websocket streaming.
- `internal/blockchain` – blocks, validation rules, NFT minting helpers.
- `internal/consensus` – PBFT service wiring, sharding utilities.
- `internal/storage` – Badger persistence, IPFS clients, memory shims.
- `internal/supabase` – auth middleware, REST client, wallet poller.
- `internal/wallet` – key generation, signing, encrypted storage.
- `pkg/utils` – shared helpers (hashing, signatures).
- `scripts/` – devnode launcher, cluster bootstrap, smoke tests.
- `docs/operator-guide.md` – advanced deployment/operations notes.

## Prerequisites
- Go 1.24+
- Bash-compatible shell for helper scripts (Windows WSL/Git Bash OK).
- Optional: public IPFS endpoint (`IPFS_API`/`DEVNODE_IPFS_API`).
- Optional: Supabase project (URL, anon key, service role key) for real auth and wallet sync.

## Quick Start
1. Install dependencies: `go mod tidy`
2. Run a single validator/dev node:
   ```bash
   go run ./cmd/devnode
   ```
   The node listens on `:8080` by default. Health probe lives at `http://localhost:8080/api/health`.
3. Stop with `Ctrl+C`. State persists in `devnode-data` unless you change `DEVNODE_DATA_DIR`.

### Using the helper script
`scripts/run-devnode.sh` loads `.env`, exports supported `DEVNODE_*` variables, and launches the node:
```bash
./scripts/run-devnode.sh --http :8080 --cluster-size 1
```

### Multi-node cluster
Start the orchestrated cluster (single process, in-memory transport):
```bash
./scripts/run-dev-cluster.sh "node-1,node-2,node-3" 1 :8080
```
Flags: `<peersCSV> <faultTolerance> <httpAddr> [additional go run flags]`. Each peer still appears in consensus routing and sharding.

## Configuration Reference
Environment variables and matching CLI flags (flags win):

| Purpose | Env Var | Flag | Default |
| --- | --- | --- | --- |
| Node identifier | `DEVNODE_NODE_ID` | `--node` | `node-1` |
| HTTP bind address | `DEVNODE_HTTP_ADDR` | `--http` | `:8080` |
| Wallet passphrase | `DEVNODE_WALLET_PASSPHRASE` | `--passphrase` | `local-passphrase` |
| Allowed CORS origins (comma list) | `DEVNODE_ALLOWED_ORIGINS` | `--origins` | `http://localhost:3000` |
| Fault tolerance (PBFT f) | `DEVNODE_FAULT_TOLERANCE` | `--fault` | `0` |
| Cluster size auto-provision | `DEVNODE_CLUSTER_SIZE` | `--cluster-size` | `1` |
| Explicit peers (comma list) | `DEVNODE_PEERS` | `--peers` | auto from cluster size |
| Seed Supabase users | `DEVNODE_SEED_USERS` | `--seed-users` | `user-123` |
| Persistent store path | `DEVNODE_DATA_DIR` | `--data-dir` | `devnode-data` |
| IPFS HTTP API | `DEVNODE_IPFS_API` / `IPFS_API` | `--ipfs-api` | in-memory fallback |
| Supabase URL | `SUPABASE_URL` | n/a | disabled if empty |
| Supabase anon key | `SUPABASE_ANON_KEY` | n/a | disabled if empty |
| Supabase service role key | `SUPABASE_SERVICE_KEY` or `SUPABASE_SERVICE_ROLE_KEY` | n/a | poller disabled if empty |
| Supabase poll interval | `SUPABASE_POLL_INTERVAL` | `--supabase-poll-interval` | 30s |

Values can live in an `.env` file loaded by the helper scripts.

## Supabase Integration
- When `SUPABASE_URL` and `SUPABASE_ANON_KEY` are set, the API uses Supabase JWT verification for authenticated routes.
- Providing `SUPABASE_SERVICE_KEY` enables the wallet poller that upserts Supabase wallet rows onto the blockchain every poll interval.
- Without Supabase creds, the dev token verifier accepts any non-empty token string (token itself maps to user ID).

## IPFS Integration
- Configure `DEVNODE_IPFS_API` (e.g. `http://127.0.0.1:5001`) for real uploads.
- When omitted, an in-memory IPFS mock stores uploaded JSON/image data for the process lifetime.

## Smoke Test
Run minimal checks against a running node (optionally boot the node inline):
```bash
./scripts/api-smoke-test.sh http://localhost:8080
START_NODE=1 ./scripts/api-smoke-test.sh https://example.ngrok.app
```
Outputs HTTP status for key endpoints and optionally tails logs if the node fails to start.

## API Overview
Base URL defaults to `http://localhost:8080`.

| Method | Path | Auth | Description |
| --- | --- | --- | --- |
| GET | `/api/health` | none | Combined status (blocks, consensus wiring, metrics, uptime). |
| GET | `/api/health/live` | none | Simple liveness probe. |
| GET | `/api/health/ready` | none | Readiness including component flags; 503 until consensus available. |
| GET | `/api/blockchain` | none | Full block list and current registry state snapshot. |
| GET | `/api/wallet/{supabaseUserID}` | none | Wallet entry keyed by Supabase user ID. |
| GET | `/api/story/{storyID}` | none | Contributions, author aggregation, minted NFTs, latest title/summary. |
| GET | `/api/nft/{tokenID}` | none | Stored NFT metadata (authors, IPFS CIDs, summary). |
| GET | `/api/nft/{tokenID}/authors` | none | Main + co-author roster for the NFT. |
| GET (WS) | `/api/events` | Origin-gated | Websocket stream of queued transactions and committed blocks. |
| POST | `/api/story/contribute` | Bearer JWT | Submit a signed story line. |
| POST | `/api/story/{storyID}/mint` | Bearer JWT | Mint story into an NFT (main author only). |

### Example Calls
```bash
# Public health check
curl -s http://localhost:8080/api/health | jq

# Lookup wallet
curl -s http://localhost:8080/api/wallet/user-123 | jq

# Fetch story data with latest title/summary
curl -s http://localhost:8080/api/story/story-42 | jq

# Submit contribution (token doubles as user ID in dev mode)
curl -s -X POST http://localhost:8080/api/story/contribute \
  -H "Authorization: Bearer user-123" \
  -H "Content-Type: application/json" \
  -d '{
    "story_id": "story-42",
    "story_line": "Validators rallied around the shard."
  }' | jq

# Mint the story (main author only)
curl -s -X POST http://localhost:8080/api/story/story-42/mint \
  -H "Authorization: Bearer user-123" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Adventure Across Chains",
    "summary": "A quest that spans validators and shards."
  }' | jq

# Websocket events (requires ws client)
wscat -H "Origin: http://localhost" -c ws://localhost:8080/api/events
```

Typical `/api/story/{id}` response after minting:
```json
{
  "story_id": "story-42",
  "title": "Adventure Across Chains",
  "summary": "A quest that spans validators and shards.",
  "contributions": [
    {"contributor_id": "user-123", "story_line": "Once upon a shard..."},
    {"contributor_id": "user-456", "story_line": "Validators rallied together."}
  ],
  "authors": [
    {"supabase_user_id": "user-123", "contribution_count": 3, "ownership_percentage": 60},
    {"supabase_user_id": "user-456", "contribution_count": 2, "ownership_percentage": 40}
  ],
  "nfts": [
    {
      "token_id": "nft_story-42_f8c1c6a2d0b3",
      "title": "Adventure Across Chains",
      "summary": "A quest that spans validators and shards.",
      "image_ipfs_cid": "bafy...image",
      "metadata_ipfs_cid": "bafy...metadata",
      "minted_at": 1729875600,
      "block_index": 3
    }
  ]
}
```

## Development Flow
1. Ensure a wallet exists for your Supabase user (auto-seeded or via poller).
2. Call `/api/story/contribute` repeatedly to build the story; contributions are signed with the contributor's wallet key.
3. When ready, the lead contributor calls `/api/story/{id}/mint` with title + summary. An NFT is minted, metadata is uploaded to IPFS, and a `mint_nft` transaction enters consensus.
4. Retrieve the minted NFT through `/api/nft/{tokenID}` or check marketplace metadata via the IPFS CID.

## Running Tests
```bash
go test ./...
```
All packages include unit tests; consensus integration tests rely on the in-memory network transport.

## Troubleshooting
- Health ready endpoint returns `degraded`: verify consensus wiring (node ID present in peer list, cluster size >= 1).
- Minting fails with `story has no authors`: ensure you contribute at least once before minting.
- `supabase: token verification failed`: supply a valid Supabase JWT, or omit credentials so dev mode accepts any token.
- IPFS CID missing: set `DEVNODE_IPFS_API` to a reachable IPFS gateway if you need persistence beyond process lifetime.

For deeper operational guidance, read `docs/operator-guide.md`.
