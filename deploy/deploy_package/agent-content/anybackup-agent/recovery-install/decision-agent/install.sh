#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROMPTS_DIR="$ROOT_DIR/prompts"
SYSTEM_PROMPT_FILE="$PROMPTS_DIR/system_prompt.md"
DOLPHIN_FILE="$PROMPTS_DIR/dolphin.txt"

usage() {
  cat <<'EOF'
Usage:
  ./install.sh [options]

Options:
  --kweaver-ip <value>         KWeaver host/IP or full base URL
  --url <value>                Compatibility alias of --kweaver-ip
  --username <value>           KWeaver username for authentication
  --password <value>           KWeaver password for authentication
  --biz-domain <value>         Business domain (default: bd_public)
  --agent-key <value>          Agent key
  --product-key <value>        Product key (default: dip)
  --category-id <value>        Optional category ID for publish
  --skip-login                 Skip kweaver auth login
  --insecure                   Pass --insecure to kweaver auth login
  --no-publish                 Create agent without publishing
  -h, --help                   Show this help

Environment variables:
  KWEAVER_IP
  KWEAVER_BASE_URL
  KWEAVER_USERNAME
  KWEAVER_PASSWORD
  KWEAVER_BUSINESS_DOMAIN
  AGENT_KEY
  AGENT_PRODUCT_KEY
  AGENT_CATEGORY_ID
  KWEAVER_SKIP_LOGIN
  KWEAVER_INSECURE
  AUTO_PUBLISH
  CONTEXTLOADER_TOOLSET_FILE

Notes:
  1. This MVP installer binds two fixed knowledge networks.
  2. The target KWeaver environment must already contain those knowledge networks.
  3. Prompt files must be customized before installation.
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

require_file() {
  [[ -f "$1" ]] || error "Required file not found: $1"
}

assert_prompt_ready() {
  local file_path="$1"
  if grep -q "TODO_REPLACE_ME" "$file_path"; then
    error "Prompt file still contains placeholder token TODO_REPLACE_ME: $file_path"
  fi
}

extract_json_value() {
  local json_input="$1"
  shift
  JSON_INPUT="$json_input" python3 - "$@" <<'PY'
import json
import os
import sys

candidates = sys.argv[1:]
text = os.environ.get("JSON_INPUT", "").strip()
if not text:
    sys.exit(1)

try:
    data = json.loads(text)
except json.JSONDecodeError:
    sys.exit(1)

def find_value(node, key):
    if isinstance(node, dict):
        if key in node and node[key] not in (None, ""):
            return node[key]
        for value in node.values():
            found = find_value(value, key)
            if found not in (None, ""):
                return found
    elif isinstance(node, list):
        for item in node:
            found = find_value(item, key)
            if found not in (None, ""):
                return found
    return None

for candidate in candidates:
    value = find_value(data, candidate)
    if value not in (None, ""):
        print(value)
        sys.exit(0)

sys.exit(1)
PY
}

find_skill_id_by_name() {
  local json_input="$1"
  local expected_name="$2"
  JSON_INPUT="$json_input" python3 - "$expected_name" <<'PY'
import json
import os
import sys

expected_name = sys.argv[1]
text = os.environ.get("JSON_INPUT", "").strip()
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
        for key in ("entries", "data", "items", "records", "skills"):
            value = node.get(key)
            if isinstance(value, list):
                return value
            if isinstance(value, dict):
                nested = iter_items(value)
                if nested:
                    return nested
    return []

for item in iter_items(data):
    if not isinstance(item, dict):
        continue
    name = item.get("name") or item.get("skill_name")
    if str(name) != expected_name:
        continue
    skill_id = item.get("id") or item.get("skill_id")
    if skill_id:
        print(skill_id)
        sys.exit(0)

sys.exit(1)
PY
}

find_kn_id_by_name() {
  local json_input="$1"
  local expected_name="$2"
  JSON_INPUT="$json_input" python3 - "$expected_name" <<'PY'
import json
import os
import sys

expected_name = sys.argv[1]
text = os.environ.get("JSON_INPUT", "").strip()
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

for item in iter_items(data):
    if not isinstance(item, dict):
        continue
    name = item.get("name") or item.get("kn_name") or item.get("knowledge_network_name")
    if str(name) != expected_name:
        continue
    kn_id = item.get("id") or item.get("kn_id") or item.get("knowledge_network_id")
    if kn_id:
        print(kn_id)
        sys.exit(0)

sys.exit(1)
PY
}

KWEAVER_BASE_URL="${KWEAVER_BASE_URL:-${KWEAVER_IP:-}}"
KWEAVER_USERNAME="${KWEAVER_USERNAME:-}"
KWEAVER_PASSWORD="${KWEAVER_PASSWORD:-}"
KWEAVER_BUSINESS_DOMAIN="${KWEAVER_BUSINESS_DOMAIN:-bd_public}"
AGENT_NAME="MySQL恢复Agent"
AGENT_PROFILE="MySQL数据库恢复智能化 Agent"
AGENT_KEY="${AGENT_KEY:-}"
AGENT_PRODUCT_KEY="${AGENT_PRODUCT_KEY:-dip}"
AGENT_CATEGORY_ID="${AGENT_CATEGORY_ID:-}"
KWEAVER_SKIP_LOGIN="${KWEAVER_SKIP_LOGIN:-0}"
KWEAVER_INSECURE="${KWEAVER_INSECURE:-0}"
AUTO_PUBLISH="${AUTO_PUBLISH:-1}"

KN_EXPERIENCE_ID="mysql_recovery_experience"
KN_EXPERIENCE_NAME="MySQL数据库恢复经验知识网络"
KN_RUN_ID="mysql_recovery_run"
KN_RUN_NAME="MySQL数据库恢复运行知识网络"
CONTEXTLOADER_TOOLSET_FILE="${CONTEXTLOADER_TOOLSET_FILE:-$ROOT_DIR/../../context-loader/toolset/context_loader_toolset.adp}"

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
    --agent-key)
      AGENT_KEY="${2:-}"
      shift 2
      ;;
    --product-key)
      AGENT_PRODUCT_KEY="${2:-}"
      shift 2
      ;;
    --category-id)
      AGENT_CATEGORY_ID="${2:-}"
      shift 2
      ;;
    --contextloader-toolset-file)
      CONTEXTLOADER_TOOLSET_FILE="${2:-}"
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
    --no-publish)
      AUTO_PUBLISH="0"
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
require_file "$SYSTEM_PROMPT_FILE"
require_file "$DOLPHIN_FILE"
require_file "$CONTEXTLOADER_TOOLSET_FILE"
assert_prompt_ready "$SYSTEM_PROMPT_FILE"
assert_prompt_ready "$DOLPHIN_FILE"

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

log "Checking fixed knowledge networks..."
kn_list_output="$(kweaver bkn list -bd "$KWEAVER_BUSINESS_DOMAIN")"
if ! experience_kn_found_id="$(find_kn_id_by_name "$kn_list_output" "$KN_EXPERIENCE_NAME")"; then
  error "Required knowledge network not found by exact name: $KN_EXPERIENCE_NAME"
fi
if [[ "$experience_kn_found_id" != "$KN_EXPERIENCE_ID" ]]; then
  error "Knowledge network ID mismatch for $KN_EXPERIENCE_NAME: expected $KN_EXPERIENCE_ID, got $experience_kn_found_id"
fi

if ! run_kn_found_id="$(find_kn_id_by_name "$kn_list_output" "$KN_RUN_NAME")"; then
  error "Required knowledge network not found by exact name: $KN_RUN_NAME"
fi
if [[ "$run_kn_found_id" != "$KN_RUN_ID" ]]; then
  error "Knowledge network ID mismatch for $KN_RUN_NAME: expected $KN_RUN_ID, got $run_kn_found_id"
fi

log "Experience KN: $KN_EXPERIENCE_NAME ($KN_EXPERIENCE_ID)"
log "Run KN: $KN_RUN_NAME ($KN_RUN_ID)"

log "Resolving ContextLoader toolbox tools..."
contextloader_tool_count="$(python3 - "$CONTEXTLOADER_TOOLSET_FILE" <<'PY'
import json
import sys
from pathlib import Path

path = Path(sys.argv[1])
data = json.loads(path.read_text(encoding="utf-8"))
names = {
    "kn_search",
    "query_object_instance",
    "query_instance_subgraph",
    "get_logic_properties_values",
    "get_action_info",
}
count = 0
for cfg in data.get("toolbox", {}).get("configs", []):
    for tool in cfg.get("tools", []):
        if tool.get("name") in names:
            count += 1
print(count)
PY
)"
if [[ "$contextloader_tool_count" -lt 5 ]]; then
  error "ContextLoader toolset is incomplete: expected at least 5 required tools, got $contextloader_tool_count"
fi
log "ContextLoader required tools: $contextloader_tool_count"

tmp_config="$(mktemp)"
cleanup() {
  rm -f "$tmp_config"
}
trap cleanup EXIT

python3 - "$SYSTEM_PROMPT_FILE" "$DOLPHIN_FILE" "$KN_EXPERIENCE_ID" "$KN_EXPERIENCE_NAME" "$KN_RUN_ID" "$KN_RUN_NAME" "$CONTEXTLOADER_TOOLSET_FILE" "$tmp_config" <<'PY'
import json
import sys
from pathlib import Path

system_prompt_path = Path(sys.argv[1])
dolphin_path = Path(sys.argv[2])
experience_id = sys.argv[3]
experience_name = sys.argv[4]
run_id = sys.argv[5]
run_name = sys.argv[6]
toolset_path = Path(sys.argv[7])
output_path = Path(sys.argv[8])

required_tool_order = [
    "kn_search",
    "query_object_instance",
    "query_instance_subgraph",
    "get_logic_properties_values",
    "get_action_info",
]

def load_contextloader_tools(path: Path):
    data = json.loads(path.read_text(encoding="utf-8"))
    by_name = {}
    for cfg in data.get("toolbox", {}).get("configs", []):
        box_id = cfg.get("box_id")
        for tool in cfg.get("tools", []):
            name = tool.get("name")
            if name in required_tool_order and box_id and tool.get("tool_id"):
                by_name[name] = {
                    "tool_id": tool["tool_id"],
                    "tool_box_id": box_id,
                    "tool_timeout": 300,
                    "tool_input": [],
                    "intervention": False,
                    "intervention_confirmation_message": "",
                    "result_process_strategies": None,
                }
    missing = [name for name in required_tool_order if name not in by_name]
    if missing:
        raise SystemExit(f"missing ContextLoader tools: {', '.join(missing)}")
    return [by_name[name] for name in required_tool_order]

config = {
    "input": {
        "fields": [
            {
                "name": "user_input",
                "type": "string",
                "desc": "",
            }
        ]
    },
    "output": {
        "default_format": "markdown",
    },
    "system_prompt": system_prompt_path.read_text(encoding="utf-8"),
    "dolphin": dolphin_path.read_text(encoding="utf-8"),
    "is_dolphin_mode": 1,
    "data_source": {
        "knowledge_network": [
            {
                "knowledge_network_id": experience_id,
                "knowledge_network_name": experience_name,
            },
            {
                "knowledge_network_id": run_id,
                "knowledge_network_name": run_name,
            }
        ]
    },
    "skills": {
        "skills": [],
        "tools": load_contextloader_tools(toolset_path),
        "mcps": [],
        "agents": [],
    },
    "plan_mode": {
        "is_enabled": True,
    },
}

output_path.write_text(
    json.dumps(config, ensure_ascii=False, indent=2) + "\n",
    encoding="utf-8",
)
PY

log "Creating decision agent..."
create_args=(
  agent create
  --name "$AGENT_NAME"
  --profile "$AGENT_PROFILE"
  --product-key "$AGENT_PRODUCT_KEY"
  --config "$tmp_config"
  -bd "$KWEAVER_BUSINESS_DOMAIN"
)
if [[ -n "$AGENT_KEY" ]]; then
  create_args+=(--key "$AGENT_KEY")
fi

create_output="$(kweaver "${create_args[@]}")"
printf '%s\n' "$create_output"

if ! AGENT_ID="$(extract_json_value "$create_output" agent_id id)"; then
  error "Failed to parse agent ID from kweaver agent create output."
fi

log "Created agent ID: $AGENT_ID"

if [[ "$AUTO_PUBLISH" == "1" ]]; then
  log "Publishing decision agent..."
  publish_args=(agent publish "$AGENT_ID")
  if [[ -n "$AGENT_CATEGORY_ID" ]]; then
    publish_args+=(--category-id "$AGENT_CATEGORY_ID")
  fi
  publish_output="$(kweaver "${publish_args[@]}")"
  printf '%s\n' "$publish_output"
fi

printf '\n'
printf 'Install completed.\n'
printf 'Agent ID: %s\n' "$AGENT_ID"
printf 'Knowledge Networks:\n'
printf '  - %s (%s)\n' "$KN_EXPERIENCE_NAME" "$KN_EXPERIENCE_ID"
printf '  - %s (%s)\n' "$KN_RUN_NAME" "$KN_RUN_ID"
printf 'ContextLoader tools:\n'
printf '  - %s\n' "$CONTEXTLOADER_TOOLSET_FILE"
