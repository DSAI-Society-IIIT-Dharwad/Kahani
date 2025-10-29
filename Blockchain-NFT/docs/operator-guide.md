# Operator Guide

## Quick Start

1. Build and start a single validator dev node:
   ```bash
   go run ./cmd/devnode
   ```
2. Override runtime settings via flags or environment variables, for example:
   ```bash
   DEVNODE_NODE_ID=node-2 DEVNODE_HTTP_ADDR=:9090 go run ./cmd/devnode --cluster-size=3
   ```
3. Seeded wallets and consensus will initialise automatically. Confirm readiness at `http://localhost:8080/api/health/ready`.

## Runtime Configuration

The dev node reads configuration from flags and matching environment variables (flags take precedence).

| Flag | Environment Variable | Description | Default |
|------|----------------------|-------------|---------|
| `--node` | `DEVNODE_NODE_ID` | Primary validator identifier. | `node-1` |
| `--http` | `DEVNODE_HTTP_ADDR` | HTTP listen address. | `:8080` |
| `--peers` | `DEVNODE_PEERS` | Comma-separated list of peer node IDs. | Auto-generated from cluster size |
| `--cluster-size` | `DEVNODE_CLUSTER_SIZE` | Auto-provision peer count when `--peers` omitted. | `1` |
| `--origins` | `DEVNODE_ALLOWED_ORIGINS` | Comma-separated CORS allowlist. | `http://localhost:3000` |
| `--passphrase` | `DEVNODE_WALLET_PASSPHRASE` | Wallet key derivation passphrase. | `local-passphrase` |
| `--fault` | `DEVNODE_FAULT_TOLERANCE` | PBFT fault tolerance parameter. | `0` |
| `--seed-users` | `DEVNODE_SEED_USERS` | Comma-separated Supabase IDs to pre-provision wallets. | `user-123` |

## Multi-Node Clusters

- Use `--peers=node-1,node-2,node-3` (or the `DEVNODE_PEERS` env var) to specify an explicit cluster membership.
- When only `--cluster-size` is supplied, node IDs are generated as `node-1..node-N`; the primary node is always the first entry.
- The harness spins up an in-memory gossip transport for all peers so they can reach consensus within a single process. Use separate terminal instances with matching peer lists to simulate distributed execution.
- Wallet and API transactions are sharded across the configured node IDs using deterministic hashing.

## Health and Observability

- `GET /api/health`: Combined status payload with block counts, consensus attachment details, metrics, and uptime.
- `GET /api/health/live`: Liveness probe suitable for container orchestrators.
- `GET /api/health/ready`: Readiness probe; returns `503` until consensus wiring and dependencies are in place.
- Structured logs (JSON) are emitted to stdout with cluster metadata.

## Troubleshooting

- **Readiness shows `degraded`**: Ensure consensus has been attached—confirm that the node ID in `--node` appears in the `--peers` list and that no other instance is holding the same identifier.
- **Consensus errors on startup**: Increase `--fault` only when enough peers are configured; negative values are rejected during bootstrap.
- **CORS blocked requests**: Update `--origins` or `DEVNODE_ALLOWED_ORIGINS` to include your client URL.
- **Wallet provisioning skipped**: Check the `--seed-users` value; any generation errors are logged at WARN level with the offending user ID.
