from __future__ import annotations

from pathlib import Path
from pathlib import PurePosixPath, PureWindowsPath

import yaml

from deploy_orchestrator.logger import get_logger
from deploy_orchestrator.schema import AgentSpec, BindingSpec, BknSpec, DeployConfig, SkillSpec


def load_config(path: str | Path) -> DeployConfig:
    file_path = Path(path)
    logger = get_logger(__name__)
    logger.info("Loading config file from %s", file_path)
    data = yaml.safe_load(file_path.read_text(encoding="utf-8")) or {}
    base_dir = file_path.resolve().parent
    logger.info("Resolving declared paths from config directory %s", base_dir)
    return DeployConfig(
        project=data["project"],
        env=data.get("env", "dev"),
        globals=data.get("globals", {}),
        skills=[SkillSpec(**_resolve_skill_item(item, base_dir)) for item in data.get("skills", [])],
        bkns=[BknSpec(**_resolve_bkn_item(item, base_dir)) for item in data.get("bkns", [])],
        agents=[AgentSpec(**_resolve_agent_item(item, base_dir)) for item in data.get("agents", [])],
        bindings=[BindingSpec(**item) for item in data.get("bindings", [])],
    )


def _resolve_skill_item(item: dict[str, object], base_dir: Path) -> dict[str, object]:
    resolved = dict(item)
    resolved["source"] = _resolve_declared_path(str(item["source"]), base_dir)
    return resolved


def _resolve_bkn_item(item: dict[str, object], base_dir: Path) -> dict[str, object]:
    resolved = dict(item)
    resolved["source"] = _resolve_declared_path(str(item["source"]), base_dir)
    return resolved


def _resolve_agent_item(item: dict[str, object], base_dir: Path) -> dict[str, object]:
    resolved = dict(item)
    config_template = item.get("config_template")
    if config_template:
        resolved["config_template"] = _resolve_declared_path(str(config_template), base_dir)
    return resolved


def _resolve_declared_path(raw_path: str, base_dir: Path) -> str:
    declared_path = Path(raw_path)
    if declared_path.is_absolute():
        return raw_path

    for pure_path in (PureWindowsPath(raw_path), PurePosixPath(raw_path)):
        if pure_path.anchor:
            return raw_path
    return str(base_dir / raw_path)
