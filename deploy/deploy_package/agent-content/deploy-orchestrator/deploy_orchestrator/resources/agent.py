from __future__ import annotations

import json
from pathlib import Path

from deploy_orchestrator.cli_runner import CliRunner
from deploy_orchestrator.context import DeploymentContext
from deploy_orchestrator.logger import get_logger
from deploy_orchestrator.renderer import render_template_file
from deploy_orchestrator.schema import AgentSpec


class AgentExecutor:
    resource_kind = "agent"

    def __init__(
        self,
        runner: CliRunner | None = None,
        workspace: str | Path | None = None,
        business_domain: str | None = None,
    ):
        self.runner = runner or CliRunner()
        self.workspace = Path(workspace or ".")
        self.business_domain = business_domain or ""
        self.logger = get_logger(__name__)

    def plan(self, spec: AgentSpec, discovered: dict[str, object] | None) -> str:
        return "update" if discovered else "create"

    def apply(
        self,
        spec: AgentSpec,
        discovered: dict[str, object] | None,
        context: DeploymentContext,
        state: dict[str, object],
    ) -> dict[str, object]:
        self.logger.info("Applying resource agent/%s", spec.name)
        payload = self._create_or_update(spec, discovered, context)
        state["agents"][spec.name] = payload
        context.set_resource("agent", spec.name, payload)
        self.logger.info("Applied resource agent/%s", spec.name)
        return payload

    def _create_or_update(
        self,
        spec: AgentSpec,
        discovered: dict[str, object] | None,
        context: DeploymentContext,
    ) -> dict[str, object]:
        rendered_config_path = self._render_config(spec, context)
        if discovered:
            payload = dict(discovered)
            update_cmd = ["kweaver", "agent", "update", str(discovered["id"])]
            update_cmd.extend(["--name", spec.name, "--profile", spec.profile])
            if spec.knowledge_network:
                kn = context.values["bkns"].get(spec.knowledge_network, {})
                if "id" in kn:
                    update_cmd.extend(["--knowledge-network-id", str(kn["id"])])
            if rendered_config_path is not None:
                update_cmd.extend(["--config-path", str(rendered_config_path)])
            if self.business_domain:
                update_cmd.extend(["-bd", self.business_domain])
            self.runner.run_checked(update_cmd)
            payload.setdefault("status", "ok")
            payload.setdefault("profile", spec.profile)
            payload.setdefault("knowledge_network", spec.knowledge_network)
            payload.setdefault("skills", list(spec.skills))
            return payload

        create_cmd = [
            "kweaver",
            "agent",
            "create",
            "--name",
            spec.name,
            "--profile",
            spec.profile,
            "--llm-id",
            spec.llm_id,
        ]
        if rendered_config_path is not None:
            create_cmd.extend(["--config", str(rendered_config_path)])
        if self.business_domain:
            create_cmd.extend(["-bd", self.business_domain])
        result = self.runner.run_checked(create_cmd)
        payload = self._parse_json(result.stdout)
        payload.setdefault("id", f"agent::{spec.name}")
        payload.setdefault("name", spec.name)
        payload.setdefault("status", "ok")
        payload.setdefault("profile", spec.profile)
        payload.setdefault("knowledge_network", spec.knowledge_network)
        payload.setdefault("skills", list(spec.skills))
        return payload

    def _render_config(self, spec: AgentSpec, context: DeploymentContext) -> Path | None:
        if not spec.config_template:
            return None
        template_path = Path(spec.config_template)
        if not template_path.is_absolute():
            template_path = self.workspace / template_path
        output_dir = self.workspace / ".deploy" / "rendered"
        output_dir.mkdir(parents=True, exist_ok=True)
        output_path = output_dir / f"{spec.name}.config.json"
        self.logger.info("Rendering agent config for %s from %s", spec.name, template_path)
        render_template_file(template_path, output_path, context.values)
        return output_path

    @staticmethod
    def _parse_json(text: str) -> dict[str, object]:
        try:
            data = json.loads(text)
        except json.JSONDecodeError:
            return {}
        return data if isinstance(data, dict) else {}
