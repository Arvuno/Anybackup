#!/usr/bin/env bash
set -euo pipefail

DB_ENGINE="${VEGA_DB_ENGINE:-postgresql}"
DB_HOST="${VEGA_DB_HOST:-}"
DB_PORT="${VEGA_DB_PORT:-}"
DB_NAME="${VEGA_DB_NAME:-ExperienceBKNDB}"
DB_USER="${VEGA_DB_USER:-}"
DB_PASSWORD="${VEGA_DB_PASSWORD:-}"
KWEAVER_BUSINESS_DOMAIN="${KWEAVER_BUSINESS_DOMAIN:-${VEGA_BUSINESS_DOMAIN:-bd_public}}"
VEGA_EXPERIENCE_CATALOG_NAME="${VEGA_EXPERIENCE_CATALOG_NAME:-${VEGA_CATALOG_NAME:-恢复经验知识网络数据连接}}"
VEGA_COMMON_CATALOG_NAME="${VEGA_COMMON_CATALOG_NAME:-CommonServiceDB-client}"
VEGA_MULTI_STORAGE_CATALOG_NAME="${VEGA_MULTI_STORAGE_CATALOG_NAME:-MultiStorageSvcMgmServiceDB-storageservice}"
VEGA_STORAGE_RES_CATALOG_NAME="${VEGA_STORAGE_RES_CATALOG_NAME:-StorageResMgmServiceDB-poolv8}"
VEGA_HYPER_BACKUP_CATALOG_NAME="${VEGA_HYPER_BACKUP_CATALOG_NAME:-HyperBackupMgmServiceDB-protectobject}"
VEGA_HYPER_JOB_CATALOG_NAME="${VEGA_HYPER_JOB_CATALOG_NAME:-HyperJobWorkerServiceDB-job}"
SKIP_VEGA_CATALOG="${VEGA_SKIP_CATALOG:-0}"

TMP_DIR=""
MYSQL_DEFAULTS_FILE=""

usage() {
  cat <<'EOF'
Usage:
  ./uninstall.sh <host> <port> <database> <username> <password> [--engine postgresql|mysql|mariadb]
  ./uninstall.sh <host> <port> <username> <password> [--engine postgresql|mysql|mariadb]
  ./uninstall.sh --host <host> --port <port> --username <username> --password <password> [--database ExperienceBKNDB] [--engine postgresql|mysql|mariadb]

Options:
  --engine <value>      Database engine. Supported: postgresql, mysql, mariadb. Default: postgresql.
  --host <value>        Database host or host:port.
  --port <value>        Database port. Default: 5432 for postgresql, 3306 for mysql/mariadb.
  --database <value>    Database name to remove. Default: ExperienceBKNDB.
  --username <value>    Database user.
  --password <value>    Database password.
  --biz-domain <value>  KWeaver business domain for Vega catalog deletion. Default: bd_public.
  --skip-vega-catalog   Skip deleting the Vega data connection.
  -h, --help            Show this help.

Environment variables:
  VEGA_DB_ENGINE
  VEGA_DB_HOST
  VEGA_DB_PORT
  VEGA_DB_NAME
  VEGA_DB_USER
  VEGA_DB_PASSWORD
  VEGA_EXPERIENCE_CATALOG_NAME
  VEGA_COMMON_CATALOG_NAME
  VEGA_MULTI_STORAGE_CATALOG_NAME
  VEGA_STORAGE_RES_CATALOG_NAME
  VEGA_HYPER_BACKUP_CATALOG_NAME
  VEGA_HYPER_JOB_CATALOG_NAME
  VEGA_BUSINESS_DOMAIN
  KWEAVER_BUSINESS_DOMAIN
  VEGA_SKIP_CATALOG

Notes:
  1. The core inputs are database host/IP, port, database name, username, and password.
  2. Database name defaults to ExperienceBKNDB when omitted.
  3. Database engine defaults to postgresql. Use --engine mysql or --engine mariadb to switch.
  4. Database port defaults to 5432 for postgresql and 3306 for mysql/mariadb when omitted.
  5. The uninstaller deletes the Vega data connections named 恢复经验知识网络数据连接, CommonServiceDB-client, MultiStorageSvcMgmServiceDB-storageservice, StorageResMgmServiceDB-poolv8, HyperBackupMgmServiceDB-protectobject, and HyperJobWorkerServiceDB-job.
  6. It then removes the five recovery experience tables and drops the target database.
EOF
}

log() {
  printf '[INFO] %s\n' "$*"
}

error() {
  printf '[ERROR] %s\n' "$*" >&2
  exit 1
}

cleanup() {
  if [[ -n "$TMP_DIR" && -d "$TMP_DIR" ]]; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT

require_command() {
  command -v "$1" >/dev/null 2>&1 || error "Required command not found: $1"
}

validate_database_name() {
  [[ "$DB_NAME" =~ ^[A-Za-z_][A-Za-z0-9_]*$ ]] || error "Database name must match ^[A-Za-z_][A-Za-z0-9_]*$"
}

validate_port() {
  [[ "$DB_PORT" =~ ^[0-9]+$ ]] || error "Database port must be numeric"
}

parse_host_port() {
  if [[ -z "$DB_PORT" && "$DB_HOST" == *:* && "$DB_HOST" != *"]"* ]]; then
    DB_PORT="${DB_HOST##*:}"
    DB_HOST="${DB_HOST%:*}"
  fi
}

default_port() {
  if [[ -n "$DB_PORT" ]]; then
    return 0
  fi

  case "$DB_ENGINE" in
    postgresql)
      DB_PORT="5432"
      ;;
    mysql|mariadb)
      DB_PORT="3306"
      ;;
  esac
}

normalize_engine() {
  case "$DB_ENGINE" in
    postgres|postgresql|pg)
      DB_ENGINE="postgresql"
      ;;
    mysql)
      DB_ENGINE="mysql"
      ;;
    mariadb)
      DB_ENGINE="mariadb"
      ;;
    *)
      error "Unsupported engine: $DB_ENGINE"
      ;;
  esac
}

parse_args() {
  local positional=()

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
      --biz-domain)
        [[ $# -ge 2 ]] || error "--biz-domain requires a value"
        KWEAVER_BUSINESS_DOMAIN="$2"
        shift 2
        ;;
      --skip-vega-catalog)
        SKIP_VEGA_CATALOG="1"
        shift
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      --)
        shift
        while [[ $# -gt 0 ]]; do
          positional+=("$1")
          shift
        done
        ;;
      -*)
        error "Unknown option: $1"
        ;;
      *)
        positional+=("$1")
        shift
        ;;
    esac
  done

  if [[ ${#positional[@]} -gt 0 ]]; then
    [[ ${#positional[@]} -eq 3 || ${#positional[@]} -eq 4 || ${#positional[@]} -ge 5 ]] || error "Positional mode requires: <host> <port> <username> <password> or <host> <port> <database> <username> <password>"
    DB_HOST="${positional[0]}"
    if [[ ${#positional[@]} -eq 3 ]]; then
      DB_USER="${positional[1]}"
      DB_PASSWORD="${positional[2]}"
    elif [[ "${positional[1]}" =~ ^[0-9]+$ ]]; then
      DB_PORT="${positional[1]}"
      if [[ ${#positional[@]} -eq 4 ]]; then
        DB_USER="${positional[2]}"
        DB_PASSWORD="${positional[3]}"
      else
        DB_NAME="${positional[2]}"
        DB_USER="${positional[3]}"
        DB_PASSWORD="${positional[4]}"
      fi
      if [[ ${#positional[@]} -ge 6 ]]; then
        DB_ENGINE="${positional[5]}"
      fi
    else
      DB_NAME="${positional[1]}"
      DB_USER="${positional[2]}"
      DB_PASSWORD="${positional[3]}"
      if [[ ${#positional[@]} -ge 5 ]]; then
        DB_ENGINE="${positional[4]}"
      fi
    fi
  fi
}

assert_inputs() {
  [[ -n "$DB_HOST" ]] || error "Missing database host"
  [[ -n "$DB_NAME" ]] || error "Missing database name"
  [[ -n "$DB_USER" ]] || error "Missing database username"
  [[ -n "$DB_PASSWORD" ]] || error "Missing database password"

  validate_database_name
  parse_host_port
  normalize_engine
  default_port
  validate_port
}

run_psql() {
  local database="$1"
  shift

  if [[ -n "$DB_PORT" ]]; then
    PGPASSWORD="$DB_PASSWORD" PGCLIENTENCODING=UTF8 PGPORT="$DB_PORT" psql -h "$DB_HOST" -U "$DB_USER" -d "$database" -v ON_ERROR_STOP=1 "$@"
  else
    PGPASSWORD="$DB_PASSWORD" PGCLIENTENCODING=UTF8 psql -h "$DB_HOST" -U "$DB_USER" -d "$database" -v ON_ERROR_STOP=1 "$@"
  fi
}

postgresql_database_exists() {
  local exists
  exists="$(run_psql postgres -qAt -v db_name="$DB_NAME" -c "SELECT 1 FROM pg_database WHERE datname = :'db_name';")"
  [[ "$exists" == "1" ]]
}

uninstall_postgresql() {
  require_command psql

  if ! postgresql_database_exists; then
    log "PostgreSQL database does not exist, skipping: $DB_NAME"
    return 0
  fi

  log "Dropping PostgreSQL recovery experience tables from $DB_NAME"
  run_psql "$DB_NAME" <<'SQL'
DROP TABLE IF EXISTS
  "availability_checkpoint_template",
  "fault_pattern",
  "recovery_capability",
  "recovery_strategy_template",
  "risk_rule";
SQL

  log "Dropping PostgreSQL database: $DB_NAME"
  run_psql postgres -v db_name="$DB_NAME" <<'SQL'
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE datname = :'db_name'
  AND pid <> pg_backend_pid();

DROP DATABASE :"db_name";
SQL
}

prepare_mysql_defaults_file() {
  TMP_DIR="$(mktemp -d)"
  MYSQL_DEFAULTS_FILE="$TMP_DIR/mysql.cnf"

  {
    printf '[client]\n'
    printf 'host=%s\n' "$DB_HOST"
    if [[ -n "$DB_PORT" ]]; then
      printf 'port=%s\n' "$DB_PORT"
    fi
    printf 'user=%s\n' "$DB_USER"
    printf 'password=%s\n' "$DB_PASSWORD"
    printf 'default-character-set=utf8mb4\n'
  } >"$MYSQL_DEFAULTS_FILE"
  chmod 600 "$MYSQL_DEFAULTS_FILE"
}

mysql_exec() {
  mysql --defaults-extra-file="$MYSQL_DEFAULTS_FILE" "$@"
}

mysql_exec_db() {
  mysql_exec "$DB_NAME" "$@"
}

mysql_database_exists() {
  local exists
  exists="$(mysql_exec --batch --skip-column-names -e "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '$DB_NAME';")"
  [[ "$exists" == "$DB_NAME" ]]
}

uninstall_mysql() {
  require_command mysql
  prepare_mysql_defaults_file

  if ! mysql_database_exists; then
    log "MySQL/MariaDB database does not exist, skipping: $DB_NAME"
    return 0
  fi

  log "Dropping MySQL/MariaDB recovery experience tables from $DB_NAME"
  mysql_exec_db <<'SQL'
DROP TABLE IF EXISTS
  `availability_checkpoint_template`,
  `fault_pattern`,
  `recovery_capability`,
  `recovery_strategy_template`,
  `risk_rule`;
SQL

  log "Dropping MySQL/MariaDB database: $DB_NAME"
  mysql_exec -e "DROP DATABASE IF EXISTS \`$DB_NAME\`;"
}

find_vega_catalog_ids_by_name() {
  local catalog_name="$1"

  require_command python3

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

delete_vega_data_connection_by_name() {
  local catalog_name="$1"

  local ids
  ids="$(find_vega_catalog_ids_by_name "$catalog_name")"
  if [[ -z "$ids" ]]; then
    log "Vega data connection does not exist, skipping: $catalog_name"
    return 0
  fi

  local catalog_id
  while IFS= read -r catalog_id; do
    [[ -n "$catalog_id" ]] || continue
    log "Deleting Vega data connection: $catalog_name ($catalog_id)"
    kweaver vega catalog delete "$catalog_id" -y -bd "$KWEAVER_BUSINESS_DOMAIN"
  done <<<"$ids"
}

delete_vega_data_connections() {
  if [[ "$SKIP_VEGA_CATALOG" == "1" || "$SKIP_VEGA_CATALOG" == "true" ]]; then
    log "Skipping Vega data connection deletion."
    return 0
  fi

  require_command kweaver

  delete_vega_data_connection_by_name "$VEGA_EXPERIENCE_CATALOG_NAME"
  delete_vega_data_connection_by_name "$VEGA_COMMON_CATALOG_NAME"
  delete_vega_data_connection_by_name "$VEGA_MULTI_STORAGE_CATALOG_NAME"
  delete_vega_data_connection_by_name "$VEGA_STORAGE_RES_CATALOG_NAME"
  delete_vega_data_connection_by_name "$VEGA_HYPER_BACKUP_CATALOG_NAME"
  delete_vega_data_connection_by_name "$VEGA_HYPER_JOB_CATALOG_NAME"
}

main() {
  parse_args "$@"
  assert_inputs

  log "Uninstalling Vega recovery experience data"
  log "Engine: $DB_ENGINE"
  log "Host: $DB_HOST"
  log "Port: $DB_PORT"
  log "Database: $DB_NAME"

  delete_vega_data_connections

  case "$DB_ENGINE" in
    postgresql)
      uninstall_postgresql
      ;;
    mysql|mariadb)
      uninstall_mysql
      ;;
  esac

  log "Vega data uninstall completed."
}

main "$@"
