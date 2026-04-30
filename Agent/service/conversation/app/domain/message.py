from dataclasses import dataclass
from enum import StrEnum

from app.domain.shared.errors import DomainError, ErrorReason


class MessageStatus(StrEnum):
    RECEIVED = "received"
    PERSISTED = "persisted"
    PUBLISHED = "published"
    PROCESSING = "processing"
    RESPONDED = "responded"
    FAILED = "failed"


ALLOWED_MESSAGE_STATUS_TRANSITIONS: frozenset[tuple[MessageStatus, MessageStatus]] = frozenset(
    {
        (MessageStatus.RECEIVED, MessageStatus.PERSISTED),
        (MessageStatus.PERSISTED, MessageStatus.PUBLISHED),
        (MessageStatus.PUBLISHED, MessageStatus.PROCESSING),
        (MessageStatus.PERSISTED, MessageStatus.RESPONDED),
        (MessageStatus.PUBLISHED, MessageStatus.RESPONDED),
        (MessageStatus.PROCESSING, MessageStatus.RESPONDED),
        (MessageStatus.RECEIVED, MessageStatus.FAILED),
        (MessageStatus.PERSISTED, MessageStatus.FAILED),
        (MessageStatus.PUBLISHED, MessageStatus.FAILED),
        (MessageStatus.PROCESSING, MessageStatus.FAILED),
    }
)


@dataclass(slots=True)
class ConversationMessage:
    message_id: int
    conversation_id: int
    status: MessageStatus

    def transition_to(self, target: MessageStatus) -> None:
        if (self.status, target) not in ALLOWED_MESSAGE_STATUS_TRANSITIONS:
            raise DomainError(ErrorReason.INVALID_STATUS_TRANSITION)
        self.status = target
