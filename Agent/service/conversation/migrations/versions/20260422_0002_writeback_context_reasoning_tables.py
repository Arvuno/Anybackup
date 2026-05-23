from __future__ import annotations

import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

revision = "20260422_0002"
down_revision = "20260421_0001"
branch_labels = None
depends_on = None


def json_document_type() -> sa.TypeEngine[object]:
    return sa.JSON().with_variant(postgresql.JSONB(), "postgresql")


def upgrade() -> None:
    op.create_table(
        "t_conversation_writeback_idempotency",
        sa.Column("f_writeback_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_idempotency_key", sa.String(length=128), nullable=False),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_output_id", sa.String(length=128), nullable=True),
        sa.Column("f_stream_id", sa.String(length=128), nullable=True),
        sa.Column("f_request_hash", sa.String(length=128), nullable=False),
        sa.Column("f_result_status", sa.String(length=32), nullable=False),
        sa.Column("f_result_message_id", sa.BigInteger(), nullable=True),
        sa.Column("f_reject_code", sa.String(length=80), nullable=True),
        sa.Column("f_reject_reason", sa.String(length=300), nullable=True),
        sa.Column("f_trace_id", sa.String(length=128), nullable=True),
        sa.Column("f_correlation_id", sa.String(length=128), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("f_idempotency_key"),
        sa.UniqueConstraint(
            "f_output_id",
            name="uq_t_conversation_writeback_idempotency_output_id",
        ),
    )
    op.create_index(
        "idx_t_conversation_writeback_conversation",
        "t_conversation_writeback_idempotency",
        ["f_conversation_id", "f_created_time"],
    )

    op.create_table(
        "t_conversation_context_delta",
        sa.Column("f_context_delta_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_turn_id", sa.BigInteger(), nullable=True),
        sa.Column("f_source_message_id", sa.BigInteger(), nullable=True),
        sa.Column("f_base_snapshot_version", sa.BigInteger(), nullable=True),
        sa.Column("f_delta_payload", json_document_type(), nullable=False),
        sa.Column("f_merge_status", sa.String(length=32), nullable=False),
        sa.Column("f_created_by_agent", sa.String(length=120), nullable=True),
        sa.Column("f_trace_id", sa.String(length=128), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
    )
    op.create_index(
        "idx_t_conversation_context_delta_merge",
        "t_conversation_context_delta",
        ["f_merge_status", "f_created_time"],
    )
    op.create_index(
        "idx_t_conversation_context_delta_conversation",
        "t_conversation_context_delta",
        ["f_conversation_id", "f_created_time"],
    )

    op.create_table(
        "t_conversation_reasoning_trace",
        sa.Column("f_reasoning_trace_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_source_message_id", sa.BigInteger(), nullable=True),
        sa.Column("f_trace_payload", json_document_type(), nullable=False),
        sa.Column("f_core_agent_run_id", sa.String(length=128), nullable=True),
        sa.Column("f_created_by_agent", sa.String(length=120), nullable=True),
        sa.Column("f_trace_id", sa.String(length=128), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
    )
    op.create_index(
        "idx_t_conversation_reasoning_trace_conversation",
        "t_conversation_reasoning_trace",
        ["f_conversation_id", "f_created_time"],
    )


def downgrade() -> None:
    op.drop_index(
        "idx_t_conversation_reasoning_trace_conversation",
        table_name="t_conversation_reasoning_trace",
    )
    op.drop_table("t_conversation_reasoning_trace")
    op.drop_index(
        "idx_t_conversation_context_delta_conversation",
        table_name="t_conversation_context_delta",
    )
    op.drop_index(
        "idx_t_conversation_context_delta_merge",
        table_name="t_conversation_context_delta",
    )
    op.drop_table("t_conversation_context_delta")
    op.drop_index(
        "idx_t_conversation_writeback_conversation",
        table_name="t_conversation_writeback_idempotency",
    )
    op.drop_table("t_conversation_writeback_idempotency")
