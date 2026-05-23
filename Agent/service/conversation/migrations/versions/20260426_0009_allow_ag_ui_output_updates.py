from __future__ import annotations

from alembic import op
from sqlalchemy import inspect

revision = "20260426_0009"
down_revision = "20260424_0008"
branch_labels = None
depends_on = None

_TABLE_NAME = "t_conversation_writeback_idempotency"
_OUTPUT_UNIQUE_NAME = "uq_t_conversation_writeback_idempotency_output_id"


def upgrade() -> None:
    constraint_name = _find_output_unique_constraint_name()
    if constraint_name is None:
        return
    with op.batch_alter_table(_TABLE_NAME) as batch_op:
        batch_op.drop_constraint(constraint_name, type_="unique")


def downgrade() -> None:
    with op.batch_alter_table(_TABLE_NAME) as batch_op:
        batch_op.create_unique_constraint(_OUTPUT_UNIQUE_NAME, ["f_output_id"])


def _find_output_unique_constraint_name() -> str | None:
    inspector = inspect(op.get_bind())
    for constraint in inspector.get_unique_constraints(_TABLE_NAME):
        if constraint.get("column_names") == ["f_output_id"]:
            name = constraint.get("name")
            if isinstance(name, str) and name:
                return name
    return None
