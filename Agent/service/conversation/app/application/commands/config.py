from dataclasses import dataclass, field


@dataclass(frozen=True, slots=True)
class UpdateConversationCommand:
    conversation_id: int
    title: str
    request_id: str | None = None


@dataclass(frozen=True, slots=True)
class CopyConversationConfigCommand:
    source_conversation_id: int
    title: str | None = None
    copy_tags: bool = True
    copy_scenario_binding: bool = True
    additional_tags: tuple[str, ...] = field(default_factory=tuple)
    idempotency_key: str | None = None
    request_id: str | None = None
