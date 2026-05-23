#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALLER="$(find "${SCRIPT_DIR}" -maxdepth 1 -type f -name 'anybackup-foundation-*.bin' | head -n 1 || true)"

if [[ -z "${INSTALLER}" ]]; then
  echo "[foundation] no Foundation installer found under ${SCRIPT_DIR}"
  echo "[foundation] add the real offline package and rerun this script"
  exit 0
fi

echo "[foundation] installer discovered: ${INSTALLER}"
echo "[foundation] automatic Foundation installation is not implemented in this repackaged delivery yet."
