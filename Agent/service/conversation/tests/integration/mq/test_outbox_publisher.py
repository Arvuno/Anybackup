import asyncio
from collections.abc import Iterator
from pathlib import Path
from typing import Any

import pytest
from sqlalchemy import select
from sqlalchemy.ext.asyncio import create_async_engine

from app.application.use_cases.outbox import OutboxPublisherWorker
from app.infrastructure.locking.redis.lock import RedisGlobalLock
from app.infrastructure.persistence.sqlalchemy.models import Base, ConversationMqOutboxModel
from app.infrastructure.persistence.sqlalchemy.session import create_async_session_factory
from app.infrastructure.persistence.sqlalchemy.unit_of_work import SqlAlchemyUnitOfWork


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
async def test_outbox_publisher_publishes_pending_event_and_marks_published(
    database_url: str,
) -> None:
    await _insert_outbox(database_url, outbox_id=1, event_id="evt-001")
    publisher = FakePublisher()
    worker = _worker(database_url, publisher=publisher, lock=FakeLock(acquired=True))

    result = await worker.run_once(worker_id="worker-001", now_ms=1_800_000_000_000)

    row = await _get_outbox(database_url, 1)
    assert result.lock_acquired is True
    assert result.scanned == 1
    assert result.published == 1
    assert row.f_status == "published"
    assert row.f_attempt_count == 0
    assert len(publisher.messages) == 1
    message = publisher.messages[0]
    assert message.exchange == "conversation.message.events"
    assert message.routing_key == "conversation.message.sent.v1"
    assert message.body["event_id"] == "evt-001"
    assert message.body["event_type"] == "conversation.message.sent"
    assert set(message.body) == {
        "event_id",
        "event_type",
        "occurred_at",
        "source_service",
        "payload",
    }
    assert message.body["payload"] == {
        "conversation_id": "100",
        "message_id": "200",
        "content": "restore database",
    }
    assert message.headers["trace_id"] == "trace-001"
    assert message.headers["correlation_id"] == "corr-001"
    assert "traceparent" in message.headers


@pytest.mark.asyncio
async def test_outbox_publisher_skips_when_redis_lock_is_held(database_url: str) -> None:
    await _insert_outbox(database_url, outbox_id=1, event_id="evt-001")
    publisher = FakePublisher()
    worker = _worker(database_url, publisher=publisher, lock=FakeLock(acquired=False))

    result = await worker.run_once(worker_id="worker-001", now_ms=1_800_000_000_000)

    row = await _get_outbox(database_url, 1)
    assert result.lock_acquired is False
    assert result.scanned == 0
    assert publisher.messages == []
    assert row.f_status == "pending"


@pytest.mark.asyncio
async def test_outbox_publisher_retries_when_broker_publish_fails(database_url: str) -> None:
    await _insert_outbox(database_url, outbox_id=1, event_id="evt-001")
    worker = _worker(
        database_url,
        publisher=FakePublisher(error=RuntimeError("broker unavailable")),
        lock=FakeLock(acquired=True),
        max_attempts=5,
        retry_delay_ms=2_000,
    )

    result = await worker.run_once(worker_id="worker-001", now_ms=1_800_000_000_000)

    row = await _get_outbox(database_url, 1)
    assert result.failed == 1
    assert row.f_status == "retry"
    assert row.f_attempt_count == 1
    assert row.f_next_retry_time == 1_800_000_002_000
    assert row.f_last_error_code == "PUBLISH_FAILED"


@pytest.mark.asyncio
async def test_outbox_publisher_marks_dlq_after_max_attempts(database_url: str) -> None:
    await _insert_outbox(
        database_url,
        outbox_id=1,
        event_id="evt-001",
        attempt_count=4,
    )
    worker = _worker(
        database_url,
        publisher=FakePublisher(error=RuntimeError("broker unavailable")),
        lock=FakeLock(acquired=True),
        max_attempts=5,
    )

    result = await worker.run_once(worker_id="worker-001", now_ms=1_800_000_000_000)

    row = await _get_outbox(database_url, 1)
    assert result.dlq == 1
    assert row.f_status == "dlq"
    assert row.f_attempt_count == 5
    assert row.f_last_error_code == "PUBLISH_FAILED"


@pytest.mark.asyncio
async def test_outbox_publisher_rejects_forbidden_payload_fields(database_url: str) -> None:
    await _insert_outbox(
        database_url,
        outbox_id=1,
        event_id="evt-001",
        payload={"platform_session_id": "leaked", "message": {"content": "restore"}},
    )
    publisher = FakePublisher()
    worker = _worker(database_url, publisher=publisher, lock=FakeLock(acquired=True))

    result = await worker.run_once(worker_id="worker-001", now_ms=1_800_000_000_000)

    row = await _get_outbox(database_url, 1)
    assert result.dlq == 1
    assert publisher.messages == []
    assert row.f_status == "dlq"
    assert row.f_last_error_code == "FORBIDDEN_PAYLOAD"


@pytest.mark.asyncio
async def test_redis_global_lock_uses_set_nx_px_and_lua_release() -> None:
    redis = FakeRedis()
    lock = RedisGlobalLock(redis, owner_token_factory=lambda: "owner-token")

    lease = await lock.acquire(
        "conversation_service:lock:outbox:publisher",
        ttl_ms=30_000,
    )
    assert lease is not None
    await lease.release()

    assert redis.set_calls == [
        {
            "name": "conversation_service:lock:outbox:publisher",
            "value": "owner-token",
            "nx": True,
            "px": 30_000,
        }
    ]
    assert redis.eval_calls == [
        {
            "numkeys": 1,
            "keys": ["conversation_service:lock:outbox:publisher"],
            "args": ["owner-token"],
        }
    ]


def _worker(
    database_url: str,
    *,
    publisher: "FakePublisher",
    lock: "FakeLock",
    max_attempts: int = 5,
    retry_delay_ms: int = 1_000,
) -> OutboxPublisherWorker:
    engine = create_async_engine(database_url)
    session_factory = create_async_session_factory(engine)
    return OutboxPublisherWorker(
        unit_of_work_factory=lambda: SqlAlchemyUnitOfWork(session_factory),
        publisher=publisher,
        lock=lock,
        exchange_name="conversation.message.events",
        batch_size=100,
        max_attempts=max_attempts,
        retry_delay_ms=retry_delay_ms,
        lock_ttl_ms=30_000,
    )


async def _insert_outbox(
    database_url: str,
    *,
    outbox_id: int,
    event_id: str,
    payload: dict[str, Any] | None = None,
    attempt_count: int = 0,
) -> None:
    engine = create_async_engine(database_url)
    async with engine.begin() as connection:
        await connection.execute(
            ConversationMqOutboxModel.__table__.insert().values(
                f_outbox_id=outbox_id,
                f_event_id=event_id,
                f_event_type="conversation.message.sent.v1",
                f_routing_key="conversation.message.sent.v1",
                f_conversation_id=100,
                f_message_id=200,
                f_payload=payload
                or {
                    "conversation_id": "100",
                    "message_id": "200",
                    "content": "restore database",
                },
                f_status="pending",
                f_attempt_count=attempt_count,
                f_next_retry_time=None,
                f_last_error_code=None,
                f_trace_id="trace-001",
                f_correlation_id="corr-001",
                f_idempotency_key="idem-001",
                f_created_time=1_800_000_000_000,
                f_updated_time=1_800_000_000_000,
            )
        )
    await engine.dispose()


async def _get_outbox(database_url: str, outbox_id: int) -> ConversationMqOutboxModel:
    engine = create_async_engine(database_url)
    session_factory = create_async_session_factory(engine)
    async with session_factory() as session:
        row = await session.scalar(
            select(ConversationMqOutboxModel).where(
                ConversationMqOutboxModel.f_outbox_id == outbox_id
            )
        )
    await engine.dispose()
    assert row is not None
    return row


class FakePublisher:
    def __init__(self, error: Exception | None = None) -> None:
        self.messages: list[object] = []
        self._error = error

    async def publish(self, message: object) -> None:
        if self._error is not None:
            raise self._error
        self.messages.append(message)


class FakeLease:
    def __init__(self) -> None:
        self.released = False

    async def release(self) -> None:
        self.released = True


class FakeLock:
    def __init__(self, *, acquired: bool) -> None:
        self.acquired = acquired
        self.acquire_calls: list[dict[str, object]] = []

    async def acquire(self, key: str, ttl_ms: int) -> FakeLease | None:
        self.acquire_calls.append({"key": key, "ttl_ms": ttl_ms})
        if not self.acquired:
            return None
        return FakeLease()


class FakeRedis:
    def __init__(self) -> None:
        self.set_calls: list[dict[str, object]] = []
        self.eval_calls: list[dict[str, object]] = []

    async def set(self, name: str, value: str, *, nx: bool, px: int) -> bool:
        self.set_calls.append({"name": name, "value": value, "nx": nx, "px": px})
        return True

    async def eval(self, script: str, numkeys: int, *keys_and_args: str) -> int:
        del script
        self.eval_calls.append(
            {
                "numkeys": numkeys,
                "keys": list(keys_and_args[:numkeys]),
                "args": list(keys_and_args[numkeys:]),
            }
        )
        return 1
