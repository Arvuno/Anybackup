#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEFAULT_SKILLS_ROOT="$SCRIPT_DIR/skills"

usage() {
  cat <<'EOF'
Usage:
  ./install.sh --base-url <url> [options]

Options:
  --base-url <url>      KWeaver service base URL.
  --biz-domain <value>  Business domain ID (default: bd_public).
  --skills-root <dir>   Directory containing skill directories.
  --skill-dir <dir>     Skill directory to import. Can be repeated.
  --state-file <path>   Write imported skill id mapping to this JSON file.
  --source <value>      KWeaver skill source value (default: custom).
  --publish             Set imported/reused skills to published.
  --no-publish          Do not change skill publish status.
  --kweaver-config <p>  KWeaver deploy config used for publish schema checks.
  --resource-namespace <ns>
                        KWeaver resource namespace (default: resource).
  --skip-publish-schema-check
                        Do not create missing KWeaver skill publish tables.
  --insecure            Tell KWeaver CLI to skip TLS verification where supported.
  -h, --help            Show this help.

Environment variables:
  KWEAVER_BASE_URL
  KWEAVER_BUSINESS_DOMAIN
  ANYBACKUP_SKILLS_ROOT
  ANYBACKUP_SKILL_DIRS       Comma-separated skill directories.
  ANYBACKUP_SKILLS_STATE_FILE
  ANYBACKUP_SKILLS_SOURCE    Default: custom.
  ANYBACKUP_SKILLS_PUBLISH   1 or 0. Default: 1.
  ANYBACKUP_SKILLS_ENSURE_PUBLISH_SCHEMA
                             1 or 0. Default: 1.
  KWEAVER_CONFIG_PATH        Default: /opt/v9-alpha-deploy/kweaver-config/config.yaml.
  KWEAVER_RESOURCE_NAMESPACE Default: resource.
  KWEAVER_INSECURE           1 or 0.
EOF
}

error() {
  printf '[ERROR] %s\n' "$*" >&2
  exit 1
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || error "Required command not found: $1"
}

KWEAVER_BASE_URL="${KWEAVER_BASE_URL:-}"
KWEAVER_BUSINESS_DOMAIN="${KWEAVER_BUSINESS_DOMAIN:-bd_public}"
SKILLS_ROOT="${ANYBACKUP_SKILLS_ROOT:-$DEFAULT_SKILLS_ROOT}"
SKILL_DIRS=()
if [[ -n "${ANYBACKUP_SKILL_DIRS:-}" ]]; then
  IFS=',' read -r -a SKILL_DIRS <<< "$ANYBACKUP_SKILL_DIRS"
fi
STATE_FILE="${ANYBACKUP_SKILLS_STATE_FILE:-}"
SKILL_SOURCE="${ANYBACKUP_SKILLS_SOURCE:-custom}"
PUBLISH="${ANYBACKUP_SKILLS_PUBLISH:-1}"
ENSURE_PUBLISH_SCHEMA="${ANYBACKUP_SKILLS_ENSURE_PUBLISH_SCHEMA:-1}"
KWEAVER_CONFIG_PATH="${KWEAVER_CONFIG_PATH:-/opt/v9-alpha-deploy/kweaver-config/config.yaml}"
KWEAVER_RESOURCE_NAMESPACE="${KWEAVER_RESOURCE_NAMESPACE:-resource}"
KWEAVER_INSECURE="${KWEAVER_INSECURE:-0}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --base-url|--url)
      KWEAVER_BASE_URL="${2:-}"
      shift 2
      ;;
    --biz-domain)
      KWEAVER_BUSINESS_DOMAIN="${2:-}"
      shift 2
      ;;
    --skills-root)
      SKILLS_ROOT="${2:-}"
      shift 2
      ;;
    --skill-dir)
      SKILL_DIRS+=("${2:-}")
      shift 2
      ;;
    --state-file)
      STATE_FILE="${2:-}"
      shift 2
      ;;
    --source)
      SKILL_SOURCE="${2:-}"
      shift 2
      ;;
    --publish)
      PUBLISH="1"
      shift
      ;;
    --no-publish)
      PUBLISH="0"
      shift
      ;;
    --kweaver-config)
      KWEAVER_CONFIG_PATH="${2:-}"
      shift 2
      ;;
    --resource-namespace)
      KWEAVER_RESOURCE_NAMESPACE="${2:-}"
      shift 2
      ;;
    --skip-publish-schema-check)
      ENSURE_PUBLISH_SCHEMA="0"
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

require_command python3
require_command kweaver

[[ -n "$KWEAVER_BASE_URL" ]] || error "--base-url or KWEAVER_BASE_URL is required."
[[ -n "$KWEAVER_BUSINESS_DOMAIN" ]] || error "--biz-domain or KWEAVER_BUSINESS_DOMAIN is required."

if [[ ${#SKILL_DIRS[@]} -eq 0 ]]; then
  [[ -d "$SKILLS_ROOT" ]] || error "Skills root not found: $SKILLS_ROOT"
  while IFS= read -r dir; do
    SKILL_DIRS+=("$dir")
  done < <(find "$SKILLS_ROOT" -mindepth 1 -maxdepth 1 -type d | sort)
fi

[[ ${#SKILL_DIRS[@]} -gt 0 ]] || error "No skill directories found."
for dir in "${SKILL_DIRS[@]}"; do
  [[ -d "$dir" ]] || error "Skill directory not found: $dir"
  [[ -f "$dir/SKILL.md" ]] || error "Missing SKILL.md in skill directory: $dir"
done

python3 - "$KWEAVER_BASE_URL" "$KWEAVER_BUSINESS_DOMAIN" "$STATE_FILE" "$SKILL_SOURCE" "$PUBLISH" "$KWEAVER_INSECURE" "$ENSURE_PUBLISH_SCHEMA" "$KWEAVER_CONFIG_PATH" "$KWEAVER_RESOURCE_NAMESPACE" "${SKILL_DIRS[@]}" <<'PY'
import json
import os
import re
import subprocess
import sys
import tempfile
import zipfile
from pathlib import Path

(
    base_url,
    biz_domain,
    state_file,
    source,
    publish_flag,
    insecure_flag,
    ensure_publish_schema_flag,
    kweaver_config_path,
    resource_namespace,
    *skill_dirs,
) = sys.argv[1:]
publish = publish_flag == "1"
insecure = insecure_flag == "1"
ensure_publish_schema = ensure_publish_schema_flag == "1"

env = os.environ.copy()
env["KWEAVER_BASE_URL"] = base_url
env["KWEAVER_BUSINESS_DOMAIN"] = biz_domain
if insecure:
    env["KWEAVER_TLS_INSECURE"] = "1"
    env["NODE_TLS_REJECT_UNAUTHORIZED"] = "0"

help_proc = subprocess.run(
    ["kweaver", "--help"],
    text=True,
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
    env=env,
)
supports_base_url = "--base-url" in (help_proc.stdout + help_proc.stderr)


def log(message: str) -> None:
    print(f"[INFO] {message}", flush=True)


def run_kweaver(args, check=True):
    cmd = ["kweaver"]
    if supports_base_url:
        cmd.extend(["--base-url", base_url])
    cmd.extend(args)
    proc = subprocess.run(cmd, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, env=env)
    if check and proc.returncode != 0:
        raise RuntimeError(
            "KWeaver CLI failed: "
            + " ".join(cmd)
            + f"\nrc={proc.returncode}\nstdout={proc.stdout}\nstderr={proc.stderr}"
        )
    return proc


def run_cmd(args, *, check=True, input_text=None, safe_name="command"):
    proc = subprocess.run(
        args,
        input=input_text,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        env=env,
    )
    if check and proc.returncode != 0:
        raise RuntimeError(
            f"{safe_name} failed"
            + f"\nrc={proc.returncode}\nstdout={proc.stdout}\nstderr={proc.stderr}"
        )
    return proc


def parse_json_payload(text: str):
    text = text.strip()
    if not text:
        return None
    decoder = json.JSONDecoder()
    starts = [idx for idx in (text.find("{"), text.find("[")) if idx >= 0]
    if not starts:
        return None
    start = min(starts)
    payload, _ = decoder.raw_decode(text[start:])
    return payload


def list_items(payload):
    if isinstance(payload, list):
        return payload
    if isinstance(payload, dict):
        for path in (
            ("data", "items"),
            ("data", "records"),
            ("data", "list"),
            ("items",),
            ("records",),
            ("entries",),
            ("data",),
        ):
            cur = payload
            for part in path:
                if not isinstance(cur, dict) or part not in cur:
                    cur = None
                    break
                cur = cur[part]
            if isinstance(cur, list):
                return cur
    return []


def item_id(item: dict) -> str:
    for key in ("id", "skill_id", "skillId"):
        value = item.get(key)
        if value:
            return str(value)
    return ""


def item_name(item: dict) -> str:
    for key in ("name", "skill_name", "skillName"):
        value = item.get(key)
        if value:
            return str(value)
    return ""


def item_status(item: dict) -> str:
    for key in ("status", "skill_status", "skillStatus"):
        value = item.get(key)
        if value:
            return str(value)
    return ""


def skill_name(skill_dir: Path) -> str:
    text = (skill_dir / "SKILL.md").read_text(encoding="utf-8")
    match = re.search(r"(?m)^name:\s*(.+?)\s*$", text)
    if not match:
        raise RuntimeError(f"Unable to read skill name from {skill_dir / 'SKILL.md'}")
    return match.group(1).strip().strip('"').strip("'")


def zip_skill_dir(skill_dir: Path, output_dir: Path) -> Path:
    zip_path = output_dir / f"{skill_dir.name}.zip"
    with zipfile.ZipFile(zip_path, "w", compression=zipfile.ZIP_DEFLATED) as archive:
        for path in sorted(skill_dir.rglob("*")):
            if path.is_file():
                archive.write(path, path.relative_to(skill_dir).as_posix())
    return zip_path


def find_skill_by_name(name: str):
    proc = run_kweaver(["skill", "list", "--name", name, "-bd", biz_domain])
    payload = parse_json_payload(proc.stdout)
    for item in list_items(payload):
        if isinstance(item, dict) and item_name(item) == name and item_id(item):
            return item
    return None


def resolve_registered_skill(name: str, stdout: str) -> dict:
    payload = parse_json_payload(stdout)
    candidates = []
    if isinstance(payload, dict):
        candidates.append(payload)
        for key in ("data", "item", "skill"):
            value = payload.get(key)
            if isinstance(value, dict):
                candidates.append(value)
    for item in candidates:
        sid = item_id(item)
        if sid:
            return item
    item = find_skill_by_name(name)
    if item:
        return item
    raise RuntimeError(f"Unable to resolve imported skill id for {name}")


def yaml_value(value: str) -> str:
    value = value.strip()
    if " #" in value:
        value = value.split(" #", 1)[0].rstrip()
    if len(value) >= 2 and value[0] == value[-1] and value[0] in ("'", '"'):
        return value[1:-1]
    return value


def read_rds_config(path: str) -> dict:
    config = Path(path)
    if not config.exists():
        raise RuntimeError(f"KWeaver config file not found: {path}")

    lines = config.read_text(encoding="utf-8").splitlines()
    for idx, line in enumerate(lines):
        stripped = line.strip()
        if not stripped or stripped.startswith("#") or stripped != "rds:":
            continue

        base_indent = len(line) - len(line.lstrip())
        block = {}
        for child in lines[idx + 1 :]:
            child_stripped = child.strip()
            if not child_stripped or child_stripped.startswith("#"):
                continue
            indent = len(child) - len(child.lstrip())
            if indent <= base_indent:
                break
            if ":" not in child_stripped:
                continue
            key, value = child_stripped.split(":", 1)
            block[key.strip()] = yaml_value(value)

        if all(block.get(key) for key in ("user", "password", "database")):
            return block

    raise RuntimeError(f"Unable to resolve rds.user/password/database from {path}")


def discover_mariadb_pod(namespace: str) -> str:
    proc = run_cmd(
        ["kubectl", "get", "pods", "-n", namespace, "-o", "json"],
        safe_name="kubectl get mariadb pods",
    )
    payload = json.loads(proc.stdout)
    candidates = []
    for item in payload.get("items", []):
        metadata = item.get("metadata", {})
        status = item.get("status", {})
        name = metadata.get("name", "")
        if "mariadb" not in name:
            continue
        phase = status.get("phase", "")
        ready = False
        for condition in status.get("conditions", []):
            if condition.get("type") == "Ready" and condition.get("status") == "True":
                ready = True
                break
        candidates.append((phase == "Running" and ready, name))
    candidates.sort(reverse=True)
    if candidates:
        return candidates[0][1]
    raise RuntimeError(f"Unable to find a MariaDB pod in namespace {namespace}")


def ensure_skill_publish_schema() -> bool:
    if not publish or not ensure_publish_schema:
        return False

    rds = read_rds_config(kweaver_config_path)
    pod = discover_mariadb_pod(resource_namespace)
    sql = r"""
CREATE TABLE IF NOT EXISTS `t_skill_release` (
    `f_id` bigint AUTO_INCREMENT NOT NULL,
    `f_skill_id` varchar(40) NOT NULL DEFAULT '',
    `f_name` varchar(255) NOT NULL DEFAULT '',
    `f_description` longtext DEFAULT NULL,
    `f_skill_content` longtext DEFAULT NULL,
    `f_version` varchar(40) NOT NULL DEFAULT '',
    `f_status` varchar(40) NOT NULL DEFAULT 'published',
    `f_source` varchar(50) NOT NULL DEFAULT '',
    `f_extend_info` longtext DEFAULT NULL,
    `f_dependencies` longtext DEFAULT NULL,
    `f_file_manifest` longtext DEFAULT NULL,
    `f_create_user` varchar(50) NOT NULL DEFAULT '',
    `f_create_time` bigint(20) NOT NULL DEFAULT 0,
    `f_update_user` varchar(50) NOT NULL DEFAULT '',
    `f_update_time` bigint(20) NOT NULL DEFAULT 0,
    `f_release_desc` varchar(255) NOT NULL DEFAULT '',
    `f_release_user` varchar(50) NOT NULL DEFAULT '',
    `f_release_time` bigint(20) NOT NULL DEFAULT 0,
    `f_delete_user` varchar(50) NOT NULL DEFAULT '',
    `f_delete_time` bigint(20) NOT NULL DEFAULT 0,
    `f_category` varchar(50) DEFAULT '',
    `f_is_deleted` boolean DEFAULT 0,
    `f_creation_type` varchar(20) NOT NULL DEFAULT 'custom',
    `f_config_source` varchar(40) NOT NULL DEFAULT '',
    `f_component_id` varchar(40) NOT NULL DEFAULT '',
    PRIMARY KEY (`f_id`),
    UNIQUE KEY `uk_skill_release` (`f_skill_id`, `f_version`) USING BTREE,
    KEY `idx_skill_id_create_time` (`f_skill_id`, `f_create_time`) USING BTREE,
    KEY `idx_status_update_time` (`f_status`, `f_update_time`) USING BTREE,
    KEY `idx_category_update_time` (`f_category`, `f_update_time`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'Skill release table compatibility';

CREATE TABLE IF NOT EXISTS `t_skill_release_history` (
    `f_id` bigint AUTO_INCREMENT NOT NULL,
    `f_skill_id` varchar(40) NOT NULL DEFAULT '',
    `f_skill_release` longtext DEFAULT NULL,
    `f_version` varchar(40) NOT NULL DEFAULT '',
    `f_skill_version` varchar(40) NOT NULL DEFAULT '',
    `f_release_desc` varchar(255) NOT NULL DEFAULT '',
    `f_create_user` varchar(50) NOT NULL DEFAULT '',
    `f_create_time` bigint(20) NOT NULL DEFAULT 0,
    `f_update_user` varchar(50) NOT NULL DEFAULT '',
    `f_update_time` bigint(20) NOT NULL DEFAULT 0,
    `f_release_user` varchar(50) NOT NULL DEFAULT '',
    `f_release_time` bigint(20) NOT NULL DEFAULT 0,
    `f_metadata_version` varchar(40) NOT NULL DEFAULT '',
    `f_metadata_type` varchar(40) NOT NULL DEFAULT '',
    PRIMARY KEY (`f_id`),
    UNIQUE KEY `uk_skill_release_history` (`f_skill_id`, `f_version`) USING BTREE,
    KEY `idx_skill_id_create_time` (`f_skill_id`, `f_create_time`) USING BTREE,
    KEY `idx_skill_id_metadata_version` (`f_skill_id`, `f_metadata_version`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = 'Skill release history table compatibility';
"""
    run_cmd(
        [
            "kubectl",
            "exec",
            "-i",
            "-n",
            resource_namespace,
            pod,
            "--",
            "env",
            f"MYSQL_PWD={rds['password']}",
            "mariadb",
            f"-u{rds['user']}",
            "-D",
            rds["database"],
        ],
        input_text=sql,
        safe_name="ensure KWeaver skill publish schema",
    )
    log(f"Ensured KWeaver skill publish schema in {resource_namespace}/{pod}")
    return True


def publish_skill(skill_id: str, status: str) -> str:
    if not publish:
        return "skipped"
    if status == "published":
        return "already_published"
    proc = run_kweaver(["skill", "set-status", skill_id, "published", "-bd", biz_domain], check=False)
    if proc.returncode == 0:
        return "published"
    fallback = run_kweaver(["skill", "status", skill_id, "published", "-bd", biz_domain], check=False)
    if fallback.returncode == 0:
        return "published"
    raise RuntimeError(
        f"Failed to publish skill {skill_id}\n"
        f"set-status stdout={proc.stdout}\nset-status stderr={proc.stderr}\n"
        f"status stdout={fallback.stdout}\nstatus stderr={fallback.stderr}"
    )


state = {
    "base_url": base_url,
    "business_domain": biz_domain,
    "publish_schema_checked": False,
    "skills": {},
}

state["publish_schema_checked"] = ensure_skill_publish_schema()

for raw_dir in skill_dirs:
    skill_dir = Path(raw_dir).resolve()
    name = skill_name(skill_dir)
    existing = find_skill_by_name(name)
    action = "reused"
    if existing:
        skill = existing
        log(f"Reuse existing skill: {name} ({item_id(skill)})")
    else:
        log(f"Register skill: {name} from {skill_dir}")
        with tempfile.TemporaryDirectory(prefix="anybackup-skill-") as tmpdir:
            zip_path = zip_skill_dir(skill_dir, Path(tmpdir))
            proc = run_kweaver(
                [arg for arg in [
                    "skill",
                    "register",
                    "--zip-file",
                    str(zip_path),
                    *(["--source", source] if source else []),
                    "-bd",
                    biz_domain,
                ] if arg]
            )
        skill = resolve_registered_skill(name, proc.stdout)
        action = "registered"
        log(f"Registered skill: {name} ({item_id(skill)})")
    sid = item_id(skill)
    publish_result = publish_skill(sid, item_status(skill))
    if publish_result == "published":
        log(f"Published skill: {name} ({sid})")
    elif publish_result == "already_published":
        log(f"Skill already published: {name} ({sid})")
    state["skills"][name] = {
        "id": sid,
        "source_dir": str(skill_dir),
        "action": action,
        "published": publish_result in ("published", "already_published"),
        "publish_action": publish_result,
    }

if state_file:
    path = Path(state_file)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(state, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    log(f"Wrote skill state: {path}")

print(json.dumps(state, ensure_ascii=False, indent=2))
PY
