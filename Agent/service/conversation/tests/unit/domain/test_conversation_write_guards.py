import pytest

from app.domain.conversation import Conversation, ConversationStatus, InteractionStatus
from app.domain.shared.errors import DomainError, ErrorReason


@pytest.mark.parametrize(
    ("status", "expected_reason"),
    [
        (ConversationStatus.ARCHIVED, ErrorReason.CONVERSATION_ARCHIVED),
        (ConversationStatus.EXPIRED, ErrorReason.CONVERSATION_EXPIRED),
    ],
)
def test_archived_and_expired_conversations_reject_user_writes(
    status: ConversationStatus,
    expected_reason: ErrorReason,
) -> None:
    conversation = Conversation(conversation_id=1, status=status)

    with pytest.raises(DomainError) as exc_info:
        conversation.ensure_user_message_allowed()

    assert exc_info.value.reason is expected_reason


@pytest.mark.parametrize(
    "interaction_status",
    [
        InteractionStatus.THINKING,
        InteractionStatus.CLARIFYING,
        InteractionStatus.EXECUTING,
    ],
)
def test_in_progress_interaction_rejects_new_user_message(
    interaction_status: InteractionStatus,
) -> None:
    conversation = Conversation(
        conversation_id=1,
        status=ConversationStatus.ACTIVE,
        interaction_status=interaction_status,
    )

    with pytest.raises(DomainError) as exc_info:
        conversation.ensure_user_message_allowed()

    assert exc_info.value.reason is ErrorReason.CONVERSATION_BUSY


def test_idle_active_conversation_allows_user_message() -> None:
    conversation = Conversation(
        conversation_id=1,
        status=ConversationStatus.ACTIVE,
        interaction_status=InteractionStatus.IDLE,
    )

    conversation.ensure_user_message_allowed()


def test_active_conversation_allows_metadata_update_during_in_progress_interaction() -> None:
    conversation = Conversation(
        conversation_id=1,
        status=ConversationStatus.ACTIVE,
        interaction_status=InteractionStatus.EXECUTING,
    )

    conversation.ensure_metadata_update_allowed()


@pytest.mark.parametrize(
    ("status", "expected_reason"),
    [
        (ConversationStatus.ARCHIVED, ErrorReason.CONVERSATION_ARCHIVED),
        (ConversationStatus.EXPIRED, ErrorReason.CONVERSATION_EXPIRED),
    ],
)
def test_archived_and_expired_conversations_reject_agent_visible_content(
    status: ConversationStatus,
    expected_reason: ErrorReason,
) -> None:
    conversation = Conversation(conversation_id=1, status=status)

    with pytest.raises(DomainError) as exc_info:
        conversation.ensure_agent_visible_content_allowed()

    assert exc_info.value.reason is expected_reason


def test_child_object_must_belong_to_conversation() -> None:
    conversation = Conversation(conversation_id=1, status=ConversationStatus.ACTIVE)

    with pytest.raises(DomainError) as exc_info:
        conversation.ensure_child_belongs(child_conversation_id=2)

    assert exc_info.value.reason is ErrorReason.CHILD_CONVERSATION_MISMATCH


@pytest.mark.parametrize(
    ("status", "expected_reason"),
    [
        (ConversationStatus.ARCHIVED, ErrorReason.CONVERSATION_ARCHIVED),
        (ConversationStatus.EXPIRED, ErrorReason.CONVERSATION_EXPIRED),
    ],
)
def test_archived_and_expired_conversations_reject_metadata_update(
    status: ConversationStatus,
    expected_reason: ErrorReason,
) -> None:
    conversation = Conversation(conversation_id=1, status=status)

    with pytest.raises(DomainError) as exc_info:
        conversation.ensure_metadata_update_allowed()

    assert exc_info.value.reason is expected_reason
