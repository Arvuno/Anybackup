#!/usr/bin/env bash
set -euo pipefail
shopt -s nullglob

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PACKAGE_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

EXPERIENCE_BKN_DIR="${EXPERIENCE_BKN_DIR:-$PACKAGE_ROOT/recovery-bkn/recovery_experience/database}"
RUN_BKN_DIR="${RUN_BKN_DIR:-$PACKAGE_ROOT/recovery-bkn/recovery_run/database}"

EXPERIENCE_KN_ID="mysql_recovery_experience"
EXPERIENCE_KN_NAME="MySQL数据库恢复经验知识网络"
RUN_KN_ID="mysql_recovery_run"
RUN_KN_NAME="MySQL数据库恢复运行知识网络"

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
  --index-small-model <value>  Vector index small model ID/name (default: disable vector index)
  --skip-login                 Skip kweaver auth login
  --insecure                   Pass --insecure to kweaver auth login
  -h, --help                   Show this help

Environment variables:
  KWEAVER_IP
  KWEAVER_BASE_URL
  KWEAVER_USERNAME
  KWEAVER_PASSWORD
  KWEAVER_BUSINESS_DOMAIN
  KWEAVER_INDEX_SMALL_MODEL
  KWEAVER_SKIP_LOGIN
  KWEAVER_INSECURE

Notes:
  1. This installer imports two fixed BKN directories from examples/bkn.
  2. After each import, it configures indexes for every data property.
  3. The imported networks are:
     - MySQL数据库恢复经验知识网络 (mysql_recovery_experience)
     - MySQL数据库恢复运行知识网络 (mysql_recovery_run)
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

require_directory() {
  [[ -d "$1" ]] || error "Required directory not found: $1"
}

normalize_base_url() {
  local raw="$1"
  if [[ "$raw" == http://* || "$raw" == https://* ]]; then
    printf '%s\n' "$raw"
    return 0
  fi
  printf 'https://%s\n' "$raw"
}

extract_kn_id() {
  local json_input="$1"
  JSON_INPUT="$json_input" python3 - <<'PY'
import json
import os
import sys

text = os.environ.get("JSON_INPUT", "").strip()
if not text:
    sys.exit(1)

try:
    data = json.loads(text)
except json.JSONDecodeError:
    sys.exit(1)

def find_value(node):
    if isinstance(node, dict):
        for key in ("knowledge_network_id", "kn_id", "id"):
            value = node.get(key)
            if value not in (None, ""):
                return value
        for value in node.values():
            found = find_value(value)
            if found not in (None, ""):
                return found
    elif isinstance(node, list):
        for item in node:
            found = find_value(item)
            if found not in (None, ""):
                return found
    return None

value = find_value(data)
if value in (None, ""):
    sys.exit(1)

print(value)
PY
}

build_recovery_data_view_map() {
  KWEAVER_BUSINESS_DOMAIN_VALUE="$KWEAVER_BUSINESS_DOMAIN" python3 <<'PY'
import json
import os
import subprocess
import sys

expected = {
    "ExperienceBKNDB.public.availability_checkpoint_template": [
        "availability_checkpoint_template",
    ],
    "ExperienceBKNDB.public.fault_pattern": [
        "fault_pattern",
    ],
    "ExperienceBKNDB.public.recovery_capability": [
        "recovery_capability",
    ],
    "ExperienceBKNDB.public.recovery_strategy_template": [
        "recovery_strategy_template",
    ],
    "ExperienceBKNDB.public.risk_rule": [
        "risk_rule",
    ],
    "CommonServiceDB.client": [
        "client",
    ],
    "HyperBackupMgmServiceDB.protect_object": [
        "protect_object",
    ],
    "HyperJobWorkerServiceDB.job": [
        "job",
    ],
    "StorageResMgmServiceDB.pool_v8": [
        "storagePool",
        "pool_v8",
    ],
    "MultiStorageSvcMgmServiceDB.console_storage_svc": [
        "storageService",
        "console_storage_svc",
    ],
}

business_domain = os.environ.get("KWEAVER_BUSINESS_DOMAIN_VALUE", "bd_public")
result = subprocess.run(
    [
        "kweaver",
        "dataview",
        "list",
        "--limit",
        "1000",
        "-bd",
        business_domain,
        "--pretty",
    ],
    check=True,
    text=True,
    stdout=subprocess.PIPE,
)
data = json.loads(result.stdout)
if isinstance(data, list):
    entries = data
elif isinstance(data, dict):
    entries = data.get("entries") or data.get("data") or data.get("items") or []
else:
    entries = []

views_by_name = {}
for item in entries:
    if not isinstance(item, dict):
        continue
    name = item.get("name")
    view_id = item.get("id")
    if name and view_id:
        views_by_name[name] = {"id": view_id, "name": name}

mapping = {}
missing = []
for resource_name, aliases in expected.items():
    view = views_by_name.get(resource_name)
    if view is None:
        missing.append(resource_name)
        continue
    for alias in [resource_name, *aliases]:
        mapping[alias] = view

if missing:
    print(
        "Missing required KWeaver data views for recovery BKN: " + ", ".join(missing),
        file=sys.stderr,
    )
    sys.exit(1)

print(json.dumps(mapping, ensure_ascii=False, separators=(",", ":")))
PY
}

build_indexed_object_type_body() {
  local input_path="$1"
  local output_path="$2"
  local small_model="$3"

  python3 - "$input_path" "$output_path" "$small_model" <<'PY'
import json
import os
import sys

input_path, output_path, small_model = sys.argv[1:4]
data_view_map = json.loads(os.environ.get("KWEAVER_DATA_VIEW_MAP_JSON") or "{}")

with open(input_path, "r", encoding="utf-8") as f:
    data = json.load(f)

if isinstance(data, dict) and isinstance(data.get("entries"), list) and data["entries"]:
    entry = data["entries"][0]
elif isinstance(data, dict) and isinstance(data.get("data"), list) and data["data"]:
    entry = data["data"][0]
else:
    entry = data

if not isinstance(entry, dict):
    raise SystemExit("Unexpected object type response shape.")

for key in ("status", "creator", "updater", "create_time", "update_time", "module_type", "kn_id"):
    entry.pop(key, None)

data_source = entry.get("data_source")
if isinstance(data_source, dict) and data_source.get("type") == "data_view":
    resolved = None
    for lookup_key in (data_source.get("id"), data_source.get("name")):
        if lookup_key in data_view_map:
            resolved = data_view_map[lookup_key]
            break
    if isinstance(resolved, dict):
        data_source["id"] = resolved["id"]
        data_source["name"] = resolved["name"]

properties = entry.get("data_properties")
if not isinstance(properties, list):
    raise SystemExit("Object type response does not include data_properties.")

for prop in properties:
    if not isinstance(prop, dict):
        continue

    index_config = prop.get("index_config")
    if not isinstance(index_config, dict):
        index_config = {}

    keyword_config = index_config.get("keyword_config")
    if not isinstance(keyword_config, dict):
        keyword_config = {}
    keyword_config.update({
        "enabled": True,
        "length": 512,
        "max_length": 512,
        "max_length_bytes": 512,
        "ignore_above_len": 512,
    })

    fulltext_config = index_config.get("fulltext_config")
    if not isinstance(fulltext_config, dict):
        fulltext_config = {}
    fulltext_config.update({
        "enabled": True,
        "tokenizer": "standard",
        "analyzer": "standard",
    })

    vector_config = index_config.get("vector_config")
    if not isinstance(vector_config, dict):
        vector_config = {}
    if small_model:
        vector_config.update({
            "enabled": True,
            "model_type": "small",
            "model_id": small_model,
            "model_name": small_model,
            "small_model": small_model,
            "small_model_id": small_model,
        })
    else:
        vector_config["enabled"] = False
        for key in ("model_id", "model_name", "small_model", "small_model_id"):
            vector_config.pop(key, None)

    index_config["keyword_config"] = keyword_config
    index_config["fulltext_config"] = fulltext_config
    index_config["vector_config"] = vector_config
    prop["index_config"] = index_config

with open(output_path, "w", encoding="utf-8") as f:
    json.dump(entry, f, ensure_ascii=False, separators=(",", ":"))
PY
}

configure_object_type_indexes() {
  local kn_id="$1"
  local ot_id="$2"

  local get_output
  local put_body
  get_output="$(mktemp)"
  put_body="$(mktemp)"

  if ! kweaver bkn object-type get "$kn_id" "$ot_id" -bd "$KWEAVER_BUSINESS_DOMAIN" --pretty >"$get_output"; then
    rm -f "$get_output" "$put_body"
    error "Failed to get object type for index configuration: $kn_id/$ot_id"
  fi

  if ! build_indexed_object_type_body "$get_output" "$put_body" "$KWEAVER_INDEX_SMALL_MODEL"; then
    rm -f "$get_output" "$put_body"
    error "Failed to build index configuration body: $kn_id/$ot_id"
  fi

  local body
  body="$(python3 - "$put_body" <<'PY'
import sys

with open(sys.argv[1], "r", encoding="utf-8") as f:
    print(f.read(), end="")
PY
)"

  log "Configuring indexes for object type: $kn_id/$ot_id"
  if ! kweaver bkn object-type update "$kn_id" "$ot_id" "$body" -bd "$KWEAVER_BUSINESS_DOMAIN" --pretty >/dev/null; then
    rm -f "$get_output" "$put_body"
    error "Failed to update index configuration: $kn_id/$ot_id"
  fi

  rm -f "$get_output" "$put_body"
}

configure_bkn_indexes() {
  local label="$1"
  local directory="$2"
  local kn_id="$3"

  local object_type_files=("$directory"/object_types/*.bkn)
  if [[ ${#object_type_files[@]} -eq 0 ]]; then
    log "No object types found for index configuration: $label"
    return 0
  fi

  log "Configuring indexes for $label..."
  for object_type_file in "${object_type_files[@]}"; do
    local filename
    local ot_id
    filename="$(basename "$object_type_file")"
    ot_id="${filename%.bkn}"
    configure_object_type_indexes "$kn_id" "$ot_id"
  done
}

import_bkn() {
  local label="$1"
  local directory="$2"
  local expected_id="$3"

  log "Validating $label..."
  kweaver bkn validate "$directory"

  log "Importing $label..."
  local output
  output="$(kweaver bkn push "$directory" -bd "$KWEAVER_BUSINESS_DOMAIN" --pretty)"
  printf '%s\n' "$output"

  local imported_id
  if ! imported_id="$(extract_kn_id "$output")"; then
    error "Failed to parse knowledge network ID from import output for $label"
  fi

  if [[ "$imported_id" != "$expected_id" ]]; then
    error "Imported knowledge network ID mismatch for $label: expected $expected_id, got $imported_id"
  fi

  configure_bkn_indexes "$label" "$directory" "$imported_id"
}

KWEAVER_BASE_URL="${KWEAVER_BASE_URL:-${KWEAVER_IP:-}}"
KWEAVER_USERNAME="${KWEAVER_USERNAME:-}"
KWEAVER_PASSWORD="${KWEAVER_PASSWORD:-}"
KWEAVER_BUSINESS_DOMAIN="${KWEAVER_BUSINESS_DOMAIN:-bd_public}"
KWEAVER_INDEX_SMALL_MODEL="${KWEAVER_INDEX_SMALL_MODEL:-}"
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
    --index-small-model)
      KWEAVER_INDEX_SMALL_MODEL="${2:-}"
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
require_directory "$EXPERIENCE_BKN_DIR"
require_directory "$RUN_BKN_DIR"

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

KWEAVER_DATA_VIEW_MAP_JSON="$(build_recovery_data_view_map)"
export KWEAVER_DATA_VIEW_MAP_JSON

import_bkn "$EXPERIENCE_KN_NAME" "$EXPERIENCE_BKN_DIR" "$EXPERIENCE_KN_ID"
import_bkn "$RUN_KN_NAME" "$RUN_BKN_DIR" "$RUN_KN_ID"

printf '\n'
printf 'BKN install completed.\n'
printf 'Imported knowledge networks:\n'
printf '  - %s (%s)\n' "$EXPERIENCE_KN_NAME" "$EXPERIENCE_KN_ID"
printf '  - %s (%s)\n' "$RUN_KN_NAME" "$RUN_KN_ID"
