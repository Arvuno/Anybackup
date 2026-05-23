from types import TracebackType
from typing import Protocol

from app.application.ports.repositories import (
    CandidateSelectionRepository,
    ContextDeltaRepository,
    ContextSnapshotRepository,
    ConversationRepository,
    MessageRepository,
    MqOutboxRepository,
    ReasoningTraceRepository,
    StatusEventRepository,
    WritebackIdempotencyRepository,
)


class UnitOfWork(Protocol):
    conversations: ConversationRepository
    messages: MessageRepository
    status_events: StatusEventRepository
    outbox: MqOutboxRepository
    writebacks: WritebackIdempotencyRepository
    context_deltas: ContextDeltaRepository
    context_snapshots: ContextSnapshotRepository
    reasoning_traces: ReasoningTraceRepository
    candidate_selections: CandidateSelectionRepository

    async def __aenter__(self) -> "UnitOfWork":
        raise NotImplementedError

    async def __aexit__(
        self,
        exc_type: type[BaseException] | None,
        exc: BaseException | None,
        traceback: TracebackType | None,
    ) -> None:
        raise NotImplementedError

    async def commit(self) -> None:
        raise NotImplementedError

    async def rollback(self) -> None:
        raise NotImplementedError
