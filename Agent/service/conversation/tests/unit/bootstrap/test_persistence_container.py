from dependency_injector import providers
from sqlalchemy.ext.asyncio import async_sessionmaker

from app.bootstrap.container import Container
from app.bootstrap.settings import Settings
from app.infrastructure.id_generator.snowflake import SnowflakeIdGenerator
from app.infrastructure.persistence.sqlalchemy.unit_of_work import SqlAlchemyUnitOfWork


def test_container_exposes_persistence_and_id_providers() -> None:
    container = Container()

    with container.settings.override(
        providers.Object(
            Settings(
                database_url="sqlite+aiosqlite:///:memory:",
                snowflake_node_id=3,
                snowflake_epoch_ms=1_735_689_600_000,
            )
        )
    ):
        id_generator = container.id_generator()
        session_factory = container.async_session_factory()
        unit_of_work = container.unit_of_work()

    assert isinstance(id_generator, SnowflakeIdGenerator)
    assert isinstance(session_factory, async_sessionmaker)
    assert isinstance(unit_of_work, SqlAlchemyUnitOfWork)
