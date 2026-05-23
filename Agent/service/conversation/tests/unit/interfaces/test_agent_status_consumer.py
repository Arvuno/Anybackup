import pytest

from app.application.commands.agent_events import CoreAgentStatusEventCommand
from app.interfaces.mq import agent_status_consumer
from app.interfaces.mq.agent_status_consumer import RabbitMqCoreAgentStatusConsumer


@pytest.mark.asyncio
async def test_status_consumer_binds_core_agent_status_exchange(monkeypatch) -> None:
    connection = FakeConnection()

    async def fake_connect_robust(url: str) -> FakeConnection:
        connection.url = url
        return connection

    monkeypatch.setattr(agent_status_consumer.aio_pika, "connect_robust", fake_connect_robust)
    consumer = RabbitMqCoreAgentStatusConsumer(
        rabbitmq_url="amqp://guest:guest@localhost:5672/",
        exchange_name="core_agent.run_status.events",
        queue_name="conversation.core_agent.run_status",
        prefetch_count=10,
        handler=NoopHandler(),
    )

    await consumer.start()

    assert connection.url == "amqp://guest:guest@localhost:5672/"
    assert connection.channel_instance.qos_calls == [10]
    assert connection.channel_instance.declared_exchanges == [
        {"name": "core_agent.run_status.events", "durable": True}
    ]
    assert connection.channel_instance.declared_queues == [
        {"name": "conversation.core_agent.run_status", "durable": True}
    ]
    assert connection.channel_instance.queue.bind_calls == [
        {"exchange": "core_agent.run_status.events", "routing_key": "core_agent.run.accepted.v1"},
        {"exchange": "core_agent.run_status.events", "routing_key": "core_agent.run.completed.v1"},
        {"exchange": "core_agent.run_status.events", "routing_key": "core_agent.run.failed.v1"},
    ]
    assert connection.channel_instance.queue.consume_calls == [False]


class NoopHandler:
    async def handle(self, command: CoreAgentStatusEventCommand) -> object:
        return command


class FakeConnection:
    def __init__(self) -> None:
        self.url = ""
        self.channel_instance = FakeChannel()

    async def channel(self) -> "FakeChannel":
        return self.channel_instance

    async def close(self) -> None:
        pass


class FakeChannel:
    def __init__(self) -> None:
        self.qos_calls: list[int] = []
        self.declared_exchanges: list[dict[str, object]] = []
        self.declared_queues: list[dict[str, object]] = []
        self.queue = FakeQueue()

    async def set_qos(self, *, prefetch_count: int) -> None:
        self.qos_calls.append(prefetch_count)

    async def declare_exchange(self, name: str, exchange_type, *, durable: bool):
        del exchange_type
        self.declared_exchanges.append({"name": name, "durable": durable})
        return FakeExchange(name=name)

    async def declare_queue(self, name: str, *, durable: bool) -> "FakeQueue":
        self.declared_queues.append({"name": name, "durable": durable})
        return self.queue


class FakeExchange:
    def __init__(self, *, name: str) -> None:
        self.name = name


class FakeQueue:
    def __init__(self) -> None:
        self.bind_calls: list[dict[str, str]] = []
        self.consume_calls: list[bool] = []

    async def bind(self, exchange: FakeExchange, *, routing_key: str) -> None:
        self.bind_calls.append({"exchange": exchange.name, "routing_key": routing_key})

    async def consume(self, callback, *, no_ack: bool) -> None:
        del callback
        self.consume_calls.append(no_ack)

