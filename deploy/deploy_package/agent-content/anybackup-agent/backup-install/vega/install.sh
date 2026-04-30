#!/usr/bin/env bash
set -euo pipefail

DB_ENGINE="${VEGA_DB_ENGINE:-postgresql}"
DB_HOST="${VEGA_DB_HOST:-}"
DB_PORT="${VEGA_DB_PORT:-5432}"
DB_NAME="${VEGA_DB_NAME:-backup_rule_knowledge}"
DB_SCHEMA="${VEGA_DB_SCHEMA:-backup_rule_knowledge}"
DB_USER="${VEGA_DB_USER:-}"
DB_PASSWORD="${VEGA_DB_PASSWORD:-}"
CATALOG_NAME="${VEGA_CATALOG_NAME:-backup_rule_knowledge}"
KWEAVER_BUSINESS_DOMAIN="${KWEAVER_BUSINESS_DOMAIN:-${VEGA_BUSINESS_DOMAIN:-bd_public}}"
DISCOVER_CATALOG="${VEGA_DISCOVER_CATALOG:-1}"

usage() {
  cat <<'EOF'
Usage:
  ./install.sh --host <host> --port <port> --database <database> --schema <schema> --username <username> --password <password> [options]

Options:
  --engine <value>        Database engine. Supported: postgresql. Default: postgresql.
  --host <value>          Database host or host:port.
  --port <value>          Database port. Default: 5432.
  --database <value>      Database name. Default: backup_rule_knowledge.
  --schema <value>        PostgreSQL schema. Default: backup_rule_knowledge.
  --username <value>      Database user.
  --password <value>      Database password.
  --catalog-name <value>  Vega catalog name. Default: backup_rule_knowledge.
  --biz-domain <value>    KWeaver business domain. Default: bd_public.
  --no-discover           Create or update catalog without running discover.
  -h, --help              Show this help.

Notes:
  1. This installer does not create or import the PostgreSQL database. Run
     backup-install/rule-knowledge/install.sh first.
  2. Re-running is idempotent: existing catalog names are updated, otherwise a
     new catalog is created.
  3. BKN data_view IDs use logical table names, so Vega discover must expose
     data_view names matching the PostgreSQL table names.
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

parse_host_port() {
  if [[ "$DB_HOST" == *:* && "$DB_HOST" != *"]"* ]]; then
    DB_PORT="${DB_HOST##*:}"
    DB_HOST="${DB_HOST%:*}"
  fi
}

normalize_engine() {
  case "$DB_ENGINE" in
    postgres|postgresql|pg)
      DB_ENGINE="postgresql"
      ;;
    *)
      error "Unsupported engine for backup rule knowledge Vega catalog: $DB_ENGINE"
      ;;
  esac
}

validate_port() {
  [[ "$DB_PORT" =~ ^[0-9]+$ ]] || error "Database port must be numeric"
}

parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --engine)
        [[ $# -ge 2 ]] || error "--engine requires a value"
        DB_ENGINE="$2"
        shift 2
        ;;
      --host)
        [[ $# -ge 2 ]] || error "--host requires a value"
        DB_HOST="$2"
        shift 2
        ;;
      --port)
        [[ $# -ge 2 ]] || error "--port requires a value"
        DB_PORT="$2"
        shift 2
        ;;
      --database|--db-name)
        [[ $# -ge 2 ]] || error "$1 requires a value"
        DB_NAME="$2"
        shift 2
        ;;
      --schema)
        [[ $# -ge 2 ]] || error "--schema requires a value"
        DB_SCHEMA="$2"
        shift 2
        ;;
      --username|--user)
        [[ $# -ge 2 ]] || error "$1 requires a value"
        DB_USER="$2"
        shift 2
        ;;
      --password)
        [[ $# -ge 2 ]] || error "--password requires a value"
        DB_PASSWORD="$2"
        shift 2
        ;;
      --catalog-name)
        [[ $# -ge 2 ]] || error "--catalog-name requires a value"
        CATALOG_NAME="$2"
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
  [[ -n "$DB_HOST" ]] || error "Missing database host"
  [[ -n "$DB_NAME" ]] || error "Missing database name"
  [[ -n "$DB_SCHEMA" ]] || error "Missing database schema"
  [[ -n "$DB_USER" ]] || error "Missing database username"
  [[ -n "$DB_PASSWORD" ]] || error "Missing database password"
  [[ -n "$CATALOG_NAME" ]] || error "Missing Vega catalog name"

  parse_host_port
  normalize_engine
  validate_port
  require_command python3
  require_command kweaver
}

build_vega_connector_config() {
  VEGA_DB_HOST="$DB_HOST" \
  VEGA_DB_PORT="$DB_PORT" \
  VEGA_DB_NAME="$DB_NAME" \
  VEGA_DB_SCHEMA="$DB_SCHEMA" \
  VEGA_DB_USER="$DB_USER" \
  VEGA_DB_PASSWORD="$DB_PASSWORD" \
  python3 <<'PY'
import json
import os

host = os.environ["VEGA_DB_HOST"]
port = int(os.environ["VEGA_DB_PORT"])
database = os.environ["VEGA_DB_NAME"]
schema = os.environ["VEGA_DB_SCHEMA"]

config = {
    "host": host,
    "port": port,
    "username": os.environ["VEGA_DB_USER"],
    "password": os.environ["VEGA_DB_PASSWORD"],
    "databases": [database],
    "database": database,
    "database_name": database,
    "schema": schema,
    "schema_name": schema,
    "connect_protocol": "jdbc",
    "jdbc_url": f"jdbc:postgresql://{host}:{port}/{database}",
}

print(json.dumps(config, ensure_ascii=False, separators=(",", ":")))
PY
}

find_vega_catalog_ids_by_name() {
  local catalog_name="$1"
  local output

  output="$(kweaver vega catalog list --limit 100 -bd "$KWEAVER_BUSINESS_DOMAIN" --pretty)"

  VEGA_CATALOG_LIST_JSON="$output" \
  VEGA_CATALOG_NAME="$catalog_name" \
  python3 <<'PY'
import json
import os

text = os.environ.get("VEGA_CATALOG_LIST_JSON", "").strip()
target_name = os.environ["VEGA_CATALOG_NAME"]
if not text:
    raise SystemExit(0)

data = json.loads(text)
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
    if item.get("name") != target_name:
        continue
    catalog_id = item.get("id") or item.get("catalog_id")
    if catalog_id:
        print(catalog_id)
PY
}

extract_vega_catalog_id() {
  VEGA_CATALOG_JSON="$1" python3 <<'PY'
import json
import os

text = os.environ.get("VEGA_CATALOG_JSON", "").strip()
if not text:
    raise SystemExit(1)

data = json.loads(text)

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
if catalog_id in (None, ""):
    raise SystemExit(1)

print(catalog_id)
PY
}

discover_vega_catalog() {
  local catalog_id="$1"

  if [[ "$DISCOVER_CATALOG" == "0" || "$DISCOVER_CATALOG" == "false" ]]; then
    log "Skipping Vega catalog discover: $CATALOG_NAME ($catalog_id)"
    return 0
  fi

  log "Discovering Vega catalog: $CATALOG_NAME ($catalog_id)"
  kweaver vega catalog discover "$catalog_id" --wait -bd "$KWEAVER_BUSINESS_DOMAIN" --pretty
}

create_or_update_vega_catalog() {
  local connector_config
  connector_config="$(build_vega_connector_config)"

  local existing_id
  existing_id="$(find_vega_catalog_ids_by_name "$CATALOG_NAME" | sed -n '1p')"

  if [[ -n "$existing_id" ]]; then
    log "Updating existing Vega catalog: $CATALOG_NAME ($existing_id)"
    kweaver vega catalog update "$existing_id" \
      --name "$CATALOG_NAME" \
      --connector-type postgresql \
      --connector-config "$connector_config" \
      --description "Backup rule knowledge PostgreSQL data connection" \
      -bd "$KWEAVER_BUSINESS_DOMAIN" \
      --pretty
    discover_vega_catalog "$existing_id"
    return 0
  fi

  log "Creating Vega catalog: $CATALOG_NAME"
  local create_output
  create_output="$(kweaver vega catalog create \
    --name "$CATALOG_NAME" \
    --connector-type postgresql \
    --connector-config "$connector_config" \
    --description "Backup rule knowledge PostgreSQL data connection" \
    -bd "$KWEAVER_BUSINESS_DOMAIN" \
    --pretty)"
  printf '%s\n' "$create_output"

  local created_id
  created_id="$(extract_vega_catalog_id "$create_output" || true)"
  if [[ -z "$created_id" ]]; then
    created_id="$(find_vega_catalog_ids_by_name "$CATALOG_NAME" | sed -n '1p')"
  fi
  [[ -n "$created_id" ]] || error "Vega catalog was created, but its ID could not be resolved for discover."

  discover_vega_catalog "$created_id"
}

main() {
  parse_args "$@"
  assert_inputs

  log "Installing backup rule knowledge Vega catalog"
  log "Host: $DB_HOST"
  log "Port: $DB_PORT"
  log "Database: $DB_NAME"
  log "Schema: $DB_SCHEMA"
  log "Catalog: $CATALOG_NAME"
  log "Business domain: $KWEAVER_BUSINESS_DOMAIN"

  create_or_update_vega_catalog

  log "Backup rule knowledge Vega catalog install completed."
}

main "$@"
