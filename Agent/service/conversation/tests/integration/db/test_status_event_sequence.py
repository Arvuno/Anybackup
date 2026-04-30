import pytest
from sqlalchemy.ext.asyncio import async_sessionmaker, create_async_engine

from app.infrastructure.persistence.sqlalchemy.models import Base, ConversationStatusEventModel
from app.infrastructure.persistence.sqlalchemy.repositories import SqlAlchemyStatusEventRepository


@pytest.mark.asyncio
async def test_status_event_repository_returns_next_sequence() -> None:
    engine = create_async_engine("sqlite+aiosqlite:///:memory:")
    async with engine.begin() as connection:
        await connection.run_sync(Base.metadata.create_all)

    session_factory = async_sessionmaker(engine, expire_on_commit=False)
    async with session_factory() as session:
        session.add(
            ConversationStatusEventModel(
                f_status_event_id=1,
                f_conversation_id=10,
                f_sequence=1,
                f_event_type="message.created",
                f_event_version="v1",
                f_message_id=None,
                f_payload={},
                f_visible_to_user=True,
                f_trace_id=None,
                f_correlation_id=None,
                f_created_time=1_800_000_000_000,
                f_updated_time=1_800_000_000_000,
            )
        )
        await session.commit()

    async with session_factory() as session:
        repository = SqlAlchemyStatusEventRepository(session)

        assert await repository.next_sequence(conversation_id=10) == 2
        assert await repository.next_sequence(conversation_id=20) == 1

    await engine.dispose()
