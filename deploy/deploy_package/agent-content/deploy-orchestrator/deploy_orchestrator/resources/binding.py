from __future__ import annotations

import json
from pathlib import Path

from deploy_orchestrator.cli_runner import CliRunner
from deploy_orchestrator.context import DeploymentContext
from deploy_orchestrator.logger import get_logger
from deploy_orchestrator.schema import BindingSpec


class BindingExecutor:
    resource_kind = "binding"

    def __init__(
        self,
        runner: CliRunner | None = None,
        business_domain: str | None = None,
        workspace: str | Path | None = None,
    ):
        self.runner = runner or CliRunner()
        self.business_domain = business_domain or ""
        self.workspace = Path(workspace or ".")
        self.logger = get_logger(__name__)

    def plan(self, spec: BindingSpec, discovered: dict[str, object] | None) -> str:
        return "skip" if discovered else "create"

    def apply(
        self,
        spec: BindingSpec,
        discovered: dict[str, object] | None,
        context: DeploymentContext,
        state: dict[str, object],
    ) -> dict[str, object]:
        self.logger.info("Applying resource binding/%s", spec.name)
        if spec.type == "skill":
            payload = self._apply_skill_binding(spec, context)
        elif spec.type == "agent_member":
            payload = self._apply_agent_member_binding(spec, context)
        else:
            raise ValueError(f"binding 类型暂不支持: {spec.type}")
        state["bindings"][spec.name] = payload
        context.set_resource("binding", spec.name, payload)
        self.logger.info("Applied resource binding/%s", spec.name)
        return payload

    def _apply_skill_binding(self, spec: BindingSpec, context: DeploymentContext) -> dict[str, object]:
        if not spec.skill:
            raise ValueError("skill binding 缺少 skill 字段")
        agent = context.values["agents"].get(spec.agent)
        skill = context.values["skills"].get(spec.skill)
        if not agent or "id" not in agent:
            raise ValueError(f"未找到 agent 上下文: {spec.agent}")
        if not skill or "id" not in skill:
            raise ValueError(f"未找到 skill 上下文: {spec.skill}")

        agent_id = str(agent["id"])
        skill_id = str(skill["id"])
        list_cmd = ["kweaver", "agent", "skill", "list", agent_id, "--compact"]
        if self.business_domain:
            list_cmd.extend(["-bd", self.business_domain])
        result = self.runner.run_checked(list_cmd)
        attached = self._parse_skill_list(result.stdout)
        if skill_id not in attached:
            add_cmd = ["kweaver", "agent", "skill", "add", agent_id, skill_id]
            if self.business_domain:
                add_cmd.extend(["-bd", self.business_domain])
            self.runner.run_checked(add_cmd)

        return {
            "id": f"binding::{spec.name}",
            "status": "ok",
            "type": spec.type,
            "agent": spec.agent,
            "skill": spec.skill,
            "agent_id": agent_id,
            "skill_id": skill_id,
        }

    def _apply_agent_member_binding(self, spec: BindingSpec, context: DeploymentContext) -> dict[str, object]:
        if not spec.member_agent:
            raise ValueError("agent_member binding 缺少 member_agent 字段")
        if not spec.config_path:
            raise ValueError("agent_member binding 缺少 config_path 字段")
        id_field = spec.id_field or "agent_id"

        agent = context.values["agents"].get(spec.agent)
        member_agent = context.values["agents"].get(spec.member_agent)
        if not agent or "id" not in agent:
            raise ValueError(f"未找到 agent 上下文: {spec.agent}")
        if not member_agent or "id" not in member_agent:
            raise ValueError(f"未找到 member_agent 上下文: {spec.member_agent}")

        agent_id = str(agent["id"])
        member_agent_id = str(member_agent["id"])

        get_cmd = ["kweaver", "agent", "get", agent_id]
        if self.business_domain:
            get_cmd.extend(["-bd", self.business_domain])
        result = self.runner.run_checked(get_cmd)
        current = self._parse_object(result.stdout)
        config = current.get("config", {})
        if not isinstance(config, dict):
            config = {}

        new_config, changed = self._ensure_member_entry(config, spec.config_path, id_field, member_agent_id)
        if changed:
            rendered_dir = self.workspace / ".deploy" / "bindings"
            rendered_dir.mkdir(parents=True, exist_ok=True)
            config_path = rendered_dir / f"{spec.name}.config.json"
            config_path.write_text(json.dumps(new_config, ensure_ascii=False, indent=2), encoding="utf-8")

            update_cmd = ["kweaver", "agent", "update", agent_id, "--config-path", str(config_path)]
            if self.business_domain:
                update_cmd.extend(["-bd", self.business_domain])
            self.runner.run_checked(update_cmd)

        return {
            "id": f"binding::{spec.name}",
            "status": "ok",
            "type": spec.type,
            "agent": spec.agent,
            "member_agent": spec.member_agent,
            "agent_id": agent_id,
            "member_agent_id": member_agent_id,
            "config_path": spec.config_path,
            "id_field": id_field,
        }

    @staticmethod
    def _parse_skill_list(stdout: str) -> set[str]:
        text = stdout.strip() or "[]"
        data = json.loads(text)
        if not isinstance(data, list):
            return set()
        values: set[str] = set()
        for item in data:
            if isinstance(item, dict) and "skill_id" in item:
                values.add(str(item["skill_id"]))
        return values

    @staticmethod
    def _parse_object(stdout: str) -> dict[str, object]:
        text = stdout.strip() or "{}"
        data = json.loads(text)
        return data if isinstance(data, dict) else {}

    @staticmethod
    def _ensure_member_entry(
        config: dict[str, object],
        path: list[str],
        id_field: str,
        member_id: str,
    ) -> tuple[dict[str, object], bool]:
        new_config = json.loads(json.dumps(config))
        cursor: dict[str, object] = new_config
        for key in path[:-1]:
            value = cursor.get(key)
            if value is None:
                value = {}
                cursor[key] = value
            if not isinstance(value, dict):
                raise ValueError(f"配置路径冲突: {'.'.join(path)}")
            cursor = value

        leaf_key = path[-1]
        leaf = cursor.get(leaf_key)
        if leaf is None:
            leaf = []
            cursor[leaf_key] = leaf
        if not isinstance(leaf, list):
            raise ValueError(f"配置路径冲突: {'.'.join(path)}")

        existing = {str(item.get(id_field)) for item in leaf if isinstance(item, dict) and id_field in item}
        if member_id in existing:
            return new_config, False
        leaf.append({id_field: member_id})
        return new_config, True
