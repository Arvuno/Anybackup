import pytest

from app.domain.conversation import Conversation, ConversationStatus, InteractionStatus
from app.domain.shared.errors import DomainError, ErrorReason


def test_interaction_status_allows_defined_transitions() -> None:
    conversation = Conversation(
        conversation_id=1,
        status=ConversationStatus.ACTIVE,
        interaction_status=InteractionStatus.IDLE,
    )

    conversation.transition_interaction_to(InteractionStatus.THINKING)
    conversation.transition_interaction_to(InteractionStatus.CLARIFYING)
    conversation.transition_interaction_to(InteractionStatus.THINKING)
    conversation.transition_interaction_to(InteractionStatus.EXECUTING)
    conversation.transition_interaction_to(InteractionStatus.COMPLETED)
    conversation.transition_interaction_to(InteractionStatus.IDLE)

    assert conversation.interaction_status is InteractionStatus.IDLE


def test_interaction_status_contract_does_not_expose_streaming() -> None:
    assert "streaming" not in {status.value for status in InteractionStatus}


def test_interaction_status_rejects_undefined_transition() -> None:
    conversation = Conversation(
        conversation_id=1,
        status=ConversationStatus.ACTIVE,
        interaction_status=InteractionStatus.IDLE,
    )

    with pytest.raises(DomainError) as exc_info:
        conversation.transition_interaction_to(InteractionStatus.EXECUTING)

    assert exc_info.value.reason is ErrorReason.INVALID_STATUS_TRANSITION
