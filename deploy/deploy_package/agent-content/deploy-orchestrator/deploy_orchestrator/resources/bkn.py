from __future__ import annotations

import json

from deploy_orchestrator.cli_runner import CliRunner
from deploy_orchestrator.context import DeploymentContext
from deploy_orchestrator.logger import get_logger
from deploy_orchestrator.schema import BknSpec


class BknExecutor:
    resource_kind = "bkn"

    def __init__(self, runner: CliRunner | None = None, business_domain: str | None = None):
        self.runner = runner or CliRunner()
        self.business_domain = business_domain or ""
        self.logger = get_logger(__name__)

    def plan(self, spec: BknSpec, discovered: dict[str, object] | None) -> str:
        return "update" if discovered else "create"

    def apply(
        self,
        spec: BknSpec,
        discovered: dict[str, object] | None,
        context: DeploymentContext,
        state: dict[str, object],
    ) -> dict[str, object]:
        self.logger.info("Applying resource bkn/%s", spec.name)
        payload = self._push(spec, discovered)
        state["bkns"][spec.name] = payload
        context.set_resource("bkn", spec.name, payload)
        self.logger.info("Applied resource bkn/%s", spec.name)
        return payload

    def _push(self, spec: BknSpec, discovered: dict[str, object] | None) -> dict[str, object]:
        if spec.validate:
            validate_cmd = ["kweaver", "bkn", "validate", spec.source]
            if self.business_domain:
                validate_cmd.extend(["-bd", self.business_domain])
            self.runner.run_checked(validate_cmd)

        push_cmd = ["kweaver", "bkn", "push", spec.source, "--branch", spec.branch]
        if self.business_domain:
            push_cmd.extend(["-bd", self.business_domain])
        result = self.runner.run_checked(push_cmd)

        payload = discovered or self._parse_json(result.stdout)
        payload.setdefault("id", f"bkn::{spec.name}")
        payload.setdefault("name", spec.name)
        payload.setdefault("status", "ok")
        payload.setdefault("source", spec.source)
        payload.setdefault("branch", spec.branch)
        return payload

    @staticmethod
    def _parse_json(text: str) -> dict[str, object]:
        try:
            data = json.loads(text)
        except json.JSONDecodeError:
            return {}
        return data if isinstance(data, dict) else {}
