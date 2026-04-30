from __future__ import annotations

from deploy_orchestrator.cli_runner import CliRunner
from deploy_orchestrator.resources.agent import AgentExecutor
from deploy_orchestrator.resources.binding import BindingExecutor
from deploy_orchestrator.resources.bkn import BknExecutor
from deploy_orchestrator.resources.skill import SkillExecutor


def build_resource_executors(
    runner: CliRunner | None = None,
    workspace: str | None = None,
    business_domain: str | None = None,
) -> dict[str, object]:
    shared_runner = runner or CliRunner()
    return {
        "skill": SkillExecutor(shared_runner, business_domain=business_domain),
        "bkn": BknExecutor(shared_runner, business_domain=business_domain),
        "agent": AgentExecutor(shared_runner, workspace=workspace, business_domain=business_domain),
        "binding": BindingExecutor(shared_runner, business_domain=business_domain, workspace=workspace),
    }
