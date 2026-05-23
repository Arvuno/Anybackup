import logging
from collections.abc import Callable
from dataclasses import replace
from datetime import UTC, datetime
from time import time
from typing import Any

from app.application.models.conversation import MqOutboxRecord
from app.application.models.outbox import OutboxPublishMessage, OutboxPublishRunResult
from app.application.ports.locking import GlobalLock
from app.application.ports.messaging import EventPublisher
from app.application.ports.unit_of_work import UnitOfWork
from app.domain.message import MessageStatus

logger = logging.getLogger(__name__)

OUTBOX_LOCK_KEY = "conversation_service:lock:outbox:publisher"
PUBLISH_FAILED = "PUBLISH_FAILED"
FORBIDDEN_PAYLOAD = "FORBIDDEN_PAYLOAD"
FORBIDDEN_FIELD_NAMES = frozenset(
    {
        "platform_session_id",
        "kweaver_session_id",
        "password",
        "token",
        "secret",
        "private_key",
        "connection_string",
        "system_prompt",
        "internal_reasoning",
    }
)


class ForbiddenOutboxPayloadError(ValueError):
    pass


class OutboxPublisherWorker:
    def __init__(
        self,
        *,
        unit_of_work_factory: Callable[[], UnitOfWork],
        publisher: EventPublisher,
        lock: GlobalLock,
        exchange_name: str,
        batch_size: int,
        max_attempts: int,
        retry_delay_ms: int,
        lock_ttl_ms: int,
    ) -> None:
        self._unit_of_work_factory = unit_of_work_factory
        self._publisher = publisher
        self._lock = lock
        self._exchange_name = exchange_name
        self._batch_size = batch_size
        self._max_attempts = max_attempts
        self._retry_delay_ms = retry_delay_ms
        self._lock_ttl_ms = lock_ttl_ms

    async def run_once(
        self,
        *,
        worker_id: str,
        now_ms: int | None = None,
    ) -> OutboxPublishRunResult:
        now = now_ms or _current_time_ms()
        logger.info(
            "outbox_publisher_enter",
            extra={"worker_id": worker_id, "lock_key": OUTBOX_LOCK_KEY},
        )
        lease = await self._lock.acquire(OUTBOX_LOCK_KEY, self._lock_ttl_ms)
        if lease is None:
            logger.info(
                "outbox_publisher_lock_skipped",
                extra={"worker_id": worker_id, "lock_key": OUTBOX_LOCK_KEY},
            )
            return OutboxPublishRunResult(lock_acquired=False)

        scanned = 0
        published = 0
        failed = 0
        dlq = 0
        try:
            async with self._unit_of_work_factory() as unit_of_work:
                events = await unit_of_work.outbox.list_publishable(
                    limit=self._batch_size,
                    now_ms=now,
                )
            scanned = len(events)
            logger.info(
                "outbox_batch_scanned",
                extra={"worker_id": worker_id, "scanned": scanned},
            )

            for event in events:
                try:
                    message = build_outbox_publish_message(event, self._exchange_name)
                    assert_no_forbidden_payload(message.body)
                    await self._publisher.publish(message)
                except ForbiddenOutboxPayloadError:
                    dlq += 1
                    await self._mark_dlq(event, now, FORBIDDEN_PAYLOAD)
                    logger.warning(
                        "outbox_payload_rejected",
                        extra={"event_id": event.event_id, "error_code": FORBIDDEN_PAYLOAD},
                    )
                except Exception:
                    next_attempt_count = event.attempt_count + 1
                    if next_attempt_count >= self._max_attempts:
                        dlq += 1
                        await self._mark_dlq(event, now, PUBLISH_FAILED)
                    else:
                        failed += 1
                        await self._mark_retry(event, next_attempt_count, now)
                    logger.warning(
                        "outbox_publish_failed",
                        extra={"event_id": event.event_id, "attempt": next_attempt_count},
                    )
                else:
                    published += 1
                    await self._mark_published(event, now)
                    logger.info(
                        "outbox_published",
                        extra={"event_id": event.event_id, "routing_key": event.routing_key},
                    )
        finally:
            await lease.release()

        return OutboxPublishRunResult(
            lock_acquired=True,
            scanned=scanned,
            published=published,
            failed=failed,
            dlq=dlq,
        )

    async def _mark_published(self, event: MqOutboxRecord, now_ms: int) -> None:
        async with self._unit_of_work_factory() as unit_of_work:
            if event.message_id is not None:
                message = await unit_of_work.messages.get_by_id(event.message_id)
                if message is not None and message.status is MessageStatus.PERSISTED:
                    await unit_of_work.messages.update(
                        replace(
                            message,
                            status=MessageStatus.PUBLISHED,
                            updated_time=now_ms,
                        )
                    )
            await unit_of_work.outbox.mark_published(outbox_id=event.outbox_id, now_ms=now_ms)

    async def _mark_retry(
        self,
        event: MqOutboxRecord,
        attempt_count: int,
        now_ms: int,
    ) -> None:
        async with self._unit_of_work_factory() as unit_of_work:
            await unit_of_work.outbox.mark_retry(
                outbox_id=event.outbox_id,
                attempt_count=attempt_count,
                next_retry_time=now_ms + self._retry_delay_ms,
                error_code=PUBLISH_FAILED,
                now_ms=now_ms,
            )

    async def _mark_dlq(self, event: MqOutboxRecord, now_ms: int, error_code: str) -> None:
        async with self._unit_of_work_factory() as unit_of_work:
            await unit_of_work.outbox.mark_dlq(
                outbox_id=event.outbox_id,
                attempt_count=event.attempt_count + 1,
                error_code=error_code,
                now_ms=now_ms,
            )


def build_outbox_publish_message(
    event: MqOutboxRecord,
    exchange_name: str,
) -> OutboxPublishMessage:
    event_type, _event_version = _split_event_type(event.event_type)
    body: dict[str, Any] = {
        "event_id": event.event_id,
        "event_type": event_type,
        "occurred_at": _ms_to_iso(event.created_time),
        "source_service": "conversation_service",
        "payload": event.payload,
    }
    headers = {
        "trace_id": event.trace_id,
        "correlation_id": event.correlation_id,
        "traceparent": _traceparent(event.trace_id),
    }
    return OutboxPublishMessage(
        exchange=exchange_name,
        routing_key=event.routing_key,
        body=body,
        headers=headers,
        message_id=event.event_id,
        correlation_id=event.correlation_id,
    )


def assert_no_forbidden_payload(value: Any) -> None:
    _scan_forbidden_payload(value, path="$")


def _scan_forbidden_payload(value: Any, *, path: str) -> None:
    if isinstance(value, dict):
        for key, item in value.items():
            key_text = str(key).lower()
            if any(forbidden in key_text for forbidden in FORBIDDEN_FIELD_NAMES):
                raise ForbiddenOutboxPayloadError(f"forbidden payload field: {path}.{key}")
            _scan_forbidden_payload(item, path=f"{path}.{key}")
    elif isinstance(value, list):
        for index, item in enumerate(value):
            _scan_forbidden_payload(item, path=f"{path}[{index}]")


def _split_event_type(event_type: str) -> tuple[str, str]:
    if event_type.endswith(".v1"):
        return event_type[:-3], "v1"
    return event_type, "v1"


def _traceparent(trace_id: str) -> str:
    sanitized = "".join(ch for ch in trace_id.lower() if ch in "0123456789abcdef")
    if len(sanitized) >= 32:
        trace_hex = sanitized[:32]
    else:
        trace_hex = sanitized.ljust(32, "0")
    return f"00-{trace_hex}-0000000000000000-01"


def _ms_to_iso(value: int) -> str:
    return datetime.fromtimestamp(value / 1000, UTC).isoformat().replace("+00:00", "Z")


def _current_time_ms() -> int:
    return int(time() * 1000)
