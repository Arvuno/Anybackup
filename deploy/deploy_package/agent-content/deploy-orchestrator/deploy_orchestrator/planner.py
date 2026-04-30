from __future__ import annotations

from dataclasses import dataclass
from heapq import heappop, heappush
from pathlib import Path

from deploy_orchestrator.logger import get_logger
from deploy_orchestrator.renderer import extract_template_references_from_file
from deploy_orchestrator.schema import DeployConfig


@dataclass
class PlanItem:
    resource_kind: str
    resource_name: str


def build_execution_plan(config: DeployConfig, base_path: str | Path | None = None) -> list[PlanItem]:
    logger = get_logger(__name__)
    logger.info("Building execution plan for project=%s env=%s", config.project, config.env)
    items: list[PlanItem] = []
    all_items: list[PlanItem] = []
    all_items.extend(PlanItem("skill", item.name) for item in config.skills)
    all_items.extend(PlanItem("bkn", item.name) for item in config.bkns)
    all_items.extend(PlanItem("agent", item.name) for item in config.agents)
    all_items.extend(PlanItem("binding", item.name) for item in config.bindings)

    priority = {"skill": 0, "bkn": 1, "agent": 2, "binding": 3}
    item_map = {(item.resource_kind, item.resource_name): item for item in all_items}
    order_map = {(item.resource_kind, item.resource_name): index for index, item in enumerate(all_items)}
    edges: dict[tuple[str, str], set[tuple[str, str]]] = {key: set() for key in item_map}
    indegree: dict[tuple[str, str], int] = {key: 0 for key in item_map}

    def add_edge(src: tuple[str, str], dst: tuple[str, str]) -> None:
        if src not in item_map or dst not in item_map or src == dst:
            return
        if dst not in edges[src]:
            edges[src].add(dst)
            indegree[dst] += 1

    for agent in config.agents:
        agent_key = ("agent", agent.name)
        for dep in agent.depends_on:
            add_edge(("agent", dep), agent_key)
        if agent.knowledge_network:
            add_edge(("bkn", agent.knowledge_network), agent_key)
        for skill_name in agent.skills:
            add_edge(("skill", skill_name), agent_key)
        for kind, name in _extract_agent_template_dependencies(agent, base_path):
            add_edge((kind, name), agent_key)

    for binding in config.bindings:
        binding_key = ("binding", binding.name)
        add_edge(("agent", binding.agent), binding_key)
        if binding.member_agent:
            add_edge(("agent", binding.member_agent), binding_key)
        if binding.skill:
            add_edge(("skill", binding.skill), binding_key)

    heap: list[tuple[int, int, tuple[str, str]]] = []
    for key, count in indegree.items():
        if count == 0:
            heappush(heap, (priority[key[0]], order_map[key], key))

    while heap:
        _, _, key = heappop(heap)
        items.append(item_map[key])
        for neighbor in sorted(edges[key], key=lambda item: (priority[item[0]], order_map[item])):
            indegree[neighbor] -= 1
            if indegree[neighbor] == 0:
                heappush(heap, (priority[neighbor[0]], order_map[neighbor], neighbor))

    if len(items) != len(all_items):
        remaining = [f"{kind}/{name}" for (kind, name), count in indegree.items() if count > 0]
        logger.error("Detected dependency cycle: %s", " -> ".join(remaining))
        raise ValueError(f"检测到循环依赖: {' -> '.join(remaining)}")

    logger.info("Built execution plan with %s items", len(items))
    return items


def _extract_agent_template_dependencies(agent, base_path: str | Path | None) -> list[tuple[str, str]]:
    if not agent.config_template:
        return []
    template_path = Path(agent.config_template)
    if not template_path.is_absolute():
        template_path = Path(base_path or ".") / template_path
    references = extract_template_references_from_file(template_path)
    dependencies: list[tuple[str, str]] = []
    dependencies.extend(("skill", name) for name in sorted(references["skills"]))
    dependencies.extend(("bkn", name) for name in sorted(references["bkns"]))
    dependencies.extend(("agent", name) for name in sorted(references["agents"]))
    return dependencies
