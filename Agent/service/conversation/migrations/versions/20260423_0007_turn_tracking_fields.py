from __future__ import annotations

import sqlalchemy as sa
from alembic import op

revision = "20260423_0007"
down_revision = "20260422_0006"
branch_labels = None
depends_on = None


def upgrade() -> None:
    with op.batch_alter_table("t_conversation") as batch_op:
        batch_op.add_column(sa.Column("f_active_turn_id", sa.BigInteger(), nullable=True))

    with op.batch_alter_table("t_conversation_status_event") as batch_op:
        batch_op.add_column(sa.Column("f_turn_id", sa.BigInteger(), nullable=True))


def downgrade() -> None:
    with op.batch_alter_table("t_conversation_status_event") as batch_op:
        batch_op.drop_column("f_turn_id")

    with op.batch_alter_table("t_conversation") as batch_op:
        batch_op.drop_column("f_active_turn_id")
