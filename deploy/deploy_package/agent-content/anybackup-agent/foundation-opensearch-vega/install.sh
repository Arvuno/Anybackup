#!/usr/bin/env bash
set -euo pipefail

CATALOG_NAME="${FOUNDATION_OPENSEARCH_CATALOG_NAME:-FoundationOpenSearch}"
HOST="${FOUNDATION_OPENSEARCH_PROXY_HOST:-v9-infra-foundation-opensearch-proxy.v9-system.svc.cluster.local}"
PORT="${FOUNDATION_OPENSEARCH_PROXY_PORT:-9896}"
USERNAME="${FOUNDATION_OPENSEARCH_USERNAME:-admin}"
PASSWORD="${FOUNDATION_OPENSEARCH_PASSWORD:-}"
INDEX_PATTERN="${FOUNDATION_OPENSEARCH_INDEX_PATTERN:-*}"
KWEAVER_BUSINESS_DOMAIN="${KWEAVER_BUSINESS_DOMAIN:-${VEGA_BUSINESS_DOMAIN:-bd_public}}"
DISCOVER_CATALOG="${FOUNDATION_OPENSEARCH_DISCOVER:-1}"
TMP_DIR="${FOUNDATION_OPENSEARCH_TMP_DIR:-/tmp/foundation-opensearch-vega}"

usage() {
  cat <<'USAGE'
Usage:
  install.sh --password <password> [options]

Options:
  --catalog-name NAME   Vega catalog name. Default: FoundationOpenSearch
  --host HOST           HTTP proxy host. Default: v9-infra-foundation-opensearch-proxy.v9-system.svc.cluster.local
  --port PORT           HTTP proxy port. Default: 9896
  --username USER       OpenSearch username. Default: admin
  --password PASSWORD   OpenSearch password. Required.
  --index-pattern TEXT  OpenSearch index pattern. Default: *
  --biz-domain ID       KWeaver business domain. Default: bd_public
  --no-discover         Skip Vega discover after create/update.
  -h, --help            Show this help.

This installer creates or updates a KWeaver Vega opensearch catalog that points
to the internal Foundation OpenSearch HTTP proxy. The proxy performs only
HTTP-to-HTTPS conversion; credentials stay in KWeaver's connector config.
USAGE
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

parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --catalog-name)
        [[ $# -ge 2 ]] || error "--catalog-name requires a value"
        CATALOG_NAME="$2"
        shift 2
        ;;
      --host)
        [[ $# -ge 2 ]] || error "--host requires a value"
        HOST="$2"
        shift 2
        ;;
      --port)
        [[ $# -ge 2 ]] || error "--port requires a value"
        PORT="$2"
        shift 2
        ;;
      --username|--user)
        [[ $# -ge 2 ]] || error "$1 requires a value"
        USERNAME="$2"
        shift 2
        ;;
      --password)
        [[ $# -ge 2 ]] || error "--password requires a value"
        PASSWORD="$2"
        shift 2
        ;;
      --index-pattern)
        [[ $# -ge 2 ]] || error "--index-pattern requires a value"
        INDEX_PATTERN="$2"
        shift 2
        ;;
      --biz-domain)
        [[ $# -ge 2 ]] || error "--biz-domain requires a value"
        KWEAVER_BUSINESS_DOMAIN="$2"
        shift 2
        ;;
      --no-discover)
        DISCOVER_CATALOG="0"
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
}

assert_inputs() {
  [[ -n "$CATALOG_NAME" ]] || error "Missing catalog name"
  [[ -n "$HOST" ]] || error "Missing OpenSearch proxy host"
  [[ "$PORT" =~ ^[0-9]+$ ]] || error "OpenSearch proxy port must be numeric"
  [[ -n "$USERNAME" ]] || error "Missing OpenSearch username"
  [[ -n "$PASSWORD" ]] || error "Missing OpenSearch password"
  [[ -n "$KWEAVER_BUSINESS_DOMAIN" ]] || error "Missing KWeaver business domain"

  require_command python3
  require_command kweaver
}

build_connector_config() {
  FOUNDATION_OPENSEARCH_HOST="$HOST" \
  FOUNDATION_OPENSEARCH_PORT="$PORT" \
  FOUNDATION_OPENSEARCH_USERNAME="$USERNAME" \
  FOUNDATION_OPENSEARCH_PASSWORD="$PASSWORD" \
  FOUNDATION_OPENSEARCH_INDEX_PATTERN="$INDEX_PATTERN" \
  python3 <<'PY'
import json
import os

config = {
    "host": os.environ["FOUNDATION_OPENSEARCH_HOST"],
    "port": int(os.environ["FOUNDATION_OPENSEARCH_PORT"]),
    "username": os.environ["FOUNDATION_OPENSEARCH_USERNAME"],
    "password": os.environ["FOUNDATION_OPENSEARCH_PASSWORD"],
    "index_pattern": os.environ["FOUNDATION_OPENSEARCH_INDEX_PATTERN"],
}

print(json.dumps(config, ensure_ascii=False, separators=(",", ":")))
PY
}

find_catalog_id_by_name() {
  local catalog_name="$1"
  local list_file="$2"

  FOUNDATION_OPENSEARCH_CATALOG_NAME="$catalog_name" \
  FOUNDATION_OPENSEARCH_CATALOG_LIST_FILE="$list_file" \
  python3 <<'PY'
import json
import os

catalog_name = os.environ["FOUNDATION_OPENSEARCH_CATALOG_NAME"]
list_file = os.environ["FOUNDATION_OPENSEARCH_CATALOG_LIST_FILE"]

with open(list_file, encoding="utf-8") as fp:
    data = json.load(fp)

if isinstance(data, list):
    entries = data
elif isinstance(data, dict):
    entries = (
        data.get("entries")
        or data.get("data")
        or data.get("items")
        or data.get("catalogs")
        or []
    )
else:
    entries = []

for item in entries:
    if not isinstance(item, dict):
        continue
    if item.get("name") != catalog_name:
        continue
    catalog_id = item.get("id") or item.get("catalog_id")
    if catalog_id:
        print(catalog_id)
    break
PY
}

extract_catalog_id() {
  local output_file="$1"

  FOUNDATION_OPENSEARCH_OUTPUT_FILE="$output_file" python3 <<'PY'
import json
import os

with open(os.environ["FOUNDATION_OPENSEARCH_OUTPUT_FILE"], encoding="utf-8") as fp:
    data = json.load(fp)

def find_catalog_id(node):
    if isinstance(node, dict):
        for key in ("id", "catalog_id"):
            value = node.get(key)
            if value not in (None, ""):
                return value
        for value in node.values():
            found = find_catalog_id(value)
            if found not in (None, ""):
                return found
    elif isinstance(node, list):
        for item in node:
            found = find_catalog_id(item)
            if found not in (None, ""):
                return found
    return None

catalog_id = find_catalog_id(data)
if catalog_id:
    print(catalog_id)
PY
}

create_or_update_catalog() {
  local connector_config="$1"
  local list_file="$TMP_DIR/catalog-list.json"
  local write_file="$TMP_DIR/catalog-write.json"
  local catalog_id

  export KWEAVER_BUSINESS_DOMAIN

  kweaver vega catalog list --limit 500 --pretty >"$list_file"
  catalog_id="$(find_catalog_id_by_name "$CATALOG_NAME" "$list_file")"

  if [[ -n "$catalog_id" ]]; then
    log "Updating Vega OpenSearch catalog: $CATALOG_NAME ($catalog_id)"
    kweaver vega catalog update "$catalog_id" \
      --name "$CATALOG_NAME" \
      --connector-type opensearch \
      --connector-config "$connector_config" \
      --tags foundation,opensearch,v9-infra \
      --description "Foundation OpenSearch through v9-infra HTTP proxy" \
      >"$write_file"
  else
    log "Creating Vega OpenSearch catalog: $CATALOG_NAME"
    kweaver vega catalog create \
      --name "$CATALOG_NAME" \
      --connector-type opensearch \
      --connector-config "$connector_config" \
      --tags foundation,opensearch,v9-infra \
      --description "Foundation OpenSearch through v9-infra HTTP proxy" \
      >"$write_file"
    catalog_id="$(extract_catalog_id "$write_file")"
  fi

  [[ -n "$catalog_id" ]] || error "Unable to resolve Vega catalog ID"
  printf '%s\n' "$catalog_id" >"$TMP_DIR/catalog-id.txt"
}

verify_catalog() {
  local catalog_id
  catalog_id="$(cat "$TMP_DIR/catalog-id.txt")"

  export KWEAVER_BUSINESS_DOMAIN

  log "Testing Vega OpenSearch catalog connection: $CATALOG_NAME ($catalog_id)"
  kweaver vega catalog test-connection "$catalog_id" >"$TMP_DIR/catalog-test.json"

  if [[ "$DISCOVER_CATALOG" != "0" && "$DISCOVER_CATALOG" != "false" ]]; then
    log "Discovering Vega OpenSearch catalog resources: $CATALOG_NAME ($catalog_id)"
    kweaver vega catalog discover "$catalog_id" --wait >"$TMP_DIR/catalog-discover.json"
  else
    log "Skipping Vega OpenSearch catalog discover: $CATALOG_NAME ($catalog_id)"
  fi

  log "Foundation OpenSearch Vega catalog is ready: $CATALOG_NAME ($catalog_id)"
}

main() {
  parse_args "$@"
  assert_inputs
  mkdir -p "$TMP_DIR"

  log "Business domain: $KWEAVER_BUSINESS_DOMAIN"
  log "OpenSearch proxy endpoint: http://$HOST:$PORT"
  log "Catalog: $CATALOG_NAME"

  local connector_config
  connector_config="$(build_connector_config)"
  create_or_update_catalog "$connector_config"
  verify_catalog
}

main "$@"
