#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ANSIBLE_DIR="${ROOT_DIR}/ansible"
COMMON_LIB="${ROOT_DIR}/scripts/lib/ansible-common.sh"
INVENTORY_PATH="${ANSIBLE_DIR}/inventory.ini"
LOCAL_MODE="false"
TEMP_INVENTORY_PATH=""

source "${COMMON_LIB}"

cleanup_temp_inventory() {
  if [[ -n "${TEMP_INVENTORY_PATH}" && -f "${TEMP_INVENTORY_PATH}" ]]; then
    rm -f "${TEMP_INVENTORY_PATH}"
  fi
}

trap cleanup_temp_inventory EXIT

usage() {
  cat <<'EOF'
Usage:
  bash scripts/verify-sandbox-overlay-local.sh --local
  bash scripts/verify-sandbox-overlay-local.sh --inventory /path/to/inventory.ini

This is a non-destructive dry-run. It detects the current KWeaver sandbox image,
checks live sandbox imagePullPolicy values, renders the overlay Dockerfile, and
skips docker build, image import, and pod deletion.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --inventory)
      shift
      if [[ $# -eq 0 ]]; then
        echo "ERROR: --inventory requires a value" >&2
        usage
        exit 1
      fi
      INVENTORY_PATH="$1"
      ;;
    --local)
      LOCAL_MODE="true"
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "ERROR: unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
  shift
done

if [[ "${LOCAL_MODE}" == "true" && "${INVENTORY_PATH}" != "${ANSIBLE_DIR}/inventory.ini" ]]; then
  echo "ERROR: --local cannot be combined with --inventory" >&2
  usage
  exit 1
fi

ensure_ansible_playbook

if [[ "${LOCAL_MODE}" == "true" ]]; then
  TEMP_INVENTORY_PATH="$(create_local_inventory)"
  INVENTORY_PATH="${TEMP_INVENTORY_PATH}"
fi

ansible-playbook \
  -i "${INVENTORY_PATH}" \
  "${ANSIBLE_DIR}/site.yml" \
  -e "deployment_profile=full" \
  -e "sandbox_overlay_enabled=true" \
  -e "sandbox_overlay_dry_run=true" \
  -e "sandbox_overlay_restart_pods=false" \
  --tags sandbox-overlay
