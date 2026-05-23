from dataclasses import dataclass
from typing import Any


@dataclass(frozen=True, slots=True)
class CreateConversationCommand:
    initial_message_content: str
    initial_message_client_id: str | None
    title: str | None
    scenario_binding: dict[str, Any] | None
    tags: tuple[str, ...]
    source: str
    idempotency_key: str | None
    request_id: str | None


@dataclass(frozen=True, slots=True)
class SendUserMessageCommand:
    conversation_id: int
    content: str
    client_message_id: str | None
    parent_message_id: int | None
    idempotency_key: str | None
    request_id: str | None
    submission_type: str = "user_message"
