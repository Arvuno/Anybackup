from __future__ import annotations

import json
import logging
from threading import Thread

import pika


logger = logging.getLogger(__name__)


class RabbitMQConsumerPool:
    def __init__(self, consumers):
        self.consumers = tuple(consumers)

    def run(self) -> None:
        if len(self.consumers) == 1:
            self.consumers[0].run()
            return
        logger.info("rabbitmq consumer pool starting", extra={"consumer_count": len(self.consumers)})
        threads = [
            Thread(target=consumer.run, name=f"core-agent-mq-consumer-{index}", daemon=False)
            for index, consumer in enumerate(self.consumers, start=1)
        ]
        for thread in threads:
            thread.start()
        for thread in threads:
            thread.join()


class RabbitMQConsumer:
    def __init__(
        self,
        *,
        rabbitmq_url: str,
        exchange: str,
        queue: str,
        relay_service,
        publisher,
        session=None,
        exchange_type: str = "topic",
        routing_keys: tuple[str, ...] | None = None,
        connection_factory=None,
        channel=None,
    ):
        self.rabbitmq_url = rabbitmq_url
        self.exchange = exchange
        self.exchange_type = exchange_type
        self.queue = queue
        self.relay_service = relay_service
        self.publisher = publisher
        self.session = session
        self.routing_keys = routing_keys or (
            "conversation.message.sent.v1",
            "conversation.message.cancel_requested.v1",
        )
        self.connection_factory = connection_factory
        self.connection = None
        self.channel = channel

    def handle_delivery(self, body: bytes) -> None:
        # 先发布 accepted，再进入业务处理，这样上游能尽早知道消息已被服务接收。
        raw_event = json.loads(body)
        payload = raw_event["payload"]
        conversation_id = payload["conversation_id"]
        message_id = payload.get("message_id")
        logger.info(
            "rabbitmq delivery received",
            extra={"event_type": raw_event["event_type"], "conversation_id": conversation_id},
        )
        self.publisher.publish_accepted(conversation_id, message_id=message_id)
        try:
            if raw_event["event_type"] == "conversation.message.cancel_requested":
                result = self.relay_service.handle_cancel(raw_event)
            else:
                result = self.relay_service.handle_message(raw_event)
        except Exception as exc:
            if self.session is not None:
                # MQ 消费失败时要先回滚本地事务，避免把半成品映射或事件记录提交出去。
                self.session.rollback()
                logger.exception("session rollback executed after message failure", extra={"conversation_id": conversation_id})
            self.publisher.publish_failed(conversation_id, str(exc), message_id=message_id)
            if self.session is not None:
                # failed 状态本身也是业务事实，需要单独提交，不能跟前面的失败事务绑在一起。
                self.session.commit()
                logger.info("failed status committed", extra={"conversation_id": conversation_id})
            raise
        self.publisher.publish_completed(conversation_id, result.status, message_id=message_id)
        if self.session is not None:
            # 只有业务处理和状态发布都完成后，才统一提交这次消费对应的数据库事实。
            self.session.commit()
            logger.info("successful processing committed", extra={"conversation_id": conversation_id})

    def run(self) -> None:
        while True:
            try:
                self._start_consuming_once()
            except (pika.exceptions.AMQPError, OSError) as exc:
                # RabbitMQ 长连接可能在长耗时业务处理后被 broker 重置；外层循环负责重建连接继续消费。
                logger.warning("rabbitmq consumer connection lost, reconnecting", exc_info=exc)
                self._drop_channel()
                continue
            return

    def _start_consuming_once(self) -> None:
        channel = self._ensure_channel()
        channel.exchange_declare(
            exchange=self.exchange,
            exchange_type=self.exchange_type,
            durable=True,
        )
        channel.queue_declare(queue=self.queue, durable=True)
        for routing_key in self.routing_keys:
            channel.queue_bind(
                exchange=self.exchange,
                queue=self.queue,
                routing_key=routing_key,
            )
            logger.info(
                "rabbitmq queue bound",
                extra={"exchange": self.exchange, "queue": self.queue, "routing_key": routing_key},
            )
        channel.basic_qos(prefetch_count=1)
        channel.basic_consume(queue=self.queue, on_message_callback=self.on_message)
        logger.info(
            "rabbitmq consumer started",
            extra={"exchange": self.exchange, "queue": self.queue, "routing_keys": ",".join(self.routing_keys)},
        )
        channel.start_consuming()

    def on_message(self, channel, method, properties, body: bytes) -> None:
        try:
            self.handle_delivery(body)
        except Exception:
            # 当前阶段失败消息直接进入 nack 且不重回队列，避免单条坏消息阻塞整个最小原型。
            try:
                channel.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
            except Exception:
                # 如果消费连接本身也已经断开，交给 pika 外层连接恢复或容器重启处理。
                logger.exception("rabbitmq delivery nack failed", extra={"delivery_tag": method.delivery_tag})
                self._drop_channel()
                return
            logger.exception("rabbitmq delivery failed and was nacked", extra={"delivery_tag": method.delivery_tag})
            return
        try:
            channel.basic_ack(delivery_tag=method.delivery_tag)
        except Exception:
            # 长耗时业务调用可能超过 RabbitMQ 心跳窗口，完成后 ack 时消费连接已被 broker 重置。
            # 此时业务事实和 completed 状态已经发布，不能让进程因为 ack 异常退出。
            logger.exception("rabbitmq delivery ack failed after successful processing", extra={"delivery_tag": method.delivery_tag})
            self._drop_channel()
            return
        logger.info("rabbitmq delivery acked", extra={"delivery_tag": method.delivery_tag})

    def _ensure_channel(self):
        if self.channel is not None:
            return self.channel
        if self.connection_factory is not None:
            self.connection = self.connection_factory()
        else:
            parameters = pika.URLParameters(self.rabbitmq_url)
            self.connection = pika.BlockingConnection(parameters)
        # 真实运行时这里会建立和 RabbitMQ 的阻塞连接；测试里则可以用注入的 fake channel 替代。
        logger.info("rabbitmq consumer connection created", extra={"rabbitmq_url": self.rabbitmq_url, "queue": self.queue})
        self.channel = self.connection.channel()
        return self.channel

    def _drop_channel(self) -> None:
        # 连接断开后丢弃旧引用，下一轮消费循环重新声明 exchange、queue 和绑定。
        self.channel = None
        self.connection = None
