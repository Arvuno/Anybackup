import json
import logging

import pika

from core_agent_service.bootstrap import build_app
from core_agent_service.config.settings import Settings
from core_agent_service.infrastructure.mq.consumer import RabbitMQConsumer
from core_agent_service.infrastructure.mq.consumer import RabbitMQConsumerPool
from core_agent_service.infrastructure.mq.publisher import StatusPublisher


class FakeRelayService:
    def __init__(self):
        self.events: list[dict] = []

    def handle_message(self, raw_event: dict):
        self.events.append(raw_event)
        return type("RelayResult", (), {"status": "completed"})()

    def handle_cancel(self, raw_event: dict):
        self.events.append(raw_event)
        return type("RelayResult", (), {"status": "completed"})()


class FakeSession:
    def __init__(self):
        self.commits = 0
        self.rollbacks = 0

    def commit(self):
        self.commits += 1

    def rollback(self):
        self.rollbacks += 1


def test_build_app_wires_consumer_and_publisher(tmp_path):
    settings = Settings(
        database_url=f"sqlite+pysqlite:///{tmp_path / 'runtime.db'}",
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        rabbitmq_exchange="conversation.message.events",
        rabbitmq_queue="core_agent.message.events",
        rabbitmq_status_exchange="core_agent.run_status.events",
        rabbitmq_status_routing_key=None,
        kweaver_base_url="https://example.com",
        kweaver_decision_agent_id="agent_001",
        kweaver_probe_on_startup=False,
        core_agent_log_file=str(tmp_path / "core_agent_service.log"),
    )

    app = build_app(settings)

    assert app["consumer"].exchange == "conversation.message.events"
    assert app["consumer"].queue == "core_agent.message.events"
    assert app["publisher"].exchange == "core_agent.run_status.events"
    assert app["publisher"].status_routing_key is None
    assert app["publisher"].accepted_key == "core_agent.run.accepted.v1"


def test_consumer_publishes_accepted_and_completed():
    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        channel=FakeChannel(),
    )
    relay_service = FakeRelayService()
    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=relay_service,
        publisher=publisher,
    )

    consumer.handle_delivery(
        json.dumps(
            {
                "event_id": "evt_001",
                "event_type": "conversation.message.sent",
                "occurred_at": "2026-04-22T10:00:00Z",
                "source_service": "conversation_service",
                "payload": {"conversation_id": "conv_001", "message_id": "msg_001", "turn_id": "turn_001", "content": "hello"},
            }
        ).encode("utf-8")
    )

    assert [item[0] for item in publisher.published] == [
        "core_agent.run.accepted.v1",
        "core_agent.run.completed.v1",
    ]
    published_payloads = [json.loads(item[1])["payload"] for item in publisher.published]
    assert published_payloads == [
        {"conversation_id": "conv_001", "message_id": "msg_001", "content": "accepted"},
        {"conversation_id": "conv_001", "message_id": "msg_001", "content": "completed"},
    ]


class FakeChannel:
    def __init__(self, fail_publish_once: bool = False, fail_ack: bool = False, fail_start_consuming: bool = False):
        self.calls: list[tuple[str, tuple, dict]] = []
        self.callback = None
        self.fail_publish_once = fail_publish_once
        self.fail_ack = fail_ack
        self.fail_start_consuming = fail_start_consuming
        self.is_closed = False

    def exchange_declare(self, *args, **kwargs):
        self.calls.append(("exchange_declare", args, kwargs))

    def queue_declare(self, *args, **kwargs):
        self.calls.append(("queue_declare", args, kwargs))

    def queue_bind(self, *args, **kwargs):
        self.calls.append(("queue_bind", args, kwargs))

    def basic_qos(self, *args, **kwargs):
        self.calls.append(("basic_qos", args, kwargs))

    def basic_consume(self, *args, **kwargs):
        self.calls.append(("basic_consume", args, kwargs))
        self.callback = kwargs["on_message_callback"]

    def start_consuming(self):
        self.calls.append(("start_consuming", tuple(), {}))
        if self.fail_start_consuming:
            self.fail_start_consuming = False
            self.is_closed = True
            raise pika.exceptions.ConnectionWrongStateError("connection closed")

    def basic_publish(self, *args, **kwargs):
        self.calls.append(("basic_publish", args, kwargs))
        if self.fail_publish_once:
            self.fail_publish_once = False
            self.is_closed = True
            raise pika.exceptions.StreamLostError("connection reset by peer")

    def basic_ack(self, *args, **kwargs):
        self.calls.append(("basic_ack", args, kwargs))
        if self.fail_ack:
            raise pika.exceptions.StreamLostError("connection reset by peer")

    def basic_nack(self, *args, **kwargs):
        self.calls.append(("basic_nack", args, kwargs))


class FakeConnection:
    def __init__(self, channel: FakeChannel):
        self._channel = channel

    def channel(self):
        return self._channel


class FakeConnectionFactory:
    def __init__(self, channels: list[FakeChannel]):
        self.channels = channels
        self.index = 0

    def __call__(self):
        channel = self.channels[self.index]
        self.index += 1
        return FakeConnection(channel)


def test_consumer_run_uses_real_channel_setup():
    fake_channel = FakeChannel()
    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=FakeRelayService(),
        publisher=StatusPublisher(exchange="core_agent.run_status.events", status_routing_key=None),
        connection_factory=lambda: FakeConnection(fake_channel),
    )

    consumer.run()

    call_names = [name for name, _, _ in fake_channel.calls]
    assert "exchange_declare" in call_names
    assert "queue_declare" in call_names
    assert call_names.count("queue_bind") == 2
    assert "basic_consume" in call_names
    assert "start_consuming" in call_names


def test_consumer_run_reconnects_after_connection_loss():
    first_channel = FakeChannel(fail_start_consuming=True)
    second_channel = FakeChannel()
    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=FakeRelayService(),
        publisher=StatusPublisher(exchange="core_agent.run_status.events", status_routing_key=None),
        connection_factory=FakeConnectionFactory([first_channel, second_channel]),
    )

    consumer.run()

    assert [name for name, _, _ in first_channel.calls].count("start_consuming") == 1
    assert [name for name, _, _ in second_channel.calls].count("start_consuming") == 1


def test_consumer_pool_runs_all_consumers_in_parallel():
    calls: list[str] = []

    class FakeConsumer:
        def __init__(self, name: str):
            self.name = name

        def run(self):
            calls.append(self.name)

    pool = RabbitMQConsumerPool([FakeConsumer("worker-1"), FakeConsumer("worker-2")])

    pool.run()

    assert sorted(calls) == ["worker-1", "worker-2"]


def test_publisher_uses_status_exchange_and_event_routing_key():
    fake_channel = FakeChannel()
    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        channel=fake_channel,
    )

    publisher.publish_completed("conv_001", "completed")

    call_names = [name for name, _, _ in fake_channel.calls]
    assert "exchange_declare" in call_names
    assert "basic_publish" in call_names
    publish_call = next(call for call in fake_channel.calls if call[0] == "basic_publish")
    assert publish_call[2]["exchange"] == "core_agent.run_status.events"
    assert publish_call[2]["routing_key"] == "core_agent.run.completed.v1"


def test_publisher_logs_completed_status_payload(caplog):
    fake_channel = FakeChannel()
    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        channel=fake_channel,
    )

    with caplog.at_level(logging.INFO, logger="core_agent_service.infrastructure.mq.publisher"):
        publisher.publish_completed("conv_001", "completed", message_id="msg_001")

    log_text = "\n".join(record.getMessage() for record in caplog.records)
    assert "rabbitmq status event publish payload" in log_text
    assert '"routing_key": "core_agent.run.completed.v1"' in log_text
    assert '"event_type": "core_agent.run.completed"' in log_text
    assert '"conversation_id": "conv_001"' in log_text
    assert '"message_id": "msg_001"' in log_text
    assert '"content": "completed"' in log_text


def test_publisher_reconnects_and_retries_after_stream_loss():
    first_channel = FakeChannel(fail_publish_once=True)
    second_channel = FakeChannel()
    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        connection_factory=FakeConnectionFactory([first_channel, second_channel]),
    )

    publisher.publish_accepted("conv_retry_001")

    assert [name for name, _, _ in first_channel.calls].count("basic_publish") == 1
    assert [name for name, _, _ in second_channel.calls].count("basic_publish") == 1


def test_consumer_acks_message_after_success():
    fake_channel = FakeChannel()
    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        channel=FakeChannel(),
    )
    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=FakeRelayService(),
        publisher=publisher,
    )
    method = type("Method", (), {"delivery_tag": "tag-001"})()

    consumer.on_message(
        fake_channel,
        method,
        None,
        json.dumps(
            {
                "event_id": "evt_001",
                "event_type": "conversation.message.sent",
                "occurred_at": "2026-04-22T10:00:00Z",
                "source_service": "conversation_service",
                "payload": {"conversation_id": "conv_001", "message_id": "msg_001", "turn_id": "turn_001", "content": "hello"},
            }
        ).encode("utf-8"),
    )

    assert [name for name, _, _ in fake_channel.calls].count("basic_ack") == 1


def test_consumer_does_not_raise_when_ack_connection_is_lost():
    fake_channel = FakeChannel(fail_ack=True)
    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        channel=FakeChannel(),
    )
    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=FakeRelayService(),
        publisher=publisher,
    )
    method = type("Method", (), {"delivery_tag": "tag-ack-lost"})()

    consumer.on_message(
        fake_channel,
        method,
        None,
        json.dumps(
            {
                "event_id": "evt_ack_lost",
                "event_type": "conversation.message.sent",
                "occurred_at": "2026-04-22T10:00:00Z",
                "source_service": "conversation_service",
                "payload": {
                    "conversation_id": "conv_ack_lost",
                    "message_id": "msg_ack_lost",
                    "turn_id": "turn_ack_lost",
                    "content": "hello",
                },
            }
        ).encode("utf-8"),
    )

    assert [name for name, _, _ in fake_channel.calls].count("basic_ack") == 1
    assert [item[0] for item in publisher.published] == [
        "core_agent.run.accepted.v1",
        "core_agent.run.completed.v1",
    ]


def test_consumer_nacks_message_after_failure():
    fake_channel = FakeChannel()

    class FailingRelayService:
        def handle_message(self, raw_event: dict):
            raise RuntimeError("boom")

        def handle_cancel(self, raw_event: dict):
            raise RuntimeError("boom")

    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=FailingRelayService(),
        publisher=StatusPublisher(
            rabbitmq_url="amqp://guest:guest@localhost:5672/",
            exchange="core_agent.run_status.events",
            status_routing_key=None,
            channel=FakeChannel(),
        ),
    )
    method = type("Method", (), {"delivery_tag": "tag-001"})()

    consumer.on_message(
        fake_channel,
        method,
        None,
        json.dumps(
            {
                "event_id": "evt_001",
                "event_type": "conversation.message.sent",
                "occurred_at": "2026-04-22T10:00:00Z",
                "source_service": "conversation_service",
                "payload": {"conversation_id": "conv_001", "message_id": "msg_001", "turn_id": "turn_001", "content": "hello"},
            }
        ).encode("utf-8"),
    )

    assert [name for name, _, _ in fake_channel.calls].count("basic_nack") == 1


def test_consumer_publishes_failed_and_keeps_callback_alive_after_kweaver_error():
    fake_channel = FakeChannel()
    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        channel=FakeChannel(),
    )

    class KWeaverTimeoutRelayService:
        def handle_message(self, raw_event: dict):
            raise RuntimeError("KWeaver request failed with HTTP 504")

        def handle_cancel(self, raw_event: dict):
            raise RuntimeError("KWeaver request failed with HTTP 504")

    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=KWeaverTimeoutRelayService(),
        publisher=publisher,
    )
    method = type("Method", (), {"delivery_tag": "tag-504"})()

    consumer.on_message(
        fake_channel,
        method,
        None,
        json.dumps(
            {
                "event_id": "evt_504",
                "event_type": "conversation.message.sent",
                "occurred_at": "2026-04-22T10:00:00Z",
                "source_service": "conversation_service",
                "payload": {"conversation_id": "conv_504", "message_id": "msg_504", "turn_id": "turn_504", "content": "hello"},
            }
        ).encode("utf-8"),
    )

    assert [name for name, _, _ in fake_channel.calls].count("basic_nack") == 1
    assert [name for name, _, _ in fake_channel.calls].count("basic_ack") == 0
    assert [item[0] for item in publisher.published] == [
        "core_agent.run.accepted.v1",
        "core_agent.run.failed.v1",
    ]
    failed_event = json.loads(publisher.published[-1][1])
    assert failed_event["event_type"] == "core_agent.run.failed"
    assert failed_event["payload"] == {
        "conversation_id": "conv_504",
        "message_id": "msg_504",
        "content": "KWeaver request failed with HTTP 504",
    }


def test_consumer_commits_session_after_success():
    fake_channel = FakeChannel()
    fake_session = FakeSession()
    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        channel=FakeChannel(),
    )
    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=FakeRelayService(),
        publisher=publisher,
        session=fake_session,
    )
    method = type("Method", (), {"delivery_tag": "tag-001"})()

    consumer.on_message(
        fake_channel,
        method,
        None,
        json.dumps(
            {
                "event_id": "evt_010",
                "event_type": "conversation.message.sent",
                "occurred_at": "2026-04-22T10:00:00Z",
                "source_service": "conversation_service",
                "payload": {"conversation_id": "conv_commit_001", "message_id": "msg_010", "turn_id": "turn_010", "content": "hello"},
            }
        ).encode("utf-8"),
    )

    assert fake_session.commits == 1
    assert fake_session.rollbacks == 0


def test_consumer_rolls_back_then_commits_failed_status():
    fake_channel = FakeChannel()
    fake_session = FakeSession()

    class FailingRelayService:
        def handle_message(self, raw_event: dict):
            raise RuntimeError("boom")

        def handle_cancel(self, raw_event: dict):
            raise RuntimeError("boom")

    publisher = StatusPublisher(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="core_agent.run_status.events",
        status_routing_key=None,
        channel=FakeChannel(),
    )
    consumer = RabbitMQConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange="conversation.message.events",
        queue="core_agent.message.events",
        relay_service=FailingRelayService(),
        publisher=publisher,
        session=fake_session,
    )
    method = type("Method", (), {"delivery_tag": "tag-011"})()

    consumer.on_message(
        fake_channel,
        method,
        None,
        json.dumps(
            {
                "event_id": "evt_011",
                "event_type": "conversation.message.sent",
                "occurred_at": "2026-04-22T10:00:00Z",
                "source_service": "conversation_service",
                "payload": {"conversation_id": "conv_commit_002", "message_id": "msg_011", "turn_id": "turn_011", "content": "hello"},
            }
        ).encode("utf-8"),
    )

    assert fake_session.rollbacks == 1
    assert fake_session.commits == 1
