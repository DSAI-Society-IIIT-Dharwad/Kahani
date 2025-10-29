#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ -f "${ROOT_DIR}/.env" ]; then
  # shellcheck disable=SC2046
  set -a
  # shellcheck disable=SC1091
  . "${ROOT_DIR}/.env"
  set +a
fi

if [ -n "${DEVNODE_DATA_DIR:-}" ]; then
  mkdir -p "${DEVNODE_DATA_DIR}"
fi

for var in \
  DEVNODE_NODE_ID \
  DEVNODE_HTTP_ADDR \
  DEVNODE_ALLOWED_ORIGINS \
  DEVNODE_FAULT_TOLERANCE \
  DEVNODE_SEED_USERS \
  DEVNODE_CLUSTER_SIZE \
  DEVNODE_PEERS \
  DEVNODE_DATA_DIR \
  DEVNODE_IPFS_API; do
  value="${!var-}"
  if [ -n "$value" ]; then
    export "$var"
  fi
done

cd "${ROOT_DIR}"

go run ./cmd/devnode "$@"
