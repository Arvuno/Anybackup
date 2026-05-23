import pytest

from app.domain.message import ConversationMessage, MessageStatus
from app.domain.shared.errors import DomainError, ErrorReason


def test_message_status_allows_normal_agent_processing_flow() -> None:
    message = ConversationMessage(message_id=1, conversation_id=1, status=MessageStatus.RECEIVED)

    message.transition_to(MessageStatus.PERSISTED)
    message.transition_to(MessageStatus.PUBLISHED)
    message.transition_to(MessageStatus.PROCESSING)
    message.transition_to(MessageStatus.RESPONDED)

    assert message.status is MessageStatus.RESPONDED


@pytest.mark.parametrize(
    "source_status",
    [
        MessageStatus.PERSISTED,
        MessageStatus.PUBLISHED,
        MessageStatus.PROCESSING,
    ],
)
def test_message_status_allows_direct_response_from_visible_pre_terminal_states(
    source_status: MessageStatus,
) -> None:
    message = ConversationMessage(message_id=1, conversation_id=1, status=source_status)

    message.transition_to(MessageStatus.RESPONDED)

    assert message.status is MessageStatus.RESPONDED


def test_message_status_contract_does_not_expose_streaming() -> None:
    assert "streaming" not in {status.value for status in MessageStatus}


@pytest.mark.parametrize(
    "source_status",
    [
        MessageStatus.RECEIVED,
        MessageStatus.PERSISTED,
        MessageStatus.PUBLISHED,
        MessageStatus.PROCESSING,
    ],
)
def test_message_status_can_fail_from_incomplete_states(source_status: MessageStatus) -> None:
    message = ConversationMessage(message_id=1, conversation_id=1, status=source_status)

    message.transition_to(MessageStatus.FAILED)

    assert message.status is MessageStatus.FAILED


def test_message_status_rejects_undefined_transition() -> None:
    message = ConversationMessage(message_id=1, conversation_id=1, status=MessageStatus.RECEIVED)

    with pytest.raises(DomainError) as exc_info:
        message.transition_to(MessageStatus.RESPONDED)

    assert exc_info.value.reason is ErrorReason.INVALID_STATUS_TRANSITION
