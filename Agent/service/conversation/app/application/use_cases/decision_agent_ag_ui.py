import hashlib
import json
import logging
from collections.abc import Callable
from dataclasses import replace
from time import time
from typing import Any

from app.application.commands.agent_events import DecisionAgentAgUiEventCommand
from app.application.models.conversation import (
    ConversationMessageRecord,
    ConversationRecord,
    ConversationStatusEventRecord,
)
from app.application.models.writeback import (
    ConversationContextDeltaRecord,
    ConversationWritebackIdempotencyRecord,
    DecisionAgentAgUiEventResult,
)
from app.application.ports.id_generator import IdGenerator
from app.application.ports.unit_of_work import UnitOfWork
from app.domain.conversation import ConversationStatus, InteractionStatus
from app.domain.message import MessageStatus
from app.domain.shared.errors import DomainError, ErrorReason

logger = logging.getLogger(__name__)


class DecisionAgentAgUiEventHandler:
    def __init__(
        self,
        *,
        unit_of_work_factory: Callable[[], UnitOfWork],
        id_generator: IdGenerator,
    ) -> None:
        self._unit_of_work_factory = unit_of_work_factory
        self._id_generator = id_generator

    async def handle(
        self,
        command: DecisionAgentAgUiEventCommand,
    ) -> DecisionAgentAgUiEventResult:
        logger.info(
            "decision_agent_ag_ui_event_handle_enter",
            extra={
                "event_id": command.event_id,
                "conversation_id": command.conversation_id,
                "turn_id": command.turn_id,
                "message_id": command.message_id,
                "sequence": command.sequence,
            },
        )
        request_hash = calculate_ag_ui_event_hash(command)
        output_id = _ag_ui_output_id(
            command.conversation_id,
            command.turn_id,
            command.message_id,
            command.sequence,
        )

        async with self._unit_of_work_factory() as unit_of_work:
            existing = await unit_of_work.writebacks.get_by_key(command.event_id)
            if existing is not None:
                return await _load_idempotent_result(unit_of_work, existing)

            conversation = await unit_of_work.conversations.get_record_by_id(
                command.conversation_id
            )
            if conversation is None:
                raise DomainError(ErrorReason.CONVERSATION_NOT_FOUND)
            source_message = await unit_of_work.messages.get_by_id(command.turn_id)
            if (
                source_message is None
                or source_message.conversation_id != command.conversation_id
                or source_message.role != "user"
            ):
                raise DomainError(ErrorReason.CHILD_CONVERSATION_MISMATCH)
            target_message = await unit_of_work.messages.get_by_id(command.message_id)
            if (
                target_message is not None
                and (
                    target_message.conversation_id != command.conversation_id
                    or target_message.role != "assistant"
                    or target_message.turn_id != command.turn_id
                )
            ):
                raise DomainError(ErrorReason.CHILD_CONVERSATION_MISMATCH)

            reject_reason = _conversation_reject_reason(conversation)
            if reject_reason is not None:
                await _persist_rejection(
                    unit_of_work,
                    id_generator=self._id_generator,
                    command=command,
                    output_id=output_id,
                    request_hash=request_hash,
                    reject_reason=reject_reason,
                )
                await unit_of_work.commit()
                raise DomainError(reject_reason)

            if not command.ag_ui.strip():
                await _persist_rejection(
                    unit_of_work,
                    id_generator=self._id_generator,
                    command=command,
                    output_id=output_id,
                    request_hash=request_hash,
                    reject_reason=ErrorReason.INVALID_STATUS_TRANSITION,
                    reject_detail="payload.ag_ui must be a non-empty Markdown string",
                )
                await unit_of_work.commit()
                raise DomainError(ErrorReason.INVALID_STATUS_TRANSITION)

            sequence_record = await unit_of_work.writebacks.get_by_conversation_output_id(
                conversation_id=command.conversation_id,
                output_id=output_id,
            )
            max_sequence = await unit_of_work.writebacks.max_accepted_output_sequence(
                command.conversation_id,
                turn_id=command.turn_id,
                message_id=command.message_id,
            )
            is_update = sequence_record is not None
            if not is_update and command.sequence != max_sequence + 1:
                await _persist_rejection(
                    unit_of_work,
                    id_generator=self._id_generator,
                    command=command,
                    output_id=output_id,
                    request_hash=request_hash,
                    reject_reason=ErrorReason.CONVERSATION_WRITEBACK_STALE,
                )
                await unit_of_work.commit()
                raise DomainError(ErrorReason.CONVERSATION_WRITEBACK_STALE)

            now_ms = command.occurred_time or _current_time_ms()
            should_add_message = target_message is None
            if target_message is None:
                target_message = ConversationMessageRecord(
                    message_id=command.message_id,
                    conversation_id=command.conversation_id,
                    parent_message_id=None,
                    turn_id=command.turn_id,
                    role="assistant",
                    content_type="rich_content",
                    content=None,
                    status=MessageStatus.PERSISTED,
                    trace_id=command.trace_id,
                    correlation_id=command.correlation_id,
                    created_time=now_ms,
                    updated_time=now_ms,
                )
            rich_payload = _rich_payload(command)
            target_interaction_status = InteractionStatus.COMPLETED
            next_sequence = await unit_of_work.status_events.next_sequence(
                command.conversation_id
            )
            source_message_id = command.turn_id
            source_turn_id = command.turn_id
            if target_interaction_status in {
                InteractionStatus.COMPLETED,
                InteractionStatus.ERROR,
            }:
                next_sequence = await _complete_processing_user_message(
                    unit_of_work,
                    id_generator=self._id_generator,
                    command=command,
                    source_message_id=source_message_id,
                    now_ms=now_ms,
                    next_sequence=next_sequence,
                    target_interaction_status=target_interaction_status,
                )

            message = replace(
                target_message,
                content_type="rich_content",
                content=command.content,
                rich_payload=rich_payload,
                status=MessageStatus.RESPONDED,
                trace_id=command.trace_id or target_message.trace_id,
                correlation_id=command.correlation_id or target_message.correlation_id,
                updated_time=now_ms,
            )
            updated_conversation = replace(
                conversation,
                interaction_status=target_interaction_status,
                active_turn_id=(
                    None
                    if target_interaction_status
                    in {InteractionStatus.COMPLETED, InteractionStatus.ERROR}
                    else conversation.active_turn_id
                ),
                updated_time=now_ms,
                last_active_time=now_ms,
            )
            if should_add_message:
                await unit_of_work.messages.add(message)
            else:
                await unit_of_work.messages.update(message)
            await unit_of_work.conversations.update_record(updated_conversation)
            context_delta = _context_delta_record(
                id_generator=self._id_generator,
                command=command,
                turn_id=source_turn_id,
                now_ms=now_ms,
            )
            if context_delta is not None:
                await unit_of_work.context_deltas.add(context_delta)

            status_event = ConversationStatusEventRecord(
                status_event_id=self._id_generator.next_id(),
                conversation_id=command.conversation_id,
                message_id=message.message_id,
                turn_id=source_turn_id,
                event_type="rich_content.updated" if is_update else "rich_content.created",
                sequence=next_sequence,
                interaction_status=updated_conversation.interaction_status,
                message_status=message.status,
                title="Decision Agent Markdown content",
                detail=command.content,
                payload={
                    "decision_agent_event_id": command.event_id,
                    "decision_agent_event_type": command.event_type,
                    "source_service": command.source_service,
                    "turn_id": str(command.turn_id),
                    "ag_ui_sequence": command.sequence,
                    "message": _message_payload(message),
                    "rich_payload": rich_payload,
                },
                rich_payload=rich_payload,
                trace_id=command.trace_id,
                correlation_id=command.correlation_id,
                created_time=now_ms,
                updated_time=now_ms,
            )
            interaction_event = _interaction_status_event(
                id_generator=self._id_generator,
                command=command,
                turn_id=source_turn_id,
                sequence=next_sequence + 1,
                previous_conversation=conversation,
                updated_conversation=updated_conversation,
                message_status=message.status,
                now_ms=now_ms,
            )
            idempotency = ConversationWritebackIdempotencyRecord(
                writeback_id=self._id_generator.next_id(),
                idempotency_key=command.event_id,
                conversation_id=command.conversation_id,
                output_id=output_id,
                request_hash=request_hash,
                result_status="accepted",
                result_message_id=message.message_id,
                trace_id=command.trace_id,
                correlation_id=command.correlation_id,
                created_time=now_ms,
                updated_time=now_ms,
            )

            await unit_of_work.status_events.add(status_event)
            if interaction_event is not None:
                await unit_of_work.status_events.add(interaction_event)
            await unit_of_work.writebacks.add(idempotency)
            logger.info(
                "decision_agent_ag_ui_event_persisted",
                extra={
                    "event_id": command.event_id,
                    "message_id": command.message_id,
                    "status_event_id": (
                        interaction_event.status_event_id
                        if interaction_event is not None
                        else status_event.status_event_id
                    ),
                },
            )
            return DecisionAgentAgUiEventResult(
                result_status="accepted",
                idempotent=False,
                conversation=updated_conversation,
                message=message,
                status_event=interaction_event or status_event,
            )


async def _complete_processing_user_message(
    unit_of_work: UnitOfWork,
    *,
    id_generator: IdGenerator,
    command: DecisionAgentAgUiEventCommand,
    source_message_id: int | None,
    now_ms: int,
    next_sequence: int,
    target_interaction_status: InteractionStatus,
) -> int:
    if source_message_id is None:
        return next_sequence

    message = await unit_of_work.messages.get_by_id(source_message_id)
    if message is None or message.role != "user":
        return next_sequence
    if message.status not in {
        MessageStatus.PERSISTED,
        MessageStatus.PUBLISHED,
        MessageStatus.PROCESSING,
    }:
        return next_sequence

    updated_message = replace(
        message,
        status=MessageStatus.RESPONDED,
        updated_time=now_ms,
    )
    await unit_of_work.messages.update(updated_message)
    await unit_of_work.status_events.add(
        ConversationStatusEventRecord(
            status_event_id=id_generator.next_id(),
            conversation_id=command.conversation_id,
            message_id=updated_message.message_id,
            turn_id=updated_message.turn_id,
            event_type="message.updated",
            sequence=next_sequence,
            interaction_status=target_interaction_status,
            message_status=updated_message.status,
            title="Decision Agent response received",
            detail="Decision Agent returned user-visible Markdown content",
            payload={
                "decision_agent_event_id": command.event_id,
                "decision_agent_event_type": command.event_type,
                "message": _message_payload(updated_message),
            },
            trace_id=command.trace_id,
            correlation_id=command.correlation_id,
            created_time=now_ms,
            updated_time=now_ms,
        )
    )
    return next_sequence + 1


def calculate_ag_ui_event_hash(command: DecisionAgentAgUiEventCommand) -> str:
    body = {
        "event_id": command.event_id,
        "event_type": command.event_type,
        "source_service": command.source_service,
        "conversation_id": command.conversation_id,
        "turn_id": command.turn_id,
        "message_id": command.message_id,
        "content": command.content,
        "sequence": command.sequence,
        "ag_ui": command.ag_ui,
        "occurred_time": command.occurred_time,
    }
    encoded = json.dumps(body, ensure_ascii=False, sort_keys=True)
    return hashlib.sha256(encoded.encode("utf-8")).hexdigest()


async def _load_idempotent_result(
    unit_of_work: UnitOfWork,
    existing: ConversationWritebackIdempotencyRecord,
) -> DecisionAgentAgUiEventResult:
    if existing.result_status == "rejected":
        return DecisionAgentAgUiEventResult(
            result_status="rejected",
            idempotent=True,
            reject_code=existing.reject_code,
            reject_reason=existing.reject_reason,
        )
    conversation = await unit_of_work.conversations.get_record_by_id(existing.conversation_id)
    message = (
        await unit_of_work.messages.get_by_id(existing.result_message_id)
        if existing.result_message_id is not None
        else None
    )
    status_event = (
        await unit_of_work.status_events.get_first_by_message_id(message.message_id)
        if message is not None
        else None
    )
    return DecisionAgentAgUiEventResult(
        result_status=existing.result_status,
        idempotent=True,
        conversation=conversation,
        message=message,
        status_event=status_event,
    )


async def _persist_rejection(
    unit_of_work: UnitOfWork,
    *,
    id_generator: IdGenerator,
    command: DecisionAgentAgUiEventCommand,
    output_id: str,
    request_hash: str,
    reject_reason: ErrorReason,
    reject_detail: str | None = None,
) -> None:
    now_ms = command.occurred_time or _current_time_ms()
    await unit_of_work.writebacks.add(
        ConversationWritebackIdempotencyRecord(
            writeback_id=id_generator.next_id(),
            idempotency_key=command.event_id,
            conversation_id=command.conversation_id,
            output_id=output_id,
            request_hash=request_hash,
            result_status="rejected",
            reject_code=reject_reason.value,
            reject_reason=reject_detail or reject_reason.value,
            trace_id=command.trace_id,
            correlation_id=command.correlation_id,
            created_time=now_ms,
            updated_time=now_ms,
        )
    )


def _conversation_reject_reason(conversation: ConversationRecord) -> ErrorReason | None:
    if conversation.status is ConversationStatus.ARCHIVED:
        return ErrorReason.CONVERSATION_ARCHIVED
    if conversation.status is ConversationStatus.EXPIRED:
        return ErrorReason.CONVERSATION_EXPIRED
    return None


def _rich_payload(command: DecisionAgentAgUiEventCommand) -> dict[str, Any]:
    return {
        "payload_id": f"markdown:{command.event_id}",
        "schema_version": "conversation-rich-payload-2026-04",
        "content_summary": command.content,
        "render_fallback": {"type": "text", "text": command.content},
        "ag_ui": command.ag_ui,
    }


def _message_payload(message: ConversationMessageRecord) -> dict[str, Any]:
    return {
        "message_id": str(message.message_id),
        "conversation_id": str(message.conversation_id),
        "parent_message_id": (
            str(message.parent_message_id) if message.parent_message_id is not None else None
        ),
        "turn_id": str(message.turn_id) if message.turn_id is not None else None,
        "role": message.role,
        "content_type": message.content_type,
        "content": message.content,
        "status": message.status.value,
    }


def _context_delta_record(
    *,
    id_generator: IdGenerator,
    command: DecisionAgentAgUiEventCommand,
    turn_id: int | None,
    now_ms: int,
) -> ConversationContextDeltaRecord | None:
    payload = _context_delta_payload(command=command)
    if payload is None:
        return None
    return ConversationContextDeltaRecord(
        context_delta_id=id_generator.next_id(),
        conversation_id=command.conversation_id,
        turn_id=turn_id,
        source_message_id=command.turn_id,
        base_snapshot_version=None,
        delta_payload=payload,
        merge_status="pending",
        created_by_agent=command.source_service,
        trace_id=command.trace_id,
        created_time=now_ms,
        updated_time=now_ms,
    )


def _context_delta_payload(
    *,
    command: DecisionAgentAgUiEventCommand,
) -> dict[str, Any] | None:
    summary_delta = command.content.strip()
    if not summary_delta:
        return None
    return {"summary_delta": summary_delta}


def _ag_ui_output_id(conversation_id: int, turn_id: int, message_id: int, sequence: int) -> str:
    return f"markdown:{conversation_id}:{turn_id}:{message_id}:{sequence}"


def _interaction_status_event(
    *,
    id_generator: IdGenerator,
    command: DecisionAgentAgUiEventCommand,
    turn_id: int | None,
    sequence: int,
    previous_conversation: ConversationRecord,
    updated_conversation: ConversationRecord,
    message_status: MessageStatus,
    now_ms: int,
) -> ConversationStatusEventRecord | None:
    if (
        previous_conversation.interaction_status == updated_conversation.interaction_status
        and previous_conversation.active_turn_id == updated_conversation.active_turn_id
    ):
        return None
    return ConversationStatusEventRecord(
        status_event_id=id_generator.next_id(),
        conversation_id=command.conversation_id,
        message_id=None,
        turn_id=turn_id,
        event_type="interaction.status_changed",
        sequence=sequence,
        interaction_status=updated_conversation.interaction_status,
        message_status=message_status,
        title="Interaction status changed",
        detail=command.content,
        payload={
            "decision_agent_event_id": command.event_id,
            "decision_agent_event_type": command.event_type,
            "source_service": command.source_service,
            "turn_id": str(command.turn_id),
            "ag_ui_sequence": command.sequence,
            "active_turn_id": (
                str(updated_conversation.active_turn_id)
                if updated_conversation.active_turn_id is not None
                else None
            ),
        },
        trace_id=command.trace_id,
        correlation_id=command.correlation_id,
        created_time=now_ms,
        updated_time=now_ms,
    )


def _current_time_ms() -> int:
    return int(time() * 1000)
