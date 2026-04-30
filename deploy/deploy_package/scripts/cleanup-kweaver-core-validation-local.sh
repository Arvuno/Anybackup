#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "${ROOT_DIR}/scripts/lib/kweaver-core-local-common.sh"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/cleanup-kweaver-core-validation-local.sh
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

require_root
require_cmd kubectl helm python3 rpm yum systemctl

log_info "Cleaning namespaces, ingress classes, legacy Proton services, and deployment residue."
cleanup_kweaver_core_validation
log_info "Cleanup finished."
