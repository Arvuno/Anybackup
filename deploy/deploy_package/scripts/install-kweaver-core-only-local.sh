#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "${ROOT_DIR}/scripts/lib/kweaver-core-local-common.sh"

usage() {
  cat <<'EOF'
Usage:
  ./scripts/install-kweaver-core-only-local.sh
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

require_root
require_cmd kubectl helm git python3 awk grep sed hostname
ensure_package_assets_present
prepare_runtime_dirs
stage_runtime_assets

log_info "Resolved public access host: ${PUBLIC_IP}"
log_info "Stage 1/4: package preflight passed."

log_info "Stage 2/4: deploying v9_infra."
deploy_v9_infra

log_info "Stage 3/4: cloning and installing KWeaver Core online."
install_kweaver_core

log_info "Stage 4/4: running final verification."
verify_kweaver_core_only

log_info "Core-only online deployment completed."
