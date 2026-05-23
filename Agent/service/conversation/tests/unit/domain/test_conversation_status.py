import pytest

from app.domain.conversation import Conversation, ConversationStatus
from app.domain.shared.errors import DomainError, ErrorReason


def test_conversation_status_allows_defined_transitions() -> None:
    conversation = Conversation(
        conversation_id=1,
        status=ConversationStatus.CREATED,
    )

    conversation.transition_to(ConversationStatus.ACTIVE)
    conversation.transition_to(ConversationStatus.PAUSED)
    conversation.transition_to(ConversationStatus.ACTIVE)
    conversation.transition_to(ConversationStatus.ARCHIVED)
    conversation.transition_to(ConversationStatus.ACTIVE)
    conversation.transition_to(ConversationStatus.ARCHIVED)
    conversation.transition_to(ConversationStatus.EXPIRED)

    assert conversation.status is ConversationStatus.EXPIRED


def test_conversation_status_rejects_undefined_transition() -> None:
    conversation = Conversation(
        conversation_id=1,
        status=ConversationStatus.CREATED,
    )

    with pytest.raises(DomainError) as exc_info:
        conversation.transition_to(ConversationStatus.EXPIRED)

    assert exc_info.value.reason is ErrorReason.INVALID_STATUS_TRANSITION
