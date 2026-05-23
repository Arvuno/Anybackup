from __future__ import annotations

from dataclasses import dataclass, field


@dataclass
class SkillSpec:
    name: str
    source: str
    register_mode: str = "zip"


@dataclass
class BknSpec:
    name: str
    source: str
    validate: bool = True
    branch: str = "main"


@dataclass
class AgentSpec:
    name: str
    profile: str
    llm_id: str
    knowledge_network: str | None = None
    skills: list[str] = field(default_factory=list)
    depends_on: list[str] = field(default_factory=list)
    config_template: str | None = None


@dataclass
class BindingSpec:
    name: str
    type: str
    agent: str
    member_agent: str | None = None
    skill: str | None = None
    config_path: list[str] = field(default_factory=list)
    id_field: str | None = None


@dataclass
class DeployConfig:
    project: str
    env: str = "dev"
    globals: dict[str, object] = field(default_factory=dict)
    skills: list[SkillSpec] = field(default_factory=list)
    bkns: list[BknSpec] = field(default_factory=list)
    agents: list[AgentSpec] = field(default_factory=list)
    bindings: list[BindingSpec] = field(default_factory=list)
