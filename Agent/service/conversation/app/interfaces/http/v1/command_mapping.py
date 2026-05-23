from app.application.commands.conversation import CreateConversationCommand, SendUserMessageCommand
from app.interfaces.http.v1.schemas import (
    ClarificationResponseRequest,
    CreateConversationRequest,
    UserMessageRequest,
)


def build_create_conversation_command(
    *,
    request: CreateConversationRequest,
    idempotency_key_header: str | None,
    request_id: str | None,
) -> CreateConversationCommand:
    return CreateConversationCommand(
        initial_message_content=request.initial_message.content,
        initial_message_client_id=request.initial_message.client_message_id,
        title=request.title,
        scenario_binding=(
            request.scenario_binding.model_dump(mode="json")
            if request.scenario_binding is not None
            else None
        ),
        tags=tuple(request.tags),
        source=request.source,
        idempotency_key=(
            idempotency_key_header
            or request.idempotency_key
            or request.initial_message.idempotency_key
        ),
        request_id=request_id,
    )


def build_send_user_message_command(
    *,
    conversation_id: int,
    request: UserMessageRequest,
    idempotency_key_header: str | None,
    request_id: str | None,
) -> SendUserMessageCommand:
    return SendUserMessageCommand(
        conversation_id=conversation_id,
        content=request.content,
        client_message_id=request.client_message_id,
        parent_message_id=_optional_int(request.parent_message_id),
        idempotency_key=idempotency_key_header or request.idempotency_key,
        request_id=request_id,
    )


def build_send_clarification_response_command(
    *,
    conversation_id: int,
    request: ClarificationResponseRequest,
    idempotency_key_header: str | None,
    request_id: str | None,
) -> SendUserMessageCommand:
    return SendUserMessageCommand(
        conversation_id=conversation_id,
        content=_clarification_response_content(request),
        client_message_id=request.client_message_id,
        parent_message_id=_optional_int(request.message_id),
        idempotency_key=idempotency_key_header or request.idempotency_key,
        request_id=request_id,
        submission_type="clarification_response",
    )


def _optional_int(raw_value: str | None) -> int | None:
    if raw_value is None:
        return None
    return int(raw_value)


def _clarification_response_content(request: ClarificationResponseRequest) -> str:
    segments = ["clarification_response"]
    if request.clarification_id:
        segments.append(request.clarification_id.strip())
    if request.selected_value:
        segments.append(request.selected_value.strip())
    if request.free_text:
        segments.append(request.free_text.strip())
    return " ".join(segment for segment in segments if segment)
