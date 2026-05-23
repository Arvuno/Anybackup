#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 ]]; then
  echo "usage: $0 <sandbox-image>"
  echo "example: $0 swr.example.com/kweaver/python-basic:latest"
  exit 64
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
SANDBOX_IMAGE="$1"
WORK_DIR="${ROOT_DIR}/.work/sandbox"

mkdir -p "${WORK_DIR}"
cp "${ROOT_DIR}/bin/foundation-cli" "${WORK_DIR}/foundation-cli"

cat > "${WORK_DIR}/Dockerfile" <<EOF
FROM ${SANDBOX_IMAGE}
COPY foundation-cli /usr/local/bin/foundation-cli
RUN chmod +x /usr/local/bin/foundation-cli
EOF

echo "[inject-cli] building overlay image with the same tag: ${SANDBOX_IMAGE}"
docker build -t "${SANDBOX_IMAGE}" "${WORK_DIR}"
echo "[inject-cli] exporting/importing the image is environment-specific; continue with docker save / ctr import on the target node."
