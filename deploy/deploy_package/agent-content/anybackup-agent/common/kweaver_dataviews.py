#!/usr/bin/env python3
from __future__ import annotations

import argparse
import hashlib
import json
import subprocess
import sys
from typing import Any


def run(args: list[str], *, check: bool = True) -> subprocess.CompletedProcess[str]:
    result = subprocess.run(args, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if result.stderr:
        print(result.stderr, file=sys.stderr, end="")
    if check and result.returncode != 0:
        raise SystemExit(result.returncode)
    return result


def run_json(args: list[str], *, check: bool = True) -> Any:
    result = run(args, check=check)
    text = result.stdout.strip()
    if not text:
        return None
    return json.loads(text)


def entries(payload: Any) -> list[dict[str, Any]]:
    if isinstance(payload, list):
        return [item for item in payload if isinstance(item, dict)]
    if isinstance(payload, dict):
        for key in ("entries", "data", "items"):
            value = payload.get(key)
            if isinstance(value, list):
                return [item for item in value if isinstance(item, dict)]
    return []


def normalize_engine(engine: str) -> str:
    value = engine.lower()
    if value in {"postgres", "postgresql", "pg"}:
        return "postgresql"
    if value in {"mysql", "mariadb"}:
        return "mysql"
    return value


def find_datasource_id(name: str, business_domain: str) -> str:
    payload = run_json(
        [
            "kweaver",
            "ds",
            "list",
            "--keyword",
            name,
            "-bd",
            business_domain,
            "--pretty",
        ]
    )
    for item in entries(payload):
        if item.get("name") == name and item.get("id"):
            return str(item["id"])
    return ""


def ensure_datasource(args: argparse.Namespace) -> str:
    existing = find_datasource_id(args.datasource_name, args.biz_domain)
    if existing:
        print(f"[INFO] Reusing KWeaver datasource: {args.datasource_name} ({existing})")
        return existing

    print(f"[INFO] Creating KWeaver datasource: {args.datasource_name}")
    connect_args = [
        "kweaver",
        "ds",
        "connect",
        normalize_engine(args.engine),
        args.host,
        str(args.port),
        args.database,
        "--account",
        args.username,
        "--password",
        args.password,
        "--name",
        args.datasource_name,
        "-bd",
        args.biz_domain,
        "--pretty",
    ]
    if args.schema:
        connect_args.extend(["--schema", args.schema])

    result = run(connect_args, check=False)
    if result.returncode != 0:
        existing = find_datasource_id(args.datasource_name, args.biz_domain)
        if existing:
            return existing
        print(result.stdout, file=sys.stderr, end="")
        raise SystemExit(result.returncode)

    datasource_id = find_datasource_id(args.datasource_name, args.biz_domain)
    if not datasource_id:
        print(f"[ERROR] Unable to resolve datasource ID for {args.datasource_name}", file=sys.stderr)
        raise SystemExit(1)
    return datasource_id


def find_dataview_id(name: str, datasource_id: str, business_domain: str) -> str:
    payload = run_json(
        [
            "kweaver",
            "dataview",
            "find",
            "--name",
            name,
            "--exact",
            "--datasource-id",
            datasource_id,
            "-bd",
            business_domain,
            "--pretty",
        ]
    )
    for item in entries(payload):
        if item.get("name") == name and item.get("id"):
            return str(item["id"])
    return ""


def table_matches(candidate: str, table: str) -> bool:
    return candidate == table or candidate.endswith(f".{table}") or table.endswith(f".{candidate}")


def get_table_metadata(datasource_id: str, table: str, business_domain: str) -> dict[str, Any]:
    last_payload: Any = None
    for attempt in range(1, 5):
        metadata_path = (
            f"/api/data-connection/v1/metadata/data-source/{datasource_id}"
            f"?limit=-1&keyword={table}"
        )
        payload = run_json(
            [
                "kweaver",
                "call",
                metadata_path,
                "-bd",
                business_domain,
                "--pretty",
            ]
        )
        last_payload = payload
        for item in entries(payload):
            name = str(item.get("name") or "")
            if table_matches(name, table):
                return item

    print(
        f"[WARN] Table {table} was not discovered in datasource {datasource_id}; "
        f"creating a minimal data view without scanned fields. Last payload: {last_payload}",
        file=sys.stderr,
    )
    return {"name": table, "columns": []}


def fields_from_table(table_metadata: dict[str, Any]) -> list[dict[str, str]]:
    raw_columns = table_metadata.get("columns") or table_metadata.get("fields") or []
    fields: list[dict[str, str]] = []
    if isinstance(raw_columns, list):
        for column in raw_columns:
            if not isinstance(column, dict):
                continue
            name = column.get("name") or column.get("field_name")
            if not name:
                continue
            fields.append(
                {
                    "name": str(name),
                    "type": str(column.get("type") or column.get("field_type") or "varchar"),
                }
            )
    return fields


def create_dataview(
    datasource_id: str,
    table: str,
    view_name: str,
    business_domain: str,
) -> str:
    table_metadata = get_table_metadata(datasource_id, table, business_domain)
    view_id = hashlib.md5(f"{datasource_id}:{table}".encode()).hexdigest()[:35]
    body = [
        {
            "id": view_id,
            "name": view_name,
            "technical_name": str(table_metadata.get("name") or table),
            "type": "atomic",
            "query_type": "SQL",
            "data_source_id": datasource_id,
            "group_id": datasource_id,
            "fields": fields_from_table(table_metadata),
        }
    ]
    result = run(
        [
            "kweaver",
            "call",
            "/api/mdl-data-model/v1/data-views",
            "-X",
            "POST",
            "-d",
            json.dumps(body, ensure_ascii=False, separators=(",", ":")),
            "-bd",
            business_domain,
            "--pretty",
        ],
        check=False,
    )
    if result.returncode != 0:
        existing = find_dataview_id(view_name, datasource_id, business_domain)
        if existing:
            return existing
        print(result.stdout, file=sys.stderr, end="")
        raise SystemExit(result.returncode)

    existing = find_dataview_id(view_name, datasource_id, business_domain)
    return existing or view_id


def parse_views(values: list[str]) -> list[tuple[str, str]]:
    views: list[tuple[str, str]] = []
    for value in values:
        if "=" in value:
            table, view_name = value.split("=", 1)
        else:
            table = value
            view_name = value
        table = table.strip()
        view_name = view_name.strip()
        if not table or not view_name:
            raise SystemExit(f"Invalid --view value: {value}")
        views.append((table, view_name))
    return views


def main() -> int:
    parser = argparse.ArgumentParser(description="Ensure KWeaver datasource and atomic data views")
    parser.add_argument("--biz-domain", default="bd_public")
    parser.add_argument("--datasource-name", required=True)
    parser.add_argument("--engine", required=True)
    parser.add_argument("--host", required=True)
    parser.add_argument("--port", required=True)
    parser.add_argument("--database", required=True)
    parser.add_argument("--schema", default="")
    parser.add_argument("--username", required=True)
    parser.add_argument("--password", required=True)
    parser.add_argument("--view", action="append", default=[], help="table or table=view-name")
    args = parser.parse_args()

    datasource_id = ensure_datasource(args)
    mapping: dict[str, str] = {}
    for table, view_name in parse_views(args.view):
        dataview_id = find_dataview_id(view_name, datasource_id, args.biz_domain)
        if not dataview_id:
            print(f"[INFO] Creating KWeaver data view: {view_name} ({table})")
            dataview_id = create_dataview(datasource_id, table, view_name, args.biz_domain)
        else:
            print(f"[INFO] Reusing KWeaver data view: {view_name} ({dataview_id})")
        mapping[view_name] = dataview_id

    print(json.dumps({"datasource_id": datasource_id, "data_views": mapping}, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
