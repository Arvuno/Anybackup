import logging

from sqlalchemy import create_engine, inspect, text
from sqlalchemy.orm import sessionmaker

from core_agent_service.infrastructure.db.base import Base
from core_agent_service.infrastructure.db import models as _db_models  # noqa: F401


logger = logging.getLogger(__name__)


def create_session_factory(database_url: str):
    engine = create_engine(database_url, future=True)
    return sessionmaker(bind=engine, expire_on_commit=False)


def init_db(database_url: str) -> None:
    engine = create_engine(database_url, future=True)
    try:
        Base.metadata.create_all(engine)
        _ensure_runtime_schema(engine)
    finally:
        engine.dispose()


def _ensure_runtime_schema(engine) -> None:
    # create_all 不会修改已存在的旧表；这里仅补齐最小中转原型当前需要的兼容列。
    inspector = inspect(engine)
    table_names = set(inspector.get_table_names())

    if "inbound_events" in table_names:
        column_names = {column["name"] for column in inspector.get_columns("inbound_events")}
        if "message_id" not in column_names:
            logger.info("database schema migration start", extra={"table": "inbound_events", "column": "message_id"})
            with engine.begin() as connection:
                connection.execute(text("ALTER TABLE inbound_events ADD COLUMN message_id VARCHAR(128)"))
                connection.execute(text("CREATE INDEX IF NOT EXISTS ix_inbound_events_message_id ON inbound_events (message_id)"))
            logger.info("database schema migration completed", extra={"table": "inbound_events", "column": "message_id"})

    if "conversation_mappings" in table_names:
        mapping_column_names = {column["name"] for column in inspector.get_columns("conversation_mappings")}
        if "decision_agent_id" not in mapping_column_names:
            logger.info("database schema migration start", extra={"table": "conversation_mappings", "column": "decision_agent_id"})
            with engine.begin() as connection:
                connection.execute(text("ALTER TABLE conversation_mappings ADD COLUMN decision_agent_id VARCHAR(128)"))
                connection.execute(text("CREATE INDEX IF NOT EXISTS ix_conversation_mappings_decision_agent_id ON conversation_mappings (decision_agent_id)"))
            logger.info("database schema migration completed", extra={"table": "conversation_mappings", "column": "decision_agent_id"})
