#!/usr/bin/env bash
set -euo pipefail

EXPERIENCE_KN_ID="mysql_recovery_experience"
EXPERIENCE_KN_NAME="MySQL数据库恢复经验知识网络"
RUN_KN_ID="mysql_recovery_run"
RUN_KN_NAME="MySQL数据库恢复运行知识网络"

usage() {
  cat <<'EOF'
Usage:
  ./uninstall.sh [options]

Options:
  --kweaver-ip <value>         KWeaver host/IP or full base URL
  --url <value>                Compatibility alias of --kweaver-ip
  --username <value>           KWeaver username for authentication
  --password <value>           KWeaver password for authentication
  --biz-domain <value>         Business domain (default: bd_public)
  --skip-login                 Skip kweaver auth login
  --insecure                   Pass --insecure to kweaver auth login
  -h, --help                   Show this help

Environment variables:
  KWEAVER_IP
  KWEAVER_BASE_URL
  KWEAVER_USERNAME
  KWEAVER_PASSWORD
  KWEAVER_BUSINESS_DOMAIN
  KWEAVER_SKIP_LOGIN
  KWEAVER_INSECURE

Notes:
  1. This uninstaller deletes two fixed knowledge networks by exact name and expected ID.
  2. Missing networks are reported and skipped.
EOF
}

log() {
  printf '[INFO] %s\n' "$*"
}

error() {
  printf '[ERROR] %s\n' "$*" >&2
  exit 1
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || error "Required command not found: $1"
}

normalize_base_url() {
  local raw="$1"
  if [[ "$raw" == http://* || "$raw" == https://* ]]; then
    printf '%s\n' "$raw"
    return 0
  fi
  printf 'https://%s\n' "$raw"
}

find_kn_id_by_name() {
  local json_input="$1"
  local expected_name="$2"
  printf '%s' "$json_input" | python3 - "$expected_name" <<'PY'
import json
import sys

expected_name = sys.argv[1]
text = sys.stdin.read().strip()
if not text:
    sys.exit(1)

try:
    data = json.loads(text)
except json.JSONDecodeError:
    sys.exit(1)

def iter_items(node):
    if isinstance(node, list):
        return node
    if isinstance(node, dict):
        for key in ("entries", "data", "items", "records", "knowledge_networks"):
            value = node.get(key)
            if isinstance(value, list):
                return value
            if isinstance(value, dict):
                nested = iter_items(value)
                if nested:
                    return nested
    return []

matches = []
for item in iter_items(data):
    if not isinstance(item, dict):
        continue
    name = item.get("name") or item.get("kn_name") or item.get("knowledge_network_name")
    if str(name) != expected_name:
        continue
    kn_id = item.get("id") or item.get("kn_id") or item.get("knowledge_network_id")
    if kn_id:
        matches.append(str(kn_id))

if len(matches) != 1:
    sys.exit(1)

print(matches[0])
PY
}

delete_bkn() {
  local expected_name="$1"
  local expected_id="$2"

  log "Resolving knowledge network: $expected_name"
  local list_output
  list_output="$(kweaver bkn list --name "$expected_name" -bd "$KWEAVER_BUSINESS_DOMAIN")"

  local resolved_id
  if ! resolved_id="$(find_kn_id_by_name "$list_output" "$expected_name")"; then
    log "Knowledge network not found, skip delete: $expected_name"
    return 0
  fi

  if [[ "$resolved_id" != "$expected_id" ]]; then
    error "Knowledge network ID mismatch for $expected_name: expected $expected_id, got $resolved_id"
  fi

  log "Deleting knowledge network: $expected_name ($resolved_id)"
  kweaver bkn delete "$resolved_id" --yes -bd "$KWEAVER_BUSINESS_DOMAIN"
}

KWEAVER_BASE_URL="${KWEAVER_BASE_URL:-${KWEAVER_IP:-}}"
KWEAVER_USERNAME="${KWEAVER_USERNAME:-}"
KWEAVER_PASSWORD="${KWEAVER_PASSWORD:-}"
KWEAVER_BUSINESS_DOMAIN="${KWEAVER_BUSINESS_DOMAIN:-bd_public}"
KWEAVER_SKIP_LOGIN="${KWEAVER_SKIP_LOGIN:-0}"
KWEAVER_INSECURE="${KWEAVER_INSECURE:-0}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --kweaver-ip)
      KWEAVER_BASE_URL="${2:-}"
      shift 2
      ;;
    --url)
      KWEAVER_BASE_URL="${2:-}"
      shift 2
      ;;
    --username)
      KWEAVER_USERNAME="${2:-}"
      shift 2
      ;;
    --password)
      KWEAVER_PASSWORD="${2:-}"
      shift 2
      ;;
    --biz-domain)
      KWEAVER_BUSINESS_DOMAIN="${2:-}"
      shift 2
      ;;
    --skip-login)
      KWEAVER_SKIP_LOGIN="1"
      shift
      ;;
    --insecure)
      KWEAVER_INSECURE="1"
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      error "Unknown argument: $1"
      ;;
  esac
done

require_command "kweaver"
require_command "python3"

if [[ "$KWEAVER_SKIP_LOGIN" != "1" ]]; then
  [[ -n "$KWEAVER_BASE_URL" ]] || error "KWEAVER_IP or KWEAVER_BASE_URL is required when login is enabled."
  [[ -n "$KWEAVER_USERNAME" ]] || error "KWEAVER_USERNAME is required when login is enabled."
  [[ -n "$KWEAVER_PASSWORD" ]] || error "KWEAVER_PASSWORD is required when login is enabled."
  KWEAVER_BASE_URL="$(normalize_base_url "$KWEAVER_BASE_URL")"

  login_args=(auth login "$KWEAVER_BASE_URL" -u "$KWEAVER_USERNAME" -p "$KWEAVER_PASSWORD")
  if [[ "$KWEAVER_INSECURE" == "1" ]]; then
    login_args+=(--insecure)
  fi

  log "Logging in to KWeaver..."
  kweaver "${login_args[@]}"
fi

delete_bkn "$RUN_KN_NAME" "$RUN_KN_ID"
delete_bkn "$EXPERIENCE_KN_NAME" "$EXPERIENCE_KN_ID"

printf '\n'
printf 'BKN uninstall completed.\n'
