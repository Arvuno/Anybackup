from __future__ import annotations

import json
from uuid import uuid4
import logging
import re

import pika

from core_agent_service.domain.statuses import build_status_event


logger = logging.getLogger(__name__)


_SECRET_PATTERNS = (
    re.compile(r"(?i)(ak|sk|token|password|passwd|secret)(\s*[=:]\s*)([^,\s，}]+)"),
    re.compile(r"(?i)(amqps?://[^:\s/@]+:)([^@\s]+)(@)"),
)


class StatusPublisher:
    # 三类状态事件是当前最小中转版唯一允许向外发布的结果。
    accepted_key = "core_agent.run.accepted.v1"
    completed_key = "core_agent.run.completed.v1"
    failed_key = "core_agent.run.failed.v1"

    def __init__(
        self,
        *,
        rabbitmq_url: str | None = None,
        exchange: str,
        exchange_type: str = "topic",
        status_routing_key: str | None = None,
        outbound_event_repository=None,
        connection_factory=None,
        channel=None,
    ):
        self.rabbitmq_url = rabbitmq_url
        self.exchange = exchange
        self.exchange_type = exchange_type
        self.status_routing_key = status_routing_key
        self.outbound_event_repository = outbound_event_repository
        self.connection_factory = connection_factory
        self.connection = None
        self.channel = channel
        self.published: list[tuple[str, str]] = []

    def publish_accepted(
        self,
        conversation_id: str,
        content: str | None = None,
        message_id: str | None = None,
    ) -> dict:
        return self._publish(
            routing_key=self.accepted_key,
            event_type="core_agent.run.accepted",
            conversation_id=conversation_id,
            content=content or "accepted",
            message_id=message_id,
        )

    def publish_completed(
        self,
        conversation_id: str,
        content: str | None = None,
        message_id: str | None = None,
    ) -> dict:
        return self._publish(
            routing_key=self.completed_key,
            event_type="core_agent.run.completed",
            conversation_id=conversation_id,
            content=content or "completed",
            message_id=message_id,
        )

    def publish_failed(
        self, conversation_id: str, error_message: str, message_id: str | None = None
    ) -> dict:
        return self._publish(
            routing_key=self.failed_key,
            event_type="core_agent.run.failed",
            conversation_id=conversation_id,
            content=error_message,
            message_id=message_id,
        )

    def _publish(
        self,
        *,
        routing_key: str,
        event_type: str,
        conversation_id: str,
        content: str,
        message_id: str | None = None,
    ) -> dict:
        event = build_status_event(
            event_id=f"evt_status_{uuid4().hex[:12]}",
            event_type=event_type,
            conversation_id=conversation_id,
            content=content,
            message_id=message_id,
        )
        delivery_routing_key = self.status_routing_key or routing_key
        try:
            self._publish_to_rabbitmq(
                event, delivery_routing_key, event_type, conversation_id
            )
        except (pika.exceptions.AMQPError, OSError) as exc:
            # RabbitMQ 可能在连接空闲或网络抖动后重置 publisher 连接；状态发布可用同一事件安全重试一次。
            logger.warning(
                "rabbitmq status publisher connection lost, retrying once",
                extra={
                    "routing_key": delivery_routing_key,
                    "event_type": event_type,
                    "conversation_id": conversation_id,
                },
                exc_info=exc,
            )
            self._drop_channel()
            self._publish_to_rabbitmq(
                event, delivery_routing_key, event_type, conversation_id
            )

        if self.outbound_event_repository is not None:
            self.outbound_event_repository.save(
                event_id=event.event_id,
                event_type=event.event_type,
                occurred_at=event.occurred_at,
                conversation_id=conversation_id,
                content=content,
            )
        logger.info(
            "rabbitmq status event published",
            extra={
                "routing_key": delivery_routing_key,
                "event_id": event.event_id,
                "conversation_id": conversation_id,
            },
        )
        return event.model_dump(mode="json", exclude_none=True)

    def _publish_to_rabbitmq(
        self, event, routing_key: str, event_type: str, conversation_id: str
    ) -> None:
        channel = self._ensure_channel()
        # 先组装并记录统一事件结构，再发 MQ，这样测试、持久化和运行态排障都能复用同一份数据。
        publish_payload = {
            "exchange": self.exchange,
            "routing_key": routing_key,
            "event_type": event_type,
            "conversation_id": conversation_id,
            "body": event.model_dump(mode="json", exclude_none=True),
        }
        logger.info(
            "rabbitmq status event publish payload: %s",
            json.dumps(self._sanitize_for_log(publish_payload), ensure_ascii=False),
        )
        logger.info(
            "publishing rabbitmq status event",
            extra={
                "routing_key": routing_key,
                "event_type": event_type,
                "conversation_id": conversation_id,
            },
        )
        channel.basic_publish(
            exchange=self.exchange,
            routing_key=routing_key,
            body=event.model_dump_json(exclude_none=True),
            properties=pika.BasicProperties(
                content_type="application/json",
                delivery_mode=2,
            ),
        )
        self.published.append((routing_key, event.model_dump_json(exclude_none=True)))

    def _ensure_channel(self):
        if self.channel is not None and not getattr(self.channel, "is_closed", False):
            try:
                if not self.exchange:
                    return self.channel
                self.channel.exchange_declare(
                    exchange=self.exchange,
                    exchange_type=self.exchange_type,
                    durable=True,
                )
                return self.channel
            except (pika.exceptions.AMQPError, OSError):
                # Channel 可能已失效，丢弃后重建
                logger.warning(
                    "rabbitmq publisher channel verification failed, will rebuild connection"
                )
                self._drop_channel()

        if self.connection_factory is not None:
            self.connection = self.connection_factory()
        else:
            if not self.rabbitmq_url:
                raise ValueError(
                    "rabbitmq_url is required for real RabbitMQ publishing"
                )
            parameters = pika.URLParameters(self.rabbitmq_url)
            self.connection = pika.BlockingConnection(parameters)
        # publisher 自己维护独立连接，避免和 consumer 共享 channel 造成生命周期耦合。
        logger.info(
            "rabbitmq publisher connection created",
            extra={"rabbitmq_url": self.rabbitmq_url, "exchange": self.exchange},
        )
        self.channel = self.connection.channel()
        if self.exchange:
            self.channel.exchange_declare(
                exchange=self.exchange,
                exchange_type=self.exchange_type,
                durable=True,
            )
            logger.info(
                "rabbitmq publisher exchange declared",
                extra={"exchange": self.exchange, "exchange_type": self.exchange_type},
            )
        else:
            # 默认 Exchange 由 RabbitMQ 内置提供，不能显式声明；routing key 直接使用目标队列名。
            logger.info(
                "rabbitmq publisher using default exchange",
                extra={"routing_key": self.status_routing_key},
            )
        return self.channel

    def _drop_channel(self) -> None:
        # 出错后丢弃旧连接引用，下一次发布重新创建 publisher 独立连接。
        self.channel = None
        self.connection = None

    @classmethod
    def _sanitize_for_log(cls, value):
        if isinstance(value, dict):
            return {key: cls._sanitize_for_log(item) for key, item in value.items()}
        if isinstance(value, list):
            return [cls._sanitize_for_log(item) for item in value]
        if isinstance(value, str):
            sanitized = value
            for pattern in _SECRET_PATTERNS:
                sanitized = pattern.sub(cls._mask_secret_match, sanitized)
            return sanitized
        return value

    @staticmethod
    def _mask_secret_match(match) -> str:
        # 第一类是 ak=xxx 形式，第二类是 amqp://user:pass@host 形式。
        middle = match.group(2)
        if middle.strip().startswith(("=", ":")):
            return f"{match.group(1)}{middle}***"
        return f"{match.group(1)}***{match.group(3)}"
