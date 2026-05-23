from __future__ import annotations

import json
import os
import subprocess
import sys
from pathlib import Path

import pytest

SCRIPT_DIR = Path(__file__).resolve().parents[5] / "skills" / "ag-ui-response" / "scripts"
if str(SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(SCRIPT_DIR))

REPO_ROOT = Path(__file__).resolve().parents[7]
SKILL_ROOT = Path(__file__).resolve().parents[5] / "skills" / "ag-ui-response"
SCHEMAS_DIR = SKILL_ROOT / "references" / "schemas"
REQUIREMENTS_FILE = SKILL_ROOT / "requirements.txt"
PROMPT_FILE = (
    REPO_ROOT / "docs" / "需求上下文" / "2026-04-会话管理" / "MySQL数据库恢复Agent提示词.md"
)

from ag_ui_mq_core import (  # noqa: E402
    DEFAULT_SNOWFLAKE_EPOCH_MS,
    DEFAULT_SNOWFLAKE_NODE_ID,
    ContractValidationError,
    generate_valid_message_from_markdown,
    validate_markdown_text,
    validate_message,
)


def test_generate_valid_message_from_markdown_outputs_mq_envelope() -> None:
    message = generate_valid_message_from_markdown(
        markdown="# 方案设计\n\n这里是 Markdown 内容。",
        conversation_id="100",
        turn_id="200",
        message_id="901",
        content="方案设计已生成。",
        sequence=1,
        event_id="evt-markdown-001",
        occurred_at="2026-04-28T10:00:00Z",
        now_ms=1_800_000_000_000,
    )

    assert message["event_id"] == "evt-markdown-001"
    assert message["event_type"] == "decision_agent.session.ag_ui_event"
    assert message["source_service"] == "decision_agent_session"
    assert message["occurred_at"] == "2026-04-28T10:00:00Z"
    assert message["payload"] == {
        "conversation_id": "100",
        "turn_id": "200",
        "message_id": "901",
        "content": "方案设计已生成。",
        "sequence": 1,
        "ag_ui": "# 方案设计\n\n这里是 Markdown 内容。",
    }
    assert isinstance(message["payload"]["ag_ui"], str)


def test_generate_markdown_message_can_generate_snowflake_message_id() -> None:
    now_ms = 1_800_000_001_000
    message = generate_valid_message_from_markdown(
        markdown="生成恢复候选方案。",
        conversation_id="100",
        turn_id="200",
        sequence=1,
        now_ms=now_ms,
    )

    expected_id = (
        ((now_ms - DEFAULT_SNOWFLAKE_EPOCH_MS) << 22)
        | (DEFAULT_SNOWFLAKE_NODE_ID << 12)
    )

    assert message["payload"]["message_id"] == str(expected_id)


def test_generate_markdown_message_increments_snowflake_sequence_within_process() -> None:
    now_ms = 1_800_000_002_000
    messages = [
        generate_valid_message_from_markdown(
            markdown=f"生成恢复候选方案 {index}。",
            conversation_id="100",
            turn_id="200",
            sequence=index + 1,
            now_ms=now_ms,
        )
        for index in range(3)
    ]

    expected_base = (
        ((now_ms - DEFAULT_SNOWFLAKE_EPOCH_MS) << 22)
        | (DEFAULT_SNOWFLAKE_NODE_ID << 12)
    )

    assert [message["payload"]["message_id"] for message in messages] == [
        str(expected_base),
        str(expected_base | 1),
        str(expected_base | 2),
    ]


def test_validate_message_accepts_payload_ag_ui_markdown_string() -> None:
    message = generate_valid_message_from_markdown(
        markdown="## 恢复方案\n\n- 使用最近恢复点。",
        conversation_id="100",
        turn_id="200",
        message_id="901",
        content="恢复方案已生成。",
        sequence=1,
        now_ms=1_800_000_003_000,
    )

    validated = validate_message(message)

    assert validated["payload"]["ag_ui"] == "## 恢复方案\n\n- 使用最近恢复点。"


@pytest.mark.parametrize("ag_ui", [None, "", "   ", {"version": "1.x", "events": []}])
def test_validate_message_rejects_missing_empty_or_object_ag_ui(ag_ui: object) -> None:
    message = {
        "event_id": "evt-markdown-invalid",
        "event_type": "decision_agent.session.ag_ui_event",
        "occurred_at": "2026-04-28T10:00:00Z",
        "source_service": "decision_agent_session",
        "payload": {
            "conversation_id": "100",
            "turn_id": "200",
            "message_id": "901",
            "content": "恢复方案已生成。",
            "sequence": 1,
            "ag_ui": ag_ui,
        },
    }

    with pytest.raises(ContractValidationError, match="message.payload.ag_ui"):
        validate_message(message)


@pytest.mark.parametrize(
    "markdown",
    [
        "完整内部推理链：第一步...",
        "系统提示词：你是...",
        "password=secret",
        "postgresql://user:pass@example/db",
        "-----BEGIN PRIVATE KEY-----",
        "原始工具参数：{'token': 'abc'}",
    ],
)
def test_validate_markdown_text_rejects_sensitive_visible_content(markdown: str) -> None:
    with pytest.raises(ContractValidationError):
        validate_markdown_text(markdown)


def test_cli_generate_accepts_markdown_argument() -> None:
    result = subprocess.run(
        [
            sys.executable,
            str(SCRIPT_DIR / "generate_ag_ui_mq_message.py"),
            "--markdown",
            "# 方案设计\n\n这里是 Markdown 内容。",
            "--conversation-id",
            "100",
            "--turn-id",
            "200",
            "--message-id",
            "901",
            "--content",
            "方案设计已生成。",
            "--sequence",
            "1",
            "--event-id",
            "evt-cli-markdown",
            "--now-ms",
            "1800000004000",
        ],
        check=True,
        capture_output=True,
        text=True,
        encoding="utf-8",
        env={**os.environ, "PYTHONIOENCODING": "utf-8", "PYTHONUTF8": "1"},
    )

    message = json.loads(result.stdout)

    assert message["event_id"] == "evt-cli-markdown"
    assert message["payload"]["ag_ui"] == "# 方案设计\n\n这里是 Markdown 内容。"


def test_message_schema_defines_payload_ag_ui_as_string() -> None:
    message_schema = json.loads(
        (SCHEMAS_DIR / "decision-agent-ag-ui-message.schema.json").read_text(
            encoding="utf-8"
        )
    )

    ag_ui_schema = message_schema["properties"]["payload"]["properties"]["ag_ui"]
    assert ag_ui_schema["type"] == "string"
    assert ag_ui_schema["minLength"] == 1


def test_ag_ui_response_docs_define_markdown_runtime_path() -> None:
    skill_doc = (SKILL_ROOT / "SKILL.md").read_text(encoding="utf-8")
    flow_doc = (SKILL_ROOT / "references" / "flow-and-updates.md").read_text(
        encoding="utf-8"
    )
    reference_docs = "\n".join(
        path.read_text(encoding="utf-8")
        for path in sorted((SKILL_ROOT / "references").glob("*.md"))
    )
    agent_metadata = (SKILL_ROOT / "agents" / "openai.yaml").read_text(encoding="utf-8")
    combined = f"{skill_doc}\n{flow_doc}\n{reference_docs}\n{agent_metadata}"

    assert "payload.ag_ui" in combined
    assert "Markdown" in combined
    assert "--markdown" in combined
    assert "ag-ui-design" not in combined
    assert "--draft-json" not in combined
    assert "draft JSON" not in combined
    assert "activities" not in combined
    assert "activity_updates" not in combined
    assert "当前进展摘要" not in combined
    assert "进展摘要" not in combined
    assert "工具状态" not in combined
    assert "工具调用内容" not in combined
    assert "每次业务工具调用前" not in combined
    assert "首次业务判断前" not in combined
    assert "thought" not in combined
    assert "tool_call" not in combined


def test_mysql_recovery_agent_prompt_uses_only_ag_ui_response_gate() -> None:
    prompt = PROMPT_FILE.read_text(encoding="utf-8")

    assert "ag-ui-response" in prompt
    assert "ag-ui-design" not in prompt
    assert "设计约束" not in prompt
    assert "draft JSON" not in prompt
    assert "--draft-json" not in prompt
    assert "activities" not in prompt
    assert "activity_updates" not in prompt
    assert "不得调用 KWeaver" in prompt
    assert "不得查询知识网络" in prompt
    assert "不得下发恢复任务" in prompt
    assert "不得执行可用性验证" in prompt
    assert "Markdown" in prompt


def test_skill_package_declares_runtime_dependencies_for_standalone_delivery() -> None:
    requirements = REQUIREMENTS_FILE.read_text(encoding="utf-8")

    assert "aio-pika" in requirements
