import pytest
from sqlalchemy.ext.asyncio import async_sessionmaker, create_async_engine

from app.domain.conversation import Conversation, ConversationStatus
from app.infrastructure.persistence.sqlalchemy.models import Base
from app.infrastructure.persistence.sqlalchemy.unit_of_work import SqlAlchemyUnitOfWork


@pytest.mark.asyncio
async def test_unit_of_work_commits_conversation_repository_changes() -> None:
    engine = create_async_engine("sqlite+aiosqlite:///:memory:")
    async with engine.begin() as connection:
        await connection.run_sync(Base.metadata.create_all)

    session_factory = async_sessionmaker(engine, expire_on_commit=False)

    async with SqlAlchemyUnitOfWork(session_factory) as unit_of_work:
        await unit_of_work.conversations.add(
            Conversation(
                conversation_id=1,
                owner_user_id="user-001",
                title="DB restore",
                status=ConversationStatus.ACTIVE,
                last_active_time=1_800_000_000_000,
                created_time=1_800_000_000_000,
                updated_time=1_800_000_000_000,
            )
        )

    async with SqlAlchemyUnitOfWork(session_factory) as unit_of_work:
        loaded = await unit_of_work.conversations.get_by_id(1)

    await engine.dispose()

    assert loaded is not None
    assert loaded.title == "DB restore"
