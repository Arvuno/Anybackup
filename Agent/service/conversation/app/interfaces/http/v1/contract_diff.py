from pathlib import Path
from typing import Any

import yaml


def find_missing_contract_paths(
    *,
    upstream_contract_path: Path,
    generated_openapi: dict[str, Any],
    service_prefix: str,
) -> set[str]:
    upstream = yaml.safe_load(upstream_contract_path.read_text(encoding="utf-8"))
    upstream_paths = set(upstream.get("paths", {}))
    generated_paths = set(generated_openapi.get("paths", {}))

    expected_generated_paths = {f"{service_prefix}{path}" for path in upstream_paths}
    return expected_generated_paths - generated_paths
