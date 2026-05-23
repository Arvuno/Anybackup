from __future__ import annotations

import json

from deploy_orchestrator.cli_runner import CliRunner
from deploy_orchestrator.context import DeploymentContext
from deploy_orchestrator.logger import get_logger
from deploy_orchestrator.schema import SkillSpec


class SkillExecutor:
    resource_kind = "skill"

    def __init__(self, runner: CliRunner | None = None, business_domain: str | None = None):
        self.runner = runner or CliRunner()
        self.business_domain = business_domain or ""
        self.logger = get_logger(__name__)

    def plan(self, spec: SkillSpec, discovered: dict[str, object] | None) -> str:
        return "skip" if discovered else "create"

    def apply(
        self,
        spec: SkillSpec,
        discovered: dict[str, object] | None,
        context: DeploymentContext,
        state: dict[str, object],
    ) -> dict[str, object]:
        self.logger.info("Applying resource skill/%s", spec.name)
        payload = discovered or self._register(spec)
        state["skills"][spec.name] = payload
        context.set_resource("skill", spec.name, payload)
        self.logger.info("Applied resource skill/%s", spec.name)
        return payload

    def _register(self, spec: SkillSpec) -> dict[str, object]:
        command = ["kweaver", "skill", "register", "--zip-file", spec.source]
        if self.business_domain:
            command.extend(["-bd", self.business_domain])
        result = self.runner.run_checked(command)
        payload = self._parse_json(result.stdout)
        payload.setdefault("id", f"skill::{spec.name}")
        payload.setdefault("name", spec.name)
        payload.setdefault("status", "ok")
        payload.setdefault("source", spec.source)
        return payload

    @staticmethod
    def _parse_json(text: str) -> dict[str, object]:
        try:
            data = json.loads(text)
        except json.JSONDecodeError:
            return {}
        return data if isinstance(data, dict) else {}
