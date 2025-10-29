#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: api-smoke-test.sh [base_url]

Runs a lightweight smoke test against a running dev node.

Environment variables:
  SMOKE_USER_ID     Wallet user ID to query (default: user-123).
  START_NODE        Set to 1 to auto-start run-devnode.sh (default: 0).
  HEALTH_WAIT       Seconds to wait for dev node health check (default: 30).

Examples:
  ./scripts/api-smoke-test.sh http://localhost:8080
  SUPABASE_BEARER="token" ./scripts/api-smoke-test.sh https://example.ngrok.app
EOF
}

if [[ ${1:-} == "-h" || ${1:-} == "--help" ]]; then
  usage
  exit 0
fi

BASE_URL=${1:-http://localhost:8080}
USER_ID=${SMOKE_USER_ID:-user-123}
START_NODE=${START_NODE:-0}
HEALTH_WAIT=${HEALTH_WAIT:-30}
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DEVNODE_PID=""
DEVNODE_LOG=""

log() {
  printf '==> %s\n' "$1"
}

cleanup() {
  if [[ -n $DEVNODE_PID ]]; then
    log "Stopping dev node (PID $DEVNODE_PID)"
    kill "$DEVNODE_PID" >/dev/null 2>&1 || true
    wait "$DEVNODE_PID" >/dev/null 2>&1 || true
  fi

  if [[ -n $DEVNODE_LOG && -f $DEVNODE_LOG ]]; then
    log "Dev node logs: $DEVNODE_LOG"
  fi
}

request() {
  local method=$1
  local path=$2
  local body=${3-}
  local expect=${4-}
  local url="${BASE_URL}${path}"

  local tmp
  tmp=$(mktemp)

  local args=("-sS" "-o" "$tmp" "-w" "%{http_code}" "-X" "$method")
  args+=("-H" "Accept: application/json")
  if [[ -n $body ]]; then
    args+=("-H" "Content-Type: application/json" "--data" "$body")
  fi

  local status
  if ! status=$(curl "${args[@]}" "$url"); then
    printf 'Request %s %s failed to connect\n' "$method" "$path" >&2
    rm -f "$tmp"
    exit 1
  fi

  if [[ -n $expect && $status != "$expect" ]]; then
    printf 'Request %s %s failed: expected %s got %s\nResponse:\n' "$method" "$path" "$expect" "$status" >&2
    cat "$tmp" >&2
    rm -f "$tmp"
    exit 1
  fi

  LAST_STATUS=$status
  cat "$tmp"
  echo
  rm -f "$tmp"
}

if [[ "${START_NODE}" == "1" ]]; then
  log "Starting dev node via run-devnode.sh"
  DEVNODE_LOG=$(mktemp -t devnode-smoke-XXXX.log)
  "${ROOT_DIR}/scripts/run-devnode.sh" >"$DEVNODE_LOG" 2>&1 &
  DEVNODE_PID=$!
  trap cleanup EXIT INT TERM

  log "Dev node started (PID $DEVNODE_PID); logs: $DEVNODE_LOG"
  log "Waiting for dev node health endpoint (timeout ${HEALTH_WAIT}s)"

  ready=0
  for ((i = 0; i < HEALTH_WAIT; i++)); do
    if curl -fsS "${BASE_URL}/api/health" >/dev/null 2>&1; then
      ready=1
      break
    fi
    sleep 1
  done

  if [[ $ready -ne 1 ]]; then
    log "Dev node failed to become healthy in ${HEALTH_WAIT}s"
    tail -n 40 "$DEVNODE_LOG" || true
    exit 1
  fi

  log "Dev node is healthy"
fi

log "Base URL: ${BASE_URL}"

log "Checking /api/health"
request GET "/api/health" '' 200 >/dev/null
log "Health status: ${LAST_STATUS}"

log "Checking /api/blockchain"
request GET "/api/blockchain" '' 200 >/dev/null
log "Blockchain status: ${LAST_STATUS}"

log "Checking wallet for ${USER_ID}"
request GET "/api/wallet/${USER_ID}" '' '' >/dev/null || true
log "Wallet status: ${LAST_STATUS:-n/a}"

log "Listing NFTs requires prior minting; skipping authenticated flows"

log "Smoke test complete"
