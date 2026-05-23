#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage:
  ./uninstall.sh [options]

Options:
  --agent-id <value>        Agent ID to remove
  --biz-domain <value>      Business domain (default: bd_public)
  --kweaver-ip <value>      KWeaver host/IP or full base URL
  --url <value>             Compatibility alias of --kweaver-ip
  --username <value>        KWeaver username for authentication
  --password <value>        KWeaver password for authentication
  --skip-login              Skip kweaver auth login
  --insecure                Pass --insecure to kweaver auth login
  -h, --help                Show this help

Environment variables:
  AGENT_ID
  KWEAVER_IP
  KWEAVER_BUSINESS_DOMAIN
  KWEAVER_BASE_URL
  KWEAVER_USERNAME
  KWEAVER_PASSWORD
  KWEAVER_SKIP_LOGIN
  KWEAVER_INSECURE

Notes:
  1. If AGENT_ID is omitted, the script resolves the agent by the fixed name MySQL恢复Agent.
  2. This script only removes the decision agent.
  3. The linked knowledge network is not deleted.
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

extract_agent_id_by_exact_name() {
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
        for key in ("entries", "data", "items", "records"):
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
    name = item.get("name") or item.get("agent_name") or item.get("title")
    if str(name) != expected_name:
        continue
    agent_id = item.get("id") or item.get("agent_id")
    if agent_id:
        matches.append(str(agent_id))

if len(matches) != 1:
    sys.exit(1)

print(matches[0])
PY
}

AGENT_ID="${AGENT_ID:-}"
AGENT_NAME="MySQL恢复Agent"
KWEAVER_BUSINESS_DOMAIN="${KWEAVER_BUSINESS_DOMAIN:-bd_public}"
KWEAVER_BASE_URL="${KWEAVER_BASE_URL:-${KWEAVER_IP:-}}"
KWEAVER_USERNAME="${KWEAVER_USERNAME:-}"
KWEAVER_PASSWORD="${KWEAVER_PASSWORD:-}"
KWEAVER_SKIP_LOGIN="${KWEAVER_SKIP_LOGIN:-0}"
KWEAVER_INSECURE="${KWEAVER_INSECURE:-0}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --agent-id)
      AGENT_ID="${2:-}"
      shift 2
      ;;
    --kweaver-ip)
      KWEAVER_BASE_URL="${2:-}"
      shift 2
      ;;
    --biz-domain)
      KWEAVER_BUSINESS_DOMAIN="${2:-}"
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

if [[ -z "$AGENT_ID" ]]; then
  log "Resolving installed decision agent by fixed name..."
  personal_agents_output="$(kweaver agent personal-list --name "$AGENT_NAME" -bd "$KWEAVER_BUSINESS_DOMAIN" --pretty)"
  if ! AGENT_ID="$(extract_agent_id_by_exact_name "$personal_agents_output" "$AGENT_NAME")"; then
    error "Failed to resolve a unique agent by exact name: $AGENT_NAME"
  fi
fi

log "Target agent ID: $AGENT_ID"

log "Trying to unpublish the agent before delete..."
if ! kweaver agent unpublish "$AGENT_ID"; then
  log "Unpublish failed or was not required. Continue with delete."
fi

log "Deleting decision agent..."
kweaver agent delete "$AGENT_ID" -y

printf '\n'
printf 'Uninstall completed.\n'
printf 'Deleted Agent ID: %s\n' "$AGENT_ID"
