from pathlib import Path


WORKTREE_ROOT = Path(__file__).resolve().parents[2]


def _find_workspace_agents_file() -> Path:
    # 当前源码已迁移到 AnyBackup 子目录，协作约束仍沉淀在 v9_work 工作区根目录。
    for parent in WORKTREE_ROOT.parents:
        agents = parent / "AGENTS.md"
        if agents.exists():
            return agents
    raise FileNotFoundError("AGENTS.md was not found in workspace ancestors")


def _read(path: Path) -> str:
    return path.read_text(encoding="utf-8")


def test_repo_docs_lock_minimal_relay_contract():
    agents = _read(_find_workspace_agents_file())
    readme = _read(WORKTREE_ROOT / "README.md")
    requirements = _read(WORKTREE_ROOT / "docs" / "01-需求文档.md")
    tech = _read(WORKTREE_ROOT / "docs" / "02-技术选型文档.md")
    architecture = _read(WORKTREE_ROOT / "docs" / "03-架构设计文档.md")
    implementation = _read(WORKTREE_ROOT / "docs" / "04-实现设计文档.md")
    testing = _read(WORKTREE_ROOT / "docs" / "05-测试设计文档.md")
    evolution = _read(WORKTREE_ROOT / "docs" / "06-开发与演进说明.md")
    diagrams = _read(WORKTREE_ROOT / "docs" / "07-系统架构与交互图.md")
    interaction = _read(WORKTREE_ROOT / "docs" / "服务交互定义.md")
    mq_contract_detail = _read(WORKTREE_ROOT / "docs" / "核心智能体MQ消息契约.md")
    message_format = _read(WORKTREE_ROOT / "docs" / "消息格式定义.md")
    upstream = _read(WORKTREE_ROOT / "docs" / "contracts" / "upstream-mq-contract.md")
    downstream = _read(WORKTREE_ROOT / "docs" / "contracts" / "downstream-sdk-relay-contract.md")
    storage = _read(WORKTREE_ROOT / "docs" / "storage" / "minimal-schema.md")
    runbook = _read(WORKTREE_ROOT / "docs" / "runbooks" / "local-dev.md")
    combined = (
        agents
        + readme
        + requirements
        + tech
        + architecture
        + implementation
        + testing
        + evolution
        + diagrams
        + interaction
        + mq_contract_detail
        + message_format
        + upstream
        + downstream
        + storage
        + runbook
    )

    for text in (agents, readme):
        assert "最小中转版" in text
        assert "kweaver-sdk" in text
        assert "pip install kweaver-sdk" in text

    assert "建设目标" in requirements
    assert "Python" in tech
    assert "PostgreSQL" in tech
    assert "RabbitMQ" in tech
    assert "SQLAlchemy" in tech
    assert "kweaver-sdk" in tech
    assert "总体架构" in architecture
    assert "RelayService" in implementation
    assert "测试目标" in testing
    assert "开发共识" in evolution
    assert "总体架构图" in diagrams
    assert "消息发送主链路时序图" in diagrams
    assert "Core Agent Service" in interaction
    assert "conversation.message.sent.v1" in mq_contract_detail
    assert "conversation.message.events" in mq_contract_detail
    assert "core_agent.message.events" in mq_contract_detail
    assert "core_agent.run.accepted.v1" in mq_contract_detail
    assert "core_agent.run_status.events" in mq_contract_detail
    assert "conversation.core_agent.run_status" in mq_contract_detail
    assert "AG-UI 内容链路" in mq_contract_detail
    assert "重复触发 Decision Agent Session" in mq_contract_detail
    assert "payload.message_id" in mq_contract_detail
    assert "query" in mq_contract_detail
    assert "MQ 最小消息格式" in message_format
    assert "payload.message_id" in message_format

    assert "conversation.message.sent.v1" in upstream
    assert "conversation.message.cancel_requested.v1" in upstream
    assert "payload.message_id" in upstream
    assert "event_id" in upstream
    assert "不得重复触发 Decision Agent Session" in upstream
    assert "core_agent.run.accepted.v1" in downstream
    assert "core_agent.run.completed.v1" in downstream
    assert "core_agent.run.failed.v1" in downstream
    assert "core_agent.run_status.events" in downstream
    assert "conversation.core_agent.run_status" in downstream
    assert "message_id" in downstream
    assert "query" in downstream
    assert "只携带 `rabbitmq_url`" in downstream
    assert "conversation_id" in storage
    assert "decision_conversation_id" in storage
    assert "message_id" in storage
    assert "RabbitMQ" in runbook
    assert "taskId" not in combined
    assert "CoreAgentConversationOutput" not in combined
    assert "userView" not in combined
    assert "machineContext" not in combined
