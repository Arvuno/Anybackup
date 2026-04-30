import asyncio
import subprocess
from pathlib import Path

import pytest
from alembic import command
from alembic.config import Config
from sqlalchemy import inspect
from sqlalchemy.ext.asyncio import create_async_engine

pytest.importorskip("testcontainers.postgres")
from testcontainers.postgres import PostgresContainer

SERVICE_ROOT = Path(__file__).resolve().parents[3]


def test_alembic_upgrade_head_runs_against_postgres_container(monkeypatch) -> None:
    _skip_if_docker_is_unavailable()

    with PostgresContainer("postgres:16", driver="asyncpg") as postgres:
        database_url = postgres.get_connection_url()
        monkeypatch.setenv("DATABASE_URL", database_url)

        config = Config(str(SERVICE_ROOT / "alembic.ini"))
        command.upgrade(config, "head")

        table_names = asyncio.run(_load_table_names(database_url))

    assert "t_conversation" in table_names
    assert "t_conversation_status_event" in table_names


async def _load_table_names(database_url: str) -> set[str]:
    engine = create_async_engine(database_url)
    try:
        async with engine.connect() as connection:
            return await connection.run_sync(
                lambda sync_connection: set(inspect(sync_connection).get_table_names())
            )
    finally:
        await engine.dispose()


def _skip_if_docker_is_unavailable() -> None:
    result = subprocess.run(
        ["docker", "info"],
        check=False,
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )
    if result.returncode != 0:
        pytest.skip("Docker is not available for Testcontainers PostgreSQL")
