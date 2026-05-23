from __future__ import annotations

import sqlalchemy as sa
from alembic import op

revision = "20260422_0005"
down_revision = "20260422_0004"
branch_labels = None
depends_on = None


def upgrade() -> None:
    op.create_table(
        "t_conversation_candidate_selection",
        sa.Column("f_selection_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_reasoning_trace_id", sa.BigInteger(), nullable=False),
        sa.Column("f_candidate_option_id", sa.String(length=128), nullable=False),
        sa.Column("f_action", sa.String(length=16), nullable=False),
        sa.Column("f_comment", sa.String(length=500), nullable=True),
        sa.Column("f_idempotency_key", sa.String(length=128), nullable=False),
        sa.Column("f_created_by_user_id", sa.String(length=128), nullable=False),
        sa.Column("f_trace_id", sa.String(length=128), nullable=True),
        sa.Column("f_correlation_id", sa.String(length=128), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("f_idempotency_key"),
        sa.UniqueConstraint("f_conversation_id", "f_reasoning_trace_id", "f_candidate_option_id"),
    )
    op.create_index(
        "idx_t_conversation_candidate_selection_conversation",
        "t_conversation_candidate_selection",
        ["f_conversation_id", "f_created_time"],
    )


def downgrade() -> None:
    op.drop_index(
        "idx_t_conversation_candidate_selection_conversation",
        table_name="t_conversation_candidate_selection",
    )
    op.drop_table("t_conversation_candidate_selection")
