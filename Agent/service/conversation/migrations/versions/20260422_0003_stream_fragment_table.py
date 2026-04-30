from __future__ import annotations

import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

revision = "20260422_0003"
down_revision = "20260422_0002"
branch_labels = None
depends_on = None


def json_document_type() -> sa.TypeEngine[object]:
    return sa.JSON().with_variant(postgresql.JSONB(), "postgresql")


def upgrade() -> None:
    op.create_table(
        "t_conversation_stream_fragment",
        sa.Column("f_stream_fragment_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_stream_id", sa.String(length=128), nullable=False),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_message_id", sa.BigInteger(), nullable=True),
        sa.Column("f_phase", sa.String(length=16), nullable=False),
        sa.Column("f_sequence", sa.BigInteger(), nullable=False),
        sa.Column("f_payload", json_document_type(), nullable=False),
        sa.Column("f_idempotency_key", sa.String(length=128), nullable=True),
        sa.Column("f_trace_id", sa.String(length=128), nullable=True),
        sa.Column("f_correlation_id", sa.String(length=128), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("f_stream_id", "f_sequence"),
        sa.UniqueConstraint("f_idempotency_key"),
    )
    op.create_index(
        "idx_t_conversation_stream_fragment_conversation",
        "t_conversation_stream_fragment",
        ["f_conversation_id", "f_created_time"],
    )


def downgrade() -> None:
    op.drop_index(
        "idx_t_conversation_stream_fragment_conversation",
        table_name="t_conversation_stream_fragment",
    )
    op.drop_table("t_conversation_stream_fragment")
