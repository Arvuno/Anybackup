import logging
from datetime import UTC, datetime
from typing import Any, Protocol

import aio_pika
from aio_pika import ExchangeType

from app.application.commands.agent_events import CoreAgentStatusEventCommand
from app.application.use_cases.agent_status import CoreAgentStatusEventHandler
from app.infrastructure.messaging.rabbitmq.publisher import decode_message_body

logger = logging.getLogger(__name__)

CORE_AGENT_STATUS_ROUTING_KEYS = (
    "core_agent.run.accepted.v1",
    "core_agent.run.completed.v1",
    "core_agent.run.failed.v1",
)


class IncomingMessage(Protocol):
    body: bytes

    async def ack(self) -> None:
        raise NotImplementedError

    async def reject(self, *, requeue: bool) -> None:
        raise NotImplementedError


class CoreAgentStatusHandler(Protocol):
    async def handle(self, command: CoreAgentStatusEventCommand) -> object:
        raise NotImplementedError


class CoreAgentStatusMessageConsumer:
    def __init__(self, *, handler: CoreAgentStatusHandler) -> None:
        self._handler = handler

    async def process_message(self, message: IncomingMessage) -> None:
        try:
            command = _command_from_body(decode_message_body(message.body))
            logger.info(
                "core_agent_status_consume_enter",
                extra={
                    "event_id": command.event_id,
                    "event_type": command.event_type,
                    "conversation_id": command.conversation_id,
                    "message_id": command.message_id,
                },
            )
            await self._handler.handle(command)
            await message.ack()
        except Exception:
            logger.exception("core_agent_status_consume_failed")
            await message.reject(requeue=False)


class RabbitMqCoreAgentStatusConsumer:
    def __init__(
        self,
        *,
        rabbitmq_url: str,
        exchange_name: str,
        queue_name: str,
        prefetch_count: int,
        handler: CoreAgentStatusEventHandler,
        routing_keys: tuple[str, ...] = CORE_AGENT_STATUS_ROUTING_KEYS,
    ) -> None:
        self._rabbitmq_url = rabbitmq_url
        self._exchange_name = exchange_name
        self._queue_name = queue_name
        self._prefetch_count = prefetch_count
        self._routing_keys = routing_keys
        self._message_consumer = CoreAgentStatusMessageConsumer(handler=handler)
        self._connection: aio_pika.abc.AbstractRobustConnection | None = None
        self._channel: aio_pika.abc.AbstractChannel | None = None

    async def start(self) -> None:
        logger.info(
            "core_agent_status_consumer_start_enter",
            extra={"queue_name": self._queue_name, "prefetch_count": self._prefetch_count},
        )
        connection = await aio_pika.connect_robust(self._rabbitmq_url)
        channel = await connection.channel()
        await channel.set_qos(prefetch_count=self._prefetch_count)
        exchange = await channel.declare_exchange(
            self._exchange_name,
            ExchangeType.TOPIC,
            durable=True,
        )
        queue = await channel.declare_queue(self._queue_name, durable=True)
        for routing_key in self._routing_keys:
            await queue.bind(exchange, routing_key=routing_key)
        await queue.consume(self._message_consumer.process_message, no_ack=False)
        self._connection = connection
        self._channel = channel

    async def close(self) -> None:
        if self._connection is not None:
            await self._connection.close()
        self._connection = None
        self._channel = None


def _command_from_body(body: dict[str, Any]) -> CoreAgentStatusEventCommand:
    payload = body.get("payload", {})
    if not isinstance(payload, dict):
        raise ValueError("core agent status payload must be an object")
    return CoreAgentStatusEventCommand(
        event_id=str(body["event_id"]),
        event_type=str(body["event_type"]),
        event_version=str(body.get("event_version", "v1")),
        conversation_id=int(payload["conversation_id"]),
        payload=payload,
        message_id=_optional_int(payload.get("message_id")),
        trace_id=str(body.get("trace_id") or ""),
        correlation_id=str(body.get("correlation_id") or ""),
        occurred_time=_occurred_at_to_ms(body.get("occurred_at")),
    )


def _optional_int(value: object) -> int | None:
    if value is None:
        return None
    if isinstance(value, bool):
        raise ValueError("integer value cannot be boolean")
    if isinstance(value, int):
        return value
    if isinstance(value, str):
        return int(value)
    raise ValueError("integer value must be int or string")


def _occurred_at_to_ms(value: object) -> int | None:
    if value is None:
        return None
    if not isinstance(value, str):
        raise ValueError("occurred_at must be a string")
    parsed = datetime.fromisoformat(value.replace("Z", "+00:00"))
    if parsed.tzinfo is None:
        parsed = parsed.replace(tzinfo=UTC)
    return int(parsed.timestamp() * 1000)
