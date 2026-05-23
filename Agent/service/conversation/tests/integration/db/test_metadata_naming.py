from sqlalchemy import BigInteger, UniqueConstraint

from app.infrastructure.persistence.sqlalchemy.models import Base

BUSINESS_TABLES = {
    "t_conversation",
    "t_conversation_message",
    "t_conversation_status_event",
    "t_conversation_mq_outbox",
}


def test_business_table_names_use_required_prefix() -> None:
    assert BUSINESS_TABLES.issubset(set(Base.metadata.tables))

    for table_name in BUSINESS_TABLES:
        assert table_name.startswith("t_")


def test_business_columns_use_required_prefix_and_time_columns() -> None:
    for table_name in BUSINESS_TABLES:
        table = Base.metadata.tables[table_name]
        column_names = {column.name for column in table.columns}

        assert "f_created_time" in column_names
        assert "f_updated_time" in column_names
        assert all(column_name.startswith("f_") for column_name in column_names)


def test_primary_keys_are_int64_bigint_columns() -> None:
    for table_name in BUSINESS_TABLES:
        table = Base.metadata.tables[table_name]
        primary_key_columns = list(table.primary_key.columns)

        assert len(primary_key_columns) == 1
        assert isinstance(primary_key_columns[0].type, BigInteger)


def test_status_event_has_conversation_sequence_unique_constraint() -> None:
    table = Base.metadata.tables["t_conversation_status_event"]

    unique_column_sets = {
        tuple(column.name for column in constraint.columns)
        for constraint in table.constraints
        if isinstance(constraint, UniqueConstraint)
    }

    assert ("f_conversation_id", "f_sequence") in unique_column_sets
