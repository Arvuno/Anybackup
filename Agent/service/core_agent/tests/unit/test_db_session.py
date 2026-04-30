from sqlalchemy import create_engine, inspect, text

from core_agent_service.infrastructure.db.session import init_db


def test_init_db_adds_message_id_to_existing_inbound_events_table(tmp_path):
    database_path = tmp_path / "runtime.db"
    database_url = f"sqlite+pysqlite:///{database_path}"
    engine = create_engine(database_url, future=True)
    with engine.begin() as connection:
        # 模拟线上已有旧表结构，验证启动初始化会补齐新增的 message_id 列。
        connection.execute(
            text(
                """
                CREATE TABLE inbound_events (
                    event_id VARCHAR(128) PRIMARY KEY,
                    event_type VARCHAR(128) NOT NULL,
                    occurred_at DATETIME NOT NULL,
                    conversation_id VARCHAR(128) NOT NULL,
                    content TEXT NOT NULL,
                    processing_status VARCHAR(32) NOT NULL,
                    error_message TEXT,
                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
                    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
                )
                """
            )
        )
    engine.dispose()

    init_db(database_url)

    inspector = inspect(create_engine(database_url, future=True))
    columns = {column["name"] for column in inspector.get_columns("inbound_events")}
    indexes = {index["name"] for index in inspector.get_indexes("inbound_events")}

    assert "message_id" in columns
    assert "ix_inbound_events_message_id" in indexes


def test_init_db_adds_decision_agent_id_to_existing_mapping_table(tmp_path):
    database_path = tmp_path / "runtime_mapping.db"
    database_url = f"sqlite+pysqlite:///{database_path}"
    engine = create_engine(database_url, future=True)
    with engine.begin() as connection:
        # 模拟线上已有旧表结构，验证启动初始化会补齐目标 Agent 映射列。
        connection.execute(
            text(
                """
                CREATE TABLE conversation_mappings (
                    id INTEGER PRIMARY KEY,
                    conversation_id VARCHAR(128) NOT NULL UNIQUE,
                    decision_conversation_id VARCHAR(128) NOT NULL,
                    mapping_status VARCHAR(32) NOT NULL,
                    created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
                    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
                )
                """
            )
        )
    engine.dispose()

    init_db(database_url)

    inspector = inspect(create_engine(database_url, future=True))
    columns = {column["name"] for column in inspector.get_columns("conversation_mappings")}
    indexes = {index["name"] for index in inspector.get_indexes("conversation_mappings")}

    assert "decision_agent_id" in columns
    assert "ix_conversation_mappings_decision_agent_id" in indexes
