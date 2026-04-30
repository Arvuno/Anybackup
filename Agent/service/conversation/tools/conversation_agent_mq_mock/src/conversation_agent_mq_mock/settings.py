from __future__ import annotations

import os
from dataclasses import dataclass

from conversation_agent_mq_mock.messages import (
    DEFAULT_AG_UI_EXCHANGE,
    DEFAULT_AG_UI_ROUTING_KEY,
    DEFAULT_CORE_STATUS_QUEUE,
)


@dataclass(frozen=True, slots=True)
class Settings:
    rabbitmq_url: str
    conversation_exchange: str = "conversation.agent.events"
    conversation_routing_key: str = "conversation.message.sent.v1"
    mock_input_queue: str = "core-agent-mock.message.events"
    core_status_queue: str = DEFAULT_CORE_STATUS_QUEUE
    ag_ui_exchange: str = DEFAULT_AG_UI_EXCHANGE
    ag_ui_routing_key: str = DEFAULT_AG_UI_ROUTING_KEY
    ag_ui_queue: str = "conversation.decision_agent.ag_ui"
    delay_ms: int = 800
    prefetch_count: int = 10
    worker_id: str = "conversation-agent-mq-mock"

    @classmethod
    def from_env(cls) -> Settings:
        return cls(
            rabbitmq_url=_required_env("RABBITMQ_URL"),
            conversation_exchange=os.getenv(
                "CONVERSATION_EXCHANGE",
                "conversation.agent.events",
            ),
            conversation_routing_key=os.getenv(
                "CONVERSATION_ROUTING_KEY",
                "conversation.message.sent.v1",
            ),
            mock_input_queue=os.getenv("MOCK_INPUT_QUEUE", "core-agent-mock.message.events"),
            core_status_queue=os.getenv("CORE_STATUS_QUEUE", DEFAULT_CORE_STATUS_QUEUE),
            ag_ui_exchange=os.getenv("AGUI_EXCHANGE", DEFAULT_AG_UI_EXCHANGE),
            ag_ui_routing_key=os.getenv("AGUI_ROUTING_KEY", DEFAULT_AG_UI_ROUTING_KEY),
            ag_ui_queue=os.getenv("AGUI_QUEUE", "conversation.decision_agent.ag_ui"),
            delay_ms=int(os.getenv("MOCK_DELAY_MS", "800")),
            prefetch_count=int(os.getenv("MOCK_PREFETCH_COUNT", "10")),
            worker_id=os.getenv("MOCK_WORKER_ID", "conversation-agent-mq-mock"),
        )


def _required_env(name: str) -> str:
    value = os.getenv(name)
    if not value:
        raise RuntimeError(f"{name} is required")
    return value
