#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ANSIBLE_DIR="${ROOT_DIR}/ansible"
COMMON_LIB="${ROOT_DIR}/scripts/lib/ansible-common.sh"
LOCAL_CLEANUP_SCRIPT="${ROOT_DIR}/scripts/cleanup-kweaver-core-validation-local.sh"
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
  ./scripts/cleanup-kweaver-core-validation.sh [--inventory /path/to/inventory.ini]
  ./scripts/cleanup-kweaver-core-validation.sh --local
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

if [[ "${LOCAL_MODE}" == "true" ]]; then
  exec "${LOCAL_CLEANUP_SCRIPT}"
fi

ensure_ansible_playbook

echo "=== 清理 109 验证环境上的中间件 / Ingress / KWeaver Core 残留 ==="
ansible-playbook \
  -i "${INVENTORY_PATH}" \
  "${ANSIBLE_DIR}/site.yml" \
  -e "deployment_profile=kweaver-core-only" \
  --tags cleanup-core-validation

echo "=== 清理完成 ==="
echo "说明：该脚本会清理 validation 范围内的 K8s 资源，以及本次离线部署使用的 Proton 节点侧残留。"
