#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

PEERS_DEFAULT="node-1,node-2,node-3,node-4"
FAULT_DEFAULT=1
HTTP_DEFAULT=":8080"

PEERS_INPUT="${1:-$PEERS_DEFAULT}"
if [ "$#" -ge 1 ]; then
	shift
fi

FAULT_INPUT="${1:-$FAULT_DEFAULT}"
if [ "$#" -ge 1 ]; then
	shift
fi

HTTP_INPUT="${1:-$HTTP_DEFAULT}"
if [ "$#" -ge 1 ]; then
	shift
fi

IFS=',' read -r -a __cluster_nodes <<<"${PEERS_INPUT}"
CLUSTER_SIZE=${#__cluster_nodes[@]}
if [ "${CLUSTER_SIZE}" -lt 1 ]; then
	CLUSTER_SIZE=1
fi

export DEVNODE_PEERS="${PEERS_INPUT}"
export DEVNODE_FAULT_TOLERANCE="${FAULT_INPUT}"
export DEVNODE_CLUSTER_SIZE="${CLUSTER_SIZE}"
export DEVNODE_HTTP_ADDR="${HTTP_INPUT}"

exec "${ROOT_DIR}/scripts/run-devnode.sh" "$@"
