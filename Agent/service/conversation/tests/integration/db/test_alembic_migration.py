from pathlib import Path

from alembic import command
from alembic.config import Config
from sqlalchemy import create_engine, inspect

SERVICE_ROOT = Path(__file__).resolve().parents[3]


def test_alembic_upgrade_head_creates_core_tables(tmp_path, monkeypatch) -> None:
    database_file = tmp_path / "conversation_service.sqlite"
    async_url = f"sqlite+aiosqlite:///{database_file}"
    sync_url = f"sqlite:///{database_file}"
    monkeypatch.setenv("DATABASE_URL", async_url)

    config = Config(str(SERVICE_ROOT / "alembic.ini"))
    command.upgrade(config, "head")

    engine = create_engine(sync_url)
    try:
        inspector = inspect(engine)
        table_names = set(inspector.get_table_names())
        writeback_uniques = inspector.get_unique_constraints(
            "t_conversation_writeback_idempotency"
        )
    finally:
        engine.dispose()

    assert {
        "alembic_version",
        "t_conversation",
        "t_conversation_message",
        "t_conversation_status_event",
        "t_conversation_mq_outbox",
    }.issubset(table_names)
    assert not any(
        constraint.get("column_names") == ["f_output_id"]
        for constraint in writeback_uniques
    )
