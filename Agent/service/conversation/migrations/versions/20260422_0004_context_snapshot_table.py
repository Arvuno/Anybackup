from __future__ import annotations

import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

revision = "20260422_0004"
down_revision = "20260422_0003"
branch_labels = None
depends_on = None


def json_document_type() -> sa.TypeEngine[object]:
    return sa.JSON().with_variant(postgresql.JSONB(), "postgresql")


def upgrade() -> None:
    op.create_table(
        "t_conversation_context_snapshot",
        sa.Column("f_context_snapshot_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_snapshot_version", sa.BigInteger(), nullable=False),
        sa.Column("f_short_summary", sa.Text(), nullable=False),
        sa.Column("f_structured_state", json_document_type(), nullable=False),
        sa.Column("f_last_message_id", sa.BigInteger(), nullable=True),
        sa.Column("f_status", sa.String(length=32), nullable=False),
        sa.Column("f_created_by", sa.String(length=80), nullable=False),
        sa.Column("f_trace_id", sa.String(length=128), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("f_conversation_id", "f_snapshot_version"),
    )
    op.create_index(
        "idx_t_conversation_context_snapshot_current",
        "t_conversation_context_snapshot",
        ["f_conversation_id", "f_status", "f_snapshot_version"],
    )


def downgrade() -> None:
    op.drop_index(
        "idx_t_conversation_context_snapshot_current",
        table_name="t_conversation_context_snapshot",
    )
    op.drop_table("t_conversation_context_snapshot")
