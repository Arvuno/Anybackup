import logging
from collections.abc import Callable
from dataclasses import dataclass
from time import time
from typing import Any

from app.application.commands.reasoning import CandidateSelectionCommand
from app.application.models.conversation import AuthenticatedUser, MqOutboxRecord, Page
from app.application.models.writeback import (
    ConversationCandidateSelectionRecord,
    ConversationReasoningTraceRecord,
)
from app.application.ports.id_generator import IdGenerator
from app.application.ports.unit_of_work import UnitOfWork
from app.application.use_cases.access import (
    ensure_conversation_owner,
    ensure_conversation_owner_id,
)
from app.domain.conversation import ConversationStatus
from app.domain.shared.errors import DomainError, ErrorReason

logger = logging.getLogger(__name__)


@dataclass(frozen=True, slots=True)
class ReasoningTraceListResult:
    items: tuple[ConversationReasoningTraceRecord, ...]
    page: Page


@dataclass(frozen=True, slots=True)
class CandidateSelectionResult:
    selection: ConversationCandidateSelectionRecord
    idempotent: bool


class ListReasoningTracesHandler:
    def __init__(self, *, unit_of_work_factory: Callable[[], UnitOfWork]) -> None:
        self._unit_of_work_factory = unit_of_work_factory

    async def handle(
        self,
        conversation_id: int,
        user: AuthenticatedUser,
        *,
        limit: int,
    ) -> ReasoningTraceListResult | None:
        logger.info(
            "reasoning_trace_list_enter",
            extra={"conversation_id": conversation_id, "user_id": user.user_id},
        )
        async with self._unit_of_work_factory() as unit_of_work:
            conversation = await unit_of_work.conversations.get_record_by_id(conversation_id)
            if conversation is None:
                return None
            ensure_conversation_owner(conversation, user)
            traces = await unit_of_work.reasoning_traces.list_by_conversation(
                conversation_id,
                limit=limit,
            )
            return ReasoningTraceListResult(
                items=traces,
                page=Page(next_cursor=None, has_more=False, limit=limit),
            )


class ConfirmCandidateSelectionHandler:
    def __init__(
        self,
        *,
        unit_of_work_factory: Callable[[], UnitOfWork],
        id_generator: IdGenerator,
    ) -> None:
        self._unit_of_work_factory = unit_of_work_factory
        self._id_generator = id_generator

    async def handle(self, command: CandidateSelectionCommand) -> CandidateSelectionResult:
        logger.info(
            "candidate_selection_enter",
            extra={
                "conversation_id": command.conversation_id,
                "reasoning_trace_id": command.reasoning_trace_id,
                "candidate_option_id": command.candidate_option_id,
                "action": command.action,
            },
        )
        async with self._unit_of_work_factory() as unit_of_work:
            existing = await unit_of_work.candidate_selections.get_by_idempotency_key(
                command.idempotency_key
            )
            if existing is not None:
                conversation = await unit_of_work.conversations.get_record_by_id(
                    existing.conversation_id
                )
                if conversation is None:
                    raise DomainError(ErrorReason.CONVERSATION_NOT_FOUND)
                ensure_conversation_owner_id(conversation, command.user_id)
                return CandidateSelectionResult(selection=existing, idempotent=True)
            conversation = await unit_of_work.conversations.get_record_by_id(
                command.conversation_id
            )
            reasoning_trace = await unit_of_work.reasoning_traces.get_by_id(
                command.reasoning_trace_id
            )
            if conversation is None:
                raise DomainError(ErrorReason.CONVERSATION_NOT_FOUND)
            ensure_conversation_owner_id(conversation, command.user_id)
            if conversation.status is ConversationStatus.ARCHIVED:
                raise DomainError(ErrorReason.CONVERSATION_ARCHIVED)
            if conversation.status is ConversationStatus.EXPIRED:
                raise DomainError(ErrorReason.CONVERSATION_EXPIRED)
            if (
                reasoning_trace is None
                or reasoning_trace.conversation_id != command.conversation_id
            ):
                raise DomainError(ErrorReason.CHILD_CONVERSATION_MISMATCH)

            now_ms = _current_time_ms()
            selection_id = self._id_generator.next_id()
            outbox_id = self._id_generator.next_id()
            trace_id = command.request_id or f"candidate-selection-{selection_id}"
            correlation_id = f"conversation-{command.conversation_id}-selection-{selection_id}"
            selection = ConversationCandidateSelectionRecord(
                selection_id=selection_id,
                conversation_id=command.conversation_id,
                reasoning_trace_id=command.reasoning_trace_id,
                candidate_option_id=command.candidate_option_id,
                action=command.action,
                comment=command.comment,
                idempotency_key=command.idempotency_key,
                created_by_user_id=command.user_id,
                trace_id=trace_id,
                correlation_id=correlation_id,
                created_time=now_ms,
                updated_time=now_ms,
            )
            outbox = MqOutboxRecord(
                outbox_id=outbox_id,
                event_id=f"conversation.candidate_selection.created.{selection_id}",
                event_type="conversation.candidate_selection.created.v1",
                routing_key="conversation.candidate_selection.created.v1",
                conversation_id=command.conversation_id,
                message_id=None,
                payload={
                    "conversation_id": str(command.conversation_id),
                    "selection_id": str(selection_id),
                    "reasoning_trace_id": str(command.reasoning_trace_id),
                    "candidate_option_id": command.candidate_option_id,
                    "action": command.action,
                    "comment": command.comment,
                    "trace_id": trace_id,
                    "correlation_id": correlation_id,
                },
                status="pending",
                trace_id=trace_id,
                correlation_id=correlation_id,
                idempotency_key=command.idempotency_key,
                created_time=now_ms,
                updated_time=now_ms,
            )
            await unit_of_work.candidate_selections.add(selection)
            await unit_of_work.outbox.add(outbox)
            return CandidateSelectionResult(selection=selection, idempotent=False)


def visible_reasoning_payload(record: ConversationReasoningTraceRecord) -> dict[str, Any]:
    cleaned = _remove_forbidden_reasoning_fields(record.trace_payload)
    if isinstance(cleaned, dict):
        return cleaned
    return {}


def _remove_forbidden_reasoning_fields(value: Any, *, container_key: str | None = None) -> Any:
    forbidden = {"internal_reasoning", "system_prompt", "raw_tool_result"}
    if isinstance(value, dict):
        return {
            key: _remove_forbidden_reasoning_fields(item, container_key=key)
            for key, item in value.items()
            if key not in forbidden
            and not (key == "prompt" and container_key != "pending_confirmations")
        }
    if isinstance(value, list):
        return [
            _remove_forbidden_reasoning_fields(item, container_key=container_key)
            for item in value
        ]
    return value


def _current_time_ms() -> int:
    return int(time() * 1000)
