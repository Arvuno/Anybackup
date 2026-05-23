import asyncio
import json
from collections.abc import Iterator
from pathlib import Path

import pytest
from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import create_async_engine

from app.application.commands.agent_events import DecisionAgentAgUiEventCommand
from app.application.use_cases.decision_agent_ag_ui import DecisionAgentAgUiEventHandler
from app.domain.conversation import InteractionStatus
from app.domain.message import MessageStatus
from app.domain.shared.errors import DomainError, ErrorReason
from app.infrastructure.persistence.sqlalchemy.models import (
    Base,
    ConversationContextDeltaModel,
    ConversationMessageModel,
    ConversationModel,
    ConversationReasoningTraceModel,
    ConversationStatusEventModel,
    ConversationWritebackIdempotencyModel,
)
from app.infrastructure.persistence.sqlalchemy.session import create_async_session_factory
from app.infrastructure.persistence.sqlalchemy.unit_of_work import SqlAlchemyUnitOfWork
from app.interfaces.mq.decision_agent_ag_ui_consumer import DecisionAgentAgUiMessageConsumer


@pytest.fixture()
def database_url(tmp_path: Path) -> Iterator[str]:
    database_path = tmp_path / "conversation.db"
    url = f"sqlite+aiosqlite:///{database_path.as_posix()}"

    async def create_schema() -> None:
        engine = create_async_engine(url)
        async with engine.begin() as connection:
            await connection.run_sync(Base.metadata.create_all)
        await engine.dispose()

    asyncio.run(create_schema())
    yield url


@pytest.mark.asyncio
async def test_markdown_event_creates_assistant_message_and_closes_turn(
    database_url: str,
) -> None:
    await _insert_processing_conversation(database_url)
    handler = _handler(database_url)

    result = await handler.handle(_command(event_id="markdown-event-001", sequence=1))

    assert result.idempotent is False
    assert result.result_status == "accepted"
    assistant_messages = await _assistant_messages(database_url)
    assert len(assistant_messages) == 1
    assert assistant_messages[0].f_message_id == 901
    assert assistant_messages[0].f_parent_message_id is None
    assert assistant_messages[0].f_turn_id == 200
    assert assistant_messages[0].f_content_type == "rich_content"
    assert assistant_messages[0].f_content == "Generated two recovery candidates."
    assert assistant_messages[0].f_status == MessageStatus.RESPONDED.value
    assert assistant_messages[0].f_rich_payload == {
        "payload_id": "markdown:markdown-event-001",
        "schema_version": "conversation-rich-payload-2026-04",
        "content_summary": "Generated two recovery candidates.",
        "render_fallback": {"type": "text", "text": "Generated two recovery candidates."},
        "ag_ui": "# Recovery candidates\n\n- Candidate A\n- Candidate B",
    }
    conversation = await _get_conversation(database_url)
    assert conversation.f_interaction_status == InteractionStatus.COMPLETED.value
    assert conversation.f_active_turn_id is None
    user_message = await _message(database_url, 200)
    assert user_message.f_status == MessageStatus.RESPONDED.value
    event_types = await _event_types(database_url)
    assert event_types[-3:] == [
        "message.updated",
        "rich_content.created",
        "interaction.status_changed",
    ]
    assert await _count(database_url, ConversationReasoningTraceModel.f_reasoning_trace_id) == 0
    assert await _count(database_url, ConversationWritebackIdempotencyModel.f_writeback_id) == 1


@pytest.mark.asyncio
async def test_markdown_replaces_same_message_id_and_sequence_with_new_event_id(
    database_url: str,
) -> None:
    await _insert_processing_conversation(database_url)
    handler = _handler(database_url)

    await handler.handle(_command(event_id="markdown-event-create", sequence=1))
    await handler.handle(
        _command(
            event_id="markdown-event-replace",
            sequence=1,
            content="Updated recovery candidates.",
            markdown="## Updated candidates\n\nUse candidate C.",
        )
    )

    assistant_messages = await _assistant_messages(database_url)
    assert len(assistant_messages) == 1
    assert assistant_messages[0].f_message_id == 901
    assert assistant_messages[0].f_content == "Updated recovery candidates."
    assert assistant_messages[0].f_rich_payload["payload_id"] == "markdown:markdown-event-replace"
    assert assistant_messages[0].f_rich_payload["ag_ui"] == (
        "## Updated candidates\n\nUse candidate C."
    )
    event_types = await _event_types(database_url)
    assert "rich_content.updated" in event_types
    assert await _count(database_url, ConversationWritebackIdempotencyModel.f_writeback_id) == 2


@pytest.mark.asyncio
async def test_markdown_event_accepts_next_sequence_for_new_assistant_message(
    database_url: str,
) -> None:
    await _insert_processing_conversation(database_url)
    handler = _handler(database_url)

    await handler.handle(_command(event_id="markdown-event-one", sequence=1, message_id=901))
    await handler.handle(
        _command(
            event_id="markdown-event-two",
            sequence=1,
            message_id=902,
            content="Second Markdown response.",
            markdown="## Second response\n\nFollow-up content.",
        )
    )

    assistant_messages = await _assistant_messages(database_url)
    assert [message.f_message_id for message in assistant_messages] == [901, 902]
    assert [message.f_content_type for message in assistant_messages] == [
        "rich_content",
        "rich_content",
    ]
    assert await _count(database_url, ConversationWritebackIdempotencyModel.f_writeback_id) == 2


@pytest.mark.asyncio
async def test_duplicate_markdown_event_is_idempotent(database_url: str) -> None:
    await _insert_processing_conversation(database_url)
    handler = _handler(database_url)
    command = _command(event_id="markdown-event-duplicate", sequence=1)

    first = await handler.handle(command)
    second = await handler.handle(command)

    assert first.idempotent is False
    assert second.idempotent is True
    assert len(await _assistant_messages(database_url)) == 1
    assert await _count(database_url, ConversationWritebackIdempotencyModel.f_writeback_id) == 1


@pytest.mark.asyncio
async def test_markdown_event_rejects_sequence_gap(database_url: str) -> None:
    await _insert_processing_conversation(database_url)
    handler = _handler(database_url)

    with pytest.raises(DomainError) as exc_info:
        await handler.handle(_command(event_id="markdown-event-gap", sequence=2))

    assert exc_info.value.reason is ErrorReason.CONVERSATION_WRITEBACK_STALE
    assert len(await _assistant_messages(database_url)) == 0
    assert await _count(database_url, ConversationWritebackIdempotencyModel.f_writeback_id) == 1


@pytest.mark.asyncio
async def test_markdown_event_rejects_archived_conversation(database_url: str) -> None:
    await _insert_processing_conversation(database_url, conversation_status="archived")
    handler = _handler(database_url)

    with pytest.raises(DomainError) as exc_info:
        await handler.handle(_command(event_id="markdown-event-archived", sequence=1))

    assert exc_info.value.reason is ErrorReason.CONVERSATION_ARCHIVED
    assert len(await _assistant_messages(database_url)) == 0
    assert await _count(database_url, ConversationWritebackIdempotencyModel.f_writeback_id) == 1


@pytest.mark.asyncio
async def test_markdown_event_persists_context_delta_summary(database_url: str) -> None:
    await _insert_processing_conversation(database_url)
    handler = _handler(database_url)

    await handler.handle(_command(event_id="markdown-event-context", sequence=1))

    context_delta = await _latest_context_delta(database_url)
    assert context_delta.f_conversation_id == 100
    assert context_delta.f_turn_id == 200
    assert context_delta.f_source_message_id == 200
    assert context_delta.f_merge_status == "pending"
    assert context_delta.f_delta_payload == {
        "summary_delta": "Generated two recovery candidates."
    }


@pytest.mark.asyncio
async def test_consumer_acks_after_successful_markdown_processing(database_url: str) -> None:
    await _insert_processing_conversation(database_url)
    consumer = DecisionAgentAgUiMessageConsumer(handler=_handler(database_url))
    incoming = FakeIncomingMessage(_event_body(event_id="markdown-consumer-accepted", sequence=1))

    await consumer.process_message(incoming)

    assert incoming.ack_calls == 1
    assert incoming.reject_calls == []


@pytest.mark.asyncio
async def test_consumer_rejects_object_ag_ui_payload_to_dlq(database_url: str) -> None:
    await _insert_processing_conversation(database_url)
    consumer = DecisionAgentAgUiMessageConsumer(handler=_handler(database_url))
    incoming = FakeIncomingMessage(
        _event_body(
            event_id="markdown-consumer-invalid",
            sequence=1,
            ag_ui={"version": "1.x", "events": []},
        )
    )

    await consumer.process_message(incoming)

    assert incoming.ack_calls == 0
    assert incoming.reject_calls == [{"requeue": False}]


@pytest.mark.asyncio
async def test_consumer_rejects_to_dlq_when_markdown_processing_raises() -> None:
    consumer = DecisionAgentAgUiMessageConsumer(handler=FailingHandler())
    incoming = FakeIncomingMessage(_event_body(event_id="markdown-consumer-failed", sequence=1))

    await consumer.process_message(incoming)

    assert incoming.ack_calls == 0
    assert incoming.reject_calls == [{"requeue": False}]


def _handler(database_url: str) -> DecisionAgentAgUiEventHandler:
    engine = create_async_engine(database_url)
    session_factory = create_async_session_factory(engine)
    return DecisionAgentAgUiEventHandler(
        unit_of_work_factory=lambda: SqlAlchemyUnitOfWork(session_factory),
        id_generator=FakeIdGenerator(start=1_000),
    )


def _command(
    *,
    event_id: str,
    sequence: int,
    content: str = "Generated two recovery candidates.",
    markdown: str = "# Recovery candidates\n\n- Candidate A\n- Candidate B",
    message_id: int = 901,
    turn_id: int = 200,
) -> DecisionAgentAgUiEventCommand:
    return DecisionAgentAgUiEventCommand(
        event_id=event_id,
        event_type="decision_agent.session.ag_ui_event",
        source_service="decision_agent_session",
        conversation_id=100,
        turn_id=turn_id,
        message_id=message_id,
        content=content,
        sequence=sequence,
        ag_ui=markdown,
        trace_id="trace-agui",
        correlation_id="corr-agui",
        occurred_time=1_800_000_000_100,
    )


def _event_body(
    *,
    event_id: str,
    sequence: int,
    ag_ui: object = "# Recovery candidates\n\n- Candidate A\n- Candidate B",
) -> bytes:
    return json.dumps(
        {
            "event_id": event_id,
            "event_type": "decision_agent.session.ag_ui_event",
            "occurred_at": "2026-04-23T10:00:00Z",
            "source_service": "decision_agent_session",
            "trace_id": "trace-agui",
            "correlation_id": "corr-agui",
            "payload": {
                "conversation_id": "100",
                "turn_id": "200",
                "message_id": "901",
                "content": "Generated two recovery candidates.",
                "sequence": sequence,
                "ag_ui": ag_ui,
            },
        }
    ).encode("utf-8")


async def _insert_processing_conversation(
    database_url: str,
    *,
    conversation_status: str = "active",
    message_status: str = "processing",
) -> None:
    engine = create_async_engine(database_url)
    async with engine.begin() as connection:
        await connection.execute(
            ConversationModel.__table__.insert().values(
                f_conversation_id=100,
                f_owner_user_id="user-001",
                f_tenant_id=None,
                f_title="restore backup",
                f_display_summary=None,
                f_status=conversation_status,
                f_interaction_status="executing",
                f_active_turn_id=200,
                f_scenario_binding=None,
                f_tags=[],
                f_retention_policy="conversation_default_v1",
                f_legal_hold=False,
                f_last_active_time=1_800_000_000_000,
                f_archived_time=None,
                f_archived_by=None,
                f_archive_reason=None,
                f_expires_time=None,
                f_expired_time=None,
                f_purge_after_time=None,
                f_purged_time=None,
                f_version=1,
                f_created_time=1_800_000_000_000,
                f_updated_time=1_800_000_000_000,
            )
        )
        await connection.execute(
            ConversationMessageModel.__table__.insert().values(
                f_message_id=200,
                f_conversation_id=100,
                f_parent_message_id=None,
                f_turn_id=200,
                f_role="user",
                f_content_type="text",
                f_content="restore backup",
                f_rich_payload=None,
                f_status=message_status,
                f_client_message_id="client-001",
                f_idempotency_key="idem-001",
                f_trace_id="trace-001",
                f_correlation_id="corr-001",
                f_error_code=None,
                f_created_time=1_800_000_000_000,
                f_updated_time=1_800_000_000_000,
            )
        )
        await connection.execute(
            ConversationStatusEventModel.__table__.insert().values(
                f_status_event_id=300,
                f_conversation_id=100,
                f_sequence=1,
                f_event_type="message.updated",
                f_event_version="v1",
                f_message_id=200,
                f_turn_id=200,
                f_payload={"message_status": "processing"},
                f_visible_to_user=True,
                f_trace_id="trace-001",
                f_correlation_id="corr-001",
                f_created_time=1_800_000_000_000,
                f_updated_time=1_800_000_000_000,
            )
        )
    await engine.dispose()


async def _assistant_messages(database_url: str) -> list[ConversationMessageModel]:
    engine = create_async_engine(database_url)
    session_factory = create_async_session_factory(engine)
    async with session_factory() as session:
        rows = await session.scalars(
            select(ConversationMessageModel)
            .where(ConversationMessageModel.f_role == "assistant")
            .order_by(ConversationMessageModel.f_message_id.asc())
        )
    await engine.dispose()
    return list(rows.all())


async def _message(database_url: str, message_id: int) -> ConversationMessageModel:
    engine = create_async_engine(database_url)
    session_factory = create_async_session_factory(engine)
    async with session_factory() as session:
        row = await session.scalar(
            select(ConversationMessageModel).where(
                ConversationMessageModel.f_message_id == message_id
            )
        )
    await engine.dispose()
    assert row is not None
    return row


async def _get_conversation(database_url: str) -> ConversationModel:
    engine = create_async_engine(database_url)
    session_factory = create_async_session_factory(engine)
    async with session_factory() as session:
        row = await session.scalar(
            select(ConversationModel).where(ConversationModel.f_conversation_id == 100)
        )
    await engine.dispose()
    assert row is not None
    return row


async def _event_types(database_url: str) -> list[str]:
    engine = create_async_engine(database_url)
    session_factory = create_async_session_factory(engine)
    async with session_factory() as session:
        rows = await session.scalars(
            select(ConversationStatusEventModel.f_event_type).order_by(
                ConversationStatusEventModel.f_sequence.asc()
            )
        )
    await engine.dispose()
    return list(rows.all())


async def _latest_context_delta(database_url: str) -> ConversationContextDeltaModel:
    engine = create_async_engine(database_url)
    session_factory = create_async_session_factory(engine)
    async with session_factory() as session:
        row = await session.scalar(
            select(ConversationContextDeltaModel).order_by(
                ConversationContextDeltaModel.f_created_time.desc()
            )
        )
    await engine.dispose()
    assert row is not None
    return row


async def _count(database_url: str, column: object) -> int:
    engine = create_async_engine(database_url)
    async with engine.connect() as connection:
        count = await connection.scalar(select(func.count(column)))
    await engine.dispose()
    return int(count or 0)


class FakeIdGenerator:
    def __init__(self, *, start: int) -> None:
        self._next = start

    def next_id(self) -> int:
        current = self._next
        self._next += 1
        return current


class FakeIncomingMessage:
    def __init__(self, body: bytes) -> None:
        self.body = body
        self.ack_calls = 0
        self.reject_calls: list[dict[str, bool]] = []

    async def ack(self) -> None:
        self.ack_calls += 1

    async def reject(self, *, requeue: bool) -> None:
        self.reject_calls.append({"requeue": requeue})


class FailingHandler:
    async def handle(self, command: DecisionAgentAgUiEventCommand) -> object:
        del command
        raise RuntimeError("boom")
