from app.application.models.conversation import AuthenticatedUser, ConversationRecord
from app.domain.shared.errors import DomainError, ErrorReason


def ensure_conversation_owner(
    conversation: ConversationRecord,
    user: AuthenticatedUser,
) -> None:
    ensure_conversation_owner_id(conversation, user.user_id)


def ensure_conversation_owner_id(
    conversation: ConversationRecord,
    user_id: str,
) -> None:
    if conversation.owner_user_id != user_id:
        raise DomainError(ErrorReason.CONVERSATION_ACCESS_DENIED)
