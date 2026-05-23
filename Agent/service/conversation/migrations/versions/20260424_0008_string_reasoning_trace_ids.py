from __future__ import annotations

import sqlalchemy as sa
from alembic import op

revision = "20260424_0008"
down_revision = "20260423_0007"
branch_labels = None
depends_on = None


def upgrade() -> None:
    with op.batch_alter_table("t_conversation_reasoning_trace") as batch_op:
        batch_op.alter_column(
            "f_reasoning_trace_id",
            existing_type=sa.BigInteger(),
            type_=sa.String(length=64),
            existing_nullable=False,
        )

    with op.batch_alter_table("t_conversation_candidate_selection") as batch_op:
        batch_op.alter_column(
            "f_reasoning_trace_id",
            existing_type=sa.BigInteger(),
            type_=sa.String(length=64),
            existing_nullable=False,
        )


def downgrade() -> None:
    with op.batch_alter_table("t_conversation_candidate_selection") as batch_op:
        batch_op.alter_column(
            "f_reasoning_trace_id",
            existing_type=sa.String(length=64),
            type_=sa.BigInteger(),
            existing_nullable=False,
        )

    with op.batch_alter_table("t_conversation_reasoning_trace") as batch_op:
        batch_op.alter_column(
            "f_reasoning_trace_id",
            existing_type=sa.String(length=64),
            type_=sa.BigInteger(),
            existing_nullable=False,
        )
