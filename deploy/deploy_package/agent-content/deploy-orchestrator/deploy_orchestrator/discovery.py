from __future__ import annotations

import json

from deploy_orchestrator.cli_runner import CliRunner
from deploy_orchestrator.logger import get_logger


class DiscoveryError(RuntimeError):
    pass


class AmbiguousResourceError(DiscoveryError):
    pass


class DiscoveryService:
    def __init__(self, runner: CliRunner | None = None, business_domain: str | None = None):
        self.runner = runner or CliRunner()
        self.business_domain = business_domain or ""
        self.logger = get_logger(__name__)

    def find_by_name(self, resource_kind: str, name: str) -> dict[str, object] | None:
        self.logger.info("Discovering resource %s/%s", resource_kind, name)
        result = self.runner.run_checked(self._build_list_command(resource_kind, name))
        items = self._parse_items(result.stdout)
        matched = [item for item in items if str(item.get("name", "")) == name]
        if not matched:
            self.logger.info("Resource %s/%s was not found", resource_kind, name)
            return None
        if len(matched) > 1:
            self.logger.error("Resource %s/%s is ambiguous", resource_kind, name)
            raise AmbiguousResourceError(f"发现多个同名 {resource_kind}: {name}")
        self.logger.info("Resource %s/%s resolved to id=%s", resource_kind, name, matched[0].get("id"))
        return matched[0]

    def _build_list_command(self, resource_kind: str, name: str) -> list[str]:
        if resource_kind == "skill":
            command = ["kweaver", "skill", "list", "--name", name, "--compact"]
        elif resource_kind == "bkn":
            command = ["kweaver", "bkn", "list", "--name-pattern", name]
        elif resource_kind == "agent":
            command = ["kweaver", "agent", "list", "--name", name]
        elif resource_kind == "binding":
            return ["python", "-c", "print('[]')"]
        else:
            raise DiscoveryError(f"不支持的资源类型: {resource_kind}")
        if self.business_domain:
            command.extend(["-bd", self.business_domain])
        return command

    @staticmethod
    def _parse_items(stdout: str) -> list[dict[str, object]]:
        text = stdout.strip() or "[]"
        data = json.loads(text)
        if isinstance(data, list):
            return [item for item in data if isinstance(item, dict)]
        if isinstance(data, dict):
            entries = data.get("entries") or data.get("data") or []
            if isinstance(entries, list):
                return [item for item in entries if isinstance(item, dict)]
        return []
