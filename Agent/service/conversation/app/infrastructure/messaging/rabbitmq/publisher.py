import json
import logging
from typing import Any

import aio_pika
from aio_pika import DeliveryMode, ExchangeType, Message

from app.application.models.outbox import OutboxPublishMessage

logger = logging.getLogger(__name__)


class RabbitMqEventPublisher:
    def __init__(self, *, rabbitmq_url: str, exchange_name: str) -> None:
        self._rabbitmq_url = rabbitmq_url
        self._exchange_name = exchange_name
        self._connection: aio_pika.abc.AbstractRobustConnection | None = None
        self._channel: aio_pika.abc.AbstractChannel | None = None
        self._exchange: aio_pika.abc.AbstractExchange | None = None

    async def publish(self, message: OutboxPublishMessage) -> None:
        exchange = await self._get_exchange()
        logger.info(
            "rabbitmq_publish_enter",
            extra={"event_id": message.message_id, "routing_key": message.routing_key},
        )
        await exchange.publish(
            Message(
                body=json.dumps(message.body, ensure_ascii=False).encode("utf-8"),
                content_type="application/json",
                delivery_mode=DeliveryMode.PERSISTENT,
                message_id=message.message_id,
                correlation_id=message.correlation_id,
                headers=dict(message.headers),
            ),
            routing_key=message.routing_key,
            mandatory=True,
        )

    async def close(self) -> None:
        if self._connection is not None:
            await self._connection.close()
        self._connection = None
        self._channel = None
        self._exchange = None

    async def _get_exchange(self) -> aio_pika.abc.AbstractExchange:
        if self._exchange is not None:
            return self._exchange
        connection = await aio_pika.connect_robust(self._rabbitmq_url)
        channel = await connection.channel(publisher_confirms=True)
        exchange = await channel.declare_exchange(
            self._exchange_name,
            ExchangeType.TOPIC,
            durable=True,
        )
        self._connection = connection
        self._channel = channel
        self._exchange = exchange
        return exchange


def decode_message_body(body: bytes) -> dict[str, Any]:
    raw = json.loads(body.decode("utf-8"))
    if not isinstance(raw, dict):
        raise ValueError("RabbitMQ message body must be a JSON object")
    return raw
