from collections.abc import Callable
from types import TracebackType

from sqlalchemy.ext.asyncio import AsyncSession

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
from app.application.ports.unit_of_work import UnitOfWork
from app.infrastructure.persistence.sqlalchemy.repositories import (
    SqlAlchemyCandidateSelectionRepository,
    SqlAlchemyContextDeltaRepository,
    SqlAlchemyContextSnapshotRepository,
    SqlAlchemyConversationRepository,
    SqlAlchemyMessageRepository,
    SqlAlchemyMqOutboxRepository,
    SqlAlchemyReasoningTraceRepository,
    SqlAlchemyStatusEventRepository,
    SqlAlchemyWritebackIdempotencyRepository,
)


class SqlAlchemyUnitOfWork(UnitOfWork):
    def __init__(self, session_factory: Callable[[], AsyncSession]) -> None:
        self._session_factory = session_factory
        self._session: AsyncSession | None = None
        self.conversations: ConversationRepository
        self.messages: MessageRepository
        self.status_events: StatusEventRepository
        self.outbox: MqOutboxRepository
        self.writebacks: WritebackIdempotencyRepository
        self.context_deltas: ContextDeltaRepository
        self.context_snapshots: ContextSnapshotRepository
        self.reasoning_traces: ReasoningTraceRepository
        self.candidate_selections: CandidateSelectionRepository

    async def __aenter__(self) -> "SqlAlchemyUnitOfWork":
        self._session = self._session_factory()
        self.conversations = SqlAlchemyConversationRepository(self._session)
        self.messages = SqlAlchemyMessageRepository(self._session)
        self.status_events = SqlAlchemyStatusEventRepository(self._session)
        self.outbox = SqlAlchemyMqOutboxRepository(self._session)
        self.writebacks = SqlAlchemyWritebackIdempotencyRepository(self._session)
        self.context_deltas = SqlAlchemyContextDeltaRepository(self._session)
        self.context_snapshots = SqlAlchemyContextSnapshotRepository(self._session)
        self.reasoning_traces = SqlAlchemyReasoningTraceRepository(self._session)
        self.candidate_selections = SqlAlchemyCandidateSelectionRepository(self._session)
        return self

    async def __aexit__(
        self,
        exc_type: type[BaseException] | None,
        exc: BaseException | None,
        traceback: TracebackType | None,
    ) -> None:
        if self._session is None:
            return
        if exc_type is None:
            await self.commit()
        else:
            await self.rollback()
        await self._session.close()

    async def commit(self) -> None:
        if self._session is None:
            raise RuntimeError("unit of work is not active")
        await self._session.commit()

    async def rollback(self) -> None:
        if self._session is None:
            raise RuntimeError("unit of work is not active")
        await self._session.rollback()
