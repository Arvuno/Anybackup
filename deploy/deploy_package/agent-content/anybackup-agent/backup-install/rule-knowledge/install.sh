#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PACKAGE_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
DEFAULT_SQL_FILE="$PACKAGE_ROOT/backup/backup_rule_knowledge_network/data/backup_rule_knowledge_pg.sql"

K8S_NAMESPACE="${K8S_NAMESPACE:-v9-system}"
POSTGRES_POD="${POSTGRES_POD:-v9-infra-postgres-0}"
POSTGRES_USERNAME="${POSTGRES_USERNAME:-v9}"
DATABASE_NAME="${DATABASE_NAME:-backup_rule_knowledge}"
SQL_FILE="${SQL_FILE:-$DEFAULT_SQL_FILE}"
EXPECTED_SCHEMA="${EXPECTED_SCHEMA:-backup_rule_knowledge}"
MIN_TABLE_COUNT="${MIN_TABLE_COUNT:-13}"

usage() {
  cat <<'EOF'
Usage:
  ./install.sh [options]

Options:
  --namespace <value>      Kubernetes namespace of the PostgreSQL pod.
  --pod <value>            PostgreSQL pod name.
  --username <value>       PostgreSQL username.
  --database <value>       Target database name.
  --sql-file <value>       SQL file to import.
  --schema <value>         Schema expected after import.
  --min-table-count <num>  Minimum expected table count in the schema.
  -h, --help               Show this help.

Environment variables mirror the option names:
  K8S_NAMESPACE, POSTGRES_POD, POSTGRES_USERNAME, DATABASE_NAME,
  SQL_FILE, EXPECTED_SCHEMA, MIN_TABLE_COUNT
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

require_file() {
  [[ -f "$1" ]] || error "Required file not found: $1"
}

validate_identifier() {
  local label="$1"
  local value="$2"
  [[ "$value" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]] || error "$label must be a PostgreSQL identifier: $value"
}

psql_postgres() {
  kubectl exec -n "$K8S_NAMESPACE" "$POSTGRES_POD" -- \
    psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USERNAME" -d postgres "$@"
}

psql_target() {
  kubectl exec -n "$K8S_NAMESPACE" "$POSTGRES_POD" -- \
    psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USERNAME" -d "$DATABASE_NAME" "$@"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --namespace)
      K8S_NAMESPACE="${2:-}"
      shift 2
      ;;
    --pod)
      POSTGRES_POD="${2:-}"
      shift 2
      ;;
    --username)
      POSTGRES_USERNAME="${2:-}"
      shift 2
      ;;
    --database)
      DATABASE_NAME="${2:-}"
      shift 2
      ;;
    --sql-file)
      SQL_FILE="${2:-}"
      shift 2
      ;;
    --schema)
      EXPECTED_SCHEMA="${2:-}"
      shift 2
      ;;
    --min-table-count)
      MIN_TABLE_COUNT="${2:-}"
      shift 2
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

require_command kubectl
require_file "$SQL_FILE"
validate_identifier "database" "$DATABASE_NAME"
validate_identifier "schema" "$EXPECTED_SCHEMA"
[[ "$MIN_TABLE_COUNT" =~ ^[0-9]+$ ]] || error "min table count must be numeric: $MIN_TABLE_COUNT"

log "Checking PostgreSQL pod: $K8S_NAMESPACE/$POSTGRES_POD"
kubectl get pod -n "$K8S_NAMESPACE" "$POSTGRES_POD" >/dev/null

db_exists="$(psql_postgres -Atc "select 1 from pg_database where datname = '$DATABASE_NAME';")"
if [[ "$db_exists" != "1" ]]; then
  log "Creating database: $DATABASE_NAME"
  kubectl exec -n "$K8S_NAMESPACE" "$POSTGRES_POD" -- \
    createdb -U "$POSTGRES_USERNAME" "$DATABASE_NAME"
else
  log "Database already exists: $DATABASE_NAME"
fi

log "Importing SQL: $SQL_FILE"
kubectl exec -i -n "$K8S_NAMESPACE" "$POSTGRES_POD" -- \
  psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USERNAME" -d "$DATABASE_NAME" <"$SQL_FILE"

table_count="$(psql_target -Atc "select count(*) from information_schema.tables where table_schema = '$EXPECTED_SCHEMA';")"
if (( table_count < MIN_TABLE_COUNT )); then
  error "Expected at least $MIN_TABLE_COUNT tables in schema $EXPECTED_SCHEMA, got $table_count"
fi

printf '\n'
printf 'Backup rule knowledge data import completed.\n'
printf 'Database: %s\n' "$DATABASE_NAME"
printf 'Schema: %s\n' "$EXPECTED_SCHEMA"
printf 'Tables: %s\n' "$table_count"
