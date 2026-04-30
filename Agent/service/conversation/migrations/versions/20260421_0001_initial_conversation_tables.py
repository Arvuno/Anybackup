from __future__ import annotations

import sqlalchemy as sa
from alembic import op
from sqlalchemy.dialects import postgresql

revision = "20260421_0001"
down_revision = None
branch_labels = None
depends_on = None


def json_document_type() -> sa.TypeEngine[object]:
    return sa.JSON().with_variant(postgresql.JSONB(), "postgresql")


def upgrade() -> None:
    op.create_table(
        "t_conversation",
        sa.Column("f_conversation_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_owner_user_id", sa.String(length=128), nullable=False),
        sa.Column("f_tenant_id", sa.String(length=128), nullable=True),
        sa.Column("f_title", sa.String(length=120), nullable=False),
        sa.Column("f_display_summary", sa.String(length=500), nullable=True),
        sa.Column("f_status", sa.String(length=32), nullable=False),
        sa.Column("f_interaction_status", sa.String(length=32), nullable=False),
        sa.Column("f_scenario_binding", json_document_type(), nullable=True),
        sa.Column("f_tags", json_document_type(), nullable=True),
        sa.Column("f_retention_policy", sa.String(length=80), nullable=False),
        sa.Column("f_legal_hold", sa.Boolean(), nullable=False),
        sa.Column("f_last_active_time", sa.BigInteger(), nullable=False),
        sa.Column("f_archived_time", sa.BigInteger(), nullable=True),
        sa.Column("f_archived_by", sa.String(length=16), nullable=True),
        sa.Column("f_archive_reason", sa.String(length=200), nullable=True),
        sa.Column("f_expires_time", sa.BigInteger(), nullable=True),
        sa.Column("f_expired_time", sa.BigInteger(), nullable=True),
        sa.Column("f_purge_after_time", sa.BigInteger(), nullable=True),
        sa.Column("f_purged_time", sa.BigInteger(), nullable=True),
        sa.Column("f_version", sa.BigInteger(), nullable=False),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
    )
    op.create_index(
        "idx_t_conversation_owner_status_active",
        "t_conversation",
        ["f_owner_user_id", "f_status", "f_last_active_time"],
    )
    op.create_index(
        "idx_t_conversation_retention_scan",
        "t_conversation",
        ["f_status", "f_legal_hold", "f_last_active_time"],
    )

    op.create_table(
        "t_conversation_message",
        sa.Column("f_message_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_parent_message_id", sa.BigInteger(), nullable=True),
        sa.Column("f_turn_id", sa.BigInteger(), nullable=True),
        sa.Column("f_role", sa.String(length=24), nullable=False),
        sa.Column("f_content_type", sa.String(length=32), nullable=False),
        sa.Column("f_content", sa.Text(), nullable=True),
        sa.Column("f_rich_payload", json_document_type(), nullable=True),
        sa.Column("f_status", sa.String(length=32), nullable=False),
        sa.Column("f_client_message_id", sa.String(length=80), nullable=True),
        sa.Column("f_idempotency_key", sa.String(length=128), nullable=True),
        sa.Column("f_trace_id", sa.String(length=128), nullable=True),
        sa.Column("f_correlation_id", sa.String(length=128), nullable=True),
        sa.Column("f_error_code", sa.String(length=80), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("f_conversation_id", "f_client_message_id"),
        sa.UniqueConstraint("f_conversation_id", "f_idempotency_key"),
    )
    op.create_index(
        "idx_t_conversation_message_page",
        "t_conversation_message",
        ["f_conversation_id", "f_created_time"],
    )

    op.create_table(
        "t_conversation_status_event",
        sa.Column("f_status_event_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_sequence", sa.BigInteger(), nullable=False),
        sa.Column("f_event_type", sa.String(length=80), nullable=False),
        sa.Column("f_event_version", sa.String(length=16), nullable=False),
        sa.Column("f_message_id", sa.BigInteger(), nullable=True),
        sa.Column("f_payload", json_document_type(), nullable=False),
        sa.Column("f_visible_to_user", sa.Boolean(), nullable=False),
        sa.Column("f_trace_id", sa.String(length=128), nullable=True),
        sa.Column("f_correlation_id", sa.String(length=128), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("f_conversation_id", "f_sequence"),
    )
    op.create_index(
        "idx_t_conversation_status_event_cursor",
        "t_conversation_status_event",
        ["f_conversation_id", "f_sequence"],
    )

    op.create_table(
        "t_conversation_mq_outbox",
        sa.Column("f_outbox_id", sa.BigInteger(), primary_key=True),
        sa.Column("f_event_id", sa.String(length=128), nullable=False),
        sa.Column("f_event_type", sa.String(length=120), nullable=False),
        sa.Column("f_routing_key", sa.String(length=160), nullable=False),
        sa.Column("f_conversation_id", sa.BigInteger(), nullable=False),
        sa.Column("f_message_id", sa.BigInteger(), nullable=True),
        sa.Column("f_payload", json_document_type(), nullable=False),
        sa.Column("f_status", sa.String(length=32), nullable=False),
        sa.Column("f_attempt_count", sa.BigInteger(), nullable=False),
        sa.Column("f_next_retry_time", sa.BigInteger(), nullable=True),
        sa.Column("f_last_error_code", sa.String(length=80), nullable=True),
        sa.Column("f_trace_id", sa.String(length=128), nullable=False),
        sa.Column("f_correlation_id", sa.String(length=128), nullable=False),
        sa.Column("f_idempotency_key", sa.String(length=128), nullable=True),
        sa.Column("f_created_time", sa.BigInteger(), nullable=False),
        sa.Column("f_updated_time", sa.BigInteger(), nullable=False),
        sa.UniqueConstraint("f_event_id"),
        sa.UniqueConstraint("f_idempotency_key"),
    )
    op.create_index(
        "idx_t_conversation_mq_outbox_retry",
        "t_conversation_mq_outbox",
        ["f_status", "f_next_retry_time"],
    )


def downgrade() -> None:
    op.drop_index("idx_t_conversation_mq_outbox_retry", table_name="t_conversation_mq_outbox")
    op.drop_table("t_conversation_mq_outbox")
    op.drop_index(
        "idx_t_conversation_status_event_cursor",
        table_name="t_conversation_status_event",
    )
    op.drop_table("t_conversation_status_event")
    op.drop_index("idx_t_conversation_message_page", table_name="t_conversation_message")
    op.drop_table("t_conversation_message")
    op.drop_index("idx_t_conversation_retention_scan", table_name="t_conversation")
    op.drop_index("idx_t_conversation_owner_status_active", table_name="t_conversation")
    op.drop_table("t_conversation")
